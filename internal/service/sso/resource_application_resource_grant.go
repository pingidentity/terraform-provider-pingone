package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResourceGrantResource serviceClientType

type ApplicationResourceGrantResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	ResourceId    types.String `tfsdk:"resource_id"`
	ResourceName  types.String `tfsdk:"resource_name"`
	Scopes        types.Set    `tfsdk:"scopes"`
	ScopeNames    types.Set    `tfsdk:"scope_names"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationResourceGrantResource{}
	_ resource.ResourceWithConfigure   = &ApplicationResourceGrantResource{}
	_ resource.ResourceWithImportState = &ApplicationResourceGrantResource{}
)

// New Object
func NewApplicationResourceGrantResource() resource.Resource {
	return &ApplicationResourceGrantResource{}
}

// Metadata
func (r *ApplicationResourceGrantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_resource_grant"
}

// Schema.
func (r *ApplicationResourceGrantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice**: This parameter is deprecated and will be made read-only in a future release.  This attribute should be replaced with the `resource_name` parameter instead.  The ID of the resource to assign the resource attribute to.",
	).ExactlyOneOf([]string{"resource_id", "resource_name"}).AppendMarkdownString("Must be a valid PingOne resource ID.").RequiresReplace()

	resourceNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource to assign to the application.  The built-in OpenID Connect resource name is `openid` and the built-in PingOne API resource anem is `PingOne API`.",
	).ExactlyOneOf([]string{"resource_id", "resource_name"}).RequiresReplace()

	scopesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice**: This parameter is deprecated and will be made read-only in a future release.  This attribute should be replaced with the `scope_names` parameter instead.  A list of IDs of the scopes associated with this grant.  When using the `openid` resource, the `openid` scope should not be included.",
	).ExactlyOneOf([]string{"scopes", "scope_names"})

	scopeNamesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of scopes by name that should be associated with this grant.  For example, `profile`, `email` etc.  When using the `openid` resource, the `openid` scope should not be included.",
	).ExactlyOneOf([]string{"scopes", "scope_names"})

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a resource grant for administrator defined applications or built-in system applications configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the application resource grant in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to create the resource grant for.  The value for `application_id` may come from the `id` attribute of the `pingone_application` or `pingone_system_application` resources or data sources."),
			),

			"resource_id": schema.StringAttribute{
				Description:         resourceIdDescription.Description,
				MarkdownDescription: resourceIdDescription.MarkdownDescription,
				DeprecationMessage:  "This parameter is deprecated and will be made read-only in a future release.  This attribute should be replaced with the `resource_name` parameter instead.",
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("resource_id"),
						path.MatchRoot("resource_name"),
					),
				},
			},

			"resource_name": schema.StringAttribute{
				Description:         resourceNameDescription.Description,
				MarkdownDescription: resourceNameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("resource_id"),
						path.MatchRoot("resource_name"),
					),
				},
			},

			"scopes": schema.SetAttribute{
				Description:         scopesDescription.Description,
				MarkdownDescription: scopesDescription.MarkdownDescription,
				DeprecationMessage:  "This parameter is deprecated and will be made read-only in a future release.  This attribute should be replaced with the `scope_names` parameter instead.",
				Optional:            true,
				Computed:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.ValueStringsAre(
						verify.P1ResourceIDValidator(),
					),
					setvalidator.ExactlyOneOf(
						path.MatchRoot("scopes"),
						path.MatchRoot("scope_names"),
					),
				},
			},

			"scope_names": schema.SetAttribute{
				Description:         scopeNamesDescription.Description,
				MarkdownDescription: scopeNamesDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(attrMinLength),
					),
					setvalidator.ExactlyOneOf(
						path.MatchRoot("scopes"),
						path.MatchRoot("scope_names"),
					),
				},
			},
		},
	}
}

func (r *ApplicationResourceGrantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResourceGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

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

	resource, resourceScopes, d := plan.getResource(ctx, r.Client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource, resourceScopes)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant := plan.expand(*resource, resourceScopes)

	// Run the API call
	var grantResponse *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ApplicationResourceGrantsApi.CreateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
		},
		"CreateApplicationGrant",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&grantResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource response
	resourceResponse, d := fetchResourceFromID(ctx, r.Client, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource scopes response
	scopeIds := make([]string, 0, len(grantResponse.Scopes))
	for _, scope := range grantResponse.Scopes {
		scopeIds = append(scopeIds, scope.GetId())
	}

	resourceScopesResponse, d := fetchResourceScopesFromIDs(ctx, r.Client, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), scopeIds)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(grantResponse, resourceResponse, resourceScopesResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourceGrantResourceModel

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

	// Run the API call
	var grantResponse *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplicationGrant",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&grantResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if grantResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Get the resource response
	resourceResponse, d := fetchResourceFromID(ctx, r.Client, data.EnvironmentId.ValueString(), grantResponse.Resource.GetId())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource scopes response
	scopeIds := make([]string, 0, len(grantResponse.Scopes))
	for _, scope := range grantResponse.Scopes {
		scopeIds = append(scopeIds, scope.GetId())
	}

	resourceScopesResponse, d := fetchResourceScopesFromIDs(ctx, r.Client, data.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), scopeIds)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(grantResponse, resourceResponse, resourceScopesResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResourceGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

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

	resource, resourceScopes, d := plan.getResource(ctx, r.Client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource, resourceScopes)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant := plan.expand(*resource, resourceScopes)

	// Run the API call
	var grantResponse *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
		},
		"UpdateApplicationGrant",
		framework.DefaultCustomError,
		nil,
		&grantResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource response
	resourceResponse, d := fetchResourceFromID(ctx, r.Client, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource scopes response
	scopeIds := make([]string, 0, len(grantResponse.Scopes))
	for _, scope := range grantResponse.Scopes {
		scopeIds = append(scopeIds, scope.GetId())
	}

	resourceScopesResponse, d := fetchResourceScopesFromIDs(ctx, r.Client, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), scopeIds)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(grantResponse, resourceResponse, resourceScopesResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourceGrantResourceModel

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

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := r.Client.ApplicationResourceGrantsApi.DeleteApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteApplicationGrant",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResourceGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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
			Label:     "resource_grant_id",
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

func (p *ApplicationResourceGrantResourceModel) getResource(ctx context.Context, apiClient *management.APIClient) (*management.Resource, []management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var d diag.Diagnostics

	var resource *management.Resource
	if !p.ResourceId.IsNull() && !p.ResourceId.IsUnknown() {
		resource, d = fetchResourceFromID(ctx, apiClient, p.EnvironmentId.ValueString(), p.ResourceId.ValueString())
	}

	if !p.ResourceName.IsNull() && !p.ResourceName.IsUnknown() {
		resource, d = fetchResourceFromName(ctx, apiClient, p.EnvironmentId.ValueString(), p.ResourceName.ValueString())
	}

	resourceScopes := make([]management.ResourceScope, 0)
	if resource != nil && !p.Scopes.IsNull() && !p.Scopes.IsUnknown() {

		var scopeIds []string
		diags.Append(p.Scopes.ElementsAs(ctx, &scopeIds, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		resourceScopes, d = fetchResourceScopesFromIDs(ctx, apiClient, p.EnvironmentId.ValueString(), resource.GetId(), scopeIds)
	}

	if resource != nil && !p.ScopeNames.IsNull() && !p.ScopeNames.IsUnknown() {

		var scopeNames []string
		diags.Append(p.ScopeNames.ElementsAs(ctx, &scopeNames, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		resourceScopes, d = fetchResourceScopesFromNames(ctx, apiClient, p.EnvironmentId.ValueString(), resource.GetId(), scopeNames)
	}

	diags.Append(d...)

	if resource == nil {
		diags.AddError(
			"Invalid resource",
			"Cannot create an application resource grant as the resource could not be found.",
		)
	}

	if len(resourceScopes) == 0 {
		diags.AddError(
			"Invalid scopes",
			"Cannot create an application resource grant as the scopes could not be found.",
		)
	}

	if diags.HasError() {
		return nil, nil, diags
	}

	return resource, resourceScopes, diags
}

func (p *ApplicationResourceGrantResourceModel) expand(resource management.Resource, resourceScopes []management.ResourceScope) *management.ApplicationResourceGrant {

	resourceObj := management.NewApplicationResourceGrantResource(resource.GetId())

	scopes := make([]management.ApplicationResourceGrantScopesInner, 0, len(resourceScopes))
	for _, scope := range resourceScopes {
		scopes = append(scopes, management.ApplicationResourceGrantScopesInner{
			Id: scope.GetId(),
		})
	}

	data := management.NewApplicationResourceGrant(*resourceObj, scopes)

	return data
}

func (p *ApplicationResourceGrantResourceModel) validate(resource management.Resource, resourceScopes []management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the `openid` scope from the `openid` resource is not in the list
	if v, ok := resource.GetNameOk(); ok && *v == "openid" && len(resourceScopes) > 0 {
		for _, resourceScope := range resourceScopes {
			if resourceScopeName, ok := resourceScope.GetNameOk(); ok && *resourceScopeName == "openid" {
				diags.AddError(
					"Invalid scope",
					"Cannot create an application resource grant with the `openid` scope.  This scope is automatically applied and should be removed from the `scopes` parameter.",
				)
				break
			}
		}
	}

	return diags
}

func (p *ApplicationResourceGrantResourceModel) toState(apiObject *management.ApplicationResourceGrant, resourceApiObject *management.Resource, resourceScopesApiObjects []management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.ResourceId = framework.StringOkToTF(resourceApiObject.GetIdOk())
	p.ResourceName = framework.StringOkToTF(resourceApiObject.GetNameOk())
	p.ApplicationId = framework.StringOkToTF(apiObject.Application.GetIdOk())

	if _, ok := apiObject.GetScopesOk(); ok {

		scopeIds := make([]string, 0, len(resourceScopesApiObjects))
		scopeNames := make([]string, 0, len(resourceScopesApiObjects))

		for _, scope := range resourceScopesApiObjects {
			scopeIds = append(scopeIds, scope.GetId())
			scopeNames = append(scopeNames, scope.GetName())
		}

		p.Scopes = framework.StringSetToTF(scopeIds)
		p.ScopeNames = framework.StringSetToTF(scopeNames)

	} else {
		p.Scopes = types.SetNull(types.StringType)
		p.ScopeNames = types.SetNull(types.StringType)
	}

	return diags
}
