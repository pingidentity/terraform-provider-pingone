package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type MFADevicePolicyResource serviceClientType

type MFADevicePolicyResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                  types.String                 `tfsdk:"name"`
	Authentication        types.Object                 `tfsdk:"authentication"`
	NewDeviceNotification types.String                 `tfsdk:"new_device_notification"`
	Sms                   types.Object                 `tfsdk:"sms"`
	Voice                 types.Object                 `tfsdk:"voice"`
	Email                 types.Object                 `tfsdk:"email"`
	Mobile                types.Object                 `tfsdk:"mobile"`
	Totp                  types.Object                 `tfsdk:"totp"`
	Fido2                 types.Object                 `tfsdk:"fido2"`
}

type MFADevicePolicyAuthenticationResourceModel struct {
	DeviceSelection types.String `tfsdk:"device_selection"`
}

type MFADevicePolicySmsResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyVoiceResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyEmailResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyTotpResourceModel MFADevicePolicyOfflineDeviceResourceModel

type MFADevicePolicyOfflineDeviceResourceModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	Otp             types.Object `tfsdk:"otp"`
	PairingDisabled types.Bool   `tfsdk:"pairing_disabled"`
}

type MFADevicePolicyOfflineDeviceOtpResourceModel struct {
	Failure  types.Object `tfsdk:"failure"`
	Lifetime types.Object `tfsdk:"lifetime"`
}

type MFADevicePolicyOtpResourceModel struct {
	Failure types.Object `tfsdk:"failure"`
}

type MFADevicePolicyFailureResourceModel struct {
	CoolDown types.Object `tfsdk:"cool_down"`
	Count    types.Int64  `tfsdk:"count"`
}

type MFADevicePolicyCooldownResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyPushTimeoutResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyLockDurationResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyPairingKeyLifetimeResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyTimePeriodResourceModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type MFADevicePolicyFido2ResourceModel struct {
	Enabled         types.Bool                   `tfsdk:"enabled"`
	Fido2PolicyId   pingonetypes.ResourceIDValue `tfsdk:"fido2_policy_id"`
	PairingDisabled types.Bool                   `tfsdk:"pairing_disabled"`
}

