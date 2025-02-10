// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
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

func ResourceGatewayCredential() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne gateway credentials.",

		CreateContext: resourceGatewayCredentialCreate,
		ReadContext:   resourceGatewayCredentialRead,
		DeleteContext: resourceGatewayCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGatewayCredentialImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the gateway credential in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"gateway_id": {
				Description:      "The ID of the gateway to associate the credential with.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"created_at": {
				Description: "A date that specifies the date the credential was created in Coordinated Universal Time (UTC).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_used_at": {
				Description: "A date that specifies the date the credential was last used in UTC.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"credential": {
				Description: "A string that specifies the signed JWT for the gateway credential.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceGatewayCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.GatewayCredentialsApi.CreateGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string)).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"CreateGatewayCredential",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.GatewayCredential)

	d.SetId(respObject.GetId())
	d.Set("credential", respObject.GetCredential())

	return resourceGatewayCredentialRead(ctx, d, meta)
}

func resourceGatewayCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.GatewayCredentialsApi.ReadOneGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"ReadOneGatewayCredential",
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

	respObject := resp.(*management.GatewayCredential)

	if v, ok := respObject.GetCreatedAtOk(); ok {
		d.Set("created_at", v.Format(time.RFC3339))
	} else {
		d.Set("created_at", nil)
	}

	if v, ok := respObject.GetLastUsedAtOk(); ok {
		d.Set("last_used_at", v.Format(time.RFC3339))
	} else {
		d.Set("last_used_at", nil)
	}

	return diags
}

func resourceGatewayCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := apiClient.GatewayCredentialsApi.DeleteGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), nil, fR, fErr)
		},
		"DeleteGatewayCredential",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGatewayCredentialImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "gateway_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "gateway_credential_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.Set("gateway_id", attributes["gateway_id"])
	d.SetId(attributes["gateway_credential_id"])

	resourceGatewayCredentialRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
