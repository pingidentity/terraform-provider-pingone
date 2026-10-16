// Copyright © 2026 Ping Identity Corporation

package pingcli

import (
	"fmt"
	"os"
	"slices"
)

// Attribute names in the provider schema.
const (
	AttrConfigProfile   = "config_profile"
	AttrConfigGrantType = "config_grant_type"
)

// Environment variables for the provider attributes. Ping CLI reads the same
// names, so an exported Ping CLI setup is picked up automatically.
const (
	EnvConfigProfile   = "PINGCLI_PROFILE"
	EnvConfigGrantType = "PINGCLI_PINGONE_AUTHENTICATION_TYPE"
)

// ConflictEnvVars are the PingOne credential environment variables checked
// when deciding whether Ping CLI authentication conflicts with an existing
// credential setup. A slice, not named constants, to avoid gosec G101
// flagging credential-looking identifiers.
var ConflictEnvVars = []string{
	"PINGONE_CLIENT_ID",
	"PINGONE_CLIENT_SECRET",
	"PINGONE_ENVIRONMENT_ID",
	"PINGONE_API_ACCESS_TOKEN",
}

// maskedEnvVars are the environment variables unset when Ping CLI
// authentication is enabled: the credential variables above, which the SDKs
// read directly, plus the new SDK's own grant type and token storage
// variables. Those reach the SDK through envconfig after the provider builds
// its configuration, so setting the builder values does not help: a
// PINGONE_AUTH_GRANT_TYPE of client_credentials makes the SDK's token source
// skip the static Ping CLI token, and the storage variables would redirect it
// to the SDK's own keychain flow.
var maskedEnvVars = append(slices.Clone(ConflictEnvVars),
	"PINGONE_AUTH_GRANT_TYPE",
	"PINGONE_STORAGE_TYPE",
	"PINGONE_STORAGE_NAME",
	"PINGONE_STORAGE_OPTIONS",
	"PINGONE_STORAGE_OPTIONAL_SUFFIX",
)

// MaskCredentialEnv unsets maskedEnvVars in this process so the SDKs cannot
// combine them with the Ping CLI token. The caller's shell is unaffected.
func MaskCredentialEnv() {
	for _, name := range maskedEnvVars {
		_ = os.Unsetenv(name)
	}
}

// UnknownValueSummary and UnknownValueDetail report a Ping CLI attribute unknown
// at configure time; treating it as unset could select the wrong session.
const UnknownValueSummary = "Unknown Ping CLI configuration value"

func UnknownValueDetail(attribute string) string {
	return fmt.Sprintf(
		"The %s value is not known when the provider is configured, so the Ping CLI session cannot be selected.  Provider configuration cannot depend on values that are only known after apply; use a literal value, an input variable or the matching environment variable instead.",
		attribute,
	)
}

// DescriptionConfigProfile returns the provider attribute description, shared
// by every muxed provider flavour so the schemas cannot drift apart.
func DescriptionConfigProfile() string {
	return fmt.Sprintf(
		"The Ping CLI profile whose keychain authentication session is used by the provider.  Setting this attribute makes the provider authenticate with the access token that Ping CLI stores in the OS keychain after `%s`, instead of requiring `client_id` and `client_secret`.  The session is looked up in the keychain by profile name; the Ping CLI configuration file is not read.  Note that Ping CLI stores an authentication session under the profile that was active when the login was performed, which may differ from the profile named with `--profile`.  Default value can be set with the `%s` environment variable.  Must not be configured with `client_id`, `client_secret`, `environment_id` or `api_access_token`.  Must be configured with `region_code` (or the `PINGONE_REGION_CODE` environment variable).",
		"pingcli pingone auth login", EnvConfigProfile,
	)
}

// DescriptionConfigGrantType returns the provider attribute description, shared
// by every muxed provider flavour so the schemas cannot drift apart.
func DescriptionConfigGrantType() string {
	return fmt.Sprintf(
		"The grant type that Ping CLI used to log in, which is part of the keychain entry that stores the authentication session.  Options are `%s`, `%s` and `%s`.  When unset, the provider checks each grant type for the profile named by `%s` and uses the valid session that expires last; set this attribute when a profile holds sessions for more than one grant type.  Default value can be set with the `%s` environment variable.  Must not be configured with `client_id`, `client_secret`, `environment_id` or `api_access_token`.  Must be configured with `%s` (or the `%s` environment variable) and `region_code` (or the `PINGONE_REGION_CODE` environment variable).",
		GrantTypeClientCredentials, GrantTypeAuthorizationCode, GrantTypeDeviceCode,
		AttrConfigProfile, EnvConfigGrantType, AttrConfigProfile, EnvConfigProfile,
	)
}

// ConflictDescriptionSuffix is appended to the descriptions of the credential
// attributes in every flavour, so the docs and the validation rules agree.
const ConflictDescriptionSuffix = "  Must not be configured with `config_profile` or `config_grant_type`."
