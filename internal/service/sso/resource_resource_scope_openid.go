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
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ResourceScopeOpenIDResource serviceClientType

type ResourceScopeOpenIDResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ResourceId    types.String `tfsdk:"resource_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	MappedClaims  types.Set    `tfsdk:"mapped_claims"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ResourceScopeOpenIDResource{}
	_ resource.ResourceWithConfigure   = &ResourceScopeOpenIDResource{}
	_ resource.ResourceWithImportState = &ResourceScopeOpenIDResource{}
)

// New Object
func NewResourceScopeOpenIDResource() resource.Resource {
	return &ResourceScopeOpenIDResource{}
}

// Metadata
func (r *ResourceScopeOpenIDResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_scope_openid"
}

// Schema
func (r *ResourceScopeOpenIDResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource scope.  Predefined scopes of `address`, `email`, `openid`, `phone` and `profile` can be overridden, and new scopes can be defined.  E.g. `myawesomescope`",
	)

	mappedClaimsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of custom resource attribute IDs.  This property does not control predefined OpenID Connect (OIDC) mappings, such as the `email` claim in the OIDC `email` scope or the `name` claim in the `profile` scope. You can create custom attributes, and these custom attributes can be added to `mapped_claims` and will display in the response.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource scopes for the OpenID Connect resource.  Predefined scopes can be overridden, and new scopes can be defined.",

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
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description to apply to the resource scope.  The description can only be set when defining new scopes.").Description,
				Optional:    true,
				Computed:    true,
			},

			"mapped_claims": schema.SetAttribute{
				Description:         mappedClaimsDescription.Description,
				MarkdownDescription: mappedClaimsDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						verify.P1ResourceIDValidator(),
					),
				},
			},

			"resource_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the OpenID Connect resource.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *ResourceScopeOpenIDResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.Client = preparedClient
}

func (r *ResourceScopeOpenIDResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceScopeOpenIDResourceModel

	if r.Client == nil {
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

	resource, d := fetchResourceFromName(ctx, r.Client, plan.EnvironmentId.ValueString(), "openid", false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ResourceId = framework.StringOkToTF(resource.GetIdOk())

	// Build the model for the API
	resourceScope, d := plan.expand(ctx, r.Client, *resource)
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
				fO, fR, fErr := r.Client.ResourceScopesApi.UpdateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId(), *v).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceScope-OpenID-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&resourceScopeResponse,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ResourceScopesApi.CreateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId()).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateResourceScope-OpenID-Create",
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

func (r *ResourceScopeOpenIDResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceScopeOpenIDResourceModel

	if r.Client == nil {
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

	resource, d := fetchResourceFromName(ctx, r.Client, data.EnvironmentId.ValueString(), "openid", true)
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
			fO, fR, fErr := r.Client.ResourceScopesApi.ReadOneResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneResourceScope-OpenID",
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

func (r *ResourceScopeOpenIDResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceScopeOpenIDResourceModel

	if r.Client == nil {
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

	resource, d := fetchResourceFromName(ctx, r.Client, plan.EnvironmentId.ValueString(), "openid", false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceScope, d := plan.expand(ctx, r.Client, *resource)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceScopeResponse *management.ResourceScope
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ResourceScopesApi.UpdateResourceScope(ctx, plan.EnvironmentId.ValueString(), resource.GetId(), plan.Id.ValueString()).ResourceScope(*resourceScope).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateResourceScope-OpenID",
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

func (r *ResourceScopeOpenIDResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceScopeOpenIDResourceModel

	if r.Client == nil {
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

	resource, d := fetchResourceFromName(ctx, r.Client, data.EnvironmentId.ValueString(), "openid", true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resource == nil {
		return
	}

	if m, err := regexp.MatchString("^(address|email|openid|phone|profile)$", data.Name.ValueString()); err == nil && m {

		resourceScope, d := fetchResourceScopeFromName(ctx, r.Client, data.EnvironmentId.ValueString(), resource.GetId(), data.Name.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resourceScope.SetMappedClaims([]string{})

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ResourceScopesApi.UpdateResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).ResourceScope(*resourceScope).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceScope-OpenID-Delete",
			framework.DefaultCustomError,
			nil,
			nil,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := r.Client.ResourceScopesApi.DeleteResourceScope(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, data.EnvironmentId.ValueString(), nil, fR, fErr)
			},
			"DeleteResourceScope-OpenID-Delete",
			framework.CustomErrorResourceNotFoundWarning,
			nil,
			nil,
		)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceScopeOpenIDResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (p *ResourceScopeOpenIDResourceModel) expand(ctx context.Context, apiClient *management.APIClient, resource management.Resource) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ResourceScope

	newScope := true
	if m, err := regexp.MatchString("^(address|email|openid|phone|profile)$", p.Name.ValueString()); err == nil && m {
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

	if !p.MappedClaims.IsNull() && !p.MappedClaims.IsUnknown() {

		var plan []string
		diags.Append(p.MappedClaims.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetMappedClaims(plan)
	}

	return data, diags
}

func (p *ResourceScopeOpenIDResourceModel) validate(resource management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the `openid` scope from the `openid` resource is not in the list
	if v, ok := resource.GetTypeOk(); ok && *v != management.ENUMRESOURCETYPE_OPENID_CONNECT {
		diags.AddError(
			"Invalid resource",
			"This resource cannot control scopes for resources that are of type PingOne API or OpenID Connect.  Please ensure that the resource in the `resource_id` parameter is a custom resource or consider using the `pingone_resource_scope_openid` or `pingone_resource_scope_pingone_api` provider resources.",
		)
	}

	return diags
}

func (p *ResourceScopeOpenIDResourceModel) toState(apiObject *management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.MappedClaims = framework.StringSetOkToTF(apiObject.GetMappedClaimsOk())

	return diags
}
