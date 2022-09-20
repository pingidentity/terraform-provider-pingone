package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceApplication() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne applications",

		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "A string that specifies the name of the application.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A string that specifies the description of the application.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "A boolean that specifies whether the application is enabled in the environment.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tags": {
				Description: "An array that specifies the list of labels associated with the application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"PING_FED_CONNECTION_INTEGRATION"}, false),
				},
			},
			"login_page_url": {
				Description:      "A string that specifies the custom login page URL for the application. If you set the `login_page_url` property for applications in an environment that sets a custom domain, the URL should include the top-level domain and at least one additional domain level. **Warning** To avoid issues with third-party cookies in some browsers, a custom domain must be used, giving your PingOne environment the same parent domain as your authentication application. For more information about custom domains, see Custom domains.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
			},
			"icon": {
				Description: "The HREF and the ID for the application icon.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:      "The ID for the application icon.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"href": {
							Description:      "The HREF for the application icon.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
			"access_control_role_type": {
				Description:  "A string that specifies the user role required to access the application. Options are `ADMIN_USERS_ONLY`. A user is an admin user if the user has one or more of the following roles Organization Admin, Environment Admin, Identity Data Admin, or Client Application Developer.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADMIN_USERS_ONLY"}, false),
			},
			"access_control_group_options": {
				Description: "Group access control settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "A string that specifies the group type required to access the application. Options are `ANY_GROUP` (the actor must belong to at least one group listed in the `groups` property) and `ALL_GROUPS` (the actor must belong to all groups listed in the `groups` property).",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ANY_GROUP", "ALL_GROUPS"}, false),
						},
						"groups": {
							Description: "A set that specifies the group IDs for the groups the actor must belong to for access to the application.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
			},
			"oidc_options": {
				Description:  "OIDC/OAuth application specific settings.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"oidc_options", "saml_options"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "A string that specifies the type associated with the application.",
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"WEB_APP", "NATIVE_APP", "SINGLE_PAGE_APP", "WORKER", "CUSTOM"}, false),
						},
						"home_page_url": {
							Description:      "A string that specifies the custom home page URL for the application.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"grant_types": {
							Description: "A list that specifies the grant type for the authorization request. Options are `AUTHORIZATION_CODE`, `IMPLICIT`, `REFRESH_TOKEN`, `CLIENT_CREDENTIALS`.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"AUTHORIZATION_CODE", "IMPLICIT", "REFRESH_TOKEN", "CLIENT_CREDENTIALS"}, false),
							},
							Required: true,
						},
						"response_types": {
							Description: "A list that specifies the code or token type returned by an authorization request. Options are `TOKEN`, `ID_TOKEN`, and `CODE`. Note that `CODE` cannot be used in an authorization request with `TOKEN` or `ID_TOKEN` because PingOne does not currently support OIDC hybrid flows.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"TOKEN", "ID_TOKEN", "CODE"}, false),
							},
							Optional: true,
							Computed: true,
						},
						"token_endpoint_authn_method": {
							Description:  "A string that specifies the client authentication methods supported by the token endpoint.  Options are `NONE`, `CLIENT_SECRET_BASIC`, `CLIENT_SECRET_POST`",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "CLIENT_SECRET_BASIC", "CLIENT_SECRET_POST"}, false),
						},
						"pkce_enforcement": {
							Description:  "A string that specifies how `PKCE` request parameters are handled on the authorize request. Options are `OPTIONAL` PKCE `code_challenge` is optional and any code challenge method is acceptable. `REQUIRED` PKCE `code_challenge` is required and any code challenge method is acceptable. `S256_REQUIRED` PKCE `code_challenge` is required and the `code_challenge_method` must be `S256`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "OPTIONAL",
							ValidateFunc: validation.StringInSlice([]string{"OPTIONAL", "REQUIRED", "S256_REQUIRED"}, false),
						},
						"redirect_uris": {
							Description: "A string that specifies the callback URI for the authentication response.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
							},
							Optional: true,
						},
						"post_logout_redirect_uris": {
							Description: "A string that specifies the URLs that the browser can be redirected to after logout.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
							},
							Optional: true,
						},
						"refresh_token_duration": {
							Description:  "An integer that specifies the lifetime in seconds of the refresh token. If a value is not provided, the default value is 2592000, or 30 days. Valid values are between 60 and 2147483647. If the refresh_token_rolling_duration property is specified for the application, then this property must be less than or equal to the value of refreshTokenRollingDuration. After this property is set, the value cannot be nullified. This value is used to generate the value for the exp claim when minting a new refresh token.",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2592000,
							ValidateFunc: validation.IntBetween(60, 2147483647),
						},
						"refresh_token_rolling_duration": {
							Description:  "An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required. If a value is not provided, the default value is 15552000, or 180 days. Valid values are between 60 and 2147483647. After this property is set, the value cannot be nullified. This value is used to generate the value for the exp claim when minting a new refresh token.",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      15552000,
							ValidateFunc: validation.IntBetween(60, 2147483647),
						},
						"client_id": {
							Description: "A string that specifies the application ID used to authenticate to the authorization server.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"client_secret": {
							Description: "A string that specifies the application secret ID used to authenticate to the authorization server.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"support_unsigned_request_object": {
							Description: "A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If false or null (default), an unsigned request object is not allowed.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"mobile_app": {
							Description: "Mobile application integration settings for `NATIVE_APP` type applications.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bundle_id": {
										Description:      "A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the bundle_id property is unique per environment, and once defined, is immutable.  this setting overrides the top-level `bundle_id` field",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"package_name": {
										Description:      "A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.  this setting overrides the top-level `package_name` field.",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"passcode_refresh_seconds": {
										Description:  "The amount of time a passcode should be displayed before being replaced with a new passcode - must be between 30 and 60.",
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IntBetween(30, 60),
									},
									"integrity_detection": {
										Description: "Mobile application integrity detection settings.",
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Description: "A boolean that specifies whether device integrity detection takes place on mobile devices.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"cache_duration": {
													Description: "Every attestation request entails a certain time tradeoff. You can choose to cache successful integrity detection calls for a predefined duration, between a minimum of 1 minute and a maximum of 48 hours. If integrity detection is ENABLED, the cache duration must be set.",
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"amount": {
																Description: "An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.",
																Type:        schema.TypeInt,
																Required:    true,
															},
															"units": {
																Description:  "A string that specifies the time units of the mobile.integrityDetection.cacheDuration.amount.  Possible values are `MINUTES`, `HOURS` and defaults to `MINUTES`",
																Type:         schema.TypeString,
																Optional:     true,
																Default:      "MINUTES",
																ValidateFunc: validation.StringInSlice([]string{"MINUTES", "HOURS"}, false),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"bundle_id": {
							Description:      "A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable; any change will force recreation of the applicationr resource.",
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"package_name": {
							Description:      "A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.",
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
				},
			},
			"saml_options": {
				Description:  "SAML application specific settings.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"oidc_options", "saml_options"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "A string that specifies the type associated with the application.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "WEB_APP",
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"WEB_APP", "CUSTOM"}, false),
						},
						"acs_urls": {
							Description: "A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"assertion_duration": {
							Description: "An integer that specifies the assertion validity duration in seconds.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"assertion_signed_enabled": {
							Description: "A boolean that specifies whether the SAML assertion itself should be signed.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"idp_signing_key_id": {
							Description:      "An ID for the certificate key pair to be used by the identity provider to sign assertions and responses. If this property is omitted, the default signing certificate for the environment is used.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"nameid_format": {
							Description: "A string that specifies the format of the Subject NameID attibute in the SAML assertion.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"response_is_signed": {
							Description: "A boolean that specifies whether the SAML assertion response itself should be signed.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"slo_binding": {
							Description:  "A string that specifies the binding protocol to be used for the logout response. Options are `HTTP_REDIRECT` or `HTTP_POST`. The default is `HTTP_POST`; existing configurations with no data default to `HTTP_POST`.",
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"HTTP_REDIRECT", "HTTP_POST"}, false),
							Optional:     true,
							Default:      "HTTP_POST",
						},
						"slo_endpoint": {
							Description: "A string that specifies the logout endpoint URL. This is an optional property. However, if a sloEndpoint logout endpoint URL is not defined, logout actions result in an error.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"slo_response_endpoint": {
							Description: "A string that specifies the endpoint URL to submit the logout response. If a value is not provided, the sloEndpoint property value is used to submit SLO response.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"sp_entity_id": {
							Description: "A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sp_verification_certificate_ids": {
							Description: "A list that specifies the certificate IDs used to verify the service provider signature.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
							},
						},
					},
				},
			},
		},
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	applicationRequest := &management.CreateApplicationRequest{}

	if _, ok := d.GetOk("oidc_options"); ok {
		var application *management.ApplicationOIDC
		application, diags = expandApplicationOIDC(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationOIDC = application
	}

	if _, ok := d.GetOk("saml_options"); ok {
		var application *management.ApplicationSAML
		application, diags = expandApplicationSAML(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationSAML = application
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApi.CreateApplication(ctx, d.Get("environment_id").(string)).CreateApplicationRequest(*applicationRequest).Execute()
		},
		"CreateApplication",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CreateApplication201Response)

	if respObject.ApplicationOIDC != nil && respObject.ApplicationOIDC.GetId() != "" {
		d.SetId(respObject.ApplicationOIDC.GetId())
	} else if respObject.ApplicationSAML != nil && respObject.ApplicationSAML.GetId() != "" {
		d.SetId(respObject.ApplicationSAML.GetId())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Cannot determine application ID from API response for application: %s", d.Get("name")),
			Detail:   fmt.Sprintf("Full response object: %v\n", resp),
		})

		return diags
	}

	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApi.ReadOneApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneApplication",
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

	respObject := resp.(*management.ReadOneApplication200Response)

	if respObject.ApplicationOIDC != nil && respObject.ApplicationOIDC.GetId() != "" {

		respSecret, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.ApplicationSecretApi.ReadApplicationSecret(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			},
			"ReadApplicationSecret",
			sdk.CustomErrorResourceNotFoundWarning,
			func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

				// The secret may take a short time to propagate
				if r.StatusCode == 404 {
					tflog.Warn(ctx, "Application secret not found, available for retry")
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

		application := respObject.ApplicationOIDC

		d.Set("name", application.GetName())
		d.Set("enabled", application.GetEnabled())

		if v, ok := application.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		if v, ok := application.GetTagsOk(); ok {
			d.Set("tags", v)
		} else {
			d.Set("tags", nil)
		}

		if v, ok := application.GetLoginPageUrlOk(); ok {
			d.Set("login_page_url", v)
		} else {
			d.Set("login_page_url", nil)
		}

		if v, ok := application.GetIconOk(); ok {
			d.Set("icon", flattenIcon(v))
		} else {
			d.Set("icon", nil)
		}

		if v, ok := application.GetAccessControlOk(); ok {

			if j, ok := v.Role.GetTypeOk(); ok {
				d.Set("access_control_role_type", string(*j))
			} else {
				d.Set("access_control_role_type", nil)
			}

			if j, ok := v.GetGroupOk(); ok {

				groups := make([]string, 0)
				for _, k := range j.GetGroups() {
					groups = append(groups, k.GetId())
				}

				groupObj := map[string]interface{}{
					"type":   j.GetType(),
					"groups": groups,
				}

				groupsObj := make([]interface{}, 0)

				d.Set("access_control_group_options", append(groupsObj, groupObj))
			} else {
				d.Set("access_control_group_options", nil)
			}
		} else {
			d.Set("access_control_role_type", nil)
			d.Set("access_control_group_options", nil)
		}

		v, diags := flattenOIDCOptions(application, respSecret.(*management.ApplicationSecret))
		if diags.HasError() {
			return diags
		}

		d.Set("oidc_options", v)

	} else if respObject.ApplicationSAML != nil && respObject.ApplicationSAML.GetId() != "" {

		application := respObject.ApplicationSAML

		d.Set("name", application.GetName())
		d.Set("enabled", application.GetEnabled())

		if v, ok := application.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		if v, ok := application.GetTagsOk(); ok {
			d.Set("tags", v)
		} else {
			d.Set("tags", nil)
		}

		if v, ok := application.GetLoginPageUrlOk(); ok {
			d.Set("login_page_url", v)
		} else {
			d.Set("login_page_url", nil)
		}

		if v, ok := application.GetIconOk(); ok {
			d.Set("icon", flattenIcon(v))
		} else {
			d.Set("icon", nil)
		}

		if v, ok := application.GetAccessControlOk(); ok {

			if j, ok := v.Role.GetTypeOk(); ok {
				d.Set("access_control_role_type", string(*j))
			} else {
				d.Set("access_control_role_type", nil)
			}

			if j, ok := v.GetGroupOk(); ok {

				groups := make([]string, 0)
				for _, k := range j.GetGroups() {
					groups = append(groups, k.GetId())
				}

				groupObj := map[string]interface{}{
					"type":   j.GetType(),
					"groups": groups,
				}

				groupsObj := make([]interface{}, 0)

				d.Set("access_control_group_options", append(groupsObj, groupObj))
			} else {
				d.Set("access_control_group_options", nil)
			}
		} else {
			d.Set("access_control_role_type", nil)
			d.Set("access_control_group_options", nil)
		}

		d.Set("saml_options", flattenSAMLOptions(application))

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Cannot determine application ID from API response for application: %s", d.Get("name")),
			Detail:   fmt.Sprintf("Full response object: %v\n", resp),
		})

		return diags
	}

	return diags
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	applicationRequest := &management.UpdateApplicationRequest{}

	if _, ok := d.GetOk("oidc_options"); ok {
		var application *management.ApplicationOIDC
		application, diags = expandApplicationOIDC(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationOIDC = application
	}

	if _, ok := d.GetOk("saml_options"); ok {
		var application *management.ApplicationSAML
		application, diags = expandApplicationSAML(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationSAML = application
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApi.UpdateApplication(ctx, d.Get("environment_id").(string), d.Id()).UpdateApplicationRequest(*applicationRequest).Execute()
		},
		"UpdateApplication",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ApplicationsApi.DeleteApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteApplication",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceApplicationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID\"", d.Id())
	}

	environmentID, applicationID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(applicationID)

	resourceApplicationRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

