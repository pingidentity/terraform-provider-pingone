package mfa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type MFAPolicyResource struct {
	client     *mfa.APIClient
	mgmtClient *management.APIClient
	region     model.RegionMapping
}

type mfaPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	Name            types.String `tfsdk:"name"`
	DeviceSelection types.String `tfsdk:"device_selection"`
	SMS             types.List   `tfsdk:"sms"`
	Voice           types.List   `tfsdk:"voice"`
	Email           types.List   `tfsdk:"email"`
	Mobile          types.List   `tfsdk:"mobile"`
	Totp            types.List   `tfsdk:"totp"`
	SecurityKey     types.List   `tfsdk:"security_key"`
	Platform        types.List   `tfsdk:"platform"`
}

type mfaPolicyResourceOfflineDeviceModel struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	OTPLifetimeDuration        types.Int64  `tfsdk:"otp_lifetime_duration"`
	OTPLifetimeTimeunit        types.String `tfsdk:"otp_lifetime_timeunit"`
	OTPFailureCount            types.Int64  `tfsdk:"otp_failure_count"`
	OTPFailureCooldownDuration types.Int64  `tfsdk:"otp_failure_cooldown_duration"`
	OTPFailureCooldownTimeunit types.String `tfsdk:"otp_failure_cooldown_timeunit"`
}

type mfaPolicyResourceMobileModel struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	OTPFailureCount            types.Int64  `tfsdk:"otp_failure_count"`
	OTPFailureCooldownDuration types.Int64  `tfsdk:"otp_failure_cooldown_duration"`
	OTPFailureCooldownTimeunit types.String `tfsdk:"otp_failure_cooldown_timeunit"`
	Application                types.Set    `tfsdk:"application"`
}

type mfaPolicyResourceMobileApplicationModel struct {
	Id                                   types.String `tfsdk:"id"`
	PushEnabled                          types.Bool   `tfsdk:"push_enabled"`
	PushTimeoutDuration                  types.Int64  `tfsdk:"push_timeout_duration"`
	PushTimeoutTimeunit                  types.String `tfsdk:"push_timeout_timeunit"`
	OTPEnabled                           types.Bool   `tfsdk:"otp_enabled"`
	DeviceAuthorizationEnabled           types.Bool   `tfsdk:"device_authorization_enabled"`
	DeviceAuthorizationExtraVerification types.String `tfsdk:"device_authorization_extra_verification"`
	AutoEnrollmentEnabled                types.Bool   `tfsdk:"auto_enrollment_enabled"`
	IntegrityDetection                   types.String `tfsdk:"integrity_detection"`
	PairingKeyLifetimeDuration           types.Int64  `tfsdk:"pairing_key_lifetime_duration"`
	PairingKeyLifetimeTimeunit           types.String `tfsdk:"pairing_key_lifetime_timeunit"`
	PushLimit                            types.List   `tfsdk:"push_limit"`
}

type mfaPolicyResourceMobileApplicationPushLimitModel struct {
	Count                types.Int64  `tfsdk:"count"`
	LockDurationDuration types.Int64  `tfsdk:"lock_duration_duration"`
	LockDurationTimeunit types.String `tfsdk:"lock_duration_timeunit"`
	TimePeriodDuration   types.Int64  `tfsdk:"time_period_duration"`
	TimePeriodTimeunit   types.String `tfsdk:"time_period_timeunit"`
}

type mfaPolicyResourceTotpModel struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	OTPFailureCount            types.Int64  `tfsdk:"otp_failure_count"`
	OTPFailureCooldownDuration types.Int64  `tfsdk:"otp_failure_cooldown_duration"`
	OTPFailureCooldownTimeunit types.String `tfsdk:"otp_failure_cooldown_timeunit"`
}

type mfaPolicyResourceFidoDeviceModel struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	FIDOPolicyID types.String `tfsdk:"fido_policy_id"`
}

