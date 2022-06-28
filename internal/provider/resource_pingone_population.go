package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go"
)

func resourcePingOnePopulation() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne populations",

		CreateContext: resourcePingOnePopulationCreate,
		ReadContext:   resourcePingOnePopulationRead,
		UpdateContext: resourcePingOnePopulationUpdate,
		DeleteContext: resourcePingOnePopulationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOnePopulationImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description: "The ID of the environment to create the population in.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the population",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description to apply to the population",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourcePingOnePopulationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	description := ""

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	resp, _, err := pingOnePopulationCreate(ctx, apiClient, d.Get("environment_id").(string), d.Get("name").(string), description)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetId())

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := pingOnePopulationRead(ctx, apiClient, d.Get("environment_id").(string), d.Id())
	if err != nil {
		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Population %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", resp.GetName())
	d.Set("description", resp.GetDescription())

	return diags
}

func resourcePingOnePopulationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	description := ""

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	_, _, err := pingOnePopulationUpdate(ctx, apiClient, d.Get("environment_id").(string), d.Id(), d.Get("name").(string), description)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ManagementAPIsPopulationsApi.DeletePopulation(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsPopulationsApi.DeletePopulation``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOnePopulationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
	}

	environmentID, populationID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(populationID)

	resourcePingOnePopulationRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func pingOnePopulationCreate(ctx context.Context, apiClient *pingone.APIClient, environmentID string, name string, description string) (*pingone.Population, *http.Response, error) {
	var diags diag.Diagnostics

	log.Printf("[INFO] Creating PingOne Population: name %s, environment: %s", name, environmentID)

	population := *pingone.NewPopulation(name) // Population |  (optional)

	if description != "" {
		population.SetDescription(description)
	}

	resp, r, err := apiClient.ManagementAPIsPopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsPopulationsApi.CreatePopulation``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return nil, r, err
	}

	return resp, r, nil
}

func pingOnePopulationRead(ctx context.Context, apiClient *pingone.APIClient, environmentID string, populationID string) (*pingone.Population, *http.Response, error) {
	var diags diag.Diagnostics

	log.Printf("[INFO] Reading PingOne Population: populationID %s", populationID)

	resp, r, err := apiClient.ManagementAPIsPopulationsApi.ReadOnePopulation(ctx, environmentID, populationID).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsPopulationsApi.ReadOnePopulation``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return nil, r, err
	}

	return resp, r, nil
}

func pingOnePopulationUpdate(ctx context.Context, apiClient *pingone.APIClient, environmentID string, populationID string, name string, description string) (*pingone.Population, *http.Response, error) {
	var diags diag.Diagnostics

	log.Printf("[INFO] Updating PingOne Population: name %s", name)

	population := *pingone.NewPopulation(name) // Population |  (optional)

	if description != "" {
		population.SetDescription(description)
	} else {
		population.SetDescription("1")
	}

	_, r, err := apiClient.ManagementAPIsPopulationsApi.UpdatePopulation(ctx, environmentID, populationID).Population(population).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsPopulationsApi.UpdatePopulation``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return nil, r, err
	}

	return nil, r, nil
}
