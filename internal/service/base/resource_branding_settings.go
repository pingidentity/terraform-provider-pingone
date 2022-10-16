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

func ResourceBrandingSettings() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the PingOne branding settings for an environment.",

		CreateContext: resourceBrandingSettingsCreate,
		ReadContext:   resourceBrandingSettingsRead,
		UpdateContext: resourceBrandingSettingsUpdate,
		DeleteContext: resourceBrandingSettingsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBrandingSettingsImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to set branding settings for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"company_name": {
				Description:      "The company name associated with the specified environment.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"logo_image": {
				Description: "The HREF and the ID for the company logo.",
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
		},
	}
}

func resourceBrandingSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	brandingSettings := *management.NewBrandingSettings()

	if v, ok := d.GetOk("company_name"); ok {
		brandingSettings.SetCompanyName(v.(string))
	} else {
		brandingSettings.SetCompanyName("")
	}

	if v, ok := d.GetOk("logo_image"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			brandingSettings.SetLogo(*management.NewBrandingSettingsLogo(attrs["href"].(string), attrs["id"].(string)))
		}
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, d.Get("environment_id").(string)).BrandingSettings(brandingSettings).Execute()
		},
		"UpdateBrandingSettings-Create",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.BrandingSettings)

	d.SetId(respObject.GetId())

	return resourceBrandingSettingsRead(ctx, d, meta)
}

func resourceBrandingSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingSettingsApi.ReadBrandingSettings(ctx, d.Get("environment_id").(string)).Execute()
		},
		"ReadBrandingSettings",
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

	respObject := resp.(*management.BrandingSettings)

	if v, ok := respObject.GetCompanyNameOk(); ok {
		d.Set("company_name", v)
	} else {
		d.Set("company_name", nil)
	}

	if v, ok := respObject.GetLogoOk(); ok {
		d.Set("logo_image", flattenBrandingSettingsLogo(v))
	} else {
		d.Set("logo_image", nil)
	}

	return diags
}

func resourceBrandingSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	brandingSettings := *management.NewBrandingSettings()

	if v, ok := d.GetOk("company_name"); ok {
		brandingSettings.SetCompanyName(v.(string))
	} else {
		brandingSettings.SetCompanyName("")
	}

	if v, ok := d.GetOk("logo_image"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			brandingSettings.SetLogo(*management.NewBrandingSettingsLogo(attrs["href"].(string), attrs["id"].(string)))
		}
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, d.Get("environment_id").(string)).BrandingSettings(brandingSettings).Execute()
		},
		"UpdateBrandingSettings-Update",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceBrandingSettingsRead(ctx, d, meta)
}

func resourceBrandingSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	brandingSettings := *management.NewBrandingSettings()
	brandingSettings.SetCompanyName("")

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, d.Get("environment_id").(string)).BrandingSettings(brandingSettings).Execute()
		},
		"UpdateBrandingSettings-Delete",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceBrandingSettingsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 1
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID\"", d.Id())
	}

	environmentID := attributes[0]

	d.Set("environment_id", environmentID)
	d.SetId(environmentID)

	resourceBrandingSettingsRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func flattenBrandingSettingsLogo(s *management.BrandingSettingsLogo) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}
