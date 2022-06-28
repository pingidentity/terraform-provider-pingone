package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go"
)

func datasourcePingOneEnvironment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne environment data",

		ReadContext: datasourcePingOneEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:   "The ID of the environment",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Description:   "The name of the environment",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"environment_id"},
			},
			"description": {
				Description: "A description of the environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the environment to create.  Options are SANDBOX for a development/testing environment and PRODUCTION for environments that require protection from deletion.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"region": {
				Description: "The region to create the environment in.  Should be consistent with the PingOne organisation region",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"license_id": {
				Description: "An ID of a valid license to apply to the environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_population_id": {
				Description: "The ID of the environment's default population",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_population_name": {
				Description: "The name of the environment's default population",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_population_description": {
				Description: "A description to apply to the environment's default population",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	var resp pingone.Environment
	if v, ok := d.GetOk("name"); ok {

		limit := int32(1000)
		filter := fmt.Sprintf("name sw \"%s\"", v.(string)) // need the eq filter
		respList, r, err := apiClient.ManagementAPIsEnvironmentsApi.ReadAllEnvironments(ctx).Limit(limit).Filter(filter).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.ReadAllEnvironments``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}

		resp = respList.Embedded.GetEnvironments()[0]
		log.Printf("Environment found %s", resp.Name)

	} else if v, ok := d.GetOk("environment_id"); ok {

		resp, r, err := apiClient.ManagementAPIsEnvironmentsApi.ReadOneEnvironment(ctx, v.(string)).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.ReadOneEnvironment``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
		log.Printf("Environment found %s", resp.Name)
	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Neither environment_id or name are set"),
			Detail:   fmt.Sprintf("Neither environment_id or name are set"),
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("environment_id", resp.GetId())
	d.Set("name", resp.GetName())
	d.Set("description", resp.GetDescription())
	d.Set("type", resp.GetType())
	d.Set("region", resp.GetRegion())
	d.Set("license_id", resp.GetLicense().Id)

	return diags
}
