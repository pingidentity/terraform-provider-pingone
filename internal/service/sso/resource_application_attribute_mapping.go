package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	customboolvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationAttributeMappingResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ApplicationAttributeMappingResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	EnvironmentId         types.String `tfsdk:"environment_id"`
	ApplicationId         types.String `tfsdk:"application_id"`
	Name                  types.String `tfsdk:"name"`
	Required              types.Bool   `tfsdk:"required"`
	Value                 types.String `tfsdk:"value"`
	MappingType           types.String `tfsdk:"mapping_type"`
	OIDCScopes            types.Set    `tfsdk:"oidc_scopes"`
	OIDCIDTokenEnabled    types.Bool   `tfsdk:"oidc_id_token_enabled"`
	OIDCUserinfoEnabled   types.Bool   `tfsdk:"oidc_userinfo_enabled"`
	SAMLSubjectNameformat types.String `tfsdk:"saml_subject_nameformat"`
}

type coreApplicationAttributeType struct {
	name         string
	defaultValue string
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationAttributeMappingResource{}
	_ resource.ResourceWithConfigure   = &ApplicationAttributeMappingResource{}
	_ resource.ResourceWithImportState = &ApplicationAttributeMappingResource{}
)

// New Object
func NewApplicationAttributeMappingResource() resource.Resource {
	return &ApplicationAttributeMappingResource{}
}

var (
	applicationCoreAttrMetadata = map[management.EnumApplicationProtocol][]coreApplicationAttributeType{
		management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT: {
			{
				name:         "sub",
				defaultValue: "${user.id}",
			},
		},
		management.ENUMAPPLICATIONPROTOCOL_SAML: {
			{
				name:         "saml_subject",
				defaultValue: "${user.id}",
			},
		},
	}
)

// Metadata
func (r *ApplicationAttributeMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_attribute_mapping"
}

// Schema.
func (r *ApplicationAttributeMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	reservedNames := []string{"acr", "amr", "at_hash", "aud", "auth_time", "azp", "client_id", "exp", "iat", "iss", "jti", "nbf", "nonce", "org", "scope", "sid"}
	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the name of attribute and must be unique within an application. For SAML applications, the `saml_subject` name is a case-insensitive name which indicates the mapping to be used for the subject in an assertion and can be overridden. For OpenID Connect applications, the `sub` name indicates the mapping to be used for the subject in the token and can be overridden.  The following OpenID Connect names are reserved and cannot be used: `%s`.", strings.Join(reservedNames, "`, `")),
	)

	requiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify whether a mapping value is required for this attribute. If `true`, a value must be set and a non-empty value must be available in the SAML assertion or ID token. If overriding a core attribute mapping (`saml_subject` for SAML applications and `sub` for OpenID Connect applications), then this value must be set to `true`.  Defaults to `false`.",
	)

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the string constants or expression for mapping the attribute path against a specific source. The expression format is `${<source>.<attribute_path>}`. The only supported source is user (for example, `${user.id}`).  When defining attribute mapping values in Terraform, the expression must be escaped (for example `value = \"$${user.id}}\"`)",
	)

	mappingTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the mapping type of the attribute. Options are `CORE`, `SCOPE`, and `CUSTOM`.",
	)

	//oidc block
	oidcScopesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"OIDC resource scope IDs that this attribute mapping is available for exclusively. This setting overrides any global OIDC resource scopes that contain an attribute mapping with the same name. The list can contain only scope IDs that have been granted for the application through the `/grants` endpoint. At least one scope ID is expected.",
	)

	oidcIdTokenEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Whether the attribute mapping should be available in the ID Token. This property is applicable only when the application's `protocol` property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.",
	)

	oidcUserinfoEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Whether the attribute mapping should be available through the `/as/userinfo` endpoint. This property is applicable only when the application's protocol property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.",
	)

	//saml
	samlsubjectNameformatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A URI reference representing the classification of the attribute, which helps the service provider interpret the attribute format.  This property is applicable only when the application's protocol property is `SAML` and the name is the `saml_subject` core attribute.  Examples include `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`, `urn:oasis:names:tc:SAML:2.0:attrname-format:uri`, `urn:oasis:names:tc:SAML:2.0:attrname-format:basic`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage custom attribute mappings for administrator defined applications configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the application attribute mapping in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to create the attribute mapping for."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"required": schema.BoolAttribute{
				Description:         requiredDescription.Description,
				MarkdownDescription: requiredDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

			"oidc_scopes": schema.SetAttribute{
				Description:         oidcScopesDescription.Description,
				MarkdownDescription: oidcScopesDescription.MarkdownDescription,
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.ValueStringsAre(verify.P1ResourceIDValidator()),
				},
			},

			"oidc_id_token_enabled": schema.BoolAttribute{
				Description:         oidcIdTokenEnabledDescription.Description,
				MarkdownDescription: oidcIdTokenEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(
						types.BoolValue(true),
						types.BoolValue(true),
						path.MatchRelative().AtParent().AtName("oidc_userinfo_enabled"),
					),
				},
			},

			"oidc_userinfo_enabled": schema.BoolAttribute{
				Description:         oidcUserinfoEnabledDescription.Description,
				MarkdownDescription: oidcUserinfoEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(
						types.BoolValue(true),
						types.BoolValue(true),
						path.MatchRelative().AtParent().AtName("oidc_id_token_enabled")),
				},
			},

			"saml_subject_nameformat": schema.StringAttribute{
				Description:         samlsubjectNameformatDescription.Description,
				MarkdownDescription: samlsubjectNameformatDescription.MarkdownDescription,
				Optional:            true,
			},
		},
	}
}

