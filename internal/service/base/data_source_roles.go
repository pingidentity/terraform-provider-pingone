package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type RolesDataSource serviceClientType

type RolesDataSourceModel struct {
	Id  pingonetypes.ResourceIDValue `tfsdk:"id"`
	Ids types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &RolesDataSource{}
)

// New Object
func NewRolesDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

// Metadata
func (r *RolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

// Schema
func (r *RolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a list of role IDs in the active PingOne tenant.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of roles that have been successfully retrieved.",
			)),
		},
	}
}

func (r *RolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}

}

func (r *RolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RolesDataSourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var entityArray *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.RolesApi.ReadAllRoles(ctx).Execute()
		},
		"ReadAllRoles",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(entityArray.Embedded.GetRoles())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *RolesDataSourceModel) toState(v []management.Role) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range v {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	if p.Id.IsNull() {
		p.Id = framework.PingOneResourceIDToTF(uuid.New().String())
	}

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
