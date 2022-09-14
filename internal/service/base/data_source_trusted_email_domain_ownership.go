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

func DatasourceTrustedEmailDomainOwnership() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve Trusted Email Domain ownership status.",

		ReadContext: datasourcePingOneTrustedEmailDomainOwnershipRead,

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
			"region": {
				Description: "The regions collection specifies the properties for the 4 AWS SES regions that are used for sending email for the environment. The regions are determined by the geography where this environment was provisioned (North America, Canada, Europe & Asia-Pacific).",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the region.",
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
				},
			},
		},
	}
}

func datasourcePingOneTrustedEmailDomainOwnershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.TrustedEmailDomainsApi.ReadTrustedEmailDomainOwnershipStatus(ctx, d.Get("environment_id").(string), d.Get("trusted_email_domain_id").(string)).Execute()
		},
		"ReadTrustedEmailDomainOwnershipStatus",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.EmailDomainOwnershipStatus)

	d.SetId(d.Get("trusted_email_domain_id").(string))

	d.Set("type", respObject.GetType())
	d.Set("region", flattenOwnershipRegions(respObject.GetRegions()))

	return diags
}

func flattenOwnershipRegions(c []management.EmailDomainOwnershipStatusRegionsInner) []interface{} {

	items := make([]interface{}, 0)

	for _, v := range c {
		// Required
		items = append(items, map[string]interface{}{
			"name":   v.GetName(),
			"status": string(v.GetStatus()),
			"key":    v.GetKey(),
			"value":  v.GetValue(),
		})

	}

	return items

}
