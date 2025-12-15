// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	clientconfig "github.com/pingidentity/pingone-go-client/config"
	"github.com/pingidentity/pingone-go-client/oauth2"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/davinci"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Ensure PingOneProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &pingOneProvider{}
)

// PingOneProvider defines the provider implementation.
type pingOneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// pingOneProviderModel describes the provider data model.
type pingOneProviderModel struct {
	ClientID         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
	EnvironmentID    types.String `tfsdk:"environment_id"`
	APIAccessToken   types.String `tfsdk:"api_access_token"`
	ConfigPath       types.String `tfsdk:"config_path"`
	ConfigProfile    types.String `tfsdk:"config_profile"`
	AppendUserAgent  types.String `tfsdk:"append_user_agent"`
	RegionCode       types.String `tfsdk:"region_code"`
	ServiceEndpoints types.List   `tfsdk:"service_endpoints"`
	GlobalOptions    types.List   `tfsdk:"global_options"`
	HTTPProxy        types.String `tfsdk:"http_proxy"`
}

type pingOneProviderGlobalOptionsModel struct {
	Population types.List `tfsdk:"population"`
}

type pingOneProviderGlobalOptionsPopulationModel struct {
	ContainsUsersForceDelete types.Bool `tfsdk:"contains_users_force_delete"`
}

type pingOneProviderServiceEndpointsModel struct {
	AuthHostname types.String `tfsdk:"auth_hostname"`
	APIHostname  types.String `tfsdk:"api_hostname"`
}

func (p *pingOneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingone"
	resp.Version = p.version
}

