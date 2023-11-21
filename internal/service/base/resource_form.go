package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type FormResource serviceClientType

type formResourceModel struct {
	Id                types.String `tfsdk:"id"`
	EnvironmentId     types.String `tfsdk:"environment_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Category          types.String `tfsdk:"category"`
	Cols              types.Int64  `tfsdk:"cols"`
	Components        types.Object `tfsdk:"components"`
	FieldTypes        types.Set    `tfsdk:"field_types"`
	LanguageBundle    types.Map    `tfsdk:"language_bundle"`
	MarkOptional      types.Bool   `tfsdk:"mark_optional"`
	MarkRequired      types.Bool   `tfsdk:"mark_required"`
	TranslationMethod types.String `tfsdk:"translation_method"`
}

type formComponentsResourceModel struct {
	Fields types.Set `tfsdk:"fields"`
}

type formComponentsFieldResourceModel struct {
	Position types.Object `tfsdk:"position"`
	Type     types.String `tfsdk:"type"`
	// The form fields
	FieldText              types.Object `tfsdk:"field_text"`
	FieldPassword          types.Object `tfsdk:"field_password"`
	FieldPasswordVerify    types.Object `tfsdk:"field_password_verify"`
	FieldRadio             types.Object `tfsdk:"field_radio"`
	FieldCheckbox          types.Object `tfsdk:"field_checkbox"`
	FieldDropdown          types.Object `tfsdk:"field_dropdown"`
	FieldCombobox          types.Object `tfsdk:"field_combobox"`
	FieldDivider           types.Object `tfsdk:"field_divider"`
	FieldEmptyField        types.Object `tfsdk:"field_empty_field"`
	FieldTextblob          types.Object `tfsdk:"field_textblob"`
	FieldSlateTextblob     types.Object `tfsdk:"field_slate_textblob"`
	FieldSubmitButton      types.Object `tfsdk:"field_submit_button"`
	FieldErrorDisplay      types.Object `tfsdk:"field_error_display"`
	FieldFlowLink          types.Object `tfsdk:"field_flow_link"`
	FieldFlowButton        types.Object `tfsdk:"field_flow_button"`
	FieldRecaptchaV2       types.Object `tfsdk:"field_recaptcha_v2"`
	FieldQrCode            types.Object `tfsdk:"field_qr_code"`
	FieldSocialLoginButton types.Object `tfsdk:"field_social_login_button"`
}

type formComponentsFieldPositionResourceModel struct {
	Col   types.Int64 `tfsdk:"col"`
	Row   types.Int64 `tfsdk:"row"`
	Width types.Int64 `tfsdk:"width"`
}

// TEXT, PASSWORD, RADIO, CHECKBOX, DROPDOWN
type formComponentsFieldElementResourceModel struct {
	AttributeDisabled            types.Bool   `tfsdk:"attribute_disabled"`
	Key                          types.String `tfsdk:"key"`
	LabelMode                    types.String `tfsdk:"label_mode"`
	Layout                       types.String `tfsdk:"layout"`
	Options                      types.Set    `tfsdk:"options"`
	Required                     types.Bool   `tfsdk:"required"`
	Validation                   types.Object `tfsdk:"validation"`
	OtherOptionEnabled           types.Bool   `tfsdk:"other_option_enabled"`
	OtherOptionKey               types.String `tfsdk:"other_option_key"`
	OtherOptionLabel             types.String `tfsdk:"other_option_label"`
	OtherOptionInputLabel        types.String `tfsdk:"other_option_input_label"`
	OtherOptionAttributeDisabled types.Bool   `tfsdk:"other_option_attribute_disabled"`
}

type formComponentsFieldTextResourceModel formComponentsFieldElementResourceModel
type formComponentsFieldPasswordResourceModel formComponentsFieldElementResourceModel
type formComponentsFieldRadioResourceModel formComponentsFieldElementResourceModel
type formComponentsFieldCheckboxResourceModel formComponentsFieldElementResourceModel
type formComponentsFieldDropdownResourceModel formComponentsFieldElementResourceModel

type formComponentsFieldElementOptionResourceModel struct {
	Label types.String `tfsdk:"label"`
	Value types.String `tfsdk:"value"`
}

type formComponentsFieldElementValidationResourceModel struct {
	Regex        types.String `tfsdk:"regex"`
	Type         types.String `tfsdk:"type"`
	ErrorMessage types.String `tfsdk:"error_message"`
}

// DIVIDER, PARAGRAPH, EMPTY_FIELD, ERROR_DISPLAY, (TEXTBLOB, SLATE_TEXTBLOB)?
type formComponentsFieldItemResourceModel struct {
	Content types.String `tfsdk:"content"`
}

type formComponentsFieldDividerResourceModel formComponentsFieldItemResourceModel
type formComponentsFieldParagraphResourceModel formComponentsFieldItemResourceModel
type formComponentsFieldEmptyFieldResourceModel formComponentsFieldItemResourceModel
type formComponentsFieldErrorDisplayResourceModel formComponentsFieldItemResourceModel
type formComponentsFieldTextblobResourceModel formComponentsFieldItemResourceModel
type formComponentsFieldSlateTextblobResourceModel formComponentsFieldItemResourceModel

// PASSWORD_VERIFY
type formComponentsFieldPasswordVerifyResourceModel struct {
	LabelPasswordVerify types.String `tfsdk:"label_password_verify"`
}

// SUBMIT_BUTTON, FLOW_BUTTON
type formComponentsFieldButtonResourceModel struct {
	Key    types.String `tfsdk:"key"`
	Label  types.String `tfsdk:"label"`
	Styles types.Object `tfsdk:"styles"`
}

type formComponentsFieldSubmitButtonResourceModel formComponentsFieldButtonResourceModel
type formComponentsFieldFlowButtonResourceModel formComponentsFieldButtonResourceModel

type formComponentsFieldButtonStylesResourceModel struct {
	Width           types.Int64  `tfsdk:"width"`
	Alignment       types.String `tfsdk:"alignment"`
	BackgroundColor types.String `tfsdk:"background_color"`
	TextColor       types.String `tfsdk:"text_color"`
	BorderColor     types.String `tfsdk:"border_color"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// FLOW_LINK
type formComponentsFieldFlowLinkResourceModel struct {
	Key    types.String `tfsdk:"key"`
	Label  types.String `tfsdk:"label"`
	Styles types.Object `tfsdk:"styles"`
}

type formComponentsFieldFlowLinkStylesResourceModel struct {
	HorizontalAlignment types.String `tfsdk:"horizontal_alignment"`
	TextColor           types.String `tfsdk:"text_color"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}

// RECAPTCHA_V2
type formComponentsFieldRecaptchaV2ResourceModel struct {
	Key       types.String `tfsdk:"key"`
	Size      types.String `tfsdk:"size"`
	Theme     types.String `tfsdk:"theme"`
	Alignment types.String `tfsdk:"alignment"`
}

// QR_CODE
type formComponentsFieldQrCodeResourceModel struct {
	QrCodeType types.String `tfsdk:"qr_code_type"`
	Alignment  types.String `tfsdk:"alignment"`
	ShowBorder types.Bool   `tfsdk:"show_border"`
}

// SOCIAL_LOGIN_BUTTON
type formComponentsFieldSocialLoginButtonResourceModel struct {
	Label      types.String `tfsdk:"label"`
	Styles     types.Object `tfsdk:"styles"`
	IdpType    types.String `tfsdk:"idp_type"`
	IdpName    types.String `tfsdk:"idp_name"`
	IdpId      types.String `tfsdk:"idp_id"`
	IdpEnabled types.Bool   `tfsdk:"idp_enabled"`
	IconSrc    types.String `tfsdk:"icon_src"`
	Width      types.Int64  `tfsdk:"width"`
}

type formComponentsFieldSocialLoginButtonStylesResourceModel struct {
	HorizontalAlignment types.String `tfsdk:"horizontal_alignment"`
	TextColor           types.String `tfsdk:"text_color"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}

// Framework interfaces
var (
	_ resource.Resource                = &FormResource{}
	_ resource.ResourceWithConfigure   = &FormResource{}
	_ resource.ResourceWithImportState = &FormResource{}
)

// New Object
func NewFormResource() resource.Resource {
	return &FormResource{}
}

// Metadata
func (r *FormResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form"
}

// Schema.
func (r *FormResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const colsMinValue = 0
	const colsMaxValue = 4
	const rowMaxValue = 50

	formFieldNames := []string{
		"field_text",
		"field_password",
		"field_password_verify",
		"field_radio",
		"field_checkbox",
		"field_dropdown",
		"field_combobox",
		"field_divider",
		"field_empty_field",
		"field_textblob",
		"field_slate_textblob",
		"field_submit_button",
		"field_error_display",
		"field_flow_link",
		"field_flow_button",
		"field_recaptcha_v2",
		"field_qr_code",
		"field_social_login_button",
	}

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the form name, which must be provided and must be unique within an environment.",
	)

	descriptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the description of the form.",
	)

	categoryDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of form.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMFORMCATEGORY_CUSTOM): "allows the form to be built with fields that do not map specifically to the PingOne directory attributes",
	},
	).DefaultValue(string(management.ENUMFORMCATEGORY_CUSTOM))

	colsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the number of columns in the form (min = `%d`; max = `%d`).", colsMinValue, colsMaxValue),
	).DefaultValue("UNKNOWN")

	componentsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the form configuration elements.",
	)

	componentsFieldsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies the form fields that make up the form.",
	)

	componentsFieldsPositionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the position of the form field in the form.",
	)

	componentsFieldsPositionColDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the column position of the form field in the form  (min = `%d`; max = `%d`).", colsMinValue, colsMaxValue),
	)

	componentsFieldsPositionRowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the row position of the form field in the form (maximum number is `%d`).", rowMaxValue),
	)

	componentsFieldsPositionWidthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the width of the form field in the form (in percentage).",
	)

	componentsFieldsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of form field.",
	).AllowedValuesEnum(management.AllowedEnumFormFieldTypeEnumValues)

	fieldTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the field types in the form.",
	).AllowedValuesEnum(management.AllowedEnumFormFieldTypeEnumValues)

	languageBundleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An map of strings that provides i18n keys to their translations. This object includes both the keys and their default translations. The PingOne language management service finds this object, and creates the new keys for translation for this form.",
	)

	markOptionalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether optional fields are highlighted in the rendered form.",
	)

	markRequiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether required fields are highlighted in the rendered form.",
	)

	translationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how to translate the text strings in the form.",
	).AllowedValuesEnum(management.AllowedEnumFormTranslationMethodEnumValues).DefaultValue("UNKNOWN")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne forms for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the form in."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description:         descriptionDescription.Description,
				MarkdownDescription: descriptionDescription.MarkdownDescription,
				Optional:            true,
			},

			"category": schema.StringAttribute{
				Description:         categoryDescription.Description,
				MarkdownDescription: categoryDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(management.ENUMFORMCATEGORY_CUSTOM)),
			},

			"cols": schema.Int64Attribute{
				Description:         colsDescription.Description,
				MarkdownDescription: colsDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Int64{
					int64validator.Between(colsMinValue, colsMaxValue),
				},
			},

			"components": schema.SingleNestedAttribute{
				Description:         componentsDescription.Description,
				MarkdownDescription: componentsDescription.MarkdownDescription,
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"fields": schema.SetNestedAttribute{
						Description:         componentsFieldsDescription.Description,
						MarkdownDescription: componentsFieldsDescription.MarkdownDescription,
						Required:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"position": schema.SingleNestedAttribute{
									Description:         componentsFieldsPositionDescription.Description,
									MarkdownDescription: componentsFieldsPositionDescription.MarkdownDescription,
									Required:            true,

									Attributes: map[string]schema.Attribute{
										"col": schema.Int64Attribute{
											Description:         componentsFieldsPositionColDescription.Description,
											MarkdownDescription: componentsFieldsPositionColDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.Int64{
												int64validator.Between(colsMinValue, colsMaxValue),
											},
										},

										"row": schema.Int64Attribute{
											Description:         componentsFieldsPositionRowDescription.Description,
											MarkdownDescription: componentsFieldsPositionRowDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.Int64{
												int64validator.AtMost(rowMaxValue),
											},
										},

										"width": schema.Int64Attribute{
											Description:         componentsFieldsPositionWidthDescription.Description,
											MarkdownDescription: componentsFieldsPositionWidthDescription.MarkdownDescription,
											Optional:            true,
										},
									},
								},

								"type": schema.StringAttribute{
									Description:         componentsFieldsTypeDescription.Description,
									MarkdownDescription: componentsFieldsTypeDescription.MarkdownDescription,
									Computed:            true,
								},

								// The form fields
								"field_text": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `TEXT` form field type."),

									formFieldElementSchemaAttributes(management.ENUMFORMFIELDTYPE_TEXT),

									formFieldNames,
								),

								"field_password": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `PASSWORD` form field type."),

									formFieldElementSchemaAttributes(management.ENUMFORMFIELDTYPE_PASSWORD),

									formFieldNames,
								),

								"field_password_verify": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `PASSWORD_VERIFY` form field type."),

									formFieldElementPasswordVerifySchemaAttributes(),

									formFieldNames,
								),

								"field_radio": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `RADIO` form field type."),

									formFieldElementSchemaAttributes(management.ENUMFORMFIELDTYPE_RADIO),

									formFieldNames,
								),

								"field_checkbox": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `CHECKBOX` form field type."),

									formFieldElementSchemaAttributes(management.ENUMFORMFIELDTYPE_CHECKBOX),

									formFieldNames,
								),

								"field_dropdown": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `DROPDOWN` form field type."),

									formFieldElementSchemaAttributes(management.ENUMFORMFIELDTYPE_DROPDOWN),

									formFieldNames,
								),

								"field_combobox": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `COMBOBOX` form field type."),

									map[string]schema.Attribute{}, //

									formFieldNames,
								),

								"field_divider": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `DIVIDER` form field type."),

									formFieldItemSchemaAttributes(),

									formFieldNames,
								),

								"field_empty_field": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `EMPTY_FIELD` form field type."),

									formFieldItemSchemaAttributes(),

									formFieldNames,
								),

								"field_textblob": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `TEXTBLOB` form field type."),

									formFieldItemSchemaAttributes(),

									formFieldNames,
								),

								"field_slate_textblob": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `SLATE_TEXTBLOB` form field type."),

									formFieldItemSchemaAttributes(),

									formFieldNames,
								),

								"field_submit_button": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `SUBMIT_BUTTON` form field type."),

									formFieldButtonSchemaAttributes(),

									formFieldNames,
								),

								"field_error_display": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `ERROR_DISPLAY` form field type."),

									formFieldItemSchemaAttributes(),

									formFieldNames,
								),

								"field_flow_link": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `FLOW_LINK` form field type."),

									formFieldFlowLinkSchemaAttributes(),

									formFieldNames,
								),

								"field_flow_button": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `FLOW_BUTTON` form field type."),

									formFieldButtonSchemaAttributes(),

									formFieldNames,
								),

								"field_recaptcha_v2": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `RECAPTCHA_V2` form field type."),

									formFieldRecaptchaV2SchemaAttributes(),

									formFieldNames,
								),

								"field_qr_code": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `QR_CODE` form field type."),

									formFieldQrCodeSchemaAttributes(),

									formFieldNames,
								),

								"field_social_login_button": formFieldSchema(
									framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies options for the `SOCIAL_LOGIN_BUTTON` form field type."),

									formFieldSocialLoginButtonSchemaAttributes(),

									formFieldNames,
								),
							},
						},

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},
				},
			},

			"field_types": schema.SetAttribute{
				Description:         fieldTypesDescription.Description,
				MarkdownDescription: fieldTypesDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,

				// Validators: []validator.Set{
				// 	setvalidator.ValueStringsAre(
				// 		stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormFieldTypeEnumValues)...),
				// 	),
				// },
			},

			"language_bundle": schema.MapAttribute{
				Description:         languageBundleDescription.Description,
				MarkdownDescription: languageBundleDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,
			},

			"mark_optional": schema.BoolAttribute{
				Description:         markOptionalDescription.Description,
				MarkdownDescription: markOptionalDescription.MarkdownDescription,
				Required:            true,
			},

			"mark_required": schema.BoolAttribute{
				Description:         markRequiredDescription.Description,
				MarkdownDescription: markRequiredDescription.MarkdownDescription,
				Required:            true,
			},

			"translation_method": schema.StringAttribute{
				Description:         translationMethodDescription.Description,
				MarkdownDescription: translationMethodDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormTranslationMethodEnumValues)...),
				},
			},
		},
	}
}

