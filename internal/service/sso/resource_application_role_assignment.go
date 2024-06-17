package sso

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationRoleAssignmentResource serviceClientType

type ApplicationRoleAssignmentResourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId       pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	RoleId              pingonetypes.ResourceIDValue `tfsdk:"role_id"`
	ScopeEnvironmentId  pingonetypes.ResourceIDValue `tfsdk:"scope_environment_id"`
	ScopeOrganizationId pingonetypes.ResourceIDValue `tfsdk:"scope_organization_id"`
	ScopePopulationId   pingonetypes.ResourceIDValue `tfsdk:"scope_population_id"`
	ReadOnly            types.Bool                   `tfsdk:"read_only"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &ApplicationRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &ApplicationRoleAssignmentResource{}
)

// New Object
func NewApplicationRoleAssignmentResource() resource.Resource {
	return &ApplicationRoleAssignmentResource{}
}

// Metadata
func (r *ApplicationRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_role_assignment"
}

// Schema.
func (r *ApplicationRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	attributes := map[string]schema.Attribute{
		"id": framework.Attr_ID(),

		"environment_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the application to assign the admin role to."),
		),

		"application_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an application to assign an admin role to."),
		),

		"role_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an admin role to assign to the application."),
		),

		"read_only": schema.BoolAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A flag to show whether the admin role assignment is read only or can be changed.").Description,
			Computed:    true,
		},
	}

	utils.MergeSchemaAttributeMaps(attributes, service.RoleAssignmentScopeSchema(), true)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to administrator defined applications.",

		Attributes: attributes,
	}
}

func (r *ApplicationRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationRoleAssignmentResourceModel

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
	applicationRoleAssignment, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check the application can have roles assigned
	application, d := fetchApplication(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if application == nil {
		resp.Diagnostics.AddError(
			"Invalid parameter value - Application not found",
			fmt.Sprintf("The application ID provided (%s) does not exist in the environment.", plan.ApplicationId.ValueString()),
		)

		return
	}

	if !checkApplicationTypeForRoleAssignment(*application) {
		resp.Diagnostics.AddError(
			"Invalid parameter value - Unmappable application type",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `OPENID_CONNECT` or `SAML` type.  Roles cannot be mapped to this application.", plan.ApplicationId.ValueString()),
		)

		return
	}

	// Run the API call
	var response *management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationRoleAssignmentsApi.CreateApplicationRoleAssignment(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).RoleAssignment(*applicationRoleAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplicationRoleAssignment",
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

func (r *ApplicationRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationRoleAssignmentResourceModel

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

	// Check the application can have roles assigned
	application, d := fetchApplication(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if application == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if !checkApplicationTypeForRoleAssignment(*application) {
		resp.Diagnostics.AddError(
			"Invalid parameter value - Unmappable application type",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `OPENID_CONNECT` or `SAML` type.  Roles cannot be mapped to this application.", data.ApplicationId.ValueString()),
		)

		return
	}

	// Run the API call
	var response *management.RoleAssignment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationRoleAssignmentsApi.ReadOneApplicationRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplicationRoleAssignment",
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

func (r *ApplicationRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ApplicationRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationRoleAssignmentResourceModel

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

	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{
			"400",
		},
		Target: []string{
			"204",
			"404",
		},
		Refresh: func() (interface{}, string, error) {
			var fR *http.Response
			var fErr error
			resp.Diagnostics.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					fR, fErr = r.Client.ManagementAPIClient.ApplicationRoleAssignmentsApi.DeleteApplicationRoleAssignment(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
				},
				"DeleteApplicationRoleAssignment",
				framework.CustomErrorResourceNotFoundWarning,
				nil,
				nil,
			)...)

			statusCode := ""
			if fR != nil {
				base := 10
				statusCode = strconv.FormatInt(int64(fR.StatusCode), base)
			}

			return resp, statusCode, fErr
		},
		Timeout:                   20 * time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Role Assignment Delete Error",
			fmt.Sprintf("Error waiting for application role assignment (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}

}

func (r *ApplicationRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (p *ApplicationRoleAssignmentResourceModel) expand() (*management.RoleAssignment, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopeID, scopeType, d := service.ExpandRoleAssignmentScope(p.ScopeEnvironmentId, p.ScopeOrganizationId, p.ScopePopulationId)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	applicationRoleAssignmentRole := *management.NewRoleAssignmentRole(p.RoleId.ValueString())
	applicationRoleAssignmentScope := *management.NewRoleAssignmentScope(scopeID, management.EnumRoleAssignmentScopeType(scopeType))
	data := management.NewRoleAssignment(applicationRoleAssignmentRole, applicationRoleAssignmentScope)

	return data, diags
}

func (p *ApplicationRoleAssignmentResourceModel) toState(apiObject *management.RoleAssignment) diag.Diagnostics {
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

	p.ScopeEnvironmentId, p.ScopeOrganizationId, p.ScopePopulationId = service.RoleAssignmentScopeOkToTF(apiObject.GetScopeOk())

	return diags
}

func fetchApplication(ctx context.Context, apiClient *management.APIClient, environmentId, applicationId string, warnIfNotFound bool) (*management.ReadOneApplication200Response, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError
	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	var response *management.ReadOneApplication200Response
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.ReadOneApplication(ctx, environmentId, applicationId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
		},
		"ReadOneApplication",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if response == nil {
		return nil, diags
	}

	return response, diags
}

func checkApplicationTypeForRoleAssignment(application management.ReadOneApplication200Response) bool {
	if application.ApplicationOIDC != nil && application.ApplicationOIDC.GetId() != "" {
		return true
	}

	if application.ApplicationSAML != nil && application.ApplicationSAML.GetId() != "" {
		return true
	}

	return false
}
