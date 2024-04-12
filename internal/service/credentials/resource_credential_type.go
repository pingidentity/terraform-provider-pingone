package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	customstringvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialTypeResource serviceClientType

type CredentialTypeResourceModel struct {
	Id                 types.String      `tfsdk:"id"`
	EnvironmentId      types.String      `tfsdk:"environment_id"`
	IssuerId           types.String      `tfsdk:"issuer_id"`
	CardType           types.String      `tfsdk:"card_type"`
	CardDesignTemplate types.String      `tfsdk:"card_design_template"`
	Description        types.String      `tfsdk:"description"`
	Metadata           types.Object      `tfsdk:"metadata"`
	RevokeOnDelete     types.Bool        `tfsdk:"revoke_on_delete"`
	Title              types.String      `tfsdk:"title"`
	CreatedAt          timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt          timetypes.RFC3339 `tfsdk:"updated_at"`
}

type MetadataModel struct {
	BackgroundImage  types.String `tfsdk:"background_image"`
	BgOpacityPercent types.Int64  `tfsdk:"bg_opacity_percent"`
	CardColor        types.String `tfsdk:"card_color"`
	Columns          types.Int64  `tfsdk:"columns"`
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"`
	LogoImage        types.String `tfsdk:"logo_image"`
	Name             types.String `tfsdk:"name"`
	Fields           types.List   `tfsdk:"fields"`
}

type FieldsModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Title       types.String `tfsdk:"title"`
	FileSupport types.String `tfsdk:"file_support"`
	IsVisible   types.Bool   `tfsdk:"is_visible"`
	Attribute   types.String `tfsdk:"attribute"`
	Value       types.String `tfsdk:"value"`
}

