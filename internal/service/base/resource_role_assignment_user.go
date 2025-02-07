// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type RoleAssignmentUserResource serviceClientType

// Framework interfaces
var (
	_ resource.Resource              = &RoleAssignmentUserResource{}
	_ resource.ResourceWithConfigure = &RoleAssignmentUserResource{}
)

// New Object
func NewRoleAssignmentUserResource() resource.Resource {
	return &RoleAssignmentUserResource{}
}

// Metadata
func (r *RoleAssignmentUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_assignment_user"
}

// Schema.
func (r *RoleAssignmentUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	attributes := map[string]schema.Attribute{
		"id": framework.Attr_ID(),

		"environment_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the user to assign the admin role to."),
		),

		"user_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an user to assign an admin role to."),
		),

		"role_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an admin role to assign to the user."),
		),

		"read_only": schema.BoolAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A flag to show whether the admin role assignment is read only or can be changed.").Description,
			Computed:    true,
		},
	}

	utils.MergeSchemaAttributeMaps(attributes, service.RoleAssignmentScopeSchema(), true)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to administrator defined users.",

		DeprecationMessage: "Use the pingone_user_role_assignment resource going forward.  This resource will be removed in the next major release.",

		Attributes: attributes,
	}
}

func (r *RoleAssignmentUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleAssignmentUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Resource deprecated",
		"The `pingone_role_assignment_user` resource is deprecated and cannot be used to create new user admin role assignments. Use the `pingone_user_role_assignment` resource instead.")

}

func (r *RoleAssignmentUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserRoleAssignmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.UserId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneUserRoleAssignment",
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

func (r *RoleAssignmentUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *RoleAssignmentUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserRoleAssignmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	if data.ReadOnly.Equal(types.BoolValue(true)) {
		resp.Diagnostics.AddError(
			"Cannot destroy read only role assignment",
			fmt.Sprintf("Role assignment %s cannot be deleted as it is read only", data.Id.ValueString()),
		)

		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.UserRoleAssignmentsApi.DeleteUserRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.UserId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteUserRoleAssignment",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
