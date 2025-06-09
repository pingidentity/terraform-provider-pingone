// Copyright © 2025 Ping Identity Corporation

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type APIServiceResource serviceClientType

type APIServiceResourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AccessControl       types.Object                 `tfsdk:"access_control"`
	AuthorizationServer types.Object                 `tfsdk:"authorization_server"`
	BaseURLs            types.Set                    `tfsdk:"base_urls"`
	Directory           types.Object                 `tfsdk:"directory"`
	Name                types.String                 `tfsdk:"name"`
	PolicyId            pingonetypes.ResourceIDValue `tfsdk:"policy_id"`
}

type APIServiceAccessControlResourceModel struct {
	Custom types.Object `tfsdk:"custom"`
}

type APIServiceAccessControlCustomResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type APIServiceAuthorizationServerResourceModel struct {
	ResourceId pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	Type       types.String                 `tfsdk:"type"`
}

type APIServiceDirectoryResourceModel struct {
	Type types.String `tfsdk:"type"`
}

var (
	apiServiceAccessControlTFObjectTypes = map[string]attr.Type{
		"custom": types.ObjectType{AttrTypes: apiServiceAccessControlCustomTFObjectTypes},
	}

	apiServiceAccessControlCustomTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	apiServiceAuthorizationServerTFObjectTypes = map[string]attr.Type{
		"resource_id": pingonetypes.ResourceIDType{},
		"type":        types.StringType,
	}

	apiServiceDirectoryTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                   = &APIServiceResource{}
	_ resource.ResourceWithConfigure      = &APIServiceResource{}
	_ resource.ResourceWithValidateConfig = &APIServiceResource{}
	_ resource.ResourceWithImportState    = &APIServiceResource{}
)

// New Object
func NewAPIServiceResource() resource.Resource {
	return &APIServiceResource{}
}

// Metadata
func (r *APIServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_api_service"
}

// Schema.
func (r *APIServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const attrBaseUrlsMaxLength = 256

	accessControlCustomEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, if set to `true`, means the custom policy will be used for the endpoint.",
	).DefaultValue(false).RequiresReplace()

	authorizationServerResourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the UUID of the custom PingOne resource. The resource defines the characteristics of the OAuth 2.0 access tokens used to get access to the APIs on the API service such as the audience and scopes. This property must identify a PingOne resource with a `type` property value of `CUSTOM`.",
	)

	authorizationServerTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the type of authorization server that will issue access tokens. Must be the same value as the `directory.type` field. If `%s`, the `resource` field must not be provided.", string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_EXTERNAL)),
	).AllowedValuesEnum(authorize.AllowedEnumAPIServerAuthorizationServerTypeEnumValues).DefaultValue(string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_PINGONE_SSO)).RequiresReplace()

	baseUrlsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the possible base URLs that an end-user will use to access the APIs hosted on the customer's API service. Multiple base URLs may be specified to support cases where the same API may be available from multiple URLs (for example, from a user-friendly domain URL and an internal domain URL). Base URLs must be valid absolute URLs with the `https` or `http` scheme. If the path component is non-empty, it must not end in a trailing slash. The path must not contain empty backslash, dot, or double-dot segments. It must not have a query or fragment present, and the host portion of the authority must be a DNS hostname or valid IP (IPv4 or IPv6). The length must be less than or equal to 256 characters.",
	)

	directoryDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A container object for fields related to the user directory used to issue access tokens for accessing the APIs. If not provided, `directory.type` will default to `PINGONE_SSO`.",
	)

	directoryTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of directory that will be used to issue access tokens.",
	).AllowedValuesEnum(authorize.AllowedEnumAPIServerAuthorizationServerTypeEnumValues).DefaultValue(string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_PINGONE_SSO))

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage API Services for PingOne Authorize in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create and manage the API Service in."),
			),

			"access_control": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies properties related to access control settings of the API service.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					apiServiceAccessControlTFObjectTypes,
					map[string]attr.Value{
						"custom": types.ObjectValueMust(
							apiServiceAccessControlCustomTFObjectTypes,
							map[string]attr.Value{
								"enabled": types.BoolValue(false),
							},
						),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"custom": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that defines if the operation will use custom policy rather than the \"Group\" or \"Scope\" access control requirement.").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         accessControlCustomEnabledDescription.Description,
								MarkdownDescription: accessControlCustomEnabledDescription.MarkdownDescription,
								Required:            true,

								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.RequiresReplace(),
								},
							},
						},
					},
				},
			},

			"authorization_server": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies properties related to the authorization server that will issue access tokens used to access the APIs.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"resource_id": schema.StringAttribute{
						Description:         authorizationServerResourceIdDescription.Description,
						MarkdownDescription: authorizationServerResourceIdDescription.MarkdownDescription,
						Optional:            true,

						CustomType: pingonetypes.ResourceIDType{},

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_PINGONE_SSO)),
								path.MatchRelative().AtParent().AtName("type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_EXTERNAL)),
								path.MatchRelative().AtParent().AtName("type"),
							),
						},
					},

					"type": schema.StringAttribute{
						Description:         authorizationServerTypeDescription.Description,
						MarkdownDescription: authorizationServerTypeDescription.MarkdownDescription,
						Required:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAPIServerAuthorizationServerTypeEnumValues)...),
						},
					},
				},
			},

			"base_urls": schema.SetAttribute{
				Description:         baseUrlsDescription.Description,
				MarkdownDescription: baseUrlsDescription.MarkdownDescription,
				Required:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(attrBaseUrlsMaxLength),
					),
				},
			},

			"directory": schema.SingleNestedAttribute{
				Description:         directoryDescription.Description,
				MarkdownDescription: directoryDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					apiServiceDirectoryTFObjectTypes,
					map[string]attr.Value{
						"type": types.StringValue(string(authorize.ENUMAPISERVERAUTHORIZATIONSERVERTYPE_PINGONE_SSO)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         directoryTypeDescription.Description,
						MarkdownDescription: directoryTypeDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAPIServerAuthorizationServerTypeEnumValues)...),
						},
					},
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the API service resource name. The name value must be unique among all API services, and it must be a valid resource name.").Description,
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

