package sso

import (
	"context"
	"fmt"
	"net/http"

	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceResourceAttribute() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne resource attribute data",

		ReadContext: datasourcePingOneResourceAttributeRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_id": {
				Description:      "The ID of the resource that the resource attribute is assigned to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_attribute_id": {
				Description:      "The ID of the resource attribute.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the resource attribute.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_attribute_id"},
			},
			"type": {
				Description: "A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "A string that specifies the value of the custom resource attribute.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id_token_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available in the ID Token.  Only applies to resources that are of type `OPENID_CONNECT`.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"userinfo_enabled": {
				Description: "A boolean that specifies whether the attribute mapping should be available through the /as/userinfo endpoint.  Only applies to resources that are of type `OPENID_CONNECT`.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneResourceAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp *management.ResourceAttribute

	if v, ok := d.GetOk("name"); ok {

		resp, diags = fetchResourceAttributeFromName(ctx, apiClient, d.Get("environment_id").(string), d.Get("resource_id").(string), v.(string))
		if diags.HasError() {
			return diags
		}

	} else if v, ok2 := d.GetOk("resource_attribute_id"); ok2 {

		resourceAttrResp, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ResourceAttributesApi.ReadOneResourceAttribute(ctx, d.Get("environment_id").(string), d.Get("resource_id").(string), v.(string)).Execute()
			},
			"ReadOneResourceAttribute",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		resp = resourceAttrResp.(*management.ResourceAttribute)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither resource_attribute_id or name are set",
			Detail:   "Neither resource_attribute_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("resource_attribute_id", resp.GetId())
	d.Set("name", resp.GetName())
	d.Set("value", resp.GetValue())
	d.Set("type", resp.GetType())

	if v, ok := resp.GetIdTokenOk(); ok {
		d.Set("id_token_enabled", v)
	}

	if v, ok := resp.GetUserInfoOk(); ok {
		d.Set("userinfo_enabled", v)
	}

	return diags
}

// Replace with fetchResourceAttributeFromName_Framework when migrating to plugin framework
func fetchResourceAttributeFromName(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceAttributeName string) (*management.ResourceAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resp *management.ResourceAttribute

	respList, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.ReadAllResourceAttributes(ctx, environmentID, resourceID).Execute()
		},
		"ReadAllResourceAttributes",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	if resourceAttributes, ok := respList.(*management.EntityArray).Embedded.GetAttributesOk(); ok {

		found := false
		for _, resourceAttribute := range resourceAttributes {

			if resourceAttribute.ResourceAttribute.GetName() == resourceAttributeName {
				resp = resourceAttribute.ResourceAttribute
				found = true
				break
			}
		}

		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot find resource attribute %s", resourceAttributeName),
			})

			return nil, diags
		}

	}

	return resp, diags
}

func fetchResourceAttributeFromName_Framework(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceAttributeName string) (*management.ResourceAttribute, frameworkdiag.Diagnostics) {
	var diags frameworkdiag.Diagnostics

	var resp *management.ResourceAttribute

	respList, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ResourceAttributesApi.ReadAllResourceAttributes(ctx, environmentID, resourceID).Execute()
		},
		"ReadAllResourceAttributes",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(diags...)
	if diags.HasError() {
		return nil, diags
	}

	if resourceAttributes, ok := respList.(*management.EntityArray).Embedded.GetAttributesOk(); ok {

		found := false
		for _, resourceAttribute := range resourceAttributes {

			if resourceAttribute.ResourceAttribute.GetName() == resourceAttributeName {
				resp = resourceAttribute.ResourceAttribute
				found = true
				break
			}
		}

		if !found {
			diags.AddError(
				fmt.Sprintf("Cannot find resource attribute %s", resourceAttributeName),
				"The resource attribute cannot be found by the provided name.",
			)

			return nil, diags
		}

	}

	return resp, diags
}
