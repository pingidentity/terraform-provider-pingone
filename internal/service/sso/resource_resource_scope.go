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
				Description: "A list that specifies the user schema attributes that can be read or updated for the specified PingOne access control scope. The value is an array of schema attribute paths (such as `username`, `name.given`, `shirtSize`) that the scope controls. This property is supported only for the `p1:read:user`, `p1:update:user` and `p1:read:user:{suffix}` and `p1:update:user:{suffix}` scopes. No other PingOne platform scopes allow this behavior. Any attributes not listed in the attribute array are excluded from the read or update action. The wildcard path (`*`) in the array includes all attributes and cannot be used in conjunction with any other user schema attribute paths.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"schema_attributes", "mapped_claims"},
			},
			"mapped_claims": {
				Description: "A list of custom resource attribute IDs. This property applies only for the resource with its type property set to `OPENID_CONNECT`. Moreover, this property does not display predefined OpenID Connect (OIDC) mappings, such as the `email` claim in the OIDC `email` scope or the `name` claim in the `profile` scope. You can create custom attributes, and these custom attributes can be added to `mapped_claims` and will display in the response.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				},
				ConflictsWith: []string{"schema_attributes", "mapped_claims"},
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
		if v1, ok := v.(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			items := make([]string, 0)

			for _, item := range v1.List() {
				items = append(items, item.(string))
			}

			resourceScope.SetSchemaAttributes(items)
		}
	}

	if v, ok := d.GetOk("mapped_claims"); ok {
		if v1, ok := v.(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			items := make([]string, 0)

			for _, item := range v1.List() {
				items = append(items, item.(string))
			}

			resourceScope.SetMappedClaims(items)
		}
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.CreateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).ResourceScope(resourceScope).Execute()
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

	if v, ok := respObject.GetSchemaAttributesOk(); ok {
		d.Set("schema_attributes", v)
	} else {
		d.Set("schema_attributes", nil)
	}

	if v, ok := respObject.GetMappedClaimsOk(); ok {
		d.Set("mapped_claims", v)
	} else {
		d.Set("mapped_claims", nil)
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
		if v1, ok := v.(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			items := make([]string, 0)

			for _, item := range v1.List() {
				items = append(items, item.(string))
			}

			resourceScope.SetSchemaAttributes(items)
		}
	}

	if v, ok := d.GetOk("mapped_claims"); ok {
		if v1, ok := v.(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			items := make([]string, 0)

			for _, item := range v1.List() {
				items = append(items, item.(string))
			}

			resourceScope.SetMappedClaims(items)
		}
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).ResourceScope(resourceScope).Execute()
		},
		"UpdateResourceScope",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
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

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteResourceScope",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
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