func formFieldSchema(description framework.SchemaAttributeDescription, attributes map[string]schema.Attribute, exactlyOneOfBlockNames []string) schema.Attribute {
	description = description.ExactlyOneOf(exactlyOneOfBlockNames).RequiresReplaceBlock()

	exactlyOneOfPaths := make([]path.Expression, len(exactlyOneOfBlockNames))
	for i, blockName := range exactlyOneOfBlockNames {
		exactlyOneOfPaths[i] = path.MatchRelative().AtParent().AtName(blockName)
	}

	return schema.SingleNestedAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Optional:            true,

		Attributes: attributes,

		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},

		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				exactlyOneOfPaths...,
			),
		},
	}
}

func formFieldElementSchemaAttributes(fieldType management.EnumFormFieldType) map[string]schema.Attribute {

	layoutRequired := false
	validationRequired := false
	optionsRequired := false

	switch fieldType {
	case management.ENUMFORMFIELDTYPE_CHECKBOX, management.ENUMFORMFIELDTYPE_RADIO:
		layoutRequired = true
		optionsRequired = true
	case management.ENUMFORMFIELDTYPE_DROPDOWN:
		optionsRequired = true
	case management.ENUMFORMFIELDTYPE_TEXT:
		validationRequired = true
	}

	attributeDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the linked directory attribute is disabled.",
	).RequiresReplace()

	keyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an identifier for the field component.",
	)

	labelModeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how the field is rendered.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLabelModeEnumValues)

	layoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies layout attributes for radio button and checkbox fields.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLayoutEnumValues)

	optionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array of strings that specifies the unique list of options.",
	)

	requiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the field is required.",
	)

	validationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object containing validation data for the field.",
	)

	validationRegexDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a validation regular expression. The expression must be a valid regular expression string. This is a required property when the validation type is `CUSTOM`.",
	)

	validationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the validation type",
	).AllowedValuesEnum(management.AllowedEnumFormElementValidationTypeEnumValues)

	validationErrorMessageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the error message to be displayed when the field validation fails.",
	)

	otherOptionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the end user can type an entry that is not in a predefined list.",
	)

	otherOptionKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies whether the form identifies that the choice is a custom choice not from a predefined list.",
	)

	otherOptionLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the label for a custom or \"other\" choice in a list.",
	)

	otherOptionInputLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the label for the other option in drop-down controls.",
	)

	otherOptionAttributeDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the directory attribute option is disabled. Set to `true` if it references a PingOne directory attribute.",
	)

	return map[string]schema.Attribute{
		"attribute_disabled": schema.BoolAttribute{
			Description:         attributeDisabledDescription.Description,
			MarkdownDescription: attributeDisabledDescription.MarkdownDescription,
			Optional:            true,

			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},

		"key": schema.StringAttribute{
			Description:         keyDescription.Description,
			MarkdownDescription: keyDescription.MarkdownDescription,
			Required:            true,
		},

		"label_mode": schema.StringAttribute{
			Description:         labelModeDescription.Description,
			MarkdownDescription: labelModeDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementLabelModeEnumValues)...),
			},
		},

		"layout": schema.StringAttribute{
			Description:         layoutDescription.Description,
			MarkdownDescription: layoutDescription.MarkdownDescription,
			Required:            layoutRequired,
			Optional:            !layoutRequired,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementLayoutEnumValues)...),
			},
		},

		"options": schema.SetAttribute{
			Description:         optionsDescription.Description,
			MarkdownDescription: optionsDescription.MarkdownDescription,
			Required:            optionsRequired,
			Optional:            !optionsRequired,

			ElementType: types.StringType,
		},

		"required": schema.BoolAttribute{
			Description:         requiredDescription.Description,
			MarkdownDescription: requiredDescription.MarkdownDescription,
			Required:            true,
		},

		"validation": schema.SingleNestedAttribute{
			Description:         validationDescription.Description,
			MarkdownDescription: validationDescription.MarkdownDescription,
			Required:            validationRequired,
			Optional:            !validationRequired,

			Attributes: map[string]schema.Attribute{
				"regex": schema.StringAttribute{
					Description:         validationRegexDescription.Description,
					MarkdownDescription: validationRegexDescription.MarkdownDescription,
					Optional:            true,
				},

				"type": schema.StringAttribute{
					Description:         validationTypeDescription.Description,
					MarkdownDescription: validationTypeDescription.MarkdownDescription,
					Required:            validationRequired,
					Optional:            !validationRequired,

					Validators: []validator.String{
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementValidationTypeEnumValues)...),
					},
				},

				"error_message": schema.StringAttribute{
					Description:         validationErrorMessageDescription.Description,
					MarkdownDescription: validationErrorMessageDescription.MarkdownDescription,
					Optional:            true,
				},
			},
		},

		"other_option_enabled": schema.BoolAttribute{
			Description:         otherOptionEnabledDescription.Description,
			MarkdownDescription: otherOptionEnabledDescription.MarkdownDescription,
			Optional:            true,
		},

		"other_option_key": schema.StringAttribute{
			Description:         otherOptionKeyDescription.Description,
			MarkdownDescription: otherOptionKeyDescription.MarkdownDescription,
			Optional:            true,
		},

		"other_option_label": schema.StringAttribute{
			Description:         otherOptionLabelDescription.Description,
			MarkdownDescription: otherOptionLabelDescription.MarkdownDescription,
			Optional:            true,
		},

		"other_option_input_label": schema.StringAttribute{
			Description:         otherOptionInputLabelDescription.Description,
			MarkdownDescription: otherOptionInputLabelDescription.MarkdownDescription,
			Optional:            true,
		},

		"other_option_attribute_disabled": schema.BoolAttribute{
			Description:         otherOptionAttributeDisabledDescription.Description,
			MarkdownDescription: otherOptionAttributeDisabledDescription.MarkdownDescription,
			Optional:            true,
		},
	}
}

