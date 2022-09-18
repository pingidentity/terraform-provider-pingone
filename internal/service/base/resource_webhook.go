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

func ResourceWebhook() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Webhooks / Data Subscriptions.",

		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWebhookImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the webhook in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "A string that specifies the webhook name.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "A boolean that specifies whether a created or updated webhook should be active or suspended. A suspended state (`\"enabled\":false`) accumulates all matched events, but these events are not delivered until the webhook becomes active again (`\"enabled\":true`). For suspended webhooks, events accumulate for a maximum of two weeks. Events older than two weeks are deleted. Restarted webhooks receive the saved events (up to two weeks from the restart date).",
				Optional:    true,
				Default:     false,
			},
			"http_endpoint_url": {
				Description:      "A string that specifies a valid HTTPS URL to which event messages are sent.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
			},
			"http_endpoint_headers": {
				Description: "A map that specifies the headers applied to the outbound request (for example, `Authorization` `Basic usernamepassword`. The purpose of these headers is for the HTTPS endpoint to authenticate the PingOne service, ensuring that the information from PingOne is from a trusted source.",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"verify_tls_certificates": {
				Type:        schema.TypeBool,
				Description: "A boolean that specifies whether a certificates should be verified. If this property's value is set to `false`, then all certificates are trusted. (Setting this property's value to false introduces a security risk.)",
				Optional:    true,
				Default:     true,
			},
			"format": {
				Description:      fmt.Sprintf("A string that specifies one of the supported webhook formats. Options are `%s`, `%s`, and `%s`.", string(management.ENUMSUBSCRIPTIONFORMAT_ACTIVITY), string(management.ENUMSUBSCRIPTIONFORMAT_SPLUNK), string(management.ENUMSUBSCRIPTIONFORMAT_NEWRELIC)),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMSUBSCRIPTIONFORMAT_ACTIVITY), string(management.ENUMSUBSCRIPTIONFORMAT_SPLUNK), string(management.ENUMSUBSCRIPTIONFORMAT_NEWRELIC)}, false)),
			},
			"filter_options": {
				Description: "A block that specifies the PingOne platform event filters to be included to trigger this webhook.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"included_action_types": {
							Description: "A non-empty list that specifies the list of action types that should be matched for the webhook.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"included_application_ids": {
							Description: "An array that specifies the list of applications (by ID) whose events are monitored by the webhook (maximum of 10 IDs in the array). If a list of applications is not provided, events are monitored for all applications in the environment.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"included_population_ids": {
							Description: "An array that specifies the list of populations (by ID) whose events are monitored by the webhook (maximum of 10 IDs in the array). This property matches events for users in the specified populations, as opposed to events generated in which the user in one of the populations is the actor.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"included_tags": {
							Description: fmt.Sprintf("An array of tags that events must have to be monitored by the webhook. If tags are not specified, all events are monitored. Currently, the available tags are `%s`. Identifies the event as the action of an administrator on other administrators.", string(management.ENUMSUBSCRIPTIONFILTERINCLUDEDTAGS_ADMIN_IDENTITY_EVENT)),
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMSUBSCRIPTIONFILTERINCLUDEDTAGS_ADMIN_IDENTITY_EVENT)}, false)),
							},
						},
					},
				},
			},
		},
	}
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	subscription, diags := expandWebhook(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SubscriptionsWebhooksApi.CreateSubscription(ctx, d.Get("environment_id").(string)).Subscription(*subscription).Execute()
		},
		"CreateSubscription",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Subscription)

	d.SetId(respObject.GetId())

	return resourceWebhookRead(ctx, d, meta)
}

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SubscriptionsWebhooksApi.ReadOneSubscription(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneSubscription",
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

	respObject := resp.(*management.Subscription)

	httpEndpoint := respObject.GetHttpEndpoint()

	d.Set("name", respObject.GetName())
	d.Set("enabled", respObject.GetEnabled())
	d.Set("http_endpoint_url", httpEndpoint.GetUrl())

	if v, ok := httpEndpoint.GetHeadersOk(); ok {
		d.Set("http_endpoint_headers", v)
	} else {
		d.Set("http_endpoint_headers", nil)
	}

	d.Set("verify_tls_certificates", respObject.GetVerifyTlsCertificates())
	d.Set("format", respObject.GetFormat())
	d.Set("filter_options", flattenWebhookFilterOptions(respObject.GetFilterOptions()))

	return diags
}

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	subscription, diags := expandWebhook(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SubscriptionsWebhooksApi.UpdateSubscription(ctx, d.Get("environment_id").(string), d.Id()).Subscription(*subscription).Execute()
		},
		"UpdateSubscription",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceWebhookRead(ctx, d, meta)
}

