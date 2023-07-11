package framework

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
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
	ForceDeleteProductionEnvironmentType types.Bool   `tfsdk:"force_delete_production_type"`
}

type pingOneProviderServiceEndpointsModel struct {
	AuthHostname          types.String `tfsdk:"auth_hostname"`
	APIHostname           types.String `tfsdk:"api_hostname"`
	AgreementMgmtHostname types.String `tfsdk:"agreement_management_hostname"`
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

	serviceEndpointsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block containing configuration items to override the service API endpoints of PingOne.",
	)

	serviceEndpointsAuthHostnameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.",
	)

	serviceEndpointsApiHostnameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.",
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
			},
		},

		Blocks: map[string]schema.Block{
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
					listvalidator.SizeAtLeast(1),
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

	if data.ForceDeleteProductionEnvironmentType.IsNull() {
		v, err := strconv.ParseBool(os.Getenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE"))
		if err != nil {
			v = false
		}
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "force_delete_production_type"), map[string]interface{}{
			"env_var":       "PINGONE_FORCE_DELETE_PRODUCTION_TYPE",
			"env_var_value": v,
		})
		data.ForceDeleteProductionEnvironmentType = basetypes.NewBoolValue(v)
	}

	config := &pingone.Config{
		ClientID:      data.ClientID.ValueString(),
		ClientSecret:  data.ClientSecret.ValueString(),
		EnvironmentID: data.EnvironmentID.ValueString(),
		AccessToken:   data.APIAccessToken.ValueString(),
		Region:        data.Region.ValueString(),
		ForceDelete:   data.ForceDeleteProductionEnvironmentType.ValueBool(),
	}

	if !data.ServiceEndpoints.IsNull() {

		var serviceEndpointsData pingOneProviderServiceEndpointsModel
		resp.Diagnostics.Append(data.ServiceEndpoints.ElementsAs(ctx, &serviceEndpointsData, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !serviceEndpointsData.AuthHostname.IsNull() {
			v := serviceEndpointsData.AuthHostname.ValueString()
			config.AuthHostnameOverride = &v
		}

		if !serviceEndpointsData.APIHostname.IsNull() {
			v := serviceEndpointsData.APIHostname.ValueString()
			config.APIHostnameOverride = &v
		}

	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := config.APIClient(ctx)
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
