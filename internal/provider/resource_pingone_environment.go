package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go"
)

func resourcePingOneEnvironment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne environments.",

		CreateContext: resourcePingOneEnvironmentCreate,
		ReadContext:   resourcePingOneEnvironmentRead,
		UpdateContext: resourcePingOneEnvironmentUpdate,
		DeleteContext: resourcePingOneEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneEnvironmentImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the environment",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description of the environment",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  "The type of the environment to create.  Options are SANDBOX for a development/testing environment and PRODUCTION for environments that require protection from deletion.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SANDBOX",
				ValidateFunc: validation.StringInSlice([]string{"PRODUCTION", "SANDBOX"}, false),
			},
			"region": {
				Description:  "The region to create the environment in.  Should be consistent with the PingOne organisation region",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NA", "EU", "ASIA", "CA"}, false),
				ForceNew:     true,
			},
			"license_id": {
				Description: "An ID of a valid license to apply to the environment",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"default_population_id": {
				Description: "The ID of the environment's default population",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_population_name": {
				Description: "The name of the environment's default population",
				Type:        schema.TypeString,
				Default:     "Default",
				Optional:    true,
			},
			"default_population_description": {
				Description: "A description to apply to the environment's default population",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

var billOfMaterialsProductElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"PING_ONE_DAVINCI", "PING_ONE_MFA", "PING_ID", "PING_ONE_RISK", "PING_ONE_VERIFY", "PING_ONE_CREDENTIALS", "PING_INTELLIGENCE", "PING_ONE_AUTHORIZE", "PING_ONE_FRAUD", "PING_FEDERATE", "PING_ACCESS", "PING_DIRECTORY", "PING_AUTHORIZE", "PING_CENTRAL"}, false),
			Required:     true,
		},
		"solution": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"WORKFORCE", "CUSTOMER"}, false),
			Optional:     true,
		},
		"console_href": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"bookmark": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"href": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	},
}

func resourcePingOneEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	var diags diag.Diagnostics

	var environmentLicense pingone.EnvironmentLicense
	if license, ok := d.GetOk("license_id"); ok {
		environmentLicense = *pingone.NewEnvironmentLicense(license.(string))
	}

	environment := *pingone.NewEnvironment(environmentLicense, d.Get("name").(string), d.Get("region").(string), d.Get("type").(string)) // Environment |  (optional)

	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	}

	resp, r, err := apiClient.ManagementAPIsEnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(environment).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.CreateEnvironmentActiveLicense``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	populationName := "Default"
	if v, ok := d.GetOk("default_population_name"); ok {
		populationName = v.(string)
	}

	populationDescription := ""
	if v, ok := d.GetOk("default_population_description"); ok {
		populationDescription = v.(string)
	}

	//Have to create a default population because of the destroy restriction on the population resource
	populationResp, _, err := pingOnePopulationCreate(ctx, apiClient, resp.GetId(), populationName, populationDescription)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetId())
	d.Set("default_population_id", populationResp.GetId())

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	resp, r, err := apiClient.ManagementAPIsEnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Environment no %s longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.ReadOneEnvironment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())
	d.Set("description", resp.GetDescription())
	d.Set("type", resp.GetType())
	d.Set("region", resp.GetRegion())
	d.Set("license_id", resp.GetLicense().Id)

	populationResp, populationR, populationErr := pingOnePopulationRead(ctx, apiClient, environmentID, populationID)
	if populationErr != nil {

		if populationR.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application Default Population no %s longer exists", populationID)
			d.Set("default_population_id", "")
			d.Set("default_population_name", "")
			d.Set("default_population_description", "")
			return diags
		}

		return diag.FromErr(populationErr)
	}

	d.Set("default_population_id", populationResp.GetId())
	d.Set("default_population_name", populationResp.GetName())
	d.Set("default_population_description", populationResp.GetDescription())

	return diags
}

func resourcePingOneEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	var environmentLicense pingone.EnvironmentLicense
	if v, ok := d.GetOk("license_id"); ok {
		environmentLicense = *pingone.NewEnvironmentLicense(v.(string))
	}

	environment := *pingone.NewEnvironment(environmentLicense, d.Get("name").(string), d.Get("region").(string), d.Get("type").(string)) // Environment |  (optional)
	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	} else {
		environment.SetDescription("")
	}

	if change := d.HasChange("type"); change {
		//If type has changed from SANDBOX -> PRODUCTION and vice versa we need a separate API call
		updateEnvironmentTypeRequest := *pingone.NewUpdateEnvironmentTypeRequest()
		_, newType := d.GetChange("type")
		updateEnvironmentTypeRequest.SetType(newType.(string))
		_, r, err := apiClient.ManagementAPIsEnvironmentsApi.UpdateEnvironmentType(ctx, environmentID).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.UpdateEnvironmentType``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	_, r, err := apiClient.ManagementAPIsEnvironmentsApi.UpdateEnvironment(ctx, environmentID).Environment(environment).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.UpdateEnvironment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	if change := d.HasChange("default_population_name") || d.HasChange("default_population_description"); change {

		populationName := "Default"
		if v, ok := d.GetOk("default_population_name"); ok {
			populationName = v.(string)
		}

		populationDescription := ""
		if v, ok := d.GetOk("default_population_description"); ok {
			populationDescription = v.(string)
		}

		_, _, err := pingOnePopulationUpdate(ctx, apiClient, environmentID, populationID, populationName, populationDescription)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ManagementAPIsEnvironmentsApi.DeleteEnvironment(ctx, d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ManagementAPIsEnvironmentsApi.DeleteEnvironment``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneEnvironmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
	}

	environmentID, populationID := attributes[0], attributes[1]

	d.SetId(environmentID)
	d.Set("default_population_id", populationID)

	resourcePingOneEnvironmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
