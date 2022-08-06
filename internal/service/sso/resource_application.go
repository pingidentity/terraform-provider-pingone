package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
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
				Description: "A boolean that specifies whether the environment is enabled in the environment.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tags": {
				Description: "An array that specifies the list of labels associated with the application.",
				Type:        schema.TypeList,
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
			"assign_actor_roles": {
				Description: "A boolean that specifies whether the permissions service should assign default roles to the application. This property is set only on the POST request.  Any update to this attribute will trigger force re-creation of the application.",
				Type:        schema.TypeBool,
				Optional:    false,
				ForceNew:    true,
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
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
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
			"access_control": {
				Description: "Define access control rules for the application.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_type": {
							Description:  "A string that specifies the user role required to access the application. Options are ADMIN_USERS_ONLY. A user is an admin user if the user has one or more of the following roles Organization Admin, Environment Admin, Identity Data Admin, or Client Application Developer.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADMIN_USERS_ONLY"}, false),
						},
						"group": {
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
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.ListOfUniqueStrings),
									},
								},
							},
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
							ValidateFunc: validation.StringInSlice([]string{"WEB_APP", "NATIVE_APP", "SINGLE_PAGE_APP", "WORKER"}, false),
						},
						"home_page_url": {
							Description:      "A string that specifies the custom home page URL for the application.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"grant_types": {
							Description: "A list that specifies the grant type for the authorization request. Options are `AUTHORIZATION_CODE`, `IMPLICIT`, `REFRESH_TOKEN`, `CLIENT_CREDENTIALS`.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"AUTHORIZATION_CODE", "IMPLICIT", "REFRESH_TOKEN", "CLIENT_CREDENTIALS"}, false),
							},
							Required: true,
						},
						"response_types": {
							Description: "A list that specifies the code or token type returned by an authorization request. Options are `TOKEN`, `ID_TOKEN`, and `CODE`. Note that `CODE` cannot be used in an authorization request with `TOKEN` or `ID_TOKEN` because PingOne does not currently support OIDC hybrid flows.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"TOKEN", "ID_TOKEN", "CODE"}, false),
							},
							Optional: true,
						},
						"token_endpoint_authn_method": {
							Description:  "A string that specifies the client authentication methods supported by the token endpoint.  Options are `NONE`, `CLIENT_SECRET_BASIC`, `CLIENT_SECRET_POST`",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "CLIENT_SECRET_BASIC", "CLIENT_SECRET_POST"}, false),
						},
						"pkce_enforcement": {
							Description:  "A string that specifies how `PKCE` request parameters are handled on the authorize request. Options are `OPTIONAL` PKCE `code_challenge` is optional and any code challenge method is acceptable. `REQUIRED` PKCE `code_challenge` is required and any code challenge method is acceptable. `S256_REQUIRED` PKCE `code_challege` is required and the `code_challenge_method` must be `S256`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "OPTIONAL",
							ValidateFunc: validation.StringInSlice([]string{"OPTIONAL", "REQUIRED", "S256_REQUIRED"}, false),
						},
						"redirect_uris": {
							Description: "A string that specifies the callback URI for the authentication response.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
							},
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.ListOfUniqueStrings),
						},
						"post_logout_redirect_uris": {
							Description: "A string that specifies the URLs that the browser can be redirected to after logout.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
							},
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.ListOfUniqueStrings),
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
							ValidateFunc: validation.StringInSlice([]string{"WEB_APP"}, false),
						},
						"acs_urls": {
							Description: "A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).",
							Type:        schema.TypeList,
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
							Description: "An ID for the certificate key pair to be used by the identity provider to sign assertions and responses. If this property is omitted, the default signing certificate for the environment is used.",
							Type:        schema.TypeString,
							Optional:    true,
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
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

	resp, r, err := apiClient.ApplicationsApplicationsApi.CreateApplication(ctx, d.Get("environment_id").(string)).CreateApplicationRequest(*applicationRequest).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationsApi.CreateApplication``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	if resp.ApplicationOIDC != nil && resp.ApplicationOIDC.GetId() != "" {
		d.SetId(resp.ApplicationOIDC.GetId())
	} else if resp.ApplicationSAML != nil && resp.ApplicationSAML.GetId() != "" {
		d.SetId(resp.ApplicationSAML.GetId())
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

	resp, r, err := apiClient.ApplicationsApplicationsApi.ReadOneApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationsApi.ReadOneApplication``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	respSecret, r, err := apiClient.ApplicationsApplicationSecretApi.ReadApplicationSecret(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application %s has no secret", d.Id())
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationsApi.ReadOneApplication``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	if resp.ApplicationOIDC != nil && resp.ApplicationOIDC.GetId() != "" {

		application := resp.ApplicationOIDC

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

		if v, ok := application.GetAssignActorRolesOk(); ok {
			d.Set("assign_actor_roles", v)
		} else {
			d.Set("assign_actor_roles", nil)
		}

		if v, ok := application.GetIconOk(); ok {
			d.Set("icon", flattenIcon(v))
		} else {
			d.Set("icon", nil)
		}

		if v, ok := application.GetAccessControlOk(); ok {
			d.Set("access_control", flattenAccessControl(v))
		} else {
			d.Set("access_control", nil)
		}

		oidcOptions, diags := flattenOIDCOptions(application, respSecret)
		if diags.HasError() {
			return diags
		}
		d.Set("oidc_options", oidcOptions)

	} else if resp.ApplicationSAML != nil && resp.ApplicationSAML.GetId() != "" {

		application := resp.ApplicationSAML

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

		if v, ok := application.GetAssignActorRolesOk(); ok {
			d.Set("assign_actor_roles", v)
		} else {
			d.Set("assign_actor_roles", nil)
		}

		if v, ok := application.GetIconOk(); ok {
			d.Set("icon", flattenIcon(v))
		} else {
			d.Set("icon", nil)
		}

		if v, ok := application.GetAccessControlOk(); ok {
			d.Set("access_control", flattenAccessControl(v))
		} else {
			d.Set("access_control", nil)
		}

		samlOptions, diags := flattenSAMLOptions(application)
		if diags.HasError() {
			return diags
		}
		d.Set("saml_options", samlOptions)

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

	_, r, err := apiClient.ApplicationsApplicationsApi.UpdateApplication(ctx, d.Get("environment_id").(string), d.Id()).UpdateApplicationRequest(*applicationRequest).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationsApi.UpdateApplication``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

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

	_, err := apiClient.ApplicationsApplicationsApi.DeleteApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationsApi.DeleteApplication``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceApplicationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
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

		grantTypes, err := expandGrantTypes(oidcOptions["grant_types"])
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot determine application `grant_types`: %v", err),
			})

			return nil, diags
		}

		// Set the object
		application = *management.NewApplicationOIDC(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT, *applicationType, grantTypes, management.EnumApplicationOIDCTokenAuthMethod(oidcOptions["token_endpoint_authn_method"].(string)))

		// set the common optional options
		applicationCommon := expandCommonOptionalAttributes(d)

		if v, ok := applicationCommon.GetDescriptionOk(); ok {
			application.SetDescription(*v)
		}

		if v, ok := applicationCommon.GetTagsOk(); ok {
			application.SetTags(v)
		}

		if v, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v)
		}

		if v, ok := applicationCommon.GetAssignActorRolesOk(); ok {
			application.SetAssignActorRoles(*v)
		}

		if v, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v)
		}

		if v, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v)
		}

		// Set the OIDC specific optional options

		if v, ok := oidcOptions["home_page_url"].(string); ok && v != "" {
			application.SetHomePageUrl(v)
		}

		if v, ok := oidcOptions["response_types"].([]string); ok && v != nil && v[0] != "" {
			obj := make([]management.EnumApplicationOIDCResponseType, 0)
			for _, j := range v {
				obj = append(obj, management.EnumApplicationOIDCResponseType(j))
			}
			application.SetResponseTypes(obj)
		}

		if v, ok := oidcOptions["pkce_enforcement"].(string); ok && v != "" {
			application.SetPkceEnforcement(management.EnumApplicationOIDCPKCEOption(v))
		}

		if v, ok := oidcOptions["redirect_uris"].([]string); ok && v != nil {
			application.SetRedirectUris(v)
		}

		if v, ok := oidcOptions["post_logout_redirect_uris"].([]string); ok && v != nil {
			application.SetPostLogoutRedirectUris(v)
		}

		if v, ok := oidcOptions["refresh_token_duration"].(int); ok {
			application.SetRefreshTokenDuration(int32(v))
		}

		if v, ok := oidcOptions["refresh_token_rolling_duration"].(int); ok {
			application.SetRefreshTokenDuration(int32(v))
		}

		if v, ok := oidcOptions["support_unsigned_request_object"].(bool); ok {
			application.SetSupportUnsignedRequestObject(v)
		}

		if v, ok := oidcOptions["mobile_app"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
			var mobile management.ApplicationOIDCAllOfMobile
			mobile, diags = expandMobile(v[0].(map[string]interface{}))
			if diags.HasError() {
				return nil, diags
			}
			application.SetMobile(mobile)
		}

		if v, ok := oidcOptions["bundle_id"].(string); ok && v != "" {
			application.SetBundleId(v)
		}

		if v, ok := oidcOptions["package_name"].(string); ok && v != "" {
			application.SetPackageName(v)
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

func expandGrantTypes(s interface{}) ([]management.EnumApplicationOIDCGrantType, error) {
	// Grant types
	grantTypesStr, ok := s.([]string)
	if !ok {
		return nil, fmt.Errorf("Cannot cast `grant_type` values to list")
	}

	grantTypes := make([]management.EnumApplicationOIDCGrantType, 0)
	for _, j := range grantTypesStr {
		grantTypes = append(grantTypes, management.EnumApplicationOIDCGrantType(j))
	}

	return grantTypes, nil
}

func expandMobile(s map[string]interface{}) (management.ApplicationOIDCAllOfMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	mobile := *management.NewApplicationOIDCAllOfMobile()
	if v, ok := s["bundle_id"].(string); ok && v != "" {
		mobile.SetBundleId(v)
	}

	if v, ok := s["package_name"].(string); ok && v != "" {
		mobile.SetPackageName(v)
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
		application = *management.NewApplicationSAML(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMAPPLICATIONPROTOCOL_SAML, *applicationType, samlOptions["acs_urls"].([]string), int32(samlOptions["assertion_duration"].(int)), samlOptions["sp_entity_id"].(string))

		// set the common optional options
		applicationCommon := expandCommonOptionalAttributes(d)

		if v, ok := applicationCommon.GetDescriptionOk(); ok {
			application.SetDescription(*v)
		}

		if v, ok := applicationCommon.GetTagsOk(); ok {
			application.SetTags(v)
		}

		if v, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v)
		}

		if v, ok := applicationCommon.GetAssignActorRolesOk(); ok {
			application.SetAssignActorRoles(*v)
		}

		if v, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v)
		}

		if v, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v)
		}

		// Set the SAML specific optional options

		if v, ok := samlOptions["assertion_signed_enabled"].(bool); ok {
			application.SetAssertionSigned(v)
		}

		if v, ok := samlOptions["idp_signing_key_id"].(string); ok && v != "" {
			application.SetIdpSigningtype(*management.NewApplicationSAMLAllOfIdpSigningtype(*management.NewApplicationSAMLAllOfIdpSigningtypeKey(v)))
		}

		if v, ok := samlOptions["nameid_format"].(string); ok && v != "" {
			application.SetNameIdFormat(v)
		}

		if v, ok := samlOptions["response_is_signed"].(bool); ok {
			application.SetResponseSigned(v)
		}

		if v, ok := samlOptions["slo_binding"].(string); ok && v != "" {
			application.SetSloBinding(management.EnumApplicationSAMLSloBinding(v))
		}

		if v, ok := samlOptions["slo_endpoint"].(string); ok && v != "" {
			application.SetSloEndpoint(v)
		}

		if v, ok := samlOptions["slo_endpoint"].(string); ok && v != "" {
			application.SetSloEndpoint(v)
		}

		if v, ok := samlOptions["slo_response_endpoint"].(string); ok && v != "" {
			application.SetSloResponseEndpoint(v)
		}

		if v, ok := samlOptions["sp_verification_certificate_ids"].([]string); ok && v != nil && v[0] != "" {
			certificates := make([]management.ApplicationSAMLAllOfSpVerificationCertificates, 0)
			for _, j := range v {
				certificate := *management.NewApplicationSAMLAllOfSpVerificationCertificates()
				certificate.SetId(j)
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
		if j, okJ := v.([]string); okJ {
			tags := make([]management.EnumApplicationTags, 0)
			for _, k := range j {
				tags = append(tags, management.EnumApplicationTags(k))
			}

			application.Tags = tags
		}
	}

	if v, ok := d.GetOk("login_page_url"); ok {
		if v != "" {
			application.SetLoginPageUrl(v.(string))
		}
	}

	if v, ok := d.GetOk("assign_actor_roles"); ok {
		application.SetAssignActorRoles(v.(bool))
	}

	if v, ok := d.GetOk("icon"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			application.SetIcon(*management.NewApplicationIcon(attrs["id"].(string), attrs["href"].(string)))
		}
	}

	if v, ok := d.GetOk("access_control"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			application.SetAccessControl(expandAccessControl(attrs))
		}
	}

	return application

}

