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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

const (
	POLICY_TYPE_PINGONE_MFA = "pingone_mfa"
	POLICY_TYPE_PINGID      = "pingid"
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
	Sms                   types.Object                 `tfsdk:"sms"`
	Voice                 types.Object                 `tfsdk:"voice"`
	Email                 types.Object                 `tfsdk:"email"`
	Mobile                types.Object                 `tfsdk:"mobile"`
	Totp                  types.Object                 `tfsdk:"totp"`
	Fido2                 types.Object                 `tfsdk:"fido2"`
	Desktop               types.Object                 `tfsdk:"desktop"`
	Yubikey               types.Object                 `tfsdk:"yubikey"`
	OathToken             types.Object                 `tfsdk:"oath_token"`
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

	// Default values for oath_token
	oathTokenDefault = types.ObjectValueMust(
		MFADevicePolicyPingIDDeviceTFObjectTypes,
		map[string]attr.Value{
			"enabled": types.BoolValue(false),
			"otp": types.ObjectValueMust(
				MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
				map[string]attr.Value{
					"failure": types.ObjectValueMust(
						MFADevicePolicyFailureTFObjectTypes,
						map[string]attr.Value{
							"count": types.Int32Value(3),
							"cool_down": types.ObjectValueMust(
								MFADevicePolicyTimePeriodTFObjectTypes,
								map[string]attr.Value{
									"duration":  types.Int32Value(2),
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
)

// Framework interfaces
var (
	_ resource.Resource                = &MFADevicePolicyDefaultResource{}
	_ resource.ResourceWithConfigure   = &MFADevicePolicyDefaultResource{}
	_ resource.ResourceWithImportState = &MFADevicePolicyDefaultResource{}
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

	totpUriParametersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map of string key:value pairs that specifies `otpauth` URI parameters. For example, if you provide a value for the `issuer` parameter, then authenticators that support that parameter will display the text you specify together with the OTP (in addition to the username). This can help users recognize which application the OTP is for. If you intend on using the same MFA policy for multiple applications, choose a name that reflects the group of applications.",
	)

	fido2PairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the FIDO2 method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	promptForNicknameOnPairingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to overwrite the default MFA device policy, or create it if it doesn't already exist.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the default MFA device policy in."),
			),

			"policy_type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of MFA device policy. Set to `PINGONE_MFA` for standard PingOne MFA environments, or `PINGID` for environments with PingID integration. This field is immutable and will trigger a replace plan if changed.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.OneOf(POLICY_TYPE_PINGONE_MFA, POLICY_TYPE_PINGID),
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name to apply to the default MFA device policy.").Description,
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
						Validators: []validator.Bool{
							boolvalidator.MustBeTrueIfPathSetToValue(
								types.StringValue(POLICY_TYPE_PINGID),
								path.MatchRoot("policy_type"),
							),
						},
					},					"applications": schema.MapNestedAttribute{
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
										"duration": schema.Int32Attribute{
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
											Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is `1` and the maximum value is `50`. If this parameter is not provided, the default value is `5`.").Description,
											Optional:    true,
											Computed:    true,

											Default: int32default.StaticInt32(mobileApplicationsPushLimitCountDefault),

											Validators: []validator.Int32{
												int32validator.Between(mobileApplicationsPushLimitCountMin, mobileApplicationsPushLimitCountMax),
											},
										},

										"lock_duration": schema.SingleNestedAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies push limit lock duration settings for the application in the policy.").Description,
											Optional:    true,

											Attributes: map[string]schema.Attribute{
												"duration": schema.Int32Attribute{
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
												"duration": schema.Int32Attribute{
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
										"duration": schema.Int32Attribute{
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
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure settings for mobile applications in the policy.").Description,
								Required:    true,

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
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure cool down settings for mobile applications in the policy.").Description,
										Required:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
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

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
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
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
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

									"count": schema.Int32Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.").Description,
										Required:    true,
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
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of FIDO2 device authentication policy settings.").Description,
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(MFADevicePolicyFido2TFObjectTypes,
					map[string]attr.Value{
						"enabled":                        types.BoolValue(false),
						"fido2_policy_id":                framework.PingOneResourceIDToTF(""),
						"pairing_disabled":               types.BoolNull(),
						"prompt_for_nickname_on_pairing": types.BoolNull(),
					}),
				),
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
						CustomType:  pingonetypes.ResourceIDType{},
					},

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"desktop": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of PingID desktop device authentication policy settings. Only applicable when `policy_type` is `PINGID`.").Description,
				Optional:    true,

				Validators: []validator.Object{
					objectvalidator.ConflictsIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGONE_MFA),
						path.MatchRoot("policy_type"),
					),
				},				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the desktop device method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure settings for desktop devices.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between 1 and 7.").Description,
										Optional:    true,

										Validators: []validator.Int32{
											int32validator.Between(1, 7),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between 1 SECONDS and 30 MINUTES.").Description,
												Required:    true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(1, 1800),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																"If `time_unit` is `SECONDS`, the allowed duration range is 1 - 1800.",
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(1, 30),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																"If `time_unit` is `MINUTES`, the allowed duration range is 1 - 30.",
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
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that, when set to `true`, prevents users from pairing new desktop devices.").Description,
						Optional:    true,
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies pairing key lifetime settings for desktop devices.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between 1 MINUTES and 48 HOURS.").Description,
								Required:    true,

								Validators: []validator.Int32{
									int32validator.Any(
										int32validator.All(
											int32validator.Between(1, 2880),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												"If `time_unit` is `MINUTES`, the allowed duration range is 1 - 2880.",
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int32validator.All(
											int32validator.Between(1, 48),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`HOURS`),
												"If `time_unit` is `HOURS`, the allowed duration range is 1 - 48.",
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
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of PingID Yubikey device authentication policy settings. Only applicable when `policy_type` is `PINGID`.").Description,
				Optional:    true,

				Validators: []validator.Object{
					objectvalidator.ConflictsIfMatchesPathValue(
						types.StringValue(POLICY_TYPE_PINGONE_MFA),
						path.MatchRoot("policy_type"),
					),
				},				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the Yubikey device method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure settings for Yubikey devices.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between 1 and 7.").Description,
										Optional:    true,

										Validators: []validator.Int32{
											int32validator.Between(1, 7),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between 1 SECONDS and 30 MINUTES.").Description,
												Required:    true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(1, 1800),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																"If `time_unit` is `SECONDS`, the allowed duration range is 1 - 1800.",
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(1, 30),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																"If `time_unit` is `MINUTES`, the allowed duration range is 1 - 30.",
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
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that, when set to `true`, prevents users from pairing new Yubikey devices.").Description,
						Optional:    true,
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies pairing key lifetime settings for Yubikey devices.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the amount of time an issued pairing key can be used until it expires.").Description,
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

					"prompt_for_nickname_on_pairing": schema.BoolAttribute{
						Description:         promptForNicknameOnPairingDescription.Description,
						MarkdownDescription: promptForNicknameOnPairingDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"oath_token": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of OATH token device authentication policy settings.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(oathTokenDefault),

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the OATH token device method is enabled or disabled in the policy.").Description,
						Optional:    true,
						Computed:    true,

						Default: booldefault.StaticBool(false),
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure settings for OATH token devices.").Description,
						Optional:    true,
						Computed:    true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							MFADevicePolicyPingIDDeviceOtpTFObjectTypes,
							map[string]attr.Value{
								"failure": types.ObjectValueMust(
									MFADevicePolicyFailureTFObjectTypes,
									map[string]attr.Value{
										"count": types.Int32Value(3),
										"cool_down": types.ObjectValueMust(
											MFADevicePolicyTimePeriodTFObjectTypes,
											map[string]attr.Value{
												"duration":  types.Int32Value(2),
												"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
											},
										),
									},
								),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between 1 and 7.").Description,
										Optional:    true,

										Validators: []validator.Int32{
											int32validator.Between(1, 7),
										},
									},

									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between 1 SECONDS and 30 MINUTES.").Description,
												Required:    true,

												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(1, 1800),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																"If `time_unit` is `SECONDS`, the allowed duration range is 1 - 1800.",
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(1, 30),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																"If `time_unit` is `MINUTES`, the allowed duration range is 1 - 30.",
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
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that, when set to `true`, prevents users from pairing new OATH token devices.").Description,
						Optional:    true,
						Computed:    true,

						Default: booldefault.StaticBool(false),
					},

					"pairing_key_lifetime": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies pairing key lifetime settings for OATH token devices.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the amount of time an issued pairing key can be used until it expires.").Description,
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

func (r *MFADevicePolicyDefaultResource) devicePolicyOfflineDeviceSchemaAttribute(descriptionMethod string) schema.SingleNestedAttribute {
	// Reuse the implementation from MFADevicePolicyResource
	mfaDevicePolicyResource := &MFADevicePolicyResource{}
	return mfaDevicePolicyResource.devicePolicyOfflineDeviceSchemaAttribute(descriptionMethod)
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
	var plan, state MFADevicePolicyDefaultResourceModel

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

	// Run the API call to check if default exists
	readResponse, d := FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API ensures a default policy always exists, so if we can't find it, something is wrong
	if readResponse == nil {
		resp.Diagnostics.AddError(
			"Default MFA Device Policy Not Found",
			"Cannot find the default MFA device policy for the environment.",
		)
		return
	}

	// Extract ID from the union type based on policy_type
	policyID, err := extractPolicyIDFromUnion(readResponse, plan.PolicyType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid policy response",
			err.Error(),
		)
		return
	}

	// Update the default policy
	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
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
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
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

	// Run the API call
	response, d := FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), true)
	resp.Diagnostics.Append(d...)
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

func (r *MFADevicePolicyDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MFADevicePolicyDefaultResourceModel

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

	// Run the API call to get the default policy ID
	readResponse, d := FetchDefaultMFADevicePolicy(ctx, r.Client.MFAAPIClient, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract ID from the union type based on policy_type
	policyID, err := extractPolicyIDFromUnion(readResponse, plan.PolicyType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid policy response",
			err.Error(),
		)
		return
	}

	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
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
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
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

	// Note: We don't actually delete the default policy or modify it.
	// We simply remove it from Terraform's state and leave it as-is in PingOne.
	// The API prevents deletion of the default policy, so this is the expected behavior.

	resp.Diagnostics.AddWarning(
		"State change warning",
		"The \"pingone_mfa_device_policy_default\" resource has been destroyed.  The default MFA device policy has been removed from Terraform's state.  The policy itself has not been removed from the PingOne service.",
	)
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

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

// extractPolicyIDFromUnion extracts the policy ID from a DeviceAuthenticationPolicy
// Note: With the flattened SDK model, there are no union fields, so we just return the ID directly
func extractPolicyIDFromUnion(response *mfa.DeviceAuthenticationPolicy, policyType string) (string, error) {
	// With flattened SDK, all policies use the same struct
	// We can determine policy type by checking for PingID-specific fields (Desktop, Yubikey, OathToken)
	// but for ID extraction we just return the ID directly
	if response == nil {
		return "", fmt.Errorf("response is nil")
	}
	return response.GetId(), nil
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

func (p *MFADevicePolicyDefaultResourceModel) expand(ctx context.Context, apiClient *management.APIClient) (mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
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
	var mobilePlan MFADevicePolicyMobileResourceModel
	diags.Append(p.Mobile.As(ctx, &mobilePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return mfa.DeviceAuthenticationPolicy{}, diags
	}
	mobile, d := mobilePlan.expand(ctx, apiClient, p.EnvironmentId.ValueString())
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
			mfa.DeviceAuthenticationPolicyPingIDDevicePairingKeyLifetime{
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

func (p *MFADevicePolicyDefaultResourceModel) toState(apiObject *mfa.DeviceAuthenticationPolicy) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	// With flattened SDK, determine policy type by checking for PingID-specific fields
	// If Desktop, Yubikey are present (and not null), it's a PingID policy
	// OathToken can be on both policy types, so we check Desktop/Yubikey first
	isPingID := false
	if desktop, ok := apiObject.GetDesktopOk(); ok && desktop != nil {
		isPingID = true
	} else if yubikey, ok := apiObject.GetYubikeyOk(); ok && yubikey != nil {
		isPingID = true
	}

	// Common fields for both policy types
	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	// Set policy type based on detection
	if isPingID {
		p.PolicyType = types.StringValue(POLICY_TYPE_PINGID)
	} else {
		p.PolicyType = types.StringValue(POLICY_TYPE_PINGONE_MFA)
	}

	// Common fields
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

	p.OathToken, d = toStateMfaDevicePolicyOathToken(apiObject.GetOathTokenOk())
	diags.Append(d...)

	// Policy type specific fields
	if isPingID {
		// PingID-specific devices
		p.Desktop, d = toStateMfaDevicePolicyPingIDDevice(apiObject.GetDesktopOk())
		diags.Append(d...)

		p.Yubikey, d = toStateMfaDevicePolicyPingIDDevice(apiObject.GetYubikeyOk())
		diags.Append(d...)

		// Set PingOneMFA-specific field to null for PingID policies
		p.Fido2 = types.ObjectNull(MFADevicePolicyFido2TFObjectTypes)
	} else {
		// PingOneMFA-specific field
		p.Fido2, d = toStateMfaDevicePolicyFido2(apiObject.GetFido2Ok())
		diags.Append(d...)

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
		"enabled":                        types.BoolValue(apiObject.GetEnabled()),
		"otp":                            otp,
		"pairing_disabled":               framework.BoolOkToTF(apiObject.GetPairingDisabledOk()),
		"pairing_key_lifetime":           pairingKeyLifetime,
		"prompt_for_nickname_on_pairing": framework.BoolOkToTF(apiObject.GetPromptForNicknameOnPairingOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func toStateMfaDevicePolicyPingIDMobile(apiObject *mfa.DeviceAuthenticationPolicyCommonMobile, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyMobileTFObjectTypes), diags
	}

	// PingID mobile is much simpler than PingOneMFA mobile - only has enabled field
	// Set all PingOneMFA-specific fields to null
	objValue, d := types.ObjectValue(MFADevicePolicyMobileTFObjectTypes, map[string]attr.Value{
		"enabled":      types.BoolValue(apiObject.GetEnabled()),
		"applications": types.ObjectNull(MFADevicePolicyMobileApplicationTFObjectTypes),
		"otp":          types.ObjectNull(MFADevicePolicyMobileOtpTFObjectTypes),
	})
	diags.Append(d...)

	return objValue, diags
}
