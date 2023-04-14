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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type IdentityProviderCoreAttributeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type IdentityProviderCoreAttributeResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	IdentityProviderId types.String `tfsdk:"identity_provider_id"`
	Name               types.String `tfsdk:"name"`
	Update             types.String `tfsdk:"update"`
	Value              types.String `tfsdk:"value"`
	MappingType        types.String `tfsdk:"mapping_type"`
}

type coreIdentityProviderAttributeType struct {
	name     string
	defaults map[management.EnumIdentityProviderExt]string
}

// Framework interfaces
var (
	_ resource.Resource                = &IdentityProviderCoreAttributeResource{}
	_ resource.ResourceWithConfigure   = &IdentityProviderCoreAttributeResource{}
	_ resource.ResourceWithImportState = &IdentityProviderCoreAttributeResource{}
)

// New Object
func NewIdentityProviderCoreAttributeResource() resource.Resource {
	return &IdentityProviderCoreAttributeResource{}
}

var (
	idpCoreAttrMetadata = []coreIdentityProviderAttributeType{
		{
			name: "username",
			defaults: map[management.EnumIdentityProviderExt]string{
				management.ENUMIDENTITYPROVIDEREXT_AMAZON:         "${providerAttributes.user_id}",
				management.ENUMIDENTITYPROVIDEREXT_APPLE:          "${providerAttributes.sub}",
				management.ENUMIDENTITYPROVIDEREXT_FACEBOOK:       "${providerAttributes.email}",
				management.ENUMIDENTITYPROVIDEREXT_GITHUB:         "${providerAttributes.id}",
				management.ENUMIDENTITYPROVIDEREXT_GOOGLE:         "${providerAttributes.emailAddress.value}",
				management.ENUMIDENTITYPROVIDEREXT_LINKEDIN:       "${providerAttributes.emailAddress}",
				management.ENUMIDENTITYPROVIDEREXT_MICROSOFT:      "${providerAttributes.id}",
				management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT: "${providerAttributes.sub}",
				management.ENUMIDENTITYPROVIDEREXT_PAYPAL:         "${providerAttributes.user_id}",
				management.ENUMIDENTITYPROVIDEREXT_SAML:           "${samlAssertion.subject}",
				management.ENUMIDENTITYPROVIDEREXT_TWITTER:        "${providerAttributes.id}",
				management.ENUMIDENTITYPROVIDEREXT_YAHOO:          "${providerAttributes.sub}",
			},
		},
	}
)

// Metadata
func (r *IdentityProviderCoreAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider_core_attribute"
}

