// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

const (
	POLICY_TYPE_PINGONE_MFA = "PING_ONE_MFA"
	POLICY_TYPE_PINGID      = "PING_ONE_ID"
)

// Types
type MFADevicePolicyDefaultResource serviceClientType

type MFADevicePolicyDefaultResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	PolicyType            types.String                 `tfsdk:"policy_type"`
	Name                  types.String                 `tfsdk:"name"`
	Authentication        types.Object                 `tfsdk:"authentication"`
	NewDeviceNotification types.String                 `tfsdk:"new_device_notification"`
	IgnoreUserLock        types.Bool                   `tfsdk:"ignore_user_lock"`
	NotificationsPolicy   types.Object                 `tfsdk:"notifications_policy"`
	RememberMe            types.Object                 `tfsdk:"remember_me"`
	Sms                   types.Object                 `tfsdk:"sms"`
	Voice                 types.Object                 `tfsdk:"voice"`
	Email                 types.Object                 `tfsdk:"email"`
	Mobile                types.Object                 `tfsdk:"mobile"`
	Totp                  types.Object                 `tfsdk:"totp"`
	Fido2                 types.Object                 `tfsdk:"fido2"`
	Desktop               types.Object                 `tfsdk:"desktop"`
	Yubikey               types.Object                 `tfsdk:"yubikey"`
	OathToken             types.Object                 `tfsdk:"oath_token"`
	UpdatedAt             timetypes.RFC3339            `tfsdk:"updated_at"`
}

type MFADevicePolicyPingIDDeviceResourceModel struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	Otp                        types.Object `tfsdk:"otp"`
	PairingDisabled            types.Bool   `tfsdk:"pairing_disabled"`
	PairingKeyLifetime         types.Object `tfsdk:"pairing_key_lifetime"`
	PromptForNicknameOnPairing types.Bool   `tfsdk:"prompt_for_nickname_on_pairing"`
}

type MFADevicePolicyDesktopResourceModel MFADevicePolicyPingIDDeviceResourceModel
type MFADevicePolicyYubikeyResourceModel MFADevicePolicyPingIDDeviceResourceModel
type MFADevicePolicyOathTokenResourceModel MFADevicePolicyPingIDDeviceResourceModel

