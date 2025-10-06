// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectplanmodifier"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso/helpers/beta"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResource serviceClientType

type applicationResourceModelV1 struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
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
	WSFedOptions              types.Object                 `tfsdk:"wsfed_options"`
}

type applicationAccessControlGroupOptionsResourceModelV1 struct {
	Type   types.String `tfsdk:"type"`
	Groups types.Set    `tfsdk:"groups"`
}

type applicationExternalLinkOptionsResourceModelV1 struct {
	HomePageUrl types.String `tfsdk:"home_page_url"`
}

type applicationOIDCOptionsResourceModelV1 struct {
	beta.ApplicationOIDCOptionsResourceModelV1
	AdditionalRefreshTokenReplayProtectionEnabled types.Bool   `tfsdk:"additional_refresh_token_replay_protection_enabled"`
	AllowWildcardsInRedirectUris                  types.Bool   `tfsdk:"allow_wildcard_in_redirect_uris"`
	CertificateBasedAuthentication                types.Object `tfsdk:"certificate_based_authentication"`
	CorsSettings                                  types.Object `tfsdk:"cors_settings"`
	DevicePathId                                  types.String `tfsdk:"device_path_id"`
	DeviceCustomVerificationUri                   types.String `tfsdk:"device_custom_verification_uri"`
	DeviceTimeout                                 types.Int32  `tfsdk:"device_timeout"`
	DevicePollingInterval                         types.Int32  `tfsdk:"device_polling_interval"`
	GrantTypes                                    types.Set    `tfsdk:"grant_types"`
	HomePageUrl                                   types.String `tfsdk:"home_page_url"`
	IdpSignoff                                    types.Bool   `tfsdk:"idp_signoff"`
	InitiateLoginUri                              types.String `tfsdk:"initiate_login_uri"`
	Jwks                                          types.String `tfsdk:"jwks"`
	JwksUrl                                       types.String `tfsdk:"jwks_url"`
	MobileApp                                     types.Object `tfsdk:"mobile_app"`
	ParRequirement                                types.String `tfsdk:"par_requirement"`
	ParTimeout                                    types.Int32  `tfsdk:"par_timeout"`
	PKCEEnforcement                               types.String `tfsdk:"pkce_enforcement"`
	PostLogoutRedirectUris                        types.Set    `tfsdk:"post_logout_redirect_uris"`
	RedirectUris                                  types.Set    `tfsdk:"redirect_uris"`
	RefreshTokenDuration                          types.Int32  `tfsdk:"refresh_token_duration"`
	RefreshTokenRollingDuration                   types.Int32  `tfsdk:"refresh_token_rolling_duration"`
	RefreshTokenRollingGracePeriodDuration        types.Int32  `tfsdk:"refresh_token_rolling_grace_period_duration"`
	RequireSignedRequestObject                    types.Bool   `tfsdk:"require_signed_request_object"`
	ResponseTypes                                 types.Set    `tfsdk:"response_types"`
	SupportUnsignedRequestObject                  types.Bool   `tfsdk:"support_unsigned_request_object"`
	TargetLinkUri                                 types.String `tfsdk:"target_link_uri"`
	TokenEndpointAuthnMethod                      types.String `tfsdk:"token_endpoint_auth_method"`
	Type                                          types.String `tfsdk:"type"`
}

type applicationCorsSettingsResourceModelV1 struct {
	Behavior types.String `tfsdk:"behavior"`
	Origins  types.Set    `tfsdk:"origins"`
}

type applicationOIDCCertificateBasedAuthenticationResourceModelV1 struct {
	KeyId pingonetypes.ResourceIDValue `tfsdk:"key_id"`
}

type applicationOIDCMobileAppResourceModelV1 struct {
	BundleId               types.String `tfsdk:"bundle_id"`
	HuaweiAppId            types.String `tfsdk:"huawei_app_id"`
	HuaweiPackageName      types.String `tfsdk:"huawei_package_name"`
	IntegrityDetection     types.Object `tfsdk:"integrity_detection"`
	PackageName            types.String `tfsdk:"package_name"`
	PasscodeRefreshSeconds types.Int32  `tfsdk:"passcode_refresh_seconds"`
	UniversalAppLink       types.String `tfsdk:"universal_app_link"`
}

