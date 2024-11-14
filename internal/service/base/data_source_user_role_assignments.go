package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
type UserRoleAssignmentsDataSource serviceClientType

type UserRoleAssignmentsDataSourceModel struct {
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	UserId          pingonetypes.ResourceIDValue `tfsdk:"user_id"`
	RoleAssignments types.Set                    `tfsdk:"role_assignments"`
}

var (
	roleAssignmentsUserRoleAssignmentTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
		"scope": types.ObjectType{
			AttrTypes: roleAssignmentsUserRoleAssignmentScopeTFObjectTypes,
		},
		"role_id":   pingonetypes.ResourceIDType{},
		"read_only": types.BoolType,
	}

	roleAssignmentsUserRoleAssignmentScopeTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"type": types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &UserRoleAssignmentsDataSource{}
)

// New Object
func NewUserRoleAssignmentsDataSource() datasource.DataSource {
	return &UserRoleAssignmentsDataSource{}
}

// Metadata
func (r *UserRoleAssignmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_role_assignments"
}

// Schema
func (r *UserRoleAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	roleAssignmentsScopesTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of the scope.",
	).AllowedValuesEnum(management.AllowedEnumRoleAssignmentScopeTypeEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the role assignments that a user has been assigned.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the admin user to retrieve role assignments for."),
			),

			"user_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the admin user to retrieve role assignments for."),
			),

			"role_assignments": schema.SetNestedAttribute{
				Description: "A set of role assignments that the user has been assigned.",
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the role assignment relationship.").Description,
							Computed:    true,

							CustomType: pingonetypes.ResourceIDType{},
						},

						"scope": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that describes the scope of the role assignment.").Description,
							Computed:    true,

							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the environment, population or organization that the role is scoped to.").Description,
									Computed:    true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"type": schema.StringAttribute{
									Description:         roleAssignmentsScopesTypeDescription.Description,
									MarkdownDescription: roleAssignmentsScopesTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},

						"role_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the admin role that has been assigned to the user.").Description,
							Computed:    true,

							CustomType: pingonetypes.ResourceIDType{},
						},

						"read_only": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the admin role assignment is read only or can be changed.").Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *UserRoleAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *UserRoleAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UserRoleAssignmentsDataSourceModel

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

	// Run the API call
	var roleAssignments []management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.UserRoleAssignmentsApi.ReadUserRoleAssignments(ctx, data.EnvironmentId.ValueString(), data.UserId.ValueString()).Execute()

			roleAssignments := make([]management.RoleAssignment, 0)

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.RoleAssignments != nil {
					roleAssignments = append(roleAssignments, pageCursor.EntityArray.Embedded.GetRoleAssignments()...)
				}

			}

			return roleAssignments, initialHttpResponse, nil
		},
		"ReadUserRoleAssignments",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&roleAssignments,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(roleAssignments)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *UserRoleAssignmentsDataSourceModel) toState(apiObject []management.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	if p.Id.IsNull() {
		p.Id = framework.PingOneResourceIDToTF(uuid.New().String())
	}

	var d diag.Diagnostics
	p.RoleAssignments, d = toStateRoleAssignments(apiObject)
	diags = append(diags, d...)

	return diags
}

func toStateRoleAssignments(apiObject []management.RoleAssignment) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: roleAssignmentsUserRoleAssignmentTFObjectTypes}

	if apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := make([]attr.Value, 0)
	for _, item := range apiObject {

		scopeObj, d := toStateRoleAssignmentsScope(item.GetScopeOk())
		diags = append(diags, d...)

		objMap := map[string]attr.Value{
			"id":        framework.PingOneResourceIDOkToTF(item.GetIdOk()),
			"scope":     scopeObj,
			"role_id":   framework.PingOneResourceIDOkToTF(item.Role.GetIdOk()),
			"read_only": framework.BoolOkToTF(item.GetReadOnlyOk()),
		}

		flattenedObj, d := types.ObjectValue(roleAssignmentsUserRoleAssignmentTFObjectTypes, objMap)
		diags = append(diags, d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateRoleAssignmentsScope(apiObject *management.RoleAssignmentScope, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(roleAssignmentsUserRoleAssignmentScopeTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"id":   framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	returnVar, d := types.ObjectValue(roleAssignmentsUserRoleAssignmentScopeTFObjectTypes, objMap)
	diags = append(diags, d...)

	return returnVar, diags
}
