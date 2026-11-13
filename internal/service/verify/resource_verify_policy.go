// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/pingone-go-sdk-legacy/pingone/model"
	"github.com/pingidentity/pingone-go-sdk-legacy/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	customobjectvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	customstringvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VerifyPolicyResource serviceClientType

type verifyPolicyResourceModel struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                   types.String                 `tfsdk:"name"`
	Default                types.Bool                   `tfsdk:"default"`
	Description            types.String                 `tfsdk:"description"`
	GovernmentId           types.Object                 `tfsdk:"government_id"`
	FacialComparison       types.Object                 `tfsdk:"facial_comparison"`
	Liveness               types.Object                 `tfsdk:"liveness"`
	Email                  types.Object                 `tfsdk:"email"`
	Phone                  types.Object                 `tfsdk:"phone"`
	Transaction            types.Object                 `tfsdk:"transaction"`
	IdentityRecordMatching types.Object                 `tfsdk:"identity_record_matching"`
	CreatedAt              timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt              timetypes.RFC3339            `tfsdk:"updated_at"`
}

type governmentIdModel struct {
	Verify         types.String `tfsdk:"verify"`
	InspectionType types.String `tfsdk:"inspection_type"`
	FailExpiredId  types.Bool   `tfsdk:"fail_expired_id"`
	ProviderAuto   types.String `tfsdk:"provider_auto"`
	ProviderManual types.String `tfsdk:"provider_manual"`
	RetryAttempts  types.Int32  `tfsdk:"retry_attempts"`
	VerifyAamva    types.Bool   `tfsdk:"verify_aamva"`
	Aadhaar        types.Object `tfsdk:"aadhaar"`
}

type facialComparisonModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type livenessnModel struct {
	Verify        types.String `tfsdk:"verify"`
	Threshold     types.String `tfsdk:"threshold"`
	RetryAttempts types.Int32  `tfsdk:"retry_attempts"`
}

type genericTimeoutModel struct {
	Duration types.Int32  `tfsdk:"duration"`
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
	Count types.Int32 `tfsdk:"count"`
}