type MFADevicePolicyMobileResourceModel struct {
	Applications types.Map    `tfsdk:"applications"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Otp          types.Object `tfsdk:"otp"`
}

type MFADevicePolicyMobileApplicationResourceModel struct {
	AutoEnrolment       types.Object `tfsdk:"auto_enrollment"`
	DeviceAuthorization types.Object `tfsdk:"device_authorization"`
	IntegrityDetection  types.String `tfsdk:"integrity_detection"`
	Otp                 types.Object `tfsdk:"otp"`
	PairingDisabled     types.Bool   `tfsdk:"pairing_disabled"`
	PairingKeyLifetime  types.Object `tfsdk:"pairing_key_lifetime"`
	Push                types.Object `tfsdk:"push"`
	PushLimit           types.Object `tfsdk:"push_limit"`
	PushTimeout         types.Object `tfsdk:"push_timeout"`
}

type MFADevicePolicyMobileApplicationAutoEnrolmentResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyMobileApplicationDeviceAuthorizationResourceModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	ExtraVerification types.String `tfsdk:"extra_verification"`
}

type MFADevicePolicyMobileApplicationOtpResourceModel MFADevicePolicyEnabledResourceModel
type MFADevicePolicyMobileApplicationPushResourceModel MFADevicePolicyEnabledResourceModel
type MFADevicePolicyEnabledResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyPushLimitResourceModel struct {
	Count        types.Int64  `tfsdk:"count"`
	LockDuration types.Object `tfsdk:"lock_duration"`
	TimePeriod   types.Object `tfsdk:"time_period"`
}

var (
	MFADevicePolicyAuthenticationTFObjectTypes = map[string]attr.Type{
		"device_selection": types.StringType,
	}

	MFADevicePolicyOfflineDeviceTFObjectTypes = map[string]attr.Type{
		"enabled":          types.BoolType,
		"otp":              types.ObjectType{AttrTypes: MFADevicePolicyOfflineDeviceOtpTFObjectTypes},
		"pairing_disabled": types.BoolType,
	}

	MFADevicePolicyOfflineDeviceOtpTFObjectTypes = map[string]attr.Type{
		"failure":  types.ObjectType{AttrTypes: MFADevicePolicyFailureTFObjectTypes},
		"lifetime": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyFailureTFObjectTypes = map[string]attr.Type{
		"cool_down": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"count":     types.Int64Type,
	}

	MFADevicePolicyTimePeriodTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int64Type,
		"time_unit": types.StringType,
	}

	MFADevicePolicyMobileTFObjectTypes = map[string]attr.Type{
		"applications": types.MapType{ElemType: types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationTFObjectTypes}},
		"enabled":      types.BoolType,
		"otp":          types.ObjectType{AttrTypes: MFADevicePolicyMobileOtpTFObjectTypes},
	}

	MFADevicePolicyMobileApplicationTFObjectTypes = map[string]attr.Type{
		"auto_enrollment":      types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes},
		"device_authorization": types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes},
		"integrity_detection":  types.StringType,
		"otp":                  types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationOtpTFObjectTypes},
		"pairing_disabled":     types.BoolType,
		"pairing_key_lifetime": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"push":                 types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationPushTFObjectTypes},
		"push_limit":           types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationPushLimitTFObjectTypes},
		"push_timeout":         types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyMobileApplicationAutoEnrolmentTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	MFADevicePolicyMobileApplicationDeviceAuthorizationTFObjectTypes = map[string]attr.Type{
		"enabled":            types.BoolType,
		"extra_verification": types.StringType,
	}

	MFADevicePolicyMobileApplicationOtpTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	MFADevicePolicyMobileApplicationPushTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	MFADevicePolicyMobileApplicationPushLimitTFObjectTypes = map[string]attr.Type{
		"count":         types.Int64Type,
		"lock_duration": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"time_period":   types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyMobileOtpTFObjectTypes = map[string]attr.Type{
		"failure": types.ObjectType{AttrTypes: MFADevicePolicyMobileOtpFailureTFObjectTypes},
	}

	MFADevicePolicyMobileOtpFailureTFObjectTypes = map[string]attr.Type{
		"count":     types.Int64Type,
		"cool_down": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyTotpTFObjectTypes = map[string]attr.Type{
		"enabled":          types.BoolType,
		"otp":              types.ObjectType{AttrTypes: MFADevicePolicyTotpOtpTFObjectTypes},
		"pairing_disabled": types.BoolType,
	}

	MFADevicePolicyTotpOtpTFObjectTypes = map[string]attr.Type{
		"failure": types.ObjectType{AttrTypes: MFADevicePolicyTotpOtpFailureTFObjectTypes},
	}

	MFADevicePolicyTotpOtpFailureTFObjectTypes = map[string]attr.Type{
		"count":     types.Int64Type,
		"cool_down": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
	}

	MFADevicePolicyFido2TFObjectTypes = map[string]attr.Type{
		"enabled":          types.BoolType,
		"fido2_policy_id":  pingonetypes.ResourceIDType{},
		"pairing_disabled": types.BoolType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &MFADevicePolicyResource{}
	_ resource.ResourceWithConfigure   = &MFADevicePolicyResource{}
	_ resource.ResourceWithImportState = &MFADevicePolicyResource{}
)

// New Object
func NewMFADevicePolicyResource() resource.Resource {
	return &MFADevicePolicyResource{}
}

// Metadata
func (r *MFADevicePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_device_policy"
}

func (r *MFADevicePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const mobileApplicationsPushLimitCountDefault = 5
	const mobileApplicationsPushLimitCountMin = 1
	const mobileApplicationsPushLimitCountMax = 50

	const mobileApplicationsPushLimitLockDurationDurationDefault = 30
	const mobileApplicationsPushLimitTimePeriodDurationDefault = 10
	const mobileApplicationsOtpFailureCoolDownDurationDefault = 2

	const mobileOtpFailureCountDefault = 3
	const mobileOtpFailureCountMin = 1
	const mobileOtpFailureCountMax = 7

	const totpOtpFailureCountDefault = 3
	const totpOtpFailureCoolDownDurationDefault = 2

	// schema descriptions and validation settings
	deviceSelectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines the device selection method.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicySelectionEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST))

	newDeviceNotificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines whether a user should be notified if a new authentication method has been added to their account.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicyNewDeviceNotificationEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE))

	mobileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that allows configuration of mobile push/OTP device authentication policy settings.  This factor requires embedding the PingOne MFA SDK into a customer facing mobile application, and configuring as a Native application using the `pingone_application` resource.",
	)

	mobileApplicationsAutoEnrollmentEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true` if you want the application to allow Auto Enrollment. Auto Enrollment means that the user can authenticate for the first time from an unpaired device, and the successful authentication will result in the pairing of the device for MFA.",
	)

	mobileApplicationsDeviceAuthorizationExtraVerificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the level of further verification when device authorization is enabled. The PingOne platform performs an extra verification check by sending a \"silent\" push notification to the customer native application, and receives a confirmation in return.  By default, the PingOne platform does not perform the extra verification check.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMMFADEVICEPOLICYMOBILEEXTRAVERIFICATION_PERMISSIVE):  "the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as successfully completed",
		string(mfa.ENUMMFADEVICEPOLICYMOBILEEXTRAVERIFICATION_RESTRICTIVE): "the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as failed",
	})

	mobileApplicationsIntegrityDetectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls how authentication or registration attempts should proceed if a device integrity check does not receive a response.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMMFADEVICEPOLICYMOBILEINTEGRITYDETECTION_PERMISSIVE):  "if you want to allow the process to continue if a device integrity check does not receive a response",
		string(mfa.ENUMMFADEVICEPOLICYMOBILEINTEGRITYDETECTION_RESTRICTIVE): "if you want to block the user if a device integrity check does not receive a response",
	})

	mobileApplicationsPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the relevant application. You can use this option if you want to phase out an existing mobile application but want to allow users to continue using the application for authentication for existing devices.",
	)

	mobileApplicationsPushLimitLockDurationDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.",
	)

	mobileApplicationsPushLimitTimePeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.",
	)

	mobileApplicationsPushTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.",
	)

	mobileOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. The minimum value is `%d`, maximum is `%d`, and the default is `%d`.", mobileOtpFailureCountMin, mobileOtpFailureCountMax, mobileOtpFailureCountDefault),
	)

	mobileOtpFailureCoolDownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. The minimum value is `2`, maximum is `30`, and the default is `2`. Note that when using the \"onetime authentication\" feature, the user is not blocked after the maximum number of failures even if you specified a block duration.",
	)

	durationTimeUnitMinsSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	mobileApplicationsPairingKeyLifetimeTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)

	durationTimeUnitSecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the type of time unit for `duration`. Currently, the only permitted value is `%s`.", mfa.ENUMTIMEUNIT_SECONDS),
	).DefaultValue(string(mfa.ENUMTIMEUNIT_SECONDS))

	totpPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the TOTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	fido2PairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the FIDO2 method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA device policies for a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the MFA device policy to manage."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the MFA policy's unique name within the environment.").Description,
				Required:    true,
			},

			"authentication": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of authentication settings in the device policy.").Description,
				Optional:    true,
				Computed:    true,

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

			"sms": r.devicePolicyOfflineDeviceSchemaAttribute("SMS OTP"),

			"voice": r.devicePolicyOfflineDeviceSchemaAttribute("voice OTP"),

			"email": r.devicePolicyOfflineDeviceSchemaAttribute("email OTP"),

			"mobile": schema.SingleNestedAttribute{
				Description:         mobileDescription.Description,
				MarkdownDescription: mobileDescription.MarkdownDescription,
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the mobile device method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"applications": schema.MapNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A map of objects that specifies settings for a configured Mobile Application.  The ID of the application should be configured as the map key.").Description,
						Optional:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"auto_enrollment": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies auto enrollment settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         mobileApplicationsAutoEnrollmentEnabledDescription.Description,
											MarkdownDescription: mobileApplicationsAutoEnrollmentEnabledDescription.MarkdownDescription,
											Required:            true,
										},
									},
								},

								"device_authorization": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies device authorization settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies the enabled or disabled state of automatic MFA for native devices paired with the user, for the specified application.").Description,
											Required:    true,
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
									},
								},

								"otp": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether OTP authentication is enabled or disabled for the application in the policy.").Description,
											Required:    true,
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
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies pairing key lifetime settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"duration": schema.Int64Attribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the amount of time an issued pairing key can be used until it expires. Minimum is 1 minute and maximum is 48 hours. If this parameter is not provided, the duration is set to 10 minutes.").Description,
											Required:    true,
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
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether push notification is enabled or disabled for the application in the policy.").Description,
											Required:    true,
										},
									},
								},

								"push_limit": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push limit settings for the application in the policy.").Description,
									Optional:    true,
									Computed:    true,

									Default: objectdefault.StaticValue(types.ObjectValueMust(
										MFADevicePolicyMobileApplicationPushLimitTFObjectTypes,
										map[string]attr.Value{
											"count": types.Int64Value(mobileApplicationsPushLimitCountDefault),
											"lock_duration": types.ObjectValueMust(
												MFADevicePolicyTimePeriodTFObjectTypes,
												map[string]attr.Value{
													"duration":  types.Int64Value(mobileApplicationsPushLimitLockDurationDurationDefault),
													"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
												},
											),
											"time_period": types.ObjectValueMust(
												MFADevicePolicyTimePeriodTFObjectTypes,
												map[string]attr.Value{
													"duration":  types.Int64Value(mobileApplicationsPushLimitTimePeriodDurationDefault),
													"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
												},
											),
										},
									)),

									Attributes: map[string]schema.Attribute{
										"count": schema.Int64Attribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is `1` and the maximum value is `50`. If this parameter is not provided, the default value is `5`.").Description,
											Optional:    true,
											Computed:    true,

											Default: int64default.StaticInt64(mobileApplicationsPushLimitCountDefault),

											Validators: []validator.Int64{
												int64validator.Between(mobileApplicationsPushLimitCountMin, mobileApplicationsPushLimitCountMax),
											},
										},

										"lock_duration": schema.SingleNestedAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push limit lock duration settings for the application in the policy.").Description,
											Optional:    true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int64Attribute{
													Description:         mobileApplicationsPushLimitLockDurationDurationDescription.Description,
													MarkdownDescription: mobileApplicationsPushLimitLockDurationDurationDescription.MarkdownDescription,
													Required:            true,
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
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push limit time period settings for the application in the policy.").Description,
											Optional:    true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int64Attribute{
													Description:         mobileApplicationsPushLimitTimePeriodDurationDescription.Description,
													MarkdownDescription: mobileApplicationsPushLimitTimePeriodDurationDescription.MarkdownDescription,
													Required:            true,
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
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push timeout settings for the application in the policy.").Description,
									Optional:    true,

									Attributes: map[string]schema.Attribute{
										"duration": schema.Int64Attribute{
											Description:         mobileApplicationsPushTimeoutDurationDescription.Description,
											MarkdownDescription: mobileApplicationsPushTimeoutDurationDescription.MarkdownDescription,
											Required:            true,
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
							},
						},
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP settings for mobile applications in the policy.").Description,
						Optional:    true,
						Computed:    true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyMobileOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyMobileOtpFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int64Value(mobileOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int64Value(mobileApplicationsOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure settings for mobile applications in the policy.").Description,
								Required:    true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description:         mobileOtpFailureCountDescription.Description,
										MarkdownDescription: mobileOtpFailureCountDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int64default.StaticInt64(mobileOtpFailureCountDefault),

										Validators: []validator.Int64{
											int64validator.Between(mobileOtpFailureCountMin, mobileOtpFailureCountMax),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure cool down settings for mobile applications in the policy.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description:         mobileOtpFailureCoolDownDurationDescription.Description,
												MarkdownDescription: mobileOtpFailureCoolDownDurationDescription.MarkdownDescription,
												Required:            true,
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
				},
			},

			"totp": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the TOTP method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         totpPairingDisabledDescription.Description,
						MarkdownDescription: totpPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP settings.").Description,
						Optional:    true,
						Computed:    true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyTotpOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyTotpOtpFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int64Value(totpOtpFailureCountDefault),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int64Value(totpOtpFailureCoolDownDurationDefault),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.").Description,
												Required:    true,
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

									"count": schema.Int64Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.").Description,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},

			"fido2": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of FIDO2 device authentication policy settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the FIDO2 method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         fido2PairingDisabledDescription.Description,
						MarkdownDescription: fido2PairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"fido2_policy_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the resource UUID that represents the FIDO2 policy in PingOne. This property can be null / left undefined. When null, the environment's default FIDO2 Policy is used.  Must be a valid PingOne resource ID.").Description,
						Optional:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},
		},
	}
}