type MFADevicePolicyDefaultMobileResourceModel struct {
	Applications               types.Map    `tfsdk:"applications"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	Otp                        types.Object `tfsdk:"otp"`
	PromptForNicknameOnPairing types.Bool   `tfsdk:"prompt_for_nickname_on_pairing"`
}

type MFADevicePolicyDefaultMobileApplicationResourceModel struct {
	AutoEnrolment                   types.Object `tfsdk:"auto_enrollment"`
	BiometricsEnabled               types.Bool   `tfsdk:"biometrics_enabled"`
	DeviceAuthorization             types.Object `tfsdk:"device_authorization"`
	IntegrityDetection              types.String `tfsdk:"integrity_detection"`
	IpPairingConfiguration          types.Object `tfsdk:"ip_pairing_configuration"`
	Otp                             types.Object `tfsdk:"otp"`
	PairingDisabled                 types.Bool   `tfsdk:"pairing_disabled"`
	PairingKeyLifetime              types.Object `tfsdk:"pairing_key_lifetime"`
	Push                            types.Object `tfsdk:"push"`
	PushLimit                       types.Object `tfsdk:"push_limit"`
	PushTimeout                     types.Object `tfsdk:"push_timeout"`
	NewRequestDurationConfiguration types.Object `tfsdk:"new_request_duration_configuration"`
	Type                            types.String `tfsdk:"type"`
}

type MFADevicePolicyDefaultMobileApplicationOtpResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyDefaultMobileApplicationPushResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	NumberMatching types.Object `tfsdk:"number_matching"`
}

type MFADevicePolicyDefaultMobileApplicationPushNumberMatchingResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyMobileApplicationPushLimitResourceModel struct {
	Count        types.Int32  `tfsdk:"count"`
	LockDuration types.Object `tfsdk:"lock_duration"`
	TimePeriod   types.Object `tfsdk:"time_period"`
}

type MFADevicePolicyMobileApplicationNewRequestDurationConfigurationResourceModel struct {
	DeviceTimeout types.Object `tfsdk:"device_timeout"`
	TotalTimeout  types.Object `tfsdk:"total_timeout"`
}

type MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutResourceModel struct {
	Duration types.Int32  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

// TF Object type definitions for PingID devices
var (
	MFADevicePolicyPingIDDeviceOtpTFObjectTypes = map[string]attr.Type{
		"failure": types.ObjectType{AttrTypes: MFADevicePolicyFailureTFObjectTypes},
	}

	MFADevicePolicyPingIDDeviceTFObjectTypes = map[string]attr.Type{
		"enabled":                        types.BoolType,
		"otp":                            types.ObjectType{AttrTypes: MFADevicePolicyPingIDDeviceOtpTFObjectTypes},
		"pairing_disabled":               types.BoolType,
		"pairing_key_lifetime":           types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"prompt_for_nickname_on_pairing": types.BoolType,
	}

	// TF Object type definitions for mobile applications in default resource
	MFADevicePolicyDefaultMobileApplicationPushNumberMatchingTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	MFADevicePolicyDefaultMobileApplicationPushTFObjectTypes = map[string]attr.Type{
		"enabled":         types.BoolType,
		"number_matching": types.ObjectType{AttrTypes: MFADevicePolicyDefaultMobileApplicationPushNumberMatchingTFObjectTypes},
	}

	MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int32Type,
		"time_unit": types.StringType,
	}

	MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes = map[string]attr.Type{
		"device_timeout": types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes},
		"total_timeout":  types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes},
	}

	MFADevicePolicyDefaultMobileApplicationTFObjectTypes = map[string]attr.Type{
		"auto_enrollment":                    types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes},
		"biometrics_enabled":                 types.BoolType,
		"device_authorization":               types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes},
		"integrity_detection":                types.StringType,
		"otp":                                types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationOtpTFObjectTypes},
		"pairing_disabled":                   types.BoolType,
		"pairing_key_lifetime":               types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"push":                               types.ObjectType{AttrTypes: MFADevicePolicyDefaultMobileApplicationPushTFObjectTypes},
		"push_limit":                         types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationPushLimitTFObjectTypes},
		"push_timeout":                       types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"new_request_duration_configuration": types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes},
		"type":                               types.StringType,
		"ip_pairing_configuration":           types.ObjectType{AttrTypes: MFADevicePolicyMobileIpPairingConfigurationTFObjectTypes},
	}

	MFADevicePolicyMobileIpPairingConfigurationTFObjectTypes = map[string]attr.Type{
		"any_ip_address":          types.BoolType,
		"only_these_ip_addresses": types.SetType{ElemType: types.StringType},
	}

	MFADevicePolicyRememberMeWebTFObjectTypes = map[string]attr.Type{
		"enabled":   types.BoolType,
		"life_time": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyRememberMeTFObjectTypes = map[string]attr.Type{
		"web": types.ObjectType{AttrTypes: MFADevicePolicyRememberMeWebTFObjectTypes},
	}

	MFADevicePolicyNotificationsPolicyTFObjectTypes = map[string]attr.Type{
		"id": types.StringType,
	}

	MFADevicePolicyDefaultMobileTFObjectTypes = map[string]attr.Type{
		"applications":                   types.MapType{ElemType: types.ObjectType{AttrTypes: MFADevicePolicyDefaultMobileApplicationTFObjectTypes}},
		"enabled":                        types.BoolType,
		"otp":                            types.ObjectType{AttrTypes: MFADevicePolicyMobileOtpTFObjectTypes},
		"prompt_for_nickname_on_pairing": types.BoolType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &MFADevicePolicyDefaultResource{}
	_ resource.ResourceWithConfigure   = &MFADevicePolicyDefaultResource{}
	_ resource.ResourceWithImportState = &MFADevicePolicyDefaultResource{}
	_ resource.ResourceWithModifyPlan  = &MFADevicePolicyDefaultResource{}
)

// New Object
func NewMFADevicePolicyDefaultResource() resource.Resource {
	return &MFADevicePolicyDefaultResource{}
}

// Metadata
func (r *MFADevicePolicyDefaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_device_policy_default"
}

// Schema
func (r *MFADevicePolicyDefaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const mobileApplicationsPushLimitCountDefault = 5
	const mobileApplicationsPushLimitCountMin = 1
	const mobileApplicationsPushLimitCountMax = 50

	const mobileApplicationsPushLimitLockDurationDurationDefault = 30
	const mobileApplicationsPushLimitLockDurationDurationMinMinutes = 1
	const mobileApplicationsPushLimitLockDurationDurationMaxMinutes = 120
	const mobileApplicationsPushLimitLockDurationDurationMinSeconds = 60
	const mobileApplicationsPushLimitLockDurationDurationMaxSeconds = 7200

	const mobileApplicationsPushLimitTimePeriodDurationDefault = 10
	const mobileApplicationsPushLimitTimePeriodDurationMinMinutes = 1
	const mobileApplicationsPushLimitTimePeriodDurationMaxMinutes = 120
	const mobileApplicationsPushLimitTimePeriodDurationMinSeconds = 60
	const mobileApplicationsPushLimitTimePeriodDurationMaxSeconds = 7200

	const mobileApplicationsPushTimeoutDurationMin = 1
	const mobileApplicationsPushTimeoutDurationMax = 120

	const mobileApplicationsOtpFailureCoolDownDurationDefault = 2
	const mobileOtpFailureCoolDownDurationDefault = 2
	const mobileOtpFailureCoolDownDurationMinMinutes = 2
	const mobileOtpFailureCoolDownDurationMaxMinutes = 30
	const mobileOtpFailureCoolDownDurationMinSeconds = 120
	const mobileOtpFailureCoolDownDurationMaxSeconds = 1800

	const mobileOtpFailureCountDefault = 3
	const mobileOtpFailureCountMin = 1
	const mobileOtpFailureCountMax = 7

	const totpOtpFailureCountDefault = 3
	const totpOtpFailureCoolDownDurationDefault = 2

	const rememberMeWebLifeTimeDurationDefault = 30
	const rememberMeWebLifeTimeDurationMinMinutes = 1
	const rememberMeWebLifeTimeDurationMaxMinutes = 129600
	const rememberMeWebLifeTimeDurationMinHours = 1
	const rememberMeWebLifeTimeDurationMaxHours = 2160
	const rememberMeWebLifeTimeDurationMinDays = 1
	const rememberMeWebLifeTimeDurationMaxDays = 90

	const mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDefault = 25
	const mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMin = 15
	const mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMax = 75

	const mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDefault = 40
	const mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMin = 30
	const mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMax = 90

	const pingidDeviceOtpFailureCountDefault = 3
	const pingidDeviceOtpFailureCountMin = 1
	const pingidDeviceOtpFailureCountMax = 7

	const pingidDeviceOtpFailureCoolDownDurationDefault = 2
	const pingidDeviceOtpFailureCoolDownDurationMinSeconds = 1
	const pingidDeviceOtpFailureCoolDownDurationMaxSeconds = 1800
	const pingidDeviceOtpFailureCoolDownDurationMinMinutes = 1
	const pingidDeviceOtpFailureCoolDownDurationMaxMinutes = 30

	const pingidDevicePairingKeyLifetimeDurationMinMinutes = 1
	const pingidDevicePairingKeyLifetimeDurationMaxMinutes = 2880
	const pingidDevicePairingKeyLifetimeDurationMinHours = 1
	const pingidDevicePairingKeyLifetimeDurationMaxHours = 48

	// Default value for remember_me
	rememberMeDefault := types.ObjectValueMust(
		MFADevicePolicyRememberMeTFObjectTypes,
		map[string]attr.Value{
			"web": types.ObjectValueMust(
				MFADevicePolicyRememberMeWebTFObjectTypes,
				map[string]attr.Value{
					"enabled": types.BoolValue(false),
					"life_time": types.ObjectValueMust(
						MFADevicePolicyTimePeriodTFObjectTypes,
						map[string]attr.Value{
							"duration":  types.Int32Value(rememberMeWebLifeTimeDurationDefault),
							"time_unit": types.StringValue(string(mfa.ENUMTIMEUNITREMEMBERMEWEBLIFETIME_MINUTES)),
						},
					),
				},
			),
		},
	)

	// Default values for oath_token
	oathTokenDefault := types.ObjectValueMust(
		MFADevicePolicyPingIDDeviceTFObjectTypes,
		map[string]attr.Value{
			"enabled": types.BoolValue(false),
			"otp": types.ObjectValueMust(
				MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
				map[string]attr.Value{
					"failure": types.ObjectValueMust(
						MFADevicePolicyFailureTFObjectTypes,
						map[string]attr.Value{
							"count": types.Int32Value(pingidDeviceOtpFailureCountDefault),
							"cool_down": types.ObjectValueMust(
								MFADevicePolicyTimePeriodTFObjectTypes,
								map[string]attr.Value{
									"duration":  types.Int32Value(pingidDeviceOtpFailureCoolDownDurationDefault),
									"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
								},
							),
						},
					),
				},
			),
			"pairing_disabled":               types.BoolValue(false),
			"pairing_key_lifetime":           types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes),
			"prompt_for_nickname_on_pairing": types.BoolValue(false),
		},
	)

	// schema descriptions and validation settings
	deviceSelectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines the device selection method.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicySelectionEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST))

	newDeviceNotificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines whether a user should be notified if a new authentication method has been added to their account.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicyNewDeviceNotificationEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE))

	ignoreUserLockDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, allows PingOne to skip the account lock check during MFA authentication.",
	).DefaultValue(false)

	updatedAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time the resource was last updated.",
	)

	notificationsPolicyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the notification policy to use for this MFA device policy. If not specified, the default notification policy for the environment will be used.  **Note:** When destroying this resource, the `notifications_policy` will be unset (set to null) to release any dependencies, allowing the referenced notification policy to be deleted if needed.",
	)

	notificationsPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the notification policy to use.",
	)

	mobileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of mobile push/OTP device authentication policy settings. This factor requires embedding the PingOne MFA SDK into a customer facing mobile application, and configuring as a Native application using the `pingone_application` resource.",
	)

	mobileApplicationsAutoEnrollmentEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, allows Auto Enrollment for the application.",
	)

	mobileApplicationsDeviceAuthorizationExtraVerificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the level of further verification when device authorization is enabled.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicyMobileExtraVerificationEnumValues)

	mobileApplicationsIntegrityDetectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that controls how authentication or registration attempts should proceed if a device integrity check does not receive a response.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicyMobileIntegrityDetectionEnumValues)

	mobileApplicationsPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the relevant application. You can use this option if you want to phase out an existing mobile application but want to allow users to continue using the application for authentication for existing devices.",
	)

	mobileApplicationsPushLimitLockDurationDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. Defaults to `30` minutes.",
	)

	mobileApplicationsPushLimitTimePeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. Defaults to `30` minutes.",
	)

	mobileApplicationsPushTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `40` seconds and the maximum value is `150` seconds. Defaults to `40` seconds.",
	)

	mobileOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. The minimum value is `%d`, maximum is `%d`, and the default is `%d`.", mobileOtpFailureCountMin, mobileOtpFailureCountMax, mobileOtpFailureCountDefault),
	)

	mobileOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. The minimum value is `2`, maximum is `30`.",
	).DefaultValue(mobileOtpFailureCoolDownDurationDefault)

	mobileIpPairingConfigurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows you to restrict device pairing to specific IP addresses. Only applicable for PingID policies.",
	)

	mobileIpPairingConfigurationAnyIpAddressDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `false`, restricts device pairing to specific IP addresses defined in `only_these_ip_addresses`.",
	).DefaultValue(true)

	mobileIpPairingConfigurationOnlyTheseIpAddressesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of IP addresses or address ranges from which users can pair their devices. This parameter is required when `any_ip_address` is set to `false`. Each item in the array must be in CIDR notation, for example, `192.168.1.1/32` or `10.0.0.0/8`.",
	)

	mobileApplicationsNewRequestDurationConfigurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that configures timeout settings for authentication request notifications. Only applicable for PingID policies.",
	)

	mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that specifies the maximum time a notification can remain pending before it is displayed to the user. Value must be between `%d` and `%d` seconds.", mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMin, mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMax),
	).DefaultValue(mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDefault)

	mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that specifies the total time an authentication request notification has to be handled by the user before timing out. The `total_timeout.duration` must exceed `device_timeout.duration` by at least 15 seconds.  Value must be between `%d` and `%d` seconds.", mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMin, mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMax),
	).DefaultValue(mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDefault)

	mobileApplicationsNewRequestDurationConfigurationTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the timeout duration in seconds.",
	).AllowedValuesEnum(mfa.ENUMTIMEUNIT_SECONDS)

	mobileApplicationsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the application type. Only applicable when `policy_type` is `PINGID`. Must be set to `pingIdAppConfig`.",
	)

	rememberMeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies 'remember me' settings so that users do not have to authenticate when accessing applications from a device they have used already.",
	)

	rememberMeWebDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that contains the 'remember me' settings for accessing applications from a browser.",
	)

	rememberMeWebEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, enables the 'remember me' option in the MFA policy.",
	)

	rememberMeWebLifeTimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that defines the period during which users will not have to authenticate if they are accessing applications from a device they have used before. The 'remember me' period can be anywhere from `1` minute to `90` days.",
	)

	rememberMeWebLifeTimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that, used in conjunction with `time_unit`, defines the 'remember me' period.",
	)

	rememberMeWebLifeTimeTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time unit to use for the 'remember me' period.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitRememberMeWebLifeTimeEnumValues)

	durationTimeUnitMinsSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	mobileApplicationsPairingKeyLifetimeTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)

	mobileApplicationsPairingKeyLifetimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies pairing key lifetime settings for the application in the policy. Defaults to 10 minutes for PingOne MFA policies and 48 hours for PingID policies.",
	)

	durationTimeUnitSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).DefaultValue(string(mfa.ENUMTIMEUNIT_SECONDS)).AllowedValuesEnum(mfa.ENUMTIMEUNIT_SECONDS)

	totpPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the TOTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	totpUriParametersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map of string key:value pairs that specifies `otpauth` URI parameters. For example, if you provide a value for the `issuer` parameter, then authenticators that support that parameter will display the text you specify together with the OTP (in addition to the username). This can help users recognize which application the OTP is for. If you intend on using the same MFA policy for multiple applications, choose a name that reflects the group of applications.",
	)

	fido2PairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the FIDO2 method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	promptForNicknameOnPairingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.",
	)

	policyTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of MFA device policy.",
	).AllowedValues(POLICY_TYPE_PINGONE_MFA, POLICY_TYPE_PINGID)

	environmentIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the environment to manage the default MFA device policy in.",
	)

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name to apply to the default MFA device policy.",
	)

	authenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of authentication settings in the device policy.",
	)

	mobileEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the mobile device method is enabled or disabled in the policy.",
	)

	mobileApplicationsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map of objects that specifies settings for configured Mobile Applications. The ID of the application should be configured as the map key.",
	)

	mobileApplicationsAutoEnrollmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies auto enrollment settings for the application in the policy.",
	)

	mobileApplicationsBiometricsEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether biometric authentication methods (such as fingerprint or facial recognition) are enabled for MFA. Only applicable for PingID policies.",
	)

	mobileApplicationsDeviceAuthorizationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies device authorization settings for the application in the policy.",
	)

	mobileApplicationsDeviceAuthorizationEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies the enabled or disabled state of automatic MFA for native devices paired with the user, for the specified application.",
	)

	mobileApplicationsOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP settings for the application in the policy.",
	)

	mobileApplicationsOtpEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether OTP authentication is enabled or disabled for the application in the policy.",
	)

	mobileApplicationsPairingKeyLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the amount of time an issued pairing key can be used until it expires. Minimum is `%d` minute and maximum is `%d` hours. If this parameter is not provided, the duration is set to `10` minutes.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxHours),
	)

	mobileApplicationsPushDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies push settings for the application in the policy.",
	)

	mobileApplicationsPushEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether push notification is enabled or disabled for the application in the policy.",
	)

	mobileApplicationsPushNumberMatchingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that configures number matching for push notifications.",
	)

	mobileApplicationsPushNumberMatchingEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, requires the authenticating user to select a number that was displayed to them on the accessing device.",
	)

	mobileApplicationsPushLimitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies push limit settings for the application in the policy.",
	)

	mobileApplicationsPushLimitCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is `1` and the maximum value is `50`.",
	).DefaultValue(mobileApplicationsPushLimitCountDefault)

	mobileApplicationsPushLimitLockDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies push limit lock duration settings for the application in the policy.",
	)

	mobileApplicationsPushLimitTimePeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies push limit time period settings for the application in the policy.",
	)

	mobileApplicationsPushTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies push timeout settings for the application in the policy.",
	)

	mobileOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP settings for mobile applications in the policy.",
	)

	mobileOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure settings for mobile applications in the policy.",
	)

	mobileOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure cool down settings for mobile applications in the policy.",
	)

	totpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of TOTP device authentication policy settings.",
	)

	totpEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the TOTP method is enabled or disabled in the policy.",
	)

	totpOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of TOTP OTP settings.",
	)

	totpOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of TOTP OTP failure settings.",
	)

	totpOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of TOTP OTP failure cool down settings.",
	)

	totpOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
	)

	totpOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
	)

	fido2Description := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of FIDO2 device authentication policy settings.",
	)

	fido2EnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the FIDO2 method is enabled or disabled in the policy.",
	)

	fido2PolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the resource UUID that represents the FIDO2 policy in PingOne. When null, the environment's default FIDO2 Policy is used.",
	)

	desktopDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of PingID desktop device authentication policy settings. Only applicable when `policy_type` is `PINGID`.",
	)

	desktopEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the desktop device method is enabled or disabled in the policy.",
	)

	desktopOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure settings for desktop devices.",
	)

	desktopOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of OTP failure settings.",
	)

	desktopOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between %d and %d.", pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
	)

	desktopOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure cool down settings.",
	)

	desktopOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `%d` seconds and `%d` minutes.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
	)

	desktopPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new desktop devices.",
	)

	desktopPairingKeyLifetimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies pairing key lifetime settings for desktop devices.",
	)

	desktopPairingKeyLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between %d minutes and %d hours.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxHours),
	)

	yubikeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of PingID Yubikey device authentication policy settings. Only applicable when `policy_type` is `PINGID`.",
	)

	yubikeyEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the Yubikey device method is enabled or disabled in the policy.",
	)

	yubikeyOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure settings for Yubikey devices.",
	)

	yubikeyOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of OTP failure settings.",
	)

	yubikeyOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between %d and %d.", pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
	)

	yubikeyOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure cool down settings.",
	)

	yubikeyOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `%d` seconds and `%d` minutes.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
	)

	yubikeyPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new Yubikey devices.",
	)

	yubikeyPairingKeyLifetimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies pairing key lifetime settings for Yubikey devices.",
	)

	yubikeyPairingKeyLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between %d minutes and %d hours.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxHours),
	)

	oathTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of OATH token device authentication policy settings.",
	)

	oathTokenEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the OATH token device method is enabled or disabled in the policy.",
	)

	oathTokenOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure settings for OATH token devices.",
	)

	oathTokenOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of OTP failure settings.",
	)

	oathTokenOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between `%d` and `%d`.", pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
	)

	oathTokenOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies OTP failure cool down settings.",
	)

	oathTokenOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `%d` seconds and `%d` minutes.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
	)

	oathTokenPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new OATH token devices.",
	)

	oathTokenPairingKeyLifetimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies pairing key lifetime settings for OATH token devices.",
	)

	oathTokenPairingKeyLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between `%d` minutes and `%d` hours.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxHours),
	)

	resourceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Resource to overwrite the default MFA device policy. Valid for both PingOne MFA and PingID integrated environments.",
	)

	resp.Schema = schema.Schema{
		Description:         resourceDescription.Description,
		MarkdownDescription: resourceDescription.MarkdownDescription,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				environmentIdDescription,
			),

			"policy_type": schema.StringAttribute{
				Description:         policyTypeDescription.Description,
				MarkdownDescription: policyTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(POLICY_TYPE_PINGONE_MFA, POLICY_TYPE_PINGID),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
			},

			"authentication": schema.SingleNestedAttribute{
				Description:         authenticationDescription.Description,
				MarkdownDescription: authenticationDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					MFADevicePolicyAuthenticationTFObjectTypes,
					map[string]attr.Value{
						"device_selection": types.StringValue(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"device_selection": schema.StringAttribute{
						Description:         deviceSelectionDescription.Description,
						MarkdownDescription: deviceSelectionDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicySelectionEnumValues)...),
						},
					},
				},
			},

			"new_device_notification": schema.StringAttribute{
				Description:         newDeviceNotificationDescription.Description,
				MarkdownDescription: newDeviceNotificationDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicyNewDeviceNotificationEnumValues)...),
				},
			},

			"ignore_user_lock": schema.BoolAttribute{
				Description:         ignoreUserLockDescription.Description,
				MarkdownDescription: ignoreUserLockDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"updated_at": schema.StringAttribute{
				Description:         updatedAtDescription.Description,
				MarkdownDescription: updatedAtDescription.MarkdownDescription,
				Computed:            true,

				CustomType: timetypes.RFC3339Type{},
			},

			"notifications_policy": schema.SingleNestedAttribute{
				Description:         notificationsPolicyDescription.Description,
				MarkdownDescription: notificationsPolicyDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         notificationsPolicyIdDescription.Description,
						MarkdownDescription: notificationsPolicyIdDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							verify.P1ResourceIDValidator(),
						},
					},
				},
			},

			"remember_me": schema.SingleNestedAttribute{
				Description:         rememberMeDescription.Description,
				MarkdownDescription: rememberMeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(rememberMeDefault),

				Attributes: map[string]schema.Attribute{
					"web": schema.SingleNestedAttribute{
						Description:         rememberMeWebDescription.Description,
						MarkdownDescription: rememberMeWebDescription.MarkdownDescription,
						Required:            true,

						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         rememberMeWebEnabledDescription.Description,
								MarkdownDescription: rememberMeWebEnabledDescription.MarkdownDescription,
								Required:            true,
							},

							"life_time": schema.SingleNestedAttribute{
								Description:         rememberMeWebLifeTimeDescription.Description,
								MarkdownDescription: rememberMeWebLifeTimeDescription.MarkdownDescription,
								Required:            true,

								Attributes: map[string]schema.Attribute{
									"duration": schema.Int32Attribute{
										Description:         rememberMeWebLifeTimeDurationDescription.Description,
										MarkdownDescription: rememberMeWebLifeTimeDurationDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.Int32{
											int32validator.Any(
												int32validator.All(
													int32validator.Between(rememberMeWebLifeTimeDurationMinMinutes, rememberMeWebLifeTimeDurationMaxMinutes),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d (1 minute to 90 days).", rememberMeWebLifeTimeDurationMinMinutes, rememberMeWebLifeTimeDurationMaxMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int32validator.All(
													int32validator.Between(rememberMeWebLifeTimeDurationMinHours, rememberMeWebLifeTimeDurationMaxHours),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`HOURS`),
														fmt.Sprintf("If `time_unit` is `HOURS`, the allowed duration range is %d - %d (1 hour to 90 days).", rememberMeWebLifeTimeDurationMinHours, rememberMeWebLifeTimeDurationMaxHours),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int32validator.All(
													int32validator.Between(rememberMeWebLifeTimeDurationMinDays, rememberMeWebLifeTimeDurationMaxDays),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`DAYS`),
														fmt.Sprintf("If `time_unit` is `DAYS`, the allowed duration range is %d - %d.", rememberMeWebLifeTimeDurationMinDays, rememberMeWebLifeTimeDurationMaxDays),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
											),
										},
									},

									"time_unit": schema.StringAttribute{
										Description:         rememberMeWebLifeTimeTimeUnitDescription.Description,
										MarkdownDescription: rememberMeWebLifeTimeTimeUnitDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitRememberMeWebLifeTimeEnumValues)...),
										},
									},
								},
							},
						},
					},
				},
			},

			"sms": r.devicePolicyOfflineDeviceSchemaAttribute("SMS OTP"),

			"voice": r.devicePolicyOfflineDeviceSchemaAttribute("voice OTP"),

			"email": r.devicePolicyOfflineDeviceSchemaAttribute("email OTP"),

			"mobile": schema.SingleNestedAttribute{
				Description:         mobileDescription.Description,
				MarkdownDescription: mobileDescription.MarkdownDescription,
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         mobileEnabledDescription.Description,
						MarkdownDescription: mobileEnabledDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.Bool{
							boolvalidator.MustBeTrueIfPathSetToValue(
								types.StringValue(POLICY_TYPE_PINGID),
								path.MatchRoot("policy_type"),
							),
						},
					},
					"applications": schema.MapNestedAttribute{
						Description:         mobileApplicationsDescription.Description,
						MarkdownDescription: mobileApplicationsDescription.MarkdownDescription,
						Optional:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"auto_enrollment": schema.SingleNestedAttribute{
									Description:         mobileApplicationsAutoEnrollmentDescription.Description,
									MarkdownDescription: mobileApplicationsAutoEnrollmentDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.Object{
										objectvalidator.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
										objectvalidator.ConflictsIfMatchesPathValue(types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
									},

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         mobileApplicationsAutoEnrollmentEnabledDescription.Description,
											MarkdownDescription: mobileApplicationsAutoEnrollmentEnabledDescription.MarkdownDescription,
											Required:            true,
										},
									},
								},
								"biometrics_enabled": schema.BoolAttribute{
									Description:         mobileApplicationsBiometricsEnabledDescription.Description,
									MarkdownDescription: mobileApplicationsBiometricsEnabledDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Validators: []validator.Bool{
										boolvalidator.ConflictsIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
									},
								},
								"device_authorization": schema.SingleNestedAttribute{
									Description:         mobileApplicationsDeviceAuthorizationDescription.Description,
									MarkdownDescription: mobileApplicationsDeviceAuthorizationDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.Object{
										objectvalidator.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
										objectvalidator.ConflictsIfMatchesPathValue(types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
									},

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         mobileApplicationsDeviceAuthorizationEnabledDescription.Description,
											MarkdownDescription: mobileApplicationsDeviceAuthorizationEnabledDescription.MarkdownDescription,
											Required:            true,
										},

										"extra_verification": schema.StringAttribute{
											Description:         mobileApplicationsDeviceAuthorizationExtraVerificationDescription.Description,
											MarkdownDescription: mobileApplicationsDeviceAuthorizationExtraVerificationDescription.MarkdownDescription,
											Optional:            true,

											Validators: []validator.String{
												stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicyMobileExtraVerificationEnumValues)...),
											},
										},
									},
								},

								"integrity_detection": schema.StringAttribute{
									Description:         mobileApplicationsIntegrityDetectionDescription.Description,
									MarkdownDescription: mobileApplicationsIntegrityDetectionDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicyMobileIntegrityDetectionEnumValues)...),
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
									},
								},
								"otp": schema.SingleNestedAttribute{
									Description:         mobileApplicationsOtpDescription.Description,
									MarkdownDescription: mobileApplicationsOtpDescription.MarkdownDescription,
									Required:            true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         mobileApplicationsOtpEnabledDescription.Description,
											MarkdownDescription: mobileApplicationsOtpEnabledDescription.MarkdownDescription,
											Required:            true,
										},
									},
								},
								"pairing_disabled": schema.BoolAttribute{
									Description:         mobileApplicationsPairingDisabledDescription.Description,
									MarkdownDescription: mobileApplicationsPairingDisabledDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Default: booldefault.StaticBool(false),
								},
								"pairing_key_lifetime": schema.SingleNestedAttribute{
									Description:         mobileApplicationsPairingKeyLifetimeDescription.Description,
									MarkdownDescription: mobileApplicationsPairingKeyLifetimeDescription.MarkdownDescription,
									Optional:            true,

									Attributes: map[string]schema.Attribute{
										"duration": schema.Int32Attribute{
											Description:         mobileApplicationsPairingKeyLifetimeDurationDescription.Description,
											MarkdownDescription: mobileApplicationsPairingKeyLifetimeDurationDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.Int32{
												int32validator.Any(
													int32validator.All(
														int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
														int32validatorinternal.RegexMatchesPathValue(
															regexp.MustCompile(`MINUTES`),
															fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
															path.MatchRelative().AtParent().AtName("time_unit"),
														),
													),
													int32validator.All(
														int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
														int32validatorinternal.RegexMatchesPathValue(
															regexp.MustCompile(`HOURS`),
															fmt.Sprintf("If `time_unit` is `HOURS`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
															path.MatchRelative().AtParent().AtName("time_unit"),
														),
													),
												),
											},
										},

										"time_unit": schema.StringAttribute{
											Description:         mobileApplicationsPairingKeyLifetimeTimeUnitDescription.Description,
											MarkdownDescription: mobileApplicationsPairingKeyLifetimeTimeUnitDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.String{
												stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)...),
											},
										},
									},
								},

								"push": schema.SingleNestedAttribute{
									Description:         mobileApplicationsPushDescription.Description,
									MarkdownDescription: mobileApplicationsPushDescription.MarkdownDescription,
									Optional:            true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         mobileApplicationsPushEnabledDescription.Description,
											MarkdownDescription: mobileApplicationsPushEnabledDescription.MarkdownDescription,
											Required:            true,
										},

										"number_matching": schema.SingleNestedAttribute{
											Description:         mobileApplicationsPushNumberMatchingDescription.Description,
											MarkdownDescription: mobileApplicationsPushNumberMatchingDescription.MarkdownDescription,
											Optional:            true,
											Computed:            true,

											Default: objectdefault.StaticValue(types.ObjectValueMust(
												MFADevicePolicyDefaultMobileApplicationPushNumberMatchingTFObjectTypes,
												map[string]attr.Value{
													"enabled": types.BoolValue(false),
												},
											)),

											Attributes: map[string]schema.Attribute{
												"enabled": schema.BoolAttribute{
													Description:         mobileApplicationsPushNumberMatchingEnabledDescription.Description,
													MarkdownDescription: mobileApplicationsPushNumberMatchingEnabledDescription.MarkdownDescription,
													Required:            true,
												},
											},
										},
									},
								},
								"push_limit": schema.SingleNestedAttribute{
									Description:         mobileApplicationsPushLimitDescription.Description,
									MarkdownDescription: mobileApplicationsPushLimitDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Default: objectdefault.StaticValue(types.ObjectValueMust(
										MFADevicePolicyMobileApplicationPushLimitTFObjectTypes,
										map[string]attr.Value{
											"count": types.Int32Value(mobileApplicationsPushLimitCountDefault),
											"lock_duration": types.ObjectValueMust(
												MFADevicePolicyTimePeriodTFObjectTypes,
												map[string]attr.Value{
													"duration":  types.Int32Value(mobileApplicationsPushLimitLockDurationDurationDefault),
													"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
												},
											),
											"time_period": types.ObjectValueMust(
												MFADevicePolicyTimePeriodTFObjectTypes,
												map[string]attr.Value{
													"duration":  types.Int32Value(mobileApplicationsPushLimitTimePeriodDurationDefault),
													"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
												},
											),
										},
									)),

									Attributes: map[string]schema.Attribute{
										"count": schema.Int32Attribute{
											Description:         mobileApplicationsPushLimitCountDescription.Description,
											MarkdownDescription: mobileApplicationsPushLimitCountDescription.MarkdownDescription,
											Optional:            true,
											Computed:            true,

											Default: int32default.StaticInt32(mobileApplicationsPushLimitCountDefault),

											Validators: []validator.Int32{
												int32validator.Between(mobileApplicationsPushLimitCountMin, mobileApplicationsPushLimitCountMax),
											},
										},

										"lock_duration": schema.SingleNestedAttribute{
											Description:         mobileApplicationsPushLimitLockDurationDescription.Description,
											MarkdownDescription: mobileApplicationsPushLimitLockDurationDescription.MarkdownDescription,
											Optional:            true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int32Attribute{
													Description:         mobileApplicationsPushLimitLockDurationDurationDescription.Description,
													MarkdownDescription: mobileApplicationsPushLimitLockDurationDurationDescription.MarkdownDescription,
													Required:            true,

													Validators: []validator.Int32{
														int32validator.Any(
															int32validator.All(
																int32validator.Between(mobileApplicationsPushLimitLockDurationDurationMinSeconds, mobileApplicationsPushLimitLockDurationDurationMaxSeconds),
																int32validatorinternal.RegexMatchesPathValue(
																	regexp.MustCompile(`SECONDS`),
																	fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", mobileApplicationsPushLimitLockDurationDurationMinSeconds, mobileApplicationsPushLimitLockDurationDurationMaxSeconds),
																	path.MatchRelative().AtParent().AtName("time_unit"),
																),
															),
															int32validator.All(
																int32validator.Between(mobileApplicationsPushLimitLockDurationDurationMinMinutes, mobileApplicationsPushLimitLockDurationDurationMaxMinutes),
																int32validatorinternal.RegexMatchesPathValue(
																	regexp.MustCompile(`MINUTES`),
																	fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", mobileApplicationsPushLimitLockDurationDurationMinMinutes, mobileApplicationsPushLimitLockDurationDurationMaxMinutes),
																	path.MatchRelative().AtParent().AtName("time_unit"),
																),
															),
														),
													},
												},

												"time_unit": schema.StringAttribute{
													Description:         durationTimeUnitMinsSecondsDescription.Description,
													MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
													Required:            true,

													Validators: []validator.String{
														stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
													},
												},
											},
										},

										"time_period": schema.SingleNestedAttribute{
											Description:         mobileApplicationsPushLimitTimePeriodDescription.Description,
											MarkdownDescription: mobileApplicationsPushLimitTimePeriodDescription.MarkdownDescription,
											Optional:            true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int32Attribute{
													Description:         mobileApplicationsPushLimitTimePeriodDurationDescription.Description,
													MarkdownDescription: mobileApplicationsPushLimitTimePeriodDurationDescription.MarkdownDescription,
													Required:            true,

													Validators: []validator.Int32{
														int32validator.Any(
															int32validator.All(
																int32validator.Between(mobileApplicationsPushLimitTimePeriodDurationMinSeconds, mobileApplicationsPushLimitTimePeriodDurationMaxSeconds),
																int32validatorinternal.RegexMatchesPathValue(
																	regexp.MustCompile(`SECONDS`),
																	fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", mobileApplicationsPushLimitTimePeriodDurationMinSeconds, mobileApplicationsPushLimitTimePeriodDurationMaxSeconds),
																	path.MatchRelative().AtParent().AtName("time_unit"),
																),
															),
															int32validator.All(
																int32validator.Between(mobileApplicationsPushLimitTimePeriodDurationMinMinutes, mobileApplicationsPushLimitTimePeriodDurationMaxMinutes),
																int32validatorinternal.RegexMatchesPathValue(
																	regexp.MustCompile(`MINUTES`),
																	fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", mobileApplicationsPushLimitTimePeriodDurationMinMinutes, mobileApplicationsPushLimitTimePeriodDurationMaxMinutes),
																	path.MatchRelative().AtParent().AtName("time_unit"),
																),
															),
														),
													},
												},

												"time_unit": schema.StringAttribute{
													Description:         durationTimeUnitMinsSecondsDescription.Description,
													MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
													Required:            true,

													Validators: []validator.String{
														stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
													},
												},
											},
										},
									},
								},

								"push_timeout": schema.SingleNestedAttribute{
									Description:         mobileApplicationsPushTimeoutDescription.Description,
									MarkdownDescription: mobileApplicationsPushTimeoutDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.Object{
										objectvalidator.ConflictsIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
									},

									Attributes: map[string]schema.Attribute{
										"duration": schema.Int32Attribute{
											Description:         mobileApplicationsPushTimeoutDurationDescription.Description,
											MarkdownDescription: mobileApplicationsPushTimeoutDurationDescription.MarkdownDescription,
											Required:            true,

											Validators: []validator.Int32{
												int32validator.Between(mobileApplicationsPushTimeoutDurationMin, mobileApplicationsPushTimeoutDurationMax),
											},
										},

										"time_unit": schema.StringAttribute{
											Description:         durationTimeUnitSecondsDescription.Description,
											MarkdownDescription: durationTimeUnitSecondsDescription.MarkdownDescription,
											Optional:            true,
											Computed:            true,

											Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_SECONDS)),

											Validators: []validator.String{
												stringvalidator.OneOf(string(mfa.ENUMTIMEUNIT_SECONDS)),
											},
										},
									},
								},

								"new_request_duration_configuration": schema.SingleNestedAttribute{
									Description:         mobileApplicationsNewRequestDurationConfigurationDescription.Description,
									MarkdownDescription: mobileApplicationsNewRequestDurationConfigurationDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Validators: []validator.Object{
										objectvalidator.ConflictsIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
										objectvalidator.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
									}, Attributes: map[string]schema.Attribute{
										"device_timeout": schema.SingleNestedAttribute{
											Description:         mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDescription.Description,
											MarkdownDescription: mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDescription.MarkdownDescription,
											Required:            true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int32Attribute{
													Description:         mobileApplicationsNewRequestDurationConfigurationTimeoutDurationDescription.Description,
													MarkdownDescription: mobileApplicationsNewRequestDurationConfigurationTimeoutDurationDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: int32default.StaticInt32(mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutDefault),

													Validators: []validator.Int32{
														int32validator.Between(mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMin, mobileApplicationsNewRequestDurationConfigurationDeviceTimeoutMax),
													},
												}, "time_unit": schema.StringAttribute{
													Description:         durationTimeUnitSecondsDescription.Description,
													MarkdownDescription: durationTimeUnitSecondsDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_SECONDS)),

													Validators: []validator.String{
														stringvalidator.OneOf(string(mfa.ENUMTIMEUNIT_SECONDS)),
													},
												},
											},
										},

										"total_timeout": schema.SingleNestedAttribute{
											Description:         mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDescription.Description,
											MarkdownDescription: mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDescription.MarkdownDescription,
											Required:            true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int32Attribute{
													Description:         mobileApplicationsNewRequestDurationConfigurationTimeoutDurationDescription.Description,
													MarkdownDescription: mobileApplicationsNewRequestDurationConfigurationTimeoutDurationDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: int32default.StaticInt32(mobileApplicationsNewRequestDurationConfigurationTotalTimeoutDefault),

													Validators: []validator.Int32{
														int32validator.Between(mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMin, mobileApplicationsNewRequestDurationConfigurationTotalTimeoutMax),
													},
												}, "time_unit": schema.StringAttribute{
													Description:         durationTimeUnitSecondsDescription.Description,
													MarkdownDescription: durationTimeUnitSecondsDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_SECONDS)),

													Validators: []validator.String{
														stringvalidator.OneOf(string(mfa.ENUMTIMEUNIT_SECONDS)),
													},
												},
											},
										},
									},
								},

								"type": schema.StringAttribute{
									Description:         mobileApplicationsTypeDescription.Description,
									MarkdownDescription: mobileApplicationsTypeDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
										stringvalidatorinternal.ConflictsIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
									},
								},

								"ip_pairing_configuration": schema.SingleNestedAttribute{
									Description:         mobileIpPairingConfigurationDescription.Description,
									MarkdownDescription: mobileIpPairingConfigurationDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.Object{
										objectvalidator.ConflictsIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGONE_MFA),
											path.MatchRoot("policy_type"),
										),
										objectvalidator.IsRequiredIfMatchesPathValue(
											types.StringValue(POLICY_TYPE_PINGID),
											path.MatchRoot("policy_type"),
										),
									},

									Attributes: map[string]schema.Attribute{
										"any_ip_address": schema.BoolAttribute{
											Description:         mobileIpPairingConfigurationAnyIpAddressDescription.Description,
											MarkdownDescription: mobileIpPairingConfigurationAnyIpAddressDescription.MarkdownDescription,
											Optional:            true,
											Computed:            true,

											Default: booldefault.StaticBool(true),
										},

										"only_these_ip_addresses": schema.SetAttribute{
											Description:         mobileIpPairingConfigurationOnlyTheseIpAddressesDescription.Description,
											MarkdownDescription: mobileIpPairingConfigurationOnlyTheseIpAddressesDescription.MarkdownDescription,
											ElementType:         types.StringType,
											Optional:            true,

											Validators: []validator.Set{
												setvalidator.ValueStringsAre(
													stringvalidator.RegexMatches(regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}/\d{1,2}$`), "Expected value to be in CIDR notation (e.g., 192.168.0.1/24 or 10.0.0.5/32)"),
												),
												setvalidatorinternal.IsRequiredIfMatchesPathBoolValue(
													types.BoolValue(false),
													path.MatchRelative().AtParent().AtName("any_ip_address"),
												),
											},
										},
									},
								},
							},
						},
					},
					"otp": schema.SingleNestedAttribute{
						Description:         mobileOtpDescription.Description,
						MarkdownDescription: mobileOtpDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyMobileOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyMobileOtpFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(mobileOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(mobileApplicationsOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description:         mobileOtpFailureDescription.Description,
								MarkdownDescription: mobileOtpFailureDescription.MarkdownDescription,
								Required:            true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description:         mobileOtpFailureCountDescription.Description,
										MarkdownDescription: mobileOtpFailureCountDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int32default.StaticInt32(mobileOtpFailureCountDefault),

										Validators: []validator.Int32{
											int32validator.Between(mobileOtpFailureCountMin, mobileOtpFailureCountMax),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description:         mobileOtpFailureCoolDownDescription.Description,
										MarkdownDescription: mobileOtpFailureCoolDownDescription.MarkdownDescription,
										Required:            true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         mobileOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: mobileOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(mobileOtpFailureCoolDownDurationMinSeconds, mobileOtpFailureCoolDownDurationMaxSeconds),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", mobileOtpFailureCoolDownDurationMinSeconds, mobileOtpFailureCoolDownDurationMaxSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(mobileOtpFailureCoolDownDurationMinMinutes, mobileOtpFailureCoolDownDurationMaxMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", mobileOtpFailureCoolDownDurationMinMinutes, mobileOtpFailureCoolDownDurationMaxMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},

											"time_unit": schema.StringAttribute{
												Description:         durationTimeUnitMinsSecondsDescription.Description,
												MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
						},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"totp": schema.SingleNestedAttribute{
				Description:         totpDescription.Description,
				MarkdownDescription: totpDescription.MarkdownDescription,
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         totpEnabledDescription.Description,
						MarkdownDescription: totpEnabledDescription.MarkdownDescription,
						Required:            true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         totpPairingDisabledDescription.Description,
						MarkdownDescription: totpPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"otp": schema.SingleNestedAttribute{
						Description:         totpOtpDescription.Description,
						MarkdownDescription: totpOtpDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyTotpOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyTotpOtpFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(totpOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(totpOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description:         totpOtpFailureDescription.Description,
								MarkdownDescription: totpOtpFailureDescription.MarkdownDescription,
								Optional:            true,

								Attributes: map[string]schema.Attribute{
									"cool_down": schema.SingleNestedAttribute{
										Description:         totpOtpFailureCoolDownDescription.Description,
										MarkdownDescription: totpOtpFailureCoolDownDescription.MarkdownDescription,
										Optional:            true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         totpOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: totpOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},

											"time_unit": schema.StringAttribute{
												Description:         durationTimeUnitMinsSecondsDescription.Description,
												MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},

									"count": schema.Int32Attribute{
										Description:         totpOtpFailureCountDescription.Description,
										MarkdownDescription: totpOtpFailureCountDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.Int32{
											int32validator.Between(pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
										},
									},
								},
							},
						},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},

					"uri_parameters": schema.MapAttribute{
						Description:         totpUriParametersDescription.Description,
						MarkdownDescription: totpUriParametersDescription.MarkdownDescription,
						Optional:            true,

						ElementType: types.StringType,
					},
				},
			},

			"fido2": schema.SingleNestedAttribute{
				Description:         fido2Description.Description,
				MarkdownDescription: fido2Description.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(MFADevicePolicyFido2TFObjectTypes,
					map[string]attr.Value{
						"enabled":                        types.BoolValue(false),
						"fido2_policy_id":                framework.PingOneResourceIDToTF(""),
						"pairing_disabled":               types.BoolValue(false),
						"prompt_for_nickname_on_pairing": types.BoolNull(),
					}),
				),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         fido2EnabledDescription.Description,
						MarkdownDescription: fido2EnabledDescription.MarkdownDescription,
						Required:            true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         fido2PairingDisabledDescription.Description,
						MarkdownDescription: fido2PairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"fido2_policy_id": schema.StringAttribute{
						Description:         fido2PolicyIdDescription.Description,
						MarkdownDescription: fido2PolicyIdDescription.MarkdownDescription,
						Optional:            true,
						CustomType:          pingonetypes.ResourceIDType{},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"desktop": schema.SingleNestedAttribute{
				Description:         desktopDescription.Description,
				MarkdownDescription: desktopDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidator.IsRequiredIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGID),
						path.MatchRoot("policy_type"),
					),
					objectvalidator.ConflictsIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGONE_MFA),
						path.MatchRoot("policy_type"),
					),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         desktopEnabledDescription.Description,
						MarkdownDescription: desktopEnabledDescription.MarkdownDescription,
						Required:            true,
					},

					"otp": schema.SingleNestedAttribute{
						Description:         desktopOtpDescription.Description,
						MarkdownDescription: desktopOtpDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(pingidDeviceOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(pingidDeviceOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),
						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description:         desktopOtpFailureDescription.Description,
								MarkdownDescription: desktopOtpFailureDescription.MarkdownDescription,
								Optional:            true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description:         desktopOtpFailureCountDescription.Description,
										MarkdownDescription: desktopOtpFailureCountDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.Int32{
											int32validator.Between(pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description:         desktopOtpFailureCoolDownDescription.Description,
										MarkdownDescription: desktopOtpFailureCoolDownDescription.MarkdownDescription,
										Optional:            true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         desktopOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: desktopOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},

											"time_unit": schema.StringAttribute{
												Description:         durationTimeUnitMinsSecondsDescription.Description,
												MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
						},
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         desktopPairingDisabledDescription.Description,
						MarkdownDescription: desktopPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description:         desktopPairingKeyLifetimeDescription.Description,
						MarkdownDescription: desktopPairingKeyLifetimeDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description:         desktopPairingKeyLifetimeDurationDescription.Description,
								MarkdownDescription: desktopPairingKeyLifetimeDurationDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.Int32{
									int32validator.Any(
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`HOURS`),
												fmt.Sprintf("If `time_unit` is `HOURS`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
									),
								},
							},

							"time_unit": schema.StringAttribute{
								Description:         mobileApplicationsPairingKeyLifetimeTimeUnitDescription.Description,
								MarkdownDescription: mobileApplicationsPairingKeyLifetimeTimeUnitDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)...),
								},
							},
						},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"yubikey": schema.SingleNestedAttribute{
				Description:         yubikeyDescription.Description,
				MarkdownDescription: yubikeyDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidator.IsRequiredIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGID),
						path.MatchRoot("policy_type"),
					),
					objectvalidator.ConflictsIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGONE_MFA),
						path.MatchRoot("policy_type"),
					),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         yubikeyEnabledDescription.Description,
						MarkdownDescription: yubikeyEnabledDescription.MarkdownDescription,
						Required:            true,
					},

					"otp": schema.SingleNestedAttribute{
						Description:         yubikeyOtpDescription.Description,
						MarkdownDescription: yubikeyOtpDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(pingidDeviceOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(pingidDeviceOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),
						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description:         yubikeyOtpFailureDescription.Description,
								MarkdownDescription: yubikeyOtpFailureDescription.MarkdownDescription,
								Optional:            true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description:         yubikeyOtpFailureCountDescription.Description,
										MarkdownDescription: yubikeyOtpFailureCountDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.Int32{
											int32validator.Between(pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description:         yubikeyOtpFailureCoolDownDescription.Description,
										MarkdownDescription: yubikeyOtpFailureCoolDownDescription.MarkdownDescription,
										Optional:            true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         yubikeyOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: yubikeyOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},

											"time_unit": schema.StringAttribute{
												Description:         durationTimeUnitMinsSecondsDescription.Description,
												MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
						},
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         yubikeyPairingDisabledDescription.Description,
						MarkdownDescription: yubikeyPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description:         yubikeyPairingKeyLifetimeDescription.Description,
						MarkdownDescription: yubikeyPairingKeyLifetimeDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description:         yubikeyPairingKeyLifetimeDurationDescription.Description,
								MarkdownDescription: yubikeyPairingKeyLifetimeDurationDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.Int32{
									int32validator.Any(
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`HOURS`),
												fmt.Sprintf("If `time_unit` is `HOURS`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
									),
								},
							},

							"time_unit": schema.StringAttribute{
								Description:         mobileApplicationsPairingKeyLifetimeTimeUnitDescription.Description,
								MarkdownDescription: mobileApplicationsPairingKeyLifetimeTimeUnitDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)...),
								},
							},
						},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"oath_token": schema.SingleNestedAttribute{
				Description:         oathTokenDescription.Description,
				MarkdownDescription: oathTokenDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(oathTokenDefault),

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         oathTokenEnabledDescription.Description,
						MarkdownDescription: oathTokenEnabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"otp": schema.SingleNestedAttribute{
						Description:         oathTokenOtpDescription.Description,
						MarkdownDescription: oathTokenOtpDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(pingidDeviceOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(pingidDeviceOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description:         oathTokenOtpFailureDescription.Description,
								MarkdownDescription: oathTokenOtpFailureDescription.MarkdownDescription,
								Optional:            true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description:         oathTokenOtpFailureCountDescription.Description,
										MarkdownDescription: oathTokenOtpFailureCountDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.Int32{
											int32validator.Between(pingidDeviceOtpFailureCountMin, pingidDeviceOtpFailureCountMax),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description:         oathTokenOtpFailureCoolDownDescription.Description,
										MarkdownDescription: oathTokenOtpFailureCoolDownDescription.MarkdownDescription,
										Optional:            true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         oathTokenOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: oathTokenOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinSeconds, pingidDeviceOtpFailureCoolDownDurationMaxSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDeviceOtpFailureCoolDownDurationMinMinutes, pingidDeviceOtpFailureCoolDownDurationMaxMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},

											"time_unit": schema.StringAttribute{
												Description:         durationTimeUnitMinsSecondsDescription.Description,
												MarkdownDescription: durationTimeUnitMinsSecondsDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
						},
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         oathTokenPairingDisabledDescription.Description,
						MarkdownDescription: oathTokenPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description:         oathTokenPairingKeyLifetimeDescription.Description,
						MarkdownDescription: oathTokenPairingKeyLifetimeDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description:         oathTokenPairingKeyLifetimeDurationDescription.Description,
								MarkdownDescription: oathTokenPairingKeyLifetimeDurationDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.Int32{
									int32validator.Any(
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinMinutes, pingidDevicePairingKeyLifetimeDurationMaxMinutes),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int32validator.All(
											int32validator.Between(pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`HOURS`),
												fmt.Sprintf("If `time_unit` is `HOURS`, the allowed duration range is %d - %d.", pingidDevicePairingKeyLifetimeDurationMinHours, pingidDevicePairingKeyLifetimeDurationMaxHours),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
									),
								},
							},

							"time_unit": schema.StringAttribute{
								Description:         mobileApplicationsPairingKeyLifetimeTimeUnitDescription.Description,
								MarkdownDescription: mobileApplicationsPairingKeyLifetimeTimeUnitDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)...),
								},
							},
						},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func (r *MFADevicePolicyDefaultResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"State change warning",
			"A destroy plan has been detected for the \"pingone_mfa_device_policy_default\" resource.  The default MFA device policy will be removed from Terraform's state.  The policy itself will not be removed from the PingOne service and will retain its current configuration.",
		)
	}
}

