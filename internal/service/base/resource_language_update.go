package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceLanguageUpdate() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to complete the configuration of customer-defined and system-defined languages in PingOne.",

		CreateContext: resourceLanguageUpdateCreate,
		ReadContext:   resourceLanguageUpdateRead,
		UpdateContext: resourceLanguageUpdateUpdate,
		DeleteContext: resourceLanguageUpdateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLanguageUpdateImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment that contains the language to manage.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"language_id": {
				Description:      "The ID of the langauge in PingOne to update.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"locale": {
				Description: "An ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The language name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Specifies whether this language is enabled for the environment.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"default": {
				Description: "Specifies whether this language is the default for the environment.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"customer_added": {
				Description: "Specifies whether this language was added by a customer administrator.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourceLanguageUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId(d.Get("language_id").(string))

	return resourceLanguageUpdateUpdate(ctx, d, meta)
}

func resourceLanguageUpdateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.ReadOneLanguage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneLanguage",
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

	respObject := resp.(*management.Language)

	d.Set("locale", respObject.GetLocale())
	d.Set("enabled", respObject.GetEnabled())
	d.Set("default", respObject.GetDefault())

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetCustomerAddedOk(); ok {
		d.Set("customer_added", v)
	} else {
		d.Set("customer_added", nil)
	}

	return diags
}

func resourceLanguageUpdateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.ReadOneLanguage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneLanguage",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	language := resp.(*management.Language)

	language.SetEnabled(d.Get("enabled").(bool))
	language.SetDefault(false)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.UpdateLanguage(ctx, d.Get("environment_id").(string), d.Id()).Language(*language).Execute()
		},
		"UpdateLanguage",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if v, ok := d.Get("default").(bool); ok && v && d.Get("enabled").(bool) {

		language.SetDefault(v)

		_, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.LanguagesApi.UpdateLanguage(ctx, d.Get("environment_id").(string), d.Id()).Language(*language).Execute()
			},
			"UpdateLanguage",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}

	}

	return resourceLanguageUpdateRead(ctx, d, meta)
}

func resourceLanguageUpdateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.ReadOneLanguage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneLanguage",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	language := resp.(*management.Language)

	language.SetEnabled(false)
	language.SetDefault(false)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.UpdateLanguage(ctx, d.Get("environment_id").(string), d.Id()).Language(*language).Execute()
		},
		"UpdateLanguage",
		sdk.DefaultCustomError,
		func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

			// Gateway errors
			if r.StatusCode >= 502 && r.StatusCode <= 504 {
				tflog.Warn(ctx, "Gateway error detected, available for retry")
				return true
			}

			if p1error != nil {
				var err error

				// Permissions may not have propagated by this point
				if m, err := regexp.MatchString("^The default language cannot be disabled.", p1error.GetMessage()); err == nil && m {
					tflog.Warn(ctx, "Default language collision .. attempting reset ..")
					err = resetDefaultLanguage(ctx, apiClient, d.Get("environment_id").(string))
					if err != nil {
						tflog.Error(ctx, fmt.Sprintf("%s", err))
						return false
					}

					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}

			}

			return false
		},
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceLanguageUpdateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/languageID\"", d.Id())
	}

	environmentID, languageID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(languageID)

	resourceLanguageUpdateRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func resetDefaultLanguage(ctx context.Context, apiClient *management.APIClient, environmentID string) error {

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.ReadLanguages(ctx, environmentID).Execute()
		},
		"ReadLanguages",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return fmt.Errorf("Error looking up languages. Resp: %+v", resp)
	}

	respObject := resp.(*management.EntityArray)

	var enLanguage management.Language

	if languages, ok := respObject.Embedded.GetLanguagesOk(); ok {
		fmt.Printf("\n\nHERE1!!!! %+v", languages)

		found := false
		for _, language := range languages {

			// Set back to En
			if language.Language.GetLocale() == "en" {
				enLanguage = *language.Language
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Cannot find language `en`")
		}

	}

	fmt.Printf("\n\nHERE!!!! %+v", enLanguage)

	return nil
}