func (r *APIServiceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data APIServiceResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	resp.Diagnostics.Append(data.validateAPIServiceAuthzServerType(ctx, true)...)
}

func (r *APIServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APIServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state APIServiceResourceModel

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

	resp.Diagnostics.Append(plan.validateAPIServiceAuthzServerType(ctx, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	apiService, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.APIServer
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServersApi.CreateAPIServer(ctx, plan.EnvironmentId.ValueString()).APIServer(*apiService).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAPIServer",
		legacysdk.DefaultCustomError,
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

func (r *APIServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *APIServiceResourceModel

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
	var response *authorize.APIServer
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServersApi.ReadOneAPIServer(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAPIServer",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *APIServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state APIServiceResourceModel

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

	resp.Diagnostics.Append(plan.validateAPIServiceAuthzServerType(ctx, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	apiService, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.APIServer
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.APIServersApi.UpdateAPIServer(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).APIServer(*apiService).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAPIServer",
		legacysdk.DefaultCustomError,
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

func (r *APIServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *APIServiceResourceModel

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.AuthorizeAPIClient.APIServersApi.DeleteAPIServer(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAPIServer",
		legacysdk.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *APIServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "api_service_id",
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

func (p *APIServiceResourceModel) validateAPIServiceAuthzServerType(ctx context.Context, allowUnknown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if !allowUnknown && p.AuthorizationServer.IsUnknown() {
		diags.AddAttributeError(
			path.Root("authorization_server"),
			"Parameter should be declared",
			"The `authorization_server` parameter is unknown at the time of validation but must be declared.",
		)
	}

	if !allowUnknown && p.Directory.IsUnknown() {
		diags.AddAttributeError(
			path.Root("directory"),
			"Parameter should be declared",
			"The `directory` parameter is unknown at the time of validation but must be declared.",
		)
	}

	if !p.AuthorizationServer.IsNull() && !p.Directory.IsNull() {

		var authzServerPlan APIServiceAuthorizationServerResourceModel
		diags.Append(p.AuthorizationServer.As(ctx, &authzServerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		if !allowUnknown && authzServerPlan.Type.IsUnknown() {
			diags.AddAttributeError(
				path.Root("authorization_server").AtName("type"),
				"Parameter should be declared",
				"The `authorization_server.type` parameter is unknown at the time of validation but must be declared.",
			)
		}

		var directoryPlan APIServiceDirectoryResourceModel
		diags.Append(p.Directory.As(ctx, &directoryPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		if !allowUnknown && directoryPlan.Type.IsUnknown() {
			diags.AddAttributeError(
				path.Root("directory").AtName("type"),
				"Parameter should be declared",
				"The `directory.type` parameter is unknown at the time of validation but must be declared.",
			)
		}

		if !authzServerPlan.Type.IsNull() && !directoryPlan.Type.IsNull() && !authzServerPlan.Type.Equal(directoryPlan.Type) {
			diags.AddAttributeError(
				path.Root("authorization_server").AtName("type"),
				"Parameter conflict",
				fmt.Sprintf("The `authorization_server.type` (value `%s`) and `directory.type` (value `%s`) parameters must be set to the same value.", authzServerPlan.Type.ValueString(), directoryPlan.Type.ValueString()),
			)
			diags.AddAttributeError(
				path.Root("directory").AtName("type"),
				"Parameter conflict",
				fmt.Sprintf("The `authorization_server.type` (value `%s`) and `directory.type` (value `%s`) parameters must be set to the same value.", authzServerPlan.Type.ValueString(), directoryPlan.Type.ValueString()),
			)
		}
	}

	return diags

}

func (p *APIServiceResourceModel) expand(ctx context.Context) (*authorize.APIServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	var authzServerPlan APIServiceAuthorizationServerResourceModel
	diags.Append(p.AuthorizationServer.As(ctx, &authzServerPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	var baseUrlsPlan []types.String
	diags.Append(p.BaseURLs.ElementsAs(ctx, &baseUrlsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	baseUrls, d := framework.TFTypeStringSliceToStringSlice(baseUrlsPlan, path.Root("base_urls"))
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	authorizationServer := authzServerPlan.expand()

	data := authorize.NewAPIServer(
		*authorizationServer,
		baseUrls,
		p.Name.ValueString(),
	)

	if !p.AccessControl.IsNull() && !p.AccessControl.IsUnknown() {
		var accessControlPlan APIServiceAccessControlResourceModel
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

	if !p.Directory.IsNull() && !p.Directory.IsUnknown() {
		var directoryPlan APIServiceDirectoryResourceModel
		diags.Append(p.Directory.As(ctx, &directoryPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		directory := directoryPlan.expand()

		data.SetDirectory(*directory)
	}

	return data, diags
}

func (p *APIServiceAuthorizationServerResourceModel) expand() *authorize.APIServerAuthorizationServer {

	data := authorize.NewAPIServerAuthorizationServer(authorize.EnumAPIServerAuthorizationServerType(p.Type.ValueString()))

	if !p.ResourceId.IsNull() && !p.ResourceId.IsUnknown() {
		data.SetResource(*authorize.NewAPIServerAuthorizationServerResource(p.ResourceId.ValueString()))
	}

	return data
}

func (p *APIServiceAccessControlResourceModel) expand(ctx context.Context) (*authorize.APIServerAccessControl, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAPIServerAccessControl()

	if !p.Custom.IsNull() && !p.Custom.IsUnknown() {
		var customPlan APIServiceAccessControlCustomResourceModel
		diags.Append(p.Custom.As(ctx, &customPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		custom := customPlan.expand()

		data.SetCustom(*custom)
	}

	return data, diags
}

func (p *APIServiceAccessControlCustomResourceModel) expand() *authorize.APIServerAccessControlCustom {
	data := authorize.NewAPIServerAccessControlCustom()

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	return data
}

func (p *APIServiceDirectoryResourceModel) expand() *authorize.APIServerDirectory {

	data := authorize.NewAPIServerDirectory(authorize.EnumAPIServerAuthorizationServerType(p.Type.ValueString()))

	return data
}

func (p *APIServiceResourceModel) toState(apiObject *authorize.APIServer) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())

	p.AccessControl, d = apiServiceAccessControlOkToTF(apiObject.GetAccessControlOk())
	diags.Append(d...)

	p.AuthorizationServer, d = apiServiceAuthorizationServerOkToTF(apiObject.GetAuthorizationServerOk())
	diags.Append(d...)

	p.BaseURLs = framework.StringSetOkToTF(apiObject.GetBaseUrlsOk())

	p.Directory, d = apiServiceDirectoryOkToTF(apiObject.GetDirectoryOk())
	diags.Append(d...)

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	if v, ok := apiObject.GetPolicyOk(); ok {
		p.PolicyId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.PolicyId = pingonetypes.NewResourceIDNull()
	}

	return diags
}

func apiServiceAccessControlOkToTF(apiObject *authorize.APIServerAccessControl, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceAccessControlTFObjectTypes), diags
	}

	custom, d := apiServiceAccessControlCustomOkToTF(apiObject.GetCustomOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(apiServiceAccessControlTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceAccessControlTFObjectTypes, map[string]attr.Value{
		"custom": custom,
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceAccessControlCustomOkToTF(apiObject *authorize.APIServerAccessControlCustom, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceAccessControlCustomTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceAccessControlCustomTFObjectTypes, map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceAuthorizationServerOkToTF(apiObject *authorize.APIServerAuthorizationServer, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceAuthorizationServerTFObjectTypes), diags
	}

	resourceId := pingonetypes.NewResourceIDNull()

	if v, ok := apiObject.GetResourceOk(); ok {
		resourceId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	objValue, d := types.ObjectValue(apiServiceAuthorizationServerTFObjectTypes, map[string]attr.Value{
		"resource_id": resourceId,
		"type":        framework.EnumOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func apiServiceDirectoryOkToTF(apiObject *authorize.APIServerDirectory, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(apiServiceDirectoryTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(apiServiceDirectoryTFObjectTypes, map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
