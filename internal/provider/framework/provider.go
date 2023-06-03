package framework

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	client  *pingone.Client
}

// pingOneProviderModel describes the provider data model.
type pingOneProviderModel struct {
	ClientID                             types.String `tfsdk:"client_id"`
	ClientSecret                         types.String `tfsdk:"client_secret"`
	EnvironmentID                        types.String `tfsdk:"environment_id"`
	APIAccessToken                       types.String `tfsdk:"api_access_token"`
	Region                               types.String `tfsdk:"region"`
	ForceDeleteProductionEnvironmentType types.Bool   `tfsdk:"force_delete_production_type"`
}

func (p *pingOneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingone"
	resp.Version = p.version
}

func (p *pingOneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.",
				Optional:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.",
				Optional:            true,
			},
			"api_access_token": schema.StringAttribute{
				MarkdownDescription: "The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.",
				Optional:            true,
			},
			"force_delete_production_type": schema.BoolAttribute{
				MarkdownDescription: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
				Optional:            true,
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
	if data.ClientID.IsNull() {
		data.ClientID = basetypes.NewStringValue(os.Getenv("PINGONE_CLIENT_ID"))
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "client_id"), map[string]interface{}{
			"env_var":       "PINGONE_CLIENT_ID",
			"env_var_value": data.ClientID.String(),
		})
	}

	if data.ClientSecret.IsNull() {
		data.ClientSecret = basetypes.NewStringValue(os.Getenv("PINGONE_CLIENT_SECRET"))
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "client_secret"), map[string]interface{}{
			"env_var": "PINGONE_CLIENT_SECRET",
			"env_var_value": func() string {
				if len(data.ClientSecret.ValueString()) > 0 {
					return "***"
				}
				return ""

			}(),
		})
	}

	if data.EnvironmentID.IsNull() {
		data.EnvironmentID = basetypes.NewStringValue(os.Getenv("PINGONE_ENVIRONMENT_ID"))
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "environment_id"), map[string]interface{}{
			"env_var":       "PINGONE_ENVIRONMENT_ID",
			"env_var_value": data.EnvironmentID.String(),
		})
	}

	if data.APIAccessToken.IsNull() {
		data.APIAccessToken = basetypes.NewStringValue(os.Getenv("PINGONE_API_ACCESS_TOKEN"))
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "api_access_token"), map[string]interface{}{
			"env_var": "PINGONE_API_ACCESS_TOKEN",
			"env_var_value": func() string {
				if len(data.APIAccessToken.ValueString()) > 0 {
					return "***"
				}
				return ""

			}(),
		})
	}

	if data.Region.IsNull() {
		data.Region = basetypes.NewStringValue(os.Getenv("PINGONE_REGION"))
		tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "region"), map[string]interface{}{
			"env_var":       "PINGONE_REGION",
			"env_var_value": data.Region.String(),
		})
	}

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

	// Example client configuration for data sources and resources
	config := &pingone.Config{
		ClientID:      data.ClientID.ValueString(),
		ClientSecret:  data.ClientSecret.ValueString(),
		EnvironmentID: data.EnvironmentID.ValueString(),
		AccessToken:   data.APIAccessToken.ValueString(),
		Region:        data.Region.ValueString(),
		ForceDelete:   data.ForceDeleteProductionEnvironmentType.ValueBool(),
	}

	err := config.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Parameter validation error",
			fmt.Sprintf("%s", err))
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
