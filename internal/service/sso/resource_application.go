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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceApplication() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage administrator defined applications in PingOne.",

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
				Description: fmt.Sprintf("An array that specifies the list of labels associated with the application.  Options are: `%s`", string(management.ENUMAPPLICATIONTAGS_PING_FED_CONNECTION_INTEGRATION)),
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONTAGS_PING_FED_CONNECTION_INTEGRATION)}, false)),
				},
				ConflictsWith: []string{"external_link_options", "saml_options"},
			},
			"login_page_url": {
				Description:      "A string that specifies the custom login page URL for the application. If you set the `login_page_url` property for applications in an environment that sets a custom domain, the URL should include the top-level domain and at least one additional domain level. **Warning** To avoid issues with third-party cookies in some browsers, a custom domain must be used, giving your PingOne environment the same parent domain as your authentication application. For more information about custom domains, see Custom domains.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".")),
				ConflictsWith:    []string{"external_link_options"},
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
				Description:      fmt.Sprintf("A string that specifies the user role required to access the application. A user is an admin user if the user has one or more of the following roles Organization Admin, Environment Admin, Identity Data Admin, or Client Application Developer. Options are `%s`.", string(management.ENUMAPPLICATIONACCESSCONTROLTYPE_ADMIN_USERS_ONLY)),
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONACCESSCONTROLTYPE_ADMIN_USERS_ONLY)}, false)),
				ConflictsWith:    []string{"external_link_options"},
			},
			"access_control_group_options": {
				Description: "Group access control settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:      "A string that specifies the group type required to access the application. Options are `ANY_GROUP` (the actor must belong to at least one group listed in the `groups` property) and `ALL_GROUPS` (the actor must belong to all groups listed in the `groups` property).",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ANY_GROUP", "ALL_GROUPS"}, false)),
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
			"hidden_from_app_portal": {
				Description: "A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"external_link_options": {
				Description:  "External link application specific settings.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"oidc_options", "saml_options", "external_link_options"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"home_page_url": {
							Description:      "A string that specifies the custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
					},
				},
			},
			"oidc_options": {
				Description:  "OIDC/OAuth application specific settings.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"oidc_options", "saml_options", "external_link_options"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:      fmt.Sprintf("A string that specifies the type associated with the application.  Options are `%s`, `%s`, `%s`, `%s`, `%s` and `%s`.", string(management.ENUMAPPLICATIONTYPE_WEB_APP), string(management.ENUMAPPLICATIONTYPE_NATIVE_APP), string(management.ENUMAPPLICATIONTYPE_SINGLE_PAGE_APP), string(management.ENUMAPPLICATIONTYPE_WORKER), string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP), string(management.ENUMAPPLICATIONTYPE_SERVICE)),
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONTYPE_WEB_APP), string(management.ENUMAPPLICATIONTYPE_NATIVE_APP), string(management.ENUMAPPLICATIONTYPE_SINGLE_PAGE_APP), string(management.ENUMAPPLICATIONTYPE_WORKER), string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP), string(management.ENUMAPPLICATIONTYPE_SERVICE)}, false)),
						},
						"home_page_url": {
							Description:      "A string that specifies the custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".")),
						},
						"initiate_login_uri": {
							Description:      "A string that specifies the URI to use for third-parties to begin the sign-on process for the application. If specified, PingOne redirects users to this URI to initiate SSO to PingOne. The application is responsible for implementing the relevant OIDC flow when the initiate login URI is requested. This property is required if you want the application to appear in the PingOne Application Portal. See the OIDC specification section of [Initiating Login from a Third Party](https://openid.net/specs/openid-connect-core-1_0.html#ThirdPartyInitiatedLogin) for more information.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".")),
						},
						"target_link_uri": {
							Description:      "The URI for the application. If specified, PingOne will redirect application users to this URI after a user is authenticated. In the PingOne admin console, this becomes the value of the `target_link_uri` parameter used for the Initiate Single Sign-On URL field.  Both `http://` and `https://` URLs are permitted as well as custom mobile native schema (e.g., `org.bxretail.app://target`).",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://target`).")),
						},
						"grant_types": {
							Description: fmt.Sprintf("A list that specifies the grant type for the authorization request.  Options are `%s`, `%s`, `%s` and `%s`.", string(management.ENUMAPPLICATIONOIDCGRANTTYPE_AUTHORIZATION_CODE), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_IMPLICIT), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_CLIENT_CREDENTIALS)),
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONOIDCGRANTTYPE_AUTHORIZATION_CODE), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_IMPLICIT), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN), string(management.ENUMAPPLICATIONOIDCGRANTTYPE_CLIENT_CREDENTIALS)}, false)),
							},
							Required: true,
						},
						"response_types": {
							Description: fmt.Sprintf("A list that specifies the code or token type returned by an authorization request.  Note that `%s` cannot be used in an authorization request with `%s` or `%s` because PingOne does not currently support OIDC hybrid flows.", management.ENUMAPPLICATIONOIDCRESPONSETYPE_CODE, management.ENUMAPPLICATIONOIDCRESPONSETYPE_TOKEN, management.ENUMAPPLICATIONOIDCRESPONSETYPE_ID_TOKEN),
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_CODE), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_TOKEN), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_ID_TOKEN)}, false)),
							},
							Optional: true,
						},
						"token_endpoint_authn_method": {
							Description:      fmt.Sprintf("A string that specifies the client authentication methods supported by the token endpoint.  Options are `%s`, `%s` and `%s`.", string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_NONE), string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC), string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_CLIENT_SECRET_POST)),
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_NONE), string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC), string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_CLIENT_SECRET_POST)}, false)),
						},
						"par_requirement": {
							Description:      fmt.Sprintf("A string that specifies whether pushed authorization requests (PAR) are required.  Options are `%s` and `%s`.", string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL), string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_REQUIRED)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL),
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL), string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_REQUIRED)}, false)),
						},
						"par_timeout": {
							Description:      "An integer that specifies the pushed authorization request (PAR) timeout in seconds.  If a value is not provided, the default value is `60`.  Valid values are between `1` and `600`.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          60,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 600)),
						},
						"pkce_enforcement": {
							Description:      fmt.Sprintf("A string that specifies how `PKCE` request parameters are handled on the authorize request.  Options are `%s`, `%s` and `%s`.", string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL), string(management.ENUMAPPLICATIONOIDCPKCEOPTION_REQUIRED), string(management.ENUMAPPLICATIONOIDCPKCEOPTION_S256_REQUIRED)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL),
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL), string(management.ENUMAPPLICATIONOIDCPKCEOPTION_REQUIRED), string(management.ENUMAPPLICATIONOIDCPKCEOPTION_S256_REQUIRED)}, false)),
						},
						"redirect_uris": {
							Description: "A list of strings that specifies the allowed callback URIs for the authentication response.    The provided URLs are expected to use the `https://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  The `http` schema is only permitted where the host is `localhost` or `127.0.0.1`.",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(\S+:\/\/).+)`), "Expected value to have a url with schema of \"https\" or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".")),
							},
							Optional: true,
						},
						"allow_wildcards_in_redirect_uris": {
							Description: "A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"post_logout_redirect_uris": {
							Description: "A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://logout`).")),
							},
							Optional: true,
						},
						"refresh_token_duration": {
							Description:      "An integer that specifies the lifetime in seconds of the refresh token. If a value is not provided, the default value is `2592000`, or 30 days. Valid values are between `60` and `2147483647`. If the `refresh_token_rolling_duration` property is specified for the application, then this property value must be less than or equal to the value of `refresh_token_rolling_duration`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          2592000,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(60, 2147483647)),
						},
						"refresh_token_rolling_duration": {
							Description:      "An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required. If a value is not provided, the default value is `15552000`, or 180 days. Valid values are between `60` and `2147483647`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          15552000,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(60, 2147483647)),
						},
						"refresh_token_rolling_grace_period_duration": {
							Description:      "The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens. This is useful in the case of network errors on the client. Valid values are between `0` and `86400` seconds. `Null` is treated the same as `0`.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 86400)),
						},
						"additional_refresh_token_replay_protection_enabled": {
							Description: "A boolean that, when set to `true` (the default), if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens. Setting this to null equates to a `false` setting.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
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
						"certificate_based_authentication": {
							Description: "Certificate based authentication settings. This parameter block can only be set where the application's `type` parameter is set to `NATIVE_APP`.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Description:      "A string that represents a PingOne ID for the issuance certificate key.  The key must be of type `ISSUANCE`.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
									},
								},
							},
						},
						"support_unsigned_request_object": {
							Description: "A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If false or null (default), an unsigned request object is not allowed.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"require_signed_request_object": {
							Description: "A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null (default), a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"cors_settings": resourceApplicationSchemaCorsSettings(),
						"mobile_app": {
							Description: fmt.Sprintf("Mobile application integration settings for `%s` type applications.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// If we have an app configured, then we need to compute the diff.
								if _, ok := d.GetOk("oidc_options.0.mobile_app.0.bundle_id"); ok {
									return false
								}

								if _, ok := d.GetOk("oidc_options.0.mobile_app.0.package_name"); ok {
									return false
								}

								if _, ok := d.GetOk("oidc_options.0.mobile_app.0.huawei_app_id"); ok {
									return false
								}

								if _, ok := d.GetOk("oidc_options.0.mobile_app.0.huawei_package_name"); ok {
									return false
								}

								// If no app configured, we can suppress the diff.
								return true
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bundle_id": {
										Description:      "A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable.  Changing this value will trigger a replacement plan of this resource.",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"package_name": {
										Description:      "A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.  Changing this value will trigger a replacement plan of this resource.",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"huawei_app_id": {
										Description:      "The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_package_name`.  Changing this value will trigger a replacement plan of this resource.",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										RequiredWith:     []string{"oidc_options.0.mobile_app.0.huawei_app_id", "oidc_options.0.mobile_app.0.huawei_package_name"},
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"huawei_package_name": {
										Description:      "The package name associated with the application, for push notifications in native apps. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_app_id`.  Changing this value will trigger a replacement plan of this resource.",
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										RequiredWith:     []string{"oidc_options.0.mobile_app.0.huawei_app_id", "oidc_options.0.mobile_app.0.huawei_package_name"},
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"passcode_refresh_seconds": {
										Description:      "The amount of time a passcode should be displayed before being replaced with a new passcode - must be between 30 and 60.",
										Type:             schema.TypeInt,
										Optional:         true,
										Default:          30,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(30, 60)),
									},
									"universal_app_link": {
										Description:      "A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code. The URI prefix can be set to a universal link with a valid value (which can be a URL address that starts with `HTTP://` or `HTTPS://`, such as `https://www.bxretail.org`), or an app schema, which is just a string and requires no special validation.",
										Type:             schema.TypeString,
										Optional:         true,
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
												"excluded_platforms": {
													Description: fmt.Sprintf("You can enable device integrity checking separately for Android and iOS by setting `enabled` to `true` and then using `excluded_platforms` to specify the OS where you do not want to use device integrity checking. The values to use are `%s` and `%s` (all upper case). Note that this is implemented as an array even though currently you can only include a single value.  If `%s` is not included in this list, the `google_play` attribute block must be configured.", string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE), string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_IOS), string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE)),
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Elem: &schema.Schema{
														Type:             schema.TypeString,
														ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE), string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_IOS)}, false)),
													},
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
																Description:      fmt.Sprintf("A string that specifies the time units of the `amount` parameter. Options are `%s` and `%s`.", string(management.ENUMDURATIONUNITMINSHOURS_MINUTES), string(management.ENUMDURATIONUNITMINSHOURS_HOURS)),
																Type:             schema.TypeString,
																Optional:         true,
																Default:          string(management.ENUMDURATIONUNITMINSHOURS_MINUTES),
																ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMDURATIONUNITMINSHOURS_MINUTES), string(management.ENUMDURATIONUNITMINSHOURS_HOURS)}, false)),
															},
														},
													},
												},
												"google_play": {
													Description: "Required when `excluded_platforms` is unset or does not include `GOOGLE`.  A single block that describes Google Play Integrity API credential settings for Android device integrity detection.",
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"decryption_key": {
																Description: "Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.",
																Type:        schema.TypeString,
																Optional:    true,
																Sensitive:   true,
																DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
																	return old == "DUMMY_SUPPRESS_VALUE"
																},
																ConflictsWith: []string{"oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.service_account_credentials_json"},
															},
															"service_account_credentials_json": {
																Description: "Contents of the JSON file that represents your Service Account Credentials. This parameter must be provided if you have set `verification_type` to `GOOGLE`.  Cannot be set with `decryption_key` or `verification_key`.",
																Type:        schema.TypeString,
																Optional:    true,
																Sensitive:   true,
																DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
																	return old == "DUMMY_SUPPRESS_VALUE"
																},
																ConflictsWith: []string{"oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.decryption_key", "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.verification_key"},
															},
															"verification_key": {
																Description: "Play Integrity verdict signature verification key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.",
																Type:        schema.TypeString,
																Optional:    true,
																Sensitive:   true,
																DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
																	return old == "DUMMY_SUPPRESS_VALUE"
																},
																ConflictsWith: []string{"oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.service_account_credentials_json"},
															},
															"verification_type": {
																Description:      "The type of verification that should be used. The possible values are `GOOGLE` and `INTERNAL`. Using internal verification will not count against your Google API call quota. The value you select for this attribute determines what other parameters you must provide. When set to `GOOGLE`, you must provide `service_account_credentials_json`. When set to `INTERNAL`, you must provide both `decryption_key` and `verification_key`.",
																Type:             schema.TypeString,
																Required:         true,
																ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE), string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL)}, false)),
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
							Description:      "**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead. A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							Deprecated:       "This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead.",
						},
						"package_name": {
							Description:      "**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead. A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							Deprecated:       "This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead.",
						},
					},
				},
			},
			"saml_options": {
				Description:  "SAML application specific settings.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"oidc_options", "saml_options", "external_link_options"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"home_page_url": {
							Description:      "A string that specifies the custom home page URL for the application.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"type": {
							Description:      fmt.Sprintf("A string that specifies the type associated with the application.  Options are `%s` and `%s`.", string(management.ENUMAPPLICATIONTYPE_WEB_APP), string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(management.ENUMAPPLICATIONTYPE_WEB_APP),
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONTYPE_WEB_APP), string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP)}, false)),
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
							Computed:         true,
							ConflictsWith:    []string{"saml_options.0.idp_signing_key"},
							Deprecated:       "The `idp_signing_key_id` attribute is deprecated and will be removed in the next major release.  Please use the `idp_signing_key` block going forward.",
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"idp_signing_key": {
							Description:   "SAML application assertion/response signing key settings.  Use with `assertion_signed_enabled` to enable assertion signing and/or `response_is_signed` to enable response signing.  It's highly recommended, and best practice, to define signing key settings for the configured SAML application.  However if this property is omitted, the default signing certificate for the environment is used.  This parameter will become a required field in the next major release of the provider.",
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"saml_options.0.idp_signing_key_id"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Description:      fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)}, false)),
									},
									"key_id": {
										Description:      "An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
									},
								},
							},
						},
						"enable_requested_authn_context": {
							Description: "A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
							Type:        schema.TypeBool,
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
							Description:      fmt.Sprintf("A string that specifies the binding protocol to be used for the logout response. Options are `%s` and `%s`.  Existing configurations with no data default to `%s`.", string(management.ENUMAPPLICATIONSAMLSLOBINDING_REDIRECT), string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST), string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST)),
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONSAMLSLOBINDING_REDIRECT), string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST)}, false)),
							Optional:         true,
							Default:          string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST),
						},
						"slo_endpoint": {
							Description:      "A string that specifies the logout endpoint URL. This is an optional property. However, if a logout endpoint URL is not defined, logout actions result in an error.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"slo_response_endpoint": {
							Description:      "A string that specifies the endpoint URL to submit the logout response. If a value is not provided, the `slo_endpoint` property value is used to submit SLO response.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"slo_window": {
							Description:      "An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request. The minimum value is `1` hour and the maximum is `24` hours.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 24)),
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
						"cors_settings": resourceApplicationSchemaCorsSettings(),
					},
				},
			},
		},
	}
}

