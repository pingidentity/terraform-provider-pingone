// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

type brandingThemeResourceModelV1 struct {
	Id                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId        pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                 types.String                 `tfsdk:"name"`
	Template             types.String                 `tfsdk:"template"`
	Default              types.Bool                   `tfsdk:"default"`
	Logo                 types.Object                 `tfsdk:"logo"`
	BackgroundImage      types.Object                 `tfsdk:"background_image"`
	BackgroundColor      types.String                 `tfsdk:"background_color"`
	UseDefaultBackground types.Bool                   `tfsdk:"use_default_background"`
	BodyTextColor        types.String                 `tfsdk:"body_text_color"`
	ButtonColor          types.String                 `tfsdk:"button_color"`
	ButtonTextColor      types.String                 `tfsdk:"button_text_color"`
	CardColor            types.String                 `tfsdk:"card_color"`
	FooterText           types.String                 `tfsdk:"footer_text"`
	HeadingTextColor     types.String                 `tfsdk:"heading_text_color"`
	LinkTextColor        types.String                 `tfsdk:"link_text_color"`
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

	templateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The template name of the branding theme associated with the environment.",
	).AllowedValuesEnum(management.AllowedEnumBrandingThemeTemplateEnumValues)

	backgroundExactlyOneOfRelativePaths := []string{
		"background_image",
		"background_color",
		"use_default_background",
	}

	backgroundColorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The background color for the theme. It must be a valid hexadecimal color code.",
	).ExactlyOneOf(backgroundExactlyOneOfRelativePaths)

	useDefaultBackgroundDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify that the background should be set to the theme template's default.",
	).ExactlyOneOf(backgroundExactlyOneOfRelativePaths)

	backgroundImageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the HREF and ID for the background image.",
	).ExactlyOneOf(backgroundExactlyOneOfRelativePaths)

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
					boolplanmodifier.UseStateForUnknown(),
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
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("background_color"),
						path.MatchRelative().AtParent().AtName("use_default_background"),
						path.MatchRelative().AtParent().AtName("background_image"),
					),
				},
			},

			"use_default_background": schema.BoolAttribute{
				Description:         useDefaultBackgroundDescription.Description,
				MarkdownDescription: useDefaultBackgroundDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),

				Validators: []validator.Bool{
					boolvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("background_color"),
						path.MatchRelative().AtParent().AtName("use_default_background"),
						path.MatchRelative().AtParent().AtName("background_image"),
					),
				},
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

	if backgroundType == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_IMAGE {
		configuration.SetBackgroundImage(*management.NewBrandingThemeConfigurationBackgroundImage(background.Href.ValueString(), background.Id.ValueString()))
	}

	if backgroundType == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_COLOR {
		configuration.SetBackgroundColor(backgroundColour)
	}

	if !p.FooterText.IsNull() && !p.FooterText.IsUnknown() {
		configuration.SetFooter(p.FooterText.ValueString())
	}

	data := management.NewBrandingTheme(
		configuration,
		false,
		management.EnumBrandingThemeTemplate(p.Template.ValueString()),
	)

	return data, diags
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
	}

	return diags
}
