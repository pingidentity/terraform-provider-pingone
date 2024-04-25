package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationFlowPolicyAssignmentResource serviceClientType

type ApplicationFlowPolicyAssignmentResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	FlowPolicyId  davincitypes.ResourceIDValue `tfsdk:"flow_policy_id"`
	Priority      types.Int64                  `tfsdk:"priority"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationFlowPolicyAssignmentResource{}
	_ resource.ResourceWithConfigure   = &ApplicationFlowPolicyAssignmentResource{}
	_ resource.ResourceWithImportState = &ApplicationFlowPolicyAssignmentResource{}
)

// New Object
func NewApplicationFlowPolicyAssignmentResource() resource.Resource {
	return &ApplicationFlowPolicyAssignmentResource{}
}

// Metadata
func (r *ApplicationFlowPolicyAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_flow_policy_assignment"
}

// Schema
func (r *ApplicationFlowPolicyAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a DaVinci flow policy assignment for an administrator defined application or built-in system applications configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the application flow policy assignment in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to create the flow policy assignment for. The value for `application_id` may come from the `id` attribute of the `pingone_application` or `pingone_system_application` resources or data sources."),
			),

			"flow_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the DaVinci flow policy to associate.").Description,
				Required:    true,

				CustomType: davincitypes.ResourceIDType{},
			},

			"priority": schema.Int64Attribute{
				Description: "The order in which the policy referenced by this assignment is evaluated during an authentication flow relative to other policies. An assignment with a lower priority will be evaluated first.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
		},
	}
}

func (r *ApplicationFlowPolicyAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationFlowPolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationFlowPolicyAssignmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	applicationFlowPolicyAssignment := plan.expand()

	// Run the API call
	var response *management.FlowPolicyAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationFlowPolicyAssignmentsApi.CreateFlowPolicyAssignment(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).FlowPolicyAssignment(*applicationFlowPolicyAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateFlowPolicyAssignment",
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

func (r *ApplicationFlowPolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationFlowPolicyAssignmentResourceModel

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
	var response *management.FlowPolicyAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationFlowPolicyAssignmentsApi.ReadOneFlowPolicyAssignment(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneFlowPolicyAssignment",
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

func (r *ApplicationFlowPolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationFlowPolicyAssignmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	applicationFlowPolicyAssignment := plan.expand()

	// Run the API call
	var response *management.FlowPolicyAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationFlowPolicyAssignmentsApi.UpdateFlowPolicyAssignment(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).FlowPolicyAssignment(*applicationFlowPolicyAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateFlowPolicyAssignment",
		framework.DefaultCustomError,
		nil,
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

func (r *ApplicationFlowPolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationFlowPolicyAssignmentResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.ApplicationFlowPolicyAssignmentsApi.DeleteFlowPolicyAssignment(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteFlowPolicyAssignment",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationFlowPolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "flow_policy_assignment_id",
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

func (p *ApplicationFlowPolicyAssignmentResourceModel) expand() *management.FlowPolicyAssignment {

	flowPolicy := management.NewFlowPolicyAssignmentFlowPolicy(p.FlowPolicyId.ValueString())
	data := management.NewFlowPolicyAssignment(*flowPolicy, int32(p.Priority.ValueInt64()))

	return data
}

func (p *ApplicationFlowPolicyAssignmentResourceModel) toState(apiObject *management.FlowPolicyAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.ApplicationId = framework.PingOneResourceIDToTF(*apiObject.GetApplication().Id)
	p.FlowPolicyId = framework.DaVinciResourceIDToTF(apiObject.GetFlowPolicy().Id)
	p.Priority = framework.Int32OkToTF(apiObject.GetPriorityOk())

	return diags
}
