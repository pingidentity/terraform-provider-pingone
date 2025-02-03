// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type ResourceScopeDataSource serviceClientType

type ResourceScopeDataSourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceType     types.String                 `tfsdk:"resource_type"`
	ResourceId       pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	CustomResourceId pingonetypes.ResourceIDValue `tfsdk:"custom_resource_id"`
	ResourceScopeId  pingonetypes.ResourceIDValue `tfsdk:"resource_scope_id"`
	Name             types.String                 `tfsdk:"name"`
	Description      types.String                 `tfsdk:"description"`
	SchemaAttributes types.Set                    `tfsdk:"schema_attributes"`
	MappedClaims     types.Set                    `tfsdk:"mapped_claims"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &ResourceScopeDataSource{}
)

// New Object
func NewResourceScopeDataSource() datasource.DataSource {
	return &ResourceScopeDataSource{}
}

// Metadata
func (r *ResourceScopeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_scope"
}

// Schema
func (r *ResourceScopeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	resourceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The type of the resource to select.  When the value is set to `%s`, `custom_resource_id` must be specified.", string(management.ENUMRESOURCETYPE_CUSTOM)),
	).AllowedValuesEnum(management.AllowedEnumResourceTypeEnumValues)

	customResourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the ID of the custom resource to select.  Must be a valid PingOne resource ID.  Required if `resource_type` is set to `%s`, but cannot be set if `resource_type` is set to `%s` or `%s`.", string(management.ENUMRESOURCETYPE_CUSTOM), string(management.ENUMRESOURCETYPE_OPENID_CONNECT), string(management.ENUMRESOURCETYPE_PINGONE_API)),
	)

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the resource granted to the application.",
	)

	resourceScopeIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the resource scope.",
	).ExactlyOneOf([]string{"resource_scope_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource scope.",
	).ExactlyOneOf([]string{"resource_scope_id", "name"})

	schemaAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the user schema attributes that can be read or updated for the specified PingOne access control scope. The value is an array of schema attribute paths (such as `username`, `name.given`, `shirtSize`) that the scope controls. This property is supported only for the `p1:read:user`, `p1:update:user` and `p1:read:user:{suffix}` and `p1:update:user:{suffix}` scopes. No other PingOne platform scopes allow this behavior. Any attributes not listed in the attribute array are excluded from the read or update action. The wildcard path (`*`) in the array includes all attributes and cannot be used in conjunction with any other user schema attribute path.",
	)

	mappedClaimsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of custom resource attribute IDs. This property applies only for the resource with its type property set to `OPENID_CONNECT`. Moreover, this property does not display predefined OpenID Connect (OIDC) mappings, such as the `email` claim in the OIDC `email` scope or the `name` claim in the `profile` scope. You can create custom attributes, and these custom attributes can be added to `mapped_claims` and will display in the response.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne OAuth 2.0 resource scope data.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the resource scope."),
			),

			"resource_type": schema.StringAttribute{
				Description:         resourceTypeDescription.Description,
				MarkdownDescription: resourceTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumResourceTypeEnumValues)...),
				},
			},

			"custom_resource_id": schema.StringAttribute{
				Description:         customResourceIdDescription.Description,
				MarkdownDescription: customResourceIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

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

			"resource_scope_id": schema.StringAttribute{
				Description:         resourceScopeIdDescription.Description,
				MarkdownDescription: resourceScopeIdDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("resource_scope_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the resource scope.").Description,
				Computed:    true,
			},

			"schema_attributes": schema.SetAttribute{
				Description:         schemaAttributesDescription.Description,
				MarkdownDescription: schemaAttributesDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,
			},

			"mapped_claims": schema.SetAttribute{
				Description:         mappedClaimsDescription.Description,
				MarkdownDescription: mappedClaimsDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,
			},
		},
	}
}

func (r *ResourceScopeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *ResourceScopeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ResourceScopeDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resource *management.Resource
	var d diag.Diagnostics

	switch data.ResourceType.ValueString() {
	case string(management.ENUMRESOURCETYPE_CUSTOM):
		resource, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.CustomResourceId.ValueString(), false)
		resp.Diagnostics.Append(d...)
	case string(management.ENUMRESOURCETYPE_OPENID_CONNECT), string(management.ENUMRESOURCETYPE_PINGONE_API):
		resource, d = fetchResourceByType(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), management.EnumResourceType(data.ResourceType.ValueString()), false)
		resp.Diagnostics.Append(d...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var resourceScope *management.ResourceScope

	if !data.Name.IsNull() {

		var d diag.Diagnostics
		resourceScope, d = fetchResourceScopeFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), resource.GetId(), data.Name.ValueString(), true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.ResourceScopeId.IsNull() {

		var d diag.Diagnostics
		resourceScope, d = fetchResourceScopeFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), resource.GetId(), data.ResourceScopeId.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested resource scope. resource_scope_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resourceScope, resource)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ResourceScopeDataSourceModel) toState(apiObject *management.ResourceScope, resourceApiObject *management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.ResourceId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.ResourceType = framework.EnumOkToTF(resourceApiObject.GetTypeOk())
	p.ResourceScopeId = framework.PingOneResourceIDToTF(apiObject.GetId())

	p.CustomResourceId = pingonetypes.NewResourceIDNull()
	if resourceApiObject.GetType() == management.ENUMRESOURCETYPE_CUSTOM {
		p.CustomResourceId = framework.PingOneResourceIDOkToTF(apiObject.Resource.GetIdOk())
	}

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.SchemaAttributes = framework.StringSetOkToTF(apiObject.GetSchemaAttributesOk())
	p.MappedClaims = framework.StringSetOkToTF(apiObject.GetMappedClaimsOk())

	return diags
}

func fetchResourceScopeFromID(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceScopeID string) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resourceScope *management.ResourceScope
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ResourceScopesApi.ReadOneResourceScope(ctx, environmentID, resourceID, resourceScopeID).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"ReadOneResourceScope",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resourceScope,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return resourceScope, diags
}

func fetchResourceScopeFromName(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceScopeName string, warnIfNotFound bool) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	resourceScopes, d := fetchResourceScopesFromNames(ctx, apiClient, environmentID, resourceID, []string{resourceScopeName}, warnIfNotFound)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if len(resourceScopes) == 0 {
		if !warnIfNotFound {
			diags.AddError(
				"Cannot find resource scope from ID",
				fmt.Sprintf("The resource scope %s for resource %s in environment %s cannot be found", resourceScopeName, resourceID, environmentID),
			)
		} else {
			diags.AddWarning(
				"Cannot find resource scope from ID",
				fmt.Sprintf("The resource scope %s for resource %s in environment %s cannot be found", resourceScopeName, resourceID, environmentID),
			)
		}
		return nil, diags
	}

	if len(resourceScopes) > 1 {
		diags.AddError(
			"Multiple resource scopes found from ID",
			fmt.Sprintf("Multiple resource scopes %s for resource %s in environment %s were found", resourceScopeName, resourceID, environmentID),
		)
		return nil, diags
	}

	return &resourceScopes[0], diags
}

func fetchResourceScopesFromIDs(ctx context.Context, apiClient *management.APIClient, environmentID string, resourceID string, resourceScopeIDs []string, warnIfNotFound bool) ([]management.ResourceScope, diag.Diagnostics) {
	return fetchResourceScopesFromIDOrNameSlice(ctx, apiClient, environmentID, resourceID, resourceScopeIDs, false, warnIfNotFound)
}

func fetchResourceScopesFromNames(ctx context.Context, apiClient *management.APIClient, environmentID string, resourceID string, resourceScopeNames []string, warnIfNotFound bool) ([]management.ResourceScope, diag.Diagnostics) {
	return fetchResourceScopesFromIDOrNameSlice(ctx, apiClient, environmentID, resourceID, resourceScopeNames, true, warnIfNotFound)
}

func fetchResourceScopesFromIDOrNameSlice(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID string, resourceScopeInputList []string, byName bool, warnIfNotFound bool) ([]management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError
	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	var resourceScopesMap map[string]management.ResourceScope
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.ResourceScopesApi.ReadAllResourceScopes(ctx, environmentID, resourceID).Execute()

			var initialHttpResponse *http.Response

			foundResourceScopesMap := make(map[string]management.ResourceScope)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if resourceScopes, ok := pageCursor.EntityArray.Embedded.GetScopesOk(); ok {

					for _, resourceScope := range resourceScopes {
						if byName {
							foundResourceScopesMap[strings.ToLower(resourceScope.GetName())] = resourceScope
						} else {
							foundResourceScopesMap[strings.ToLower(resourceScope.GetId())] = resourceScope
						}
					}

				}
			}

			return foundResourceScopesMap, initialHttpResponse, nil
		},
		"ReadAllResourceScopes",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&resourceScopesMap,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if resourceScopesMap == nil {
		if warnIfNotFound {
			diags.AddWarning(
				"Cannot find resource scopes",
				fmt.Sprintf("The resource scopes for environment %s cannot be found", environmentID),
			)
		} else {
			diags.AddError(
				"Cannot find resource scopes",
				fmt.Sprintf("The resource scopes for environment %s cannot be found", environmentID),
			)
		}
		return nil, diags
	}

	var foundResourceScopes []management.ResourceScope

	for _, resourceScopeInputItem := range resourceScopeInputList {
		if v, ok := resourceScopesMap[strings.ToLower(resourceScopeInputItem)]; ok {
			foundResourceScopes = append(foundResourceScopes, v)
		} else {
			diags.AddError(
				"Cannot find resource scope",
				fmt.Sprintf("The resource scope %s for resource %s in environment %s cannot be found", resourceScopeInputItem, resourceID, environmentID),
			)
		}
	}

	if diags.HasError() {
		return nil, diags
	}

	return foundResourceScopes, diags
}