// Schema.
func (r *IdentityProviderCoreAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	nameAttributeValues := make([]string, 0)

	// Loop the core attrs for the application type
	for _, coreAttr := range idpCoreAttrMetadata {
		nameAttributeValues = append(nameAttributeValues, coreAttr.name)
	}
	nameDescriptionFmt := fmt.Sprintf("A string that specifies the name of the PingOne core directory attribute to map the Identity Provider attribute value to.  The following are valid values: `%s`.", strings.Join(nameAttributeValues, "`, `"))
	nameDescription := framework.SchemaDescription{
		MarkdownDescription: nameDescriptionFmt,
		Description:         strings.ReplaceAll(nameDescriptionFmt, "`", "\""),
	}

	updateDescriptionFmt := fmt.Sprintf("Indicates whether to update the user attribute in the directory with the non-empty mapped value from the IdP. Options are `%s` (only update the user attribute if it has an empty value); `%s` (always update the user attribute value).", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS))
	updateDescription := framework.SchemaDescription{
		MarkdownDescription: updateDescriptionFmt,
		Description:         strings.ReplaceAll(updateDescriptionFmt, "`", "\""),
	}

	valueDescriptionFmt := "A placeholder referring to the attribute (or attributes) from the provider. Placeholders must be valid for the attributes returned by the IdP type and use the `${}` syntax (for example, `${email}`). For SAML, any placeholder is acceptable, and it is mapped against the attributes available in the SAML assertion after authentication. The `${samlAssertion.subject}` placeholder is a special reserved placeholder used to refer to the subject name ID in the SAML assertion response."
	valueDescription := framework.SchemaDescription{
		MarkdownDescription: valueDescriptionFmt,
		Description:         strings.ReplaceAll(valueDescriptionFmt, "`", "\""),
	}

	mappingTypeDescriptionFmt := fmt.Sprintf("The mapping type. Options are `%s` (This attribute is required by the schema and cannot be removed. The name and update properties cannot be changed.) or `%s` (All user-created attributes are of this type.)", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CORE), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CUSTOM))
	mappingTypeDescription := framework.SchemaDescription{
		MarkdownDescription: mappingTypeDescriptionFmt,
		Description:         strings.ReplaceAll(mappingTypeDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a core attribute mapping for identity providers configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the identity provider attribute in."},
			),

			"identity_provider_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the identity provider to create the attribute mapping for."},
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.OneOf(nameAttributeValues...),
				},
			},

			"update": schema.StringAttribute{
				Description:         updateDescription.Description,
				MarkdownDescription: updateDescription.MarkdownDescription,
				Computed:            true,
			},

			"value": schema.StringAttribute{
				Description:         valueDescription.Description,
				MarkdownDescription: valueDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"mapping_type": schema.StringAttribute{
				Description:         mappingTypeDescription.Description,
				MarkdownDescription: mappingTypeDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *IdentityProviderCoreAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentityProviderCoreAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state IdentityProviderCoreAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute()

	resp.Diagnostics.Append(plan.validate(isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	idpAttributeMapping, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString(), idpAttributeMapping.GetId()).IdentityProviderAttribute(*idpAttributeMapping).Execute()
		},
		"UpdateIdentityProviderAttribute",
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
	resp.Diagnostics.Append(state.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IdentityProviderCoreAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityProviderCoreAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneIdentityProviderAttribute",
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
	resp.Diagnostics.Append(data.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderCoreAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IdentityProviderCoreAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute()

	resp.Diagnostics.Append(plan.validate(isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	idpAttributeMapping, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString(), plan.Id.ValueString()).IdentityProviderAttribute(*idpAttributeMapping).Execute()
		},
		"UpdateIdentityProviderAttribute",
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
	resp.Diagnostics.Append(state.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IdentityProviderCoreAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityProviderCoreAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var d diag.Diagnostics

	idpType, d := data.getIdentityProviderType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if coreAttributeData, ok := data.isCoreAttribute(); ok {

		data.Value = framework.StringToTF(coreAttributeData.defaults[*idpType])

		idpAttributeMapping, d := data.expand(ctx, r.client)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		_, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).IdentityProviderAttribute(*idpAttributeMapping).Execute()
			},
			"UpdateIdentityProviderAttribute",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
		)

		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Unexpected core attribute",
			"The core attribute identified is unexpected.  Please report this issue to the provider maintainers.",
		)
		return
	}

}

