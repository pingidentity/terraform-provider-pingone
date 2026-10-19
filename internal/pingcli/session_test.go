// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/zalando/go-keyring"
)

// sessionAccount derives the Ping CLI v1.x keychain account name the same way
// the production code does, so tests assert against the real key layout.
func sessionAccount(t *testing.T, grant, profile string) string {
	t.Helper()

	h := sha256.Sum256([]byte("pingone:" + grant))

	return fmt.Sprintf("auth_session_token-%x_pingone_%s_%s", h[:4], grant, profile)
}

// seedKeychainSession stores a single-entry session in the mock keychain.
func seedKeychainSession(t *testing.T, account, accessToken string, expiry time.Time) {
	t.Helper()

	tokenJSON, err := json.Marshal(map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expiry":       expiry.Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Fatalf("failed to marshal token: %v", err)
	}

	if err := keyring.Set("pingcli", account, string(tokenJSON)); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}
}

func TestLoad_SingleEntrySession(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")
	seedKeychainSession(t, account, "test-access-token", time.Now().Add(time.Hour))

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "test-access-token" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "test-access-token")
	}
	if result.Profile != "default" {
		t.Errorf("Profile = %q, want %q", result.Profile, "default")
	}
	if result.GrantType != pingcli.GrantTypeClientCredentials {
		t.Errorf("GrantType = %q, want %q", result.GrantType, pingcli.GrantTypeClientCredentials)
	}
	if result.Expiry.IsZero() {
		t.Error("Expiry is zero, want the seeded expiry")
	}
}

func TestLoad_FindsSessionForEveryGrantType(t *testing.T) {
	for _, grant := range []string{
		pingcli.GrantTypeClientCredentials,
		pingcli.GrantTypeAuthorizationCode,
		pingcli.GrantTypeDeviceCode,
	} {
		t.Run(grant, func(t *testing.T) {
			keyring.MockInit()

			seedKeychainSession(t, sessionAccount(t, grant, "dev"), "token-"+grant, time.Now().Add(time.Hour))

			result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "dev"})
			if err != nil {
				t.Fatalf("Load returned error: %v", err)
			}

			if result.GrantType != grant {
				t.Errorf("GrantType = %q, want %q", result.GrantType, grant)
			}
			if result.AccessToken != "token-"+grant {
				t.Errorf("AccessToken = %q, want %q", result.AccessToken, "token-"+grant)
			}
		})
	}
}

func TestLoad_PrefersLongestLivedSession(t *testing.T) {
	keyring.MockInit()

	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeClientCredentials, "default"), "short-lived", time.Now().Add(10*time.Minute))
	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeDeviceCode, "default"), "long-lived", time.Now().Add(2*time.Hour))

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "long-lived" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "long-lived")
	}
}

func TestLoad_SkipsExpiredSessionWhenAnotherIsValid(t *testing.T) {
	keyring.MockInit()

	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeClientCredentials, "default"), "expired-token", time.Now().Add(-time.Hour))
	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeAuthorizationCode, "default"), "valid-token", time.Now().Add(time.Hour))

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "valid-token" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "valid-token")
	}
}

func TestLoad_PinnedGrantTypeIgnoresOtherGrantTypes(t *testing.T) {
	keyring.MockInit()

	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeClientCredentials, "default"), "short-lived", time.Now().Add(10*time.Minute))
	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeDeviceCode, "default"), "long-lived", time.Now().Add(2*time.Hour))

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default", PingCLIGrantType: pingcli.GrantTypeClientCredentials})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "short-lived" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "short-lived")
	}
	if result.GrantType != pingcli.GrantTypeClientCredentials {
		t.Errorf("GrantType = %q, want %q", result.GrantType, pingcli.GrantTypeClientCredentials)
	}
}

func TestLoad_PinnedGrantTypeWithoutSession(t *testing.T) {
	keyring.MockInit()

	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeDeviceCode, "default"), "device-token", time.Now().Add(time.Hour))

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default", PingCLIGrantType: pingcli.GrantTypeAuthorizationCode})
	if err == nil {
		t.Fatal("Load returned nil error, want a not-found error")
	}

	assertErrorContains(t, err, "no Ping CLI authentication session", "'default'", "'authorization_code'", "pingcli pingone auth login")
}

func TestLoad_UnsupportedGrantType(t *testing.T) {
	keyring.MockInit()

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default", PingCLIGrantType: "password"})
	if err == nil {
		t.Fatal("Load returned nil error, want an unsupported grant type error")
	}

	assertErrorContains(t, err, "not supported", "'password'")
}

func TestLoad_OtherProfileSessionIsIgnored(t *testing.T) {
	keyring.MockInit()

	seedKeychainSession(t, sessionAccount(t, pingcli.GrantTypeClientCredentials, "other"), "other-token", time.Now().Add(time.Hour))

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "dev"})
	if err == nil {
		t.Fatal("Load returned nil error, want a not-found error")
	}

	assertErrorContains(t, err, "no Ping CLI authentication session", "'dev'")
}

func TestLoad_SplitSession(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeAuthorizationCode, "default")
	expiry := time.Now().Add(time.Hour)

	markerJSON, err := json.Marshal(map[string]any{
		"split":        true,
		"token_expiry": expiry.Format(time.RFC3339Nano),
		"has_refresh":  true,
	})
	if err != nil {
		t.Fatalf("failed to marshal marker: %v", err)
	}

	if err := keyring.Set("pingcli", account, string(markerJSON)); err != nil {
		t.Fatalf("failed to seed split marker: %v", err)
	}
	if err := keyring.Set("pingcli", account+"#access", "split-access-token"); err != nil {
		t.Fatalf("failed to seed split access entry: %v", err)
	}

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "split-access-token" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "split-access-token")
	}
	if result.Expiry.IsZero() {
		t.Error("Expiry is zero, want the marker expiry")
	}
}