func (r *MFADevicePolicyDefaultResource) devicePolicyOfflineDeviceSchemaAttribute(descriptionMethod string) schema.SingleNestedAttribute {

	const otpFailureCountDefault = 3
	const otpFailureCoolDownDurationDefault = 0
	const otpLifetimeDurationDefault = 30
	const otpOtpLengthDefault = 6

	const otpFailureCountMin = 1
	const otpFailureCountMax = 7
	const otpFailureCoolDownDurationMin = 0
	const otpFailureCoolDownDurationMax = 30
	const otpLifetimeDurationMin = 1
	const otpLifetimeDurationMax = 120
	const otpOtpLengthMin = 6
	const otpOtpLengthMax = 10

	devicePolicyOfflineDeviceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that allows configuration of %s device authentication policy settings.", descriptionMethod),
	)

	devicePolicyOfflineDeviceEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A boolean that specifies whether the %s method is enabled or disabled in the policy.", descriptionMethod),
	)

	pairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A boolean that, when set to `true`, prevents users from pairing new devices with the %s method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.", descriptionMethod),
	).DefaultValue(false)

	devicePolicyOfflineDeviceOtpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that allows configuration of %s settings.", descriptionMethod),
	)

	devicePolicyOfflineDeviceOtpFailureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that allows configuration of %s failure settings.", descriptionMethod),
	)

	devicePolicyOfflineDeviceOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that allows configuration of %s failure cool down settings.", descriptionMethod),
	)

	devicePolicyOfflineDeviceOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
	)

	otpCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	devicePolicyOfflineDeviceOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. The minimum value is `%d` and the maximum value is `%d`.", otpFailureCountMin, otpFailureCountMax),
	)

	devicePolicyOfflineDeviceOtpLifetimeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A single object that allows configuration of %s lifetime settings.", descriptionMethod),
	)

	devicePolicyOfflineDeviceOtpLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the duration (number of time units) that the passcode is valid before it expires. The minimum value is `%d` and the maximum value is `%d`.", otpLifetimeDurationMin, otpLifetimeDurationMax),
	)

	otpOtpLengthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the length of the OTP that is shown to users.  Minimum length is `%d` digits and maximum is `%d` digits.", otpOtpLengthMin, otpOtpLengthMax),
	).DefaultValue(otpOtpLengthDefault)

	promptForNicknameOnPairingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.",
	)

	return schema.SingleNestedAttribute{
		Description:         devicePolicyOfflineDeviceDescription.Description,
		MarkdownDescription: devicePolicyOfflineDeviceDescription.MarkdownDescription,
		Required:            true,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description:         devicePolicyOfflineDeviceEnabledDescription.Description,
				MarkdownDescription: devicePolicyOfflineDeviceEnabledDescription.MarkdownDescription,
				Required:            true,
			},

			"pairing_disabled": schema.BoolAttribute{
				Description:         pairingDisabledDescription.Description,
				MarkdownDescription: pairingDisabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"otp": schema.SingleNestedAttribute{
				Description:         devicePolicyOfflineDeviceOtpDescription.Description,
				MarkdownDescription: devicePolicyOfflineDeviceOtpDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					MFADevicePolicyOfflineDeviceOtpTFObjectTypes,
					map[string]attr.Value{
						"failure": types.ObjectValueMust(
							MFADevicePolicyFailureTFObjectTypes,
							map[string]attr.Value{
								"count": types.Int32Value(otpFailureCountDefault),
								"cool_down": types.ObjectValueMust(
									MFADevicePolicyTimePeriodTFObjectTypes,
									map[string]attr.Value{
										"duration":  types.Int32Value(otpFailureCoolDownDurationDefault),
										"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
									},
								),
							},
						),
						"lifetime": types.ObjectValueMust(
							MFADevicePolicyTimePeriodTFObjectTypes,
							map[string]attr.Value{
								"duration":  types.Int32Value(otpLifetimeDurationDefault),
								"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
							},
						),
						"otp_length": types.Int32Value(otpOtpLengthDefault),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"failure": schema.SingleNestedAttribute{
						Description:         devicePolicyOfflineDeviceOtpFailureDescription.Description,
						MarkdownDescription: devicePolicyOfflineDeviceOtpFailureDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"cool_down": schema.SingleNestedAttribute{
								Description:         devicePolicyOfflineDeviceOtpFailureCoolDownDescription.Description,
								MarkdownDescription: devicePolicyOfflineDeviceOtpFailureCoolDownDescription.MarkdownDescription,
								Required:            true,

								Attributes: map[string]schema.Attribute{
									"duration": schema.Int32Attribute{
										Description:         devicePolicyOfflineDeviceOtpFailureCoolDownDurationDescription.Description,
										MarkdownDescription: devicePolicyOfflineDeviceOtpFailureCoolDownDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int32{
											int32validator.Between(otpFailureCoolDownDurationMin, otpFailureCoolDownDurationMax),
										},
									},

									"time_unit": schema.StringAttribute{
										Description:         otpCoolDownDescription.Description,
										MarkdownDescription: otpCoolDownDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
										},
									},
								},
							},

							"count": schema.Int32Attribute{
								Description:         devicePolicyOfflineDeviceOtpFailureCountDescription.Description,
								MarkdownDescription: devicePolicyOfflineDeviceOtpFailureCountDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.Int32{
									int32validator.Between(otpFailureCountMin, otpFailureCountMax),
								},
							},
						},
					},

					"lifetime": schema.SingleNestedAttribute{
						Description:         devicePolicyOfflineDeviceOtpLifetimeDescription.Description,
						MarkdownDescription: devicePolicyOfflineDeviceOtpLifetimeDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyTimePeriodTFObjectTypes,
							map[string]attr.Value{
								"duration":  types.Int32Value(otpLifetimeDurationDefault),
								"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description:         devicePolicyOfflineDeviceOtpLifetimeDurationDescription.Description,
								MarkdownDescription: devicePolicyOfflineDeviceOtpLifetimeDurationDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.Int32{
									int32validator.Between(otpLifetimeDurationMin, otpLifetimeDurationMax),
								},
							},

							"time_unit": schema.StringAttribute{
								Description:         otpCoolDownDescription.Description,
								MarkdownDescription: otpCoolDownDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
								},
							},
						},
					},
					"otp_length": schema.Int32Attribute{
						Description:         otpOtpLengthDescription.Description,
						MarkdownDescription: otpOtpLengthDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int32default.StaticInt32(otpOtpLengthDefault),

						Validators: []validator.Int32{
							int32validator.Between(otpOtpLengthMin, otpOtpLengthMax),
						},
					},
				},
			},

			"prompt_for_nickname_on_pairing": schema.BoolAttribute{
				Description:         promptForNicknameOnPairingDescription.Description,
				MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
				Optional:            true,
			},
		},
	}
}

