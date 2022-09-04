package base

import (
	"context"
	"encoding/base64"
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

func DatasourceCertificateExport() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to export the public certificate (in PEM and DER file encoding) from a Key pair stored in PingOne.",

		ReadContext: datasourcePingOneCertificateExportRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"key_id": {
				Description:      "The ID of the key to export the public certificate from.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"pkcs7_file_base64": {
				Description: "The public certificate in PKCS7 DER file format, base64 encoded.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pem_file": {
				Description: "The public certificate in X509 PEM file format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneCertificateExportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := certificateExport(ctx, apiClient, d.Get("environment_id").(string), d.Get("key_id").(string), management.ENUMGETKEYACCEPTHEADER_X_PKCS7_CERTIFICATES)
	if diags.HasError() {
		return diags
	}

	d.Set("pkcs7_file_base64", base64.StdEncoding.EncodeToString(resp.([]byte)))

	respPem, diags := certificateExport(ctx, apiClient, d.Get("environment_id").(string), d.Get("key_id").(string), management.ENUMGETKEYACCEPTHEADER_X_X509_CA_CERT)
	if diags.HasError() {
		return diags
	}

	d.Set("pem_file", strings.TrimSuffix(respPem.(string), "\n"))
	d.SetId(d.Get("key_id").(string))

	return diags
}

func certificateExport(ctx context.Context, apiClient *management.APIClient, environmentID, keyID string, exportFileType management.EnumGetKeyAcceptHeader) (interface{}, diag.Diagnostics) {
	return sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CertificateManagementApi.GetKey(ctx, environmentID, keyID).Accept(exportFileType).Execute()
		},
		"GetKey",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
}
