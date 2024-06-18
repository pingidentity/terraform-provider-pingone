package sso

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

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
	customboolvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ResourceAttributeResource serviceClientType

type ResourceAttributeResourceModel struct {
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId      pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	ResourceName    types.String                 `tfsdk:"resource_name"`
	Name            types.String                 `tfsdk:"name"`
	Type            types.String                 `tfsdk:"type"`
	Value           types.String                 `tfsdk:"value"`
	IDTokenEnabled  types.Bool                   `tfsdk:"id_token_enabled"`
	UserinfoEnabled types.Bool                   `tfsdk:"userinfo_enabled"`
}

type coreResourceAttributeType struct {
	name         string
	defaultValue string
}

// Framework interfaces
var (
	_ resource.Resource                = &ResourceAttributeResource{}
	_ resource.ResourceWithConfigure   = &ResourceAttributeResource{}
	_ resource.ResourceWithImportState = &ResourceAttributeResource{}
)

// New Object
func NewResourceAttributeResource() resource.Resource {
	return &ResourceAttributeResource{}
}

var (
	resourceCoreAttrMetadata = map[management.EnumResourceType][]coreResourceAttributeType{
		management.ENUMRESOURCETYPE_CUSTOM: {
			{
				name:         "sub",
				defaultValue: "${user.id}",
			},
		},
	}
)

// Metadata
func (r *ResourceAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_attribute"
}