// OIDC
func expandApplicationOIDC(d *schema.ResourceData) (*management.ApplicationOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	var application management.ApplicationOIDC

	if v, ok := d.Get("oidc_options").([]interface{}); ok && len(v) > 0 && v[0] != nil {

		oidcOptions := v[0].(map[string]interface{})

		var applicationType *management.EnumApplicationType
		applicationType, err := expandApplicationType(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot determine application `type`: %v", err),
			})

			return nil, diags
		}

		grantTypes, _ := expandGrantTypes(oidcOptions["grant_types"].(*schema.Set))

		// Set the object
		application = *management.NewApplicationOIDC(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, *applicationType, grantTypes, management.EnumApplicationOIDCTokenAuthMethod(oidcOptions["token_endpoint_authn_method"].(string)))

		// set the common optional options
		applicationCommon := expandCommonOptionalAttributes(d)

		if v1, ok := applicationCommon.GetDescriptionOk(); ok {
			application.SetDescription(*v1)
		}

		if v1, ok := applicationCommon.GetTagsOk(); ok {
			application.SetTags(v1)
		}

		if v1, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v1)
		}

		if v1, ok := applicationCommon.GetAssignActorRolesOk(); ok {
			application.SetAssignActorRoles(*v1)
		}

		if v1, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v1)
		}

		if v1, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v1)
		}

		// Set the OIDC specific optional options

		if v1, ok := oidcOptions["home_page_url"].(string); ok && v1 != "" {
			application.SetHomePageUrl(v1)
		}

		if v1, ok := oidcOptions["response_types"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]management.EnumApplicationOIDCResponseType, 0)
			for _, j := range v1.List() {
				obj = append(obj, management.EnumApplicationOIDCResponseType(j.(string)))
			}
			application.SetResponseTypes(obj)
		}

		if v1, ok := oidcOptions["pkce_enforcement"].(string); ok && v1 != "" {
			if application.GetType() == management.ENUMAPPLICATIONTYPE_WEB_APP || application.GetType() == management.ENUMAPPLICATIONTYPE_SINGLE_PAGE_APP || application.GetType() == management.ENUMAPPLICATIONTYPE_CUSTOM_APP {
				application.SetPkceEnforcement(management.EnumApplicationOIDCPKCEOption(v1))
			}
		}

		if v1, ok := oidcOptions["redirect_uris"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]string, 0)
			for _, j := range v1.List() {
				obj = append(obj, j.(string))
			}
			application.SetRedirectUris(obj)
		}

		if v1, ok := oidcOptions["post_logout_redirect_uris"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]string, 0)
			for _, j := range v1.List() {
				obj = append(obj, j.(string))
			}
			application.SetPostLogoutRedirectUris(obj)
		}

		if v1, ok := oidcOptions["refresh_token_duration"].(int); ok {
			//if refreshTokenEnabled {
			application.SetRefreshTokenDuration(int32(v1))
			//} else {
			//	diags = append(diags, diag.Diagnostic{
			//		Severity: diag.Warning,
			//		Summary:  fmt.Sprintf("`refresh_token_duration` has no effect when the %s grant type is not set", management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN),
			//	})
			//}
		}

		if v1, ok := oidcOptions["refresh_token_rolling_duration"].(int); ok {
			//if refreshTokenEnabled {
			application.SetRefreshTokenRollingDuration(int32(v1))
			//} else {
			//	diags = append(diags, diag.Diagnostic{
			//		Severity: diag.Warning,
			//		Summary:  fmt.Sprintf("`refresh_token_rolling_duration` has no effect when the %s grant type is not set", management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN),
			//	})
			//}
		}

		if v1, ok := oidcOptions["support_unsigned_request_object"].(bool); ok {
			application.SetSupportUnsignedRequestObject(v1)
		}

		if v1, ok := oidcOptions["mobile_app"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			var mobile *management.ApplicationOIDCAllOfMobile
			mobile, diags = expandMobile(v1[0].(map[string]interface{}))
			if diags.HasError() {
				return nil, diags
			}
			application.SetMobile(*mobile)
		}

		if v1, ok := oidcOptions["bundle_id"].(string); ok && v1 != "" {
			application.SetBundleId(v1)
		}

		if v1, ok := oidcOptions["package_name"].(string); ok && v1 != "" {
			application.SetPackageName(v1)
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("OIDC options not available for application: %s", d.Get("name")),
		})

		return nil, diags
	}

	return &application, diags
}

