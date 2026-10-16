// Copyright © 2026 Ping Identity Corporation

package frameworklegacysdk_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/frameworklegacysdk"
	"github.com/zalando/go-keyring"
)

// unknownValue marks an attribute that configValue builds as unknown.
const unknownValue = "\x00unknown"

var testCtx = context.Background()

// schemaOf runs the provider Schema method and returns the schema.
func schemaOf(t *testing.T, p provider.Provider) providerschema.Schema {
	t.Helper()

	resp := &provider.SchemaResponse{}
	p.Schema(testCtx, provider.SchemaRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("schema returned error diagnostics: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// configValue builds the full provider config object value. Attributes absent
// from values are null; blocks are always null lists.
func configValue(t *testing.T, schema providerschema.Schema, values map[string]string) tftypes.Value {
	t.Helper()

	attrTypes := map[string]tftypes.Type{}
	attrValues := map[string]tftypes.Value{}

	for name, attr := range schema.Attributes {
		attrType := attr.GetType().TerraformType(testCtx)
		attrTypes[name] = attrType

		if v, ok := values[name]; ok && v == unknownValue {
			attrValues[name] = tftypes.NewValue(attrType, tftypes.UnknownValue)
		} else if ok {
			attrValues[name] = tftypes.NewValue(attrType, v)
		} else {
			attrValues[name] = tftypes.NewValue(attrType, nil)
		}
	}

	for name, block := range schema.Blocks {
		nestedTypes := map[string]tftypes.Type{}

		// The provider only declares ListNestedBlock declarations; use the
		// interface accessor so the test does not depend on the concrete type.
		if listBlock, ok := block.(providerschema.ListNestedBlock); ok {
			for attrName, attr := range listBlock.NestedObject.Attributes {
				nestedTypes[attrName] = attr.GetType().TerraformType(testCtx)
			}
			for blockName, nestedBlock := range listBlock.NestedObject.Blocks {
				// global_options.population: recurse into the nested block's
				// own attributes for the element object type.
				if nestedList, ok := nestedBlock.(providerschema.ListNestedBlock); ok {
					innerTypes := map[string]tftypes.Type{}
					for innerName, innerAttr := range nestedList.NestedObject.Attributes {
						innerTypes[innerName] = innerAttr.GetType().TerraformType(testCtx)
					}
					nestedTypes[blockName] = tftypes.List{ElementType: tftypes.Object{AttributeTypes: innerTypes}}
				}
			}
		}

		attrTypes[name] = tftypes.List{ElementType: tftypes.Object{AttributeTypes: nestedTypes}}
		attrValues[name] = tftypes.NewValue(attrTypes[name], nil)
	}

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	return tftypes.NewValue(objectType, attrValues)
}

// configureProvider runs the provider Configure against a config built from
// the real schema and returns the first error diagnostic text, if any.
func configureProvider(t *testing.T, values map[string]string) error {
	t.Helper()

	p := frameworklegacysdk.New("test")()

	schema := schemaOf(t, p)

	req := provider.ConfigureRequest{
		Config: tfsdk.Config{
			Schema: schema,
			Raw:    configValue(t, schema, values),
		},
	}
	resp := &provider.ConfigureResponse{}

	p.Configure(testCtx, req, resp)

	for _, diagnostic := range resp.Diagnostics {
		if diagnostic.Severity() == diag.SeverityError {
			return fmt.Errorf("%s: %s", diagnostic.Summary(), diagnostic.Detail())
		}
	}

	return nil
}

// assertCredentialModeSelected fails if err shows Ping CLI authentication was
// involved. The legacy client fetches its worker token eagerly, so these tests use
// an unreachable proxy and expect only a transport failure.
func assertCredentialModeSelected(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		return
	}

	for _, unwanted := range []string{"Conflicting provider authentication configuration", "Ping CLI", "Missing region code"} {
		if contains(err.Error(), unwanted) {
			t.Fatalf("credentials configuration was handled as Ping CLI authentication: %v", err)
		}
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

func blankCredentialEnv(t *testing.T) {
	t.Helper()

	for _, name := range []string{
		"PINGONE_CLIENT_ID", "PINGONE_CLIENT_SECRET", "PINGONE_ENVIRONMENT_ID",
		"PINGONE_API_ACCESS_TOKEN", "PINGONE_REGION_CODE", "PINGONE_REGION",
		"PINGCLI_PROFILE", "PINGCLI_PINGONE_AUTHENTICATION_TYPE", "PINGONE_AUTH_GRANT_TYPE",
		"PINGONE_STORAGE_TYPE", "PINGONE_STORAGE_NAME",
	} {
		// Unset, don't blank: the SDK's envconfig treats an empty-but-set variable
		// as a value and overrides the provider's configuration.
		t.Setenv(name, "placeholder")
		os.Unsetenv(name)
	}

	pingcli.ResetMemo()
	keyring.MockInit()
}

// seedPingCLISession seeds a mock keychain session for the profile.
func seedPingCLISession(t *testing.T, grant, profile, accessToken string) {
	t.Helper()

	h := sha256.Sum256([]byte("pingone:" + grant))
	account := fmt.Sprintf("auth_session_token-%x_pingone_%s_%s", h[:4], grant, profile)

	tokenJSON, err := json.Marshal(map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expiry":       time.Now().Add(time.Hour).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Fatalf("failed to marshal token: %v", err)
	}

	if err := keyring.Set("pingcli", account, string(tokenJSON)); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}

}

func TestFrameworkLegacySDK_Configure_PingCLISession(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "legacy-token")

	err := configureProvider(t, map[string]string{
		"region_code":    "NA",
		"config_profile": "default",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIConflict(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "token")

	err := configureProvider(t, map[string]string{
		"region_code":    "NA",
		"config_profile": "default",
		"client_id":      "some-client",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a conflict error")
	}

	if !contains(err.Error(), "Conflicting provider authentication configuration") {
		t.Errorf("error = %q, want a conflict message", err.Error())
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIImplicitEnvOnly(t *testing.T) {
	blankCredentialEnv(t)

	// The documented env-only workflow: no HCL attributes at all, the region
	// and profile come from the environment.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "env-only-token")
	t.Setenv("PINGCLI_PROFILE", "default")
	t.Setenv("PINGONE_REGION_CODE", "NA")

	err := configureProvider(t, map[string]string{})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIMissingRegion(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "token")

	err := configureProvider(t, map[string]string{
		"config_profile": "default",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a missing-region error")
	}

	if !contains(err.Error(), "Missing region code") {
		t.Errorf("error = %q, want a missing-region message", err.Error())
	}
}

func TestFrameworkLegacySDK_Configure_PingCLINoSession(t *testing.T) {
	blankCredentialEnv(t)

	err := configureProvider(t, map[string]string{
		"region_code":    "NA",
		"config_profile": "default",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a no-session error")
	}

	if !contains(err.Error(), "no Ping CLI authentication session") {
		t.Errorf("error = %q, want a no-session message", err.Error())
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIGrantTypePinned(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeDeviceCode, "default", "device-token")

	err := configureProvider(t, map[string]string{
		"region_code":       "NA",
		"config_profile":    "default",
		"config_grant_type": pingcli.GrantTypeDeviceCode,
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIGrantTypeMismatch(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeDeviceCode, "default", "device-token")

	err := configureProvider(t, map[string]string{
		"region_code":       "NA",
		"config_profile":    "default",
		"config_grant_type": pingcli.GrantTypeClientCredentials,
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a no-session error")
	}

	if !contains(err.Error(), "no Ping CLI authentication session") {
		t.Errorf("error = %q, want a no-session message", err.Error())
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIGrantTypeWithoutProfile(t *testing.T) {
	blankCredentialEnv(t)

	err := configureProvider(t, map[string]string{
		"region_code":       "NA",
		"config_grant_type": pingcli.GrantTypeDeviceCode,
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a missing-profile error")
	}

	if !contains(err.Error(), "Missing Ping CLI profile") {
		t.Errorf("error = %q, want a missing-profile message", err.Error())
	}
}

func TestFrameworkLegacySDK_Configure_ClientCredentialsFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	err := configureProvider(t, map[string]string{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
		"http_proxy":     "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestFrameworkLegacySDK_Configure_ClientCredentialsFromEnv(t *testing.T) {
	blankCredentialEnv(t)

	t.Setenv("PINGONE_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	t.Setenv("PINGONE_CLIENT_SECRET", "secret")
	t.Setenv("PINGONE_ENVIRONMENT_ID", "22222222-2222-2222-2222-222222222222")

	err := configureProvider(t, map[string]string{
		"region_code": "NA",
		"http_proxy":  "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestFrameworkLegacySDK_Configure_ClientCredentialsWinOverExportedPingCLIEnv(t *testing.T) {
	blankCredentialEnv(t)

	// A Ping CLI setup exported in the shell must never break a credentials
	// configuration.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "unused-token")
	t.Setenv("PINGCLI_PROFILE", "default")
	t.Setenv("PINGCLI_PINGONE_AUTHENTICATION_TYPE", pingcli.GrantTypeClientCredentials)

	err := configureProvider(t, map[string]string{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
		"http_proxy":     "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestFrameworkLegacySDK_Configure_AccessTokenFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	err := configureProvider(t, map[string]string{
		"api_access_token": "static-token",
		"region_code":      "NA",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIWinsOverExportedCredentialEnv(t *testing.T) {
	blankCredentialEnv(t)

	// PINGONE_ENVIRONMENT_ID in particular is commonly exported. An explicit
	// config_profile must win over all of them rather than report a conflict.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "env-wins-token")
	t.Setenv("PINGONE_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	t.Setenv("PINGONE_CLIENT_SECRET", "secret")
	t.Setenv("PINGONE_ENVIRONMENT_ID", "22222222-2222-2222-2222-222222222222")

	err := configureProvider(t, map[string]string{
		"region_code":    "NA",
		"config_profile": "default",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFrameworkLegacySDK_Configure_PingCLIUnknownValues(t *testing.T) {
	for _, attribute := range []string{"config_profile", "config_grant_type"} {
		t.Run(attribute, func(t *testing.T) {
			blankCredentialEnv(t)

			// An unknown value must be an error, not silently treated as unset.
			t.Setenv("PINGCLI_PROFILE", "default")

			values := map[string]string{"region_code": "NA", attribute: unknownValue}
			if attribute == "config_grant_type" {
				values["config_profile"] = "default"
			}

			err := configureProvider(t, values)
			if err == nil {
				t.Fatal("configure returned nil error, want an unknown value error")
			}

			if !contains(err.Error(), "Unknown Ping CLI configuration value") {
				t.Errorf("error = %q, want an unknown value message", err.Error())
			}
		})
	}
}

func TestFrameworkLegacySDK_Configure_CredentialModeUnaffected(t *testing.T) {
	blankCredentialEnv(t)

	// With no Ping CLI attributes or credentials, client construction fails
	// validation, which shows Configure took the normal path.
	err := configureProvider(t, map[string]string{
		"region_code": "NA",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a validation error for missing credentials")
	}
}
