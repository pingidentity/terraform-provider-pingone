package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne groups",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the group in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the group.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"population_id": {
				Description:      "The ID of the population that the group should be assigned to.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"user_filter": {
				Description: "A SCIM filter to dynamically assign users to the group.  Examples are found in the [PingOne online documentation](https://docs.pingidentity.com/bundle/pingone/page/kti1564020489340.html).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"external_id": {
				Description: "A user defined ID that represents the counterpart group in an external system.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	group := *management.NewGroup(d.Get("name").(string)) // Group |  (optional)

	if v, ok := d.GetOk("description"); ok {
		group.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("population_id"); ok {
		groupPopulation := *management.NewGroupPopulation(v.(string)) // NewGroupPopulation |  (optional)
		group.SetPopulation(groupPopulation)
	}

	if v, ok := d.GetOk("user_filter"); ok {
		group.SetUserFilter(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		group.SetExternalId(v.(string))
	}

	resp, r, err := apiClient.GroupsApi.CreateGroup(ctx, d.Get("environment_id").(string)).Group(group).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `GroupsApi.CreateGroup``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.GroupsApi.ReadOneGroup(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Group %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `GroupsApi.ReadOneGroup``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := resp.GetPopulationOk(); ok {
		d.Set("population_id", v.GetId())
	} else {
		d.Set("population_id", nil)
	}

	if v, ok := resp.GetUserFilterOk(); ok {
		d.Set("user_filter", v)
	} else {
		d.Set("user_filter", nil)
	}

	if v, ok := resp.GetExternalIdOk(); ok {
		d.Set("external_id", v)
	} else {
		d.Set("external_id", nil)
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	group := *management.NewGroup(d.Get("name").(string)) // Group |  (optional)

	if v, ok := d.GetOk("description"); ok {
		group.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("population_id"); ok {
		groupPopulation := *management.NewGroupPopulation(v.(string)) // NewGroupPopulation |  (optional)
		group.SetPopulation(groupPopulation)
	}

	if v, ok := d.GetOk("user_filter"); ok {
		group.SetUserFilter(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		group.SetExternalId(v.(string))
	}

	_, r, err := apiClient.GroupsApi.UpdateGroup(ctx, d.Get("environment_id").(string), d.Id()).Group(group).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `GroupsApi.UpdateGroup``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.GroupsApi.DeleteGroup(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `GroupsApi.DeleteGroup``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/groupID\"", d.Id())
	}

	environmentID, groupID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(groupID)

	resourceGroupRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
