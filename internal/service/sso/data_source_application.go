package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationDataSource serviceClientType

type applicationDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	EnvironmentId             types.String `tfsdk:"environment_id"`
	ApplicationId             types.String `tfsdk:"application_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	Tags                      types.Set    `tfsdk:"tags"`
	LoginPageUrl              types.String `tfsdk:"login_page_url"`
	Icon                      types.List   `tfsdk:"icon"`
	AccessControlRoleType     types.String `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.List   `tfsdk:"access_control_group_options"`
	HiddenFromAppPortal       types.Bool   `tfsdk:"hidden_from_app_portal"`
	ExternalLinkOptions       types.List   `tfsdk:"external_link_options"`
	OIDCOptions               types.List   `tfsdk:"oidc_options"`
	SAMLOptions               types.List   `tfsdk:"saml_options"`
}

var (
	applicationCorsSettingsTFObjectTypes = map[string]attr.Type{
		"behavior": types.StringType,
		"origins":  types.SetType{ElemType: types.StringType},
	}

	applicationOidcOptionsTFObjectTypes = map[string]attr.Type{
		"type":                             types.StringType,
		"home_page_url":                    types.StringType,
		"initiate_login_uri":               types.StringType,
		"jwks":                             types.StringType,
		"jwks_url":                         types.StringType,
		"target_link_uri":                  types.StringType,
		"grant_types":                      types.SetType{ElemType: types.StringType},
		"response_types":                   types.SetType{ElemType: types.StringType},
		"token_endpoint_authn_method":      types.StringType,
		"par_requirement":                  types.StringType,
		"par_timeout":                      types.Int64Type,
		"pkce_enforcement":                 types.StringType,
		"redirect_uris":                    types.SetType{ElemType: types.StringType},
		"allow_wildcards_in_redirect_uris": types.BoolType,
		"post_logout_redirect_uris":        types.SetType{ElemType: types.StringType},
		"refresh_token_duration":           types.Int64Type,
		"refresh_token_rolling_duration":   types.Int64Type,
		"refresh_token_rolling_grace_period_duration":        types.Int64Type,
		"additional_refresh_token_replay_protection_enabled": types.BoolType,
		"client_id":                        types.StringType,
		"client_secret":                    types.StringType,
		"certificate_based_authentication": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcOptionsCertificateAuthenticationTFObjectTypes}},
		"support_unsigned_request_object":  types.BoolType,
		"require_signed_request_object":    types.BoolType,
		"cors_settings":                    types.ListType{ElemType: types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}},
		"mobile_app":                       types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppTFObjectTypes}},
	}

	applicationOidcMobileAppTFObjectTypes = map[string]attr.Type{
		"bundle_id":                types.StringType,
		"package_name":             types.StringType,
		"huawei_app_id":            types.StringType,
		"huawei_package_name":      types.StringType,
		"passcode_refresh_seconds": types.Int64Type,
		"universal_app_link":       types.StringType,
		"integrity_detection":      types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionTFObjectTypes}},
	}

	applicationOidcMobileAppIntegrityDetectionTFObjectTypes = map[string]attr.Type{
		"enabled":            types.BoolType,
		"excluded_platforms": types.SetType{ElemType: types.StringType},
		"cache_duration":     types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes}},
		"google_play":        types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes}},
	}

	applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes = map[string]attr.Type{
		"amount": types.Int64Type,
		"units":  types.StringType,
	}

	applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes = map[string]attr.Type{
		"decryption_key":                   types.StringType,
		"service_account_credentials_json": types.StringType,
		"verification_key":                 types.StringType,
		"verification_type":                types.StringType,
	}

	applicationOidcOptionsCertificateAuthenticationTFObjectTypes = map[string]attr.Type{
		"key_id": types.StringType,
	}

	applicationSamlOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url":                   types.StringType,
		"type":                            types.StringType,
		"acs_urls":                        types.SetType{ElemType: types.StringType},
		"assertion_duration":              types.Int64Type,
		"assertion_signed_enabled":        types.BoolType,
		"default_target_url":              types.StringType,
		"idp_signing_key":                 types.ListType{ElemType: types.ObjectType{AttrTypes: applicationSamlOptionsIdpSigningKeyTFObjectTypes}},
		"enable_requested_authn_context":  types.BoolType,
		"nameid_format":                   types.StringType,
		"response_is_signed":              types.BoolType,
		"slo_binding":                     types.StringType,
		"slo_endpoint":                    types.StringType,
		"slo_response_endpoint":           types.StringType,
		"slo_window":                      types.Int64Type,
		"sp_entity_id":                    types.StringType,
		"sp_verification_certificate_ids": types.SetType{ElemType: types.StringType},
		"sp_verification":                 types.ListType{ElemType: types.ObjectType{AttrTypes: applicationSamlOptionsSpVerificationTFObjectTypes}},
		"cors_settings":                   types.ListType{ElemType: types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}},
	}

	applicationSamlOptionsIdpSigningKeyTFObjectTypes = map[string]attr.Type{
		"algorithm": types.StringType,
		"key_id":    types.StringType,
	}

	applicationSamlOptionsSpVerificationTFObjectTypes = map[string]attr.Type{
		"certificate_ids":      types.SetType{ElemType: types.StringType},
		"authn_request_signed": types.BoolType,
	}

	applicationExternalLinkOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url": types.StringType,
	}

	applicationIconTFObjectTypes = map[string]attr.Type{
		"id":   types.StringType,
		"href": types.StringType,
	}

	applicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"groups": types.SetType{ElemType: types.StringType},
	}
)

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

	oidcJwksDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	)

	oidcJwksUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
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

	samlDefaultTargetUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).",
	)

	samlSpVerificationCertificateIds := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release.  Please use the `sp_verification.certificate_ids` attribute going forward.  A list that specifies the certificate IDs used to verify the service provider signature.",
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
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
					validation.P1ResourceIDValidator(),
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
			"icon": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The HREF and the ID for the application icon.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
			},
			"access_control_role_type": schema.StringAttribute{
				Description:         accessControlRoleTypeDescription.Description,
				MarkdownDescription: accessControlRoleTypeDescription.MarkdownDescription,
				Computed:            true,
			},
			"access_control_group_options": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Group access control settings.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
			},
			"hidden_from_app_portal": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.").Description,
				Computed:    true,
			},
			"external_link_options": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("External link application specific settings.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"home_page_url": schema.StringAttribute{
							Description:         externalLinkHomePageURLDescription.Description,
							MarkdownDescription: externalLinkHomePageURLDescription.MarkdownDescription,
							Computed:            true,
						},
					},
				},
			},
			"oidc_options": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("OIDC/OAuth application specific settings.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
						"token_endpoint_authn_method": schema.StringAttribute{
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
						"allow_wildcards_in_redirect_uris": schema.BoolAttribute{
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
						"client_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID used to authenticate to the authorization server.").Description,
							Computed:    true,
						},
						"client_secret": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application secret ID used to authenticate to the authorization server.").Description,
							Computed:    true,
							Sensitive:   true,
						},
						"certificate_based_authentication": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Certificate based authentication settings.").Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key_id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents a PingOne ID for the issuance certificate key.").Description,
										Computed:    true,
									},
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
						"mobile_app": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Mobile application integration settings.").Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
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
									"integrity_detection": schema.ListNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Mobile application integrity detection settings.").Description,
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
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
												"cache_duration": schema.ListNestedAttribute{
													Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates the caching duration of successful integrity detection calls.").Description,
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
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
												},
												"google_play": schema.ListNestedAttribute{
													Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that describes Google Play Integrity API credential settings for Android device integrity detection.").Description,
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"decryption_key": schema.StringAttribute{
																Description:         googlePlayDecryptionKeyDescription.Description,
																MarkdownDescription: googlePlayDecryptionKeyDescription.MarkdownDescription,
																Computed:            true,
																Sensitive:           true,
															},
															"service_account_credentials_json": schema.StringAttribute{
																Description: framework.SchemaAttributeDescriptionFromMarkdown("Contents of the JSON file that represents your Service Account Credentials.").Description,
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
						},
					},
				},
			},
			"saml_options": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("SAML application specific settings.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
						"idp_signing_key": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("SAML application assertion/response signing key settings.").Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
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
						"sp_entity_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.").Description,
							Computed:    true,
						},
						"sp_verification_certificate_ids": schema.SetAttribute{
							Description:         samlSpVerificationCertificateIds.Description,
							MarkdownDescription: samlSpVerificationCertificateIds.MarkdownDescription,
							ElementType:         types.StringType,
							DeprecationMessage:  "The `sp_verification_certificate_ids` attribute is deprecated and will be removed in the next major release.  Please use the `sp_verification.certificate_ids` attribute going forward.",
							Computed:            true,
						},
						"sp_verification": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single list item that specifies SP signature verification settings.").Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
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
						},
						"cors_settings": datasourceApplicationSchemaCorsSettings(),
					},
				},
			},
		},
	}
}

