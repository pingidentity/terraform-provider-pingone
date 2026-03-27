// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type VerifyPolicyDataSource serviceClientType

type verifyPolicyDataSourceModel struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	VerifyPolicyId         pingonetypes.ResourceIDValue `tfsdk:"verify_policy_id"`
	Name                   types.String                 `tfsdk:"name"`
	Default                types.Bool                   `tfsdk:"default"`
	Description            types.String                 `tfsdk:"description"`
	GovernmentId           types.Object                 `tfsdk:"government_id"`
	FacialComparison       types.Object                 `tfsdk:"facial_comparison"`
	Liveness               types.Object                 `tfsdk:"liveness"`
	Email                  types.Object                 `tfsdk:"email"`
	Phone                  types.Object                 `tfsdk:"phone"`
	Transaction            types.Object                 `tfsdk:"transaction"`
	Voice                  types.Object                 `tfsdk:"voice"`
	IdentityRecordMatching types.Object                 `tfsdk:"identity_record_matching"`
	CreatedAt              timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt              timetypes.RFC3339            `tfsdk:"updated_at"`
}

var (
	genericTimeoutDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int32Type,
		"time_unit": types.StringType,
	}

	aadhaarDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	governmentIdDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"verify":          types.StringType,
		"inspection_type": types.StringType,
		"fail_expired_id": types.BoolType,
		"provider_auto":   types.StringType,
		"provider_manual": types.StringType,
		"retry_attempts":  types.Int32Type,
		"verify_aamva":    types.BoolType,
		"aadhaar":         types.ObjectType{AttrTypes: aadhaarDataSourceServiceTFObjectTypes},
	}

	facialComparisonDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	livenessDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"verify":         types.StringType,
		"threshold":      types.StringType,
		"retry_attempts": types.Int32Type,
	}

	deviceDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"verify":            types.StringType,
		"create_mfa_device": types.BoolType,
		"otp":               types.ObjectType{AttrTypes: otpDataSourceServiceTFObjectTypes},
	}

	otpDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"attempts":     types.ObjectType{AttrTypes: otpAttemptsDataSourceServiceTFObjectTypes},
		"deliveries":   types.ObjectType{AttrTypes: otpDeliveriesDataSourceServiceTFObjectTypes},
		"lifetime":     types.ObjectType{AttrTypes: genericTimeoutDataSourceServiceTFObjectTypes},
		"notification": types.ObjectType{AttrTypes: otpNotificationDataSourceServiceTFObjectTypes},
	}

	otpAttemptsDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"count": types.Int32Type,
	}

	otpDeliveriesDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"count":    types.Int32Type,
		"cooldown": types.ObjectType{AttrTypes: genericTimeoutDataSourceServiceTFObjectTypes},
	}

	otpNotificationDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"template_name": types.StringType,
		"variant_name":  types.StringType,
	}

	transactionDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"timeout":              types.ObjectType{AttrTypes: genericTimeoutDataSourceServiceTFObjectTypes},
		"data_collection":      types.ObjectType{AttrTypes: dataCollectionDataSourceServiceTFObjectTypes},
		"data_collection_only": types.BoolType,
	}

	dataCollectionDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"timeout": types.ObjectType{AttrTypes: genericTimeoutDataSourceServiceTFObjectTypes},
	}

	voiceDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"verify":               types.StringType,
		"enrollment":           types.BoolType,
		"comparison_threshold": types.StringType,
		"liveness_threshold":   types.StringType,
		"text_dependent":       types.ObjectType{AttrTypes: textDependentDataSourceServiceTFObjectTypes},
		"reference_data":       types.ObjectType{AttrTypes: referenceDataDataSourceServiceTFObjectTypes},
	}

	textDependentDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"samples":         types.Int32Type,
		"voice_phrase_id": pingonetypes.ResourceIDType{},
	}

	referenceDataDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"retain_original_recordings": types.BoolType,
		"update_on_reenrollment":     types.BoolType,
		"update_on_verification":     types.BoolType,
	}

	identityRecordMatchingDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"address":     types.ObjectType{AttrTypes: identityRecordMatchingFieldDataSourceServiceTFObjectTypes},
		"birth_date":  types.ObjectType{AttrTypes: identityRecordMatchingFieldDataSourceServiceTFObjectTypes},
		"family_name": types.ObjectType{AttrTypes: identityRecordMatchingFieldDataSourceServiceTFObjectTypes},
		"given_name":  types.ObjectType{AttrTypes: identityRecordMatchingFieldDataSourceServiceTFObjectTypes},
		"name":        types.ObjectType{AttrTypes: identityRecordMatchingFieldDataSourceServiceTFObjectTypes},
	}

	identityRecordMatchingFieldDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"field_required": types.BoolType,
		"threshold":      types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &VerifyPolicyDataSource{}
)