func resourceApplicationSchemaCorsSettings() *schema.Schema {
	return &schema.Schema{
		Description: "A single block that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings.",
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"behavior": {
					Description:      fmt.Sprintf("A string that specifies the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `%s` (rejects all CORS requests) and `%s` (rejects all CORS requests except those listed in `origins`).", string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS), string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS)),
					Type:             schema.TypeString,
					Required:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS), string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS)}, false)),
				},
				"origins": {
					Description: fmt.Sprintf("A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value must be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 or IPv6 address.  Subdomains may use the wildcard (`*`) to match any string.  Must be non-empty when `behavior` is `%s` and must be omitted or empty when `behavior` is `%s`.  Limited to 20 values.", string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS), string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS)),
					Type:        schema.TypeSet,
					MaxItems:    20,
					Elem: &schema.Schema{
						Type:             schema.TypeString,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(\S+:\/\/).+)`), "Expected value to have a url with schema of \"https\" or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".")),
					},
					Optional: true,
				},
			},
		},
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

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

	if _, ok := d.GetOk("external_link_options"); ok {
		var application *management.ApplicationExternalLink
		application, diags = expandApplicationExternalLink(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationExternalLink = application
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.CreateApplication(ctx, d.Get("environment_id").(string)).CreateApplicationRequest(*applicationRequest).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"CreateApplication",
		applicationWriteCustomError,
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
	} else if respObject.ApplicationExternalLink != nil && respObject.ApplicationExternalLink.GetId() != "" {
		d.SetId(respObject.ApplicationExternalLink.GetId())
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

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.ReadOneApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
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

			func() (any, *http.Response, error) {
				fO, fR, fErr := apiClient.ApplicationSecretApi.ReadApplicationSecret(ctx, d.Get("environment_id").(string), d.Id()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
			},
			"ReadApplicationSecret",
			sdk.CustomErrorResourceNotFoundWarning,
			func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

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

		if v, ok := application.GetHiddenFromAppPortalOk(); ok {
			d.Set("hidden_from_app_portal", v)
		} else {
			d.Set("hidden_from_app_portal", nil)
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

		if v, ok := application.GetHiddenFromAppPortalOk(); ok {
			d.Set("hidden_from_app_portal", v)
		} else {
			d.Set("hidden_from_app_portal", nil)
		}

		d.Set("saml_options", flattenSAMLOptions(application))

	} else if respObject.ApplicationExternalLink != nil && respObject.ApplicationExternalLink.GetId() != "" {

		application := respObject.ApplicationExternalLink

		d.Set("name", application.GetName())
		d.Set("enabled", application.GetEnabled())

		if v, ok := application.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		if v, ok := application.GetIconOk(); ok {
			d.Set("icon", flattenIcon(v))
		} else {
			d.Set("icon", nil)
		}

		if v, ok := application.GetAccessControlOk(); ok {

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
			d.Set("access_control_group_options", nil)
		}

		if v, ok := application.GetHiddenFromAppPortalOk(); ok {
			d.Set("hidden_from_app_portal", v)
		} else {
			d.Set("hidden_from_app_portal", nil)
		}

		externalLinkOpts := make([]interface{}, 0)

		d.Set("external_link_options", append(externalLinkOpts, map[string]interface{}{
			"home_page_url": application.GetHomePageUrl(),
		}))

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

	if _, ok := d.GetOk("external_link_options"); ok {
		var application *management.ApplicationExternalLink
		application, diags = expandApplicationExternalLink(d)
		if diags.HasError() {
			return diags
		}
		applicationRequest.ApplicationExternalLink = application
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.UpdateApplication(ctx, d.Get("environment_id").(string), d.Id()).UpdateApplicationRequest(*applicationRequest).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"UpdateApplication",
		applicationWriteCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := apiClient.ApplicationsApi.DeleteApplication(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), nil, fR, fErr)
		},
		"DeleteApplication",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceApplicationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.SetId(attributes["application_id"])

	resourceApplicationRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func applicationWriteCustomError(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Wildcards in redirect URis
	m, err := regexp.MatchString("^Wildcards are not allowed in redirect URIs.", error.GetMessage())
	if err != nil {
		diags = diag.FromErr(fmt.Errorf("Invalid regexp: Wildcards are not allowed in redirect URIs."))
		return diags
	}
	if m {
		diags = diag.FromErr(fmt.Errorf("Current configuration is invalid as wildcards are not allowed in redirect URIs.  Wildcards can be enabled by setting `allow_wildcards_in_redirect_uris` to `true`."))

		return diags
	}

	return nil
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

		if v1, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v1)
		}

		if v1, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v1)
		}

		if v1, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v1)
		}

		if v1, ok := applicationCommon.GetHiddenFromAppPortalOk(); ok {
			application.SetHiddenFromAppPortal(*v1)
		}

		// Set the OIDC specific optional options

		if v1, ok := oidcOptions["cors_settings"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			corsSettings := expandCorsSettings(v1[0].(map[string]interface{}))
			application.SetCorsSettings(*corsSettings)
		}

		if v1, ok := oidcOptions["home_page_url"].(string); ok && v1 != "" {
			application.SetHomePageUrl(v1)
		}

		if v1, ok := oidcOptions["initiate_login_uri"].(string); ok && v1 != "" {
			application.SetInitiateLoginUri(v1)
		}

		if v1, ok := oidcOptions["target_link_uri"].(string); ok && v1 != "" {
			application.SetTargetLinkUri(v1)
		}

		if v1, ok := oidcOptions["response_types"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]management.EnumApplicationOIDCResponseType, 0)
			for _, j := range v1.List() {
				obj = append(obj, management.EnumApplicationOIDCResponseType(j.(string)))
			}
			application.SetResponseTypes(obj)
		}

		if v1, ok := oidcOptions["par_requirement"].(string); ok && v1 != "" {
			application.SetParRequirement(management.EnumApplicationOIDCPARRequirement(v1))
		}

		if v1, ok := oidcOptions["par_timeout"].(int); ok {
			application.SetParTimeout(int32(v1))
		}

		if v1, ok := oidcOptions["pkce_enforcement"].(string); ok && v1 != "" {
			application.SetPkceEnforcement(management.EnumApplicationOIDCPKCEOption(v1))
		}

		if v1, ok := oidcOptions["redirect_uris"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]string, 0)
			for _, j := range v1.List() {
				obj = append(obj, j.(string))
			}
			application.SetRedirectUris(obj)
		}

		if v1, ok := oidcOptions["allow_wildcards_in_redirect_uris"].(bool); ok {
			application.SetAllowWildcardInRedirectUris(v1)
		}

		if v1, ok := oidcOptions["post_logout_redirect_uris"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			obj := make([]string, 0)
			for _, j := range v1.List() {
				obj = append(obj, j.(string))
			}
			application.SetPostLogoutRedirectUris(obj)
		}

		if v1, ok := oidcOptions["refresh_token_duration"].(int); ok {
			application.SetRefreshTokenDuration(int32(v1))
		}

		if v1, ok := oidcOptions["refresh_token_rolling_duration"].(int); ok {
			application.SetRefreshTokenRollingDuration(int32(v1))
		}

		if v1, ok := oidcOptions["refresh_token_rolling_grace_period_duration"].(int); ok {
			application.SetRefreshTokenRollingGracePeriodDuration(int32(v1))
		}

		if v1, ok := oidcOptions["additional_refresh_token_replay_protection_enabled"].(bool); ok {
			application.SetAdditionalRefreshTokenReplayProtectionEnabled(v1)
		}

		if v, ok := oidcOptions["tags"]; ok {
			if j, okJ := v.([]interface{}); okJ {
				tags := make([]management.EnumApplicationTags, 0)
				for _, k := range j {
					tags = append(tags, management.EnumApplicationTags(k.(string)))
				}

				application.Tags = tags
			}
		}

		application.SetAssignActorRoles(false)

		if v1, ok := oidcOptions["certificate_based_authentication"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			if *applicationType != management.ENUMAPPLICATIONTYPE_NATIVE_APP {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("`certificate_based_authentication` can only be set with applications that have a `type` value of `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
				})
				return nil, diags
			}

			application.SetKerberos(*expandKerberos(v1[0].(map[string]interface{})))
		}

		if v1, ok := oidcOptions["support_unsigned_request_object"].(bool); ok {
			application.SetSupportUnsignedRequestObject(v1)
		}

		if v1, ok := oidcOptions["require_signed_request_object"].(bool); ok {
			application.SetRequireSignedRequestObject(v1)
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

func expandKerberos(s map[string]interface{}) *management.ApplicationOIDCAllOfKerberos {

	key := management.NewApplicationOIDCAllOfKerberosKey(s["key_id"].(string))
	kerberos := management.NewApplicationOIDCAllOfKerberos(*key)

	return kerberos
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

	if v, ok := s["huawei_app_id"].(string); ok && v != "" {
		mobile.SetHuaweiAppId(v)
	}

	if v, ok := s["huawei_package_name"].(string); ok && v != "" {
		mobile.SetHuaweiPackageName(v)
	}

	if v, ok := s["passcode_refresh_seconds"].(int); ok {
		mobile.SetPasscodeRefreshDuration(*management.NewApplicationOIDCAllOfMobilePasscodeRefreshDuration(int32(v), management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS))
	}

	if v, ok := s["universal_app_link"].(string); ok && v != "" {
		mobile.SetUriPrefix(v)
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

		googleVerificationIncluded := true
		if j, okJ := obj["excluded_platforms"].([]interface{}); okJ && len(j) > 0 && j[0] != nil {
			list := make([]management.EnumMobileIntegrityDetectionPlatform, 0)

			for _, platform := range j {
				list = append(list, management.EnumMobileIntegrityDetectionPlatform(platform.(string)))
				if platform == string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE) {
					googleVerificationIncluded = false
				}
			}

			integrityDetection.SetExcludedPlatforms(list)
		}

		if j, okJ := obj["cache_duration"].([]interface{}); okJ && len(j) > 0 && j[0] != nil {
			integrityDetection.SetCacheDuration(expandMobileIntegrityCacheDuration(j[0].(map[string]interface{})))
		} else {
			if integrityDetection.GetMode() == management.ENUMENABLEDSTATUS_ENABLED {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Attribute block `cache_duration` is required when the mobile integrity check is enabled in the application",
				})

				return nil, diags
			}
		}

		if j, okJ := obj["google_play"].([]interface{}); okJ && len(j) > 0 && j[0] != nil {
			googlePlay, d := expandMobileIntegrityGooglePlay(j[0].(map[string]interface{}))
			if d.HasError() {
				return nil, d
			}

			integrityDetection.SetGooglePlay(*googlePlay)
		} else {

			if integrityDetection.GetMode() == management.ENUMENABLEDSTATUS_ENABLED && googleVerificationIncluded {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Attribute block `google_play` is required when the mobile integrity check is enabled in the application and `excluded_platforms` is unset, or `excluded_platforms` is not configured with `GOOGLE`.",
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

func expandMobileIntegrityGooglePlay(s map[string]interface{}) (*management.ApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay, diag.Diagnostics) {
	var diags diag.Diagnostics

	obj := management.NewApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay()

	obj.SetVerificationType(management.EnumApplicationNativeGooglePlayVerificationType(s["verification_type"].(string)))

	if obj.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE {
		if v, ok := s["service_account_credentials_json"].(string); ok && v != "" {
			obj.SetServiceAccountCredentials(v)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Attribute `service_account_credentials_json` is required when the `verification_type` is set to `GOOGLE`.",
			})

			return nil, diags
		}
	}

	if obj.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL {
		if v, ok := s["decryption_key"].(string); ok && v != "" {
			obj.SetDecryptionKey(v)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Attribute `decryption_key` is required when the `verification_type` is set to `INTERNAL`.",
			})

			return nil, diags
		}

		if v, ok := s["verification_key"].(string); ok && v != "" {
			obj.SetVerificationKey(v)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Attribute `verification_key` is required when the `verification_type` is set to `INTERNAL`.",
			})

			return nil, diags
		}
	}

	return obj, diags
}

