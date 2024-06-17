package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type RoleDataSource serviceClientType

type RoleDataSourceModel struct {
	Id           pingonetypes.ResourceIDValue `tfsdk:"id"`
	RoleId       pingonetypes.ResourceIDValue `tfsdk:"role_id"`
	Name         types.String                 `tfsdk:"name"`
	Description  types.String                 `tfsdk:"description"`
	ApplicableTo types.Set                    `tfsdk:"applicable_to"`
	Permissions  types.Set                    `tfsdk:"permissions"`
}

var (
	rolePermissionTFObjectTypes = map[string]attr.Type{
		"id":          types.StringType,
		"classifier":  types.StringType,
		"description": types.StringType,
	}
)

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

	roleIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the role to retrieve.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf([]string{"name", "role_id"})

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the role to look up.",
	).AllowedValuesEnum(management.AllowedEnumRoleNameEnumValues).ExactlyOneOf([]string{"name", "role_id"})

	applicableToDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the applicable scopes that the role can be assigned to.",
	).AllowedValuesEnum(management.AllowedEnumRoleAssignmentScopeTypeEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne admin role data for a tenant.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"role_id": schema.StringAttribute{
				Description:         roleIdDescription.Description,
				MarkdownDescription: roleIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				MarkdownDescription: nameDescription.MarkdownDescription,
				Description:         nameDescription.Description,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("role_id")),
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumRoleNameEnumValues)...),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The description of the role.").Description,
				Computed:    true,
			},

			"applicable_to": schema.SetAttribute{
				Description:         applicableToDescription.Description,
				MarkdownDescription: applicableToDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,
			},

			"permissions": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of strings that represent permissions that have been assigned to the role.").Description,
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the permission.").Description,
							Computed:    true,
						},

						"classifier": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the resource for which the permission is applicable.").Description,
							Computed:    true,
						},

						"description": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the permission and what the permission enables.").Description,
							Computed:    true,
						},
					},
				},
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

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	var role *management.Role

	if !data.Name.IsNull() {

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

				if string(roleItem.Role.GetName()) == data.Name.ValueString() {
					role = roleItem.Role
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
	} else if !data.RoleId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ManagementAPIClient.RolesApi.ReadOneRole(ctx, data.RoleId.ValueString()).Execute()
			},
			"ReadOneRole",
			framework.DefaultCustomError,
			retryEnvironmentDefault,
			&role,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested role. role_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(role)...)
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

	p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	p.Name = framework.EnumOkToTF(v.GetNameOk())
	p.Description = framework.StringOkToTF(v.GetDescriptionOk())
	p.ApplicableTo = framework.EnumSetOkToTF(v.GetApplicableToOk())

	permissions, d := toStateRolePermissions(v.GetPermissions())
	diags.Append(d...)
	p.Permissions = permissions

	return diags
}

func toStateRolePermissions(rolePermission []management.RolePermissionsInner) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: rolePermissionTFObjectTypes}

	if len(rolePermission) == 0 {
		return types.SetValueMust(tfObjType, []attr.Value{}), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range rolePermission {

		service := map[string]attr.Value{
			"id":          framework.StringOkToTF(v.GetIdOk()),
			"classifier":  framework.StringOkToTF(v.GetClassifierOk()),
			"description": framework.StringOkToTF(v.GetDescriptionOk()),
		}

		flattenedObj, d := types.ObjectValue(rolePermissionTFObjectTypes, service)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}
