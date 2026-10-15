// Copyright © 2026 Ping Identity Corporation

package pingcli

import "strings"

// Inputs carries the provider configuration values needed to decide whether
// Ping CLI authentication is enabled.
type Inputs struct {
	// ConfigProfile is the HCL value of config_profile (empty when unset).
	ConfigProfile string
	// ConfigGrantType is the HCL value of config_grant_type (empty when unset).
	ConfigGrantType string
	// RegionCode is the resolved region_code value (HCL or environment).
	RegionCode string
	// HasHCLCredential reports whether any of client_id, client_secret,
	// environment_id or api_access_token is set in the configuration.
	HasHCLCredential bool
}

// Plan describes the outcome of resolution.
type Plan struct {
	// Enabled reports whether the provider should authenticate with a Ping CLI
	// session.
	Enabled bool
	// ConfigProfile is the resolved Ping CLI profile name; always set when the
	// plan is enabled.
	ConfigProfile string
	// ConfigGrantType is the resolved grant type (may be empty, meaning every
	// known grant type is checked).
	ConfigGrantType string
}

// Problem is a provider-facing error produced during resolution.
type Problem struct {
	Summary string
	Detail  string
}

// Resolve decides whether Ping CLI authentication is enabled and validates the
// settings. getenv is injectable for tests.
func Resolve(in Inputs, getenv func(string) string) (Plan, []Problem) {
	hclProfile := strings.TrimSpace(in.ConfigProfile)
	envProfile := strings.TrimSpace(getenv(EnvConfigProfile))

	hclGrantType := strings.TrimSpace(in.ConfigGrantType)
	envGrantType := strings.TrimSpace(getenv(EnvConfigGrantType))

	explicit := hclProfile != "" || hclGrantType != ""

	envCredential := false
	for _, name := range ConflictEnvVars {
		if strings.TrimSpace(getenv(name)) != "" {
			envCredential = true
		}
	}

	// Implicit activation needs PINGCLI_PROFILE and no credentials anywhere, so an
	// exported Ping CLI variable never breaks an existing credentials setup.
	implicit := !explicit && envProfile != "" && !in.HasHCLCredential && !envCredential

	if !explicit && !implicit {
		return Plan{}, nil
	}

	var problems []Problem

	// Exported PINGONE_* credentials never conflict with an explicit attribute: the
	// configuration wins and MaskCredentialEnv drops them. Only credentials in the
	// provider configuration contradict it.
	if in.HasHCLCredential {
		problems = append(problems, Problem{
			Summary: "Conflicting provider authentication configuration",
			Detail: "The provider cannot use both a Ping CLI authentication session and PingOne client credentials or an API access token.\n\n" +
				"Unset " + AttrConfigProfile + " and " + AttrConfigGrantType +
				" to authenticate with client credentials or an access token, or remove client_id, client_secret, environment_id and api_access_token from the provider configuration to authenticate with Ping CLI.",
		})
	}

	if strings.TrimSpace(in.RegionCode) == "" {
		problems = append(problems, Problem{
			Summary: "Missing region code",
			Detail:  "When the provider authenticates with a Ping CLI authentication session, the PingOne region must be configured with the region_code attribute or the PINGONE_REGION_CODE environment variable, as the Ping CLI configuration file does not reliably carry the region of the authentication session.",
		})
	}

	// The HCL values win; each falls back to its environment variable.
	configProfile := hclProfile
	if configProfile == "" {
		configProfile = envProfile
	}

	configGrantType := hclGrantType
	if configGrantType == "" {
		configGrantType = envGrantType
	}

	if configProfile == "" {
		problems = append(problems, Problem{
			Summary: "Missing Ping CLI profile",
			Detail:  "When the provider authenticates with a Ping CLI authentication session, the Ping CLI profile must be configured with the " + AttrConfigProfile + " attribute or the " + EnvConfigProfile + " environment variable.",
		})
	}

	if configGrantType != "" && !IsGrantType(configGrantType) {
		problems = append(problems, Problem{
			Summary: "Invalid Ping CLI grant type",
			Detail:  "The Ping CLI grant type '" + configGrantType + "' is not supported.  Use one of " + GrantTypeClientCredentials + ", " + GrantTypeAuthorizationCode + " or " + GrantTypeDeviceCode + ".",
		})
	}

	return Plan{
		Enabled:         true,
		ConfigProfile:   configProfile,
		ConfigGrantType: configGrantType,
	}, problems
}