func expandAccessControl(s map[string]interface{}) management.ApplicationAccessControl {

	accessControl := *management.NewApplicationAccessControl()

	if v, ok := s["role_type"].(string); ok && v != "" {
		accessControl.Role.SetType(v)
	}

	if v, ok := s["group"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		obj := v[0].(map[string]interface{})

		groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)
		for _, j := range obj["groups"].([]string) {
			groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(j))
		}

		accessControl.SetGroup(*management.NewApplicationAccessControlGroup(obj["type"].(string), groups))
	}

	return accessControl

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

func flattenAccessControl(s *management.ApplicationAccessControl) []interface{} {

	item := map[string]interface{}{}

	if v, ok := s.Role.GetTypeOk(); ok {
		item["role_type"] = v
	} else {
		item["role_type"] = nil
	}

	if v, ok := s.GetGroupOk(); ok {
		item["group"] = map[string]interface{}{
			"type":   v.GetGroups(),
			"groups": v.GetType(),
		}
	} else {
		item["group"] = nil
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
		"grant_types":                 application.GetGrantTypes(),
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
		var mobile interface{}
		mobile, diags = flattenMobile(v)
		if diags.HasError() {
			return nil, diags
		}

		item["mobile_app"] = mobile
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

	if v, ok := mobile.GetIntegrityDetectionOk(); ok {
		var integrityDetection interface{}
		integrityDetection, diags = flattenMobileIntegrityDetection(v)
		if diags.HasError() {
			return nil, diags
		}

		item["integrity_detection"] = integrityDetection
	} else {
		item["integrity_detection"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item), diags
}

func flattenMobileIntegrityDetection(obj *management.ApplicationOIDCAllOfMobileIntegrityDetection) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

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
		item["cache_duration"] = map[string]interface{}{
			"amount": v.GetAmount(),
			"units":  v.GetUnits(),
		}
	} else {
		item["cache_duration"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item), diags
}

func flattenSAMLOptions(application *management.ApplicationSAML) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Requried
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

	if v, ok := application.IdpSigningtype.Key.GetIdOk(); ok {
		item["idp_signing_key_id"] = v
	} else {
		item["idp_signing_key_id"] = nil
	}

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
	return append(items, item), diags

}
