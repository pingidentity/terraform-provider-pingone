package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ResourceResource serviceClientType

type ResourceResourceModel struct {
	Id                             pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                  pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                           types.String                 `tfsdk:"name"`
	Description                    types.String                 `tfsdk:"description"`
	Type                           types.String                 `tfsdk:"type"`
	Audience                       types.String                 `tfsdk:"audience"`
	AccessTokenValiditySeconds     types.Int64                  `tfsdk:"access_token_validity_seconds"`
	ApplicationPermissionsSettings types.Object                 `tfsdk:"application_permissions_settings"`
	IntrospectEndpointAuthMethod   types.String                 `tfsdk:"introspect_endpoint_auth_method"`
}

type ResourceApplicationPermissionsSettingsModel struct {
	ClaimEnabled types.Bool `tfsdk:"claim_enabled"`
}

var (
	resourceApplicationPermissionsSettingsTFObjectTypes = map[string]attr.Type{
		"claim_enabled": types.BoolType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &ResourceResource{}
	_ resource.ResourceWithConfigure   = &ResourceResource{}
	_ resource.ResourceWithModifyPlan  = &ResourceResource{}
	_ resource.ResourceWithImportState = &ResourceResource{}
)

// New Object
func NewResourceResource() resource.Resource {
	return &ResourceResource{}
}

// Metadata
func (r *ResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema.
func (r *ResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const accessTokenValiditySecondsDefault = 3600
	const accessTokenValiditySecondsMin = 300
	const accessTokenValiditySecondsMinMinutes = accessTokenValiditySecondsMin / 60
	const accessTokenValiditySecondsMax = 2592000
	const accessTokenValiditySecondsMaxDays = accessTokenValiditySecondsMax / (60 * 60 * 24)

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of resource.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMRESOURCETYPE_OPENID_CONNECT): "specifies the built-in platform resource for OpenID Connect",
		string(management.ENUMRESOURCETYPE_PINGONE_API):    "specifies the built-in platform resource for PingOne",
		string(management.ENUMRESOURCETYPE_CUSTOM):         "specifies the a resource that has been created by admin",
	}).AppendMarkdownString(fmt.Sprintf("Only the `%s` resource type can be created. `%s` specifies the built-in platform resource for OpenID Connect. `%s` specifies the built-in platform resource for PingOne.", string(management.ENUMRESOURCETYPE_CUSTOM), string(management.ENUMRESOURCETYPE_OPENID_CONNECT), string(management.ENUMRESOURCETYPE_PINGONE_API)))

	audienceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.",
	)

	accessTokenValiditySecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the number of seconds that the access token is valid.  The minimum value is `%d` seconds (%d minutes); the maximum value is `%d` seconds (%d days).", accessTokenValiditySecondsMin, accessTokenValiditySecondsMinMinutes, accessTokenValiditySecondsMax, accessTokenValiditySecondsMaxDays),
	).DefaultValue(accessTokenValiditySecondsDefault)

	introspectEndpointAuthMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the client authentication methods supported by the token endpoint",
	).AllowedValuesEnum(management.AllowedEnumResourceIntrospectEndpointAuthMethodEnumValues).DefaultValue(string(management.ENUMRESOURCEINTROSPECTENDPOINTAUTHMETHOD_CLIENT_SECRET_BASIC))

	applicationPermissionsSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies whether application permissions are added to access tokens generated by PingOne.  If not set, the default value for `claim_enabled` is `false`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 custom resources.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create and manage the resource in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the resource name, which must be provided and must be unique within an environment.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description of the resource.").Description,
				Optional:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,

				Default: stringdefault.StaticString(string(management.ENUMRESOURCETYPE_CUSTOM)),

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"audience": schema.StringAttribute{
				Description:         audienceDescription.Description,
				MarkdownDescription: audienceDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidatorinternal.ShouldNotContain("pingone", "pingidentity"),
				},
			},

			"access_token_validity_seconds": schema.Int64Attribute{
				Description:         accessTokenValiditySecondsDescription.Description,
				MarkdownDescription: accessTokenValiditySecondsDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(accessTokenValiditySecondsDefault),

				Validators: []validator.Int64{
					int64validator.Between(accessTokenValiditySecondsMin, accessTokenValiditySecondsMax),
				},
			},

			"introspect_endpoint_auth_method": schema.StringAttribute{
				Description:         introspectEndpointAuthMethodDescription.Description,
				MarkdownDescription: introspectEndpointAuthMethodDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumResourceIntrospectEndpointAuthMethodEnumValues)...),
				},

				Default: stringdefault.StaticString(string(management.ENUMRESOURCEINTROSPECTENDPOINTAUTHMETHOD_CLIENT_SECRET_BASIC)),
			},

			"application_permissions_settings": schema.SingleNestedAttribute{
				Description:         applicationPermissionsSettingsDescription.Description,
				MarkdownDescription: applicationPermissionsSettingsDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					resourceApplicationPermissionsSettingsTFObjectTypes,
					map[string]attr.Value{
						"claim_enabled": types.BoolValue(false),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"claim_enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean setting to enable application permission claims in the access token.").Description,
						Required:    true,
					},
				},
			},
		},
	}
}

