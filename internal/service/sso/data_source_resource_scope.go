package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceResourceScope() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne OAuth 2.0 resource scope data.",

		ReadContext: datasourcePingOneResourceScopeRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_id": {
				Description:      "The ID of the resource that the scope belongs to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_scope_id": {
				Description:      "The ID of the resource scope.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the resource scope.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_scope_id"},
			},
			"description": {
				Description: "A description of the resource scope.",
				Type:        schema.TypeString,
				Computed:    true,
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

func datasourcePingOneResourceScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.ResourceScope

	if v, ok := d.GetOk("name"); ok {

		respList, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceScopesApi.ReadAllResourceScopes(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).Execute()
			},
			"ReadAllResourceScopes",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		if scopes, ok := respList.(*management.EntityArray).Embedded.GetScopesOk(); ok {

			found := false
			for _, scope := range scopes {

				if scope.GetName() == v.(string) {
					resp = scope
					found = true
					break
				}
			}

			if !found {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Cannot find resource scope %s", v),
				})

				return diags
			}

		}

	} else if v, ok2 := d.GetOk("resource_scope_id"); ok2 {

		resourceResp, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceScopesApi.ReadOneResourceScope(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), v.(string)).Execute()
			},
			"ReadOneResourceScope",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		resp = *resourceResp.(*management.ResourceScope)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither resource_scope_id or name are set",
			Detail:   "Neither resource_scope_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("resource_scope_id", resp.GetId())
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
