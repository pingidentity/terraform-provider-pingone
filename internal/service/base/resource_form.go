package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
type FormResource serviceClientType

type formResourceModel struct {
	Id                pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId     pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name              types.String                 `tfsdk:"name"`
	Description       types.String                 `tfsdk:"description"`
	Category          types.String                 `tfsdk:"category"`
	Cols              types.Int64                  `tfsdk:"cols"`
	Components        types.Object                 `tfsdk:"components"`
	FieldTypes        types.Set                    `tfsdk:"field_types"`
	MarkOptional      types.Bool                   `tfsdk:"mark_optional"`
	MarkRequired      types.Bool                   `tfsdk:"mark_required"`
	TranslationMethod types.String                 `tfsdk:"translation_method"`
}

type formComponentsResourceModel struct {
	Fields types.Set `tfsdk:"fields"`
}

type formComponentsFieldResourceModel struct {
	Alignment                    types.String `tfsdk:"alignment"`
	AttributeDisabled            types.Bool   `tfsdk:"attribute_disabled"`
	Content                      types.String `tfsdk:"content"`
	Key                          types.String `tfsdk:"key"`
	Label                        types.String `tfsdk:"label"`
	LabelMode                    types.String `tfsdk:"label_mode"`
	LabelPasswordVerify          types.String `tfsdk:"label_password_verify"`
	Layout                       types.String `tfsdk:"layout"`
	Options                      types.Set    `tfsdk:"options"`
	OtherOptionAttributeDisabled types.Bool   `tfsdk:"other_option_attribute_disabled"`
	OtherOptionEnabled           types.Bool   `tfsdk:"other_option_enabled"`
	OtherOptionInputLabel        types.String `tfsdk:"other_option_input_label"`
	OtherOptionKey               types.String `tfsdk:"other_option_key"`
	OtherOptionLabel             types.String `tfsdk:"other_option_label"`
	Position                     types.Object `tfsdk:"position"`
	QrCodeType                   types.String `tfsdk:"qr_code_type"`
	Required                     types.Bool   `tfsdk:"required"`
	ShowBorder                   types.Bool   `tfsdk:"show_border"`
	ShowPasswordRequirements     types.Bool   `tfsdk:"show_password_requirements"`
	Size                         types.String `tfsdk:"size"`
	Styles                       types.Object `tfsdk:"styles"`
	Theme                        types.String `tfsdk:"theme"`
	Type                         types.String `tfsdk:"type"`
	Validation                   types.Object `tfsdk:"validation"`
}

type formComponentsFieldPositionResourceModel struct {
	Col   types.Int64 `tfsdk:"col"`
	Row   types.Int64 `tfsdk:"row"`
	Width types.Int64 `tfsdk:"width"`
}

type formComponentsFieldElementOptionsResourceModel struct {
	Value types.String `tfsdk:"value"`
	Label types.String `tfsdk:"label"`
}

type formComponentsFieldElementValidationResourceModel struct {
	Regex        types.String `tfsdk:"regex"`
	Type         types.String `tfsdk:"type"`
	ErrorMessage types.String `tfsdk:"error_message"`
}

// SUBMIT_BUTTON, FLOW_BUTTON
type formComponentsFieldButtonResourceModel struct {
	Key    types.String `tfsdk:"key"`
	Label  types.String `tfsdk:"label"`
	Styles types.Object `tfsdk:"styles"`
}

type formComponentsFieldFlowButtonResourceModel formComponentsFieldButtonResourceModel

