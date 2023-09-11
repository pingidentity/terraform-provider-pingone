package base

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceCertificate() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read metadata for certificates stored in PingOne.",

		ReadContext: datasourcePingOneCertificateRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"certificate_id": {
				Description:      "The ID of the certificate.  Either `certificate_id` or `name` must be specified.",
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"certificate_id", "name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:  "The system name of the certificate.  Either `certificate_id` or `name` must be specified.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"certificate_id", "name"},
			},
			"algorithm": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Specifies the key algorithm. Options are `%s`, `%s`, and `%s`.", string(management.ENUMCERTIFICATEKEYALGORITHM_RSA), string(management.ENUMCERTIFICATEKEYALGORITHM_EC), string(management.ENUMCERTIFICATEKEYALGORITHM_UNKNOWN)),
				Computed:    true,
			},
			"default": {
				Type:        schema.TypeBool,
				Description: "A boolean that specifies whether this is the default certificate for the specified environment.",
				Computed:    true,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Description: "The time the certificate expires.",
				Computed:    true,
			},
			"issuer_dn": {
				Type:        schema.TypeString,
				Description: "A string that specifies the distinguished name of the certificate issuer.",
				Computed:    true,
			},
			"key_length": {
				Type:        schema.TypeInt,
				Description: "An integer that specifies the key length. For RSA keys, options are `2048`, `3072`, `4096` and `7680`. For elliptical curve (EC) keys, options are `224`, `256`, `384` and `521`.",
				Computed:    true,
			},
			"serial_number": {
				Type:        schema.TypeString,
				Description: "An integer (in string data type) that specifies the serial number of the key or certificate.",
				Computed:    true,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
				Computed:    true,
			},
			"starts_at": {
				Type:        schema.TypeString,
				Description: "The time the validity period starts.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A string that specifies the status of the key. Options are `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMCERTIFICATEKEYSTATUS_VALID), string(management.ENUMCERTIFICATEKEYSTATUS_EXPIRING), string(management.ENUMCERTIFICATEKEYSTATUS_EXPIRED), string(management.ENUMCERTIFICATEKEYSTATUS_NOT_YET_VALID), string(management.ENUMCERTIFICATEKEYSTATUS_REVOKED)),
				Computed:    true,
			},
			"subject_dn": {
				Type:        schema.TypeString,
				Description: "A string that specifies the distinguished name of the subject being secured.",
				Computed:    true,
			},
			"usage_type": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A string that specifies how the certificate is used. Options are `%s`, `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYUSAGETYPE_ENCRYPTION), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SIGNING), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SSL_TLS), string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE)),
				Computed:    true,
			},
			"validity_period": {
				Type:        schema.TypeInt,
				Description: "An integer that specifies the number of days the certificate is valid.",
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var respObject management.Certificate

	if v, ok := d.GetOk("name"); ok {

		respList, diags := sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.CertificateManagementApi.GetCertificates(ctx, d.Get("environment_id").(string)).Execute()
			},
			"GetCertificates",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		respObjectList := respList.(*management.EntityArray)

		if certificates, ok := respObjectList.Embedded.GetCertificatesOk(); ok {

			found := false
			for _, certificate := range certificates {

				if certificate.GetName() == v.(string) {
					respObject = certificate
					found = true
					break
				}
			}

			if !found {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Cannot find certificate %s", v),
				})

				return diags
			}

		}

	} else if v, ok2 := d.GetOk("certificate_id"); ok2 {

		resp, diags := sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.CertificateManagementApi.GetCertificate(ctx, d.Get("environment_id").(string), v.(string)).Execute()
			},
			"GetCertificate",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		respObject = *resp.(*management.Certificate)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither certificate_id or name are set",
			Detail:   "Neither certificate_id or name are set",
		})

		return diags

	}

	serialNumber := respObject.GetSerialNumber()

	d.SetId(respObject.GetId())

	d.Set("name", respObject.GetName())
	d.Set("algorithm", string(respObject.GetAlgorithm()))
	d.Set("default", respObject.GetDefault())
	d.Set("expires_at", respObject.GetExpiresAt().Format(time.RFC3339))
	d.Set("issuer_dn", respObject.GetIssuerDN())
	d.Set("key_length", respObject.GetKeyLength())
	d.Set("serial_number", serialNumber.String())
	d.Set("signature_algorithm", string(respObject.GetSignatureAlgorithm()))
	d.Set("starts_at", respObject.GetStartsAt().Format(time.RFC3339))
	d.Set("status", string(respObject.GetStatus()))
	d.Set("subject_dn", respObject.GetSubjectDN())
	d.Set("usage_type", string(respObject.GetUsageType()))
	d.Set("validity_period", respObject.GetValidityPeriod())

	return diags
}
