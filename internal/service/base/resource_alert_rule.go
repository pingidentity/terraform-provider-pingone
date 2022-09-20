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

func ResourceAlertRule() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Alert notification rules.",

		CreateContext: resourceAlertRuleCreate,
		ReadContext:   resourceAlertRuleRead,
		UpdateContext: resourceAlertRuleUpdate,
		DeleteContext: resourceAlertRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAlertRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the alert rule in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"channel_type": {
				Description:      "The alert channel type.",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(management.ENUMALERTCHANNELTYPE_EMAIL),
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMALERTCHANNELTYPE_EMAIL)}, false)),
			},
			"addresses": {
				Description: "The email addresses to send the alerts to.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
			},
			"include_severities": {
				Description: fmt.Sprintf("Filters alerts by severity. If empty, all severities are included. Possible values are `%s`, `%s`, and `%s`.", string(management.ENUMALERTCHANNELSEVERITY_INFO), string(management.ENUMALERTCHANNELSEVERITY_WARNING), string(management.ENUMALERTCHANNELSEVERITY_ERROR)),
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMALERTCHANNELSEVERITY_INFO), string(management.ENUMALERTCHANNELSEVERITY_WARNING), string(management.ENUMALERTCHANNELSEVERITY_ERROR)}, false)),
				},
			},
			"include_alert_types": {
				Description: fmt.Sprintf("Filters alerts by alert type. If empty, all alert types are included. Possible values are `%s`, `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATED), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATING)),
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATED), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATING)}, false)),
				},
			},
			"exclude_alert_types": {
				Description: fmt.Sprintf("Administrators will not be emailed alerts of these types. If empty, no alert types are excluded. Possible values are `%s`, `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATED), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATING)),
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_CERTIFICATE_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRED), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATED), string(management.ENUMALERTCHANNELALERTTYPE_KEY_PAIR_EXPIRING), string(management.ENUMALERTCHANNELALERTTYPE_GATEWAY_VERSION_DEPRECATING)}, false)),
				},
			},
		},
	}
}

func resourceAlertRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	alertChannel := expandAlertChannelRule(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.AlertingApi.CreateAlertChannel(ctx, d.Get("environment_id").(string)).AlertChannel(*alertChannel).Execute()
		},
		"CreateAlertChannel",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.AlertChannel)

	d.SetId(respObject.GetId())

	return resourceAlertRuleRead(ctx, d, meta)
}

func resourceAlertRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.AlertingApi.ReadOneAlertChannel(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneAlertChannel",
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

	respObject := resp.(*management.AlertChannel)

	d.Set("channel_type", string(respObject.GetChannelType()))
	d.Set("addresses", respObject.GetAddresses())

	if v, ok := respObject.GetIncludeSeveritiesOk(); ok {
		d.Set("include_severities", v)
	} else {
		d.Set("include_severities", nil)
	}

	if v, ok := respObject.GetIncludeAlertTypesOk(); ok {
		d.Set("include_alert_types", v)
	} else {
		d.Set("include_alert_types", nil)
	}

	if v, ok := respObject.GetExcludeAlertTypesOk(); ok {
		d.Set("exclude_alert_types", v)
	} else {
		d.Set("exclude_alert_types", nil)
	}

	return diags
}

func resourceAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	alertChannel := expandAlertChannelRule(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.AlertingApi.UpdateAlertChannel(ctx, d.Get("environment_id").(string), d.Id()).AlertChannel(*alertChannel).Execute()
		},
		"UpdateAlertChannel",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceAlertRuleRead(ctx, d, meta)
}

func resourceAlertRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.AlertingApi.DeleteAlertChannel(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteAlertChannel",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceAlertRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/alertRuleID\"", d.Id())
	}

	environmentID, alertRuleID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(alertRuleID)

	resourceAlertRuleRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandAlertChannelRule(d *schema.ResourceData) *management.AlertChannel {

	addresses := make([]string, 0)
	for _, v := range d.Get("addresses").(*schema.Set).List() {
		addresses = append(addresses, v.(string))
	}

	returnVar := management.NewAlertChannel(management.EnumAlertChannelType(d.Get("channel_type").(string)), addresses)

	if v, ok := d.GetOk("include_severities"); ok {
		list := make([]management.EnumAlertChannelSeverity, 0)
		for _, v1 := range v.(*schema.Set).List() {
			list = append(list, management.EnumAlertChannelSeverity(v1.(string)))
		}
		if len(list) > 0 {
			returnVar.SetIncludeSeverities(list)
		}
	}

	if v, ok := d.GetOk("include_alert_types"); ok {
		list := make([]management.EnumAlertChannelAlertType, 0)
		for _, v1 := range v.(*schema.Set).List() {
			list = append(list, management.EnumAlertChannelAlertType(v1.(string)))
		}
		if len(list) > 0 {
			returnVar.SetIncludeAlertTypes(list)
		}
	}

	if v, ok := d.GetOk("exclude_alert_types"); ok {
		list := make([]management.EnumAlertChannelAlertType, 0)
		for _, v1 := range v.(*schema.Set).List() {
			list = append(list, management.EnumAlertChannelAlertType(v1.(string)))
		}
		if len(list) > 0 {
			returnVar.SetExcludeAlertTypes(list)
		}
	}

	return returnVar
}