func (r *MFADevicePolicyDefaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MFADevicePolicyDefaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MFADevicePolicyDefaultResourceModel

	if r.Client.MFAAPIClient == nil {
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

	// Update the default policy
	state, d := r.updateMFADevicePolicyDefault(ctx, plan, true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyDefaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *MFADevicePolicyDefaultResourceModel

	if r.Client.MFAAPIClient == nil {
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

	var defaultPolicy *mfa.DeviceAuthenticationPolicy
	var d diag.Diagnostics

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		var response *mfa.DeviceAuthenticationPolicy
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,
			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneDeviceAuthenticationPolicy-Default",
			legacysdk.CustomErrorResourceNotFoundWarning,
			nil,
			&response,
		)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if response != nil {
			// Check if it is still the default policy
			var isDefault bool
			if v, ok := response.GetDefaultOk(); ok {
				isDefault = *v
			}

			if isDefault {
				defaultPolicy = response
			} else {
				resp.Diagnostics.AddWarning(
					"Managed MFA Device Policy is no longer the default",
					fmt.Sprintf("The MFA Device Policy with ID '%s' is no longer the default policy for environment '%s'. The provider will now manage the new default policy. The previous policy may remain in the environment but is no longer managed by this resource.", response.GetId(), data.EnvironmentId.ValueString()),
				)
			}
		}
	}

	if defaultPolicy == nil {
		defaultPolicy, d = FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Remove from state if resource is not found
	if defaultPolicy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Determine policy type - either from state (set by ImportState) or from API
	policyType := data.PolicyType.ValueString()
	if policyType == "" {
		policyType = determinePolicyType(defaultPolicy)
	}

	// Populate state from API response
	resp.Diagnostics.Append(data.toState(defaultPolicy, policyType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *MFADevicePolicyDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MFADevicePolicyDefaultResourceModel

	if r.Client.MFAAPIClient == nil {
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

	// Update the default policy
	state, d := r.updateMFADevicePolicyDefault(ctx, plan, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyDefaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MFADevicePolicyDefaultResourceModel

	if r.Client.MFAAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If notifications_policy is set, we must unset it to allow the referenced policy to be deleted
	if !data.NotificationsPolicy.IsNull() && !data.NotificationsPolicy.IsUnknown() {
		// Fetch the default policy to get its ID
		response, d := FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), true)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response == nil {
			// Default MFA device policy not found, nothing to do
			return
		}

		// Create a copy of data with NotificationsPolicy set to null
		data.NotificationsPolicy = types.ObjectNull(data.NotificationsPolicy.AttributeTypes(ctx))

		// Build the model for the API
		mFADevicePolicy, d := data.expand(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Extract ID from the union type based on policy_type
		// We can use the ID from the fetched response directly
		policyID := response.GetId()

		// Run the API call
		var updateResponse *mfa.DeviceAuthenticationPolicy
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,
			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), policyID).DeviceAuthenticationPolicy(mFADevicePolicy).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateDeviceAuthenticationPolicy-Default-Delete",
			legacysdk.CustomErrorResourceNotFoundWarning,
			nil,
			&updateResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (r *MFADevicePolicyDefaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
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

	// Set the required attributes in state for Read to work
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes["environment_id"])...)
}

func (r *MFADevicePolicyDefaultResource) checkBOMForPolicyType(ctx context.Context, environmentID string, policyType string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map policy type to product type
	var requiredProduct *management.EnumProductType
	switch policyType {
	case POLICY_TYPE_PINGONE_MFA:
		v := management.ENUMPRODUCTTYPE_ONE_MFA
		requiredProduct = &v
	case POLICY_TYPE_PINGID:
		v := management.EnumProductType("PING_ONE_ID")
		requiredProduct = &v
	default:
		// We don't have a clear mapping for other policy types (e.g. PingID) in the BOM product types yet
		// So we skip validation for those for now
		return diags
	}

	// Fetch BOM
	bom, _, err := r.Client.ManagementAPIClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, environmentID).Execute()
	if err != nil {
		diags.AddError(
			"Error reading Bill of Materials",
			fmt.Sprintf("Could not read Bill of Materials for environment %s to validate policy type support: %s", environmentID, err.Error()),
		)
		return diags
	}

	// Check if product exists
	found := false
	if bom != nil {
		for _, product := range bom.GetProducts() {
			if product.Type == *requiredProduct {
				found = true
				break
			}
		}
	}

	if !found {
		diags.AddError(
			"Unsupported Policy Type",
			fmt.Sprintf("The environment %s does not support the policy type '%s'. The required product '%s' is not enabled in the Bill of Materials.", environmentID, policyType, string(*requiredProduct)),
		)
	}

	return diags
}

func (r *MFADevicePolicyDefaultResource) updateMFADevicePolicyDefault(ctx context.Context, plan MFADevicePolicyDefaultResourceModel, isCreate bool) (MFADevicePolicyDefaultResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state MFADevicePolicyDefaultResourceModel

	if r.Client.MFAAPIClient == nil {
		diags.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return state, diags
	}

	// Validate policy type against BOM
	if !plan.PolicyType.IsUnknown() && !plan.PolicyType.IsNull() {
		diags.Append(r.checkBOMForPolicyType(ctx, plan.EnvironmentId.ValueString(), plan.PolicyType.ValueString())...)
		if diags.HasError() {
			return state, diags
		}
	}

	// Build the model for the API
	mFADevicePolicy, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return state, diags
	}

	// Run the API call to check if default exists
	var readResponse *mfa.DeviceAuthenticationPolicy

	if !isCreate && !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		var response *mfa.DeviceAuthenticationPolicy
		diags.Append(legacysdk.ParseResponse(
			ctx,
			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneDeviceAuthenticationPolicy-Default",
			legacysdk.CustomErrorResourceNotFoundWarning,
			nil,
			&response,
		)...)

		if diags.HasError() {
			return state, diags
		}

		if response != nil {
			// Check if it is still the default policy
			var isDefault bool
			if v, ok := response.GetDefaultOk(); ok {
				isDefault = *v
			}

			if isDefault {
				readResponse = response
			}
		}
	}

	if readResponse == nil {
		readResponse, d = FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), false)
		diags.Append(d...)
		if diags.HasError() {
			return state, diags
		}
	}

	// The API ensures a default policy always exists, so if we can't find it, something is wrong
	if readResponse == nil {
		diags.AddError(
			"Default MFA Device Policy Not Found",
			"The default MFA device policy could not be found for the environment.",
		)
		return state, diags
	}

	// Extract ID from the union type based on policy_type
	policyID := readResponse.GetId()

	// Update the default policy
	var response *mfa.DeviceAuthenticationPolicy

	diags.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, plan.EnvironmentId.ValueString(), policyID).DeviceAuthenticationPolicy(mFADevicePolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDeviceAuthenticationPolicy-Default",
		legacysdk.DefaultCustomError,
		nil,
		&response,
	)...)
	if diags.HasError() {
		return state, diags
	}

	// Populate state from API response
	diags.Append(state.toState(response, plan.PolicyType.ValueString())...)
	if diags.HasError() {
		return state, diags
	}

	// PolicyType is not returned by API, preserve it from plan
	state.PolicyType = plan.PolicyType

	return state, diags
}