type applicationOIDCMobileAppIntegrityDetectionResourceModelV1 struct {
	CacheDuration     types.Object `tfsdk:"cache_duration"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	ExcludedPlatforms types.Set    `tfsdk:"excluded_platforms"`
	GooglePlay        types.Object `tfsdk:"google_play"`
}

type applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV1 struct {
	Amount types.Int32  `tfsdk:"amount"`
	Units  types.String `tfsdk:"units"`
}

type applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV1 struct {
	DecryptionKey                 types.String         `tfsdk:"decryption_key"`
	ServiceAccountCredentialsJson jsontypes.Normalized `tfsdk:"service_account_credentials_json"`
	VerificationKey               types.String         `tfsdk:"verification_key"`
	VerificationType              types.String         `tfsdk:"verification_type"`
}

type applicationSAMLOptionsResourceModelV1 struct {
	AcsUrls                     types.Set    `tfsdk:"acs_urls"`
	AssertionDuration           types.Int32  `tfsdk:"assertion_duration"`
	AssertionSignedEnabled      types.Bool   `tfsdk:"assertion_signed_enabled"`
	CorsSettings                types.Object `tfsdk:"cors_settings"`
	DefaultTargetUrl            types.String `tfsdk:"default_target_url"`
	EnableRequestedAuthnContext types.Bool   `tfsdk:"enable_requested_authn_context"`
	HomePageUrl                 types.String `tfsdk:"home_page_url"`
	IdpSigningKey               types.Object `tfsdk:"idp_signing_key"`
	NameIdFormat                types.String `tfsdk:"nameid_format"`
	ResponseIsSigned            types.Bool   `tfsdk:"response_is_signed"`
	SessionNotOnOrAfterDuration types.Int32  `tfsdk:"session_not_on_or_after_duration"`
	SloBinding                  types.String `tfsdk:"slo_binding"`
	SloEndpoint                 types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint         types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                   types.Int32  `tfsdk:"slo_window"`
	SpEncryption                types.Object `tfsdk:"sp_encryption"`
	SpEntityId                  types.String `tfsdk:"sp_entity_id"`
	SpVerification              types.Object `tfsdk:"sp_verification"`
	Type                        types.String `tfsdk:"type"`
	VirtualServerIdSettings     types.Object `tfsdk:"virtual_server_id_settings"`
}

type applicationOptionsIdpSigningKeyResourceModelV1 struct {
	Algorithm types.String                 `tfsdk:"algorithm"`
	KeyId     pingonetypes.ResourceIDValue `tfsdk:"key_id"`
}

type applicationSAMLOptionsSpEncryptionResourceModelV1 struct {
	Algorithm   types.String `tfsdk:"algorithm"`
	Certificate types.Object `tfsdk:"certificate"`
}

type applicationSAMLOptionsSpEncryptionCertificateResourceModelV1 struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type applicationSAMLOptionsSpVerificationResourceModelV1 struct {
	CertificateIds     types.Set  `tfsdk:"certificate_ids"`
	AuthnRequestSigned types.Bool `tfsdk:"authn_request_signed"`
}

type applicationSAMLOptionsVirtualServerIdSettingsResourceModelV1 struct {
	Enabled          types.Bool `tfsdk:"enabled"`
	VirtualServerIds types.List `tfsdk:"virtual_server_ids"`
}

type applicationSAMLOptionsVirtualServerIdSettingsVirtualServerIdsResourceModelV1 struct {
	VsId    types.String `tfsdk:"vs_id"`
	Default types.Bool   `tfsdk:"default"`
}

type applicationWSFedOptionsResourceModelV1 struct {
	AudienceRestriction         types.String `tfsdk:"audience_restriction"`
	CorsSettings                types.Object `tfsdk:"cors_settings"`
	DomainName                  types.String `tfsdk:"domain_name"`
	IdpSigningKey               types.Object `tfsdk:"idp_signing_key"`
	Kerberos                    types.Object `tfsdk:"kerberos"`
	ReplyUrl                    types.String `tfsdk:"reply_url"`
	SloEndpoint                 types.String `tfsdk:"slo_endpoint"`
	SubjectNameIdentifierFormat types.String `tfsdk:"subject_name_identifier_format"`
	Type                        types.String `tfsdk:"type"`
}

type applicationWSFedKerberosResourceModelV1 struct {
	Gateways types.Set `tfsdk:"gateways"`
}

type applicationWSFedKerberosGatewayResourceModelV1 struct {
	Id       pingonetypes.ResourceIDValue `tfsdk:"id"`
	Type     types.String                 `tfsdk:"type"`
	UserType types.Object                 `tfsdk:"user_type"`
}

type applicationWSFedGatewayUserTypeRersourceModelV1 struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

var (
	applicationCorsSettingsTFObjectTypes = map[string]attr.Type{
		"behavior": types.StringType,
		"origins":  types.SetType{ElemType: types.StringType},
	}

	applicationIdpSigningKeyTFObjectTypes = map[string]attr.Type{
		"algorithm": types.StringType,
		"key_id":    pingonetypes.ResourceIDType{},
	}

	applicationOidcOptionsTFObjectTypes = utils.MergeAttributeTypeMapsRtn(
		beta.ApplicationOidcOptionsTFObjectTypes,
		map[string]attr.Type{
			"additional_refresh_token_replay_protection_enabled": types.BoolType,
			"allow_wildcard_in_redirect_uris":                    types.BoolType,
			"certificate_based_authentication":                   types.ObjectType{AttrTypes: applicationOidcOptionsCertificateAuthenticationTFObjectTypes},
			"cors_settings":                                      types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes},
			"device_path_id":                                     types.StringType,
			"device_custom_verification_uri":                     types.StringType,
			"device_timeout":                                     types.Int32Type,
			"device_polling_interval":                            types.Int32Type,
			"grant_types":                                        types.SetType{ElemType: types.StringType},
			"home_page_url":                                      types.StringType,
			"idp_signoff":                                        types.BoolType,
			"initiate_login_uri":                                 types.StringType,
			"jwks_url":                                           types.StringType,
			"jwks":                                               types.StringType,
			"mobile_app":                                         types.ObjectType{AttrTypes: applicationOidcMobileAppTFObjectTypes},
			"par_requirement":                                    types.StringType,
			"par_timeout":                                        types.Int32Type,
			"pkce_enforcement":                                   types.StringType,
			"post_logout_redirect_uris":                          types.SetType{ElemType: types.StringType},
			"redirect_uris":                                      types.SetType{ElemType: types.StringType},
			"refresh_token_duration":                             types.Int32Type,
			"refresh_token_rolling_duration":                     types.Int32Type,
			"refresh_token_rolling_grace_period_duration":        types.Int32Type,
			"require_signed_request_object":                      types.BoolType,
			"response_types":                                     types.SetType{ElemType: types.StringType},
			"support_unsigned_request_object":                    types.BoolType,
			"target_link_uri":                                    types.StringType,
			"token_endpoint_auth_method":                         types.StringType,
			"type":                                               types.StringType,
		},
	)

	applicationOidcMobileAppTFObjectTypes = map[string]attr.Type{
		"bundle_id":                types.StringType,
		"huawei_app_id":            types.StringType,
		"huawei_package_name":      types.StringType,
		"integrity_detection":      types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionTFObjectTypes},
		"package_name":             types.StringType,
		"passcode_refresh_seconds": types.Int32Type,
		"universal_app_link":       types.StringType,
	}

	applicationOidcMobileAppIntegrityDetectionTFObjectTypes = map[string]attr.Type{
		"cache_duration":     types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes},
		"enabled":            types.BoolType,
		"excluded_platforms": types.SetType{ElemType: types.StringType},
		"google_play":        types.ObjectType{AttrTypes: applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes},
	}

	applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes = map[string]attr.Type{
		"amount": types.Int32Type,
		"units":  types.StringType,
	}

	applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes = map[string]attr.Type{
		"decryption_key":                   types.StringType,
		"service_account_credentials_json": jsontypes.NormalizedType{},
		"verification_key":                 types.StringType,
		"verification_type":                types.StringType,
	}

	applicationOidcOptionsCertificateAuthenticationTFObjectTypes = map[string]attr.Type{
		"key_id": pingonetypes.ResourceIDType{},
	}

	applicationSamlOptionsTFObjectTypes = map[string]attr.Type{
		"acs_urls":                         types.SetType{ElemType: types.StringType},
		"assertion_duration":               types.Int32Type,
		"assertion_signed_enabled":         types.BoolType,
		"cors_settings":                    types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes},
		"enable_requested_authn_context":   types.BoolType,
		"home_page_url":                    types.StringType,
		"idp_signing_key":                  types.ObjectType{AttrTypes: applicationIdpSigningKeyTFObjectTypes},
		"default_target_url":               types.StringType,
		"nameid_format":                    types.StringType,
		"response_is_signed":               types.BoolType,
		"session_not_on_or_after_duration": types.Int32Type,
		"slo_binding":                      types.StringType,
		"slo_endpoint":                     types.StringType,
		"slo_response_endpoint":            types.StringType,
		"slo_window":                       types.Int32Type,
		"sp_encryption":                    types.ObjectType{AttrTypes: applicationSamlOptionsSpEncryptionTFObjectTypes},
		"sp_entity_id":                     types.StringType,
		"sp_verification":                  types.ObjectType{AttrTypes: applicationSamlOptionsSpVerificationTFObjectTypes},
		"type":                             types.StringType,
		"virtual_server_id_settings":       types.ObjectType{AttrTypes: applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes},
	}

	applicationSamlOptionsSpEncryptionTFObjectTypes = map[string]attr.Type{
		"algorithm":   types.StringType,
		"certificate": types.ObjectType{AttrTypes: applicationSamlOptionsSpEncryptionCertificateTFObjectTypes},
	}

	applicationSamlOptionsSpEncryptionCertificateTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	applicationSamlOptionsSpVerificationTFObjectTypes = map[string]attr.Type{
		"authn_request_signed": types.BoolType,
		"certificate_ids":      types.SetType{ElemType: pingonetypes.ResourceIDType{}},
	}

	applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes = map[string]attr.Type{
		"enabled":            types.BoolType,
		"virtual_server_ids": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationSamlOptionsVirtualServerIdSettingsVirtualServerIdsTFObjectTypes}},
	}

	applicationSamlOptionsVirtualServerIdSettingsVirtualServerIdsTFObjectTypes = map[string]attr.Type{
		"vs_id":   types.StringType,
		"default": types.BoolType,
	}

	applicationExternalLinkOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url": types.StringType,
	}

	applicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"groups": types.SetType{ElemType: pingonetypes.ResourceIDType{}},
		"type":   types.StringType,
	}

	applicationWsfedOptionsTFObjectTypes = map[string]attr.Type{
		"audience_restriction":           types.StringType,
		"cors_settings":                  types.ObjectType{AttrTypes: applicationCorsSettingsTFObjectTypes},
		"domain_name":                    types.StringType,
		"idp_signing_key":                types.ObjectType{AttrTypes: applicationIdpSigningKeyTFObjectTypes},
		"kerberos":                       types.ObjectType{AttrTypes: applicationWsfedOptionsKerberosTFObjectTypes},
		"reply_url":                      types.StringType,
		"slo_endpoint":                   types.StringType,
		"subject_name_identifier_format": types.StringType,
		"type":                           types.StringType,
	}

	applicationWsfedOptionsKerberosTFObjectTypes = map[string]attr.Type{
		"gateways": types.SetType{ElemType: types.ObjectType{AttrTypes: applicationWsfedOptionsKerberosGatewayTFObjectTypes}},
	}

	applicationWsfedOptionsKerberosGatewayTFObjectTypes = map[string]attr.Type{
		"id":        pingonetypes.ResourceIDType{},
		"type":      types.StringType,
		"user_type": types.ObjectType{AttrTypes: applicationWsfedOptionsKerberosGatewayUserTypeTFObjectTypes},
	}

	applicationWsfedOptionsKerberosGatewayUserTypeTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}
)

// Framework interfaces
var (
	_ resource.Resource                   = &ApplicationResource{}
	_ resource.ResourceWithConfigure      = &ApplicationResource{}
	_ resource.ResourceWithImportState    = &ApplicationResource{}
	_ resource.ResourceWithValidateConfig = &ApplicationResource{}
	_ resource.ResourceWithUpgradeState   = &ApplicationResource{}
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
	).AllowedValuesEnum(management.AllowedEnumApplicationTagsEnumValues).ConflictsWith([]string{"external_link_options", "saml_options", "wsfed_options"})

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

	appTypesExactlyOneOf := []string{"external_link_options", "oidc_options", "saml_options", "wsfed_options"}

	externalLinkOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies External link application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf).RequiresReplaceNestedAttributes()

	externalLinkOptionsHomePageURLDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.",
	)

	oidcOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OIDC/OAuth application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf).RequiresReplaceNestedAttributes()

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

	const oidcOptionsDevicePathIdMin = 1
	const oidcOptionsDevicePathIdMax = 50
	oidcOptionsDevicePathIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies a unique identifier within an environment for a device authorization grant flow to provide a short identifier to the application. This property is ignored when the `device_custom_verification_uri` property is configured. The string can contain any letters, numbers, and some special characters (regex: `a-zA-Z0-9_-`). It can have a length of no more than `%[2]d` characters (min/max=`%[1]d`/`%[2]d`).", oidcOptionsDevicePathIdMin, oidcOptionsDevicePathIdMax),
	)

	oidcOptionsDeviceCustomVerificationUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an optional custom verification URI that is returned for the `/device_authorization` endpoint.",
	)

	const oidcOptionsDeviceTimeoutDefault = 600
	const oidcOptionsDeviceTimeoutMin = 1
	const oidcOptionsDeviceTimeoutMax = 3600
	oidcOptionsDeviceTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the length of time (in seconds) that the `userCode` and `deviceCode` returned by the `/device_authorization` endpoint are valid. This property is required only for applications in which the `grant_types` property is set to `%[1]s`. The default value is `%[2]d` seconds. It can have a value of no more than `%[4]d` seconds (min/max=`%[3]d`/`%[4]d`).", management.ENUMAPPLICATIONOIDCGRANTTYPE_DEVICE_CODE, oidcOptionsDeviceTimeoutDefault, oidcOptionsDeviceTimeoutMin, oidcOptionsDeviceTimeoutMax),
	)

	const oidcOptionsDevicePollingIntervalDefault = 5
	const oidcOptionsDevicePollingIntervalMin = 1
	const oidcOptionsDevicePollingIntervalMax = 60
	oidcOptionsDevicePollingIntervalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the frequency (in seconds) for the client to poll the `/as/token` endpoint. This property is required only for applications in which the `grant_types` property is set to `%[1]s`. The default value is `%[2]d` seconds. It can have a value of no more than `%[4]d` seconds (min/max=`%[3]d`/`%[4]d`).", management.ENUMAPPLICATIONOIDCGRANTTYPE_DEVICE_CODE, oidcOptionsDevicePollingIntervalDefault, oidcOptionsDevicePollingIntervalMin, oidcOptionsDevicePollingIntervalMax),
	)

	oidcOptionsIdpSignoffDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean flag to allow signoff without access to the session token cookie.",
	).DefaultValue(false)

	oidcOptionsInitiateLoginUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the URI to use for third-parties to begin the sign-on process for the application. If specified, PingOne redirects users to this URI to initiate SSO to PingOne. The application is responsible for implementing the relevant OIDC flow when the initiate login URI is requested. This property is required if you want the application to appear in the PingOne Application Portal. See the OIDC specification section of [Initiating Login from a Third Party](https://openid.net/specs/openid-connect-core-1_0.html#ThirdPartyInitiatedLogin) for more information.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.",
	)

	oidcOptionsJwksDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	).ConflictsWith([]string{"jwks_url"})

	oidcOptionsJwksUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).",
	).ConflictsWith([]string{"jwks"})

	oidcOptionsTargetLinkUriDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URI for the application. If specified, PingOne will redirect application users to this URI after a user is authenticated. In the PingOne admin console, this becomes the value of the `target_link_uri` parameter used for the Initiate Single Sign-On URL field.  Both `http://` and `https://` URLs are permitted as well as custom mobile native schema (e.g., `org.bxretail.app://target`).",
	)

	oidcOptionsGrantTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list that specifies the grant type for the authorization request.",
	).AllowedValuesEnum(management.AllowedEnumApplicationOIDCGrantTypeEnumValues)

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
		fmt.Sprintf("An integer that specifies the lifetime in seconds of the refresh token. Valid values are between `%d` and `%d`. If the `refresh_token_rolling_duration` property is specified for the application, then this property value must be less than or equal to the value of `refresh_token_rolling_duration`. After this property is set, the value cannot be nullified - this will reset the value back to the default. This value is used to generate the value for the exp claim when minting a new refresh token.", oidcOptionsRefreshTokenDurationMin, oidcOptionsRefreshTokenDurationMax),
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

	oidcOptionsSupportUnsignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If `false` or null, an unsigned request object is not allowed.",
	).DefaultValue(false)

	oidcOptionsRequireSignedRequestObjectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
	).DefaultValue(false)

	oidcOptionsCertificateBasedAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that specifies Certificate based authentication settings. This parameter block can only be set where the application's `type` parameter is set to `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
	)

	oidcOptionsCertificateBasedAuthenticationKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents a PingOne ID for the issuance certificate key.  The key must be of type `ISSUANCE`.  Must be a valid PingOne Resource ID.",
	)

	oidcOptionsMobileAppDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that specifies Mobile application integration settings for `%s` type applications.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
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
		"A single object that specifies settings for the caching duration of successful integrity detection calls.  Every attestation request entails a certain time tradeoff. You can choose to cache successful integrity detection calls for a predefined duration, between a minimum of 1 minute and a maximum of 48 hours. If integrity detection is ENABLED, the cache duration must be set.",
	)

	oidcOptionsMobileAppIntegrityDetectionCacheDurationUnitsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time units of the cache `amount` parameter.",
	).AllowedValuesEnum(management.AllowedEnumDurationUnitMinsHoursEnumValues).DefaultValue(string(management.ENUMDURATIONUNITMINSHOURS_MINUTES))

	oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that describes Google Play Integrity API credential settings for Android device integrity detection.  Required when `excluded_platforms` is unset or does not include `%s`.", management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE),
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
		"A single object that specifies SAML application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf).RequiresReplaceNestedAttributes()

	samlOptionsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type associated with the application.",
	).AllowedValues(
		string(management.ENUMAPPLICATIONTYPE_WEB_APP),
		string(management.ENUMAPPLICATIONTYPE_CUSTOM_APP),
	).DefaultValue(string(management.ENUMAPPLICATIONTYPE_WEB_APP)).RequiresReplace()

	samlOptionsAssertionSignedEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML assertion itself should be signed.",
	).DefaultValue(true)

	samlOptionsDefaultTargetUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).",
	)

	samlOptionsEnableRequestedAuthnContextDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
	)

	samlOptionsResponseIsSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML assertion response itself should be signed.",
	).DefaultValue(false)

	const samlOptionsSessionNotOnOrAfterDurationMin = 60
	samlOptionsSessionNotOnOrAfterDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies a value for if the SAML application requires a different `SessionNotOnOrAfter` attribute value within the `AuthnStatement` element than the `NotOnOrAfter` value set by the `assertion_duration` property.",
	)

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

	samlOptionsIdpSigningKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"SAML application assertion/response signing key settings.  Use with `assertion_signed_enabled` to enable assertion signing and/or `response_is_signed` to enable response signing.  It's highly recommended, and best practice, to define signing key settings for the configured SAML application.  However if this property is omitted, the default signing certificate for the environment is used.  This parameter will become a required field in the next major release of the provider.",
	)

	idpSigningKeyAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
	)

	samlSpEncryptionAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The algorithm to use when encrypting assertions.",
	).AllowedValuesEnum(management.AllowedEnumCertificateKeyEncryptionAlgorithmEnumValues)

	samlOptionsSpVerificationAuthnRequestSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the Authn Request signing should be enforced.",
	).DefaultValue(false)

	wsfedOptionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies WS-Fed application specific settings.",
	).ExactlyOneOf(appTypesExactlyOneOf).RequiresReplaceNestedAttributes()

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne application (SAML, OpenID Connect, External Link, WS-Fed) in an environment.",

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
						path.MatchRelative().AtParent().AtName("wsfed_options"),
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

			"icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings for the application icon.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID for the application icon.  Must be a valid PingOne Resource ID.").Description,
						Required:    true,

						CustomType: pingonetypes.ResourceIDType{},
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

			"access_control_group_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies group access control settings.").Description,
				Optional:    true,

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

						ElementType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"external_link_options": schema.SingleNestedAttribute{
				Description:         externalLinkOptionsDescription.Description,
				MarkdownDescription: externalLinkOptionsDescription.MarkdownDescription,
				Optional:            true,

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

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
						path.MatchRelative().AtParent().AtName("wsfed_options"),
					),
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"oidc_options": schema.SingleNestedAttribute{
				Description:         oidcOptionsDescription.Description,
				MarkdownDescription: oidcOptionsDescription.MarkdownDescription,
				Optional:            true,

				Attributes: utils.MergeResourceSchemaAttributeMapsRtn(
					beta.ResourceSchemaItems(),
					map[string]schema.Attribute{
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

						"device_path_id": schema.StringAttribute{
							Description:         oidcOptionsDevicePathIdDescription.Description,
							MarkdownDescription: oidcOptionsDevicePathIdDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9_-]*`), "The string can contain any letters, numbers, underscore and dash characters"),
								stringvalidator.LengthBetween(oidcOptionsDevicePathIdMin, oidcOptionsDevicePathIdMax),
							},
						},

						"device_custom_verification_uri": schema.StringAttribute{
							Description:         oidcOptionsDeviceCustomVerificationUriDescription.Description,
							MarkdownDescription: oidcOptionsDeviceCustomVerificationUriDescription.MarkdownDescription,
							Optional:            true,
						},

						"device_timeout": schema.Int32Attribute{
							Description:         oidcOptionsDeviceTimeoutDescription.Description,
							MarkdownDescription: oidcOptionsDeviceTimeoutDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int32default.StaticInt32(oidcOptionsDeviceTimeoutDefault),

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsDeviceTimeoutMin, oidcOptionsDeviceTimeoutMax),
							},
						},

						"device_polling_interval": schema.Int32Attribute{
							Description:         oidcOptionsDevicePollingIntervalDescription.Description,
							MarkdownDescription: oidcOptionsDevicePollingIntervalDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int32default.StaticInt32(oidcOptionsDevicePollingIntervalDefault),

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsDevicePollingIntervalMin, oidcOptionsDevicePollingIntervalMax),
							},
						},

						"idp_signoff": schema.BoolAttribute{
							Description:         oidcOptionsIdpSignoffDescription.Description,
							MarkdownDescription: oidcOptionsIdpSignoffDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
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

						"token_endpoint_auth_method": schema.StringAttribute{
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

						"par_timeout": schema.Int32Attribute{
							Description:         oidcOptionsParTimeoutDescription.Description,
							MarkdownDescription: oidcOptionsParTimeoutDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int32default.StaticInt32(oidcOptionsParTimeoutDefault),

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsParTimeoutMin, oidcOptionsParTimeoutMax),
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

						"allow_wildcard_in_redirect_uris": schema.BoolAttribute{
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

						"refresh_token_duration": schema.Int32Attribute{
							Description:         oidcOptionsRefreshTokenDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int32default.StaticInt32(oidcOptionsRefreshTokenDurationDefault),

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsRefreshTokenDurationMin, oidcOptionsRefreshTokenDurationMax),
							},
						},

						"refresh_token_rolling_duration": schema.Int32Attribute{
							Description:         oidcOptionsRefreshTokenRollingDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenRollingDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int32default.StaticInt32(oidcOptionsRefreshTokenRollingDurationDefault),

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsRefreshTokenRollingDurationMin, oidcOptionsRefreshTokenRollingDurationMax),
							},
						},

						"refresh_token_rolling_grace_period_duration": schema.Int32Attribute{
							Description:         oidcOptionsRefreshTokenRollingGracePeriodDurationDescription.Description,
							MarkdownDescription: oidcOptionsRefreshTokenRollingGracePeriodDurationDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.Int32{
								int32validator.Between(oidcOptionsRefreshTokenRollingGracePeriodDurationMin, oidcOptionsRefreshTokenRollingGracePeriodDurationMax),
							},
						},

						"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
							Description:         oidcOptionsAdditionalRefreshTokenReplayProtectionEnabledDescription.Description,
							MarkdownDescription: oidcOptionsAdditionalRefreshTokenReplayProtectionEnabledDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(true),
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

						"certificate_based_authentication": schema.SingleNestedAttribute{
							Description:         oidcOptionsCertificateBasedAuthenticationDescription.Description,
							MarkdownDescription: oidcOptionsCertificateBasedAuthenticationDescription.MarkdownDescription,
							Optional:            true,

							Attributes: map[string]schema.Attribute{
								"key_id": schema.StringAttribute{
									Description:         oidcOptionsCertificateBasedAuthenticationKeyIdDescription.Description,
									MarkdownDescription: oidcOptionsCertificateBasedAuthenticationKeyIdDescription.MarkdownDescription,
									Required:            true,

									CustomType: pingonetypes.ResourceIDType{},
								},
							},
						},

						"mobile_app": schema.SingleNestedAttribute{
							Description:         oidcOptionsMobileAppDescription.Description,
							MarkdownDescription: oidcOptionsMobileAppDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

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

								"passcode_refresh_seconds": schema.Int32Attribute{
									Description:         oidcOptionsMobileAppPasscodeRefreshSecondsDescription.Description,
									MarkdownDescription: oidcOptionsMobileAppPasscodeRefreshSecondsDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Default: int32default.StaticInt32(oidcOptionsMobileAppPasscodeRefreshSecondsDefault),

									Validators: []validator.Int32{
										int32validator.Between(oidcOptionsMobileAppPasscodeRefreshSecondsMin, oidcOptionsMobileAppPasscodeRefreshSecondsMax),
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

								"integrity_detection": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies mobile application integrity detection settings.").Description,
									Optional:    true,
									Computed:    true,

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

										"cache_duration": schema.SingleNestedAttribute{
											Description:         oidcOptionsMobileAppIntegrityDetectionCacheDurationDescription.Description,
											MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionCacheDurationDescription.MarkdownDescription,
											Optional:            true,

											Attributes: map[string]schema.Attribute{
												"amount": schema.Int32Attribute{
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

										"google_play": schema.SingleNestedAttribute{
											Description:         oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription.Description,
											MarkdownDescription: oidcOptionsMobileAppIntegrityDetectionGooglePlayDescription.MarkdownDescription,
											Optional:            true,

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

													CustomType: jsontypes.NormalizedType{},

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
									},
								},
							},
						},

						"cors_settings": resourceApplicationSchemaCorsSettings(),
					}),

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
						path.MatchRelative().AtParent().AtName("wsfed_options"),
					),
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"saml_options": schema.SingleNestedAttribute{
				Description:         samlOptionsDescription.Description,
				MarkdownDescription: samlOptionsDescription.MarkdownDescription,
				Optional:            true,

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

					"assertion_duration": schema.Int32Attribute{
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

					"session_not_on_or_after_duration": schema.Int32Attribute{
						Description:         samlOptionsSessionNotOnOrAfterDurationDescription.Description,
						MarkdownDescription: samlOptionsSessionNotOnOrAfterDurationDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Int32{
							int32validator.AtLeast(samlOptionsSessionNotOnOrAfterDurationMin),
						},
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
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Expected value to have a url with schema of \"https\"."),
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

					"slo_window": schema.Int32Attribute{
						Description:         samlOptionsSloWindowDescription.Description,
						MarkdownDescription: samlOptionsSloWindowDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Int32{
							int32validator.Between(samlOptionsSloWindowMin, samlOptionsSloWindowMax),
						},
					},

					"sp_encryption": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings for PingOne to encrypt SAML assertions to be sent to the application. Assertions are not encrypted by default.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description:         samlSpEncryptionAlgorithmDescription.Description,
								MarkdownDescription: samlSpEncryptionAlgorithmDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeyEncryptionAlgorithmEnumValues)...),
								},
							},

							"certificate": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the certificate settings used to encrypt SAML assertions.").Description,
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique identifier of the encryption public certificate that has been uploaded to PingOne.").Description,
										Required:    true,

										CustomType: pingonetypes.ResourceIDType{},
									},
								},
							},
						},
					},

					"sp_entity_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.").Description,
						Required:    true,
					},

					"idp_signing_key": schema.SingleNestedAttribute{
						Description:         samlOptionsIdpSigningKeyDescription.Description,
						MarkdownDescription: samlOptionsIdpSigningKeyDescription.MarkdownDescription,
						Required:            true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description:         idpSigningKeyAlgorithmDescription.Description,
								MarkdownDescription: idpSigningKeyAlgorithmDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeySignagureAlgorithmEnumValues)...),
								},
							},

							"key_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.  Must be a valid PingOne resource ID.").Description,
								Required:    true,

								CustomType: pingonetypes.ResourceIDType{},
							},
						},
					},

					"sp_verification": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object item that specifies SP signature verification settings.").Description,
						Optional:    true,

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
								ElementType: pingonetypes.ResourceIDType{},
								Required:    true,
							},
						},
					},

					"cors_settings": resourceApplicationSchemaCorsSettings(),

					"virtual_server_id_settings": schema.SingleNestedAttribute{
						Description: "Contains the virtual server ID or IDs to be used.",
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Indicates whether the virtual server ID or IDs specified are to be used. Defaults to `false`.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"virtual_server_ids": schema.ListNestedAttribute{
								Description: "Required if `enabled` is `true`. Contains the list of virtual server ID or IDs to be used.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"vs_id": schema.StringAttribute{
											Description: "This must be a valid SAML entity ID.",
											Required:    true,
										},
										"default": schema.BoolAttribute{
											Description: "Indicates whether the virtual server identified by the associated `vs_id` is to be used as the default virtual server. Defaults to `false`.",
											Optional:    true,
											Computed:    true,
											Default:     booldefault.StaticBool(false),
										},
									},
								},
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.UniqueValues(),
									// Ensure that the virtual server ID settings are only required if the `enabled` field is set to true.
									listvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue("true"),
										path.MatchRelative().AtParent().AtName("enabled"),
									),
								},
							},
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
						path.MatchRelative().AtParent().AtName("wsfed_options"),
					),
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"wsfed_options": schema.SingleNestedAttribute{
				Description:         wsfedOptionsDescription.Description,
				MarkdownDescription: wsfedOptionsDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"audience_restriction": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "The service provider ID. The default value is \"urn:federation:MicrosoftOnline\".",
						MarkdownDescription: "The service provider ID. The default value is `urn:federation:MicrosoftOnline`.",
						Default:             stringdefault.StaticString("urn:federation:MicrosoftOnline"),
					},
					"cors_settings": resourceApplicationSchemaCorsSettings(),
					"domain_name": schema.StringAttribute{
						Required:    true,
						Description: "The federated domain name (for example, the Azure custom domain).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsDomain, "Must be a valid domain name"),
						},
					},
					"idp_signing_key": schema.SingleNestedAttribute{
						Description: "Contains the information about the signing of requests by the identity provider (IdP).",
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description:         idpSigningKeyAlgorithmDescription.Description,
								MarkdownDescription: idpSigningKeyAlgorithmDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeySignagureAlgorithmEnumValues)...),
								},
							},

							"key_id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.  Must be a valid PingOne resource ID.").Description,
								Required:    true,

								CustomType: pingonetypes.ResourceIDType{},
							},
						},
					},
					"kerberos": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"gateways": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:    true,
											Description: "The UUID of the LDAP gateway. Must be a valid PingOne resource ID.",
											CustomType:  pingonetypes.ResourceIDType{},
										},
										"type": schema.StringAttribute{
											Optional:    true,
											Computed:    true,
											Default:     stringdefault.StaticString("LDAP"),
											Description: "The gateway type. This must be \"LDAP\".",
											Validators: []validator.String{
												stringvalidator.OneOf(
													"LDAP",
												),
											},
										},
										"user_type": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Optional:            true,
													Description:         "The UUID of a user type in the list of \"userTypes\" for the LDAP gateway. Must be a valid PingOne resource ID.",
													MarkdownDescription: "The UUID of a user type in the list of `userTypes` for the LDAP gateway. Must be a valid PingOne resource ID.",
													CustomType:          pingonetypes.ResourceIDType{},
												},
											},
											Required:    true,
											Description: "The object reference to the user type in the list of \"userTypes\" for the LDAP gateway.",
										},
									},
								},
								Optional:    true,
								Description: "The LDAP gateway properties.",
							},
						},
						Optional:    true,
						Description: "The Kerberos authentication settings. Leave this out of the configuration to disable Kerberos authentication.",
					},
					"reply_url": schema.StringAttribute{
						Required:    true,
						Description: "The URL that the replying party (such as, Office365) uses to accept submissions of RequestSecurityTokenResponse messages that are a result of SSO requests.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Expected value to have a url with schema of \"http\" or \"https\"."),
						},
					},
					"slo_endpoint": schema.StringAttribute{
						Optional:    true,
						Description: "The single logout endpoint URL.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Expected value to have a url with schema of \"http\" or \"https\"."),
						},
					},
					"subject_name_identifier_format": schema.StringAttribute{
						Optional:            true,
						Description:         "The format to use for the SubjectNameIdentifier element. Options are \"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified\", \"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress\".",
						MarkdownDescription: "The format to use for the SubjectNameIdentifier element. Options are `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified`, `urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
								"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
							),
						},
					},
					"type": schema.StringAttribute{
						Required:            true,
						Description:         "A string that specifies the type associated with the application. This is a required property. Options are \"WEB_APP\", \"CUSTOM_APP\".",
						MarkdownDescription: "A string that specifies the type associated with the application. This is a required property. Options are `WEB_APP`, `CUSTOM_APP`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"WEB_APP",
								"CUSTOM_APP",
							),
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("external_link_options"),
						path.MatchRelative().AtParent().AtName("oidc_options"),
						path.MatchRelative().AtParent().AtName("saml_options"),
						path.MatchRelative().AtParent().AtName("wsfed_options"),
					),
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},
		},
	}
}