// ModifyPlan
func (r *ResourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var namePlan, audiencePlan types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("audience"), &audiencePlan)...)

	if audiencePlan.IsNull() {
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &namePlan)...)

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("audience"), namePlan)...)
	}
}

func (r *ResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceResourceModel

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

	// Build the model for the API
	resource, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Resource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourcesApi.CreateResource(ctx, plan.EnvironmentId.ValueString()).Resource(*resource).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateResource",
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

func (r *ResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceResourceModel

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
	var response *management.Resource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourcesApi.ReadOneResource(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneResource",
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

func (r *ResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceResourceModel

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

	// Build the model for the API
	resource, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Resource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourcesApi.UpdateResource(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Resource(*resource).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateResource",
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

func (r *ResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceResourceModel

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
			fR, fErr := r.Client.ManagementAPIClient.ResourcesApi.DeleteResource(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteResource",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "resource_id",
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

func (p *ResourceResourceModel) expand(ctx context.Context) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := *management.NewResource(
		p.Name.ValueString(),
	)

	data.SetType(management.EnumResourceType(p.Type.ValueString()))

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Audience.IsNull() && !p.Audience.IsUnknown() {
		data.SetAudience(p.Audience.ValueString())
	} else {
		data.SetAudience(p.Name.ValueString())
	}

	if !p.AccessTokenValiditySeconds.IsNull() && !p.AccessTokenValiditySeconds.IsUnknown() {
		data.SetAccessTokenValiditySeconds(int32(p.AccessTokenValiditySeconds.ValueInt64()))
	}

	if !p.IntrospectEndpointAuthMethod.IsNull() && !p.IntrospectEndpointAuthMethod.IsUnknown() {
		data.SetIntrospectEndpointAuthMethod(management.EnumResourceIntrospectEndpointAuthMethod(p.IntrospectEndpointAuthMethod.ValueString()))
	}

	if !p.ApplicationPermissionsSettings.IsNull() && !p.ApplicationPermissionsSettings.IsUnknown() {
		var plan ResourceApplicationPermissionsSettingsModel
		diags.Append(p.ApplicationPermissionsSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		applicationPermissionsSettings := management.NewResourceApplicationPermissionsSettings()

		if !plan.ClaimEnabled.IsNull() && !plan.ClaimEnabled.IsUnknown() {
			applicationPermissionsSettings.SetClaimEnabled(plan.ClaimEnabled.ValueBool())
		}

		data.SetApplicationPermissionsSettings(*applicationPermissionsSettings)
	}

	return &data, diags
}

func (p *ResourceResourceModel) toState(apiObject *management.Resource) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Audience = framework.StringOkToTF(apiObject.GetAudienceOk())
	p.AccessTokenValiditySeconds = framework.Int32OkToTF(apiObject.GetAccessTokenValiditySecondsOk())

	p.ApplicationPermissionsSettings, d = resourceApplicationPermissionsSettingsOk(apiObject.GetApplicationPermissionsSettingsOk())
	diags.Append(d...)

	p.IntrospectEndpointAuthMethod = framework.EnumOkToTF(apiObject.GetIntrospectEndpointAuthMethodOk())

	return diags
}

func resourceApplicationPermissionsSettingsOk(apiObject *management.ResourceApplicationPermissionsSettings, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(resourceApplicationPermissionsSettingsTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"claim_enabled": framework.BoolOkToTF(apiObject.GetClaimEnabledOk()),
	}

	returnVar, d := types.ObjectValue(resourceApplicationPermissionsSettingsTFObjectTypes, objMap)
	diags.Append(d...)

	return returnVar, diags
}