var (
	mfaPolicyOfflineDeviceTFObjectTypes = map[string]attr.Type{
		"enabled":                       types.BoolType,
		"otp_lifetime_duration":         types.Int64Type,
		"otp_lifetime_timeunit":         types.StringType,
		"otp_failure_count":             types.Int64Type,
		"otp_failure_cooldown_duration": types.Int64Type,
		"otp_failure_cooldown_timeunit": types.StringType,
	}

	mfaPolicyMobileTFObjectTypes = map[string]attr.Type{
		"enabled":                       types.BoolType,
		"otp_failure_count":             types.Int64Type,
		"otp_failure_cooldown_duration": types.Int64Type,
		"otp_failure_cooldown_timeunit": types.StringType,
		"application":                   types.SetType{ElemType: types.ObjectType{AttrTypes: mfaPolicyMobileApplicationTFObjectTypes}},
	}

	mfaPolicyMobileApplicationTFObjectTypes = map[string]attr.Type{
		"id":                           types.StringType,
		"push_enabled":                 types.BoolType,
		"push_timeout_duration":        types.Int64Type,
		"push_timeout_timeunit":        types.StringType,
		"otp_enabled":                  types.BoolType,
		"device_authorization_enabled": types.BoolType,
		"device_authorization_extra_verification": types.StringType,
		"auto_enrollment_enabled":                 types.BoolType,
		"integrity_detection":                     types.StringType,
		"pairing_key_lifetime_duration":           types.Int64Type,
		"pairing_key_lifetime_timeunit":           types.StringType,
		"push_limit":                              types.ListType{ElemType: types.ObjectType{AttrTypes: mfaPolicyMobileApplicationPushLimitTFObjectTypes}},
	}

	mfaPolicyMobileApplicationPushLimitTFObjectTypes = map[string]attr.Type{
		"count":                  types.Int64Type,
		"lock_duration_duration": types.Int64Type,
		"lock_duration_timeunit": types.StringType,
		"time_period_duration":   types.Int64Type,
		"time_period_timeunit":   types.StringType,
	}

	mfaPolicyTotpTFObjectTypes = map[string]attr.Type{
		"enabled":                       types.BoolType,
		"otp_failure_count":             types.Int64Type,
		"otp_failure_cooldown_duration": types.Int64Type,
		"otp_failure_cooldown_timeunit": types.StringType,
	}

	mfaPolicyFidoDeviceTFObjectTypes = map[string]attr.Type{
		"enabled":        types.BoolType,
		"fido_policy_id": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &MFAPolicyResource{}
	_ resource.ResourceWithConfigure   = &MFAPolicyResource{}
	_ resource.ResourceWithImportState = &MFAPolicyResource{}
)

// New Object
func NewMFAPolicyResource() resource.Resource {
	return &MFAPolicyResource{}
}

// Metadata
func (r *MFAPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_policy"
}

// Schema.
func (r *MFAPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	const mobileOtpFailureCountDefault = 3

	const mobileOtpFailureCooldownDurationDefault = 2

	const totpOtpFailureCountDefault = 3

	const totpOtpFailureCooldownDurationDefault = 2

	const mobileApplicationPushTimeoutDurationDefault = 40
	const mobileApplicationPushTimeoutDurationMin = 40
	const mobileApplicationPushTimeoutDurationMax = 150

	const mobileApplicationPairingKeyLifetimeDurationDefault = 10
	const mobileApplicationPairingKeyLifetimeDurationMin = 1

	const mobileApplicationPushLimitCountDefault = 5
	const mobileApplicationPushLimitCountMin = 1
	const mobileApplicationPushLimitCountMax = 50

	const mobileApplicationPushLimitLockDurationDurationDefault = 30

	const mobileApplicationPushLimitTimePeriodDurationDefault = 10

	deviceSelectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines the device selection method.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicySelectionEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST))

	authenticatorPaths := []string{"sms", "voice", "email", "mobile", "totp", "security_key", "platform"}

	// SMS
	smsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"SMS OTP authentication policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	// Voice
	voiceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Voice OTP authentication policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	// Email
	emailDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Email OTP authentication policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	// Mobile
	mobileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Mobile authenticator device policy settings.  This factor requires embedding the PingOne MFA SDK into a customer facing mobile application, and configuring as a Native application using the `pingone_application` resource.",
	).ExactlyOneOf(authenticatorPaths)

	mobileOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
	).DefaultValue(fmt.Sprintf("%d", mobileOtpFailureCountDefault))

	mobileOtpFailureCooldownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
	).DefaultValue(fmt.Sprintf("%d", mobileOtpFailureCooldownDurationDefault))

	mobileOtpFailureCooldownTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of time unit for `otp_failure_cooldown_duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	mobileApplicationPushTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the amount of time (in seconds) a user has to respond to a push notification before it expires. Minimum is `40` seconds and maximum is `150` seconds.",
	).DefaultValue(fmt.Sprintf("%d", mobileApplicationPushTimeoutDurationDefault))

	mobileApplicationPushTimoutTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The time unit for the `push_timeout_duration` parameter. Currently, the only permitted value is `SECONDS`.",
	)

	mobileApplicationDeviceAuthorizationExtraVerificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the level of further verification when `device_authorization_enabled` is true. The PingOne platform performs an extra verification check by sending a \"silent\" push notification to the customer native application, and receives a confirmation in return.",
	).AllowedValuesComplex(map[string]string{
		"permissive":  "the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as successfully completed",
		"restrictive": "the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as failed",
	})

	mobileApplicationAutoEnrollmentEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Set to `true` if you want the application to allow Auto Enrollment. Auto Enrollment means that the user can authenticate for the first time from an unpaired device, and the successful authentication will result in the pairing of the device for MFA.",
	)

	mobileApplicationIntegrityDetectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls how authentication or registration attempts should proceed if a device integrity check does not receive a response.",
	).AllowedValuesComplex(map[string]string{
		"permissive":  "if you want to allow the process to continue",
		"restrictive": "if you want to block the user in such situations",
	})

	mobileApplicationPairingKeyLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The amount of time an issued pairing key can be used until it expires. Minimum is `1` minute and maximum is `48` hours.",
	).DefaultValue(fmt.Sprintf("%d", mobileApplicationPairingKeyLifetimeDurationDefault))

	mobileApplicationPairingKeyLifetimeTimeunit := framework.SchemaAttributeDescriptionFromMarkdown(
		"The time unit for the `pairing_key_lifetime_duration` parameter.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues).DefaultValue(string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_MINUTES))

	mobileApplicationPushLimitCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is `1` and the maximum value is `50`.",
	).DefaultValue(fmt.Sprintf("%d", mobileApplicationPushLimitCountDefault))

	mobileApplicationPushLimitLockDurationDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes.",
	).DefaultValue(fmt.Sprintf("%d", mobileApplicationPushLimitLockDurationDurationDefault))

	mobileApplicationPushLimitLockDurationTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The time unit for the `lock_duration_duration` parameter.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	mobileApplicationPushLimitTimePeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The time period in which the push notifications are counted towards the defined limit. The minimum value is `1` minute and the maximum value is `120` minutes.",
	).DefaultValue(fmt.Sprintf("%d", mobileApplicationPushLimitTimePeriodDurationDefault))

	mobileApplicationPushLimitTimePeriodTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The time unit for the `time_period_duration` parameter.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	// TOTP
	totpDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"TOTP authenticator policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	totpOtpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
	).DefaultValue(fmt.Sprintf("%d", totpOtpFailureCountDefault))

	totpOtpFailureCooldownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
	).DefaultValue(fmt.Sprintf("%d", totpOtpFailureCooldownDurationDefault))

	totpOtpFailureCooldownTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of time unit for `otp_failure_cooldown_duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	// Security Key
	securityKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Security key (FIDO2) authentication policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	// Platform
	platformDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Platform biometrics authentication policy settings.",
	).ExactlyOneOf(authenticatorPaths)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA Policies in a PingOne Environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the MFA policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the MFA policy's name.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

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

		Blocks: map[string]schema.Block{

			"sms": schema.ListNestedBlock{
				Description:         smsDescription.Description,
				MarkdownDescription: smsDescription.MarkdownDescription,

				NestedObject: offlineDeviceResourceSchema(),

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"voice": schema.ListNestedBlock{
				Description:         voiceDescription.Description,
				MarkdownDescription: voiceDescription.MarkdownDescription,

				NestedObject: offlineDeviceResourceSchema(),

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"email": schema.ListNestedBlock{
				Description:         emailDescription.Description,
				MarkdownDescription: emailDescription.MarkdownDescription,

				NestedObject: offlineDeviceResourceSchema(),

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"mobile": schema.ListNestedBlock{
				Description:         mobileDescription.Description,
				MarkdownDescription: mobileDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Enabled or disabled in the policy.").Description,
							Required:    true,
						},

						"otp_failure_count": schema.Int64Attribute{
							Description:         mobileOtpFailureCountDescription.Description,
							MarkdownDescription: mobileOtpFailureCountDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(int64(mobileOtpFailureCountDefault)),
						},

						"otp_failure_cooldown_duration": schema.Int64Attribute{
							Description:         mobileOtpFailureCooldownDurationDescription.Description,
							MarkdownDescription: mobileOtpFailureCooldownDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(int64(mobileOtpFailureCooldownDurationDefault)),
						},

						"otp_failure_cooldown_timeunit": schema.StringAttribute{
							Description:         mobileOtpFailureCooldownTimeunitDescription.Description,
							MarkdownDescription: mobileOtpFailureCooldownTimeunitDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
								stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("otp_failure_cooldown_duration")),
							},
						},
					},

					Blocks: map[string]schema.Block{
						"application": schema.SetNestedBlock{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Settings for a configured Mobile Application.").Description,

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": framework.Attr_LinkID(
										framework.SchemaAttributeDescriptionFromMarkdown("The mobile application's ID.  Mobile applications are configured with the `pingone_application` resource, as an OIDC `NATIVE` type."),
									),

									"push_enabled": schema.BoolAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies whether push notification is enabled or disabled for the policy.").Description,
										Required:    true,
									},

									"push_timeout_duration": schema.Int64Attribute{
										Description:         mobileApplicationPushTimeoutDurationDescription.Description,
										MarkdownDescription: mobileApplicationPushTimeoutDurationDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int64default.StaticInt64(int64(mobileApplicationPushTimeoutDurationDefault)),

										Validators: []validator.Int64{
											int64validator.AtLeast(int64(mobileApplicationPushTimeoutDurationMin)),
											int64validator.AtMost(int64(mobileApplicationPushTimeoutDurationMax)),
										},
									},

									"push_timeout_timeunit": schema.StringAttribute{
										Description:         mobileApplicationPushTimoutTimeunitDescription.Description,
										MarkdownDescription: mobileApplicationPushTimoutTimeunitDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_SECONDS)),

										Validators: []validator.String{
											stringvalidator.OneOf(string(mfa.ENUMTIMEUNIT_SECONDS)),
											stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("push_timeout_duration")),
										},
									},

									"otp_enabled": schema.BoolAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies whether OTP authentication is enabled or disabled for the policy.").Description,
										Required:    true,
									},

									"device_authorization_enabled": schema.BoolAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies the enabled or disabled state of automatic MFA for native devices paired with the user, for the specified application.").Description,
										Optional:    true,
									},

									"device_authorization_extra_verification": schema.StringAttribute{
										Description:         mobileApplicationDeviceAuthorizationExtraVerificationDescription.Description,
										MarkdownDescription: mobileApplicationDeviceAuthorizationExtraVerificationDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.String{
											stringvalidator.OneOf("permissive", "restrictive"),
										},
									},

									"auto_enrollment_enabled": schema.BoolAttribute{
										Description:         mobileApplicationAutoEnrollmentEnabledDescription.Description,
										MarkdownDescription: mobileApplicationAutoEnrollmentEnabledDescription.MarkdownDescription,
										Optional:            true,
									},

									"integrity_detection": schema.StringAttribute{
										Description:         mobileApplicationIntegrityDetectionDescription.Description,
										MarkdownDescription: mobileApplicationIntegrityDetectionDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.String{
											stringvalidator.OneOf("permissive", "restrictive"),
										},
									},

									"pairing_key_lifetime_duration": schema.Int64Attribute{
										Description:         mobileApplicationPairingKeyLifetimeDurationDescription.Description,
										MarkdownDescription: mobileApplicationPairingKeyLifetimeDurationDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: int64default.StaticInt64(int64(mobileApplicationPairingKeyLifetimeDurationDefault)),

										Validators: []validator.Int64{
											int64validator.AtLeast(int64(mobileApplicationPairingKeyLifetimeDurationMin)),
										},
									},

									"pairing_key_lifetime_timeunit": schema.StringAttribute{
										Description:         mobileApplicationPairingKeyLifetimeTimeunit.Description,
										MarkdownDescription: mobileApplicationPairingKeyLifetimeTimeunit.MarkdownDescription,
										Optional:            true,
										Computed:            true,

										Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_MINUTES)),

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitPairingKeyLifetimeEnumValues)...),
											stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("pairing_key_lifetime_duration")),
										},
									},
								},

								Blocks: map[string]schema.Block{
									"push_limit": schema.ListNestedBlock{
										Description: "A single block that describes mobile application push limit settings.",

										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"count": schema.Int64Attribute{
													Description:         mobileApplicationPushLimitCountDescription.Description,
													MarkdownDescription: mobileApplicationPushLimitCountDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: int64default.StaticInt64(int64(mobileApplicationPushLimitCountDefault)),

													Validators: []validator.Int64{
														int64validator.AtLeast(int64(mobileApplicationPushLimitCountMin)),
														int64validator.AtMost(int64(mobileApplicationPushLimitCountMax)),
													},
												},

												"lock_duration_duration": schema.Int64Attribute{
													Description:         mobileApplicationPushLimitLockDurationDurationDescription.Description,
													MarkdownDescription: mobileApplicationPushLimitLockDurationDurationDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: int64default.StaticInt64(int64(mobileApplicationPushLimitLockDurationDurationDefault)),
												},

												"lock_duration_timeunit": schema.StringAttribute{
													Description:         mobileApplicationPushLimitLockDurationTimeunitDescription.Description,
													MarkdownDescription: mobileApplicationPushLimitLockDurationTimeunitDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

													Validators: []validator.String{
														stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
														stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("lock_duration_duration")),
													},
												},

												"time_period_duration": schema.Int64Attribute{
													Description:         mobileApplicationPushLimitTimePeriodDurationDescription.Description,
													MarkdownDescription: mobileApplicationPushLimitTimePeriodDurationDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: int64default.StaticInt64(int64(mobileApplicationPushLimitTimePeriodDurationDefault)),
												},

												"time_period_timeunit": schema.StringAttribute{
													Description:         mobileApplicationPushLimitTimePeriodTimeunitDescription.Description,
													MarkdownDescription: mobileApplicationPushLimitTimePeriodTimeunitDescription.MarkdownDescription,
													Optional:            true,
													Computed:            true,

													Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

													Validators: []validator.String{
														stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
														stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("time_period_duration")),
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
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"totp": schema.ListNestedBlock{
				Description:         totpDescription.Description,
				MarkdownDescription: totpDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Enabled or disabled in the policy.").Description,
							Required:    true,
						},

						"otp_failure_count": schema.Int64Attribute{
							Description:         totpOtpFailureCountDescription.Description,
							MarkdownDescription: totpOtpFailureCountDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(int64(totpOtpFailureCountDefault)),
						},

						"otp_failure_cooldown_duration": schema.Int64Attribute{
							Description:         totpOtpFailureCooldownDurationDescription.Description,
							MarkdownDescription: totpOtpFailureCooldownDurationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: int64default.StaticInt64(int64(totpOtpFailureCooldownDurationDefault)),
						},

						"otp_failure_cooldown_timeunit": schema.StringAttribute{
							Description:         totpOtpFailureCooldownTimeunitDescription.Description,
							MarkdownDescription: totpOtpFailureCooldownTimeunitDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
								stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("otp_failure_cooldown_duration")),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"security_key": schema.ListNestedBlock{
				Description:         securityKeyDescription.Description,
				MarkdownDescription: securityKeyDescription.MarkdownDescription,

				NestedObject: fidoDeviceResourceSchema(),

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},

			"platform": schema.ListNestedBlock{
				Description:         platformDescription.Description,
				MarkdownDescription: platformDescription.MarkdownDescription,

				NestedObject: fidoDeviceResourceSchema(),

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.AtLeastOneOf(
						path.MatchRoot("sms"),
						path.MatchRoot("voice"),
						path.MatchRoot("email"),
						path.MatchRoot("mobile"),
						path.MatchRoot("totp"),
						path.MatchRoot("security_key"),
						path.MatchRoot("platform"),
					),
				},
			},
		},
	}
}