func expandGrantTypes(s *schema.Set) ([]management.EnumApplicationOIDCGrantType, bool) {
	grantTypes := make([]management.EnumApplicationOIDCGrantType, 0)
	refreshToken := false
	for _, j := range s.List() {
		grantTypes = append(grantTypes, management.EnumApplicationOIDCGrantType(j.(string)))
		if management.EnumApplicationOIDCGrantType(j.(string)) == management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN {
			refreshToken = true
		}
	}

	return grantTypes, refreshToken
}

func expandMobile(s map[string]interface{}) (*management.ApplicationOIDCAllOfMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	mobile := management.NewApplicationOIDCAllOfMobile()
	if v, ok := s["bundle_id"].(string); ok && v != "" {
		mobile.SetBundleId(v)
	}

	if v, ok := s["package_name"].(string); ok && v != "" {
		mobile.SetPackageName(v)
	}

	if v, ok := s["passcode_refresh_seconds"].(int); ok {
		mobile.SetPasscodeRefreshDuration(*management.NewApplicationOIDCAllOfMobilePasscodeRefreshDuration(int32(v), management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS))
	}

	if v, ok := s["integrity_detection"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {

		obj := v[0].(map[string]interface{})

		integrityDetection := *management.NewApplicationOIDCAllOfMobileIntegrityDetection()

		if j, okJ := obj["enabled"].(bool); okJ {
			var mode management.EnumEnabledStatus
			if j {
				mode = management.ENUMENABLEDSTATUS_ENABLED
			} else {
				mode = management.ENUMENABLEDSTATUS_DISABLED
			}
			integrityDetection.SetMode(mode)
		}

		if j, okJ := obj["cache_duration"].([]interface{}); okJ && len(j) > 0 && j[0] != nil {
			integrityDetection.SetCacheDuration(expandMobileIntegrityCacheDuration(j[0].(map[string]interface{})))
		} else {
			if integrityDetection.GetMode() == management.ENUMENABLEDSTATUS_ENABLED {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Attribute `cache_duration` is required when the mobile integrity check is enabled in the application",
				})

				return nil, diags
			}
		}

		mobile.SetIntegrityDetection(integrityDetection)
	}

	return mobile, diags
}

