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

func ResourceCertificateSigningResponse() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to import a CA issued response to a downloaded certificate signing request (CSR).",

		CreateContext: resourceCertificateSigningResponseCreate,
		ReadContext:   resourceCertificateSigningResponseRead,
		DeleteContext: resourceCertificateSigningResponseDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment that contains the key to which the CSR corresponds.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"key_id": {
				Description:      "The system name of the key.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"pem_ca_response_file": {
				Description: "A PEM encoded file that has been provided by the signing authority in response to the key's CSR.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The system name of the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Specifies the key algorithm. Options are `%s`, `%s`, and `%s`.", string(management.ENUMCERTIFICATEKEYALGORITHM_RSA), string(management.ENUMCERTIFICATEKEYALGORITHM_EC), string(management.ENUMCERTIFICATEKEYALGORITHM_UNKNOWN)),
				Computed:    true,
			},
			"default": {
				Type:        schema.TypeBool,
				Description: "A boolean that specifies whether this is the default key for the specified environment.",
				Computed:    true,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Description: "The time the key resource expires.",
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
				Description: fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
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
				Description: "An integer that specifies the number of days the key is valid.",
				Computed:    true,
			},
		},
	}
}

func resourceCertificateSigningResponseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	archive := []byte(d.Get("pem_ca_response_file").(string))

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.ImportCSRResponse(ctx, d.Get("environment_id").(string), d.Get("key_id").(string)).File(&archive).Execute()
		},
		"ImportCSRResponse",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Certificate)

	d.SetId(respObject.GetId())

	return resourceCertificateSigningResponseRead(ctx, d, meta)
}

func resourceCertificateSigningResponseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.GetKey(ctx, d.Get("environment_id").(string), d.Id()).Accept(management.ENUMGETKEYACCEPTHEADER_JSON).Execute()
		},
		"GetKey",
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

	respObject := resp.(*management.Certificate)

	serialNumber := respObject.GetSerialNumber()

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

func resourceCertificateSigningResponseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
