package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResource serviceClientType

type ApplicationResourceModel struct {
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
		"type":                                        types.StringType,
		"home_page_url":                               types.StringType,
		"initiate_login_uri":                          types.StringType,
		"target_link_uri":                             types.StringType,
		"grant_types":                                 types.SetType{ElemType: types.StringType},
		"response_types":                              types.SetType{ElemType: types.StringType},
		"token_endpoint_authn_method":                 types.StringType,
		"par_requirement":                             types.StringType,
		"par_timeout":                                 types.Int64Type,
		"pkce_enforcement":                            types.StringType,
		"redirect_uris":                               types.SetType{ElemType: types.StringType},
		"allow_wildcards_in_redirect_uris":            types.BoolType,
		"post_logout_redirect_uris":                   types.SetType{ElemType: types.StringType},
		"refresh_token_duration":                      types.Int64Type,
		"refresh_token_rolling_duration":              types.Int64Type,
		"refresh_token_rolling_grace_period_duration": types.Int64Type,
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
	_ resource.Resource = &ApplicationResource{}
)

// New Object
func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

// Metadata
func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the application.",
	)

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the application is enabled in the environment.",
	).DefaultValue(false)

	tagsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array of strings that specifies the list of labels associated with the application.",
	).AllowedValuesEnum(management.AllowedEnumApplicationTagsEnumValues).ConflictsWith([]string{"external_link_options", "saml_options"})

	loginPageUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom login page URL for the application. If you set the `login_page_url` property for applications in an environment that sets a custom domain, the URL should include the top-level domain and at least one additional domain level. **Warning** To avoid issues with third-party cookies in some browsers, a custom domain must be used, giving your PingOne environment the same parent domain as your authentication application. For more information about custom domains, see Custom domains.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	accessControlRoleTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user role required to access the application.  A user is an admin user if the user has one or more admin roles assigned, such as `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.",
	).AllowedValuesEnum(management.AllowedEnumApplicationAccessControlTypeEnumValues)

	accessControlGroupOptionsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the group type required to access the application.",
	).AllowedValuesComplex(map[string]string{
		"ANY_GROUP":  "the actor must belong to at least one group listed in the `groups` property",
		"ALL_GROUPS": "the actor must belong to all groups listed in the `groups` property",
	})

	hiddenFromAppPortalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.",
	).DefaultValue(false)

	iconHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the URL for the application icon.  Both `http://` and `https://` are permitted.",
	)

	appTypesExactlyOneOf := []string{"external_link_options", "oidc_options", "saml_options"}

	externalLinkOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies External link application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf)

	externalLinkHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.",
	)

	oidcOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies OIDC/OAuth application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf)

	oidcOptionsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type associated with the application.",
	).AllowedValues(
		string(management.ENUMAPPLICATIONTYPE_WEB_APP),
		string(management.ENUMAPPLICATIONTYPE_NATIVE_APP),
		string(management.ENUMAPPLICATIONTYPE_SINGLE_PAGE_APP),
		string(management.ENUMAPPLICATIONTYPE_WORKER),
		string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP),
		string(management.ENUMAPPLICATIONTYPE_SERVICE),
	).RequiresReplace()

	oidcHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcInitiateLoginUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the URI to use for third-parties to begin the sign-on process for the application. If specified, PingOne redirects users to this URI to initiate SSO to PingOne. The application is responsible for implementing the relevant OIDC flow when the initiate login URI is requested. This property is required if you want the application to appear in the PingOne Application Portal. See the OIDC specification section of [Initiating Login from a Third Party](https://openid.net/specs/openid-connect-core-1_0.html#ThirdPartyInitiatedLogin) for more information.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcTargetLinkUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URI for the application. If specified, PingOne will redirect application users to this URI after a user is authenticated. In the PingOne admin console, this becomes the value of the `target_link_uri` parameter used for the Initiate Single Sign-On URL field.  Both `http://` and `https://` URLs are permitted as well as custom mobile native schema (e.g., `org.bxretail.app://target`).",
	)

	oidcGrantTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the grant type for the authorization request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCGrantTypeEnumValues)

	oidcResponseTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the code or token type returned by an authorization request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCResponseTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("Note that `%s` cannot be used in an authorization request with `%s` or `%s` because PingOne does not currently support OIDC hybrid flows.", string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_CODE), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_TOKEN), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_ID_TOKEN)),
	)

	oidcTokenEndpointAuthnMethod := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the client authentication methods supported by the token endpoint.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCTokenAuthMethodEnumValues)

	oidcParRequirementDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies whether pushed authorization requests (PAR) are required.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCPARRequirementEnumValues).DefaultValue(string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL))

	oidcParTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the pushed authorization request (PAR) timeout in seconds.  If a value is not provided, the default value is `60`.  Valid values are between `1` and `600`.",
	).DefaultValue(60)

	oidcOptionsPKCEEnforcementDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how `PKCE` request parameters are handled on the authorize request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCPKCEOptionEnumValues).DefaultValue(string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL))

	oidcRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of strings that specifies the allowed callback URIs for the authentication response.    The provided URLs are expected to use the `https://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  The `http` schema is only permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcAllowWildcardsInRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).",
	).DefaultValue(false)

	oidcPostLogoutRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).",
	)

	oidcRefreshTokenDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the lifetime in seconds of the refresh token. If a value is not provided, the default value is `2592000`, or 30 days. Valid values are between `60` and `2147483647`. If the `refresh_token_rolling_duration` property is specified for the application, then this property value must be less than or equal to the value of `refresh_token_rolling_duration`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.",
	).DefaultValue(2592000)

	oidcRefreshTokenRollingDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required. If a value is not provided, the default value is `15552000`, or 180 days. Valid values are between `60` and `2147483647`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.",
	).DefaultValue(15552000)

	oidcRefreshTokenRollingGracePeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens. This is useful in the case of network errors on the client. Valid values are between `0` and `86400` seconds. `Null` is treated the same as `0`.",
	)

	oidcAdditionalRefreshTokenReplayProtectionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true` (the default), if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens. Setting this to null equates to a `false` setting.",
	).DefaultValue(true)

	oidcSupportUnsignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If `false` or null, an unsigned request object is not allowed.",
	).DefaultValue(false)

	oidcRequireSignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
	).DefaultValue(false)

	oidcBundleIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead. A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
	)

	oidcPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead. A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
	)

	oidcCertificateBasedAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single block that specifies Certificate based authentication settings. This parameter block can only be set where the application's `type` parameter is set to `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
	)

	oidcCertificateBasedAuthenticationKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents a PingOne ID for the issuance certificate key.  The key must be of type `ISSUANCE`.  Must be a valid PingOne Resource ID.",
	)

	oidcMobileAppDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single block that specifies Mobile application integration settings for `%s` type applications.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
	)

	oidcMobileAppBundleIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable.",
	).RequiresReplace()

	oidcMobileAppPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.",
	).RequiresReplace()

	oidcMobileAppHuaweiAppIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_package_name`.",
	).RequiresReplace()

	oidcMobileAppHuaweiPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The package name associated with the application, for push notifications in native apps. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_app_id`.",
	).RequiresReplace()

	oidcMobileAppPasscodeRefreshSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The amount of time a passcode should be displayed before being replaced with a new passcode - must be between `30` and `60` seconds.",
	).DefaultValue(30)

	oidcMobileAppUniversalAppLinkDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code. The URI prefix can be set to a universal link with a valid value (which can be a URL address that starts with `HTTP://` or `HTTPS://`, such as `https://www.bxretail.org`), or an app schema, which is just a string and requires no special validation.",
	)


	googlePlayDecryptionKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.",
	)

	samlOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies SAML application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf)

	samlEnableRequestedAuthnContextDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
	)

	samlSpVerificationCertificateIds := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release.  Please use the `sp_verification.certificate_ids` attribute going forward.  A list that specifies the certificate IDs used to verify the service provider signature.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne application (SAML, OpenID Connect, External Link) in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The PingOne resource ID of the environment to create and manage the application in."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the application.").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"tags": schema.SetAttribute{
				Description:         tagsDescription.Description,
				MarkdownDescription: tagsDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationTagsEnumValues)...),
					),
					setvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
					),
				},
			},

			"login_page_url": schema.StringAttribute{
				Description:         loginPageUrlDescription.Description,
				MarkdownDescription: loginPageUrlDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`),
						"Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\".",
					),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("external_link_options"),
					),
				},
			},

			"access_control_role_type": schema.StringAttribute{
				Description:         accessControlRoleTypeDescription.Description,
				MarkdownDescription: accessControlRoleTypeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationAccessControlTypeEnumValues)...),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("external_link_options"),
					),
				},
			},

			"hidden_from_app_portal": schema.BoolAttribute{
				Description:         hiddenFromAppPortalDescription.Description,
				MarkdownDescription: hiddenFromAppPortalDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},
		},

		Blocks: map[string]schema.Block{
			"icon": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies settings for the application icon.").Description,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID for the application icon.  Must be a valid PingOne Resource ID.").Description,
							Required:    true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
						},

						"href": schema.StringAttribute{
							Description:         iconHrefDescription.Description,
							MarkdownDescription: iconHrefDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Value must be a valid URL with `http://` or `https://` prefix."),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"access_control_group_options": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies group access control settings.").Description,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         accessControlGroupOptionsTypeDescription.Description,
							MarkdownDescription: accessControlGroupOptionsTypeDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationAccessControlGroupTypeEnumValues)...),
							},
						},

						"groups": schema.SetAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A set that specifies the group IDs for the groups the actor must belong to for access to the application.  Values must be valid PingOne Resource IDs.").Description,
							Required:    true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									verify.P1ResourceIDValidator(),
								),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"external_link_options": schema.ListNestedBlock{
				Description:         externalLinkOptionsDescription.Description,
				MarkdownDescription: externalLinkOptionsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"home_page_url": schema.StringAttribute{
							Description:         externalLinkHomePageURLDescription.Description,
							MarkdownDescription: externalLinkHomePageURLDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Value must be a valid URL with `http://` or `https://` prefix."),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
					),
				},
			},

			"oidc_options": schema.ListNestedBlock{
				Description:         oidcOptionsDescription.Description,
				MarkdownDescription: oidcOptionsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         oidcOptionsTypeDescription.Description,
							MarkdownDescription: oidcOptionsTypeDescription.MarkdownDescription,
							Required:            true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},

							Validators: []validator.String{
								stringvalidator.OneOf(
									string(management.ENUMAPPLICATIONTYPE_WEB_APP),
									string(management.ENUMAPPLICATIONTYPE_NATIVE_APP),
									string(management.ENUMAPPLICATIONTYPE_SINGLE_PAGE_APP),
									string(management.ENUMAPPLICATIONTYPE_WORKER),
									string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP),
									string(management.ENUMAPPLICATIONTYPE_SERVICE),
								),
							},
						},

						"home_page_url": schema.StringAttribute{
							Description:         oidcHomePageURLDescription.Description,
							MarkdownDescription: oidcHomePageURLDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
							},
						},

						"initiate_login_uri": schema.StringAttribute{
							Description:         oidcInitiateLoginUriDescription.Description,
							MarkdownDescription: oidcInitiateLoginUriDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
							},
						},

						"target_link_uri": schema.StringAttribute{
							Description:         oidcTargetLinkUriDescription.Description,
							MarkdownDescription: oidcTargetLinkUriDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://target`)."),
							},
						},

						"grant_types": schema.SetAttribute{
							Description:         oidcGrantTypesDescription.Description,
							MarkdownDescription: oidcGrantTypesDescription.MarkdownDescription,
							Required:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCGrantTypeEnumValues)...),
								),
							},
						},

						"response_types": schema.SetAttribute{
							Description:         oidcResponseTypesDescription.Description,
							MarkdownDescription: oidcResponseTypesDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCResponseTypeEnumValues)...),
								),
							},
						},

						"token_endpoint_authn_method": schema.StringAttribute{
							Description:         oidcTokenEndpointAuthnMethod.Description,
							MarkdownDescription: oidcTokenEndpointAuthnMethod.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCTokenAuthMethodEnumValues)...),
							},
						},

						"par_requirement": schema.StringAttribute{
							Description:         oidcParRequirementDescription.Description,
							MarkdownDescription: oidcParRequirementDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCPARRequirementEnumValues)...),
							},
						},

						"par_timeout": schema.Int64Attribute{
							Description:         oidcParTimeoutDescription.Description,
							MarkdownDescription: oidcParTimeoutDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(60),

							Validators: []validator.Int64{
								int64validator.Between(0, 600),
							},
						},

						"pkce_enforcement": schema.StringAttribute{
							Description:         oidcOptionsPKCEEnforcementDescription.Description,
							MarkdownDescription: oidcOptionsPKCEEnforcementDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCPKCEOptionEnumValues)...),
							},
						},

						"redirect_uris": schema.SetAttribute{
							Description:         oidcRedirectUrisDescription.Description,
							MarkdownDescription: oidcRedirectUrisDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(\S+:\/\/).+)`), "Expected value to have a url with schema of \"https\" or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
								),
							},
						},

						"allow_wildcards_in_redirect_uris": schema.BoolAttribute{
							Description:         oidcAllowWildcardsInRedirectUrisDescription.Description,
							MarkdownDescription: oidcAllowWildcardsInRedirectUrisDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"post_logout_redirect_uris": schema.SetAttribute{
							Description:         oidcPostLogoutRedirectUrisDescription.Description,
							MarkdownDescription: oidcPostLogoutRedirectUrisDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://logout`)."),
								),
							},
						},

						"refresh_token_duration": schema.Int64Attribute{
							Description:         oidcRefreshTokenDurationDescription.Description,
							MarkdownDescription: oidcRefreshTokenDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(2592000),

							Validators: []validator.Int64{
								int64validator.Between(60, 2147483647),
							},
						},

						"refresh_token_rolling_duration": schema.Int64Attribute{
							Description:         oidcRefreshTokenRollingDurationDescription.Description,
							MarkdownDescription: oidcRefreshTokenRollingDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(15552000),

							Validators: []validator.Int64{
								int64validator.Between(60, 2147483647),
							},
						},

						"refresh_token_rolling_grace_period_duration": schema.Int64Attribute{
							Description:         oidcRefreshTokenRollingGracePeriodDurationDescription.Description,
							MarkdownDescription: oidcRefreshTokenRollingGracePeriodDurationDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.Int64{
								int64validator.Between(60, 86400),
							},
						},

						"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
							Description:         oidcAdditionalRefreshTokenReplayProtectionEnabledDescription.Description,
							MarkdownDescription: oidcAdditionalRefreshTokenReplayProtectionEnabledDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(true),
						},

						"client_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID used to authenticate to the authorization server.").Description,
							Computed:    true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},

						"client_secret": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application secret ID used to authenticate to the authorization server.").Description,
							Computed:    true,
							Sensitive:   true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},

						"support_unsigned_request_object": schema.BoolAttribute{
							Description:         oidcSupportUnsignedRequestObjectDescription.Description,
							MarkdownDescription: oidcSupportUnsignedRequestObjectDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"require_signed_request_object": schema.BoolAttribute{
							Description:         oidcRequireSignedRequestObjectDescription.Description,
							MarkdownDescription: oidcRequireSignedRequestObjectDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"bundle_id": schema.StringAttribute{
							Description:         oidcBundleIdDescription.Description,
							MarkdownDescription: oidcBundleIdDescription.MarkdownDescription,
							DeprecationMessage:  "This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead.",
							Optional:            true,
							Computed:            true,

							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},

						"package_name": schema.StringAttribute{
							Description:         oidcPackageNameDescription.Description,
							MarkdownDescription: oidcPackageNameDescription.MarkdownDescription,
							DeprecationMessage:  "This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead.",
							Optional:            true,
							Computed:            true,

							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},
					},

					Blocks: map[string]schema.Block{
						"certificate_based_authentication": schema.ListNestedBlock{
							Description:         oidcCertificateBasedAuthenticationDescription.Description,
							MarkdownDescription: oidcCertificateBasedAuthenticationDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"key_id": schema.StringAttribute{
										Description:         oidcCertificateBasedAuthenticationKeyIdDescription.Description,
										MarkdownDescription: oidcCertificateBasedAuthenticationKeyIdDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.String{
											verify.P1ResourceIDValidator(),
										},
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},

						"mobile_app": schema.ListNestedBlock{
							Description:         oidcMobileAppDescription.Description,
							MarkdownDescription: oidcMobileAppDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"bundle_id": schema.StringAttribute{
										Description:         oidcMobileAppBundleIdDescription.Description,
										MarkdownDescription: oidcMobileAppBundleIdDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},

									"package_name": schema.StringAttribute{
										Description:         oidcMobileAppPackageNameDescription.Description,
										MarkdownDescription: oidcMobileAppPackageNameDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},

									"huawei_app_id": schema.StringAttribute{
										Description:         oidcMobileAppHuaweiAppIdDescription.Description,
										MarkdownDescription: oidcMobileAppHuaweiAppIdDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
											stringvalidator.AlsoRequires(
												path.MatchRelative().AtParent().AtName("huawei_app_id"),
												path.MatchRelative().AtParent().AtName("huawei_package_name"),
											),
										},
									},

									"huawei_package_name": schema.StringAttribute{
										Description:         oidcMobileAppHuaweiPackageNameDescription.Description,
										MarkdownDescription: oidcMobileAppHuaweiPackageNameDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
											stringvalidator.AlsoRequires(
												path.MatchRelative().AtParent().AtName("huawei_app_id"),
												path.MatchRelative().AtParent().AtName("huawei_package_name"),
											),
										},
									},

									"passcode_refresh_seconds": schema.Int64Attribute{
										Description:         oidcMobileAppPasscodeRefreshSecondsDescription.Description,
										MarkdownDescription: oidcMobileAppPasscodeRefreshSecondsDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int64default.StaticInt64(30),

										Validators: []validator.Int64{
											int64validator.Between(30, 60),
										},
									},

									"universal_app_link": schema.StringAttribute{
										Description:         oidcMobileAppUniversalAppLinkDescription.Description,
										MarkdownDescription: oidcMobileAppUniversalAppLinkDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
								},

								Blocks: map[string]schema.Block{
									"integrity_detection": schema.ListNestedBlock{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies mobile application integrity detection settings.").Description,

										NestedObject: schema.NestedBlockObject{
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

					Blocks: map[string]schema.Block{
						"cors_settings": resourceApplicationSchemaCorsSettings(),
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
					),
				},
			},

			"saml_options": schema.ListNestedBlock{
				Description:         samlOptionsDescription.Description,
				MarkdownDescription: samlOptionsDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{
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
							Computed:            true,

							ElementType:        types.StringType,
							DeprecationMessage: "The `sp_verification_certificate_ids` attribute is deprecated and will be removed in the next major release.  Please use the `sp_verification.certificate_ids` attribute going forward.",
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
					},

					Blocks: map[string]schema.Block{
						"cors_settings": resourceApplicationSchemaCorsSettings(),
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
					),
				},
			},
		},
	}
}

