package framework

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/agreementmanagement"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/verify"
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
	ClientID                             types.String `tfsdk:"client_id"`
	ClientSecret                         types.String `tfsdk:"client_secret"`
	EnvironmentID                        types.String `tfsdk:"environment_id"`
	APIAccessToken                       types.String `tfsdk:"api_access_token"`
	Region                               types.String `tfsdk:"region"`
	ServiceEndpoints                     types.List   `tfsdk:"service_endpoints"`
	GlobalOptions                        types.List   `tfsdk:"global_options"`
	ForceDeleteProductionEnvironmentType types.Bool   `tfsdk:"force_delete_production_type"`
	HTTPProxy                            types.String `tfsdk:"http_proxy"`
}

type pingOneProviderGlobalOptionsModel struct {
	Environment types.List `tfsdk:"environment"`
	Population  types.List `tfsdk:"population"`
}

type pingOneProviderGlobalOptionsEnvironmentModel struct {
	ProductionTypeForceDelete types.Bool `tfsdk:"production_type_force_delete"`
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

	regionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.",
	)

	forceDeleteProductionTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
	)

	globalOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing configuration items to override API behaviours in PingOne.",
	)

	globalOptionsEnvironmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing global configuration items to override environment resource settings in PingOne.",
	)

	globalOptionsEnvironmentProductionTypeForceDeleteDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
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

			"region": schema.StringAttribute{
				Description:         regionDescription.Description,
				MarkdownDescription: regionDescription.MarkdownDescription,
				Optional:            true,
			},

			"force_delete_production_type": schema.BoolAttribute{
				Description:         forceDeleteProductionTypeDescription.Description,
				MarkdownDescription: forceDeleteProductionTypeDescription.MarkdownDescription,
				Optional:            true,
				DeprecationMessage:  "This parameter is deprecated and will be removed in the next major release. Use the `global_options.environment.production_type_force_delete` block going forward.",
			},

			"http_proxy": schema.StringAttribute{
				Description:         httpProxyDescription.Description,
				MarkdownDescription: httpProxyDescription.MarkdownDescription,
				Optional:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"global_options": schema.ListNestedBlock{
				Description:         globalOptionsDescription.Description,
				MarkdownDescription: globalOptionsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Blocks: map[string]schema.Block{
						"environment": schema.ListNestedBlock{
							Description:         globalOptionsEnvironmentDescription.Description,
							MarkdownDescription: globalOptionsEnvironmentDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"production_type_force_delete": schema.BoolAttribute{
										Description:         globalOptionsEnvironmentProductionTypeForceDeleteDescription.Description,
										MarkdownDescription: globalOptionsEnvironmentProductionTypeForceDeleteDescription.MarkdownDescription,
										Optional:            true,
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},

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

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the defaults
	tflog.Info(ctx, "[v6] Provider setting defaults..")
	debugLogMessage := "[v6] Provider parameter %s missing, defaulting to environment variable"

	globalOptions := &client.GlobalOptions{
		Environment: &client.EnvironmentOptions{
			ProductionTypeForceDelete: false,
		},
		Population: &client.PopulationOptions{
			ContainsUsersForceDelete: false,
		},
	}

	if v, err := strconv.ParseBool(os.Getenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE")); err == nil && v {
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "force_delete_production_type"), map[string]interface{}{
			"env_var":       "PINGONE_FORCE_DELETE_PRODUCTION_TYPE",
			"env_var_value": v,
		})
		globalOptions.Environment.ProductionTypeForceDelete = v
	}

	config := &pingone.Config{
		ClientID:      data.ClientID.ValueString(),
		ClientSecret:  data.ClientSecret.ValueString(),
		EnvironmentID: data.EnvironmentID.ValueString(),
		AccessToken:   data.APIAccessToken.ValueString(),
		Region:        data.Region.ValueString(),
		GlobalOptions: globalOptions,
	}

	if !data.HTTPProxy.IsNull() {
		v := data.HTTPProxy.ValueString()
		config.ProxyURL = &v
	}

	deprecatedForceDeleteSet := false
	if !data.ForceDeleteProductionEnvironmentType.IsNull() {
		globalOptions.Environment.ProductionTypeForceDelete = data.ForceDeleteProductionEnvironmentType.ValueBool()
		deprecatedForceDeleteSet = true
	}

	if !data.GlobalOptions.IsNull() {

		var globalOptionsData []pingOneProviderGlobalOptionsModel
		resp.Diagnostics.Append(data.GlobalOptions.ElementsAs(ctx, &globalOptionsData, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(globalOptionsData) > 0 {
			if !globalOptionsData[0].Environment.IsNull() {

				var globalOptionsEnvironmentData []pingOneProviderGlobalOptionsEnvironmentModel
				resp.Diagnostics.Append(globalOptionsData[0].Environment.ElementsAs(ctx, &globalOptionsEnvironmentData, false)...)
				if resp.Diagnostics.HasError() {
					return
				}

				if len(globalOptionsEnvironmentData) > 0 {
					if !globalOptionsEnvironmentData[0].ProductionTypeForceDelete.IsNull() {
						if deprecatedForceDeleteSet {
							resp.Diagnostics.AddAttributeError(
								path.Root("force_delete_production_type"),
								fmt.Sprintf("Invalid provider configuration"),
								fmt.Sprintf("Cannot set both `force_delete_production_type` and `global_options.environment.production_type_force_delete` in the PingOne provider configuration.  Please unset `force_delete_production_type` and use `global_options.environment.production_type_force_delete` going forward."),
							)
							return
						}

						globalOptions.Environment.ProductionTypeForceDelete = globalOptionsEnvironmentData[0].ProductionTypeForceDelete.ValueBool()
					}
				}
			}

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

	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := config.APIClient(ctx, p.version)
	if err != nil {
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
	v = append(v, agreementmanagement.Resources()...)
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
	v = append(v, agreementmanagement.DataSources()...)
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
