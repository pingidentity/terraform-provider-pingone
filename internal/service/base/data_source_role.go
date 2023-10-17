package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type RoleDataSource serviceClientType

type RoleDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &RoleDataSource{}
)

// New Object
func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// Metadata
func (r *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema
func (r *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const minAttrLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne admin role data for a tenant.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The name of the role to look up.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(minAttrLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The description of the role.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RoleDataSourceModel

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

	var role management.Role

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

	if roles, ok := entityArray.Embedded.GetRolesOk(); ok {

		found := false
		for _, roleItem := range roles {

			if string(roleItem.GetName()) == data.Name.ValueString() {
				role = roleItem
				found = true
				break
			}
		}

		if !found {
			resp.Diagnostics.AddError(
				"Cannot find role from name",
				fmt.Sprintf("The role %s cannot be found in the tenant", data.Name.String()),
			)
			return
		}

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&role)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *RoleDataSourceModel) toState(v *management.Role) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(v.GetIdOk())
	p.Name = framework.EnumOkToTF(v.GetNameOk())
	p.Description = framework.StringOkToTF(v.GetDescriptionOk())

	return diags
}
