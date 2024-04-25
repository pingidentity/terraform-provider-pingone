package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ResourceScopeResource serviceClientType

type ResourceScopeResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId    pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ResourceScopeResource{}
	_ resource.ResourceWithConfigure   = &ResourceScopeResource{}
	_ resource.ResourceWithImportState = &ResourceScopeResource{}
)

// New Object
func NewResourceScopeResource() resource.Resource {
	return &ResourceScopeResource{}
}

// Metadata
func (r *ResourceScopeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_scope"
}

// Schema.
func (r *ResourceScopeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 resource scopes for custom resources.  This resource cannot manage PingOne API or OpenID Connect scopes.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the resource scope in."),
			),

			"resource_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the resource to assign the resource scope to."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The name of the resource scope.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description to apply to the resource scope.").Description,
				Optional:    true,
			},
		},
	}
}

func (r *ResourceScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceScopeResourceModel

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

	resource, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceScope := plan.expand()

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.CreateResourceScope(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString()).ResourceScope(*resourceScope).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateResourceScope",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resourceScopeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceScopeResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceScopeResourceModel

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

	resource, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resource == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Validate the plan
	resp.Diagnostics.Append(data.validate(*resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.ReadOneResourceScope(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneResourceScope",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&resourceScopeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resourceScopeResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resourceScopeResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceScopeResourceModel

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

	resource, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceScope := plan.expand()

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.UpdateResourceScope(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Id.ValueString()).ResourceScope(*resourceScope).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateResourceScope",
		framework.DefaultCustomError,
		nil,
		&resourceScopeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceScopeResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceScopeResourceModel

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
			fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.DeleteResourceScope(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteResourceScope",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "resource_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "resource_scope_id",
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

func (p *ResourceScopeResourceModel) expand() *management.ResourceScope {
	data := management.NewResourceScope(p.Name.ValueString())

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	return data
}

func (p *ResourceScopeResourceModel) validate(resource management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the `openid` scope from the `openid` resource is not in the list
	if v, ok := resource.GetTypeOk(); ok && *v != management.ENUMRESOURCETYPE_CUSTOM {
		diags.AddError(
			"Invalid resource",
			"This resource cannot control scopes for resources that are of type PingOne API or OpenID Connect.  Please ensure that the resource in the `resource_id` parameter is a custom resource or consider using the `pingone_resource_scope_openid` or `pingone_resource_scope_pingone_api` provider resources.",
		)
	}

	return diags
}

func (p *ResourceScopeResourceModel) toState(apiObject *management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	return diags
}
