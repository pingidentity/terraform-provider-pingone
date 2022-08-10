package sso

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceApplicationResourceGrant() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a resource grant for an application configured in PingOne.",

		CreateContext: resourcePingOneApplicationResourceGrantCreate,
		ReadContext:   resourcePingOneApplicationResourceGrantRead,
		UpdateContext: resourcePingOneApplicationResourceGrantUpdate,
		DeleteContext: resourcePingOneApplicationResourceGrantDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationResourceGrantImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application resource grant in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the resource grant for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"resource_id": {
				Description:      "The ID of the protected resource associated with this grant.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"scopes": {
				Description: "A list of IDs of the scopes associated with this grant.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePingOneApplicationResourceGrantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource := *management.NewApplicationResourceGrantResource(d.Get("resource_id").(string))
	scopes := expandApplicationResourceGrant(d.Get("scopes").([]interface{}))

	applicationResourceGrant := *management.NewApplicationResourceGrant(resource, scopes)

	resp, r, err := apiClient.ApplicationsApplicationResourceGrantsApi.CreateApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).ApplicationResourceGrant(applicationResourceGrant).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationResourceGrantsApi.CreateApplicationGrant``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneApplicationResourceGrantRead(ctx, d, meta)
}

func resourcePingOneApplicationResourceGrantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.ApplicationsApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application Grant %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationResourceGrantsApi.ReadOneApplicationGrant``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("resource_id", resp.Resource.GetId())
	d.Set("scopes", flattenAppResourceGrantScopes(resp.GetScopes()))

	return diags
}

func resourcePingOneApplicationResourceGrantUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource := *management.NewApplicationResourceGrantResource(d.Get("resource_id").(string))
	scopes := expandApplicationResourceGrant(d.Get("scopes").([]interface{}))

	applicationResourceGrant := *management.NewApplicationResourceGrant(resource, scopes)

	_, r, err := apiClient.ApplicationsApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).ApplicationResourceGrant(applicationResourceGrant).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationResourceGrantsApi.UpdateApplicationGrant``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourcePingOneApplicationResourceGrantRead(ctx, d, meta)
}

func resourcePingOneApplicationResourceGrantDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ApplicationsApplicationResourceGrantsApi.DeleteApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationResourceGrantsApi.DeleteApplicationGrant``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneApplicationResourceGrantImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/grantID\"", d.Id())
	}

	environmentID, applicationID, grantID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(grantID)

	resourcePingOneApplicationResourceGrantRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandApplicationResourceGrant(scopesIn []interface{}) []management.ApplicationResourceGrantScopesInner {

	scopes := make([]management.ApplicationResourceGrantScopesInner, 0, len(scopesIn))
	for _, scope := range scopesIn {
		scopes = append(scopes, management.ApplicationResourceGrantScopesInner{
			Id: scope.(string),
		})
	}

	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i].GetId() < scopes[j].GetId()
	})

	return scopes
}

func flattenAppResourceGrantScopes(in []management.ApplicationResourceGrantScopesInner) []string {

	items := make([]string, 0, len(in))
	for _, v := range in {

		items = append(items, v.GetId())
	}

	sort.Strings(items)
	return items
}
