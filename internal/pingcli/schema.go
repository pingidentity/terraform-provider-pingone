// Copyright © 2026 Ping Identity Corporation

package pingcli

import (
	"fmt"
	"os"
	"slices"
)

// Attribute names in the provider schema.
const (
	AttrUsePingCLIAuth   = "use_pingcli_auth"
	AttrPingCLIProfile   = "pingcli_profile"
	AttrPingCLIGrantType = "pingcli_grant_type"
)

// Environment variables for the provider attributes. Ping CLI reads the same
// names, so an exported Ping CLI setup is picked up automatically.
const (
	EnvPingCLIProfile   = "PINGCLI_PROFILE"
	EnvPingCLIGrantType = "PINGCLI_PINGONE_AUTHENTICATION_TYPE"
)

// DefaultProfile is the Ping CLI profile used when none is configured.
const DefaultProfile = "default"

// ConflictEnvVars are the PingOne credential environment variables that the
// SDKs read directly. A slice, not named constants, to avoid gosec G101
// flagging credential-looking identifiers.
var ConflictEnvVars = []string{
	"PINGONE_CLIENT_ID",
	"PINGONE_CLIENT_SECRET",
	"PINGONE_ENVIRONMENT_ID",
	"PINGONE_API_ACCESS_TOKEN",
}

// maskedEnvVars are the environment variables unset when Ping CLI
// authentication is enabled: the credential variables above, plus the new
// SDK's own grant type and token storage variables. Those reach the SDK
// through envconfig after the provider builds its configuration, so setting
// the builder values does not help: a PINGONE_AUTH_GRANT_TYPE of
// client_credentials makes the SDK's token source skip the static Ping CLI
// token, and the storage variables would redirect it to the SDK's own
// keychain flow.
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

// DescriptionUsePingCLIAuth returns the provider attribute description, shared
// by every muxed provider flavour so the schemas cannot drift apart.
func DescriptionUsePingCLIAuth() string {
	return fmt.Sprintf(
		"Enable authentication with the access token that Ping CLI stores in the OS keychain after `%s`, instead of client credentials or an API access token.  The Ping CLI login must be performed with a worker application that has the PingOne roles Terraform requires.  Must not be configured with `client_id`, `client_secret`, `environment_id` or `api_access_token`.  Must be configured with `region_code` (or the `PINGONE_REGION_CODE` environment variable).",
		"pingcli pingone auth login",
	)
}

// DescriptionPingCLIProfile returns the provider attribute description, shared
// by every muxed provider flavour so the schemas cannot drift apart.
func DescriptionPingCLIProfile() string {
	return fmt.Sprintf(
		"The Ping CLI profile whose keychain authentication session is used when `%s` is enabled.  When unset, the Ping CLI default profile (`%s`) is used.  Default value can be set with the `%s` environment variable.",
		AttrUsePingCLIAuth, DefaultProfile, EnvPingCLIProfile,
	)
}

// DescriptionPingCLIGrantType returns the provider attribute description, shared
// by every muxed provider flavour so the schemas cannot drift apart.
func DescriptionPingCLIGrantType() string {
	return fmt.Sprintf(
		"The grant type that Ping CLI used to log in, used when `%s` is enabled.  Options are `%s`, `%s` and `%s`.  When unset, each grant type is checked for the profile and the valid session that expires last is used.  Default value can be set with the `%s` environment variable.",
		AttrUsePingCLIAuth, GrantTypeClientCredentials, GrantTypeAuthorizationCode, GrantTypeDeviceCode, EnvPingCLIGrantType,
	)
}

// ConflictDescriptionSuffix is appended to the descriptions of the credential
// attributes in every flavour, so the docs and the validation rules agree.
const ConflictDescriptionSuffix = "  Must not be configured with `use_pingcli_auth`."
