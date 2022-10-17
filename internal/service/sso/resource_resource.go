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

func ResourceResource() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 resources",

		CreateContext: resourceResourceCreate,
		ReadContext:   resourceResourceRead,
		UpdateContext: resourceResourceUpdate,
		DeleteContext: resourceResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the resource in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the resource.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description: "A string that specifies the type of resource. Options are `OPENID_CONNECT`, `PINGONE_API`, and `CUSTOM`. Only the `CUSTOM` resource type can be created. `OPENID_CONNECT` specifies the built-in platform resource for OpenID Connect. `PINGONE_API` specifies the built-in platform resource for PingOne.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"audience": {
				Description:      "A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringNotInSlice([]string{"pingone", "pingidentity"}, true)),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {

					if d.Get("name").(string) == old && new == "" {
						return true
					}

					return false
				},
			},
			"access_token_validity_seconds": {
				Description:      "An integer that specifies the number of seconds that the access token is valid. If a value is not specified, the default is 3600. The minimum value is 300 seconds (5 minutes); the maximum value is 2592000 seconds (30 days).",
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          3600,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(300, 2592000)),
			},
			"introspect_endpoint_auth_method": {
				Description:      fmt.Sprintf("The client authentication methods supported by the token endpoint. Options are `NONE`, `CLIENT_SECRET_BASIC`, and `CLIENT_SECRET_POST`."),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{`NONE`, `CLIENT_SECRET_BASIC`, `CLIENT_SECRET_POST`}, false)),
				Default:          "CLIENT_SECRET_BASIC",
			},
		},
	}
}

func resourceResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource := *management.NewResource(d.Get("name").(string)) // Resource |  (optional)
	resource.SetType(management.ENUMRESOURCETYPE_CUSTOM)

	if v, ok := d.GetOk("description"); ok {
		resource.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("audience"); ok {
		resource.SetAudience(v.(string))
	}

	if v, ok := d.GetOk("access_token_validity_seconds"); ok {
		resource.SetAccessTokenValiditySeconds(int32(v.(int)))
	}

	if v, ok := d.GetOk("introspect_endpoint_auth_method"); ok {
		resource.SetIntrospectEndpointAuthMethod(v.(string))
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourcesApi.CreateResource(ctx, d.Get("environment_id").(string)).Resource(resource).Execute()
		},
		"CreateResource",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Resource)

	d.SetId(respObject.GetId())

	return resourceResourceRead(ctx, d, meta)
}

func resourceResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourcesApi.ReadOneResource(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneResource",
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

	respObject := resp.(*management.Resource)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := respObject.GetTypeOk(); ok {
		d.Set("type", string(*v))
	} else {
		d.Set("type", nil)
	}

	if v, ok := respObject.GetAudienceOk(); ok {
		d.Set("audience", v)
	} else {
		d.Set("audience", nil)
	}

	if v, ok := respObject.GetAccessTokenValiditySecondsOk(); ok {
		d.Set("access_token_validity_seconds", v)
	} else {
		d.Set("access_token_validity_seconds", nil)
	}

	if v, ok := respObject.GetIntrospectEndpointAuthMethodOk(); ok {
		d.Set("introspect_endpoint_auth_method", v)
	} else {
		d.Set("introspect_endpoint_auth_method", nil)
	}

	return diags
}

func resourceResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resource := *management.NewResource(d.Get("name").(string)) // Resource |  (optional)

	if v, ok := d.GetOk("type"); ok {
		resource.SetType(management.EnumResourceType(v.(string)))
	}

	if v, ok := d.GetOk("description"); ok {
		resource.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("audience"); ok {
		resource.SetAudience(v.(string))
	} else {
		resource.SetAudience(d.Get("name").(string))
	}

	if v, ok := d.GetOk("access_token_validity_seconds"); ok {
		resource.SetAccessTokenValiditySeconds(int32(v.(int)))
	}

	if v, ok := d.GetOk("introspect_endpoint_auth_method"); ok {
		resource.SetIntrospectEndpointAuthMethod(v.(string))
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourcesApi.UpdateResource(ctx, d.Get("environment_id").(string), d.Id()).Resource(resource).Execute()
		},
		"UpdateResource",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceResourceRead(ctx, d, meta)
}

func resourceResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ResourcesApi.DeleteResource(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteResource",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceResourceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/resourceID\"", d.Id())
	}

	environmentID, resourceID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(resourceID)

	resourceResourceRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