func formFieldItemSchemaAttributes() map[string]schema.Attribute {
	contentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the field content (for example, HTML.)",
	)

	return map[string]schema.Attribute{
		"content": schema.StringAttribute{
			Description:         contentDescription.Description,
			MarkdownDescription: contentDescription.MarkdownDescription,
			Optional:            true,
		},
	}
}

func formFieldElementPasswordVerifySchemaAttributes() map[string]schema.Attribute {
	labelPasswordVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that when a second field for verifies password is used, this poperty specifies the field label for that verify field.",
	)

	return map[string]schema.Attribute{
		"label_password_verify": schema.StringAttribute{
			Description:         labelPasswordVerifyDescription.Description,
			MarkdownDescription: labelPasswordVerifyDescription.MarkdownDescription,
			Optional:            true,
		},
	}
}

func formFieldButtonSchemaAttributes() map[string]schema.Attribute {
	keyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an identifier for the field component.",
	)

	labelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button label.",
	)

	stylesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes style settings for the button.",
	)

	stylesWidthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the button width. Set as a percentage.",
	)

	stylesAlignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	stylesBackgroundColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button background color. The value must be a valid hexadecimal color.",
	)

	stylesTextColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button text color. The value must be a valid hexadecimal color.",
	)

	stylesBorderColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button border color. The value must be a valid hexadecimal color.",
	)

	stylesEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the button is enabled.",
	)

	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Description:         keyDescription.Description,
			MarkdownDescription: keyDescription.MarkdownDescription,
			Required:            true,
		},

		"label": schema.StringAttribute{
			Description:         labelDescription.Description,
			MarkdownDescription: labelDescription.MarkdownDescription,
			Required:            true,
		},

		"styles": schema.SingleNestedAttribute{
			Description:         stylesDescription.Description,
			MarkdownDescription: stylesDescription.MarkdownDescription,
			Optional:            true,

			Attributes: map[string]schema.Attribute{
				"width": schema.Int64Attribute{
					Description:         stylesWidthDescription.Description,
					MarkdownDescription: stylesWidthDescription.MarkdownDescription,
					Optional:            true,
				},

				"alignment": schema.StringAttribute{
					Description:         stylesAlignmentDescription.Description,
					MarkdownDescription: stylesAlignmentDescription.MarkdownDescription,
					Optional:            true,

					Validators: []validator.String{
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormItemAlignmentEnumValues)...),
					},
				},

				"background_color": schema.StringAttribute{
					Description:         stylesBackgroundColorDescription.Description,
					MarkdownDescription: stylesBackgroundColorDescription.MarkdownDescription,
					Optional:            true,
				},

				"text_color": schema.StringAttribute{
					Description:         stylesTextColorDescription.Description,
					MarkdownDescription: stylesTextColorDescription.MarkdownDescription,
					Optional:            true,
				},

				"border_color": schema.StringAttribute{
					Description:         stylesBorderColorDescription.Description,
					MarkdownDescription: stylesBorderColorDescription.MarkdownDescription,
					Optional:            true,
				},

				"enabled": schema.BoolAttribute{
					Description:         stylesEnabledDescription.Description,
					MarkdownDescription: stylesEnabledDescription.MarkdownDescription,
					Optional:            true,
				},
			},
		},
	}
}