// Schema.
func (r *ResourceAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the resource that the attribute is assigned to.",
	)

	resourceNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource to assign the resource attribute to.  The built-in OpenID Connect resource name is `openid`.",
	).RequiresReplace()

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the name of the resource attribute to map a value for. When the resource's type property is `OPENID_CONNECT`, the following are reserved names and cannot be used: %s.  The resource will also override the default configured values for a resource, rather than creating new attributes.  For resources of type `CUSTOM`, the `sub` name is overridden.  For resources of type `OPENID_CONNECT`, the following names are overridden: %s.", verify.IllegalOIDCAttributeNameString(), verify.OverrideOIDCAttributeNameString()),
	)

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the value of the custom resource attribute. This value can be a placeholder that references an attribute in the user schema, expressed as `${user.path.to.value}`, or it can be an expression, or a static string. Placeholders must be valid, enabled attributes in the environment’s user schema. Examples of valid values are: `${user.email}`, `${user.name.family}`, and `myClaimValueString`.  Note that definition in HCL requires escaping with the `$` character when defining attribute paths, for example `value = \"$${user.email}\"`.",
	)

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified).",
	)

	idTokenEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the attribute mapping should be available in the ID Token.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`.",
	)

	userinfoEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the attribute mapping should be available through the /as/userinfo endpoint.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource attributes in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the resource attribute in."),
			),

			"resource_id": schema.StringAttribute{
				Description:         resourceIdDescription.Description,
				MarkdownDescription: resourceIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"resource_name": schema.StringAttribute{
				Description:         resourceNameDescription.Description,
				MarkdownDescription: resourceNameDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"value": schema.StringAttribute{
				Description:         valueDescription.Description,
				MarkdownDescription: valueDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"id_token_enabled": schema.BoolAttribute{
				Description:         idTokenEnabledDescription.Description,
				MarkdownDescription: idTokenEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(
						types.BoolValue(true),
						types.BoolValue(true),
						path.MatchRelative().AtParent().AtName("userinfo_enabled"),
					),
				},
			},

			"userinfo_enabled": schema.BoolAttribute{
				Description:         userinfoEnabledDescription.Description,
				MarkdownDescription: userinfoEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(
						types.BoolValue(true),
						types.BoolValue(true),
						path.MatchRelative().AtParent().AtName("id_token_enabled"),
					),
				},
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ResourceAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	resourceResponse, d := fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceName.ValueString(), false)

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ResourceId = framework.PingOneResourceIDOkToTF(resourceResponse.GetIdOk())

	_, isCoreAttribute := plan.isCoreAttribute(resourceResponse.GetType())
	isOverriddenAttribute := plan.isOverriddenAttribute(resourceResponse.GetType())

	resp.Diagnostics.Append(plan.validate(resourceResponse.GetType())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute, d := plan.expand(ctx, r.Client.ManagementAPIClient, resourceResponse.GetType(), isCoreAttribute || isOverriddenAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceAttributeResponse *management.ResourceAttribute
	if !isCoreAttribute && !isOverriddenAttribute {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.CreateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), resourceResponse.GetId()).ResourceAttribute(*resourceAttribute).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&resourceAttributeResponse,
		)...)
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), resourceResponse.GetId(), resourceAttribute.GetId()).ResourceAttribute(*resourceAttribute).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&resourceAttributeResponse,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceAttributeResourceModel

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

	var resourceResponse *management.Resource
	var d diag.Diagnostics
	if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
		resourceResponse, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), true)
	} else {
		resourceResponse, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceName.ValueString(), true)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resourceResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Run the API call
	var resourceAttributeResponse *management.ResourceAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.ReadOneResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneResourceAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&resourceAttributeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resourceAttributeResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	var resourceResponse *management.Resource
	var d diag.Diagnostics
	if !plan.ResourceId.IsNull() && !plan.ResourceId.IsUnknown() {
		resourceResponse, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), false)
	} else {
		resourceResponse, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceName.ValueString(), false)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(resourceResponse.GetType())
	isOverriddenAttribute := plan.isOverriddenAttribute(resourceResponse.GetType())

	resp.Diagnostics.Append(plan.validate(resourceResponse.GetType())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute, d := plan.expand(ctx, r.Client.ManagementAPIClient, resourceResponse.GetType(), isCoreAttribute || isOverriddenAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceAttributeResponse *management.ResourceAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), resourceResponse.GetId(), plan.Id.ValueString()).ResourceAttribute(*resourceAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateResourceAttribute",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resourceAttributeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceAttributeResourceModel

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

	var resource *management.Resource
	var d diag.Diagnostics
	if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
		resource, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), true)
	} else {
		resource, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceName.ValueString(), true)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resource == nil {
		return
	}

	// Run the API call
	if data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_PREDEFINED))) || data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_CORE))) {

		// defaults
		var resourceMapping *management.ResourceAttribute
		if data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_PREDEFINED))) {
			defaultValues := map[string]string{
				"address.country":        "${user.address.countryCode}",
				"address.formatted":      "",
				"address.locality":       "${user.address.locality}",
				"address.postal_code":    "${user.address.postalCode}",
				"address.region":         "${user.address.region}",
				"address.street_address": "${user.address.streetAddress}",
				"birthdate":              "",
				"email_verified":         "",
				"email":                  "${user.email}",
				"family_name":            "${user.name.family}",
				"gender":                 "",
				"given_name":             "${user.name.given}",
				"locale":                 "${user.locale}",
				"middle_name":            "${user.name.middle}",
				"name":                   "${user.name.formatted}",
				"nickname":               "${user.nickname}",
				"phone_number_verified":  "",
				"phone_number":           "${user.primaryPhone}",
				"picture":                "${user.photo.href}",
				"preferred_username":     "${user.username}",
				"profile":                "",
				"updated_at":             "${#datetime.toUnixTimestamp(user.updatedAt)}",
				"website":                "",
				"zoneinfo":               "${user.timezone}",
			}

			data.Value = framework.StringToTF(defaultValues[data.Name.ValueString()])

		} else if data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_CORE))) {

			coreAttributeData, ok := data.isCoreAttribute(resource.GetType())
			if !ok {
				resp.Diagnostics.AddError(
					"Core attribute mismatch error",
					"The provider cannot determine the core attribute values to reset.  Please raise this issue to the provider maintainers.",
				)
				return
			}

			data.Value = framework.StringToTF(coreAttributeData.defaultValue)
		}

		resourceMapping, d = data.expand(ctx, r.Client.ManagementAPIClient, resource.GetType(), true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).ResourceAttribute(*resourceMapping).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)
	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := r.Client.ManagementAPIClient.ResourceAttributesApi.DeleteResourceAttribute(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
			},
			"DeleteResourceAttribute",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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
			Label:     "resource_attribute_id",
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

func (p *ResourceAttributeResourceModel) validate(resourceType management.EnumResourceType) diag.Diagnostics {
	var diags diag.Diagnostics

	if resourceType != management.ENUMRESOURCETYPE_OPENID_CONNECT && resourceType != management.ENUMRESOURCETYPE_CUSTOM {
		diags.AddError(
			"Invalid parameter value - Invalid resource type",
			fmt.Sprintf("The `resource_id` provided (%s) is neither %s or %s type.  Attributes can only be created for CUSTOM and OPENID_CONNECT (openid) resources.", p.ResourceId.ValueString(), string(management.ENUMRESOURCETYPE_CUSTOM), string(management.ENUMRESOURCETYPE_OPENID_CONNECT)),
		)
		return diags
	}

	if resourceType != management.ENUMRESOURCETYPE_OPENID_CONNECT {

		// If we have defined values that we shouldn't.  If unknown, deal with when the values are known.
		if (!p.IDTokenEnabled.IsNull() && !p.IDTokenEnabled.IsUnknown()) ||
			(!p.UserinfoEnabled.IsNull() && !p.UserinfoEnabled.IsUnknown()) {
			diags.AddError(
				"Invalid parameter value - Parameter doesn't apply to resource type",
				fmt.Sprintf("The `resource_id` provided (%s) is of type `%s`.  The `id_token_enabled` and `userinfo_enabled` attributes do not apply to this resource type.", p.ResourceId.ValueString(), resourceType),
			)
			return diags
		}
	}

	if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
		if slices.Contains(verify.IllegalOIDCattributeNamesList(), p.Name.ValueString()) {
			diags.AddError(
				fmt.Sprintf("Invalid attribute name `%s` for the configured OpenID Connect resource.", p.Name.ValueString()),
				fmt.Sprintf("The `resource_id` provided (%s) is of type `%s`. The attribute name provided, `%s`, is reserved and cannot be used for this resource type.", p.ResourceId.ValueString(), resourceType, p.Name.ValueString()),
			)
			return diags
		}
	}

	return diags
}

func (p *ResourceAttributeResourceModel) isCoreAttribute(resourceType management.EnumResourceType) (*coreResourceAttributeType, bool) {

	// Evaluate against the core attribute
	if v, ok := resourceCoreAttrMetadata[resourceType]; ok {
		// Loop the core attrs for the resource type
		for _, coreAttr := range v {
			if strings.EqualFold(p.Name.ValueString(), coreAttr.name) {
				// We're a core attribute
				return &coreAttr, true
			}
		}
	}

	return nil, false
}

func (p *ResourceAttributeResourceModel) isOverriddenAttribute(resourceType management.EnumResourceType) bool {

	if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
		if slices.Contains(verify.OverrideOIDCAttributeNameList(), p.Name.ValueString()) {
			return true
		}
	}

	return false
}

func (p *ResourceAttributeResourceModel) expand(ctx context.Context, apiClient *management.APIClient, resourceType management.EnumResourceType, overrideExisting bool) (*management.ResourceAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ResourceAttribute

	if overrideExisting {

		var d diag.Diagnostics
		data, d = fetchResourceAttributeFromName_Framework(ctx, apiClient, p.EnvironmentId.ValueString(), p.ResourceId.ValueString(), p.Name.ValueString())
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetValue(p.Value.ValueString())

	} else {
		data = management.NewResourceAttribute(p.Name.ValueString(), p.Value.ValueString())

		if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
			if !p.IDTokenEnabled.IsNull() && !p.IDTokenEnabled.IsUnknown() {
				data.SetIdToken(p.IDTokenEnabled.ValueBool())
			} else {
				data.SetIdToken(true)
			}

			if !p.UserinfoEnabled.IsNull() && !p.UserinfoEnabled.IsUnknown() {
				data.SetUserInfo(p.UserinfoEnabled.ValueBool())
			} else {
				data.SetUserInfo(true)
			}
		}
	}

	return data, diags
}

func (p *ResourceAttributeResourceModel) toState(apiObject *management.ResourceAttribute, resourceApiObject *management.Resource) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || resourceApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.ResourceId = framework.PingOneResourceIDOkToTF(resourceApiObject.GetIdOk())
	p.ResourceName = framework.StringOkToTF(resourceApiObject.GetNameOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.IDTokenEnabled = framework.BoolOkToTF(apiObject.GetIdTokenOk())
	p.UserinfoEnabled = framework.BoolOkToTF(apiObject.GetUserInfoOk())

	return diags
}
