package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceBrandingTheme() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne branding themes for an environment.",

		CreateContext: resourceBrandingThemeCreate,
		ReadContext:   resourceBrandingThemeRead,
		UpdateContext: resourceBrandingThemeUpdate,
		DeleteContext: resourceBrandingThemeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBrandingThemeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to set branding settings for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "The name of the branding theme.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"template": {
				Description:      fmt.Sprintf("The template name of the branding theme associated with the environment. Options are `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMBRANDINGTHEMETEMPLATE_DEFAULT), string(management.ENUMBRANDINGTHEMETEMPLATE_FOCUS), string(management.ENUMBRANDINGTHEMETEMPLATE_MURAL), string(management.ENUMBRANDINGTHEMETEMPLATE_SLATE), string(management.ENUMBRANDINGTHEMETEMPLATE_SPLIT)),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMBRANDINGTHEMETEMPLATE_DEFAULT), string(management.ENUMBRANDINGTHEMETEMPLATE_FOCUS), string(management.ENUMBRANDINGTHEMETEMPLATE_MURAL), string(management.ENUMBRANDINGTHEMETEMPLATE_SLATE), string(management.ENUMBRANDINGTHEMETEMPLATE_SPLIT)}, false)),
			},
			"default": {
				Description: "Specifies whether this theme is the environment's default branding configuration.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"logo": {
				Description: "The HREF and the ID for the company logo, for this branding template.  If not set, the environment's default logo (set with the `pingone_branding_settings` resource) will be applied.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:      "The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"href": {
							Description:      "The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
			"background_image": {
				Description:  "The HREF and the ID for the background image.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"background_image", "background_color", "use_default_background"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:      "The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"href": {
							Description:      "The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
			"background_color": {
				Description:  "The background color for the theme. It must be a valid hexadecimal color code.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"background_image", "background_color", "use_default_background"},
			},
			"use_default_background": {
				Description:  "A boolean to specify that the background should be set to the theme template's default.",
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"background_image", "background_color", "use_default_background"},
			},
			"body_text_color": {
				Description: "The body text color for the theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"button_color": {
				Description: "The button color for the theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"button_text_color": {
				Description: "The button text color for the branding theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"card_color": {
				Description: "The card color for the branding theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"footer_text": {
				Description: "The text to be displayed in the footer of the branding theme.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"heading_text_color": {
				Description: "The heading text color for the branding theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"link_text_color": {
				Description: "The hyperlink text color for the branding theme. It must be a valid hexadecimal color code.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceBrandingThemeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	brandingTheme := expandBrandingTheme(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingThemesApi.CreateBrandingTheme(ctx, d.Get("environment_id").(string)).BrandingTheme(*brandingTheme).Execute()
		},
		"CreateBrandingTheme",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.BrandingTheme)

	d.SetId(respObject.GetId())

	return resourceBrandingThemeRead(ctx, d, meta)
}

func resourceBrandingThemeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingThemesApi.ReadOneBrandingTheme(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneBrandingTheme",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.BrandingTheme)

	d.Set("template", string(respObject.GetTemplate()))
	d.Set("default", respObject.GetDefault())

	if v, ok := respObject.GetConfigurationOk(); ok {
		d.Set("name", v.GetName())

		if v1, ok := v.GetLogoOk(); ok {
			d.Set("logo", flattenBrandingThemeLogo(v1))
		} else {
			d.Set("logo", nil)
		}

		if v1, ok := v.GetBackgroundImageOk(); ok {
			d.Set("background_image", flattenBrandingThemeBackgroundImage(v1))
		} else {
			d.Set("background_image", nil)
		}

		if v1, ok := v.GetBackgroundColorOk(); ok {
			d.Set("background_color", v1)
		} else {
			d.Set("background_color", nil)
		}

		if v1, ok := v.GetBackgroundTypeOk(); ok && *v1 == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_DEFAULT {
			d.Set("use_default_background", true)
		} else {
			d.Set("use_default_background", false)
		}

		if v1, ok := v.GetBodyTextColorOk(); ok {
			d.Set("body_text_color", v1)
		} else {
			d.Set("body_text_color", nil)
		}

		if v1, ok := v.GetButtonColorOk(); ok {
			d.Set("button_color", v1)
		} else {
			d.Set("button_color", nil)
		}

		if v1, ok := v.GetButtonTextColorOk(); ok {
			d.Set("button_text_color", v1)
		} else {
			d.Set("button_text_color", nil)
		}

		if v1, ok := v.GetCardColorOk(); ok {
			d.Set("card_color", v1)
		} else {
			d.Set("card_color", nil)
		}

		if v1, ok := v.GetFooterOk(); ok {
			d.Set("footer_text", v1)
		} else {
			d.Set("footer_text", nil)
		}

		if v1, ok := v.GetHeadingTextColorOk(); ok {
			d.Set("heading_text_color", v1)
		} else {
			d.Set("heading_text_color", nil)
		}

		if v1, ok := v.GetLinkTextColorOk(); ok {
			d.Set("link_text_color", v1)
		} else {
			d.Set("link_text_color", nil)
		}

	} else {
		d.Set("name", nil)
		d.Set("logo", nil)
		d.Set("background_image", nil)
		d.Set("background_color", nil)
		d.Set("use_default_background", nil)
		d.Set("body_text_color", nil)
		d.Set("button_color", nil)
		d.Set("button_text_color", nil)
		d.Set("card_color", nil)
		d.Set("footer_text", nil)
		d.Set("heading_text_color", nil)
		d.Set("link_text_color", nil)
	}

	return diags
}

func resourceBrandingThemeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	brandingTheme := expandBrandingTheme(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingThemesApi.UpdateBrandingTheme(ctx, d.Get("environment_id").(string), d.Id()).BrandingTheme(*brandingTheme).Execute()
		},
		"UpdateBrandingTheme",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceBrandingThemeRead(ctx, d, meta)
}

func resourceBrandingThemeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.BrandingThemesApi.DeleteBrandingTheme(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteBrandingTheme",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceBrandingThemeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/brandingThemeID\"", d.Id())
	}

	environmentID, brandingThemeID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(brandingThemeID)

	resourceBrandingThemeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandBrandingTheme(d *schema.ResourceData) *management.BrandingTheme {

	logoType := management.ENUMBRANDINGLOGOTYPE_NONE
	var logo map[string]interface{}

	if v, ok := d.GetOk("logo"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			logo = j[0].(map[string]interface{})
			logoType = management.ENUMBRANDINGLOGOTYPE_IMAGE
		}
	}

	backgroundType := management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_NONE
	var background map[string]interface{}
	if v, ok := d.GetOk("background_image"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			background = j[0].(map[string]interface{})
			backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_IMAGE
		}
	}

	var backgroundColour string
	if v, ok := d.GetOk("background_color"); ok {
		backgroundColour = v.(string)
		backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_COLOR
	}

	if v, ok := d.GetOk("use_default_background"); ok && v.(bool) {
		backgroundType = management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_DEFAULT
	}

	configuration := *management.NewBrandingThemeConfiguration(
		backgroundType,
		d.Get("body_text_color").(string),
		d.Get("button_color").(string),
		d.Get("button_text_color").(string),
		d.Get("card_color").(string),
		d.Get("heading_text_color").(string),
		d.Get("link_text_color").(string),
		logoType,
	)

	configuration.SetName(d.Get("name").(string))

	if logoType == management.ENUMBRANDINGLOGOTYPE_IMAGE {
		configuration.SetLogo(*management.NewBrandingThemeConfigurationLogo(logo["href"].(string), logo["id"].(string)))
	}

	if backgroundType == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_IMAGE {
		configuration.SetBackgroundImage(*management.NewBrandingThemeConfigurationBackgroundImage(background["href"].(string), background["id"].(string)))
	}

	if backgroundType == management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_COLOR {
		configuration.SetBackgroundColor(backgroundColour)
	}

	if v, ok := d.GetOk("footer_text"); ok {
		configuration.SetFooter(v.(string))
	}

	brandingTheme := management.NewBrandingTheme(
		configuration,
		false,
		management.EnumBrandingThemeTemplate(d.Get("template").(string)),
	)

	return brandingTheme
}

func flattenBrandingThemeBackgroundImage(s *management.BrandingThemeConfigurationBackgroundImage) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenBrandingThemeLogo(s *management.BrandingThemeConfigurationLogo) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}