func expandMobileIntegrityCacheDuration(s map[string]interface{}) management.ApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration {

	obj := *management.NewApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration()
	obj.SetAmount(int32(s["amount"].(int)))
	obj.SetUnits(management.EnumDurationUnitMinsHours(s["units"].(string)))

	return obj

}

// SAML
func expandApplicationSAML(d *schema.ResourceData) (*management.ApplicationSAML, diag.Diagnostics) {
	var diags diag.Diagnostics

	var application management.ApplicationSAML

	if v, ok := d.Get("saml_options").([]interface{}); ok && len(v) > 0 && v[0] != nil {

		samlOptions := v[0].(map[string]interface{})

		var applicationType *management.EnumApplicationType
		applicationType, err := expandApplicationType(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot determine application `type`: %v", err),
			})

			return nil, diags
		}

		// Set the object
		acsUrls := make([]string, 0)
		for _, v := range samlOptions["acs_urls"].(*schema.Set).List() {
			acsUrls = append(acsUrls, v.(string))
		}
		application = *management.NewApplicationSAML(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMAPPLICATIONPROTOCOL_SAML, *applicationType, acsUrls, int32(samlOptions["assertion_duration"].(int)), samlOptions["sp_entity_id"].(string))

		// set the common optional options
		applicationCommon := expandCommonOptionalAttributes(d)

		if v1, ok := applicationCommon.GetDescriptionOk(); ok {
			application.SetDescription(*v1)
		}

		if v1, ok := applicationCommon.GetTagsOk(); ok {
			application.SetTags(v1)
		}

		if v1, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v1)
		}

		if v1, ok := applicationCommon.GetAssignActorRolesOk(); ok {
			application.SetAssignActorRoles(*v1)
		}

		if v1, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v1)
		}

		if v1, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v1)
		}

		// Set the SAML specific optional options

		if v1, ok := samlOptions["assertion_signed_enabled"].(bool); ok {
			application.SetAssertionSigned(v1)
		}

		if v1, ok := samlOptions["idp_signing_key_id"].(string); ok && v1 != "" {
			application.SetIdpSigning(*management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(v1)))
		}

		if v1, ok := samlOptions["nameid_format"].(string); ok && v1 != "" {
			application.SetNameIdFormat(v1)
		}

		if v1, ok := samlOptions["response_is_signed"].(bool); ok {
			application.SetResponseSigned(v1)
		}

		if v1, ok := samlOptions["slo_binding"].(string); ok && v1 != "" {
			application.SetSloBinding(management.EnumApplicationSAMLSloBinding(v1))
		}

		if v1, ok := samlOptions["slo_endpoint"].(string); ok && v1 != "" {
			application.SetSloEndpoint(v1)
		}

		if v1, ok := samlOptions["slo_endpoint"].(string); ok && v1 != "" {
			application.SetSloEndpoint(v1)
		}

		if v1, ok := samlOptions["slo_response_endpoint"].(string); ok && v1 != "" {
			application.SetSloResponseEndpoint(v1)
		}

		if v1, ok := samlOptions["sp_verification_certificate_ids"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			certificates := make([]management.ApplicationSAMLAllOfSpVerificationCertificates, 0)
			for _, j := range v1.List() {
				certificate := *management.NewApplicationSAMLAllOfSpVerificationCertificates(j.(string))
				certificates = append(certificates, certificate)
			}

			application.SetSpVerification(*management.NewApplicationSAMLAllOfSpVerification(certificates))
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("SAML options not available for application: %s", d.Get("name")),
		})

		return nil, diags
	}

	return &application, diags
}

