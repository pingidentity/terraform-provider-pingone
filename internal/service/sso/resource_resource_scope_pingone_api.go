package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
type ResourceScopePingOneAPIResource serviceClientType

type ResourceScopePingOneAPIResourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId       pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	Name             types.String                 `tfsdk:"name"`
	Description      types.String                 `tfsdk:"description"`
	SchemaAttributes types.Set                    `tfsdk:"schema_attributes"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ResourceScopePingOneAPIResource{}
	_ resource.ResourceWithConfigure   = &ResourceScopePingOneAPIResource{}
	_ resource.ResourceWithImportState = &ResourceScopePingOneAPIResource{}
)

// New Object
func NewResourceScopePingOneAPIResource() resource.Resource {
	return &ResourceScopePingOneAPIResource{}
}

// Metadata
func (r *ResourceScopePingOneAPIResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_scope_pingone_api"
}

// Schema
func (r *ResourceScopePingOneAPIResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource scope.  Predefined scopes of `p1:read:user` and `p1:update:user` can be overridden, and new scopes can be defined as subscopes in the format `p1:read:user:{suffix}` or `p1:update:user:{suffix}`.  E.g. `p1:read:user:newscope` or `p1:update:user:newscope`",
	)

	schemaAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the user schema attributes that can be read or updated for the specified PingOne access control scope. The value is an array of schema attribute paths (such as `username`, `name.given`, `shirtSize`) that the scope controls. This property is supported only for the `p1:read:user`, `p1:update:user` and `p1:read:user:{suffix}` and `p1:update:user:{suffix}` scopes. No other PingOne platform scopes allow this behavior. Any attributes not listed in the attribute array are excluded from the read or update action. The wildcard path (`*`) in the array includes all attributes and cannot be used in conjunction with any other user schema attribute paths.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource scopes for the PingOne API resource.  Predefined scopes of `p1:read:user` and `p1:update:user` can be overridden, and new scopes can be defined as subscopes in the format `p1:read:user:{suffix}` or `p1:update:user:{suffix}`.  E.g. `p1:read:user:newscope` or `p1:update:user:newscope`.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the resource scope in."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^p1:(read|update):user(:{1}[a-zA-Z0-9]+)*$`),
						"Resource scope name must be either `p1:read:user`, `p1:update:user`, `p1:read:user:{suffix}` or `p1:update:user:{suffix}`",
					),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description to apply to the resource scope.  The description can only be set when defining new scopes.").Description,
				Optional:    true,
				Computed:    true,
			},

			"schema_attributes": schema.SetAttribute{
				Description:         schemaAttributesDescription.Description,
				MarkdownDescription: schemaAttributesDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(attrMinLength),
					),
				},
			},

			"resource_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the PingOne API resource.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (r *ResourceScopePingOneAPIResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceScopePingOneAPIResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceScopePingOneAPIResourceModel

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

	resource, d := fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), "PingOne API", false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ResourceId = framework.PingOneResourceIDOkToTF(resource.GetIdOk())

	// Build the model for the API
	resourceScope, d := plan.expand(ctx, r.Client.ManagementAPIClient, *resource)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	if v, ok := resourceScope.GetIdOk(); ok {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.UpdateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId(), *v).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceScope-PingOneAPI-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&resourceScopeResponse,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.CreateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId()).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateResourceScope-PingOneAPI-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&resourceScopeResponse,
		)...)

	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceScopeResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceScopePingOneAPIResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceScopePingOneAPIResourceModel

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

	resource, d := fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), "PingOne API", true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resource == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.ReadOneResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneResourceScope-PingOneAPI",
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

func (r *ResourceScopePingOneAPIResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceScopePingOneAPIResourceModel

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

	resource, d := fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), "PingOne API", false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ResourceId = framework.PingOneResourceIDOkToTF(resource.GetIdOk())

	// Build the model for the API
	resourceScope, d := plan.expand(ctx, r.Client.ManagementAPIClient, *resource)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.UpdateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId(), plan.Id.ValueString()).ResourceScope(*resourceScope).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateResourceScope-PingOneAPI",
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

func (r *ResourceScopePingOneAPIResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceScopePingOneAPIResourceModel

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

	resource, d := fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), "PingOne API", true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resource == nil {
		return
	}

	m, err := regexp.MatchString("^p1:(read|update):user$", data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid resource scope name",
			fmt.Sprintf("Cannot determine if the resource scope is a predefined scope: %s", err),
		)
		return
	}

	if m {

		resourceScope, d := fetchResourceScopeFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), resource.GetId(), data.Name.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resourceScope.SetSchemaAttributes([]string{"*"})

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.UpdateResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceScope-PingOneAPI-Delete",
			framework.DefaultCustomError,
			nil,
			nil,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := r.Client.ManagementAPIClient.ResourceScopesApi.DeleteResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
			},
			"DeleteResourceScope-PingOneAPI-Delete",
			framework.CustomErrorResourceNotFoundWarning,
			nil,
			nil,
		)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceScopePingOneAPIResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
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

func (p *ResourceScopePingOneAPIResourceModel) expand(ctx context.Context, apiClient *management.APIClient, resource management.Resource) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ResourceScope

	newScope := true
	m, err := regexp.MatchString("^p1:(read|update):user$", p.Name.ValueString())
	if err != nil {
		diags.AddError(
			"Invalid resource scope name",
			fmt.Sprintf("Cannot determine if the resource scope is a predefined scope: %s", err),
		)
		return nil, diags
	}

	if m {
		newScope = false

		data, diags = fetchResourceScopeFromName(ctx, apiClient, p.EnvironmentId.ValueString(), resource.GetId(), p.Name.ValueString())
		if diags.HasError() {
			return nil, diags
		}

	} else {
		data = management.NewResourceScope(p.Name.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		if newScope {
			data.SetDescription(p.Description.ValueString())
		} else {
			diags.AddError(
				"Invalid attribute value",
				"Cannot update the description of a predefined scope.",
			)
		}
	}

	if !p.SchemaAttributes.IsNull() && !p.SchemaAttributes.IsUnknown() {

		var plan []string
		diags.Append(p.SchemaAttributes.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetSchemaAttributes(plan)
	}

	return data, diags
}

func (p *ResourceScopePingOneAPIResourceModel) validate(resource management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the `openid` scope from the `openid` resource is not in the list
	if v, ok := resource.GetTypeOk(); ok && *v != management.ENUMRESOURCETYPE_PINGONE_API {
		diags.AddError(
			"Invalid resource",
			"This resource cannot control scopes for resources that are of type PingOne API or PingOneAPI Connect.  Please ensure that the resource in the `resource_id` parameter is a custom resource or consider using the `pingone_resource_scope_openid` or `pingone_resource_scope_pingone_api` provider resources.",
		)
	}

	return diags
}

func (p *ResourceScopePingOneAPIResourceModel) toState(apiObject *management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())

	if v, ok := apiObject.GetResourceOk(); ok {
		p.ResourceId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.ResourceId = pingonetypes.NewResourceIDNull()
	}

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.SchemaAttributes = framework.StringSetOkToTF(apiObject.GetSchemaAttributesOk())

	return diags
}
