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

func ResourceLanguage() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne languages.",

		CreateContext: resourceLanguageCreate,
		ReadContext:   resourceLanguageRead,
		UpdateContext: resourceLanguageUpdate,
		DeleteContext: resourceLanguageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLanguageImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the key in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "The user-defined language name.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"locale": {
				Description:      "An ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(verify.IsoList(), false)),
			},
			"enabled": {
				Description: "Specifies whether this language is enabled for the environment. This property value must be set to false when creating a language.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"default": {
				Description: "Specifies whether this language is the default for the environment. This property value must be set to `false` when creating a language resource. It can be set to `true` only after the language is enabled and after the localization of an agreement resource is complete when agreements are used for the environment.",
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

func resourceLanguageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	language := *management.NewLanguage(d.Get("default").(bool), false, d.Get("locale").(string))
	language.SetName(d.Get("name").(string))

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.CreateLanguage(ctx, d.Get("environment_id").(string)).Language(language).Execute()
		},
		"CreateLanguage",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Language)

	d.SetId(respObject.GetId())

	if d.Get("enabled").(bool) {

		language.SetEnabled(true)

		_, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.LanguagesApi.UpdateLanguage(ctx, d.Get("environment_id").(string), respObject.GetId()).Language(language).Execute()
			},
			"UpdateLanguage",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

	}

	return resourceLanguageRead(ctx, d, meta)
}

func resourceLanguageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceLanguageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	language := *management.NewLanguage(d.Get("default").(bool), d.Get("enabled").(bool), d.Get("locale").(string))
	language.SetName(d.Get("name").(string))

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LanguagesApi.UpdateLanguage(ctx, d.Get("environment_id").(string), d.Id()).Language(language).Execute()
		},
		"UpdateLanguage",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceLanguageRead(ctx, d, meta)
}

func resourceLanguageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.LanguagesApi.DeleteLanguage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteLanguage",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceLanguageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/languageID\"", d.Id())
	}

	environmentID, languageID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(languageID)

	resourceLanguageRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
