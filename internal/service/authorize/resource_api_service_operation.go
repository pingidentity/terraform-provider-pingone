package authorize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type APIServiceOperationResource serviceClientType

type APIServiceOperationResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	APIServiceId  pingonetypes.ResourceIDValue `tfsdk:"api_service_id"`
	AccessControl types.Object                 `tfsdk:"access_control"`
	Methods       types.Set                    `tfsdk:"methods"`
	Paths         types.Set                    `tfsdk:"paths"`
	Name          types.String                 `tfsdk:"name"`
	PolicyId      pingonetypes.ResourceIDValue `tfsdk:"policy_id"`
}

type APIServiceOperationAccessControlResourceModel struct {
	Group      types.Object `tfsdk:"group"`
	Permission types.Object `tfsdk:"permission"`
	Scope      types.Object `tfsdk:"scope"`
}

type APIServiceOperationAccessControlGroupResourceModel struct {
	Groups types.Set `tfsdk:"groups"`
}

type APIServiceOperationAccessControlGroupGroupsResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceOperationAccessControlPermissionResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceOperationAccessControlScopeResourceModel struct {
	MatchType types.String `tfsdk:"match_type"`
	Scopes    types.Set    `tfsdk:"scopes"`
}

type APIServiceOperationAccessControlScopeScopesResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type APIServiceOperationPathResourceModel struct {
	Pattern types.String `tfsdk:"pattern"`
	Type    types.String `tfsdk:"type"`
}

