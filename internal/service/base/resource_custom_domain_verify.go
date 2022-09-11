package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceCustomDomainVerify() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Custom Domain verification.",

		CreateContext: resourceCustomDomainVerifyCreate,
		ReadContext:   resourceCustomDomainVerifyRead,
		DeleteContext: resourceCustomDomainVerifyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the certificate in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"custom_domain_id": {
				Description:      "The ID of the custom domain to verify.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"domain_name": {
				Type:        schema.TypeString,
				Description: "A string that specifies the domain name in use.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A string that specifies the status of the custom domain. Options are `%s`, `%s` and `%s`.", string(management.ENUMCUSTOMDOMAINSTATUS_ACTIVE), string(management.ENUMCUSTOMDOMAINSTATUS_VERIFICATION_REQUIRED), string(management.ENUMCUSTOMDOMAINSTATUS_SSL_CERTIFICATE_REQUIRED)),
				Computed:    true,
			},
		},
	}
}

func resourceCustomDomainVerifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	timeoutValue := 60

	resp, diags = sdk.ParseResponseWithCustomTimeout(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CustomDomainsApi.UpdateDomain(ctx, d.Get("environment_id").(string), d.Get("custom_domain_id").(string)).ContentType(management.ENUMCUSTOMDOMAINPOSTHEADER_DOMAIN_NAME_VERIFYJSON).Execute()
		},
		"UpdateDomain",
		func(error management.P1Error) diag.Diagnostics {

			// Cannot validate against the authoritative name service
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				m, _ := regexp.MatchString("^Error response from authoritative name servers: NXDOMAIN", details[0].GetMessage())
				if m {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Cannot verify the domain - %s", details[0].GetMessage()),
						Detail:   `Please check the domain authority exists or is reachable.`,
					})
					return diags
				}

				m, _ = regexp.MatchString("^No CNAME records found", details[0].GetMessage())
				if m {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Cannot verify the domain - %s", details[0].GetMessage()),
						Detail:   `Please check the domain authority has the correct CNAME value set (hint: if using the "pingone_custom_domain" resource, the CNAME value to use is returned in the "canonical_name" attribute.)`,
					})
					return diags
				}
			}

			return nil
		},
		customDomainRetryConditions,
		time.Duration(timeoutValue)*time.Minute, // 60 mins
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CustomDomain)

	d.SetId(respObject.GetId())

	return resourceCustomDomainVerifyRead(ctx, d, meta)
}

func resourceCustomDomainVerifyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CustomDomainsApi.ReadOneDomain(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneDomain",
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

	respObject := resp.(*management.CustomDomain)

	d.Set("domain_name", respObject.GetDomainName())
	d.Set("status", string(respObject.GetStatus()))

	return diags
}

func resourceCustomDomainVerifyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func customDomainRetryConditions(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

	if p1error != nil {
		var err error

		// Permissions may not have propagated by this point
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

		// add retry time for DNS propegating
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			// perhaps it's the DNS authority
			if m, err := regexp.MatchString("^Error response from authoritative name servers: NXDOMAIN", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("Cannot verify the domain - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

			// perhaps it's the CNAME
			if m, err := regexp.MatchString("^No CNAME records found", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("Cannot verify the domain - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}
		}

	}

	return false
}