func (r *IdentityProviderCoreAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/identity_provider_id/identity_provider_core_attribute_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identity_provider_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *IdentityProviderCoreAttributeResourceModel) getIdentityProviderType(ctx context.Context, apiClient *management.APIClient) (*management.EnumIdentityProviderExt, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get application type and verify against the set params
	resp, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, p.EnvironmentId.ValueString(), p.IdentityProviderId.ValueString()).Execute()
		},
		"ReadOneIdentityProvider",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	respObject := resp.(*management.IdentityProvider)

	var idpType *management.EnumIdentityProviderExt
	if respObject.IdentityProviderApple != nil && respObject.IdentityProviderApple.GetId() != "" {
		idpType = &respObject.IdentityProviderApple.Type
	} else if respObject.IdentityProviderClientIDClientSecret != nil && respObject.IdentityProviderClientIDClientSecret.GetId() != "" {
		idpType = &respObject.IdentityProviderClientIDClientSecret.Type
	} else if respObject.IdentityProviderFacebook != nil && respObject.IdentityProviderFacebook.GetId() != "" {
		idpType = &respObject.IdentityProviderFacebook.Type
	} else if respObject.IdentityProviderOIDC != nil && respObject.IdentityProviderOIDC.GetId() != "" {
		idpType = &respObject.IdentityProviderOIDC.Type
	} else if respObject.IdentityProviderPaypal != nil && respObject.IdentityProviderPaypal.GetId() != "" {
		idpType = &respObject.IdentityProviderPaypal.Type
	} else if respObject.IdentityProviderSAML != nil && respObject.IdentityProviderSAML.GetId() != "" {
		idpType = &respObject.IdentityProviderSAML.Type
	} else {
		diags.AddError(
			"Invalid parameter value - Unmappable identity provider type",
			fmt.Sprintf("The identity provider ID provided (%s) relates to an unknown type.  Attributes cannot be mapped to this identity provider.", p.IdentityProviderId.ValueString()),
		)
		return nil, diags
	}

	return idpType, diags
}

func (p *IdentityProviderCoreAttributeResourceModel) validate(isCoreAttribute bool) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that we're a core attribute
	if !isCoreAttribute {
		diags.AddError(
			"Invalid parameter value - Invalid attribute name",
			fmt.Sprintf("The attribute provided (%s) is not a core attribute.  Use the `pingone_identity_provider_attribute_mapping` resource to map custom attribute values.", p.Name.ValueString()),
		)
		return diags
	}

	return diags
}

func (p *IdentityProviderCoreAttributeResourceModel) isCoreAttribute() (*coreIdentityProviderAttributeType, bool) {

	// Loop the core attrs for the application type
	for _, coreAttr := range idpCoreAttrMetadata {
		if strings.ToUpper(p.Name.ValueString()) == strings.ToUpper(coreAttr.name) {
			// We're a core attribute
			return &coreAttr, true
		}
	}

	return nil, false
}

func (p *IdentityProviderCoreAttributeResourceModel) expand(ctx context.Context, apiClient *management.APIClient) (*management.IdentityProviderAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.IdentityProviderAttribute

	// fetch the attribute that already exists
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderAttributesApi.ReadAllIdentityProviderAttributes(ctx, p.EnvironmentId.ValueString(), p.IdentityProviderId.ValueString()).Execute()
		},
		"ReadAllIdentityProviderAttributes",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(diags...)
	if diags.HasError() {
		return nil, diags
	}

	if attributes, ok := response.(*management.EntityArray).Embedded.GetAttributesOk(); ok {

		found := false
		for _, attribute := range attributes {

			if strings.ToUpper(attribute.IdentityProviderAttribute.GetName()) == strings.ToUpper(p.Name.ValueString()) {
				data = attribute.IdentityProviderAttribute
				found = true
				break
			}
		}

		if !found {
			diags.AddError(
				"Core attribute not found",
				fmt.Sprintf("The configured attribute name (\"%s\") is identified as a core attribute, but the attribute cannot be found in the platform.  Please raise this issue with the provider maintainers.", p.Name.ValueString()),
			)

			return nil, diags
		}

	}

	data.SetValue(p.Value.ValueString())

	return data, diags
}

func (p *IdentityProviderCoreAttributeResourceModel) toState(apiObject *management.IdentityProviderAttribute) diag.Diagnostics {
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
	p.Update = IdentityProviderCoreAttributeMappingUpdateOkToTF(apiObject.GetUpdateOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.MappingType = IdentityProviderCoreAttributeMappingTypeOkToTF(apiObject.GetMappingTypeOk())

	return diags
}

func IdentityProviderCoreAttributeMappingTypeOkToTF(v *management.EnumIdentityProviderAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func IdentityProviderCoreAttributeMappingUpdateOkToTF(v *management.EnumIdentityProviderAttributeMappingUpdate, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
