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

func ResourceResourceScope() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 resource scopes.",

		CreateContext: resourceResourceScopeCreate,
		ReadContext:   resourceResourceScopeRead,
		UpdateContext: resourceResourceScopeUpdate,
		DeleteContext: resourceResourceScopeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceScopeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the resource scope in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"resource_id": {
				Description:      "The ID of the resource to assign the resource scope to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the resource scope.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the resource scope.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"schema_attributes": {
				Description: "A list that specifies the user schema attributes that can be read or updated for the specified PingOne access control scope. The value is an array of schema attribute paths (such as username, name.given, shirtSize) that the scope controls. This property is supported only for the `p1:read:user`, `p1:update:user` and `p1:read:user:{suffix}` and `p1:update:user:{suffix}` scopes. No other PingOne platform scopes allow this behavior. Any attributes not listed in the attribute array are excluded from the read or update action. The wildcard path (*) in the array includes all attributes and cannot be used in conjunction with any other user schema attribute path.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceResourceScopeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resourceScope := *management.NewResourceScope(d.Get("name").(string)) // ResourceScope |  (optional)

	if v, ok := d.GetOk("description"); ok {
		resourceScope.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("schema_attributes"); ok {
		resourceScope.SetSchemaAttributes(v.([]string))
	}

	resp, r, err := apiClient.ResourcesResourceScopesApi.CreateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).ResourceScope(resourceScope).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ResourcesResourceScopesApi.CreateResourceScope``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourceResourceScopeRead(ctx, d, meta)
}

func resourceResourceScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.ResourcesResourceScopesApi.ReadOneResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Resource %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ResourcesResourceScopesApi.ReadOneResourceScope``: %v", err),
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

	if v, ok := resp.GetSchemaAttributesOk(); ok {
		d.Set("schema_attributes", v)
	} else {
		d.Set("schema_attributes", nil)
	}

	return diags
}

func resourceResourceScopeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resourceScope := *management.NewResourceScope(d.Get("name").(string)) // Resource |  (optional)

	if v, ok := d.GetOk("description"); ok {
		resourceScope.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("schema_attributes"); ok {
		resourceScope.SetSchemaAttributes(v.([]string))
	}

	_, r, err := apiClient.ResourcesResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).ResourceScope(resourceScope).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ResourcesResourceScopesApi.UpdateResourceScope``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceResourceScopeRead(ctx, d, meta)
}

func resourceResourceScopeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ResourcesResourceScopesApi.DeleteResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ResourcesResourceScopesApi.DeleteResourceScope``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceResourceScopeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/resourceID/resourceScopeID\"", d.Id())
	}

	environmentID, resourceID, resourceScopeID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("resource_id", resourceID)
	d.SetId(resourceScopeID)

	resourceResourceScopeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