// New Object
func NewVerifyPolicyDataSource() datasource.DataSource {
	return &VerifyPolicyDataSource{}
}

// Metadata
func (r *VerifyPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verify_policy"
}

func (r *VerifyPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	const attrMinDuration = 0
	const attrMaxDurationSeconds = 1800
	const attrMaxDurationMinutes = 30

	const attrMinVoiceSamples = 3
	const attrMaxVoiceSamples = 5

	const attrMinLifetimeDurationSeconds = 60
	const attrMaxLifetimeDurationSeconds = 1800
	const attrMinLifetimeDurationMinutes = 1
	const attrMaxLifetimeDurationMinutes = 30

	const attrMinRetryAttempts = 0
	const attrMaxRetryAttempts = 3

	// defaults
	const defaultNotificationTemplate = "email_phone_verification"

	const defaultVerify = verify.ENUMVERIFY_DISABLED
	const defaultThreshold = verify.ENUMTHRESHOLD_MEDIUM

	const defaultOTPEmailDuration = 10
	const defaultOTPEmailTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultOTPPhoneDuration = 5
	const defaultOTPPhoneTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultOTPCooldownDuration = 30
	const defaultOTPCooldownTimeUnit = verify.ENUMTIMEUNIT_SECONDS

	const defaultTransactionDuration = 30
	const defaultTransactionDataCollectionDuration = 15
	const defaultTransactionTimeUnit = verify.ENUMTIMEUNIT_MINUTES

	const defaultProviderAuto = verify.ENUMPROVIDERAUTO_MITEK
	const defaultProviderManual = verify.ENUMPROVIDERAUTO_MITEK

	const defaultAadhaarEnabled = false

	dataSourceExactlyOneOfRelativePaths := []string{
		"verify_policy_id",
		"name",
		"default",
	}

	verifyPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifier (UUID) associated with the verify policy.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	nametDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the verification policy displayed in PingOne Admin UI.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Set value to `true` to return the default verify policy. There is only one default policy per environment.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

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

	voiceVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls the verification requirements for a Voice verification.",
	).AllowedValuesEnum(verify.AllowedEnumVerifyEnumValues).DefaultValue(string(defaultVerify))

	voiceEnrollmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls if the transaction performs voice enrollment (`TRUE`) or voice verification (`FALSE`).",
	)

	voiceTexttDependentSamplesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Number of voice samples to collect. The allowed range is `%d - %d`", attrMinVoiceSamples, attrMaxVoiceSamples),
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

	retryAttemptsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Number of retries permitted when submitting images.  The allowed range is `%d - %d`.", attrMinRetryAttempts, attrMaxRetryAttempts),
	)

	providerAutoDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Provider to use for the automatic verification service.",
	).AllowedValuesEnum(verify.AllowedEnumProviderAutoEnumValues).DefaultValue(string(defaultProviderAuto))

	providerManualDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Provider to use for the manual verification service.",
	).AllowedValuesEnum(verify.AllowedEnumProviderManualEnumValues).DefaultValue(string(defaultProviderManual))

	identityRecordMatchingThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Threshold for successful comparison.",
	).AllowedValuesEnum(verify.AllowedEnumThresholdEnumValues)

	aadhaarDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Aadhaar configuration for India-based government Aadhaar documents;`facial_comparison.verify` must be `REQUIRED` to enable.",
	)

	aadhaarEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Whether Aadhaar verification is enabled.",
	).DefaultValue(defaultAadhaarEnabled)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve the default PingOne Verify Policy or to find a PingOne Verify Policy by its Verify Policy Id or Name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify policy exists."),
			),

			"verify_policy_id": schema.StringAttribute{
				Description:         verifyPolicyIdDescription.Description,
				MarkdownDescription: verifyPolicyIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("default"),
					),
				},
			},

			"name": schema.StringAttribute{
				Description:         nametDescription.Description,
				MarkdownDescription: nametDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("verify_policy_id"),
						path.MatchRelative().AtParent().AtName("default"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.Bool{
					boolvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("verify_policy_id"),
						path.MatchRelative().AtParent().AtName("name"),
					),
				},
			},

			"description": schema.StringAttribute{
				Description: "Description of the verification policy displayed in PingOne Admin UI, 1-1024 characters.",
				Computed:    true,
			},

			"government_id": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements for a government-issued identity document, which includes a photograph.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         governmentIdVerifyDescription.Description,
						MarkdownDescription: governmentIdVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
					"inspection_type": schema.StringAttribute{
						Description:         governmentIdInspectionTypeDescription.Description,
						MarkdownDescription: governmentIdInspectionTypeDescription.MarkdownDescription,
						Computed:            true,
					},
					"fail_expired_id": schema.BoolAttribute{
						Description: "When enabled, Government ID verification fails if the document is expired.",
						Computed:    true,
					},
					"provider_auto": schema.StringAttribute{
						Description:         providerAutoDescription.Description,
						MarkdownDescription: providerAutoDescription.MarkdownDescription,
						Computed:            true,
					},
					"provider_manual": schema.StringAttribute{
						Description:         providerManualDescription.Description,
						MarkdownDescription: providerManualDescription.MarkdownDescription,
						Computed:            true,
					},
					"retry_attempts": schema.Int32Attribute{
						Description:         retryAttemptsDescription.Description,
						MarkdownDescription: retryAttemptsDescription.MarkdownDescription,
						Computed:            true,
					},
					"verify_aamva": schema.BoolAttribute{
						Description: "When enabled, the AAMVA DLDV system is used to validate identity documents issued by participating states.",
						Computed:    true,
					},
					"aadhaar": schema.SingleNestedAttribute{
						Description:         aadhaarDescription.Description,
						MarkdownDescription: aadhaarDescription.MarkdownDescription,
						Computed:            true,

						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         aadhaarEnabledDescription.Description,
								MarkdownDescription: aadhaarEnabledDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
				},
			},

			"facial_comparison": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         facialComparisonVerifyDescription.Description,
						MarkdownDescription: facialComparisonVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
					"threshold": schema.StringAttribute{
						Description:         facialComparisonThresholdDescription.Description,
						MarkdownDescription: facialComparisonThresholdDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"liveness": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
					"threshold": schema.StringAttribute{
						Description:         livenessThresholdDescription.Description,
						MarkdownDescription: livenessThresholdDescription.MarkdownDescription,
						Computed:            true,
					},
					"retry_attempts": schema.Int32Attribute{
						Description:         retryAttemptsDescription.Description,
						MarkdownDescription: retryAttemptsDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"email": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to validate an email address using a one-time password (OTP).",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description: "When enabled, PingOne Verify registers the email address with PingOne MFA as a verified MFA device.",
						Computed:    true,
					},
					"otp": schema.SingleNestedAttribute{
						Description: "SMS/Voice/Email one-time password (OTP) configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Computed:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Computed:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Computed:            true,
											},
											"time_unit": schema.StringAttribute{
												Description:         otpDeliveriesCooldownTimeUnitDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownTimeUnitDescription.MarkdownDescription,
												Computed:            true,
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "The length of time for which the OTP is valid.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int32Attribute{
										Description:         otpLifeTimeEmailDurationDescription.Description,
										MarkdownDescription: otpLifeTimeEmailDurationDescription.MarkdownDescription,
										Computed:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         otpLifetimeEmailTimeUnitDescription.Description,
										MarkdownDescription: otpLifetimeEmailTimeUnitDescription.MarkdownDescription,
										Computed:            true,
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "OTP notification template configuration.",
								Computed:    true,

								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         otpNotificationTemplateDescription.Description,
										MarkdownDescription: otpNotificationTemplateDescription.MarkdownDescription,
										Computed:            true,
									},
									"variant_name": schema.StringAttribute{
										Description: "Name of the template variant to use to pass a one-time passcode (OTP).",
										Computed:    true,
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         deviceVerifyDescription.Description,
						MarkdownDescription: deviceVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"phone": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements to validate a mobile phone number using a one-time password (OTP).",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description: "When enabled, PingOne Verify registers the mobile phone with PingOne MFA as a verified MFA device.",
						Computed:    true,
					},
					"otp": schema.SingleNestedAttribute{
						Description: "SMS/Voice/Email one-time password (OTP) configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "OTP attempts configuration.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP failures.",
										Computed:    true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "OTP delivery configuration.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int32Attribute{
										Description: "Allowed maximum number of OTP deliveries.",
										Computed:    true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Cooldown (waiting period between OTP attempts) configuration.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int32Attribute{
												Description:         otpDeliveriesCooldownDurationDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownDurationDescription.MarkdownDescription,
												Computed:            true,
											},
											"time_unit": schema.StringAttribute{
												Description:         otpDeliveriesCooldownTimeUnitDescription.Description,
												MarkdownDescription: otpDeliveriesCooldownTimeUnitDescription.MarkdownDescription,
												Computed:            true,
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "The length of time for which the OTP is valid.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int32Attribute{
										Description:         otpLifeTimePhoneDurationDescription.Description,
										MarkdownDescription: otpLifeTimePhoneDurationDescription.MarkdownDescription,
										Computed:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         otpLifetimePhoneTimeUnitDescription.Description,
										MarkdownDescription: otpLifetimePhoneTimeUnitDescription.MarkdownDescription,
										Computed:            true,
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "OTP notification template configuration.",
								Computed:    true,

								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         otpNotificationTemplateDescription.Description,
										MarkdownDescription: otpNotificationTemplateDescription.MarkdownDescription,
										Computed:            true,
									},
									"variant_name": schema.StringAttribute{
										Description: "Name of the template variant to use to pass a one-time passcode (OTP).",
										Computed:    true,
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         deviceVerifyDescription.Description,
						MarkdownDescription: deviceVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"transaction": schema.SingleNestedAttribute{
				Description: "Defines the requirements for transactions invoked by the policy.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"timeout": schema.SingleNestedAttribute{
						Description: "Object for transaction timeout.",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int32Attribute{
								Description:         transactionTimeoutDurationDescription.Description,
								MarkdownDescription: transactionTimeoutDurationDescription.MarkdownDescription,
								Computed:            true,
							},
							"time_unit": schema.StringAttribute{
								Description:         transactionTimeoutTimeUnitDescription.Description,
								MarkdownDescription: transactionTimeoutTimeUnitDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
					"data_collection": schema.SingleNestedAttribute{
						Description: "Object for data collection timeout definition.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"timeout": schema.SingleNestedAttribute{
								Description: "Object for data collection timeout.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int32Attribute{
										Description:         dataCollectionDurationDescription.Description,
										MarkdownDescription: dataCollectionDurationDescription.MarkdownDescription,
										Computed:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         dataCollectionTimeoutTimeUnitDescription.Description,
										MarkdownDescription: dataCollectionTimeoutTimeUnitDescription.MarkdownDescription,
										Computed:            true,
									},
								},
							},
						},
					},
					"data_collection_only": schema.BoolAttribute{
						Description:         dataCollectionOnlyDescription.Description,
						MarkdownDescription: dataCollectionOnlyDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"voice": schema.SingleNestedAttribute{
				Description:        "**[Deprecation notice: This field is deprecated and will be removed in a future release. Please use alternative verification methods.]** Defines the requirements for transactions invoked by the policy.",
				DeprecationMessage: "Deprecation notice: This field is deprecated and will be removed in a future release. Please use alternative verification methods.",
				Computed:           true, Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         voiceVerifyDescription.Description,
						MarkdownDescription: voiceVerifyDescription.MarkdownDescription,
						Computed:            true,
					},
					"enrollment": schema.BoolAttribute{
						Description:         voiceEnrollmentDescription.Description,
						MarkdownDescription: voiceEnrollmentDescription.MarkdownDescription,
						Computed:            true,
					},
					"comparison_threshold": schema.StringAttribute{
						Description:         voiceComparisonThresholdDescription.Description,
						MarkdownDescription: voiceComparisonThresholdDescription.MarkdownDescription,
						Computed:            true,
					},
					"liveness_threshold": schema.StringAttribute{
						Description:         voiceLivenessThresholdDescription.Description,
						MarkdownDescription: voiceLivenessThresholdDescription.MarkdownDescription,
						Computed:            true,
					},
					"text_dependent": schema.SingleNestedAttribute{
						Description: "Object for configuration of text dependent voice verification.",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"samples": schema.Int32Attribute{
								Description:         voiceTexttDependentSamplesDescription.Description,
								MarkdownDescription: voiceTexttDependentSamplesDescription.MarkdownDescription,
								Computed:            true,
							},
							"voice_phrase_id": schema.StringAttribute{
								Description: "	Identifier (UUID) of the voice phrase to use.",
								Computed:    true,

								CustomType: pingonetypes.ResourceIDType{},
							},
						},
					},
					"reference_data": schema.SingleNestedAttribute{
						Description: "Object for configuration of voice recording reference data.",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"retain_original_recordings": schema.BoolAttribute{
								Description: "Controls if the service stores the original voice recordings.",
								Computed:    true,
							},
							"update_on_reenrollment": schema.BoolAttribute{
								Description:         referenceDataUpdateOnEnrollmentDescription.Description,
								MarkdownDescription: referenceDataUpdateOnEnrollmentDescription.MarkdownDescription,
								Computed:            true,
							},
							"update_on_verification": schema.BoolAttribute{
								Description:         referenceDataUpdateOnVerificationDescription.Description,
								MarkdownDescription: referenceDataUpdateOnVerificationDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
				},
			},

			"identity_record_matching": schema.SingleNestedAttribute{
				Description: "Defines the verification requirements for identity record matching.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"address": schema.SingleNestedAttribute{
						Description: "Configuration for address verification.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Computed:    true,
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
					"birth_date": schema.SingleNestedAttribute{
						Description: "Configuration for birth date verification.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Computed:    true,
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
					"family_name": schema.SingleNestedAttribute{
						Description: "Configuration for family name verification.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Computed:    true,
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
					"given_name": schema.SingleNestedAttribute{
						Description: "Configuration for given name verification.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Computed:    true,
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
					"name": schema.SingleNestedAttribute{
						Description: "Configuration for full name verification.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"field_required": schema.BoolAttribute{
								Description: "Whether the field is required.",
								Computed:    true,
							},
							"threshold": schema.StringAttribute{
								Description:         identityRecordMatchingThresholdDescription.Description,
								MarkdownDescription: identityRecordMatchingThresholdDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},
				},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify policy was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify policy was updated. Can be null.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *VerifyPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *VerifyPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *verifyPolicyDataSourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var verifyPolicy *verify.VerifyPolicy

	if !data.VerifyPolicyId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.VerifyAPIClient.VerifyPoliciesApi.ReadOneVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.VerifyPolicyId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneVerifyPolicy",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&verifyPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.VerifyAPIClient.VerifyPoliciesApi.ReadAllVerifyPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if verifyPolicies, ok := pageCursor.EntityArray.Embedded.GetVerifyPoliciesOk(); ok {

						for _, verifyPolicyItem := range verifyPolicies {

							if strings.EqualFold(verifyPolicyItem.GetName(), data.Name.ValueString()) {
								return &verifyPolicyItem, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllVerifyPolicies",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&verifyPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if verifyPolicy == nil {
			resp.Diagnostics.AddError(
				"Cannot find verify policy from name",
				fmt.Sprintf("The verify policy name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else if data.Default.ValueBool() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.VerifyAPIClient.VerifyPoliciesApi.ReadAllVerifyPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if verifyPolicies, ok := pageCursor.EntityArray.Embedded.GetVerifyPoliciesOk(); ok {

						for _, verifyPolicyItem := range verifyPolicies {

							if verifyPolicyItem.GetDefault() {
								return &verifyPolicyItem, pageCursor.HTTPResponse, nil
							}
						}

					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllVerifyPolicies",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&verifyPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if verifyPolicy == nil {
			resp.Diagnostics.AddError(
				"Cannot find default verify policy",
				fmt.Sprintf("The default verify policy for environment %s cannot be found", data.EnvironmentId.String()),
			)
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Verify Policy: verify_policy_id, name, or default argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(verifyPolicy)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *verifyPolicyDataSourceModel) toState(apiObject *verify.VerifyPolicy) diag.Diagnostics {
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
	p.VerifyPolicyId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
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

	p.IdentityRecordMatching, d = p.toStateIdentityRecordMatching(apiObject.GetIdentityRecordMatchingOk())
	diags.Append(d...)

	return diags
}

func (p *verifyPolicyDataSourceModel) toStateGovernmentId(apiObject *verify.GovernmentIdConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(governmentIdDataSourceServiceTFObjectTypes), diags
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

	aadhaarObject := types.ObjectNull(aadhaarDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetAadhaarOk(); ok {
		var d diag.Diagnostics

		aadhaarAttrs := map[string]attr.Value{
			"enabled": framework.BoolOkToTF(v.GetEnabledOk()),
		}

		objValue, d := types.ObjectValue(aadhaarDataSourceServiceTFObjectTypes, aadhaarAttrs)
		diags.Append(d...)

		aadhaarObject = objValue
	}

	objValue, d := types.ObjectValue(governmentIdDataSourceServiceTFObjectTypes, map[string]attr.Value{
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

func (p *verifyPolicyDataSourceModel) toStateFacialComparison(apiObject *verify.FacialComparisonConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(facialComparisonDataSourceServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(facialComparisonDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"verify":    framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold": framework.EnumOkToTF(apiObject.GetThresholdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyDataSourceModel) toStateLiveness(apiObject *verify.LivenessConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(livenessDataSourceServiceTFObjectTypes), diags
	}

	retryAttempts := types.Int32Null()
	if v, ok := apiObject.GetRetryOk(); ok {
		if t, ok := v.GetAttemptsOk(); ok {
			retryAttempts = framework.Int32ToTF(*t)
		}
	}

	objValue, d := types.ObjectValue(livenessDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"verify":         framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold":      framework.EnumOkToTF(apiObject.GetThresholdOk()),
		"retry_attempts": retryAttempts,
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyDataSourceModel) toStateTransaction(apiObject *verify.TransactionConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(transactionDataSourceServiceTFObjectTypes), diags
	}

	transactionTimeout := types.ObjectNull(genericTimeoutDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetTimeoutOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"duration":  framework.Int32OkToTF(v.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(v.GetTimeUnitOk()),
		}

		objValue, d := types.ObjectValue(genericTimeoutDataSourceServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionTimeout = objValue
	}

	transactionDataCollection := types.ObjectNull(dataCollectionDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetDataCollectionOk(); ok {
		var d diag.Diagnostics

		transactionDataCollectionTimeout := types.ObjectNull(genericTimeoutDataSourceServiceTFObjectTypes)
		if t, ok := v.GetTimeoutOk(); ok {
			o := map[string]attr.Value{
				"duration":  framework.Int32OkToTF(t.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(t.GetTimeUnitOk()),
			}

			objValue, d := types.ObjectValue(genericTimeoutDataSourceServiceTFObjectTypes, o)
			diags.Append(d...)

			transactionDataCollectionTimeout = objValue
		}

		o := map[string]attr.Value{
			"timeout": transactionDataCollectionTimeout,
		}

		objValue, d := types.ObjectValue(dataCollectionDataSourceServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionDataCollection = objValue

	}

	objValue, d := types.ObjectValue(transactionDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"timeout":              transactionTimeout,
		"data_collection":      transactionDataCollection,
		"data_collection_only": framework.BoolOkToTF(apiObject.GetDataCollectionOnlyOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyDataSourceModel) toStateDevice(apiObject *verify.OTPDeviceConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(deviceDataSourceServiceTFObjectTypes), diags
	}

	otp := types.ObjectNull(otpDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetOtpOk(); ok {
		var d diag.Diagnostics

		attempts := types.ObjectNull(otpAttemptsDataSourceServiceTFObjectTypes)
		if t, ok := v.GetAttemptsOk(); ok {
			o := map[string]attr.Value{
				"count": framework.Int32OkToTF(t.GetCountOk()),
			}

			objValue, d := types.ObjectValue(otpAttemptsDataSourceServiceTFObjectTypes, o)
			diags.Append(d...)

			attempts = objValue
		}

		deliveries := types.ObjectNull(otpDeliveriesDataSourceServiceTFObjectTypes)
		if t, ok := v.GetDeliveriesOk(); ok {

			cooldown := types.ObjectNull(genericTimeoutDataSourceServiceTFObjectTypes)
			if c, ok := t.GetCooldownOk(); ok {
				o := map[string]attr.Value{
					"duration":  framework.Int32OkToTF(c.GetDurationOk()),
					"time_unit": framework.EnumOkToTF(c.GetTimeUnitOk()),
				}
				objValue, d := types.ObjectValue(genericTimeoutDataSourceServiceTFObjectTypes, o)
				diags.Append(d...)

				cooldown = objValue
			}

			o := map[string]attr.Value{
				"count":    framework.Int32OkToTF(t.GetCountOk()),
				"cooldown": cooldown,
			}
			objValue, d := types.ObjectValue(otpDeliveriesDataSourceServiceTFObjectTypes, o)
			diags.Append(d...)

			deliveries = objValue
		}

		lifetime := types.ObjectNull(genericTimeoutDataSourceServiceTFObjectTypes)
		if t, ok := v.GetLifeTimeOk(); ok {
			o := map[string]attr.Value{
				"duration":  framework.Int32OkToTF(t.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(t.GetTimeUnitOk()),
			}

			objValue, d := types.ObjectValue(genericTimeoutDataSourceServiceTFObjectTypes, o)
			diags.Append(d...)

			lifetime = objValue
		}

		notification := types.ObjectNull(otpNotificationDataSourceServiceTFObjectTypes)
		if t, ok := v.GetNotificationOk(); ok {
			o := map[string]attr.Value{
				"template_name": framework.StringOkToTF(t.GetTemplateNameOk()),
				"variant_name":  framework.StringOkToTF(t.GetVariantNameOk()),
			}

			objValue, d := types.ObjectValue(otpNotificationDataSourceServiceTFObjectTypes, o)
			diags.Append(d...)

			notification = objValue
		}

		o := map[string]attr.Value{
			"attempts":     attempts,
			"lifetime":     lifetime,
			"deliveries":   deliveries,
			"notification": notification,
		}
		objValue, d := types.ObjectValue(otpDataSourceServiceTFObjectTypes, o)
		diags.Append(d...)

		otp = objValue
	}

	objValue, d := types.ObjectValue(deviceDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"verify":            framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"create_mfa_device": framework.BoolOkToTF(apiObject.GetCreateMfaDeviceOk()),
		"otp":               otp,
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyDataSourceModel) toStateVoice(apiObject *verify.VoiceConfiguration, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(voiceDataSourceServiceTFObjectTypes), diags
	}

	textDependent := types.ObjectNull(textDependentDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetTextDependentOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"samples":         framework.Int32OkToTF(v.GetSamplesOk()),
			"voice_phrase_id": framework.PingOneResourceIDToTF(v.GetPhrase().Id),
		}
		objValue, d := types.ObjectValue(textDependentDataSourceServiceTFObjectTypes, o)
		diags.Append(d...)

		textDependent = objValue
	}

	referenceData := types.ObjectNull(referenceDataDataSourceServiceTFObjectTypes)
	if v, ok := apiObject.GetReferenceDataOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"retain_original_recordings": framework.BoolOkToTF(v.GetRetainOriginalRecordingsOk()),
			"update_on_reenrollment":     framework.BoolOkToTF(v.GetUpdateOnReenrollmentOk()),
			"update_on_verification":     framework.BoolOkToTF(v.GetUpdateOnVerificationOk()),
		}
		objValue, d := types.ObjectValue(referenceDataDataSourceServiceTFObjectTypes, o)
		diags.Append(d...)

		referenceData = objValue
	}

	objValue, d := types.ObjectValue(voiceDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"verify":               framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"enrollment":           framework.BoolOkToTF(apiObject.GetEnrollmentOk()),
		"comparison_threshold": framework.EnumOkToTF(apiObject.GetComparison().Threshold, ok),
		"liveness_threshold":   framework.EnumOkToTF(apiObject.GetLiveness().Threshold, ok),
		"text_dependent":       textDependent,
		"reference_data":       referenceData,
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyDataSourceModel) toStateIdentityRecordMatching(apiObject *verify.IdentityRecordMatching, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityRecordMatchingDataSourceServiceTFObjectTypes), diags
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

	objValue, d := types.ObjectValue(identityRecordMatchingDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"address":     address,
		"birth_date":  birthDate,
		"family_name": familyName,
		"given_name":  givenName,
		"name":        name,
	})
	diags.Append(d...)

	return objValue, diags
}

// toStateIdentityRecordMatchingField converts any identity record matching field type to TF state for data source
func (p *verifyPolicyDataSourceModel) toStateIdentityRecordMatchingField(apiObject IdentityRecordMatchingField, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityRecordMatchingFieldDataSourceServiceTFObjectTypes), diags
	}

	fieldRequired, fieldRequiredOk := apiObject.GetFieldRequiredOk()
	threshold, thresholdOk := apiObject.GetThresholdOk()

	objValue, d := types.ObjectValue(identityRecordMatchingFieldDataSourceServiceTFObjectTypes, map[string]attr.Value{
		"field_required": framework.BoolOkToTF(fieldRequired, fieldRequiredOk),
		"threshold":      framework.EnumOkToTF(threshold, thresholdOk),
	})
	diags.Append(d...)

	return objValue, diags
}
