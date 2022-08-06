package sso

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
// 	"github.com/patrickcping/pingone-go-sdk-v2/management"
// 	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
// )

// func ResourceSystemApplication() *schema.Resource {
// 	return &schema.Resource{

// 		// This description is used by the documentation generator and the language server.
// 		Description: "Resource to create and manage PingOne populations",

// 		CreateContext: resourcePingOneSystemApplicationCreate,
// 		ReadContext:   resourcePingOneSystemApplicationRead,
// 		UpdateContext: resourcePingOneSystemApplicationUpdate,
// 		DeleteContext: resourcePingOneSystemApplicationDelete,

// 		Importer: &schema.ResourceImporter{
// 			StateContext: resourcePingOneSystemApplicationImport,
// 		},

// 		Schema: map[string]*schema.Schema{
// 			"environment_id": {
// 				Description:      "The ID of the environment to create the population in.",
// 				Type:             schema.TypeString,
// 				Required:         true,
// 				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
// 				ForceNew:         true,
// 			},
// 			"name": {
// 				Description:      "The name of the population.",
// 				Type:             schema.TypeString,
// 				Required:         true,
// 				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
// 			},
// 			"description": {
// 				Description: "A description to apply to the population.",
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 			},
// 			"password_policy_id": {
// 				Description: "The ID of a password policy to assign to the population.",
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func resourcePingOneSystemApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	p1Client := meta.(*client.Client)
// 	apiClient := p1Client.API.ManagementAPIClient
// 	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
// 		"suffix": p1Client.API.Region.URLSuffix,
// 	})

// 	population := *management.NewSystemApplication(d.Get("name").(string)) // SystemApplication |  (optional)

// 	if v, ok := d.GetOk("description"); ok {
// 		population.SetDescription(v.(string))
// 	}

// 	if v, ok := d.GetOk("password_policy_id"); ok {
// 		populationPasswordPolicy := *management.NewSystemApplicationPasswordPolicy(v.(string))
// 		population.SetPasswordPolicy(populationPasswordPolicy)
// 	}

// 	resp, _, err := PingOneSystemApplicationCreate(ctx, apiClient, d.Get("environment_id").(string), population)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(resp.GetId())

// 	return resourcePingOneSystemApplicationRead(ctx, d, meta)
// }

// func resourcePingOneSystemApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	p1Client := meta.(*client.Client)
// 	apiClient := p1Client.API.ManagementAPIClient
// 	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
// 		"suffix": p1Client.API.Region.URLSuffix,
// 	})
// 	var diags diag.Diagnostics

// 	resp, r, err := PingOneSystemApplicationRead(ctx, apiClient, d.Get("environment_id").(string), d.Id())
// 	if err != nil {
// 		if r.StatusCode == 404 {
// 			log.Printf("[INFO] PingOne SystemApplication %s no longer exists", d.Id())
// 			d.SetId("")
// 			return nil
// 		}
// 		return diag.FromErr(err)
// 	}

// 	d.Set("name", resp.GetName())

// 	if v, ok := resp.GetDescriptionOk(); ok {
// 		d.Set("description", v)
// 	} else {
// 		d.Set("description", nil)
// 	}

// 	if v, ok := resp.GetPasswordPolicyOk(); ok {
// 		d.Set("password_policy_id", v.GetId())
// 	} else {
// 		d.Set("password_policy_id", nil)
// 	}

// 	return diags
// }

// func resourcePingOneSystemApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	p1Client := meta.(*client.Client)
// 	apiClient := p1Client.API.ManagementAPIClient
// 	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
// 		"suffix": p1Client.API.Region.URLSuffix,
// 	})

// 	population := *management.NewSystemApplication(d.Get("name").(string)) // SystemApplication |  (optional)

// 	if v, ok := d.GetOk("description"); ok {
// 		population.SetDescription(v.(string))
// 	}

// 	if v, ok := d.GetOk("password_policy_id"); ok {
// 		populationPasswordPolicy := *management.NewSystemApplicationPasswordPolicy(v.(string))
// 		population.SetPasswordPolicy(populationPasswordPolicy)
// 	}

// 	_, _, err := PingOneSystemApplicationUpdate(ctx, apiClient, d.Get("environment_id").(string), d.Id(), population)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	return resourcePingOneSystemApplicationRead(ctx, d, meta)
// }

// func resourcePingOneSystemApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	p1Client := meta.(*client.Client)
// 	apiClient := p1Client.API.ManagementAPIClient
// 	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
// 		"suffix": p1Client.API.Region.URLSuffix,
// 	})
// 	var diags diag.Diagnostics

// 	_, err := apiClient.SystemApplicationsApi.DeleteSystemApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
// 	if err != nil {
// 		diags = append(diags, diag.Diagnostic{
// 			Severity: diag.Error,
// 			Summary:  fmt.Sprintf("Error when calling `SystemApplicationsApi.DeleteSystemApplication``: %v", err),
// 		})

// 		return diags
// 	}

// 	return nil
// }

// func resourcePingOneSystemApplicationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
// 	attributes := strings.SplitN(d.Id(), "/", 2)

// 	if len(attributes) != 2 {
// 		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
// 	}

// 	environmentID, populationID := attributes[0], attributes[1]

// 	d.Set("environment_id", environmentID)
// 	d.SetId(populationID)

// 	resourcePingOneSystemApplicationRead(ctx, d, meta)

// 	return []*schema.ResourceData{d}, nil
// }

// func PingOneSystemApplicationCreate(ctx context.Context, apiClient *management.APIClient, environmentID string, population management.SystemApplication) (*management.SystemApplication, *http.Response, error) {

// 	resp, r, err := apiClient.SystemApplicationsApi.CreateSystemApplication(ctx, environmentID).SystemApplication(population).Execute()
// 	if (err != nil) || (r.StatusCode != 201) {

// 		return nil, r, err
// 	}

// 	return resp, r, nil
// }

// func PingOneSystemApplicationRead(ctx context.Context, apiClient *management.APIClient, environmentID string, populationID string) (*management.SystemApplication, *http.Response, error) {

// 	resp, r, err := apiClient.SystemApplicationsApi.ReadOneSystemApplication(ctx, environmentID, populationID).Execute()
// 	if err != nil {

// 		return nil, r, err
// 	}

// 	return resp, r, nil
// }

// func PingOneSystemApplicationUpdate(ctx context.Context, apiClient *management.APIClient, environmentID string, populationID string, population management.SystemApplication) (*management.SystemApplication, *http.Response, error) {

// 	_, r, err := apiClient.SystemApplicationsApi.UpdateSystemApplication(ctx, environmentID, populationID).SystemApplication(population).Execute()
// 	if err != nil {

// 		return nil, r, err
// 	}

// 	return nil, r, nil
// }
