package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type ApplicationAttributeMappingResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ApplicationAttributeMappingResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"name"`
	Required      types.Bool   `tfsdk:"required"`
	Value         types.String `tfsdk:"value"`
	MappingType   types.String `tfsdk:"mapping_type"`
	OIDCOptions   types.List   `tfsdk:"oidc_mapping_options"`
	SAMLOptions   types.List   `tfsdk:"saml_mapping_options"`
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
	coreValueNames = map[string]struct {
		applicationType string
		defaultValue    string
	}{
		"sub": {
			applicationType: string(management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT),
			defaultValue:    "${user.id}",
		},
		"saml_subject": {
			applicationType: string(management.ENUMAPPLICATIONPROTOCOL_SAML),
			defaultValue:    "${user.id}",
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

	reservedNames := []string{"acr", "amr", "at_hash", "aud", "auth_time", "azp", "client_id", "exp", "iat", "iss", "jti", "nbf", "nonce", "org", "scope", "sid", "sub"}
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

	//saml block
	samlsubjectNameformatDescriptionFmt := "A URI reference representing the classification of the attribute. Helps the service provider interpret the attribute format.  Options are `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`, `urn:oasis:names:tc:SAML:2.0:attrname-format:uri`, `urn:oasis:names:tc:SAML:2.0:attrname-format:basic`."
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
							_, isCoreValueNameState := coreValueNames[req.StateValue.ValueString()]
							_, isCoreValueNamePlan := coreValueNames[req.PlanValue.ValueString()]

							resp.RequiresReplace = isCoreValueNameState != isCoreValueNamePlan
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
		},

		Blocks: map[string]schema.Block{
			"oidc_mapping_options": schema.ListNestedBlock{
				Description: "A single block containing attribute mapping options specific to OpenID Connect applications.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"scopes": schema.SetAttribute{
							Description:         oidcScopesDescription.Description,
							MarkdownDescription: oidcScopesDescription.MarkdownDescription,
							Required:            true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(attrMinLength),
							},
						},

						"id_token_enabled": schema.BoolAttribute{
							Description:         oidcIdTokenEnabledDescription.Description,
							MarkdownDescription: oidcIdTokenEnabledDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(true),
						},

						"userinfo_enabled": schema.BoolAttribute{
							Description:         oidcUserinfoEnabledDescription.Description,
							MarkdownDescription: oidcUserinfoEnabledDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(true),
						},
					},
				},
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("saml_mapping_options")),
				},
			},

			"saml_mapping_options": schema.ListNestedBlock{
				Description: "A single block containing attribute mapping options specific to SAML applications.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"saml_subject_nameformat": schema.StringAttribute{
							Description:         samlsubjectNameformatDescription.Description,
							MarkdownDescription: samlsubjectNameformatDescription.MarkdownDescription,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("oidc_mapping_options")),
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

	var coreAttribute *management.ApplicationAttributeMapping
	var d diag.Diagnostics

	if v, ok := coreValueNames[plan.Name.ValueString()]; ok {
		coreAttribute, d = validateAttributeAgainstApplicationType(ctx, r.client, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Name.ValueString(), v.applicationType)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the model for the API
	applicationAttributeMapping := plan.expand(coreAttribute)

	var response interface{}

	if coreAttribute == nil {
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
				return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), coreAttribute.GetId()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
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

	var coreAttribute *management.ApplicationAttributeMapping
	if plan.MappingType.ValueString() == string(management.ENUMATTRIBUTEMAPPINGTYPE_CORE) {
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).Execute()
			},
			"ReadOneApplicationAttributeMapping",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		coreAttribute = response.(*management.ApplicationAttributeMapping)
	}

	// Build the model for the API
	applicationAttributeMapping := plan.expand(coreAttribute)

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

		coreAttribute := response.(*management.ApplicationAttributeMapping)

		if v, ok := coreValueNames[coreAttribute.Name]; ok {

			coreAttribute.SetValue(v.defaultValue)

			_, d = framework.ParseResponse(
				ctx,

				func() (interface{}, *http.Response, error) {
					return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).ApplicationAttributeMapping(*coreAttribute).Execute()
				},
				"UpdateApplicationAttributeMapping",
				framework.CustomErrorInvalidValue,
				sdk.DefaultCreateReadRetryable,
			)
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

func (p *ApplicationAttributeMappingResourceModel) expand(coreAttribute *management.ApplicationAttributeMapping) *management.ApplicationAttributeMapping {

	var data *management.ApplicationAttributeMapping

	if coreAttribute == nil {
		data = management.NewApplicationAttributeMapping(p.Name.ValueString(), p.Required.ValueBool(), p.Value.ValueString())
	} else {
		data.SetValue(p.Id.ValueString())
	}

	return data
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

	return diags
}

func ApplicationAttributeMappingMappingTypeOkToTF(v *management.EnumAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