// determinePolicyType determines the policy type from an API response
// by checking for PingID-specific fields
func determinePolicyType(response *mfa.DeviceAuthenticationPolicy) string {
	if response == nil {
		return POLICY_TYPE_PINGONE_MFA
	}

	// Check for PingID-specific fields
	if response.Desktop != nil || response.Yubikey != nil {
		return POLICY_TYPE_PINGID
	}

	return POLICY_TYPE_PINGONE_MFA
}

func FetchDefaultMFADevicePolicy(ctx context.Context, apiClient *mfa.APIClient, managementAPIClient *management.APIClient, environmentID string, warnOnNotFound bool) (*mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	defaultTimeout := 30 * time.Second
	return FetchDefaultMFADevicePolicyWithTimeout(ctx, apiClient, managementAPIClient, environmentID, warnOnNotFound, defaultTimeout)
}

func FetchDefaultMFADevicePolicyWithTimeout(ctx context.Context, apiClient *mfa.APIClient, managementAPIClient *management.APIClient, environmentID string, warnOnNotFound bool, timeout time.Duration) (*mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := legacysdk.DefaultCustomError
	if warnOnNotFound {
		errorFunction = legacysdk.CustomErrorResourceNotFoundWarning
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			"false",
		},
		Target: []string{
			"true",
			"err",
		},
		Refresh: func() (interface{}, string, error) {

			// Run the API call
			var defaultMFADevicePolicy *mfa.DeviceAuthenticationPolicy
			diags.Append(legacysdk.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					pagedIterator := apiClient.DeviceAuthenticationPolicyApi.ReadDeviceAuthenticationPolicies(ctx, environmentID).Execute()

					var initialHttpResponse *http.Response

					for pageCursor, err := range pagedIterator {
						if err != nil {
							return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, managementAPIClient, environmentID, nil, pageCursor.HTTPResponse, err)
						}

						if initialHttpResponse == nil {
							initialHttpResponse = pageCursor.HTTPResponse
						}

						if policies, ok := pageCursor.EntityArray.Embedded.GetDeviceAuthenticationPoliciesOk(); ok {

							for _, policyItem := range policies {
								// With flattened SDK, check default flag directly on policyItem
								var isDefault bool
								if v, ok := policyItem.GetDefaultOk(); ok {
									isDefault = *v
								}

								if isDefault {
									defaultMFADevicePolicy = &policyItem
									break
								}
							}
						}

						if defaultMFADevicePolicy != nil {
							break
						}
					}

					return nil, initialHttpResponse, nil
				},
				"ReadDeviceAuthenticationPolicies-FetchDefaultMFADevicePolicy",
				errorFunction,
				sdk.DefaultCreateReadRetryable,
				&defaultMFADevicePolicy,
			)...)
			if diags.HasError() {
				return nil, "err", fmt.Errorf("Error reading MFA device policies")
			}

			tflog.Debug(ctx, "Find default MFA device policy attempt", map[string]interface{}{
				"policy": defaultMFADevicePolicy,
				"result": strings.ToLower(strconv.FormatBool(defaultMFADevicePolicy != nil)),
			})

			return defaultMFADevicePolicy, strings.ToLower(strconv.FormatBool(defaultMFADevicePolicy != nil)), nil
		},
		Timeout:                   timeout,
		Delay:                     1 * time.Second,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	policy, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		tflog.Warn(ctx, "Cannot find default MFA device policy for the environment", map[string]interface{}{
			"environment": environmentID,
			"err":         err,
		})

		return nil, diags
	}

	returnVar := policy.(*mfa.DeviceAuthenticationPolicy)

	return returnVar, diags
}