func resourceApplicationSchemaCorsSettings() schema.ListNestedBlock {

	listDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings.",
	)

	behaviorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS):       "rejects all CORS requests",
		string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS): "rejects all CORS requests except those listed in `origins`",
	})

	originsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value must be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Must be non-empty when `behavior` is `%s` and must be omitted or empty when `behavior` is `%s`.  Limited to 20 values.", string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS), string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS)),
	)

	return schema.ListNestedBlock{
		Description:         listDescription.Description,
		MarkdownDescription: listDescription.MarkdownDescription,

		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"behavior": schema.StringAttribute{
					Description:         behaviorDescription.Description,
					MarkdownDescription: behaviorDescription.MarkdownDescription,
					Required:            true,

					Validators: []validator.String{
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationCorsSettingsBehaviorEnumValues)...),
					},
				},

				"origins": schema.SetAttribute{
					Description:         originsDescription.Description,
					MarkdownDescription: originsDescription.MarkdownDescription,
					Optional:            true,

					ElementType: types.StringType,

					Validators: []validator.Set{
						setvalidator.SizeAtMost(20),
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(https?:\/\/)?(localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|([\*a-zA-Z0-9-]+\.)+[a-zA-Z]{2,})(:\d{1,5})?$`),
								"Expected value to be a URL (with schema of \"http\" or \"https\") without a path.  Subdomains may use a wildcard to match any string",
							),
						),
					},
				},
			},
		},

		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	application, _ := plan.expand()

	// Run the API call
	var response *management.CreateApplication201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.CreateApplication(ctx, plan.EnvironmentId.ValueString()).CreateApplicationRequest(*application).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretResponse *management.ApplicationSecret
	if response.ApplicationOIDC != nil && response.ApplicationOIDC.GetId() != "" {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadApplicationSecret",
			framework.DefaultCustomError,
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
			&secretResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(nil, response, secretResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var secretResponse *management.ApplicationSecret
	if response.ApplicationOIDC != nil && response.ApplicationOIDC.GetId() != "" {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadApplicationSecret",
			framework.CustomErrorResourceNotFoundWarning,
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
			&secretResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response, nil, secretResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	_, application := plan.expand()

	// Run the API call
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateApplicationRequest(*application).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretResponse *management.ApplicationSecret
	if response.ApplicationOIDC != nil && response.ApplicationOIDC.GetId() != "" {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadApplicationSecret",
			framework.DefaultCustomError,
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
			&secretResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response, nil, secretResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.DeleteApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteApplication",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "application_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *ApplicationResourceModel) expand() (*management.CreateApplicationRequest, *management.UpdateApplicationRequest) {

	data := management.NewApplication(p.Name.ValueString())

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.PopulationId.IsNull() && !p.PopulationId.IsUnknown() {
		data.SetPopulation(
			*management.NewApplicationPopulation(p.PopulationId.ValueString()),
		)
	}

	if !p.UserFilter.IsNull() && !p.UserFilter.IsUnknown() {
		data.SetUserFilter(p.UserFilter.ValueString())
	}

	if !p.ExternalId.IsNull() && !p.ExternalId.IsUnknown() {
		data.SetExternalId(p.ExternalId.ValueString())
	}

	return data
}

func (p *ApplicationResourceModel) toState(apiObject *management.ReadOneApplication200Response, apiObjectCreate *management.CreateApplication201Response, apiSecretObject *management.ApplicationSecret) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	if v, ok := apiObject.GetPopulationOk(); ok {
		p.PopulationId = framework.StringOkToTF(v.GetIdOk())
	} else {
		p.PopulationId = types.StringNull()
	}

	p.UserFilter = framework.StringOkToTF(apiObject.GetUserFilterOk())
	p.ExternalId = framework.StringOkToTF(apiObject.GetExternalIdOk())

	return diags
}
