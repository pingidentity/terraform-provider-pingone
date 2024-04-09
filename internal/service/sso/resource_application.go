package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

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
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResource serviceClientType

type ApplicationResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	EnvironmentId             types.String `tfsdk:"environment_id"`
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

type ApplicationAccessControlGroupOptionsResourceModel struct {
	Type   types.String `tfsdk:"type"`
	Groups types.Set    `tfsdk:"groups"`
}

type ApplicationExternalLinkOptionsResourceModel struct {
	HomePageUrl types.String `tfsdk:"home_page_url"`
}

type ApplicationOIDCOptionsResourceModel struct {
	AdditionalRefreshTokenReplayProtectionEnabled types.Bool   `tfsdk:"additional_refresh_token_replay_protection_enabled"`
	AllowWildcardsInRedirectUris                  types.Bool   `tfsdk:"allow_wildcards_in_redirect_uris"`
	BundleId                                      types.String `tfsdk:"bundle_id"`
	CertificateBasedAuthentication                types.List   `tfsdk:"certificate_based_authentication"`
	ClientId                                      types.String `tfsdk:"client_id"`
	ClientSecret                                  types.String `tfsdk:"client_secret"`
	CorsSettings                                  types.List   `tfsdk:"cors_settings"`
	GrantTypes                                    types.Set    `tfsdk:"grant_types"`
	HomePageUrl                                   types.String `tfsdk:"home_page_url"`
	InitiateLoginUri                              types.String `tfsdk:"initiate_login_uri"`
	MobileApp                                     types.List   `tfsdk:"mobile_app"`
	PackageName                                   types.String `tfsdk:"package_name"`
	ParRequirement                                types.String `tfsdk:"par_requirement"`
	ParTimeout                                    types.Int64  `tfsdk:"par_timeout"`
	PKCEEnforcement                               types.String `tfsdk:"pkce_enforcement"`
	PostLogoutRedirectUris                        types.Set    `tfsdk:"post_logout_redirect_uris"`
	RedirectUris                                  types.Set    `tfsdk:"redirect_uris"`
	RefreshTokenDuration                          types.Int64  `tfsdk:"refresh_token_duration"`
	RefreshTokenRollingDuration                   types.Int64  `tfsdk:"refresh_token_rolling_duration"`
	RefreshTokenRollingGracePeriodDuration        types.Int64  `tfsdk:"refresh_token_rolling_grace_period_duration"`
	RequireSignedRequestObject                    types.Bool   `tfsdk:"require_signed_request_object"`
	ResponseTypes                                 types.Set    `tfsdk:"response_types"`
	SupportUnsignedRequestObject                  types.Bool   `tfsdk:"support_unsigned_request_object"`
	Jwks                                          types.String `tfsdk:"jwks"`
	JwksUrl                                       types.String `tfsdk:"jwks_url"`
	TargetLinkUri                                 types.String `tfsdk:"target_link_uri"`
	TokenEndpointAuthnMethod                      types.String `tfsdk:"token_endpoint_authn_method"`
	Type                                          types.String `tfsdk:"type"`
}

type ApplicationCorsSettingsResourceModel struct {
	Behavior types.String `tfsdk:"behavior"`
	Origins  types.Set    `tfsdk:"origins"`
}

type ApplicationOIDCCertificateBasedAuthenticationResourceModel struct {
	KeyId types.String `tfsdk:"key_id"`
}

type ApplicationOIDCMobileAppResourceModel struct {
	BundleId               types.String `tfsdk:"bundle_id"`
	HuaweiAppId            types.String `tfsdk:"huawei_app_id"`
	HuaweiPackageName      types.String `tfsdk:"huawei_package_name"`
	IntegrityDetection     types.List   `tfsdk:"integrity_detection"`
	PackageName            types.String `tfsdk:"package_name"`
	PasscodeRefreshSeconds types.Int64  `tfsdk:"passcode_refresh_seconds"`
	UniversalAppLink       types.String `tfsdk:"universal_app_link"`
}

type ApplicationOIDCMobileAppIntegrityDetectionResourceModel struct {
	CacheDuration     types.List `tfsdk:"cache_duration"`
	Enabled           types.Bool `tfsdk:"enabled"`
	ExcludedPlatforms types.Set  `tfsdk:"excluded_platforms"`
	GooglePlay        types.List `tfsdk:"google_play"`
}

type ApplicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModel struct {
	Amount types.Int64  `tfsdk:"amount"`
	Units  types.String `tfsdk:"units"`
}

type ApplicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModel struct {
	DecryptionKey                 types.String `tfsdk:"decryption_key"`
	ServiceAccountCredentialsJson types.String `tfsdk:"service_account_credentials_json"`
	VerificationKey               types.String `tfsdk:"verification_key"`
	VerificationType              types.String `tfsdk:"verification_type"`
}

type ApplicationSAMLOptionsResourceModel struct {
	AcsUrls                      types.Set    `tfsdk:"acs_urls"`
	AssertionDuration            types.Int64  `tfsdk:"assertion_duration"`
	AssertionSignedEnabled       types.Bool   `tfsdk:"assertion_signed_enabled"`
	CorsSettings                 types.List   `tfsdk:"cors_settings"`
	EnableRequestedAuthnContext  types.Bool   `tfsdk:"enable_requested_authn_context"`
	HomePageUrl                  types.String `tfsdk:"home_page_url"`
	IdpSigningKey                types.List   `tfsdk:"idp_signing_key"`
	IdpSigningKeyId              types.String `tfsdk:"idp_signing_key_id"`
	DefaultTargetUrl             types.String `tfsdk:"default_target_url"`
	NameIdFormat                 types.String `tfsdk:"nameid_format"`
	ResponseIsSigned             types.Bool   `tfsdk:"response_is_signed"`
	SloBinding                   types.String `tfsdk:"slo_binding"`
	SloEndpoint                  types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint          types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                    types.Int64  `tfsdk:"slo_window"`
	SpEntityId                   types.String `tfsdk:"sp_entity_id"`
	SpVerification               types.List   `tfsdk:"sp_verification"`
	SpVerificationCertificateIds types.Set    `tfsdk:"sp_verification_certificate_ids"`
	Type                         types.String `tfsdk:"type"`
}

type ApplicationSAMLOptionsIdpSigningKeyResourceModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
	KeyId     types.String `tfsdk:"key_id"`
}

type ApplicationSAMLOptionsSpVerificationResourceModel struct {
	CertificateIds     types.Set  `tfsdk:"certificate_ids"`
	AuthnRequestSigned types.Bool `tfsdk:"authn_request_signed"`
}

