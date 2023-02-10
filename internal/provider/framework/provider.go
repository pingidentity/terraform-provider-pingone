package framework

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
)

// Ensure PingOneProvider satisfies various provider interfaces.
var _ provider.Provider = &PingOneProvider{}

// PingOneProvider defines the provider implementation.
type PingOneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	client  *pingone.Client
}

// PingOneProviderModel describes the provider data model.
type PingOneProviderModel struct {
	ClientID                             types.String `tfsdk:"client_id"`
	ClientSecret                         types.String `tfsdk:"client_secret"`
	EnvironmentID                        types.String `tfsdk:"environment_id"`
	APIAccessToken                       types.String `tfsdk:"api_access_token"`
	Region                               types.String `tfsdk:"region"`
	ForceDeleteProductionEnvironmentType types.Bool   `tfsdk:"force_delete_production_type"`
}

func (p *PingOneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingone"
	resp.Version = p.version
}

func (p *PingOneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_id")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_secret")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("api_access_token")), // path.MatchRoot("other_attr")),
				},
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_id")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_secret")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("api_access_token")),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_id")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("client_secret")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("api_access_token")),
				},
			},
			"api_access_token": schema.StringAttribute{
				MarkdownDescription: "The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("client_id")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("client_secret")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("client_id")), // path.MatchRoot("other_attr")),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("client_id")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("client_secret")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("client_id")), // path.MatchRoot("other_attr")),
					stringvalidator.OneOf(model.RegionsAvailableList()...),
				},
			},
			"force_delete_production_type": schema.BoolAttribute{
				MarkdownDescription: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *PingOneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PingOneProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the defaults
	if data.ClientID.IsNull() {
		data.ClientID = basetypes.NewStringValue(os.Getenv("PINGONE_CLIENT_ID"))
	}

	if data.ClientSecret.IsNull() {
		data.ClientSecret = basetypes.NewStringValue(os.Getenv("PINGONE_CLIENT_SECRET"))
	}

	if data.EnvironmentID.IsNull() {
		data.EnvironmentID = basetypes.NewStringValue(os.Getenv("PINGONE_ENVIRONMENT_ID"))
	}

	if data.APIAccessToken.IsNull() {
		data.APIAccessToken = basetypes.NewStringValue(os.Getenv("PINGONE_API_ACCESS_TOKEN"))
	}

	if data.Region.IsNull() {
		data.Region = basetypes.NewStringValue(os.Getenv("PINGONE_REGION"))
	}

	if data.ForceDeleteProductionEnvironmentType.IsNull() {
		v, err := strconv.ParseBool(os.Getenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE"))
		if err != nil {
			v = false
		}
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

	client, err := config.APIClient(ctx)
	if err != nil {
		return
	}

	p.client = client
}

func (p *PingOneProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//base.NotificationSettingsResource,
		// base.NotificationPolicyResource,
		// base.PhoneDeliverySettingsResource,
		base.TrustedEmailAddressResource,
	} // define resources here
}

func (p *PingOneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{} // define resources here
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PingOneProvider{
			version: version,
		}
	}
}
