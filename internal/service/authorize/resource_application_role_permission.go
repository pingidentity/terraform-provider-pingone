package authorize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationRolePermissionResource serviceClientType

type ApplicationRolePermissionResourceModel struct {
	EnvironmentId                   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationRoleId               pingonetypes.ResourceIDValue `tfsdk:"application_role_id"`
	ApplicationResourcePermissionId pingonetypes.ResourceIDValue `tfsdk:"application_resource_permission_id"`
	Action                          types.String                 `tfsdk:"action"`
	Description                     types.String                 `tfsdk:"description"`
	Key                             types.String                 `tfsdk:"key"`
	Resource                        types.Object                 `tfsdk:"resource"`
}

type ApplicationRolePermissionPermissionResourceResourceModel struct {
	Id   pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name types.String                 `tfsdk:"name"`
}

var (
	applicationRolePermissionPermissionResourceTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"name": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationRolePermissionResource{}
	_ resource.ResourceWithConfigure   = &ApplicationRolePermissionResource{}
	_ resource.ResourceWithImportState = &ApplicationRolePermissionResource{}
)

// New Object
func NewApplicationRolePermissionResource() resource.Resource {
	return &ApplicationRolePermissionResource{}
}

// Metadata
func (r *ApplicationRolePermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_application_role_permission"
}

func (r *ApplicationRolePermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage application role permissions in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the application role permission in."),
			),

			"application_role_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application role to configure the application role permission for."),
			),

			"application_resource_permission_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application resource permission to assign to the application role."),
			),

			"action": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the action associated with this permission.").Description,
				Computed:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the description associated with this permission.").Description,
				Computed:    true,
			},

			"key": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the composite key associated with this permission assignment.").Description,
				Computed:    true,
			},

			"resource": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes the assigned application resource.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the ID of the application resource associated with this permission.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the name of the application resource associated with this permission.").Description,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *ApplicationRolePermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationRolePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationRolePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationRolePermission := plan.expand()

	// Run the API call
	var response *authorize.ApplicationRolePermission
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.ApplicationRolePermissionsApi.CreateApplicationRolePermission(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationRoleId.ValueString()).ApplicationRolePermission(*applicationRolePermission).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplicationRolePermission",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationRolePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationRolePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.ApplicationRolePermission
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.AuthorizeAPIClient.ApplicationRolePermissionsApi.ReadApplicationRolePermissions(ctx, data.EnvironmentId.ValueString(), data.ApplicationRoleId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Permissions != nil {
					for _, permission := range pageCursor.EntityArray.Embedded.Permissions {
						if v := permission.ApplicationRolePermission; v != nil && v.GetId() == data.ApplicationResourcePermissionId.ValueString() {
							return v, pageCursor.HTTPResponse, nil
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadApplicationRolePermissions",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationRolePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ApplicationRolePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationRolePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.AuthorizeAPIClient.ApplicationRolePermissionsApi.DeleteApplicationRolePermission(ctx, data.EnvironmentId.ValueString(), data.ApplicationRoleId.ValueString(), data.ApplicationResourcePermissionId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteApplicationRolePermission",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationRolePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_role_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_resource_permission_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *ApplicationRolePermissionResourceModel) expand() *authorize.ApplicationRolePermission {

	// Main object
	data := authorize.NewApplicationRolePermission(
		p.ApplicationResourcePermissionId.ValueString(),
	)

	return data
}

func (p *ApplicationRolePermissionResourceModel) toState(apiObject *authorize.ApplicationRolePermission) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.ApplicationResourcePermissionId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Action = framework.StringOkToTF(apiObject.GetActionOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Key = framework.StringOkToTF(apiObject.GetKeyOk())

	p.Resource, d = toStateApplicationRolePermissionPermissionResourceOkToTF(apiObject.GetResourceOk())
	diags.Append(d...)

	return diags
}

func toStateApplicationRolePermissionPermissionResourceOkToTF(apiObject *authorize.ApplicationRolePermissionResource, ok bool) (types.Object, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationRolePermissionPermissionResourceTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"id":   framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	}

	returnVar, d := types.ObjectValue(applicationRolePermissionPermissionResourceTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}
