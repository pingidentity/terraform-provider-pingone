// Copyright © 2026 Ping Identity Corporation

// Package pingcli authenticates the provider with the access token that Ping CLI
// (https://github.com/pingidentity/pingcli) stores in the OS keychain at login.
// It derives the keychain account from the profile name and grant type; it never
// reads the Ping CLI config file, runs pingcli, or refreshes or writes sessions.
package pingcli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pingidentity/pingone-go-client/oauth2"
	"github.com/zalando/go-keyring"
)

const (
	// keychainServiceName is the service under which Ping CLI stores
	// authentication sessions in the OS keychain.
	keychainServiceName = "pingcli"

	// keychainAccountPrefix prefixes the session key to form the keychain
	// account name.
	keychainAccountPrefix = "auth_session_"

	// splitEntrySuffix names the entry holding the access token of a session Ping
	// CLI split across entries (v1.9.0+).
	splitEntrySuffix = "#access"

	// serviceName is the Ping CLI service the session belongs to; it is part
	// of the session key.
	serviceName = "pingone"

	// expirySkew is subtracted from the expiry so a nearly expired token is not
	// used for a long apply.
	expirySkew = 30 * time.Second
)

// Grant types Ping CLI can log in with. The session key includes the grant type,
// which is unknown without the config file, so each one is probed.
const (
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeDeviceCode        = "device_code"
)

// grantTypes lists the grant types probed, in order of preference when more
// than one holds a session with the same expiry.
var grantTypes = []string{
	GrantTypeClientCredentials,
	GrantTypeAuthorizationCode,
	GrantTypeDeviceCode,
}

// IsGrantType reports whether grantType is one Ping CLI can log in with.
func IsGrantType(grantType string) bool {
	return slices.Contains(grantTypes, grantType)
}

// profileNamePattern restricts profile names to the characters Ping CLI itself
// allows. It also keeps profile names from escaping the token file directory.
var profileNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// SessionResult is a loaded Ping CLI authentication session.
type SessionResult struct {
	// AccessToken is the bearer token issued by the Ping CLI login.
	AccessToken string
	// Expiry is when the access token expires. A zero value means the stored
	// session did not carry a usable expiry.
	Expiry time.Time
	// GrantType is the grant type of the session, taken from the keychain
	// account under which it was found.
	GrantType string
	// Profile is the profile name under which the session was found.
	Profile string
}

// splitMarker mirrors the marker Ping CLI writes when a session is too large for
// one keychain entry; the access token then lives in a separate entry.
type splitMarker struct {
	Split       bool      `json:"split"`
	TokenExpiry time.Time `json:"token_expiry"`
}