// Common

func expandCommonOptionalAttributes(d *schema.ResourceData) management.Application {

	application := management.Application{}

	if v, ok := d.GetOk("description"); ok {
		application.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		if j, okJ := v.([]interface{}); okJ {
			tags := make([]management.EnumApplicationTags, 0)
			for _, k := range j {
				tags = append(tags, management.EnumApplicationTags(k.(string)))
			}

			application.Tags = tags
		}
	}

	if v, ok := d.GetOk("login_page_url"); ok {
		if v != "" {
			application.SetLoginPageUrl(v.(string))
		}
	}

	application.SetAssignActorRoles(false)

	if v, ok := d.GetOk("icon"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			application.SetIcon(*management.NewApplicationIcon(attrs["id"].(string), attrs["href"].(string)))
		}
	}

	accessControl := *management.NewApplicationAccessControl()
	accessControlCount := 0

	if v, ok := d.GetOk("access_control_role_type"); ok {
		accessControl.SetRole(*management.NewApplicationAccessControlRole(management.EnumApplicationAccessControlType(v.(string))))
		accessControlCount += 1
	}

	if v, ok := d.GetOk("access_control_group_options"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			obj := j[0].(map[string]interface{})

			groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)
			for _, j := range obj["groups"].(*schema.Set).List() {
				groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(j.(string)))
			}

			accessControl.SetGroup(*management.NewApplicationAccessControlGroup(obj["type"].(string), groups))

			accessControlCount += 1
		}
	}

	if accessControlCount > 0 {
		application.SetAccessControl(accessControl)
	}

	return application

}

