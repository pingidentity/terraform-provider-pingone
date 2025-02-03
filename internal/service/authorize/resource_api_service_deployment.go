// Copyright © 2025 Ping Identity Corporation

package authorize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type APIServiceDeploymentResource serviceClientType

type APIServiceDeploymentResourceModel struct {
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	APIServiceId              pingonetypes.ResourceIDValue `tfsdk:"api_service_id"`
	AuthorizationVersion      types.Object                 `tfsdk:"authorization_version"`
	DecisionEndpoint          types.Object                 `tfsdk:"decision_endpoint"`
	DeployedAt                timetypes.RFC3339            `tfsdk:"deployed_at"`
	Policy                    types.Object                 `tfsdk:"policy"`
	Status                    types.Object                 `tfsdk:"status"`
	RedeploymentTriggerValues types.Map                    `tfsdk:"redeployment_trigger_values"`
}

type APIServiceDeploymentAuthorizationVersionResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceDeploymentDecisionEndpointResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceDeploymentPolicyResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceDeploymentStatusResourceModel struct {
	Code  types.String `tfsdk:"code"`
	Error types.Object `tfsdk:"error"`
}

type APIServiceDeploymentStatusErrorResourceModel struct {
	Id      pingonetypes.ResourceIDValue `tfsdk:"id"`
	Code    types.String                 `tfsdk:"code"`
	Message types.String                 `tfsdk:"message"`
}

var (
	apiServiceDeploymentAuthorizationVersionTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceDeploymentDecisionEndpointTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceDeploymentPolicyTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceDeploymentStatusTFObjectTypes = map[string]attr.Type{
		"code":  types.StringType,
		"error": types.ObjectType{AttrTypes: apiServiceDeploymentStatusErrorTFObjectTypes},
	}

	apiServiceDeploymentStatusErrorTFObjectTypes = map[string]attr.Type{
		"id":      pingonetypes.ResourceIDType{},
		"code":    types.StringType,
		"message": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &APIServiceDeploymentResource{}
	_ resource.ResourceWithConfigure   = &APIServiceDeploymentResource{}
	_ resource.ResourceWithImportState = &APIServiceDeploymentResource{}
	_ resource.ResourceWithModifyPlan  = &APIServiceDeploymentResource{}
)

// New Object
func NewAPIServiceDeploymentResource() resource.Resource {
	return &APIServiceDeploymentResource{}
}

// Metadata
func (r *APIServiceDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_api_service_deployment"
}

// Schema.
func (r *APIServiceDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	statusCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that describes the deployment status code. For possible values, see [Deployment status codes](https://apidocs.pingidentity.com/pingone/platform/v1/api/#service-deployment-status-codes).",
	)

	statusErrorCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that describes a general fault code that identifies the the type of error. See [Error codes](https://apidocs.pingidentity.com/pingone/platform/v1/api/#error-codes).",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the deployment of an API Service for PingOne Authorize in an environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to deploy the API Service in."),
			),

			"api_service_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the API service to deploy."),
			),

			"authorization_version": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes properties related to the authorization version that relates to the API service that has been deployed.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the UUID of the last deployed policy authorization version. This is present only if custom polcies are enabled and the API service has been deployed at least once.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"decision_endpoint": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes properties related to the decision endpoint that relates to the API service that has been deployed.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the UUID of the decision endpoint.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"deployed_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An RFC3339 compliant date/time string that specifies the time of the most recent successful deployment. The field will be null if the API service has never been successfully deployed.").Description,
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"policy": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes properties related to the root policy that relates to the API service that has been deployed.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the UUID of the root policy.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"status": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes properties related to the status of the API service deployment.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         statusCodeDescription.Description,
						MarkdownDescription: statusCodeDescription.MarkdownDescription,
						Computed:            true,
					},

					"error": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes properties related to an error status of the API service deployment.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the UUID of the root policy.").Description,
								Computed:    true,

								CustomType: pingonetypes.ResourceIDType{},
							},

							"code": schema.StringAttribute{
								Description:         statusErrorCodeDescription.Description,
								MarkdownDescription: statusErrorCodeDescription.MarkdownDescription,
								Computed:            true,
							},

							"message": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes a short, human-readable description of the error.").Description,
								Computed:    true,
							},
						},
					},
				},
			},

			"redeployment_trigger_values": schema.MapAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A meta-argument map of values that, if any values are changed, will force redeployment.  Adding values to and removing values from the map will not trigger a deployment.").Description,
				Optional:    true,

				ElementType: types.StringType,
			},
		},
	}
}

