// Copyright © 2026 Ping Identity Corporation

package pingcli

import "strings"

// Inputs carries the provider configuration values needed to decide whether
// Ping CLI authentication is enabled.
type Inputs struct {
	// UsePingCLIAuth is the HCL value of use_pingcli_auth.
	UsePingCLIAuth bool
	// PingCLIProfile is the HCL value of pingcli_profile (empty when unset).
	PingCLIProfile string
	// PingCLIGrantType is the HCL value of pingcli_grant_type (empty when
	// unset).
	PingCLIGrantType string
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
	// PingCLIProfile is the resolved Ping CLI profile name.
	PingCLIProfile string
	// PingCLIGrantType is the resolved grant type (may be empty, meaning every
	// known grant type is checked).
	PingCLIGrantType string
}

// Problem is a provider-facing error produced during resolution.
type Problem struct {
	Summary string
	Detail  string
}

// Resolve decides whether Ping CLI authentication is enabled and validates the
// settings. getenv is injectable for tests.
func Resolve(in Inputs, getenv func(string) string) (Plan, []Problem) {
	if !in.UsePingCLIAuth {
		return Plan{}, nil
	}

	var problems []Problem

	// The credential attributes also take their values from environment
	// variables, so an exported credential is a conflict just like one in the
	// configuration.
	var envCredentials []string
	for _, name := range ConflictEnvVars {
		if strings.TrimSpace(getenv(name)) != "" {
			envCredentials = append(envCredentials, name)
		}
	}

	if in.HasHCLCredential || len(envCredentials) > 0 {
		detail := "The provider cannot use both a Ping CLI authentication session and PingOne client credentials or an API access token.\n\n"
		if in.HasHCLCredential {
			detail += "Remove client_id, client_secret, environment_id and api_access_token from the provider configuration. "
		}
		if len(envCredentials) > 0 {
			detail += "Unset the " + strings.Join(envCredentials, ", ") + " environment variable(s). "
		}
		detail += "Alternatively, unset " + AttrUsePingCLIAuth + " to authenticate with client credentials or an access token."

		problems = append(problems, Problem{
			Summary: "Conflicting provider authentication configuration",
			Detail:  detail,
		})
	}

	if strings.TrimSpace(in.RegionCode) == "" {
		problems = append(problems, Problem{
			Summary: "Missing region code",
			Detail:  "When the provider authenticates with a Ping CLI authentication session, the PingOne region must be configured with the region_code attribute or the PINGONE_REGION_CODE environment variable, as the Ping CLI configuration file does not reliably carry the region of the authentication session.",
		})
	}

	profile := strings.TrimSpace(in.PingCLIProfile)
	if profile == "" {
		profile = strings.TrimSpace(getenv(EnvPingCLIProfile))
	}
	if profile == "" {
		profile = DefaultProfile
	}

	grantType := strings.TrimSpace(in.PingCLIGrantType)
	if grantType == "" {
		grantType = strings.TrimSpace(getenv(EnvPingCLIGrantType))
	}
	if grantType != "" && !IsGrantType(grantType) {
		problems = append(problems, Problem{
			Summary: "Invalid Ping CLI grant type",
			Detail:  "The Ping CLI grant type '" + grantType + "' is not supported.  Use one of " + GrantTypeClientCredentials + ", " + GrantTypeAuthorizationCode + " or " + GrantTypeDeviceCode + ".",
		})
	}

	return Plan{
		Enabled:          true,
		PingCLIProfile:   profile,
		PingCLIGrantType: grantType,
	}, problems
}