func expandApplicationType(s interface{}) (*management.EnumApplicationType, error) {
	var applicationType management.EnumApplicationType

	if j, ok := s.([]interface{})[0].(map[string]interface{})["type"].(string); ok {
		applicationType = management.EnumApplicationType(j)
	} else {
		return nil, fmt.Errorf("Cannot determine the application type")
	}

	return &applicationType, nil

}

func flattenIcon(s *management.ApplicationIcon) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenOIDCOptions(application *management.ApplicationOIDC, secret *management.ApplicationSecret) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Required
	item := map[string]interface{}{
		"client_id":                   application.GetId(),
		"type":                        application.GetType(),
		"grant_types":                 flattenGrantTypes(application),
		"token_endpoint_authn_method": application.GetTokenEndpointAuthMethod(),
	}

	// Optional
	if v, ok := application.GetHomePageUrlOk(); ok {
		item["home_page_url"] = v
	} else {
		item["home_page_url"] = nil
	}

	if v, ok := application.GetResponseTypesOk(); ok {
		item["response_types"] = v
	} else {
		item["response_types"] = nil
	}

	if v, ok := application.GetPkceEnforcementOk(); ok {
		item["pkce_enforcement"] = v
	} else {
		item["pkce_enforcement"] = nil
	}

	if v, ok := application.GetRedirectUrisOk(); ok {
		item["redirect_uris"] = v
	} else {
		item["redirect_uris"] = nil
	}

	if v, ok := application.GetPostLogoutRedirectUrisOk(); ok {
		item["post_logout_redirect_uris"] = v
	} else {
		item["post_logout_redirect_uris"] = nil
	}

	if v, ok := application.GetRefreshTokenDurationOk(); ok {
		item["refresh_token_duration"] = v
	} else {
		item["refresh_token_duration"] = nil
	}

	if v, ok := application.GetRefreshTokenRollingDurationOk(); ok {
		item["refresh_token_rolling_duration"] = v
	} else {
		item["refresh_token_rolling_duration"] = nil
	}

	if v, ok := secret.GetSecretOk(); ok {
		item["client_secret"] = v
	} else {
		item["client_secret"] = nil
	}

	if v, ok := application.GetSupportUnsignedRequestObjectOk(); ok {
		item["support_unsigned_request_object"] = v
	} else {
		item["support_unsigned_request_object"] = nil
	}

	if v, ok := application.GetMobileOk(); ok {
		j, diags := flattenMobile(v)
		if diags.HasError() {
			return nil, diags
		}
		item["mobile_app"] = j
	} else {
		item["mobile_app"] = nil
	}

	if v, ok := application.GetBundleIdOk(); ok {
		item["bundle_id"] = v
	} else {
		item["bundle_id"] = nil
	}

	if v, ok := application.GetPackageNameOk(); ok {
		item["package_name"] = v
	} else {
		item["package_name"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item), diags

}

