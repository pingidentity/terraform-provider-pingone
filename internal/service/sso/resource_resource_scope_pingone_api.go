package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceResourceScopePingOneAPI() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource scopes for the PingOne API resource.  Predefined scopes of `p1:read:user` and `p1:update:user` can be overridden, and new scopes can be defined as subscopes in the format `p1:read:user:{suffix}` or `p1:update:user:{suffix}`.  E.g. `p1:read:user:newscope` or `p1:update:user:newscope`.",

		CreateContext: resourceResourceScopePingOneAPICreate,
		ReadContext:   resourceResourceScopePingOneAPIRead,
		UpdateContext: resourceResourceScopePingOneAPIUpdate,
		DeleteContext: resourceResourceScopePingOneAPIDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceScopePingOneAPIImport,
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
				Description:      "The name of the resource scope.  Predefined scopes of `p1:read:user` and `p1:update:user` can be overridden, and new scopes can be defined as subscopes in the format `p1:read:user:{suffix}` or `p1:update:user:{suffix}`.  E.g. `p1:read:user:newscope` or `p1:update:user:newscope`",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^p1:(read|update):user(:{1}[a-zA-Z0-9]+)*$`), "Resource scope name must be either `p1:read:user`, `p1:update:user`, `p1:read:user:{suffix}` or `p1:update:user:{suffix}`")),
			},
			"description": {
				Description: "A description to apply to the resource scope.  The description can only be set when defining new scopes.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"schema_attributes": {
				Description: "A list that specifies the user schema attributes that can be read or updated for the specified PingOne access control scope. The value is an array of schema attribute paths (such as `username`, `name.given`, `shirtSize`) that the scope controls. This property is supported only for the `p1:read:user`, `p1:update:user` and `p1:read:user:{suffix}` and `p1:update:user:{suffix}` scopes. No other PingOne platform scopes allow this behavior. Any attributes not listed in the attribute array are excluded from the read or update action. The wildcard path (`*`) in the array includes all attributes and cannot be used in conjunction with any other user schema attribute paths.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_id": {
				Description: "The ID of the PingOne API resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceResourceScopePingOneAPICreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resource, diags := getPingOneAPIResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope, diags := expandResourceScopePingOneAPI(ctx, apiClient, d, resource.GetId())
	if diags.HasError() {
		return diags
	}

	var resp interface{}

	if v, ok := resourceScope.GetIdOk(); ok {

		resp, diags = sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), *v).ResourceScope(*resourceScope).Execute()
			},
			"UpdateResourceScope-PingOneAPI-Create",
			sdk.DefaultCustomError,
			nil,
		)

	} else {

		resp, diags = sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ResourceScopesApi.CreateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId()).ResourceScope(*resourceScope).Execute()
			},
			"CreateResourceScope-PingOneAPI",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	}

	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.ResourceScope)

	d.SetId(respObject.GetId())

	return resourceResourceScopePingOneAPIRead(ctx, d, meta)
}

func resourceResourceScopePingOneAPIRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resource, diags := getPingOneAPIResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourceScopesApi.ReadOneResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).Execute()
		},
		"ReadOneResourceScope-PingOneAPI",
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

	d.Set("resource_id", resource.GetId())

	return diags
}

func resourceResourceScopePingOneAPIUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resource, diags := getPingOneAPIResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	resourceScope, diags := expandResourceScopePingOneAPI(ctx, apiClient, d, resource.GetId())
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).ResourceScope(*resourceScope).Execute()
		},
		"UpdateResourceScope-PingOneAPI",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceResourceScopePingOneAPIRead(ctx, d, meta)
}

func resourceResourceScopePingOneAPIDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resource, diags := getPingOneAPIResource(ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return diags
	}

	if m, err := regexp.MatchString("^p1:(read|update):user$", d.Get("name").(string)); err == nil && m {

		resourceScope, diags := fetchResourceScopeFromName(ctx, apiClient, d.Get("environment_id").(string), resource.GetId(), d.Get("name").(string))
		if diags.HasError() {
			return diags
		}

		resourceScope.SetSchemaAttributes([]string{"*"})

		_, diags = sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ResourceScopesApi.UpdateResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).ResourceScope(*resourceScope).Execute()
			},
			"UpdateResourceScope-PingOneAPI-Delete",
			sdk.DefaultCustomError,
			nil,
		)
		if diags.HasError() {
			return diags
		}

	} else {
		_, diags = sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				r, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, d.Get("environment_id").(string), resource.GetId(), d.Id()).Execute()
				return nil, r, err
			},
			"DeleteResourceScope-PingOneAPI",
			sdk.CustomErrorResourceNotFoundWarning,
			nil,
		)
		if diags.HasError() {
			return diags
		}
	}

	return diags
}

func resourceResourceScopePingOneAPIImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "resource_scope_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.SetId(attributes["resource_scope_id"])

	resourceResourceScopePingOneAPIRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandResourceScopePingOneAPI(ctx context.Context, apiClient *management.APIClient, d *schema.ResourceData, resourceID string) (*management.ResourceScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resourceScope *management.ResourceScope

	newScope := true
	if m, err := regexp.MatchString("^p1:(read|update):user$", d.Get("name").(string)); err == nil && m {
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

	if v, ok := d.GetOk("schema_attributes"); ok {

		if v1, ok := v.(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			items := make([]string, 0)

			for _, item := range v1.List() {
				items = append(items, item.(string))
			}

			resourceScope.SetSchemaAttributes(items)
		}

	}

	return resourceScope, diags
}
