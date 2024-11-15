package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type UserApplicationRoleAssignmentResource serviceClientType

type UserApplicationRoleAssignmentResourceModel struct {
	EnvironmentId     pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	UserId            pingonetypes.ResourceIDValue `tfsdk:"user_id"`
	ApplicationRoleId pingonetypes.ResourceIDValue `tfsdk:"application_role_id"`
	Name              types.String                 `tfsdk:"name"`
}

// Framework interfaces
var (
	_ resource.Resource                = &UserApplicationRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &UserApplicationRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &UserApplicationRoleAssignmentResource{}
)

// New Object
func NewUserApplicationRoleAssignmentResource() resource.Resource {
	return &UserApplicationRoleAssignmentResource{}
}

// Metadata
func (r *UserApplicationRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_application_role_assignment"
}

// Schema.
func (r *UserApplicationRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an application role assignment to a user in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the user to assign the application role to."),
			),

			"user_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the user to assign the application role to."),
			),

			"application_role_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application role to assign to the user."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the name of the application role assigned to the user.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *UserApplicationRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserApplicationRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state UserApplicationRoleAssignmentResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	roleAssignment := plan.expand()

	// Run the API call
	var response *management.UserApplicationRoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UserApplicationRoleAssignmentsApi.CreateUserApplicationRoleAssignment(ctx, plan.EnvironmentId.ValueString(), plan.UserId.ValueString()).UserApplicationRoleAssignment(*roleAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateUserApplicationRoleAssignment",
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

func (r *UserApplicationRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserApplicationRoleAssignmentResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var roleAssignmentObj *management.UserApplicationRoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.UserApplicationRoleAssignmentsApi.ReadUserApplicationRoleAssignments(ctx, data.EnvironmentId.ValueString(), data.UserId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				for _, roleAssignment := range pageCursor.EntityArray.Embedded.GetRoles() {
					if v, ok := roleAssignment.UserApplicationRoleAssignment.GetIdOk(); ok && v != nil && *v == data.ApplicationRoleId.ValueString() {
						return roleAssignment.UserApplicationRoleAssignment, pageCursor.HTTPResponse, nil
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadUserApplicationRoleAssignments",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&roleAssignmentObj,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if roleAssignmentObj == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(roleAssignmentObj)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserApplicationRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *UserApplicationRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserApplicationRoleAssignmentResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.UserApplicationRoleAssignmentsApi.DeleteUserApplicationRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.UserId.ValueString(), data.ApplicationRoleId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteUserApplicationRoleAssignment",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserApplicationRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "user_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_role_id",
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

func (p *UserApplicationRoleAssignmentResourceModel) expand() *management.UserApplicationRoleAssignment {

	data := management.NewUserApplicationRoleAssignment(p.ApplicationRoleId.ValueString())

	return data
}

func (p *UserApplicationRoleAssignmentResourceModel) toState(apiObject *management.UserApplicationRoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.ApplicationRoleId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}