func formFieldFlowLinkSchemaAttributes() map[string]schema.Attribute {
	keyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an identifier for the field component.",
	)

	labelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the link label.",
	)

	stylesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes style settings for the flow link.",
	)

	stylesHorizontalAlignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the link alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	stylesTextColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the link text color.",
	)

	stylesEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the link is enabled.",
	)

	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Description:         keyDescription.Description,
			MarkdownDescription: keyDescription.MarkdownDescription,
			Required:            true,
		},

		"label": schema.StringAttribute{
			Description:         labelDescription.Description,
			MarkdownDescription: labelDescription.MarkdownDescription,
			Required:            true,
		},

		"styles": schema.SingleNestedAttribute{
			Description:         stylesDescription.Description,
			MarkdownDescription: stylesDescription.MarkdownDescription,
			Optional:            true,

			Attributes: map[string]schema.Attribute{
				"horizontal_alignment": schema.StringAttribute{
					Description:         stylesHorizontalAlignmentDescription.Description,
					MarkdownDescription: stylesHorizontalAlignmentDescription.MarkdownDescription,
					Optional:            true,

					Validators: []validator.String{
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormItemAlignmentEnumValues)...),
					},
				},

				"text_color": schema.StringAttribute{
					Description:         stylesTextColorDescription.Description,
					MarkdownDescription: stylesTextColorDescription.MarkdownDescription,
					Optional:            true,
				},

				"enabled": schema.StringAttribute{
					Description:         stylesEnabledDescription.Description,
					MarkdownDescription: stylesEnabledDescription.MarkdownDescription,
					Optional:            true,
				},
			},
		},
	}
}