func resourceApplicationSchemaCorsSettings() schema.SingleNestedAttribute {

	listDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings.",
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

	return schema.SingleNestedAttribute{
		Description:         listDescription.Description,
		MarkdownDescription: listDescription.MarkdownDescription,
		Optional:            true,

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
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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
	var data applicationResourceModelV1

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.validate(ctx, true)...)
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state applicationResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	resp.Diagnostics.Append(plan.validate(ctx, false)...)
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.CreateApplication(ctx, plan.EnvironmentId.ValueString()).CreateApplicationRequest(*application).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	} else if createResponse.ApplicationWSFED != nil && createResponse.ApplicationWSFED.GetId() != "" {
		applicationId = createResponse.ApplicationWSFED.GetId()
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Cannot determine application ID from API response for application: %s", plan.Name.ValueString()),
			fmt.Sprintf("Full response object: %v\n", resp),
		)
	}

	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, plan.EnvironmentId.ValueString(), applicationId).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *applicationResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state applicationResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	resp.Diagnostics.Append(plan.validate(ctx, false)...)
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateApplicationRequest(*application).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
		applicationWriteCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *applicationResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.DeleteApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteApplication",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func applicationWriteCustomError(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Wildcards in redirect URis
		m, err := regexp.MatchString("^Wildcards are not allowed in redirect URIs.", p1Error.GetMessage())
		if err != nil {
			diags.AddError("API Validation error", "Cannot match error string for wildcard in redirect URIs")
			return diags
		}
		if m {
			diags.AddError("Invalid configuration", "Current configuration is invalid as wildcards are not allowed in redirect URIs.  Wildcards can be enabled by setting `allow_wildcard_in_redirect_uris` to `true`.")

			return diags
		}
	}

	diags.Append(legacysdk.DefaultCustomError(r, p1Error)...)
	return diags
}

