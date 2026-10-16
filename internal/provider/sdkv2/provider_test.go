// Copyright © 2026 Ping Identity Corporation

package sdkv2_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
	"github.com/zalando/go-keyring"
)

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

// blankCredentialEnv unsets every authentication-related variable, clears the
// session memo and resets the mock keychain.
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

// runConfigure executes the SDKv2 provider configure function against raw
// attribute data and returns the client and the first error diagnostic.
func runConfigure(t *testing.T, raw map[string]any) (any, error) {
	t.Helper()

	p := sdkv2.New("test")()

	d := schema.TestResourceDataRaw(t, p.Schema, raw)

	client, diags := sdkv2.ConfigureForTest(context.Background(), d)

	var firstErr error
	for _, d := range diags {
		if d.Severity == diag.Error && firstErr == nil {
			firstErr = fmt.Errorf("%s: %s", d.Summary, d.Detail)
		}
	}

	return client, firstErr
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

func TestSDKv2Configure_PingCLISession(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "sdkv2-token")

	client, err := runConfigure(t, map[string]any{
		"region_code":    "NA",
		"config_profile": "default",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}

	if client == nil {
		t.Fatal("configure returned nil client")
	}
}

func TestSDKv2Configure_PingCLIConflictWithHCLCredential(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "token")

	_, err := runConfigure(t, map[string]any{
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

func TestSDKv2Configure_PingCLIImplicitEnvOnly(t *testing.T) {
	blankCredentialEnv(t)

	// The documented env-only workflow: no HCL attributes at all, the region
	// and profile come from the environment.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "env-only-token")
	t.Setenv("PINGCLI_PROFILE", "default")
	t.Setenv("PINGONE_REGION_CODE", "NA")

	client, err := runConfigure(t, map[string]any{})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}

	if client == nil {
		t.Fatal("configure returned nil client")
	}
}

func TestSDKv2Configure_PingCLIMissingRegion(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "token")

	_, err := runConfigure(t, map[string]any{
		"config_profile": "default",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a missing-region error")
	}

	if !contains(err.Error(), "Missing region code") {
		t.Errorf("error = %q, want a missing-region message", err.Error())
	}
}

func TestSDKv2Configure_PingCLINoSession(t *testing.T) {
	blankCredentialEnv(t)

	_, err := runConfigure(t, map[string]any{
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

func TestSDKv2Configure_PingCLIGrantTypePinned(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeDeviceCode, "default", "device-token")

	client, err := runConfigure(t, map[string]any{
		"region_code":       "NA",
		"config_profile":    "default",
		"config_grant_type": pingcli.GrantTypeDeviceCode,
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}

	if client == nil {
		t.Fatal("configure returned nil client")
	}
}

func TestSDKv2Configure_PingCLIGrantTypeMismatch(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeDeviceCode, "default", "device-token")

	_, err := runConfigure(t, map[string]any{
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

func TestSDKv2Configure_PingCLIGrantTypeWithoutProfile(t *testing.T) {
	blankCredentialEnv(t)

	_, err := runConfigure(t, map[string]any{
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

func TestSDKv2Configure_ClientCredentialsFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	_, err := runConfigure(t, map[string]any{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
		"http_proxy":     "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestSDKv2Configure_ClientCredentialsFromEnv(t *testing.T) {
	blankCredentialEnv(t)

	t.Setenv("PINGONE_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	t.Setenv("PINGONE_CLIENT_SECRET", "secret")
	t.Setenv("PINGONE_ENVIRONMENT_ID", "22222222-2222-2222-2222-222222222222")

	_, err := runConfigure(t, map[string]any{
		"region_code": "NA",
		"http_proxy":  "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestSDKv2Configure_ClientCredentialsWinOverExportedPingCLIEnv(t *testing.T) {
	blankCredentialEnv(t)

	// A Ping CLI setup exported in the shell must never break a credentials
	// configuration.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "unused-token")
	t.Setenv("PINGCLI_PROFILE", "default")
	t.Setenv("PINGCLI_PINGONE_AUTHENTICATION_TYPE", pingcli.GrantTypeClientCredentials)

	_, err := runConfigure(t, map[string]any{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
		"http_proxy":     "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

func TestSDKv2Configure_AccessTokenFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	_, err := runConfigure(t, map[string]any{
		"api_access_token": "static-token",
		"region_code":      "NA",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}
}

func TestSDKv2Configure_PingCLIWinsOverExportedCredentialEnv(t *testing.T) {
	blankCredentialEnv(t)

	// PINGONE_ENVIRONMENT_ID in particular is commonly exported. An explicit
	// config_profile must win over all of them rather than report a conflict.
	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "env-wins-token")
	t.Setenv("PINGONE_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	t.Setenv("PINGONE_CLIENT_SECRET", "secret")
	t.Setenv("PINGONE_ENVIRONMENT_ID", "22222222-2222-2222-2222-222222222222")

	client, err := runConfigure(t, map[string]any{
		"region_code":    "NA",
		"config_profile": "default",
	})
	if err != nil {
		t.Fatalf("configure returned error: %v", err)
	}

	if client == nil {
		t.Fatal("configure returned nil client")
	}
}

func TestSDKv2Configure_CredentialModeUnaffected(t *testing.T) {
	blankCredentialEnv(t)

	// With no Ping CLI attributes or credentials, client construction fails
	// validation, which shows Configure took the normal path.
	_, err := runConfigure(t, map[string]any{
		"region_code": "NA",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a validation error for missing credentials")
	}
}

// configureWithUnknownAttribute runs the SDKv2 provider configure function
// against a ResourceData whose raw config marks attributeName as unknown — the
// shape Terraform sends when the value can only be known after apply.
// schema.TestResourceDataRaw cannot express unknowns, so the raw config is
// built with cty and attached through InstanceState.RawConfig.
func configureWithUnknownAttribute(t *testing.T, attributeName string) error {
	t.Helper()

	p := sdkv2.New("test")()

	schemaType, ok := map[string]cty.Type{
		pingcli.AttrConfigProfile:   cty.String,
		pingcli.AttrConfigGrantType: cty.String,
	}[attributeName]
	if !ok {
		t.Fatalf("attribute %q is not a known Ping CLI attribute", attributeName)
	}

	rawConfig := cty.ObjectVal(map[string]cty.Value{
		attributeName: cty.UnknownVal(schemaType),
	})

	d := (&schema.Resource{Schema: p.Schema}).Data(&terraform.InstanceState{RawConfig: rawConfig})

	_, diags := sdkv2.ConfigureForTest(context.Background(), d)

	for _, d := range diags {
		if d.Severity == diag.Error {
			return fmt.Errorf("%s: %s", d.Summary, d.Detail)
		}
	}

	return nil
}

func TestSDKv2Configure_PingCLIUnknownValues(t *testing.T) {
	for _, attribute := range []string{pingcli.AttrConfigProfile, pingcli.AttrConfigGrantType} {
		t.Run(attribute, func(t *testing.T) {
			blankCredentialEnv(t)

			// An unknown value must be an error, not silently treated as unset.
			t.Setenv("PINGCLI_PROFILE", "default")

			err := configureWithUnknownAttribute(t, attribute)
			if err == nil {
				t.Fatal("configure returned nil error, want an unknown value error")
			}

			if !contains(err.Error(), "Unknown Ping CLI configuration value") {
				t.Errorf("error = %q, want an unknown value message", err.Error())
			}
		})
	}
}
