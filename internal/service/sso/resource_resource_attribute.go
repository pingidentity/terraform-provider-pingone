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
	"golang.org/x/exp/slices"
)

func ResourceResourceAttribute() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource attributes in PingOne.",

		CreateContext: resourceResourceAttributeCreate,
		ReadContext:   resourceResourceAttributeRead,
		UpdateContext: resourceResourceAttributeUpdate,
		DeleteContext: resourceResourceAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the resource attribute in.",
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
				Description: fmt.Sprintf("A string that specifies the name of the custom resource attribute to be included in the access token. When the resource's type property is `OPENID_CONNECT`, the following are reserved names and cannot be used: %s.  When the resource's type property is `OPENID_CONNECT`, using the following names will override the default configured values, rather than creating new attributes: %s.", verify.IllegalOIDCAttributeNameString(), verify.OverrideOIDCAttributeNameString()),
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description:      "A string that specifies the value of the custom resource attribute. This value can be a placeholder that references an attribute in the user schema, expressed as “${user.path.to.value}”, or it can be a static string. Placeholders must be valid, enabled attributes in the environment’s user schema. Examples of valid values are: `${user.email}`, `${user.name.family}`, and `myClaimValueString`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"id_token_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available in the ID Token.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"userinfo_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available through the /as/userinfo endpoint.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false.",
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

	resourceAttribute := *management.NewResourceAttribute(d.Get("name").(string), d.Get("value").(string)) // ResourceAttribute |  (optional)

	attributeID, resourceType, diags := validateAttributeAgainstResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Get("name").(string))
	if diags.HasError() {
		return diags
	}

	if v, ok := d.GetOk("id_token_enabled"); ok {
		if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
			resourceAttribute.SetIdToken(v.(bool))
		}
	}

	if v, ok := d.GetOk("userinfo_enabled"); ok {
		if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
			resourceAttribute.SetUserInfo(v.(bool))
		}
	}

	var resp interface{}

	if attributeID == nil {
		resp, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceAttributesApi.CreateResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string)).ResourceAttribute(resourceAttribute).Execute()
			},
			"CreateResourceAttribute",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	} else {

		resp, diags = sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), *attributeID).ResourceAttribute(resourceAttribute).Execute()
			},
			"UpdateResourceAttribute",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)
	}
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
	}

	if v, ok := respObject.GetUserInfoOk(); ok {
		d.Set("userinfo_enabled", v)
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

	resourceAttribute := *management.NewResourceAttribute(d.Get("name").(string), d.Get("value").(string)) // ResourceAttribute |  (optional)

	_, resourceType, diags := validateAttributeAgainstResourceType(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Get("name").(string))
	if diags.HasError() {
		return diags
	}

	if v, ok := d.GetOk("id_token_enabled"); ok {
		if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
			resourceAttribute.SetIdToken(v.(bool))
		}
	}

	if v, ok := d.GetOk("userinfo_enabled"); ok {
		if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
			resourceAttribute.SetUserInfo(v.(bool))
		}
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.UpdateResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), d.Id()).ResourceAttribute(resourceAttribute).Execute()
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

func validateAttributeAgainstResourceType(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceAttributeName string) (*string, management.EnumResourceType, diag.Diagnostics) {
	var diags diag.Diagnostics

	respObject, diags := fetchResource(ctx, apiClient, environmentID, resourceID)
	if diags.HasError() {
		return nil, management.ENUMRESOURCETYPE_CUSTOM, diags
	}

	if respObject.GetType() == management.ENUMRESOURCETYPE_OPENID_CONNECT {
		if slices.Contains(verify.IllegalOIDCattributeNamesList(), resourceAttributeName) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Invalid attribute name `%s` for the configured OpenID Connect resource.", resourceAttributeName),
				Detail:   fmt.Sprintf("The attribute name provided, `%s`, cannot be used for resource ID `%s`, which is of type `OPENID_CONNECT`.", resourceAttributeName, resourceID),
			})
			return nil, respObject.GetType(), diags
		}

		if slices.Contains(verify.OverrideOIDCAttributeNameList(), resourceAttributeName) {

			resourceAttribute, diags := fetchResourceAttributeFromName(ctx, apiClient, environmentID, resourceID, resourceAttributeName)
			if diags.HasError() {
				return nil, respObject.GetType(), diags
			}

			return resourceAttribute.Id, respObject.GetType(), diags
		}
	}

	return nil, respObject.GetType(), diags
}
