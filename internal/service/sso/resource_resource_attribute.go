package sso

import (
	"context"
	"fmt"
	"net/http"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	customboolvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
	"golang.org/x/exp/slices"
)

// Types
type ResourceAttributeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ResourceAttributeResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	ResourceId      types.String `tfsdk:"resource_id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
	IDTokenEnabled  types.Bool   `tfsdk:"id_token_enabled"`
	UserinfoEnabled types.Bool   `tfsdk:"userinfo_enabled"`
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

			"resource_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the resource to assign the resource attribute to."),
			),

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

func (r *ResourceAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	resourceType, d := plan.getResourceType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(*resourceType)
	isOverriddenAttribute := plan.isOverriddenAttribute(*resourceType)

	resp.Diagnostics.Append(plan.validate(*resourceType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute, d := plan.expand(ctx, r.client, *resourceType, isCoreAttribute || isOverriddenAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response interface{}
	if !isCoreAttribute && !isOverriddenAttribute {
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ResourceAttributesApi.CreateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString()).ResourceAttribute(*resourceAttribute).Execute()
			},
			"CreateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	} else {
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), resourceAttribute.GetId()).ResourceAttribute(*resourceAttribute).Execute()
			},
			"UpdateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceAttributeResourceModel

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
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ResourceAttributesApi.ReadOneResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneResourceAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	resourceType, d := plan.getResourceType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(*resourceType)
	isOverriddenAttribute := plan.isOverriddenAttribute(*resourceType)

	resp.Diagnostics.Append(plan.validate(*resourceType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute, d := plan.expand(ctx, r.client, *resourceType, isCoreAttribute || isOverriddenAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Id.ValueString()).ResourceAttribute(*resourceAttribute).Execute()
		},
		"UpdateResourceAttribute",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceAttributeResourceModel

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
	var d diag.Diagnostics
	if data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_PREDEFINED))) || data.Type.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_CORE))) {

		resourceType, d := data.getResourceType(ctx, r.client)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

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

			coreAttributeData, ok := data.isCoreAttribute(*resourceType)
			if !ok {
				resp.Diagnostics.AddError(
					"Core attribute mismatch error",
					"The provider cannot determine the core attribute values to reset.  Please raise this issue to the provider maintainers.",
				)
				return
			}

			data.Value = framework.StringToTF(coreAttributeData.defaultValue)
		}

		resourceMapping, d = data.expand(ctx, r.client, *resourceType, true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		_, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ResourceAttributesApi.UpdateResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).ResourceAttribute(*resourceMapping).Execute()
			},
			"UpdateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(d...)
	} else {

		_, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				r, err := r.client.ResourceAttributesApi.DeleteResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
				return nil, r, err
			},
			"DeleteResourceAttribute",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(d...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/resource_id/resource_attribute_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *ResourceAttributeResourceModel) getResourceType(ctx context.Context, apiClient *management.APIClient) (*management.EnumResourceType, diag.Diagnostics) {
	var diags diag.Diagnostics

	respObject, d := fetchResource_Framework(ctx, apiClient, p.EnvironmentId.ValueString(), p.ResourceId.ValueString())
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return respObject.GetType().Ptr(), diags
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

func (p *ResourceAttributeResourceModel) toState(apiObject *management.ResourceAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.Type = ResourceAttributeTypeOkToTF(apiObject.GetTypeOk())
	p.IDTokenEnabled = framework.BoolOkToTF(apiObject.GetIdTokenOk())
	p.UserinfoEnabled = framework.BoolOkToTF(apiObject.GetUserInfoOk())

	return diags
}

func ResourceAttributeTypeOkToTF(v *management.EnumResourceAttributeType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