func (r *MFADevicePolicyResource) devicePolicyOfflineDeviceSchemaAttribute(descriptionMethod string) schema.SingleNestedAttribute {

	const otpFailureCountDefault = 3
	const otpFailureCoolDownDurationDefault = 0
	const otpLifetimeDurationDefault = 30

	pairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A boolean that, when set to `true`, prevents users from pairing new devices with the %s method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.", descriptionMethod),
	).DefaultValue(false)

	otpCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	return schema.SingleNestedAttribute{
		Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s device authentication policy settings.", descriptionMethod)).Description,
		Required:    true,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A boolean that specifies whether the %s method is enabled or disabled in the policy.", descriptionMethod)).Description,
				Required:    true,
			},

			"pairing_disabled": schema.BoolAttribute{
				Description:         pairingDisabledDescription.Description,
				MarkdownDescription: pairingDisabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"otp": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s settings.", descriptionMethod)).Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					MFADevicePolicyOfflineDeviceOtpTFObjectTypes,
					map[string]attr.Value{
						"failure": types.ObjectValueMust(
							MFADevicePolicyFailureTFObjectTypes,
							map[string]attr.Value{
								"count": types.Int64Value(otpFailureCountDefault),
								"cool_down": types.ObjectValueMust(
									MFADevicePolicyTimePeriodTFObjectTypes,
									map[string]attr.Value{
										"duration":  types.Int64Value(otpFailureCoolDownDurationDefault),
										"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
									},
								),
							},
						),
						"lifetime": types.ObjectValueMust(
							MFADevicePolicyTimePeriodTFObjectTypes,
							map[string]attr.Value{
								"duration":  types.Int64Value(otpLifetimeDurationDefault),
								"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
							},
						),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"failure": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s failure settings.", descriptionMethod)).Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"cool_down": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s failure cool down settings.", descriptionMethod)).Description,
								Required:    true,

								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.").Description,
										Required:    true,
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

							"count": schema.Int64Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.").Description,
								Required:    true,
							},
						},
					},

					"lifetime": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s lifetime settings.", descriptionMethod)).Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int64Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) that the passcode is valid before it expires.").Description,
								Required:    true,
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
				},
			},
		},
	}
}