func flattenGrantTypes(application *management.ApplicationOIDC) []string {

	grantTypes := application.GetGrantTypes()

	returnGrants := make([]string, 0)
	for _, v := range grantTypes {
		returnGrants = append(returnGrants, string(v))
	}
	return returnGrants
}

func flattenMobile(mobile *management.ApplicationOIDCAllOfMobile) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	item := map[string]interface{}{}

	if v, ok := mobile.GetBundleIdOk(); ok {
		item["bundle_id"] = v
	} else {
		item["bundle_id"] = nil
	}

	if v, ok := mobile.GetPackageNameOk(); ok {
		item["package_name"] = v
	} else {
		item["package_name"] = nil
	}

	if v, ok := mobile.GetPasscodeRefreshDurationOk(); ok {
		item["passcode_refresh_seconds"] = v.GetDuration()
		if j, okJ := v.GetTimeUnitOk(); okJ && *j != management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Expecting time unit of %s for attribute `passcode_refresh_seconds`, got %v", management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS, j),
			})
			return nil, diags
		}
	} else {
		item["passcode_refresh_seconds"] = nil
	}

	if v, ok := mobile.GetIntegrityDetectionOk(); ok {
		item["integrity_detection"] = flattenMobileIntegrityDetection(v)
	} else {
		item["integrity_detection"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item), diags
}