// ModifyPlan
func (r *APIServiceDeploymentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state types.Map
	var planValues, stateValues map[string]attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("redeployment_trigger_values"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues = plan.Elements()

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("redeployment_trigger_values"), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateValues = state.Elements()

	for k, v := range planValues {
		if stateValue, ok := stateValues[k]; ok && (v == types.StringUnknown() || !stateValue.Equal(v)) {
			resp.RequiresReplace = path.Paths{path.Root("redeployment_trigger_values")}
			break
		}
	}

}

func (r *APIServiceDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APIServiceDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state APIServiceDeploymentResourceModel

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

	// Run the API call
	var response *authorize.APIServerDeployment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerDeploymentApi.DeployAPIServer(ctx, plan.EnvironmentId.ValueString(), plan.APIServiceId.ValueString()).ContentType("application/vnd.pingidentity.apiserver.deploy+json").Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"DeployAPIServer",
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

func (r *APIServiceDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *APIServiceDeploymentResourceModel

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
	var response *authorize.APIServerDeployment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerDeploymentApi.ReadDeploymentStatus(ctx, data.EnvironmentId.ValueString(), data.APIServiceId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadDeploymentStatus",
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

func (r *APIServiceDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state APIServiceDeploymentResourceModel

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

	// Run the API call
	var response *authorize.APIServerDeployment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerDeploymentApi.ReadDeploymentStatus(ctx, plan.EnvironmentId.ValueString(), plan.APIServiceId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadDeploymentStatus",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	state = plan

	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *APIServiceDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *APIServiceDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "api_service_id",
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

func (p *APIServiceDeploymentResourceModel) toState(apiObject *authorize.APIServerDeployment) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.AuthorizationVersion, d = apiServiceDeploymentAuthorizationVersionOkToTF(apiObject.GetAuthorizationVersionOk())
	diags.Append(d...)

	p.DecisionEndpoint, d = apiServiceDeploymentDecisionEndpointOkToTF(apiObject.GetDecisionEndpointOk())
	diags.Append(d...)

	p.DeployedAt = framework.TimeOkToTF(apiObject.GetDeployedAtOk())

	p.Policy, d = apiServiceDeploymentPolicyOkToTF(apiObject.GetPolicyOk())
	diags.Append(d...)

	p.Status, d = apiServiceDeploymentStatusOkToTF(apiObject.GetStatusOk())
	diags.Append(d...)

	return diags
}

func apiServiceDeploymentAuthorizationVersionOkToTF(apiObject *authorize.APIServerDeploymentAuthorizationVersion, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDeploymentAuthorizationVersionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDeploymentAuthorizationVersionTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceDeploymentDecisionEndpointOkToTF(apiObject *authorize.APIServerDeploymentDecisionEndpoint, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDeploymentDecisionEndpointTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDeploymentDecisionEndpointTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceDeploymentPolicyOkToTF(apiObject *authorize.APIServerDeploymentPolicy, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDeploymentPolicyTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDeploymentPolicyTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceDeploymentStatusOkToTF(apiObject *authorize.APIServerDeploymentStatus, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDeploymentStatusTFObjectTypes), diags
	}

	errorObj, d := apiServiceDeploymentStatusErrorOkToTF(apiObject.GetErrorOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceDeploymentStatusTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDeploymentStatusTFObjectTypes, map[string]attr.Value{
		"code":  framework.StringOkToTF(apiObject.GetCodeOk()),
		"error": errorObj,
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceDeploymentStatusErrorOkToTF(apiObject *authorize.APIServerDeploymentStatusError, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDeploymentStatusErrorTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDeploymentStatusErrorTFObjectTypes, map[string]attr.Value{
		"id":      framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"code":    framework.StringOkToTF(apiObject.GetCodeOk()),
		"message": framework.StringOkToTF(apiObject.GetMessageOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
