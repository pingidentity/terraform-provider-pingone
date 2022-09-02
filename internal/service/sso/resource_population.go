package sso

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
				Description:      "The ID of the environment to create the population in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the population.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the population.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password_policy_id": {
				Description:      "The ID of a password policy to assign to the population.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
		},
	}
}

func resourcePingOnePopulationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	population := *management.NewPopulation(d.Get("name").(string)) // Population |  (optional)

	if v, ok := d.GetOk("description"); ok {
		population.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("password_policy_id"); ok {
		populationPasswordPolicy := *management.NewPopulationPasswordPolicy(v.(string))
		population.SetPasswordPolicy(populationPasswordPolicy)
	}

	resp, diags := PingOnePopulationCreate(ctx, apiClient, d.Get("environment_id").(string), population)
	if diags.HasError() {
		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := PingOnePopulationRead(ctx, apiClient, d.Get("environment_id").(string), d.Id())
	if diags.HasError() {
		return diags
	}

	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := resp.GetPasswordPolicyOk(); ok {
		d.Set("password_policy_id", v.GetId())
	} else {
		d.Set("password_policy_id", nil)
	}

	return diags
}

func resourcePingOnePopulationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	population := *management.NewPopulation(d.Get("name").(string)) // Population |  (optional)

	if v, ok := d.GetOk("description"); ok {
		population.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("password_policy_id"); ok {
		populationPasswordPolicy := *management.NewPopulationPasswordPolicy(v.(string))
		population.SetPasswordPolicy(populationPasswordPolicy)
	}

	_, diags := PingOnePopulationUpdate(ctx, apiClient, d.Get("environment_id").(string), d.Id(), population)
	if diags.HasError() {
		return diags
	}

	return resourcePingOnePopulationRead(ctx, d, meta)
}

func resourcePingOnePopulationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.PopulationsApi.DeletePopulation(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeletePopulation",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOnePopulationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
	}

	environmentID, populationID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(populationID)

	resourcePingOnePopulationRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func PingOnePopulationCreate(ctx context.Context, apiClient *management.APIClient, environmentID string, population management.Population) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
		},
		"CreatePopulation",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	respObject := resp.(*management.Population)

	return respObject, diags
}

func PingOnePopulationRead(ctx context.Context, apiClient *management.APIClient, environmentID string, populationID string) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PopulationsApi.ReadOnePopulation(ctx, environmentID, populationID).Execute()
		},
		"ReadOnePopulation",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	respObject := resp.(*management.Population)

	return respObject, diags
}

func PingOnePopulationUpdate(ctx context.Context, apiClient *management.APIClient, environmentID string, populationID string, population management.Population) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PopulationsApi.UpdatePopulation(ctx, environmentID, populationID).Population(population).Execute()
		},
		"UpdatePopulation",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	respObject := resp.(*management.Population)

	return respObject, diags
}
