package verify

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	int64validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int64validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VerifyPolicyResource serviceClientType

type verifyPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Name             types.String `tfsdk:"name"`
	Default          types.Bool   `tfsdk:"default"`
	Description      types.String `tfsdk:"description"`
	GovernmentId     types.Object `tfsdk:"government_id"`
	FacialComparison types.Object `tfsdk:"facial_comparison"`
	Liveness         types.Object `tfsdk:"liveness"`
	Email            types.Object `tfsdk:"email"`
	Phone            types.Object `tfsdk:"phone"`
	Transaction      types.Object `tfsdk:"transaction"`
	Voice            types.Object `tfsdk:"voice"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

type governmentIdModel struct {
	Verify types.String `tfsdk:"verify"`
}

type facialComparisonModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type livenessnModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type genericTimeoutModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type deviceModel struct {
	CreateMfaDevice types.Bool   `tfsdk:"create_mfa_device"`
	OTP             types.Object `tfsdk:"otp"`
	Verify          types.String `tfsdk:"verify"`
}

type otpConfigurationModel struct {
	Attempts     types.Object `tfsdk:"attempts"`
	Deliveries   types.Object `tfsdk:"deliveries"`
	LifeTime     types.Object `tfsdk:"lifetime"`
	Notification types.Object `tfsdk:"notification"`
}

type otpAttemptsModel struct {
	Count types.Int64 `tfsdk:"count"`
}

type otpDeliveriesModel struct {
	Count    types.Int64  `tfsdk:"count"`
	Cooldown types.Object `tfsdk:"cooldown"`
}

type otpNotificationModel struct {
	TemplateName types.String `tfsdk:"template_name"`
	VariantName  types.String `tfsdk:"variant_name"`
}

type transactionModel struct {
	Timeout            types.Object `tfsdk:"timeout"`
	DataCollection     types.Object `tfsdk:"data_collection"`
	DataCollectionOnly types.Bool   `tfsdk:"data_collection_only"`
}

type transactionDataCollectionModel struct {
	Timeout types.Object `tfsdk:"timeout"`
}

type voiceModel struct {
	Verify              types.String `tfsdk:"verify"`
	Enrollment          types.Bool   `tfsdk:"enrollment"`
	ComparisonThreshold types.String `tfsdk:"comparison_threshold"`
	LivenessThreshold   types.String `tfsdk:"liveness_threshold"`
	TextDependent       types.Object `tfsdk:"text_dependent"`
	ReferenceData       types.Object `tfsdk:"reference_data"`
}

type textDependentModel struct {
	Samples  types.Int64  `tfsdk:"samples"`
	PhraseId types.String `tfsdk:"voice_phrase_id"`
}

type referenceDataModel struct {
	RetainOriginalRecordings types.Bool `tfsdk:"retain_original_recordings"`
	UpdateOnReenrollment     types.Bool `tfsdk:"update_on_reenrollment"`
	UpdateOnRVerification    types.Bool `tfsdk:"update_on_verification"`
}

var (
	genericTimeoutServiceTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int64Type,
		"time_unit": types.StringType,
	}

	governmentIdServiceTFObjectTypes = map[string]attr.Type{
		"verify": types.StringType,
	}

	facialComparisonServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	livenessServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	deviceServiceTFObjectTypes = map[string]attr.Type{
		"verify":            types.StringType,
		"create_mfa_device": types.BoolType,
		"otp":               types.ObjectType{AttrTypes: otpServiceTFObjectTypes},
	}

	otpServiceTFObjectTypes = map[string]attr.Type{
		"attempts":     types.ObjectType{AttrTypes: otpAttemptsServiceTFObjectTypes},
		"deliveries":   types.ObjectType{AttrTypes: otpDeliveriesServiceTFObjectTypes},
		"lifetime":     types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
		"notification": types.ObjectType{AttrTypes: otpNotificationServiceTFObjectTypes},
	}

	otpAttemptsServiceTFObjectTypes = map[string]attr.Type{
		"count": types.Int64Type,
	}

	otpDeliveriesServiceTFObjectTypes = map[string]attr.Type{
		"count":    types.Int64Type,
		"cooldown": types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
	}

	otpNotificationServiceTFObjectTypes = map[string]attr.Type{
		"template_name": types.StringType,
		"variant_name":  types.StringType,
	}

	transactionServiceTFObjectTypes = map[string]attr.Type{
		"timeout":              types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
		"data_collection":      types.ObjectType{AttrTypes: dataCollectionServiceTFObjectTypes},
		"data_collection_only": types.BoolType,
	}

	dataCollectionServiceTFObjectTypes = map[string]attr.Type{
		"timeout": types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
	}

	voiceServiceTFObjectTypes = map[string]attr.Type{
		"verify":               types.StringType,
		"enrollment":           types.BoolType,
		"comparison_threshold": types.StringType,
		"liveness_threshold":   types.StringType,
		"text_dependent":       types.ObjectType{AttrTypes: textDependentServiceTFObjectTypes},
		"reference_data":       types.ObjectType{AttrTypes: referenceDataServiceTFObjectTypes},
	}

	textDependentServiceTFObjectTypes = map[string]attr.Type{
		"samples":         types.Int64Type,
		"voice_phrase_id": types.StringType,
	}

	referenceDataServiceTFObjectTypes = map[string]attr.Type{
		"retain_original_recordings": types.BoolType,
		"update_on_reenrollment":     types.BoolType,
		"update_on_verification":     types.BoolType,
	}

	verifyPolicyOptions = []validator.Object{
		objectvalidator.AtLeastOneOf(
			path.MatchRelative().AtParent().AtName("government_id"),
			path.MatchRelative().AtParent().AtName("facial_comparison"),
			path.MatchRelative().AtParent().AtName("liveness"),
			path.MatchRelative().AtParent().AtName("email"),
			path.MatchRelative().AtParent().AtName("phone"),
			path.MatchRelative().AtParent().AtName("voice"),
		),
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &VerifyPolicyResource{}
	_ resource.ResourceWithConfigure   = &VerifyPolicyResource{}
	_ resource.ResourceWithImportState = &VerifyPolicyResource{}
)

// New Object
func NewVerifyPolicyResource() resource.Resource {
	return &VerifyPolicyResource{}
}

// Metadata
func (r *VerifyPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verify_policy"
}

func (r *VerifyPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMaxLength = 1024

	const attrMinDuration = 0
	const attrMaxDurationSeconds = 1800
	const attrMaxDurationMinutes = 30

	const attrMinVoiceSamples = 3
	const attrMaxVoiceSamples = 5

	const attrMinLifetimeDurationSeconds = 60
	const attrMaxLifetimeDurationSeconds = 1800
	const attrMinLifetimeDurationMinutes = 1
	const attrMaxLifetimeDurationMinutes = 30

	// defaults
	const defaultNotificationTemplate = "email_phone_verification"

	const defaultVerify = verify.ENUMVERIFY_DISABLED
	const defaultThreshold = verify.ENUMTHRESHOLD_MEDIUM
	const defaultOTPAttemptsCount = 5
	const defaultOTPDeliveryCount = 3

	const defaultOTPEmailDuration = 10
	const defaultOTPEmailTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultOTPPhoneDuration = 5
	const defaultOTPPhoneTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultOTPCooldownDuration = 30
	const defaultOTPCooldownTimeUnit = verify.ENUMTIMEUNIT_SECONDS

	const defaultTransactionDuration = 30
	const defaultTransactionDataCollectionDuration = 15
	const defaultTransactionTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultVoiceSamples = 3
	// P1 Platform does not set a traditional UUID as the default phrase ID value
	const defaultVoicePhraseId = "exceptional_experiences"

	const defaultBoolFalse = false
	const defaultBoolTrue = true

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies whether this is the environment's default verify policy.",
	)

	governmentIdVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls Government ID verification requirements.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	facialComparisonVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls Facial Comparison verification requirements.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	facialComparisonThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Facial Comparison threshold requirements.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues).DefaultValue(string(defaultThreshold))

	livenessVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls Liveness Check verification requirements.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	livenessThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Liveness Check threshold requirements.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues).DefaultValue(string(defaultThreshold))

	deviceVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls the verification requirements for an Email or Phone verification.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	voiceVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls the verification requirements for a Voice verification.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	voiceEnrollmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls if the transaction performs voice enrollment (`TRUE`) or voice verification (`FALSE`).",
	)

	voiceTexttDependentSamplesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Number of voice samples to collect. The allowed range is `%d - %d`.", attrMinVoiceSamples, attrMaxVoiceSamples),
	)

	voicePhraseIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of a defined `voice_phrase` to associate with the policy.",
	)
	voiceComparisonThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Comparison threshold requirements.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues).DefaultValue(string(defaultThreshold))

	voiceLivenessThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Liveness threshold requirements.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues).DefaultValue(string(defaultThreshold))

	referenceDataUpdateOnEnrollmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls updates to user's voice reference data (voice recordings) upon user re-enrollment. If `TRUE`, new data adds to existing data. If `FALSE`, new data replaces existing data.",
	)

	referenceDataUpdateOnVerificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls updates to user's voice reference data (voice recordings) upon user verification. If `TRUE`, new data adds to existing data. If `FALSE`, new voice recordings are not retained as reference data.",
	)

	otpLifeTimeEmailDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Lifetime of the OTP delivered via email.\n" +
			fmt.Sprintf("    - If `lifetime.time_unit` is `MINUTES`, the allowed range is `%d - %d`.\n", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes) +
			fmt.Sprintf("    - If `lifetime.time_unit` is `SECONDS`, the allowed range is `%d - %d`.\n", attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds) +
			fmt.Sprintf("    - Defaults to `%d %s`.\n", defaultOTPEmailDuration, defaultOTPEmailTimeUnit),
	)

	otpLifetimeEmailTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Time unit of the OTP (Email) duration lifetime.",
	).AllowedValuesEnum(verify.AllowedEnumTimeUnitEnumValues).DefaultValue(string(defaultOTPEmailTimeUnit))

	otpLifeTimePhoneDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Lifetime of the OTP delivered via phone (SMS).\n" +
			fmt.Sprintf("    - If `lifetime.time_unit` is `MINUTES`, the allowed range is `%d - %d`.\n", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes) +
			fmt.Sprintf("    - If `lifetime.time_unit` is `SECONDS`, the allowed range is `%d - %d`.\n", attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds) +
			fmt.Sprintf("    - Defaults to `%d %s`.\n", defaultOTPPhoneDuration, defaultOTPPhoneTimeUnit),
	)

	otpLifetimePhoneTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Time unit of the OTP (SMS) duration lifetime.",
	).AllowedValuesEnum(verify.AllowedEnumTimeUnitEnumValues).DefaultValue(string(defaultOTPPhoneTimeUnit))

	otpDeliveriesCooldownDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Cooldown duration.\n" +
			fmt.Sprintf("    - If `cooldown.time_unit` is `MINUTES`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationMinutes) +
			fmt.Sprintf("    - If `cooldown.time_unit` is `SECONDS`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationSeconds) +
			fmt.Sprintf("    - Defaults to `%d %s`.\n", defaultOTPCooldownDuration, defaultOTPCooldownTimeUnit),
	)

	otpDeliveriesCooldownTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Time unit of the cooldown duration configuration.",
	).AllowedValuesEnum(verify.AllowedEnumTimeUnitEnumValues).DefaultValue(string(defaultOTPCooldownTimeUnit))

	otpNotificationTemplateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Name of the template to use to pass a one-time passcode (OTP). The default value of `%s` is static. Use the `notification.variant_name` property to define an alternate template.", defaultNotificationTemplate),
	)

	transactionTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Length of time before the transaction expires.\n" +
			fmt.Sprintf("    - If `transaction.timeout.time_unit` is `MINUTES`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationMinutes) +
			fmt.Sprintf("    - If `transaction.timeout.time_unit` is `SECONDS`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationSeconds) +
			fmt.Sprintf("    - Defaults to `%d %s`.\n", defaultTransactionDuration, defaultTransactionTimeUnit),
	)

	transactionTimeoutTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Time unit of transaction timeout.",
	).AllowedValuesEnum(verify.AllowedEnumTimeUnitEnumValues).DefaultValue(string(defaultTransactionTimeUnit))

	dataCollectionDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Length of time before the data collection transaction expires.\n" +
			fmt.Sprintf("    - If `transaction.data_collection.timeout.time_unit` is `MINUTES`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationMinutes) +
			fmt.Sprintf("    - If `transaction.data_collection.timeout.time_unit` is `SECONDS`, the allowed range is `%d - %d`.\n", attrMinDuration, attrMaxDurationSeconds) +
			fmt.Sprintf("    - Defaults to `%d %s`.\n\n", defaultTransactionDataCollectionDuration, defaultTransactionTimeUnit) +
			"    ~> When setting or changing timeouts in the transaction configuration object, `transaction.data_collection.timeout.duration` must be less than or equal to `transaction.timeout.duration`.\n",
	)

	dataCollectionTimeoutTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Time unit of data collection timeout.",
	).AllowedValuesEnum(verify.AllowedEnumTimeUnitEnumValues).DefaultValue(string(defaultTransactionTimeUnit))

	dataCollectionOnlyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"When `true`, collects documents specified in the policy without determining their validity; defaults to `false`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to configure the requirements to verify a user, including the parameters for verification.\n\n" +
			"A verify policy defines which of the following one or more checks are performed for a verification transaction and configures the parameters of each check. " +
			"If a type is optional, then the transaction can be processed with or without the documents for that type. If the documents are provided for that type and the optional type verification fails, it will not cause the entire transaction to fail.\n\n" +
			"Verify policies can perform any of the following checks:\n" +
			"- Government identity document - Validate a government-issued identity document, which includes a photograph.\n" +
			"- Facial comparison - Compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph.\n" +
			"- Liveness - Inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask.\n" +
			"- Email - Receive a one-time password (OTP) on an email address and return the OTP to the service.\n" +
			"- Phone - Receive a one-time password (OTP) on a mobile phone and return the OTP to the service.\n" +
			"- Voice - Compare a voice recording to a previously submitted reference voice recording.\n\n ",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify policy exists."),
			),

			"name": schema.StringAttribute{
				Description: "Name of the verification policy displayed in PingOne Admin UI.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"description": schema.StringAttribute{
				Description: "Description of the verification policy displayed in PingOne Admin UI, 1-1024 characters.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(attrMinLength, attrMaxLength),
				},
			},

			"government_id": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements for a government-issued identity document, which includes a photograph.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					governmentIdServiceTFObjectTypes,
					map[string]attr.Value{
						"verify": types.StringValue(string(defaultVerify)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         governmentIdVerifyDescription.Description,
						MarkdownDescription: governmentIdVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"facial_comparison": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					facialComparisonServiceTFObjectTypes,
					map[string]attr.Value{
						"verify":    types.StringValue(string(defaultVerify)),
						"threshold": types.StringValue(string(defaultThreshold)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         facialComparisonVerifyDescription.Description,
						MarkdownDescription: facialComparisonVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
					"threshold": schema.StringAttribute{
						Description:         facialComparisonThresholdDescription.Description,
						MarkdownDescription: facialComparisonThresholdDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"liveness": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					livenessServiceTFObjectTypes,
					map[string]attr.Value{
						"verify":    types.StringValue(string(defaultVerify)),
						"threshold": types.StringValue(string(defaultThreshold)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
					"threshold": schema.StringAttribute{
						Description:         livenessThresholdDescription.Description,
						MarkdownDescription: livenessThresholdDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"email": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to validate an email address using a one-time password (OTP).",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"count": types.Int64Value(defaultOTPAttemptsCount),
					}
					attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int64Value(defaultOTPEmailDuration),
						"time_unit": types.StringValue(string(defaultOTPEmailTimeUnit)),
					}
					lifetimeObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"template_name": types.StringValue(defaultNotificationTemplate),
						"variant_name":  types.StringNull(),
					}
					notificationObjValue, d := types.ObjectValue(otpNotificationServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int64Value(defaultOTPCooldownDuration),
						"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
					}
					cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"count":    types.Int64Value(defaultOTPDeliveryCount),
						"cooldown": cooldownObjValue,
					}
					deliveriesObjValue, d := types.ObjectValue(otpDeliveriesServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"attempts":     attemptsObjValue,
						"lifetime":     lifetimeObjValue,
						"deliveries":   deliveriesObjValue,
						"notification": notificationObjValue,
					}
					otpObjValue, d := types.ObjectValue(otpServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"verify":            types.StringValue(string(defaultVerify)),
						"create_mfa_device": types.BoolValue(defaultBoolFalse),
						"otp":               otpObjValue,
					}
					objValue, d := types.ObjectValue(deviceServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description: "When enabled, PingOne Verify registers the email address with PingOne MFA as a verified MFA device.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(defaultBoolFalse),
					},
					"otp": schema.SingleNestedAttribute{
						Description: "SMS/Voice/Email one-time password (OTP) configuration.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Required:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Required:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Required:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.Int64{
													int64validator.Any(
														int64validator.All(
															int64validator.Between(attrMinDuration, attrMaxDurationMinutes),
															int64validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int64validator.All(
															int64validator.Between(attrMinDuration, attrMaxDurationSeconds),
															int64validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},
											"time_unit": schema.StringAttribute{
												Description:         otpDeliveriesCooldownTimeUnitDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownTimeUnitDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "The length of time for which the OTP is valid.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         otpLifeTimeEmailDurationDescription.Description,
										MarkdownDescription: otpLifeTimeEmailDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int64{
											int64validator.Any(
												int64validator.All(
													int64validator.Between(attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int64validator.All(
													int64validator.Between(attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`SECONDS`),
														fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
											),
										},
									},
									"time_unit": schema.StringAttribute{
										Description:         otpLifetimeEmailTimeUnitDescription.Description,
										MarkdownDescription: otpLifetimeEmailTimeUnitDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
										},
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "OTP notification template configuration.",
								Optional:    true,
								Computed:    true,

								Default: objectdefault.StaticValue(types.ObjectValueMust(
									otpNotificationServiceTFObjectTypes,
									map[string]attr.Value{
										"template_name": types.StringValue(defaultNotificationTemplate),
										"variant_name":  types.StringNull(),
									},
								)),

								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         otpNotificationTemplateDescription.Description,
										MarkdownDescription: otpNotificationTemplateDescription.MarkdownDescription,
										Computed:            true,
										Default:             stringdefault.StaticString(defaultNotificationTemplate),
									},
									"variant_name": schema.StringAttribute{
										Description: "Name of the template variant to use to pass a one-time passcode (OTP).",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         deviceVerifyDescription.Description,
						MarkdownDescription: deviceVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"phone": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to validate a mobile phone number using a one-time password (OTP).",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"count": types.Int64Value(defaultOTPAttemptsCount),
					}
					attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int64Value(defaultOTPPhoneDuration),
						"time_unit": types.StringValue(string(defaultOTPPhoneTimeUnit)),
					}
					lifetimeObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"template_name": types.StringValue(defaultNotificationTemplate),
						"variant_name":  types.StringNull(),
					}
					notificationObjValue, d := types.ObjectValue(otpNotificationServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int64Value(defaultOTPCooldownDuration),
						"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
					}
					cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"count":    types.Int64Value(defaultOTPDeliveryCount),
						"cooldown": cooldownObjValue,
					}
					deliveriesObjValue, d := types.ObjectValue(otpDeliveriesServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"attempts":     attemptsObjValue,
						"lifetime":     lifetimeObjValue,
						"deliveries":   deliveriesObjValue,
						"notification": notificationObjValue,
					}
					otpObjValue, d := types.ObjectValue(otpServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"verify":            types.StringValue(string(defaultVerify)),
						"create_mfa_device": types.BoolValue(defaultBoolFalse),
						"otp":               otpObjValue,
					}
					objValue, d := types.ObjectValue(deviceServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description: "When enabled, PingOne Verify registers the mobile phone with PingOne MFA as a verified MFA device.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(defaultBoolFalse),
					},
					"otp": schema.SingleNestedAttribute{
						Description: "SMS/Voice/Email one-time password (OTP) configuration.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Required:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Required:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Required:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.Int64{
													int64validator.Any(
														int64validator.All(
															int64validator.Between(attrMinDuration, attrMaxDurationMinutes),
															int64validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int64validator.All(
															int64validator.Between(attrMinDuration, attrMaxDurationSeconds),
															int64validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`SECONDS`),
																fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationSeconds),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
													),
												},
											},
											"time_unit": schema.StringAttribute{
												Description:         otpDeliveriesCooldownTimeUnitDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownTimeUnitDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "The length of time for which the OTP is valid.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         otpLifeTimePhoneDurationDescription.Description,
										MarkdownDescription: otpLifeTimePhoneDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int64{
											int64validator.Any(
												int64validator.All(
													int64validator.Between(attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int64validator.All(
													int64validator.Between(attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`SECONDS`),
														fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
											),
										},
									},
									"time_unit": schema.StringAttribute{
										Description:         otpLifetimePhoneTimeUnitDescription.Description,
										MarkdownDescription: otpLifetimePhoneTimeUnitDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
										},
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "OTP notification template configuration.",
								Optional:    true,
								Computed:    true,

								Default: objectdefault.StaticValue(types.ObjectValueMust(
									otpNotificationServiceTFObjectTypes,
									map[string]attr.Value{
										"template_name": types.StringValue(defaultNotificationTemplate),
										"variant_name":  types.StringNull(),
									},
								)),

								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         otpNotificationTemplateDescription.Description,
										MarkdownDescription: otpNotificationTemplateDescription.MarkdownDescription,
										Computed:            true,
										Default:             stringdefault.StaticString(defaultNotificationTemplate),
									},
									"variant_name": schema.StringAttribute{
										Description: "Name of the template variant to use to pass a one-time passcode (OTP).",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         deviceVerifyDescription.Description,
						MarkdownDescription: deviceVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"transaction": schema.SingleNestedAttribute{
				Description: "Defines the requirements for transactions invoked by the policy.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"duration":  framework.Int32ToTF(defaultTransactionDuration),
						"time_unit": framework.EnumOkToTF(defaultTransactionTimeUnit, true),
					}
					timeoutObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  framework.Int32ToTF(defaultTransactionDataCollectionDuration),
						"time_unit": framework.EnumOkToTF(defaultTransactionTimeUnit, true),
					}
					dataCollectionTimeoutObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"timeout": dataCollectionTimeoutObjValue,
					}
					dataCollectionObjValue, d := types.ObjectValue(dataCollectionServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					dataCollectionBool := new(bool)
					*dataCollectionBool = false
					objValue, d := types.ObjectValue(transactionServiceTFObjectTypes, map[string]attr.Value{
						"timeout":              timeoutObjValue,
						"data_collection":      dataCollectionObjValue,
						"data_collection_only": framework.BoolOkToTF(dataCollectionBool, true),
					})
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"timeout": schema.SingleNestedAttribute{
						Description: "Object for transaction timeout.",
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int64Attribute{
								Description:         transactionTimeoutDurationDescription.Description,
								MarkdownDescription: transactionTimeoutDurationDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.Int64{
									int64validator.Any(
										int64validator.All(
											int64validator.Between(attrMinDuration, attrMaxDurationMinutes),
											int64validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int64validator.All(
											int64validator.Between(attrMinDuration, attrMaxDurationSeconds),
											int64validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`SECONDS`),
												fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationSeconds),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
									),
								},
							},
							"time_unit": schema.StringAttribute{
								Description:         transactionTimeoutTimeUnitDescription.Description,
								MarkdownDescription: transactionTimeoutTimeUnitDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
									stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("time_unit")),
								},
							},
						},
					},
					"data_collection": schema.SingleNestedAttribute{
						Description: "Object for data collection timeout definition.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"timeout": schema.SingleNestedAttribute{
								Description: "Object for data collection timeout.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         dataCollectionDurationDescription.Description,
										MarkdownDescription: dataCollectionDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int64{
											int64validatorinternal.IsLessThanEqualToPathValue(
												path.MatchRoot("transaction").AtName("timeout").AtName("duration"),
											),
											int64validator.Any(
												int64validator.All(
													int64validator.Between(attrMinDuration, attrMaxDurationMinutes),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int64validator.All(
													int64validator.Between(attrMinDuration, attrMaxDurationSeconds),
													int64validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`SECONDS`),
														fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationSeconds),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
											),
										},
									},
									"time_unit": schema.StringAttribute{
										Description:         dataCollectionTimeoutTimeUnitDescription.Description,
										MarkdownDescription: dataCollectionTimeoutTimeUnitDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumTimeUnitEnumValues)...),
											stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("time_unit")),
										},
									},
								},
							},
						},
					},
					"data_collection_only": schema.BoolAttribute{
						Description:         dataCollectionOnlyDescription.Description,
						MarkdownDescription: dataCollectionOnlyDescription.MarkdownDescription,
						Optional:            true,
					},
				},
			},

			"voice": schema.SingleNestedAttribute{
				Description: "Defines the requirements for transactions invoked by the policy.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"samples":         types.Int64Value(defaultVoiceSamples),
						"voice_phrase_id": types.StringValue(defaultVoicePhraseId),
					}
					textDependentObjValue, d := types.ObjectValue(textDependentServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"retain_original_recordings": types.BoolValue(defaultBoolFalse),
						"update_on_reenrollment":     types.BoolValue(defaultBoolTrue),
						"update_on_verification":     types.BoolValue(defaultBoolTrue),
					}
					referenceDataObjValue, d := types.ObjectValue(referenceDataServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					objValue, d := types.ObjectValue(voiceServiceTFObjectTypes, map[string]attr.Value{
						"verify":               types.StringValue(string(defaultVerify)),
						"enrollment":           types.BoolValue(defaultBoolFalse),
						"comparison_threshold": types.StringValue(string(defaultThreshold)),
						"liveness_threshold":   types.StringValue(string(defaultThreshold)),
						"text_dependent":       textDependentObjValue,
						"reference_data":       referenceDataObjValue,
					})
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         voiceVerifyDescription.Description,
						MarkdownDescription: voiceVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
					"enrollment": schema.BoolAttribute{
						Description:         voiceEnrollmentDescription.Description,
						MarkdownDescription: voiceEnrollmentDescription.MarkdownDescription,
						Required:            true,
					},
					"comparison_threshold": schema.StringAttribute{
						Description:         voiceComparisonThresholdDescription.Description,
						MarkdownDescription: voiceComparisonThresholdDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
					"liveness_threshold": schema.StringAttribute{
						Description:         voiceLivenessThresholdDescription.Description,
						MarkdownDescription: voiceLivenessThresholdDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
					"text_dependent": schema.SingleNestedAttribute{
						Description: "Object for configuration of text dependent voice verification.",
						Optional:    true,
						Computed:    true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							textDependentServiceTFObjectTypes,
							map[string]attr.Value{
								"samples":         types.Int64Value(defaultVoiceSamples),
								"voice_phrase_id": types.StringValue(defaultVoicePhraseId),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"samples": schema.Int64Attribute{
								Description:         voiceTexttDependentSamplesDescription.Description,
								MarkdownDescription: voiceTexttDependentSamplesDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.Int64{
									int64validator.Between(attrMinVoiceSamples, attrMaxVoiceSamples),
								},
							},
							"voice_phrase_id": schema.StringAttribute{
								Description:         voicePhraseIdDescription.Description,
								MarkdownDescription: voicePhraseIdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.Any(
										validation.P1ResourceIDValidator(),
									),
								},
							},
						},
					},
					"reference_data": schema.SingleNestedAttribute{
						Description: "Object for configuration of voice recording reference data.",
						Optional:    true,

						Computed: true,

						Default: objectdefault.StaticValue(types.ObjectValueMust(
							referenceDataServiceTFObjectTypes,
							map[string]attr.Value{
								"retain_original_recordings": types.BoolValue(defaultBoolFalse),
								"update_on_reenrollment":     types.BoolValue(defaultBoolTrue),
								"update_on_verification":     types.BoolValue(defaultBoolTrue),
							},
						)),

						Attributes: map[string]schema.Attribute{
							"retain_original_recordings": schema.BoolAttribute{
								Description: "Controls if the service stores the original voice recordings.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"update_on_reenrollment": schema.BoolAttribute{
								Description:         referenceDataUpdateOnEnrollmentDescription.Description,
								MarkdownDescription: referenceDataUpdateOnEnrollmentDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"update_on_verification": schema.BoolAttribute{
								Description:         referenceDataUpdateOnVerificationDescription.Description,
								MarkdownDescription: referenceDataUpdateOnVerificationDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
				},

				Validators: verifyPolicyOptions,
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify policy was created.",
				Computed:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify policy was updated. Can be null.",
				Computed:    true,
			},
		},
	}
}

func (r *VerifyPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.Client = preparedClient
}

func (r *VerifyPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.Client == nil {
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
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VerifyPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.VerifyPoliciesApi.CreateVerifyPolicy(ctx, plan.EnvironmentId.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"CreateVerifyPolicy",
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
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *verifyPolicyResourceModel

	if r.Client == nil {
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
	var response *verify.VerifyPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.VerifyPoliciesApi.ReadOneVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneVerifyPolicy",
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

func (r *VerifyPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.Client == nil {
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
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VerifyPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.VerifyPoliciesApi.UpdateVerifyPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"UpdateVerifyPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
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

func (r *VerifyPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *verifyPolicyResourceModel

	if r.Client == nil {
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
			r, err := r.Client.VerifyPoliciesApi.DeleteVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteVerifyPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VerifyPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: validation.P1ResourceIDRegexp,
		},
		{
			Label:     "verify_policy_id",
			Regexp:    validation.P1ResourceIDRegexp,
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

func (p *verifyPolicyResourceModel) expand(ctx context.Context) (*verify.VerifyPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := verify.NewVerifyPolicyWithDefaults()

	// Government Id Verification Object
	if !p.GovernmentId.IsNull() && !p.GovernmentId.IsUnknown() {

		var governmentId governmentIdModel
		d := p.GovernmentId.As(ctx, &governmentId, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyGovernmentId, d := governmentId.expandgovernmentIdModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetGovernmentId(*verifyGovernmentId)
	}

	// Facial Comparison Verification Object
	if !p.FacialComparison.IsNull() && !p.FacialComparison.IsUnknown() {

		var facialComparison facialComparisonModel
		d := p.FacialComparison.As(ctx, &facialComparison, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyFacialComparison, d := facialComparison.expandFacialComparisonModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetFacialComparison(*verifyFacialComparison)
	}

	// Liveness Verification Object
	if !p.Liveness.IsNull() && !p.Liveness.IsUnknown() {

		var liveness livenessnModel
		d := p.Liveness.As(ctx, &liveness, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyLiveness, d := liveness.expandLivenessModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetLiveness(*verifyLiveness)
	}

	// Transaction Object
	if !p.Transaction.IsNull() && !p.Transaction.IsUnknown() {

		var transaction transactionModel
		d := p.Transaction.As(ctx, &transaction, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		transactionSettings, d := transaction.expandTransactionModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetTransaction(*transactionSettings)
	}

	// Email Object
	if !p.Email.IsNull() && !p.Email.IsUnknown() {

		var emailConfiguration deviceModel
		d := p.Email.As(ctx, &emailConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		emailSettings, d := emailConfiguration.expandDevice(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetEmail(*emailSettings)
	}

	// Phone Object
	if !p.Phone.IsNull() && !p.Phone.IsUnknown() {

		var phoneConfiguration deviceModel
		d := p.Phone.As(ctx, &phoneConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		phoneSettings, d := phoneConfiguration.expandDevice(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPhone(*phoneSettings)
	}

	// Voice Object
	if !p.Voice.IsNull() && !p.Voice.IsUnknown() {

		var voice voiceModel
		d := p.Voice.As(ctx, &voice, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		voiceSettings, d := voice.expandVoiceModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetVoice(*voiceSettings)
	}

	// Top-level arguments
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	// Verify policies managed via TF currently cannot be set to the default policy due to a potential lock situation or state management problem.
	// The verify policy will also have default set to false.
	data.SetDefault(false)

	return data, diags
}

func (p *governmentIdModel) expandgovernmentIdModel() (*verify.GovernmentIdConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyGovernmentId := verify.NewGovernmentIdConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyGovernmentId.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if verifyGovernmentId == nil {
		diags.AddError(
			"Unexpected Value",
			"GovernmentId configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyGovernmentId, diags

}

func (p *facialComparisonModel) expandFacialComparisonModel() (*verify.FacialComparisonConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyFacialComparison := verify.NewFacialComparisonConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyFacialComparison.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.Threshold.IsNull() && !p.Threshold.IsUnknown() {
		verifyFacialComparison.SetThreshold(verify.EnumThreshold(p.Threshold.ValueString()))
	}

	if verifyFacialComparison == nil {
		diags.AddError(
			"Unexpected Value",
			"Facial Comparison configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyFacialComparison, diags

}

func (p *livenessnModel) expandLivenessModel() (*verify.LivenessConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyLiveness := verify.NewLivenessConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyLiveness.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.Threshold.IsNull() && !p.Threshold.IsUnknown() {
		verifyLiveness.SetThreshold(verify.EnumThreshold(p.Threshold.ValueString()))
	}

	if verifyLiveness == nil {
		diags.AddError(
			"Unexpected Value",
			"Liveness configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyLiveness, diags

}

func (p *transactionModel) expandTransactionModel(ctx context.Context) (*verify.TransactionConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	transactionSettings := verify.NewTransactionConfigurationWithDefaults()

	if !p.DataCollectionOnly.IsNull() && !p.DataCollectionOnly.IsUnknown() {
		transactionSettings.SetDataCollectionOnly(p.DataCollectionOnly.ValueBool())
	}

	if !p.DataCollection.IsNull() && !p.DataCollection.IsUnknown() {
		var dataCollection transactionDataCollectionModel
		d := p.DataCollection.As(ctx, &dataCollection, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		var genericTimeout genericTimeoutModel
		d = dataCollection.Timeout.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		dataCollectionTimeout := verify.NewTransactionConfigurationDataCollectionTimeoutWithDefaults()
		if !genericTimeout.TimeUnit.IsNull() && !genericTimeout.TimeUnit.IsUnknown() {
			dataCollectionTimeout.SetTimeUnit(verify.EnumTimeUnit(genericTimeout.TimeUnit.ValueString()))
		}

		if !genericTimeout.Duration.IsNull() && !genericTimeout.Duration.IsUnknown() {
			dataCollectionTimeout.SetDuration(int32(genericTimeout.Duration.ValueInt64()))
		}

		transactionDataCollection := verify.NewTransactionConfigurationDataCollection(*dataCollectionTimeout)
		transactionSettings.SetDataCollection(*transactionDataCollection)
	}

	if !p.Timeout.IsNull() && !p.Timeout.IsUnknown() {
		var genericTimeout genericTimeoutModel
		d := p.Timeout.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		transactionTimeout := verify.NewTransactionConfigurationTimeoutWithDefaults()
		transactionTimeout.SetTimeUnit(verify.EnumTimeUnit(genericTimeout.TimeUnit.ValueString()))
		transactionTimeout.SetDuration(int32(genericTimeout.Duration.ValueInt64()))

		transactionSettings.SetTimeout(*transactionTimeout)
	}

	if transactionSettings == nil {
		diags.AddError(
			"Unexpected Value",
			"Transaction configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return transactionSettings, diags

}

func (p *deviceModel) expandDevice(ctx context.Context) (*verify.OTPDeviceConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	deviceSettings := verify.NewOTPDeviceConfigurationWithDefaults()

	if !p.CreateMfaDevice.IsNull() && !p.CreateMfaDevice.IsUnknown() {
		deviceSettings.SetCreateMfaDevice(p.CreateMfaDevice.ValueBool())
	}

	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		deviceSettings.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.OTP.IsNull() && !p.OTP.IsUnknown() {
		var otpConfiguration otpConfigurationModel
		d := p.OTP.As(ctx, &otpConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		otpSettings := verify.NewOTPDeviceConfigurationOtpWithDefaults()

		// OTP Attempts
		var attempts otpAttemptsModel
		d = otpConfiguration.Attempts.As(ctx, &attempts, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		otpAttempts := verify.NewOTPDeviceConfigurationOtpAttemptsWithDefaults()
		if !attempts.Count.IsNull() && !attempts.Count.IsUnknown() {
			otpAttempts.SetCount(int32(attempts.Count.ValueInt64()))
		}
		otpSettings.SetAttempts(*otpAttempts)

		// OTP Deliveries (also has cooldown object)
		var deliveries otpDeliveriesModel
		d = otpConfiguration.Deliveries.As(ctx, &deliveries, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		otpDeliveries := verify.NewOTPDeviceConfigurationOtpDeliveriesWithDefaults()
		if !deliveries.Count.IsNull() && !deliveries.Count.IsUnknown() {
			otpDeliveries.SetCount(int32(deliveries.Count.ValueInt64()))
		}

		if !deliveries.Cooldown.IsNull() && !deliveries.Cooldown.IsUnknown() {
			var cooldown genericTimeoutModel
			d = deliveries.Cooldown.As(ctx, &cooldown, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			deliveriesCooldown := verify.NewOTPDeviceConfigurationOtpDeliveriesCooldownWithDefaults()
			if !cooldown.Duration.IsNull() && !cooldown.Duration.IsUnknown() {
				deliveriesCooldown.SetDuration(int32(cooldown.Duration.ValueInt64()))
			}
			if !cooldown.TimeUnit.IsNull() && !cooldown.TimeUnit.IsUnknown() {
				deliveriesCooldown.SetTimeUnit(verify.EnumTimeUnit(cooldown.TimeUnit.ValueString()))
			}

			otpDeliveries.SetCooldown(*deliveriesCooldown)
		}
		otpSettings.SetDeliveries(*otpDeliveries)

		// OTP LifeTime (generic timeout model)
		var genericTimeout genericTimeoutModel
		d = otpConfiguration.LifeTime.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		otpLifeTime := verify.NewOTPDeviceConfigurationOtpLifeTimeWithDefaults()
		if !genericTimeout.TimeUnit.IsNull() && !genericTimeout.TimeUnit.IsUnknown() {
			otpLifeTime.SetTimeUnit(verify.EnumTimeUnit(genericTimeout.TimeUnit.ValueString()))
		}

		if !genericTimeout.Duration.IsNull() && !genericTimeout.Duration.IsUnknown() {
			otpLifeTime.SetDuration(int32(genericTimeout.Duration.ValueInt64()))
		}
		otpSettings.SetLifeTime(*otpLifeTime)

		// OTP Notification
		var notification otpNotificationModel
		d = otpConfiguration.Notification.As(ctx, &notification, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		otpNotification := verify.NewOTPDeviceConfigurationOtpNotificationWithDefaults()
		if !notification.TemplateName.IsNull() && !notification.TemplateName.IsUnknown() {
			otpNotification.SetTemplateName(notification.TemplateName.ValueString())
		}

		if !notification.VariantName.IsNull() && !notification.VariantName.IsUnknown() {
			otpNotification.SetVariantName(notification.VariantName.ValueString())
		}
		otpSettings.SetNotification(*otpNotification)

		// Complete OTP Object
		deviceSettings.SetOtp(*otpSettings)
	}

	if deviceSettings == nil {
		diags.AddError(
			"Unexpected Value",
			"Device configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return deviceSettings, diags

}

func (p *voiceModel) expandVoiceModel(ctx context.Context) (*verify.VoiceConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	voiceSettings := verify.NewVoiceConfigurationWithDefaults()

	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		voiceSettings.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.Enrollment.IsNull() && !p.Enrollment.IsUnknown() {
		voiceSettings.SetEnrollment(p.Enrollment.ValueBool())
	}

	if !p.ComparisonThreshold.IsNull() && !p.ComparisonThreshold.IsUnknown() {
		comparisonThreshold := verify.NewVoiceConfigurationThreshold(verify.EnumThreshold(p.ComparisonThreshold.ValueString()))
		voiceSettings.SetComparison(*comparisonThreshold)
	}

	if !p.LivenessThreshold.IsNull() && !p.LivenessThreshold.IsUnknown() {
		livenessThreshold := verify.NewVoiceConfigurationThreshold(verify.EnumThreshold(p.LivenessThreshold.ValueString()))
		voiceSettings.SetLiveness(*livenessThreshold)
	}

	if !p.TextDependent.IsNull() && !p.TextDependent.IsUnknown() {
		var textDependent textDependentModel
		d := p.TextDependent.As(ctx, &textDependent, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		textDependentObject := verify.NewVoiceConfigurationTextDependentWithDefaults()
		if !textDependent.PhraseId.IsNull() && !textDependent.PhraseId.IsUnknown() {
			textDependentPhrase := verify.NewVoiceConfigurationTextDependentPhrase(textDependent.PhraseId.ValueString())
			textDependentObject.SetPhrase(*textDependentPhrase)
		}

		if !textDependent.Samples.IsNull() && !textDependent.Samples.IsUnknown() {
			textDependentObject.SetSamples(int32(textDependent.Samples.ValueInt64()))
		}
		voiceSettings.SetTextDependent(*textDependentObject)
	}

	if !p.ReferenceData.IsNull() && !p.ReferenceData.IsUnknown() {
		var referenceData referenceDataModel
		d := p.ReferenceData.As(ctx, &referenceData, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		referenceDataObject := verify.NewVoiceConfigurationReferenceDataWithDefaults()
		if !referenceData.RetainOriginalRecordings.IsNull() && !referenceData.RetainOriginalRecordings.IsUnknown() {
			referenceDataObject.SetRetainOriginalRecordings(referenceData.RetainOriginalRecordings.ValueBool())
		}

		if !referenceData.UpdateOnReenrollment.IsNull() && !referenceData.UpdateOnReenrollment.IsUnknown() {
			referenceDataObject.SetUpdateOnReenrollment(referenceData.UpdateOnReenrollment.ValueBool())
		}

		if !referenceData.UpdateOnRVerification.IsNull() && !referenceData.UpdateOnRVerification.IsUnknown() {
			referenceDataObject.SetUpdateOnVerification(referenceData.UpdateOnRVerification.ValueBool())
		}

		voiceSettings.SetReferenceData(*referenceDataObject)
	}

	if voiceSettings == nil {
		diags.AddError(
			"Unexpected Value",
			"Voice configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return voiceSettings, diags

}

func (p *verifyPolicyResourceModel) toState(apiObject *verify.VerifyPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	var d diag.Diagnostics
	p.GovernmentId, d = p.toStateGovernmentId(apiObject.GetGovernmentIdOk())
	diags.Append(d...)

	p.FacialComparison, d = p.toStateFacialComparison(apiObject.GetFacialComparisonOk())
	diags.Append(d...)

	p.Liveness, d = p.toStateLiveness(apiObject.GetLivenessOk())
	diags.Append(d...)

	p.Email, d = p.toStateDevice(apiObject.GetEmailOk())
	diags.Append(d...)

	p.Phone, d = p.toStateDevice(apiObject.GetPhoneOk())
	diags.Append(d...)

	p.Transaction, d = p.toStateTransaction(apiObject.GetTransactionOk())
	diags.Append(d...)

	p.Voice, d = p.toStateVoice(apiObject.GetVoiceOk())
	diags.Append(d...)

	return diags
}

func (p *verifyPolicyResourceModel) toStateGovernmentId(apiObject *verify.GovernmentIdConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(governmentIdDataSourceServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(governmentIdServiceTFObjectTypes, map[string]attr.Value{
		"verify": framework.EnumOkToTF(apiObject.GetVerifyOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateFacialComparison(apiObject *verify.FacialComparisonConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(facialComparisonServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(facialComparisonServiceTFObjectTypes, map[string]attr.Value{
		"verify":    framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold": framework.EnumOkToTF(apiObject.GetThresholdOk()),
	})

	diags.Append(d...)
	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateLiveness(apiObject *verify.LivenessConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(livenessServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(livenessServiceTFObjectTypes, map[string]attr.Value{
		"verify":    framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold": framework.EnumOkToTF(apiObject.GetThresholdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateTransaction(apiObject *verify.TransactionConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(transactionServiceTFObjectTypes), diags
	}

	transactionTimeout := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
	if v, ok := apiObject.GetTimeoutOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"duration":  framework.Int32OkToTF(v.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(v.GetTimeUnitOk()),
		}

		objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionTimeout = objValue
	}

	transactionDataCollection := types.ObjectNull(dataCollectionServiceTFObjectTypes)
	if v, ok := apiObject.GetDataCollectionOk(); ok {
		var d diag.Diagnostics

		transactionDataCollectionTimeout := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
		if t, ok := v.GetTimeoutOk(); ok {
			o := map[string]attr.Value{
				"duration":  framework.Int32OkToTF(t.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(t.GetTimeUnitOk()),
			}

			objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
			diags.Append(d...)

			transactionDataCollectionTimeout = objValue
		}

		o := map[string]attr.Value{
			"timeout": transactionDataCollectionTimeout,
		}

		objValue, d := types.ObjectValue(dataCollectionServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionDataCollection = objValue

	}

	objValue, d := types.ObjectValue(transactionServiceTFObjectTypes, map[string]attr.Value{
		"timeout":              transactionTimeout,
		"data_collection":      transactionDataCollection,
		"data_collection_only": framework.BoolOkToTF(apiObject.GetDataCollectionOnlyOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateDevice(apiObject *verify.OTPDeviceConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(deviceServiceTFObjectTypes), diags
	}

	otp := types.ObjectNull(otpServiceTFObjectTypes)
	if v, ok := apiObject.GetOtpOk(); ok {
		var d diag.Diagnostics

		attempts := types.ObjectNull(otpAttemptsServiceTFObjectTypes)
		if t, ok := v.GetAttemptsOk(); ok {
			o := map[string]attr.Value{
				"count": framework.Int32OkToTF(t.GetCountOk()),
			}

			objValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
			diags.Append(d...)

			attempts = objValue
		}

		deliveries := types.ObjectNull(otpDeliveriesServiceTFObjectTypes)
		if t, ok := v.GetDeliveriesOk(); ok {

			cooldown := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
			if c, ok := t.GetCooldownOk(); ok {
				o := map[string]attr.Value{
					"duration":  framework.Int32OkToTF(c.GetDurationOk()),
					"time_unit": framework.EnumOkToTF(c.GetTimeUnitOk()),
				}
				objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
				diags.Append(d...)

				cooldown = objValue
			}

			o := map[string]attr.Value{
				"count":    framework.Int32OkToTF(t.GetCountOk()),
				"cooldown": cooldown,
			}
			objValue, d := types.ObjectValue(otpDeliveriesServiceTFObjectTypes, o)
			diags.Append(d...)

			deliveries = objValue
		}

		lifetime := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
		if t, ok := v.GetLifeTimeOk(); ok {
			o := map[string]attr.Value{
				"duration":  framework.Int32OkToTF(t.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(t.GetTimeUnitOk()),
			}

			objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
			diags.Append(d...)

			lifetime = objValue
		}

		notification := types.ObjectNull(otpNotificationServiceTFObjectTypes)
		if t, ok := v.GetNotificationOk(); ok {
			o := map[string]attr.Value{
				"template_name": framework.StringOkToTF(t.GetTemplateNameOk()),
				"variant_name":  framework.StringOkToTF(t.GetVariantNameOk()),
			}

			objValue, d := types.ObjectValue(otpNotificationServiceTFObjectTypes, o)
			diags.Append(d...)

			notification = objValue
		}

		o := map[string]attr.Value{
			"attempts":     attempts,
			"lifetime":     lifetime,
			"deliveries":   deliveries,
			"notification": notification,
		}
		objValue, d := types.ObjectValue(otpServiceTFObjectTypes, o)
		diags.Append(d...)

		otp = objValue
	}

	objValue, d := types.ObjectValue(deviceServiceTFObjectTypes, map[string]attr.Value{
		"verify":            framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"create_mfa_device": framework.BoolOkToTF(apiObject.GetCreateMfaDeviceOk()),
		"otp":               otp,
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateVoice(apiObject *verify.VoiceConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(voiceServiceTFObjectTypes), diags
	}

	textDependent := types.ObjectNull(textDependentServiceTFObjectTypes)
	if v, ok := apiObject.GetTextDependentOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"samples":         framework.Int32OkToTF(v.GetSamplesOk()),
			"voice_phrase_id": framework.StringToTF(v.GetPhrase().Id),
		}
		objValue, d := types.ObjectValue(textDependentServiceTFObjectTypes, o)
		diags.Append(d...)

		textDependent = objValue
	}

	referenceData := types.ObjectNull(referenceDataServiceTFObjectTypes)
	if v, ok := apiObject.GetReferenceDataOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"retain_original_recordings": framework.BoolOkToTF(v.GetRetainOriginalRecordingsOk()),
			"update_on_reenrollment":     framework.BoolOkToTF(v.GetUpdateOnReenrollmentOk()),
			"update_on_verification":     framework.BoolOkToTF(v.GetUpdateOnVerificationOk()),
		}
		objValue, d := types.ObjectValue(referenceDataServiceTFObjectTypes, o)
		diags.Append(d...)

		referenceData = objValue
	}

	comparisonThreshold := types.StringNull()
	if v, ok := apiObject.GetComparisonOk(); ok {
		if t, ok := v.GetThresholdOk(); ok {
			comparisonThreshold = types.StringValue(utils.EnumToString(t))
		}

	}

	livenessThreshold := types.StringNull()
	if v, ok := apiObject.GetLivenessOk(); ok {
		if t, ok := v.GetThresholdOk(); ok {
			livenessThreshold = types.StringValue(utils.EnumToString(t))
		}

	}

	objValue, d := types.ObjectValue(voiceServiceTFObjectTypes, map[string]attr.Value{
		"verify":               framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"enrollment":           framework.BoolOkToTF(apiObject.GetEnrollmentOk()),
		"comparison_threshold": comparisonThreshold,
		"liveness_threshold":   livenessThreshold,
		"text_dependent":       textDependent,
		"reference_data":       referenceData,
	})
	diags.Append(d...)

	return objValue, diags
}
