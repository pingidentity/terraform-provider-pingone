package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type BrandingThemeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type brandingThemeResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	EnvironmentId        types.String `tfsdk:"environment_id"`
	Name                 types.String `tfsdk:"name"`
	Template             types.String `tfsdk:"template"`
	Default              types.Bool   `tfsdk:"default"`
	Logo                 types.List   `tfsdk:"logo"`
	BackgroundImage      types.List   `tfsdk:"background_image"`
	BackgroundColor      types.String `tfsdk:"background_color"`
	UseDefaultBackground types.Bool   `tfsdk:"use_default_background"`
	BodyTextColor        types.String `tfsdk:"body_text_color"`
	ButtonColor          types.String `tfsdk:"button_color"`
	ButtonTextColor      types.String `tfsdk:"button_text_color"`
	CardColor            types.String `tfsdk:"card_color"`
	FooterText           types.String `tfsdk:"footer_text"`
	HeadingTextColor     types.String `tfsdk:"heading_text_color"`
	LinkTextColor        types.String `tfsdk:"link_text_color"`
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
		"A single block that specifies the HREF and ID for the background image.",
	).ExactlyOneOf(backgroundExactlyOneOfRelativePaths)

	logoDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies the HREF and ID for the company logo, for this branding template.  If not set, the environment's default logo (set with the `pingone_branding_settings` resource) will be applied.",
	)

	logoIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	logoHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
	)

	backgroundImageIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the background image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	backgroundImageHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the background image file used for branding.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
	)

	resp.Schema = schema.Schema{
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

		Blocks: map[string]schema.Block{

			"logo": schema.ListNestedBlock{
				Description:         logoDescription.Description,
				MarkdownDescription: logoDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         logoIdDescription.Description,
							MarkdownDescription: logoIdDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
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

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"background_image": schema.ListNestedBlock{
				Description:         backgroundImageDescription.Description,
				MarkdownDescription: backgroundImageDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         backgroundImageIdDescription.Description,
							MarkdownDescription: backgroundImageIdDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
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

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
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

func (r *BrandingThemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingThemeResourceModel

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

	// Build the model for the API
	brandingTheme, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BrandingThemesApi.CreateBrandingTheme(ctx, plan.EnvironmentId.ValueString()).BrandingTheme(*brandingTheme).Execute()
		},
		"CreateBrandingTheme",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.BrandingTheme))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *BrandingThemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *brandingThemeResourceModel

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
			return r.client.BrandingThemesApi.ReadOneBrandingTheme(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneBrandingTheme",
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
	resp.Diagnostics.Append(data.toState(response.(*management.BrandingTheme))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BrandingThemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state brandingThemeResourceModel

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

	// Build the model for the API
	brandingTheme, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BrandingThemesApi.UpdateBrandingTheme(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).BrandingTheme(*brandingTheme).Execute()
		},
		"UpdateBrandingTheme",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.BrandingTheme))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *BrandingThemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *brandingThemeResourceModel

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
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.BrandingThemesApi.DeleteBrandingTheme(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteBrandingTheme",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingThemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/branding_theme_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *brandingThemeResourceModel) expand(ctx context.Context) (*management.BrandingTheme, diag.Diagnostics) {
	var diags diag.Diagnostics

	logoType := management.ENUMBRANDINGLOGOTYPE_NONE
	var logo imageResourceModel

	if !p.Logo.IsNull() && !p.Logo.IsUnknown() {

		var plan []imageResourceModel
		diags.Append(p.Logo.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		logo = plan[0]
		logoType = management.ENUMBRANDINGLOGOTYPE_IMAGE

	}

	backgroundType := management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_NONE
	var background imageResourceModel
	if !p.BackgroundImage.IsNull() && !p.BackgroundImage.IsUnknown() {

		var plan []imageResourceModel
		diags.Append(p.BackgroundImage.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		background = plan[0]
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

func (p *brandingThemeResourceModel) toState(apiObject *management.BrandingTheme) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
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

		logo, d := toStateImageRef(v.GetLogoOk())
		diags.Append(d...)
		p.Logo = logo

		backgroundImage, d := toStateImageRef(v.GetBackgroundImageOk())
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
