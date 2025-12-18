// Copyright © 2025 Ping Identity Corporation

package sdkv2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

// New provider function
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{

			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.",
				},
				"environment_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.",
				},
				"api_access_token": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).",
				},
				"config_path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Path to a PingCLI configuration file containing authentication credentials. When set, the provider will read authentication settings from the specified profile in this file. Default value can be set with the `PINGCLI_CONFIG` environment variable. Cannot be used together with `client_id`, `client_secret`, `environment_id`, or `api_access_token`. If set, `config_profile` can optionally specify which profile to use (defaults to the active profile in the config file).",
				},
				"config_profile": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the profile to use from the PingCLI configuration file. If not specified, uses the active profile defined in the config file. Default value can be set with the `PINGCLI_PROFILE` environment variable. Requires `config_path` to be set.",
				},
				"region_code": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The PingOne region to use, which selects the appropriate service endpoints.  Options are `AP` (for Asia-Pacific `.asia` tenants), `AU` (for Asia-Pacific `.com.au` tenants), `CA` (for Canada `.ca` tenants), `EU` (for Europe `.eu` tenants), `NA` (for North America `.com` tenants) and `SG` (for Singapore `.sg` tenants).  Default value can be set with the `PINGONE_REGION_CODE` environment variable.",
				},
				"http_proxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.",
				},
				"append_user_agent": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A custom string value to append to the end of the `User-Agent` header when making API requests to the PingOne service. Default value can be set with the `PINGONE_TF_APPEND_USER_AGENT` environment variable.",
				},
				"global_options": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "A single block containing configuration items to override API behaviours in PingOne.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"population": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "A single block containing configuration items to override population resource settings in PingOne.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"contains_users_force_delete": {
											Type:        schema.TypeBool,
											Optional:    true,
											Description: "Choose whether to force-delete populations that contain users not managed by Terraform. Useful for development and testing use cases, and only applies if the environment that contains the population is of type `SANDBOX`. The platform default is that populations cannot be removed if they contain user data. By default this parameter is set to `false`.",
										},
									},
								},
							},
						},
					},
				},
				"service_endpoints": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "A single block containing configuration items to override the service API endpoints of PingOne.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"auth_hostname": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.",
							},
							"api_hostname": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.",
							},
						},
					},
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"pingone_certificate":                 base.DatasourceCertificate(),
				"pingone_certificate_export":          base.DatasourceCertificateExport(),
				"pingone_certificate_signing_request": base.DatasourceCertificateSigningRequest(),
				"pingone_language":                    base.DatasourceLanguage(),
				"pingone_trusted_email_domain_spf":    base.DatasourceTrustedEmailDomainSPF(),

				"pingone_resource_attribute": sso.DatasourceResourceAttribute(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_authorize_decision_endpoint": authorize.ResourceDecisionEndpoint(),

				"pingone_certificate":                  base.ResourceCertificate(),
				"pingone_certificate_signing_response": base.ResourceCertificateSigningResponse(),
				"pingone_gateway_credential":           base.ResourceGatewayCredential(),
				"pingone_language":                     base.ResourceLanguage(),
				"pingone_language_update":              base.ResourceLanguageUpdate(),

				"pingone_mfa_policy": mfa.ResourceMFAPolicy(),

				"pingone_application_sign_on_policy_assignment": sso.ResourceApplicationSignOnPolicyAssignment(),
				"pingone_sign_on_policy_action":                 sso.ResourceSignOnPolicyAction(),
			},
		}

		p.ConfigureContextFunc = configure(version)

		return p
	}
}