func datasourceApplicationSchemaCorsSettings() schema.ListNestedAttribute {

	listDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings.",
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

	return schema.ListNestedAttribute{
		Description:         listDescription.Description,
		MarkdownDescription: listDescription.MarkdownDescription,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
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

	if r.Client.ManagementAPIClient == nil {
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

	// Secondary call required to obtain kerberos secret for OIDC applications
	var secret *management.ApplicationSecret
	if application.ApplicationOIDC != nil && application.ApplicationOIDC.GetId() != "" {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, *application.ApplicationOIDC.GetEnvironment().Id, application.ApplicationOIDC.GetId()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			},
			"ReadApplicationSecret",
			framework.DefaultCustomError,
			applicationOIDCSecretDataSourceRetryConditions,
			&secret,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(application, secret)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *applicationDataSourceModel) toState(apiObject *management.ReadOneApplication200Response, secret *management.ApplicationSecret) diag.Diagnostics {
	var diags diag.Diagnostics
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
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.iconOkToTF(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.accessControlOkToTF(v.GetAccessControlOk())
		diags.Append(d...)

		p.ExternalLinkOptions, d = p.toStateExternalLinkOptions(v)
		diags.Append(d...)

	case *management.ApplicationOIDC:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.Tags = framework.EnumSetOkToTF(v.GetTagsOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.iconOkToTF(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.accessControlOkToTF(v.GetAccessControlOk())
		diags.Append(d...)

		p.OIDCOptions, d = p.toStateOIDCOptions(v, secret)
		diags.Append(d...)

	case *management.ApplicationSAML:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.iconOkToTF(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.accessControlOkToTF(v.GetAccessControlOk())
		diags.Append(d...)

		p.SAMLOptions, d = p.toStateSAMLOptions(v)
		diags.Append(d...)
	}

	return diags
}

func (p *applicationDataSourceModel) toStateExternalLinkOptions(apiObject *management.ApplicationExternalLink) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: applicationExternalLinkOptionsTFObjectTypes}

	if apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	externalLinkOptions := map[string]attr.Value{
		"home_page_url": framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationExternalLinkOptionsTFObjectTypes, externalLinkOptions)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) iconOkToTF(apiObject *management.ApplicationIcon, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: applicationIconTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	applicationIcon := map[string]attr.Value{
		"id":   framework.StringOkToTF(apiObject.GetIdOk()),
		"href": framework.StringOkToTF(apiObject.GetHrefOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationIconTFObjectTypes, applicationIcon)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) accessControlOkToTF(apiObject *management.ApplicationAccessControl, ok bool) (basetypes.StringValue, basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objAccessControlGroupOptions := types.ListNull(types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes})
	accessControlRoleType := types.StringNull()
	accessControlGroupOption := []attr.Value{}

	if !ok || apiObject == nil {
		return accessControlRoleType, objAccessControlGroupOptions, diags
	}

	accessControlRoleType = framework.EnumOkToTF(apiObject.Role.GetTypeOk())

	if group, ok := apiObject.GetGroupOk(); ok {
		var d diag.Diagnostics
		groups := make([]string, 0)
		for _, v := range group.GetGroups() {
			groups = append(groups, v.GetId())
		}

		groupObj := map[string]attr.Value{
			"type":   framework.EnumOkToTF(group.GetTypeOk()),
			"groups": framework.StringSetToTF(groups),
		}

		groupsObj, d := types.ObjectValue(applicationAccessControlGroupOptionsTFObjectTypes, groupObj)
		diags.Append(d...)

		accessControlGroupOption = append(accessControlGroupOption, groupsObj)

		objAccessControlGroupOptions, d = types.ListValue(types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes}, accessControlGroupOption)
		diags.Append(d...)
	}

	return accessControlRoleType, objAccessControlGroupOptions, diags

}

func (p *applicationDataSourceModel) toStateOIDCOptions(apiObject *management.ApplicationOIDC, secret *management.ApplicationSecret) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	objOIDCOptions := []attr.Value{}

	if apiObject == nil {
		return types.ListNull(types.ObjectType{AttrTypes: applicationOidcOptionsTFObjectTypes}), diags
	}

	kerberoObj, d := p.applicationOidcKerberosOkToTF(apiObject.GetKerberosOk())
	diags.Append(d...)

	// CORS Settings
	corsSettingsObj := types.ListNull(types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes})
	if v, ok := apiObject.GetCorsSettingsOk(); ok {
		corsSettings := map[string]attr.Value{
			"behavior": framework.EnumOkToTF(v.GetBehaviorOk()),
			"origins":  framework.StringSetOkToTF(v.GetOriginsOk()),
		}

		corsSettingsFlattenedObj, d := types.ObjectValue(applicationCorsSettingsTFObjectTypes, corsSettings)
		diags.Append(d...)

		corsSettingsObj, d = types.ListValue(types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}, append([]attr.Value{}, corsSettingsFlattenedObj))
		diags.Append(d...)
	}

	mobileObj, d := p.applicationOidcMobileOkToTF(apiObject.GetMobileOk())
	diags.Append(d...)

	// Build OIDC Options
	oidcOptions := map[string]attr.Value{
		"client_id":                        framework.StringOkToTF(apiObject.GetIdOk()),
		"type":                             framework.EnumOkToTF(apiObject.GetTypeOk()),
		"grant_types":                      framework.EnumSetOkToTF(apiObject.GetGrantTypesOk()),
		"token_endpoint_authn_method":      framework.EnumOkToTF(apiObject.GetTokenEndpointAuthMethodOk()),
		"par_requirement":                  framework.EnumOkToTF(apiObject.GetParRequirementOk()),
		"par_timeout":                      framework.Int32OkToTF(apiObject.GetParTimeoutOk()),
		"home_page_url":                    framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"initiate_login_uri":               framework.StringOkToTF(apiObject.GetInitiateLoginUriOk()),
		"jwks":                             framework.StringOkToTF(apiObject.GetJwksOk()),
		"jwks_url":                         framework.StringOkToTF(apiObject.GetJwksUrlOk()),
		"target_link_uri":                  framework.StringOkToTF(apiObject.GetTargetLinkUriOk()),
		"response_types":                   framework.EnumSetOkToTF(apiObject.GetResponseTypesOk()),
		"pkce_enforcement":                 framework.EnumOkToTF(apiObject.GetPkceEnforcementOk()),
		"redirect_uris":                    framework.StringSetOkToTF(apiObject.GetRedirectUrisOk()),
		"allow_wildcards_in_redirect_uris": framework.BoolOkToTF(apiObject.GetAllowWildcardInRedirectUrisOk()),
		"post_logout_redirect_uris":        framework.StringSetOkToTF(apiObject.GetPostLogoutRedirectUrisOk()),
		"refresh_token_duration":           framework.Int32OkToTF(apiObject.GetRefreshTokenDurationOk()),
		"refresh_token_rolling_duration":   framework.Int32OkToTF(apiObject.GetRefreshTokenRollingDurationOk()),
		"refresh_token_rolling_grace_period_duration":        framework.Int32OkToTF(apiObject.GetRefreshTokenRollingGracePeriodDurationOk()),
		"additional_refresh_token_replay_protection_enabled": framework.BoolOkToTF(apiObject.GetAdditionalRefreshTokenReplayProtectionEnabledOk()),
		"client_secret":                    framework.StringOkToTF(secret.GetSecretOk()),
		"certificate_based_authentication": kerberoObj,
		"support_unsigned_request_object":  framework.BoolOkToTF(apiObject.GetSupportUnsignedRequestObjectOk()),
		"require_signed_request_object":    framework.BoolOkToTF(apiObject.GetRequireSignedRequestObjectOk()),
		"cors_settings":                    corsSettingsObj,
		"mobile_app":                       mobileObj,
	}

	oidcObject, d := types.ObjectValue(applicationOidcOptionsTFObjectTypes, oidcOptions)
	diags.Append(d...)

	objOIDCOptions = append(objOIDCOptions, oidcObject)

	returnVar, d := types.ListValue(types.ObjectType{AttrTypes: applicationOidcOptionsTFObjectTypes}, objOIDCOptions)
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) applicationOidcKerberosOkToTF(apiObject *management.ApplicationOIDCAllOfKerberos, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: applicationOidcOptionsCertificateAuthenticationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	kerberos := map[string]attr.Value{}
	if v, ok := apiObject.GetKeyOk(); ok {
		kerberos["key_id"] = framework.StringOkToTF(v.GetIdOk())
	} else {
		kerberos["key_id"] = types.StringNull()
	}

	flattenedObj, d := types.ObjectValue(applicationOidcOptionsCertificateAuthenticationTFObjectTypes, kerberos)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) applicationOidcMobileOkToTF(apiObject *management.ApplicationOIDCAllOfMobile, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	passcodeRefreshSeconds := types.Int64Null()
	if v, ok := apiObject.GetPasscodeRefreshDurationOk(); ok {
		passcodeRefreshSeconds = framework.Int32OkToTF(v.GetDurationOk())

		if j, okJ := v.GetTimeUnitOk(); okJ && *j != management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS {
			diags.AddError(
				"Invalid time unit",
				fmt.Sprintf("Expecting time unit of %s for attribute `passcode_refresh_seconds`, got %v", management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS, j),
			)

			return types.ListNull(tfObjType), diags
		}
	}

	integrityDetection, d := p.applicationOidcMobileIntegrityDetectionOkToTF(apiObject.GetIntegrityDetectionOk())
	diags.Append(d...)

	mobileApp := map[string]attr.Value{
		"bundle_id":                framework.StringOkToTF(apiObject.GetBundleIdOk()),
		"package_name":             framework.StringOkToTF(apiObject.GetPackageNameOk()),
		"huawei_app_id":            framework.StringOkToTF(apiObject.GetHuaweiAppIdOk()),
		"huawei_package_name":      framework.StringOkToTF(apiObject.GetHuaweiPackageNameOk()),
		"passcode_refresh_seconds": passcodeRefreshSeconds,
		"universal_app_link":       framework.StringOkToTF(apiObject.GetUriPrefixOk()),
		"integrity_detection":      integrityDetection,
	}

	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppTFObjectTypes, mobileApp)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) applicationOidcMobileIntegrityDetectionOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetection, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	// Cache Duration
	cacheDuration := map[string]attr.Value{}
	if v, ok := apiObject.GetCacheDurationOk(); ok {
		cacheDuration = map[string]attr.Value{
			"amount": framework.Int32OkToTF(v.GetAmountOk()),
			"units":  framework.EnumOkToTF(v.GetUnitsOk()),
		}
	}
	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes, cacheDuration)
	diags.Append(d...)

	cacheDurationObj, d := types.ListValue(
		types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes},
		append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	// Google Play
	dummySuppressValue := "DUMMY_SUPPRESS_VALUE"
	googlePlay := map[string]attr.Value{}
	if v, ok := apiObject.GetGooglePlayOk(); ok {
		googlePlay["verification_type"] = framework.EnumOkToTF(v.GetVerificationTypeOk())

		if v.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL {
			googlePlay["decryption_key"] = framework.StringToTF(dummySuppressValue)
			googlePlay["verification_key"] = framework.StringToTF(dummySuppressValue)
		} else {
			googlePlay["decryption_key"] = types.StringNull()
			googlePlay["verification_key"] = types.StringNull()
		}

		if v.GetVerificationType() == management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE {
			googlePlay["service_account_credentials_json"] = framework.StringToTF(dummySuppressValue)
		} else {
			googlePlay["service_account_credentials_json"] = types.StringNull()
		}
	}

	flattenedObj, d = types.ObjectValue(applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes, googlePlay)
	diags.Append(d...)

	googlePlayObj, d := types.ListValue(
		types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes},
		append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	// Enabled
	enabledBool := false
	if v, ok := apiObject.GetModeOk(); ok {
		if *v == management.ENUMENABLEDSTATUS_ENABLED {
			enabledBool = true
		}
	}

	// Build Main Object
	integrityDetection := map[string]attr.Value{
		"enabled":            framework.BoolOkToTF(&enabledBool, ok),
		"excluded_platforms": framework.EnumSetOkToTF(apiObject.GetExcludedPlatformsOk()),
		"cache_duration":     cacheDurationObj,
		"google_play":        googlePlayObj,
	}

	flattenedObj, d = types.ObjectValue(applicationOidcMobileAppIntegrityDetectionTFObjectTypes, integrityDetection)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) toStateSAMLOptions(apiObject *management.ApplicationSAML) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	objSAMLOptions := []attr.Value{}

	if apiObject == nil {
		return types.ListNull(types.ObjectType{AttrTypes: applicationSamlOptionsTFObjectTypes}), diags
	}

	// IdP Signing Key
	idpSigningKey := map[string]attr.Value{}
	if v, ok := apiObject.GetIdpSigningOk(); ok {
		idpSigningKey["algorithm"] = framework.EnumOkToTF(v.GetAlgorithmOk())

		if j, ok := v.GetKeyOk(); ok {
			idpSigningKey["key_id"] = framework.StringOkToTF(j.GetIdOk())
		}
	} else {
		idpSigningKey["algorithm"] = types.StringNull()
		idpSigningKey["key_id"] = types.StringNull()
	}
	flattenedObj, d := types.ObjectValue(applicationSamlOptionsIdpSigningKeyTFObjectTypes, idpSigningKey)
	diags.Append(d...)

	idpSigningKeyObj, d := types.ListValue(types.ObjectType{AttrTypes: applicationSamlOptionsIdpSigningKeyTFObjectTypes}, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	// SP Verification
	var idList []string
	spVerification := map[string]attr.Value{}
	if v, ok := apiObject.GetSpVerificationOk(); ok {
		spVerification["authn_request_signed"] = framework.BoolOkToTF(v.GetAuthnRequestSignedOk())

		if v1, ok := v.GetCertificatesOk(); ok {

			idList = make([]string, 0)
			for _, j := range v1 {
				idList = append(idList, j.GetId())
			}
		}

		spVerification["certificate_ids"] = framework.StringSetToTF(idList)
	} else {
		spVerification["authn_request_signed"] = types.BoolNull()
		spVerification["certificate_ids"] = types.SetNull(types.StringType)
	}
	spVerificationFlattenedObj, d := types.ObjectValue(applicationSamlOptionsSpVerificationTFObjectTypes, spVerification)
	diags.Append(d...)

	spVerificationObj, d := types.ListValue(types.ObjectType{AttrTypes: applicationSamlOptionsSpVerificationTFObjectTypes}, append([]attr.Value{}, spVerificationFlattenedObj))
	diags.Append(d...)

	// CORS Settings
	corsSettingsObj := types.ListNull(types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes})
	if v, ok := apiObject.GetCorsSettingsOk(); ok {
		corsSettings := map[string]attr.Value{
			"behavior": framework.EnumOkToTF(v.GetBehaviorOk()),
			"origins":  framework.StringSetOkToTF(v.GetOriginsOk()),
		}

		corsSettingsFlattenedObj, d := types.ObjectValue(applicationCorsSettingsTFObjectTypes, corsSettings)
		diags.Append(d...)

		corsSettingsObj, d = types.ListValue(types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}, append([]attr.Value{}, corsSettingsFlattenedObj))
		diags.Append(d...)
	}

	// Build Main Object
	samlOptions := map[string]attr.Value{
		"type":                            framework.EnumOkToTF(apiObject.GetTypeOk()),
		"acs_urls":                        framework.EnumSetOkToTF(apiObject.GetAcsUrlsOk()),
		"assertion_duration":              framework.Int32OkToTF(apiObject.GetAssertionDurationOk()),
		"sp_entity_id":                    framework.StringOkToTF(apiObject.GetSpEntityIdOk()),
		"home_page_url":                   framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"assertion_signed_enabled":        framework.BoolOkToTF(apiObject.GetAssertionSignedOk()),
		"default_target_url":              framework.StringOkToTF(apiObject.GetDefaultTargetUrlOk()),
		"idp_signing_key":                 idpSigningKeyObj,
		"enable_requested_authn_context":  framework.BoolOkToTF(apiObject.GetEnableRequestedAuthnContextOk()),
		"nameid_format":                   framework.StringOkToTF(apiObject.GetNameIdFormatOk()),
		"response_is_signed":              framework.BoolOkToTF(apiObject.GetResponseSignedOk()),
		"slo_binding":                     framework.EnumOkToTF(apiObject.GetSloBindingOk()),
		"slo_endpoint":                    framework.StringOkToTF(apiObject.GetSloEndpointOk()),
		"slo_response_endpoint":           framework.StringOkToTF(apiObject.GetSloResponseEndpointOk()),
		"slo_window":                      framework.Int32OkToTF(apiObject.GetSloWindowOk()),
		"sp_verification_certificate_ids": framework.StringSetToTF(idList),
		"sp_verification":                 spVerificationObj,
		"cors_settings":                   corsSettingsObj,
	}

	samlObject, d := types.ObjectValue(applicationSamlOptionsTFObjectTypes, samlOptions)
	diags.Append(d...)

	objSAMLOptions = append(objSAMLOptions, samlObject)

	returnVar, d := types.ListValue(types.ObjectType{AttrTypes: applicationSamlOptionsTFObjectTypes}, objSAMLOptions)
	diags.Append(d...)

	return returnVar, diags
}

func applicationOIDCSecretDataSourceRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	// The secret may take a short time to propagate
	if r.StatusCode == 404 {
		tflog.Warn(ctx, "Application secret not found, available for retry")
		return true
	}

	if p1error != nil {

		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

	}

	return false
}
