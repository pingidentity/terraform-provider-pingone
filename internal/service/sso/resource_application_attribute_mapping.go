package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

type coreAttributeType struct {
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
	applicationCoreAttrMetadata = map[string][]coreAttributeType{
		string(management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT): {
			{
				name:         "sub",
				defaultValue: "${user.id}",
			},
		},
		string(management.ENUMAPPLICATIONPROTOCOL_SAML): {
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
	nameDescriptionFmt := fmt.Sprintf("A string that specifies the name of attribute and must be unique within an application. For SAML applications, the `samlAssertion.subject` name is a reserved case-insensitive name which indicates the mapping to be used for the subject in an assertion. For OpenID Connect applications, the following names are reserved and cannot be used: `%s`.", strings.Join(reservedNames, "`, `"))
	nameDescription := framework.SchemaDescription{
		MarkdownDescription: nameDescriptionFmt,
		Description:         strings.ReplaceAll(nameDescriptionFmt, "`", "\""),
	}

	requiredDescriptionFmt := "A boolean to specify whether a mapping value is required for this attribute. If `true`, a value must be set and a non-empty value must be available in the SAML assertion or ID token. Defaults to `false`."
	requiredDescription := framework.SchemaDescription{
		MarkdownDescription: requiredDescriptionFmt,
		Description:         strings.ReplaceAll(requiredDescriptionFmt, "`", "\""),
	}

	valueDescriptionFmt := "A string that specifies the string constants or expression for mapping the attribute path against a specific source. The expression format is `${<source>.<attribute_path>}`. The only supported source is user (for example, `${user.id}`)."
	valueDescription := framework.SchemaDescription{
		MarkdownDescription: valueDescriptionFmt,
		Description:         strings.ReplaceAll(valueDescriptionFmt, "`", "\""),
	}

	mappingTypeDescriptionFmt := "A string that specifies the mapping type of the attribute. Options are `CORE`, `SCOPE`, and `CUSTOM`."
	mappingTypeDescription := framework.SchemaDescription{
		MarkdownDescription: mappingTypeDescriptionFmt,
		Description:         strings.ReplaceAll(mappingTypeDescriptionFmt, "`", "\""),
	}

	//oidc block
	oidcScopesDescriptionFmt := "OIDC resource scope IDs that this attribute mapping is available for exclusively. This setting overrides any global OIDC resource scopes that contain an attribute mapping with the same name. The list can contain only scope IDs that have been granted for the application through the `/grants` endpoint. At least one scope ID is expected."
	oidcScopesDescription := framework.SchemaDescription{
		MarkdownDescription: oidcScopesDescriptionFmt,
		Description:         strings.ReplaceAll(oidcScopesDescriptionFmt, "`", "\""),
	}

	oidcIdTokenEnabledDescriptionFmt := "Whether the attribute mapping should be available in the ID Token. This property is applicable only when the application's `protocol` property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`."
	oidcIdTokenEnabledDescription := framework.SchemaDescription{
		MarkdownDescription: oidcIdTokenEnabledDescriptionFmt,
		Description:         strings.ReplaceAll(oidcIdTokenEnabledDescriptionFmt, "`", "\""),
	}

	oidcuserinfoEnabledDescriptionFmt := "Whether the attribute mapping should be available through the `/as/userinfo` endpoint. This property is applicable only when the application's protocol property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`."
	oidcUserinfoEnabledDescription := framework.SchemaDescription{
		MarkdownDescription: oidcuserinfoEnabledDescriptionFmt,
		Description:         strings.ReplaceAll(oidcuserinfoEnabledDescriptionFmt, "`", "\""),
	}

	//saml
	samlsubjectNameformatDescriptionFmt := "A URI reference representing the classification of the attribute, which helps the service provider interpret the attribute format.  This property is applicable only when the application's protocol property is `SAML`.  Options are `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`, `urn:oasis:names:tc:SAML:2.0:attrname-format:uri`, `urn:oasis:names:tc:SAML:2.0:attrname-format:basic`."
	samlsubjectNameformatDescription := framework.SchemaDescription{
		MarkdownDescription: samlsubjectNameformatDescriptionFmt,
		Description:         strings.ReplaceAll(samlsubjectNameformatDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for applications configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the application attribute mapping in."},
			),

			"application_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the application to create the attribute mapping for."},
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							for _, v := range applicationCoreAttrMetadata {
								// Loop the core attrs for the application type
								for _, coreAttr := range v {
									if strings.ToUpper(req.StateValue.ValueString()) == strings.ToUpper(coreAttr.name) ||
										strings.ToUpper(req.PlanValue.ValueString()) == strings.ToUpper(coreAttr.name) {
										// State is a core attribute
										resp.RequiresReplace = true
										return
									}
								}
							}
						},
						"The resource must be replaced if changing between core and custom attribute types.",
						"The resource must be replaced if changing between core and custom attribute types.",
					),
				},
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
					setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("saml_subject_nameformat")),
					setvalidator.ValueStringsAre(verify.P1ResourceIDValidator()),
				},
			},

			"oidc_id_token_enabled": schema.BoolAttribute{
				Description:         oidcIdTokenEnabledDescription.Description,
				MarkdownDescription: oidcIdTokenEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(path.MatchRelative().AtParent().AtName("oidc_userinfo_enabled")),
					boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("saml_subject_nameformat")),
					boolvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("oidc_scopes")),
				},
			},

			"oidc_userinfo_enabled": schema.BoolAttribute{
				Description:         oidcUserinfoEnabledDescription.Description,
				MarkdownDescription: oidcUserinfoEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				Validators: []validator.Bool{
					customboolvalidator.AtLeastOneOfMustBeTrue(path.MatchRelative().AtParent().AtName("oidc_id_token_enabled")),
					boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("saml_subject_nameformat")),
					boolvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("oidc_scopes")),
				},
			},

			"saml_subject_nameformat": schema.StringAttribute{
				Description:         samlsubjectNameformatDescription.Description,
				MarkdownDescription: samlsubjectNameformatDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("oidc_scopes")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("oidc_id_token_enabled")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("oidc_userinfo_enabled")),
				},
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