func (r *MFADevicePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MFADevicePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state MFADevicePolicyResourceModel

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

	// Build the model for the API
	mFADevicePolicy, d := plan.expandCreate(ctx, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicyPostResponse
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.CreateDeviceAuthenticationPolicies(ctx, plan.EnvironmentId.ValueString()).DeviceAuthenticationPolicyPost(*mFADevicePolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateDeviceAuthenticationPolicies",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toStateCreate(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *MFADevicePolicyResourceModel

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

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneDeviceAuthenticationPolicy",
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MFADevicePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MFADevicePolicyResourceModel

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

	// Build the model for the API
	mFADevicePolicy, d := plan.expand(ctx, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).DeviceAuthenticationPolicy(*mFADevicePolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDeviceAuthenticationPolicy",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *MFADevicePolicyResourceModel

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

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteDeviceAuthenticationPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MFADevicePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "mfa_device_policy_id",
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

func (p *MFADevicePolicyResourceModel) expand(ctx context.Context, apiClient *management.APIClient) (*mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	// SMS
	var smsPlan MFADevicePolicySmsResourceModel
	diags.Append(p.Sms.As(ctx, &smsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	sms, d := smsPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Voice
	var voicePlan MFADevicePolicyVoiceResourceModel
	diags.Append(p.Voice.As(ctx, &voicePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	voice, d := voicePlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Email
	var emailPlan MFADevicePolicyEmailResourceModel
	diags.Append(p.Email.As(ctx, &emailPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	email, d := emailPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Mobile
	var mobilePlan MFADevicePolicyMobileResourceModel
	diags.Append(p.Mobile.As(ctx, &mobilePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	mobile, d := mobilePlan.expand(ctx, apiClient, p.EnvironmentId.ValueString())
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// TOTP
	var totpPlan MFADevicePolicyTotpResourceModel
	diags.Append(p.Totp.As(ctx, &totpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	totp, d := totpPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Main object
	data := mfa.NewDeviceAuthenticationPolicy(
		p.Name.ValueString(),
		*sms,
		*voice,
		*email,
		*mobile,
		*totp,
		false,
		false,
	)

	// FIDO2
	if !p.Fido2.IsNull() && !p.Fido2.IsUnknown() {
		var fido2Plan MFADevicePolicyFido2ResourceModel
		diags.Append(p.Fido2.As(ctx, &fido2Plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		fido2 := fido2Plan.expand()

		data.SetFido2(*fido2)
	}

	// Authentication
	if !p.Authentication.IsNull() && !p.Authentication.IsUnknown() {
		var authenticationPlan MFADevicePolicyAuthenticationResourceModel
		diags.Append(p.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetAuthentication(
			*mfa.NewDeviceAuthenticationPolicyAuthentication(
				mfa.EnumMFADevicePolicySelection(authenticationPlan.DeviceSelection.ValueString()),
			),
		)
	}

	// New Device Notification
	if !p.NewDeviceNotification.IsNull() && !p.NewDeviceNotification.IsUnknown() {
		data.SetNewDeviceNotification(
			mfa.EnumMFADevicePolicyNewDeviceNotification(p.NewDeviceNotification.ValueString()),
		)
	}

	return data, diags
}

func (p *MFADevicePolicyResourceModel) expandCreate(ctx context.Context, apiClient *management.APIClient) (*mfa.DeviceAuthenticationPolicyPost, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, diags := p.expand(ctx, apiClient)

	return &mfa.DeviceAuthenticationPolicyPost{
		DeviceAuthenticationPolicy: data,
	}, diags
}

func (p *MFADevicePolicySmsResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDevice, diag.Diagnostics) {
	data := MFADevicePolicyOfflineDeviceResourceModel(*p)
	return data.expand(ctx)
}

func (p *MFADevicePolicyVoiceResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDevice, diag.Diagnostics) {
	data := MFADevicePolicyOfflineDeviceResourceModel(*p)
	return data.expand(ctx)
}

func (p *MFADevicePolicyEmailResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDevice, diag.Diagnostics) {
	data := MFADevicePolicyOfflineDeviceResourceModel(*p)
	return data.expand(ctx)
}

func (p *MFADevicePolicyOfflineDeviceResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDevice, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	// OTP
	var otpPlan MFADevicePolicyOfflineDeviceOtpResourceModel
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

	data := mfa.NewDeviceAuthenticationPolicyOfflineDevice(
		p.Enabled.ValueBool(),
		*otp,
	)

	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	return data, diags
}

func (p *MFADevicePolicyOfflineDeviceOtpResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDeviceOtp, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Lifetime
	var lifetimePlan MFADevicePolicyPairingKeyLifetimeResourceModel
	diags.Append(p.Lifetime.As(ctx, &lifetimePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	lifetime := mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpLifeTime(
		int32(lifetimePlan.Duration.ValueInt64()),
		mfa.EnumTimeUnit(lifetimePlan.TimeUnit.ValueString()),
	)

	// Failure
	var failurePlan MFADevicePolicyFailureResourceModel
	diags.Append(p.Failure.As(ctx, &failurePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	// Failure Cool Down
	var failureCooldownPlan MFADevicePolicyCooldownResourceModel
	diags.Append(failurePlan.CoolDown.As(ctx, &failureCooldownPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	failure := mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
		int32(failurePlan.Count.ValueInt64()),
		*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(
			int32(failureCooldownPlan.Duration.ValueInt64()),
			mfa.EnumTimeUnit(failureCooldownPlan.TimeUnit.ValueString()),
		),
	)

	data := mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtp(
		*lifetime,
		*failure,
	)

	return data, diags
}

func (p *MFADevicePolicyMobileResourceModel) expand(ctx context.Context, apiClient *management.APIClient, environmentId string) (*mfa.DeviceAuthenticationPolicyMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Otp
	var otpPlan MFADevicePolicyOtpResourceModel
	diags.Append(p.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	// Otp Failure
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

	otp := mfa.NewDeviceAuthenticationPolicyMobileOtp(
		*failure,
	)

	// Main object
	data := mfa.NewDeviceAuthenticationPolicyMobile(
		p.Enabled.ValueBool(),
		*otp,
	)

	// Applications
	if !p.Applications.IsNull() && !p.Applications.IsUnknown() {
		applicationsPlan := make(map[string]MFADevicePolicyMobileApplicationResourceModel, len(p.Applications.Elements()))
		diags.Append(p.Applications.ElementsAs(ctx, &applicationsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		applications := make([]mfa.DeviceAuthenticationPolicyMobileApplicationsInner, 0)

		for applicationId, applicationPlan := range applicationsPlan {
			application, d := applicationPlan.expand(ctx, apiClient, environmentId, applicationId)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			if application != nil {
				applications = append(applications, *application)
			}
		}

		data.SetApplications(applications)
	}

	return data, diags
}

func (p *MFADevicePolicyFailureResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailure, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Cooldown
	var cooldownPlan MFADevicePolicyCooldownResourceModel
	diags.Append(p.CoolDown.As(ctx, &cooldownPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	data := mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
		int32(p.Count.ValueInt64()),
		*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(
			int32(cooldownPlan.Duration.ValueInt64()),
			mfa.EnumTimeUnit(cooldownPlan.TimeUnit.ValueString()),
		),
	)

	return data, diags
}

func (p *MFADevicePolicyMobileApplicationResourceModel) expand(ctx context.Context, apiClient *management.APIClient, environmentId, applicationId string) (*mfa.DeviceAuthenticationPolicyMobileApplicationsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	application, d := checkApplicationForMobileApp(ctx, apiClient, environmentId, applicationId)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := mfa.NewDeviceAuthenticationPolicyMobileApplicationsInner(
		applicationId,
	)

	// Auto enrollment
	if !p.AutoEnrolment.IsNull() && !p.AutoEnrolment.IsUnknown() {
		var plan MFADevicePolicyMobileApplicationAutoEnrolmentResourceModel
		diags.Append(p.AutoEnrolment.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetAutoEnrollment(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerAutoEnrollment(
				plan.Enabled.ValueBool(),
			),
		)
	}

	// Device authorisation
	if !p.DeviceAuthorization.IsNull() && !p.DeviceAuthorization.IsUnknown() {
		var plan MFADevicePolicyMobileApplicationDeviceAuthorizationResourceModel
		diags.Append(p.DeviceAuthorization.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		deviceAuthorization := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerDeviceAuthorization(
			plan.Enabled.ValueBool(),
		)

		if !plan.ExtraVerification.IsNull() && !plan.ExtraVerification.IsUnknown() {
			deviceAuthorization.SetExtraVerification(mfa.EnumMFADevicePolicyMobileExtraVerification(plan.ExtraVerification.ValueString()))
		}

		data.SetDeviceAuthorization(deviceAuthorization)
	}

	// Integrity detection
	if p.IntegrityDetection.IsNull() && application.GetMobile().IntegrityDetection.GetMode() == management.ENUMENABLEDSTATUS_ENABLED {
		diags.AddError(
			"Invalid mobile application integrity detection setting",
			fmt.Sprintf("An application ID, %s, configured as the map key in `mobile.applications` has integrity detection enabled. This policy must specify the level of integrity detection in the `mobile.application.integrity_detection` parameter.", applicationId),
		)
		return nil, diags
	}

	if !p.IntegrityDetection.IsNull() && !p.IntegrityDetection.IsUnknown() {
		data.SetIntegrityDetection(mfa.EnumMFADevicePolicyMobileIntegrityDetection(p.IntegrityDetection.ValueString()))

		if application.GetMobile().IntegrityDetection.GetMode() != management.ENUMENABLEDSTATUS_ENABLED {
			// error - this has no effect
			diags.AddError(
				"Mobile application integrity detection setting has no effect",
				fmt.Sprintf("An application ID, %s, configured as the map key in `mobile.applications` has integrity detection disabled. Setting the `mobile.application.integrity_detection` parameter has no effect.", applicationId),
			)

			return nil, diags
		}
	}

	// OTP
	if !p.Otp.IsNull() && !p.Otp.IsUnknown() {
		var plan MFADevicePolicyMobileApplicationOtpResourceModel
		diags.Append(p.Otp.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetOtp(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerOtp(
			plan.Enabled.ValueBool(),
		))
	}

	// Pairing Disabled
	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	// Pairing Key Lifetime
	if !p.PairingKeyLifetime.IsNull() && !p.PairingKeyLifetime.IsUnknown() {
		var plan MFADevicePolicyPairingKeyLifetimeResourceModel
		diags.Append(p.PairingKeyLifetime.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPairingKeyLifetime(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPairingKeyLifetime(
				int32(plan.Duration.ValueInt64()),
				mfa.EnumTimeUnitPairingKeyLifetime(plan.TimeUnit.ValueString()),
			),
		)
	}

	// Push
	if !p.Push.IsNull() && !p.Push.IsUnknown() {
		var plan MFADevicePolicyMobileApplicationPushResourceModel
		diags.Append(p.Push.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPush(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPush(
				plan.Enabled.ValueBool(),
			),
		)
	}

	// Push Limit
	if !p.PushLimit.IsNull() && !p.PushLimit.IsUnknown() {
		var plan MFADevicePolicyPushLimitResourceModel
		diags.Append(p.PushLimit.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		pushLimit, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPushLimit(*pushLimit)
	}

	// Push Timeout
	if !p.PushTimeout.IsNull() && !p.PushTimeout.IsUnknown() {
		var plan MFADevicePolicyPushTimeoutResourceModel
		diags.Append(p.PushTimeout.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPushTimeout(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushTimeout(
				int32(plan.Duration.ValueInt64()),
				mfa.EnumTimeUnitPushTimeout(plan.TimeUnit.ValueString()),
			),
		)
	}

	return data, diags
}

func (p *MFADevicePolicyPushLimitResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushLimit, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimit()

	if !p.Count.IsNull() && !p.Count.IsUnknown() {
		data.SetCount(int32(p.Count.ValueInt64()))
	}

	if !p.LockDuration.IsNull() && !p.LockDuration.IsUnknown() {
		var plan MFADevicePolicyLockDurationResourceModel
		diags.Append(p.LockDuration.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetLockDuration(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimitLockDuration(
				int32(plan.Duration.ValueInt64()),
				mfa.EnumTimeUnit(plan.TimeUnit.ValueString()),
			),
		)
	}

	if !p.TimePeriod.IsNull() && !p.TimePeriod.IsUnknown() {
		var plan MFADevicePolicyTimePeriodResourceModel
		diags.Append(p.TimePeriod.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetTimePeriod(
			*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimitTimePeriod(
				int32(plan.Duration.ValueInt64()),
				mfa.EnumTimeUnit(plan.TimeUnit.ValueString()),
			),
		)
	}

	return data, diags
}

func (p *MFADevicePolicyTotpResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyTotp, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Otp
	var otpPlan MFADevicePolicyOtpResourceModel
	diags.Append(p.Otp.As(ctx, &otpPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	// Otp Failure
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

	otp := mfa.NewDeviceAuthenticationPolicyTotpOtp(
		*failure,
	)

	data := mfa.NewDeviceAuthenticationPolicyTotp(
		p.Enabled.ValueBool(),
		*otp,
	)

	// Pairing Disabled
	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	return data, diags
}

func (p *MFADevicePolicyFido2ResourceModel) expand() *mfa.DeviceAuthenticationPolicyFido2 {

	data := mfa.NewDeviceAuthenticationPolicyFido2(
		p.Enabled.ValueBool(),
	)

	if !p.PairingDisabled.IsNull() && !p.PairingDisabled.IsUnknown() {
		data.SetPairingDisabled(p.PairingDisabled.ValueBool())
	}

	if !p.Fido2PolicyId.IsNull() && !p.Fido2PolicyId.IsUnknown() {
		data.SetFido2PolicyId(p.Fido2PolicyId.ValueString())
	}

	return data
}

func (p *MFADevicePolicyResourceModel) toState(apiObject *mfa.DeviceAuthenticationPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	p.Authentication, d = toStateMfaDevicePolicyAuthentication(apiObject.GetAuthenticationOk())
	diags.Append(d...)

	p.NewDeviceNotification = framework.EnumOkToTF(apiObject.GetNewDeviceNotificationOk())

	p.Sms, d = toStateMfaDevicePolicySms(apiObject.GetSmsOk())
	diags.Append(d...)

	p.Voice, d = toStateMfaDevicePolicyVoice(apiObject.GetVoiceOk())
	diags.Append(d...)

	p.Email, d = toStateMfaDevicePolicyEmail(apiObject.GetEmailOk())
	diags.Append(d...)

	p.Mobile, d = toStateMfaDevicePolicyMobile(apiObject.GetMobileOk())
	diags.Append(d...)

	p.Totp, d = toStateMfaDevicePolicyTotp(apiObject.GetTotpOk())
	diags.Append(d...)

	p.Fido2, d = toStateMfaDevicePolicyFido2(apiObject.GetFido2Ok())
	diags.Append(d...)

	return diags
}

func (p *MFADevicePolicyResourceModel) toStateCreate(apiObject *mfa.DeviceAuthenticationPolicyPostResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	return p.toState(apiObject.DeviceAuthenticationPolicy)
}

func toStateMfaDevicePolicyAuthentication(apiObject *mfa.DeviceAuthenticationPolicyAuthentication, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyAuthenticationTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"device_selection": framework.EnumOkToTF(apiObject.GetDeviceSelectionOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyAuthenticationTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyOfflineDevice(apiObject *mfa.DeviceAuthenticationPolicyOfflineDevice, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyOfflineDeviceTFObjectTypes), nil
	}

	otp, d := toStateMfaDevicePolicyOfflineDeviceOtp(apiObject.GetOtpOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyOfflineDeviceTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":              otp,
		"pairing_disabled": framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyOfflineDeviceTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyOfflineDeviceOtp(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyOfflineDeviceOtpTFObjectTypes), nil
	}

	failure, d := toStateMfaDevicePolicyOfflineDeviceOtpFailure(apiObject.GetFailureOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyOfflineDeviceOtpTFObjectTypes), diags
	}

	lifetime, d := toStateMfaDevicePolicyOfflineDeviceOtpLifeTime(apiObject.GetLifeTimeOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyOfflineDeviceOtpTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"failure":  failure,
		"lifetime": lifetime,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyOfflineDeviceOtpTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyOfflineDeviceOtpFailure(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailure, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyFailureTFObjectTypes), nil
	}

	coolDown, d := toStateMfaDevicePolicyOfflineDeviceOtpFailureCoolDown(apiObject.GetCoolDownOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyFailureTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"cool_down": coolDown,
		"count":     framework.Int32OkToTF(apiObject.GetCountOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyFailureTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyOfflineDeviceOtpFailureCoolDown(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyOfflineDeviceOtpLifeTime(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpLifeTime, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicySms(apiObject *mfa.DeviceAuthenticationPolicyOfflineDevice, ok bool) (types.Object, diag.Diagnostics) {
	return toStateMfaDevicePolicyOfflineDevice(apiObject, ok)
}

func toStateMfaDevicePolicyVoice(apiObject *mfa.DeviceAuthenticationPolicyOfflineDevice, ok bool) (types.Object, diag.Diagnostics) {
	return toStateMfaDevicePolicyOfflineDevice(apiObject, ok)
}

func toStateMfaDevicePolicyEmail(apiObject *mfa.DeviceAuthenticationPolicyOfflineDevice, ok bool) (types.Object, diag.Diagnostics) {
	return toStateMfaDevicePolicyOfflineDevice(apiObject, ok)
}

func toStateMfaDevicePolicyMobile(apiObject *mfa.DeviceAuthenticationPolicyMobile, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileTFObjectTypes), nil
	}

	applications, d := toStateMfaDevicePolicyMobileApplications(apiObject.GetApplicationsOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileTFObjectTypes), diags
	}

	otp, d := toStateMfaDevicePolicyMobileOtp(apiObject.GetOtpOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"applications": applications,
		"enabled":      framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":          otp,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplications(apiObject []mfa.DeviceAuthenticationPolicyMobileApplicationsInner, ok bool) (types.Map, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: MFADevicePolicyMobileApplicationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.MapNull(tfObjType), nil
	}

	objectList := map[string]attr.Value{}
	for _, application := range apiObject {

		autoEnrolment, d := toStateMfaDevicePolicyMobileApplicationsAutoEnrolment(application.GetAutoEnrollmentOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		deviceAuthorization, d := toStateMfaDevicePolicyMobileApplicationsDeviceAuthorization(application.GetDeviceAuthorizationOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		otp, d := toStateMfaDevicePolicyMobileApplicationsOtp(application.GetOtpOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		pairingKeyLifetime, d := toStateMfaDevicePolicyMobileApplicationsPairingKeyLifetime(application.GetPairingKeyLifetimeOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		push, d := toStateMfaDevicePolicyMobileApplicationsPush(application.GetPushOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		pushLimit, d := toStateMfaDevicePolicyMobileApplicationsPushLimit(application.GetPushLimitOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		pushTimeout, d := toStateMfaDevicePolicyMobileApplicationsPushTimeout(application.GetPushTimeoutOk())
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(tfObjType), diags
		}

		o := map[string]attr.Value{
			"auto_enrollment":      autoEnrolment,
			"device_authorization": deviceAuthorization,
			"integrity_detection":  framework.EnumOkToTF(application.GetIntegrityDetectionOk()),
			"otp":                  otp,
			"pairing_disabled":     framework.BoolOkToTF(application.GetPairingDisabledOk()),
			"pairing_key_lifetime": pairingKeyLifetime,
			"push":                 push,
			"push_limit":           pushLimit,
			"push_timeout":         pushTimeout,
		}

		objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationTFObjectTypes, o)
		diags.Append(d...)

		objectList[application.GetId()] = objValue
	}

	returnVar, d := types.MapValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateMfaDevicePolicyMobileApplicationsAutoEnrolment(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerAutoEnrollment, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsDeviceAuthorization(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerDeviceAuthorization, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsOtp(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerOtp, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsPush(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPush, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyMobileApplicationPushTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyMobileApplicationsPushLimit(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushLimit, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes), nil
	}

	lockDuration, d := toStateMfaDevicePolicyMobileApplicationsPushLimitLockDuration(apiObject.GetLockDurationOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyMobileApplicationPushLimitTFObjectTypes), diags
	}

	timePeriod, d := toStateMfaDevicePolicyMobileApplicationsPushLimitTimePeriod(apiObject.GetTimePeriodOk())
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

func toStateMfaDevicePolicyMobileApplicationsPushLimitLockDuration(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushLimitLockDuration, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsPushLimitTimePeriod(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushLimitTimePeriod, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsPairingKeyLifetime(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPairingKeyLifetime, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileApplicationsPushTimeout(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushTimeout, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyMobileOtp(apiObject *mfa.DeviceAuthenticationPolicyMobileOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileOtpTFObjectTypes), nil
	}

	failure, d := toStateMfaDevicePolicyMobileOtpFailure(apiObject.GetFailureOk())
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

func toStateMfaDevicePolicyMobileOtpFailure(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailure, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileOtpFailureTFObjectTypes), nil
	}

	coolDown, d := toStateMfaDevicePolicyMobileOtpFailureCooldown(apiObject.GetCoolDownOk())
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

func toStateMfaDevicePolicyMobileOtpFailureCooldown(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyTotp(apiObject *mfa.DeviceAuthenticationPolicyTotp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTotpTFObjectTypes), nil
	}

	otp, d := toStateMfaDevicePolicyTotpOtp(apiObject.GetOtpOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyTotpTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp":              otp,
		"pairing_disabled": framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTotpTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyTotpOtp(apiObject *mfa.DeviceAuthenticationPolicyTotpOtp, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTotpOtpTFObjectTypes), nil
	}

	failure, d := toStateMfaDevicePolicyTotpOtpFailure(apiObject.GetFailureOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyTotpOtpTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"failure": failure,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTotpOtpTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyTotpOtpFailure(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailure, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyTotpOtpFailureTFObjectTypes), nil
	}

	coolDown, d := toStateMfaDevicePolicyTotpOtpFailureCooldown(apiObject.GetCoolDownOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(MFADevicePolicyTotpOtpFailureTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"count":     framework.Int32OkToTF(apiObject.GetCountOk()),
		"cool_down": coolDown,
	}

	objValue, d := types.ObjectValue(MFADevicePolicyTotpOtpFailureTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyTotpOtpFailureCooldown(apiObject *mfa.DeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown, ok bool) (types.Object, diag.Diagnostics) {
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

func toStateMfaDevicePolicyFido2(apiObject *mfa.DeviceAuthenticationPolicyFido2, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyFido2TFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled":          framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"fido2_policy_id":  framework.PingOneResourceIDOkToTF(apiObject.GetFido2PolicyIdOk()),
		"pairing_disabled": framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyFido2TFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func checkApplicationForMobileApp(ctx context.Context, apiClient *management.APIClient, environmentId, applicationId string) (*management.ApplicationOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	var response *management.ReadOneApplication200Response
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.ReadOneApplication(ctx, environmentId, applicationId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
		},
		"ReadOneApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if response == nil {
		diags.AddError(
			"Application not found",
			fmt.Sprintf("An application ID, %s, configured as the map key in the `mobile.applications` set does not exist", applicationId),
		)

		return nil, diags
	}

	var oidcObject *management.ApplicationOIDC

	// check if oidc and native
	if (response.ApplicationOIDC == nil) || (response.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_NATIVE_APP && response.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_CUSTOM_APP) {
		diags.AddError(
			"Invalid application type",
			fmt.Sprintf("An application ID, %s, configured as the map key in `mobile.applications` is not of type OIDC.  To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.", applicationId),
		)
		return nil, diags
	} else {
		oidcObject = response.ApplicationOIDC
	}

	// check if mobile set and package/bundle set
	if _, ok := response.ApplicationOIDC.GetMobileOk(); !ok {
		diags.AddError(
			"Missing application configuration",
			fmt.Sprintf("An application ID, %s, configured as the map key in `mobile.applications` does not contain mobile application configuration.  To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.", applicationId),
		)

		return nil, diags
	}

	if v, ok := response.ApplicationOIDC.GetMobileOk(); ok {

		_, bundleIDOk := v.GetBundleIdOk()
		_, packageNameOk := v.GetPackageNameOk()
		_, huaweiAppIdOk := v.GetHuaweiAppIdOk()

		if !bundleIDOk && !packageNameOk && !huaweiAppIdOk {
			diags.AddError(
				"Missing application configuration",
				fmt.Sprintf("An application ID, %s, configured as the map key in `mobile.applications` does not contain mobile application configuration.  To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.", applicationId),
			)

			return nil, diags
		}
	}

	return oidcObject, diags
}
