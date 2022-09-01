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

func ResourceGateway() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne gateways.",

		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceGatewayRead,
		UpdateContext: resourceGatewayUpdate,
		DeleteContext: resourceGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGatewayImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the gateway in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the gateway resource.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the gateway resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  fmt.Sprintf("The type of gateway resource. Options are `%s` and `%s`.", string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION)),
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION)}, false),
			},
			"enabled": {
				Description: "Indicates whether the gateway is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	gateway := *management.NewGateway(d.Get("name").(string), management.EnumGatewayType(d.Get("type").(string)), d.Get("enabled").(bool)) // Gateway |  (optional)

	if v, ok := d.GetOk("description"); ok {
		gateway.SetDescription(v.(string))
	}

	var gatewayRequest management.CreateGatewayRequest
	gatewayRequest.Gateway = &gateway

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.CreateGateway(ctx, d.Get("environment_id").(string)).CreateGatewayRequest(gatewayRequest).Execute()
		},
		"CreateGateway",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CreateGateway201Response)

	d.SetId(respObject.Gateway.GetId())

	return resourceGatewayRead(ctx, d, meta)
}

func resourceGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.ReadOneGateway(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneGateway",
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

	respObject := resp.(*management.CreateGateway201Response).Gateway

	d.Set("name", respObject.GetName())
	d.Set("type", respObject.GetType())
	d.Set("enabled", respObject.GetEnabled())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	return diags
}

func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	gateway := *management.NewGateway(d.Get("name").(string), management.EnumGatewayType(d.Get("type").(string)), d.Get("enabled").(bool)) // Gateway |  (optional)

	if v, ok := d.GetOk("description"); ok {
		gateway.SetDescription(v.(string))
	}

	var gatewayRequest management.CreateGatewayRequest
	gatewayRequest.Gateway = &gateway

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.UpdateGateway(ctx, d.Get("environment_id").(string), d.Id()).CreateGatewayRequest(gatewayRequest).Execute()
		},
		"UpdateGateway",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceGatewayRead(ctx, d, meta)
}

func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.GatewaysApi.DeleteGateway(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteGateway",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/gatewayID\"", d.Id())
	}

	environmentID, gatewayID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(gatewayID)

	resourceGatewayRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
