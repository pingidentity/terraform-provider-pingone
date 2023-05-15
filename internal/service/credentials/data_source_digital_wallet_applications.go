package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type DigitalWalletApplicationsDataSource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type DigitalWalletApplicationsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &DigitalWalletApplicationsDataSource{}
)

// New Object
func NewDigitalWalletApplicationsDataSource() datasource.DataSource {
	return &DigitalWalletApplicationsDataSource{}
}

// Metadata
func (r *DigitalWalletApplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_digital_wallet_applications"
}

func (r *DigitalWalletApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a list of multiple PingOne Credentials digital wallet applications.  Filtering the list by SCIM or data filter currently is not supported.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "PingOne environment identifier (UUID) in which the credential digital wallet app exists."},
			),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaDescription{
				Description: "The list of resulting IDs of digital walliet applications that have been successfully retrieved.",
			}),
		},
	}
}

func (r *DigitalWalletApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *DigitalWalletApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DigitalWalletApplicationsDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DigitalWalletAppsApi.ReadAllDigitalWalletApps(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadAllDigitalWalletApplications",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Error if not found
	if response == nil {
		resp.Diagnostics.AddError(
			"Cannot find digital wallet applications",
			fmt.Sprintf("The list of digital wallet ids for environment %s cannot be found.", data.EnvironmentId.String()),
		)
		return
	}

	entityArray := response.(*credentials.EntityArray)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), entityArray.Embedded.GetDigitalWalletApplications())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (p *DigitalWalletApplicationsDataSourceModel) toState(environmentID string, digitalWalletApplications []credentials.DigitalWalletApplication) diag.Diagnostics {
	var diags diag.Diagnostics

	if digitalWalletApplications == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range digitalWalletApplications {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	p.Id = framework.StringToTF(environmentID)
	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
