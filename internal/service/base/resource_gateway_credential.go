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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewayCredentialsApi.CreateGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string)).Execute()
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewayCredentialsApi.ReadOneGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string), d.Id()).Execute()
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

	d.Set("created_at", respObject.GetCreatedAt())
	d.Set("last_used_at", respObject.GetLastUsedAt())

	return diags
}

func resourceGatewayCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.GatewayCredentialsApi.DeleteGatewayCredential(ctx, d.Get("environment_id").(string), d.Get("gateway_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteGatewayCredential",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGatewayCredentialImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 3 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/gatewayID/gatewayCredentialID\"", d.Id())
	}

	environmentID, gatewayID, gatewayCredentialID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("gateway_id", gatewayID)
	d.SetId(gatewayCredentialID)

	resourceGatewayCredentialRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
