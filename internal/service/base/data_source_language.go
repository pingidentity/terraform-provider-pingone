package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceLanguage() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne language data",

		ReadContext: datasourcePingOneLanguageRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"language_id": {
				Description:      "The ID of the language in PingOne to update.",
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"locale", "language_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"locale": {
				Description:      "An ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).",
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"locale", "language_id"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(verify.FullIsoList(), false)),
			},
			"name": {
				Description: "The language name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Specifies whether this language is enabled for the environment.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"default": {
				Description: "Specifies whether this language is the default for the environment.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"customer_added": {
				Description: "Specifies whether this language was added by a customer administrator.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneLanguageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp *management.Language

	if v, ok := d.GetOk("locale"); ok {
		resp, diags = findLanguageByLocale(ctx, apiClient, d.Get("environment_id").(string), v.(string))
		if diags.HasError() {
			return diags
		}

	} else if v, ok2 := d.GetOk("language_id"); ok2 {

		languageResp, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.LanguagesApi.ReadOneLanguage(ctx, d.Get("environment_id").(string), v.(string)).Execute()
			},
			"ReadOneLanguage",
			sdk.DefaultCustomError,
			nil,
		)
		if diags.HasError() {
			return diags
		}

		resp = languageResp.(*management.Language)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither language_id or locale are set",
			Detail:   "Neither language_id or locale are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("language_id", resp.GetId())
	d.Set("locale", resp.GetLocale())
	d.Set("enabled", resp.GetEnabled())
	d.Set("default", resp.GetDefault())

	d.Set("name", resp.GetName())

	if v, ok := resp.GetCustomerAddedOk(); ok {
		d.Set("customer_added", v)
	} else {
		d.Set("customer_added", nil)
	}

	return diags
}

func findLanguageByLocale(ctx context.Context, apiClient *management.APIClient, environmentID, locale string) (*management.Language, diag.Diagnostics) {
	var diags diag.Diagnostics

	respList, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.ReadLanguages(ctx, environmentID).Execute()
		},
		"ReadAllLanguages",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return nil, diags
	}

	respObject := respList.(*management.EntityArray)

	var resp *management.Language
	if languages, ok := respObject.Embedded.GetLanguagesOk(); ok {

		found := false
		for _, language := range languages {

			if language.Language.GetLocale() == locale {
				resp = language.Language
				found = true
				break
			}
		}

		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot find language by locale %s", locale),
			})

			return nil, diags
		}

	}

	return resp, diags
}