func configure(version string) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		var config client.Config

		// Check if PingCLI config is being used
		configPath := strings.TrimSpace(d.Get("config_path").(string))
		if configPath == "" {
			configPath = strings.TrimSpace(os.Getenv("PINGCLI_CONFIG"))
		}
		// Expand tilde in config_path like ~/.pingcli/config.yaml
		if configPath != "" && strings.HasPrefix(configPath, "~") {
			if home, herr := os.UserHomeDir(); herr == nil && home != "" {
				configPath = filepath.Join(home, strings.TrimPrefix(configPath, "~"))
			}
		}

		configProfile := strings.TrimSpace(d.Get("config_profile").(string))
		if configProfile == "" {
			configProfile = strings.TrimSpace(os.Getenv("PINGCLI_PROFILE"))
		}

		usingPingCLIConfig := configPath != ""

		// If using PingCLI config, load it
		if usingPingCLIConfig {
			profileConfig, err := pingcli.LoadProfileConfig(configPath, configProfile)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("failed to load PingCLI configuration: %w", err))
			}

			// If no profile was specified, get the active profile name
			if configProfile == "" {
				pingCliConfig, err := pingcli.NewConfig(configPath)
				if err != nil {
					return nil, diag.FromErr(fmt.Errorf("failed to load PingCLI configuration: %w", err))
				}
				configProfile, err = pingCliConfig.GetActiveProfile()
				if err != nil {
					return nil, diag.FromErr(fmt.Errorf("failed to get active profile: %w", err))
				}
			}

			storedToken, err := pingcli.LoadStoredToken(profileConfig, configProfile)
			if err != nil || storedToken == nil || !storedToken.Valid() {
				// First fallback: explicit API access token supplied in provider config
				if v, ok := d.Get("api_access_token").(string); ok && strings.TrimSpace(v) != "" {
					config.AccessToken = strings.TrimSpace(v)
					// Clear other fields to enforce bearer-only mode
					config.ClientID = ""
					config.ClientSecret = ""
					config.EnvironmentID = ""
					// Carry region code from profile or provider for endpoint selection
					if rv := strings.TrimSpace(profileConfig.RegionCode); rv != "" {
						rc := management.EnumRegionCode(rv)
						config.RegionCode = &rc
					} else if rv, ok := d.Get("region_code").(string); ok && strings.TrimSpace(rv) != "" {
						rc := management.EnumRegionCode(strings.TrimSpace(rv))
						config.RegionCode = &rc
					}
					// Second fallback: if client credentials are present in the PingCLI profile, use them regardless of grant type
				} else if strings.TrimSpace(profileConfig.ClientID) != "" && strings.TrimSpace(profileConfig.ClientSecret) != "" && strings.TrimSpace(profileConfig.EnvironmentID) != "" {
					config.ClientID = strings.TrimSpace(profileConfig.ClientID)
					config.ClientSecret = strings.TrimSpace(profileConfig.ClientSecret)
					config.EnvironmentID = strings.TrimSpace(profileConfig.EnvironmentID)
					if v := strings.TrimSpace(profileConfig.RegionCode); v != "" {
						rc := management.EnumRegionCode(v)
						config.RegionCode = &rc
					}
					// Third fallback: provider-supplied client credentials even when config_path is set
				} else if pid, pok := d.Get("client_id").(string); pok && strings.TrimSpace(pid) != "" {
					psecret, sok := d.Get("client_secret").(string)
					penv, eok := d.Get("environment_id").(string)
					if sok && eok && strings.TrimSpace(psecret) != "" && strings.TrimSpace(penv) != "" {
						config.ClientID = strings.TrimSpace(pid)
						config.ClientSecret = strings.TrimSpace(psecret)
						config.EnvironmentID = strings.TrimSpace(penv)
						if rv, ok := d.Get("region_code").(string); ok && strings.TrimSpace(rv) != "" {
							rc := management.EnumRegionCode(strings.TrimSpace(rv))
							config.RegionCode = &rc
						}
					} else {
						return nil, diag.Errorf("PingCLI configuration requires a valid stored token. Please run 'pingcli login' first. Error: %v", err)
					}
				} else {
					return nil, diag.Errorf("PingCLI configuration requires a valid stored token. Please run 'pingcli login' first. Error: %v", err)
				}
			}

			// If we have a stored token, prefer bearer-only mode.
			if storedToken != nil && storedToken.Valid() {
				config.AccessToken = storedToken.AccessToken
				config.ClientID = ""
				config.ClientSecret = ""
				config.EnvironmentID = ""
				if profileConfig.RegionCode != "" {
					regionCode := management.EnumRegionCode(profileConfig.RegionCode)
					config.RegionCode = &regionCode
				}
			}
		} else {
			// Use explicit credentials from provider configuration or environment variables
			// Set the defaults
			if v, ok := d.Get("client_id").(string); ok && v != "" {
				config.ClientID = v
			}

			if v, ok := d.Get("client_secret").(string); ok && v != "" {
				config.ClientSecret = v
			}

			if v, ok := d.Get("environment_id").(string); ok && v != "" {
				config.EnvironmentID = v
			}

			if v, ok := d.Get("api_access_token").(string); ok && v != "" {
				// When a static token is provided, clear client credentials and environment ID
				// to satisfy underlying SDK validation rules.
				config.AccessToken = v
				config.ClientID = ""
				config.ClientSecret = ""
				config.EnvironmentID = ""
			}

			if v, ok := d.Get("region_code").(string); ok && v != "" {
				regionCode := management.EnumRegionCode(v)
				config.RegionCode = &regionCode
			}
		}

		config.GlobalOptions = &client.GlobalOptions{
			Population: &client.PopulationOptions{
				ContainsUsersForceDelete: false,
			},
		}

		if v, ok := d.Get("http_proxy").(string); ok && v != "" {
			config.ProxyURL = &v
		}

		if v, ok := d.Get("service_endpoints").([]interface{}); ok && len(v) > 0 && v[0] != nil {
			vp, ok := v[0].(map[string]interface{})
			if !ok {
				return nil, diag.Errorf("service_endpoints must be a map.  This is always an error in the provider code, please raise an issue with the provider maintainers.")
			}

			if v, ok := vp["auth_hostname"].(string); ok && v != "" {
				config.AuthHostnameOverride = &v
			}

			if v, ok := vp["api_hostname"].(string); ok && v != "" {
				config.APIHostnameOverride = &v
			}
		}

		if v, ok := d.Get("global_options").([]interface{}); ok && len(v) > 0 && v[0] != nil {
			vp, ok := v[0].(map[string]interface{})
			if !ok {
				return nil, diag.Errorf("global_options must be a map.  This is always an error in the provider code, please raise an issue with the provider maintainers.")
			}

			if v1, ok := vp["population"].([]interface{}); ok && len(v1) > 0 && v1[0] != nil {
				v1p, ok := v1[0].(map[string]interface{})
				if !ok {
					return nil, diag.Errorf("global_options.population must be a map.  This is always an error in the provider code, please raise an issue with the provider maintainers.")
				}

				if v2, ok := v1p["contains_users_force_delete"].(bool); ok {
					config.GlobalOptions.Population.ContainsUsersForceDelete = v2
				}
			}
		}

		if v, ok := d.Get("append_user_agent").(string); ok && v != "" {
			config.UserAgentAppend = &v
		} else if v := strings.TrimSpace(os.Getenv("PINGONE_TF_APPEND_USER_AGENT")); v != "" {
			config.UserAgentAppend = &v
		}

		client, err := config.APIClient(ctx, version)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
