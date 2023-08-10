package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResourceGrantResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ApplicationResourceGrantResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	ResourceId    types.String `tfsdk:"resource_id"`
	Scopes        types.Set    `tfsdk:"scopes"`
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

			"resource_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the protected resource associated with this grant."),
			),

			"scopes": schema.SetAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of IDs of the scopes associated with this grant.  When using the `openid` resource, the `openid` scope should not be included.").Description,
				Required:    true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.ValueStringsAre(
						verify.P1ResourceIDValidator(),
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *ApplicationResourceGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

	if r.client == nil {
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

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(ctx, r.client)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationResourceGrantsApi.CreateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
		},
		"CreateApplicationGrant",
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

func (r *ApplicationResourceGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourceGrantResourceModel

	if r.client == nil {
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
	var response *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplicationGrant",
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

func (r *ApplicationResourceGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourceGrantResourceModel

	if r.client == nil {
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

	// Validate the plan
	resp.Diagnostics.Append(plan.validate(ctx, r.client)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResourceGrant, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationResourceGrant
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationResourceGrant(*applicationResourceGrant).Execute()
		},
		"UpdateApplicationGrant",
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

func (r *ApplicationResourceGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourceGrantResourceModel

	if r.client == nil {
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
			r, err := r.client.ApplicationResourceGrantsApi.DeleteApplicationGrant(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
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
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/application_id/resource_grant_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *ApplicationResourceGrantResourceModel) expand(ctx context.Context) (*management.ApplicationResourceGrant, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := management.NewApplicationResourceGrantResource(p.ResourceId.ValueString())

	var scopesPlan []string
	diags.Append(p.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	scopes := make([]management.ApplicationResourceGrantScopesInner, 0, len(scopesPlan))
	for _, scope := range scopesPlan {
		scopes = append(scopes, management.ApplicationResourceGrantScopesInner{
			Id: scope,
		})
	}

	data := management.NewApplicationResourceGrant(*resource, scopes)

	return data, diags
}

func (p *ApplicationResourceGrantResourceModel) validate(ctx context.Context, apiClient *management.APIClient) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that the `openid` scope from the `openid` resource is not in the list
	var resource *management.Resource
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourcesApi.ReadOneResource(ctx, p.EnvironmentId.ValueString(), p.ResourceId.ValueString()).Execute()
		},
		"ReadOneResource",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resource,
	)...)
	if diags.HasError() {
		return diags
	}

	if v, ok := resource.GetNameOk(); ok && *v == "openid" {
		var entityArray *management.EntityArray
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ResourceScopesApi.ReadAllResourceScopes(ctx, p.EnvironmentId.ValueString(), p.ResourceId.ValueString()).Execute()
			},
			"ReadAllResourceScopes",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)

		if diags.HasError() {
			return diags
		}

		if resourceScopes, ok := entityArray.Embedded.GetScopesOk(); ok {
			openidScope := ""
			for _, resourceScope := range resourceScopes {
				if resourceScopeName, ok := resourceScope.GetNameOk(); ok && *resourceScopeName == "openid" {
					openidScope = resourceScope.GetId()
					break
				}
			}

			if openidScope != "" {
				var scopesPlan []string
				diags.Append(p.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
				if diags.HasError() {
					return diags
				}

				for _, scope := range scopesPlan {
					if scope == openidScope {
						diags.AddError(
							"Invalid scope",
							"Cannot create an application resource grant with the `openid` scope.  This scope is automatically applied and should be removed from the `scopes` parameter.",
						)
						break
					}
				}
			}
		}
	}

	return diags
}

func (p *ApplicationResourceGrantResourceModel) toState(apiObject *management.ApplicationResourceGrant) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.ResourceId = framework.StringOkToTF(apiObject.Resource.GetIdOk())
	p.ApplicationId = framework.StringOkToTF(apiObject.Application.GetIdOk())

	if v, ok := apiObject.GetScopesOk(); ok {
		items := make([]string, 0, len(v))
		for _, scope := range v {
			items = append(items, scope.GetId())
		}
		p.Scopes = framework.StringSetToTF(items)
	} else {
		p.Scopes = types.SetNull(types.StringType)
	}

	return diags
}
