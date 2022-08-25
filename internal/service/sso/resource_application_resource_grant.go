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
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the resource grant for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_id": {
				Description:      "The ID of the protected resource associated with this grant.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scopes": {
				Description: "A list of IDs of the scopes associated with this grant.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				},
				Set: schema.HashString,
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
	scopes := expandApplicationResourceGrant(d.Get("scopes").(*schema.Set))

	applicationResourceGrant := *management.NewApplicationResourceGrant(resource, scopes)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationResourceGrantsApi.CreateApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).ApplicationResourceGrant(applicationResourceGrant).Execute()
		},
		"CreateApplicationGrant",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.ApplicationResourceGrant)

	d.SetId(respObject.GetId())

	return resourcePingOneApplicationResourceGrantRead(ctx, d, meta)
}

func resourcePingOneApplicationResourceGrantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
		},
		"ReadOneApplicationGrant",
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

	respObject := resp.(*management.ApplicationResourceGrant)

	d.Set("resource_id", respObject.Resource.GetId())
	d.Set("scopes", flattenAppResourceGrantScopes(respObject.GetScopes()))

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
	scopes := expandApplicationResourceGrant(d.Get("scopes").(*schema.Set))

	applicationResourceGrant := *management.NewApplicationResourceGrant(resource, scopes)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationResourceGrantsApi.UpdateApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).ApplicationResourceGrant(applicationResourceGrant).Execute()
		},
		"UpdateApplicationGrant",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
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

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ApplicationResourceGrantsApi.DeleteApplicationGrant(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteApplicationGrant",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOneApplicationResourceGrantImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

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

func expandApplicationResourceGrant(scopesIn *schema.Set) []management.ApplicationResourceGrantScopesInner {

	scopes := make([]management.ApplicationResourceGrantScopesInner, 0, len(scopesIn.List()))
	for _, scope := range scopesIn.List() {
		scopes = append(scopes, management.ApplicationResourceGrantScopesInner{
			Id: scope.(string),
		})
	}

	return scopes
}

func flattenAppResourceGrantScopes(in []management.ApplicationResourceGrantScopesInner) []string {

	items := make([]string, 0, len(in))
	for _, v := range in {

		items = append(items, v.GetId())
	}

	return items
}
