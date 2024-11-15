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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResourceGrantResource serviceClientType

type ApplicationResourceGrantResourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId    pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	ResourceType     types.String                 `tfsdk:"resource_type"`
	ResourceId       pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	CustomResourceId pingonetypes.ResourceIDValue `tfsdk:"custom_resource_id"`
	Scopes           types.Set                    `tfsdk:"scopes"`
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

	resourceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The type of the resource to configure the grant for. When the value is set to `%s`, `custom_resource_id` must be specified.", string(management.ENUMRESOURCETYPE_CUSTOM)),
	).AllowedValuesEnum(management.AllowedEnumResourceTypeEnumValues).RequiresReplace()

	customResourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the ID of the custom resource to be granted to the application.  Must be a valid PingOne resource ID.  Required if `resource_type` is set to `%s`, but cannot be set if `resource_type` is set to `%s` or `%s`.", string(management.ENUMRESOURCETYPE_CUSTOM), string(management.ENUMRESOURCETYPE_OPENID_CONNECT), string(management.ENUMRESOURCETYPE_PINGONE_API)),
	).RequiresReplace()

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the resource granted to the application.",
	)

	scopesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of IDs of the scopes associated with this grant.  Values must be valid PingOne resource IDs.",
	)

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

			"resource_type": schema.StringAttribute{
				Description:         resourceTypeDescription.Description,
				MarkdownDescription: resourceTypeDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumResourceTypeEnumValues)...),
				},
			},

			"custom_resource_id": schema.StringAttribute{
				Description:         customResourceIdDescription.Description,
				MarkdownDescription: customResourceIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(management.ENUMRESOURCETYPE_CUSTOM)),
						path.MatchRelative().AtParent().AtName("resource_type"),
					),
					stringvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(management.ENUMRESOURCETYPE_OPENID_CONNECT)),
						path.MatchRelative().AtParent().AtName("resource_type"),
					),
					stringvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(management.ENUMRESOURCETYPE_PINGONE_API)),
						path.MatchRelative().AtParent().AtName("resource_type"),
					),
				},
			},

			"resource_id": schema.StringAttribute{
				Description:         resourceIdDescription.Description,
				MarkdownDescription: resourceIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"scopes": schema.SetAttribute{
				Description:         scopesDescription.Description,
				MarkdownDescription: scopesDescription.MarkdownDescription,
				Required:            true,

				ElementType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *ApplicationResourceGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

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

	resource, resourceScopes, d := plan.getResourceWithScopes(ctx, r.Client.ManagementAPIClient, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Application
	application, d := plan.getApplication(ctx, r.Client.ManagementAPIClient, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemApplication := false

	if application.ApplicationPingOnePortal != nil || application.ApplicationPingOneSelfService != nil {
		systemApplication = true
	}

	if application.ApplicationPingOneAdminConsole != nil {
		resp.Diagnostics.AddError(
			"Invalid application",
			"Cannot create an application resource grant for the PingOne Admin Console application.",
		)
		return
	}

	// Get the resourceGrant if it exists
	var replaceResourceGrant *management.ApplicationResourceGrant
	if systemApplication {
		replaceResourceGrant, d = plan.getResourceGrant(ctx, r.Client.ManagementAPIClient, resource.GetId())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource, resourceScopes)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant, d := plan.expand(ctx, *resource, replaceResourceGrant)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var grantResponse *management.ApplicationResourceGrant

	if replaceResourceGrant != nil {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), applicationResourceGrant.GetId()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateApplicationGrant-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&grantResponse,
		)...)
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourceGrantsApi.CreateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateApplicationGrant-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&grantResponse,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the resource response
	resourceResponse, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(grantResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourceGrantResourceModel

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
	var grantResponse *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
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
	resourceResponse, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(grantResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResourceGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

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

	resource, resourceScopes, d := plan.getResourceWithScopes(ctx, r.Client.ManagementAPIClient, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Application
	application, d := plan.getApplication(ctx, r.Client.ManagementAPIClient, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemApplication := false

	if application.ApplicationPingOnePortal != nil || application.ApplicationPingOneSelfService != nil {
		systemApplication = true
	}

	if application.ApplicationPingOneAdminConsole != nil {
		resp.Diagnostics.AddError(
			"Invalid application",
			"Cannot create an application resource grant for the PingOne Admin Console application.",
		)
		return
	}

	// Get the resourceGrant if it exists
	var replaceResourceGrant *management.ApplicationResourceGrant
	if systemApplication {
		replaceResourceGrant, d = plan.getResourceGrant(ctx, r.Client.ManagementAPIClient, resource.GetId())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(*resource, resourceScopes)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant, d := plan.expand(ctx, *resource, replaceResourceGrant)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var grantResponse *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	resourceResponse, d := fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), grantResponse.Resource.GetId(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(grantResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourceGrantResourceModel

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
			fR, fErr := r.Client.ManagementAPIClient.ApplicationResourceGrantsApi.DeleteApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
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

func (p *ApplicationResourceGrantResourceModel) getResourceWithScopes(ctx context.Context, apiClient *management.APIClient, warnIfNotFound bool) (*management.Resource, []management.ResourceScope, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	var resource *management.Resource

	switch p.ResourceType.ValueString() {
	case string(management.ENUMRESOURCETYPE_CUSTOM):
		resource, d = fetchResourceFromID(ctx, apiClient, p.EnvironmentId.ValueString(), p.CustomResourceId.ValueString(), warnIfNotFound)
		diags.Append(d...)
	case string(management.ENUMRESOURCETYPE_OPENID_CONNECT), string(management.ENUMRESOURCETYPE_PINGONE_API):
		resource, d = fetchResourceByType(ctx, apiClient, p.EnvironmentId.ValueString(), management.EnumResourceType(p.ResourceType.ValueString()), warnIfNotFound)
		diags.Append(d...)
	}
	if diags.HasError() {
		return nil, nil, diags
	}

	if resource == nil {
		return nil, nil, diags
	}

	resourceScopes := make([]management.ResourceScope, 0)
	if !p.Scopes.IsNull() && !p.Scopes.IsUnknown() {

		var scopesPlan []pingonetypes.ResourceIDValue
		diags.Append(p.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		scopes, d := framework.TFTypePingOneResourceIDSliceToStringSlice(scopesPlan, path.Root("scopes"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}

		resourceScopes, d = fetchResourceScopesFromIDs(ctx, apiClient, p.EnvironmentId.ValueString(), resource.GetId(), scopes, false)
		diags.Append(d...)
	}
	if diags.HasError() {
		return nil, nil, diags
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

func (p *ApplicationResourceGrantResourceModel) getApplication(ctx context.Context, apiClient *management.APIClient, warnIfNotFound bool) (*management.ReadOneApplication200Response, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError
	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	var application *management.ReadOneApplication200Response
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.ReadOneApplication(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, p.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&application,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return application, diags
}

func (p *ApplicationResourceGrantResourceModel) getResourceGrant(ctx context.Context, apiClient *management.APIClient, resourceID string) (*management.ApplicationResourceGrant, diag.Diagnostics) {
	var diags diag.Diagnostics

	var applicationGrant *management.ApplicationResourceGrant
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.ApplicationResourceGrantsApi.ReadAllApplicationGrants(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, p.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				for _, grant := range pageCursor.EntityArray.Embedded.GetGrants() {
					if grant.Resource.GetId() == resourceID {
						return &grant, pageCursor.HTTPResponse, nil
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllApplicationGrants",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&applicationGrant,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return applicationGrant, diags
}

func (p *ApplicationResourceGrantResourceModel) expand(ctx context.Context, resource management.Resource, replaceResourceGrant *management.ApplicationResourceGrant) (*management.ApplicationResourceGrant, diag.Diagnostics) {
	var diags diag.Diagnostics

	resourceObj := management.NewApplicationResourceGrantResource(resource.GetId())

	var scopesPlan []pingonetypes.ResourceIDValue
	diags.Append(p.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	scopesStr, d := framework.TFTypePingOneResourceIDSliceToStringSlice(scopesPlan, path.Root("scopes"))
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	scopes := make([]management.ApplicationResourceGrantScopesInner, 0, len(scopesPlan))
	for _, scope := range scopesStr {
		scopes = append(scopes, management.ApplicationResourceGrantScopesInner{
			Id: scope,
		})
	}

	var data *management.ApplicationResourceGrant

	if replaceResourceGrant != nil {
		data = replaceResourceGrant
		data.SetResource(*resourceObj)
		data.SetScopes(scopes)
	} else {
		data = management.NewApplicationResourceGrant(*resourceObj, scopes)
	}

	return data, diags
}

func (p *ApplicationResourceGrantResourceModel) validate(resource management.Resource, resourceScopes []management.ResourceScope) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the scopes relate to the resource
	for _, resourceScope := range resourceScopes {
		if resourceS, ok := resourceScope.GetResourceOk(); ok && resourceS.GetId() != resource.GetId() {
			diags.AddError(
				"Invalid scope",
				fmt.Sprintf("Cannot create an application resource grant as the scope %s does not relate to the resource %s.", resourceScope.GetId(), resource.GetId()),
			)

		}
	}

	return diags
}

func (p *ApplicationResourceGrantResourceModel) toState(apiObject *management.ApplicationResourceGrant, resourceApiObject *management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || resourceApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.ApplicationId = framework.PingOneResourceIDOkToTF(apiObject.Application.GetIdOk())
	p.ResourceId = framework.PingOneResourceIDOkToTF(resourceApiObject.GetIdOk())
	p.ResourceType = framework.EnumOkToTF(resourceApiObject.GetTypeOk())

	p.Scopes = types.SetNull(types.StringType)
	if scopes, ok := apiObject.GetScopesOk(); ok {
		scopesList := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			scopesList = append(scopesList, scope.GetId())
		}
		p.Scopes = framework.PingOneResourceIDSetToTF(scopesList)
	}

	p.CustomResourceId = pingonetypes.NewResourceIDNull()
	if resourceApiObject.GetType() == management.ENUMRESOURCETYPE_CUSTOM {
		p.CustomResourceId = framework.PingOneResourceIDOkToTF(apiObject.Resource.GetIdOk())
	}

	return diags
}