// Load returns the freshest valid session for the profile, checking the requested
// grant type or all known ones. It only reads: a missing or expired session is
// an error with guidance, never refreshed.
func Load(ctx context.Context, plan Plan) (SessionResult, error) {
	// Trim like Resolve does, so entry points that build a Plan directly (the
	// acctest pre-checks) accept the same values as the provider path.
	profile := strings.TrimSpace(plan.PingCLIProfile)
	if profile == "" {
		return SessionResult{}, fmt.Errorf(
			"no Ping CLI profile was specified; set %s in the provider configuration or the %s environment variable",
			AttrPingCLIProfile, EnvPingCLIProfile,
		)
	}
	if !profileNamePattern.MatchString(profile) {
		return SessionResult{}, fmt.Errorf(
			"the Ping CLI profile name '%s' contains unsupported characters; a profile name can only contain letters, numbers, '-' and '_', so set %s or %s to a valid profile name",
			profile, AttrPingCLIProfile, EnvPingCLIProfile,
		)
	}

	candidates := grantTypes
	if configGrantType := strings.TrimSpace(plan.PingCLIGrantType); configGrantType != "" {
		if !IsGrantType(configGrantType) {
			return SessionResult{}, fmt.Errorf(
				"the Ping CLI grant type '%s' is not supported; use one of %s, %s or %s",
				configGrantType, GrantTypeClientCredentials, GrantTypeAuthorizationCode, GrantTypeDeviceCode,
			)
		}
		candidates = []string{configGrantType}
	}

	var (
		best         SessionResult
		found        bool
		expired      bool
		unreadable   bool
		unrecognised bool
		hardErr      error
	)

	for _, grantType := range candidates {
		account := keychainAccountPrefix + sessionKey(serviceName, grantType, profile)

		accessToken, expiry, err := loadKeychainSession(account)
		if err != nil {
			switch {
			case errors.Is(err, keyring.ErrNotFound):
			case errors.Is(err, errUnrecognisedEntry):
				unrecognised = true
			case hardErr == nil:
				hardErr = err
			}
			continue
		}

		if accessToken == "" {
			unreadable = true
			continue
		}

		if !expiry.IsZero() && time.Now().After(expiry.Add(-expirySkew)) {
			expired = true
			continue
		}

		// Prefer the session that stays valid longest; an unknown expiry ranks
		// below any known one.
		if !found || expiry.After(best.Expiry) {
			best = SessionResult{AccessToken: accessToken, Expiry: expiry, GrantType: grantType, Profile: profile}
			found = true
		}
	}

	switch {
	case found:
		tflog.Debug(ctx, "Loaded Ping CLI authentication session", map[string]any{
			"profile":    profile,
			"grant_type": best.GrantType,
			"expiry":     best.Expiry.Format(time.RFC3339),
		})

		return best, nil
	case hardErr != nil:
		return SessionResult{}, fmt.Errorf(
			"failed to read the Ping CLI authentication session from the OS keychain (%w); the provider only reads sessions stored in the OS keychain, so run Terraform where one is available, or authenticate with client credentials or an access token instead",
			hardErr,
		)
	case expired:
		return SessionResult{}, fmt.Errorf(
			"the Ping CLI authentication session in the OS keychain for profile '%s' has expired; run 'pingcli pingone auth login' to authenticate again",
			profile,
		)
	case unrecognised:
		return SessionResult{}, fmt.Errorf(
			"the Ping CLI keychain entry for profile '%s' has a format this provider does not recognise (it has neither an access token nor the split-session marker); it may have been written by a newer Ping CLI than this provider supports, so update the provider, or authenticate with client credentials or an access token instead",
			profile,
		)
	case unreadable:
		return SessionResult{}, fmt.Errorf(
			"the Ping CLI authentication session in the OS keychain for profile '%s' is empty or unreadable; run 'pingcli pingone auth login' to authenticate",
			profile,
		)
	case plan.PingCLIGrantType != "":
		return SessionResult{}, fmt.Errorf(
			"no Ping CLI authentication session was found in the OS keychain for profile '%s' and grant type '%s'; run 'pingcli pingone auth login' to authenticate, and note that Ping CLI stores the session under the profile that was active when the login was performed",
			profile, plan.PingCLIGrantType,
		)
	default:
		return SessionResult{}, fmt.Errorf(
			"no Ping CLI authentication session was found in the OS keychain for profile '%s'; run 'pingcli pingone auth login' to authenticate, and note that Ping CLI stores the session under the profile that was active when the login was performed",
			profile,
		)
	}
}

// sessionKey derives Ping CLI v1.x's session key:
// token-<first 4 bytes of sha256("service:grantType") in hex>_service_grantType_profile.
// The keychain account is this prefixed with auth_session_.
func sessionKey(service, grantType, profile string) string {
	h := sha256.Sum256([]byte(service + ":" + grantType))

	return fmt.Sprintf("token-%x_%s_%s_%s", h[:4], service, grantType, profile)
}

// errUnrecognisedEntry marks a keychain entry with neither an access token nor the
// split-session marker, as an unknown Ping CLI layout would look.
var errUnrecognisedEntry = errors.New("the keychain entry has neither an access token nor a split-session marker")

// loadKeychainSession reads a session through the SDK's keychain storage,
// handling both the single-entry and the split layout (Ping CLI v1.9.0+).
func loadKeychainSession(account string) (string, time.Time, error) {
	storage, err := oauth2.NewKeychainStorage(keychainServiceName, account)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to prepare keychain access: %w", err)
	}

	token, err := storage.LoadToken()
	switch {
	case err == nil:
		if token != nil && token.AccessToken != "" {
			return token.AccessToken, token.Expiry, nil
		}
		// Fall through: an entry without a token may be a split marker.
	case errors.Is(err, keyring.ErrNotFound):
		return "", time.Time{}, err
	default:
		return "", time.Time{}, err
	}

	// Split-session layout: the main entry holds a marker without tokens and
	// the access token lives in its own entry.
	raw, err := keyring.Get(keychainServiceName, account)
	if err != nil {
		return "", time.Time{}, err
	}

	var marker splitMarker
	if err := json.Unmarshal([]byte(raw), &marker); err != nil || !marker.Split {
		return "", time.Time{}, errUnrecognisedEntry
	}

	accessRaw, err := keyring.Get(keychainServiceName, account+splitEntrySuffix)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessRaw, marker.TokenExpiry, nil
}