func (p *MFADevicePolicyDefaultResourceModel) expand(ctx context.Context) (mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	// Get policy type to handle divergences
	policyType := p.PolicyType.ValueString()

	// SMS
	var smsPlan MFADevicePolicySmsResourceModel
	diags.Append(p.Sms.As(ctx, &smsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	sms, d := smsPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}

	// Voice
	var voicePlan MFADevicePolicyVoiceResourceModel
	diags.Append(p.Voice.As(ctx, &voicePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	voice, d := voicePlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}

	// Email
	var emailPlan MFADevicePolicyEmailResourceModel
	diags.Append(p.Email.As(ctx, &emailPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	email, d := emailPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}

	// Mobile
	var mobilePlan MFADevicePolicyDefaultMobileResourceModel
	diags.Append(p.Mobile.As(ctx, &mobilePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	mobile, d := expandMobileForDefault(ctx, mobilePlan, policyType)
	diags.Append(d...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}

	// TOTP
	var totpPlan MFADevicePolicyTotpResourceModel
	diags.Append(p.Totp.As(ctx, &totpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	totp, d := totpPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}

	// Main object - build policy with flattened SDK
	data := mfa.NewDeviceAuthenticationPolicy(
		p.Name.ValueString(),
		*sms,
		*voice,
		*email,
		*mobile,
		*totp,
		false, // default
		false, // forSignOnPolicy
	)

	// Always set default to true for the default policy
	data.SetDefault(true)

	// FIDO2 - available for both policy types
	if !p.Fido2.IsNull() && !p.Fido2.IsUnknown() {
		var fido2Plan MFADevicePolicyFido2ResourceModel
		diags.Append(p.Fido2.As(ctx, &fido2Plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		fido2 := fido2Plan.expand()
		data.SetFido2(*fido2)
	}

	// Desktop - only for PingID
	if policyType == POLICY_TYPE_PINGID {
		if !p.Desktop.IsNull() && !p.Desktop.IsUnknown() {
			var desktopPlan MFADevicePolicyDesktopResourceModel
			diags.Append(p.Desktop.As(ctx, &desktopPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			desktop, d := desktopPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			data.SetDesktop(*desktop)
		}

		// Yubikey - only for PingID
		if !p.Yubikey.IsNull() && !p.Yubikey.IsUnknown() {
			var yubikeyPlan MFADevicePolicyYubikeyResourceModel
			diags.Append(p.Yubikey.As(ctx, &yubikeyPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			yubikey, d := yubikeyPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			data.SetYubikey(*yubikey)
		}
	}

	// Authentication - both policy types
	if !p.Authentication.IsNull() && !p.Authentication.IsUnknown() {
		var authenticationPlan MFADevicePolicyAuthenticationResourceModel
		diags.Append(p.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		data.SetAuthentication(
			*mfa.NewDeviceAuthenticationPolicyCommonAuthentication(
				mfa.EnumMFADevicePolicySelection(authenticationPlan.DeviceSelection.ValueString()),
			),
		)
	}

	// New Device Notification - both policy types
	if !p.NewDeviceNotification.IsNull() && !p.NewDeviceNotification.IsUnknown() {
		data.SetNewDeviceNotification(
			mfa.EnumMFADevicePolicyNewDeviceNotification(p.NewDeviceNotification.ValueString()),
		)
	}

	// Ignore User Lock - both policy types
	if !p.IgnoreUserLock.IsNull() && !p.IgnoreUserLock.IsUnknown() {
		data.SetIgnoreUserLock(p.IgnoreUserLock.ValueBool())
	}

	// NotificationsPolicy - both policy types
	if !p.NotificationsPolicy.IsNull() && !p.NotificationsPolicy.IsUnknown() {
		var notificationsPolicyPlan struct {
			Id types.String `tfsdk:"id"`
		}
		diags.Append(p.NotificationsPolicy.As(ctx, &notificationsPolicyPlan, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		data.SetNotificationsPolicy(
			*mfa.NewDeviceAuthenticationPolicyCommonNotificationsPolicy(notificationsPolicyPlan.Id.ValueString()),
		)
	}

	// RememberMe - both policy types
	if !p.RememberMe.IsNull() && !p.RememberMe.IsUnknown() {
		var rememberMePlan struct {
			Web types.Object `tfsdk:"web"`
		}
		diags.Append(p.RememberMe.As(ctx, &rememberMePlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		if !rememberMePlan.Web.IsNull() && !rememberMePlan.Web.IsUnknown() {
			var webPlan struct {
				Enabled  types.Bool   `tfsdk:"enabled"`
				LifeTime types.Object `tfsdk:"life_time"`
			}
			diags.Append(rememberMePlan.Web.As(ctx, &webPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			// Parse LifeTime
			var lifeTimePlan struct {
				Duration types.Int32  `tfsdk:"duration"`
				TimeUnit types.String `tfsdk:"time_unit"`
			}
			diags.Append(webPlan.LifeTime.As(ctx, &lifeTimePlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return mfa.DeviceAuthenticationPolicy{}, diags
			}

			lifeTime := mfa.DeviceAuthenticationPolicyCommonRememberMeWebLifeTime{}
			lifeTime.SetDuration(lifeTimePlan.Duration.ValueInt32())
			lifeTime.SetTimeUnit(mfa.EnumTimeUnitRememberMeWebLifeTime(lifeTimePlan.TimeUnit.ValueString()))

			web := mfa.NewDeviceAuthenticationPolicyCommonRememberMeWeb(
				webPlan.Enabled.ValueBool(),
				lifeTime,
			)

			rememberMe := mfa.NewDeviceAuthenticationPolicyCommonRememberMe(*web)
			data.SetRememberMe(*rememberMe)
		}
	}

	// OathToken - both policy types
	if !p.OathToken.IsNull() && !p.OathToken.IsUnknown() {
		var oathTokenPlan MFADevicePolicyOathTokenResourceModel
		diags.Append(p.OathToken.As(ctx, &oathTokenPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		oathToken, d := oathTokenPlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return mfa.DeviceAuthenticationPolicy{}, diags
		}

		data.SetOathToken(*oathToken)
	} else {
		tflog.Debug(ctx, "oath_token is null or unknown, NOT sending to API")
	}

	return *data, diags
}

func (p *MFADevicePolicyDesktopResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyPingIDDevice, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	// OTP
	var otpPlan MFADevicePolicyPingIDDeviceOtpResourceModel
	diags.Append(p.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	otp, d := otpPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := mfa.NewDeviceAuthenticationPolicyPingIDDevice(
		p.Enabled.ValueBool(),
		*otp,
	)

	// Pairing Disabled
	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	// Pairing Key Lifetime
	if !p.PairingKeyLifetime.IsNull() && !p.PairingKeyLifetime.IsUnknown() {
		var pairingKeyLifetimePlan MFADevicePolicyTimePeriodResourceModel
		diags.Append(p.PairingKeyLifetime.As(ctx, &pairingKeyLifetimePlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPairingKeyLifetime(
			mfa.DeviceAuthenticationPolicyPingIDDevicePairingKeyLifetime{
				Duration: pairingKeyLifetimePlan.Duration.ValueInt32(),
				TimeUnit: mfa.EnumTimeUnitPairingKeyLifetime(pairingKeyLifetimePlan.TimeUnit.ValueString()),
			},
		)
	}

	// Prompt for Nickname on Pairing
	if !p.PromptForNicknameOnPairing.IsNull() && !p.PromptForNicknameOnPairing.IsUnknown() {
		data.SetPromptForNicknameOnPairing(p.PromptForNicknameOnPairing.ValueBool())
	}

	return data, diags
}

// Yubikey uses the same expand logic as Desktop
func (p *MFADevicePolicyYubikeyResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyPingIDDevice, diag.Diagnostics) {
	return (*MFADevicePolicyDesktopResourceModel)(p).expand(ctx)
}

// OathToken has similar structure but different SDK type
func (p *MFADevicePolicyOathTokenResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOathToken, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	// OTP
	var otpPlan MFADevicePolicyPingIDDeviceOtpResourceModel
	diags.Append(p.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	otp, d := otpPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := mfa.NewDeviceAuthenticationPolicyOathToken(
		p.Enabled.ValueBool(),
		*otp,
	)

	// Pairing Disabled
	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	// Pairing Key Lifetime
	if !p.PairingKeyLifetime.IsNull() && !p.PairingKeyLifetime.IsUnknown() {
		var pairingKeyLifetimePlan MFADevicePolicyTimePeriodResourceModel
		diags.Append(p.PairingKeyLifetime.As(ctx, &pairingKeyLifetimePlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPairingKeyLifetime(
			mfa.DeviceAuthenticationPolicyOathTokenPairingKeyLifetime{
				Duration: pairingKeyLifetimePlan.Duration.ValueInt32(),
				TimeUnit: mfa.EnumTimeUnit(pairingKeyLifetimePlan.TimeUnit.ValueString()),
			},
		)
	}

	// Prompt for Nickname on Pairing
	if !p.PromptForNicknameOnPairing.IsNull() && !p.PromptForNicknameOnPairing.IsUnknown() {
		data.SetPromptForNicknameOnPairing(p.PromptForNicknameOnPairing.ValueBool())
	}

	return data, diags
}

// PingID Device OTP model (for otp field within desktop/yubikey/oathToken)
type MFADevicePolicyPingIDDeviceOtpResourceModel struct {
	Failure types.Object `tfsdk:"failure"`
}

func (p *MFADevicePolicyPingIDDeviceOtpResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyPingIDDeviceOtp, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	var failurePlan MFADevicePolicyFailureResourceModel
	diags.Append(p.Failure.As(ctx, &failurePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	failure, d := failurePlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	otp := mfa.NewDeviceAuthenticationPolicyPingIDDeviceOtp(
		*failure,
	)

	return otp, diags
}

func expandMobileForDefault(ctx context.Context, mobilePlan MFADevicePolicyDefaultMobileResourceModel, policyType string) (*mfa.DeviceAuthenticationPolicyCommonMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Handle OTP (required)
	var otpPlan MFADevicePolicyOtpResourceModel
	diags.Append(mobilePlan.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	var otpFailurePlan MFADevicePolicyFailureResourceModel
	diags.Append(otpPlan.Failure.As(ctx, &otpFailurePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	failure, d := otpFailurePlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	otp := mfa.NewDeviceAuthenticationPolicyCommonMobileOtp(*failure)

	// Build the mobile object
	mobile := mfa.NewDeviceAuthenticationPolicyCommonMobile(mobilePlan.Enabled.ValueBool(), *otp)

	// Handle prompt for nickname
	if !mobilePlan.PromptForNicknameOnPairing.IsNull() && !mobilePlan.PromptForNicknameOnPairing.IsUnknown() {
		mobile.SetPromptForNicknameOnPairing(mobilePlan.PromptForNicknameOnPairing.ValueBool())
	}

	if !mobilePlan.Applications.IsNull() && !mobilePlan.Applications.IsUnknown() {
		planApps := make(map[string]MFADevicePolicyDefaultMobileApplicationResourceModel)
		diags.Append(mobilePlan.Applications.ElementsAs(ctx, &planApps, false)...)
		if diags.HasError() {
			return nil, diags
		}

		applications := make([]mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInner, 0, len(planApps))
		for appId, appPlan := range planApps {

			app := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInner(appId)

			if !appPlan.AutoEnrolment.IsNull() && !appPlan.AutoEnrolment.IsUnknown() {
				var autoEnrolPlan MFADevicePolicyMobileApplicationAutoEnrolmentResourceModel
				diags.Append(appPlan.AutoEnrolment.As(ctx, &autoEnrolPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					autoEnrol := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerAutoEnrollment(autoEnrolPlan.Enabled.ValueBool())
					app.SetAutoEnrollment(*autoEnrol)
				}
			}

			if !appPlan.BiometricsEnabled.IsNull() && !appPlan.BiometricsEnabled.IsUnknown() {
				app.SetBiometricsEnabled(appPlan.BiometricsEnabled.ValueBool())
			}

			if !appPlan.DeviceAuthorization.IsNull() && !appPlan.DeviceAuthorization.IsUnknown() {
				var devAuthPlan MFADevicePolicyMobileApplicationDeviceAuthorizationResourceModel
				diags.Append(appPlan.DeviceAuthorization.As(ctx, &devAuthPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					devAuth := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerDeviceAuthorization(devAuthPlan.Enabled.ValueBool())
					if !devAuthPlan.ExtraVerification.IsNull() && !devAuthPlan.ExtraVerification.IsUnknown() {
						devAuth.SetExtraVerification(mfa.EnumMFADevicePolicyMobileExtraVerification(devAuthPlan.ExtraVerification.ValueString()))
					}
					app.SetDeviceAuthorization(*devAuth)
				}
			}

			if !appPlan.IntegrityDetection.IsNull() && !appPlan.IntegrityDetection.IsUnknown() {
				app.SetIntegrityDetection(mfa.EnumMFADevicePolicyMobileIntegrityDetection(appPlan.IntegrityDetection.ValueString()))
			}

			// For PingID policies, we must send the type to get PingID-specific fields back from the API
			tflog.Debug(ctx, "DEBUG: Application type value in expand", map[string]interface{}{
				"app_id":          appId,
				"policy_type":     policyType,
				"type_value":      appPlan.Type.ValueString(),
				"type_is_null":    appPlan.Type.IsNull(),
				"type_is_unknown": appPlan.Type.IsUnknown(),
			})
			// Only send type for PingID policies - the API will reject it for PingOne MFA policies
			if policyType == POLICY_TYPE_PINGID && !appPlan.Type.IsNull() && !appPlan.Type.IsUnknown() {
				app.SetType(mfa.EnumPingIDApplicationType(appPlan.Type.ValueString()))
			}

			if !appPlan.PairingDisabled.IsNull() && !appPlan.PairingDisabled.IsUnknown() {
				app.SetPairingDisabled(appPlan.PairingDisabled.ValueBool())
			}

			if !appPlan.Otp.IsNull() && !appPlan.Otp.IsUnknown() {
				var otpPlan MFADevicePolicyDefaultMobileApplicationOtpResourceModel
				diags.Append(appPlan.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					otp := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerOtp(otpPlan.Enabled.ValueBool())
					app.SetOtp(*otp)
				}
			}

			if !appPlan.PairingKeyLifetime.IsNull() && !appPlan.PairingKeyLifetime.IsUnknown() {
				var pklPlan MFADevicePolicyTimePeriodResourceModel
				diags.Append(appPlan.PairingKeyLifetime.As(ctx, &pklPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					pairingKeyLifetime := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPairingKeyLifetime(
						pklPlan.Duration.ValueInt32(),
						mfa.EnumTimeUnitPairingKeyLifetime(pklPlan.TimeUnit.ValueString()),
					)
					app.SetPairingKeyLifetime(*pairingKeyLifetime)
				}
			}

			if !appPlan.Push.IsNull() && !appPlan.Push.IsUnknown() {
				var pushPlan MFADevicePolicyDefaultMobileApplicationPushResourceModel
				diags.Append(appPlan.Push.As(ctx, &pushPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					push := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPush(pushPlan.Enabled.ValueBool())

					// Expand number_matching
					if !pushPlan.NumberMatching.IsNull() && !pushPlan.NumberMatching.IsUnknown() {
						var nmPlan MFADevicePolicyDefaultMobileApplicationPushNumberMatchingResourceModel
						diags.Append(pushPlan.NumberMatching.As(ctx, &nmPlan, basetypes.ObjectAsOptions{
							UnhandledNullAsEmpty:    false,
							UnhandledUnknownAsEmpty: false,
						})...)
						if !diags.HasError() {
							numberMatching := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPushNumberMatching(nmPlan.Enabled.ValueBool())
							push.SetNumberMatching(*numberMatching)
						}
					}

					app.SetPush(*push)
				}
			}

			if !appPlan.PushLimit.IsNull() && !appPlan.PushLimit.IsUnknown() {
				var pushLimitPlan MFADevicePolicyMobileApplicationPushLimitResourceModel
				diags.Append(appPlan.PushLimit.As(ctx, &pushLimitPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					pushLimit := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimit()

					if !pushLimitPlan.Count.IsNull() && !pushLimitPlan.Count.IsUnknown() {
						pushLimit.SetCount(pushLimitPlan.Count.ValueInt32())
					}

					if !pushLimitPlan.LockDuration.IsNull() && !pushLimitPlan.LockDuration.IsUnknown() {
						var lockDurationPlan MFADevicePolicyTimePeriodResourceModel
						diags.Append(pushLimitPlan.LockDuration.As(ctx, &lockDurationPlan, basetypes.ObjectAsOptions{
							UnhandledNullAsEmpty:    false,
							UnhandledUnknownAsEmpty: false,
						})...)
						if !diags.HasError() {
							lockDuration := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimitLockDuration(
								lockDurationPlan.Duration.ValueInt32(),
								mfa.EnumTimeUnit(lockDurationPlan.TimeUnit.ValueString()),
							)
							pushLimit.SetLockDuration(*lockDuration)
						}
					}

					if !pushLimitPlan.TimePeriod.IsNull() && !pushLimitPlan.TimePeriod.IsUnknown() {
						var timePeriodPlan MFADevicePolicyTimePeriodResourceModel
						diags.Append(pushLimitPlan.TimePeriod.As(ctx, &timePeriodPlan, basetypes.ObjectAsOptions{
							UnhandledNullAsEmpty:    false,
							UnhandledUnknownAsEmpty: false,
						})...)
						if !diags.HasError() {
							timePeriod := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimitTimePeriod(
								timePeriodPlan.Duration.ValueInt32(),
								mfa.EnumTimeUnit(timePeriodPlan.TimeUnit.ValueString()),
							)
							pushLimit.SetTimePeriod(*timePeriod)
						}
					}

					app.SetPushLimit(*pushLimit)
				}
			}

			if !appPlan.PushTimeout.IsNull() && !appPlan.PushTimeout.IsUnknown() {
				var pushTimeoutPlan MFADevicePolicyTimePeriodResourceModel
				diags.Append(appPlan.PushTimeout.As(ctx, &pushTimeoutPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					pushTimeout := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerPushTimeout(
						pushTimeoutPlan.Duration.ValueInt32(),
						mfa.EnumTimeUnitPushTimeout(pushTimeoutPlan.TimeUnit.ValueString()),
					)
					app.SetPushTimeout(*pushTimeout)
				}
			}

			if !appPlan.NewRequestDurationConfiguration.IsNull() && !appPlan.NewRequestDurationConfiguration.IsUnknown() {
				var nrdcPlan MFADevicePolicyMobileApplicationNewRequestDurationConfigurationResourceModel
				diags.Append(appPlan.NewRequestDurationConfiguration.As(ctx, &nrdcPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})...)
				if !diags.HasError() {
					var deviceTimeoutPlan, totalTimeoutPlan MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutResourceModel

					diags.Append(nrdcPlan.DeviceTimeout.As(ctx, &deviceTimeoutPlan, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    false,
						UnhandledUnknownAsEmpty: false,
					})...)

					diags.Append(nrdcPlan.TotalTimeout.As(ctx, &totalTimeoutPlan, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    false,
						UnhandledUnknownAsEmpty: false,
					})...)

					if !diags.HasError() {
						deviceTimeout := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfigurationDeviceTimeout(
							deviceTimeoutPlan.Duration.ValueInt32(),
							mfa.EnumTimeUnitSeconds(deviceTimeoutPlan.TimeUnit.ValueString()),
						)
						totalTimeout := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfigurationTotalTimeout(
							totalTimeoutPlan.Duration.ValueInt32(),
							mfa.EnumTimeUnitSeconds(totalTimeoutPlan.TimeUnit.ValueString()),
						)
						newRequestDurationConfig := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfiguration(*deviceTimeout, *totalTimeout)
						app.SetNewRequestDurationConfiguration(*newRequestDurationConfig)
					}
				}
			}

			if !appPlan.IpPairingConfiguration.IsNull() && !appPlan.IpPairingConfiguration.IsUnknown() {
				ipConfig := mfa.NewDeviceAuthenticationPolicyCommonMobileApplicationsInnerIpPairingConfiguration()

				ipConfigAttrs := appPlan.IpPairingConfiguration.Attributes()

				if anyIPAttr, exists := ipConfigAttrs["any_ip_address"]; exists {
					if anyIPVal, ok := anyIPAttr.(types.Bool); ok && !anyIPVal.IsNull() && !anyIPVal.IsUnknown() {
						ipConfig.SetAnyIPAdress(anyIPVal.ValueBool())
					}
				}

				if ipListAttr, exists := ipConfigAttrs["only_these_ip_addresses"]; exists {
					if ipListVal, ok := ipListAttr.(types.Set); ok && !ipListVal.IsNull() && !ipListVal.IsUnknown() {
						var ipAddresses []string
						diags.Append(ipListVal.ElementsAs(ctx, &ipAddresses, false)...)
						if !diags.HasError() && len(ipAddresses) > 0 {
							ipConfig.SetOnlyTheseIpAddresses(ipAddresses)
						}
					}
				}

				app.SetIpPairingConfiguration(*ipConfig)
			}

			applications = append(applications, *app)
		}

		mobile.SetApplications(applications)
	}

	return mobile, diags
}

func (p *MFADevicePolicyDefaultResourceModel) toState(apiObject *mfa.DeviceAuthenticationPolicy, policyType string) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	// Use provided policy type instead of detecting it
	// This ensures we respect the user's configuration and don't cause inconsistencies
	isPingID := (policyType == POLICY_TYPE_PINGID)

	p.PolicyType = types.StringValue(policyType)

	// Common fields for both policy types
	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	// Common fields
	p.Authentication, d = toStateMfaDevicePolicyAuthentication(apiObject.GetAuthenticationOk())
	diags.Append(d...)

	p.NewDeviceNotification = framework.EnumOkToTF(apiObject.GetNewDeviceNotificationOk())

	p.IgnoreUserLock = framework.BoolOkToTF(apiObject.GetIgnoreUserLockOk())

	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	p.NotificationsPolicy, d = toStateMfaDevicePolicyNotificationsPolicy(apiObject.GetNotificationsPolicyOk())
	diags.Append(d...)

	p.RememberMe, d = toStateMfaDevicePolicyRememberMe(apiObject.GetRememberMeOk())
	diags.Append(d...)

	p.Sms, d = toStateMfaDevicePolicySms(apiObject.GetSmsOk())
	diags.Append(d...)

	p.Voice, d = toStateMfaDevicePolicyVoice(apiObject.GetVoiceOk())
	diags.Append(d...)

	p.Email, d = toStateMfaDevicePolicyEmail(apiObject.GetEmailOk())
	diags.Append(d...)

	mobileApiObj, mobileOk := apiObject.GetMobileOk()
	p.Mobile, d = toStateMfaDevicePolicyMobileForDefault(mobileApiObj, mobileOk, policyType)
	diags.Append(d...)

	p.Totp, d = toStateMfaDevicePolicyTotp(apiObject.GetTotpOk())
	diags.Append(d...)

	p.OathToken, d = toStateMfaDevicePolicyOathToken(apiObject.GetOathTokenOk())
	diags.Append(d...)

	p.Fido2, d = toStateMfaDevicePolicyFido2(apiObject.GetFido2Ok())
	diags.Append(d...)

	// Policy type specific fields
	if isPingID {
		// PingID-specific devices
		p.Desktop, d = toStateMfaDevicePolicyPingIDDevice(apiObject.GetDesktopOk())
		diags.Append(d...)

		p.Yubikey, d = toStateMfaDevicePolicyPingIDDevice(apiObject.GetYubikeyOk())
		diags.Append(d...)
	} else {
		// Set PingID-specific fields to null for PingOneMFA policies
		p.Desktop = types.ObjectNull(MFADevicePolicyPingIDDeviceTFObjectTypes)
		p.Yubikey = types.ObjectNull(MFADevicePolicyPingIDDeviceTFObjectTypes)
	}

	return diags
}

func toStateMfaDevicePolicyPingIDDevice(apiObject *mfa.DeviceAuthenticationPolicyPingIDDevice, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: MFADevicePolicyPingIDDeviceTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ObjectNull(tfObjType.AttrTypes), diags
	}

	// OTP
	otp, d := toStateMfaDevicePolicyPingIDDeviceOtp(apiObject.GetOtpOk())
	diags.Append(d...)

	// Pairing Key Lifetime
	var pairingKeyLifetime types.Object
	if pkl, ok := apiObject.GetPairingKeyLifetimeOk(); ok && pkl != nil {
		pairingKeyLifetime, d = types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, map[string]attr.Value{
			"duration":  framework.Int32OkToTF(pkl.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(pkl.GetTimeUnitOk()),
		})
		diags.Append(d...)
	} else {
		pairingKeyLifetime = types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes)
	}

	objValue, d := types.ObjectValue(MFADevicePolicyPingIDDeviceTFObjectTypes, map[string]attr.Value{
		"enabled":                        framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":                            otp,
		"pairing_disabled":               framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
		"pairing_key_lifetime":           pairingKeyLifetime,
		"prompt_for_nickname_on_pairing": framework.BoolOkToTF(apiObject.GetPromptForNicknameOnPairingOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyPingIDDeviceOtp(apiObject *mfa.DeviceAuthenticationPolicyPingIDDeviceOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags, d diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: MFADevicePolicyPingIDDeviceOtpTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ObjectNull(tfObjType.AttrTypes), diags
	}

	var failure types.Object
	if f, ok := apiObject.GetFailureOk(); ok && f != nil {
		var coolDown types.Object
		if cd, cdOk := f.GetCoolDownOk(); cdOk && cd != nil {
			coolDown, d = types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, map[string]attr.Value{
				"duration":  framework.Int32OkToTF(cd.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(cd.GetTimeUnitOk()),
			})
			diags.Append(d...)
		} else {
			coolDown = types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes)
		}
		failure, d = types.ObjectValue(MFADevicePolicyFailureTFObjectTypes, map[string]attr.Value{
			"count":     framework.Int32OkToTF(f.GetCountOk()),
			"cool_down": coolDown,
		})
		diags.Append(d...)
	} else {
		failure = types.ObjectNull(MFADevicePolicyFailureTFObjectTypes)
	}

	objValue, d := types.ObjectValue(MFADevicePolicyPingIDDeviceOtpTFObjectTypes, map[string]attr.Value{
		"failure": failure,
	})
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyOathToken(apiObject *mfa.DeviceAuthenticationPolicyOathToken, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: MFADevicePolicyPingIDDeviceTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ObjectNull(tfObjType.AttrTypes), diags
	}

	// OTP
	otp, d := toStateMfaDevicePolicyPingIDDeviceOtp(apiObject.GetOtpOk())
	diags.Append(d...)

	// Pairing Key Lifetime
	var pairingKeyLifetime types.Object
	if pairingKeyLifetimeAPI, ok := apiObject.GetPairingKeyLifetimeOk(); ok && pairingKeyLifetimeAPI != nil {
		// Convert the PingID-specific pairing key lifetime type
		pairingKeyLifetime, d = types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, map[string]attr.Value{
			"duration":  framework.Int32OkToTF(pairingKeyLifetimeAPI.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(pairingKeyLifetimeAPI.GetTimeUnitOk()),
		})
		diags.Append(d...)
	} else {
		pairingKeyLifetime = types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes)
	}

	objValue, d := types.ObjectValue(MFADevicePolicyPingIDDeviceTFObjectTypes, map[string]attr.Value{
		"enabled":                        framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":                            otp,
		"pairing_disabled":               framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
		"pairing_key_lifetime":           pairingKeyLifetime,
		"prompt_for_nickname_on_pairing": framework.BoolOkToTF(apiObject.GetPromptForNicknameOnPairingOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

// Mobile toState functions for default resource with number_matching support
func toStateMfaDevicePolicyMobileForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobile, ok bool, policyType string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyDefaultMobileTFObjectTypes), nil
	}

	appsApiObj, appsOk := apiObject.GetApplicationsOk()
	applications, d := toStateMfaDevicePolicyMobileApplicationsForDefault(appsApiObj, appsOk, policyType)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyDefaultMobileTFObjectTypes), diags
	}

	otp, d := toStateMfaDevicePolicyMobileOtpForDefault(apiObject.GetOtpOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyDefaultMobileTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"applications":                   applications,
		"enabled":                        framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":                            otp,
		"prompt_for_nickname_on_pairing": framework.BoolOkToTF(apiObject.GetPromptForNicknameOnPairingOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyDefaultMobileTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

// toStateMfaDevicePolicyMobileApplicationsForDefault converts API mobile applications to Terraform state.
// policyType is used to conditionally ignore fields that conflict with the policy type.
func toStateMfaDevicePolicyMobileApplicationsForDefault(apiObject []mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInner, ok bool, policyType string) (types.Map, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: MFADevicePolicyDefaultMobileApplicationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.MapNull(tfObjType), nil
	}

	isPingID := (policyType == POLICY_TYPE_PINGID)

	objectMap := make(map[string]attr.Value)
	for _, application := range apiObject {
		// Debug: log what the API returned for each application
		biometricsVal, biometricsOk := application.GetBiometricsEnabledOk()
		typeVal, typeOk := application.GetTypeOk()
		nrdcVal, nrdcOk := application.GetNewRequestDurationConfigurationOk()
		tflog.Debug(context.Background(), "Mobile application from API", map[string]interface{}{
			"id":                                     application.GetId(),
			"biometrics_enabled_value":               biometricsVal,
			"biometrics_enabled_ok":                  biometricsOk,
			"type_value":                             typeVal,
			"type_ok":                                typeOk,
			"new_request_duration_configuration_ok":  nrdcOk,
			"new_request_duration_configuration_val": nrdcVal,
		})

		// For PingID policies, auto_enrollment and device_authorization conflict - keep them null
		var autoEnrolment types.Object
		var deviceAuthorization types.Object
		if isPingID {
			autoEnrolment = types.ObjectNull(MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes)
			deviceAuthorization = types.ObjectNull(MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes)
		} else {
			autoEnrolment, d = toStateMfaDevicePolicyMobileApplicationsAutoEnrolmentForDefault(application.GetAutoEnrollmentOk())
			diags.Append(d...)
			if diags.HasError() {
				return types.MapNull(tfObjType), diags
			}

			deviceAuthorization, d = toStateMfaDevicePolicyMobileApplicationsDeviceAuthorizationForDefault(application.GetDeviceAuthorizationOk())
			diags.Append(d...)
			if diags.HasError() {
				return types.MapNull(tfObjType), diags
			}
		}

		otp, d := toStateMfaDevicePolicyMobileApplicationsOtpForDefault(application.GetOtpOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		pairingKeyLifetime, d := toStateMfaDevicePolicyMobileApplicationsPairingKeyLifetimeForDefault(application.GetPairingKeyLifetimeOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		push, d := toStateMfaDevicePolicyMobileApplicationsPushForDefault(application.GetPushOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		pushLimit, d := toStateMfaDevicePolicyMobileApplicationsPushLimitForDefault(application.GetPushLimitOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		newRequestDurationConfiguration, d := toStateMfaDevicePolicyMobileApplicationsNewRequestDurationConfigurationForDefault(application.GetNewRequestDurationConfigurationOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		// Handle policy-type specific fields
		// For PingID: biometrics_enabled, type, new_request_duration_configuration, ip_pairing_configuration are valid
		// For PingOne MFA: auto_enrollment, device_authorization, push_timeout are valid
		var biometricsEnabled types.Bool
		var typeAttr types.String
		var pushTimeout types.Object
		var ipPairingConfiguration types.Object

		if isPingID {
			// For PingID policies, the API may not return biometrics_enabled, type, or new_request_duration_configuration
			// Preserve these values from the prior state if API didn't return them
			biometricsEnabled = framework.BoolOkToTF(application.GetBiometricsEnabledOk())
			typeAttr = framework.EnumOkToTF(application.GetTypeOk())

			// push_timeout conflicts with PingID - keep it null
			pushTimeout = types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes)

			// Handle ip_pairing_configuration from API response (PingID only)
			ipPairingConfiguration = types.ObjectNull(MFADevicePolicyMobileIpPairingConfigurationTFObjectTypes)
			if ipPairingConfigAPI, ipOk := application.GetIpPairingConfigurationOk(); ipOk && ipPairingConfigAPI != nil {
				var ipAddressesSet types.Set
				if ipAddresses, addrOk := ipPairingConfigAPI.GetOnlyTheseIpAddressesOk(); addrOk && ipAddresses != nil && len(ipAddresses) > 0 {
					ipElements := make([]attr.Value, len(ipAddresses))
					for i, ip := range ipAddresses {
						ipElements[i] = types.StringValue(ip)
					}
					ipAddressesSet, d = types.SetValue(types.StringType, ipElements)
					diags.Append(d...)
				} else {
					ipAddressesSet = types.SetNull(types.StringType)
				}

				ipPairingConfigMap := map[string]attr.Value{
					"any_ip_address":          framework.BoolOkToTF(ipPairingConfigAPI.GetAnyIPAdressOk()),
					"only_these_ip_addresses": ipAddressesSet,
				}
				ipPairingConfiguration, d = types.ObjectValue(MFADevicePolicyMobileIpPairingConfigurationTFObjectTypes, ipPairingConfigMap)
				diags.Append(d...)
			}
		} else {
			// For PingOne MFA policies, PingID-specific fields should be null
			biometricsEnabled = types.BoolNull()
			typeAttr = types.StringNull()
			newRequestDurationConfiguration = types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes)
			ipPairingConfiguration = types.ObjectNull(MFADevicePolicyMobileIpPairingConfigurationTFObjectTypes)

			// push_timeout is valid for PingOne MFA
			pushTimeout, d = toStateMfaDevicePolicyMobileApplicationsPushTimeoutForDefault(application.GetPushTimeoutOk())
			diags.Append(d...)
			if diags.HasError() {
				return types.MapNull(tfObjType), diags
			}
		}

		o := map[string]attr.Value{
			"auto_enrollment":                    autoEnrolment,
			"biometrics_enabled":                 biometricsEnabled,
			"device_authorization":               deviceAuthorization,
			"integrity_detection":                framework.EnumOkToTF(application.GetIntegrityDetectionOk()),
			"ip_pairing_configuration":           ipPairingConfiguration,
			"otp":                                otp,
			"pairing_disabled":                   framework.BoolOkToTF(application.GetPairingDisabledOk()),
			"pairing_key_lifetime":               pairingKeyLifetime,
			"push":                               push,
			"push_limit":                         pushLimit,
			"push_timeout":                       pushTimeout,
			"new_request_duration_configuration": newRequestDurationConfiguration,
			"type":                               typeAttr,
		}

		objValue, d := types.ObjectValue(MFADevicePolicyDefaultMobileApplicationTFObjectTypes, o)
		diags.Append(d...)

		objectMap[application.GetId()] = objValue
	}

	returnVar, d := types.MapValue(tfObjType, objectMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateMfaDevicePolicyMobileApplicationsAutoEnrolmentForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerAutoEnrollment, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsDeviceAuthorizationForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerDeviceAuthorization, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled":            framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"extra_verification": framework.EnumOkToTF(apiObject.GetExtraVerificationOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsOtpForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationOtpTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationOtpTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPush, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushTFObjectTypes), nil
	}

	// Handle number_matching
	numberMatching := types.ObjectNull(MFADevicePolicyDefaultMobileApplicationPushNumberMatchingTFObjectTypes)
	if nm, nmOk := apiObject.GetNumberMatchingOk(); nmOk && nm != nil {
		nmMap := map[string]attr.Value{
			"enabled": framework.BoolOkToTF(nm.GetEnabledOk()),
		}
		var d diag.Diagnostics
		numberMatching, d = types.ObjectValue(MFADevicePolicyDefaultMobileApplicationPushNumberMatchingTFObjectTypes, nmMap)
		diags.Append(d...)
	}

	o := map[string]attr.Value{
		"enabled":         framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"number_matching": numberMatching,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyDefaultMobileApplicationPushTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushLimitForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimit, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes), nil
	}

	lockDuration, d := toStateMfaDevicePolicyMobileApplicationsPushLimitLockDurationForDefault(apiObject.GetLockDurationOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes), diags
	}

	timePeriod, d := toStateMfaDevicePolicyMobileApplicationsPushLimitTimePeriodForDefault(apiObject.GetTimePeriodOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"count":         framework.Int32OkToTF(apiObject.GetCountOk()),
		"lock_duration": lockDuration,
		"time_period":   timePeriod,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushLimitLockDurationForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimitLockDuration, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushLimitTimePeriodForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPushLimitTimePeriod, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPairingKeyLifetimeForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPairingKeyLifetime, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushTimeoutForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerPushTimeout, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsNewRequestDurationConfigurationForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfiguration, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes), nil
	}

	deviceTimeout, d := toStateMfaDevicePolicyMobileApplicationsNewRequestDurationConfigurationTimeoutForDefault(&apiObject.DeviceTimeout)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes), diags
	}

	totalTimeout, d := toStateMfaDevicePolicyMobileApplicationsNewRequestDurationConfigurationTimeoutForDefault(&apiObject.TotalTimeout)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"device_timeout": deviceTimeout,
		"total_timeout":  totalTimeout,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsNewRequestDurationConfigurationTimeoutForDefault(apiObject interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes), nil
	}

	var duration *int32
	var durationOk bool
	var timeUnit *mfa.EnumTimeUnitSeconds
	var timeUnitOk bool

	// Handle both DeviceTimeout and TotalTimeout types
	switch v := apiObject.(type) {
	case *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfigurationDeviceTimeout:
		duration, durationOk = v.GetDurationOk()
		timeUnit, timeUnitOk = v.GetTimeUnitOk()
	case *mfa.DeviceAuthenticationPolicyCommonMobileApplicationsInnerNewRequestDurationConfigurationTotalTimeout:
		duration, durationOk = v.GetDurationOk()
		timeUnit, timeUnitOk = v.GetTimeUnitOk()
	default:
		return types.ObjectNull(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(duration, durationOk),
		"time_unit": framework.EnumOkToTF(timeUnit, timeUnitOk),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationNewRequestDurationConfigurationTimeoutTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileOtpForDefault(apiObject *mfa.DeviceAuthenticationPolicyCommonMobileOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileOtpTFObjectTypes), nil
	}

	failure, d := toStateMfaDevicePolicyMobileOtpFailureForDefault(apiObject.GetFailureOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileOtpTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"failure": failure,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileOtpTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileOtpFailureForDefault(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailure, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileOtpFailureTFObjectTypes), nil
	}

	coolDown, d := toStateMfaDevicePolicyMobileOtpFailureCooldownForDefault(apiObject.GetCoolDownOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileOtpFailureTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"count":     framework.Int32OkToTF(apiObject.GetCountOk()),
		"cool_down": coolDown,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileOtpFailureTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileOtpFailureCooldownForDefault(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyRememberMe(apiObject *mfa.DeviceAuthenticationPolicyCommonRememberMe, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyRememberMeTFObjectTypes), nil
	}

	web, d := toStateMfaDevicePolicyRememberMeWeb(apiObject.GetWebOk())
	diags.Append(d...)

	o := map[string]attr.Value{
		"web": web,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyRememberMeTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyRememberMeWeb(apiObject *mfa.DeviceAuthenticationPolicyCommonRememberMeWeb, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyRememberMeWebTFObjectTypes), nil
	}

	lifeTime, d := toStateMfaDevicePolicyRememberMeWebLifeTime(apiObject.GetLifeTimeOk())
	diags.Append(d...)

	o := map[string]attr.Value{
		"enabled":   framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"life_time": lifeTime,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyRememberMeWebTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyRememberMeWebLifeTime(apiObject *mfa.DeviceAuthenticationPolicyCommonRememberMeWebLifeTime, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTimePeriodTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTimePeriodTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyNotificationsPolicy(apiObject *mfa.DeviceAuthenticationPolicyCommonNotificationsPolicy, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyNotificationsPolicyTFObjectTypes), nil
	}

	id, idOk := apiObject.GetIdOk()
	if !idOk {
		return types.ObjectNull(MFADevicePolicyNotificationsPolicyTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"id": types.StringValue(*id),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyNotificationsPolicyTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}