func (p *applicationResourceModelV1) validate(ctx context.Context, allowUnknown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	var oidcPlan *applicationOIDCOptionsResourceModelV1
	diags.Append(p.OIDCOptions.As(ctx, &oidcPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: allowUnknown,
	})...)
	if diags.HasError() {
		return diags
	}

	if oidcPlan != nil {
		diags.Append(oidcPlan.validateCertificateBasedAuthentication(allowUnknown)...)
		diags.Append(oidcPlan.validateWildcardInRedirectUri(ctx, allowUnknown)...)
	}

	var samlPlan *applicationSAMLOptionsResourceModelV1
	diags.Append(p.SAMLOptions.As(ctx, &samlPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: allowUnknown,
	})...)
	if diags.HasError() {
		return diags
	}

	if samlPlan != nil {
		diags.Append(samlPlan.validateVirtualServerIdSettings(ctx, allowUnknown)...)
	}

	return diags
}

func (p *applicationOIDCOptionsResourceModelV1) validateCertificateBasedAuthentication(allowUnknown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if p.CertificateBasedAuthentication.IsUnknown() && !allowUnknown {
		diags.AddAttributeError(
			path.Root("oidc_options").AtName("certificate_based_authentication"),
			"Invalid configuration",
			"Current configuration is invalid as the `oidc_options.certificate_based_authentication` value is unknown, cannot validate.",
		)
	}

	// Certificate based authentication
	if !p.CertificateBasedAuthentication.IsNull() && !p.CertificateBasedAuthentication.IsUnknown() {
		if !p.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_NATIVE_APP))) {
			diags.AddAttributeError(
				path.Root("oidc_options").AtName("certificate_based_authentication"),
				"Invalid configuration",
				fmt.Sprintf("`certificate_based_authentication` can only be set with OIDC applications that have a `type` value of `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
			)
		}
	}

	return diags
}