var (
	applicationCorsSettingsTFObjectTypes = map[string]attr.Type{
		"behavior": types.StringType,
		"origins":  types.SetType{ElemType: types.StringType},
	}

	applicationOidcOptionsTFObjectTypes = map[string]attr.Type{
		"additional_refresh_token_replay_protection_enabled": types.BoolType,
		"allow_wildcards_in_redirect_uris":                   types.BoolType,
		"bundle_id":                                          types.StringType,
		"certificate_based_authentication":                   types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcOptionsCertificateAuthenticationTFObjectTypes}},
		"client_id":                                          types.StringType,
		"client_secret":                                      types.StringType,
		"cors_settings":                                      types.ListType{ElemType: types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}},
		"grant_types":                                        types.SetType{ElemType: types.StringType},
		"home_page_url":                                      types.StringType,
		"initiate_login_uri":                                 types.StringType,
		"mobile_app":                                         types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppTFObjectTypes}},
		"package_name":                                       types.StringType,
		"par_requirement":                                    types.StringType,
		"par_timeout":                                        types.Int64Type,
		"pkce_enforcement":                                   types.StringType,
		"post_logout_redirect_uris":                          types.SetType{ElemType: types.StringType},
		"redirect_uris":                                      types.SetType{ElemType: types.StringType},
		"refresh_token_duration":                             types.Int64Type,
		"refresh_token_rolling_duration":                     types.Int64Type,
		"refresh_token_rolling_grace_period_duration":        types.Int64Type,
		"require_signed_request_object":                      types.BoolType,
		"response_types":                                     types.SetType{ElemType: types.StringType},
		"support_unsigned_request_object":                    types.BoolType,
		"jwks":                                               types.StringType,
		"jwks_url":                                           types.StringType,
		"target_link_uri":                                    types.StringType,
		"token_endpoint_authn_method":                        types.StringType,
		"type":                                               types.StringType,
	}

	applicationOidcMobileAppTFObjectTypes = map[string]attr.Type{
		"bundle_id":                types.StringType,
		"huawei_app_id":            types.StringType,
		"huawei_package_name":      types.StringType,
		"integrity_detection":      types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionTFObjectTypes}},
		"package_name":             types.StringType,
		"passcode_refresh_seconds": types.Int64Type,
		"universal_app_link":       types.StringType,
	}

	applicationOidcMobileAppIntegrityDetectionTFObjectTypes = map[string]attr.Type{
		"cache_duration":     types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes}},
		"enabled":            types.BoolType,
		"excluded_platforms": types.SetType{ElemType: types.StringType},
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
		"acs_urls":                        types.SetType{ElemType: types.StringType},
		"assertion_duration":              types.Int64Type,
		"assertion_signed_enabled":        types.BoolType,
		"cors_settings":                   types.ListType{ElemType: types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}},
		"enable_requested_authn_context":  types.BoolType,
		"home_page_url":                   types.StringType,
		"idp_signing_key_id":              types.StringType,
		"idp_signing_key":                 types.ListType{ElemType: types.ObjectType{AttrTypes: applicationSamlOptionsIdpSigningKeyTFObjectTypes}},
		"default_target_url":              types.StringType,
		"nameid_format":                   types.StringType,
		"response_is_signed":              types.BoolType,
		"slo_binding":                     types.StringType,
		"slo_endpoint":                    types.StringType,
		"slo_response_endpoint":           types.StringType,
		"slo_window":                      types.Int64Type,
		"sp_entity_id":                    types.StringType,
		"sp_verification_certificate_ids": types.SetType{ElemType: types.StringType},
		"sp_verification":                 types.ListType{ElemType: types.ObjectType{AttrTypes: applicationSamlOptionsSpVerificationTFObjectTypes}},
		"type":                            types.StringType,
	}

	applicationSamlOptionsIdpSigningKeyTFObjectTypes = map[string]attr.Type{
		"algorithm": types.StringType,
		"key_id":    types.StringType,
	}

	applicationSamlOptionsSpVerificationTFObjectTypes = map[string]attr.Type{
		"authn_request_signed": types.BoolType,
		"certificate_ids":      types.SetType{ElemType: types.StringType},
	}

	applicationExternalLinkOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url": types.StringType,
	}

	applicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"groups": types.SetType{ElemType: types.StringType},
		"type":   types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                   = &ApplicationResource{}
	_ resource.ResourceWithConfigure      = &ApplicationResource{}
	_ resource.ResourceWithImportState    = &ApplicationResource{}
	_ resource.ResourceWithValidateConfig = &ApplicationResource{}
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

	externalLinkOptionsHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
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

	oidcOptionsHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcOptionsInitiateLoginUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the URI to use for third-parties to begin the sign-on process for the application. If specified, PingOne redirects users to this URI to initiate SSO to PingOne. The application is responsible for implementing the relevant OIDC flow when the initiate login URI is requested. This property is required if you want the application to appear in the PingOne Application Portal. See the OIDC specification section of [Initiating Login from a Third Party](https://openid.net/specs/openid-connect-core-1_0.html#ThirdPartyInitiatedLogin) for more information.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcOptionsJwksDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	).ConflictsWith([]string{"jwks_url"})

	oidcOptionsJwksUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	).ConflictsWith([]string{"jwks"})

	oidcOptionsTargetLinkUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URI for the application. If specified, PingOne will redirect application users to this URI after a user is authenticated. In the PingOne admin console, this becomes the value of the `target_link_uri` parameter used for the Initiate Single Sign-On URL field.  Both `http://` and `https://` URLs are permitted as well as custom mobile native schema (e.g., `org.bxretail.app://target`).",
	)

	oidcOptionsGrantTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the grant type for the authorization request.",
	).AllowedValuesEnum([]string{
		string(management.ENUMAPPLICATIONOIDCGRANTTYPE_AUTHORIZATION_CODE),
		string(management.ENUMAPPLICATIONOIDCGRANTTYPE_IMPLICIT),
		string(management.ENUMAPPLICATIONOIDCGRANTTYPE_REFRESH_TOKEN),
		string(management.ENUMAPPLICATIONOIDCGRANTTYPE_CLIENT_CREDENTIALS),
	})

	oidcOptionsResponseTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the code or token type returned by an authorization request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCResponseTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("Note that `%s` cannot be used in an authorization request with `%s` or `%s` because PingOne does not currently support OIDC hybrid flows.", string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_CODE), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_TOKEN), string(management.ENUMAPPLICATIONOIDCRESPONSETYPE_ID_TOKEN)),
	)

	oidcOptionsTokenEndpointAuthnMethod := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the client authentication methods supported by the token endpoint.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCTokenAuthMethodEnumValues).AppendMarkdownString(fmt.Sprintf("When `%s` is configured, either `jwks` or `jwks_url` must also be configured.", string(management.ENUMAPPLICATIONOIDCTOKENAUTHMETHOD_PRIVATE_KEY_JWT)))

	oidcOptionsParRequirementDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies whether pushed authorization requests (PAR) are required.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCPARRequirementEnumValues).DefaultValue(string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL))

	const oidcOptionsParTimeoutDefault = 60
	const oidcOptionsParTimeoutMin = 1
	const oidcOptionsParTimeoutMax = 600
	oidcOptionsParTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the pushed authorization request (PAR) timeout in seconds.  Valid values are between `%d` and `%d`.", oidcOptionsParTimeoutMin, oidcOptionsParTimeoutMax),
	).DefaultValue(oidcOptionsParTimeoutDefault)

	oidcOptionsPKCEEnforcementDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how `PKCE` request parameters are handled on the authorize request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCPKCEOptionEnumValues).DefaultValue(string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL))

	oidcOptionsRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of strings that specifies the allowed callback URIs for the authentication response.    The provided URLs are expected to use the `https://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  The `http` schema is only permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcOptionsAllowWildcardsInRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).",
	).DefaultValue(false)

	oidcOptionsPostLogoutRedirectUrisDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).",
	)

	const oidcOptionsRefreshTokenDurationDefault = 2592000
	const oidcOptionsRefreshTokenDurationMin = 60
	const oidcOptionsRefreshTokenDurationMax = 2147483647
	oidcOptionsRefreshTokenDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the lifetime in seconds of the refresh token. Valid values are between `%d` and `%d`. If the `refresh_token_rolling_duration` property is specified for the application, then this property value must be less than or equal to the value of `refresh_token_rolling_duration`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.", oidcOptionsRefreshTokenDurationMin, oidcOptionsRefreshTokenDurationMax),
	).DefaultValue(oidcOptionsRefreshTokenDurationDefault)

	const oidcOptionsRefreshTokenRollingDurationDefault = 15552000
	const oidcOptionsRefreshTokenRollingDurationMin = 60
	const oidcOptionsRefreshTokenRollingDurationMax = 2147483647
	oidcOptionsRefreshTokenRollingDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required. Valid values are between `%d` and `%d`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.", oidcOptionsRefreshTokenRollingDurationMin, oidcOptionsRefreshTokenRollingDurationMax),
	).DefaultValue(oidcOptionsRefreshTokenRollingDurationDefault)

	const oidcOptionsRefreshTokenRollingGracePeriodDurationMin = 0
	const oidcOptionsRefreshTokenRollingGracePeriodDurationMax = 86400
	oidcOptionsRefreshTokenRollingGracePeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens. This is useful in the case of network errors on the client. Valid values are between `%d` and `%d` seconds. `Null` is treated the same as `0`.", oidcOptionsRefreshTokenRollingGracePeriodDurationMin, oidcOptionsRefreshTokenRollingGracePeriodDurationMax),
	)

	oidcOptionsAdditionalRefreshTokenReplayProtectionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true` (the default), if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens. Setting this to null equates to a `false` setting.",
	).DefaultValue(true)

	oidcOptionsClientSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the application secret ID used to authenticate to the authorization server.\n\n~> The `client_secret` cannot be rotated in this resource.  The `pingone_application_secret` resource should be used to control rotation of the `client_secret` value.  If using the `pingone_application_secret` resource, use of this attribute is likely to conflict with that resource.  In this case, the `pingone_application_secret.secret` attribute should be used instead.",
	)

	oidcOptionsSupportUnsignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If `false` or null, an unsigned request object is not allowed.",
	).DefaultValue(false)

	oidcOptionsRequireSignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
	).DefaultValue(false)

	oidcOptionsBundleIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead. A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
	)

	oidcOptionsPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead. A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable; any change will force recreation of the application resource.",
	)

	oidcOptionsCertificateBasedAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single block that specifies Certificate based authentication settings. This parameter block can only be set where the application's `type` parameter is set to `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
	)

	oidcOptionsCertificateBasedAuthenticationKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents a PingOne ID for the issuance certificate key.  The key must be of type `ISSUANCE`.  Must be a valid PingOne Resource ID.",
	)

	oidcOptionsMobileAppDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single block that specifies Mobile application integration settings for `%s` type applications.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
	)

	oidcOptionsMobileAppBundleIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable.",
	).RequiresReplace()

	oidcOptionsMobileAppPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.",
	).RequiresReplace()

	oidcOptionsMobileAppHuaweiAppIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_package_name`.",
	).RequiresReplace()

	oidcOptionsMobileAppHuaweiPackageNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The package name associated with the application, for push notifications in native apps. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_app_id`.",
	).RequiresReplace()

	const oidcOptionsMobileAppPasscodeRefreshSecondsDefault = 30
	const oidcOptionsMobileAppPasscodeRefreshSecondsMin = 30
	const oidcOptionsMobileAppPasscodeRefreshSecondsMax = 60
	oidcOptionsMobileAppPasscodeRefreshSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The amount of time a passcode should be displayed before being replaced with a new passcode - must be between `%d` and `%d` seconds.", oidcOptionsMobileAppPasscodeRefreshSecondsMin, oidcOptionsMobileAppPasscodeRefreshSecondsMax),
	).DefaultValue(oidcOptionsMobileAppPasscodeRefreshSecondsDefault)

	oidcOptionsMobileAppUniversalAppLinkDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code. The URI prefix can be set to a universal link with a valid value (which can be a URL address that starts with `HTTP://` or `HTTPS://`, such as `https://www.bxretail.org`), or an app schema, which is just a string and requires no special validation.",
	)

	oidcOptionsMobileAppIntegrityDetectionEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether device integrity detection takes place on mobile devices.",
	).DefaultValue(false)

	oidcOptionsMobileAppIntegrityDetectionExcludedPlatformsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("You can enable device integrity checking separately for Android and iOS by setting `enabled` to `true` and then using `excluded_platforms` to specify the OS where you do not want to use device integrity checking. The values to use are `%s` and `%s` (all upper case). Note that this is implemented as an array even though currently you can only include a single value.  If `%s` is not included in this list, the `google_play` attribute block must be configured.", string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE), string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_IOS), string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE)),
	)

	oidcOptionsMobileAppIntegrityDetectionCacheDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies settings for the caching duration of successful integrity detection calls.  Every attestation request entails a certain time tradeoff. You can choose to cache successful integrity detection calls for a predefined duration, between a minimum of 1 minute and a maximum of 48 hours. If integrity detection is ENABLED, the cache duration must be set.",
	)

	oidcOptionsMobileAppIntegrityDetectionCacheDurationUnitsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time units of the cache `amount` parameter.",
	).AllowedValuesEnum(management.AllowedEnumDurationUnitMinsHoursEnumValues).DefaultValue(string(management.ENUMDURATIONUNITMINSHOURS_MINUTES))

	oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single block that describes Google Play Integrity API credential settings for Android device integrity detection.  Required when `excluded_platforms` is unset or does not include `%s`.", management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE),
	)

	oidcOptionsMobileAppIntegrityDetectionGooglePlayDecryptionKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `%s`.", string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL)),
	).ConflictsWith([]string{"service_account_credentials_json"})

	oidcOptionsMobileAppIntegrityDetectionGooglePlayServiceAccountCredentialsJsonDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Contents of the JSON file that represents your Service Account Credentials. This parameter must be provided if you have set `verification_type` to `%s`.", string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE)),
	).ConflictsWith([]string{"decryption_key", "verification_key"})

	oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Play Integrity verdict signature verification key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `%s`.", string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL)),
	).ConflictsWith([]string{"service_account_credentials_json"})

	oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of verification that should be used.",
	).AllowedValuesEnum(management.AllowedEnumApplicationNativeGooglePlayVerificationTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("Using internal verification will not count against your Google API call quota. The value you select for this attribute determines what other parameters you must provide. When set to `%s`, you must provide `service_account_credentials_json`. When set to `%s`, you must provide both `decryption_key` and `verification_key`.", string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_GOOGLE), string(management.ENUMAPPLICATIONNATIVEGOOGLEPLAYVERIFICATIONTYPE_INTERNAL)),
	)

	samlOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single block that specifies SAML application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf)

	samlOptionsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type associated with the application.",
	).AllowedValues(
		string(management.ENUMAPPLICATIONTYPE_WEB_APP),
		string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP),
	).DefaultValue(string(management.ENUMAPPLICATIONTYPE_WEB_APP)).RequiresReplace()

	samlOptionsAssertionSignedEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML assertion itself should be signed.",
	).DefaultValue(true)

	samlOptionsIdpSigningKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release.  Please use the `idp_signing_key` block going forward.  An ID for the certificate key pair to be used by the identity provider to sign assertions and responses. If this property is omitted, the default signing certificate for the environment is used.",
	)

	samlOptionsDefaultTargetUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).",
	)

	samlOptionsEnableRequestedAuthnContextDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
	)

	samlOptionsResponseIsSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML assertion response itself should be signed.",
	).DefaultValue(false)

	samlOptionsSloBindingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the binding protocol to be used for the logout response.",
	).AllowedValuesEnum(management.AllowedEnumApplicationSAMLSloBindingEnumValues).AppendMarkdownString(
		fmt.Sprintf("Existing configurations with no data default to `%s`.", string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST)),
	).DefaultValue(string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST))

	samlOptionsSloResponseEndpointDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the endpoint URL to submit the logout response. If a value is not provided, the `slo_endpoint` property value is used to submit SLO response.",
	)

	const samlOptionsSloWindowMin = 0
	const samlOptionsSloWindowMax = 24
	samlOptionsSloWindowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request.  The minimum value is `%d` hour and the maximum is `%d` hours.", samlOptionsSloWindowMin, samlOptionsSloWindowMax),
	)

	samlOptionsSpVerificationCertificateIds := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Notice** This field is deprecated and will be removed in a future release.  Please use the `sp_verification.certificate_ids` parameter going forward.  A list that specifies the certificate IDs used to verify the service provider signature.",
	)

	samlOptionsIdpSigningKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"SAML application assertion/response signing key settings.  Use with `assertion_signed_enabled` to enable assertion signing and/or `response_is_signed` to enable response signing.  It's highly recommended, and best practice, to define signing key settings for the configured SAML application.  However if this property is omitted, the default signing certificate for the environment is used.  This parameter will become a required field in the next major release of the provider.",
	)

	samlOptionsIdpSigningKeyAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
	)

	samlOptionsSpVerificationAuthnRequestSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the Authn Request signing should be enforced.",
	).DefaultValue(false)

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
					setvalidator.SizeAtLeast(attrMinLength),
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
							Description:         externalLinkOptionsHomePageURLDescription.Description,
							MarkdownDescription: externalLinkOptionsHomePageURLDescription.MarkdownDescription,
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
							Description:         oidcOptionsHomePageURLDescription.Description,
							MarkdownDescription: oidcOptionsHomePageURLDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
							},
						},

						"initiate_login_uri": schema.StringAttribute{
							Description:         oidcOptionsInitiateLoginUriDescription.Description,
							MarkdownDescription: oidcOptionsInitiateLoginUriDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(https:\/\/).*)`), "Expected value to have a url with schema of \"https\".  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
							},
						},

						"jwks": schema.StringAttribute{
							Description:         oidcOptionsJwksDescription.Description,
							MarkdownDescription: oidcOptionsJwksDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("jwks_url"),
									path.MatchRelative().AtParent().AtName("jwks"),
								),
							},
						},

						"jwks_url": schema.StringAttribute{
							Description:         oidcOptionsJwksUrlDescription.Description,
							MarkdownDescription: oidcOptionsJwksUrlDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("jwks_url"),
									path.MatchRelative().AtParent().AtName("jwks"),
								),
								stringvalidator.RegexMatches(regexp.MustCompile(`^(https:\/\/).*`), "Expected value to have a url with schema of \"https\"."),
							},
						},

						"target_link_uri": schema.StringAttribute{
							Description:         oidcOptionsTargetLinkUriDescription.Description,
							MarkdownDescription: oidcOptionsTargetLinkUriDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://target`)."),
							},
						},

						"grant_types": schema.SetAttribute{
							Description:         oidcOptionsGrantTypesDescription.Description,
							MarkdownDescription: oidcOptionsGrantTypesDescription.MarkdownDescription,
							Required:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCGrantTypeEnumValues)...),
								),
							},
						},

						"response_types": schema.SetAttribute{
							Description:         oidcOptionsResponseTypesDescription.Description,
							MarkdownDescription: oidcOptionsResponseTypesDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCResponseTypeEnumValues)...),
								),
							},
						},

						"token_endpoint_authn_method": schema.StringAttribute{
							Description:         oidcOptionsTokenEndpointAuthnMethod.Description,
							MarkdownDescription: oidcOptionsTokenEndpointAuthnMethod.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCTokenAuthMethodEnumValues)...),
							},
						},

						"par_requirement": schema.StringAttribute{
							Description:         oidcOptionsParRequirementDescription.Description,
							MarkdownDescription: oidcOptionsParRequirementDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationOIDCPARRequirementEnumValues)...),
							},
						},

						"par_timeout": schema.Int64Attribute{
							Description:         oidcOptionsParTimeoutDescription.Description,
							MarkdownDescription: oidcOptionsParTimeoutDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(oidcOptionsParTimeoutDefault),

							Validators: []validator.Int64{
								int64validator.Between(oidcOptionsParTimeoutMin, oidcOptionsParTimeoutMax),
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
							Description:         oidcOptionsRedirectUrisDescription.Description,
							MarkdownDescription: oidcOptionsRedirectUrisDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^(http:\/\/((localhost)|(127\.0\.0\.1))(:[0-9]+)?(\/?(.+))?$|(\S+:\/\/).+)`), "Expected value to have a url with schema of \"https\" or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  \"http\" urls are permitted when using localhost hosts \"localhost\" and \"127.0.0.1\"."),
								),
							},
						},

						"allow_wildcards_in_redirect_uris": schema.BoolAttribute{
							Description:         oidcOptionsAllowWildcardsInRedirectUrisDescription.Description,
							MarkdownDescription: oidcOptionsAllowWildcardsInRedirectUrisDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"post_logout_redirect_uris": schema.SetAttribute{
							Description:         oidcOptionsPostLogoutRedirectUrisDescription.Description,
							MarkdownDescription: oidcOptionsPostLogoutRedirectUrisDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^(\S+:\/\/).+`), "Expected value to have a url with schema of \"https\", \"http\" or a custom mobile native schema (e.g., `org.bxretail.app://logout`)."),
								),
							},
						},

						"refresh_token_duration": schema.Int64Attribute{
							Description:         oidcOptionsRefreshTokenDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(oidcOptionsRefreshTokenDurationDefault),

							Validators: []validator.Int64{
								int64validator.Between(oidcOptionsRefreshTokenDurationMin, oidcOptionsRefreshTokenDurationMax),
							},
						},

						"refresh_token_rolling_duration": schema.Int64Attribute{
							Description:         oidcOptionsRefreshTokenRollingDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenRollingDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(oidcOptionsRefreshTokenRollingDurationDefault),

							Validators: []validator.Int64{
								int64validator.Between(oidcOptionsRefreshTokenRollingDurationMin, oidcOptionsRefreshTokenRollingDurationMax),
							},
						},

						"refresh_token_rolling_grace_period_duration": schema.Int64Attribute{
							Description:         oidcOptionsRefreshTokenRollingGracePeriodDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenRollingGracePeriodDurationDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.Int64{
								int64validator.Between(oidcOptionsRefreshTokenRollingGracePeriodDurationMin, oidcOptionsRefreshTokenRollingGracePeriodDurationMax),
							},
						},

						"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
							Description:         oidcOptionsAdditionalRefreshTokenReplayProtectionEnabledDescription.Description,
							MarkdownDescription: oidcOptionsAdditionalRefreshTokenReplayProtectionEnabledDescription.MarkdownDescription,
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
							Description:         oidcOptionsClientSecretDescription.Description,
							MarkdownDescription: oidcOptionsClientSecretDescription.MarkdownDescription,
							Computed:            true,
							Sensitive:           true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},

						"support_unsigned_request_object": schema.BoolAttribute{
							Description:         oidcOptionsSupportUnsignedRequestObjectDescription.Description,
							MarkdownDescription: oidcOptionsSupportUnsignedRequestObjectDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"require_signed_request_object": schema.BoolAttribute{
							Description:         oidcOptionsRequireSignedRequestObjectDescription.Description,
							MarkdownDescription: oidcOptionsRequireSignedRequestObjectDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"bundle_id": schema.StringAttribute{
							Description:         oidcOptionsBundleIdDescription.Description,
							MarkdownDescription: oidcOptionsBundleIdDescription.MarkdownDescription,
							DeprecationMessage:  "This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead.",
							Optional:            true,
							Computed:            true,

							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},

						"package_name": schema.StringAttribute{
							Description:         oidcOptionsPackageNameDescription.Description,
							MarkdownDescription: oidcOptionsPackageNameDescription.MarkdownDescription,
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
							Description:         oidcOptionsCertificateBasedAuthenticationDescription.Description,
							MarkdownDescription: oidcOptionsCertificateBasedAuthenticationDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"key_id": schema.StringAttribute{
										Description:         oidcOptionsCertificateBasedAuthenticationKeyIdDescription.Description,
										MarkdownDescription: oidcOptionsCertificateBasedAuthenticationKeyIdDescription.MarkdownDescription,
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
							Description:         oidcOptionsMobileAppDescription.Description,
							MarkdownDescription: oidcOptionsMobileAppDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"bundle_id": schema.StringAttribute{
										Description:         oidcOptionsMobileAppBundleIdDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppBundleIdDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},

									"package_name": schema.StringAttribute{
										Description:         oidcOptionsMobileAppPackageNameDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppPackageNameDescription.MarkdownDescription,
										Optional:            true,

										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},

									"huawei_app_id": schema.StringAttribute{
										Description:         oidcOptionsMobileAppHuaweiAppIdDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppHuaweiAppIdDescription.MarkdownDescription,
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
										Description:         oidcOptionsMobileAppHuaweiPackageNameDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppHuaweiPackageNameDescription.MarkdownDescription,
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
										Description:         oidcOptionsMobileAppPasscodeRefreshSecondsDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppPasscodeRefreshSecondsDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int64default.StaticInt64(oidcOptionsMobileAppPasscodeRefreshSecondsDefault),

										Validators: []validator.Int64{
											int64validator.Between(oidcOptionsMobileAppPasscodeRefreshSecondsMin, oidcOptionsMobileAppPasscodeRefreshSecondsMax),
										},
									},

									"universal_app_link": schema.StringAttribute{
										Description:         oidcOptionsMobileAppUniversalAppLinkDescription.Description,
										MarkdownDescription: oidcOptionsMobileAppUniversalAppLinkDescription.MarkdownDescription,
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
													Description:         oidcOptionsMobileAppIntegrityDetectionEnabledDescription.Description,
													MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionEnabledDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: booldefault.StaticBool(false),
												},

												"excluded_platforms": schema.SetAttribute{
													Description:         oidcOptionsMobileAppIntegrityDetectionExcludedPlatformsDescription.Description,
													MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionExcludedPlatformsDescription.MarkdownDescription,
													Optional:            true,

													ElementType: types.StringType,

													Validators: []validator.Set{
														setvalidator.SizeAtMost(1),
														setvalidator.ValueStringsAre(
															stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumMobileIntegrityDetectionPlatformEnumValues)...),
														),
													},
												},
											},

											Blocks: map[string]schema.Block{
												"cache_duration": schema.ListNestedBlock{
													Description:         oidcOptionsMobileAppIntegrityDetectionCacheDurationDescription.Description,
													MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionCacheDurationDescription.MarkdownDescription,

													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"amount": schema.Int64Attribute{
																Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.").Description,
																Required:    true,
															},

															"units": schema.StringAttribute{
																Description:         oidcOptionsMobileAppIntegrityDetectionCacheDurationUnitsDescription.Description,
																MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionCacheDurationUnitsDescription.MarkdownDescription,
																Optional:            true,
																Computed:            true,

																Default: stringdefault.StaticString(string(management.ENUMDURATIONUNITMINSHOURS_MINUTES)),

																Validators: []validator.String{
																	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumDurationUnitMinsHoursEnumValues)...),
																},
															},
														},
													},

													Validators: []validator.List{
														listvalidator.SizeAtMost(1),
													},
												},

												"google_play": schema.ListNestedBlock{
													Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription.Description,
													MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription.MarkdownDescription,

													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"decryption_key": schema.StringAttribute{
																Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayDecryptionKeyDescription.Description,
																MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayDecryptionKeyDescription.MarkdownDescription,
																Optional:            true,
																Sensitive:           true,

																Validators: []validator.String{
																	stringvalidator.ConflictsWith(
																		path.MatchRelative().AtParent().AtName("service_account_credentials_json"),
																	),
																},
															},

															"service_account_credentials_json": schema.StringAttribute{
																Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayServiceAccountCredentialsJsonDescription.Description,
																MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayServiceAccountCredentialsJsonDescription.MarkdownDescription,
																Optional:            true,
																Sensitive:           true,

																Validators: []validator.String{
																	stringvalidator.ConflictsWith(
																		path.MatchRelative().AtParent().AtName("decryption_key"),
																		path.MatchRelative().AtParent().AtName("verification_key"),
																	),
																},
															},

															"verification_key": schema.StringAttribute{
																Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationKeyDescription.Description,
																MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationKeyDescription.MarkdownDescription,
																Optional:            true,
																Sensitive:           true,

																Validators: []validator.String{
																	stringvalidator.ConflictsWith(
																		path.MatchRelative().AtParent().AtName("service_account_credentials_json"),
																	),
																},
															},

															"verification_type": schema.StringAttribute{
																Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationTypeDescription.Description,
																MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayVerificationTypeDescription.MarkdownDescription,
																Required:            true,

																Validators: []validator.String{
																	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationNativeGooglePlayVerificationTypeEnumValues)...),
																},
															},
														},
													},

													Validators: []validator.List{
														listvalidator.SizeAtMost(1),
													},
												},
											},
										},

										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
										},
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},

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
							Optional:    true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Expected value to have a url with schema of \"https\"."),
							},
						},

						"type": schema.StringAttribute{
							Description:         samlOptionsTypeDescription.Description,
							MarkdownDescription: samlOptionsTypeDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONTYPE_WEB_APP)),

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},

							Validators: []validator.String{
								stringvalidator.OneOf(
									string(management.ENUMAPPLICATIONTYPE_WEB_APP),
									string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP),
								),
							},
						},

						"acs_urls": schema.SetAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).").Description,
							Required:    true,

							ElementType: types.StringType,
						},

						"assertion_duration": schema.Int64Attribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the assertion validity duration in seconds.").Description,
							Required:    true,
						},

						"assertion_signed_enabled": schema.BoolAttribute{
							Description:         samlOptionsAssertionSignedEnabledDescription.Description,
							MarkdownDescription: samlOptionsAssertionSignedEnabledDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(true),
						},

						"idp_signing_key_id": schema.StringAttribute{
							Description:         samlOptionsIdpSigningKeyIdDescription.Description,
							MarkdownDescription: samlOptionsIdpSigningKeyIdDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							DeprecationMessage: "The `idp_signing_key_id` attribute is deprecated and will be removed in the next major release.  Please use the `idp_signing_key` block going forward.",

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("idp_signing_key"),
								),
							},
						},

						"default_target_url": schema.StringAttribute{
							Description:         samlOptionsDefaultTargetUrlDescription.Description,
							MarkdownDescription: samlOptionsDefaultTargetUrlDescription.MarkdownDescription,
							Optional:            true,
						},

						"enable_requested_authn_context": schema.BoolAttribute{
							Description:         samlOptionsEnableRequestedAuthnContextDescription.Description,
							MarkdownDescription: samlOptionsEnableRequestedAuthnContextDescription.MarkdownDescription,
							Optional:            true,
						},

						"nameid_format": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the format of the Subject NameID attibute in the SAML assertion.").Description,
							Optional:    true,
						},

						"response_is_signed": schema.BoolAttribute{
							Description:         samlOptionsResponseIsSignedDescription.Description,
							MarkdownDescription: samlOptionsResponseIsSignedDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"slo_binding": schema.StringAttribute{
							Description:         samlOptionsSloBindingDescription.Description,
							MarkdownDescription: samlOptionsSloBindingDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumApplicationSAMLSloBindingEnumValues)...),
							},
						},

						"slo_endpoint": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the logout endpoint URL. This is an optional property. However, if a logout endpoint URL is not defined, logout actions result in an error.").Description,
							Optional:    true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Expected value to have a url with schema of \"http\" or \"https\"."),
							},
						},

						"slo_response_endpoint": schema.StringAttribute{
							Description:         samlOptionsSloResponseEndpointDescription.Description,
							MarkdownDescription: samlOptionsSloResponseEndpointDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Expected value to have a url with schema of \"http\" or \"https\"."),
							},
						},

						"slo_window": schema.Int64Attribute{
							Description:         samlOptionsSloWindowDescription.Description,
							MarkdownDescription: samlOptionsSloWindowDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.Int64{
								int64validator.Between(samlOptionsSloWindowMin, samlOptionsSloWindowMax),
							},
						},

						"sp_entity_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.").Description,
							Required:    true,
						},

						"sp_verification_certificate_ids": schema.SetAttribute{
							Description:         samlOptionsSpVerificationCertificateIds.Description,
							MarkdownDescription: samlOptionsSpVerificationCertificateIds.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							DeprecationMessage: "The `sp_verification_certificate_ids` parameter is deprecated and will be removed in the next major release.  Please use the `sp_verification.certificate_ids` parameter going forward.",

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									verify.P1ResourceIDValidator(),
								),
								setvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("sp_verification"),
								),
							},
						},
					},

					Blocks: map[string]schema.Block{
						"idp_signing_key": schema.ListNestedBlock{
							Description:         samlOptionsIdpSigningKeyDescription.Description,
							MarkdownDescription: samlOptionsIdpSigningKeyDescription.MarkdownDescription,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"algorithm": schema.StringAttribute{
										Description:         samlOptionsIdpSigningKeyAlgorithmDescription.Description,
										MarkdownDescription: samlOptionsIdpSigningKeyAlgorithmDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeySignagureAlgorithmEnumValues)...),
										},
									},

									"key_id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.  Must be a valid PingOne resource ID.").Description,
										Required:    true,

										Validators: []validator.String{
											verify.P1ResourceIDValidator(),
										},
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
								listvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("idp_signing_key_id"),
								),
							},
						},

						"sp_verification": schema.ListNestedBlock{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block item that specifies SP signature verification settings.").Description,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"authn_request_signed": schema.BoolAttribute{
										Description:         samlOptionsSpVerificationAuthnRequestSignedDescription.Description,
										MarkdownDescription: samlOptionsSpVerificationAuthnRequestSignedDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: booldefault.StaticBool(false),
									},

									"certificate_ids": schema.SetAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A list that specifies the certificate IDs used to verify the service provider signature.  Values must be valid PingOne resource IDs.").Description,
										ElementType: types.StringType,
										Required:    true,

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

	const originsMax = 20
	originsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value must be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Must be non-empty when `behavior` is `%s` and must be omitted or empty when `behavior` is `%s`.  Limited to %d values.", string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_SPECIFIC_ORIGINS), string(management.ENUMAPPLICATIONCORSSETTINGSBEHAVIOR_NO_ORIGINS), originsMax),
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
						setvalidator.SizeAtMost(originsMax),
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

