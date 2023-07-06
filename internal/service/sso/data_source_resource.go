package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceResource() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne OAuth 2.0 resource data",

		ReadContext: datasourcePingOneResourceRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"resource_id": {
				Description:      "The ID of the resource.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the resource.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_id"},
			},
			"description": {
				Description: "A description of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "A string that specifies the type of resource. Options are `OPENID_CONNECT`, `PINGONE_API`, and `CUSTOM`. Only the `CUSTOM` resource type can be created. `OPENID_CONNECT` specifies the built-in platform resource for OpenID Connect. `PINGONE_API` specifies the built-in platform resource for PingOne.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"audience": {
				Description: "A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"access_token_validity_seconds": {
				Description: "An integer that specifies the number of seconds that the access token is valid.  The minimum value is 300 seconds (5 minutes); the maximum value is 2592000 seconds (30 days).",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"introspect_endpoint_auth_method": {
				Description: fmt.Sprintf("The client authentication methods supported by the token endpoint. Options are `%s`, `%s`, and `%s`.", string(management.ENUMRESOURCEINTROSPECTENDPOINTAUTHMETHOD_NONE), string(management.ENUMRESOURCEINTROSPECTENDPOINTAUTHMETHOD_CLIENT_SECRET_BASIC), string(management.ENUMRESOURCEINTROSPECTENDPOINTAUTHMETHOD_CLIENT_SECRET_POST)),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"client_secret": {
				Description: "An auto-generated resource client secret. Possible characters are `a-z`, `A-Z`, `0-9`, `-`, `.`, `_`, `~`. The secret has a minimum length of 64 characters per SHA-512 requirements when using the HS512 algorithm to sign ID tokens using the secret as the key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func datasourcePingOneResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp *management.Resource

	if v, ok := d.GetOk("name"); ok {

		resp, diags = fetchResourceFromName(ctx, apiClient, d.Get("environment_id").(string), v.(string))
		if diags.HasError() {
			return diags
		}

	} else if v, ok2 := d.GetOk("resource_id"); ok2 {

		resourceResp, diags := sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return apiClient.ResourcesApi.ReadOneResource(ctx, d.Get("environment_id").(string), v.(string)).Execute()
			},
			"ReadOneResource",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		resp = resourceResp.(*management.Resource)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither resource_id or name are set",
			Detail:   "Neither resource_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("resource_id", resp.GetId())
	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := resp.GetTypeOk(); ok {
		d.Set("type", string(*v))

		if *v == management.ENUMRESOURCETYPE_CUSTOM {
			respSecret, diags := sdk.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					return apiClient.ResourceClientSecretApi.ReadResourceSecret(ctx, d.Get("environment_id").(string), d.Id()).Execute()
				},
				"ReadResourceSecret",
				sdk.CustomErrorResourceNotFoundWarning,
				func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

					// The secret may take a short time to propagate
					if r.StatusCode == 404 {
						tflog.Warn(ctx, "Resource secret not found, available for retry")
						return true
					}

					if p1error != nil {
						var err error

						// Permissions may not have propagated by this point
						if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
							tflog.Warn(ctx, "Insufficient PingOne privileges detected")
							return true
						}
						if err != nil {
							tflog.Warn(ctx, "Cannot match error string for retry")
							return false
						}

					}

					return false
				},
			)
			if diags.HasError() {
				return diags
			}

			respSecretObj := *respSecret.(*management.ResourceSecret)

			if v, ok := respSecretObj.GetSecretOk(); ok {
				d.Set("client_secret", v)
			} else {
				d.Set("client_secret", nil)
			}
		} else {
			d.Set("client_secret", nil)
		}
	} else {
		d.Set("type", nil)
		d.Set("client_secret", nil)
	}

	if v, ok := resp.GetAudienceOk(); ok {
		d.Set("audience", v)
	} else {
		d.Set("audience", nil)
	}

	if v, ok := resp.GetAccessTokenValiditySecondsOk(); ok {
		d.Set("access_token_validity_seconds", v)
	} else {
		d.Set("access_token_validity_seconds", nil)
	}

	if v, ok := resp.GetIntrospectEndpointAuthMethodOk(); ok {
		d.Set("introspect_endpoint_auth_method", string(*v))
	} else {
		d.Set("introspect_endpoint_auth_method", nil)
	}

	return diags
}

func fetchResourceFromName(ctx context.Context, apiClient *management.APIClient, environmentID, resourceName string) (*management.Resource, diag.Diagnostics) {

	var resp *management.Resource

	respList, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourcesApi.ReadAllResources(ctx, environmentID).Execute()
		},
		"ReadAllResources",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	if resources, ok := respList.(*management.EntityArray).Embedded.GetResourcesOk(); ok {

		found := false
		for _, resource := range resources {

			resource := resource // fix for exportloopref lint

			if resource.GetName() == resourceName {
				resp = &resource
				found = true
				break
			}
		}

		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot find resource %s", resourceName),
			})

			return nil, diags
		}

	}

	return resp, diags
}