var (
	metadataServiceTFObjectTypes = map[string]attr.Type{
		"background_image":   types.StringType,
		"bg_opacity_percent": types.Int64Type,
		"card_color":         types.StringType,
		"columns":            types.Int64Type,
		"description":        types.StringType,
		"text_color":         types.StringType,
		"version":            types.Int64Type,
		"logo_image":         types.StringType,
		"name":               types.StringType,
		"fields":             types.ListType{ElemType: types.ObjectType{AttrTypes: innerFieldsServiceTFObjectTypes}},
	}

	innerFieldsServiceTFObjectTypes = map[string]attr.Type{
		"id":           types.StringType,
		"type":         types.StringType,
		"title":        types.StringType,
		"file_support": types.StringType,
		"is_visible":   types.BoolType,
		"attribute":    types.StringType,
		"value":        types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &CredentialTypeResource{}
	_ resource.ResourceWithConfigure   = &CredentialTypeResource{}
	_ resource.ResourceWithImportState = &CredentialTypeResource{}
)

// New Object
func NewCredentialTypeResource() resource.Resource {
	return &CredentialTypeResource{}
}

// Metadata
func (r *CredentialTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_type"
}

func (r *CredentialTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMinColumns = 1
	const attrMaxColumns = 3
	const attrMinPercent = 0
	const attrMaxPercent = 100

	titleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Title of the credential. Verification sites are expected to be able to request the issued credential from the compatible wallet app using the title.  This value aligns to `${cardTitle}` in the `card_design_template`.",
	)

	credentialDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A description of the credential type. This value aligns to `${cardSubtitle}` in the `card_design_template`.",
	)

	issuerIdDescriptipion := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of the issuer of the credential, which is the `id` of the `credential_issuer_profile` defined in the `environment`.",
	)

	revokeOnDeleteDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether a user's issued verifiable credentials are automatically revoked when a `credential_type`, `user`, or `environment` is deleted.",
	).DefaultValue("true")

	fieldsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"In a credential, the information is stored as key-value pairs where `fields` defines those key-value pairs. Effectively, `fields.title` is the key and its value is `fields.value` or extracted from the PingOne Directory attribute named in `fields.attribute`.",
	)

	fieldsIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifier of the field formatted as `<fields.type> -> <fields.title>`.",
	)

	fieldsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the type of data in the credential field.",
	).AllowedValuesEnum(credentials.AllowedEnumCredentialTypeMetaDataFieldsTypeEnumValues)

	fieldsFileSupportDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies how an image is stored in the credential field.",
	).AllowedValuesEnum(credentials.AllowedEnumCredentialTypeMetaDataFieldsFileSupportEnumValues)

	fieldsAttributeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the PingOne Directory attribute. Present if `field.type` is `Directory Attribute`.",
	)

	fieldsValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The text to appear on the credential for a `field.type` of `Alphanumeric Text`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the credential types used by compatible wallet applications.\n\n" +
			framework.SchemaAttributeDescriptionFromMarkdown("~> You must ensure that any fields used in the `card_design_template` are defined appropriately in `metadata.fields` or errors occur when you attempt to create a credential of that type.").MarkdownDescription,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential type exists."),
			),

			"issuer_id": schema.StringAttribute{
				Description:         issuerIdDescriptipion.Description,
				MarkdownDescription: issuerIdDescriptipion.MarkdownDescription,
				Computed:            true,
			},

			"title": schema.StringAttribute{
				Description:         titleDescription.Description,
				MarkdownDescription: titleDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					customstringvalidator.IsRequiredIfRegexMatchesPathValue(
						regexp.MustCompile(`\${cardTitle}`),
						framework.SchemaAttributeDescriptionFromMarkdown("The title argument is required because the ${cardTitle} element is defined in the `card_design_template`.").MarkdownDescription,
						path.MatchRoot("card_design_template"),
					),
					customstringvalidator.RegexMatchesPathValue(
						regexp.MustCompile(`\${cardTitle}`),
						"The title argument is defined but the card_design_template does not have a ${cardTitle} element.",
						path.MatchRoot("card_design_template"),
					),
				},
			},

			"description": schema.StringAttribute{
				Description:         credentialDescription.Description,
				MarkdownDescription: credentialDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					customstringvalidator.IsRequiredIfRegexMatchesPathValue(
						regexp.MustCompile(`\${cardSubtitle}`),
						framework.SchemaAttributeDescriptionFromMarkdown("The description argument is required because the ${cardSubtitle} element is defined in the `card_design_template`.").MarkdownDescription,
						path.MatchRoot("card_design_template"),
					),
					customstringvalidator.RegexMatchesPathValue(
						regexp.MustCompile(`\${cardSubtitle}`),
						"The description argument is defined but the card_design_template does not have a ${cardSubtitle} element.",
						path.MatchRoot("card_design_template"),
					),
				},
			},

			"card_type": schema.StringAttribute{
				Description: "A descriptor of the credential type. Can be non-identity types such as proof of employment or proof of insurance.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"card_design_template": schema.StringAttribute{
				Description: "An SVG formatted image containing placeholders for the credentials fields that need to be displayed in the image.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^<svg.*>[\s\S]*<\/svg>\s*$`), "expected value to contain a valid PingOne Credentials SVG card template."),
				},
			},

			"revoke_on_delete": schema.BoolAttribute{
				Description:         revokeOnDeleteDescription.Description,
				MarkdownDescription: revokeOnDeleteDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},

			"metadata": schema.SingleNestedAttribute{
				Description: "Contains the names, data types, and other metadata related to the credential.",
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"background_image": schema.StringAttribute{
						Description: "The URL or fully qualified path to the image file used for the credential background.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.  Image size must not exceed 50 KB.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${backgroundImage}`),
								"The metadata.background_image argument is required because the ${backgroundImage} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${backgroundImage}`),
								"The metadata.background_image argument is defined but the card_design_template does not have a ${backgroundImage} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"bg_opacity_percent": schema.Int64Attribute{
						Description: "A numnber indicating the percent opacity of the background image in the credential. High percentage opacity may make text on the credential difficult to read.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(attrMinPercent, attrMaxPercent),
						},
					},

					"card_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`),
								"expected value to contain a valid 6-digit hexadecimal color code, prefixed with a hash (#) symbol."),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${cardColor}`),
								"The metadata.card_color argument is required because the ${cardColor} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${cardColor}`),
								"The metadata.card_color argument is defined but the card_design_template does not have a ${cardColor} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"columns": schema.Int64Attribute{
						Description: "Indicates a number (between 1-3) of columns to display visible fields on the credential.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(attrMinColumns, attrMaxColumns),
						},
					},

					"description": schema.StringAttribute{
						Description: "Description of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"logo_image": schema.StringAttribute{
						Description: "The URL or fully qualified path to the image file used for the credential logo.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.  Image size must not exceed 25 KB.",
						Optional:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${logoImage}`),
								"The metadata.logo_image argument is required because the ${logoImage} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${logoImage}`),
								"The metadata.logo_image argument is defined but the card_design_template does not have a ${logoImage} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"name": schema.StringAttribute{
						Description: "Name of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"text_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential text.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(
									`^#([A-Fa-f0-9]{6})$`),
								"expected value to contain a valid 6-digit hexadecimal color code, prefixed with a hash (#) symbol."),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${textColor}`),
								"The metadata.text_color argument is required because the ${textColor} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${textColor}`),
								"The metadata.text_color argument is defined but the card_design_template does not have a ${textColor} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},
					"version": schema.Int64Attribute{
						Description: "Number version of this credential.",
						Computed:    true,
					},

					"fields": schema.ListNestedAttribute{
						Description:         fieldsDescription.Description,
						MarkdownDescription: fieldsDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(attrMinLength),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         fieldsIdDescription.Description,
									MarkdownDescription: fieldsIdDescription.MarkdownDescription,
									Computed:            true,
								},
								"type": schema.StringAttribute{
									Description:         fieldsTypeDescription.Description,
									MarkdownDescription: fieldsTypeDescription.MarkdownDescription,
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(credentials.AllowedEnumCredentialTypeMetaDataFieldsTypeEnumValues)...),
									},
								},
								"title": schema.StringAttribute{
									Description: "Descriptive text when showing the field.",
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
									},
								},
								"file_support": schema.StringAttribute{
									Description:         fieldsFileSupportDescription.Description,
									MarkdownDescription: fieldsFileSupportDescription.MarkdownDescription,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.All(
											customstringvalidator.RegexMatchesPathValue(
												regexp.MustCompile(string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE)),
												fmt.Sprintf("The fields.file_support argument is only applicable when fields.type has a value of %s.", string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE)),
												path.MatchRelative().AtParent().AtName("type"),
											),
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(credentials.AllowedEnumCredentialTypeMetaDataFieldsFileSupportEnumValues)...),
										),
									},
								},
								"is_visible": schema.BoolAttribute{
									Description: "Specifies whether the field should be visible to viewers of the credential.",
									Optional:    true,
								},
								"attribute": schema.StringAttribute{
									Description:         fieldsAttributeDescription.Description,
									MarkdownDescription: fieldsAttributeDescription.MarkdownDescription,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("value")),
										customstringvalidator.IsRequiredIfMatchesPathValue(basetypes.NewStringValue(string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE)), path.MatchRelative().AtParent().AtName("type")),
									},
								},
								"value": schema.StringAttribute{
									Description:         fieldsValueDescription.Description,
									MarkdownDescription: fieldsValueDescription.MarkdownDescription,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("attribute")),
										customstringvalidator.IsRequiredIfMatchesPathValue(basetypes.NewStringValue(string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT)), path.MatchRelative().AtParent().AtName("type")),
									},
								},
							},
						},
					},
				},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the object was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the object was updated. Can be null.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CredentialTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialTypeResourceModel

	if r.Client.CredentialsAPIClient == nil {
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
	credentialType, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	timeoutValue := 15

	var response *credentials.CredentialType
	resp.Diagnostics.Append(framework.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.CreateCredentialType(ctx, plan.EnvironmentId.ValueString()).CredentialType(*credentialType).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateCredentialType",
		framework.DefaultCustomError,
		credentialTypeRetryConditions,
		&response,
		time.Duration(timeoutValue)*time.Minute, // 15 mins
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

func (r *CredentialTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialTypeResourceModel

	if r.Client.CredentialsAPIClient == nil {
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
	var response *credentials.CredentialType
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.ReadOneCredentialType(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneCredentialType",
		framework.CustomErrorResourceNotFoundWarning,
		credentialTypeRetryConditions,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found or soft deleted
	if response == nil || response.DeletedAt != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialTypeResourceModel

	if r.Client.CredentialsAPIClient == nil {
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
	credentialType, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.CredentialType
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.UpdateCredentialType(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).CredentialType(*credentialType).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateCredentialType",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CredentialTypeResourceModel

	if r.Client.CredentialsAPIClient == nil {
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
			fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.DeleteCredentialType(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteCredentialType",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CredentialTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "credential_type_id",
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

func (p *CredentialTypeResourceModel) expand(ctx context.Context) (*credentials.CredentialType, diag.Diagnostics) {
	var diags diag.Diagnostics

	credentialTypeMetaData := credentials.NewCredentialTypeMetaData()

	// expand credential type metadata and metadata.fields
	if !p.Metadata.IsNull() && !p.Metadata.IsUnknown() {
		var metadata MetadataModel
		d := p.Metadata.As(ctx, &metadata, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		credentialTypeMetaData, d = metadata.expandMetaDataModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	data := credentials.NewCredentialType(p.CardDesignTemplate.ValueString(), *credentialTypeMetaData, p.Title.ValueString())

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.CardType.IsNull() && !p.CardType.IsUnknown() {
		data.SetCardType(p.CardType.ValueString())
	}

	if !p.RevokeOnDelete.IsNull() && !p.RevokeOnDelete.IsUnknown() {
		onDeleteObject := credentials.NewCredentialTypeOnDelete()
		onDeleteObject.SetRevokeIssuedCredentials(p.RevokeOnDelete.ValueBool())

		data.SetOnDelete(*onDeleteObject)
	}

	return data, diags
}

func (p *MetadataModel) expandMetaDataModel(ctx context.Context) (*credentials.CredentialTypeMetaData, diag.Diagnostics) {
	var diags diag.Diagnostics

	cardMetadata := credentials.NewCredentialTypeMetaDataWithDefaults()

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		cardMetadata.SetName(p.Name.ValueString())
	}

	if !p.BackgroundImage.IsNull() && !p.BackgroundImage.IsUnknown() {
		cardMetadata.SetBackgroundImage(p.BackgroundImage.ValueString())
	}

	if !p.BgOpacityPercent.IsNull() && !p.BgOpacityPercent.IsUnknown() {
		cardMetadata.SetBgOpacityPercent(int32(p.BgOpacityPercent.ValueInt64()))
	}

	if !p.CardColor.IsNull() && !p.CardColor.IsUnknown() {
		cardMetadata.SetCardColor(p.CardColor.ValueString())
	}

	if !p.Columns.IsNull() && !p.Columns.IsUnknown() {
		cardMetadata.SetColumns(int32(p.Columns.ValueInt64()))
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		cardMetadata.SetDescription(p.Description.ValueString())
	}

	if !p.LogoImage.IsNull() && !p.LogoImage.IsUnknown() {
		cardMetadata.SetLogoImage(p.LogoImage.ValueString())
	}

	if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
		cardMetadata.SetTextColor(p.TextColor.ValueString())
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		cardMetadata.SetVersion(int32(p.Version.ValueInt64()))
	}

	// expand fields
	if !p.Fields.IsNull() && !p.Fields.IsUnknown() {
		var innerFields []FieldsModel
		fields := make([]credentials.CredentialTypeMetaDataFieldsInner, 0)
		diags.Append(p.Fields.ElementsAs(ctx, &innerFields, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, v := range innerFields {
			field, d := v.expandFields()
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			fields = append(fields, *field)
		}
		// complete the meta data object
		cardMetadata.SetFields(fields)
	}

	return cardMetadata, diags
}

func (p *FieldsModel) expandFields() (*credentials.CredentialTypeMetaDataFieldsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	innerFields := credentials.NewCredentialTypeMetaDataFieldsInnerWithDefaults()

	attrType := credentials.EnumCredentialTypeMetaDataFieldsType(p.Type.ValueString())
	innerFields.SetType(attrType)

	attrId := p.Type.ValueString() + " -> " + p.Title.ValueString() // construct id per P1Creds API recommendations
	innerFields.SetId(attrId)

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT {
		innerFields.SetValue(p.Value.ValueString())
	}

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE {
		innerFields.SetAttribute(p.Attribute.ValueString())

		if !p.FileSupport.IsNull() && !p.FileSupport.IsUnknown() {
			innerFields.SetFileSupport(credentials.EnumCredentialTypeMetaDataFieldsFileSupport(p.FileSupport.ValueString()))
		}

	}

	if !p.Title.IsNull() && !p.Title.IsUnknown() {
		innerFields.SetTitle(p.Title.ValueString())
	}

	if !p.IsVisible.IsNull() && !p.IsVisible.IsUnknown() {
		innerFields.SetIsVisible(p.IsVisible.ValueBool())
	}

	if innerFields == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Metadata.Fields object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}

	return innerFields, diags
}

func (p *CredentialTypeResourceModel) toState(apiObject *credentials.CredentialType) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	// credential attributes
	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.IssuerId = framework.StringToTF(*apiObject.GetIssuer().Id)
	p.Title = framework.StringOkToTF(apiObject.GetTitleOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.CardType = framework.StringOkToTF(apiObject.GetCardTypeOk())
	p.CardDesignTemplate = framework.StringOkToTF(apiObject.GetCardDesignTemplateOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	revokeOnDelete := types.BoolNull()
	if v, ok := apiObject.GetOnDeleteOk(); ok {
		revokeOnDelete = framework.BoolOkToTF(v.GetRevokeIssuedCredentialsOk())
	}
	p.RevokeOnDelete = revokeOnDelete

	// credential metadata
	metadata, d := toStateMetadata(apiObject.GetMetadataOk())
	diags.Append(d...)
	p.Metadata = metadata

	return diags
}

func toStateMetadata(metadata *credentials.CredentialTypeMetaData, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// core metadata object
	metadataMap := map[string]attr.Value{
		"background_image":   framework.StringOkToTF(metadata.GetBackgroundImageOk()),
		"bg_opacity_percent": framework.Int32OkToTF(metadata.GetBgOpacityPercentOk()),
		"card_color":         framework.StringOkToTF(metadata.GetCardColorOk()),
		"columns":            framework.Int32OkToTF(metadata.GetColumnsOk()),
		"description":        framework.StringOkToTF(metadata.GetDescriptionOk()),
		"text_color":         framework.StringOkToTF(metadata.GetTextColorOk()),
		"version":            framework.Int32OkToTF(metadata.GetVersionOk()),
		"logo_image":         framework.StringOkToTF(metadata.GetLogoImageOk()),
		"name":               framework.StringOkToTF(metadata.GetNameOk()),
	}

	// metadata fields objects
	fields, d := toStateFields(metadata.GetFieldsOk())
	diags.Append(d...)

	metadataMap["fields"] = fields
	flattenedObj, d := types.ObjectValue(metadataServiceTFObjectTypes, metadataMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFields(innerFields []credentials.CredentialTypeMetaDataFieldsInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfInnerObjType := types.ObjectType{AttrTypes: innerFieldsServiceTFObjectTypes}
	innerflattenedList := []attr.Value{}
	for _, v := range innerFields {

		fieldsMap := map[string]attr.Value{
			"id":           framework.StringOkToTF(v.GetIdOk()),
			"type":         framework.EnumOkToTF(v.GetTypeOk()),
			"title":        framework.StringOkToTF(v.GetTitleOk()),
			"file_support": framework.EnumOkToTF(v.GetFileSupportOk()),
			"is_visible":   framework.BoolOkToTF(v.GetIsVisibleOk()),
			"attribute":    framework.StringOkToTF(v.GetAttributeOk()),
			"value":        framework.StringOkToTF(v.GetValueOk()),
		}
		innerflattenedObj, d := types.ObjectValue(innerFieldsServiceTFObjectTypes, fieldsMap)
		diags.Append(d...)

		innerflattenedList = append(innerflattenedList, innerflattenedObj)
	}
	fields, d := types.ListValue(tfInnerObjType, innerflattenedList)
	diags.Append(d...)

	return fields, diags
}

func credentialTypeRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Credential Issuer Profile's keys may not have propagated yet. Rare, but possible.
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			// detected issuer profile not fully deployed yet
			if m, err := regexp.MatchString("^issuerProfile must exist before creating credentialTypes", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("IssuerProfile (prerequisite) has not finished provisioning - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}
		}

		// detected credentials service not fully deployed yet
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}
	}

	return false
}