func offlineDeviceResourceSchema() schema.NestedBlockObject {

	const otpLifetimeDefault = 30
	const otpFailureDefault = 3
	const otpFailureCooldownDefault = 0

	otpLifetimeDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines duration (number of time units) that the passcode is valid before it expires.",
	).DefaultValue(fmt.Sprintf("%d", otpLifetimeDefault))

	otpLifetimeTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of time unit for `otp_lifetime_duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	otpFailureCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
	).DefaultValue(fmt.Sprintf("%d", otpFailureDefault))

	otpFailureCooldownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Note that when using the \"onetime authentication\" feature, the user is not blocked after the maximum number of failures even if you specified a block duration.",
	).DefaultValue(fmt.Sprintf("%d", otpFailureCooldownDefault))

	otpFailureCooldownTimeunitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of time unit for `otp_failure_cooldown_duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Enabled or disabled in the policy.").Description,
				Required:    true,
			},

			"otp_lifetime_duration": schema.Int64Attribute{
				Description:         otpLifetimeDurationDescription.Description,
				MarkdownDescription: otpLifetimeDurationDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(int64(otpLifetimeDefault)),
			},

			"otp_lifetime_timeunit": schema.StringAttribute{
				Description:         otpLifetimeTimeunitDescription.Description,
				MarkdownDescription: otpLifetimeTimeunitDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("otp_lifetime_duration")),
				},
			},

			"otp_failure_count": schema.Int64Attribute{
				Description:         otpFailureCountDescription.Description,
				MarkdownDescription: otpFailureCountDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(int64(otpFailureDefault)),
			},

			"otp_failure_cooldown_duration": schema.Int64Attribute{
				Description:         otpFailureCooldownDurationDescription.Description,
				MarkdownDescription: otpFailureCooldownDurationDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(int64(otpFailureCooldownDefault)),
			},

			"otp_failure_cooldown_timeunit": schema.StringAttribute{
				Description:         otpFailureCooldownTimeunitDescription.Description,
				MarkdownDescription: otpFailureCooldownTimeunitDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("otp_failure_cooldown_duration")),
				},
			},
		},
	}
}

