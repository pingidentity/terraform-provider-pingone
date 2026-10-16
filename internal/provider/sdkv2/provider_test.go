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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// blankCredentialEnv unsets every authentication-related variable and resets
// the mock keychain. Unset, don't blank: the SDK's envconfig treats an
// empty-but-set variable as a value and overrides the provider's configuration.
func blankCredentialEnv(t *testing.T) {
	t.Helper()

	for _, name := range []string{
		"PINGONE_CLIENT_ID", "PINGONE_CLIENT_SECRET", "PINGONE_ENVIRONMENT_ID",
		"PINGONE_API_ACCESS_TOKEN", "PINGONE_REGION_CODE", "PINGONE_REGION",
		"PINGCLI_PROFILE", "PINGCLI_PINGONE_AUTHENTICATION_TYPE", "PINGONE_AUTH_GRANT_TYPE",
		"PINGONE_STORAGE_TYPE", "PINGONE_STORAGE_NAME",
	} {
		t.Setenv(name, "placeholder")
		os.Unsetenv(name)
	}

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

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

// The unit tests below cover the provider's own decisions around Ping CLI
// authentication: what enables it and what conflicts with it. Session lookup
// and the live behaviour are covered by the pingcli package tests and the
// Ping CLI acceptance tests.

func TestSDKv2Configure_PingCLISession(t *testing.T) {
	blankCredentialEnv(t)

	seedPingCLISession(t, pingcli.GrantTypeClientCredentials, "default", "sdkv2-token")

	client, err := runConfigure(t, map[string]any{
		"region_code":      "NA",
		"use_pingcli_auth": true,
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
		"region_code":      "NA",
		"use_pingcli_auth": true,
		"client_id":        "some-client",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a conflict error")
	}

	if !contains(err.Error(), "Conflicting provider authentication configuration") {
		t.Errorf("error = %q, want a conflict message", err.Error())
	}
}

func TestSDKv2Configure_ClientCredentialsFromHCL(t *testing.T) {
	blankCredentialEnv(t)

	// The legacy SDK fetches its worker token eagerly, so the configuration is
	// pointed at an unreachable proxy: the expected failure is the transport,
	// never Ping CLI resolution.
	_, err := runConfigure(t, map[string]any{
		"client_id":      "11111111-1111-1111-1111-111111111111",
		"client_secret":  "secret",
		"environment_id": "22222222-2222-2222-2222-222222222222",
		"region_code":    "NA",
		"http_proxy":     "http://127.0.0.1:1",
	})
	assertCredentialModeSelected(t, err)
}

// assertCredentialModeSelected fails if err shows Ping CLI authentication was
// involved.
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

func TestSDKv2Configure_CredentialModeUnaffected(t *testing.T) {
	blankCredentialEnv(t)

	// With no Ping CLI attributes and no credentials, client construction
	// fails with a validation error, which proves Configure ran the normal
	// path.
	_, err := runConfigure(t, map[string]any{
		"region_code": "NA",
	})
	if err == nil {
		t.Fatal("configure returned nil error, want a validation error for missing credentials")
	}
}
