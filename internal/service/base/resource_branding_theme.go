// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type BrandingThemeResource serviceClientType

var brandingThemeLocalizedTextTFObjectTypes = map[string]attr.Type{
	"content":          types.MapType{ElemType: types.ObjectType{AttrTypes: brandingThemeLocalizedTextContentTFObjectTypes}},
	"content_type":     types.StringType,
	"default_language": types.StringType,
	"enabled":          types.BoolType,
}

var brandingThemeLocalizedTextContentTFObjectTypes = map[string]attr.Type{
	"input_text": types.StringType,
}

type brandingThemeResourceModelV1 struct {
	Id                          pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId               pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                        types.String                 `tfsdk:"name"`
	Template                    types.String                 `tfsdk:"template"`
	Default                     types.Bool                   `tfsdk:"default"`
	Logo                        types.Object                 `tfsdk:"logo"`
	BackgroundImage             types.Object                 `tfsdk:"background_image"`
	BackgroundColor             types.String                 `tfsdk:"background_color"`
	UseDefaultBackground        types.Bool                   `tfsdk:"use_default_background"`
	BodyTextColor               types.String                 `tfsdk:"body_text_color"`
	ButtonColor                 types.String                 `tfsdk:"button_color"`
	ButtonTextColor             types.String                 `tfsdk:"button_text_color"`
	CardColor                   types.String                 `tfsdk:"card_color"`
	FooterText                  types.String                 `tfsdk:"footer_text"`
	HeadingTextColor            types.String                 `tfsdk:"heading_text_color"`
	LinkTextColor               types.String                 `tfsdk:"link_text_color"`
	ApplicationBackgroundColor  types.String                 `tfsdk:"application_background_color"`
	BodyTextSize                types.String                 `tfsdk:"body_text_size"`
	BodyTextWeight              types.String                 `tfsdk:"body_text_weight"`
	ButtonBorderColor           types.String                 `tfsdk:"button_border_color"`
	ButtonCornerRadius          types.String                 `tfsdk:"button_corner_radius"`
	ButtonBorderRadius          types.String                 `tfsdk:"button_border_radius"`
	ButtonBorderWidth           types.String                 `tfsdk:"button_border_width"`
	ButtonHoverStateBorderColor types.String                 `tfsdk:"button_hover_state_border_color"`
	ButtonHoverStateFillColor   types.String                 `tfsdk:"button_hover_state_fill_color"`
	ButtonHoverStateTextColor   types.String                 `tfsdk:"button_hover_state_text_color"`
	ButtonTextSize              types.String                 `tfsdk:"button_text_size"`
	ButtonTextWeight            types.String                 `tfsdk:"button_text_weight"`
	CardBorderColor             types.String                 `tfsdk:"card_border_color"`
	CardBorderWidth             types.String                 `tfsdk:"card_border_width"`
	CardCornerRadius            types.String                 `tfsdk:"card_corner_radius"`
	CardHorizontalAlignment     types.String                 `tfsdk:"card_horizontal_alignment"`
	CardShadow                  types.String                 `tfsdk:"card_shadow"`
	CardVerticalAlignment       types.String                 `tfsdk:"card_vertical_alignment"`
	FocusRectangleColor         types.String                 `tfsdk:"focus_rectangle_color"`
	FooterLocalized             types.Object                 `tfsdk:"footer_localized"`
	GlobalFont                  types.String                 `tfsdk:"global_font"`
	Header                      types.String                 `tfsdk:"header"`
	HeaderLocalized             types.Object                 `tfsdk:"header_localized"`
	HeaderBackgroundColor       types.String                 `tfsdk:"header_background_color"`
	InputBorderWidth            types.String                 `tfsdk:"input_border_width"`
	InputBoxBorderColor         types.String                 `tfsdk:"input_box_border_color"`
	InputCornerRadius           types.String                 `tfsdk:"input_corner_radius"`
	InputLabelPosition          types.String                 `tfsdk:"input_label_position"`
	InputLabelTextColor         types.String                 `tfsdk:"input_label_text_color"`
	InputLabelTextSize          types.String                 `tfsdk:"input_label_text_size"`
	InputLabelTextWeight        types.String                 `tfsdk:"input_label_text_weight"`
	InputValueTextColor         types.String                 `tfsdk:"input_value_text_color"`
	InputValueTextSize          types.String                 `tfsdk:"input_value_text_size"`
	InputValueTextWeight        types.String                 `tfsdk:"input_value_text_weight"`
	LinkTextHoverColor          types.String                 `tfsdk:"link_text_hover_color"`
	LinkTextSize                types.String                 `tfsdk:"link_text_size"`
	LinkTextWeight              types.String                 `tfsdk:"link_text_weight"`
	LogoHeight                  types.String                 `tfsdk:"logo_height"`
	SubTitleTextColor           types.String                 `tfsdk:"sub_title_text_color"`
	SubTitleTextSize            types.String                 `tfsdk:"sub_title_text_size"`
	SubTitleTextWeight          types.String                 `tfsdk:"sub_title_text_weight"`
	TitleTextColor              types.String                 `tfsdk:"title_text_color"`
	TitleTextSize               types.String                 `tfsdk:"title_text_size"`
	TitleTextWeight             types.String                 `tfsdk:"title_text_weight"`
}

type brandingThemeLocalizedTextResourceModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	DefaultLanguage types.String `tfsdk:"default_language"`
	ContentType     types.String `tfsdk:"content_type"`
	Content         types.Map    `tfsdk:"content"`
}

type brandingThemeLocalizedTextContentResourceModel struct {
	InputText types.String `tfsdk:"input_text"`
}

// Framework interfaces
var (
	_ resource.Resource                = &BrandingThemeResource{}
	_ resource.ResourceWithConfigure   = &BrandingThemeResource{}
	_ resource.ResourceWithImportState = &BrandingThemeResource{}
)

// New Object
func NewBrandingThemeResource() resource.Resource {
	return &BrandingThemeResource{}
}

// Metadata
func (r *BrandingThemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branding_theme"
}

// Schema.
func (r *BrandingThemeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const defaultBoolFalse = false

	templateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The template name of the branding theme associated with the environment.",
	).AllowedValuesEnum(management.AllowedEnumBrandingThemeTemplateEnumValues)

	backgroundColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The background color for the theme. It must be a valid hexadecimal color code.",
	).ConflictsWith([]string{"background_image"})

	useDefaultBackgroundDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify that the background should be set to the theme template's default.",
	).DefaultValue(bool(defaultBoolFalse))

	backgroundImageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the HREF and ID for the background image.",
	).ConflictsWith([]string{"background_color"})

	logoDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the HREF and ID for the company logo, for this branding template.  If not set, the environment's default logo (set with the `pingone_branding_settings` resource) will be applied.",
	)

	logoIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	logoHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.",
	)

	backgroundImageIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the background image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	backgroundImageHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the background image file used for branding.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.",
	)

	localizedContentTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The content type for the localized text.",
	).AllowedValuesEnum(management.AllowedEnumBrandingThemeLocalizedContentTypeEnumValues)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne branding themes for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to set branding settings for."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique name of the branding theme.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"template": schema.StringAttribute{
				Description:         templateDescription.Description,
				MarkdownDescription: templateDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumBrandingThemeTemplateEnumValues)...),
				},
			},

			"default": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies whether this theme is the environment's default branding configuration.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"logo": schema.SingleNestedAttribute{
				Description:         logoDescription.Description,
				MarkdownDescription: logoDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         logoIdDescription.Description,
						MarkdownDescription: logoIdDescription.MarkdownDescription,
						Required:            true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"href": schema.StringAttribute{
						Description:         logoHrefDescription.Description,
						MarkdownDescription: logoHrefDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},
				},
			},

			"background_image": schema.SingleNestedAttribute{
				Description:         backgroundImageDescription.Description,
				MarkdownDescription: backgroundImageDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("background_color"),
					),
				},

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         backgroundImageIdDescription.Description,
						MarkdownDescription: backgroundImageIdDescription.MarkdownDescription,
						Required:            true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"href": schema.StringAttribute{
						Description:         backgroundImageHrefDescription.Description,
						MarkdownDescription: backgroundImageHrefDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},
				},
			},

			"background_color": schema.StringAttribute{
				Description:         backgroundColorDescription.Description,
				MarkdownDescription: backgroundColorDescription.MarkdownDescription,
				Optional:            true,
			},

			"use_default_background": schema.BoolAttribute{
				Description:         useDefaultBackgroundDescription.Description,
				MarkdownDescription: useDefaultBackgroundDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"body_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The body text color for the theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button color for the theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button text color for the branding theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"card_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card color for the branding theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"footer_text": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The text to be displayed in the footer of the branding theme.").Description,
				Optional:    true,
			},

			"heading_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The heading text color for the branding theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"link_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The hyperlink text color for the branding theme. It must be a valid hexadecimal color code.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"application_background_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The application background color for the theme. It must be a valid hexadecimal color code.  Note that this property is not used by DaVinci forms.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"header": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The header for the theme.  For example, `<h1>Welcome to PingOne</h1>`.").Description,
				Optional:    true,
			},

			"header_background_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The header background color for the theme. It must be a valid hexadecimal color code.  Note that this property is not used by DaVinci forms.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"header_localized": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The localization object to specify language translations for the form header.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies whether localization is enabled for the header.").Description,
						Optional:    true,
						Computed:    true,

						Default: booldefault.StaticBool(false),
					},

					"default_language": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The localization language code for the header's default language (for example, `en`).").Description,
						Required:    true,
					},

					"content_type": schema.StringAttribute{
						Description:         localizedContentTypeDescription.Description,
						MarkdownDescription: localizedContentTypeDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumBrandingThemeLocalizedContentTypeEnumValues)...),
						},
					},

					"content": schema.MapNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A map of language codes to the localized text for the language options.  The map key is the language code (for example, `en`).").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"input_text": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("The header text in the specified language.").Description,
									Required:    true,
								},
							},
						},
					},
				},
			},

			"footer_localized": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The localization object to specify language translations for the form footer.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies whether localization is enabled for the footer.").Description,
						Optional:    true,
						Computed:    true,

						Default: booldefault.StaticBool(false),
					},

					"default_language": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The localization language code for the footer's default language (for example, `en`).").Description,
						Required:    true,
					},

					"content_type": schema.StringAttribute{
						Description:         localizedContentTypeDescription.Description,
						MarkdownDescription: localizedContentTypeDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumBrandingThemeLocalizedContentTypeEnumValues)...),
						},
					},

					"content": schema.MapNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A map of language codes to the localized text for the language options.  The map key is the language code (for example, `en`).").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"input_text": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("The footer text in the specified language.").Description,
									Required:    true,
								},
							},
						},
					},
				},
			},

			"body_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The body text size for the theme.").Description,
				Optional:    true,
			},

			"body_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The body text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"button_border_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button border color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_corner_radius": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button corner radius for the theme (range from 0-25px).").Description,
				Optional:    true,
			},

			"button_border_radius": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button border radius for the theme (range from 0-25px).").Description,
				Optional:    true,
			},

			"button_border_width": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button border width for the theme (range from 0-4px).").Description,
				Optional:    true,
			},

			"button_hover_state_border_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button hover state border color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_hover_state_fill_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button hover state fill color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_hover_state_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button hover state text color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"button_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button text size for the theme.").Description,
				Optional:    true,
			},

			"button_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The button text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"card_border_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card border color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"card_border_width": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card border width for the theme (range from 0-4px).").Description,
				Optional:    true,
			},

			"card_corner_radius": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card corner radius for the theme (range from 0-25px).").Description,
				Optional:    true,
			},

			"card_horizontal_alignment": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card horizontal alignment for the theme.").Description,
				Optional:    true,
			},

			"card_shadow": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card shadow for the theme.").Description,
				Optional:    true,
			},

			"card_vertical_alignment": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The card vertical alignment for the theme.").Description,
				Optional:    true,
			},

			"focus_rectangle_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The focus rectangle color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"global_font": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The global font for the theme.  The default value is `Helvetica Neue, Helvetica, sans-serif`.").Description,
				Optional:    true,
			},

			"input_border_width": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input border width for the theme (range from 0-4px).").Description,
				Optional:    true,
			},

			"input_box_border_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input box border color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"input_corner_radius": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input corner radius for the theme (range from 0-25px).").Description,
				Optional:    true,
			},

			"input_label_position": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input label position for the theme.").Description,
				Optional:    true,
			},

			"input_label_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input label text color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"input_label_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input label text size for the theme.").Description,
				Optional:    true,
			},

			"input_label_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input label text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"input_value_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input value text color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"input_value_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input value text size for the theme.").Description,
				Optional:    true,
			},

			"input_value_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The input value text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"link_text_hover_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The link text hover color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"link_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The link text size for the theme.").Description,
				Optional:    true,
			},

			"link_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The link text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"logo_height": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The logo height assigned to the image (range from 0-100px).").Description,
				Optional:    true,
			},

			"sub_title_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The subtitle text color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"sub_title_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The subtitle text size for the theme.").Description,
				Optional:    true,
			},

			"sub_title_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The subtitle text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},

			"title_text_color": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The title text color for the theme. It must be a valid hexadecimal color code.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.HexColorCode, "Value must be a valid hex color code."),
				},
			},

			"title_text_size": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The title text size for the theme.").Description,
				Optional:    true,
			},

			"title_text_weight": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The title text weight for the theme (range from 100-900).").Description,
				Optional:    true,
			},
		},
	}
}

