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

func ResourceResourceAttribute() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne OAuth 2.0 resource scope attributes.",

		CreateContext: resourceResourceAttributeCreate,
		ReadContext:   resourceResourceAttributeRead,
		UpdateContext: resourceResourceAttributeUpdate,
		DeleteContext: resourceResourceAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the resource scope attribute in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"resource_id": {
				Description:      "The ID of the resource to assign the resource attribute to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "A string that specifies the name of the custom resource attribute to be included in the access token. The following are reserved names and cannot be used. These reserved names are applicable only when the resource's type property is `OPENID_CONNECT`: `acr`, `amr`, `aud`, `auth_time`, `client_id`, `env`, `exp`, `iat`, `iss`, `jti`, `org`, `p1.*` (any name starting with the p1. prefix), `scope`, `sid`, `sub`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringNotInSlice([]string{"acr", "amr", "aud", "auth_time", "client_id", "env", "exp", "iat", "iss", "jti", "org", "scope", "sid", "sub"}, false)),
			},
			"type": {
				Description: "A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description:      "A string that specifies the value of the custom resource attribute. This value can be a placeholder that references an attribute in the user schema, expressed as “${user.path.to.value}”, or it can be a static string. Placeholders must be valid, enabled attributes in the environment’s user schema. Examples fo valid values are: `${user.email}`, `${user.name.family}`, and `myClaimValueString`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"id_token_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available in the ID Token. This property is applicable only when the application's protocol property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `user_info_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"user_info_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available through the `/as/userinfo` endpoint. This property is applicable only when the application's protocol property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `user_info_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceResourceAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resourceScope := *management.NewResourceAttribute(d.Get("name").(string), d.Get("value").(string)) // ResourceAttribute |  (optional)

	if v, ok := d.GetOk("id_token_enabled"); ok {
		resourceScope.SetIdToken(v.(bool))
	} else {
		resourceScope.SetIdToken(false)
	}

	if v, ok := d.GetOk("user_info_enabled"); ok {
		resourceScope.SetUserInfo(v.(bool))
	} else {
		resourceScope.SetUserInfo(false)
	}

	if !resourceScope.GetIdToken() && !resourceScope.GetUserInfo() {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "`id_token_enabled` and `user_info_enabled` cannot both be false.",
		})

		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.CreateResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).ResourceAttribute(resourceScope).Execute()
		},
		"CreateResourceAttribute",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.ResourceAttribute)

	d.SetId(respObject.GetId())

	return resourceResourceAttributeRead(ctx, d, meta)
}

func resourceResourceAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.ReadOneResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
		},
		"ReadOneResourceAttribute",
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

	respObject := resp.(*management.ResourceAttribute)

	d.Set("name", respObject.GetName())
	d.Set("value", respObject.GetValue())
	d.Set("type", respObject.GetType())

	if v, ok := respObject.GetIdTokenOk(); ok {
		d.Set("id_token_enabled", v)
	} else {
		d.Set("id_token_enabled", nil)
	}

	if v, ok := respObject.GetUserInfoOk(); ok {
		d.Set("user_info_enabled", v)
	} else {
		d.Set("user_info_enabled", nil)
	}

	return diags
}

func resourceResourceAttributeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resourceScope := *management.NewResourceAttribute(d.Get("name").(string), d.Get("value").(string)) // ResourceAttribute |  (optional)

	if v, ok := d.GetOk("id_token_enabled"); ok {
		resourceScope.SetIdToken(v.(bool))
	} else {
		resourceScope.SetIdToken(false)
	}

	if v, ok := d.GetOk("user_info_enabled"); ok {
		resourceScope.SetUserInfo(v.(bool))
	} else {
		resourceScope.SetUserInfo(false)
	}

	if !resourceScope.GetIdToken() && !resourceScope.GetUserInfo() {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "`id_token_enabled` and `user_info_enabled` cannot both be false.",
		})

		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).ResourceAttribute(resourceScope).Execute()
		},
		"UpdateResourceAttribute",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceResourceAttributeRead(ctx, d, meta)
}

func resourceResourceAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ResourceAttributesApi.DeleteResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteResourceAttribute",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceResourceAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/resourceID/resourceAttributeID\"", d.Id())
	}

	environmentID, resourceID, resourceAttributeID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("resource_id", resourceID)
	d.SetId(resourceAttributeID)

	resourceResourceAttributeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