var (
	apiServiceOperationAccessControlTFObjectTypes = map[string]attr.Type{
		"group":      types.ObjectType{AttrTypes: apiServiceOperationAccessControlGroupTFObjectTypes},
		"permission": types.ObjectType{AttrTypes: apiServiceOperationAccessControlPermissionTFObjectTypes},
		"scope":      types.ObjectType{AttrTypes: apiServiceOperationAccessControlScopeTFObjectTypes},
	}

	apiServiceOperationAccessControlGroupTFObjectTypes = map[string]attr.Type{
		"groups": types.SetType{ElemType: types.ObjectType{AttrTypes: apiServiceOperationAccessControlGroupGroupsTFObjectTypes}},
	}

	apiServiceOperationAccessControlGroupGroupsTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceOperationAccessControlPermissionTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceOperationAccessControlScopeTFObjectTypes = map[string]attr.Type{
		"match_type": types.StringType,
		"scopes":     types.SetType{ElemType: types.ObjectType{AttrTypes: apiServiceOperationAccessControlScopeScopesTFObjectTypes}},
	}

	apiServiceOperationAccessControlScopeScopesTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	apiServiceOperationPathTFObjectTypes = map[string]attr.Type{
		"pattern": types.StringType,
		"type":    types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &APIServiceOperationResource{}
	_ resource.ResourceWithConfigure   = &APIServiceOperationResource{}
	_ resource.ResourceWithImportState = &APIServiceOperationResource{}
)

// New Object
func NewAPIServiceOperationResource() resource.Resource {
	return &APIServiceOperationResource{}
}

// Metadata
func (r *APIServiceOperationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_api_service_operation"
}

// Schema.
func (r *APIServiceOperationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const attrGroupGroupsMaxLength = 25
	const attrMethodsMaxLength = 10
	const attrPathsMaxLength = 10
	const attrPathsPatternMaxLength = 2048

	accessControlScopeMatchTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the match type of the scope rule.",
	).AllowedValuesComplex(map[string]string{
		string(authorize.ENUMAPISERVEROPERATIONMATCHTYPE_ALL): "the client must be authorized with all scopes configured in the `scopes` array to obtain access",
		string(authorize.ENUMAPISERVEROPERATIONMATCHTYPE_ANY): "the client must be authorized with one or more of the scopes configured in the `scopes` array to obtain access",
	})

	accessControlScopeScopesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specify the scopes that define the access requirements for the operation. The client must be authorized with `ANY` or `ALL` of the scopes to be granted access, depending on the `match_type` field value.",
	)

	methodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The methods that define the operation. No duplicates are allowed. Each element must be a valid HTTP token, according to [RFC 7230](https://datatracker.ietf.org/doc/html/rfc7230), and cannot exceed 64 characters. An empty array is not valid. To indicate that an operation is defined for every method, the `methods` array should be set to null. The `methods` array is limited to 10 entries.",
	).AllowedValuesEnum(authorize.AllowedEnumAPIServerOperationMethodEnumValues)

	pathsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies the paths that define the operation. The same literal pattern is not allowed within the same operation (the pattern of a `paths` element must be unique as compared to all other patterns in the same `paths` array). However, the same literal pattern is allowed in different operations (for example, OperationA, `/path1`, OperationB, `/path1` is valid). This set is limited to 10 entries.",
	)

	pathsPatternDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the pattern used to identify the path or paths for the operation. The semantics of the pattern are determined by the type. For any type, the pattern can contain characters that are otherwise invalid in a URL path. Invalid characters are handled by performing matching against a percent-decoded HTTP request target path. This allows an administrator to configure patterns without worrying about percent encoding special characters.\n" +
			"When the `type` is `PARAMETER`, the syntax outlined in the table below is enforced. Additionally, the pattern must contain a wildcard, double wildcard or named parameter. When the `type` is `EXACT`, the pattern can be any byte sequence except for ASCII control characters such as line feeds or carriage returns. The length of the pattern cannot exceed 2048 characters. The path pattern must not contain empty path segments such as `/../`, `//`, and `/./`.",
	)

	pathsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the pattern.",
	).AllowedValuesComplex(map[string]string{
		string(authorize.ENUMAPISERVEROPERATIONPATHPATTERNTYPE_EXACT):     "the verbatim pattern is compared against the path from the request using a case-sensitive comparison",
		string(authorize.ENUMAPISERVEROPERATIONPATHPATTERNTYPE_PARAMETER): "the pattern is compared against the path from the request using a case-sensitive comparison, using the syntax below to encode wildcards and named parameters",
	})

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage API Operations for API Services with PingOne Authorize in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create and manage the API Service in."),
			),

			"api_service_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the API service to create and manage the API Service operation for."),
			),

			"access_control": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies properties related to access control settings of the API service operation.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"group": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that defines the group membership requirements for the operation.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"groups": schema.SetNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of objects that define the access requirements for the operation. The end user must be a member of one or more of these groups to gain access to the operation. The ID must reference a group that exists at the time the data is persisted. There is no referential integrity between a group and this configuration. If a group is subsequently deleted, the access control configuration will continue to reference that group. The set must not contain more than 25 elements.").Description,
								Required:    true,

								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the UUID that represents the ID of the PingOne group. Must be a valid PingOne resource ID.").Description,
											Required:    true,

											CustomType: pingonetypes.ResourceIDType{},
										},
									},
								},

								Validators: []validator.Set{
									setvalidator.SizeAtMost(attrGroupGroupsMaxLength),
								},
							},
						},
					},

					"permission": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that defines permission requirements for the operation.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application permission ID that defines the access requirements for the operation. The end user must be entitled to the specified application permission to gain access to the operation.  Must be a valid PingOne resource ID.").Description,
								Required:    true,

								CustomType: pingonetypes.ResourceIDType{},
							},
						},
					},

					"scope": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that defines scope membership requirements for the operation.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"match_type": schema.StringAttribute{
								Description:         accessControlScopeMatchTypeDescription.Description,
								MarkdownDescription: accessControlScopeMatchTypeDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAPIServerOperationMatchTypeEnumValues)...),
								},
							},

							"scopes": schema.SetNestedAttribute{
								Description:         accessControlScopeScopesDescription.Description,
								MarkdownDescription: accessControlScopeScopesDescription.MarkdownDescription,
								Required:            true,

								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the scope.  Must be a valid PingOne resource ID.").Description,
											Required:    true,

											CustomType: pingonetypes.ResourceIDType{},
										},
									},
								},
							},
						},
					},
				},
			},

			"methods": schema.SetAttribute{
				Description:         methodsDescription.Description,
				MarkdownDescription: methodsDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.SizeAtMost(attrMethodsMaxLength),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAPIServerOperationMethodEnumValues)...),
					),
				},
			},

			"paths": schema.SetNestedAttribute{
				Description:         pathsDescription.Description,
				MarkdownDescription: pathsDescription.MarkdownDescription,
				Required:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pattern": schema.StringAttribute{
							Description:         pathsPatternDescription.Description,
							MarkdownDescription: pathsPatternDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.LengthAtMost(attrPathsPatternMaxLength),
							},
						},

						"type": schema.StringAttribute{
							Description:         pathsTypeDescription.Description,
							MarkdownDescription: pathsTypeDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAPIServerOperationPathPatternTypeEnumValues)...),
							},
						},
					},
				},

				Validators: []validator.Set{
					setvalidator.SizeAtMost(attrPathsMaxLength),
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the API service operation name.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the ID of the root policy.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (r *APIServiceOperationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APIServiceOperationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state APIServiceOperationResourceModel

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

	// Build the model for the API
	apiServiceOperation, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.APIServerOperation
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerOperationsApi.CreateAPIServerOperation(ctx, plan.EnvironmentId.ValueString(), plan.APIServiceId.ValueString()).APIServerOperation(*apiServiceOperation).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAPIServerOperation",
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

func (r *APIServiceOperationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *APIServiceOperationResourceModel

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
	var response *authorize.APIServerOperation
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerOperationsApi.ReadOneAPIServerOperation(ctx, data.EnvironmentId.ValueString(), data.APIServiceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAPIServerOperation",
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

func (r *APIServiceOperationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state APIServiceOperationResourceModel

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

	// Build the model for the API
	apiServiceOperation, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.APIServerOperation
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServerOperationsApi.UpdateAPIServerOperation(ctx, plan.EnvironmentId.ValueString(), plan.APIServiceId.ValueString(), plan.Id.ValueString()).APIServerOperation(*apiServiceOperation).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAPIServerOperation",
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

func (r *APIServiceOperationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *APIServiceOperationResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.AuthorizeAPIClient.APIServerOperationsApi.DeleteAPIServerOperation(ctx, data.EnvironmentId.ValueString(), data.APIServiceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAPIServerOperation",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *APIServiceOperationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "api_service_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "api_service_operation_id",
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

func (p *APIServiceOperationResourceModel) expand(ctx context.Context) (*authorize.APIServerOperation, diag.Diagnostics) {
	var diags diag.Diagnostics

	var pathsPlan []APIServiceOperationPathResourceModel
	diags.Append(p.Paths.ElementsAs(ctx, &pathsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	paths := make([]authorize.APIServerOperationPathsInner, 0)

	for _, pathPlan := range pathsPlan {
		paths = append(paths, *pathPlan.expand())
	}

	data := authorize.NewAPIServerOperation(
		p.Name.ValueString(),
		paths,
	)

	if !p.AccessControl.IsNull() && !p.AccessControl.IsUnknown() {
		var accessControlPlan APIServiceOperationAccessControlResourceModel
		diags.Append(p.AccessControl.As(ctx, &accessControlPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		accessControl, d := accessControlPlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetAccessControl(*accessControl)
	}

	if !p.Methods.IsNull() && !p.Methods.IsUnknown() {
		var methodsPlan []string
		diags.Append(p.Methods.ElementsAs(ctx, &methodsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		methods := make([]authorize.EnumAPIServerOperationMethod, 0)

		for _, methodPlan := range methodsPlan {
			methods = append(methods, authorize.EnumAPIServerOperationMethod(methodPlan))
		}

		data.SetMethods(methods)
	}

	return data, diags
}

func (p *APIServiceOperationPathResourceModel) expand() *authorize.APIServerOperationPathsInner {

	data := authorize.NewAPIServerOperationPathsInner(
		p.Pattern.ValueString(),
		authorize.EnumAPIServerOperationPathPatternType(p.Type.ValueString()),
	)

	return data
}

func (p *APIServiceOperationAccessControlResourceModel) expand(ctx context.Context) (*authorize.APIServerOperationAccessControl, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAPIServerOperationAccessControl()

	if !p.Group.IsNull() && !p.Group.IsUnknown() {
		var groupPlan APIServiceOperationAccessControlGroupResourceModel
		diags.Append(p.Group.As(ctx, &groupPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		group, d := groupPlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetGroup(*group)
	}

	if !p.Permission.IsNull() && !p.Permission.IsUnknown() {
		var permissionPlan APIServiceOperationAccessControlPermissionResourceModel
		diags.Append(p.Permission.As(ctx, &permissionPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPermission(*authorize.NewAPIServerOperationAccessControlPermission(permissionPlan.Id.ValueString()))
	}

	if !p.Scope.IsNull() && !p.Scope.IsUnknown() {
		var scopePlan APIServiceOperationAccessControlScopeResourceModel
		diags.Append(p.Scope.As(ctx, &scopePlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		scope, d := scopePlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetScope(*scope)
	}

	return data, diags
}

func (p *APIServiceOperationAccessControlGroupResourceModel) expand(ctx context.Context) (*authorize.APIServerOperationAccessControlGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupsPlan []APIServiceOperationAccessControlGroupGroupsResourceModel
	diags.Append(p.Groups.ElementsAs(ctx, &groupsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	paths := make([]authorize.APIServerOperationAccessControlGroupGroupsInner, 0)

	for _, groupPlan := range groupsPlan {
		paths = append(paths, *authorize.NewAPIServerOperationAccessControlGroupGroupsInner(groupPlan.Id.ValueString()))
	}

	data := authorize.NewAPIServerOperationAccessControlGroup(
		paths,
	)
	return data, diags
}

func (p *APIServiceOperationAccessControlScopeResourceModel) expand(ctx context.Context) (*authorize.APIServerOperationAccessControlScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var scopesPlan []APIServiceOperationAccessControlScopeScopesResourceModel
	diags.Append(p.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	scopes := make([]authorize.APIServerOperationAccessControlScopeScopesInner, 0)

	for _, scopePlan := range scopesPlan {
		scopes = append(scopes, *authorize.NewAPIServerOperationAccessControlScopeScopesInner(scopePlan.Id.ValueString()))
	}

	data := authorize.NewAPIServerOperationAccessControlScope(
		scopes,
	)

	if !p.MatchType.IsNull() && !p.MatchType.IsUnknown() {
		data.SetMatchType(authorize.EnumAPIServerOperationMatchType(p.MatchType.ValueString()))
	}

	return data, diags
}

func (p *APIServiceOperationResourceModel) toState(apiObject *authorize.APIServerOperation) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())

	p.AccessControl, d = apiServiceOperationAccessControlOkToTF(apiObject.GetAccessControlOk())
	diags.Append(d...)

	p.Methods = framework.EnumSetOkToTF(apiObject.GetMethodsOk())

	p.Paths, d = apiServiceOperationPathsOkToTF(apiObject.GetPathsOk())
	diags.Append(d...)

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	if v, ok := apiObject.GetPolicyOk(); ok {
		p.PolicyId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.PolicyId = pingonetypes.NewResourceIDNull()
	}

	return diags
}

func apiServiceOperationAccessControlOkToTF(apiObject *authorize.APIServerOperationAccessControl, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceOperationAccessControlTFObjectTypes), diags
	}

	group, d := apiServiceOperationAccessControlGroupOkToTF(apiObject.GetGroupOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceOperationAccessControlTFObjectTypes), diags
	}

	permission, d := apiServiceOperationAccessControlPermissionOkToTF(apiObject.GetPermissionOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceOperationAccessControlTFObjectTypes), diags
	}

	scope, d := apiServiceOperationAccessControlScopeOkToTF(apiObject.GetScopeOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceOperationAccessControlTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceOperationAccessControlTFObjectTypes, map[string]attr.Value{
		"group":      group,
		"permission": permission,
		"scope":      scope,
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceOperationAccessControlGroupOkToTF(apiObject *authorize.APIServerOperationAccessControlGroup, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceOperationAccessControlGroupTFObjectTypes), diags
	}

	groups, d := apiServiceOperationAccessControlGroupGroupsOkToTF(apiObject.GetGroupsOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceOperationAccessControlGroupTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceOperationAccessControlGroupTFObjectTypes, map[string]attr.Value{
		"groups": groups,
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceOperationAccessControlScopeOkToTF(apiObject *authorize.APIServerOperationAccessControlScope, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceOperationAccessControlScopeTFObjectTypes), diags
	}

	scopes, d := apiServiceOperationAccessControlScopeScopesOkToTF(apiObject.GetScopesOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceOperationAccessControlScopeTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceOperationAccessControlScopeTFObjectTypes, map[string]attr.Value{
		"match_type": framework.EnumOkToTF(apiObject.GetMatchTypeOk()),
		"scopes":     scopes,
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceOperationAccessControlScopeScopesOkToTF(apiObject []authorize.APIServerOperationAccessControlScopeScopesInner, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: apiServiceOperationAccessControlScopeScopesTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := types.ObjectValue(apiServiceOperationAccessControlScopeScopesTFObjectTypes, map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(v.GetIdOk()),
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func apiServiceOperationAccessControlPermissionOkToTF(apiObject *authorize.APIServerOperationAccessControlPermission, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceOperationAccessControlPermissionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceOperationAccessControlPermissionTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceOperationAccessControlGroupGroupsOkToTF(apiObject []authorize.APIServerOperationAccessControlGroupGroupsInner, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: apiServiceOperationAccessControlGroupGroupsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := types.ObjectValue(apiServiceOperationAccessControlGroupGroupsTFObjectTypes, map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(v.GetIdOk()),
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func apiServiceOperationPathsOkToTF(apiObject []authorize.APIServerOperationPathsInner, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: apiServiceOperationPathTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := types.ObjectValue(apiServiceOperationPathTFObjectTypes, map[string]attr.Value{
			"pattern": framework.StringOkToTF(v.GetPatternOk()),
			"type":    framework.EnumOkToTF(v.GetTypeOk()),
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}