func (r *ApplicationAttributeMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *ApplicationAttributeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationAttributeMappingResourceModel

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

	var d diag.Diagnostics

	applicationType, d := plan.getApplicationType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(*applicationType)

	resp.Diagnostics.Append(plan.validate(applicationType, isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationAttributeMapping, d := plan.expand(ctx, r.client, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationAttributeMapping
	if !isCoreAttribute {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.CreateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
			},
			"CreateApplicationAttributeMapping",
			framework.CustomErrorInvalidValue,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), applicationAttributeMapping.GetId()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
			},
			"UpdateApplicationAttributeMapping",
			framework.CustomErrorInvalidValue,
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

func (r *ApplicationAttributeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationAttributeMappingResourceModel

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
	var response *management.ApplicationAttributeMapping
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplicationAttributeMapping",
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

func (r *ApplicationAttributeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationAttributeMappingResourceModel

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

	applicationType, d := plan.getApplicationType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(*applicationType)

	resp.Diagnostics.Append(plan.validate(applicationType, isCoreAttribute)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationAttributeMapping, d := plan.expand(ctx, r.client, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationAttributeMapping
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
		},
		"UpdateApplicationAttributeMapping",
		framework.CustomErrorInvalidValue,
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

func (r *ApplicationAttributeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationAttributeMappingResourceModel

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
	if data.MappingType.Equal(types.StringValue(string(management.ENUMRESOURCEATTRIBUTETYPE_CORE))) {

		applicationType, d := data.getApplicationType(ctx, r.client)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		coreAttributeData, ok := data.isCoreAttribute(*applicationType)
		if !ok {
			resp.Diagnostics.AddError(
				"Core attribute mismatch error",
				"The provider cannot determine the core attribute values to reset.  Please raise this issue to the provider maintainers.",
			)
			return
		}

		data.Value = framework.StringToTF(coreAttributeData.defaultValue)

		applicationAttributeMapping, d := data.expand(ctx, r.client, true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
			},
			"UpdateApplicationAttributeMapping",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)

	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				r, err := r.client.ApplicationAttributeMappingApi.DeleteApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
				return nil, r, err
			},
			"DeleteApplicationAttributeMapping",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationAttributeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/application_id/attribute_mapping_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *ApplicationAttributeMappingResourceModel) getApplicationType(ctx context.Context, apiClient *management.APIClient) (*management.EnumApplicationProtocol, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get application type and verify against the set params
	var respObject *management.ReadOneApplication200Response
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationsApi.ReadOneApplication(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()
		},
		"ReadOneApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&respObject,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	var applicationType *management.EnumApplicationProtocol
	if respObject.ApplicationOIDC != nil && respObject.ApplicationOIDC.GetId() != "" {
		applicationType = &respObject.ApplicationOIDC.Protocol
	} else if respObject.ApplicationSAML != nil && respObject.ApplicationSAML.GetId() != "" {
		applicationType = &respObject.ApplicationSAML.Protocol
	} else {
		diags.AddError(
			"Invalid parameter value - Unmappable application type",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `%s` or `%s` type.  Attributes cannot be mapped to this application.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, management.ENUMAPPLICATIONPROTOCOL_SAML),
		)
		return nil, diags
	}

	return applicationType, diags
}

func (p *ApplicationAttributeMappingResourceModel) validate(applicationType *management.EnumApplicationProtocol, isCoreAttribute bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if isCoreAttribute && p.Required.Equal(types.BoolValue(false)) {
		diags.AddError(
			"Invalid parameter value",
			fmt.Sprintf("The attribute name provided (%s) is a core attribute.  The `required` attribute must be set to `true`.", p.Name.ValueString()),
		)
		return diags
	}

	if *applicationType != management.ENUMAPPLICATIONPROTOCOL_SAML {

		// Check for SAML attributes that shouldn't be here
		if !p.SAMLSubjectNameformat.IsNull() {
			diags.AddError(
				"Invalid parameter value - Parameter doesn't apply to application type",
				fmt.Sprintf("The application ID provided (%s) is of type `%s`.  The `saml_subject_nameformat` attribute does not apply to this application type.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT),
			)
			return diags
		}
	}

	if *applicationType == management.ENUMAPPLICATIONPROTOCOL_SAML && !p.Name.Equal(types.StringValue("saml_subject")) {
		// Check for SAML attributes that shouldn't be here
		if !p.SAMLSubjectNameformat.IsNull() {
			diags.AddError(
				"Invalid parameter value - Attribute cannot be set",
				fmt.Sprintf("The attribute name provided (%s) is not a SAML subject attribute.  The `saml_subject_nameformat` attribute only applies to the `saml_subject` attribute.", p.Name.ValueString()),
			)
			return diags
		}
	}

	if *applicationType == management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT && isCoreAttribute {
		if (!p.OIDCScopes.IsNull() && !p.OIDCScopes.IsUnknown()) ||
			(!p.OIDCIDTokenEnabled.IsNull() && !p.OIDCIDTokenEnabled.IsUnknown()) ||
			(!p.OIDCUserinfoEnabled.IsNull() && !p.OIDCUserinfoEnabled.IsUnknown()) {
			diags.AddError(
				"Invalid parameter value - Parameter doesn't apply to core attributes",
				fmt.Sprintf("The application attribute name provided (%s) is a core attribute.  The `oidc_scopes`, `oidc_id_token_enabled` and `oidc_userinfo_enabled` attributes do not apply to this mapping type.", p.Name.ValueString()),
			)
			return diags
		}
	}

	if *applicationType != management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT {

		// If we have defined values that we shouldn't.  If unknown, deal with when the values are known.
		if (!p.OIDCScopes.IsNull() && !p.OIDCScopes.IsUnknown()) ||
			(!p.OIDCIDTokenEnabled.IsNull() && !p.OIDCIDTokenEnabled.IsUnknown()) ||
			(!p.OIDCUserinfoEnabled.IsNull() && !p.OIDCUserinfoEnabled.IsUnknown()) {
			diags.AddError(
				"Invalid parameter value - Parameter doesn't apply to application type",
				fmt.Sprintf("The application ID provided (%s) is of type `%s`.  The `oidc_scopes`, `oidc_id_token_enabled` and `oidc_userinfo_enabled` attributes do not apply to this application type.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_SAML),
			)
			return diags
		}

	}

	if *applicationType != management.ENUMAPPLICATIONPROTOCOL_SAML && *applicationType != management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT {
		diags.AddError(
			"Invalid parameter value - Unmappable application type",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `%s` or `%s` type.  Attributes cannot be mapped to this application.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, management.ENUMAPPLICATIONPROTOCOL_SAML),
		)
		return diags
	}

	return diags
}

func (p *ApplicationAttributeMappingResourceModel) isCoreAttribute(applicationType management.EnumApplicationProtocol) (*coreApplicationAttributeType, bool) {

	// Evaluate against the core attribute
	if v, ok := applicationCoreAttrMetadata[applicationType]; ok {
		// Loop the core attrs for the application type
		for _, coreAttr := range v {
			if strings.EqualFold(p.Name.ValueString(), coreAttr.name) {
				// We're a core attribute
				return &coreAttr, true
			}
		}
	}

	return nil, false
}

func (p *ApplicationAttributeMappingResourceModel) expand(ctx context.Context, apiClient *management.APIClient, overrideExisting bool) (*management.ApplicationAttributeMapping, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationAttributeMapping

	if overrideExisting {

		var response *management.EntityArray
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ApplicationAttributeMappingApi.ReadAllApplicationAttributeMappings(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()
			},
			"ReadAllApplicationAttributeMappings",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if diags.HasError() {
			return nil, diags
		}

		if attributes, ok := response.Embedded.GetAttributesOk(); ok {

			found := false
			for _, attribute := range attributes {

				if strings.EqualFold(attribute.ApplicationAttributeMapping.GetName(), p.Name.ValueString()) {
					data = attribute.ApplicationAttributeMapping
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

		if !p.SAMLSubjectNameformat.IsNull() {
			data.SetNameFormat(p.SAMLSubjectNameformat.ValueString())
		}

	} else {

		data = management.NewApplicationAttributeMapping(p.Name.ValueString(), p.Required.ValueBool(), p.Value.ValueString())

		if !p.Required.IsNull() && !p.Required.IsUnknown() {
			data.SetRequired(p.Required.ValueBool())
		} else {
			data.SetRequired(false)
		}

		if !p.OIDCScopes.IsNull() && !p.OIDCScopes.IsUnknown() {
			scopesSet, d := p.OIDCScopes.ToSetValue(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			scopesPointerSlice := framework.TFSetToStringSlice(ctx, scopesSet)

			if len(scopesPointerSlice) > 0 {
				scopesSlice := make([]string, 0)
				for i := range scopesPointerSlice {
					scopesSlice = append(scopesSlice, *scopesPointerSlice[i])
				}
				data.SetOidcScopes(scopesSlice)
			}
		}

		if !p.OIDCIDTokenEnabled.IsNull() && !p.OIDCIDTokenEnabled.IsUnknown() {
			data.SetIdToken(p.OIDCIDTokenEnabled.ValueBool())
		} else {
			data.SetIdToken(true)
		}

		if !p.OIDCUserinfoEnabled.IsNull() && !p.OIDCUserinfoEnabled.IsUnknown() {
			data.SetUserInfo(p.OIDCUserinfoEnabled.ValueBool())
		} else {
			data.SetUserInfo(true)
		}
	}

	return data, diags
}

func (p *ApplicationAttributeMappingResourceModel) toState(apiObject *management.ApplicationAttributeMapping) diag.Diagnostics {
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
	p.Required = framework.BoolOkToTF(apiObject.GetRequiredOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.MappingType = ApplicationAttributeMappingMappingTypeOkToTF(apiObject.GetMappingTypeOk())

	p.OIDCScopes = framework.StringSetOkToTF(apiObject.GetOidcScopesOk())
	p.OIDCIDTokenEnabled = framework.BoolOkToTF(apiObject.GetIdTokenOk())
	p.OIDCUserinfoEnabled = framework.BoolOkToTF(apiObject.GetUserInfoOk())

	return diags
}

func ApplicationAttributeMappingMappingTypeOkToTF(v *management.EnumAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