func formFieldRecaptchaV2SchemaAttributes() map[string]schema.Attribute {
	keyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an identifier for the field component.",
	)

	sizeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the reCAPTCHA size.",
	).AllowedValuesEnum(management.AllowedEnumFormRecaptchaV2SizeEnumValues)

	themeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the reCAPTCHA theme.",
	).AllowedValuesEnum(management.AllowedEnumFormRecaptchaV2ThemeEnumValues)

	alignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the reCAPTCHA alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Description:         keyDescription.Description,
			MarkdownDescription: keyDescription.MarkdownDescription,
			Required:            true,
		},

		"size": schema.StringAttribute{
			Description:         sizeDescription.Description,
			MarkdownDescription: sizeDescription.MarkdownDescription,
			Required:            true,
		},

		"theme": schema.StringAttribute{
			Description:         themeDescription.Description,
			MarkdownDescription: themeDescription.MarkdownDescription,
			Required:            true,
		},

		"alignment": schema.StringAttribute{
			Description:         alignmentDescription.Description,
			MarkdownDescription: alignmentDescription.MarkdownDescription,
			Required:            true,
		},
	}
}

func formFieldQrCodeSchemaAttributes() map[string]schema.Attribute {
	qrCodeTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the QR Code type.",
	)

	alignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the QR Code alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	showBorderDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies the border visibility.",
	)

	return map[string]schema.Attribute{
		"qr_code_type": schema.StringAttribute{
			Description:         qrCodeTypeDescription.Description,
			MarkdownDescription: qrCodeTypeDescription.MarkdownDescription,
			Required:            true,
		},

		"alignment": schema.StringAttribute{
			Description:         alignmentDescription.Description,
			MarkdownDescription: alignmentDescription.MarkdownDescription,
			Required:            true,
		},

		"show_border": schema.BoolAttribute{
			Description:         showBorderDescription.Description,
			MarkdownDescription: showBorderDescription.MarkdownDescription,
			Required:            true,
		},
	}
}

func formFieldSocialLoginButtonSchemaAttributes() map[string]schema.Attribute {
	labelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the social login button label.",
	)

	stylesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes style settings for the social login button.",
	)

	stylesHorizontalAlignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the social login button alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	stylesTextColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the social login button text color.",
	)

	stylesEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the social login button is enabled.",
	)

	idpTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the external identity provider type.",
	).AllowedValuesEnum(management.AllowedEnumFormSocialLoginIdpTypeEnumValues)

	idpNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the external identity provider name.",
	)

	idpIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the external identity provider's ID.",
	)

	idpEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the external identity provider is enabled.",
	)

	iconSrcDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the HTTP link (URL format) for the external identity provider's icon.",
	)

	widthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the button width. Set as a percentage.",
	)

	return map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description:         labelDescription.Description,
			MarkdownDescription: labelDescription.MarkdownDescription,
			Required:            true,
		},

		"styles": schema.SingleNestedAttribute{
			Description:         stylesDescription.Description,
			MarkdownDescription: stylesDescription.MarkdownDescription,
			Optional:            true,

			Attributes: map[string]schema.Attribute{
				"horizontal_alignment": schema.StringAttribute{
					Description:         stylesHorizontalAlignmentDescription.Description,
					MarkdownDescription: stylesHorizontalAlignmentDescription.MarkdownDescription,
					Optional:            true,

					Validators: []validator.String{
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormItemAlignmentEnumValues)...),
					},
				},

				"text_color": schema.StringAttribute{
					Description:         stylesTextColorDescription.Description,
					MarkdownDescription: stylesTextColorDescription.MarkdownDescription,
					Optional:            true,
				},

				"enabled": schema.StringAttribute{
					Description:         stylesEnabledDescription.Description,
					MarkdownDescription: stylesEnabledDescription.MarkdownDescription,
					Optional:            true,
				},
			},
		},

		"idp_type": schema.StringAttribute{
			Description:         idpTypeDescription.Description,
			MarkdownDescription: idpTypeDescription.MarkdownDescription,
			Required:            true,
		},

		"idp_name": schema.StringAttribute{
			Description:         idpNameDescription.Description,
			MarkdownDescription: idpNameDescription.MarkdownDescription,
			Required:            true,
		},

		"idp_id": schema.StringAttribute{
			Description:         idpIdDescription.Description,
			MarkdownDescription: idpIdDescription.MarkdownDescription,
			Required:            true,
		},

		"idp_enabled": schema.BoolAttribute{
			Description:         idpEnabledDescription.Description,
			MarkdownDescription: idpEnabledDescription.MarkdownDescription,
			Required:            true,
		},

		"icon_src": schema.BoolAttribute{
			Description:         iconSrcDescription.Description,
			MarkdownDescription: iconSrcDescription.MarkdownDescription,
			Required:            true,
		},

		"width": schema.Int64Attribute{
			Description:         widthDescription.Description,
			MarkdownDescription: widthDescription.MarkdownDescription,
			Optional:            true,
		},
	}
}