func (r *ApplicationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if !data.OIDCOptions.IsNull() && !data.OIDCOptions.IsUnknown() {
		var plan []ApplicationOIDCOptionsResourceModel
		resp.Diagnostics.Append(data.OIDCOptions.ElementsAs(ctx, &plan, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(plan) > 0 {
			// Certificate based authentication
			if !plan[0].CertificateBasedAuthentication.IsNull() && !plan[0].CertificateBasedAuthentication.IsUnknown() {
				if !plan[0].Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_NATIVE_APP))) {
					resp.Diagnostics.AddAttributeError(
						path.Root("oidc_options").AtName("certificate_based_authentication"),
						"Invalid configuration",
						fmt.Sprintf("`certificate_based_authentication` can only be set with OIDC applications that have a `type` value of `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
					)
				}
			}

			// Wildcards in redirect URIs
			if !plan[0].RedirectUris.IsNull() && !plan[0].RedirectUris.IsUnknown() {
				var uris []string
				resp.Diagnostics.Append(plan[0].RedirectUris.ElementsAs(ctx, &uris, false)...)

				if resp.Diagnostics.HasError() {
					return
				}

				for _, uri := range uris {
					if strings.Contains(uri, "*") {
						if plan[0].AllowWildcardsInRedirectUris.IsNull() || plan[0].AllowWildcardsInRedirectUris.Equal(types.BoolValue(false)) {
							resp.Diagnostics.AddAttributeError(
								path.Root("oidc_options").AtName("redirect_uris"),
								"Invalid configuration",
								"Current configuration is invalid as wildcards are not allowed in redirect URIs.  Wildcards can be enabled by setting `allow_wildcards_in_redirect_uris` to `true`.",
							)
							break
						}

					}
				}
			}
		}
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
	application, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var createResponse *management.CreateApplication201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.CreateApplication(ctx, plan.EnvironmentId.ValueString()).CreateApplicationRequest(*application).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplication",
		applicationWriteCustomError,
		sdk.DefaultCreateReadRetryable,
		&createResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var applicationId string
	if createResponse.ApplicationOIDC != nil && createResponse.ApplicationOIDC.GetId() != "" {
		applicationId = createResponse.ApplicationOIDC.GetId()
	} else if createResponse.ApplicationSAML != nil && createResponse.ApplicationSAML.GetId() != "" {
		applicationId = createResponse.ApplicationSAML.GetId()
	} else if createResponse.ApplicationExternalLink != nil && createResponse.ApplicationExternalLink.GetId() != "" {
		applicationId = createResponse.ApplicationExternalLink.GetId()
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Cannot determine application ID from API response for application: %s", plan.Name.ValueString()),
			fmt.Sprintf("Full response object: %v\n", resp),
		)
	}

	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, plan.EnvironmentId.ValueString(), applicationId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
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
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, plan.EnvironmentId.ValueString(), applicationId).Execute()
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
	resp.Diagnostics.Append(state.toState(response, secretResponse)...)
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
	resp.Diagnostics.Append(data.toState(response, secretResponse)...)
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
	application, d := plan.expandUpdate(ctx)
	resp.Diagnostics = append(resp.Diagnostics, d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateApplicationRequest(*application).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
		applicationWriteCustomError,
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
	resp.Diagnostics.Append(state.toState(response, secretResponse)...)
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

func applicationWriteCustomError(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Wildcards in redirect URis
	m, err := regexp.MatchString("^Wildcards are not allowed in redirect URIs.", error.GetMessage())
	if err != nil {
		diags.AddError("API Validation error", "Cannot match error string for wildcard in redirect URIs")
		return diags
	}
	if m {
		diags.AddError("Invalid configuration", "Current configuration is invalid as wildcards are not allowed in redirect URIs.  Wildcards can be enabled by setting `allow_wildcards_in_redirect_uris` to `true`.")

		return diags
	}

	return framework.DefaultCustomError(error)
}

func (p *ApplicationResourceModel) expandCreate(ctx context.Context) (*management.CreateApplicationRequest, diag.Diagnostics) {
	var d, diags diag.Diagnostics

	data := &management.CreateApplicationRequest{}

	if !p.OIDCOptions.IsNull() && !p.OIDCOptions.IsUnknown() {
		data.ApplicationOIDC, d = p.expandApplicationOIDC(ctx)
		diags = append(diags, d...)
	}

	if !p.SAMLOptions.IsNull() && !p.SAMLOptions.IsUnknown() {
		data.ApplicationSAML, d = p.expandApplicationSAML(ctx)
		diags = append(diags, d...)
	}

	if !p.ExternalLinkOptions.IsNull() && !p.ExternalLinkOptions.IsUnknown() {
		data.ApplicationExternalLink, d = p.expandApplicationExternalLink(ctx)
		diags = append(diags, d...)
	}

	return data, diags
}

func (p *ApplicationResourceModel) expandUpdate(ctx context.Context) (*management.UpdateApplicationRequest, diag.Diagnostics) {
	var d, diags diag.Diagnostics

	data := &management.UpdateApplicationRequest{}

	if !p.OIDCOptions.IsNull() && !p.OIDCOptions.IsUnknown() {
		data.ApplicationOIDC, d = p.expandApplicationOIDC(ctx)
		diags = append(diags, d...)
	}

	if !p.SAMLOptions.IsNull() && !p.SAMLOptions.IsUnknown() {
		data.ApplicationSAML, d = p.expandApplicationSAML(ctx)
		diags = append(diags, d...)
	}

	if !p.ExternalLinkOptions.IsNull() && !p.ExternalLinkOptions.IsUnknown() {
		data.ApplicationExternalLink, d = p.expandApplicationExternalLink(ctx)
		diags = append(diags, d...)
	}

	return data, diags
}

func (p *ApplicationCorsSettingsResourceModel) expand() (*management.ApplicationCorsSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationCorsSettings(management.EnumApplicationCorsSettingsBehavior(p.Behavior.ValueString()))

	if !p.Origins.IsNull() && !p.Origins.IsUnknown() {
		var origins []string
		d := p.Origins.ElementsAs(context.Background(), &origins, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetOrigins(origins)
	}

	return data, diags
}

func (p *ApplicationResourceModel) expandApplicationOIDC(ctx context.Context) (*management.ApplicationOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationOIDC

	if !p.OIDCOptions.IsNull() && !p.OIDCOptions.IsUnknown() {
		var plan []ApplicationOIDCOptionsResourceModel
		d := p.OIDCOptions.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `oidc_options` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		grantTypes := make([]management.EnumApplicationOIDCGrantType, 0)

		var grantTypesPlan []string

		diags.Append(planItem.GrantTypes.ElementsAs(ctx, &grantTypesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, v := range grantTypesPlan {
			grantTypes = append(grantTypes, management.EnumApplicationOIDCGrantType(v))
		}

		data = management.NewApplicationOIDC(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT,
			management.EnumApplicationType(planItem.Type.ValueString()),
			grantTypes,
			management.EnumApplicationOIDCTokenAuthMethod(planItem.TokenEndpointAuthnMethod.ValueString()),
		)

		applicationCommon, d := p.expandApplicationCommon(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Description = applicationCommon.Description
		data.LoginPageUrl = applicationCommon.LoginPageUrl
		data.Icon = applicationCommon.Icon
		data.AccessControl = applicationCommon.AccessControl
		data.HiddenFromAppPortal = applicationCommon.HiddenFromAppPortal

		if !planItem.CorsSettings.IsNull() && !planItem.CorsSettings.IsUnknown() {
			var corsPlan []ApplicationCorsSettingsResourceModel

			diags.Append(planItem.CorsSettings.ElementsAs(ctx, &corsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(corsPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `oidc_options.cors_settings` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			corsSettings, d := corsPlan[0].expand()
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetCorsSettings(*corsSettings)
		}

		if !planItem.HomePageUrl.IsNull() && !planItem.HomePageUrl.IsUnknown() {
			data.SetHomePageUrl(planItem.HomePageUrl.ValueString())
		}

		if !planItem.InitiateLoginUri.IsNull() && !planItem.InitiateLoginUri.IsUnknown() {
			data.SetInitiateLoginUri(planItem.InitiateLoginUri.ValueString())
		}

		if !planItem.Jwks.IsNull() && !planItem.Jwks.IsUnknown() {
			data.SetJwks(planItem.Jwks.ValueString())
		}

		if !planItem.JwksUrl.IsNull() && !planItem.JwksUrl.IsUnknown() {
			data.SetJwksUrl(planItem.JwksUrl.ValueString())
		}

		if !planItem.TargetLinkUri.IsNull() && !planItem.TargetLinkUri.IsUnknown() {
			data.SetTargetLinkUri(planItem.TargetLinkUri.ValueString())
		}

		if !planItem.ResponseTypes.IsNull() && !planItem.ResponseTypes.IsUnknown() {
			var responseTypesPlan []string

			diags.Append(planItem.ResponseTypes.ElementsAs(ctx, &responseTypesPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			obj := make([]management.EnumApplicationOIDCResponseType, 0)

			for _, v := range responseTypesPlan {
				obj = append(obj, management.EnumApplicationOIDCResponseType(v))
			}
			data.SetResponseTypes(obj)
		}

		if !planItem.ParRequirement.IsNull() && !planItem.ParRequirement.IsUnknown() {
			data.SetParRequirement(management.EnumApplicationOIDCPARRequirement(planItem.ParRequirement.ValueString()))
		}

		if !planItem.ParTimeout.IsNull() && !planItem.ParTimeout.IsUnknown() {
			data.SetParTimeout(int32(planItem.ParTimeout.ValueInt64()))
		}

		if !planItem.PKCEEnforcement.IsNull() && !planItem.PKCEEnforcement.IsUnknown() {
			data.SetPkceEnforcement(management.EnumApplicationOIDCPKCEOption(planItem.PKCEEnforcement.ValueString()))
		}

		if !planItem.RedirectUris.IsNull() && !planItem.RedirectUris.IsUnknown() {
			var redirectUrisPlan []string

			diags.Append(planItem.RedirectUris.ElementsAs(ctx, &redirectUrisPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetRedirectUris(redirectUrisPlan)
		}

		if !planItem.AllowWildcardsInRedirectUris.IsNull() && !planItem.AllowWildcardsInRedirectUris.IsUnknown() {
			data.SetAllowWildcardInRedirectUris(planItem.AllowWildcardsInRedirectUris.ValueBool())
		}

		if !planItem.PostLogoutRedirectUris.IsNull() && !planItem.PostLogoutRedirectUris.IsUnknown() {
			var postLogoutRedirectUrisPlan []string

			diags.Append(planItem.PostLogoutRedirectUris.ElementsAs(ctx, &postLogoutRedirectUrisPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetPostLogoutRedirectUris(postLogoutRedirectUrisPlan)
		}

		if !planItem.RefreshTokenDuration.IsNull() && !planItem.RefreshTokenDuration.IsUnknown() {
			data.SetRefreshTokenDuration(int32(planItem.RefreshTokenDuration.ValueInt64()))
		}

		if !planItem.RefreshTokenRollingDuration.IsNull() && !planItem.RefreshTokenRollingDuration.IsUnknown() {
			data.SetRefreshTokenRollingDuration(int32(planItem.RefreshTokenRollingDuration.ValueInt64()))
		}

		if !planItem.RefreshTokenRollingGracePeriodDuration.IsNull() && !planItem.RefreshTokenRollingGracePeriodDuration.IsUnknown() {
			data.SetRefreshTokenRollingGracePeriodDuration(int32(planItem.RefreshTokenRollingGracePeriodDuration.ValueInt64()))
		}

		if !planItem.AdditionalRefreshTokenReplayProtectionEnabled.IsNull() && !planItem.AdditionalRefreshTokenReplayProtectionEnabled.IsUnknown() {
			data.SetAdditionalRefreshTokenReplayProtectionEnabled(planItem.AdditionalRefreshTokenReplayProtectionEnabled.ValueBool())
		}

		if !p.Tags.IsNull() && !p.Tags.IsUnknown() {
			var tagsPlan []string

			diags.Append(p.Tags.ElementsAs(ctx, &tagsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			tags := make([]management.EnumApplicationTags, 0)

			for _, v := range tagsPlan {
				tags = append(tags, management.EnumApplicationTags(v))
			}

			data.Tags = tags

		}

		data.SetAssignActorRoles(false)

		if !planItem.CertificateBasedAuthentication.IsNull() && !planItem.CertificateBasedAuthentication.IsUnknown() {
			if !planItem.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_NATIVE_APP))) {
				diags.AddError(
					"Invalid configuration",
					fmt.Sprintf("`certificate_based_authentication` can only be set with applications that have a `type` value of `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
				)

				return nil, diags
			}

			var kerberosPlan []ApplicationOIDCCertificateBasedAuthenticationResourceModel

			diags.Append(planItem.CertificateBasedAuthentication.ElementsAs(ctx, &kerberosPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(kerberosPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `oidc_options.certificate_based_authentication` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			data.SetKerberos(*management.NewApplicationOIDCAllOfKerberos(*management.NewApplicationOIDCAllOfKerberosKey(kerberosPlan[0].KeyId.ValueString())))
		}

		if !planItem.SupportUnsignedRequestObject.IsNull() && !planItem.SupportUnsignedRequestObject.IsUnknown() {
			data.SetSupportUnsignedRequestObject(planItem.SupportUnsignedRequestObject.ValueBool())
		}

		if !planItem.RequireSignedRequestObject.IsNull() && !planItem.RequireSignedRequestObject.IsUnknown() {
			data.SetRequireSignedRequestObject(planItem.RequireSignedRequestObject.ValueBool())
		}

		if !planItem.MobileApp.IsNull() && !planItem.MobileApp.IsUnknown() {
			var mobileAppPlan []ApplicationOIDCMobileAppResourceModel

			diags.Append(planItem.MobileApp.ElementsAs(ctx, &mobileAppPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(mobileAppPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `oidc_options.mobile_app` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			mobile, d := mobileAppPlan[0].expand(ctx)
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetMobile(*mobile)
		}

		if !planItem.BundleId.IsNull() && !planItem.BundleId.IsUnknown() {
			data.SetBundleId(planItem.BundleId.ValueString())
		}

		if !planItem.PackageName.IsNull() && !planItem.PackageName.IsUnknown() {
			data.SetPackageName(planItem.PackageName.ValueString())
		}
	}

	return data, diags
}

func (p *ApplicationOIDCMobileAppResourceModel) expand(ctx context.Context) (*management.ApplicationOIDCAllOfMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationOIDCAllOfMobile()

	if !p.BundleId.IsNull() && !p.BundleId.IsUnknown() {
		data.SetBundleId(p.BundleId.ValueString())
	}

	if !p.HuaweiAppId.IsNull() && !p.HuaweiAppId.IsUnknown() {
		data.SetHuaweiAppId(p.HuaweiAppId.ValueString())
	}

	if !p.HuaweiPackageName.IsNull() && !p.HuaweiPackageName.IsUnknown() {
		data.SetHuaweiPackageName(p.HuaweiPackageName.ValueString())
	}

	if !p.IntegrityDetection.IsNull() && !p.IntegrityDetection.IsUnknown() {

		var integrityDetectionPlan []ApplicationOIDCMobileAppIntegrityDetectionResourceModel
		diags.Append(p.IntegrityDetection.ElementsAs(ctx, &integrityDetectionPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(integrityDetectionPlan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `oidc_options.mobile_app` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		integrityDetection, d := integrityDetectionPlan[0].expand(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetIntegrityDetection(*integrityDetection)
	}

	if !p.PackageName.IsNull() && !p.PackageName.IsUnknown() {
		data.SetPackageName(p.PackageName.ValueString())
	}

	if !p.PasscodeRefreshSeconds.IsNull() && !p.PasscodeRefreshSeconds.IsUnknown() {
		data.SetPasscodeRefreshDuration(*management.NewApplicationOIDCAllOfMobilePasscodeRefreshDuration(
			int32(p.PasscodeRefreshSeconds.ValueInt64()),
			management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS,
		))
	}

	if !p.UniversalAppLink.IsNull() && !p.UniversalAppLink.IsUnknown() {
		data.SetUriPrefix(p.UniversalAppLink.ValueString())
	}

	return data, diags
}

func (p *ApplicationOIDCMobileAppIntegrityDetectionResourceModel) expand(ctx context.Context) (*management.ApplicationOIDCAllOfMobileIntegrityDetection, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationOIDCAllOfMobileIntegrityDetection()

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		var mode management.EnumEnabledStatus
		if p.Enabled.ValueBool() {
			mode = management.ENUMENABLEDSTATUS_ENABLED
		} else {
			mode = management.ENUMENABLEDSTATUS_DISABLED
		}
		data.SetMode(mode)
	}

	googleVerificationIncluded := true

	if !p.ExcludedPlatforms.IsNull() && !p.ExcludedPlatforms.IsUnknown() {
		var excludedPlatformsPlan []string

		diags.Append(p.ExcludedPlatforms.ElementsAs(ctx, &excludedPlatformsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		excludedPlatforms := make([]management.EnumMobileIntegrityDetectionPlatform, 0)

		for _, v := range excludedPlatformsPlan {
			excludedPlatforms = append(excludedPlatforms, management.EnumMobileIntegrityDetectionPlatform(v))
			if v == string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE) {
				googleVerificationIncluded = false
			}
		}

		data.SetExcludedPlatforms(excludedPlatforms)
	}

	if !p.GooglePlay.IsNull() && !p.GooglePlay.IsUnknown() {

		var googlePlayPlan []ApplicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModel
		diags.Append(p.GooglePlay.ElementsAs(ctx, &googlePlayPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(googlePlayPlan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `oidc_options.mobile_app.integrity_detection.google_play` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		googlePlay := management.NewApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay()

		if !googlePlayPlan[0].DecryptionKey.IsNull() && !googlePlayPlan[0].DecryptionKey.IsUnknown() {
			googlePlay.SetDecryptionKey(googlePlayPlan[0].DecryptionKey.ValueString())
		}

		if !googlePlayPlan[0].ServiceAccountCredentialsJson.IsNull() && !googlePlayPlan[0].ServiceAccountCredentialsJson.IsUnknown() {
			googlePlay.SetServiceAccountCredentials(googlePlayPlan[0].ServiceAccountCredentialsJson.ValueString())
		}

		if !googlePlayPlan[0].VerificationKey.IsNull() && !googlePlayPlan[0].VerificationKey.IsUnknown() {
			googlePlay.SetVerificationKey(googlePlayPlan[0].VerificationKey.ValueString())
		}

		if !googlePlayPlan[0].VerificationType.IsNull() && !googlePlayPlan[0].VerificationType.IsUnknown() {
			googlePlay.SetVerificationType(management.EnumApplicationNativeGooglePlayVerificationType(googlePlayPlan[0].VerificationType.ValueString()))
		}

		data.SetGooglePlay(*googlePlay)
	} else {
		if googleVerificationIncluded {
			diags.AddError(
				"Invalid configuration",
				"The `oidc_options.mobile_app.integrity_detection.google_play` is required when the mobile integrity check is enabled in the application and `excluded_platforms` is unset, or `excluded_platforms` is not configured with `GOOGLE`.",
			)
		}

	}

	if !p.CacheDuration.IsNull() && !p.CacheDuration.IsUnknown() {

		var cacheDurationPlan []ApplicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModel
		diags.Append(p.CacheDuration.ElementsAs(ctx, &cacheDurationPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(cacheDurationPlan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `oidc_options.mobile_app.integrity_detection.cache_duration` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		cacheDuration := management.NewApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration()

		if !cacheDurationPlan[0].Amount.IsNull() && !cacheDurationPlan[0].Amount.IsUnknown() {
			cacheDuration.SetAmount(int32(cacheDurationPlan[0].Amount.ValueInt64()))
		}

		if !cacheDurationPlan[0].Units.IsNull() && !cacheDurationPlan[0].Units.IsUnknown() {
			cacheDuration.SetUnits(management.EnumDurationUnitMinsHours(cacheDurationPlan[0].Units.ValueString()))
		}

		data.SetCacheDuration(*cacheDuration)
	}

	return data, diags
}

func (p *ApplicationResourceModel) expandApplicationSAML(ctx context.Context) (*management.ApplicationSAML, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationSAML

	if !p.SAMLOptions.IsNull() && !p.SAMLOptions.IsUnknown() {
		var plan []ApplicationSAMLOptionsResourceModel
		d := p.SAMLOptions.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `saml_options` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		var acsUrls []string

		diags.Append(planItem.AcsUrls.ElementsAs(ctx, &acsUrls, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data = management.NewApplicationSAML(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_SAML,
			management.EnumApplicationType(planItem.Type.ValueString()),
			acsUrls,
			int32(planItem.AssertionDuration.ValueInt64()),
			planItem.SpEntityId.ValueString(),
		)

		applicationCommon, d := p.expandApplicationCommon(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Description = applicationCommon.Description
		data.LoginPageUrl = applicationCommon.LoginPageUrl
		data.Icon = applicationCommon.Icon
		data.AccessControl = applicationCommon.AccessControl
		data.HiddenFromAppPortal = applicationCommon.HiddenFromAppPortal

		// SAML specific options
		if !planItem.CorsSettings.IsNull() && !planItem.CorsSettings.IsUnknown() {
			var corsPlan []ApplicationCorsSettingsResourceModel

			diags.Append(planItem.CorsSettings.ElementsAs(ctx, &corsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(corsPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `saml_options.cors_settings` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			corsSettings, d := corsPlan[0].expand()
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetCorsSettings(*corsSettings)
		}

		if !planItem.HomePageUrl.IsNull() && !planItem.HomePageUrl.IsUnknown() {
			data.SetHomePageUrl(planItem.HomePageUrl.ValueString())
		}

		if !planItem.AssertionSignedEnabled.IsNull() && !planItem.AssertionSignedEnabled.IsUnknown() {
			data.SetAssertionSigned(planItem.AssertionSignedEnabled.ValueBool())
		}

		if !planItem.IdpSigningKeyId.IsNull() && !planItem.IdpSigningKeyId.IsUnknown() {
			data.SetIdpSigning(*management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(planItem.IdpSigningKeyId.ValueString())))
		}

		if !planItem.IdpSigningKey.IsNull() && !planItem.IdpSigningKey.IsUnknown() {

			var idpSigningOptionsPlan []ApplicationSAMLOptionsIdpSigningKeyResourceModel

			diags.Append(planItem.IdpSigningKey.ElementsAs(ctx, &idpSigningOptionsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(idpSigningOptionsPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `saml_options.idp_signing_key` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			idpSigning := *management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(idpSigningOptionsPlan[0].KeyId.ValueString()))
			idpSigning.SetAlgorithm(management.EnumCertificateKeySignagureAlgorithm(idpSigningOptionsPlan[0].Algorithm.ValueString()))

			data.SetIdpSigning(idpSigning)
		}

		if !planItem.EnableRequestedAuthnContext.IsNull() && !planItem.EnableRequestedAuthnContext.IsUnknown() {
			data.SetEnableRequestedAuthnContext(planItem.EnableRequestedAuthnContext.ValueBool())
		}

		if !planItem.DefaultTargetUrl.IsNull() && !planItem.DefaultTargetUrl.IsUnknown() {
			data.SetDefaultTargetUrl(planItem.DefaultTargetUrl.ValueString())
		}

		if !planItem.NameIdFormat.IsNull() && !planItem.NameIdFormat.IsUnknown() {
			data.SetNameIdFormat(planItem.NameIdFormat.ValueString())
		}

		if !planItem.ResponseIsSigned.IsNull() && !planItem.ResponseIsSigned.IsUnknown() {
			data.SetResponseSigned(planItem.ResponseIsSigned.ValueBool())
		}

		if !planItem.SloBinding.IsNull() && !planItem.SloBinding.IsUnknown() {
			data.SetSloBinding(management.EnumApplicationSAMLSloBinding(planItem.SloBinding.ValueString()))
		}

		if !planItem.SloEndpoint.IsNull() && !planItem.SloEndpoint.IsUnknown() {
			data.SetSloEndpoint(planItem.SloEndpoint.ValueString())
		}

		if !planItem.SloResponseEndpoint.IsNull() && !planItem.SloResponseEndpoint.IsUnknown() {
			data.SetSloResponseEndpoint(planItem.SloResponseEndpoint.ValueString())
		}

		if !planItem.SloWindow.IsNull() && !planItem.SloWindow.IsUnknown() {
			data.SetSloWindow(int32(planItem.SloWindow.ValueInt64()))
		}

		if !planItem.SpVerificationCertificateIds.IsNull() && !planItem.SpVerificationCertificateIds.IsUnknown() {
			var certificateIdsPlan []string

			diags.Append(planItem.SpVerificationCertificateIds.ElementsAs(ctx, &certificateIdsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			certificates := make([]management.ApplicationSAMLAllOfSpVerificationCertificates, 0)
			for _, v := range certificateIdsPlan {
				certificate := *management.NewApplicationSAMLAllOfSpVerificationCertificates(v)
				certificates = append(certificates, certificate)
			}

			data.SetSpVerification(*management.NewApplicationSAMLAllOfSpVerification(certificates))
		}

		if !planItem.SpVerification.IsNull() && !planItem.SpVerification.IsUnknown() {
			var spVerificationPlan []ApplicationSAMLOptionsSpVerificationResourceModel

			diags.Append(planItem.SpVerification.ElementsAs(ctx, &spVerificationPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(spVerificationPlan) == 0 {
				diags.AddError(
					"Invalid configuration",
					"The `saml_options.sp_verification` block is declared but has no configuration.  Please report this to the provider maintainers.",
				)
			}

			certificates := make([]management.ApplicationSAMLAllOfSpVerificationCertificates, 0)
			if !spVerificationPlan[0].CertificateIds.IsNull() && !spVerificationPlan[0].CertificateIds.IsUnknown() {
				var certificateIdsPlan []string

				diags.Append(spVerificationPlan[0].CertificateIds.ElementsAs(ctx, &certificateIdsPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}
				for _, v := range certificateIdsPlan {
					certificate := *management.NewApplicationSAMLAllOfSpVerificationCertificates(v)
					certificates = append(certificates, certificate)
				}
			}

			spVerification := management.NewApplicationSAMLAllOfSpVerification(certificates)

			if !spVerificationPlan[0].AuthnRequestSigned.IsNull() && !spVerificationPlan[0].AuthnRequestSigned.IsUnknown() {
				spVerification.SetAuthnRequestSigned(spVerificationPlan[0].AuthnRequestSigned.ValueBool())
			}

			data.SetSpVerification(*spVerification)
		}
	}

	return data, diags
}

func (p *ApplicationResourceModel) expandApplicationExternalLink(ctx context.Context) (*management.ApplicationExternalLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationExternalLink

	if !p.ExternalLinkOptions.IsNull() && !p.ExternalLinkOptions.IsUnknown() {
		var plan []ApplicationExternalLinkOptionsResourceModel
		d := p.ExternalLinkOptions.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `external_link_options` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		data = management.NewApplicationExternalLink(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_EXTERNAL_LINK,
			management.ENUMAPPLICATIONTYPE_PORTAL_LINK_APP,
			planItem.HomePageUrl.ValueString(),
		)

		applicationCommon, d := p.expandApplicationCommon(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.Description = applicationCommon.Description
		data.LoginPageUrl = applicationCommon.LoginPageUrl
		data.Icon = applicationCommon.Icon
		data.AccessControl = applicationCommon.AccessControl
		data.HiddenFromAppPortal = applicationCommon.HiddenFromAppPortal

	}

	return data, diags
}

func (p *ApplicationResourceModel) expandApplicationCommon(ctx context.Context) (*management.Application, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.Application{}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.LoginPageUrl.IsNull() && !p.LoginPageUrl.IsUnknown() {
		data.SetLoginPageUrl(p.LoginPageUrl.ValueString())
	}

	if !p.Icon.IsNull() && !p.Icon.IsUnknown() {
		var plan []service.ImageResourceModel
		d := p.Icon.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `icon` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)

			return nil, diags
		}

		iconPlanItem := plan[0]

		data.SetIcon(*management.NewApplicationIcon(
			iconPlanItem.Id.ValueString(),
			iconPlanItem.Href.ValueString(),
		))
	}

	accessControl := *management.NewApplicationAccessControl()
	accessControlCount := 0

	if !p.AccessControlRoleType.IsNull() && !p.AccessControlRoleType.IsUnknown() {
		accessControl.SetRole(*management.NewApplicationAccessControlRole(management.EnumApplicationAccessControlType(p.AccessControlRoleType.ValueString())))
		accessControlCount += 1
	}

	if !p.AccessControlGroupOptions.IsNull() && !p.AccessControlGroupOptions.IsUnknown() {
		var plan []ApplicationAccessControlGroupOptionsResourceModel
		d := p.AccessControlGroupOptions.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `access_control_group_options` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)

			return nil, diags
		}

		planItem := plan[0]

		groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)

		var groupsPlan []string

		diags.Append(planItem.Groups.ElementsAs(ctx, &groupsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, group := range groupsPlan {
			groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(group))
		}

		accessControl.SetGroup(*management.NewApplicationAccessControlGroup(
			management.EnumApplicationAccessControlGroupType(planItem.Type.ValueString()),
			groups,
		))

		accessControlCount += 1
	}

	if accessControlCount > 0 {
		data.SetAccessControl(accessControl)
	}

	if !p.HiddenFromAppPortal.IsNull() && !p.HiddenFromAppPortal.IsUnknown() {
		data.SetHiddenFromAppPortal(p.HiddenFromAppPortal.ValueBool())
	}

	return &data, diags
}

func (p *ApplicationResourceModel) toState(apiObject *management.ReadOneApplication200Response, apiSecretObject *management.ApplicationSecret) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics
	if v := apiObject.ApplicationExternalLink; v != nil {
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ListNull(types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes})
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
		p.OIDCOptions = types.ListNull(types.ObjectType{AttrTypes: applicationOidcOptionsTFObjectTypes})
		p.SAMLOptions = types.ListNull(types.ObjectType{AttrTypes: applicationSamlOptionsTFObjectTypes})

		p.ExternalLinkOptions, d = applicationExternalLinkOptionsToTF(v)
		diags = append(diags, d...)
	}

	if v := apiObject.ApplicationOIDC; v != nil {
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ListNull(types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes})
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

		p.OIDCOptions, d = applicationOidcOptionsToTF(v, apiSecretObject)
		diags = append(diags, d...)

		p.SAMLOptions = types.ListNull(types.ObjectType{AttrTypes: applicationSamlOptionsTFObjectTypes})
		p.ExternalLinkOptions = types.ListNull(types.ObjectType{AttrTypes: applicationExternalLinkOptionsTFObjectTypes})
	}

	if v := apiObject.ApplicationSAML; v != nil {
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringOkToTF(v.Environment.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())

		p.AccessControlRoleType = types.StringNull()
		p.AccessControlGroupOptions = types.ListNull(types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes})
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
		p.OIDCOptions = types.ListNull(types.ObjectType{AttrTypes: applicationOidcOptionsTFObjectTypes})

		p.SAMLOptions, d = applicationSamlOptionsToTF(v)
		diags = append(diags, d...)

		p.ExternalLinkOptions = types.ListNull(types.ObjectType{AttrTypes: applicationExternalLinkOptionsTFObjectTypes})
	}

	return diags
}

func applicationExternalLinkOptionsToTF(apiObject *management.ApplicationExternalLink) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationExternalLinkOptionsTFObjectTypes}

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"home_page_url": framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationExternalLinkOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationAccessControlGroupOptionsToTF(apiObject *management.ApplicationAccessControlGroup, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationAccessControlGroupOptionsTFObjectTypes}

	if !ok && apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	if v, ok := apiObject.GetGroupsOk(); ok {
		groups := make([]string, 0)

		for _, group := range v {
			groups = append(groups, group.GetId())
		}

		attributesMap["groups"] = framework.StringSetToTF(groups)
	}

	flattenedObj, d := types.ObjectValue(applicationAccessControlGroupOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationCorsSettingsOkToTF(apiObject *management.ApplicationCorsSettings, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"behavior": framework.EnumOkToTF(apiObject.GetBehaviorOk()),
		"origins":  framework.StringSetOkToTF(apiObject.GetOriginsOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationCorsSettingsTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsToTF(apiObject *management.ApplicationOIDC, apiObjectSecret *management.ApplicationSecret) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcOptionsTFObjectTypes}

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ListNull(tfObjType), diags
	}

	kerberos, d := applicationOidcOptionsCertificateBasedAuthenticationToTF(apiObject.GetKerberosOk())
	diags.Append(d...)

	corsSettings, d := applicationCorsSettingsOkToTF(apiObject.GetCorsSettingsOk())
	diags.Append(d...)

	mobileApp, d := applicationMobileAppOkToTF(apiObject.GetMobileOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"additional_refresh_token_replay_protection_enabled": framework.BoolOkToTF(apiObject.GetAdditionalRefreshTokenReplayProtectionEnabledOk()),
		"allow_wildcards_in_redirect_uris":                   framework.BoolOkToTF(apiObject.GetAllowWildcardInRedirectUrisOk()),
		"bundle_id":                                          framework.StringOkToTF(apiObject.GetBundleIdOk()),
		"certificate_based_authentication":                   kerberos,
		"client_id":                                          framework.StringOkToTF(apiObject.GetIdOk()),
		"client_secret":                                      types.StringNull(),
		"cors_settings":                                      corsSettings,
		"grant_types":                                        framework.EnumSetOkToTF(apiObject.GetGrantTypesOk()),
		"home_page_url":                                      framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"initiate_login_uri":                                 framework.StringOkToTF(apiObject.GetInitiateLoginUriOk()),
		"mobile_app":                                         mobileApp,
		"package_name":                                       framework.StringOkToTF(apiObject.GetPackageNameOk()),
		"par_requirement":                                    framework.EnumOkToTF(apiObject.GetParRequirementOk()),
		"par_timeout":                                        framework.Int32OkToTF(apiObject.GetParTimeoutOk()),
		"pkce_enforcement":                                   framework.EnumOkToTF(apiObject.GetPkceEnforcementOk()),
		"post_logout_redirect_uris":                          framework.StringSetOkToTF(apiObject.GetPostLogoutRedirectUrisOk()),
		"redirect_uris":                                      framework.StringSetOkToTF(apiObject.GetRedirectUrisOk()),
		"refresh_token_duration":                             framework.Int32OkToTF(apiObject.GetRefreshTokenDurationOk()),
		"refresh_token_rolling_duration":                     framework.Int32OkToTF(apiObject.GetRefreshTokenRollingDurationOk()),
		"refresh_token_rolling_grace_period_duration":        framework.Int32OkToTF(apiObject.GetRefreshTokenRollingGracePeriodDurationOk()),
		"require_signed_request_object":                      framework.BoolOkToTF(apiObject.GetRequireSignedRequestObjectOk()),
		"response_types":                                     framework.EnumSetOkToTF(apiObject.GetResponseTypesOk()),
		"support_unsigned_request_object":                    framework.BoolOkToTF(apiObject.GetSupportUnsignedRequestObjectOk()),
		"jwks":                                               framework.StringOkToTF(apiObject.GetJwksOk()),
		"jwks_url":                                           framework.StringOkToTF(apiObject.GetJwksUrlOk()),
		"target_link_uri":                                    framework.StringOkToTF(apiObject.GetTargetLinkUriOk()),
		"token_endpoint_authn_method":                        framework.EnumOkToTF(apiObject.GetTokenEndpointAuthMethodOk()),
		"type":                                               framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	if apiObjectSecret != nil {
		attributesMap["client_secret"] = framework.StringOkToTF(apiObjectSecret.GetSecretOk())
	}

	flattenedObj, d := types.ObjectValue(applicationOidcOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsCertificateBasedAuthenticationToTF(apiObject *management.ApplicationOIDCAllOfKerberos, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcOptionsCertificateAuthenticationTFObjectTypes}

	if !ok && apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"key_id": types.StringNull(),
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		attributesMap["key_id"] = framework.StringOkToTF(v.GetIdOk())
	}

	flattenedObj, d := types.ObjectValue(applicationOidcOptionsCertificateAuthenticationTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppOkToTF(apiObject *management.ApplicationOIDCAllOfMobile, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	integrityDetection, d := applicationMobileAppIntegrityDetectionOkToTF(apiObject.GetIntegrityDetectionOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"bundle_id":                framework.StringOkToTF(apiObject.GetBundleIdOk()),
		"huawei_app_id":            framework.StringOkToTF(apiObject.GetHuaweiAppIdOk()),
		"huawei_package_name":      framework.StringOkToTF(apiObject.GetHuaweiPackageNameOk()),
		"integrity_detection":      integrityDetection,
		"package_name":             framework.StringOkToTF(apiObject.GetPackageNameOk()),
		"passcode_refresh_seconds": types.Int64Null(),
		"universal_app_link":       framework.StringOkToTF(apiObject.GetUriPrefixOk()),
	}

	if v, ok := apiObject.GetPasscodeRefreshDurationOk(); ok && v != nil {

		if vU, ok := v.GetTimeUnitOk(); ok && *vU != management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS {
			diags.AddError(
				"Unsupported passcode refresh duration time unit",
				fmt.Sprintf("Expecting time unit of %s for attribute `passcode_refresh_seconds`, got %v", management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS, vU),
			)

			return types.ListNull(tfObjType), diags
		}

		attributesMap["passcode_refresh_seconds"] = framework.Int32OkToTF(v.GetDurationOk())
	}

	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetection, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	cacheDuration, d := applicationMobileAppIntegrityDetectionCacheDurationOkToTF(apiObject.GetCacheDurationOk())
	diags.Append(d...)

	googlePlay, d := applicationMobileAppIntegrityDetectionGooglePlayOkToTF(apiObject.GetGooglePlayOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"cache_duration":     cacheDuration,
		"enabled":            types.BoolNull(),
		"excluded_platforms": framework.EnumSetOkToTF(apiObject.GetExcludedPlatformsOk()),
		"google_play":        googlePlay,
	}

	if v, ok := apiObject.GetModeOk(); ok {
		if *v == management.ENUMENABLEDSTATUS_ENABLED {
			attributesMap["enabled"] = types.BoolValue(true)
		} else {
			attributesMap["enabled"] = types.BoolValue(false)
		}
	}

	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionCacheDurationOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"amount": framework.Int32OkToTF(apiObject.GetAmountOk()),
		"units":  framework.EnumOkToTF(apiObject.GetUnitsOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionGooglePlayOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"decryption_key":                   framework.StringOkToTF(apiObject.GetDecryptionKeyOk()),
		"service_account_credentials_json": framework.StringOkToTF(apiObject.GetServiceAccountCredentialsOk()),
		"verification_key":                 framework.StringOkToTF(apiObject.GetVerificationKeyOk()),
		"verification_type":                framework.EnumOkToTF(apiObject.GetVerificationTypeOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlOptionsToTF(apiObject *management.ApplicationSAML) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationSamlOptionsTFObjectTypes}

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ListNull(tfObjType), diags
	}

	corsSettings, d := applicationCorsSettingsOkToTF(apiObject.GetCorsSettingsOk())
	diags.Append(d...)

	idpSigningKey, idpSigningKeyId, d := applicationSamlIdpSigningKeyOkToTF(apiObject.GetIdpSigningOk())
	diags.Append(d...)

	spVerification, spVerificationCertificateIds, d := applicationSamlSpVerificationOkToTF(apiObject.GetSpVerificationOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"acs_urls":                        framework.StringSetOkToTF(apiObject.GetAcsUrlsOk()),
		"assertion_duration":              framework.Int32OkToTF(apiObject.GetAssertionDurationOk()),
		"assertion_signed_enabled":        framework.BoolOkToTF(apiObject.GetAssertionSignedOk()),
		"cors_settings":                   corsSettings,
		"enable_requested_authn_context":  framework.BoolOkToTF(apiObject.GetEnableRequestedAuthnContextOk()),
		"home_page_url":                   framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"idp_signing_key_id":              idpSigningKeyId,
		"idp_signing_key":                 idpSigningKey,
		"default_target_url":              framework.StringOkToTF(apiObject.GetDefaultTargetUrlOk()),
		"nameid_format":                   framework.StringOkToTF(apiObject.GetNameIdFormatOk()),
		"response_is_signed":              framework.BoolOkToTF(apiObject.GetResponseSignedOk()),
		"slo_binding":                     framework.EnumOkToTF(apiObject.GetSloBindingOk()),
		"slo_endpoint":                    framework.StringOkToTF(apiObject.GetSloEndpointOk()),
		"slo_response_endpoint":           framework.StringOkToTF(apiObject.GetSloResponseEndpointOk()),
		"slo_window":                      framework.Int32OkToTF(apiObject.GetSloWindowOk()),
		"sp_entity_id":                    framework.StringOkToTF(apiObject.GetSpEntityIdOk()),
		"sp_verification_certificate_ids": spVerificationCertificateIds,
		"sp_verification":                 spVerification,
		"type":                            framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	flattenedObj, d := types.ObjectValue(applicationSamlOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlIdpSigningKeyOkToTF(apiObject *management.ApplicationSAMLAllOfIdpSigning, ok bool) (types.List, types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationSamlOptionsIdpSigningKeyTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), types.StringNull(), diags
	}

	keyId := types.StringNull()

	attributesMap := map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
		"key_id":    types.StringNull(),
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		keyId = framework.StringOkToTF(v.GetIdOk())
		attributesMap["key_id"] = keyId
	}

	flattenedObj, d := types.ObjectValue(applicationSamlOptionsIdpSigningKeyTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, keyId, diags
}

func applicationSamlSpVerificationOkToTF(apiObject *management.ApplicationSAMLAllOfSpVerification, ok bool) (types.List, types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: applicationSamlOptionsSpVerificationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), types.SetNull(types.StringType), diags
	}

	certificateIds := make([]string, 0)
	if v, ok := apiObject.GetCertificatesOk(); ok {
		for _, certificate := range v {
			certificateIds = append(certificateIds, certificate.GetId())
		}
	}

	certificateIdsList := framework.StringSetToTF(certificateIds)

	attributesMap := map[string]attr.Value{
		"authn_request_signed": framework.BoolOkToTF(apiObject.GetAuthnRequestSignedOk()),
		"certificate_ids":      certificateIdsList,
	}

	flattenedObj, d := types.ObjectValue(applicationSamlOptionsSpVerificationTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, certificateIdsList, diags
}