func fidoDeviceResourceSchema() schema.NestedBlockObject {

	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Enabled or disabled in the policy.").Description,
				Required:    true,
			},

			"fido_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Specifies the FIDO policy ID. This property can be null. When null, the environment's default FIDO Policy is used.").Description,
				Optional:    true,

				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},
		},
	}
}

func (r *MFAPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region

	preparedMgmtClient, err := prepareMgmtClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.mgmtClient = preparedMgmtClient
}

func (r *MFAPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state mfaPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	mfaPolicy, d := plan.expand(ctx, r.mgmtClient, plan.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DeviceAuthenticationPolicyApi.CreateDeviceAuthenticationPolicies(ctx, plan.EnvironmentId.ValueString()).DeviceAuthenticationPolicy(*mfaPolicy).Execute()
		},
		"CreateDeviceAuthenticationPolicies",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*mfa.DeviceAuthenticationPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFAPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *mfaPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneDeviceAuthenticationPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*mfa.DeviceAuthenticationPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MFAPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mfaPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	mfaPolicy, d := plan.expand(ctx, r.mgmtClient, plan.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).DeviceAuthenticationPolicy(*mfaPolicy).Execute()
		},
		"UpdateDeviceAuthenticationPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*mfa.DeviceAuthenticationPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFAPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *mfaPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteDeviceAuthenticationPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MFAPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/mfa_device_policy_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *mfaPolicyResourceModel) expand(ctx context.Context, apiClient *management.APIClient, environmentID string) (*mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// SMS
	var smsPlan []mfaPolicyResourceOfflineDeviceModel
	diags.Append(p.SMS.ElementsAs(ctx, &smsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataSms, d := smsPlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Voice
	var voicePlan []mfaPolicyResourceOfflineDeviceModel
	diags.Append(p.Voice.ElementsAs(ctx, &voicePlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataVoice, d := voicePlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Email
	var emailPlan []mfaPolicyResourceOfflineDeviceModel
	diags.Append(p.Email.ElementsAs(ctx, &emailPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataEmail, d := emailPlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Mobile
	var mobilePlan []mfaPolicyResourceMobileModel
	diags.Append(p.Mobile.ElementsAs(ctx, &mobilePlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataMobile, d := mobilePlan[0].expand(ctx, apiClient, environmentID)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// TOTP
	var totpPlan []mfaPolicyResourceTotpModel
	diags.Append(p.Totp.ElementsAs(ctx, &totpPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataTotp, d := totpPlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Security Key
	var securityKeyPlan []mfaPolicyResourceFidoDeviceModel
	diags.Append(p.SecurityKey.ElementsAs(ctx, &securityKeyPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataSecurityKey, d := securityKeyPlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Platform
	var platformPlan []mfaPolicyResourceFidoDeviceModel
	diags.Append(p.Platform.ElementsAs(ctx, &platformPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	dataPlatform, d := platformPlan[0].expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := mfa.NewDeviceAuthenticationPolicy(
		p.Name.ValueString(),
		*dataSms,
		*dataVoice,
		*dataEmail,
		*dataMobile,
		*dataTotp,
		*dataSecurityKey,
		*dataPlatform,
		false,
		false,
	)

	if !p.DeviceSelection.IsNull() && !p.DeviceSelection.IsUnknown() {
		data.SetAuthentication(
			*mfa.NewDeviceAuthenticationPolicyAuthentication(
				mfa.EnumMFADevicePolicySelection(p.DeviceSelection.ValueString()),
			),
		)
	}

	return data, diags
}

func (p *mfaPolicyResourceOfflineDeviceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyOfflineDevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := mfa.NewDeviceAuthenticationPolicyOfflineDevice(p.Enabled.ValueBool(),
		*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpLifeTime(
				int32(p.OTPLifetimeDuration.ValueInt64()),
				mfa.EnumTimeUnit(p.OTPLifetimeTimeunit.ValueString()),
			),
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(p.OTPFailureCount.ValueInt64()),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(
					int32(p.OTPFailureCooldownDuration.ValueInt64()),
					mfa.EnumTimeUnit(p.OTPFailureCooldownTimeunit.ValueString()),
				),
			),
		),
	)

	return data, diags
}

func (p *mfaPolicyResourceMobileModel) expand(ctx context.Context, apiClient *management.APIClient, environmentID string) (*mfa.DeviceAuthenticationPolicyMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	otpStepSizeDuration := 30

	data := mfa.NewDeviceAuthenticationPolicyMobile(p.Enabled.ValueBool(),
		*mfa.NewDeviceAuthenticationPolicyMobileOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(p.OTPFailureCount.ValueInt64()),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(
					int32(p.OTPFailureCooldownDuration.ValueInt64()),
					mfa.EnumTimeUnit(p.OTPFailureCooldownTimeunit.ValueString()),
				),
			),
			*mfa.NewDeviceAuthenticationPolicyMobileOtpWindow(
				*mfa.NewDeviceAuthenticationPolicyMobileOtpWindowStepSize(
					int32(otpStepSizeDuration),
					mfa.ENUMTIMEUNIT_SECONDS,
				),
			),
		),
	)

	if !p.Application.IsNull() && !p.Application.IsUnknown() {
		var applicationsPlan []mfaPolicyResourceMobileApplicationModel
		diags.Append(p.Application.ElementsAs(ctx, &applicationsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		applications := make([]mfa.DeviceAuthenticationPolicyMobileApplicationsInner, 0)

		for _, applicationPlan := range applicationsPlan {

			item := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInner(applicationPlan.Id.ValueString())

			diags.Append(checkApplicationForMobileApp(ctx, apiClient, environmentID, applicationPlan.Id.ValueString())...)
			if diags.HasError() {
				return nil, diags
			}

			if !applicationPlan.PushEnabled.IsNull() && !applicationPlan.PushEnabled.IsUnknown() {
				item.SetPush(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPush(applicationPlan.PushEnabled.ValueBool()))
			}

			if !applicationPlan.PushTimeoutDuration.IsNull() && !applicationPlan.PushTimeoutDuration.IsUnknown() {
				item.SetPushTimeout(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushTimeout(
					int32(applicationPlan.PushTimeoutDuration.ValueInt64()),
					mfa.EnumTimeUnitPushTimeout(applicationPlan.PushTimeoutTimeunit.ValueString())),
				)
			}

			if !applicationPlan.OTPEnabled.IsNull() && !applicationPlan.OTPEnabled.IsUnknown() {
				item.SetOtp(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerOtp(applicationPlan.OTPEnabled.ValueBool()))
			}

			deviceAuthz := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerDeviceAuthorization(applicationPlan.DeviceAuthorizationEnabled.ValueBool())

			if !applicationPlan.DeviceAuthorizationExtraVerification.IsNull() && !applicationPlan.DeviceAuthorizationExtraVerification.IsUnknown() {
				deviceAuthz.SetExtraVerification(mfa.EnumMFADevicePolicyMobileExtraVerification(applicationPlan.DeviceAuthorizationExtraVerification.ValueString()))
			}

			item.SetDeviceAuthorization(deviceAuthz)

			if !applicationPlan.AutoEnrollmentEnabled.IsNull() && !applicationPlan.AutoEnrollmentEnabled.IsUnknown() {
				item.SetAutoEnrollment(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerAutoEnrollment(applicationPlan.AutoEnrollmentEnabled.ValueBool()))
			}

			if !applicationPlan.IntegrityDetection.IsNull() && !applicationPlan.IntegrityDetection.IsUnknown() {
				item.SetIntegrityDetection(mfa.EnumMFADevicePolicyMobileIntegrityDetection(applicationPlan.IntegrityDetection.ValueString()))
			}

			applications = append(applications, item)
		}

		data.SetApplications(applications)
	}

	return data, diags
}

func (p *mfaPolicyResourceTotpModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyTotp, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := mfa.NewDeviceAuthenticationPolicyTotp(p.Enabled.ValueBool(),
		*mfa.NewDeviceAuthenticationPolicyTotpOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(p.OTPFailureCount.ValueInt64()),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(
					int32(p.OTPFailureCooldownDuration.ValueInt64()),
					mfa.EnumTimeUnit(p.OTPFailureCooldownTimeunit.ValueString()),
				),
			),
		),
	)

	return data, diags
}

func (p *mfaPolicyResourceFidoDeviceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyFIDODevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := mfa.NewDeviceAuthenticationPolicyFIDODevice(p.Enabled.ValueBool())

	if !p.FIDOPolicyID.IsNull() && !p.FIDOPolicyID.IsUnknown() {
		data.SetFidoPolicyId(p.FIDOPolicyID.ValueString())
	}

	return data, diags
}

func (p *mfaPolicyResourceModel) toState(apiObject *mfa.DeviceAuthenticationPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	if v, ok := apiObject.GetAuthenticationOk(); ok {
		p.DeviceSelection = framework.EnumOkToTF(v.GetDeviceSelectionOk())
	} else {
		p.DeviceSelection = types.StringNull()
	}

	var d diag.Diagnostics
	p.SMS, d = offlineDeviceOkToTF(apiObject.GetSmsOk())
	diags.Append(d...)

	p.Voice, d = offlineDeviceOkToTF(apiObject.GetVoiceOk())
	diags.Append(d...)

	p.Email, d = offlineDeviceOkToTF(apiObject.GetEmailOk())
	diags.Append(d...)

	p.Mobile, d = mobileDeviceOkToTF(apiObject.GetMobileOk())
	diags.Append(d...)

	p.Totp, d = totpDeviceOkToTF(apiObject.GetTotpOk())
	diags.Append(d...)

	p.SecurityKey, d = fidoDeviceOkToTF(apiObject.GetSecurityKeyOk())
	diags.Append(d...)

	p.Platform, d = fidoDeviceOkToTF(apiObject.GetPlatformOk())
	diags.Append(d...)

	return diags
}

func offlineDeviceOkToTF(apiObject *mfa.DeviceAuthenticationPolicyOfflineDevice, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyOfflineDeviceTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{
		"enabled":                       framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp_lifetime_duration":         types.Int64Null(),
		"otp_lifetime_timeunit":         types.StringNull(),
		"otp_failure_count":             types.Int64Null(),
		"otp_failure_cooldown_duration": types.Int64Null(),
		"otp_failure_cooldown_timeunit": types.StringNull(),
	}

	if v, ok := apiObject.GetOtpOk(); ok {
		if v1, ok := v.GetLifeTimeOk(); ok {
			objectMap["otp_lifetime_duration"] = framework.Int32OkToTF(v1.GetDurationOk())
			objectMap["otp_lifetime_timeunit"] = framework.EnumOkToTF(v1.GetTimeUnitOk())
		}

		if v1, ok := v.GetFailureOk(); ok {
			objectMap["otp_failure_count"] = framework.Int32OkToTF(v1.GetCountOk())

			if v2, ok := v1.GetCoolDownOk(); ok {
				objectMap["otp_failure_cooldown_duration"] = framework.Int32OkToTF(v2.GetDurationOk())
				objectMap["otp_failure_cooldown_timeunit"] = framework.EnumOkToTF(v2.GetTimeUnitOk())
			}
		}
	}

	flattenedObj, d := types.ObjectValue(mfaPolicyOfflineDeviceTFObjectTypes, objectMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func mobileDeviceOkToTF(apiObject *mfa.DeviceAuthenticationPolicyMobile, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyMobileTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{
		"enabled":                       framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp_failure_count":             types.Int64Null(),
		"otp_failure_cooldown_duration": types.Int64Null(),
		"otp_failure_cooldown_timeunit": types.StringNull(),
		"application":                   types.SetNull(types.ObjectType{AttrTypes: mfaPolicyMobileApplicationTFObjectTypes}),
	}

	if v, ok := apiObject.GetOtpOk(); ok {
		if v1, ok := v.GetFailureOk(); ok {
			objectMap["otp_failure_count"] = framework.Int32OkToTF(v1.GetCountOk())

			if v2, ok := v1.GetCoolDownOk(); ok {
				objectMap["otp_failure_cooldown_duration"] = framework.Int32OkToTF(v2.GetDurationOk())
				objectMap["otp_failure_cooldown_timeunit"] = framework.EnumOkToTF(v2.GetTimeUnitOk())
			}
		}
	}

	applicationObj, d := mobileApplicationsOkToTF(apiObject.GetApplicationsOk())
	diags.Append(d...)
	objectMap["application"] = applicationObj

	flattenedObj, d := types.ObjectValue(mfaPolicyMobileTFObjectTypes, objectMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func mobileApplicationsOkToTF(apiObject []mfa.DeviceAuthenticationPolicyMobileApplicationsInner, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyMobileApplicationTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	list := make([]attr.Value, 0)
	for _, item := range apiObject {
		objectMap := map[string]attr.Value{
			"id":                           framework.StringOkToTF(item.GetIdOk()),
			"push_enabled":                 types.BoolNull(),
			"push_timeout_duration":        types.Int64Null(),
			"push_timeout_timeunit":        types.StringNull(),
			"otp_enabled":                  types.BoolNull(),
			"device_authorization_enabled": types.BoolNull(),
			"device_authorization_extra_verification": types.StringNull(),
			"auto_enrollment_enabled":                 types.BoolNull(),
			"integrity_detection":                     framework.EnumOkToTF(item.GetIntegrityDetectionOk()),
			"pairing_key_lifetime_duration":           types.Int64Null(),
			"pairing_key_lifetime_timeunit":           types.StringNull(),
			"push_limit":                              types.ListNull(types.ObjectType{AttrTypes: mfaPolicyMobileApplicationPushLimitTFObjectTypes}),
		}

		if v, ok := item.GetPushOk(); ok {
			objectMap["push_enabled"] = framework.BoolOkToTF(v.GetEnabledOk())
		}

		if v, ok := item.GetPushTimeoutOk(); ok {
			objectMap["push_timeout_duration"] = framework.Int32OkToTF(v.GetDurationOk())
			objectMap["push_timeout_timeunit"] = framework.EnumOkToTF(v.GetTimeUnitOk())
		}

		if v, ok := item.GetOtpOk(); ok {
			objectMap["otp_enabled"] = framework.BoolOkToTF(v.GetEnabledOk())
		}

		if v, ok := item.GetDeviceAuthorizationOk(); ok {
			objectMap["device_authorization_enabled"] = framework.BoolOkToTF(v.GetEnabledOk())
			objectMap["device_authorization_extra_verification"] = framework.EnumOkToTF(v.GetExtraVerificationOk())
		}

		if v, ok := item.GetAutoEnrollmentOk(); ok {
			objectMap["auto_enrollment_enabled"] = framework.BoolOkToTF(v.GetEnabledOk())
		}

		if v, ok := item.GetPairingKeyLifetimeOk(); ok {
			objectMap["pairing_key_lifetime_duration"] = framework.Int32OkToTF(v.GetDurationOk())
			objectMap["pairing_key_lifetime_timeunit"] = framework.EnumOkToTF(v.GetTimeUnitOk())
		}

		pushLimitObj, d := mobileApplicationsPushLimitsOkToTF(item.GetPushLimitOk())
		diags.Append(d...)
		objectMap["push_limit"] = pushLimitObj

		flattenedObj, d := types.ObjectValue(mfaPolicyMobileApplicationTFObjectTypes, objectMap)
		diags.Append(d...)

		list = append(list, flattenedObj)
	}

	return types.SetValueMust(tfObjType, list), diags
}

func mobileApplicationsPushLimitsOkToTF(apiObject *mfa.DeviceAuthenticationPolicyMobileApplicationsInnerPushLimit, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyMobileApplicationPushLimitTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{
		"count":                  framework.Int32OkToTF(apiObject.GetCountOk()),
		"lock_duration_duration": types.Int64Null(),
		"lock_duration_timeunit": types.StringNull(),
		"time_period_duration":   types.Int64Null(),
		"time_period_timeunit":   types.StringNull(),
	}

	if v, ok := apiObject.GetLockDurationOk(); ok {
		objectMap["lock_duration_duration"] = framework.Int32OkToTF(v.GetDurationOk())
		objectMap["lock_duration_timeunit"] = framework.EnumOkToTF(v.GetTimeUnitOk())
	}

	if v, ok := apiObject.GetTimePeriodOk(); ok {
		objectMap["time_period_duration"] = framework.Int32OkToTF(v.GetDurationOk())
		objectMap["time_period_timeunit"] = framework.EnumOkToTF(v.GetTimeUnitOk())
	}

	flattenedObj, d := types.ObjectValue(mfaPolicyMobileApplicationPushLimitTFObjectTypes, objectMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

}

func totpDeviceOkToTF(apiObject *mfa.DeviceAuthenticationPolicyTotp, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyTotpTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{
		"enabled":                       framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"otp_failure_count":             types.Int64Null(),
		"otp_failure_cooldown_duration": types.Int64Null(),
		"otp_failure_cooldown_timeunit": types.StringNull(),
	}

	if v, ok := apiObject.GetOtpOk(); ok {
		if v1, ok := v.GetFailureOk(); ok {
			objectMap["otp_failure_count"] = framework.Int32OkToTF(v1.GetCountOk())

			if v2, ok := v1.GetCoolDownOk(); ok {
				objectMap["otp_failure_cooldown_duration"] = framework.Int32OkToTF(v2.GetDurationOk())
				objectMap["otp_failure_cooldown_timeunit"] = framework.EnumOkToTF(v2.GetTimeUnitOk())
			}
		}
	}

	flattenedObj, d := types.ObjectValue(mfaPolicyTotpTFObjectTypes, objectMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func fidoDeviceOkToTF(apiObject *mfa.DeviceAuthenticationPolicyFIDODevice, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: mfaPolicyFidoDeviceTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectMap := map[string]attr.Value{
		"enabled":        framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"fido_policy_id": framework.StringOkToTF(apiObject.GetFidoPolicyIdOk()),
	}

	flattenedObj, d := types.ObjectValue(mfaPolicyFidoDeviceTFObjectTypes, objectMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func checkApplicationForMobileApp(ctx context.Context, apiClient *management.APIClient, environmentID, appID string) diag.Diagnostics {
	var diags diag.Diagnostics

	resp, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApi.ReadOneApplication(ctx, environmentID, appID).Execute()
		},
		"ReadOneApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		diags.AddError(
			"Application referenced in `mobile.application.id` does not exist",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with Apple, Google or Huawei app settings configured.",
		)
		return diags
	}

	respObject := resp.(*management.ReadOneApplication200Response)

	// check if oidc
	if respObject.ApplicationOIDC == nil {
		diags.AddError(
			"Application referenced in `mobile.application.id` is not of type OIDC",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with Apple, Google or Huawei app settings configured.",
		)
		return diags
	}

	// check if native
	if respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_NATIVE_APP && respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_CUSTOM_APP {
		diags.AddError(
			"Application referenced in `mobile.application.id` is OIDC, but is not the required `Native` OIDC application type",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with Apple, Google or Huawei app settings configured.",
		)
		return diags
	}

	// check if mobile set and package/bundle set
	if _, ok := respObject.ApplicationOIDC.GetMobileOk(); !ok {
		diags.AddError(
			"Application referenced in `mobile.application.id` does not contain mobile application configuration",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with Apple, Google or Huawei app settings configured.",
		)
		return diags
	}

	if v, ok := respObject.ApplicationOIDC.GetMobileOk(); ok {

		_, bundleIDOk := v.GetBundleIdOk()
		_, packageNameOk := v.GetPackageNameOk()
		_, huaweiAppIdOk := v.GetHuaweiAppIdOk()

		if !bundleIDOk && !packageNameOk && !huaweiAppIdOk {

			diags.AddError(
				"Application referenced in `mobile.application.id` does not contain mobile application configuration",
				"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with Apple, Google or Huawei app settings configured.",
			)
			return diags
		}
	}

	return diags
}