func (p *applicationOIDCOptionsResourceModelV1) validateWildcardInRedirectUri(ctx context.Context, allowUnknown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if p.RedirectUris.IsUnknown() && !allowUnknown {
		diags.AddAttributeError(
			path.Root("oidc_options").AtName("redirect_uris"),
			"Invalid configuration",
			"Current configuration is invalid as the `oidc_options.redirect_uris` value is unknown, and cannot validate wildcards.",
		)
	}

	if p.AllowWildcardsInRedirectUris.IsUnknown() && !allowUnknown {
		diags.AddAttributeError(
			path.Root("oidc_options").AtName("allow_wildcard_in_redirect_uris"),
			"Invalid configuration",
			"Current configuration is invalid as `oidc_options.allow_wildcard_in_redirect_uris` value is unknown, and cannot validate wildcards presence in `oidc_options.redirect_uris`.",
		)
	}
	if diags.HasError() {
		return diags
	}

	if !p.RedirectUris.IsNull() && !p.RedirectUris.IsUnknown() {
		var uris []types.String
		diags.Append(p.RedirectUris.ElementsAs(ctx, &uris, false)...)
		if diags.HasError() {
			return diags
		}

		for _, uri := range uris {
			if !uri.IsNull() && !uri.IsUnknown() {
				if strings.Contains(uri.ValueString(), "*") {
					if p.AllowWildcardsInRedirectUris.IsNull() || p.AllowWildcardsInRedirectUris.Equal(types.BoolValue(false)) {
						diags.AddAttributeError(
							path.Root("oidc_options").AtName("redirect_uris"),
							"Invalid configuration",
							"Current configuration is invalid as wildcards are not allowed in redirect URIs.  Wildcards can be enabled by setting `allow_wildcard_in_redirect_uris` to `true`.",
						)
						break
					}
				}
			}
		}
	}

	return diags
}