func expandCorsSettings(s map[string]interface{}) *management.ApplicationCorsSettings {
	cors := management.NewApplicationCorsSettings(management.EnumApplicationCorsSettingsBehavior(s["behavior"].(string)))

	if v2, ok := s["origins"].(*schema.Set); ok && v2 != nil && len(v2.List()) > 0 && v2.List()[0] != nil {
		obj := make([]string, 0)
		for _, j := range v2.List() {
			obj = append(obj, j.(string))
		}
		cors.SetOrigins(obj)
	}

	return cors
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

		if v1, ok := applicationCommon.GetLoginPageUrlOk(); ok {
			application.SetLoginPageUrl(*v1)
		}

		if v1, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v1)
		}

		if v1, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v1)
		}

		if v1, ok := applicationCommon.GetHiddenFromAppPortalOk(); ok {
			application.SetHiddenFromAppPortal(*v1)
		}

		// Set the SAML specific optional options

		if v1, ok := samlOptions["cors_settings"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			corsSettings := expandCorsSettings(v1[0].(map[string]interface{}))
			application.SetCorsSettings(*corsSettings)
		}

		if v1, ok := samlOptions["home_page_url"].(string); ok && v1 != "" {
			application.SetHomePageUrl(v1)
		}

		if v1, ok := samlOptions["assertion_signed_enabled"].(bool); ok {
			application.SetAssertionSigned(v1)
		}

		if v1, ok := samlOptions["idp_signing_key_id"].(string); ok && v1 != "" {
			application.SetIdpSigning(*management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(v1)))
		}

		if v1, ok := samlOptions["idp_signing_key"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {

			idpSigningOptions := v1[0].(map[string]interface{})

			idpSigning := *management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(idpSigningOptions["key_id"].(string)))
			idpSigning.SetAlgorithm(management.EnumCertificateKeySignagureAlgorithm(idpSigningOptions["algorithm"].(string)))

			application.SetIdpSigning(idpSigning)
		}

		if v1, ok := samlOptions["enable_requested_authn_context"].(bool); ok {
			application.SetEnableRequestedAuthnContext(v1)
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

		if v1, ok := samlOptions["slo_window"].(int); ok && v1 > 0 {
			application.SetSloWindow(int32(v1))
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

// External Link

func expandApplicationExternalLink(d *schema.ResourceData) (*management.ApplicationExternalLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	var application management.ApplicationExternalLink

	if v, ok := d.Get("external_link_options").([]interface{}); ok && len(v) > 0 && v[0] != nil {

		externalLinkOptions := v[0].(map[string]interface{})

		application = *management.NewApplicationExternalLink(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMAPPLICATIONPROTOCOL_EXTERNAL_LINK, management.ENUMAPPLICATIONTYPE_PORTAL_LINK_APP, externalLinkOptions["home_page_url"].(string))

		// set the common optional options
		applicationCommon := expandCommonOptionalAttributes(d)

		if v1, ok := applicationCommon.GetDescriptionOk(); ok {
			application.SetDescription(*v1)
		}

		if v1, ok := applicationCommon.GetIconOk(); ok {
			application.SetIcon(*v1)
		}

		if v1, ok := applicationCommon.GetAccessControlOk(); ok {
			application.SetAccessControl(*v1)
		}

		if v1, ok := applicationCommon.GetHiddenFromAppPortalOk(); ok {
			application.SetHiddenFromAppPortal(*v1)
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("External Link options not available for application: %s", d.Get("name")),
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

	if v, ok := d.GetOk("login_page_url"); ok {
		if v != "" {
			application.SetLoginPageUrl(v.(string))
		}
	}

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

	if v, ok := d.GetOk("hidden_from_app_portal"); ok {
		application.SetHiddenFromAppPortal(v.(bool))
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

	if v, ok := application.GetInitiateLoginUriOk(); ok {
		item["initiate_login_uri"] = v
	} else {
		item["initiate_login_uri"] = nil
	}

	if v, ok := application.GetTargetLinkUriOk(); ok {
		item["target_link_uri"] = v
	} else {
		item["target_link_uri"] = nil
	}

	if v, ok := application.GetResponseTypesOk(); ok {
		item["response_types"] = v
	} else {
		item["response_types"] = nil
	}

	if v, ok := application.GetParRequirementOk(); ok {
		item["par_requirement"] = v
	} else {
		item["par_requirement"] = nil
	}

	if v, ok := application.GetParTimeoutOk(); ok {
		item["par_timeout"] = v
	} else {
		item["par_timeout"] = nil
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

	if v, ok := application.GetAllowWildcardInRedirectUrisOk(); ok {
		item["allow_wildcards_in_redirect_uris"] = v
	} else {
		item["allow_wildcards_in_redirect_uris"] = nil
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

	if v, ok := application.GetRefreshTokenRollingGracePeriodDurationOk(); ok {
		item["refresh_token_rolling_grace_period_duration"] = v
	} else {
		item["refresh_token_rolling_grace_period_duration"] = nil
	}

	if v, ok := application.GetAdditionalRefreshTokenReplayProtectionEnabledOk(); ok {
		item["additional_refresh_token_replay_protection_enabled"] = v
	} else {
		item["additional_refresh_token_replay_protection_enabled"] = nil
	}

	if v, ok := secret.GetSecretOk(); ok {
		item["client_secret"] = v
	} else {
		item["client_secret"] = nil
	}

	if v, ok := application.GetKerberosOk(); ok {
		item["certificate_based_authentication"] = flattenKerberos(v)
	} else {
		item["certificate_based_authentication"] = nil
	}

	if v, ok := application.GetSupportUnsignedRequestObjectOk(); ok {
		item["support_unsigned_request_object"] = v
	} else {
		item["support_unsigned_request_object"] = nil
	}

	if v, ok := application.GetRequireSignedRequestObjectOk(); ok {
		item["require_signed_request_object"] = v
	} else {
		item["require_signed_request_object"] = nil
	}

	if v, ok := application.GetCorsSettingsOk(); ok {
		item["cors_settings"] = flattenCorsSettings(v)
	} else {
		item["cors_settings"] = nil
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

func flattenKerberos(kerberos *management.ApplicationOIDCAllOfKerberos) interface{} {

	item := map[string]interface{}{}

	if v, ok := kerberos.GetKeyOk(); ok {
		item["key_id"] = v.GetId()
	} else {
		item["key_id"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenCorsSettings(s *management.ApplicationCorsSettings) interface{} {

	item := map[string]interface{}{}

	if v, ok := s.GetBehaviorOk(); ok {
		item["behavior"] = v
	} else {
		item["behavior"] = nil
	}

	if v, ok := s.GetOriginsOk(); ok {
		item["origins"] = v
	} else {
		item["origins"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
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

	if v, ok := mobile.GetHuaweiAppIdOk(); ok {
		item["huawei_app_id"] = v
	} else {
		item["huawei_app_id"] = nil
	}

	if v, ok := mobile.GetHuaweiPackageNameOk(); ok {
		item["huawei_package_name"] = v
	} else {
		item["huawei_package_name"] = nil
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

	if v, ok := mobile.GetUriPrefixOk(); ok {
		item["universal_app_link"] = v
	} else {
		item["universal_app_link"] = nil
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

	if v, ok := obj.GetExcludedPlatformsOk(); ok {

		items := make([]string, 0)
		for _, platform := range v {
			items = append(items, string(platform))
		}

		item["excluded_platforms"] = items

	} else {
		item["excluded_platforms"] = nil
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

	if v, ok := obj.GetGooglePlayOk(); ok {
		googlePlay := map[string]interface{}{
			"verification_type": v.GetVerificationType(),
		}

		if v.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL {
			googlePlay["decryption_key"] = "DUMMY_SUPPRESS_VALUE"
			googlePlay["verification_key"] = "DUMMY_SUPPRESS_VALUE"
		}

		if v.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE {
			googlePlay["service_account_credentials_json"] = "DUMMY_SUPPRESS_VALUE"
		}

		googlePlays := make([]interface{}, 0)
		item["google_play"] = append(googlePlays, googlePlay)
	} else {
		item["google_play"] = nil
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
	if v, ok := application.GetHomePageUrlOk(); ok {
		item["home_page_url"] = v
	} else {
		item["home_page_url"] = nil
	}

	if v, ok := application.GetAssertionSignedOk(); ok {
		item["assertion_signed_enabled"] = v
	} else {
		item["assertion_signed_enabled"] = nil
	}

	if v, ok := application.GetIdpSigningOk(); ok {
		item["idp_signing_key"], item["idp_signing_key_id"] = flattenIdpSigningOptions(v)
	} else {
		item["idp_signing_key"] = nil
		item["idp_signing_key_id"] = nil
	}

	if v, ok := application.GetEnableRequestedAuthnContextOk(); ok {
		item["enable_requested_authn_context"] = v
	} else {
		item["enable_requested_authn_context"] = nil
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

	if v, ok := application.GetSloWindowOk(); ok {
		item["slo_window"] = v
	} else {
		item["slo_window"] = nil
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

	if v, ok := application.GetCorsSettingsOk(); ok {
		item["cors_settings"] = flattenCorsSettings(v)
	} else {
		item["cors_settings"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)

}

func flattenIdpSigningOptions(idpSigning *management.ApplicationSAMLAllOfIdpSigning) (interface{}, *string) {

	item := map[string]interface{}{}

	var signingKeyID *string

	if v, ok := idpSigning.GetAlgorithmOk(); ok {
		item["algorithm"] = string(*v)
	} else {
		item["algorithm"] = nil
	}

	item["key_id"] = nil
	if v, ok := idpSigning.GetKeyOk(); ok {
		if v1, ok := v.GetIdOk(); ok {
			item["key_id"] = v1
			signingKeyID = v1
		}
	}

	items := make([]interface{}, 0)
	return append(items, item), signingKeyID
}