func flattenMobileIntegrityDetection(obj *management.ApplicationOIDCAllOfMobileIntegrityDetection) interface{} {

	item := map[string]interface{}{}

	if v, ok := obj.GetModeOk(); ok {
		if *v == management.ENUMENABLEDSTATUS_ENABLED {
			item["enabled"] = true
		} else {
			item["enabled"] = false
		}
	} else {
		item["enabled"] = nil
	}

	if v, ok := obj.GetCacheDurationOk(); ok {
		cache := map[string]interface{}{
			"amount": v.GetAmount(),
			"units":  v.GetUnits(),
		}

		caches := make([]interface{}, 0)
		item["cache_duration"] = append(caches, cache)
	} else {
		item["cache_duration"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenSAMLOptions(application *management.ApplicationSAML) interface{} {

	// Required
	item := map[string]interface{}{
		"type":               application.GetType(),
		"acs_urls":           application.GetAcsUrls(),
		"assertion_duration": application.GetAssertionDuration(),
		"sp_entity_id":       application.GetSpEntityId(),
	}

	// Optional
	if v, ok := application.GetAssertionSignedOk(); ok {
		item["assertion_signed_enabled"] = v
	} else {
		item["assertion_signed_enabled"] = nil
	}

	var signingKeyID interface{}
	if v, ok := application.GetIdpSigningOk(); ok {
		if j, okJ := v.GetKeyOk(); okJ {
			if k, okK := j.GetIdOk(); okK {
				signingKeyID = k
			}
		}
	}
	item["idp_signing_key_id"] = signingKeyID

	if v, ok := application.GetNameIdFormatOk(); ok {
		item["nameid_format"] = v
	} else {
		item["nameid_format"] = nil
	}

	if v, ok := application.GetResponseSignedOk(); ok {
		item["response_is_signed"] = v
	} else {
		item["response_is_signed"] = nil
	}

	if v, ok := application.GetSloBindingOk(); ok {
		item["slo_binding"] = v
	} else {
		item["slo_binding"] = nil
	}

	if v, ok := application.GetSloEndpointOk(); ok {
		item["slo_endpoint"] = v
	} else {
		item["slo_endpoint"] = nil
	}

	if v, ok := application.GetSloResponseEndpointOk(); ok {
		item["slo_response_endpoint"] = v
	} else {
		item["slo_response_endpoint"] = nil
	}

	if v, ok := application.SpVerification.GetCertificatesOk(); ok {

		idList := make([]string, 0)
		for _, j := range v {
			idList = append(idList, j.GetId())
		}

		item["sp_verification_certificate_ids"] = idList
	} else {
		item["sp_verification_certificate_ids"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)

}
