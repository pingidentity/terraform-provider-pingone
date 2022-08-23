package sso

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceResource() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne OAuth 2.0 resource data",

		ReadContext: datasourcePingOneResourceRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_id": {
				Description:      "The ID of the resource.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the resource.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_id"},
			},
			"description": {
				Description: "A description of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "A string that specifies the type of resource. Options are `OPENID_CONNECT`, `PINGONE_API`, and `CUSTOM`. Only the `CUSTOM` resource type can be created. `OPENID_CONNECT` specifies the built-in platform resource for OpenID Connect. `PINGONE_API` specifies the built-in platform resource for PingOne.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"audience": {
				Description: "A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"access_token_validity_seconds": {
				Description: "An integer that specifies the number of seconds that the access token is valid.  The minimum value is 300 seconds (5 minutes); the maximum value is 2592000 seconds (30 days).",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.Resource

	if v, ok := d.GetOk("name"); ok {

		respList, r, err := apiClient.ResourcesResourcesApi.ReadAllResources(ctx, d.Get("environment_id").(string)).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ResourcesResourcesApi.ReadAllResources``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}

		if resources, ok := respList.Embedded.GetResourcesOk(); ok {

			found := false
			for _, resource := range resources {

				if resource.GetName() == v.(string) {
					resp = resource
					found = true
					break
				}
			}

			if !found {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Cannot find resource %s", v),
				})

				return diags
			}

		}

	} else if v, ok2 := d.GetOk("resource_id"); ok2 {

		resourceResp, r, err := apiClient.ResourcesResourcesApi.ReadOneResource(ctx, d.Get("environment_id").(string), v.(string)).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ResourcesResourcesApi.ReadOneResource``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}

		resp = *resourceResp

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither resource_id or name are set",
			Detail:   "Neither resource_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("resource_id", resp.GetId())
	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := resp.GetTypeOk(); ok {
		d.Set("type", string(*v))
	} else {
		d.Set("type", nil)
	}

	if v, ok := resp.GetAudienceOk(); ok {
		d.Set("audience", v)
	} else {
		d.Set("audience", nil)
	}

	if v, ok := resp.GetAccessTokenValiditySecondsOk(); ok {
		d.Set("access_token_validity_seconds", v)
	} else {
		d.Set("access_token_validity_seconds", nil)
	}

	return diags
}