func (r *ApplicationAttributeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationAttributeMappingResourceModel

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

	var d diag.Diagnostics

	applicationType, d := plan.getApplicationType(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.validate(applicationType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, isCoreAttribute := plan.isCoreAttribute(*applicationType)

	// Build the model for the API
	applicationAttributeMapping, d := plan.expand(ctx, r.client, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response interface{}

	if !isCoreAttribute {
		// Run the API call
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.CreateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
			},
			"CreateApplicationAttributeMapping",
			framework.CustomErrorInvalidValue,
			sdk.DefaultCreateReadRetryable,
		)
	} else {
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), applicationAttributeMapping.GetId()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
			},
			"UpdateApplicationAttributeMapping",
			framework.CustomErrorInvalidValue,
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
	resp.Diagnostics.Append(state.toState(response.(*management.ApplicationAttributeMapping))...)
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
			return r.client.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplicationAttributeMapping",
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
	resp.Diagnostics.Append(data.toState(response.(*management.ApplicationAttributeMapping))...)
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

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isCoreAttribute := plan.MappingType.ValueString() == string(management.ENUMATTRIBUTEMAPPINGTYPE_CORE)

	// Build the model for the API
	applicationAttributeMapping, d := plan.expand(ctx, r.client, isCoreAttribute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
		},
		"UpdateApplicationAttributeMapping",
		framework.CustomErrorInvalidValue,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ApplicationAttributeMapping))...)
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
	if data.MappingType.ValueString() == string(management.ENUMATTRIBUTEMAPPINGTYPE_CORE) {

		applicationType, d := data.getApplicationType(ctx, r.client)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		if coreAttributeData, ok := data.isCoreAttribute(*applicationType); ok {

			data.Value = framework.StringToTF(coreAttributeData.defaultValue)

			applicationAttributeMapping, d := data.expand(ctx, r.client, true)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}

			_, d = framework.ParseResponse(
				ctx,

				func() (interface{}, *http.Response, error) {
					return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
				},
				"UpdateApplicationAttributeMapping",
				framework.CustomErrorResourceNotFoundWarning,
				sdk.DefaultCreateReadRetryable,
			)
		} else {

		}

	} else {
		_, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				r, err := r.client.ApplicationAttributeMappingApi.DeleteApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
				return nil, r, err
			},
			"DeleteApplicationAttributeMapping",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
		)
	}
	resp.Diagnostics.Append(d...)
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
	resp, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApi.ReadOneApplication(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()
		},
		"ReadOneApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	respObject := resp.(*management.ReadOneApplication200Response)

	var applicationType *management.EnumApplicationProtocol
	if respObject.ApplicationOIDC != nil && respObject.ApplicationOIDC.GetId() != "" {
		applicationType = &respObject.ApplicationOIDC.Protocol
	} else if respObject.ApplicationSAML != nil && respObject.ApplicationSAML.GetId() != "" {
		applicationType = &respObject.ApplicationSAML.Protocol
	} else {
		diags.AddError(
			"Invalid parameter value",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `%s` or `%s` type.  Attributes cannot be mapped to this application.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, management.ENUMAPPLICATIONPROTOCOL_SAML),
		)
		return nil, diags
	}

	return applicationType, diags
}