func (p *applicationSAMLOptionsResourceModelV1) validateVirtualServerIdSettings(ctx context.Context, allowUnknown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if p.VirtualServerIdSettings.IsUnknown() && !allowUnknown {
		diags.AddAttributeError(
			path.Root("saml_options").AtName("virtual_server_id_settings"),
			"Invalid configuration",
			"Current configuration is invalid as the `saml_options.virtual_server_id_settings` value is unknown, cannot validate.",
		)
		return diags
	}

	if !p.VirtualServerIdSettings.IsNull() && !p.VirtualServerIdSettings.IsUnknown() {
		var vsSettings applicationSAMLOptionsVirtualServerIdSettingsResourceModelV1
		diags.Append(p.VirtualServerIdSettings.As(ctx, &vsSettings, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}
		var vsIds []applicationSAMLOptionsVirtualServerIdSettingsVirtualServerIdsResourceModelV1
		diags.Append(vsSettings.VirtualServerIds.ElementsAs(ctx, &vsIds, false)...)
		if diags.HasError() {
			return diags
		}

		defaultCount := 0
		for _, vsId := range vsIds {
			if !vsId.Default.IsNull() && vsId.Default.ValueBool() {
				defaultCount++
			}
		}

		if defaultCount == 0 {
			diags.AddAttributeError(
				path.Root("saml_options").AtName("virtual_server_id_settings").AtName("virtual_server_ids"),
				"Invalid Configuration",
				"At least one Virtual Server ID must be set as default.",
			)
		} else if defaultCount > 1 {
			diags.AddAttributeError(
				path.Root("saml_options").AtName("virtual_server_id_settings").AtName("virtual_server_ids"),
				"Invalid Configuration",
				"Only one Virtual Server ID can be set as default.",
			)
		}
	}

	return diags
}

func (p *applicationResourceModelV1) expandCreate(ctx context.Context) (*management.CreateApplicationRequest, diag.Diagnostics) {
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

	if !p.WSFedOptions.IsNull() && !p.WSFedOptions.IsUnknown() {
		data.ApplicationWSFED, d = p.expandApplicationWSFed(ctx)
		diags = append(diags, d...)
	}

	return data, diags
}

func (p *applicationResourceModelV1) expandUpdate(ctx context.Context) (*management.UpdateApplicationRequest, diag.Diagnostics) {
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

	if !p.WSFedOptions.IsNull() && !p.WSFedOptions.IsUnknown() {
		data.ApplicationWSFED, d = p.expandApplicationWSFed(ctx)
		diags = append(diags, d...)
	}

	return data, diags
}

func (p *applicationCorsSettingsResourceModelV1) expand() (*management.ApplicationCorsSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationCorsSettings(management.EnumApplicationCorsSettingsBehavior(p.Behavior.ValueString()))

	if !p.Origins.IsNull() && !p.Origins.IsUnknown() {
		var originsPlan []types.String
		d := p.Origins.ElementsAs(context.Background(), &originsPlan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		origins, d := framework.TFTypeStringSliceToStringSlice(originsPlan, path.Root("cors_settings").AtName("origins"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetOrigins(origins)
	}

	return data, diags
}

func (p *applicationResourceModelV1) expandApplicationOIDC(ctx context.Context) (*management.ApplicationOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationOIDC

	if !p.OIDCOptions.IsNull() && !p.OIDCOptions.IsUnknown() {
		var plan applicationOIDCOptionsResourceModelV1
		d := p.OIDCOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		grantTypes := make([]management.EnumApplicationOIDCGrantType, 0)

		var grantTypesPlan []types.String

		diags.Append(plan.GrantTypes.ElementsAs(ctx, &grantTypesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		grantTypesStr, d := framework.TFTypeStringSliceToStringSlice(grantTypesPlan, path.Root("oidc_options").AtName("grant_types"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		for _, v := range grantTypesStr {
			grantTypes = append(grantTypes, management.EnumApplicationOIDCGrantType(v))
		}

		data = management.NewApplicationOIDC(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_OPENID_CONNECT,
			management.EnumApplicationType(plan.Type.ValueString()),
			grantTypes,
			management.EnumApplicationOIDCTokenAuthMethod(plan.TokenEndpointAuthnMethod.ValueString()),
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

		if !plan.CorsSettings.IsNull() && !plan.CorsSettings.IsUnknown() {
			var corsPlan applicationCorsSettingsResourceModelV1

			diags.Append(plan.CorsSettings.As(ctx, &corsPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			corsSettings, d := corsPlan.expand()
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetCorsSettings(*corsSettings)
		}

		if !plan.HomePageUrl.IsNull() && !plan.HomePageUrl.IsUnknown() {
			data.SetHomePageUrl(plan.HomePageUrl.ValueString())
		}

		if !plan.DevicePathId.IsNull() && !plan.DevicePathId.IsUnknown() {
			data.SetDevicePathId(plan.DevicePathId.ValueString())
		}

		if !plan.DeviceCustomVerificationUri.IsNull() && !plan.DeviceCustomVerificationUri.IsUnknown() {
			data.SetDeviceCustomVerificationUri(plan.DeviceCustomVerificationUri.ValueString())
		}

		if !plan.DeviceTimeout.IsNull() && !plan.DeviceTimeout.IsUnknown() {
			data.SetDeviceTimeout(plan.DeviceTimeout.ValueInt32())
		}

		if !plan.DevicePollingInterval.IsNull() && !plan.DevicePollingInterval.IsUnknown() {
			data.SetDevicePollingInterval(plan.DevicePollingInterval.ValueInt32())
		}

		if !plan.IdpSignoff.IsNull() && !plan.IdpSignoff.IsUnknown() {
			data.SetIdpSignoff(plan.IdpSignoff.ValueBool())
		}

		if !plan.InitiateLoginUri.IsNull() && !plan.InitiateLoginUri.IsUnknown() {
			data.SetInitiateLoginUri(plan.InitiateLoginUri.ValueString())
		}

		if !plan.Jwks.IsNull() && !plan.Jwks.IsUnknown() {
			data.SetJwks(plan.Jwks.ValueString())
		}

		if !plan.JwksUrl.IsNull() && !plan.JwksUrl.IsUnknown() {
			data.SetJwksUrl(plan.JwksUrl.ValueString())
		}

		if !plan.TargetLinkUri.IsNull() && !plan.TargetLinkUri.IsUnknown() {
			data.SetTargetLinkUri(plan.TargetLinkUri.ValueString())
		}

		if !plan.ResponseTypes.IsNull() && !plan.ResponseTypes.IsUnknown() {
			var responseTypesPlan []types.String

			diags.Append(plan.ResponseTypes.ElementsAs(ctx, &responseTypesPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			responseTypesStr, d := framework.TFTypeStringSliceToStringSlice(responseTypesPlan, path.Root("oidc_options").AtName("response_types"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			obj := make([]management.EnumApplicationOIDCResponseType, 0)

			for _, v := range responseTypesStr {
				obj = append(obj, management.EnumApplicationOIDCResponseType(v))
			}
			data.SetResponseTypes(obj)
		}

		if !plan.ParRequirement.IsNull() && !plan.ParRequirement.IsUnknown() {
			data.SetParRequirement(management.EnumApplicationOIDCPARRequirement(plan.ParRequirement.ValueString()))
		}

		if !plan.ParTimeout.IsNull() && !plan.ParTimeout.IsUnknown() {
			data.SetParTimeout(plan.ParTimeout.ValueInt32())
		}

		if !plan.PKCEEnforcement.IsNull() && !plan.PKCEEnforcement.IsUnknown() {
			data.SetPkceEnforcement(management.EnumApplicationOIDCPKCEOption(plan.PKCEEnforcement.ValueString()))
		}

		if !plan.RedirectUris.IsNull() && !plan.RedirectUris.IsUnknown() {
			var redirectUrisPlan []types.String
			diags.Append(plan.RedirectUris.ElementsAs(ctx, &redirectUrisPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			redirectUris, d := framework.TFTypeStringSliceToStringSlice(redirectUrisPlan, path.Root("oidc_options").AtName("redirect_uris"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetRedirectUris(redirectUris)
		}

		if !plan.AllowWildcardsInRedirectUris.IsNull() && !plan.AllowWildcardsInRedirectUris.IsUnknown() {
			data.SetAllowWildcardInRedirectUris(plan.AllowWildcardsInRedirectUris.ValueBool())
		}

		if !plan.PostLogoutRedirectUris.IsNull() && !plan.PostLogoutRedirectUris.IsUnknown() {
			var postLogoutRedirectUrisPlan []types.String

			diags.Append(plan.PostLogoutRedirectUris.ElementsAs(ctx, &postLogoutRedirectUrisPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			redirectUris, d := framework.TFTypeStringSliceToStringSlice(postLogoutRedirectUrisPlan, path.Root("oidc_options").AtName("post_logout_redirect_uris"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetPostLogoutRedirectUris(redirectUris)
		}

		if !plan.RefreshTokenDuration.IsNull() && !plan.RefreshTokenDuration.IsUnknown() {
			data.SetRefreshTokenDuration(plan.RefreshTokenDuration.ValueInt32())
		}

		if !plan.RefreshTokenRollingDuration.IsNull() && !plan.RefreshTokenRollingDuration.IsUnknown() {
			data.SetRefreshTokenRollingDuration(plan.RefreshTokenRollingDuration.ValueInt32())
		}

		if !plan.RefreshTokenRollingGracePeriodDuration.IsNull() && !plan.RefreshTokenRollingGracePeriodDuration.IsUnknown() {
			data.SetRefreshTokenRollingGracePeriodDuration(plan.RefreshTokenRollingGracePeriodDuration.ValueInt32())
		}

		if !plan.AdditionalRefreshTokenReplayProtectionEnabled.IsNull() && !plan.AdditionalRefreshTokenReplayProtectionEnabled.IsUnknown() {
			data.SetAdditionalRefreshTokenReplayProtectionEnabled(plan.AdditionalRefreshTokenReplayProtectionEnabled.ValueBool())
		}

		if !p.Tags.IsNull() && !p.Tags.IsUnknown() {
			var tagsPlan []types.String

			diags.Append(p.Tags.ElementsAs(ctx, &tagsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			tagsStr, d := framework.TFTypeStringSliceToStringSlice(tagsPlan, path.Root("tags"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			tags := make([]management.EnumApplicationTags, 0)

			for _, v := range tagsStr {
				tags = append(tags, management.EnumApplicationTags(v))
			}

			data.Tags = tags

		}

		data.SetAssignActorRoles(false)

		if !plan.CertificateBasedAuthentication.IsNull() && !plan.CertificateBasedAuthentication.IsUnknown() {
			if !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_NATIVE_APP))) {
				diags.AddError(
					"Invalid configuration",
					fmt.Sprintf("`certificate_based_authentication` can only be set with applications that have a `type` value of `%s`.", management.ENUMAPPLICATIONTYPE_NATIVE_APP),
				)

				return nil, diags
			}

			var kerberosPlan applicationOIDCCertificateBasedAuthenticationResourceModelV1

			diags.Append(plan.CertificateBasedAuthentication.As(ctx, &kerberosPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetKerberos(*management.NewApplicationOIDCAllOfKerberos(*management.NewApplicationOIDCAllOfKerberosKey(kerberosPlan.KeyId.ValueString())))
		}

		if !plan.SupportUnsignedRequestObject.IsNull() && !plan.SupportUnsignedRequestObject.IsUnknown() {
			data.SetSupportUnsignedRequestObject(plan.SupportUnsignedRequestObject.ValueBool())
		}

		if !plan.RequireSignedRequestObject.IsNull() && !plan.RequireSignedRequestObject.IsUnknown() {
			data.SetRequireSignedRequestObject(plan.RequireSignedRequestObject.ValueBool())
		}

		if !plan.MobileApp.IsNull() && !plan.MobileApp.IsUnknown() {
			var mobileAppPlan applicationOIDCMobileAppResourceModelV1

			diags.Append(plan.MobileApp.As(ctx, &mobileAppPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			mobile, d := mobileAppPlan.expand(ctx)
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetMobile(*mobile)
		}

		beta.AddBeta(data, plan.ApplicationOIDCOptionsResourceModelV1)
	}

	return data, diags
}

func (p *applicationOIDCMobileAppResourceModelV1) expand(ctx context.Context) (*management.ApplicationOIDCAllOfMobile, diag.Diagnostics) {
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

		var integrityDetectionPlan applicationOIDCMobileAppIntegrityDetectionResourceModelV1
		diags.Append(p.IntegrityDetection.As(ctx, &integrityDetectionPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		integrityDetection, d := integrityDetectionPlan.expand(ctx)
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
			p.PasscodeRefreshSeconds.ValueInt32(),
			management.ENUMPASSCODEREFRESHTIMEUNIT_SECONDS,
		))
	}

	if !p.UniversalAppLink.IsNull() && !p.UniversalAppLink.IsUnknown() {
		data.SetUriPrefix(p.UniversalAppLink.ValueString())
	}

	return data, diags
}

func (p *applicationOIDCMobileAppIntegrityDetectionResourceModelV1) expand(ctx context.Context) (*management.ApplicationOIDCAllOfMobileIntegrityDetection, diag.Diagnostics) {
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

	googleVerificationIncluded := true && data.GetMode() == management.ENUMENABLEDSTATUS_ENABLED

	if !p.ExcludedPlatforms.IsNull() && !p.ExcludedPlatforms.IsUnknown() {
		var excludedPlatformsPlan []types.String

		diags.Append(p.ExcludedPlatforms.ElementsAs(ctx, &excludedPlatformsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		excludedPlatformsStr, d := framework.TFTypeStringSliceToStringSlice(excludedPlatformsPlan, path.Root("oidc_options").AtName("mobile_app").AtName("integrity_detection").AtName("excluded_platforms"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		excludedPlatforms := make([]management.EnumMobileIntegrityDetectionPlatform, 0)

		for _, v := range excludedPlatformsStr {
			excludedPlatforms = append(excludedPlatforms, management.EnumMobileIntegrityDetectionPlatform(v))
			if v == string(management.ENUMMOBILEINTEGRITYDETECTIONPLATFORM_GOOGLE) {
				googleVerificationIncluded = false
			}
		}

		data.SetExcludedPlatforms(excludedPlatforms)
	}

	if !p.GooglePlay.IsNull() && !p.GooglePlay.IsUnknown() {

		var googlePlayPlan applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV1
		diags.Append(p.GooglePlay.As(ctx, &googlePlayPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		googlePlay := management.NewApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay()

		if !googlePlayPlan.DecryptionKey.IsNull() && !googlePlayPlan.DecryptionKey.IsUnknown() {
			googlePlay.SetDecryptionKey(googlePlayPlan.DecryptionKey.ValueString())
		}

		if !googlePlayPlan.ServiceAccountCredentialsJson.IsNull() && !googlePlayPlan.ServiceAccountCredentialsJson.IsUnknown() {
			googlePlay.SetServiceAccountCredentials(googlePlayPlan.ServiceAccountCredentialsJson.ValueString())
		}

		if !googlePlayPlan.VerificationKey.IsNull() && !googlePlayPlan.VerificationKey.IsUnknown() {
			googlePlay.SetVerificationKey(googlePlayPlan.VerificationKey.ValueString())
		}

		if !googlePlayPlan.VerificationType.IsNull() && !googlePlayPlan.VerificationType.IsUnknown() {
			googlePlay.SetVerificationType(management.EnumApplicationNativeGooglePlayVerificationType(googlePlayPlan.VerificationType.ValueString()))
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

		var cacheDurationPlan applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV1
		diags.Append(p.CacheDuration.As(ctx, &cacheDurationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		cacheDuration := management.NewApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration()

		if !cacheDurationPlan.Amount.IsNull() && !cacheDurationPlan.Amount.IsUnknown() {
			cacheDuration.SetAmount(cacheDurationPlan.Amount.ValueInt32())
		}

		if !cacheDurationPlan.Units.IsNull() && !cacheDurationPlan.Units.IsUnknown() {
			cacheDuration.SetUnits(management.EnumDurationUnitMinsHours(cacheDurationPlan.Units.ValueString()))
		}

		data.SetCacheDuration(*cacheDuration)
	}

	return data, diags
}

func (p *applicationResourceModelV1) expandApplicationSAML(ctx context.Context) (*management.ApplicationSAML, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationSAML

	if !p.SAMLOptions.IsNull() && !p.SAMLOptions.IsUnknown() {
		var plan applicationSAMLOptionsResourceModelV1
		d := p.SAMLOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		var acsUrlsPlan []types.String

		diags.Append(plan.AcsUrls.ElementsAs(ctx, &acsUrlsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		acsUrls, d := framework.TFTypeStringSliceToStringSlice(acsUrlsPlan, path.Root("saml_options").AtName("acs_urls"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data = management.NewApplicationSAML(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_SAML,
			management.EnumApplicationType(plan.Type.ValueString()),
			acsUrls,
			plan.AssertionDuration.ValueInt32(),
			plan.SpEntityId.ValueString(),
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
		if !plan.CorsSettings.IsNull() && !plan.CorsSettings.IsUnknown() {
			var corsPlan applicationCorsSettingsResourceModelV1

			diags.Append(plan.CorsSettings.As(ctx, &corsPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			corsSettings, d := corsPlan.expand()
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetCorsSettings(*corsSettings)
		}

		if !plan.HomePageUrl.IsNull() && !plan.HomePageUrl.IsUnknown() {
			data.SetHomePageUrl(plan.HomePageUrl.ValueString())
		}

		if !plan.AssertionSignedEnabled.IsNull() && !plan.AssertionSignedEnabled.IsUnknown() {
			data.SetAssertionSigned(plan.AssertionSignedEnabled.ValueBool())
		}

		if !plan.IdpSigningKey.IsNull() && !plan.IdpSigningKey.IsUnknown() {

			var idpSigningOptionsPlan applicationOptionsIdpSigningKeyResourceModelV1

			diags.Append(plan.IdpSigningKey.As(ctx, &idpSigningOptionsPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			idpSigning := *management.NewApplicationSAMLAllOfIdpSigning(*management.NewApplicationSAMLAllOfIdpSigningKey(idpSigningOptionsPlan.KeyId.ValueString()))
			idpSigning.SetAlgorithm(management.EnumCertificateKeySignagureAlgorithm(idpSigningOptionsPlan.Algorithm.ValueString()))

			data.SetIdpSigning(idpSigning)
		}

		if !plan.EnableRequestedAuthnContext.IsNull() && !plan.EnableRequestedAuthnContext.IsUnknown() {
			data.SetEnableRequestedAuthnContext(plan.EnableRequestedAuthnContext.ValueBool())
		}

		if !plan.DefaultTargetUrl.IsNull() && !plan.DefaultTargetUrl.IsUnknown() {
			data.SetDefaultTargetUrl(plan.DefaultTargetUrl.ValueString())
		}

		if !plan.NameIdFormat.IsNull() && !plan.NameIdFormat.IsUnknown() {
			data.SetNameIdFormat(plan.NameIdFormat.ValueString())
		}

		if !plan.ResponseIsSigned.IsNull() && !plan.ResponseIsSigned.IsUnknown() {
			data.SetResponseSigned(plan.ResponseIsSigned.ValueBool())
		}

		if !plan.SessionNotOnOrAfterDuration.IsNull() && !plan.SessionNotOnOrAfterDuration.IsUnknown() {
			data.SetSessionNotOnOrAfterDuration(plan.SessionNotOnOrAfterDuration.ValueInt32())
		}

		if !plan.SloBinding.IsNull() && !plan.SloBinding.IsUnknown() {
			data.SetSloBinding(management.EnumApplicationSAMLSloBinding(plan.SloBinding.ValueString()))
		}

		if !plan.SloEndpoint.IsNull() && !plan.SloEndpoint.IsUnknown() {
			data.SetSloEndpoint(plan.SloEndpoint.ValueString())
		}

		if !plan.SloResponseEndpoint.IsNull() && !plan.SloResponseEndpoint.IsUnknown() {
			data.SetSloResponseEndpoint(plan.SloResponseEndpoint.ValueString())
		}

		if !plan.SloWindow.IsNull() && !plan.SloWindow.IsUnknown() {
			data.SetSloWindow(plan.SloWindow.ValueInt32())
		}

		if !plan.SpEncryption.IsNull() && !plan.SpEncryption.IsUnknown() {
			var spEncryptionPlan applicationSAMLOptionsSpEncryptionResourceModelV1

			diags.Append(plan.SpEncryption.As(ctx, &spEncryptionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			var spEncryptionCertificatePlan applicationSAMLOptionsSpEncryptionCertificateResourceModelV1

			diags.Append(spEncryptionPlan.Certificate.As(ctx, &spEncryptionCertificatePlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			spEncryption := management.NewApplicationSAMLAllOfSpEncryption(
				management.EnumCertificateKeyEncryptionAlgorithm(spEncryptionPlan.Algorithm.ValueString()),
				*management.NewApplicationSAMLAllOfSpEncryptionCertificate(spEncryptionCertificatePlan.Id.ValueString()),
			)

			data.SetSpEncryption(*spEncryption)
		}

		if !plan.SpVerification.IsNull() && !plan.SpVerification.IsUnknown() {
			var spVerificationPlan applicationSAMLOptionsSpVerificationResourceModelV1

			diags.Append(plan.SpVerification.As(ctx, &spVerificationPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			certificates := make([]management.ApplicationSAMLAllOfSpVerificationCertificates, 0)
			if !spVerificationPlan.CertificateIds.IsNull() && !spVerificationPlan.CertificateIds.IsUnknown() {
				var certificateIdsPlan []pingonetypes.ResourceIDValue

				diags.Append(spVerificationPlan.CertificateIds.ElementsAs(ctx, &certificateIdsPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				certificateIds, d := framework.TFTypePingOneResourceIDSliceToStringSlice(certificateIdsPlan, path.Root("saml_options").AtName("sp_verification").AtName("certificate_ids"))
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				for _, v := range certificateIds {
					certificate := *management.NewApplicationSAMLAllOfSpVerificationCertificates(v)
					certificates = append(certificates, certificate)
				}
			}

			spVerification := management.NewApplicationSAMLAllOfSpVerification(certificates)

			if !spVerificationPlan.AuthnRequestSigned.IsNull() && !spVerificationPlan.AuthnRequestSigned.IsUnknown() {
				spVerification.SetAuthnRequestSigned(spVerificationPlan.AuthnRequestSigned.ValueBool())
			}

			data.SetSpVerification(*spVerification)
		}

		if !plan.VirtualServerIdSettings.IsNull() && !plan.VirtualServerIdSettings.IsUnknown() {
			var vsSettingsPlan applicationSAMLOptionsVirtualServerIdSettingsResourceModelV1

			diags.Append(plan.VirtualServerIdSettings.As(ctx, &vsSettingsPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetVirtualServerIdSettings(*vsSettingsPlan.expand())
		}
	}

	return data, diags
}

func (p *applicationSAMLOptionsVirtualServerIdSettingsResourceModelV1) expand() *management.ApplicationSAMLAllOfVirtualServerIdSettings {
	vsSettings := management.NewApplicationSAMLAllOfVirtualServerIdSettings()

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		vsSettings.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.VirtualServerIds.IsNull() && !p.VirtualServerIds.IsUnknown() {
		virtualServerIdSettingsValue := []management.ApplicationSAMLAllOfVirtualServerIdSettingsVirtualServerIds{}
		for _, virtualServerIdsElement := range p.VirtualServerIds.Elements() {
			virtualServerIdsAttrs := virtualServerIdsElement.(types.Object).Attributes()
			serverIds := management.NewApplicationSAMLAllOfVirtualServerIdSettingsVirtualServerIds(virtualServerIdsAttrs["vs_id"].(types.String).ValueString())
			if !virtualServerIdsAttrs["default"].IsNull() && !virtualServerIdsAttrs["default"].IsUnknown() {
				serverIds.Default = virtualServerIdsAttrs["default"].(types.Bool).ValueBoolPointer()
			}
			virtualServerIdSettingsValue = append(virtualServerIdSettingsValue, *serverIds)
		}
		vsSettings.SetVirtualServerIds(virtualServerIdSettingsValue)
	}

	return vsSettings
}

func (p *applicationResourceModelV1) expandApplicationExternalLink(ctx context.Context) (*management.ApplicationExternalLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationExternalLink

	if !p.ExternalLinkOptions.IsNull() && !p.ExternalLinkOptions.IsUnknown() {
		var plan applicationExternalLinkOptionsResourceModelV1
		d := p.ExternalLinkOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data = management.NewApplicationExternalLink(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_EXTERNAL_LINK,
			management.ENUMAPPLICATIONTYPE_PORTAL_LINK_APP,
			plan.HomePageUrl.ValueString(),
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

func (p *applicationResourceModelV1) expandApplicationWSFed(ctx context.Context) (*management.ApplicationWSFED, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *management.ApplicationWSFED

	if !p.WSFedOptions.IsNull() && !p.WSFedOptions.IsUnknown() {
		var plan applicationWSFedOptionsResourceModelV1
		d := p.WSFedOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		var idpSigningOptionsPlan applicationOptionsIdpSigningKeyResourceModelV1
		diags.Append(plan.IdpSigningKey.As(ctx, &idpSigningOptionsPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		idpSigning := *management.NewApplicationWSFEDAllOfIdpSigning(management.EnumApplicationWSFEDIDPSigningAlgorithm(idpSigningOptionsPlan.Algorithm.ValueString()),
			*management.NewApplicationWSFEDAllOfIdpSigningKey(idpSigningOptionsPlan.KeyId.ValueString()))

		data = management.NewApplicationWSFED(
			p.Enabled.ValueBool(),
			p.Name.ValueString(),
			management.ENUMAPPLICATIONPROTOCOL_WS_FED,
			management.EnumApplicationType(plan.Type.ValueString()),
			plan.DomainName.ValueString(),
			idpSigning,
			plan.ReplyUrl.ValueString(),
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

		// WS-Fed specific options
		if !plan.AudienceRestriction.IsNull() && !plan.AudienceRestriction.IsUnknown() {
			data.SetAudienceRestriction(plan.AudienceRestriction.ValueString())
		}

		if !plan.CorsSettings.IsNull() && !plan.CorsSettings.IsUnknown() {
			var corsPlan applicationCorsSettingsResourceModelV1

			diags.Append(plan.CorsSettings.As(ctx, &corsPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			corsSettings, d := corsPlan.expand()
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetCorsSettings(*corsSettings)
		}

		if !plan.DomainName.IsNull() && !plan.DomainName.IsUnknown() {
			data.SetDomainName(plan.DomainName.ValueString())
		}

		if !plan.Kerberos.IsNull() && !plan.Kerberos.IsUnknown() {
			var kerberosPlan applicationWSFedKerberosResourceModelV1

			diags.Append(plan.Kerberos.As(ctx, &kerberosPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			kerberos, d := kerberosPlan.expand(ctx)
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}
			data.SetKerberos(*kerberos)
		}

		if !plan.ReplyUrl.IsNull() && !plan.ReplyUrl.IsUnknown() {
			data.SetReplyUrl(plan.ReplyUrl.ValueString())
		}

		if !plan.SloEndpoint.IsNull() && !plan.SloEndpoint.IsUnknown() {
			data.SetSloEndpoint(plan.SloEndpoint.ValueString())
		}

		if !plan.SubjectNameIdentifierFormat.IsNull() && !plan.SubjectNameIdentifierFormat.IsUnknown() {
			data.SetSubjectNameIdentifierFormat(management.EnumApplicationWSFEDSubjectNameIdentifierFormat(plan.SubjectNameIdentifierFormat.ValueString()))
		}

		if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
			data.SetType(management.EnumApplicationType(plan.Type.ValueString()))
		}
	}

	return data, diags
}

func (p *applicationWSFedKerberosResourceModelV1) expand(ctx context.Context) (*management.ApplicationWSFEDAllOfKerberos, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := management.NewApplicationWSFEDAllOfKerberos()

	var gateways []management.ApplicationWSFEDAllOfKerberosGateways
	for _, gateway := range p.Gateways.Elements() {
		var gatewayPlan applicationWSFedKerberosGatewayResourceModelV1
		d := gateway.(types.Object).As(ctx, &gatewayPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		gateway, d := gatewayPlan.expand(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		gateways = append(gateways, *gateway)
	}
	result.SetGateways(gateways)

	return result, diags
}

func (p *applicationWSFedKerberosGatewayResourceModelV1) expand(ctx context.Context) (*management.ApplicationWSFEDAllOfKerberosGateways, diag.Diagnostics) {
	var diags diag.Diagnostics

	var userTypePlan applicationWSFedGatewayUserTypeRersourceModelV1
	d := p.UserType.As(ctx, &userTypePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	userType := userTypePlan.expand()

	result := management.NewApplicationWSFEDAllOfKerberosGateways(
		p.Id.ValueString(),
		management.EnumApplicationWSFEDKerberosGatewayType(p.Type.ValueString()),
		*userType,
	)

	return result, diags
}

func (p *applicationWSFedGatewayUserTypeRersourceModelV1) expand() *management.ApplicationWSFEDAllOfKerberosUserType {
	result := management.NewApplicationWSFEDAllOfKerberosUserType()
	if !p.Id.IsNull() && !p.Id.IsUnknown() {
		result.SetId(p.Id.ValueString())
	}
	return result
}

func (p *applicationResourceModelV1) expandApplicationCommon(ctx context.Context) (*management.Application, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.Application{}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.LoginPageUrl.IsNull() && !p.LoginPageUrl.IsUnknown() {
		data.SetLoginPageUrl(p.LoginPageUrl.ValueString())
	}

	if !p.Icon.IsNull() && !p.Icon.IsUnknown() {
		var plan service.ImageResourceModel
		d := p.Icon.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		iconPlanItem := plan

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
		var plan applicationAccessControlGroupOptionsResourceModelV1
		d := p.AccessControlGroupOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)

		var groupsPlan []pingonetypes.ResourceIDValue

		diags.Append(plan.Groups.ElementsAs(ctx, &groupsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		groupsStr, d := framework.TFTypePingOneResourceIDSliceToStringSlice(groupsPlan, path.Root("access_control_group_options").AtName("groups"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		for _, group := range groupsStr {
			groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(group))
		}

		accessControl.SetGroup(*management.NewApplicationAccessControlGroup(
			management.EnumApplicationAccessControlGroupType(plan.Type.ValueString()),
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

func (p *applicationResourceModelV1) toState(ctx context.Context, apiObject *management.ReadOneApplication200Response) diag.Diagnostics {
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
		p.WSFedOptions = types.ObjectNull(applicationWsfedOptionsTFObjectTypes)

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
		p.WSFedOptions = types.ObjectNull(applicationWsfedOptionsTFObjectTypes)

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
		p.WSFedOptions = types.ObjectNull(applicationWsfedOptionsTFObjectTypes)

	case *management.ApplicationWSFED:
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
		p.ExternalLinkOptions = types.ObjectNull(applicationExternalLinkOptionsTFObjectTypes)

		p.WSFedOptions, d = applicationWsfedOptionsToTF(v)
		diags = append(diags, d...)
	}

	return diags
}