type otpDeliveriesModel struct {
	Count    types.Int32  `tfsdk:"count"`
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

type identityRecordMatchingModel struct {
	Address    types.Object `tfsdk:"address"`
	BirthDate  types.Object `tfsdk:"birth_date"`
	FamilyName types.Object `tfsdk:"family_name"`
	GivenName  types.Object `tfsdk:"given_name"`
	Name       types.Object `tfsdk:"name"`
}

type identityRecordMatchingFieldModel struct {
	FieldRequired types.Bool   `tfsdk:"field_required"`
	Threshold     types.String `tfsdk:"threshold"`
}

type aadhaarModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

var (
	genericTimeoutServiceTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int32Type,
		"time_unit": types.StringType,
	}

	aadhaarServiceTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	governmentIdServiceTFObjectTypes = map[string]attr.Type{
		"verify":          types.StringType,
		"inspection_type": types.StringType,
		"fail_expired_id": types.BoolType,
		"provider_auto":   types.StringType,
		"provider_manual": types.StringType,
		"retry_attempts":  types.Int32Type,
		"verify_aamva":    types.BoolType,
		"aadhaar":         types.ObjectType{AttrTypes: aadhaarServiceTFObjectTypes},
	}

	facialComparisonServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	livenessServiceTFObjectTypes = map[string]attr.Type{
		"verify":         types.StringType,
		"threshold":      types.StringType,
		"retry_attempts": types.Int32Type,
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
		"count": types.Int32Type,
	}

	otpDeliveriesServiceTFObjectTypes = map[string]attr.Type{
		"count":    types.Int32Type,
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

	identityRecordMatchingServiceTFObjectTypes = map[string]attr.Type{
		"address":     types.ObjectType{AttrTypes: identityRecordMatchingFieldServiceTFObjectTypes},
		"birth_date":  types.ObjectType{AttrTypes: identityRecordMatchingFieldServiceTFObjectTypes},
		"family_name": types.ObjectType{AttrTypes: identityRecordMatchingFieldServiceTFObjectTypes},
		"given_name":  types.ObjectType{AttrTypes: identityRecordMatchingFieldServiceTFObjectTypes},
		"name":        types.ObjectType{AttrTypes: identityRecordMatchingFieldServiceTFObjectTypes},
	}

	identityRecordMatchingFieldServiceTFObjectTypes = map[string]attr.Type{
		"field_required": types.BoolType,
		"threshold":      types.StringType,
	}

	verifyPolicyOptions = []validator.Object{
		objectvalidator.AtLeastOneOf(
			path.MatchRelative().AtParent().AtName("government_id"),
			path.MatchRelative().AtParent().AtName("facial_comparison"),
			path.MatchRelative().AtParent().AtName("liveness"),
			path.MatchRelative().AtParent().AtName("email"),
			path.MatchRelative().AtParent().AtName("phone"),
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

	const attrMinLifetimeDurationSeconds = 60
	const attrMaxLifetimeDurationSeconds = 1800
	const attrMinLifetimeDurationMinutes = 1
	const attrMaxLifetimeDurationMinutes = 30

	const attrMinRetryAttempts = 0
	const attrMaxRetryAttempts = 3

	const attrMinAadhaarDurationSeconds = 60
	const attrMaxAadhaarDurationSeconds = 1800
	const attrMinAadhaarDurationMinutes = 1
	const attrMaxAadhaarDurationMinutes = 30
	const attrMinAadhaarCount = 1
	const attrMaxAadhaarCount = 3

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

	const defaultBoolFalse = false
	const defaultBoolTrue = true

	const defaultProviderAuto = verify.ENUMPROVIDERAUTO_MITEK
	const defaultProviderManual = verify.ENUMPROVIDERAUTO_MITEK

	const defaultAadhaarEnabled = false

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies whether this is the environment's default verify policy.",
	)

	governmentIdVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls Government ID verification requirements.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	governmentIdInspectionTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Determine whether document authentication is automated, manual, or possibly both.",
	).AllowedValuesEnum(verify.AllowedEnumInspectionTypeEnumValues)

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

	retryAttemptsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Number of retries permitted when submitting images.  The allowed range is `%d - %d`.", attrMinRetryAttempts, attrMaxRetryAttempts),
	)

	providerAutoDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Provider to use for the automatic verification service.",
	).AllowedValuesEnum(verify.AllowedEnumProviderAutoEnumValues).DefaultValue(string(defaultProviderAuto))

	providerManualDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Provider to use for the manual verification service.",
	).AllowedValuesEnum(verify.AllowedEnumProviderManualEnumValues).DefaultValue(string(defaultProviderManual))

	identityRecordMatchingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Defines the verification requirements for identity record matching. If `government_id.verify` is `DISABLED`, then identity record matching is disabled.",
	)

	identityRecordMatchingThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Threshold for successful comparison.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues)

	verifyAamvaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"When enabled, the AAMVA DLDV system is used to validate identity documents issued by participating states. If license allows, defaults to `true` when `government_id.inspection_type` is `REQUIRED` or `OPTIONAL`; otherwise disabled.",
	)

	aadhaarDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Aadhaar configuration for India-based government Aadhaar documents;`facial_comparison.verify` must be `REQUIRED` to enable.",
	)

	aadhaarEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Whether Aadhaar verification is enabled.",
	).DefaultValue(defaultAadhaarEnabled)

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
			"- Identity Record Matching - Compare submitted biographic data (address, birth date, full name, given name, or family name) to an identity record.\n\n ",

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

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
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
						"verify":          types.StringValue(string(defaultVerify)),
						"inspection_type": types.StringNull(),
						"fail_expired_id": types.BoolValue(defaultBoolFalse),
						"provider_auto":   types.StringValue(string(defaultProviderAuto)),
						"provider_manual": types.StringValue(string(defaultProviderManual)),
						"retry_attempts":  types.Int32Null(),
						"verify_aamva":    types.BoolNull(),
						"aadhaar":         types.ObjectNull(aadhaarServiceTFObjectTypes),
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
					"inspection_type": schema.StringAttribute{
						Description:         governmentIdInspectionTypeDescription.Description,
						MarkdownDescription: governmentIdInspectionTypeDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Validators: []validator.String{
							stringvalidator.All(
								customstringvalidator.RegexMatchesPathValue(
									regexp.MustCompile(`REQUIRED|OPTIONAL`),
									fmt.Sprintf("The `government_id.inspection_type` argument is only applicable when `government_id.verify` does not have a value of %s.", string(verify.ENUMVERIFY_DISABLED)),
									path.MatchRelative().AtParent().AtName("verify"),
								),
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumInspectionTypeEnumValues)...),
							),
						},
					},
					"fail_expired_id": schema.BoolAttribute{
						Description: "When enabled, Government ID verification fails if the document is expired.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(defaultBoolFalse),
					},
					"provider_auto": schema.StringAttribute{
						Description:         providerAutoDescription.Description,
						MarkdownDescription: providerAutoDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(string(defaultProviderAuto)),
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumProviderAutoEnumValues)...),
						},
					},
					"provider_manual": schema.StringAttribute{
						Description:         providerManualDescription.Description,
						MarkdownDescription: providerManualDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(string(defaultProviderManual)),
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumProviderManualEnumValues)...),
						},
					},
					"retry_attempts": schema.Int32Attribute{
						Description:         retryAttemptsDescription.Description,
						MarkdownDescription: retryAttemptsDescription.MarkdownDescription,
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.Between(attrMinRetryAttempts, attrMaxRetryAttempts),
						},
					},
					"verify_aamva": schema.BoolAttribute{
						Description:         verifyAamvaDescription.Description,
						MarkdownDescription: verifyAamvaDescription.MarkdownDescription,
						Computed:            true,
						Optional:            true,
						Validators: []validator.Bool{
							boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
								types.StringValue(string(verify.ENUMVERIFY_DISABLED)),
								path.MatchRoot("government_id").AtName("verify"),
							),
						},
					},
					"aadhaar": schema.SingleNestedAttribute{
						Description:         aadhaarDescription.Description,
						MarkdownDescription: aadhaarDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         aadhaarEnabledDescription.Description,
								MarkdownDescription: aadhaarEnabledDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(defaultAadhaarEnabled),
								Validators: []validator.Bool{
									boolvalidator.All(
										boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
											types.StringValue(string(verify.ENUMVERIFY_DISABLED)),
											path.MatchRoot("facial_comparison").AtName("verify"),
										),
										boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
											types.StringValue(string(verify.ENUMVERIFY_OPTIONAL)),
											path.MatchRoot("facial_comparison").AtName("verify"),
										),
									),
								},
							},
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
						"verify":         types.StringValue(string(defaultVerify)),
						"threshold":      types.StringValue(string(defaultThreshold)),
						"retry_attempts": types.Int32Null(),
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
					"retry_attempts": schema.Int32Attribute{
						Description:         retryAttemptsDescription.Description,
						MarkdownDescription: retryAttemptsDescription.MarkdownDescription,
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.Between(attrMinRetryAttempts, attrMaxRetryAttempts),
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
						"count": types.Int32Value(defaultOTPAttemptsCount),
					}
					attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int32Value(defaultOTPEmailDuration),
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
						"duration":  types.Int32Value(defaultOTPCooldownDuration),
						"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
					}
					cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"count":    types.Int32Value(defaultOTPDeliveryCount),
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
						Computed:    true,

						Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
							o := map[string]attr.Value{
								"count": types.Int32Value(defaultOTPAttemptsCount),
							}
							attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
							resp.Diagnostics.Append(d...)

							o = map[string]attr.Value{
								"duration":  types.Int32Value(defaultOTPEmailDuration),
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
								"duration":  types.Int32Value(defaultOTPCooldownDuration),
								"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
							}
							cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
							resp.Diagnostics.Append(d...)

							o = map[string]attr.Value{
								"count":    types.Int32Value(defaultOTPDeliveryCount),
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

							return otpObjValue
						}()),

						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Required:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Required:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Required:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(attrMinDuration, attrMaxDurationMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(attrMinDuration, attrMaxDurationSeconds),
															int32validatorinternal.RegexMatchesPathValue(
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
									"duration": schema.Int32Attribute{
										Description:         otpLifeTimeEmailDurationDescription.Description,
										MarkdownDescription: otpLifeTimeEmailDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int32{
											int32validator.Any(
												int32validator.All(
													int32validator.Between(attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int32validator.All(
													int32validator.Between(attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
													int32validatorinternal.RegexMatchesPathValue(
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
						"count": types.Int32Value(defaultOTPAttemptsCount),
					}
					attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"duration":  types.Int32Value(defaultOTPPhoneDuration),
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
						"duration":  types.Int32Value(defaultOTPCooldownDuration),
						"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
					}
					cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					o = map[string]attr.Value{
						"count":    types.Int32Value(defaultOTPDeliveryCount),
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
						Computed:    true,

						Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
							o := map[string]attr.Value{
								"count": types.Int32Value(defaultOTPAttemptsCount),
							}
							attemptsObjValue, d := types.ObjectValue(otpAttemptsServiceTFObjectTypes, o)
							resp.Diagnostics.Append(d...)

							o = map[string]attr.Value{
								"duration":  types.Int32Value(defaultOTPPhoneDuration),
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
								"duration":  types.Int32Value(defaultOTPCooldownDuration),
								"time_unit": types.StringValue(string(defaultOTPCooldownTimeUnit)),
							}
							cooldownObjValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
							resp.Diagnostics.Append(d...)

							o = map[string]attr.Value{
								"count":    types.Int32Value(defaultOTPDeliveryCount),
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

							return otpObjValue
						}()),

						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Required:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Required:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Required:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.Int32{
													int32validator.Any(
														int32validator.All(
															int32validator.Between(attrMinDuration, attrMaxDurationMinutes),
															int32validatorinternal.RegexMatchesPathValue(
																regexp.MustCompile(`MINUTES`),
																fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
																path.MatchRelative().AtParent().AtName("time_unit"),
															),
														),
														int32validator.All(
															int32validator.Between(attrMinDuration, attrMaxDurationSeconds),
															int32validatorinternal.RegexMatchesPathValue(
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
									"duration": schema.Int32Attribute{
										Description:         otpLifeTimePhoneDurationDescription.Description,
										MarkdownDescription: otpLifeTimePhoneDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int32{
											int32validator.Any(
												int32validator.All(
													int32validator.Between(attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int32validator.All(
													int32validator.Between(attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
													int32validatorinternal.RegexMatchesPathValue(
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
							"duration": schema.Int32Attribute{
								Description:         transactionTimeoutDurationDescription.Description,
								MarkdownDescription: transactionTimeoutDurationDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.Int32{
									int32validator.Any(
										int32validator.All(
											int32validator.Between(attrMinDuration, attrMaxDurationMinutes),
											int32validatorinternal.RegexMatchesPathValue(
												regexp.MustCompile(`MINUTES`),
												fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
												path.MatchRelative().AtParent().AtName("time_unit"),
											),
										),
										int32validator.All(
											int32validator.Between(attrMinDuration, attrMaxDurationSeconds),
											int32validatorinternal.RegexMatchesPathValue(
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
									"duration": schema.Int32Attribute{
										Description:         dataCollectionDurationDescription.Description,
										MarkdownDescription: dataCollectionDurationDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.Int32{
											int32validatorinternal.IsLessThanEqualToPathValue(
												path.MatchRoot("transaction").AtName("timeout").AtName("duration"),
											),
											int32validator.Any(
												int32validator.All(
													int32validator.Between(attrMinDuration, attrMaxDurationMinutes),
													int32validatorinternal.RegexMatchesPathValue(
														regexp.MustCompile(`MINUTES`),
														fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinDuration, attrMaxDurationMinutes),
														path.MatchRelative().AtParent().AtName("time_unit"),
													),
												),
												int32validator.All(
													int32validator.Between(attrMinDuration, attrMaxDurationSeconds),
													int32validatorinternal.RegexMatchesPathValue(
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

			"identity_record_matching": schema.SingleNestedAttribute{
				Description:         identityRecordMatchingDescription.Description,
				MarkdownDescription: identityRecordMatchingDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"address": schema.SingleNestedAttribute{
						Description: "Configuration for address verification.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
								},
							},
						},
					},
					"birth_date": schema.SingleNestedAttribute{
						Description: "Configuration for birth date verification.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
								},
							},
						},
					},
					"family_name": schema.SingleNestedAttribute{
						Description: "Configuration for family name verification.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
								},
							},
						},
					},
					"given_name": schema.SingleNestedAttribute{
						Description: "Configuration for given name verification.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
								},
							},
						},
					},
					"name": schema.SingleNestedAttribute{
						Description: "Configuration for full name verification.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
								},
							},
						},
					},
				},

				Validators: []validator.Object{
					customobjectvalidator.ConflictsIfMatchesPathValue(
						types.StringValue(string(verify.ENUMVERIFY_DISABLED)),
						path.MatchRelative().AtParent().AtName("government_id").AtName("verify"),
					),
					customobjectvalidator.AtLeastOneChildConfigured(
						"address",
						"birth_date",
						"family_name",
						"given_name",
						"name",
					),
				},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify policy was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify policy was updated. Can be null.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *VerifyPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VerifyPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VerifyPoliciesApi.CreateVerifyPolicy(ctx, plan.EnvironmentId.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateVerifyPolicy",
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
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *verifyPolicyResourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VerifyPoliciesApi.ReadOneVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneVerifyPolicy",
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
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VerifyPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VerifyPoliciesApi.UpdateVerifyPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateVerifyPolicy",
		legacysdk.DefaultCustomError,
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

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
			fR, fErr := r.Client.VerifyAPIClient.VerifyPoliciesApi.DeleteVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteVerifyPolicy",
		verifyPolicyDeleteCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

var verifyPolicyDeleteCustomError = func(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Undeletable default verify policy
		if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
			if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
				if match, _ := regexp.MatchString("cannot delete default", v[0].GetMessage()); match {

					diags.AddWarning("Cannot delete the default verify policy", "Due to API restrictions, the provider cannot delete the default verify policy for an environment.  The policy has been removed from Terraform state but has been left in place in the PingOne service.")

					return diags
				}
			}
		}
	}

	diags.Append(legacysdk.CustomErrorResourceNotFoundWarning(r, p1Error)...)
	return diags
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

	// Top-level arguments
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	// The verify policy will also have default set to false.
	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	} else {
		data.SetDefault(false)
	}

	// Identity Record Matching Object
	if !p.IdentityRecordMatching.IsNull() && !p.IdentityRecordMatching.IsUnknown() {

		var identityRecordMatching identityRecordMatchingModel
		d := p.IdentityRecordMatching.As(ctx, &identityRecordMatching, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		identityRecordMatchingSettings, d := identityRecordMatching.expandIdentityRecordMatching(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetIdentityRecordMatching(*identityRecordMatchingSettings)
	}

	return data, diags
}

func (p *governmentIdModel) expandgovernmentIdModel() (*verify.GovernmentIdConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyGovernmentId := verify.NewGovernmentIdConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyGovernmentId.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.InspectionType.IsNull() && !p.InspectionType.IsUnknown() {
		verifyGovernmentId.SetInspectionType(verify.EnumInspectionType(p.InspectionType.ValueString()))
	}

	if !p.FailExpiredId.IsNull() && !p.FailExpiredId.IsUnknown() {
		verifyGovernmentId.SetFailExpiredId(p.FailExpiredId.ValueBool())
	}

	provider := verify.NewGovernmentIdConfigurationProviderWithDefaults()
	if !p.ProviderAuto.IsNull() && !p.ProviderAuto.IsUnknown() {
		provider.SetAuto(verify.EnumProviderAuto(p.ProviderAuto.ValueString()))
	}

	if !p.ProviderManual.IsNull() && !p.ProviderManual.IsUnknown() {
		provider.SetManual(verify.EnumProviderManual(p.ProviderManual.ValueString()))
	}

	if provider.HasAuto() || provider.HasManual() {
		verifyGovernmentId.SetProvider(*provider)
	}

	retryAttempts := verify.NewObjectRetryWithDefaults()
	if !p.RetryAttempts.IsNull() && !p.RetryAttempts.IsUnknown() {
		retryAttempts.SetAttempts(p.RetryAttempts.ValueInt32())
		verifyGovernmentId.SetRetry(*retryAttempts)
	}

	if !p.VerifyAamva.IsNull() && !p.VerifyAamva.IsUnknown() {
		verifyGovernmentId.SetVerifyAamva(p.VerifyAamva.ValueBool())
	}

	if !p.Aadhaar.IsNull() && !p.Aadhaar.IsUnknown() {
		var aadhaar aadhaarModel
		d := p.Aadhaar.As(context.Background(), &aadhaar, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		aadhaarConfig := verify.NewGovernmentIdConfigurationAadhaarWithDefaults()

		if !aadhaar.Enabled.IsNull() && !aadhaar.Enabled.IsUnknown() {
			aadhaarConfig.SetEnabled(aadhaar.Enabled.ValueBool())
		}

		verifyGovernmentId.SetAadhaar(*aadhaarConfig)
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

	retryAttempts := verify.NewObjectRetryWithDefaults()
	if !p.RetryAttempts.IsNull() && !p.RetryAttempts.IsUnknown() {
		retryAttempts.SetAttempts(p.RetryAttempts.ValueInt32())
		verifyLiveness.SetRetry(*retryAttempts)
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
			dataCollectionTimeout.SetDuration(genericTimeout.Duration.ValueInt32())
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
		transactionTimeout.SetDuration(genericTimeout.Duration.ValueInt32())

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
			otpAttempts.SetCount(attempts.Count.ValueInt32())
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
			otpDeliveries.SetCount(deliveries.Count.ValueInt32())
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
				deliveriesCooldown.SetDuration(cooldown.Duration.ValueInt32())
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
			otpLifeTime.SetDuration(genericTimeout.Duration.ValueInt32())
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

func (p *identityRecordMatchingModel) expandIdentityRecordMatching(ctx context.Context) (*verify.IdentityRecordMatching, diag.Diagnostics) {
	var diags diag.Diagnostics

	identityRecordMatching := verify.NewIdentityRecordMatchingWithDefaults()

	if !p.Address.IsNull() && !p.Address.IsUnknown() {
		var address identityRecordMatchingFieldModel
		d := p.Address.As(ctx, &address, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		addressSettings := verify.NewIdentityRecordMatchingAddress(
			address.FieldRequired.ValueBool(),
			verify.EnumThreshold(address.Threshold.ValueString()),
		)
		identityRecordMatching.SetAddress(*addressSettings)
	}

	if !p.BirthDate.IsNull() && !p.BirthDate.IsUnknown() {
		var birthDate identityRecordMatchingFieldModel
		d := p.BirthDate.As(ctx, &birthDate, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		birthDateSettings := verify.NewIdentityRecordMatchingBirthDate(
			birthDate.FieldRequired.ValueBool(),
			verify.EnumThreshold(birthDate.Threshold.ValueString()),
		)
		identityRecordMatching.SetBirthDate(*birthDateSettings)
	}

	if !p.FamilyName.IsNull() && !p.FamilyName.IsUnknown() {
		var familyName identityRecordMatchingFieldModel
		d := p.FamilyName.As(ctx, &familyName, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		familyNameSettings := verify.NewIdentityRecordMatchingFamilyName(
			familyName.FieldRequired.ValueBool(),
			verify.EnumThreshold(familyName.Threshold.ValueString()),
		)
		identityRecordMatching.SetFamilyName(*familyNameSettings)
	}

	if !p.GivenName.IsNull() && !p.GivenName.IsUnknown() {
		var givenName identityRecordMatchingFieldModel
		d := p.GivenName.As(ctx, &givenName, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		givenNameSettings := verify.NewIdentityRecordMatchingGivenName(
			givenName.FieldRequired.ValueBool(),
			verify.EnumThreshold(givenName.Threshold.ValueString()),
		)
		identityRecordMatching.SetGivenName(*givenNameSettings)
	}

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		var name identityRecordMatchingFieldModel
		d := p.Name.As(ctx, &name, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		nameSettings := verify.NewIdentityRecordMatchingName(
			name.FieldRequired.ValueBool(),
			verify.EnumThreshold(name.Threshold.ValueString()),
		)
		identityRecordMatching.SetName(*nameSettings)
	}

	return identityRecordMatching, diags
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

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
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

	p.IdentityRecordMatching, d = p.toStateIdentityRecordMatching(apiObject.GetIdentityRecordMatchingOk())
	diags.Append(d...)

	return diags
}

func (p *verifyPolicyResourceModel) toStateGovernmentId(apiObject *verify.GovernmentIdConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(governmentIdServiceTFObjectTypes), diags
	}

	retryAttempts := types.Int32Null()
	if v, ok := apiObject.GetRetryOk(); ok {
		if t, ok := v.GetAttemptsOk(); ok {
			retryAttempts = framework.Int32ToTF(*t)
		}
	}

	provider, ok := apiObject.GetProviderOk()
	if !ok {
		diags.AddError(
			"Unexpected Missing Value",
			"GovernmentId data object contained unexpected null value for the `provider` data object.  Please report this issue to the provider maintainers.")
	}

	aadhaarObject := types.ObjectNull(aadhaarServiceTFObjectTypes)
	if v, ok := apiObject.GetAadhaarOk(); ok {
		var d diag.Diagnostics

		aadhaarAttrs := map[string]attr.Value{
			"enabled": framework.BoolOkToTF(v.GetEnabledOk()),
		}

		objValue, d := types.ObjectValue(aadhaarServiceTFObjectTypes, aadhaarAttrs)
		diags.Append(d...)

		aadhaarObject = objValue
	}

	objValue, d := types.ObjectValue(governmentIdServiceTFObjectTypes, map[string]attr.Value{
		"verify":          framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"inspection_type": framework.EnumOkToTF(apiObject.GetInspectionTypeOk()),
		"fail_expired_id": framework.BoolOkToTF(apiObject.GetFailExpiredIdOk()),
		"provider_auto":   framework.EnumOkToTF(provider.GetAutoOk()),
		"provider_manual": framework.EnumOkToTF(provider.GetManualOk()),
		"retry_attempts":  retryAttempts,
		"verify_aamva":    framework.BoolOkToTF(apiObject.GetVerifyAamvaOk()),
		"aadhaar":         aadhaarObject,
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

	retryAttempts := types.Int32Null()
	if v, ok := apiObject.GetRetryOk(); ok {
		if t, ok := v.GetAttemptsOk(); ok {
			retryAttempts = framework.Int32ToTF(*t)
		}
	}

	objValue, d := types.ObjectValue(livenessServiceTFObjectTypes, map[string]attr.Value{
		"verify":         framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold":      framework.EnumOkToTF(apiObject.GetThresholdOk()),
		"retry_attempts": retryAttempts,
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

func (p *verifyPolicyResourceModel) toStateIdentityRecordMatching(apiObject *verify.IdentityRecordMatching, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityRecordMatchingServiceTFObjectTypes), diags
	}

	address, d := p.toStateIdentityRecordMatchingField(apiObject.GetAddressOk())
	diags.Append(d...)

	birthDate, d := p.toStateIdentityRecordMatchingField(apiObject.GetBirthDateOk())
	diags.Append(d...)

	familyName, d := p.toStateIdentityRecordMatchingField(apiObject.GetFamilyNameOk())
	diags.Append(d...)

	givenName, d := p.toStateIdentityRecordMatchingField(apiObject.GetGivenNameOk())
	diags.Append(d...)

	name, d := p.toStateIdentityRecordMatchingField(apiObject.GetNameOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(identityRecordMatchingServiceTFObjectTypes, map[string]attr.Value{
		"address":     address,
		"birth_date":  birthDate,
		"family_name": familyName,
		"given_name":  givenName,
		"name":        name,
	})
	diags.Append(d...)

	return objValue, diags
}

// IdentityRecordMatchingField defines the interface for identity record matching field types
type IdentityRecordMatchingField interface {
	GetFieldRequiredOk() (*bool, bool)
	GetThresholdOk() (*verify.EnumThreshold, bool)
}

// toStateIdentityRecordMatchingField converts any identity record matching field type to TF state
func (p *verifyPolicyResourceModel) toStateIdentityRecordMatchingField(apiObject IdentityRecordMatchingField, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityRecordMatchingFieldServiceTFObjectTypes), diags
	}

	fieldRequired, fieldRequiredOk := apiObject.GetFieldRequiredOk()
	threshold, thresholdOk := apiObject.GetThresholdOk()

	objValue, d := types.ObjectValue(identityRecordMatchingFieldServiceTFObjectTypes, map[string]attr.Value{
		"field_required": framework.BoolOkToTF(fieldRequired, fieldRequiredOk),
		"threshold":      framework.EnumOkToTF(threshold, thresholdOk),
	})
	diags.Append(d...)

	return objValue, diags
}
