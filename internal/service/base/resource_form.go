package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
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
	Required                     types.Bool   `tfsdk:"required"`
	ShowPasswordRequirements     types.Bool   `tfsdk:"show_password_requirements"`
	Styles                       types.Object `tfsdk:"styles"`
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

var (
	// Form Components
	formComponentsTFObjectTypes = map[string]attr.Type{
		"fields": types.SetType{ElemType: types.ObjectType{
			AttrTypes: formComponentsFieldsTFObjectTypes,
		}},
	}

	// Form Components Fields
	formComponentsFieldsTFObjectTypes = map[string]attr.Type{
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
		"required":                        types.BoolType,
		"show_password_requirements":      types.BoolType,
		"styles":                          types.ObjectType{AttrTypes: formComponentsFieldsFieldStylesTFObjectTypes},
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

	// Form Components Fields Field Password Verify
	formComponentsFieldsFieldPasswordVerifyTFObjectTypes = map[string]attr.Type{
		"label_password_verify": types.StringType,
	}

	// Form Components Fields Field Combobox
	formComponentsFieldsFieldComboboxTFObjectTypes = map[string]attr.Type{}

	// Form Components Fields Field Item
	formComponentsFieldsFieldItemTFObjectTypes = map[string]attr.Type{
		"content": types.StringType,
	}

	// Form Components Fields Field Button
	formComponentsFieldsFieldButtonTFObjectTypes = map[string]attr.Type{
		"key":    types.StringType,
		"label":  types.StringType,
		"styles": types.ObjectType{AttrTypes: formComponentsFieldsFieldStylesTFObjectTypes},
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

	// Form Components Fields Field Flow Link
	formComponentsFieldsFieldFlowLinkTFObjectTypes = map[string]attr.Type{
		"key":    types.StringType,
		"label":  types.StringType,
		"styles": types.ObjectType{AttrTypes: formComponentsFieldsFieldFlowLinkStylesTFObjectTypes},
	}

	// Form Components Fields Field Flow Link Styles
	formComponentsFieldsFieldFlowLinkStylesTFObjectTypes = map[string]attr.Type{
		"horizontal_alignment": types.StringType,
		"text_color":           types.StringType,
		"enabled":              types.BoolType,
	}

	// Form Components Fields Field Recaptcha V2
	formComponentsFieldsFieldRecaptchaV2TFObjectTypes = map[string]attr.Type{
		"key":       types.StringType,
		"size":      types.StringType,
		"theme":     types.StringType,
		"alignment": types.StringType,
	}

	// Form Components Fields Field Qr Code
	formComponentsFieldsFieldQrCodeTFObjectTypes = map[string]attr.Type{
		"qr_code_type": types.StringType,
		"alignment":    types.StringType,
		"show_border":  types.BoolType,
	}

	// Form Components Fields Field Social Login Button
	formComponentsFieldsFieldSocialLoginButtonTFObjectTypes = map[string]attr.Type{
		"label":       types.StringType,
		"styles":      types.ObjectType{AttrTypes: formComponentsFieldsFieldSocialLoginButtonStylesTFObjectTypes},
		"idp_type":    types.StringType,
		"idp_name":    types.StringType,
		"idp_id":      types.StringType,
		"idp_enabled": types.BoolType,
		"icon_src":    types.StringType,
		"width":       types.Int64Type,
	}

	// Form Components Fields Field Social Login Button Styles
	formComponentsFieldsFieldSocialLoginButtonStylesTFObjectTypes = map[string]attr.Type{
		"horizontal_alignment": types.StringType,
		"text_color":           types.StringType,
		"enabled":              types.BoolType,
	}
)

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
	const colsMinValue = 1
	const colsMaxValue = 4
	const rowMaxValue = 50
	const colMinValue = 0
	const colMaxValue = 3

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
		"A single object that specifies the position of the form field in the form.",
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
	).AllowedValuesEnum(management.AllowedEnumFormFieldTypeEnumValues)

	componentsFieldsAttributeDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the linked directory attribute is disabled.",
	).RequiresReplace()

	componentsFieldsContentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an identifier for the field component.",
	)

	componentsFieldsLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the field label.",
	)

	componentsFieldsLabelPasswordVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that when a second field for verifies password is used, this property specifies the field label for that verify field.",
	)

	componentsFieldsLabelModeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how the field is rendered.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLabelModeEnumValues)

	componentsFieldsLayoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies layout attributes for radio button and checkbox fields.",
	).AllowedValuesEnum(management.AllowedEnumFormElementLayoutEnumValues)

	componentsFieldsOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array of objects that specifies the unique list of options.",
	)

	componentsFieldsOptionsLabelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsOptionsValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsRequiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the field is required.",
	)

	componentsFieldsValidationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
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
		"",
	)

	componentsFieldsStylesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes style settings for the button.",
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
		"",
	)

	componentsFieldsStylesPaddingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsStylesPaddingTopDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsStylesPaddingRightDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsStylesPaddingBottomDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsStylesPaddingLeftDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	componentsFieldsStylesTextColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the button text color. The value must be a valid hexadecimal color.",
	)

	componentsFieldsStylesWidthUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

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
	).AllowedValuesEnum(management.AllowedEnumFormTranslationMethodEnumValues)

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
								},

								"attribute_disabled": schema.BoolAttribute{
									Description:         componentsFieldsAttributeDisabledDescription.Description,
									MarkdownDescription: componentsFieldsAttributeDisabledDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									// TODO: Validator can't be false if key is "user.username" or similar
								},

								"content": schema.StringAttribute{
									Description:         componentsFieldsContentDescription.Description,
									MarkdownDescription: componentsFieldsContentDescription.MarkdownDescription,
									Optional:            true,

									// TODO Functional validator
								},

								"key": schema.StringAttribute{
									Description:         componentsFieldsKeyDescription.Description,
									MarkdownDescription: componentsFieldsKeyDescription.MarkdownDescription,
									Optional:            true,

									// TODO Functional validator
								},

								"label": schema.StringAttribute{
									Description:         componentsFieldsLabelDescription.Description,
									MarkdownDescription: componentsFieldsLabelDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator
								},

								"label_password_verify": schema.StringAttribute{
									Description:         componentsFieldsLabelPasswordVerifyDescription.Description,
									MarkdownDescription: componentsFieldsLabelPasswordVerifyDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator
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

									// TODO functional validator

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormElementLayoutEnumValues)...),
									},
								},

								"options": schema.SetNestedAttribute{
									Description:         componentsFieldsOptionsDescription.Description,
									MarkdownDescription: componentsFieldsOptionsDescription.MarkdownDescription,
									Optional:            true,

									// TODO functional validator

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

									// TODO: Validator can't be false if key is "user.username" or similar
								},

								"validation": schema.SingleNestedAttribute{
									Description:         componentsFieldsValidationDescription.Description,
									MarkdownDescription: componentsFieldsValidationDescription.MarkdownDescription,
									Optional:            true,

									// TODO: optional validator

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
									Optional:            true,
									Computed:            true,

									// TODO: functional validator
								},

								"other_option_key": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionKeyDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionKeyDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator
								},

								"other_option_label": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionLabelDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionLabelDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator
								},

								"other_option_input_label": schema.StringAttribute{
									Description:         componentsFieldsOtherOptionInputLabelDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionInputLabelDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator
								},

								"other_option_attribute_disabled": schema.BoolAttribute{
									Description:         componentsFieldsOtherOptionAttributeDisabledDescription.Description,
									MarkdownDescription: componentsFieldsOtherOptionAttributeDisabledDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									// TODO: functional validator
								},

								"show_password_requirements": schema.BoolAttribute{
									Description:         componentsFieldsShowPasswordRequirementsDescription.Description,
									MarkdownDescription: componentsFieldsShowPasswordRequirementsDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									// TODO: functional validator
								},

								"styles": schema.SingleNestedAttribute{
									Description:         componentsFieldsStylesDescription.Description,
									MarkdownDescription: componentsFieldsStylesDescription.MarkdownDescription,
									Optional:            true,

									// TODO: functional validator

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

			"language_bundle": schema.MapAttribute{
				Description:         languageBundleDescription.Description,
				MarkdownDescription: languageBundleDescription.MarkdownDescription,
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

				"enabled": schema.BoolAttribute{
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

				"enabled": schema.BoolAttribute{
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

		"icon_src": schema.StringAttribute{
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
		var plan []string
		diags.Append(p.FieldTypes.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		fieldTypes := make([]management.EnumFormFieldType, 0)
		for _, v := range plan {
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
		data.FormFieldPassword = p.expandFieldPassword(positionData)
	case string(management.ENUMFORMFIELDTYPE_PASSWORD_VERIFY):
		data.FormFieldPasswordVerify = p.expandFieldPasswordVerify(positionData)
	case string(management.ENUMFORMFIELDTYPE_RADIO):
		data.FormFieldRadio, d = p.expandFieldRadio(ctx, positionData)
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

func (p *formComponentsFieldResourceModel) expandFieldPassword(positionData *management.FormFieldCommonPosition) *management.FormFieldPassword {

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

	return data
}

func (p *formComponentsFieldResourceModel) expandFieldPasswordVerify(positionData *management.FormFieldCommonPosition) *management.FormFieldPasswordVerify {

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

	return data
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
		diags.Append(p.Styles.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		stylesData := plan.expand()

		data.SetStyles(*stylesData)
	}

	if !p.Width.IsNull() && !p.Width.IsUnknown() {
		data.SetWidth(int32(p.Width.ValueInt64()))
	}

	return data, diags
}

func (p *formComponentsFieldSocialLoginButtonStylesResourceModel) expand() *management.FormSocialLoginButtonStyles {
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

	return data
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
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
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
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldDivider:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
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
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
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
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldEmptyField:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
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
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
			}

		case *management.FormFieldErrorDisplay:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
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
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
			}

		case *management.FormFieldFlowButton:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              types.BoolNull(),
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      types.StringNull(),
				"label_password_verify":           types.StringNull(),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          types.StringNull(),
				"options":                         types.SetNull(types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}),
				"other_option_attribute_disabled": types.BoolNull(),
				"other_option_enabled":            types.BoolNull(),
				"other_option_input_label":        types.StringNull(),
				"other_option_key":                types.StringNull(),
				"other_option_label":              types.StringNull(),
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          styles,
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
			}

		case *management.FormFieldFlowLink:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsFlowLinkStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              types.BoolNull(),
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      types.StringNull(),
				"label_password_verify":           types.StringNull(),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          types.StringNull(),
				"options":                         types.SetNull(types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}),
				"other_option_attribute_disabled": types.BoolNull(),
				"other_option_enabled":            types.BoolNull(),
				"other_option_input_label":        types.StringNull(),
				"other_option_key":                types.StringNull(),
				"other_option_label":              types.StringNull(),
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          styles,
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
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
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
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
				"content":                         types.StringNull(),
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
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldQrCode:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)
			attributeMap["position"] = position
			attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
			//attributeMap["field_qr_code"], d = formComponentsFieldsFieldQrCodeToTF(t)

		case *management.FormFieldRadio:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			validation, d := formComponentsFieldsElementValidationOkToTF(t.GetValidationOk())
			diags.Append(d...)

			options, d := formComponentsFieldsElementOptionsOkToTF(t.GetOptionsOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              framework.BoolOkToTF(t.GetAttributeDisabledOk()),
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldRecaptchaV2:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)
			attributeMap["position"] = position
			attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
			//attributeMap["field_recaptcha_v2"], d = formComponentsFieldsFieldRecaptchaV2ToTF(t)

		case *management.FormFieldSlateTextblob:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              types.BoolNull(),
				"content":                         framework.StringOkToTF(t.GetContentOk()),
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
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
			}

		case *management.FormFieldSocialLoginButton:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)
			attributeMap["position"] = position
			attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
			//attributeMap["field_social_login_button"], d = formComponentsFieldsFieldSocialLoginButtonToTF(t)

		case *management.FormFieldSubmitButton:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			styles, d := formComponentsFieldsStylesOkToTF(t.GetStylesOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              types.BoolNull(),
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      types.StringNull(),
				"label_password_verify":           types.StringNull(),
				"label":                           framework.StringOkToTF(t.GetLabelOk()),
				"layout":                          types.StringNull(),
				"options":                         types.SetNull(types.ObjectType{AttrTypes: formComponentsFieldsFieldElementOptionTFObjectTypes}),
				"other_option_attribute_disabled": types.BoolNull(),
				"other_option_enabled":            types.BoolNull(),
				"other_option_input_label":        types.StringNull(),
				"other_option_key":                types.StringNull(),
				"other_option_label":              types.StringNull(),
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          styles,
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
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
				"content":                         types.StringNull(),
				"key":                             framework.StringOkToTF(t.GetKeyOk()),
				"label_mode":                      framework.EnumOkToTF(t.GetLabelModeOk()),
				"label_password_verify":           types.StringNull(),
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
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      validation,
			}

		case *management.FormFieldTextblob:
			position, d := formComponentsFieldsPositionOkToTF(t.GetPositionOk())
			diags.Append(d...)

			attributeMap = map[string]attr.Value{
				"attribute_disabled":              types.BoolNull(),
				"content":                         framework.StringOkToTF(t.GetContentOk()),
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
				"position":                        position,
				"required":                        types.BoolNull(),
				"show_password_requirements":      types.BoolNull(),
				"styles":                          types.ObjectNull(formComponentsFieldsFieldStylesTFObjectTypes),
				"type":                            framework.EnumOkToTF(t.GetTypeOk()),
				"validation":                      types.ObjectNull(formComponentsFieldsFieldElementValidationTFObjectTypes),
			}
		}

		objValue, d := types.ObjectValue(formComponentsFieldsTFObjectTypes, attributeMap)
		diags.Append(d...)
		objectList = append(objectList, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

func formComponentsFieldsPositionOkToTF(apiObject *management.FormFieldCommonPosition, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(formComponentsFieldsPositionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(formComponentsFieldsPositionTFObjectTypes, map[string]attr.Value{
		"col":   framework.Int32OkToTF(apiObject.GetColOk()),
		"row":   framework.Int32OkToTF(apiObject.GetRowOk()),
		"width": framework.Int32OkToTF(apiObject.GetWidthOk()),
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
		"height":           framework.Int32OkToTF(apiObject.GetHeightOk()),
		"padding":          padding,
		"text_color":       framework.StringOkToTF(apiObject.GetTextColorOk()),
		"width":            framework.Int32OkToTF(apiObject.GetWidthOk()),
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
		"bottom": framework.Int32OkToTF(apiObject.GetBottomOk()),
		"left":   framework.Int32OkToTF(apiObject.GetLeftOk()),
		"right":  framework.Int32OkToTF(apiObject.GetRightOk()),
		"top":    framework.Int32OkToTF(apiObject.GetTopOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
