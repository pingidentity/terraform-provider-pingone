package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceCustomDomainSSL() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Custom Domain SSL settings.",

		CreateContext: resourceCustomDomainSSLCreate,
		ReadContext:   resourceCustomDomainSSLRead,
		DeleteContext: resourceCustomDomainSSLDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the certificate in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"custom_domain_id": {
				Description:      "The ID of the custom domain to set SSL settings for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"certificate_pem_file": {
				Description:      "A string that specifies the PEM-encoded certificate to import. The certificate must not be expired, must not be self signed and the domain must match one of the subject alternative name (SAN) values on the certificate.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"intermediate_certificates_pem_file": {
				Description:      "A string that specifies the PEM-encoded certificate chain.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"private_key_pem_file": {
				Description:      "A string that specifies the PEM-encoded, unencrypted private key that matches the certificate's public key.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
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
			"certificate_expires_at": {
				Type:        schema.TypeString,
				Description: "The time when the certificate expires.  If this property is not present, it indicates that an SSL certificate has not been setup for this custom domain.",
				Computed:    true,
			},
		},
	}
}

func resourceCustomDomainSSLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	customDomainCertificateRequest := *management.NewCustomDomainCertificateRequest(d.Get("certificate_pem_file").(string), d.Get("private_key_pem_file").(string))

	if v, ok := d.GetOk("intermediate_certificates_pem_file"); ok {
		customDomainCertificateRequest.SetIntermediateCertificates(v.(string))
	}

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CustomDomainsApi.UpdateDomain(ctx, d.Get("environment_id").(string), d.Get("custom_domain_id").(string)).ContentType(management.ENUMCUSTOMDOMAINPOSTHEADER_CERTIFICATE_IMPORTJSON).CustomDomainCertificateRequest(customDomainCertificateRequest).Execute()
		},
		"UpdateDomain",
		func(error management.P1Error) diag.Diagnostics {

			// Cannot validate against the authoritative name service
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				m, _ := regexp.MatchString("^Custom domain status must be 'SSL_CERTIFICATE_REQUIRED' or 'ACTIVE' in order to import a certificate", details[0].GetMessage())
				if m {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Cannot add SSL certificate settings to the custom domain - %s", details[0].GetMessage()),
						Detail:   `Please verify the domain first (hint: use the "pingone_custom_domain_verify" resource).)`,
					})
					return diags
				}
			}

			return nil
		},
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CustomDomain)

	d.SetId(respObject.GetId())

	return resourceCustomDomainSSLRead(ctx, d, meta)
}

func resourceCustomDomainSSLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if v, ok := respObject.GetCertificateOk(); ok {
		d.Set("certificate_expires_at", v.GetExpiresAt().Format(time.RFC3339))
	} else {
		d.Set("certificate_expires_at", nil)
	}

	return diags
}

func resourceCustomDomainSSLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