func (p *ApplicationAttributeMappingResourceModel) validate(applicationType *management.EnumApplicationProtocol) diag.Diagnostics {
	var diags diag.Diagnostics

	if *applicationType == management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT {

		if !p.SAMLSubjectNameformat.IsNull() {
			diags.AddError(
				"Invalid parameter value",
				fmt.Sprintf("The application ID provided (%s) is of type `%s`.  The `saml_subject_nameformat` attribute does not apply to this application type.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT),
			)
			return diags
		}

	} else if *applicationType == management.ENUMAPPLICATIONPROTOCOL_SAML {

		if !p.OIDCScopes.IsNull() || !p.OIDCIDTokenEnabled.IsNull() || !p.OIDCUserinfoEnabled.IsNull() {
			diags.AddError(
				"Invalid parameter value",
				fmt.Sprintf("The application ID provided (%s) is of type `%s`.  The `oidc_scopes`, `oidc_id_token_enabled` and `oidc_userinfo_enabled` attributes do not apply to this application type.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_SAML),
			)
			return diags
		}

	} else {
		diags.AddError(
			"Invalid parameter value",
			fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `%s` or `%s` type.  Attributes cannot be mapped to this application.", p.ApplicationId.ValueString(), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, management.ENUMAPPLICATIONPROTOCOL_SAML),
		)
		return diags
	}

	return diags
}

func (p *ApplicationAttributeMappingResourceModel) isCoreAttribute(applicationType management.EnumApplicationProtocol) (*coreAttributeType, bool) {

	// Evaluate against the core attribute
	if v, ok := applicationCoreAttrMetadata[string(applicationType)]; ok {
		// Loop the core attrs for the application type
		for _, coreAttr := range v {
			if strings.ToUpper(p.Name.ValueString()) == strings.ToUpper(coreAttr.name) {
				// We're a core attribute
				return &coreAttr, true
			}
		}
	}

	return nil, false
}

func (p *ApplicationAttributeMappingResourceModel) expand(ctx context.Context, apiClient *management.APIClient, isCoreAttribute bool) (*management.ApplicationAttributeMapping, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationAttributeMapping

	if !isCoreAttribute {
		data = management.NewApplicationAttributeMapping(p.Name.ValueString(), p.Required.ValueBool(), p.Value.ValueString())

		if !p.OIDCScopes.IsNull() {
			scopesSet, d := p.OIDCScopes.ToSetValue(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			scopesPointerSlice := framework.TFSetToStringSlice(ctx, scopesSet)
			scopesSlice := make([]string, 0)
			for i := range scopesPointerSlice {
				scopesSlice = append(scopesSlice, *scopesPointerSlice[i])
			}
			data.SetOidcScopes(scopesSlice)
		}

		if !p.OIDCIDTokenEnabled.IsNull() {
			data.SetIdToken(p.OIDCIDTokenEnabled.ValueBool())
		}

		if !p.OIDCUserinfoEnabled.IsNull() {
			data.SetUserInfo(p.OIDCUserinfoEnabled.ValueBool())
		}

		if !p.SAMLSubjectNameformat.IsNull() {
			data.SetNameFormat(p.SAMLSubjectNameformat.ValueString())
		}

	} else {

		// fetch the attribute that already exists
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ApplicationAttributeMappingApi.ReadAllApplicationAttributeMappings(ctx, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString()).Execute()
			},
			"ReadAllApplicationAttributeMappings",
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

				if strings.ToUpper(attribute.ApplicationAttributeMapping.GetName()) == strings.ToUpper(p.Name.ValueString()) {
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
	p.SAMLSubjectNameformat = framework.StringOkToTF(apiObject.GetNameFormatOk())

	return diags
}

func ApplicationAttributeMappingMappingTypeOkToTF(v *management.EnumAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
