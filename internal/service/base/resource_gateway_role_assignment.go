package base

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
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type GatewayRoleAssignmentResource serviceClientType

type GatewayRoleAssignmentResourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	GatewayId           pingonetypes.ResourceIDValue `tfsdk:"gateway_id"`
	RoleId              pingonetypes.ResourceIDValue `tfsdk:"role_id"`
	ScopeApplicationId  pingonetypes.ResourceIDValue `tfsdk:"scope_application_id"`
	ScopeEnvironmentId  pingonetypes.ResourceIDValue `tfsdk:"scope_environment_id"`
	ScopeOrganizationId pingonetypes.ResourceIDValue `tfsdk:"scope_organization_id"`
	ScopePopulationId   pingonetypes.ResourceIDValue `tfsdk:"scope_population_id"`
	ReadOnly            types.Bool                   `tfsdk:"read_only"`
}

// Framework interfaces
var (
	_ resource.Resource                = &GatewayRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &GatewayRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &GatewayRoleAssignmentResource{}
)

// New Object
func NewGatewayRoleAssignmentResource() resource.Resource {
	return &GatewayRoleAssignmentResource{}
}

// Metadata
func (r *GatewayRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_role_assignment"
}

// Schema.
func (r *GatewayRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	attributes := map[string]schema.Attribute{
		"id": framework.Attr_ID(),

		"environment_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the gateway to assign the admin role to."),
		),

		"gateway_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an gateway to assign an admin role to."),
		),

		"role_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an admin role to assign to the gateway."),
		),

		"read_only": schema.BoolAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A flag to show whether the admin role assignment is read only or can be changed.").Description,
			Computed:    true,
		},
	}

	utils.MergeSchemaAttributeMaps(attributes, service.RoleAssignmentScopeSchema(), true)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to administrator defined gateways.",

		Attributes: attributes,
	}
}

func (r *GatewayRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GatewayRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state GatewayRoleAssignmentResourceModel

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
	gatewayRoleAssignment, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.GatewayRoleAssignmentsApi.CreateGatewayRoleAssignment(ctx, plan.EnvironmentId.ValueString(), plan.GatewayId.ValueString()).RoleAssignment(*gatewayRoleAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateGatewayRoleAssignment",
		service.CreateRoleAssignmentErrorFunc,
		service.RoleAssignmentRetryable,
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

func (r *GatewayRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GatewayRoleAssignmentResourceModel

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
	var response *management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.GatewayRoleAssignmentsApi.ReadOneGatewayRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.GatewayId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneGatewayRoleAssignment",
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

func (r *GatewayRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *GatewayRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GatewayRoleAssignmentResourceModel

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

	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.GatewayRoleAssignmentsApi.DeleteGatewayRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.GatewayId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteGatewayRoleAssignment",
		framework.CustomErrorResourceNotFoundWarning,
		service.RoleRemovalRetryable,
		nil,
	)...)

}

func (r *GatewayRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "gateway_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "role_assignment_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
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

func (p *GatewayRoleAssignmentResourceModel) expand() (*management.RoleAssignment, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopeID, scopeType, d := service.ExpandRoleAssignmentScope(p.ScopeEnvironmentId, p.ScopeOrganizationId, p.ScopePopulationId, p.ScopeApplicationId)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	gatewayRoleAssignmentRole := *management.NewRoleAssignmentRole(p.RoleId.ValueString())
	gatewayRoleAssignmentScope := *management.NewRoleAssignmentScope(scopeID, management.EnumRoleAssignmentScopeType(scopeType))
	data := management.NewRoleAssignment(gatewayRoleAssignmentRole, gatewayRoleAssignmentScope)

	return data, diags
}

func (p *GatewayRoleAssignmentResourceModel) toState(apiObject *management.RoleAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(apiObject.Environment.GetIdOk())
	p.RoleId = framework.PingOneResourceIDOkToTF(apiObject.Role.GetIdOk())
	p.ReadOnly = framework.BoolOkToTF(apiObject.GetReadOnlyOk())

	p.ScopeEnvironmentId, p.ScopeOrganizationId, p.ScopePopulationId, p.ScopeApplicationId = service.RoleAssignmentScopeOkToTF(apiObject.GetScopeOk())

	return diags
}
