// Copyright © 2026 Ping Identity Corporation

package framework_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/framework"
	"github.com/zalando/go-keyring"
)

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

		if v, ok := values[name]; ok {
			if _, isBool := attr.(providerschema.BoolAttribute); isBool {
				attrValues[name] = tftypes.NewValue(attrType, v == "true")
			} else {
				attrValues[name] = tftypes.NewValue(attrType, v)
			}
		} else {
			attrValues[name] = tftypes.NewValue(attrType, nil)
		}
	}

	for name, block := range schema.Blocks {
		nestedTypes := map[string]tftypes.Type{}

		if listBlock, ok := block.(providerschema.ListNestedBlock); ok {
			for attrName, attr := range listBlock.NestedObject.Attributes {
				nestedTypes[attrName] = attr.GetType().TerraformType(testCtx)
			}
			for blockName, nestedBlock := range listBlock.NestedObject.Blocks {
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

	p := framework.New("test")()

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

	keyring.MockInit()
}

// seedPingCLISession seeds a mock keychain session for the profile.
func seedPingCLISession(t *testing.T, grant, profile, accessToken string) {
	t.Helper()

	keyring.MockInit()

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

// The unit tests below cover the provider's own decisions around Ping CLI
// authentication: what enables it and what conflicts with it. Session lookup
// and the live behaviour are covered by the pingcli package tests and the
// Ping CLI acceptance tests.

func TestFramework_Configure_PingCLISession(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "framework-token")

	err := configureProvider(t, map[string]string{
		"region_code":      "NA",
		"use_pingcli_auth": "true",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFramework_Configure_PingCLIConflict(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "token")

	err := configureProvider(t, map[string]string{
		"region_code":      "NA",
		"use_pingcli_auth": "true",
		"client_id":        "some-client",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a conflict error")
	}

	if !contains(err.Error(), "Conflicting provider authentication configuration") {
		t.Errorf("error = %q, want a conflict message", err.Error())
	}
}

func TestFramework_Configure_ClientCredentialsFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	err := configureProvider(t, map[string]string{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestFramework_Configure_CredentialModeUnaffected(t *testing.T) {
	blankCredentialEnv(t)

	// With no Ping CLI attributes and no credentials, client construction
	// fails with a validation error, proving Configure ran the normal path.
	err := configureProvider(t, map[string]string{
		"region_code": "NA",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want an error for missing credentials")
	}
}
