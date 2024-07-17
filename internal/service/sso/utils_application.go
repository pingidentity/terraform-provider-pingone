package sso

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

func applicationExternalLinkOptionsToTF(apiObject *management.ApplicationExternalLink) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(applicationExternalLinkOptionsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"home_page_url": framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
	}

	returnVar, d := types.ObjectValue(applicationExternalLinkOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationAccessControlGroupOptionsToTF(apiObject *management.ApplicationAccessControlGroup, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok && apiObject == nil {
		return types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	if v, ok := apiObject.GetGroupsOk(); ok {
		groups := make([]string, 0)

		for _, group := range v {
			groups = append(groups, group.GetId())
		}

		attributesMap["groups"] = framework.PingOneResourceIDSetToTF(groups)
	}

	returnVar, d := types.ObjectValue(applicationAccessControlGroupOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationCorsSettingsOkToTF(apiObject *management.ApplicationCorsSettings, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationCorsSettingsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"behavior": framework.EnumOkToTF(apiObject.GetBehaviorOk()),
		"origins":  framework.StringSetOkToTF(apiObject.GetOriginsOk()),
	}

	returnVar, d := types.ObjectValue(applicationCorsSettingsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsToTF(ctx context.Context, apiObject *management.ApplicationOIDC, stateValue applicationOIDCOptionsResourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(applicationOidcOptionsTFObjectTypes), diags
	}

	kerberos, d := applicationOidcOptionsCertificateBasedAuthenticationToTF(apiObject.GetKerberosOk())
	diags.Append(d...)

	corsSettings, d := applicationCorsSettingsOkToTF(apiObject.GetCorsSettingsOk())
	diags.Append(d...)

	mobileAppObject, ok := apiObject.GetMobileOk()
	var mobileAppState applicationOIDCMobileAppResourceModelV1
	if !stateValue.MobileApp.IsNull() && !stateValue.MobileApp.IsUnknown() {
		d := stateValue.MobileApp.As(ctx, &mobileAppState, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(applicationOidcOptionsTFObjectTypes), diags
		}
	}
	mobileApp, d := applicationMobileAppOkToTF(ctx, mobileAppObject, ok, mobileAppState)
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"additional_refresh_token_replay_protection_enabled": framework.BoolOkToTF(apiObject.GetAdditionalRefreshTokenReplayProtectionEnabledOk()),
		"allow_wildcard_in_redirect_uris":                    framework.BoolOkToTF(apiObject.GetAllowWildcardInRedirectUrisOk()),
		"certificate_based_authentication":                   kerberos,
		"client_id":                                          framework.StringOkToTF(apiObject.GetIdOk()),
		"cors_settings":                                      corsSettings,
		"device_path_id":                                     framework.StringOkToTF(apiObject.GetDevicePathIdOk()),
		"device_custom_verification_uri":                     framework.StringOkToTF(apiObject.GetDeviceCustomVerificationUriOk()),
		"device_timeout":                                     framework.Int32OkToTF(apiObject.GetDeviceTimeoutOk()),
		"device_polling_interval":                            framework.Int32OkToTF(apiObject.GetDevicePollingIntervalOk()),
		"grant_types":                                        framework.EnumSetOkToTF(apiObject.GetGrantTypesOk()),
		"home_page_url":                                      framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"initiate_login_uri":                                 framework.StringOkToTF(apiObject.GetInitiateLoginUriOk()),
		"jwks_url":                                           framework.StringOkToTF(apiObject.GetJwksUrlOk()),
		"jwks":                                               framework.StringOkToTF(apiObject.GetJwksOk()),
		"mobile_app":                                         mobileApp,
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
		"target_link_uri":                                    framework.StringOkToTF(apiObject.GetTargetLinkUriOk()),
		"token_endpoint_auth_method":                         framework.EnumOkToTF(apiObject.GetTokenEndpointAuthMethodOk()),
		"type":                                               framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsCertificateBasedAuthenticationToTF(apiObject *management.ApplicationOIDCAllOfKerberos, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok && apiObject == nil {
		return types.ObjectNull(applicationOidcOptionsCertificateAuthenticationTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"key_id": pingonetypes.ResourceIDValue{},
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		attributesMap["key_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(applicationOidcOptionsCertificateAuthenticationTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppOkToTF(ctx context.Context, apiObject *management.ApplicationOIDCAllOfMobile, ok bool, stateValue applicationOIDCMobileAppResourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppTFObjectTypes), diags
	}

	integrityDetectionObj, ok := apiObject.GetIntegrityDetectionOk()
	var integrityDetectionState applicationOIDCMobileAppIntegrityDetectionResourceModelV1
	if !stateValue.IntegrityDetection.IsNull() && !stateValue.IntegrityDetection.IsUnknown() {
		d := stateValue.IntegrityDetection.As(ctx, &integrityDetectionState, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(applicationOidcMobileAppTFObjectTypes), diags
		}
	}
	integrityDetection, d := applicationMobileAppIntegrityDetectionOkToTF(ctx, integrityDetectionObj, ok, integrityDetectionState)
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

			return types.ObjectNull(applicationOidcMobileAppTFObjectTypes), diags
		}

		attributesMap["passcode_refresh_seconds"] = framework.Int32OkToTF(v.GetDurationOk())
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionOkToTF(ctx context.Context, apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetection, ok bool, stateValue applicationOIDCMobileAppIntegrityDetectionResourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppIntegrityDetectionTFObjectTypes), diags
	}

	cacheDuration, d := applicationMobileAppIntegrityDetectionCacheDurationOkToTF(apiObject.GetCacheDurationOk())
	diags.Append(d...)

	googlePlayObject, ok := apiObject.GetGooglePlayOk()
	var googlePlayState applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV1
	if !stateValue.GooglePlay.IsNull() && !stateValue.GooglePlay.IsUnknown() {
		d := stateValue.GooglePlay.As(ctx, &googlePlayState, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(applicationOidcMobileAppIntegrityDetectionTFObjectTypes), diags
		}
	}
	googlePlay, d := applicationMobileAppIntegrityDetectionGooglePlayOkToTF(googlePlayObject, ok, googlePlayState)
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

	returnVar, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionCacheDurationOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetectionCacheDuration, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"amount": framework.Int32OkToTF(apiObject.GetAmountOk()),
		"units":  framework.EnumOkToTF(apiObject.GetUnitsOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppIntegrityDetectionGooglePlayOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetectionGooglePlay, ok bool, stateValue applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"decryption_key":                   stateValue.DecryptionKey,
		"service_account_credentials_json": stateValue.ServiceAccountCredentialsJson,
		"verification_key":                 stateValue.VerificationKey,
		"verification_type":                framework.EnumOkToTF(apiObject.GetVerificationTypeOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlOptionsToTF(apiObject *management.ApplicationSAML) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(applicationSamlOptionsTFObjectTypes), diags
	}

	corsSettings, d := applicationCorsSettingsOkToTF(apiObject.GetCorsSettingsOk())
	diags.Append(d...)

	idpSigningKey, d := applicationSamlIdpSigningKeyOkToTF(apiObject.GetIdpSigningOk())
	diags.Append(d...)

	spEncryption, d := applicationSamlSpEncryptionOkToTF(apiObject.GetSpEncryptionOk())
	diags.Append(d...)

	spVerification, d := applicationSamlSpVerificationOkToTF(apiObject.GetSpVerificationOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"acs_urls":                       framework.StringSetOkToTF(apiObject.GetAcsUrlsOk()),
		"assertion_duration":             framework.Int32OkToTF(apiObject.GetAssertionDurationOk()),
		"assertion_signed_enabled":       framework.BoolOkToTF(apiObject.GetAssertionSignedOk()),
		"cors_settings":                  corsSettings,
		"enable_requested_authn_context": framework.BoolOkToTF(apiObject.GetEnableRequestedAuthnContextOk()),
		"home_page_url":                  framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"idp_signing_key":                idpSigningKey,
		"default_target_url":             framework.StringOkToTF(apiObject.GetDefaultTargetUrlOk()),
		"nameid_format":                  framework.StringOkToTF(apiObject.GetNameIdFormatOk()),
		"response_is_signed":             framework.BoolOkToTF(apiObject.GetResponseSignedOk()),
		"slo_binding":                    framework.EnumOkToTF(apiObject.GetSloBindingOk()),
		"slo_endpoint":                   framework.StringOkToTF(apiObject.GetSloEndpointOk()),
		"slo_response_endpoint":          framework.StringOkToTF(apiObject.GetSloResponseEndpointOk()),
		"slo_window":                     framework.Int32OkToTF(apiObject.GetSloWindowOk()),
		"sp_encryption":                  spEncryption,
		"sp_entity_id":                   framework.StringOkToTF(apiObject.GetSpEntityIdOk()),
		"sp_verification":                spVerification,
		"type":                           framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlIdpSigningKeyOkToTF(apiObject *management.ApplicationSAMLAllOfIdpSigning, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationSamlOptionsIdpSigningKeyTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
		"key_id":    pingonetypes.NewResourceIDNull(),
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		attributesMap["key_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsIdpSigningKeyTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlSpEncryptionOkToTF(apiObject *management.ApplicationSAMLAllOfSpEncryption, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationSamlOptionsSpEncryptionTFObjectTypes), diags
	}

	certificate, d := applicationSamlSpEncryptionCertificateOkToTF(apiObject.GetCertificateOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"algorithm":   framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
		"certificate": certificate,
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsSpEncryptionTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlSpEncryptionCertificateOkToTF(apiObject *management.ApplicationSAMLAllOfSpEncryptionCertificate, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationSamlOptionsSpEncryptionCertificateTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsSpEncryptionCertificateTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlSpVerificationOkToTF(apiObject *management.ApplicationSAMLAllOfSpVerification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationSamlOptionsSpVerificationTFObjectTypes), diags
	}

	certificateIds := make([]string, 0)
	if v, ok := apiObject.GetCertificatesOk(); ok {
		for _, certificate := range v {
			certificateIds = append(certificateIds, certificate.GetId())
		}
	}

	attributesMap := map[string]attr.Value{
		"authn_request_signed": framework.BoolOkToTF(apiObject.GetAuthnRequestSignedOk()),
		"certificate_ids":      framework.PingOneResourceIDSetToTF(certificateIds),
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsSpVerificationTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}