type formComponentsFieldStylesResourceModel struct {
	Alignment       types.String `tfsdk:"alignment"`
	BackgroundColor types.String `tfsdk:"background_color"`
	BorderColor     types.String `tfsdk:"border_color"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Height          types.Int64  `tfsdk:"height"`
	Padding         types.Object `tfsdk:"padding"`
	TextColor       types.String `tfsdk:"text_color"`
	Width           types.Int64  `tfsdk:"width"`
	WidthUnit       types.String `tfsdk:"width_unit"`
}

type formComponentsFieldButtonStylesPaddingResourceModel struct {
	Bottom types.Int64 `tfsdk:"bottom"`
	Left   types.Int64 `tfsdk:"left"`
	Right  types.Int64 `tfsdk:"right"`
	Top    types.Int64 `tfsdk:"top"`
}

type formComponentsFieldsSchemaDef struct {
	Required []string
	Optional []string
}

var (
	// Form Components
	formComponentsTFObjectTypes = map[string]attr.Type{
		"fields": types.SetType{ElemType: types.ObjectType{
			AttrTypes: formComponentsFieldsTFObjectTypes,
		}},
	}

	// Form Components Fields
	formComponentsFieldsTFObjectTypes = map[string]attr.Type{
		"alignment":                       types.StringType,
		"attribute_disabled":              types.BoolType,
		"content":                         types.StringType,
		"key":                             types.StringType,
		"label_mode":                      types.StringType,
		"label_password_verify":           types.StringType,
		"label":                           types.StringType,
		"layout":                          types.StringType,
		"options":                         types.SetType{ElemType: types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}},
		"other_option_attribute_disabled": types.BoolType,
		"other_option_enabled":            types.BoolType,
		"other_option_input_label":        types.StringType,
		"other_option_key":                types.StringType,
		"other_option_label":              types.StringType,
		"position":                        types.ObjectType{AttrTypes: formComponentsFieldsPositionTFObjectTypes},
		"qr_code_type":                    types.StringType,
		"required":                        types.BoolType,
		"show_border":                     types.BoolType,
		"show_password_requirements":      types.BoolType,
		"size":                            types.StringType,
		"styles":                          types.ObjectType{AttrTypes: formComponentsFieldsFieldStylesTFObjectTypes},
		"theme":                           types.StringType,
		"type":                            types.StringType,
		"validation":                      types.ObjectType{AttrTypes: formComponentsFieldsFieldElementValidationTFObjectTypes},
	}

	// Form Components Fields Position
	formComponentsFieldsPositionTFObjectTypes = map[string]attr.Type{
		"col":   types.Int64Type,
		"row":   types.Int64Type,
		"width": types.Int64Type,
	}

	// Form Components Fields Field Element Option
	formComponentsFieldsFieldElementOptionTFObjectTypes = map[string]attr.Type{
		"label": types.StringType,
		"value": types.StringType,
	}

	// Form Components Fields Field Element Validation
	formComponentsFieldsFieldElementValidationTFObjectTypes = map[string]attr.Type{
		"regex":         types.StringType,
		"type":          types.StringType,
		"error_message": types.StringType,
	}

	// Form Components Fields Field Button Styles
	formComponentsFieldsFieldStylesTFObjectTypes = map[string]attr.Type{
		"alignment":        types.StringType,
		"background_color": types.StringType,
		"border_color":     types.StringType,
		"enabled":          types.BoolType,
		"height":           types.Int64Type,
		"padding":          types.ObjectType{AttrTypes: formComponentsFieldsFieldStylesPaddingTFObjectTypes},
		"text_color":       types.StringType,
		"width":            types.Int64Type,
		"width_unit":       types.StringType,
	}

	formComponentsFieldsFieldStylesPaddingTFObjectTypes = map[string]attr.Type{
		"bottom": types.Int64Type,
		"left":   types.Int64Type,
		"right":  types.Int64Type,
		"top":    types.Int64Type,
	}

	formComponentsFieldsSchemaDefMap = map[management.EnumFormFieldType]formComponentsFieldsSchemaDef{
		management.ENUMFORMFIELDTYPE_CHECKBOX: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
				"layout",
				"options",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"required",
			},
		},
		management.ENUMFORMFIELDTYPE_COMBOBOX: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
				"options",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"layout",
				"required",
			},
		},
		management.ENUMFORMFIELDTYPE_DIVIDER: {
			Required: []string{
				"type",
				"position",
			},
			Optional: []string{},
		},
		management.ENUMFORMFIELDTYPE_DROPDOWN: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
				"options",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"layout",
				"required",
			},
		},
		management.ENUMFORMFIELDTYPE_EMPTY_FIELD: {
			Required: []string{
				"type",
				"position",
			},
			Optional: []string{},
		},
		management.ENUMFORMFIELDTYPE_ERROR_DISPLAY: {
			Required: []string{
				"type",
				"position",
			},
			Optional: []string{},
		},
		management.ENUMFORMFIELDTYPE_FLOW_BUTTON: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
			},
			Optional: []string{
				"styles",
			},
		},
		management.ENUMFORMFIELDTYPE_FLOW_LINK: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
			},
			Optional: []string{
				"styles",
			},
		},
		management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"label_password_verify",
				"layout",
				"required",
				"show_password_requirements",
				"validation",
			},
		},
		management.ENUMFORMFIELDTYPE_PASSWORD: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"layout",
				"required",
				"show_password_requirements",
				"validation",
			},
		},
		management.ENUMFORMFIELDTYPE_QR_CODE: {
			Required: []string{
				"type",
				"position",
				"key",
				"qr_code_type",
				"alignment",
			},
			Optional: []string{
				"show_border",
			},
		},
		management.ENUMFORMFIELDTYPE_RADIO: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
				"layout",
				"options",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"required",
			},
		},
		management.ENUMFORMFIELDTYPE_RECAPTCHA_V2: {
			Required: []string{
				"type",
				"position",
				"size",
				"theme",
				"alignment",
			},
			Optional: []string{},
		},
		management.ENUMFORMFIELDTYPE_SLATE_TEXTBLOB: {
			Required: []string{
				"type",
				"position",
			},
			Optional: []string{
				"content",
			},
		},
		management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON: {
			Required: []string{
				"type",
				"position",
				"label",
			},
			Optional: []string{
				"styles",
			},
		},
		management.ENUMFORMFIELDTYPE_TEXT: {
			Required: []string{
				"type",
				"position",
				"key",
				"label",
				"validation",
			},
			Optional: []string{
				"attribute_disabled",
				"label_mode",
				"layout",
				"required",
			},
		},
		management.ENUMFORMFIELDTYPE_TEXTBLOB: {
			Required: []string{
				"type",
				"position",
			},
			Optional: []string{
				"content",
			},
		},
	}
)

// Framework interfaces
var (
	_ resource.Resource                   = &FormResource{}
	_ resource.ResourceWithConfigure      = &FormResource{}
	_ resource.ResourceWithImportState    = &FormResource{}
	_ resource.ResourceWithValidateConfig = &FormResource{}
	_ resource.ResourceWithModifyPlan     = &FormResource{}
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
	const colsMinValue = 1
	const colsMaxValue = 4
	const rowMaxValue = 50
	const colMinValue = 0
	const colMaxValue = 3

	supportedFormFieldTypes := []management.EnumFormFieldType{
		management.ENUMFORMFIELDTYPE_TEXT,
		management.ENUMFORMFIELDTYPE_PASSWORD,
		management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY,
		management.ENUMFORMFIELDTYPE_RADIO,
		management.ENUMFORMFIELDTYPE_CHECKBOX,
		management.ENUMFORMFIELDTYPE_DROPDOWN,
		management.ENUMFORMFIELDTYPE_COMBOBOX,
		management.ENUMFORMFIELDTYPE_DIVIDER,
		management.ENUMFORMFIELDTYPE_EMPTY_FIELD,
		management.ENUMFORMFIELDTYPE_TEXTBLOB,
		management.ENUMFORMFIELDTYPE_SLATE_TEXTBLOB,
		management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON,
		management.ENUMFORMFIELDTYPE_ERROR_DISPLAY,
		management.ENUMFORMFIELDTYPE_FLOW_LINK,
		management.ENUMFORMFIELDTYPE_FLOW_BUTTON,
		management.ENUMFORMFIELDTYPE_RECAPTCHA_V2,
		management.ENUMFORMFIELDTYPE_QR_CODE,
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
	)

	componentsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the form configuration elements.",
	)

	componentsFieldsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies the form fields that make up the form.",
	)

	componentsFieldsPositionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the position of the form field in the form.  The combination of `col` and `row` must be unique between form fields.",
	)

	componentsFieldsPositionColDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the column position of the form field in the form  (min = `%d`; max = `%d`).", colMinValue, colMaxValue),
	)

	componentsFieldsPositionRowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the row position of the form field in the form (maximum number is `%d`).", rowMaxValue),
	)

	componentsFieldsPositionWidthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the width of the form field in the form (in percentage).",
	)

	componentsFieldsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of form field.",
	).AllowedValuesEnum(supportedFormFieldTypes)

	componentsFieldsAttributeDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("attribute_disabled"),
	).AppendMarkdownString(
		"A boolean that specifies whether the linked directory attribute is disabled.",
	)

	componentsFieldsContentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("content"),
	).AppendMarkdownString(
		fmt.Sprintf("A string that specifies the field's content (for example, HTML when the field type is `%s`.)", string(management.ENUMFORMFIELDTYPE_TEXTBLOB)),
	)

	componentsFieldsKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("key"),
	).AppendMarkdownString(
		"A string that specifies an identifier for the field component.",
	)

	componentsFieldsLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("label"),
	).AppendMarkdownString(
		"A string that specifies the field label.",
	)

	componentsFieldsLabelPasswordVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("label_password_verify"),
	).AppendMarkdownString(
		"A string that when a second field for verifies password is used, this property specifies the field label for that verify field.",
	)

	componentsFieldsLabelModeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("label_mode"),
	).AppendMarkdownString(
		"A string that specifies how the field is rendered.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLabelModeEnumValues)

	componentsFieldsLayoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("layout"),
	).AppendMarkdownString(
		"A string that specifies layout attributes for radio button and checkbox fields.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLayoutEnumValues)

	componentsFieldsOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("options"),
	).AppendMarkdownString(
		"An array of objects that specifies the unique list of options.",
	)

	componentsFieldsOptionsLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the option's label in the form field that is shown to the end user.",
	)

	componentsFieldsOptionsValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the option's value in the form field that is posted as form data.",
	)

	componentsFieldsRequiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("required"),
	).AppendMarkdownString(
		"A boolean that specifies whether the field is required.",
	)

	componentsFieldsValidationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("validation"),
	).AppendMarkdownString(
		"An object containing validation data for the field.",
	)

	componentsFieldsValidationRegexDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a validation regular expression. The expression must be a valid regular expression string. This is a required property when the validation type is `CUSTOM`.",
	)

	componentsFieldsValidationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the validation type",
	).AllowedValuesEnum(management.AllowedEnumFormElementValidationTypeEnumValues)

	componentsFieldsValidationErrorMessageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the error message to be displayed when the field validation fails.  When configuring this parameter, the `regex` parameter is required.",
	)

	componentsFieldsOtherOptionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the end user can type an entry that is not in a predefined list.",
	)

	componentsFieldsOtherOptionKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies whether the form identifies that the choice is a custom choice not from a predefined list.",
	)

	componentsFieldsOtherOptionLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the label for a custom or \"other\" choice in a list.",
	)

	componentsFieldsOtherOptionInputLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the label for the other option in drop-down controls.",
	)

	componentsFieldsOtherOptionAttributeDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the directory attribute option is disabled. Set to `true` if it references a PingOne directory attribute.",
	)

	componentsFieldsShowPasswordRequirementsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("show_password_requirements"),
	).AppendMarkdownString(
		"A boolean that specifies whether to display password requirements to the user.",
	)

	componentsFieldsStylesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("styles"),
	).AppendMarkdownString(
		"A single object that describes style settings for the field.",
	)

	componentsFieldsStylesWidthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the button width. Set as a percentage.",
	)

	componentsFieldsStylesAlignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	componentsFieldsStylesBackgroundColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button background color. The value must be a valid hexadecimal color.",
	)

	componentsFieldsStylesBorderColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button border color. The value must be a valid hexadecimal color.",
	)

	componentsFieldsStylesEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the button is enabled.",
	)

	componentsFieldsStylesHeightDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies a custom height of the field (in pixels) when displayed in the form.",
	)

	componentsFieldsStylesPaddingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies custom padding styles for the field.",
	)

	componentsFieldsStylesPaddingTopDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the top padding (in pixels) to apply to the field.",
	)

	componentsFieldsStylesPaddingRightDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the right padding (in pixels) to apply to the field.",
	)

	componentsFieldsStylesPaddingBottomDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the bottom padding (in pixels) to apply to the field.",
	)

	componentsFieldsStylesPaddingLeftDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the left padding (in pixels) to apply to the field.",
	)

	componentsFieldsStylesTextColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button text color. The value must be a valid hexadecimal color.",
	)

	componentsFieldsStylesWidthUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the unit to apply to the `width` parameter.",
	).AllowedValuesEnum(management.AllowedEnumFormStylesWidthUnitEnumValues)

	componentsFieldsSizeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("size"),
	).AppendMarkdownString(
		"A string that specifies the reCAPTCHA size.",
	).AllowedValuesEnum(management.AllowedEnumFormRecaptchaV2SizeEnumValues)

	componentsFieldsThemeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("theme"),
	).AppendMarkdownString(
		"A string that specifies the reCAPTCHA theme.",
	).AllowedValuesEnum(management.AllowedEnumFormRecaptchaV2ThemeEnumValues)

	componentsFieldsAlignmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("alignment"),
	).AppendMarkdownString(
		"A string that specifies the reCAPTCHA alignment.",
	).AllowedValuesEnum(management.AllowedEnumFormItemAlignmentEnumValues)

	componentsFieldsQrCodeTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("qr_code_type"),
	).AppendMarkdownString(
		"A string that specifies the QR Code type.",
	)

	componentsFieldsShowBorderDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		formFieldValidationDocumentation("show_border"),
	).AppendMarkdownString(
		"A boolean that specifies the border visibility.",
	)

	fieldTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the field types in the form.",
	).AllowedValuesEnum(supportedFormFieldTypes)

	markOptionalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether optional fields are highlighted in the rendered form.",
	)

	markRequiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether required fields are highlighted in the rendered form.",
	)

	translationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how to translate the text strings in the form.",
	).AllowedValuesEnum(management.AllowedEnumFormTranslationMethodEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne DaVinci forms in an environment.",

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
				Required:            true,

				Validators: []validator.Int64{
					int64validator.Between(colsMinValue, colsMaxValue),
				},

				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
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
												int64validator.AtLeast(colMinValue),
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
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(supportedFormFieldTypes)...),
									},
								},

								"attribute_disabled": schema.BoolAttribute{
									Description:         componentsFieldsAttributeDisabledDescription.Description,
									MarkdownDescription: componentsFieldsAttributeDisabledDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,
								},

								"content": schema.StringAttribute{
									Description:         componentsFieldsContentDescription.Description,
									MarkdownDescription: componentsFieldsContentDescription.MarkdownDescription,
									Optional:            true,
								},

								"key": schema.StringAttribute{
									Description:         componentsFieldsKeyDescription.Description,
									MarkdownDescription: componentsFieldsKeyDescription.MarkdownDescription,
									Optional:            true,
								},

								"label": schema.StringAttribute{
									Description:         componentsFieldsLabelDescription.Description,
									MarkdownDescription: componentsFieldsLabelDescription.MarkdownDescription,
									Optional:            true,
								},

								"label_password_verify": schema.StringAttribute{
									Description:         componentsFieldsLabelPasswordVerifyDescription.Description,
									MarkdownDescription: componentsFieldsLabelPasswordVerifyDescription.MarkdownDescription,
									Optional:            true,
								},

								"label_mode": schema.StringAttribute{
									Description:         componentsFieldsLabelModeDescription.Description,
									MarkdownDescription: componentsFieldsLabelModeDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementLabelModeEnumValues)...),
									},
								},

								"layout": schema.StringAttribute{
									Description:         componentsFieldsLayoutDescription.Description,
									MarkdownDescription: componentsFieldsLayoutDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementLayoutEnumValues)...),
									},
								},

								"options": schema.SetNestedAttribute{
									Description:         componentsFieldsOptionsDescription.Description,
									MarkdownDescription: componentsFieldsOptionsDescription.MarkdownDescription,
									Optional:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"label": schema.StringAttribute{
												Description:         componentsFieldsOptionsLabelDescription.Description,
												MarkdownDescription: componentsFieldsOptionsLabelDescription.MarkdownDescription,
												Required:            true,
											},

											"value": schema.StringAttribute{
												Description:         componentsFieldsOptionsValueDescription.Description,
												MarkdownDescription: componentsFieldsOptionsValueDescription.MarkdownDescription,
												Required:            true,
											},
										},
									},
								},

								"required": schema.BoolAttribute{
									Description:         componentsFieldsRequiredDescription.Description,
									MarkdownDescription: componentsFieldsRequiredDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,
								},

								"validation": schema.SingleNestedAttribute{
									Description:         componentsFieldsValidationDescription.Description,
									MarkdownDescription: componentsFieldsValidationDescription.MarkdownDescription,
									Optional:            true,

									Attributes: map[string]schema.Attribute{
										"regex": schema.StringAttribute{
											Description:         componentsFieldsValidationRegexDescription.Description,
											MarkdownDescription: componentsFieldsValidationRegexDescription.MarkdownDescription,
											Optional:            true,

											Validators: []validator.String{
												stringvalidatorinternal.IsRequiredIfMatchesPathValue(
													types.StringValue(string(management.ENUMFORMELEMENTVALIDATIONTYPE_CUSTOM)),
													path.MatchRelative().AtParent().AtName("type"),
												),
											},
										},

										"type": schema.StringAttribute{
											Description:         componentsFieldsValidationTypeDescription.Description,
											MarkdownDescription: componentsFieldsValidationTypeDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.String{
												stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementValidationTypeEnumValues)...),
											},
										},

										"error_message": schema.StringAttribute{
											Description:         componentsFieldsValidationErrorMessageDescription.Description,
											MarkdownDescription: componentsFieldsValidationErrorMessageDescription.MarkdownDescription,
											Optional:            true,

											Validators: []validator.String{
												stringvalidator.AlsoRequires(
													path.MatchRelative().AtParent().AtName("regex"),
												),
											},
										},
									},
								},

								"other_option_enabled": schema.BoolAttribute{
									Description:         componentsFieldsOtherOptionEnabledDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionEnabledDescription.MarkdownDescription,
									Computed:            true,
								},

								"other_option_key": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionKeyDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionKeyDescription.MarkdownDescription,
									Computed:            true,

									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},

								"other_option_label": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionLabelDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionLabelDescription.MarkdownDescription,
									Computed:            true,

									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},

								"other_option_input_label": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionInputLabelDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionInputLabelDescription.MarkdownDescription,
									Computed:            true,

									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},

								"other_option_attribute_disabled": schema.BoolAttribute{
									Description:         componentsFieldsOtherOptionAttributeDisabledDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionAttributeDisabledDescription.MarkdownDescription,
									Computed:            true,
								},

								"show_password_requirements": schema.BoolAttribute{
									Description:         componentsFieldsShowPasswordRequirementsDescription.Description,
									MarkdownDescription: componentsFieldsShowPasswordRequirementsDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,
								},

								"styles": schema.SingleNestedAttribute{
									Description:         componentsFieldsStylesDescription.Description,
									MarkdownDescription: componentsFieldsStylesDescription.MarkdownDescription,
									Optional:            true,

									Attributes: map[string]schema.Attribute{
										"alignment": schema.StringAttribute{
											Description:         componentsFieldsStylesAlignmentDescription.Description,
											MarkdownDescription: componentsFieldsStylesAlignmentDescription.MarkdownDescription,
											Optional:            true,

											Validators: []validator.String{
												stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormItemAlignmentEnumValues)...),
											},
										},

										"background_color": schema.StringAttribute{
											Description:         componentsFieldsStylesBackgroundColorDescription.Description,
											MarkdownDescription: componentsFieldsStylesBackgroundColorDescription.MarkdownDescription,
											Optional:            true,
										},

										"border_color": schema.StringAttribute{
											Description:         componentsFieldsStylesBorderColorDescription.Description,
											MarkdownDescription: componentsFieldsStylesBorderColorDescription.MarkdownDescription,
											Optional:            true,
										},

										"enabled": schema.BoolAttribute{
											Description:         componentsFieldsStylesEnabledDescription.Description,
											MarkdownDescription: componentsFieldsStylesEnabledDescription.MarkdownDescription,
											Optional:            true,
											Computed:            true,

											PlanModifiers: []planmodifier.Bool{
												boolplanmodifier.UseStateForUnknown(),
											},
										},

										"height": schema.Int64Attribute{
											Description:         componentsFieldsStylesHeightDescription.Description,
											MarkdownDescription: componentsFieldsStylesHeightDescription.MarkdownDescription,
											Optional:            true,
										},

										"padding": schema.SingleNestedAttribute{
											Description:         componentsFieldsStylesPaddingDescription.Description,
											MarkdownDescription: componentsFieldsStylesPaddingDescription.MarkdownDescription,
											Optional:            true,

											Attributes: map[string]schema.Attribute{
												"top": schema.Int64Attribute{
													Description:         componentsFieldsStylesPaddingTopDescription.Description,
													MarkdownDescription: componentsFieldsStylesPaddingTopDescription.MarkdownDescription,
													Optional:            true,
												},

												"right": schema.Int64Attribute{
													Description:         componentsFieldsStylesPaddingRightDescription.Description,
													MarkdownDescription: componentsFieldsStylesPaddingRightDescription.MarkdownDescription,
													Optional:            true,
												},

												"bottom": schema.Int64Attribute{
													Description:         componentsFieldsStylesPaddingBottomDescription.Description,
													MarkdownDescription: componentsFieldsStylesPaddingBottomDescription.MarkdownDescription,
													Optional:            true,
												},

												"left": schema.Int64Attribute{
													Description:         componentsFieldsStylesPaddingLeftDescription.Description,
													MarkdownDescription: componentsFieldsStylesPaddingLeftDescription.MarkdownDescription,
													Optional:            true,
												},
											},
										},

										"text_color": schema.StringAttribute{
											Description:         componentsFieldsStylesTextColorDescription.Description,
											MarkdownDescription: componentsFieldsStylesTextColorDescription.MarkdownDescription,
											Optional:            true,
										},

										"width": schema.Int64Attribute{
											Description:         componentsFieldsStylesWidthDescription.Description,
											MarkdownDescription: componentsFieldsStylesWidthDescription.MarkdownDescription,
											Optional:            true,
										},

										"width_unit": schema.StringAttribute{
											Description:         componentsFieldsStylesWidthUnitDescription.Description,
											MarkdownDescription: componentsFieldsStylesWidthUnitDescription.MarkdownDescription,
											Optional:            true,
										},
									},
								},

								"size": schema.StringAttribute{
									Description:         componentsFieldsSizeDescription.Description,
									MarkdownDescription: componentsFieldsSizeDescription.MarkdownDescription,
									Optional:            true,
								},

								"theme": schema.StringAttribute{
									Description:         componentsFieldsThemeDescription.Description,
									MarkdownDescription: componentsFieldsThemeDescription.MarkdownDescription,
									Optional:            true,
								},

								"alignment": schema.StringAttribute{
									Description:         componentsFieldsAlignmentDescription.Description,
									MarkdownDescription: componentsFieldsAlignmentDescription.MarkdownDescription,
									Optional:            true,
								},

								"qr_code_type": schema.StringAttribute{
									Description:         componentsFieldsQrCodeTypeDescription.Description,
									MarkdownDescription: componentsFieldsQrCodeTypeDescription.MarkdownDescription,
									Optional:            true,
								},

								"show_border": schema.BoolAttribute{
									Description:         componentsFieldsShowBorderDescription.Description,
									MarkdownDescription: componentsFieldsShowBorderDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,
								},
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
			},

			"mark_optional": schema.BoolAttribute{
				Description:         markOptionalDescription.Description,
				MarkdownDescription: markOptionalDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
			},

			"mark_required": schema.BoolAttribute{
				Description:         markRequiredDescription.Description,
				MarkdownDescription: markRequiredDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
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
func formFieldValidationDocumentation(key string) string {

	requiredTypes := []string{}
	optionalTypes := []string{}

	for formFieldType, formField := range formComponentsFieldsSchemaDefMap {
		if slices.Contains(formField.Required, key) {
			requiredTypes = append(requiredTypes, string(formFieldType))
		}

		if slices.Contains(formField.Optional, key) {
			optionalTypes = append(optionalTypes, string(formFieldType))
		}
	}

	returnVar := ""

	descriptionSet := false

	if len(requiredTypes) > 0 {
		slices.Sort(requiredTypes)
		returnVar += fmt.Sprintf("**Required** when the `type` is one of `%s`", strings.Join(requiredTypes, "`, `"))
		descriptionSet = true
	}

	if len(optionalTypes) > 0 {
		slices.Sort(optionalTypes)
		if descriptionSet {
			returnVar += ", o"
		} else {
			returnVar += "O"
		}
		returnVar += fmt.Sprintf("ptional when the `type` is one of `%s`", strings.Join(optionalTypes, "`, `"))
		descriptionSet = true
	}

	if descriptionSet {
		returnVar += "."
	}

	return returnVar
}

func (r *FormResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data formResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.validate(ctx, true)...)
}

func (r *formComponentsFieldResourceModel) validateFieldSet(field string) bool {

	switch field {
	case "alignment":
		return !r.Alignment.IsNull()
	case "attribute_disabled":
		return !r.AttributeDisabled.IsNull()
	case "content":
		return !r.Content.IsNull()
	case "key":
		return !r.Key.IsNull()
	case "label":
		return !r.Label.IsNull()
	case "label_mode":
		return !r.LabelMode.IsNull()
	case "label_password_verify":
		return !r.LabelPasswordVerify.IsNull()
	case "layout":
		return !r.Layout.IsNull()
	case "options":
		return !r.Options.IsNull()
	case "other_option_attribute_disabled":
		return !r.OtherOptionAttributeDisabled.IsNull()
	case "other_option_enabled":
		return !r.OtherOptionEnabled.IsNull()
	case "other_option_input_label":
		return !r.OtherOptionInputLabel.IsNull()
	case "other_option_key":
		return !r.OtherOptionKey.IsNull()
	case "other_option_label":
		return !r.OtherOptionLabel.IsNull()
	case "position":
		return !r.Position.IsNull()
	case "qr_code_type":
		return !r.QrCodeType.IsNull()
	case "required":
		return !r.Required.IsNull()
	case "show_border":
		return !r.ShowBorder.IsNull()
	case "show_password_requirements":
		return !r.ShowPasswordRequirements.IsNull()
	case "size":
		return !r.Size.IsNull()
	case "styles":
		return !r.Styles.IsNull()
	case "theme":
		return !r.Theme.IsNull()
	case "type":
		return !r.Type.IsNull()
	case "validation":
		return !r.Validation.IsNull()
	default:
		return false
	}
}

// ModifyPlan
func (r *FormResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var data, modifiedData []formComponentsFieldResourceModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("components").AtName("fields"), &data)...)

	modifiedPlan := false
	fieldTypes := make([]string, 0)

	for _, field := range data {

		// attribute_disabled default
		if field.AttributeDisabled.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_CHECKBOX), string(management.ENUMFORMFIELDTYPE_COMBOBOX), string(management.ENUMFORMFIELDTYPE_DROPDOWN), string(management.ENUMFORMFIELDTYPE_RADIO), string(management.ENUMFORMFIELDTYPE_PASSWORD), string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY), string(management.ENUMFORMFIELDTYPE_TEXT):
				field.AttributeDisabled = types.BoolValue(false)
			default:
				field.AttributeDisabled = types.BoolNull()
			}
			modifiedPlan = true
		}

		// required default
		if field.Required.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_CHECKBOX), string(management.ENUMFORMFIELDTYPE_COMBOBOX), string(management.ENUMFORMFIELDTYPE_DROPDOWN), string(management.ENUMFORMFIELDTYPE_RADIO), string(management.ENUMFORMFIELDTYPE_PASSWORD), string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY), string(management.ENUMFORMFIELDTYPE_TEXT):
				field.Required = types.BoolValue(false)
			default:
				field.Required = types.BoolNull()
			}
			modifiedPlan = true
		}

		// show_password_requirements default
		if field.ShowPasswordRequirements.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_PASSWORD), string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY):
				field.ShowPasswordRequirements = types.BoolValue(false)
			default:
				field.ShowPasswordRequirements = types.BoolNull()
			}
			modifiedPlan = true
		}

		// show_border default
		if field.ShowBorder.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_QR_CODE):
				field.ShowBorder = types.BoolValue(false)
			default:
				field.ShowBorder = types.BoolNull()
			}
			modifiedPlan = true
		}

		// other_option_attribute_disabled default
		if field.OtherOptionAttributeDisabled.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_CHECKBOX), string(management.ENUMFORMFIELDTYPE_COMBOBOX), string(management.ENUMFORMFIELDTYPE_DROPDOWN), string(management.ENUMFORMFIELDTYPE_RADIO), string(management.ENUMFORMFIELDTYPE_PASSWORD), string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY), string(management.ENUMFORMFIELDTYPE_TEXT):
				field.OtherOptionAttributeDisabled = types.BoolValue(false)
			default:
				field.OtherOptionAttributeDisabled = types.BoolNull()
			}
			modifiedPlan = true
		}

		// other_option_enabled default
		if field.OtherOptionEnabled.IsUnknown() {
			switch field.Type.ValueString() {
			case string(management.ENUMFORMFIELDTYPE_CHECKBOX), string(management.ENUMFORMFIELDTYPE_COMBOBOX), string(management.ENUMFORMFIELDTYPE_DROPDOWN), string(management.ENUMFORMFIELDTYPE_RADIO), string(management.ENUMFORMFIELDTYPE_PASSWORD), string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY), string(management.ENUMFORMFIELDTYPE_TEXT):
				field.OtherOptionEnabled = types.BoolValue(false)
			default:
				field.OtherOptionEnabled = types.BoolNull()
			}
			modifiedPlan = true
		}

		// other_option_input_label default
		if field.OtherOptionInputLabel.IsUnknown() {
			field.OtherOptionInputLabel = types.StringNull()
			modifiedPlan = true
		}

		// other_option_key default
		if field.OtherOptionKey.IsUnknown() {
			field.OtherOptionKey = types.StringNull()
			modifiedPlan = true
		}

		// other_option_label default
		if field.OtherOptionLabel.IsUnknown() {
			field.OtherOptionLabel = types.StringNull()
			modifiedPlan = true
		}

		modifiedData = append(modifiedData, field)

		if !slices.Contains(fieldTypes, field.Type.ValueString()) {
			fieldTypes = append(fieldTypes, field.Type.ValueString())
		}
	}

	if modifiedPlan {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("components").AtName("fields"), modifiedData)...)
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("field_types"), fieldTypes)...)
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

func (p *formResourceModel) validate(ctx context.Context, allowUnknowns bool) diag.Diagnostics {
	var diags diag.Diagnostics

	var componentsPlan *formComponentsResourceModel
	diags.Append(p.Components.As(ctx, &componentsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: allowUnknowns,
	})...)
	if diags.HasError() {
		return diags
	}

	if componentsPlan != nil {

		var fieldsPlan []formComponentsFieldResourceModel
		diags.Append(componentsPlan.Fields.ElementsAs(ctx, &fieldsPlan, allowUnknowns)...)
		if diags.HasError() {
			return diags
		}

		if len(fieldsPlan) > 0 {

			hasSubmitButton := false
			submitButtonUnknown := false

			for i, field := range fieldsPlan {

				if field.Type.IsUnknown() && !allowUnknowns {
					diags.AddAttributeError(
						path.Root("components").AtName("fields"),
						"Invalid DaVinci form configuration",
						"The `type` parameter is unknown and cannot be validated.",
					)
					submitButtonUnknown = true
					continue
				}

				if !field.Type.IsNull() && !field.Type.IsUnknown() && field.Type.Equal(types.StringValue(string(management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON))) {
					hasSubmitButton = true
				}

				// Validate Position conflicts
				var positionPlan *formComponentsFieldPositionResourceModel
				diags.Append(field.Position.As(ctx, &positionPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: allowUnknowns,
				})...)
				if diags.HasError() {
					return diags
				}

				if positionPlan != nil {
					for existingPositionIndex, existingPosition := range fieldsPlan {
						if existingPositionIndex != i {
							var existingPositionPlan *formComponentsFieldPositionResourceModel
							diags.Append(existingPosition.Position.As(ctx, &existingPositionPlan, basetypes.ObjectAsOptions{
								UnhandledNullAsEmpty:    false,
								UnhandledUnknownAsEmpty: allowUnknowns,
							})...)
							if diags.HasError() {
								return diags
							}

							if existingPositionPlan != nil && existingPositionPlan.Col.Equal(positionPlan.Col) && existingPositionPlan.Row.Equal(positionPlan.Row) {
								diags.AddAttributeError(
									path.Root("components").AtName("fields"),
									"Invalid DaVinci form configuration",
									fmt.Sprintf("The combination of `col` and `row` must be unique between form fields.  The position `col`: `%d`, `row`: `%d` is duplicated.", positionPlan.Col.ValueInt64(), positionPlan.Row.ValueInt64()),
								)
							}
						}
					}
				}

				// Validate Required/Optional
				if v, ok := formComponentsFieldsSchemaDefMap[management.EnumFormFieldType(field.Type.ValueString())]; ok {
					for _, requiredField := range v.Required {
						if !field.validateFieldSet(requiredField) {
							diags.AddAttributeError(
								path.Root("components").AtName("fields"),
								"Invalid DaVinci form configuration",
								fmt.Sprintf("The `%s` field is required for the `%s` field type.", requiredField, field.Type.ValueString()),
							)
						}
					}
				} else {
					diags.AddAttributeWarning(
						path.Root("components").AtName("fields"),
						"Cannot validate form configuration",
						fmt.Sprintf("The form field type `%s` does not have Required/Optional metadata configured.  Please report this to the provider maintainers.", field.Type.ValueString()),
					)
				}

				// Validate parameters must have specific values if the `key` is a user field
				m, err := regexp.Compile(`^user\..+$`)
				if err != nil {
					diags.AddError(
						"Unexpected error",
						fmt.Sprintf("Failed to compile regex: %s.  This is always a bug in the provider.  Please report this error to the provider maintainers.", err.Error()),
					)
					return diags
				}

				if field.Key.IsUnknown() && !allowUnknowns {
					diags.AddAttributeError(
						path.Root("components").AtName("fields"),
						"Invalid DaVinci form configuration",
						"The `key` parameter is unknown and cannot be validated.",
					)
					continue
				}

				if !field.Key.IsNull() && !field.Key.IsUnknown() && m.MatchString(field.Key.ValueString()) {
					if !field.Required.IsNull() && !field.Required.IsUnknown() && !field.Required.ValueBool() {
						diags.AddAttributeError(
							path.Root("components").AtName("fields"),
							"Invalid DaVinci form configuration",
							fmt.Sprintf("The `required` parameter must be set to `true` for the `%s` field type when the `key` is a user field.", field.Type.ValueString()),
						)
					}
				}

				// Validate if PASSWORD or PASSWORDVERIFY, the validation.type must be NONE
				if field.Type.Equal(types.StringValue(string(management.ENUMFORMFIELDTYPE_PASSWORD))) || field.Type.Equal(types.StringValue(string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY))) {
					if field.Validation.IsUnknown() && !allowUnknowns {
						diags.AddAttributeError(
							path.Root("components").AtName("fields"),
							"Invalid DaVinci form configuration",
							"The `validation` parameter is unknown and cannot be validated.",
						)
						continue
					}

					if !field.Validation.IsNull() && !field.Validation.IsUnknown() {
						var vPlan formComponentsFieldElementValidationResourceModel
						diags.Append(field.Validation.As(ctx, &vPlan, basetypes.ObjectAsOptions{
							UnhandledNullAsEmpty:    false,
							UnhandledUnknownAsEmpty: false,
						})...)
						if diags.HasError() {
							return diags
						}

						if !vPlan.Type.Equal(types.StringValue(string(management.ENUMFORMELEMENTVALIDATIONTYPE_NONE))) {
							diags.AddAttributeError(
								path.Root("components").AtName("fields"),
								"Invalid DaVinci form configuration",
								fmt.Sprintf("The `validation.type` parameter must be set to `NONE` for the `%s` field type.", field.Type.ValueString()),
							)
						}
					}
				}
			}

			// Validate has submit button
			if !hasSubmitButton && !submitButtonUnknown {
				diags.AddAttributeError(
					path.Root("components").AtName("fields"),
					"Invalid DaVinci form configuration",
					"The DaVinci form is expected to contain a submit button field (`type` parameter value of `SUBMIT_BUTTON`).",
				)
			}
		}
	}

	return diags
}

func (p *formResourceModel) expand(ctx context.Context) (*management.Form, diag.Diagnostics) {
	var diags diag.Diagnostics

	var componentsPlan *formComponentsResourceModel
	diags.Append(p.Components.As(ctx, &componentsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	var componentsFieldsPlan []formComponentsFieldResourceModel
	diags.Append(componentsPlan.Fields.ElementsAs(ctx, &componentsFieldsPlan, false)...)
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
		var plan []types.String
		diags.Append(p.FieldTypes.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		forms, d := framework.TFTypeStringSliceToStringSlice(plan, path.Root("field_types"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		fieldTypes := make([]management.EnumFormFieldType, 0)
		for _, v := range forms {
			fieldTypes = append(fieldTypes, management.EnumFormFieldType(v))
		}

		data.SetFieldTypes(fieldTypes)
	}

	if !p.TranslationMethod.IsNull() && !p.TranslationMethod.IsUnknown() {
		data.SetTranslationMethod(management.EnumFormTranslationMethod(p.TranslationMethod.ValueString()))
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expand(ctx context.Context) (*management.FormField, diag.Diagnostics) {
	var d, diags diag.Diagnostics

	data := &management.FormField{}

	var positionPlan formComponentsFieldPositionResourceModel
	diags.Append(p.Position.As(ctx, &positionPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	positionData := positionPlan.expand()

	switch p.Type.ValueString() {
	case string(management.ENUMFORMFIELDTYPE_CHECKBOX):
		data.FormFieldCheckbox, d = p.expandFieldCheckbox(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_COMBOBOX):
		data.FormFieldCombobox, d = p.expandFieldCombobox(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_DIVIDER):
		data.FormFieldDivider = p.expandItemDivider(positionData)
	case string(management.ENUMFORMFIELDTYPE_DROPDOWN):
		data.FormFieldDropdown, d = p.expandFieldDropdown(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_EMPTY_FIELD):
		data.FormFieldEmptyField = p.expandItemEmptyField(positionData)
	case string(management.ENUMFORMFIELDTYPE_ERROR_DISPLAY):
		data.FormFieldErrorDisplay = p.expandItemErrorDisplay(positionData)
	case string(management.ENUMFORMFIELDTYPE_FLOW_BUTTON):
		data.FormFieldFlowButton, d = p.expandItemFlowButton(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_FLOW_LINK):
		data.FormFieldFlowLink, d = p.expandItemFlowLink(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_PASSWORD):
		data.FormFieldPassword, d = p.expandFieldPassword(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY):
		data.FormFieldPasswordVerify, d = p.expandFieldPasswordVerify(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_QR_CODE):
		data.FormFieldQrCode = p.expandItemQRCode(positionData)
	case string(management.ENUMFORMFIELDTYPE_RADIO):
		data.FormFieldRadio, d = p.expandFieldRadio(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_RECAPTCHA_V2):
		data.FormFieldRecaptchaV2 = p.expandItemRecaptchaV2(positionData)
	case string(management.ENUMFORMFIELDTYPE_SLATE_TEXTBLOB):
		data.FormFieldSlateTextblob = p.expandItemSlateTextblob(positionData)
	case string(management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON):
		data.FormFieldSubmitButton, d = p.expandFieldSubmitButton(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_TEXT):
		data.FormFieldText, d = p.expandFieldText(ctx, positionData)
	case string(management.ENUMFORMFIELDTYPE_TEXTBLOB):
		data.FormFieldTextblob = p.expandItemTextblob(positionData)
	}

	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return data, diags
}

func (p *formComponentsFieldPositionResourceModel) expand() *management.FormFieldCommonPosition {

	data := management.NewFormFieldCommonPosition(
		int32(p.Row.ValueInt64()),
		int32(p.Col.ValueInt64()),
	)

	if !p.Width.IsNull() && !p.Width.IsUnknown() {
		data.SetWidth(int32(p.Width.ValueInt64()))
	}

	return data
}

func (p *formComponentsFieldElementValidationResourceModel) expand() *management.FormElementValidation {

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

	return data
}

func (p *formComponentsFieldResourceModel) expandFieldCheckbox(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldCheckbox, diag.Diagnostics) {
	var diags diag.Diagnostics

	var optionsPlan []formComponentsFieldElementOptionsResourceModel
	diags.Append(p.Options.ElementsAs(ctx, &optionsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	options := make([]management.FormElementOption, 0)
	for _, v := range optionsPlan {
		optionsObj := *management.NewFormElementOption(
			v.Label.ValueString(),
			v.Value.ValueString(),
		)

		options = append(options, optionsObj)
	}

	data := management.NewFormFieldCheckbox(
		management.ENUMFORMFIELDTYPE_CHECKBOX,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
		management.EnumFormElementLayout(p.Layout.ValueString()),
		options,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldCombobox(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldCombobox, diag.Diagnostics) {
	var diags diag.Diagnostics

	var optionsPlan []formComponentsFieldElementOptionsResourceModel
	diags.Append(p.Options.ElementsAs(ctx, &optionsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	options := make([]management.FormElementOption, 0)
	for _, v := range optionsPlan {
		optionsObj := *management.NewFormElementOption(
			v.Label.ValueString(),
			v.Value.ValueString(),
		)

		options = append(options, optionsObj)
	}

	data := management.NewFormFieldCombobox(
		management.ENUMFORMFIELDTYPE_COMBOBOX,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
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

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandItemDivider(positionData *management.FormFieldCommonPosition) *management.FormFieldDivider {

	data := management.NewFormFieldDivider(
		management.ENUMFORMFIELDTYPE_DIVIDER,
		*positionData,
	)

	return data
}

func (p *formComponentsFieldResourceModel) expandFieldDropdown(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldDropdown, diag.Diagnostics) {
	var diags diag.Diagnostics

	var optionsPlan []formComponentsFieldElementOptionsResourceModel
	diags.Append(p.Options.ElementsAs(ctx, &optionsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	options := make([]management.FormElementOption, 0)
	for _, v := range optionsPlan {
		optionsObj := *management.NewFormElementOption(
			v.Label.ValueString(),
			v.Value.ValueString(),
		)

		options = append(options, optionsObj)
	}

	data := management.NewFormFieldDropdown(
		management.ENUMFORMFIELDTYPE_DROPDOWN,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
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

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandItemEmptyField(positionData *management.FormFieldCommonPosition) *management.FormFieldEmptyField {

	data := management.NewFormFieldEmptyField(
		management.ENUMFORMFIELDTYPE_EMPTY_FIELD,
		*positionData,
	)

	return data
}

func (p *formComponentsFieldResourceModel) expandItemErrorDisplay(positionData *management.FormFieldCommonPosition) *management.FormFieldErrorDisplay {

	data := management.NewFormFieldErrorDisplay(
		management.ENUMFORMFIELDTYPE_ERROR_DISPLAY,
		*positionData,
	)

	return data
}

func (p *formComponentsFieldResourceModel) expandItemFlowButton(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldFlowButton, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldFlowButton(
		management.ENUMFORMFIELDTYPE_FLOW_BUTTON,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldStylesResourceModel
		diags.Append(p.Styles.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx, "BUTTON")
		diags.Append(d...)

		if v, ok := stylesData.(*management.FormStyles); ok {
			data.SetStyles(*v)
		}
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandItemFlowLink(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldFlowLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldFlowLink(
		management.ENUMFORMFIELDTYPE_FLOW_LINK,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldStylesResourceModel
		diags.Append(p.Styles.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx, "FLOW_LINK")
		diags.Append(d...)

		if v, ok := stylesData.(*management.FormFlowLinkStyles); ok {
			data.SetStyles(*v)
		}
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldPassword(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldPassword, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldPassword(
		management.ENUMFORMFIELDTYPE_PASSWORD,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
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

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	if !p.ShowPasswordRequirements.IsNull() && !p.ShowPasswordRequirements.IsUnknown() {
		data.SetShowPasswordRequirements(p.ShowPasswordRequirements.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		diags.Append(p.Validation.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		if !plan.Type.Equal(types.StringValue(string(management.ENUMFORMELEMENTVALIDATIONTYPE_NONE))) {
			diags.AddAttributeError(
				path.Root("components").AtName("fields"),
				"Invalid DaVinci form configuration",
				fmt.Sprintf("The `validation.type` parameter must be set to `NONE` for the `%s` field type.", string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY)),
			)
		}

		validationData := plan.expand()

		data.SetValidation(*validationData)
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldPasswordVerify(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldPasswordVerify, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldPasswordVerify(
		management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.LabelPasswordVerify.IsNull() && !p.LabelPasswordVerify.IsUnknown() {
		data.SetLabelPasswordVerify(p.LabelPasswordVerify.ValueString())
	}

	if !p.Layout.IsNull() && !p.Layout.IsUnknown() {
		data.SetLayout(management.EnumFormElementLayout(p.Layout.ValueString()))
	}

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	if !p.ShowPasswordRequirements.IsNull() && !p.ShowPasswordRequirements.IsUnknown() {
		data.SetShowPasswordRequirements(p.ShowPasswordRequirements.ValueBool())
	}

	if !p.Validation.IsNull() && !p.Validation.IsUnknown() {
		var plan formComponentsFieldElementValidationResourceModel
		diags.Append(p.Validation.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		if !plan.Type.Equal(types.StringValue(string(management.ENUMFORMELEMENTVALIDATIONTYPE_NONE))) {
			diags.AddAttributeError(
				path.Root("components").AtName("fields"),
				"Invalid DaVinci form configuration",
				fmt.Sprintf("The `validation.type` parameter must be set to `NONE` for the `%s` field type.", string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY)),
			)
		}

		validationData := plan.expand()

		data.SetValidation(*validationData)
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldRadio(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldRadio, diag.Diagnostics) {
	var diags diag.Diagnostics

	var optionsPlan []formComponentsFieldElementOptionsResourceModel
	diags.Append(p.Options.ElementsAs(ctx, &optionsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	options := make([]management.FormElementOption, 0)
	for _, v := range optionsPlan {
		optionsObj := *management.NewFormElementOption(
			v.Label.ValueString(),
			v.Value.ValueString(),
		)

		options = append(options, optionsObj)
	}

	data := management.NewFormFieldRadio(
		management.ENUMFORMFIELDTYPE_RADIO,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
		management.EnumFormElementLayout(p.Layout.ValueString()),
		options,
	)

	if !p.AttributeDisabled.IsNull() && !p.AttributeDisabled.IsUnknown() {
		data.SetAttributeDisabled(p.AttributeDisabled.ValueBool())
	}

	if !p.LabelMode.IsNull() && !p.LabelMode.IsUnknown() {
		data.SetLabelMode(management.EnumFormElementLabelMode(p.LabelMode.ValueString()))
	}

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldSubmitButton(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldSubmitButton, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewFormFieldSubmitButton(
		management.ENUMFORMFIELDTYPE_SUBMIT_BUTTON,
		*positionData,
		p.Label.ValueString(),
	)

	if !p.Styles.IsNull() && !p.Styles.IsUnknown() {
		var plan formComponentsFieldStylesResourceModel
		diags.Append(p.Styles.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx, "BUTTON")
		diags.Append(d...)

		if v, ok := stylesData.(*management.FormStyles); ok {
			data.SetStyles(*v)
		}
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandFieldText(ctx context.Context, positionData *management.FormFieldCommonPosition) (*management.FormFieldText, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan formComponentsFieldElementValidationResourceModel
	diags.Append(p.Validation.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	validationData := plan.expand()

	data := management.NewFormFieldText(
		management.ENUMFORMFIELDTYPE_TEXT,
		*positionData,
		p.Key.ValueString(),
		p.Label.ValueString(),
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

	if !p.Required.IsNull() && !p.Required.IsUnknown() {
		data.SetRequired(p.Required.ValueBool())
	}

	return data, diags
}

func (p *formComponentsFieldResourceModel) expandItemQRCode(positionData *management.FormFieldCommonPosition) *management.FormFieldQrCode {
	data := management.NewFormFieldQrCode(
		management.ENUMFORMFIELDTYPE_QR_CODE,
		*positionData,
		p.Key.ValueString(),
		management.EnumFormQrCodeType(p.QrCodeType.ValueString()),
		management.EnumFormItemAlignment(p.Alignment.ValueString()),
	)

	if !p.ShowBorder.IsNull() && !p.ShowBorder.IsUnknown() {
		data.SetShowBorder(p.ShowBorder.ValueBool())
	}

	return data
}

func (p *formComponentsFieldResourceModel) expandItemRecaptchaV2(positionData *management.FormFieldCommonPosition) *management.FormFieldRecaptchaV2 {
	data := management.NewFormFieldRecaptchaV2(
		management.ENUMFORMFIELDTYPE_RECAPTCHA_V2,
		*positionData,
		management.EnumFormRecaptchaV2Size(p.Size.ValueString()),
		management.EnumFormRecaptchaV2Theme(p.Theme.ValueString()),
		management.EnumFormItemAlignment(p.Alignment.ValueString()),
	)

	return data
}

func (p *formComponentsFieldResourceModel) expandItemSlateTextblob(positionData *management.FormFieldCommonPosition) *management.FormFieldSlateTextblob {
	data := management.NewFormFieldSlateTextblob(
		management.ENUMFORMFIELDTYPE_SLATE_TEXTBLOB,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data
}

func (p *formComponentsFieldResourceModel) expandItemTextblob(positionData *management.FormFieldCommonPosition) *management.FormFieldTextblob {

	data := management.NewFormFieldTextblob(
		management.ENUMFORMFIELDTYPE_TEXTBLOB,
		*positionData,
	)

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	return data
}

func (p *formComponentsFieldStylesResourceModel) expand(ctx context.Context, styleType string) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if styleType == "BUTTON" {
		data := management.NewFormStyles()

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

		if !p.Height.IsNull() && !p.Height.IsUnknown() {
			data.SetHeight(int32(p.Height.ValueInt64()))
		}

		if !p.Padding.IsNull() && !p.Padding.IsUnknown() {
			var plan formComponentsFieldButtonStylesPaddingResourceModel
			diags.Append(p.Padding.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			padding := management.NewFormStylesPadding()

			if !plan.Bottom.IsNull() && !plan.Bottom.IsUnknown() {
				padding.SetBottom(int32(plan.Bottom.ValueInt64()))
			}

			if !plan.Left.IsNull() && !plan.Left.IsUnknown() {
				padding.SetLeft(int32(plan.Left.ValueInt64()))
			}

			if !plan.Right.IsNull() && !plan.Right.IsUnknown() {
				padding.SetRight(int32(plan.Right.ValueInt64()))
			}

			if !plan.Top.IsNull() && !plan.Top.IsUnknown() {
				padding.SetTop(int32(plan.Top.ValueInt64()))
			}

			data.SetPadding(*padding)
		}

		if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
			data.SetTextColor(p.TextColor.ValueString())
		}

		if !p.Width.IsNull() && !p.Width.IsUnknown() {
			data.SetWidth(int32(p.Width.ValueInt64()))
		}

		if !p.WidthUnit.IsNull() && !p.WidthUnit.IsUnknown() {
			data.SetWidthUnit(management.EnumFormStylesWidthUnit(p.WidthUnit.ValueString()))
		}

		return data, diags
	}

	if styleType == "FLOW_LINK" {
		data := management.NewFormFlowLinkStyles()

		if !p.Alignment.IsNull() && !p.Alignment.IsUnknown() {
			data.SetAlignment(management.EnumFormItemAlignment(p.Alignment.ValueString()))
		}

		if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
			data.SetEnabled(p.Enabled.ValueBool())
		}

		if !p.Padding.IsNull() && !p.Padding.IsUnknown() {
			var plan formComponentsFieldButtonStylesPaddingResourceModel
			diags.Append(p.Padding.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			padding := management.NewFormStylesPadding()

			if !plan.Bottom.IsNull() && !plan.Bottom.IsUnknown() {
				padding.SetBottom(int32(plan.Bottom.ValueInt64()))
			}

			if !plan.Left.IsNull() && !plan.Left.IsUnknown() {
				padding.SetLeft(int32(plan.Left.ValueInt64()))
			}

			if !plan.Right.IsNull() && !plan.Right.IsUnknown() {
				padding.SetRight(int32(plan.Right.ValueInt64()))
			}

			if !plan.Top.IsNull() && !plan.Top.IsUnknown() {
				padding.SetTop(int32(plan.Top.ValueInt64()))
			}

			data.SetPadding(*padding)
		}

		if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
			data.SetTextColor(p.TextColor.ValueString())
		}

		return data, diags
	}

	diags.AddError(
		"Unhandled style type",
		fmt.Sprintf("Unhandled style type %s.  This is a bug in the provider and must be reported to the provider maintainers.", styleType),
	)

	return nil, diags
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
		var plan formComponentsFieldStylesResourceModel
		diags.Append(p.Styles.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData, d := plan.expand(ctx, "BUTTON")
		diags.Append(d...)

		if v, ok := stylesData.(*management.FormStyles); ok {
			data.SetStyles(*v)
		}
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

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Category = framework.EnumOkToTF(apiObject.GetCategoryOk())
	p.Cols = framework.Int32OkToInt64TF(apiObject.GetColsOk())
	p.FieldTypes = framework.EnumSetOkToTF(apiObject.GetFieldTypesOk())
	p.MarkOptional = framework.BoolOkToTF(apiObject.GetMarkOptionalOk())
	p.MarkRequired = framework.BoolOkToTF(apiObject.GetMarkRequiredOk())
	p.TranslationMethod = framework.EnumOkToTF(apiObject.GetTranslationMethodOk())

	var d diag.Diagnostics
	p.Components, d = formComponentsOkToTF(apiObject.GetComponentsOk())
	diags.Append(d...)

	return diags
}

func formComponentsOkToTF(apiObject *management.FormComponents, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsTFObjectTypes), diags
	}

	fields, d := formComponentsFieldsOkToTF(apiObject.GetFieldsOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(formComponentsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsTFObjectTypes, map[string]attr.Value{
		"fields": fields,
	})
	diags.Append(d...)

	return objValue, diags
}

func formComponentsFieldsOkToTF(apiObject []management.FormField, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: formComponentsFieldsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	objectList := []attr.Value{}
	for _, v := range apiObject {

		attributeMap := map[string]attr.Value{}

		fieldObject := v.GetActualInstance()

		switch t := fieldObject.(type) {
		case *management.FormFieldCheckbox:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldCombobox:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldDivider:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"position": position,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldDropdown:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldEmptyField:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"position": position,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldErrorDisplay:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"position": position,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldFlowButton:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"key":      framework.StringOkToTF(t.GetKeyOk()),
				"label":    framework.StringOkToTF(t.GetLabelOk()),
				"position": position,
				"styles":   styles,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldFlowLink:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsFlowLinkStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"key":      framework.StringOkToTF(t.GetKeyOk()),
				"label":    framework.StringOkToTF(t.GetLabelOk()),
				"position": position,
				"styles":   styles,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldPassword:
			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"show_password_requirements":      framework.BoolOkToTF(t.GetShowPasswordRequirementsOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldPasswordVerify:
			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           framework.StringOkToTF(t.GetLabelPasswordVerifyOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"show_password_requirements":      framework.BoolOkToTF(t.GetShowPasswordRequirementsOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldQrCode:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"alignment":    framework.EnumOkToTF(t.GetAlignmentOk()),
				"key":          framework.StringOkToTF(t.GetKeyOk()),
				"position":     position,
				"qr_code_type": framework.EnumOkToTF(t.GetQrCodeTypeOk()),
				"show_border":  framework.BoolOkToTF(t.GetShowBorderOk()),
				"type":         framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldRadio:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldRecaptchaV2:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"alignment": framework.EnumOkToTF(t.GetAlignmentOk()),
				"position":  position,
				"size":      framework.EnumOkToTF(t.GetSizeOk()),
				"theme":     framework.EnumOkToTF(t.GetThemeOk()),
				"type":      framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldSlateTextblob:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"content":  framework.StringOkToTF(t.GetContentOk()),
				"position": position,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldSubmitButton:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"key":      framework.StringOkToTF(t.GetKeyOk()),
				"label":    framework.StringOkToTF(t.GetLabelOk()),
				"position": position,
				"styles":   styles,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}

		case *management.FormFieldText:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          framework.EnumOkToTF(t.GetLayoutOk()),
				"options":                         options,
				"other_option_attribute_disabled": framework.BoolOkToTF(t.GetOtherOptionAttributeDisabledOk()),
				"other_option_enabled":            framework.BoolOkToTF(t.GetOtherOptionEnabledOk()),
				"other_option_input_label":        framework.StringOkToTF(t.GetOtherOptionInputLabelOk()),
				"other_option_key":                framework.StringOkToTF(t.GetOtherOptionKeyOk()),
				"other_option_label":              framework.StringOkToTF(t.GetOtherOptionLabelOk()),
				"position":                        position,
				"required":                        framework.BoolOkToTF(t.GetRequiredOk()),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldTextblob:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"content":  framework.StringOkToTF(t.GetContentOk()),
				"position": position,
				"type":     framework.EnumOkToTF(t.GetTypeOk()),
			}
		}

		attributeMap = formComponentsFieldsConvertEmptyValuesToTFNulls(attributeMap)

		objValue, d := types.ObjectValue(formComponentsFieldsTFObjectTypes, attributeMap)
		diags.Append(d...)
		objectList = append(objectList, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

func formComponentsFieldsConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"alignment":                       types.StringNull(),
		"attribute_disabled":              types.BoolNull(),
		"content":                         types.StringNull(),
		"key":                             types.StringNull(),
		"label_mode":                      types.StringNull(),
		"label_password_verify":           types.StringNull(),
		"label":                           types.StringNull(),
		"layout":                          types.StringNull(),
		"options":                         types.SetNull(types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}),
		"other_option_attribute_disabled": types.BoolNull(),
		"other_option_enabled":            types.BoolNull(),
		"other_option_input_label":        types.StringNull(),
		"other_option_key":                types.StringNull(),
		"other_option_label":              types.StringNull(),
		"position":                        types.ObjectNull(formComponentsFieldsPositionTFObjectTypes),
		"qr_code_type":                    types.StringNull(),
		"required":                        types.BoolNull(),
		"show_border":                     types.BoolNull(),
		"show_password_requirements":      types.BoolNull(),
		"size":                            types.StringNull(),
		"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
		"theme":                           types.StringNull(),
		"type":                            types.StringNull(),
		"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func formComponentsFieldsPositionOkToTF(apiObject *management.FormFieldCommonPosition, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsPositionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsPositionTFObjectTypes, map[string]attr.Value{
		"col":   framework.Int32OkToInt64TF(apiObject.GetColOk()),
		"row":   framework.Int32OkToInt64TF(apiObject.GetRowOk()),
		"width": framework.Int32OkToInt64TF(apiObject.GetWidthOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func formComponentsFieldsElementValidationOkToTF(apiObject *management.FormElementValidation, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsFieldElementValidationTFObjectTypes, map[string]attr.Value{
		"regex":         framework.StringOkToTF(apiObject.GetRegexOk()),
		"type":          framework.EnumOkToTF(apiObject.GetTypeOk()),
		"error_message": framework.StringOkToTF(apiObject.GetErrorMessageOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func formComponentsFieldsElementOptionsOkToTF(apiObject []management.FormElementOption, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	objectAttrTypes := []attr.Value{}
	for _, v := range apiObject {
		objValue, d := types.ObjectValue(formComponentsFieldsFieldElementOptionTFObjectTypes, map[string]attr.Value{
			"label": framework.StringOkToTF(v.GetLabelOk()),
			"value": framework.StringOkToTF(v.GetValueOk()),
		})
		diags.Append(d...)

		objectAttrTypes = append(objectAttrTypes, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectAttrTypes)
	diags.Append(d...)

	return returnVar, diags
}

func formComponentsFieldsStylesOkToTF(apiObject *management.FormStyles, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes), diags
	}

	padding, d := formComponentsFieldsPaddingOkToTF(apiObject.GetPaddingOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsFieldStylesTFObjectTypes, map[string]attr.Value{
		"alignment":        framework.EnumOkToTF(apiObject.GetAlignmentOk()),
		"background_color": framework.StringOkToTF(apiObject.GetBackgroundColorOk()),
		"border_color":     framework.StringOkToTF(apiObject.GetBorderColorOk()),
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"height":           framework.Int32OkToInt64TF(apiObject.GetHeightOk()),
		"padding":          padding,
		"text_color":       framework.StringOkToTF(apiObject.GetTextColorOk()),
		"width":            framework.Int32OkToInt64TF(apiObject.GetWidthOk()),
		"width_unit":       framework.EnumOkToTF(apiObject.GetWidthUnitOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func formComponentsFieldsFlowLinkStylesOkToTF(apiObject *management.FormFlowLinkStyles, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes), diags
	}

	padding, d := formComponentsFieldsPaddingOkToTF(apiObject.GetPaddingOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsFieldStylesTFObjectTypes, map[string]attr.Value{
		"alignment":        framework.EnumOkToTF(apiObject.GetAlignmentOk()),
		"background_color": types.StringNull(),
		"border_color":     types.StringNull(),
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"height":           types.Int64Null(),
		"padding":          padding,
		"text_color":       framework.StringOkToTF(apiObject.GetTextColorOk()),
		"width":            types.Int64Null(),
		"width_unit":       types.StringNull(),
	})
	diags.Append(d...)

	return objValue, diags
}

func formComponentsFieldsPaddingOkToTF(apiObject *management.FormStylesPadding, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsFieldStylesPaddingTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsFieldStylesPaddingTFObjectTypes, map[string]attr.Value{
		"bottom": framework.Int32OkToInt64TF(apiObject.GetBottomOk()),
		"left":   framework.Int32OkToInt64TF(apiObject.GetLeftOk()),
		"right":  framework.Int32OkToInt64TF(apiObject.GetRightOk()),
		"top":    framework.Int32OkToInt64TF(apiObject.GetTopOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