func (p *pingOneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

	clientIDDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.",
	)

	clientSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.",
	)

	environmentIDDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.",
	)

	apiAccessTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).",
	)

	configPathDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Path to a PingCLI configuration file containing authentication credentials. When set, the provider will read authentication settings from the specified profile in this file. Default value can be set with the `PINGCLI_CONFIG` environment variable. Cannot be used together with `client_id`, `client_secret`, `environment_id`, or `api_access_token`. If set, `config_profile` can optionally specify which profile to use (defaults to the active profile in the config file).",
	)

	configProfileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the profile to use from the PingCLI configuration file. If not specified, uses the active profile defined in the config file. Default value can be set with the `PINGCLI_PROFILE` environment variable. Requires `config_path` to be set.",
	)

	regionCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The PingOne region to use, which selects the appropriate service endpoints.  Options are `AP` (for Asia-Pacific `.asia` tenants), `AU` (for Asia-Pacific `.com.au` tenants), `CA` (for Canada `.ca` tenants), `EU` (for Europe `.eu` tenants), `NA` (for North America `.com` tenants) and `SG` (for Singapore `.sg` tenants).  Default value can be set with the `PINGONE_REGION_CODE` environment variable.",
	)

	globalOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing configuration items to override API behaviours in PingOne.",
	)

	globalOptionsPopulationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing configuration items to override population resource settings in PingOne.",
	)

	globalOptionsEnvironmentContainsUsersForceDeleteDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Choose whether to force-delete populations that contain users not managed by Terraform. Useful for development and testing use cases, and only applies if the environment that contains the population is of type `SANDBOX`. The platform default is that populations cannot be removed if they contain user data. By default this parameter is set to `false`.",
	)

	serviceEndpointsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing configuration items to override the service API endpoints of PingOne.",
	)

	serviceEndpointsAuthHostnameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.",
	)

	serviceEndpointsApiHostnameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.",
	)

	httpProxyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.",
	)

	appendUserAgentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A custom string value to append to the end of the `User-Agent` header when making API requests to the PingOne service. Default value can be set with the `PINGONE_TF_APPEND_USER_AGENT` environment variable.",
	)

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description:         clientIDDescription.Description,
				MarkdownDescription: clientIDDescription.MarkdownDescription,
				Optional:            true,
			},

			"client_secret": schema.StringAttribute{
				Description:         clientSecretDescription.Description,
				MarkdownDescription: clientSecretDescription.MarkdownDescription,
				Optional:            true,
			},

			"environment_id": schema.StringAttribute{
				Description:         environmentIDDescription.Description,
				MarkdownDescription: environmentIDDescription.MarkdownDescription,
				Optional:            true,
			},

			"api_access_token": schema.StringAttribute{
				Description:         apiAccessTokenDescription.Description,
				MarkdownDescription: apiAccessTokenDescription.MarkdownDescription,
				Optional:            true,
			},

			"config_path": schema.StringAttribute{
				Description:         configPathDescription.Description,
				MarkdownDescription: configPathDescription.MarkdownDescription,
				Optional:            true,
			},

			"config_profile": schema.StringAttribute{
				Description:         configProfileDescription.Description,
				MarkdownDescription: configProfileDescription.MarkdownDescription,
				Optional:            true,
			},

			"region_code": schema.StringAttribute{
				Description:         regionCodeDescription.Description,
				MarkdownDescription: regionCodeDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(pingone.AllowedEnvironmentRegionCodeEnumValues)...),
				},
			},

			"http_proxy": schema.StringAttribute{
				Description:         httpProxyDescription.Description,
				MarkdownDescription: httpProxyDescription.MarkdownDescription,
				Optional:            true,
			},

			"append_user_agent": schema.StringAttribute{
				Description:         appendUserAgentDescription.Description,
				MarkdownDescription: appendUserAgentDescription.MarkdownDescription,
				Optional:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"global_options": schema.ListNestedBlock{
				Description:         globalOptionsDescription.Description,
				MarkdownDescription: globalOptionsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Blocks: map[string]schema.Block{
						"population": schema.ListNestedBlock{
							Description:         globalOptionsPopulationDescription.Description,
							MarkdownDescription: globalOptionsPopulationDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"contains_users_force_delete": schema.BoolAttribute{
										Description:         globalOptionsEnvironmentContainsUsersForceDeleteDescription.Description,
										MarkdownDescription: globalOptionsEnvironmentContainsUsersForceDeleteDescription.MarkdownDescription,
										Optional:            true,
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"service_endpoints": schema.ListNestedBlock{
				Description:         serviceEndpointsDescription.Description,
				MarkdownDescription: serviceEndpointsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"auth_hostname": schema.StringAttribute{
							Description:         serviceEndpointsAuthHostnameDescription.Description,
							MarkdownDescription: serviceEndpointsAuthHostnameDescription.MarkdownDescription,
							Required:            true,
						},

						"api_hostname": schema.StringAttribute{
							Description:         serviceEndpointsApiHostnameDescription.Description,
							MarkdownDescription: serviceEndpointsApiHostnameDescription.MarkdownDescription,
							Required:            true,
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (p *pingOneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "[v6] Provider configure start")
	var data pingOneProviderModel
	// Initialize a client configuration pointer used throughout setup
	var config *clientconfig.Configuration

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the defaults
	tflog.Info(ctx, "[v6] Provider setting defaults..")

	if v := strings.TrimSpace(os.Getenv("PINGONE_REGION")); v != "" {
		resp.Diagnostics.AddWarning(
			"Deprecated PINGONE_REGION environment variable",
			"The PINGONE_REGION environment variable is now deprecated and should be replaced with the PINGONE_REGION_CODE environment variable.\n\nOptions for the PINGONE_REGION_CODE environment variable are `AP` (`.asia` tenants), `AU` (`.com.au` tenants), `CA` (`.ca` tenants), `EU` (`.eu` tenants) and `NA` (`.com` tenants).",
		)
	}

	// Determine configuration source: PingCLI config file or explicit credentials
	var configPath, configProfile string
	if !data.ConfigPath.IsNull() && !data.ConfigPath.IsUnknown() {
		configPath = strings.TrimSpace(data.ConfigPath.ValueString())
	} else {
		configPath = strings.TrimSpace(os.Getenv("PINGCLI_CONFIG"))
	}

	// Expand tilde to user home for config_path like ~/.pingcli/config.yaml
	if configPath != "" && strings.HasPrefix(configPath, "~") {
		home, herr := os.UserHomeDir()
		if herr == nil && home != "" {
			configPath = filepath.Join(home, strings.TrimPrefix(configPath, "~"))
		}
	}

	if !data.ConfigProfile.IsNull() && !data.ConfigProfile.IsUnknown() {
		configProfile = strings.TrimSpace(data.ConfigProfile.ValueString())
	} else {
		configProfile = strings.TrimSpace(os.Getenv("PINGCLI_PROFILE"))
	}

	// Validate mutual exclusivity
	usingPingCLIConfig := configPath != ""
	usingExplicitCreds := !data.ClientID.IsNull() || !data.ClientSecret.IsNull() || !data.EnvironmentID.IsNull() || !data.APIAccessToken.IsNull()

	if usingPingCLIConfig && usingExplicitCreds {
		resp.Diagnostics.AddError(
			"Conflicting Configuration",
			"Cannot use both PingCLI configuration file (config_path) and explicit credentials (client_id, client_secret, environment_id, api_access_token). Please use only one authentication method.",
		)
		return
	}

	// Validate config_profile requires config_path
	if configProfile != "" && configPath == "" {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"config_profile is set but config_path is not defined. config_profile requires config_path to be set.",
		)
		return
	}

	globalOptions := &client.GlobalOptions{
		Population: &client.PopulationOptions{
			ContainsUsersForceDelete: false,
		},
	}

	var clientID, clientSecret, environmentID, regionCode string
	var grantType oauth2.GrantType
	var redirectURIPath, redirectURIPort string
	var scopes []string
	var pingCliConfig *pingcli.Config
	var activeProfile string
	// When a stored token is found (PingCLI opt-in), enforce bearer-only mode
	useBearerOnly := false

	// Load configuration from PingCLI config file if provided
	if usingPingCLIConfig {
		tflog.Info(ctx, "[v6] Loading configuration from PingCLI config file", map[string]any{
			"config_path":    configPath,
			"config_profile": configProfile,
		})

		// Load the config object so we can check storage settings later
		var err error
		pingCliConfig, err = pingcli.NewConfig(configPath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Load PingCLI Configuration",
				fmt.Sprintf("Error reading PingCLI configuration file: %s", err.Error()),
			)
			return
		}

		// Determine which profile to use
		if configProfile != "" {
			activeProfile = configProfile
		} else {
			activeProfile, err = pingCliConfig.GetActiveProfile()
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to Get Active Profile",
					fmt.Sprintf("Error getting active profile from PingCLI configuration: %s", err.Error()),
				)
				return
			}
		}

		// Get profile configuration
		profileConfig, err := pingCliConfig.GetProfileConfig(activeProfile)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Load Profile Configuration",
				fmt.Sprintf("Error loading profile '%s' from PingCLI configuration: %s", activeProfile, err.Error()),
			)
			return
		}

		// Try to load a stored token first
		storedToken, err := pingcli.LoadStoredToken(profileConfig, activeProfile)
		if err == nil && storedToken != nil && storedToken.Valid() {
			tflog.Info(ctx, "[v6] Found valid stored token from PingCLI credentials", map[string]any{
				"profile":        activeProfile,
				"token_valid":    true,
				"token_expiry":   storedToken.Expiry,
				"has_access":     storedToken.AccessToken != "",
				"has_refresh":    storedToken.RefreshToken != "",
				"region_code":    profileConfig.RegionCode,
				"environment_id": profileConfig.EnvironmentID,
			})
			// Set bearer-only mode values and skip OAuth configuration paths
			useBearerOnly = true
			// Assign environment/region for domain resolution when needed
			environmentID = strings.TrimSpace(profileConfig.EnvironmentID)
			regionCode = strings.TrimSpace(profileConfig.RegionCode)
			// Build configuration immediately with static access token
			cfg := clientconfig.NewConfiguration().WithAccessToken(storedToken.AccessToken)
			if environmentID != "" {
				cfg = cfg.WithEnvironmentID(environmentID)
			}
			// Apply region-derived top-level domain to satisfy client endpoint requirements
			if regionCode != "" {
				if tld, ok := framework.RegionTopLevelDomainFromCode(strings.ToLower(regionCode)); ok {
					cfg = cfg.WithTopLevelDomain(clientconfig.TopLevelDomain(tld))
				}
			}
			// Continue below using cfg instead of OAuth setup
			config = cfg
		} else {
			// No valid stored token, will need to authenticate
			if err != nil {
				tflog.Warn(ctx, "[v6] Could not load stored token, will use auth flow", map[string]any{
					"profile": activeProfile,
					"error":   err.Error(),
				})
			} else if storedToken == nil {
				tflog.Warn(ctx, "[v6] No stored token found, will use auth flow", map[string]any{
					"profile": activeProfile,
				})
			} else if !storedToken.Valid() {
				tflog.Warn(ctx, "[v6] Stored token is expired, will use auth flow", map[string]any{
					"profile":      activeProfile,
					"token_expiry": storedToken.Expiry,
				})
			}
		}

		// Always extract profile configuration for OAuth setup
		// The SDK will check keychain for cached tokens and handle automatic refresh
		clientID = profileConfig.ClientID
		clientSecret = profileConfig.ClientSecret
		environmentID = profileConfig.EnvironmentID
		regionCode = profileConfig.RegionCode
		scopes = profileConfig.Scopes
		redirectURIPath = profileConfig.RedirectURIPath
		redirectURIPort = profileConfig.RedirectURIPort

		// Map PingCLI auth type to SDK grant type
		switch profileConfig.GrantType {
		case "client_credentials":
			grantType = oauth2.GrantTypeClientCredentials
		case "authorization_code":
			grantType = oauth2.GrantTypeAuthorizationCode
		case "device_code":
			grantType = oauth2.GrantTypeDeviceCode
		case "worker":
			// Legacy 'worker' type maps to client credentials flow
			grantType = oauth2.GrantTypeClientCredentials
		default:
			resp.Diagnostics.AddError(
				"Unsupported Grant Type",
				fmt.Sprintf("PingCLI profile uses unsupported grant type '%s'. Supported types are: client_credentials, authorization_code, device_code", profileConfig.GrantType),
			)
			return
		}

		tflog.Debug(ctx, "[v6] Loaded PingCLI configuration", map[string]any{
			"auth_type":      profileConfig.GrantType,
			"grant_type":     grantType,
			"region_code":    profileConfig.RegionCode,
			"environment_id": profileConfig.EnvironmentID,
		})
	} else {
		// Use explicit credentials from provider configuration or environment variables
		// Check if using explicit API access token (static token without OAuth flow)
		if !data.APIAccessToken.IsNull() && !data.APIAccessToken.IsUnknown() && data.APIAccessToken.ValueString() != "" {
			// Static token mode - OAuth configuration not needed
			environmentID = data.EnvironmentID.ValueString()
			useBearerOnly = true
		} else {
			// Use OAuth flow with client credentials
			clientID = data.ClientID.ValueString()
			clientSecret = data.ClientSecret.ValueString()
			environmentID = data.EnvironmentID.ValueString()
			grantType = oauth2.GrantTypeClientCredentials
		}
	}

	// Initialize configuration; may already be set for bearer-only above
	if config == nil {
		config = clientconfig.NewConfiguration()
	}

	// Check if we have a static access token (explicit API token without OAuth)
	if useBearerOnly {
		// Static token mode - no OAuth flow
		// If APIAccessToken provided directly, ensure it's set
		if !usingPingCLIConfig && !data.APIAccessToken.IsNull() && !data.APIAccessToken.IsUnknown() && data.APIAccessToken.ValueString() != "" {
			config = config.WithAccessToken(data.APIAccessToken.ValueString())
		}
		if environmentID != "" {
			config = config.WithEnvironmentID(environmentID)
		}
	} else if grantType != "" {
		// Configure OAuth flow based on grant type
		config = config.
			WithGrantType(grantType).
			WithEnvironmentID(environmentID)

		// Set storage name for keychain-based token caching
		// Use "pingcli" when loading from PingCLI config, "terraform-provider-pingone" otherwise
		storageName := "terraform-provider-pingone"
		if usingPingCLIConfig && activeProfile != "" {
			storageName = "pingcli"
		}
		config = config.WithStorageName(storageName)

		// Use secure local storage (OS keychain) for token caching
		config = config.WithStorageType(clientconfig.StorageTypeSecureLocal)

		// Validate required fields for client credentials early to avoid opaque errors
		if grantType == oauth2.GrantTypeClientCredentials {
			if strings.TrimSpace(clientID) == "" || strings.TrimSpace(environmentID) == "" {
				missing := make([]string, 0, 2)
				if strings.TrimSpace(clientID) == "" {
					missing = append(missing, "client_id")
				}
				if strings.TrimSpace(environmentID) == "" {
					missing = append(missing, "environment_id")
				}
				resp.Diagnostics.AddError(
					"Missing required credentials",
					fmt.Sprintf("Grant type client_credentials requires %s. Review your PingCLI profile (including worker/authentication blocks) or provider inputs.", strings.Join(missing, " and ")),
				)
				return
			}
		}

		switch grantType {
		case oauth2.GrantTypeClientCredentials:
			config = config.
				WithClientID(clientID).
				WithClientSecret(clientSecret)
			if len(scopes) > 0 {
				config = config.WithClientCredentialsScopes(scopes)
			}
		case oauth2.GrantTypeAuthorizationCode:
			config = config.WithAuthorizationCodeClientID(clientID).WithEnvironmentID(environmentID)
			if len(scopes) > 0 {
				config = config.WithAuthorizationCodeScopes(scopes)
			}
			if redirectURIPath != "" || redirectURIPort != "" {
				config = config.WithAuthorizationCodeRedirectURI(clientconfig.AuthorizationCodeRedirectURI{
					Path: redirectURIPath,
					Port: redirectURIPort,
				})
			}
		case oauth2.GrantTypeDeviceCode:
			config = config.WithDeviceCodeClientID(clientID).WithEnvironmentID(environmentID)
			if len(scopes) > 0 {
				config = config.WithDeviceCodeScopes(scopes)
			}
		}
	} else {
		// No access token and no grant type - invalid configuration
		config = config.WithEnvironmentID(environmentID)
	}

	// Region code handling
	if regionCode == "" && !data.RegionCode.IsNull() && !data.RegionCode.IsUnknown() {
		regionCode = strings.TrimSpace(data.RegionCode.ValueString())
	} else if regionCode == "" {
		// Region codes are not handled automatically by the client
		regionCode = strings.TrimSpace(os.Getenv("PINGONE_REGION_CODE"))
	}
	// Region code handling provides minimal endpoint configuration for the client
	if regionCode != "" {
		regionTopLevelDomain, ok := framework.RegionTopLevelDomainFromCode(strings.ToLower(regionCode))
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid Region Code",
				fmt.Sprintf("The region code '%s' is not valid. Valid options are: %s", regionCode, strings.Join(utils.EnumSliceToStringSlice(pingone.AllowedEnvironmentRegionCodeEnumValues), ", ")),
			)
			return
		}
		config = config.WithTopLevelDomain(clientconfig.TopLevelDomain(regionTopLevelDomain))
	}

	if !data.GlobalOptions.IsNull() {

		var globalOptionsData []pingOneProviderGlobalOptionsModel
		resp.Diagnostics.Append(data.GlobalOptions.ElementsAs(ctx, &globalOptionsData, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(globalOptionsData) > 0 {
			if !globalOptionsData[0].Population.IsNull() {

				var globalOptionsPopulationData []pingOneProviderGlobalOptionsPopulationModel
				resp.Diagnostics.Append(globalOptionsData[0].Population.ElementsAs(ctx, &globalOptionsPopulationData, false)...)
				if resp.Diagnostics.HasError() {
					return
				}

				if len(globalOptionsPopulationData) > 0 {
					if !globalOptionsPopulationData[0].ContainsUsersForceDelete.IsNull() {
						globalOptions.Population.ContainsUsersForceDelete = globalOptionsPopulationData[0].ContainsUsersForceDelete.ValueBool()
					}
				}
			}
		}

	}

	var overrideApiHostname, overrideAuthHostname string
	if !data.ServiceEndpoints.IsNull() {
		var serviceEndpointsData []pingOneProviderServiceEndpointsModel
		resp.Diagnostics.Append(data.ServiceEndpoints.ElementsAs(ctx, &serviceEndpointsData, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(serviceEndpointsData) > 0 {
			if !serviceEndpointsData[0].APIHostname.IsNull() {
				overrideApiHostname = serviceEndpointsData[0].APIHostname.ValueString()
			}

			if !serviceEndpointsData[0].AuthHostname.IsNull() {
				overrideAuthHostname = serviceEndpointsData[0].AuthHostname.ValueString()
			}
		}
	}
	// Override env vars are not handled in the client
	if overrideApiHostname == "" {
		overrideApiHostname = os.Getenv("PINGONE_API_SERVICE_HOSTNAME")
	}
	if overrideAuthHostname == "" {
		overrideAuthHostname = os.Getenv("PINGONE_AUTH_SERVICE_HOSTNAME")
	}
	if overrideApiHostname != "" {
		config = config.WithAPIDomain(overrideApiHostname)
	}
	if overrideAuthHostname != "" {
		config = config.WithCustomDomain(overrideAuthHostname)
	}

	pingOneConfig := pingone.NewConfiguration(config)

	if !data.HTTPProxy.IsNull() {
		v := data.HTTPProxy.ValueString()
		pingOneConfig.ProxyURL = &v
	}

	userAgent := framework.UserAgent("", p.version)
	if !data.AppendUserAgent.IsNull() && data.AppendUserAgent.ValueString() != "" {
		userAgent = framework.UserAgent(data.AppendUserAgent.ValueString(), p.version)
	} else if v := strings.TrimSpace(os.Getenv("PINGONE_TF_APPEND_USER_AGENT")); v != "" {
		userAgent = framework.UserAgent(v, p.version)
	}
	pingOneConfig.AppendUserAgent(userAgent)

	if globalOptions.Population.ContainsUsersForceDelete {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("global_options").AtName("population").AtListIndex(0).AtName("contains_users_force_delete"),
			"Data protection notice",
			"The provider is configured to force-delete populations if they contain users.  This may result in the loss of user data.  Please ensure this configuration is intentional and that you have a backup of any data you wish to retain.",
		)
	}

	apiClient, err := pingone.NewAPIClient(pingOneConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client failed to initialize",
			fmt.Sprintf("Failed to initialize the PingOne client: %v. Please report this issue to the provider maintainers.", err),
		)
		return
	}

	var resourceConfig framework.ResourceType
	resourceConfig.Client = apiClient
	tflog.Info(ctx, "[v6] Provider initialized client")

	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig

}

func (p *pingOneProvider) Resources(ctx context.Context) []func() resource.Resource {
	v := make([]func() resource.Resource, 0)
	v = append(v, davinci.Resources()...)
	return v
}

func (p *pingOneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	v := make([]func() datasource.DataSource, 0)
	v = append(v, davinci.DataSources()...)
	return v
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pingOneProvider{
			version: version,
		}
	}
}
