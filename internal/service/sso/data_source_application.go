package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
)

// Types
type ApplicationDataSource serviceClientType

type applicationDataSourceModel struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId             pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	Name                      types.String                 `tfsdk:"name"`
	Description               types.String                 `tfsdk:"description"`
	Enabled                   types.Bool                   `tfsdk:"enabled"`
	Tags                      types.Set                    `tfsdk:"tags"`
	LoginPageUrl              types.String                 `tfsdk:"login_page_url"`
	Icon                      types.Object                 `tfsdk:"icon"`
	AccessControlRoleType     types.String                 `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.Object                 `tfsdk:"access_control_group_options"`
	HiddenFromAppPortal       types.Bool                   `tfsdk:"hidden_from_app_portal"`
	ExternalLinkOptions       types.Object                 `tfsdk:"external_link_options"`
	OIDCOptions               types.Object                 `tfsdk:"oidc_options"`
	SAMLOptions               types.Object                 `tfsdk:"saml_options"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &ApplicationDataSource{}
)

// New Object
func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// Metadata
func (r *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	applicationIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of the application.",
	).ExactlyOneOf([]string{"application_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the application.",
	).ExactlyOneOf([]string{"application_id", "name"})

	accessControlRoleTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The user role required to access the application.  A user is an admin user if the user has one or more of the following roles: `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.",
	)

	externalLinkHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.",
	)

	oidcHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcOptionsDevicePathIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that describes a unique identifier within an environment for a device authorization grant flow to provide a short identifier to the application. This property is ignored when the `device_custom_verification_uri` property is configured.",
	)

	oidcOptionsDeviceCustomVerificationUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an optional custom verification URI that is returned for the `/device_authorization` endpoint.",
	)

	oidcOptionsDeviceTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the length of time (in seconds) that the `userCode` and `deviceCode` returned by the `/device_authorization` endpoint are valid.",
	)

	oidcOptionsDevicePollingIntervalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the frequency (in seconds) for the client to poll the `/as/token` endpoint.",
	)

	oidcJwksDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	)

	oidcJwksUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	)

	oidcOptionsPKCEEnforcementDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how `PKCE` request parameters are handled on the authorize request.",
	)

	oidcAllowWildcardsInRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).",
	)

	oidcPostLogoutRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).",
	)

	oidcAdditionalRefreshTokenReplayProtectionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens.",
	).DefaultValue("true")

	oidcRequireSignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
	).DefaultValue("false")

	googlePlayDecryptionKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.",
	)

	samlEnableRequestedAuthnContextDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
	)

	samlSpEncryptionAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The algorithm to use when encrypting assertions.",
	).AllowedValuesEnum(management.AllowedEnumCertificateKeyEncryptionAlgorithmEnumValues)

	samlDefaultTargetUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve a PingOne application in an environment by ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the application exists."),
			),
			"application_id": schema.StringAttribute{
				Description:         applicationIdDescription.Description,
				MarkdownDescription: applicationIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
				},
			},
			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("application_id"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the application.").Description,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the application is enabled in the environment.").Description,
				Computed:    true,
			},
			"tags": schema.SetAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An array that specifies the list of labels associated with the application.").Description,
				ElementType: types.StringType,
				Computed:    true,
			},
			"login_page_url": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the custom login page URL for the application.").Description,
				Computed:    true,
			},
			"icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The HREF and the ID for the application icon.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID for the application icon.").Description,
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The HREF for the application icon.").Description,
						Computed:    true,
					},
				},
			},
			"access_control_role_type": schema.StringAttribute{
				Description:         accessControlRoleTypeDescription.Description,
				MarkdownDescription: accessControlRoleTypeDescription.MarkdownDescription,
				Computed:            true,
			},
			"access_control_group_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Group access control settings.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the group type required to access the application.").Description,
						Computed:    true,
					},
					"groups": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set that specifies the group IDs for the groups the actor must belong to for access to the application.").Description,
						ElementType: types.StringType,
						Computed:    true,
					},
				},
			},
			"hidden_from_app_portal": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.").Description,
				Computed:    true,
			},
			"external_link_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("External link application specific settings.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						Description:         externalLinkHomePageURLDescription.Description,
						MarkdownDescription: externalLinkHomePageURLDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},
			"oidc_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("OIDC/OAuth application specific settings.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type associated with the application.").Description,
						Computed:    true,
					},
					"home_page_url": schema.StringAttribute{
						Description:         oidcHomePageURLDescription.Description,
						MarkdownDescription: oidcHomePageURLDescription.MarkdownDescription,
						Computed:            true,
					},
					"device_path_id": schema.StringAttribute{
						Description:         oidcOptionsDevicePathIdDescription.Description,
						MarkdownDescription: oidcOptionsDevicePathIdDescription.MarkdownDescription,
						Computed:            true,
					},
					"device_custom_verification_uri": schema.StringAttribute{
						Description:         oidcOptionsDeviceCustomVerificationUriDescription.Description,
						MarkdownDescription: oidcOptionsDeviceCustomVerificationUriDescription.MarkdownDescription,
						Computed:            true,
					},
					"device_timeout": schema.Int64Attribute{
						Description:         oidcOptionsDeviceTimeoutDescription.Description,
						MarkdownDescription: oidcOptionsDeviceTimeoutDescription.MarkdownDescription,
						Computed:            true,
					},
					"device_polling_interval": schema.Int64Attribute{
						Description:         oidcOptionsDevicePollingIntervalDescription.Description,
						MarkdownDescription: oidcOptionsDevicePollingIntervalDescription.MarkdownDescription,
						Computed:            true,
					},
					"initiate_login_uri": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the URI to use for third-parties to begin the sign-on process for the application.").Description,
						Computed:    true,
					},
					"jwks": schema.StringAttribute{
						Description:         oidcJwksDescription.Description,
						MarkdownDescription: oidcJwksDescription.MarkdownDescription,
						Computed:            true,
					},
					"jwks_url": schema.StringAttribute{
						Description:         oidcJwksUrlDescription.Description,
						MarkdownDescription: oidcJwksUrlDescription.MarkdownDescription,
						Computed:            true,
					},
					"target_link_uri": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The URI for the application.").Description,
						Computed:    true,
					},
					"grant_types": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A list that specifies the grant type for the authorization request.").Description,
						ElementType: types.StringType,
						Computed:    true,
					},
					"response_types": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A list that specifies the code or token type returned by an authorization request.").Description,
						ElementType: types.StringType,
						Computed:    true,
					},
					"token_endpoint_auth_method": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the client authentication methods supported by the token endpoint.").Description,
						Computed:    true,
					},
					"par_requirement": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies whether pushed authorization requests (PAR) are required.").Description,
						Computed:    true,
					},
					"par_timeout": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the pushed authorization request (PAR) timeout in seconds.").Description,
						Computed:    true,
					},
					"pkce_enforcement": schema.StringAttribute{
						Description:         oidcOptionsPKCEEnforcementDescription.Description,
						MarkdownDescription: oidcOptionsPKCEEnforcementDescription.MarkdownDescription,
						Computed:            true,
					},
					"redirect_uris": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of strings that specifies the allowed callback URIs for the authentication response.").Description,
						ElementType: types.StringType,
						Computed:    true,
					},
					"allow_wildcard_in_redirect_uris": schema.BoolAttribute{
						Description:         oidcAllowWildcardsInRedirectUrisDescription.Description,
						MarkdownDescription: oidcAllowWildcardsInRedirectUrisDescription.MarkdownDescription,
						Computed:            true,
					},
					"post_logout_redirect_uris": schema.SetAttribute{
						Description:         oidcPostLogoutRedirectUrisDescription.Description,
						MarkdownDescription: oidcPostLogoutRedirectUrisDescription.MarkdownDescription,
						ElementType:         types.StringType,
						Computed:            true,
					},
					"refresh_token_duration": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the lifetime in seconds of the refresh token.").Description,
						Computed:    true,
					},
					"refresh_token_rolling_duration": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required.").Description,
						Computed:    true,
					},
					"refresh_token_rolling_grace_period_duration": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens.").Description,
						Computed:    true,
					},
					"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
						Description:         oidcAdditionalRefreshTokenReplayProtectionEnabledDescription.Description,
						MarkdownDescription: oidcAdditionalRefreshTokenReplayProtectionEnabledDescription.MarkdownDescription,
						Computed:            true,
					},
					"certificate_based_authentication": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("Certificate based authentication settings.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"key_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents a PingOne ID for the issuance certificate key.").Description,
								Computed:    true,
							},
						},
					},
					"support_unsigned_request_object": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the request query parameter JWT is allowed to be unsigned.").Description,
						Computed:    true,
					},
					"require_signed_request_object": schema.BoolAttribute{
						Description:         oidcRequireSignedRequestObjectDescription.Description,
						MarkdownDescription: oidcRequireSignedRequestObjectDescription.MarkdownDescription,
						Computed:            true,
					},
					"cors_settings": datasourceApplicationSchemaCorsSettings(),
					"mobile_app": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("Mobile application integration settings.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"bundle_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the bundle associated with the application, for push notifications in native apps.").Description,
								Computed:    true,
							},
							"package_name": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the package name associated with the application, for push notifications in native apps.").Description,
								Computed:    true,
							},
							"huawei_app_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery.").Description,
								Computed:    true,
							},
							"huawei_package_name": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("The package name associated with the application, for push notifications in native apps.").Description,
								Computed:    true,
							},
							"passcode_refresh_seconds": schema.Int64Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("The amount of time a passcode should be displayed before being replaced with a new passcode.").Description,
								Computed:    true,
							},
							"universal_app_link": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code.").Description,
								Computed:    true,
							},
							"integrity_detection": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("Mobile application integrity detection settings.").Description,
								Computed:    true,

								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether device integrity detection takes place on mobile devices.").Description,
										Computed:    true,
									},
									"excluded_platforms": schema.SetAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates OS excluded from device integrity checking.").Description,
										ElementType: types.StringType,
										Computed:    true,
									},
									"cache_duration": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates the caching duration of successful integrity detection calls.").Description,
										Computed:    true,

										Attributes: map[string]schema.Attribute{
											"amount": schema.Int64Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.").Description,
												Computed:    true,
											},
											"units": schema.StringAttribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the cache duration time units.").Description,
												Computed:    true,
											},
										},
									},
									"google_play": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes Google Play Integrity API credential settings for Android device integrity detection.").Description,
										Computed:    true,

										Attributes: map[string]schema.Attribute{
											"decryption_key": schema.StringAttribute{
												Description:         googlePlayDecryptionKeyDescription.Description,
												MarkdownDescription: googlePlayDecryptionKeyDescription.MarkdownDescription,
												Computed:            true,
												Sensitive:           true,
											},
											"service_account_credentials_json": schema.StringAttribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("Contents of the JSON file that represents your Service Account Credentials.").Description,
												CustomType:  jsontypes.NormalizedType{},
												Computed:    true,
												Sensitive:   true,
											},
											"verification_key": schema.StringAttribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("Play Integrity verdict signature verification key from your Google Play Services account.").Description,
												Computed:    true,
												Sensitive:   true,
											},
											"verification_type": schema.StringAttribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("The type of verification.").Description,
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"saml_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("SAML application specific settings.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the custom home page URL for the application.").Description,
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type associated with the application.").Description,
						Computed:    true,
					},
					"acs_urls": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).").Description,
						ElementType: types.StringType,
						Computed:    true,
					},
					"assertion_duration": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the assertion validity duration in seconds.").Description,
						Computed:    true,
					},
					"assertion_signed_enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the SAML assertion itself should be signed.").Description,
						Computed:    true,
					},
					"default_target_url": schema.StringAttribute{
						Description:         samlDefaultTargetUrlDescription.Description,
						MarkdownDescription: samlDefaultTargetUrlDescription.MarkdownDescription,
						Computed:            true,
					},
					"idp_signing_key": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("SAML application assertion/response signing key settings.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the signature algorithm of the key.").Description,
								Computed:    true,
							},
							"key_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.").Description,
								Computed:    true,
							},
						},
					},
					"enable_requested_authn_context": schema.BoolAttribute{
						Description:         samlEnableRequestedAuthnContextDescription.Description,
						MarkdownDescription: samlEnableRequestedAuthnContextDescription.MarkdownDescription,
						Computed:            true,
					},
					"nameid_format": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the format of the Subject NameID attibute in the SAML assertion.").Description,
						Computed:    true,
					},
					"response_is_signed": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the SAML assertion response itself should be signed.").Description,
						Computed:    true,
					},
					"slo_binding": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the binding protocol to be used for the logout response.").Description,
						Computed:    true,
					},
					"slo_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the logout endpoint URL.").Description,
						Computed:    true,
					},
					"slo_response_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the endpoint URL to submit the logout response.").Description,
						Computed:    true,
					},
					"slo_window": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request.").Description,
						Computed:    true,
					},
					"sp_encryption": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings for PingOne to encrypt SAML assertions to be sent to the application. Assertions are not encrypted by default.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description:         samlSpEncryptionAlgorithmDescription.Description,
								MarkdownDescription: samlSpEncryptionAlgorithmDescription.MarkdownDescription,
								Computed:            true,
							},
							"certificate": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the certificate settings used to encrypt SAML assertions.").Description,
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique identifier of the encryption public certificate that has been uploaded to PingOne.").Description,
										Computed:    true,

										CustomType: pingonetypes.ResourceIDType{},
									},
								},
							},
						},
					},
					"sp_entity_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.").Description,
						Computed:    true,
					},
					"sp_verification": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies SP signature verification settings.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"authn_request_signed": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the Authn Request signing should be enforced.").Description,
								Computed:    true,
							},
							"certificate_ids": schema.SetAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A list that specifies the certificate IDs used to verify the service provider signature.").Description,
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
					"cors_settings": datasourceApplicationSchemaCorsSettings(),
				},
			},
		},
	}
}

func datasourceApplicationSchemaCorsSettings() schema.SingleNestedAttribute {

	listDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings.",
	)

	behaviorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS):       "rejects all CORS requests",
		string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS): "rejects all CORS requests except those listed in `origins`",
	})

	originsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value will be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Is expected to be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and is expected to be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.",
	)

	return schema.SingleNestedAttribute{
		Description:         listDescription.Description,
		MarkdownDescription: listDescription.MarkdownDescription,
		Computed:            true,

		Attributes: map[string]schema.Attribute{
			"behavior": schema.StringAttribute{
				Description:         behaviorDescription.Description,
				MarkdownDescription: behaviorDescription.MarkdownDescription,
				Computed:            true,
			},
			"origins": schema.SetAttribute{
				Description:         originsDescription.Description,
				MarkdownDescription: originsDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,
			},
		},
	}
}

func (r *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *applicationDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var application *management.ReadOneApplication200Response

	// Application API does not support SCIM filtering
	if !data.ApplicationId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneApplication",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&application,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		var entityArray *management.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadAllApplications(ctx, data.EnvironmentId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadAllApplications",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if applications, ok := entityArray.Embedded.GetApplicationsOk(); ok {
			found := false

			var applicationObj management.ReadOneApplication200Response
			for _, applicationObj = range applications {
				applicationInstance := applicationObj.GetActualInstance()

				applicationName := ""

				switch v := applicationInstance.(type) {
				case *management.ApplicationExternalLink:
					applicationName = v.GetName()

				case *management.ApplicationOIDC:
					applicationName = v.GetName()

				case *management.ApplicationSAML:
					applicationName = v.GetName()
				}

				if applicationName == data.Name.ValueString() {
					application = &applicationObj
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find the application from name or application is not the correct type",
					fmt.Sprintf("The application name %s for environment %s cannot be found, and only %s, %s or %s application types are retrievable", data.Name.String(), data.EnvironmentId.String(),
						string(management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT), string(management.ENUMAPPLICATIONPROTOCOL_SAML), string(management.ENUMAPPLICATIONPROTOCOL_EXTERNAL_LINK)),
				)
				return
			}

		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Application: application_id or name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, application)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *applicationDataSourceModel) toState(ctx context.Context, apiObject *management.ReadOneApplication200Response) diag.Diagnostics {
	var diags, d diag.Diagnostics
	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	applicationInstance := apiObject.GetActualInstance()
	switch v := applicationInstance.(type) {
	case *management.ApplicationExternalLink:
		p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.PingOneResourceIDOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
		if vA, ok := v.GetAccessControlOk(); ok {
			if vR, ok := vA.GetRoleOk(); ok {
				p.AccessControlRoleType = framework.EnumOkToTF(vR.GetTypeOk())
			}

			p.AccessControlGroupOptions, d = applicationAccessControlGroupOptionsToTF(vA.GetGroupOk())
			diags = append(diags, d...)
		}

		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		p.Icon, d = service.ImageOkToTF(v.GetIconOk())
		diags = append(diags, d...)

		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		// Service specific attributes
		p.Tags = types.SetNull(types.StringType)
		p.OIDCOptions = types.ObjectNull(applicationOidcOptionsTFObjectTypes)
		p.SAMLOptions = types.ObjectNull(applicationSamlOptionsTFObjectTypes)

		p.ExternalLinkOptions, d = applicationExternalLinkOptionsToTF(v)
		diags = append(diags, d...)

	case *management.ApplicationOIDC:
		p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.PingOneResourceIDOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
		if vA, ok := v.GetAccessControlOk(); ok {
			if vR, ok := vA.GetRoleOk(); ok {
				p.AccessControlRoleType = framework.EnumOkToTF(vR.GetTypeOk())
			}

			p.AccessControlGroupOptions, d = applicationAccessControlGroupOptionsToTF(vA.GetGroupOk())
			diags = append(diags, d...)
		}

		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		p.Icon, d = service.ImageOkToTF(v.GetIconOk())
		diags = append(diags, d...)

		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		// Service specific attributes
		p.Tags = framework.EnumSetOkToTF(v.GetTagsOk())

		var oidcOptionsState applicationOIDCOptionsResourceModelV1
		if !p.OIDCOptions.IsNull() && !p.OIDCOptions.IsUnknown() {
			d := p.OIDCOptions.As(ctx, &oidcOptionsState, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}
		}
		p.OIDCOptions, d = applicationOidcOptionsToTF(ctx, v, oidcOptionsState)
		diags = append(diags, d...)

		p.SAMLOptions = types.ObjectNull(applicationSamlOptionsTFObjectTypes)
		p.ExternalLinkOptions = types.ObjectNull(applicationExternalLinkOptionsTFObjectTypes)

	case *management.ApplicationSAML:
		p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.PingOneResourceIDOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
		if vA, ok := v.GetAccessControlOk(); ok {
			if vR, ok := vA.GetRoleOk(); ok {
				p.AccessControlRoleType = framework.EnumOkToTF(vR.GetTypeOk())
			}

			p.AccessControlGroupOptions, d = applicationAccessControlGroupOptionsToTF(vA.GetGroupOk())
			diags = append(diags, d...)
		}

		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		p.Icon, d = service.ImageOkToTF(v.GetIconOk())
		diags = append(diags, d...)

		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		// Service specific attributes
		p.Tags = types.SetNull(types.StringType)
		p.OIDCOptions = types.ObjectNull(applicationOidcOptionsTFObjectTypes)

		p.SAMLOptions, d = applicationSamlOptionsToTF(v)
		diags = append(diags, d...)

		p.ExternalLinkOptions = types.ObjectNull(applicationExternalLinkOptionsTFObjectTypes)
	}

	return diags
}