func (r *FormResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FormResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state formResourceModel

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
	form, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.CreateForm(ctx, plan.EnvironmentId.ValueString()).Form(*form).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateForm",
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

func (r *FormResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *formResourceModel

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
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.ReadForm(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Include(management.ENUMFORMSINCLUDEPARAMETER_COMPONENTS).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadForm",
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

func (r *FormResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state formResourceModel

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
	form, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.UpdateForm(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Form(*form).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateForm",
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

func (r *FormResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *formResourceModel

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
			fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.DeleteForm(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteForm",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FormResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "form_id",
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
		pathForm := idComponent.Label

		if idComponent.PrimaryID {
			pathForm = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathForm), attributes[idComponent.Label])...)
	}
}

func (p *formResourceModel) expand(ctx context.Context) (*management.Form, diag.Diagnostics) {
	var diags diag.Diagnostics

	var componentsPlan *formComponentsResourceModel
	diags.Append(p.Components.As(ctx, componentsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	var componentsFieldsPlan []formComponentsFieldResourceModel
	diags.Append(componentsPlan.Fields.ElementsAs(ctx, componentsFieldsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	componentFields := make([]management.FormField, 0)
	for _, v := range componentsFieldsPlan {
		componentField, d := v.expand(ctx)
		diags.Append(d...)

		if componentField != nil {
			componentFields = append(componentFields, *componentField)
		}
	}

	data := management.NewForm(
		p.Name.ValueString(),
		management.EnumFormCategory(p.Category.ValueString()),
		*management.NewFormComponents(componentFields),
		p.MarkOptional.ValueBool(),
		p.MarkRequired.ValueBool(),
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Cols.IsNull() && !p.Cols.IsUnknown() {
		data.SetCols(int32(p.Cols.ValueInt64()))
	}

	if !p.FieldTypes.IsNull() && !p.FieldTypes.IsUnknown() {
		var plan []string
		diags.Append(p.FieldTypes.ElementsAs(ctx, plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		fieldTypes := make([]management.EnumFormFieldType, 0)
		for _, v := range plan {
			fieldTypes = append(fieldTypes, management.EnumFormFieldType(v))
		}

		data.SetFieldTypes(fieldTypes)
	}

	if !p.LanguageBundle.IsNull() && !p.LanguageBundle.IsUnknown() {
		var plan map[string]string
		diags.Append(p.LanguageBundle.ElementsAs(ctx, plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetLanguageBundle(plan)
	}

	if !p.TranslationMethod.IsNull() && !p.TranslationMethod.IsUnknown() {
		data.SetTranslationMethod(management.EnumFormTranslationMethod(p.TranslationMethod.ValueString()))
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expand(ctx context.Context) (*management.FormField, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := &management.FormField{}

	var positionPlan formComponentsFieldPositionResourceModel
	diags.Append(p.Position.As(ctx, positionPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	positionData, d := positionPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !p.FieldText.IsNull() && !p.FieldText.IsUnknown() {
		var plan formComponentsFieldTextResourceModel
		diags.Append(p.FieldText.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldText, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldPassword.IsNull() && !p.FieldPassword.IsUnknown() {
		var plan formComponentsFieldPasswordResourceModel
		diags.Append(p.FieldPassword.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldPassword, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldPasswordVerify.IsNull() && !p.FieldPasswordVerify.IsUnknown() {
		var plan formComponentsFieldPasswordVerifyResourceModel
		diags.Append(p.FieldPasswordVerify.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldPasswordVerify, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldRadio.IsNull() && !p.FieldRadio.IsUnknown() {
		var plan formComponentsFieldRadioResourceModel
		diags.Append(p.FieldRadio.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldRadio, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldCheckbox.IsNull() && !p.FieldCheckbox.IsUnknown() {
		var plan formComponentsFieldCheckboxResourceModel
		diags.Append(p.FieldCheckbox.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldCheckbox, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldDropdown.IsNull() && !p.FieldDropdown.IsUnknown() {
		var plan formComponentsFieldDropdownResourceModel
		diags.Append(p.FieldDropdown.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldDropdown, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldCombobox.IsNull() && !p.FieldCombobox.IsUnknown() {
		var plan formComponentsFieldElementResourceModel
		diags.Append(p.FieldCombobox.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldCombobox, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldDivider.IsNull() && !p.FieldDivider.IsUnknown() {
		var plan formComponentsFieldDividerResourceModel
		diags.Append(p.FieldDivider.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldDivider, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldEmptyField.IsNull() && !p.FieldEmptyField.IsUnknown() {
		var plan formComponentsFieldEmptyFieldResourceModel
		diags.Append(p.FieldEmptyField.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldEmptyField, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldTextblob.IsNull() && !p.FieldTextblob.IsUnknown() {
		var plan formComponentsFieldTextblobResourceModel
		diags.Append(p.FieldTextblob.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldTextblob, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldSlateTextblob.IsNull() && !p.FieldSlateTextblob.IsUnknown() {
		var plan formComponentsFieldSlateTextblobResourceModel
		diags.Append(p.FieldSlateTextblob.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldSlateTextblob, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldSubmitButton.IsNull() && !p.FieldSubmitButton.IsUnknown() {
		var plan formComponentsFieldSubmitButtonResourceModel
		diags.Append(p.FieldSubmitButton.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldSubmitButton, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldErrorDisplay.IsNull() && !p.FieldErrorDisplay.IsUnknown() {
		var plan formComponentsFieldErrorDisplayResourceModel
		diags.Append(p.FieldErrorDisplay.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldErrorDisplay, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldFlowLink.IsNull() && !p.FieldFlowLink.IsUnknown() {
		var plan formComponentsFieldFlowLinkResourceModel
		diags.Append(p.FieldFlowLink.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldFlowLink, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldFlowButton.IsNull() && !p.FieldFlowButton.IsUnknown() {
		var plan formComponentsFieldFlowButtonResourceModel
		diags.Append(p.FieldFlowButton.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldFlowButton, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldRecaptchaV2.IsNull() && !p.FieldRecaptchaV2.IsUnknown() {
		var plan formComponentsFieldRecaptchaV2ResourceModel
		diags.Append(p.FieldRecaptchaV2.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldRecaptchaV2, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldQrCode.IsNull() && !p.FieldQrCode.IsUnknown() {
		var plan formComponentsFieldQrCodeResourceModel
		diags.Append(p.FieldQrCode.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldQrCode, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if !p.FieldSocialLoginButton.IsNull() && !p.FieldSocialLoginButton.IsUnknown() {
		var plan formComponentsFieldSocialLoginButtonResourceModel
		diags.Append(p.FieldSocialLoginButton.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.FormFieldSocialLoginButton, d = plan.expand(ctx, positionData)
		diags.Append(d...)
	}

	if diags.HasError() {
		return nil, diags
	}

	return data, diags
}

func (p *formComponentsFieldPositionResourceModel) expand(ctx context.Context) (*management.FormFieldCommonPosition, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldCommonPosition(
		int32(p.Row.ValueInt64()),
		int32(p.Col.ValueInt64()),
	)

	if !p.Width.IsNull() && !p.Width.IsUnknown() {
		data.SetWidth(int32(p.Width.ValueInt64()))
	}

	return data, diags
}

func (p *formComponentsFieldElementValidationResourceModel) expand(ctx context.Context) (*management.FormElementValidation, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormElementValidation()

	if !p.Regex.IsNull() && !p.Regex.IsUnknown() {
		data.SetRegex(p.Regex.ValueString())
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(management.EnumFormElementValidationType(p.Type.ValueString()))
	}

	if !p.ErrorMessage.IsNull() && !p.ErrorMessage.IsUnknown() {
		data.SetErrorMessage(p.ErrorMessage.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldTextResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldText, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan formComponentsFieldElementValidationResourceModel
	p.Validation.As(ctx, plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})

	validationData, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewFormFieldText(
		management.ENUMFORMFIELDTYPE_TEXT,
		*positionData,
		p.Key.ValueString(),
		p.Required.ValueBool(),
		*validationData,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Layout.IsNull() && !p.Layout.IsUnknown() {
		data.SetLayout(management.EnumFormElementLayout(p.Layout.ValueString()))
	}

	if !p.Options.IsNull() && !p.Options.IsUnknown() {
		var options []string
		diags.Append(p.Options.ElementsAs(ctx, options, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetOptions(options)
	}

	if !p.OtherOptionEnabled.IsNull() && !p.OtherOptionEnabled.IsUnknown() {
		data.SetOtherOptionEnabled(p.OtherOptionEnabled.ValueBool())
	}

	if !p.OtherOptionKey.IsNull() && !p.OtherOptionKey.IsUnknown() {
		data.SetOtherOptionKey(p.OtherOptionKey.ValueString())
	}

	if !p.OtherOptionLabel.IsNull() && !p.OtherOptionLabel.IsUnknown() {
		data.SetOtherOptionlabel(p.OtherOptionLabel.ValueString())
	}

	if !p.OtherOptionInputLabel.IsNull() && !p.OtherOptionInputLabel.IsUnknown() {
		data.SetOtherOptionInputlabel(p.OtherOptionInputLabel.ValueString())
	}

	if !p.OtherOptionAttributeDisabled.IsNull() && !p.OtherOptionAttributeDisabled.IsUnknown() {
		data.SetOtherOptionAttributeDisabled(p.OtherOptionAttributeDisabled.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldPasswordResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldPassword, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldPassword(
		management.ENUMFORMFIELDTYPE_PASSWORD,
		*positionData,
		p.Key.ValueString(),
		p.Required.ValueBool(),
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Layout.IsNull() && !p.Layout.IsUnknown() {
		data.SetLayout(management.EnumFormElementLayout(p.Layout.ValueString()))
	}

	if !p.Options.IsNull() && !p.Options.IsUnknown() {
		var options []string
		diags.Append(p.Options.ElementsAs(ctx, options, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetOptions(options)
	}

	if !p.OtherOptionEnabled.IsNull() && !p.OtherOptionEnabled.IsUnknown() {
		data.SetOtherOptionEnabled(p.OtherOptionEnabled.ValueBool())
	}

	if !p.OtherOptionKey.IsNull() && !p.OtherOptionKey.IsUnknown() {
		data.SetOtherOptionKey(p.OtherOptionKey.ValueString())
	}

	if !p.OtherOptionLabel.IsNull() && !p.OtherOptionLabel.IsUnknown() {
		data.SetOtherOptionlabel(p.OtherOptionLabel.ValueString())
	}

	if !p.OtherOptionInputLabel.IsNull() && !p.OtherOptionInputLabel.IsUnknown() {
		data.SetOtherOptionInputlabel(p.OtherOptionInputLabel.ValueString())
	}

	if !p.OtherOptionAttributeDisabled.IsNull() && !p.OtherOptionAttributeDisabled.IsUnknown() {
		data.SetOtherOptionAttributeDisabled(p.OtherOptionAttributeDisabled.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		p.Validation.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})

		validationData, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Validation = validationData
	}

	return data, diags
}

func (p *formComponentsFieldPasswordVerifyResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldPasswordVerify, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldPasswordVerify(
		management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY,
		*positionData,
	)

	if !p.LabelPasswordVerify.IsNull() && !p.LabelPasswordVerify.IsUnknown() {
		data.SetLabelPasswordVerify(p.LabelPasswordVerify.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldRadioResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldRadio, diag.Diagnostics) {
	var diags diag.Diagnostics

	var options []string
	diags.Append(p.Options.ElementsAs(ctx, options, false)...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewFormFieldRadio(
		management.ENUMFORMFIELDTYPE_RADIO,
		*positionData,
		p.Key.ValueString(),
		p.Required.ValueBool(),
		management.EnumFormElementLayout(p.Layout.ValueString()),
		options,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.OtherOptionEnabled.IsNull() && !p.OtherOptionEnabled.IsUnknown() {
		data.SetOtherOptionEnabled(p.OtherOptionEnabled.ValueBool())
	}

	if !p.OtherOptionKey.IsNull() && !p.OtherOptionKey.IsUnknown() {
		data.SetOtherOptionKey(p.OtherOptionKey.ValueString())
	}

	if !p.OtherOptionLabel.IsNull() && !p.OtherOptionLabel.IsUnknown() {
		data.SetOtherOptionlabel(p.OtherOptionLabel.ValueString())
	}

	if !p.OtherOptionInputLabel.IsNull() && !p.OtherOptionInputLabel.IsUnknown() {
		data.SetOtherOptionInputlabel(p.OtherOptionInputLabel.ValueString())
	}

	if !p.OtherOptionAttributeDisabled.IsNull() && !p.OtherOptionAttributeDisabled.IsUnknown() {
		data.SetOtherOptionAttributeDisabled(p.OtherOptionAttributeDisabled.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		p.Validation.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})

		validationData, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Validation = validationData
	}

	return data, diags
}

func (p *formComponentsFieldCheckboxResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldCheckbox, diag.Diagnostics) {
	var diags diag.Diagnostics

	var options []string
	diags.Append(p.Options.ElementsAs(ctx, options, false)...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewFormFieldCheckbox(
		management.ENUMFORMFIELDTYPE_RADIO,
		*positionData,
		p.Key.ValueString(),
		p.Required.ValueBool(),
		management.EnumFormElementLayout(p.Layout.ValueString()),
		options,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Layout.IsNull() && !p.Layout.IsUnknown() {
		data.SetLayout(management.EnumFormElementLayout(p.Layout.ValueString()))
	}

	if !p.OtherOptionEnabled.IsNull() && !p.OtherOptionEnabled.IsUnknown() {
		data.SetOtherOptionEnabled(p.OtherOptionEnabled.ValueBool())
	}

	if !p.OtherOptionKey.IsNull() && !p.OtherOptionKey.IsUnknown() {
		data.SetOtherOptionKey(p.OtherOptionKey.ValueString())
	}

	if !p.OtherOptionLabel.IsNull() && !p.OtherOptionLabel.IsUnknown() {
		data.SetOtherOptionlabel(p.OtherOptionLabel.ValueString())
	}

	if !p.OtherOptionInputLabel.IsNull() && !p.OtherOptionInputLabel.IsUnknown() {
		data.SetOtherOptionInputlabel(p.OtherOptionInputLabel.ValueString())
	}

	if !p.OtherOptionAttributeDisabled.IsNull() && !p.OtherOptionAttributeDisabled.IsUnknown() {
		data.SetOtherOptionAttributeDisabled(p.OtherOptionAttributeDisabled.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		p.Validation.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})

		validationData, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Validation = validationData
	}

	return data, diags
}

func (p *formComponentsFieldDropdownResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldDropdown, diag.Diagnostics) {
	var diags diag.Diagnostics

	var options []string
	diags.Append(p.Options.ElementsAs(ctx, options, false)...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewFormFieldDropdown(
		management.ENUMFORMFIELDTYPE_RADIO,
		*positionData,
		p.Key.ValueString(),
		p.Required.ValueBool(),
		options,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Layout.IsNull() && !p.Layout.IsUnknown() {
		data.SetLayout(management.EnumFormElementLayout(p.Layout.ValueString()))
	}

	if !p.OtherOptionEnabled.IsNull() && !p.OtherOptionEnabled.IsUnknown() {
		data.SetOtherOptionEnabled(p.OtherOptionEnabled.ValueBool())
	}

	if !p.OtherOptionKey.IsNull() && !p.OtherOptionKey.IsUnknown() {
		data.SetOtherOptionKey(p.OtherOptionKey.ValueString())
	}

	if !p.OtherOptionLabel.IsNull() && !p.OtherOptionLabel.IsUnknown() {
		data.SetOtherOptionlabel(p.OtherOptionLabel.ValueString())
	}

	if !p.OtherOptionInputLabel.IsNull() && !p.OtherOptionInputLabel.IsUnknown() {
		data.SetOtherOptionInputlabel(p.OtherOptionInputLabel.ValueString())
	}

	if !p.OtherOptionAttributeDisabled.IsNull() && !p.OtherOptionAttributeDisabled.IsUnknown() {
		data.SetOtherOptionAttributeDisabled(p.OtherOptionAttributeDisabled.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		p.Validation.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})

		validationData, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Validation = validationData
	}

	return data, diags
}

func (p *formComponentsFieldElementResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldCombobox, diag.Diagnostics) {
	var diags diag.Diagnostics

	return nil, diags
}

func (p *formComponentsFieldDividerResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldDivider, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldDivider(
		management.ENUMFORMFIELDTYPE_DIVIDER,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldEmptyFieldResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldEmptyField, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldEmptyField(
		management.ENUMFORMFIELDTYPE_EMPTY_FIELD,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldTextblobResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldTextblob, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldTextblob(
		management.ENUMFORMFIELDTYPE_TEXTBLOB,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldSlateTextblobResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldSlateTextblob, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldSlateTextblob(
		management.ENUMFORMFIELDTYPE_SLATE_TEXTBLOB,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldSubmitButtonResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldSubmitButton, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldSubmitButton(
		management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldButtonStylesResourceModel
		diags.Append(p.Styles.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx)
		diags.Append(d...)

		data.SetStyles(*stylesData)
	}

	return data, diags
}

func (p *formComponentsFieldButtonStylesResourceModel) expand(ctx context.Context) (*management.FormFlowButtonStyles, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFlowButtonStyles()

	if !p.Alignment.IsNull() && !p.Alignment.IsUnknown() {
		data.SetAlignment(management.EnumFormItemAlignment(p.Alignment.ValueString()))
	}

	if !p.BackgroundColor.IsNull() && !p.BackgroundColor.IsUnknown() {
		data.SetBackgroundColor(p.BackgroundColor.ValueString())
	}

	if !p.BorderColor.IsNull() && !p.BorderColor.IsUnknown() {
		data.SetBorderColor(p.BorderColor.ValueString())
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
		data.SetTextColor(p.TextColor.ValueString())
	}

	if !p.Width.IsNull() && !p.Width.IsUnknown() {
		data.SetWidth(int32(p.Width.ValueInt64()))
	}

	return data, diags
}

func (p *formComponentsFieldErrorDisplayResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldErrorDisplay, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldErrorDisplay(
		management.ENUMFORMFIELDTYPE_ERROR_DISPLAY,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldFlowLinkResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldFlowLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldFlowLink(
		management.ENUMFORMFIELDTYPE_FLOW_LINK,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldFlowLinkStylesResourceModel
		diags.Append(p.Styles.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx)
		diags.Append(d...)

		data.SetStyles(*stylesData)
	}

	return data, diags
}

func (p *formComponentsFieldFlowLinkStylesResourceModel) expand(ctx context.Context) (*management.FormFlowLinkStyles, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFlowLinkStyles()

	if !p.HorizontalAlignment.IsNull() && !p.HorizontalAlignment.IsUnknown() {
		data.SetHorizontalAlignment(management.EnumFormItemAlignment(p.HorizontalAlignment.ValueString()))
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
		data.SetTextColor(p.TextColor.ValueString())
	}

	return data, diags
}

func (p *formComponentsFieldFlowButtonResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldFlowButton, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldFlowButton(
		management.ENUMFORMFIELDTYPE_FLOW_BUTTON,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldButtonStylesResourceModel
		diags.Append(p.Styles.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx)
		diags.Append(d...)

		data.SetStyles(*stylesData)
	}

	return data, diags
}

func (p *formComponentsFieldRecaptchaV2ResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldRecaptchaV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldRecaptchaV2(
		management.ENUMFORMFIELDTYPE_RECAPTCHA_V2,
		*positionData,
		p.Key.ValueString(),
		management.EnumFormRecaptchaV2Size(p.Size.ValueString()),
		management.EnumFormRecaptchaV2Theme(p.Theme.ValueString()),
		management.EnumFormItemAlignment(p.Alignment.ValueString()),
	)

	return data, diags
}

func (p *formComponentsFieldQrCodeResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldQrCode, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldQrCode(
		management.ENUMFORMFIELDTYPE_QR_CODE,
		*positionData,
		management.EnumFormQrCodeType(p.QrCodeType.ValueString()),
		management.EnumFormItemAlignment(p.Alignment.ValueString()),
		p.ShowBorder.ValueBool(),
	)

	return data, diags
}

func (p *formComponentsFieldSocialLoginButtonResourceModel) expand(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldSocialLoginButton, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldSocialLoginButton(
		management.ENUMFORMFIELDTYPE_SOCIAL_LOGIN_BUTTON,
		*positionData,
		p.Label.ValueString(),
		management.EnumFormSocialLoginIdpType(p.IdpType.ValueString()),
		p.IdpType.ValueString(),
		p.IdpId.ValueString(),
		p.IdpEnabled.ValueBool(),
		p.IconSrc.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldSocialLoginButtonStylesResourceModel
		diags.Append(p.Styles.As(ctx, plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx)
		diags.Append(d...)

		data.SetStyles(*stylesData)
	}

	if !p.Width.IsNull() && !p.Width.IsUnknown() {
		data.SetWidth(int32(p.Width.ValueInt64()))
	}

	return data, diags
}

func (p *formComponentsFieldSocialLoginButtonStylesResourceModel) expand(ctx context.Context) (*management.FormSocialLoginButtonStyles, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormSocialLoginButtonStyles()

	if !p.HorizontalAlignment.IsNull() && !p.HorizontalAlignment.IsUnknown() {
		data.SetHorizontalAlignment(management.EnumFormItemAlignment(p.HorizontalAlignment.ValueString()))
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
		data.SetTextColor(p.TextColor.ValueString())
	}

	return data, diags
}

func (p *formResourceModel) toState(apiObject *management.Form) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Category = framework.EnumOkToTF(apiObject.GetCategoryOk())
	p.Cols = framework.Int32OkToTF(apiObject.GetColsOk())
	p.FieldTypes = framework.EnumSetOkToTF(apiObject.GetFieldTypesOk())
	p.LanguageBundle = framework.StringMapOkToTF(apiObject.GetLanguageBundleOk())
	p.MarkOptional = framework.BoolOkToTF(apiObject.GetMarkOptionalOk())
	p.MarkRequired = framework.BoolOkToTF(apiObject.GetMarkRequiredOk())
	p.TranslationMethod = framework.EnumOkToTF(apiObject.GetTranslationMethodOk())

	p.Components = formComponentsOkToTF(apiObject.GetComponentsOk())

	return diags
}

func formComponentsOkToTF(formComponents *management.FormComponents, ok bool) types.Object {
	return types.Object{}
}
