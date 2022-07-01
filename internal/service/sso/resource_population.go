package sso

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourcePopulation() *schema.Resource {
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
				Description:  "The ID of the environment to create the population in.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},
			"name": {
				Description:  "The name of the population.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Description: "A description to apply to the population.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// "password_policy_id": {
			// 	Description: "The ID of a password policy to assign to the population",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
		},
	}
}

func resourcePingOnePopulationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	population := *pingone.NewPopulation(d.Get("name").(string)) // Population |  (optional)

	if v, ok := d.GetOk("description"); ok {
		population.SetDescription(v.(string))
	}

	// if v, ok := d.GetOk("password_policy_id"); ok {
	// 	population.SetPasswordPolicyInnerId(v.(string))
	// }

	resp, _, err := PingOnePopulationCreate(ctx, apiClient, d.Get("environment_id").(string), population)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetId())

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := PingOnePopulationRead(ctx, apiClient, d.Get("environment_id").(string), d.Id())
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
	//d.Set("password_policy_id", resp.SetPasswordPolicyInnerId())

	return diags
}

func resourcePingOnePopulationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	population := *pingone.NewPopulation(d.Get("name").(string)) // Population |  (optional)

	if v, ok := d.GetOk("description"); ok {
		population.SetDescription(v.(string))
	}

	// if v, ok := d.GetOk("password_policy_id"); ok {
	// 	population.SetPasswordPolicyInnerId(v.(string))
	// }

	_, _, err := PingOnePopulationUpdate(ctx, apiClient, d.Get("environment_id").(string), d.Id(), population)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
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

func PingOnePopulationCreate(ctx context.Context, apiClient *pingone.APIClient, environmentID string, population pingone.Population) (*pingone.Population, *http.Response, error) {

	resp, r, err := apiClient.ManagementAPIsPopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
	if (err != nil) || (r.StatusCode != 201) {

		return nil, r, err
	}

	return resp, r, nil
}

func PingOnePopulationRead(ctx context.Context, apiClient *pingone.APIClient, environmentID string, populationID string) (*pingone.Population, *http.Response, error) {

	resp, r, err := apiClient.ManagementAPIsPopulationsApi.ReadOnePopulation(ctx, environmentID, populationID).Execute()
	if err != nil {

		return nil, r, err
	}

	return resp, r, nil
}

func PingOnePopulationUpdate(ctx context.Context, apiClient *pingone.APIClient, environmentID string, populationID string, population pingone.Population) (*pingone.Population, *http.Response, error) {

	_, r, err := apiClient.ManagementAPIsPopulationsApi.UpdatePopulation(ctx, environmentID, populationID).Population(population).Execute()
	if err != nil {

		return nil, r, err
	}

	return nil, r, nil
}
