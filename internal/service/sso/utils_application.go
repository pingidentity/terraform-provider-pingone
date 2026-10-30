// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/pingone-go-sdk-legacy/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso/helpers/beta"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func applicationExternalLinkOptionsToTF(apiObject *management.ApplicationExternalLink) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(applicationExternalLinkOptionsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"external_id":   framework.StringOkToTF(apiObject.GetExternalIdOk()),
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

func applicationOidcOptionsDataSourceToTF(ctx context.Context, apiObject *management.ApplicationOIDC, stateValue applicationOIDCOptionsDataSourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	commonAttrs, d := applicationOidcOptionsCommonToTF(ctx, apiObject, stateValue.applicationOIDCOptionsCommonModelV1)
	diags.Append(d...)

	if commonAttrs == nil {
		return types.ObjectNull(applicationOidcOptionsDataSourceTFObjectTypes), diags
	}

	attributesMap := utils.MergeAttributeValueMapsRtn(
		map[string]attr.Value{
			"client_id": framework.StringOkToTF(apiObject.GetIdOk()),
		},
		commonAttrs,
	)

	returnVar, d := types.ObjectValue(applicationOidcOptionsDataSourceTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsToTF(ctx context.Context, apiObject *management.ApplicationOIDC, stateValue applicationOIDCOptionsResourceModelV1) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	commonAttrs, d := applicationOidcOptionsCommonToTF(ctx, apiObject, stateValue.applicationOIDCOptionsCommonModelV1)
	diags.Append(d...)

	if commonAttrs == nil {
		return types.ObjectNull(applicationOidcOptionsResourceTFObjectTypes), diags
	}

	attributesMap := utils.MergeAttributeValueMapsRtn(
		beta.ApplicationBetaToTF(apiObject, stateValue.ApplicationOIDCOptionsResourceModelV1),
		commonAttrs,
	)

	returnVar, d := types.ObjectValue(applicationOidcOptionsResourceTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationOidcOptionsCommonToTF(ctx context.Context, apiObject *management.ApplicationOIDC, stateValue applicationOIDCOptionsCommonModelV1) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return nil, diags
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
			return nil, diags
		}
	}
	mobileApp, d := applicationMobileAppOkToTF(ctx, mobileAppObject, ok, mobileAppState)
	diags.Append(d...)

	signingObject, signingOk := apiObject.GetSigningOk()
	signing, d := applicationOidcOptionsSigningOkToTF(signingObject, signingOk)
	diags.Append(d...)

	return map[string]attr.Value{
		"additional_refresh_token_replay_protection_enabled": framework.BoolOkToTF(apiObject.GetAdditionalRefreshTokenReplayProtectionEnabledOk()),
		"allow_wildcard_in_redirect_uris":                    framework.BoolOkToTF(apiObject.GetAllowWildcardInRedirectUrisOk()),
		"certificate_based_authentication":                   kerberos,
		"cors_settings":                                      corsSettings,
		"device_path_id":                                     framework.StringOkToTF(apiObject.GetDevicePathIdOk()),
		"device_custom_verification_uri":                     framework.StringOkToTF(apiObject.GetDeviceCustomVerificationUriOk()),
		"device_timeout":                                     framework.Int32OkToTF(apiObject.GetDeviceTimeoutOk()),
		"device_polling_interval":                            framework.Int32OkToTF(apiObject.GetDevicePollingIntervalOk()),
		"grant_types":                                        framework.EnumSetOkToTF(apiObject.GetGrantTypesOk()),
		"home_page_url":                                      framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"idp_signoff":                                        framework.BoolOkToTF(apiObject.GetIdpSignoffOk()),
		"include_typ":                                        framework.BoolOkToTF(apiObject.GetIncludeTypOk()),
		"include_x5t":                                        framework.BoolOkToTF(apiObject.GetIncludeX5tOk()),
		"initiate_login_uri":                                 framework.StringOkToTF(apiObject.GetInitiateLoginUriOk()),
		"jwks_url":                                           framework.StringOkToTF(apiObject.GetJwksUrlOk()),
		"jwks":                                               framework.StringOkToTF(apiObject.GetJwksOk()),
		"mobile_app":                                         mobileApp,
		"op_session_check_enabled":                           framework.BoolOkToTF(apiObject.GetOpSessionCheckEnabledOk()),
		"par_requirement":                                    framework.EnumOkToTF(apiObject.GetParRequirementOk()),
		"par_timeout":                                        framework.Int32OkToTF(apiObject.GetParTimeoutOk()),
		"pkce_enforcement":                                   framework.EnumOkToTF(apiObject.GetPkceEnforcementOk()),
		"post_logout_redirect_uris":                          framework.StringSetOkToTF(apiObject.GetPostLogoutRedirectUrisOk()),
		"redirect_uris":                                      framework.StringSetOkToTF(apiObject.GetRedirectUrisOk()),
		"refresh_token_duration":                             framework.Int32OkToTF(apiObject.GetRefreshTokenDurationOk()),
		"refresh_token_rolling_duration":                     framework.Int32OkToTF(apiObject.GetRefreshTokenRollingDurationOk()),
		"refresh_token_rolling_grace_period_duration":        framework.Int32OkToTF(apiObject.GetRefreshTokenRollingGracePeriodDurationOk()),
		"refresh_token_type":                                 framework.EnumOkToTF(apiObject.GetRefreshTokenTypeOk()),
		"request_scopes_for_multiple_resources_enabled":      types.BoolValue(apiObject.GetRequestScopesForMultipleResourcesEnabled()),
		"require_signed_request_object":                      framework.BoolOkToTF(apiObject.GetRequireSignedRequestObjectOk()),
		"response_types":                                     framework.EnumSetOkToTF(apiObject.GetResponseTypesOk()),
		"skip_token_session_validation":                      framework.BoolOkToTF(apiObject.GetSkipTokenSessionValidationOk()),
		"signing":                                            signing,
		"support_unsigned_request_object":                    framework.BoolOkToTF(apiObject.GetSupportUnsignedRequestObjectOk()),
		"target_link_uri":                                    framework.StringOkToTF(apiObject.GetTargetLinkUriOk()),
		"token_endpoint_auth_method":                         framework.EnumOkToTF(apiObject.GetTokenEndpointAuthMethodOk()),
		"type":                                               framework.EnumOkToTF(apiObject.GetTypeOk()),
	}, diags
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

func applicationOidcOptionsSigningOkToTF(apiObject *management.ApplicationOIDCAllOfSigning, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcOptionsSigningTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"key_rotation_policy_id": pingonetypes.ResourceIDValue{},
	}

	if v, ok := apiObject.GetKeyRotationPolicyOk(); ok {
		attributesMap["key_rotation_policy_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(applicationOidcOptionsSigningTFObjectTypes, attributesMap)
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

	allowAuthenticationDeviceLockedObj, allowAuthenticationDeviceLockedOk := apiObject.GetAllowAuthenticationDeviceLockedOk()
	allowAuthenticationDeviceLocked, d := applicationMobileAppAllowAuthenticationDeviceLockedOkToTF(allowAuthenticationDeviceLockedObj, allowAuthenticationDeviceLockedOk)
	diags.Append(d...)

	allowUnPairAndChangeDeviceFromAppObj, allowUnPairAndChangeDeviceFromAppOk := apiObject.GetAllowUnPairAndChangeDeviceFromAppOk()
	allowUnPairAndChangeDeviceFromApp, d := applicationMobileAppAllowUnPairAndChangeDeviceFromAppOkToTF(allowUnPairAndChangeDeviceFromAppObj, allowUnPairAndChangeDeviceFromAppOk)
	diags.Append(d...)

	mobileAppSecurityPinRequireObj, mobileAppSecurityPinRequireOk := apiObject.GetMobileAppSecurityPinRequireOk()
	mobileAppSecurityPinRequire, d := applicationMobileAppSecurityPinRequireOkToTF(mobileAppSecurityPinRequireObj, mobileAppSecurityPinRequireOk)
	diags.Append(d...)

	mobileDeviceBiometricsSettingObj, mobileDeviceBiometricsSettingOk := apiObject.GetMobileDeviceBiometricsSettingOk()
	mobileDeviceBiometricsSetting, d := applicationMobileAppBiometricsSettingOkToTF(mobileDeviceBiometricsSettingObj, mobileDeviceBiometricsSettingOk)
	diags.Append(d...)

	notifyOptionalMobileVersionUpdatesObj, notifyOptionalMobileVersionUpdatesOk := apiObject.GetNotifyOptionalMobileVersionUpdatesOk()
	notifyOptionalMobileVersionUpdates, d := applicationMobileAppNotifyOptionalMobileVersionUpdatesOkToTF(notifyOptionalMobileVersionUpdatesObj, notifyOptionalMobileVersionUpdatesOk)
	diags.Append(d...)

	notifyRequireMobileVersionUpdatesObj, notifyRequireMobileVersionUpdatesOk := apiObject.GetNotifyRequireMobileVersionUpdatesOk()
	notifyRequireMobileVersionUpdates, d := applicationMobileAppNotifyRequireMobileVersionUpdatesOkToTF(notifyRequireMobileVersionUpdatesObj, notifyRequireMobileVersionUpdatesOk)
	diags.Append(d...)

	otpPushNotificationEnabledObj, otpPushNotificationEnabledOk := apiObject.GetOtpPushNotificationEnabledOk()
	otpPushNotificationEnabled, d := applicationMobileAppOtpPushNotificationEnabledOkToTF(otpPushNotificationEnabledObj, otpPushNotificationEnabledOk)
	diags.Append(d...)

	showAuthenticationInformationObj, showAuthenticationInformationOk := apiObject.GetShowAuthenticationInformationOk()
	showAuthenticationInformation, d := applicationMobileAppShowAuthenticationInformationOkToTF(showAuthenticationInformationObj, showAuthenticationInformationOk)
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"allow_authentication_device_locked":       allowAuthenticationDeviceLocked,
		"allow_un_pair_and_change_device_from_app": allowUnPairAndChangeDeviceFromApp,
		"apple_team_id":                          framework.StringOkToTF(apiObject.GetAppleTeamIdOk()),
		"bundle_id":                              framework.StringOkToTF(apiObject.GetBundleIdOk()),
		"huawei_app_id":                          framework.StringOkToTF(apiObject.GetHuaweiAppIdOk()),
		"huawei_package_name":                    framework.StringOkToTF(apiObject.GetHuaweiPackageNameOk()),
		"integrity_detection":                    integrityDetection,
		"mobile_app_security_pin_require":        mobileAppSecurityPinRequire,
		"mobile_device_biometrics_setting":       mobileDeviceBiometricsSetting,
		"notify_optional_mobile_version_updates": notifyOptionalMobileVersionUpdates,
		"notify_require_mobile_version_updates":  notifyRequireMobileVersionUpdates,
		"otp_push_notification_enabled":          otpPushNotificationEnabled,
		"package_name":                           framework.StringOkToTF(apiObject.GetPackageNameOk()),
		"passcode_grace_period":                  framework.Int32OkToTF(apiObject.GetPasscodeGracePeriodOk()),
		"passcode_refresh_seconds":               types.Int32Null(),
		"ping_id_settings":                       framework.BoolOkToTF(apiObject.GetPingIDSettingsOk()),
		"push":                                   applicationMobileAppPushOkToTF(apiObject.GetPushOk()),
		"show_authentication_information":        showAuthenticationInformation,
		"universal_app_link":                     framework.StringOkToTF(apiObject.GetUriPrefixOk()),
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

	appleAppAttestSettingsObj, appleAppAttestSettingsOk := apiObject.GetAppleAppAttestSettingsOk()
	appleAppAttestSettings, d := applicationMobileAppAppleAppAttestSettingsOkToTF(appleAppAttestSettingsObj, appleAppAttestSettingsOk)
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"apple_app_attest_settings": appleAppAttestSettings,
		"cache_duration":            cacheDuration,
		"enabled":                   types.BoolNull(),
		"excluded_platforms":        framework.EnumSetOkToTF(apiObject.GetExcludedPlatformsOk()),
		"google_play":               googlePlay,
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

func applicationMobileAppAppleAppAttestSettingsOkToTF(apiObject *management.ApplicationOIDCAllOfMobileIntegrityDetectionAppleAppAttestSettings, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppAppleAppAttestSettingsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppAppleAppAttestSettingsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppPushOkToTF(apiObject *management.ApplicationOIDCAllOfMobilePush, ok bool) types.Object {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppPushTFObjectTypes)
	}

	numberMatching := types.ObjectNull(applicationOidcMobileAppPushNumberMatchingTFObjectTypes)
	if numberMatchingObj, okNumberMatching := apiObject.GetNumberMatchingOk(); okNumberMatching && numberMatchingObj != nil {
		numberMatchingValue, d := types.ObjectValue(applicationOidcMobileAppPushNumberMatchingTFObjectTypes, map[string]attr.Value{
			"type": framework.EnumOkToTF(numberMatchingObj.GetTypeOk()),
		})
		diags.Append(d...)
		numberMatching = numberMatchingValue
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppPushTFObjectTypes, map[string]attr.Value{
		"number_matching": numberMatching,
	})
	diags.Append(d...)

	return returnVar
}

func applicationMobileAppAllowAuthenticationDeviceLockedOkToTF(apiObject *management.ApplicationOIDCAllOfMobileAllowAuthenticationDeviceLocked, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppEnabledToggleTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppEnabledToggleTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppAllowUnPairAndChangeDeviceFromAppOkToTF(apiObject *management.ApplicationOIDCAllOfMobileAllowUnPairAndChangeDeviceFromApp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppEnabledToggleTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppEnabledToggleTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppOtpPushNotificationEnabledOkToTF(apiObject *management.ApplicationOIDCAllOfMobileOtpPushNotificationEnabled, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppEnabledToggleTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppEnabledToggleTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppShowAuthenticationInformationOkToTF(apiObject *management.ApplicationOIDCAllOfMobileShowAuthenticationInformation, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppEnabledToggleTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppEnabledToggleTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppNotifyOptionalMobileVersionUpdatesOkToTF(apiObject *management.ApplicationOIDCAllOfMobileNotifyOptionalMobileVersionUpdates, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppVersionUpdatesTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"android": framework.StringOkToTF(apiObject.GetAndroidOk()),
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"ios":     framework.StringOkToTF(apiObject.GetIosOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppVersionUpdatesTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppNotifyRequireMobileVersionUpdatesOkToTF(apiObject *management.ApplicationOIDCAllOfMobileNotifyRequireMobileVersionUpdates, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppVersionUpdatesTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"android": framework.StringOkToTF(apiObject.GetAndroidOk()),
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"ios":     framework.StringOkToTF(apiObject.GetIosOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppVersionUpdatesTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppSecurityPinRequireOkToTF(apiObject *management.ApplicationOIDCAllOfMobileMobileAppSecurityPinRequire, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppSecurityPinRequireTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"enabled":                 framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"require_for_all_devices": framework.BoolOkToTF(apiObject.GetRequireForAllDevicesOk()),
		"security_pin_digit":      framework.EnumOkToTF(apiObject.GetSecurityPinDigitOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppSecurityPinRequireTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationMobileAppBiometricsSettingOkToTF(apiObject *management.ApplicationOIDCAllOfMobileMobileDeviceBiometricsSetting, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationOidcMobileAppBiometricsSettingTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"device_passcode_fallback_enabled": framework.BoolOkToTF(apiObject.GetDevicePasscodeFallbackEnabledOk()),
		"face_id_consent_enabled":          framework.BoolOkToTF(apiObject.GetFaceIdConsentEnabledOk()),
		"notification_actions_enabled":     framework.BoolOkToTF(apiObject.GetNotificationActionsEnabledOk()),
		"required":                         framework.BoolOkToTF(apiObject.GetRequiredOk()),
	}

	returnVar, d := types.ObjectValue(applicationOidcMobileAppBiometricsSettingTFObjectTypes, attributesMap)
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

	typeValue := framework.EnumOkToTF(apiObject.GetTypeOk())
	if management.EnumApplicationType(typeValue.ValueString()) == management.ENUMAPPLICATIONTYPE_TEMPLATE_APP {
		diags.AddAttributeError(
			path.Root("saml_options").AtName("type"),
			"Application Type Not Supported",
			"The `saml_options.type` value `TEMPLATE_APP` (used for Catalog Applications) is not supported.\nTo proceed, select Enable Advanced Configuration in the UI or change `saml_options.type` to `WEB_APP` before using this resource in Terraform.")
	}

	virtualServerIdSettings, d := applicationSamlVirtualServerIdSettingsOkToTF(apiObject.GetVirtualServerIdSettingsOk())
	diags.Append(d...)

	attributesMap := map[string]attr.Value{
		"acs_urls":                 framework.StringSetOkToTF(apiObject.GetAcsUrlsOk()),
		"assertion_duration":       framework.Int32OkToTF(apiObject.GetAssertionDurationOk()),
		"assertion_signed_enabled": framework.BoolOkToTF(apiObject.GetAssertionSignedOk()),
		"cors_settings":            corsSettings,
		"enable_always_accept_acs_url_in_signed_authn_request": framework.BoolOkToTF(apiObject.GetEnableAlwaysAcceptAcsUrlInSignedAuthnRequestOk()),
		"enable_requested_authn_context":                       framework.BoolOkToTF(apiObject.GetEnableRequestedAuthnContextOk()),
		"home_page_url":                                        framework.StringOkToTF(apiObject.GetHomePageUrlOk()),
		"idp_signing_key":                                      idpSigningKey,
		"default_target_url":                                   framework.StringOkToTF(apiObject.GetDefaultTargetUrlOk()),
		"nameid_format":                                        framework.StringOkToTF(apiObject.GetNameIdFormatOk()),
		"response_is_signed":                                   framework.BoolOkToTF(apiObject.GetResponseSignedOk()),
		"session_not_on_or_after_duration":                     framework.Int32OkToTF(apiObject.GetSessionNotOnOrAfterDurationOk()),
		"slo_binding":                                          framework.EnumOkToTF(apiObject.GetSloBindingOk()),
		"slo_endpoint":                                         framework.StringOkToTF(apiObject.GetSloEndpointOk()),
		"slo_response_endpoint":                                framework.StringOkToTF(apiObject.GetSloResponseEndpointOk()),
		"slo_window":                                           framework.Int32OkToTF(apiObject.GetSloWindowOk()),
		"sp_encryption":                                        spEncryption,
		"sp_entity_id":                                         framework.StringOkToTF(apiObject.GetSpEntityIdOk()),
		"sp_verification":                                      spVerification,
		"type":                                                 typeValue,
		"virtual_server_id_settings":                           virtualServerIdSettings,
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationSamlIdpSigningKeyOkToTF(apiObject *management.ApplicationSAMLAllOfIdpSigning, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationIdpSigningKeyTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
		"key_id":    pingonetypes.NewResourceIDNull(),
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		attributesMap["key_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(applicationIdpSigningKeyTFObjectTypes, attributesMap)
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

func applicationSamlVirtualServerIdSettingsOkToTF(apiObject *management.ApplicationSAMLAllOfVirtualServerIdSettings, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes), diags
	}

	vsSettingsMap := map[string]attr.Value{
		"enabled":            framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"virtual_server_ids": types.ListNull(types.ObjectType{AttrTypes: applicationSamlOptionsVirtualServerIdSettingsVirtualServerIdsTFObjectTypes}),
	}

	if vsIds, ok := apiObject.GetVirtualServerIdsOk(); ok {
		vsIdsTF := make([]attr.Value, len(vsIds))
		for i, vsId := range vsIds {
			vsIdsTF[i], diags = types.ObjectValue(applicationSamlOptionsVirtualServerIdSettingsVirtualServerIdsTFObjectTypes, map[string]attr.Value{
				"vs_id":   framework.StringOkToTF(vsId.GetVsIdOk()),
				"default": framework.BoolOkToTF(vsId.GetDefaultOk()),
			})
			if diags.HasError() {
				return types.ObjectNull(applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes), diags
			}
		}
		vsSettingsMap["virtual_server_ids"], diags = types.ListValue(types.ObjectType{AttrTypes: applicationSamlOptionsVirtualServerIdSettingsVirtualServerIdsTFObjectTypes}, vsIdsTF)
		if diags.HasError() {
			return types.ObjectNull(applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes), diags
		}
	}

	returnVar, d := types.ObjectValue(applicationSamlOptionsVirtualServerIdSettingsTFObjectTypes, vsSettingsMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationWsfedOptionsToTF(apiObject *management.ApplicationWSFED) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(applicationWsfedOptionsTFObjectTypes), diags
	}

	corsSettings, d := applicationCorsSettingsOkToTF(apiObject.GetCorsSettingsOk())
	diags.Append(d...)

	idpSigningKey, d := applicationWsfedIdpSigningKeyOkToTF(apiObject.GetIdpSigningOk())
	diags.Append(d...)

	kerberos, d := applicationWsfedKerberosOkToTF(apiObject.GetKerberosOk())
	diags.Append(d...)

	validIssuanceCertificates := types.SetNull(types.ObjectType{AttrTypes: applicationWsfedOptionsValidIssuanceCertificateTFObjectTypes})
	if certificates, ok := apiObject.GetValidIssuanceCertificatesOk(); ok && certificates != nil {
		certificateValues := make([]attr.Value, 0, len(certificates))
		for _, certificate := range certificates {
			certificateValue, d := types.ObjectValue(applicationWsfedOptionsValidIssuanceCertificateTFObjectTypes, map[string]attr.Value{
				"certificate_id": framework.PingOneResourceIDOkToTF(certificate.GetIdOk()),
			})
			diags.Append(d...)
			certificateValues = append(certificateValues, certificateValue)
		}

		certificatesSet, d := types.SetValue(types.ObjectType{AttrTypes: applicationWsfedOptionsValidIssuanceCertificateTFObjectTypes}, certificateValues)
		diags.Append(d...)
		validIssuanceCertificates = certificatesSet
	}

	attributesMap := map[string]attr.Value{
		"assertion_duration":             framework.Int32OkToTF(apiObject.GetAssertionDurationOk()),
		"audience_restriction":           framework.StringOkToTF(apiObject.GetAudienceRestrictionOk()),
		"cors_settings":                  corsSettings,
		"domain_name":                    framework.StringOkToTF(apiObject.GetDomainNameOk()),
		"hybrid_join_device_attribute":   framework.StringOkToTF(apiObject.GetHybridJoinDeviceAttributeOk()),
		"idp_signing_key":                idpSigningKey,
		"kerberos":                       kerberos,
		"reply_url":                      framework.StringOkToTF(apiObject.GetReplyUrlOk()),
		"slo_endpoint":                   framework.StringOkToTF(apiObject.GetSloEndpointOk()),
		"sts_token_type":                 framework.EnumOkToTF(apiObject.GetStsTokenTypeOk()),
		"subject_name_identifier_format": framework.EnumOkToTF(apiObject.GetSubjectNameIdentifierFormatOk()),
		"type":                           framework.EnumOkToTF(apiObject.GetTypeOk()),
		"use_mex13":                      framework.BoolOkToTF(apiObject.GetUseMex13Ok()),
		"user_attribute_username_match":  framework.StringOkToTF(apiObject.GetUserAttributeUsernameMatchOk()),
		"user_attribute_x509_match":      framework.StringOkToTF(apiObject.GetUserAttributeX509MatchOk()),
		"valid_issuance_certificates":    validIssuanceCertificates,
		"ws_trust_version":               framework.EnumOkToTF(apiObject.GetWsTrustVersionOk()),
	}

	returnVar, d := types.ObjectValue(applicationWsfedOptionsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationWsfedIdpSigningKeyOkToTF(apiObject *management.ApplicationWSFEDAllOfIdpSigning, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationIdpSigningKeyTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
		"key_id":    pingonetypes.NewResourceIDNull(),
	}

	if v, ok := apiObject.GetKeyOk(); ok {
		attributesMap["key_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(applicationIdpSigningKeyTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func applicationWsfedKerberosOkToTF(apiObject *management.ApplicationWSFEDAllOfKerberos, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationWsfedOptionsKerberosTFObjectTypes), diags
	}

	var gateways []attr.Value
	for _, gateway := range apiObject.Gateways {
		gatewayObj, d := applicationWsfedKerberosGatewaysToTF(gateway)
		diags.Append(d...)
		gateways = append(gateways, gatewayObj)
	}

	gatewaysValue, d := types.SetValue(types.ObjectType{
		AttrTypes: applicationWsfedOptionsKerberosGatewayTFObjectTypes,
	}, gateways)
	diags.Append(d...)
	returnVar, d := types.ObjectValue(applicationWsfedOptionsKerberosTFObjectTypes, map[string]attr.Value{
		"gateways": gatewaysValue,
	})
	diags.Append(d...)

	return returnVar, diags
}

func applicationWsfedKerberosGatewaysToTF(apiObject management.ApplicationWSFEDAllOfKerberosGateways) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	userType, d := types.ObjectValue(applicationWsfedOptionsKerberosGatewayUserTypeTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.UserType.GetIdOk()),
	})
	diags.Append(d...)

	returnVar, d := types.ObjectValue(applicationWsfedOptionsKerberosGatewayTFObjectTypes, map[string]attr.Value{
		"id":        framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"type":      framework.EnumOkToTF(apiObject.GetTypeOk()),
		"user_type": userType,
	})
	diags.Append(d...)

	return returnVar, diags
}