func (r *BrandingThemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *BrandingThemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingThemeResourceModelV1

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
	brandingTheme, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.BrandingTheme
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.CreateBrandingTheme(ctx, plan.EnvironmentId.ValueString()).BrandingTheme(*brandingTheme).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateBrandingTheme",
		legacysdk.DefaultCustomError,
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

func (r *BrandingThemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *brandingThemeResourceModelV1

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
	var response *management.BrandingTheme
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.ReadOneBrandingTheme(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneBrandingTheme",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *BrandingThemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state brandingThemeResourceModelV1

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
	brandingTheme, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.BrandingTheme
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.UpdateBrandingTheme(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).BrandingTheme(*brandingTheme).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateBrandingTheme",
		legacysdk.DefaultCustomError,
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

func (r *BrandingThemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *brandingThemeResourceModelV1

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.DeleteBrandingTheme(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteBrandingTheme",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingThemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "branding_theme_id",
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

func (p *brandingThemeResourceModelV1) expand(ctx context.Context) (*management.BrandingTheme, diag.Diagnostics) {
	var diags diag.Diagnostics

	logoType := management.ENUMBRANDINGLOGOTYPE_NONE
	var logo service.ImageResourceModel

	if !p.Logo.IsNull() && !p.Logo.IsUnknown() {

		diags.Append(p.Logo.As(ctx, &logo, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		logoType = management.ENUMBRANDINGLOGOTYPE_IMAGE

	}

	backgroundType := management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_NONE
	var background service.ImageResourceModel
	if !p.BackgroundImage.IsNull() && !p.BackgroundImage.IsUnknown() {

		diags.Append(p.BackgroundImage.As(ctx, &background, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_IMAGE

	}

	var backgroundColour string
	if !p.BackgroundColor.IsNull() && !p.BackgroundColor.IsUnknown() {
		backgroundColour = p.BackgroundColor.ValueString()
		backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_COLOR
	}

	if !p.UseDefaultBackground.IsNull() && !p.UseDefaultBackground.IsUnknown() && p.UseDefaultBackground.Equal(types.BoolValue(true)) {
		backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_DEFAULT
	}

	configuration := *management.NewBrandingThemeConfiguration(
		backgroundType,
		p.BodyTextColor.ValueString(),
		p.ButtonColor.ValueString(),
		p.ButtonTextColor.ValueString(),
		p.CardColor.ValueString(),
		p.HeadingTextColor.ValueString(),
		p.LinkTextColor.ValueString(),
		logoType,
	)

	configuration.SetName(p.Name.ValueString())

	if logoType == management.ENUMBRANDINGLOGOTYPE_IMAGE {
		configuration.SetLogo(*management.NewBrandingThemeConfigurationLogo(logo.Href.ValueString(), logo.Id.ValueString()))
	}

	if !background.Href.IsNull() && !background.Href.IsUnknown() && !background.Id.IsNull() && !background.Id.IsUnknown() {
		configuration.SetBackgroundImage(*management.NewBrandingThemeConfigurationBackgroundImage(background.Href.ValueString(), background.Id.ValueString()))
	}

	if backgroundColour != "" {
		configuration.SetBackgroundColor(backgroundColour)
	}

	if !p.FooterText.IsNull() && !p.FooterText.IsUnknown() {
		configuration.SetFooter(p.FooterText.ValueString())
	}

	if !p.Header.IsNull() && !p.Header.IsUnknown() {
		configuration.SetHeader(p.Header.ValueString())
	}

	if !p.ApplicationBackgroundColor.IsNull() && !p.ApplicationBackgroundColor.IsUnknown() {
		configuration.SetApplicationBackgroundColor(p.ApplicationBackgroundColor.ValueString())
	}

	if !p.HeaderBackgroundColor.IsNull() && !p.HeaderBackgroundColor.IsUnknown() {
		configuration.SetHeaderBackgroundColor(p.HeaderBackgroundColor.ValueString())
	}

	if !p.BodyTextSize.IsNull() && !p.BodyTextSize.IsUnknown() {
		configuration.SetBodyTextSize(p.BodyTextSize.ValueString())
	}

	if !p.BodyTextWeight.IsNull() && !p.BodyTextWeight.IsUnknown() {
		configuration.SetBodyTextWeight(p.BodyTextWeight.ValueString())
	}

	if !p.ButtonBorderColor.IsNull() && !p.ButtonBorderColor.IsUnknown() {
		configuration.SetButtonBorderColor(p.ButtonBorderColor.ValueString())
	}

	if !p.ButtonCornerRadius.IsNull() && !p.ButtonCornerRadius.IsUnknown() {
		configuration.SetButtonCornerRadius(p.ButtonCornerRadius.ValueString())
	}

	if !p.ButtonBorderRadius.IsNull() && !p.ButtonBorderRadius.IsUnknown() {
		configuration.SetButtonBorderRadius(p.ButtonBorderRadius.ValueString())
	}

	if !p.ButtonBorderWidth.IsNull() && !p.ButtonBorderWidth.IsUnknown() {
		configuration.SetButtonBorderWidth(p.ButtonBorderWidth.ValueString())
	}

	if !p.ButtonHoverStateBorderColor.IsNull() && !p.ButtonHoverStateBorderColor.IsUnknown() {
		configuration.SetButtonHoverStateBorderColor(p.ButtonHoverStateBorderColor.ValueString())
	}

	if !p.ButtonHoverStateFillColor.IsNull() && !p.ButtonHoverStateFillColor.IsUnknown() {
		configuration.SetButtonHoverStateFillColor(p.ButtonHoverStateFillColor.ValueString())
	}

	if !p.ButtonHoverStateTextColor.IsNull() && !p.ButtonHoverStateTextColor.IsUnknown() {
		configuration.SetButtonHoverStateTextColor(p.ButtonHoverStateTextColor.ValueString())
	}

	if !p.ButtonTextSize.IsNull() && !p.ButtonTextSize.IsUnknown() {
		configuration.SetButtonTextSize(p.ButtonTextSize.ValueString())
	}

	if !p.ButtonTextWeight.IsNull() && !p.ButtonTextWeight.IsUnknown() {
		configuration.SetButtonTextWeight(p.ButtonTextWeight.ValueString())
	}

	if !p.CardBorderColor.IsNull() && !p.CardBorderColor.IsUnknown() {
		configuration.SetCardBorderColor(p.CardBorderColor.ValueString())
	}

	if !p.CardBorderWidth.IsNull() && !p.CardBorderWidth.IsUnknown() {
		configuration.SetCardBorderWidth(p.CardBorderWidth.ValueString())
	}

	if !p.CardCornerRadius.IsNull() && !p.CardCornerRadius.IsUnknown() {
		configuration.SetCardCornerRadius(p.CardCornerRadius.ValueString())
	}

	if !p.CardHorizontalAlignment.IsNull() && !p.CardHorizontalAlignment.IsUnknown() {
		configuration.SetCardHorizontalAlignment(p.CardHorizontalAlignment.ValueString())
	}

	if !p.CardShadow.IsNull() && !p.CardShadow.IsUnknown() {
		configuration.SetCardShadow(p.CardShadow.ValueString())
	}

	if !p.CardVerticalAlignment.IsNull() && !p.CardVerticalAlignment.IsUnknown() {
		configuration.SetCardVerticalAlignment(p.CardVerticalAlignment.ValueString())
	}

	if !p.FocusRectangleColor.IsNull() && !p.FocusRectangleColor.IsUnknown() {
		configuration.SetFocusRectangleColor(p.FocusRectangleColor.ValueString())
	}

	if !p.GlobalFont.IsNull() && !p.GlobalFont.IsUnknown() {
		configuration.SetGlobalFont(p.GlobalFont.ValueString())
	}

	if !p.InputBorderWidth.IsNull() && !p.InputBorderWidth.IsUnknown() {
		configuration.SetInputBorderWidth(p.InputBorderWidth.ValueString())
	}

	if !p.InputBoxBorderColor.IsNull() && !p.InputBoxBorderColor.IsUnknown() {
		configuration.SetInputBoxBorderColor(p.InputBoxBorderColor.ValueString())
	}

	if !p.InputCornerRadius.IsNull() && !p.InputCornerRadius.IsUnknown() {
		configuration.SetInputCornerRadius(p.InputCornerRadius.ValueString())
	}

	if !p.InputLabelPosition.IsNull() && !p.InputLabelPosition.IsUnknown() {
		configuration.SetInputLabelPosition(p.InputLabelPosition.ValueString())
	}

	if !p.InputLabelTextColor.IsNull() && !p.InputLabelTextColor.IsUnknown() {
		configuration.SetInputLabelTextColor(p.InputLabelTextColor.ValueString())
	}

	if !p.InputLabelTextSize.IsNull() && !p.InputLabelTextSize.IsUnknown() {
		configuration.SetInputLabelTextSize(p.InputLabelTextSize.ValueString())
	}

	if !p.InputLabelTextWeight.IsNull() && !p.InputLabelTextWeight.IsUnknown() {
		configuration.SetInputLabelTextWeight(p.InputLabelTextWeight.ValueString())
	}

	if !p.InputValueTextColor.IsNull() && !p.InputValueTextColor.IsUnknown() {
		configuration.SetInputValueTextColor(p.InputValueTextColor.ValueString())
	}

	if !p.InputValueTextSize.IsNull() && !p.InputValueTextSize.IsUnknown() {
		configuration.SetInputValueTextSize(p.InputValueTextSize.ValueString())
	}

	if !p.InputValueTextWeight.IsNull() && !p.InputValueTextWeight.IsUnknown() {
		configuration.SetInputValueTextWeight(p.InputValueTextWeight.ValueString())
	}

	if !p.LinkTextHoverColor.IsNull() && !p.LinkTextHoverColor.IsUnknown() {
		configuration.SetLinkTextHoverColor(p.LinkTextHoverColor.ValueString())
	}

	if !p.LinkTextSize.IsNull() && !p.LinkTextSize.IsUnknown() {
		configuration.SetLinkTextSize(p.LinkTextSize.ValueString())
	}

	if !p.LinkTextWeight.IsNull() && !p.LinkTextWeight.IsUnknown() {
		configuration.SetLinkTextWeight(p.LinkTextWeight.ValueString())
	}

	if !p.LogoHeight.IsNull() && !p.LogoHeight.IsUnknown() {
		configuration.SetLogoHeight(p.LogoHeight.ValueString())
	}

	if !p.SubTitleTextColor.IsNull() && !p.SubTitleTextColor.IsUnknown() {
		configuration.SetSubTitleTextColor(p.SubTitleTextColor.ValueString())
	}

	if !p.SubTitleTextSize.IsNull() && !p.SubTitleTextSize.IsUnknown() {
		configuration.SetSubTitleTextSize(p.SubTitleTextSize.ValueString())
	}

	if !p.SubTitleTextWeight.IsNull() && !p.SubTitleTextWeight.IsUnknown() {
		configuration.SetSubTitleTextWeight(p.SubTitleTextWeight.ValueString())
	}

	if !p.TitleTextColor.IsNull() && !p.TitleTextColor.IsUnknown() {
		configuration.SetTitleTextColor(p.TitleTextColor.ValueString())
	}

	if !p.TitleTextSize.IsNull() && !p.TitleTextSize.IsUnknown() {
		configuration.SetTitleTextSize(p.TitleTextSize.ValueString())
	}

	if !p.TitleTextWeight.IsNull() && !p.TitleTextWeight.IsUnknown() {
		configuration.SetTitleTextWeight(p.TitleTextWeight.ValueString())
	}

	if !p.HeaderLocalized.IsNull() && !p.HeaderLocalized.IsUnknown() {
		headerLocalized, d := p.expandLocalizedText(ctx, p.HeaderLocalized)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		configuration.SetHeaderLocalized(*headerLocalized)
	}

	if !p.FooterLocalized.IsNull() && !p.FooterLocalized.IsUnknown() {
		footerLocalized, d := p.expandLocalizedText(ctx, p.FooterLocalized)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		configuration.SetFooterLocalized(*footerLocalized)
	}

	data := management.NewBrandingTheme(
		configuration,
		false,
		management.EnumBrandingThemeTemplate(p.Template.ValueString()),
	)

	return data, diags
}

func (p *brandingThemeResourceModelV1) expandLocalizedText(ctx context.Context, localized types.Object) (*management.BrandingThemeConfigurationLocalizedText, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan brandingThemeLocalizedTextResourceModel

	diags.Append(localized.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	localizedData := management.NewBrandingThemeConfigurationLocalizedText()

	localizedData.SetEnabled(plan.Enabled.ValueBool())
	localizedData.SetDefaultLanguage(plan.DefaultLanguage.ValueString())
	localizedData.SetContentType(management.EnumBrandingThemeLocalizedContentType(plan.ContentType.ValueString()))

	content := make(map[string]management.BrandingThemeConfigurationLocalizedTextContent, len(plan.Content.Elements()))

	for languageCode, contentValue := range plan.Content.Elements() {
		contentObj, ok := contentValue.(basetypes.ObjectValue)
		if !ok {
			diags.AddError(
				"Unexpected Content Value Type",
				fmt.Sprintf("Expected an object value for the localized text content, got: %T. Please report this to the provider maintainers.", contentValue),
			)
			continue
		}

		var contentPlan brandingThemeLocalizedTextContentResourceModel

		diags.Append(contentObj.As(ctx, &contentPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		contentData := management.NewBrandingThemeConfigurationLocalizedTextContent()

		contentData.SetInputText(contentPlan.InputText.ValueString())

		content[languageCode] = *contentData
	}

	localizedData.SetContent(content)

	return localizedData, diags
}

func (p *brandingThemeResourceModelV1) toState(apiObject *management.BrandingTheme) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Template = framework.EnumOkToTF(apiObject.GetTemplateOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())

	if v, ok := apiObject.GetConfigurationOk(); ok {
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.BackgroundColor = framework.StringOkToTF(v.GetBackgroundColorOk())

		if v1, ok := v.GetBackgroundTypeOk(); ok && *v1 == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_DEFAULT {
			p.UseDefaultBackground = types.BoolValue(true)
		} else {
			p.UseDefaultBackground = types.BoolValue(false)
		}

		logo, d := service.ImageOkToTF(v.GetLogoOk())
		diags.Append(d...)
		p.Logo = logo

		backgroundImage, d := service.ImageOkToTF(v.GetBackgroundImageOk())
		diags.Append(d...)
		p.BackgroundImage = backgroundImage

		p.BodyTextColor = framework.StringOkToTF(v.GetBodyTextColorOk())
		p.ButtonColor = framework.StringOkToTF(v.GetButtonColorOk())
		p.ButtonTextColor = framework.StringOkToTF(v.GetButtonTextColorOk())
		p.CardColor = framework.StringOkToTF(v.GetCardColorOk())
		p.FooterText = framework.StringOkToTF(v.GetFooterOk())
		p.HeadingTextColor = framework.StringOkToTF(v.GetHeadingTextColorOk())
		p.LinkTextColor = framework.StringOkToTF(v.GetLinkTextColorOk())

		p.Header = framework.StringOkToTF(v.GetHeaderOk())
		p.ApplicationBackgroundColor = framework.StringOkToTF(v.GetApplicationBackgroundColorOk())
		p.HeaderBackgroundColor = framework.StringOkToTF(v.GetHeaderBackgroundColorOk())
		p.BodyTextSize = framework.StringOkToTF(v.GetBodyTextSizeOk())
		p.BodyTextWeight = framework.StringOkToTF(v.GetBodyTextWeightOk())
		p.ButtonBorderColor = framework.StringOkToTF(v.GetButtonBorderColorOk())
		p.ButtonCornerRadius = framework.StringOkToTF(v.GetButtonCornerRadiusOk())
		p.ButtonBorderRadius = framework.StringOkToTF(v.GetButtonBorderRadiusOk())
		p.ButtonBorderWidth = framework.StringOkToTF(v.GetButtonBorderWidthOk())
		p.ButtonHoverStateBorderColor = framework.StringOkToTF(v.GetButtonHoverStateBorderColorOk())
		p.ButtonHoverStateFillColor = framework.StringOkToTF(v.GetButtonHoverStateFillColorOk())
		p.ButtonHoverStateTextColor = framework.StringOkToTF(v.GetButtonHoverStateTextColorOk())
		p.ButtonTextSize = framework.StringOkToTF(v.GetButtonTextSizeOk())
		p.ButtonTextWeight = framework.StringOkToTF(v.GetButtonTextWeightOk())
		p.CardBorderColor = framework.StringOkToTF(v.GetCardBorderColorOk())
		p.CardBorderWidth = framework.StringOkToTF(v.GetCardBorderWidthOk())
		p.CardCornerRadius = framework.StringOkToTF(v.GetCardCornerRadiusOk())
		p.CardHorizontalAlignment = framework.StringOkToTF(v.GetCardHorizontalAlignmentOk())
		p.CardShadow = framework.StringOkToTF(v.GetCardShadowOk())
		p.CardVerticalAlignment = framework.StringOkToTF(v.GetCardVerticalAlignmentOk())
		p.FocusRectangleColor = framework.StringOkToTF(v.GetFocusRectangleColorOk())
		p.GlobalFont = framework.StringOkToTF(v.GetGlobalFontOk())
		p.InputBorderWidth = framework.StringOkToTF(v.GetInputBorderWidthOk())
		p.InputBoxBorderColor = framework.StringOkToTF(v.GetInputBoxBorderColorOk())
		p.InputCornerRadius = framework.StringOkToTF(v.GetInputCornerRadiusOk())
		p.InputLabelPosition = framework.StringOkToTF(v.GetInputLabelPositionOk())
		p.InputLabelTextColor = framework.StringOkToTF(v.GetInputLabelTextColorOk())
		p.InputLabelTextSize = framework.StringOkToTF(v.GetInputLabelTextSizeOk())
		p.InputLabelTextWeight = framework.StringOkToTF(v.GetInputLabelTextWeightOk())
		p.InputValueTextColor = framework.StringOkToTF(v.GetInputValueTextColorOk())
		p.InputValueTextSize = framework.StringOkToTF(v.GetInputValueTextSizeOk())
		p.InputValueTextWeight = framework.StringOkToTF(v.GetInputValueTextWeightOk())
		p.LinkTextHoverColor = framework.StringOkToTF(v.GetLinkTextHoverColorOk())
		p.LinkTextSize = framework.StringOkToTF(v.GetLinkTextSizeOk())
		p.LinkTextWeight = framework.StringOkToTF(v.GetLinkTextWeightOk())
		p.LogoHeight = framework.StringOkToTF(v.GetLogoHeightOk())
		p.SubTitleTextColor = framework.StringOkToTF(v.GetSubTitleTextColorOk())
		p.SubTitleTextSize = framework.StringOkToTF(v.GetSubTitleTextSizeOk())
		p.SubTitleTextWeight = framework.StringOkToTF(v.GetSubTitleTextWeightOk())
		p.TitleTextColor = framework.StringOkToTF(v.GetTitleTextColorOk())
		p.TitleTextSize = framework.StringOkToTF(v.GetTitleTextSizeOk())
		p.TitleTextWeight = framework.StringOkToTF(v.GetTitleTextWeightOk())

		p.HeaderLocalized, d = toStateBrandingThemeLocalizedTextOk(v.GetHeaderLocalizedOk())
		diags.Append(d...)
		p.FooterLocalized, d = toStateBrandingThemeLocalizedTextOk(v.GetFooterLocalizedOk())
		diags.Append(d...)
	} else {
		p.Name = types.StringNull()
		p.BackgroundColor = types.StringNull()
		p.UseDefaultBackground = types.BoolNull()
		p.BodyTextColor = types.StringNull()
		p.ButtonColor = types.StringNull()
		p.ButtonTextColor = types.StringNull()
		p.CardColor = types.StringNull()
		p.FooterText = types.StringNull()
		p.HeadingTextColor = types.StringNull()
		p.LinkTextColor = types.StringNull()

		p.Header = types.StringNull()
		p.ApplicationBackgroundColor = types.StringNull()
		p.HeaderBackgroundColor = types.StringNull()
		p.BodyTextSize = types.StringNull()
		p.BodyTextWeight = types.StringNull()
		p.ButtonBorderColor = types.StringNull()
		p.ButtonCornerRadius = types.StringNull()
		p.ButtonBorderRadius = types.StringNull()
		p.ButtonBorderWidth = types.StringNull()
		p.ButtonHoverStateBorderColor = types.StringNull()
		p.ButtonHoverStateFillColor = types.StringNull()
		p.ButtonHoverStateTextColor = types.StringNull()
		p.ButtonTextSize = types.StringNull()
		p.ButtonTextWeight = types.StringNull()
		p.CardBorderColor = types.StringNull()
		p.CardBorderWidth = types.StringNull()
		p.CardCornerRadius = types.StringNull()
		p.CardHorizontalAlignment = types.StringNull()
		p.CardShadow = types.StringNull()
		p.CardVerticalAlignment = types.StringNull()
		p.FocusRectangleColor = types.StringNull()
		p.GlobalFont = types.StringNull()
		p.InputBorderWidth = types.StringNull()
		p.InputBoxBorderColor = types.StringNull()
		p.InputCornerRadius = types.StringNull()
		p.InputLabelPosition = types.StringNull()
		p.InputLabelTextColor = types.StringNull()
		p.InputLabelTextSize = types.StringNull()
		p.InputLabelTextWeight = types.StringNull()
		p.InputValueTextColor = types.StringNull()
		p.InputValueTextSize = types.StringNull()
		p.InputValueTextWeight = types.StringNull()
		p.LinkTextHoverColor = types.StringNull()
		p.LinkTextSize = types.StringNull()
		p.LinkTextWeight = types.StringNull()
		p.LogoHeight = types.StringNull()
		p.SubTitleTextColor = types.StringNull()
		p.SubTitleTextSize = types.StringNull()
		p.SubTitleTextWeight = types.StringNull()
		p.TitleTextColor = types.StringNull()
		p.TitleTextSize = types.StringNull()
		p.TitleTextWeight = types.StringNull()

		p.HeaderLocalized = types.ObjectNull(brandingThemeLocalizedTextTFObjectTypes)
		p.FooterLocalized = types.ObjectNull(brandingThemeLocalizedTextTFObjectTypes)
	}

	return diags
}

func toStateBrandingThemeLocalizedTextOk(apiObject *management.BrandingThemeConfigurationLocalizedText, ok bool) (types.Object, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(brandingThemeLocalizedTextTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"default_language": framework.StringOkToTF(apiObject.GetDefaultLanguageOk()),
		"content_type":     framework.EnumOkToTF(apiObject.GetContentTypeOk()),
	}

	content, contentOk := apiObject.GetContentOk()
	o["content"], d = toStateBrandingThemeLocalizedTextContentOk(content, contentOk)
	diags.Append(d...)

	returnVar, d := types.ObjectValue(brandingThemeLocalizedTextTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func toStateBrandingThemeLocalizedTextContentOk(apiObject *map[string]management.BrandingThemeConfigurationLocalizedTextContent, ok bool) (types.Map, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: brandingThemeLocalizedTextContentTFObjectTypes}

	if !ok || apiObject == nil || len(*apiObject) == 0 {
		return types.MapNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{}
	for languageCode, content := range *apiObject {
		o := map[string]attr.Value{
			"input_text": framework.StringOkToTF(content.GetInputTextOk()),
		}

		objValue, d := types.ObjectValue(brandingThemeLocalizedTextContentTFObjectTypes, o)
		diags.Append(d...)

		objectMap[languageCode] = objValue
	}

	returnVar, d := types.MapValue(tfObjType, objectMap)
	diags.Append(d...)

	return returnVar, diags
}
