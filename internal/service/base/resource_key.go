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

func ResourceKey() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne keys",

		CreateContext: resourceKeyCreate,
		ReadContext:   resourceKeyRead,
		UpdateContext: resourceKeyUpdate,
		DeleteContext: resourceKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceKeyImport,
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
				Description:      "The system name of the key (or Common Name).",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"algorithm": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("Specifies the key algorithm. Options are `%s`, `%s`, and `%s`.", string(management.ENUMCERTIFICATEKEYALGORITHM_RSA), string(management.ENUMCERTIFICATEKEYALGORITHM_EC), string(management.ENUMCERTIFICATEKEYALGORITHM_UNKNOWN)),
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(management.ENUMCERTIFICATEKEYALGORITHM_RSA), string(management.ENUMCERTIFICATEKEYALGORITHM_EC), string(management.ENUMCERTIFICATEKEYALGORITHM_UNKNOWN)}, false),
			},
			"default": {
				Type:        schema.TypeBool,
				Description: "A boolean that specifies whether this is the default key for the specified environment.",
				Optional:    true,
				Default:     false,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Description: "The time the key resource expires.",
				Computed:    true,
			},
			"issuer_dn": {
				Type:        schema.TypeString,
				Description: "A string that specifies the distinguished name of the certificate issuer.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"key_length": {
				Type:         schema.TypeInt,
				Description:  "An integer that specifies the key length. For RSA keys, options are `2048`, `3072`, `4096` and `7680`. For elliptical curve (EC) keys, options are `224`, `256`, `384` and `521`.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{2048, 3072, 4096, 7680, 224, 256, 384, 521}),
			},
			"serial_number": {
				Type:        schema.TypeInt,
				Description: "An integer that specifies the serial number of the key or certificate.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"signature_algorithm": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s`, `%s` and `%s`", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA224WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)}, false),
			},
			"starts_at": {
				Type:        schema.TypeString,
				Description: "The time the validity period starts.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A string that specifies the status of the key. Options are `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMCERTIFICATEKEYSTATUS_VALID), string(management.ENUMCERTIFICATEKEYSTATUS_EXPIRED), string(management.ENUMCERTIFICATEKEYSTATUS_NOT_YET_VALID), string(management.ENUMCERTIFICATEKEYSTATUS_REVOKED)),
				Computed:    true,
			},
			"subject_dn": {
				Type:             schema.TypeString,
				Description:      "A string that specifies the distinguished name of the subject being secured.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"usage_type": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("A string that specifies how the certificate is used. Options are `%s`, %s, %s and `%s`.", string(management.ENUMCERTIFICATEKEYUSAGETYPE_ENCRYPTION), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SIGNING), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SSL_TLS), string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE)),
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{string(management.ENUMCERTIFICATEKEYUSAGETYPE_ENCRYPTION), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SIGNING), string(management.ENUMCERTIFICATEKEYUSAGETYPE_SSL_TLS), string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE)}, false),
			},
			"validity_period": {
				Type:         schema.TypeInt,
				Description:  "An integer that specifies the number of days the key is valid.",
				Optional:     true,
				Default:      365,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func resourceKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	certificateKey := *management.NewCertificate(
		management.EnumCertificateKeyAlgorithm(d.Get("algorithm").(string)),
		int32(d.Get("key_length").(int)),
		d.Get("name").(string),
		management.EnumCertificateKeySignagureAlgorithm(d.Get("signature_algorithm").(string)),
		d.Get("subject_dn").(string),
		management.EnumCertificateKeyUsageType(d.Get("usage_type").(string)),
		int32(d.Get("validity_period").(int)),
	)

	if v, ok := d.GetOk("default"); ok {
		certificateKey.SetDefault(v.(bool))
	}

	if v, ok := d.GetOk("issuer_dn"); ok {
		certificateKey.SetIssuerDN(v.(string))
	}

	if v, ok := d.GetOk("serial_number"); ok {
		certificateKey.SetSerialNumber(int32(v.(int)))
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.CreateKey(ctx, d.Get("environment_id").(string)).Certificate(certificateKey).Execute()
		},
		"CreateKey",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Certificate)

	d.SetId(respObject.GetId())

	return resourceKeyRead(ctx, d, meta)
}

func resourceKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.GetKey(ctx, d.Get("environment_id").(string), d.Id()).Execute()
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

	d.Set("name", respObject.GetName())
	d.Set("algorithm", string(respObject.GetAlgorithm()))
	d.Set("default", respObject.GetDefault())
	d.Set("expires_at", respObject.GetExpiresAt())
	d.Set("issuer_dn", respObject.GetIssuerDN())
	d.Set("key_length", respObject.GetKeyLength())
	d.Set("serial_number", respObject.GetSerialNumber())
	d.Set("signature_algorithm", string(respObject.GetSignatureAlgorithm()))
	d.Set("starts_at", respObject.GetStartsAt())
	d.Set("status", string(respObject.GetStatus()))
	d.Set("subject_dn", respObject.GetSubjectDN())
	d.Set("usage_type", string(respObject.GetUsageType()))
	d.Set("validity_period", string(respObject.GetValidityPeriod()))

	return diags
}

func resourceKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	keyUpdate := *management.NewCertificateKeyUpdate(d.Get("default").(bool), management.EnumCertificateKeyUsageType(d.Get("usage_type").(string)))

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.UpdateKey(ctx, d.Get("environment_id").(string), d.Id()).CertificateKeyUpdate(keyUpdate).Execute()
		},
		"UpdateKey",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceKeyRead(ctx, d, meta)
}

func resourceKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.CertificateManagementApi.DeleteKey(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteKey",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/keyID\"", d.Id())
	}

	environmentID, keyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(keyID)

	resourceKeyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
