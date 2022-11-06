package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceResourceScopeOpenID() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource scopes for the OpenID Connect resource.  Predefined scopes can be overridden, and new scopes can be defined.",

		CreateContext: resourceResourceScopeOpenIDCreate,
		ReadContext:   resourceResourceScopeOpenIDRead,
		UpdateContext: resourceResourceScopeOpenIDUpdate,
		DeleteContext: resourceResourceScopeOpenIDDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceScopeOpenIDImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the resource scope in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the resource scope.  Predefined scopes of `address`, `email`, `openid`, `phone` and `profile` can be overridden, and new scopes can be defined.  E.g. `myawesomescope`",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the resource scope.  The description can only be set when defining new scopes.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"mapped_claims": {
				Description: "A list of custom resource attribute IDs.  This property does not control predefined OpenID Connect (OIDC) mappings, such as the `email` claim in the OIDC `email` scope or the `name` claim in the `profile` scope. You can create custom attributes, and these custom attributes can be added to `mapped_claims` and will display in the response.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				},
			},
			"resource_id": {
				Description: "The ID of the OpenID Connect resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceResourceScopeOpenIDCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource, diags := getOpenIDResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope, diags := expandResourceScopeOpenID(ctx, apiClient, d, resource.GetId())
	if diags.HasError() {
		return diags
	}

	var resp interface{}

	if v, ok := resourceScope.GetIdOk(); ok {

		resp, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), *v).ResourceScope(*resourceScope).Execute()
			},
			"UpdateResourceScope-OpenID-Create",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)

	} else {

		resp, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceScopesApi.CreateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId()).ResourceScope(*resourceScope).Execute()
			},
			"CreateResourceScope-OpenID",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	}

	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.ResourceScope)

	d.SetId(respObject.GetId())

	return resourceResourceScopeOpenIDRead(ctx, d, meta)
}

func resourceResourceScopeOpenIDRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource, diags := getOpenIDResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.ReadOneResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).Execute()
		},
		"ReadOneResourceScope-OpenID",
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

	if v, ok := respObject.GetMappedClaimsOk(); ok {
		d.Set("mapped_claims", v)
	} else {
		d.Set("mapped_claims", nil)
	}

	d.Set("resource_id", resource.GetId())

	return diags
}

func resourceResourceScopeOpenIDUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource, diags := getOpenIDResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope, diags := expandResourceScopeOpenID(ctx, apiClient, d, resource.GetId())
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).ResourceScope(*resourceScope).Execute()
		},
		"UpdateResourceScope-OpenID",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceResourceScopeOpenIDRead(ctx, d, meta)
}

func resourceResourceScopeOpenIDDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource, diags := getOpenIDResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	if m, err := regexp.MatchString("^(address|email|openid|phone|profile)$", d.Get("name").(string)); err == nil && m {

		resourceScope, diags := fetchResourceScopeFromName(ctx, apiClient, d.Get("environment_id").(string), resource.GetId(), d.Get("name").(string))
		if diags.HasError() {
			return diags
		}

		resourceScope.SetMappedClaims([]string{})

		_, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).ResourceScope(*resourceScope).Execute()
			},
			"UpdateResourceScope-OpenID-Delete",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}

	} else {
		_, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				r, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).Execute()
				return nil, r, err
			},
			"DeleteResourceScope-OpenID",
			sdk.CustomErrorResourceNotFoundWarning,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}

	}

	return diags
}

func resourceResourceScopeOpenIDImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/resourceScopeID\"", d.Id())
	}

	environmentID, resourceScopeID := attributes[0], attributes[2]

	d.Set("environment_id", environmentID)
	d.SetId(resourceScopeID)

	resourceResourceScopeOpenIDRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandResourceScopeOpenID(ctx context.Context, apiClient *management.APIClient, d *schema.ResourceData, resourceID string) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resourceScope *management.ResourceScope

	newScope := true
	if m, err := regexp.MatchString("^(address|email|openid|phone|profile)$", d.Get("name").(string)); err == nil && m {
		newScope = false

		resourceScope, diags = fetchResourceScopeFromName(ctx, apiClient, d.Get("environment_id").(string), resourceID, d.Get("name").(string))
		if diags.HasError() {
			return nil, diags
		}

	} else {
		resourceScope = management.NewResourceScope(d.Get("name").(string)) // ResourceScope |  (optional)
	}

	if v, ok := d.GetOk("description"); ok {
		if newScope {
			resourceScope.SetDescription(v.(string))
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("The scope `%s` is an existing platform scope.  The description cannot be changed.", d.Get("name").(string)),
			})
			return nil, diags
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

	return resourceScope, diags
}
