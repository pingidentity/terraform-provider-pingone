package base

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceCertificate() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne certificates.",

		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		DeleteContext: resourceCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCertificateImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the certificate in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"pkcs7_file_base64": {
				Description:  "A base64 encoded PKCS7 (DER) file to import.  Either `pkcs7_file_base64` or `pem_file` must be specified.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"pkcs7_file_base64", "pem_file"},
			},
			"pem_file": {
				Description:  "The contents of a PEM encoded file to import, which should be in plain text format and not base64 encoded.  The certificate should be properly formatted for the PEM format, that includes the correct header/footer lines.  Either `pkcs7_file_base64` or `pem_file` must be specified.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"pkcs7_file_base64", "pem_file"},
			},
			"usage_type": {
				Type:             schema.TypeString,
				Description:      fmt.Sprintf("A string that specifies how the certificate is used. Options are `%s`, `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYUSAGETYPE_ENCRYPTION), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SIGNING), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SSL_TLS), string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE)),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMCERTIFICATEKEYUSAGETYPE_ENCRYPTION), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SIGNING), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SSL_TLS), string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE)}, false)),
			},
			"name": {
				Description: "The system name of the certificate.",
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
			"validity_period": {
				Type:        schema.TypeInt,
				Description: "An integer that specifies the number of days the certificate is valid.",
				Computed:    true,
			},
		},
	}
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp interface{}

	var archive []byte

	if v, ok := d.GetOk("pkcs7_file_base64"); ok {

		var err error
		archive, err = base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot base64 decode provided PKCS7 certificate file.",
			})

			return diags
		}

	}

	if v, ok := d.GetOk("pem_file"); ok {
		archive = []byte(v.(string))
	}

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.CertificateManagementApi.CreateCertificateFromFile(ctx, d.Get("environment_id").(string)).ContentType("multipart/form-data").UsageType(d.Get("usage_type").(string)).File(&archive).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"CreateCertificateFromFile",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Certificate)

	d.SetId(respObject.GetId())

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.CertificateManagementApi.GetCertificate(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"GetCertificate",
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

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := apiClient.CertificateManagementApi.DeleteCertificate(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), nil, fR, fErr)
		},
		"DeleteCertificate",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceCertificateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "certificate_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.SetId(attributes["certificate_id"])

	resourceCertificateRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
