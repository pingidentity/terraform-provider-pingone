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
	"golang.org/x/exp/slices"
)

func ResourceLanguage() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne languages.  To fully enable a created language, the `pingone_language_update` resource must be used to complete the configuration.",

		CreateContext: resourceLanguageCreate,
		ReadContext:   resourceLanguageRead,
		DeleteContext: resourceLanguageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLanguageImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the language in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"locale": {
				Description:      fmt.Sprintf("An ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).  The following language codes are reserved as they are created automatically in the environment: %s.", verify.IsoReservedListString()),
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(verify.IsoList(), false)),
			},
			"name": {
				Description: "The language name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Specifies whether this language is enabled for the environment. This property value must be set to false when creating a language.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"default": {
				Description: "Specifies whether this language is the default for the environment. This property value must be set to `false` when creating a language resource. It can be set to `true` only after the language is enabled and after the localization of an agreement resource is complete when agreements are used for the environment.",
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

func resourceLanguageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp interface{}

	language := *management.NewLanguage(false, false, d.Get("locale").(string))

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
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

	return resourceLanguageRead(ctx, d, meta)
}

func resourceLanguageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
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

	if slices.Contains(verify.ReservedIsoList(), respObject.GetLocale()) {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("The language code `%s` is reserved and cannot be imported into this provider.  Please use `pingone_language_override` for system-defined languages instead.", respObject.GetLocale()),
		})

		return diags
	}

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

func resourceLanguageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := apiClient.LanguagesApi.DeleteLanguage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteLanguage",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
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
