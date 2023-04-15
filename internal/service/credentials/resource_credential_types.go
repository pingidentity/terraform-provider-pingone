package credentials

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mjspi/pingone-neo-go-sdk/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceCredentialTypes() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne Environment's Credential Types",

		ReadContext: resourceCredentialTypesReadAll,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCredentialTypesImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment ",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
		},
	}
}

func resourceCredentialTypesReadAll(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PingOneCredentialsCredentialTypesApi.ReadAllCredentialTypes(ctx, d.Get("environment_id").(string)).Execute()
		},
		"ReadAllCredentialTypes",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*mfa.MFASettings)

	d.Set("pairing", flattenMFASettingPairing(respObject.GetPairing()))

	if v, ok := respObject.GetLockoutOk(); ok {
		d.Set("lockout", flattenMFASettingLockout(*v))
	} else {
		d.Set("lockout", nil)
	}

	if v, ok := respObject.GetAuthenticationOk(); ok {
		d.Set("authentication", flattenMFASettingAuthentication(*v))
	} else {
		d.Set("authentication", nil)
	}

	return diags
}
