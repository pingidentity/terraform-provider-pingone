// Copyright © 2026 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceCertificateSigningRequest() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to export a certificate signing request (CSR) from a PingOne Key.",

		ReadContext: datasourcePingOneCertificateSigningRequestRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"key_id": {
				Description:      "The ID of the key to export the CSR from.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"pkcs10_file_base64": {
				Description: "The Certificate Signing Request (CSR) in PKCS10 file format, base64 encoded.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pem_file": {
				Description: "The Certificate Signing Request (CSR) in PEM file format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneCertificateSigningRequestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := certificateSigningExport(ctx, apiClient, d.Get("environment_id").(string), d.Get("key_id").(string), management.ENUMCSREXPORTHEADER_PKCS10)
	if diags.HasError() {
		return diags
	}

	d.Set("pkcs10_file_base64", base64.StdEncoding.EncodeToString([]byte(resp.(string))))

	respPem, diags := certificateSigningExport(ctx, apiClient, d.Get("environment_id").(string), d.Get("key_id").(string), management.ENUMCSREXPORTHEADER_X_PEM_FILE)
	if diags.HasError() {
		return diags
	}

	d.Set("pem_file", strings.TrimSuffix(respPem.(string), "\n"))
	d.SetId(d.Get("key_id").(string))

	return diags
}

func certificateSigningExport(ctx context.Context, apiClient *management.APIClient, environmentID, keyID string, exportFileType management.EnumCSRExportHeader) (interface{}, diag.Diagnostics) {
	return sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.CertificateManagementApi.ExportCSR(ctx, environmentID, keyID).Accept(exportFileType).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"ExportCSR",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
}