func TestLoad_SplitMarkerWithoutAccessEntry(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeDeviceCode, "default")

	markerJSON, err := json.Marshal(map[string]any{
		"split":        true,
		"token_expiry": time.Now().Add(time.Hour).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Fatalf("failed to marshal marker: %v", err)
	}

	if err := keyring.Set("pingcli", account, string(markerJSON)); err != nil {
		t.Fatalf("failed to seed split marker: %v", err)
	}

	_, err = pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want a missing-access-entry error")
	}
}

func TestLoad_NoSession(t *testing.T) {
	keyring.MockInit()

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want a not-found error")
	}

	assertErrorContains(t, err, "no Ping CLI authentication session", "pingcli pingone auth login")
}

func TestLoad_ExpiredSession(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")
	seedKeychainSession(t, account, "expired-token", time.Now().Add(-time.Hour))

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want an expired error")
	}

	assertErrorContains(t, err, "expired", "pingcli pingone auth login")
}

func TestLoad_ZeroExpiryAccepted(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")

	tokenJSON, err := json.Marshal(map[string]any{
		"access_token": "no-expiry-token",
		"token_type":   "Bearer",
	})
	if err != nil {
		t.Fatalf("failed to marshal token: %v", err)
	}

	if err := keyring.Set("pingcli", account, string(tokenJSON)); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}

	result, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "no-expiry-token" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "no-expiry-token")
	}
}

func TestLoad_MissingProfile(t *testing.T) {
	keyring.MockInit()

	_, err := pingcli.Load(context.Background(), pingcli.Plan{})
	if err == nil {
		t.Fatal("Load returned nil error, want a missing-profile error")
	}

	assertErrorContains(t, err, "no Ping CLI profile was specified", "pingcli_profile", "PINGONE_TF_PINGCLI_PROFILE")
}

func TestLoad_ProfileWithUnsupportedCharacters(t *testing.T) {
	keyring.MockInit()

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "../evil"})
	if err == nil {
		t.Fatal("Load returned nil error, want an unsupported-characters error")
	}

	assertErrorContains(t, err, "unsupported characters")
}

func TestLoad_KeychainUnavailable(t *testing.T) {
	keyring.MockInitWithError(fmt.Errorf("keychain is locked"))

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want a keychain error")
	}

	assertErrorContains(t, err, "OS keychain")
}

func TestLoad_ErrorDoesNotContainToken(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")
	seedKeychainSession(t, account, "very-secret-token-value", time.Now().Add(-time.Hour))

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want an expired error")
	}

	if got := err.Error(); strings.Contains(got, "very-secret-token-value") {
		t.Errorf("error message contains the token: %q", got)
	}
}

func TestSessionKeyLayout(t *testing.T) {
	// Pin the key derivation to the names observed in real Ping CLI storage.
	tests := []struct {
		grant   string
		profile string
		want    string
	}{
		{pingcli.GrantTypeClientCredentials, "default", "auth_session_token-a1a56d96_pingone_client_credentials_default"},
		{pingcli.GrantTypeAuthorizationCode, "default", "auth_session_token-5a2aa136_pingone_authorization_code_default"},
		{pingcli.GrantTypeDeviceCode, "default", "auth_session_token-00f8f6b5_pingone_device_code_default"},
		{pingcli.GrantTypeClientCredentials, "integration", "auth_session_token-a1a56d96_pingone_client_credentials_integration"},
	}

	for _, tt := range tests {
		if got := sessionAccount(t, tt.grant, tt.profile); got != tt.want {
			t.Errorf("account for grant %q profile %q = %q, want %q", tt.grant, tt.profile, got, tt.want)
		}
	}
}

func assertErrorContains(t *testing.T, err error, wants ...string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	for _, want := range wants {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err.Error(), want)
		}
	}
}

func TestLoad_UnrecognisedKeychainEntry(t *testing.T) {
	keyring.MockInit()

	// Valid JSON that carries neither an access token nor the split marker, as
	// a Ping CLI layout this provider does not know about would.
	if err := keyring.Set("pingcli", sessionAccount(t, pingcli.GrantTypeDeviceCode, "default"), `{"session":{"nested":"token"}}`); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}

	_, err := pingcli.Load(context.Background(), pingcli.Plan{PingCLIProfile: "default"})
	if err == nil {
		t.Fatal("Load returned nil error, want an unrecognised format error")
	}

	assertErrorContains(t, err, "format this provider does not recognise", "'default'", "newer Ping CLI")
}

func TestLoad_TrimsProfileAndGrantType(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")
	seedKeychainSession(t, account, "trimmed-profile-token", time.Now().Add(time.Hour))

	// Entry points that build a Plan directly (the acctest pre-checks) may pass
	// untrimmed environment values; Load must accept the same values the
	// provider path does.
	result, err := pingcli.Load(context.Background(), pingcli.Plan{
		PingCLIProfile:   " default ",
		PingCLIGrantType: " " + pingcli.GrantTypeClientCredentials + " ",
	})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.AccessToken != "trimmed-profile-token" {
		t.Errorf("AccessToken = %q, want %q", result.AccessToken, "trimmed-profile-token")
	}
}
