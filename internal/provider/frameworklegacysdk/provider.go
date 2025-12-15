// Copyright © 2025 Ping Identity Corporation

package frameworklegacysdk

import (
	"context"
	"fmt"
	"os"
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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Ensure PingOneProvider satisfies various provider interfaces.
var _ provider.Provider = &pingOneProvider{}

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
	AppendUserAgent  types.String `tfsdk:"append_user_agent"`
	RegionCode       types.String `tfsdk:"region_code"`
	ServiceEndpoints types.List   `tfsdk:"service_endpoints"`
	GlobalOptions    types.List   `tfsdk:"global_options"`
	HTTPProxy        types.String `tfsdk:"http_proxy"`
	ConfigPath       types.String `tfsdk:"config_path"`
	ConfigProfile    types.String `tfsdk:"config_profile"`
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

	configPathDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Path to a PingCLI configuration file containing authentication credentials. When set, the provider will read authentication settings from the specified profile in this file. Default value can be set with the `PINGCLI_CONFIG` environment variable. Cannot be used together with `client_id`, `client_secret`, `environment_id`, or `api_access_token`. If set, `config_profile` can optionally specify which profile to use (defaults to the active profile in the config file).",
	)

	configProfileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the profile to use from the PingCLI configuration file. If not specified, uses the active profile defined in the config file. Default value can be set with the `PINGCLI_PROFILE` environment variable. Requires `config_path` to be set.",
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
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumRegionCodeEnumValues)...),
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
	tflog.Debug(ctx, "[v6] [legacy sdk] Provider configure start")
	var data pingOneProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if PingCLI config is being used
	var configPath string
	if !data.ConfigPath.IsNull() && !data.ConfigPath.IsUnknown() {
		configPath = strings.TrimSpace(data.ConfigPath.ValueString())
	} else {
		configPath = strings.TrimSpace(os.Getenv("PINGCLI_CONFIG"))
	}

	// If PingCLI config is being used, load it and initialize the client
	usingPingCLIConfig := configPath != ""

	// Initialize globalOptions that will be used regardless of auth method
	globalOptions := &client.GlobalOptions{
		Population: &client.PopulationOptions{
			ContainsUsersForceDelete: false,
		},
	}

	var config *pingone.Config

	if usingPingCLIConfig {
		tflog.Info(ctx, "[v6] [legacy sdk] Using PingCLI configuration")

		var configProfile string
		if !data.ConfigProfile.IsNull() && !data.ConfigProfile.IsUnknown() {
			configProfile = strings.TrimSpace(data.ConfigProfile.ValueString())
		} else {
			configProfile = strings.TrimSpace(os.Getenv("PINGCLI_PROFILE"))
		}

		// Load PingCLI configuration - this will determine the active profile if configProfile is empty
		profileConfig, err := pingcli.LoadProfileConfig(configPath, configProfile)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Load PingCLI Configuration",
				fmt.Sprintf("Error reading PingCLI configuration file: %s", err.Error()),
			)
			return
		}

		// If no profile was specified, get the active profile name for token loading
		if configProfile == "" {
			pingCliConfig, err := pingcli.NewConfig(configPath)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to Load PingCLI Configuration",
					fmt.Sprintf("Error reading PingCLI configuration file: %s", err.Error()),
				)
				return
			}
			configProfile, err = pingCliConfig.GetActiveProfile()
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to Get Active Profile",
					fmt.Sprintf("Error getting active profile from PingCLI configuration: %s", err.Error()),
				)
				return
			}
		}

		// Try to load stored token
		storedToken, err := pingcli.LoadStoredToken(profileConfig, configProfile)
		if err != nil || storedToken == nil || !storedToken.Valid() {
			resp.Diagnostics.AddError(
				"No Valid Stored Token",
				fmt.Sprintf("PingCLI configuration requires a valid stored token. Please run 'pingcli login' first. Error: %v", err),
			)
			return
		}

		tflog.Info(ctx, "[v6] [legacy sdk] Using stored token from PingCLI credentials")

		config = &pingone.Config{
			AccessToken:   storedToken.AccessToken,
			GlobalOptions: globalOptions,
		}

		tflog.Debug(ctx, "[v6] [legacy sdk] Config initialized", map[string]any{
			"has_access_token":  config.AccessToken != "",
			"has_client_id":     config.ClientID != "",
			"has_client_secret": config.ClientSecret != "",
		})

		if profileConfig.RegionCode != "" {
			rc := management.EnumRegionCode(profileConfig.RegionCode)
			config.RegionCode = &rc
		}
	} else {
		// Use explicit credentials from provider configuration or environment variables
		config = &pingone.Config{
			ClientID:      data.ClientID.ValueString(),
			ClientSecret:  data.ClientSecret.ValueString(),
			EnvironmentID: data.EnvironmentID.ValueString(),
			AccessToken:   data.APIAccessToken.ValueString(),
			GlobalOptions: globalOptions,
		}

		if !data.RegionCode.IsNull() && !data.RegionCode.IsUnknown() {
			rc := management.EnumRegionCode(data.RegionCode.ValueString())
			config.RegionCode = &rc
		}
	}

	if !data.HTTPProxy.IsNull() {
		v := data.HTTPProxy.ValueString()
		config.ProxyURL = &v
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

	if !data.ServiceEndpoints.IsNull() {

		var serviceEndpointsData []pingOneProviderServiceEndpointsModel
		resp.Diagnostics.Append(data.ServiceEndpoints.ElementsAs(ctx, &serviceEndpointsData, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(serviceEndpointsData) > 0 {
			if !serviceEndpointsData[0].AuthHostname.IsNull() {
				v := serviceEndpointsData[0].AuthHostname.ValueString()
				config.AuthHostnameOverride = &v
			}

			if !serviceEndpointsData[0].APIHostname.IsNull() {
				v := serviceEndpointsData[0].APIHostname.ValueString()
				config.APIHostnameOverride = &v
			}
		}

	}

	if !data.AppendUserAgent.IsNull() {
		config.UserAgentAppend = data.AppendUserAgent.ValueStringPointer()
	} else if v := strings.TrimSpace(os.Getenv("PINGONE_TF_APPEND_USER_AGENT")); v != "" {
		config.UserAgentAppend = &v
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if globalOptions.Population.ContainsUsersForceDelete {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("global_options").AtName("population").AtListIndex(0).AtName("contains_users_force_delete"),
			"Data protection notice",
			"The provider is configured to force-delete populations if they contain users.  This may result in the loss of user data.  Please ensure this configuration is intentional and that you have a backup of any data you wish to retain.",
		)
	}

	apiClient, err := config.APIClient(ctx, p.version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client failed to initialize",
			fmt.Sprintf("Failed to initialize the PingOne client: %v. Please report this issue to the provider maintainers.", err),
		)
		return
	}

	var resourceConfig legacysdk.ResourceType
	resourceConfig.Client = apiClient
	tflog.Info(ctx, "[v6] [legacy sdk] Provider initialized client")

	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig

}

func (p *pingOneProvider) Resources(ctx context.Context) []func() resource.Resource {
	v := make([]func() resource.Resource, 0)
	v = append(v, authorize.Resources()...)
	v = append(v, base.Resources()...)
	v = append(v, mfa.Resources()...)
	v = append(v, sso.Resources()...)
	v = append(v, risk.Resources()...)
	v = append(v, credentials.Resources()...)
	v = append(v, verify.Resources()...)
	return v
}

func (p *pingOneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	v := make([]func() datasource.DataSource, 0)
	v = append(v, authorize.DataSources()...)
	v = append(v, base.DataSources()...)
	v = append(v, mfa.DataSources()...)
	v = append(v, sso.DataSources()...)
	v = append(v, risk.DataSources()...)
	v = append(v, credentials.DataSources()...)
	v = append(v, verify.DataSources()...)
	return v
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pingOneProvider{
			version: version,
		}
	}
}