func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.SubscriptionsWebhooksApi.DeleteSubscription(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteSubscription",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceWebhookImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/webhookID\"", d.Id())
	}

	environmentID, webhookID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(webhookID)

	resourceWebhookRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandWebhook(d *schema.ResourceData) (*management.Subscription, diag.Diagnostics) {
	var diags diag.Diagnostics

	httpEndpoint := *management.NewSubscriptionHttpEndpoint(d.Get("http_endpoint_url").(string))

	if v, ok := d.GetOk("http_endpoint_headers"); ok {
		obj := v.(map[string]interface{})

		httpEndpoint.SetHeaders(obj)
	}

	filterOptions, diags := expandWebhookFilterOptions(d.Get("filter_options").([]interface{}))
	if diags.HasError() {
		return nil, diags
	}

	returnVar := management.NewSubscription(
		d.Get("enabled").(bool),
		*filterOptions,
		management.EnumSubscriptionFormat(d.Get("format").(string)),
		httpEndpoint,
		d.Get("name").(string),
		d.Get("verify_tls_certificates").(bool),
	)

	return returnVar, diags
}

func expandWebhookFilterOptions(c []interface{}) (*management.SubscriptionFilterOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	obj := c[0].(map[string]interface{})

	if obj == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot expand webhook filters - the filter options block is empty",
		})

		return nil, diags
	}

	includedActionTypes := make([]string, 0)
	for _, v := range obj["included_action_types"].(*schema.Set).List() {
		includedActionTypes = append(includedActionTypes, v.(string))
	}

	returnVar := management.NewSubscriptionFilterOptions(includedActionTypes)

	if v, ok := obj["included_application_ids"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != nil {
		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v1 := range v.List() {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v1.(string)))
		}
		returnVar.SetIncludedApplications(objList)
	}

	if v, ok := obj["included_population_ids"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != nil {
		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v1 := range v.List() {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v1.(string)))
		}
		returnVar.SetIncludedPopulations(objList)
	}

	if v, ok := obj["included_tags"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != nil {
		objList := make([]management.EnumSubscriptionFilterIncludedTags, 0)
		for _, v1 := range v.List() {
			objList = append(objList, management.EnumSubscriptionFilterIncludedTags(v1.(string)))
		}
		returnVar.SetIncludedTags(objList)
	}

	return returnVar, diags
}

func flattenWebhookFilterOptions(subscriptionFilterOptions management.SubscriptionFilterOptions) []interface{} {

	item := map[string]interface{}{
		"included_action_types": subscriptionFilterOptions.GetIncludedActionTypes(),
	}

	if v, ok := subscriptionFilterOptions.GetIncludedApplicationsOk(); ok {

		list := make([]string, 0)

		for _, v1 := range v {
			list = append(list, v1.GetId())
		}

		item["included_application_ids"] = list
	}

	if v, ok := subscriptionFilterOptions.GetIncludedPopulationsOk(); ok {

		list := make([]string, 0)

		for _, v1 := range v {
			list = append(list, v1.GetId())
		}

		item["included_population_ids"] = list
	}

	if v, ok := subscriptionFilterOptions.GetIncludedTagsOk(); ok {

		list := make([]string, 0)

		for _, v1 := range v {
			list = append(list, string(v1))
		}

		item["included_tags"] = list
	}

	return append(make([]interface{}, 0), item)
}
