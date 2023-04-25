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

func ResourceResourceScope() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 resource scopes for custom resources.  This resource cannot manage PingOne API or OpenID Connect scopes.",

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

	diags = checkResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope := expandResourceScope(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.CreateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).ResourceScope(*resourceScope).Execute()
		},
		"CreateResourceScope",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.ResourceScope)

	d.SetId(respObject.GetId())

	return resourceResourceScopeRead(ctx, d, meta)
}

func resourceResourceScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	diags = checkResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string))
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.ReadOneResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
		},
		"ReadOneResourceScope",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.ResourceScope)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
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

	diags = checkResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope := expandResourceScope(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).ResourceScope(*resourceScope).Execute()
		},
		"UpdateResourceScope",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
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

	diags = checkResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string))
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteResourceScope",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceResourceScopeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/resourceID/resourceScopeID\"", d.Id())
	}

	environmentID, resourceID, resourceScopeID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("resource_id", resourceID)
	d.SetId(resourceScopeID)

	resourceResourceScopeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func checkResourceType(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceType, diags := getResourceType(ctx, apiClient, environmentID, resourceID)
	if diags.HasError() {
		return diags
	}

	if resourceType != management.ENUMRESOURCETYPE_CUSTOM {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot control scopes for resources that are of type PingOne API or OpenID Connect.  Please ensure that the resource in the `resource_id` parameter is a custom resource.",
		})
		return diags

	}

	return diags
}

func expandResourceScope(d *schema.ResourceData) *management.ResourceScope {
	resourceScope := *management.NewResourceScope(d.Get("name").(string)) // ResourceScope |  (optional)

	if v, ok := d.GetOk("description"); ok {
		resourceScope.SetDescription(v.(string))
	}

	return &resourceScope
}
