// Code generated by ping-terraform-plugin-framework-generator

package sso

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
)

var (
	_ datasource.DataSource = &customRoleDataSource{}
)

func NewCustomRoleDataSource() datasource.DataSource {
	return &customRoleDataSource{}
}

type customRoleDataSource serviceClientType

func (r *customRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_role"
}

func (r *customRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type customRoleDataSourceModel struct {
	ApplicableTo    types.Set                    `tfsdk:"applicable_to"`
	CanAssign       types.Set                    `tfsdk:"can_assign"`
	CanBeAssignedBy types.Set                    `tfsdk:"can_be_assigned_by"`
	Description     types.String                 `tfsdk:"description"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name            types.String                 `tfsdk:"name"`
	Permissions     types.Set                    `tfsdk:"permissions"`
	RoleId          pingonetypes.ResourceIDValue `tfsdk:"role_id"`
	Type            types.String                 `tfsdk:"type"`
}

func (r *customRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve a custom administrator role in an environment, by ID or name.",
		Attributes: map[string]schema.Attribute{
			"applicable_to": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The scope types to which the role can be applied. Options are \"ENVIRONMENT\", \"ORGANIZATION\", \"POPULATION\".",
				MarkdownDescription: "The scope types to which the role can be applied. Options are `ENVIRONMENT`, `ORGANIZATION`, `POPULATION`.",
			},
			"can_assign": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of a role that can be assigned by an actor assigned the current custom role.",
							CustomType:  pingonetypes.ResourceIDType{},
						},
					},
				},
				Computed:            true,
				Description:         "A relationship that specifies if an actor is assigned the current custom role for a jurisdiction, then the actor can assign any of this set of roles to another actor for the same jurisdiction or sub-jurisdiction. This capability is derived from the \"can_be_assigned_by\" property.",
				MarkdownDescription: "A relationship that specifies if an actor is assigned the current custom role for a jurisdiction, then the actor can assign any of this set of roles to another actor for the same jurisdiction or sub-jurisdiction. This capability is derived from the `can_be_assigned_by` property.",
			},
			"can_be_assigned_by": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the role that can assign the current custom role.",
							CustomType:  pingonetypes.ResourceIDType{},
						},
					},
				},
				Computed:    true,
				Description: "A relationship that determines whether a user assigned to one of this set of roles for a jurisdiction can assign the current custom role to another user for the same jurisdiction or sub-jurisdiction.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the role.",
			},
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the custom role."),
			),
			"id": framework.Attr_ID(),
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "A string that specifies the name of the custom role to retrieve configuration for. Exactly one of \"name\" or \"role_id\" must be defined.",
				MarkdownDescription: "A string that specifies the name of the custom role to retrieve configuration for. Exactly one of `name` or `role_id` must be defined.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("role_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"permissions": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the permission assigned to this role.",
							Computed:    true,
						},
					},
				},
				Computed:    true,
				Description: "The set of permissions assigned to the role. For possible values, see the [list of available permissions](https://apidocs.pingidentity.com/pingone/platform/v1/api/#pingone-permissions-by-identifier).",
			},
			"role_id": schema.StringAttribute{
				CustomType:          pingonetypes.ResourceIDType{},
				Optional:            true,
				Description:         "A string that specifies the ID of the role to retrieve configuration for.  Must be a valid PingOne resource ID. Exactly one of \"name\" or \"role_id\" must be defined.",
				MarkdownDescription: "A string that specifies the ID of the role to retrieve configuration for.  Must be a valid PingOne resource ID. Exactly one of `name` or `role_id` must be defined.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "A value that indicates whether the role is a built-in role or a custom role. Options are \"PLATFORM\" and \"CUSTOM\". This will always be \"CUSTOM\" for custom roles.",
				MarkdownDescription: "A value that indicates whether the role is a built-in role or a custom role. Options are `PLATFORM` and `CUSTOM`. This will always be `CUSTOM` for custom roles.",
			},
		},
	}
}

func (state *customRoleDataSourceModel) readClientResponse(response *management.CustomAdminRole) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// applicable_to
	state.ApplicableTo, diags = types.SetValueFrom(context.Background(), types.StringType, response.ApplicableTo)
	respDiags.Append(diags...)
	// can_assign
	canAssignAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	canAssignElementType := types.ObjectType{AttrTypes: canAssignAttrTypes}
	var canAssignValues []attr.Value
	for _, canAssignResponseValue := range response.CanAssign {
		canAssignValue, diags := types.ObjectValue(canAssignAttrTypes, map[string]attr.Value{
			"id": types.StringValue(canAssignResponseValue.Id),
		})
		respDiags.Append(diags...)
		canAssignValues = append(canAssignValues, canAssignValue)
	}
	canAssignValue, diags := types.SetValue(canAssignElementType, canAssignValues)
	respDiags.Append(diags...)

	state.CanAssign = canAssignValue
	// can_be_assigned_by
	canBeAssignedByAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	canBeAssignedByElementType := types.ObjectType{AttrTypes: canBeAssignedByAttrTypes}
	var canBeAssignedByValues []attr.Value
	for _, canBeAssignedByResponseValue := range response.CanBeAssignedBy {
		canBeAssignedByValue, diags := types.ObjectValue(canBeAssignedByAttrTypes, map[string]attr.Value{
			"id": types.StringValue(canBeAssignedByResponseValue.Id),
		})
		respDiags.Append(diags...)
		canBeAssignedByValues = append(canBeAssignedByValues, canBeAssignedByValue)
	}
	canBeAssignedByValue, diags := types.SetValue(canBeAssignedByElementType, canBeAssignedByValues)
	respDiags.Append(diags...)

	state.CanBeAssignedBy = canBeAssignedByValue
	// description
	state.Description = types.StringPointerValue(response.Description)
	// id
	idValue := framework.PingOneResourceIDToTF(response.GetId())

	state.Id = idValue
	// name
	state.Name = types.StringValue(response.Name)
	// permissions
	permissionsAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	permissionsElementType := types.ObjectType{AttrTypes: permissionsAttrTypes}
	var permissionsValues []attr.Value
	for _, permissionsResponseValue := range response.Permissions {
		permissionsValue, diags := types.ObjectValue(permissionsAttrTypes, map[string]attr.Value{
			"id": types.StringValue(permissionsResponseValue.Id),
		})
		respDiags.Append(diags...)
		permissionsValues = append(permissionsValues, permissionsValue)
	}
	permissionsValue, diags := types.SetValue(permissionsElementType, permissionsValues)
	respDiags.Append(diags...)

	state.Permissions = permissionsValue
	// role_id
	state.RoleId = framework.PingOneResourceIDToTF(response.GetId())
	// type
	if response.Type != nil {
		state.Type = types.StringValue(string(*response.Type))
	} else {
		state.Type = types.StringNull()
	}
	return respDiags
}

func (r *customRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data customRoleDataSourceModel

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

	// Read API call logic
	var response *management.CustomAdminRole

	// Custom Role API does not support SCIM filtering currently in the SDK
	if !data.RoleId.IsNull() {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.CustomAdminRolesApi.ReadOneCustomAdminRole(ctx, data.EnvironmentId.ValueString(), data.RoleId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneCustomAdminRole",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if !data.Name.IsNull() {
		// Get all custom roles and find the one with the expected name
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.CustomAdminRolesApi.ReadAllCustomAdminRoles(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if roles, ok := pageCursor.EntityArray.Embedded.GetRolesOk(); ok {

						for _, roleObj := range roles {
							roleInstance := roleObj.CustomAdminRole
							if roleInstance == nil {
								continue
							}
							roleName := roleInstance.Name

							if roleName == data.Name.ValueString() {
								return roleInstance, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllCustomAdminRoles",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response == nil {
			resp.Diagnostics.AddError(
				"Cannot find the custom role from name or the custom role is not the correct type",
				fmt.Sprintf("The custom role name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested custom role. role_id or name must be set.",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(response)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
