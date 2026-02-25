// Copyright © 2026 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceTrustedEmailDomainSPF() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve Trusted Email Domain SPF status.",

		ReadContext: datasourcePingOneTrustedEmailDomainSPFRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"trusted_email_domain_id": {
				Description:      "A string that specifies the auto-generated ID of the email domain.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"type": {
				Description: "A string that specifies the type of DNS record.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: fmt.Sprintf("The status of the email domain ownership. Possible values are %s and %s", string(management.ENUMTRUSTEDEMAILSTATUS_ACTIVE), string(management.ENUMTRUSTEDEMAILSTATUS_VERIFICATION_REQUIRED)),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description: "Record name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "Record value.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneTrustedEmailDomainSPFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.TrustedEmailDomainsApi.ReadTrustedEmailDomainSPFStatus(ctx, d.Get("environment_id").(string), d.Get("trusted_email_domain_id").(string)).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"ReadTrustedEmailDomainSPFStatus",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.EmailDomainSPFStatus)

	d.SetId(d.Get("trusted_email_domain_id").(string))

	d.Set("type", respObject.GetType())
	d.Set("status", respObject.GetStatus())
	d.Set("key", respObject.GetKey())
	d.Set("value", respObject.GetValue())

	return diags
}
