// Copyright © 2025 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type IdentityProviderAttributeResource serviceClientType

type IdentityProviderAttributeResourceModel struct {
	Id                 pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId      pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	IdentityProviderId pingonetypes.ResourceIDValue `tfsdk:"identity_provider_id"`
	Name               types.String                 `tfsdk:"name"`
	Update             types.String                 `tfsdk:"update"`
	Value              types.String                 `tfsdk:"value"`
	MappingType        types.String                 `tfsdk:"mapping_type"`
}

type coreIdentityProviderAttributeType struct {
	name     string
	defaults map[management.EnumIdentityProviderExt]string
}

// Framework interfaces
var (
	_ resource.Resource                = &IdentityProviderAttributeResource{}
	_ resource.ResourceWithConfigure   = &IdentityProviderAttributeResource{}
	_ resource.ResourceWithImportState = &IdentityProviderAttributeResource{}
)

// New Object
func NewIdentityProviderAttributeResource() resource.Resource {
	return &IdentityProviderAttributeResource{}
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
func (r *IdentityProviderAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider_attribute"
}

// Schema.
func (r *IdentityProviderAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	reservedNames := []string{"account", "id", "created", "updated", "lifecycle", "mfaEnabled", "enabled"}
	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the name of the PingOne directory attribute to map the Identity Provider attribute value to. The attribute must not be defined as read only from the user schema or of type `COMPLEX` based on the user schema. Examples are `email`, `name.given`.  The following attributes may not be used: `%s`.", strings.Join(reservedNames, "`, `")),
	)

	updateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Indicates whether to update the user attribute in the directory with the non-empty mapped value from the IdP. Options are `%s` (only update the user attribute if it has an empty value); `%s` (always update the user attribute value). Defaults to `%s`.", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY)),
	)

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A placeholder referring to the attribute (or attributes) from the provider. Placeholders must be valid for the attributes returned by the IdP type and use the `${}` syntax (for example, `${email}`). For SAML, any placeholder is acceptable, and it is mapped against the attributes available in the SAML assertion after authentication. The `${samlAssertion.subject}` placeholder is a special reserved placeholder used to refer to the subject name ID in the SAML assertion response.\n\nFor a full list of social providers and their configurable attribute mappings, see the [Social Provider and External Identity Provider Attribute Reference](#social-provider-and-external-identity-provider-attribute-reference) below.",
	)

	mappingTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The mapping type. Options are `%s` (This attribute is required by the schema and cannot be removed. The name and update properties cannot be changed.) or `%s` (All user-created attributes are of this type.)", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CORE), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CUSTOM)),
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for identity providers configured in PingOne.\n\nFor a full list of social providers and their configurable attribute mappings, see the [Social Provider and External Identity Provider Attribute Reference](#social-provider-and-external-identity-provider-attribute-reference) below.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the identity provider attribute in."),
			),

			"identity_provider_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the identity provider to create the attribute mapping for."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.NoneOf(reservedNames...),
				},
			},

			"update": schema.StringAttribute{
				Description:         updateDescription.Description,
				MarkdownDescription: updateDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS)),
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

			"mapping_type": schema.StringAttribute{
				Description:         mappingTypeDescription.Description,
				MarkdownDescription: mappingTypeDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *IdentityProviderAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentityProviderAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state IdentityProviderAttributeResourceModel

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

	_, isCoreAttribute := plan.isCoreAttribute()

	resp.Diagnostics.Append(plan.validate(isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	identityProviderAttribute, d := plan.expand(ctx, r.Client.ManagementAPIClient, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProviderAttribute
	if !isCoreAttribute {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.CreateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString()).IdentityProviderAttribute(*identityProviderAttribute).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateIdentityProviderAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString(), identityProviderAttribute.GetId()).IdentityProviderAttribute(*identityProviderAttribute).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateIdentityProviderAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IdentityProviderAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityProviderAttributeResourceModel

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
	var response *management.IdentityProviderAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneIdentityProviderAttribute",
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

func (r *IdentityProviderAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IdentityProviderAttributeResourceModel

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

	_, isCoreAttribute := plan.isCoreAttribute()

	resp.Diagnostics.Append(plan.validate(isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	identityProviderAttributeMapping, d := plan.expand(ctx, r.Client.ManagementAPIClient, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProviderAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString(), plan.Id.ValueString()).IdentityProviderAttribute(*identityProviderAttributeMapping).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateIdentityProviderAttribute",
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

func (r *IdentityProviderAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityProviderAttributeResourceModel

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
	if data.MappingType.Equal(types.StringValue(string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CORE))) {

		idpType, d := data.getIdentityProviderType(ctx, r.Client.ManagementAPIClient)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		coreAttributeData, ok := data.isCoreAttribute()
		if !ok {
			resp.Diagnostics.AddError(
				"Core attribute mismatch error",
				"The provider cannot determine the core attribute values to reset.  Please raise this issue to the provider maintainers.",
			)
			return
		}

		data.Value = framework.StringToTF(coreAttributeData.defaults[*idpType])

		idpMapping, d := data.expand(ctx, r.Client.ManagementAPIClient, true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).IdentityProviderAttribute(*idpMapping).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateIdentityProviderAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := r.Client.ManagementAPIClient.IdentityProviderAttributesApi.DeleteIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
			},
			"DeleteIdentityProviderAttribute",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)

	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *IdentityProviderAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "identity_provider_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "identity_provider_attribute_id",
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

func (p *IdentityProviderAttributeResourceModel) getIdentityProviderType(ctx context.Context, apiClient *management.APIClient) (*management.EnumIdentityProviderExt, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get application type and verify against the set params
	var respObject *management.IdentityProvider
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, p.EnvironmentId.ValueString(), p.IdentityProviderId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, p.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneIdentityProvider",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&respObject,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	var idpType *management.EnumIdentityProviderExt
	if respObject.IdentityProviderApple != nil && respObject.IdentityProviderApple.GetId() != "" {
		idpType = &respObject.IdentityProviderApple.Type
	} else if respObject.IdentityProviderClientIDClientSecret != nil && respObject.IdentityProviderClientIDClientSecret.GetId() != "" {
		idpType = &respObject.IdentityProviderClientIDClientSecret.Type
	} else if respObject.IdentityProviderFacebook != nil && respObject.IdentityProviderFacebook.GetId() != "" {
		idpType = &respObject.IdentityProviderFacebook.Type
	} else if respObject.IdentityProviderOIDC != nil && respObject.IdentityProviderOIDC.GetId() != "" {
		idpType = &respObject.IdentityProviderOIDC.Type
	} else if respObject.IdentityProviderMicrosoft != nil && respObject.IdentityProviderMicrosoft.GetId() != "" {
		idpType = &respObject.IdentityProviderMicrosoft.Type
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

func (p *IdentityProviderAttributeResourceModel) validate(isCoreAttribute bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if isCoreAttribute {
		if !p.Update.IsNull() && !p.Update.IsUnknown() {
			diags.AddError(
				"Invalid parameter value - Parameter doesn't apply to attribute type",
				fmt.Sprintf("The `name` provided (%s) is a core attribute.  The `update` parameter cannot be set for core attributes.", p.Name.ValueString()),
			)
			return diags
		}
	}

	return diags
}

func (p *IdentityProviderAttributeResourceModel) isCoreAttribute() (*coreIdentityProviderAttributeType, bool) {

	// Loop the core attrs for the IDP type
	for _, coreAttr := range idpCoreAttrMetadata {
		if strings.EqualFold(p.Name.ValueString(), coreAttr.name) {
			// We're a core attribute
			return &coreAttr, true
		}
	}

	return nil, false
}

func (p *IdentityProviderAttributeResourceModel) expand(ctx context.Context, apiClient *management.APIClient, overrideExisting bool) (*management.IdentityProviderAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.IdentityProviderAttribute

	if overrideExisting {

		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := apiClient.IdentityProviderAttributesApi.ReadAllIdentityProviderAttributes(ctx, p.EnvironmentId.ValueString(), p.IdentityProviderId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, p.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if idpAttributes, ok := pageCursor.EntityArray.Embedded.GetAttributesOk(); ok {

						for _, idpAttribute := range idpAttributes {

							if strings.EqualFold(idpAttribute.IdentityProviderAttribute.GetName(), p.Name.ValueString()) {
								return idpAttribute.IdentityProviderAttribute, pageCursor.HTTPResponse, nil
							}
						}

					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllIdentityProviderAttributes",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&data,
		)...)
		if diags.HasError() {
			return nil, diags
		}

		if data == nil {
			diags.AddError(
				fmt.Sprintf("Cannot find identity provider attribute %s", p.Name.ValueString()),
				"The identity provider attribute cannot be found by the provided name.",
			)

			return nil, diags
		}

		data.SetValue(p.Value.ValueString())

	} else {
		data = management.NewIdentityProviderAttribute(p.Name.ValueString(), p.Value.ValueString(), management.EnumIdentityProviderAttributeMappingUpdate(p.Update.ValueString()))

		if !p.Update.IsNull() && !p.Update.IsUnknown() {
			data.SetUpdate(management.EnumIdentityProviderAttributeMappingUpdate(p.Update.ValueString()))
		} else {
			data.SetUpdate(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY)
		}
	}

	return data, diags
}

func (p *IdentityProviderAttributeResourceModel) toState(apiObject *management.IdentityProviderAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Update = IdentityProviderAttributeMappingUpdateOkToTF(apiObject.GetUpdateOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.MappingType = IdentityProviderAttributeMappingTypeOkToTF(apiObject.GetMappingTypeOk())

	return diags
}

func IdentityProviderAttributeMappingTypeOkToTF(v *management.EnumIdentityProviderAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func IdentityProviderAttributeMappingUpdateOkToTF(v *management.EnumIdentityProviderAttributeMappingUpdate, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
