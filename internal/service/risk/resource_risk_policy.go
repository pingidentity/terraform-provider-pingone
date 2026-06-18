// Copyright © 2026 Ping Identity Corporation

package risk

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type RiskPolicyResource serviceClientType

type riskPolicyResourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                types.String                 `tfsdk:"name"`
	DefaultResult       types.Object                 `tfsdk:"default_result"`
	Default             types.Bool                   `tfsdk:"default"`
	EvaluatedPredictors types.Set                    `tfsdk:"evaluated_predictors"`
	PolicyWeights       types.Object                 `tfsdk:"policy_weights"`
	PolicyScores        types.Object                 `tfsdk:"policy_scores"`
	Overrides           types.List                   `tfsdk:"overrides"`
	Mitigations         types.List                   `tfsdk:"mitigations"`
	Fallback            types.Object                 `tfsdk:"fallback"`
	Targets             types.Object                 `tfsdk:"targets"`
}

type riskPolicyResourceDefaultResultModel struct {
	Level types.String `tfsdk:"level"`
	Type  types.String `tfsdk:"type"`
}

type riskPolicyResourcePolicyModel struct {
	PolicyThresholdMedium types.Object `tfsdk:"policy_threshold_medium"`
	PolicyThresholdHigh   types.Object `tfsdk:"policy_threshold_high"`
	Predictors            types.Set    `tfsdk:"predictors"`
}

type riskPolicyResourcePolicyThresholdScoreBetweenModel struct {
	MinScore types.Int32 `tfsdk:"min_score"`
	MaxScore types.Int32 `tfsdk:"max_score"`
}

type riskPolicyResourcePolicyWeightsPredictorModel struct {
	CompactName             types.String `tfsdk:"compact_name"`
	PredictorReferenceValue types.String `tfsdk:"predictor_reference_value"`
	Weight                  types.Int32  `tfsdk:"weight"`
}

type riskPolicyResourcePolicyScoresPredictorModel struct {
	CompactName             types.String `tfsdk:"compact_name"`
	PredictorReferenceValue types.String `tfsdk:"predictor_reference_value"`
	Score                   types.Int32  `tfsdk:"score"`
}

type riskPolicyResourcePolicyOverrideModel struct {
	Name      types.String `tfsdk:"name"`
	Priority  types.Int32  `tfsdk:"priority"`
	Result    types.Object `tfsdk:"result"`
	Condition types.Object `tfsdk:"condition"`
}

type riskPolicyResourcePolicyOverrideResultModel struct {
	Value types.String `tfsdk:"value"`
	Level types.String `tfsdk:"level"`
	Type  types.String `tfsdk:"type"`
}

type riskPolicyResourcePolicyOverrideConditionModel struct {
	Type                       types.String `tfsdk:"type"`
	Equals                     types.String `tfsdk:"equals"`
	CompactName                types.String `tfsdk:"compact_name"`
	PredictorReferenceValue    types.String `tfsdk:"predictor_reference_value"`
	IPRange                    types.Set    `tfsdk:"ip_range"`
	PredictorReferenceContains types.String `tfsdk:"predictor_reference_contains"`
}

type riskPolicyResourcePolicyMitigationModel struct {
	Name                      types.String                 `tfsdk:"name"`
	Priority                  types.Int32                  `tfsdk:"priority"`
	Condition                 types.Object                 `tfsdk:"condition"`
	Action                    types.String                 `tfsdk:"action"`
	CustomAction              types.String                 `tfsdk:"custom_action"`
	MfaAuthenticationPolicyId pingonetypes.ResourceIDValue `tfsdk:"mfa_authentication_policy_id"`
	MfaRegistrationPolicyId   pingonetypes.ResourceIDValue `tfsdk:"mfa_registration_policy_id"`
	VerifyPolicyId            pingonetypes.ResourceIDValue `tfsdk:"verify_policy_id"`
}

type riskPolicyResourceMitigationFallbackModel struct {
	Action                    types.String                 `tfsdk:"action"`
	CustomAction              types.String                 `tfsdk:"custom_action"`
	MfaAuthenticationPolicyId pingonetypes.ResourceIDValue `tfsdk:"mfa_authentication_policy_id"`
	MfaRegistrationPolicyId   pingonetypes.ResourceIDValue `tfsdk:"mfa_registration_policy_id"`
	VerifyPolicyId            pingonetypes.ResourceIDValue `tfsdk:"verify_policy_id"`
}

type riskPolicyResourceTargetsModel struct {
	Condition types.Object `tfsdk:"condition"`
}

type riskPolicyResourceTargetsConditionModel struct {
	And types.List `tfsdk:"and"`
}

type riskPolicyResourceTargetsConditionAndModel struct {
	Type     types.String `tfsdk:"type"`
	List     types.List   `tfsdk:"list"`
	Contains types.String `tfsdk:"contains"`
}

var (
	policyThresholdsTFObjectTypes = map[string]attr.Type{
		"min_score": types.Int32Type,
		"max_score": types.Int32Type,
	}

	// Weights
	policyWeightsTFObjectTypes = map[string]attr.Type{
		"policy_threshold_medium": types.ObjectType{
			AttrTypes: policyThresholdsTFObjectTypes,
		},
		"policy_threshold_high": types.ObjectType{
			AttrTypes: policyThresholdsTFObjectTypes,
		},
		"predictors": types.SetType{
			ElemType: types.ObjectType{AttrTypes: policyWeightsPredictorTFObjectTypes},
		},
	}

	policyWeightsPredictorTFObjectTypes = map[string]attr.Type{
		"compact_name":              types.StringType,
		"predictor_reference_value": types.StringType,
		"weight":                    types.Int32Type,
	}

	// Scores
	policyScoresTFObjectTypes = map[string]attr.Type{
		"policy_threshold_medium": types.ObjectType{
			AttrTypes: policyThresholdsTFObjectTypes,
		},
		"policy_threshold_high": types.ObjectType{
			AttrTypes: policyThresholdsTFObjectTypes,
		},
		"predictors": types.SetType{
			ElemType: types.ObjectType{AttrTypes: policyScoresPredictorTFObjectTypes},
		},
	}

	policyScoresPredictorTFObjectTypes = map[string]attr.Type{
		"compact_name":              types.StringType,
		"predictor_reference_value": types.StringType,
		"score":                     types.Int32Type,
	}

	overridesTFObjectTypes = map[string]attr.Type{
		"name":     types.StringType,
		"priority": types.Int32Type,
		"result": types.ObjectType{
			AttrTypes: overridesResultTFObjectTypes,
		},
		"condition": types.ObjectType{
			AttrTypes: overridesConditionTFObjectTypes,
		},
	}

	overridesResultTFObjectTypes = map[string]attr.Type{
		"value": types.StringType,
		"level": types.StringType,
		"type":  types.StringType,
	}

	overridesConditionTFObjectTypes = map[string]attr.Type{
		"type":                         types.StringType,
		"equals":                       types.StringType,
		"compact_name":                 types.StringType,
		"predictor_reference_value":    types.StringType,
		"ip_range":                     types.SetType{ElemType: types.StringType},
		"predictor_reference_contains": types.StringType,
	}

	// Mitigations
	mitigationsTFObjectTypes = map[string]attr.Type{
		"name":     types.StringType,
		"priority": types.Int32Type,
		"condition": types.ObjectType{
			AttrTypes: overridesConditionTFObjectTypes,
		},
		"action":                       types.StringType,
		"custom_action":                types.StringType,
		"mfa_authentication_policy_id": pingonetypes.ResourceIDType{},
		"mfa_registration_policy_id":   pingonetypes.ResourceIDType{},
		"verify_policy_id":             pingonetypes.ResourceIDType{},
	}

	mitigationsFallbackTFObjectTypes = map[string]attr.Type{
		"action":                       types.StringType,
		"custom_action":                types.StringType,
		"mfa_authentication_policy_id": pingonetypes.ResourceIDType{},
		"mfa_registration_policy_id":   pingonetypes.ResourceIDType{},
		"verify_policy_id":             pingonetypes.ResourceIDType{},
	}

	// Targets
	targetsConditionAndTFObjectTypes = map[string]attr.Type{
		"type":     types.StringType,
		"list":     types.ListType{ElemType: types.StringType},
		"contains": types.StringType,
	}

	targetsConditionTFObjectTypes = map[string]attr.Type{
		"and": types.ListType{
			ElemType: types.ObjectType{AttrTypes: targetsConditionAndTFObjectTypes},
		},
	}

	targetsTFObjectTypes = map[string]attr.Type{
		"condition": types.ObjectType{
			AttrTypes: targetsConditionTFObjectTypes,
		},
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &RiskPolicyResource{}
	_ resource.ResourceWithConfigure   = &RiskPolicyResource{}
	_ resource.ResourceWithImportState = &RiskPolicyResource{}
	_ resource.ResourceWithModifyPlan  = &RiskPolicyResource{}
)

// New Object
func NewRiskPolicyResource() resource.Resource {
	return &RiskPolicyResource{}
}

// Metadata
func (r *RiskPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_policy"
}

// Schema
func (r *RiskPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const attrNameMaxLength = 256
	const emailAddressMaxLength = 5
	const attrDescriptionMaxLength = 1024
	const defaultWeightValue = 5

	const weightMinimumDefault = 1
	const weightMaximumDefault = 10

	const scoreMinimumDefault = 0
	const scoreMaximumDefault = 100

	defaultResultTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The default result type.",
	).AllowedValuesEnum([]risk.EnumResultType{risk.ENUMRESULTTYPE_VALUE})

	defaultResultLevelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The default result level.",
	).AllowedValuesEnum([]risk.EnumRiskLevel{risk.ENUMRISKLEVEL_LOW})

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that indicates whether this risk policy set is the environment's default risk policy set. This is used whenever an explicit policy set ID is not specified in a risk evaluation request.",
	)

	// Weighted Average Policy
	policyWeightedAverageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes settings for a risk policy using a weighted average calculation, with a final result being a risk score between `0` and `10`.",
	).ExactlyOneOf([]string{"policy_weights", "policy_scores"})

	policyWeightedAveragePredictor := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes a predictor to apply to the risk policy and its associated weight value for the overall weighted average risk calculation.",
	)

	// Scores policy
	policyScoresDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes settings for a risk policy calculated by aggregating score values, with a final result being the sum of score values from each of the configured predictors.",
	).ExactlyOneOf([]string{"policy_weights", "policy_scores"})

	policyScoresPredictor := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes a predictor to apply to the risk policy and its associated high risk / true outcome score to apply to the risk calculation.",
	)

	policyOverrideDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of policy overrides to apply to the policy.  The ordering of the overrides is important as it determines the priority of the policy override during policy evaluation.",
	)

	policyOverrideResultLevelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the risk level that should be applied to the policy evalution result when the override condition is met.",
	).AllowedValuesEnum(risk.AllowedEnumRiskLevelEnumValues)

	policyOverrideResultTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the risk result should be applied to the policy evalution result when the override condition is met.",
	).AllowedValuesEnum([]risk.EnumResultType{risk.ENUMRESULTTYPE_VALUE}).DefaultValue(string(risk.ENUMRESULTTYPE_VALUE))

	policyOverrideConditionTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the override condition to evaluate.",
	).AllowedValues(
		string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON),
		string(risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE),
	)

	policyOverrideConditionEqualsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Required when `equals` is set to `VALUE_COMPARISON`.  A string that specifies the value of the `predictor_reference_value` that must be matched for the override result to be applied to the policy evaluation.",
	)

	policyOverrideConditionCompactNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Required when `equals` is set to `VALUE_COMPARISON`.  A string that specifies the compact name of the predictor to apply to the override condition.",
	)

	policyOverrideConditionIPRangeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Required when `equals` is set to `IP_RANGE`.  A set of strings that specifies the CIDR ranges that should be evaluated against the value of the `predictor_reference_contains` attribute, that must be matched for the override result to be applied to the policy evaluation.  Values must be valid IPv4 or IPv6 CIDR ranges.",
	)

	// Mitigations
	policyMitigationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of mitigation-style policy entries to apply to the policy. Each entry pairs a condition with a single mitigation action. Mutually exclusive with `overrides`. When this block is configured, a `fallback` block must also be configured.",
	)

	policyMitigationActionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the mitigation action to apply when the condition is met.",
	).AllowedValuesEnum(risk.AllowedEnumMitigationActionEnumValues)

	policyMitigationCustomActionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom action name. Required when `action` is `CUSTOM`.",
	)

	policyMitigationMfaAuthPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the MFA (sign-on/authentication) policy to apply. Required when `action` is `MFA`.",
	)

	policyMitigationMfaRegPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the MFA registration policy to apply. Applies to MFA registration flows when `action` is `MFA`.",
	)

	policyMitigationVerifyPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the PingOne Verify policy to apply. Required when `action` is `VERIFY`.",
	)

	policyFallbackDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the required catch-all fallback mitigation entry (`result.type=MITIGATION_FALLBACK`). Required when `mitigations` is configured. Carries a single mitigation action with no condition.",
	)

	// Targets
	policyTargetsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that scopes this policy set to a subset of events (targeted policy). Pairs with `mitigations` (and the weights/scores backbone) but is mutually exclusive with `overrides`.",
	)

	policyTargetsConditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the AND-of-sub-conditions targeting condition. All sub-conditions in `and` must be satisfied for the policy set to be selected.",
	)

	policyTargetsConditionAndDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of sub-conditions that are combined with AND logic. Each entry pairs a `list` of values with the event attribute (`contains`) to check against.",
	)

	policyTargetsConditionAndTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A read-only string that identifies the sub-condition kind. Inferred by the API from `contains`.",
	).AllowedValuesEnum(risk.AllowedEnumRiskPolicySetTargetsConditionTypeEnumValues)

	policyTargetsConditionAndListDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of values to match against the event attribute specified in `contains`. Transaction types are one or more of `REGISTRATION`, `AUTHENTICATION`, `ACCESS`, `AUTHORIZATION`, `TRANSACTION`. User groups are group names. Applications are PingOne application IDs.",
	)

	policyTargetsConditionAndContainsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The event attribute checked against `list`. For transaction types use `${event.flow.type}`; for user groups use `${event.user.groups}`; for applications use `${event.targetResource.id}`.",
	)

	// Schema
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage Risk policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the risk policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique, friendly name for this policy set. Valid characters consist of any Unicode letter, mark (such as, accent, umlaut), # (numeric), / (forward slash), . (period), ' (apostrophe), _ (underscore), space, or - (hyphen). Maximum size is 256 characters.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrNameMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}0-9\/.'_\s-]+$`), " Valid characters consist of any Unicode letter, mark (such as, accent, umlaut), # (numeric), / (forward slash), . (period), ' (apostrophe), _ (underscore), space, or - (hyphen)."),
				},
			},

			"default_result": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies the default result value for the risk policy.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"level": framework.StringToTF(string(risk.ENUMRISKLEVEL_LOW)),
						"type":  types.StringUnknown(),
					}

					objValue, d := types.ObjectValue(defaultResultTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"level": schema.StringAttribute{
						Description:         defaultResultLevelDescription.Description,
						MarkdownDescription: defaultResultLevelDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumRiskLevelEnumValues)...),
						},
					},

					"type": schema.StringAttribute{
						Description:         defaultResultTypeDescription.Description,
						MarkdownDescription: defaultResultTypeDescription.MarkdownDescription,
						Computed:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseNonNullStateForUnknown(),
						},
					},
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

			"evaluated_predictors": schema.SetAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of IDs for the predictors to evaluate in this policy set.  If omitted, if this property is null, all of the licensed predictors are used.").Description,
				Optional:    true,
				Computed:    true,

				ElementType: types.StringType,

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"policy_weights": schema.SingleNestedAttribute{
				Description:         policyWeightedAverageDescription.Description,
				MarkdownDescription: policyWeightedAverageDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"policy_threshold_medium": riskPolicyThresholdSchema(
						false,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold score values that define the medium risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int32{
							int32validatorinternal.IsLessThanPathValue(
								path.MatchRoot("policy_weights").AtName("policy_threshold_high").AtName("min_score"),
							),
						},
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						false,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold score values that define the high risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int32{
							int32validatorinternal.IsGreaterThanPathValue(
								path.MatchRoot("policy_weights").AtName("policy_threshold_medium").AtName("min_score"),
							),
						},
					),

					"predictors": schema.SetNestedAttribute{
						Description:         policyWeightedAveragePredictor.Description,
						MarkdownDescription: policyWeightedAveragePredictor.MarkdownDescription,

						Required: true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"compact_name": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the compact name of the predictor to apply to the risk policy.").Description,
									Required:    true,
								},

								"predictor_reference_value": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the level to evaluate.").Description,
									Computed:    true,
								},

								"weight": schema.Int32Attribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the weight to apply to the predictor when calculating the overall risk score.").Description,
									Required:    true,

									Validators: []validator.Int32{
										int32validator.AtLeast(weightMinimumDefault),
										int32validator.AtMost(weightMaximumDefault),
									},
								},
							},
						},

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("policy_weights"),
						path.MatchRelative().AtParent().AtName("policy_scores"),
					),
				},
			},

			"policy_scores": schema.SingleNestedAttribute{
				Description:         policyScoresDescription.Description,
				MarkdownDescription: policyScoresDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"policy_threshold_medium": riskPolicyThresholdSchema(
						true,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int32{
							int32validatorinternal.IsLessThanPathValue(
								path.MatchRoot("policy_scores").AtName("policy_threshold_high").AtName("min_score"),
							),
						},
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						true,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int32{
							int32validatorinternal.IsGreaterThanPathValue(
								path.MatchRoot("policy_scores").AtName("policy_threshold_medium").AtName("min_score"),
							),
						},
					),

					"predictors": schema.SetNestedAttribute{
						Description:         policyScoresPredictor.Description,
						MarkdownDescription: policyScoresPredictor.MarkdownDescription,

						Required: true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"compact_name": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the compact name of the predictor to apply to the risk policy.").Description,
									Required:    true,
								},

								"predictor_reference_value": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the level to evaluate.").Description,
									Computed:    true,
								},

								"score": schema.Int32Attribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the score to apply to the High risk / true outcome of the predictor, to apply to the overall risk calculation.").Description,
									Required:    true,

									Validators: []validator.Int32{
										int32validator.AtLeast(scoreMinimumDefault),
										int32validator.AtMost(scoreMaximumDefault),
									},
								},
							},
						},

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("policy_weights"),
						path.MatchRelative().AtParent().AtName("policy_scores"),
					),
				},
			},

			"overrides": schema.ListNestedAttribute{
				Description:         policyOverrideDescription.Description,
				MarkdownDescription: policyOverrideDescription.MarkdownDescription,

				Optional: true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the name of the overriding risk policy in the set.").Description,
							Optional:    true,
							Computed:    true,
						},

						"priority": schema.Int32Attribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that indicates the order in which the override is applied during risk policy evaluation.  The lower the value, the higher the priority.  The priority is determined by the order in which the overrides are defined in HCL.").Description,
							Computed:    true,
						},

						"result": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains the risk result that should be applied to the policy evaluation result when the override condition is met.").Description,
							Required:    true,

							Attributes: map[string]schema.Attribute{
								"value": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An administrator defined string value that is applied to the policy evaluation result when the override condition is met.").Description,
									Optional:    true,
								},

								"level": schema.StringAttribute{
									Description:         policyOverrideResultLevelDescription.Description,
									MarkdownDescription: policyOverrideResultLevelDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumRiskLevelEnumValues)...),
									},
								},

								"type": schema.StringAttribute{
									Description:         policyOverrideResultTypeDescription.Description,
									MarkdownDescription: policyOverrideResultTypeDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Default: stringdefault.StaticString(string(risk.ENUMRESULTTYPE_VALUE)),

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice([]risk.EnumResultType{risk.ENUMRESULTTYPE_VALUE})...),
									},

									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseNonNullStateForUnknown(),
									},
								},
							},
						},

						"condition": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains the conditions to evaluate that determine whether the override result will be applied to the risk policy evaluation.").Description,
							Required:    true,

							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description:         policyOverrideConditionTypeDescription.Description,
									MarkdownDescription: policyOverrideConditionTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(
											string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON),
											string(risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE),
										),
									},
								},

								// Value comparison
								"equals": schema.StringAttribute{
									Description:         policyOverrideConditionEqualsDescription.Description,
									MarkdownDescription: policyOverrideConditionEqualsDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"compact_name": schema.StringAttribute{
									Description:         policyOverrideConditionCompactNameDescription.Description,
									MarkdownDescription: policyOverrideConditionCompactNameDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"predictor_reference_value": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the value to evaluate.").Description,
									Computed:    true,
								},

								// IP range
								"ip_range": schema.SetAttribute{
									Description:         policyOverrideConditionIPRangeDescription.Description,
									MarkdownDescription: policyOverrideConditionIPRangeDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(verify.IPv4IPv6Regexp, "Values must be valid IPv4 or IPv6 CIDR format."),
										),
										setvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"predictor_reference_contains": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the collection to evaluate.").Description,
									Computed:    true,
								},
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtLeast(attrMinLength),
					listvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("mitigations"),
						path.MatchRelative().AtParent().AtName("fallback"),
						path.MatchRelative().AtParent().AtName("targets"),
					),
				},
			},

			"mitigations": schema.ListNestedAttribute{
				Description:         policyMitigationDescription.Description,
				MarkdownDescription: policyMitigationDescription.MarkdownDescription,

				Optional: true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the name of the mitigation policy entry. Computed from the condition's compact name by the provider.").Description,
							Computed:    true,
						},

						"priority": schema.Int32Attribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that indicates the order in which the mitigation entry is applied during risk policy evaluation. The lower the value, the higher the priority. The priority is determined by the order in which the entries are defined in HCL.").Description,
							Computed:    true,
						},

						"condition": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains the conditions to evaluate that determine whether the mitigation action will be applied to the risk policy evaluation.").Description,
							Required:    true,

							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description:         policyOverrideConditionTypeDescription.Description,
									MarkdownDescription: policyOverrideConditionTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(
											string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON),
											string(risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE),
										),
									},
								},

								// Value comparison
								"equals": schema.StringAttribute{
									Description:         policyOverrideConditionEqualsDescription.Description,
									MarkdownDescription: policyOverrideConditionEqualsDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"compact_name": schema.StringAttribute{
									Description:         policyOverrideConditionCompactNameDescription.Description,
									MarkdownDescription: policyOverrideConditionCompactNameDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"predictor_reference_value": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the value to evaluate.").Description,
									Computed:    true,
								},

								// IP range
								"ip_range": schema.SetAttribute{
									Description:         policyOverrideConditionIPRangeDescription.Description,
									MarkdownDescription: policyOverrideConditionIPRangeDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(verify.IPv4IPv6Regexp, "Values must be valid IPv4 or IPv6 CIDR format."),
										),
										setvalidatorinternal.IsRequiredIfMatchesPathValue(
											basetypes.NewStringValue(string(risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},

								"predictor_reference_contains": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the attribute reference of the collection to evaluate.").Description,
									Computed:    true,
								},
							},
						},

						"action": schema.StringAttribute{
							Description:         policyMitigationActionDescription.Description,
							MarkdownDescription: policyMitigationActionDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumMitigationActionEnumValues)...),
							},
						},

						"custom_action": schema.StringAttribute{
							Description:         policyMitigationCustomActionDescription.Description,
							MarkdownDescription: policyMitigationCustomActionDescription.MarkdownDescription,
							Optional:            true,

							Validators: []validator.String{
								stringvalidatorinternal.IsRequiredIfMatchesPathValue(
									basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_CUSTOM)),
									path.MatchRelative().AtParent().AtName("action"),
								),
							},
						},

						"mfa_authentication_policy_id": schema.StringAttribute{
							Description:         policyMitigationMfaAuthPolicyIdDescription.Description,
							MarkdownDescription: policyMitigationMfaAuthPolicyIdDescription.MarkdownDescription,
							Optional:            true,
							CustomType:          pingonetypes.ResourceIDType{},

							Validators: []validator.String{
								stringvalidatorinternal.IsRequiredIfMatchesPathValue(
									basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_MFA)),
									path.MatchRelative().AtParent().AtName("action"),
								),
							},
						},

						"mfa_registration_policy_id": schema.StringAttribute{
							Description:         policyMitigationMfaRegPolicyIdDescription.Description,
							MarkdownDescription: policyMitigationMfaRegPolicyIdDescription.MarkdownDescription,
							Optional:            true,
							CustomType:          pingonetypes.ResourceIDType{},
						},

						"verify_policy_id": schema.StringAttribute{
							Description:         policyMitigationVerifyPolicyIdDescription.Description,
							MarkdownDescription: policyMitigationVerifyPolicyIdDescription.MarkdownDescription,
							Optional:            true,
							CustomType:          pingonetypes.ResourceIDType{},

							Validators: []validator.String{
								stringvalidatorinternal.IsRequiredIfMatchesPathValue(
									basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_VERIFY)),
									path.MatchRelative().AtParent().AtName("action"),
								),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("overrides"),
					),
					listvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("fallback"),
					),
				},
			},

			"fallback": schema.SingleNestedAttribute{
				Description:         policyFallbackDescription.Description,
				MarkdownDescription: policyFallbackDescription.MarkdownDescription,

				Optional: true,

				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Description:         policyMitigationActionDescription.Description,
						MarkdownDescription: policyMitigationActionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumMitigationActionEnumValues)...),
						},
					},

					"custom_action": schema.StringAttribute{
						Description:         policyMitigationCustomActionDescription.Description,
						MarkdownDescription: policyMitigationCustomActionDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_CUSTOM)),
								path.MatchRelative().AtParent().AtName("action"),
							),
						},
					},

					"mfa_authentication_policy_id": schema.StringAttribute{
						Description:         policyMitigationMfaAuthPolicyIdDescription.Description,
						MarkdownDescription: policyMitigationMfaAuthPolicyIdDescription.MarkdownDescription,
						Optional:            true,
						CustomType:          pingonetypes.ResourceIDType{},

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_MFA)),
								path.MatchRelative().AtParent().AtName("action"),
							),
						},
					},

					"mfa_registration_policy_id": schema.StringAttribute{
						Description:         policyMitigationMfaRegPolicyIdDescription.Description,
						MarkdownDescription: policyMitigationMfaRegPolicyIdDescription.MarkdownDescription,
						Optional:            true,
						CustomType:          pingonetypes.ResourceIDType{},
					},

					"verify_policy_id": schema.StringAttribute{
						Description:         policyMitigationVerifyPolicyIdDescription.Description,
						MarkdownDescription: policyMitigationVerifyPolicyIdDescription.MarkdownDescription,
						Optional:            true,
						CustomType:          pingonetypes.ResourceIDType{},

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								basetypes.NewStringValue(string(risk.ENUMMITIGATIONACTION_VERIFY)),
								path.MatchRelative().AtParent().AtName("action"),
							),
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("overrides"),
					),
					objectvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("mitigations"),
					),
				},
			},

			"targets": schema.SingleNestedAttribute{
				Description:         policyTargetsDescription.Description,
				MarkdownDescription: policyTargetsDescription.MarkdownDescription,

				Optional: true,

				Attributes: map[string]schema.Attribute{
					"condition": schema.SingleNestedAttribute{
						Description:         policyTargetsConditionDescription.Description,
						MarkdownDescription: policyTargetsConditionDescription.MarkdownDescription,
						Required:            true,

						Attributes: map[string]schema.Attribute{
							"and": schema.ListNestedAttribute{
								Description:         policyTargetsConditionAndDescription.Description,
								MarkdownDescription: policyTargetsConditionAndDescription.MarkdownDescription,
								Required:            true,

								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description:         policyTargetsConditionAndTypeDescription.Description,
											MarkdownDescription: policyTargetsConditionAndTypeDescription.MarkdownDescription,
											Computed:            true,

											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseNonNullStateForUnknown(),
											},
										},

										"list": schema.ListAttribute{
											Description:         policyTargetsConditionAndListDescription.Description,
											MarkdownDescription: policyTargetsConditionAndListDescription.MarkdownDescription,
											Required:            true,
											ElementType:         types.StringType,
										},

										"contains": schema.StringAttribute{
											Description:         policyTargetsConditionAndContainsDescription.Description,
											MarkdownDescription: policyTargetsConditionAndContainsDescription.MarkdownDescription,
											Required:            true,
										},
									},
								},
							},
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("overrides"),
					),
				},
			},
		},
	}
}

func riskPolicyThresholdSchema(useScores bool, policyThresholdsDescription framework.SchemaAttributeDescription, validators []validator.Int32) schema.SingleNestedAttribute {

	validators = append(validators, int32validator.AtLeast(1))

	policyThresholdScoresMediumScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum score to use as the lower bound value of the policy threshold.",
	)

	policyThresholdScoresHighScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.",
	)

	if !useScores {
		maxAllowedValue := 100
		denominator := 10
		validators = append(validators, int32validator.AtMost(int32(maxAllowedValue)))
		validators = append(validators, int32validatorinternal.IsDivisibleBy(int32(denominator)))

		policyThresholdScoresMediumScoreDescription = policyThresholdScoresMediumScoreDescription.AppendMarkdownString(fmt.Sprintf("For weights policies, the score values should be 10x the desired risk value in the console. For example, a risk score of `5` in the console should be entered as `50`.  The provided score must be exactly divisible by 10.  Maximum value allowed is `%d`", maxAllowedValue))
	} else {
		maxAllowedValue := 1000
		validators = append(validators, int32validator.AtMost(int32(maxAllowedValue)))
		policyThresholdScoresMediumScoreDescription = policyThresholdScoresMediumScoreDescription.AppendMarkdownString(fmt.Sprintf("Maximum value allowed is `%d`", maxAllowedValue))
	}

	return schema.SingleNestedAttribute{
		Description:         policyThresholdsDescription.Description,
		MarkdownDescription: policyThresholdsDescription.MarkdownDescription,
		Required:            true,

		Attributes: map[string]schema.Attribute{
			"min_score": schema.Int32Attribute{
				Description:         policyThresholdScoresMediumScoreDescription.Description,
				MarkdownDescription: policyThresholdScoresMediumScoreDescription.MarkdownDescription,
				Required:            true,

				Validators: validators,
			},

			"max_score": schema.Int32Attribute{
				Description:         policyThresholdScoresHighScoreDescription.Description,
				MarkdownDescription: policyThresholdScoresHighScoreDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

// ModifyPlan
func (r *RiskPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state, config riskPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Evaluated Predictors
	setEvaluatedPredictorsToUnknown := false

	if !req.State.Raw.IsNull() && config.EvaluatedPredictors.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// If the policy is different, set the evaluated predictors as unknown
		if plan.PolicyScores.IsUnknown() || plan.PolicyWeights.IsUnknown() {
			setEvaluatedPredictorsToUnknown = true
		} else if (state.PolicyScores.IsNull() && !plan.PolicyScores.IsNull()) || (state.PolicyWeights.IsNull() && !plan.PolicyWeights.IsNull()) {
			setEvaluatedPredictorsToUnknown = true
		}
	}

	// Set the max threshold score
	var rootPath, referenceValueFmt string
	var maxScore int
	flattenedPolicyList := []attr.Value{}
	var predictorAttrType map[string]attr.Type

	if !plan.PolicyWeights.IsNull() && !plan.PolicyWeights.IsUnknown() {
		rootPath = "policy_weights"
		maxScore = 100
		referenceValueFmt = "${details.aggregatedWeights.%s}"
		predictorAttrType = policyWeightsPredictorTFObjectTypes

		var predictorsPlan, predictorsState []riskPolicyResourcePolicyWeightsPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsState)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !req.State.Raw.IsNull() && config.EvaluatedPredictors.IsNull() && len(predictorsState) != len(predictorsPlan) {
			setEvaluatedPredictorsToUnknown = true
		}

		for _, predictor := range predictorsPlan {
			predictorObj := map[string]attr.Value{
				"predictor_reference_value": framework.StringToTF(fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString())),
				"compact_name":              predictor.CompactName,
				"weight":                    predictor.Weight,
			}

			if config.EvaluatedPredictors.IsNull() && !setEvaluatedPredictorsToUnknown {
				found := false
				for _, statePredictor := range predictorsState {
					if statePredictor.CompactName.Equal(predictor.CompactName) {
						found = true
						break
					}
				}

				if !found {
					setEvaluatedPredictorsToUnknown = true
				}
			}

			flattenedObj, d := types.ObjectValue(policyWeightsPredictorTFObjectTypes, predictorObj)
			resp.Diagnostics.Append(d...)

			flattenedPolicyList = append(flattenedPolicyList, flattenedObj)
		}
	}

	if !plan.PolicyScores.IsNull() && !plan.PolicyScores.IsUnknown() {
		rootPath = "policy_scores"
		maxScore = 1000
		referenceValueFmt = "${details.%s.level}"
		predictorAttrType = policyScoresPredictorTFObjectTypes

		var predictorsPlan, predictorsState []riskPolicyResourcePolicyScoresPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsState)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !req.State.Raw.IsNull() && config.EvaluatedPredictors.IsNull() && len(predictorsState) != len(predictorsPlan) {
			setEvaluatedPredictorsToUnknown = true
		}

		for _, predictor := range predictorsPlan {
			predictorObj := map[string]attr.Value{
				"predictor_reference_value": framework.StringToTF(fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString())),
				"compact_name":              predictor.CompactName,
				"score":                     predictor.Score,
			}

			if config.EvaluatedPredictors.IsNull() && !setEvaluatedPredictorsToUnknown {
				found := false
				for _, statePredictor := range predictorsState {
					if statePredictor.CompactName.Equal(predictor.CompactName) {
						found = true
						break
					}
				}

				if !found {
					setEvaluatedPredictorsToUnknown = true
				}
			}

			flattenedObj, d := types.ObjectValue(policyScoresPredictorTFObjectTypes, predictorObj)
			resp.Diagnostics.Append(d...)

			flattenedPolicyList = append(flattenedPolicyList, flattenedObj)
		}
	}

	// Set the min-max threshold scores/weights
	var policyThresholdHighMinValue *int32
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("min_score"), &policyThresholdHighMinValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_medium").AtName("max_score"), types.Int32Value(*policyThresholdHighMinValue))
	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("max_score"), types.Int32Value(int32(maxScore)))

	// Set the predictors
	plannedPredictors, d := types.SetValue(types.ObjectType{AttrTypes: predictorAttrType}, flattenedPolicyList)
	resp.Diagnostics.Append(d...)
	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("predictors"), plannedPredictors)

	// Overrides
	flattenedOverrideList := []attr.Value{}

	if !plan.Overrides.IsNull() && !plan.Overrides.IsUnknown() {
		var overridesPlan []riskPolicyResourcePolicyOverrideModel
		resp.Diagnostics.Append(plan.Overrides.ElementsAs(ctx, &overridesPlan, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		referenceValueFmt = "${details.%s.level}"
		priorityCount := 0

		for _, overridePlan := range overridesPlan {

			priorityCount++

			// The Condition
			var conditionPlan riskPolicyResourcePolicyOverrideConditionModel
			resp.Diagnostics.Append(overridePlan.Condition.As(ctx, &conditionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if resp.Diagnostics.HasError() {
				return
			}

			var predictorReferenceValue attr.Value
			var predictorReferenceContains attr.Value

			var overrideName string
			setOverrideName := true
			if !overridePlan.Name.IsNull() && !overridePlan.Name.IsUnknown() {
				overrideName = overridePlan.Name.ValueString()
				setOverrideName = false
			}

			if !conditionPlan.CompactName.IsNull() && !conditionPlan.CompactName.IsUnknown() {
				predictorReferenceValue = framework.StringToTF(fmt.Sprintf(referenceValueFmt, conditionPlan.CompactName.ValueString()))
				predictorReferenceContains = types.StringNull()

				if setOverrideName {
					overrideName = conditionPlan.CompactName.ValueString()
				}
			}

			if !conditionPlan.IPRange.IsNull() && !conditionPlan.IPRange.IsUnknown() {
				predictorReferenceContains = framework.StringToTF("${transaction.ip}")
				predictorReferenceValue = types.StringNull()

				if setOverrideName {
					overrideName = "WHITELIST"
				}
			}

			conditionMap := map[string]attr.Value{
				"type":                         conditionPlan.Type,
				"equals":                       conditionPlan.Equals,
				"compact_name":                 conditionPlan.CompactName,
				"predictor_reference_value":    predictorReferenceValue,
				"ip_range":                     conditionPlan.IPRange,
				"predictor_reference_contains": predictorReferenceContains,
			}

			conditionObj, d := types.ObjectValue(overridesConditionTFObjectTypes, conditionMap)
			resp.Diagnostics.Append(d...)

			overrideMap := map[string]attr.Value{
				"name":      types.StringValue(overrideName),
				"priority":  types.Int32Value(int32(priorityCount)),
				"result":    overridePlan.Result,
				"condition": conditionObj,
			}

			overrideObj, d := types.ObjectValue(overridesTFObjectTypes, overrideMap)
			resp.Diagnostics.Append(d...)

			flattenedOverrideList = append(flattenedOverrideList, overrideObj)
		}

		plannedOverrides, d := types.ListValue(types.ObjectType{AttrTypes: overridesTFObjectTypes}, flattenedOverrideList)
		resp.Diagnostics.Append(d...)
		resp.Plan.SetAttribute(ctx, path.Root("overrides"), plannedOverrides)
	}

	// Mitigations
	flattenedMitigationList := []attr.Value{}

	if !plan.Mitigations.IsNull() && !plan.Mitigations.IsUnknown() {
		var mitigationsPlan []riskPolicyResourcePolicyMitigationModel
		resp.Diagnostics.Append(plan.Mitigations.ElementsAs(ctx, &mitigationsPlan, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		referenceValueFmt = "${details.%s.level}"
		priorityCount := 0

		for _, mitigationPlan := range mitigationsPlan {

			priorityCount++

			// The Condition
			var conditionPlan riskPolicyResourcePolicyOverrideConditionModel
			resp.Diagnostics.Append(mitigationPlan.Condition.As(ctx, &conditionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if resp.Diagnostics.HasError() {
				return
			}

			var predictorReferenceValue attr.Value
			var predictorReferenceContains attr.Value

			var mitigationName string

			if !conditionPlan.CompactName.IsNull() && !conditionPlan.CompactName.IsUnknown() {
				predictorReferenceValue = framework.StringToTF(fmt.Sprintf(referenceValueFmt, conditionPlan.CompactName.ValueString()))
				predictorReferenceContains = types.StringNull()
				mitigationName = conditionPlan.CompactName.ValueString()
			}

			if !conditionPlan.IPRange.IsNull() && !conditionPlan.IPRange.IsUnknown() {
				predictorReferenceContains = framework.StringToTF("${transaction.ip}")
				predictorReferenceValue = types.StringNull()
				mitigationName = "WHITELIST"
			}

			conditionMap := map[string]attr.Value{
				"type":                         conditionPlan.Type,
				"equals":                       conditionPlan.Equals,
				"compact_name":                 conditionPlan.CompactName,
				"predictor_reference_value":    predictorReferenceValue,
				"ip_range":                     conditionPlan.IPRange,
				"predictor_reference_contains": predictorReferenceContains,
			}

			conditionObj, d := types.ObjectValue(overridesConditionTFObjectTypes, conditionMap)
			resp.Diagnostics.Append(d...)

			mitigationMap := map[string]attr.Value{
				"name":                         types.StringValue(mitigationName),
				"priority":                     types.Int32Value(int32(priorityCount)),
				"condition":                    conditionObj,
				"action":                       mitigationPlan.Action,
				"custom_action":                mitigationPlan.CustomAction,
				"mfa_authentication_policy_id": mitigationPlan.MfaAuthenticationPolicyId,
				"mfa_registration_policy_id":   mitigationPlan.MfaRegistrationPolicyId,
				"verify_policy_id":             mitigationPlan.VerifyPolicyId,
			}

			mitigationObj, d := types.ObjectValue(mitigationsTFObjectTypes, mitigationMap)
			resp.Diagnostics.Append(d...)

			flattenedMitigationList = append(flattenedMitigationList, mitigationObj)
		}

		plannedMitigations, d := types.ListValue(types.ObjectType{AttrTypes: mitigationsTFObjectTypes}, flattenedMitigationList)
		resp.Diagnostics.Append(d...)
		resp.Plan.SetAttribute(ctx, path.Root("mitigations"), plannedMitigations)
	}

	if setEvaluatedPredictorsToUnknown {
		resp.Plan.SetAttribute(ctx, path.Root("evaluated_predictors"), types.SetUnknown(types.StringType))
	}
}

func (r *RiskPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RiskPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state riskPolicyResourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
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
	riskPolicy, d := plan.expand(ctx, r.Client.RiskAPIClient, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var createResponse *risk.RiskPolicySet
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.CreateRiskPolicySet(ctx, plan.EnvironmentId.ValueString()).RiskPolicySet(*riskPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateRiskPolicySet",
		riskPolicyCreateUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
		&createResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We have to read it back because the API does not return the full state object on create
	var response *risk.RiskPolicySet
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), createResponse.GetId()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPolicySet",
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

func (r *RiskPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *riskPolicyResourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
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
	var response *risk.RiskPolicySet
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPolicySet",
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

func (r *RiskPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state riskPolicyResourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
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
	riskPolicy, d := plan.expand(ctx, r.Client.RiskAPIClient, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *risk.RiskPolicySet
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.UpdateRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).RiskPolicySet(*riskPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateRiskPolicySet",
		riskPolicyCreateUpdateCustomErrorHandler,
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

func (r *RiskPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *riskPolicyResourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
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
			fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.DeleteRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteRiskPolicySet",
		riskPolicyDeleteCustomError,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{
			"200",
			"403",
		},
		Target: []string{
			"404",
			"defaulted",
		},
		Refresh: func() (interface{}, string, error) {
			base := 10

			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			resp, r, err := legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil {
				if r.StatusCode == 404 {
					return risk.RiskPolicySet{}, strconv.FormatInt(int64(r.StatusCode), base), nil
				}
				return nil, strconv.FormatInt(int64(r.StatusCode), base), err
			}

			if defaultConfig, ok := fO.GetDefaultOk(); ok && *defaultConfig {
				return resp, "defaulted", nil
			}

			return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   20 * time.Minute,
		Delay:                     5 * time.Second,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	_, err := deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Risk Policy Delete Timeout",
			fmt.Sprintf("Error waiting for risk policy (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}
}

var riskPolicyDeleteCustomError = func(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Undeletable default risk policy
	if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
		if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("remove default policy", v[0].GetMessage()); match {

				diags.AddWarning("Cannot delete the default risk policy", "Due to API restrictions, the provider cannot delete the default risk policy for an environment.  The policy has been removed from Terraform state but has been left in place in the PingOne service.")

				return diags
			}
		}
	}

	diags.Append(legacysdk.CustomErrorResourceNotFoundWarning(r, p1Error)...)
	return diags
}

func (r *RiskPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "risk_policy_id",
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

func riskPolicyCreateUpdateCustomErrorHandler(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Invalid composition
	if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if target, ok := details[0].GetTargetOk(); ok && *target == "composition.condition" {
			diags.AddError(
				"Invalid \"composition.condition\" policy JSON.",
				"Please check the \"composition.condition\" policy JSON structure and contents and try again.",
			)

			return diags
		}
	}

	return diags
}

func (p *riskPolicyResourceModel) expand(ctx context.Context, apiClient *risk.APIClient, managementApiClient *management.APIClient) (*risk.RiskPolicySet, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPolicySet(p.Name.ValueString())

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	} else {
		data.SetDefault(false)
	}

	if !p.DefaultResult.IsNull() && !p.DefaultResult.IsUnknown() {
		var plan riskPolicyResourceDefaultResultModel
		diags.Append(p.DefaultResult.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		if !plan.Level.IsNull() && !plan.Level.IsUnknown() {
			data.SetDefaultResult(*risk.NewRiskPolicySetDefaultResult(risk.EnumRiskPolicyResultLevel(plan.Level.ValueString())))
		}
	}

	highPolicyCondition := risk.NewRiskPolicyCondition()
	mediumPolicyCondition := risk.NewRiskPolicyCondition()

	var plan riskPolicyResourcePolicyModel
	var d diag.Diagnostics

	var useScores bool

	if !p.PolicyWeights.IsNull() && !p.PolicyWeights.IsUnknown() {
		highPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS)
		mediumPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS)
		useScores = false

		diags.Append(p.PolicyWeights.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
	}

	if !p.PolicyScores.IsNull() && !p.PolicyScores.IsUnknown() {
		highPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES)
		mediumPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES)
		useScores = true

		diags.Append(p.PolicyScores.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
	}

	if diags.HasError() {
		return nil, diags
	}

	var predictorCompactNames []string
	highPolicyCondition, mediumPolicyCondition, predictorCompactNames, d = plan.expand(ctx, useScores, highPolicyCondition, mediumPolicyCondition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	riskPolicies := make([]risk.RiskPolicy, 0)

	// Overrides
	if !p.Overrides.IsNull() && !p.Overrides.IsUnknown() {
		var overridesPlan []riskPolicyResourcePolicyOverrideModel
		diags.Append(p.Overrides.ElementsAs(ctx, &overridesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, overridePlan := range overridesPlan {

			// The Condition
			var conditionPlan riskPolicyResourcePolicyOverrideConditionModel
			diags.Append(overridePlan.Condition.As(ctx, &conditionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			condition := risk.NewRiskPolicyCondition()

			condition.SetType(risk.EnumRiskPolicyConditionType(conditionPlan.Type.ValueString()))

			if !conditionPlan.Equals.IsNull() && !conditionPlan.Equals.IsUnknown() {
				v := conditionPlan.Equals.ValueString()
				condition.SetEquals(risk.StringAsRiskPolicyConditionEquals(&v))
			}

			if !conditionPlan.PredictorReferenceValue.IsNull() && !conditionPlan.PredictorReferenceValue.IsUnknown() {
				condition.SetValue(conditionPlan.PredictorReferenceValue.ValueString())
				if !slices.Contains(predictorCompactNames, conditionPlan.CompactName.ValueString()) {
					predictorCompactNames = append(predictorCompactNames, conditionPlan.CompactName.ValueString())
				}
			}

			if !conditionPlan.PredictorReferenceContains.IsNull() && !conditionPlan.PredictorReferenceContains.IsUnknown() {
				condition.SetContains(conditionPlan.PredictorReferenceContains.ValueString())
			}

			if !conditionPlan.IPRange.IsNull() && !conditionPlan.IPRange.IsUnknown() {
				var conditionIPRangePlan []types.String
				diags.Append(conditionPlan.IPRange.ElementsAs(ctx, &conditionIPRangePlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				conditionIPRange, d := framework.TFTypeStringSliceToStringSlice(conditionIPRangePlan, path.Root("overrides").AtName("condition").AtName("ip_range"))
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				condition.SetIpRange(conditionIPRange)
			}

			// The Result
			var resultPlan riskPolicyResourcePolicyOverrideResultModel
			diags.Append(overridePlan.Result.As(ctx, &resultPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			result := risk.NewRiskPolicyResult()
			result.SetLevel(risk.EnumRiskLevel(resultPlan.Level.ValueString()))

			if !resultPlan.Value.IsNull() && !resultPlan.Value.IsUnknown() {
				result.SetValue(resultPlan.Value.ValueString())
			}

			if !resultPlan.Type.IsNull() && !resultPlan.Type.IsUnknown() {
				result.SetType(risk.EnumResultType(resultPlan.Type.ValueString()))
			}

			op := *risk.NewRiskPolicy(
				overridePlan.Name.ValueString(),
				*result,
			)
			op.SetCondition(*condition)

			op.SetPriority(overridePlan.Priority.ValueInt32())

			riskPolicies = append(riskPolicies, op)
		}
	}

	// Mitigations
	if !p.Mitigations.IsNull() && !p.Mitigations.IsUnknown() {
		var mitigationsPlan []riskPolicyResourcePolicyMitigationModel
		diags.Append(p.Mitigations.ElementsAs(ctx, &mitigationsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, mitigationPlan := range mitigationsPlan {

			// The Condition
			var conditionPlan riskPolicyResourcePolicyOverrideConditionModel
			diags.Append(mitigationPlan.Condition.As(ctx, &conditionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			condition := risk.NewRiskPolicyCondition()

			condition.SetType(risk.EnumRiskPolicyConditionType(conditionPlan.Type.ValueString()))

			if !conditionPlan.Equals.IsNull() && !conditionPlan.Equals.IsUnknown() {
				v := conditionPlan.Equals.ValueString()
				condition.SetEquals(risk.StringAsRiskPolicyConditionEquals(&v))
			}

			if !conditionPlan.PredictorReferenceValue.IsNull() && !conditionPlan.PredictorReferenceValue.IsUnknown() {
				condition.SetValue(conditionPlan.PredictorReferenceValue.ValueString())
				if !slices.Contains(predictorCompactNames, conditionPlan.CompactName.ValueString()) {
					predictorCompactNames = append(predictorCompactNames, conditionPlan.CompactName.ValueString())
				}
			}

			if !conditionPlan.PredictorReferenceContains.IsNull() && !conditionPlan.PredictorReferenceContains.IsUnknown() {
				condition.SetContains(conditionPlan.PredictorReferenceContains.ValueString())
			}

			if !conditionPlan.IPRange.IsNull() && !conditionPlan.IPRange.IsUnknown() {
				var conditionIPRangePlan []types.String
				diags.Append(conditionPlan.IPRange.ElementsAs(ctx, &conditionIPRangePlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				conditionIPRange, d := framework.TFTypeStringSliceToStringSlice(conditionIPRangePlan, path.Root("mitigations").AtName("condition").AtName("ip_range"))
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				condition.SetIpRange(conditionIPRange)
			}

			// The Result
			mitigationInner := *risk.NewRiskPolicyResultMitigationsInner(risk.EnumMitigationAction(mitigationPlan.Action.ValueString()))

			if !mitigationPlan.CustomAction.IsNull() && !mitigationPlan.CustomAction.IsUnknown() {
				mitigationInner.SetCustomAction(mitigationPlan.CustomAction.ValueString())
			}

			if !mitigationPlan.MfaAuthenticationPolicyId.IsNull() && !mitigationPlan.MfaAuthenticationPolicyId.IsUnknown() {
				mitigationInner.SetMfaAuthenticationPolicyId(mitigationPlan.MfaAuthenticationPolicyId.ValueString())
			}

			if !mitigationPlan.MfaRegistrationPolicyId.IsNull() && !mitigationPlan.MfaRegistrationPolicyId.IsUnknown() {
				mitigationInner.SetMfaRegistrationPolicyId(mitigationPlan.MfaRegistrationPolicyId.ValueString())
			}

			if !mitigationPlan.VerifyPolicyId.IsNull() && !mitigationPlan.VerifyPolicyId.IsUnknown() {
				mitigationInner.SetVerifyPolicyId(mitigationPlan.VerifyPolicyId.ValueString())
			}

			result := risk.NewRiskPolicyResult()
			result.SetType(risk.ENUMRESULTTYPE_MITIGATION)
			result.SetMitigations([]risk.RiskPolicyResultMitigationsInner{mitigationInner})

			op := *risk.NewRiskPolicy(mitigationPlan.Name.ValueString(), *result)
			op.SetCondition(*condition)
			op.SetPriority(mitigationPlan.Priority.ValueInt32())

			riskPolicies = append(riskPolicies, op)
		}
	}

	// Fallback (MITIGATION_FALLBACK)
	if !p.Fallback.IsNull() && !p.Fallback.IsUnknown() {
		var fallbackPlan riskPolicyResourceMitigationFallbackModel
		diags.Append(p.Fallback.As(ctx, &fallbackPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		fallbackInner := *risk.NewRiskPolicyResultMitigationsInner(risk.EnumMitigationAction(fallbackPlan.Action.ValueString()))

		if !fallbackPlan.CustomAction.IsNull() && !fallbackPlan.CustomAction.IsUnknown() {
			fallbackInner.SetCustomAction(fallbackPlan.CustomAction.ValueString())
		}

		if !fallbackPlan.MfaAuthenticationPolicyId.IsNull() && !fallbackPlan.MfaAuthenticationPolicyId.IsUnknown() {
			fallbackInner.SetMfaAuthenticationPolicyId(fallbackPlan.MfaAuthenticationPolicyId.ValueString())
		}

		if !fallbackPlan.MfaRegistrationPolicyId.IsNull() && !fallbackPlan.MfaRegistrationPolicyId.IsUnknown() {
			fallbackInner.SetMfaRegistrationPolicyId(fallbackPlan.MfaRegistrationPolicyId.ValueString())
		}

		if !fallbackPlan.VerifyPolicyId.IsNull() && !fallbackPlan.VerifyPolicyId.IsUnknown() {
			fallbackInner.SetVerifyPolicyId(fallbackPlan.VerifyPolicyId.ValueString())
		}

		fallbackResult := risk.NewRiskPolicyResult()
		fallbackResult.SetType(risk.ENUMRESULTTYPE_MITIGATION_FALLBACK)
		fallbackResult.SetMitigations([]risk.RiskPolicyResultMitigationsInner{fallbackInner})

		fb := *risk.NewRiskPolicy("FALLBACK", *fallbackResult)

		riskPolicies = append(riskPolicies, fb)
	}

	// Medium Weighted Policy
	mwpResult := risk.NewRiskPolicyResult()
	mwpResult.SetLevel(risk.ENUMRISKLEVEL_MEDIUM)
	mwp := *risk.NewRiskPolicy("MEDIUM_WEIGHTED_POLICY", *mwpResult)
	mwp.SetCondition(*mediumPolicyCondition)
	mwp.SetPriority(int32(len(riskPolicies)) + 1)
	riskPolicies = append(riskPolicies, mwp)

	// High Weighted Policy
	hwpResult := risk.NewRiskPolicyResult()
	hwpResult.SetLevel(risk.ENUMRISKLEVEL_HIGH)
	hwp := *risk.NewRiskPolicy("HIGH_WEIGHTED_POLICY", *hwpResult)
	hwp.SetCondition(*highPolicyCondition)
	hwp.SetPriority(int32(len(riskPolicies)) + 1)
	riskPolicies = append(riskPolicies, hwp)

	// Set the risk policies
	data.SetRiskPolicies(riskPolicies)

	riskPolicyPredictorsIDs, d := riskPredictorFetchIDsFromCompactNames(ctx, apiClient, managementApiClient, p.EnvironmentId.ValueString(), predictorCompactNames)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	evaluatedPredictors := make([]risk.RiskPolicySetEvaluatedPredictorsInner, 0)

	if !p.EvaluatedPredictors.IsNull() && !p.EvaluatedPredictors.IsUnknown() {
		var plan []types.String
		diags.Append(p.EvaluatedPredictors.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		evaluatedPredictorsStr, d := framework.TFTypeStringSliceToStringSlice(plan, path.Root("evaluated_predictors"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		configuredEvaluatedPredictorIDs := make(map[string]bool)

		for _, predictorID := range evaluatedPredictorsStr {
			evaluatedPredictors = append(evaluatedPredictors, *risk.NewRiskPolicySetEvaluatedPredictorsInner(predictorID))
			configuredEvaluatedPredictorIDs[predictorID] = true
		}

		for _, riskPolicyPredictorID := range riskPolicyPredictorsIDs {
			if !configuredEvaluatedPredictorIDs[riskPolicyPredictorID] {
				diags.AddError(
					"A predictor in the policy set is not listed in \"evaluated_predictors\".",
					"When \"evaluated_predictors\" is defined, the IDs for predictors in the policy set must be listed in \"evaluated_predictors\".",
				)
			}
		}

		if diags.HasError() {
			return nil, diags
		}

	} else {
		for _, predictorID := range riskPolicyPredictorsIDs {
			evaluatedPredictors = append(evaluatedPredictors, *risk.NewRiskPolicySetEvaluatedPredictorsInner(predictorID))
		}
	}
	data.SetEvaluatedPredictors(evaluatedPredictors)

	// Targets
	if !p.Targets.IsNull() && !p.Targets.IsUnknown() {
		var targetsPlan riskPolicyResourceTargetsModel
		diags.Append(p.Targets.As(ctx, &targetsPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		if !targetsPlan.Condition.IsNull() && !targetsPlan.Condition.IsUnknown() {
			var conditionPlan riskPolicyResourceTargetsConditionModel
			diags.Append(targetsPlan.Condition.As(ctx, &conditionPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			sdkCondition := risk.NewRiskPolicySetTargetsCondition()

			if !conditionPlan.And.IsNull() && !conditionPlan.And.IsUnknown() {
				var andPlan []riskPolicyResourceTargetsConditionAndModel
				diags.Append(conditionPlan.And.ElementsAs(ctx, &andPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				andInners := make([]risk.RiskPolicySetTargetsConditionAndInner, 0, len(andPlan))
				for _, entry := range andPlan {
					inner := risk.NewRiskPolicySetTargetsConditionAndInner()

					if !entry.List.IsNull() && !entry.List.IsUnknown() {
						var listPlan []types.String
						diags.Append(entry.List.ElementsAs(ctx, &listPlan, false)...)
						if diags.HasError() {
							return nil, diags
						}
						listStrings, d := framework.TFTypeStringSliceToStringSlice(listPlan, path.Root("targets").AtName("condition").AtName("and").AtName("list"))
						diags.Append(d...)
						if diags.HasError() {
							return nil, diags
						}
						inner.SetList(listStrings)
					}

					if !entry.Contains.IsNull() && !entry.Contains.IsUnknown() {
						inner.SetContains(entry.Contains.ValueString())
					}

					andInners = append(andInners, *inner)
				}

				sdkCondition.SetAnd(andInners)
			}

			sdkTargets := risk.NewRiskPolicySetTargets()
			sdkTargets.SetCondition(*sdkCondition)
			data.SetTargets(*sdkTargets)
		}
	}

	return data, diags
}

func (p *riskPolicyResourcePolicyModel) expand(ctx context.Context, useScores bool, highPolicyCondition, mediumPolicyCondition *risk.RiskPolicyCondition) (*risk.RiskPolicyCondition, *risk.RiskPolicyCondition, []string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !p.PolicyThresholdMedium.IsNull() && !p.PolicyThresholdMedium.IsUnknown() {
		var plan riskPolicyResourcePolicyThresholdScoreBetweenModel
		diags.Append(p.PolicyThresholdMedium.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		mediumPolicyCondition.SetBetween(
			*risk.NewRiskPolicyConditionBetween(
				plan.MinScore.ValueInt32(),
				plan.MaxScore.ValueInt32(),
			),
		)
	}

	predictorCompactNames := make([]string, 0)

	if !p.PolicyThresholdHigh.IsNull() && !p.PolicyThresholdHigh.IsUnknown() {
		var plan riskPolicyResourcePolicyThresholdScoreBetweenModel
		diags.Append(p.PolicyThresholdHigh.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		highPolicyCondition.SetBetween(
			*risk.NewRiskPolicyConditionBetween(
				plan.MinScore.ValueInt32(),
				plan.MaxScore.ValueInt32(),
			),
		)
	}

	if !p.Predictors.IsNull() && !p.Predictors.IsUnknown() {
		if useScores {
			aggregatedScores := make([]risk.RiskPolicyConditionAggregatedScoresInner, 0)

			var predictorsPlan []riskPolicyResourcePolicyScoresPredictorModel
			diags.Append(p.Predictors.ElementsAs(ctx, &predictorsPlan, false)...)
			if diags.HasError() {
				return nil, nil, nil, diags
			}

			for _, predictor := range predictorsPlan {
				aggregatedScores = append(
					aggregatedScores,
					*risk.NewRiskPolicyConditionAggregatedScoresInner(
						predictor.PredictorReferenceValue.ValueString(),
						predictor.Score.ValueInt32(),
					),
				)

				predictorCompactNames = append(predictorCompactNames, predictor.CompactName.ValueString())
			}

			mediumPolicyCondition.SetAggregatedScores(aggregatedScores)
			highPolicyCondition.SetAggregatedScores(aggregatedScores)

		} else {
			aggregatedWeights := make([]risk.RiskPolicyConditionAggregatedWeightsInner, 0)

			var predictorsPlan []riskPolicyResourcePolicyWeightsPredictorModel
			diags.Append(p.Predictors.ElementsAs(ctx, &predictorsPlan, false)...)
			if diags.HasError() {
				return nil, nil, nil, diags
			}

			for _, predictor := range predictorsPlan {
				aggregatedWeights = append(
					aggregatedWeights,
					*risk.NewRiskPolicyConditionAggregatedWeightsInner(
						predictor.PredictorReferenceValue.ValueString(),
						predictor.Weight.ValueInt32(),
					),
				)

				predictorCompactNames = append(predictorCompactNames, predictor.CompactName.ValueString())
			}

			mediumPolicyCondition.SetAggregatedWeights(aggregatedWeights)
			highPolicyCondition.SetAggregatedWeights(aggregatedWeights)
		}
	}

	return highPolicyCondition, mediumPolicyCondition, predictorCompactNames, diags
}

func (p *riskPolicyResourceModel) toState(apiObject *risk.RiskPolicySet) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	// Default block
	p.DefaultResult = types.ObjectNull(defaultResultTFObjectTypes)
	if v, ok := apiObject.GetDefaultResultOk(); ok {
		o := map[string]attr.Value{
			"type":  framework.EnumOkToTF(v.GetTypeOk()),
			"level": framework.EnumOkToTF(v.GetLevelOk()),
		}

		objValue, d := types.ObjectValue(defaultResultTFObjectTypes, o)
		diags.Append(d...)

		p.DefaultResult = objValue
	}

	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())

	p.EvaluatedPredictors = types.SetNull(types.StringType)
	if v, ok := apiObject.GetEvaluatedPredictorsOk(); ok {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, framework.StringOkToTF(item.GetIdOk()))
		}

		var d diag.Diagnostics
		p.EvaluatedPredictors, d = types.SetValue(types.StringType, list)
		diags.Append(d...)
	}

	var d diag.Diagnostics

	r, ok := apiObject.GetRiskPoliciesOk()

	p.PolicyWeights, p.PolicyScores, p.Overrides, p.Mitigations, p.Fallback, d = p.toStatePolicy(r, ok)
	diags.Append(d...)

	// Targets
	p.Targets = types.ObjectNull(targetsTFObjectTypes)
	if targets, ok := apiObject.GetTargetsOk(); ok {
		if condition, ok := targets.GetConditionOk(); ok {
			andList := make([]attr.Value, 0)

			for _, andInner := range condition.GetAnd() {
				listValues := make([]attr.Value, 0)
				for _, s := range andInner.GetList() {
					listValues = append(listValues, types.StringValue(s))
				}

				listVal, ld := types.ListValue(types.StringType, listValues)
				diags.Append(ld...)

				andMap := map[string]attr.Value{
					"type":     framework.EnumOkToTF(andInner.GetTypeOk()),
					"list":     listVal,
					"contains": framework.StringOkToTF(andInner.GetContainsOk()),
				}

				andObj, od := types.ObjectValue(targetsConditionAndTFObjectTypes, andMap)
				diags.Append(od...)

				andList = append(andList, andObj)
			}

			andListVal, ld := types.ListValue(types.ObjectType{AttrTypes: targetsConditionAndTFObjectTypes}, andList)
			diags.Append(ld...)

			conditionMap := map[string]attr.Value{
				"and": andListVal,
			}

			conditionObj, cd := types.ObjectValue(targetsConditionTFObjectTypes, conditionMap)
			diags.Append(cd...)

			targetsMap := map[string]attr.Value{
				"condition": conditionObj,
			}

			targetsObj, td := types.ObjectValue(targetsTFObjectTypes, targetsMap)
			diags.Append(td...)

			p.Targets = targetsObj
		}
	}

	return diags
}

func (p *riskPolicyResourceModel) toStatePolicy(riskPolicies []risk.RiskPolicy, ok bool) (basetypes.ObjectValue, basetypes.ObjectValue, basetypes.ListValue, basetypes.ListValue, basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objPolicyWeightsValue := types.ObjectUnknown(policyWeightsTFObjectTypes)
	objPolicyScoresValue := types.ObjectUnknown(policyScoresTFObjectTypes)
	objOverridesValue := types.ListUnknown(types.ObjectType{AttrTypes: overridesTFObjectTypes})
	objMitigationsValue := types.ListNull(types.ObjectType{AttrTypes: mitigationsTFObjectTypes})
	objFallbackValue := types.ObjectNull(mitigationsFallbackTFObjectTypes)

	useScores := false
	useWeights := false

	if !ok || riskPolicies == nil || len(riskPolicies) < 1 {
		return objPolicyWeightsValue, objPolicyScoresValue, objOverridesValue, objMitigationsValue, objFallbackValue, diags
	}

	highMediumPolicy := map[string]attr.Value{}
	overrides := []attr.Value{}
	mitigations := []attr.Value{}

	setOverride := false
	setMitigation := false

	for _, policy := range riskPolicies {
		// Check the result type first to distinguish mitigation entries from overrides
		// (both may have VALUE_COMPARISON or IP_RANGE conditions).
		if policyResult, ok := policy.GetResultOk(); ok {
			if resultType, ok := policyResult.GetTypeOk(); ok {

				if *resultType == risk.ENUMRESULTTYPE_MITIGATION_FALLBACK {
					// Fallback entry: no condition, carries mitigations array with action.
					fallbackInners, ok := policyResult.GetMitigationsOk()
					if ok && len(fallbackInners) > 0 {
						inner := fallbackInners[0]
						fallbackMap := map[string]attr.Value{
							"action":                       framework.EnumOkToTF(inner.GetActionOk()),
							"custom_action":                framework.StringOkToTF(inner.GetCustomActionOk()),
							"mfa_authentication_policy_id": mitigationPolicyIDToTF(inner.GetMfaAuthenticationPolicyIdOk()),
							"mfa_registration_policy_id":   mitigationPolicyIDToTF(inner.GetMfaRegistrationPolicyIdOk()),
							"verify_policy_id":             mitigationPolicyIDToTF(inner.GetVerifyPolicyIdOk()),
						}
						var d diag.Diagnostics
						objFallbackValue, d = types.ObjectValue(mitigationsFallbackTFObjectTypes, fallbackMap)
						diags.Append(d...)
					}
					continue
				}

				if *resultType == risk.ENUMRESULTTYPE_MITIGATION {
					// Mitigation entry: has a condition and a mitigations array.
					setMitigation = true

					var conditionObj basetypes.ObjectValue
					if condition, ok := policy.GetConditionOk(); ok {
						var equalsString basetypes.StringValue
						if s := condition.GetEquals().String; s != nil {
							equalsString = framework.StringToTF(*s)
						} else {
							equalsString = types.StringNull()
						}

						conditionMap := map[string]attr.Value{
							"type":                         framework.EnumOkToTF(condition.GetTypeOk()),
							"equals":                       equalsString,
							"compact_name":                 riskPolicyOverrideCompactNameFromReferenceOk(condition.GetValueOk()),
							"predictor_reference_value":    framework.StringOkToTF(condition.GetValueOk()),
							"ip_range":                     framework.StringSetOkToTF(condition.GetIpRangeOk()),
							"predictor_reference_contains": framework.StringOkToTF(condition.GetContainsOk()),
						}

						var d diag.Diagnostics
						conditionObj, d = types.ObjectValue(overridesConditionTFObjectTypes, conditionMap)
						diags.Append(d...)
					} else {
						conditionObj = types.ObjectNull(overridesConditionTFObjectTypes)
					}

					mitigationInners, ok := policyResult.GetMitigationsOk()

					var actionVal basetypes.StringValue
					var customActionVal basetypes.StringValue
					var mfaAuthPolicyIDVal pingonetypes.ResourceIDValue
					var mfaRegPolicyIDVal pingonetypes.ResourceIDValue
					var verifyPolicyIDVal pingonetypes.ResourceIDValue

					if ok && len(mitigationInners) > 0 {
						inner := mitigationInners[0]
						actionVal = framework.EnumOkToTF(inner.GetActionOk())
						customActionVal = framework.StringOkToTF(inner.GetCustomActionOk())
						mfaAuthPolicyIDVal = mitigationPolicyIDToTF(inner.GetMfaAuthenticationPolicyIdOk())
						mfaRegPolicyIDVal = mitigationPolicyIDToTF(inner.GetMfaRegistrationPolicyIdOk())
						verifyPolicyIDVal = mitigationPolicyIDToTF(inner.GetVerifyPolicyIdOk())
					} else {
						actionVal = types.StringNull()
						customActionVal = types.StringNull()
						mfaAuthPolicyIDVal = pingonetypes.NewResourceIDNull()
						mfaRegPolicyIDVal = pingonetypes.NewResourceIDNull()
						verifyPolicyIDVal = pingonetypes.NewResourceIDNull()
					}

					mitigationMap := map[string]attr.Value{
						"name":                         framework.StringOkToTF(policy.GetNameOk()),
						"priority":                     framework.Int32OkToTF(policy.GetPriorityOk()),
						"condition":                    conditionObj,
						"action":                       actionVal,
						"custom_action":                customActionVal,
						"mfa_authentication_policy_id": mfaAuthPolicyIDVal,
						"mfa_registration_policy_id":   mfaRegPolicyIDVal,
						"verify_policy_id":             verifyPolicyIDVal,
					}

					mitigationObj, d := types.ObjectValue(mitigationsTFObjectTypes, mitigationMap)
					diags.Append(d...)

					mitigations = append(mitigations, mitigationObj)
					continue
				}
			}
		}

		// Non-mitigation entries: aggregated scores/weights backbone and overrides.
		if condition, ok := policy.GetConditionOk(); ok {
			if v, ok := condition.GetTypeOk(); ok && (*v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES || *v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS) {

				// Policy thresholds medium and high
				if between, ok := condition.GetBetweenOk(); ok {
					betweenObj := map[string]attr.Value{
						"min_score": framework.Int32OkToTF(between.GetMinScoreOk()),
						"max_score": framework.Int32OkToTF(between.GetMaxScoreOk()),
					}

					thresholdObj, d := types.ObjectValue(policyThresholdsTFObjectTypes, betweenObj)
					diags.Append(d...)

					if policy.Result.GetLevel() == risk.ENUMRISKLEVEL_MEDIUM {
						highMediumPolicy["policy_threshold_medium"] = thresholdObj
					}

					if policy.Result.GetLevel() == risk.ENUMRISKLEVEL_HIGH {
						highMediumPolicy["policy_threshold_high"] = thresholdObj
					}

				}

				var d diag.Diagnostics
				// Predictors
				if scores, ok := condition.GetAggregatedScoresOk(); ok && *v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES {
					useScores = true

					tfObjType := types.ObjectType{AttrTypes: policyScoresPredictorTFObjectTypes}

					if len(scores) == 0 {
						highMediumPolicy["predictors"] = types.SetValueMust(tfObjType, []attr.Value{})
					}

					flattenedList := []attr.Value{}
					for _, score := range scores {

						predictor := map[string]attr.Value{
							"predictor_reference_value": framework.StringOkToTF(score.GetValueOk()),
							"compact_name":              riskPolicyScoresCompactNameFromReferenceOk(score.GetValueOk()),
							"score":                     framework.Int32OkToTF(score.GetScoreOk()),
						}

						flattenedObj, d := types.ObjectValue(policyScoresPredictorTFObjectTypes, predictor)
						diags.Append(d...)

						flattenedList = append(flattenedList, flattenedObj)
					}

					highMediumPolicy["predictors"], d = types.SetValue(tfObjType, flattenedList)
					diags.Append(d...)
				}

				if weights, ok := condition.GetAggregatedWeightsOk(); ok && *v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS {
					useWeights = true

					tfObjType := types.ObjectType{AttrTypes: policyWeightsPredictorTFObjectTypes}

					if len(weights) == 0 {
						highMediumPolicy["predictors"] = types.SetValueMust(tfObjType, []attr.Value{})
					}

					flattenedList := []attr.Value{}
					for _, weight := range weights {

						predictor := map[string]attr.Value{
							"predictor_reference_value": framework.StringOkToTF(weight.GetValueOk()),
							"compact_name":              riskPolicyWeightsCompactNameFromReferenceOk(weight.GetValueOk()),
							"weight":                    framework.Int32OkToTF(weight.GetWeightOk()),
						}

						flattenedObj, d := types.ObjectValue(policyWeightsPredictorTFObjectTypes, predictor)
						diags.Append(d...)

						flattenedList = append(flattenedList, flattenedObj)
					}

					highMediumPolicy["predictors"], d = types.SetValue(tfObjType, flattenedList)
					diags.Append(d...)
				}
			}

			if v, ok := condition.GetTypeOk(); ok && (*v == risk.ENUMRISKPOLICYCONDITIONTYPE_VALUE_COMPARISON || *v == risk.ENUMRISKPOLICYCONDITIONTYPE_IP_RANGE) {
				setOverride = true

				resultObj := types.ObjectUnknown(overridesResultTFObjectTypes)
				if policyResult, ok := policy.GetResultOk(); ok {

					resultMap := map[string]attr.Value{
						"value": framework.StringOkToTF(policyResult.GetValueOk()),
						"level": framework.EnumOkToTF(policyResult.GetLevelOk()),
						"type":  framework.EnumOkToTF(policyResult.GetTypeOk()),
					}

					var d diag.Diagnostics
					resultObj, d = types.ObjectValue(overridesResultTFObjectTypes, resultMap)
					diags.Append(d...)
				}

				var equalsString basetypes.StringValue
				if s := condition.GetEquals().String; s != nil {
					equalsString = framework.StringToTF(*s)
				} else {
					equalsString = types.StringNull()
				}

				conditionMap := map[string]attr.Value{
					"type":                         framework.EnumOkToTF(condition.GetTypeOk()),
					"equals":                       equalsString,
					"compact_name":                 riskPolicyOverrideCompactNameFromReferenceOk(condition.GetValueOk()),
					"predictor_reference_value":    framework.StringOkToTF(condition.GetValueOk()),
					"ip_range":                     framework.StringSetOkToTF(condition.GetIpRangeOk()),
					"predictor_reference_contains": framework.StringOkToTF(condition.GetContainsOk()),
				}

				conditionObj, d := types.ObjectValue(overridesConditionTFObjectTypes, conditionMap)
				diags.Append(d...)

				overrideMap := map[string]attr.Value{
					"name":      framework.StringOkToTF(policy.GetNameOk()),
					"priority":  framework.Int32OkToTF(policy.GetPriorityOk()),
					"result":    resultObj,
					"condition": conditionObj,
				}

				overrideObj, d := types.ObjectValue(overridesTFObjectTypes, overrideMap)
				diags.Append(d...)

				overrides = append(overrides, overrideObj)
			}
		}
	}

	var d diag.Diagnostics
	if useScores {
		objPolicyScoresValue, d = types.ObjectValue(policyScoresTFObjectTypes, highMediumPolicy)
		diags.Append(d...)
	} else {
		objPolicyScoresValue = types.ObjectNull(policyScoresTFObjectTypes)
	}

	if useWeights {
		objPolicyWeightsValue, d = types.ObjectValue(policyWeightsTFObjectTypes, highMediumPolicy)
		diags.Append(d...)
	} else {
		objPolicyWeightsValue = types.ObjectNull(policyWeightsTFObjectTypes)
	}

	if setOverride {
		objOverridesValue, d = types.ListValue(types.ObjectType{AttrTypes: overridesTFObjectTypes}, overrides)
		diags.Append(d...)
	} else {
		objOverridesValue = types.ListNull(types.ObjectType{AttrTypes: overridesTFObjectTypes})
	}

	if setMitigation {
		objMitigationsValue, d = types.ListValue(types.ObjectType{AttrTypes: mitigationsTFObjectTypes}, mitigations)
		diags.Append(d...)
	} else {
		objMitigationsValue = types.ListNull(types.ObjectType{AttrTypes: mitigationsTFObjectTypes})
	}

	return objPolicyWeightsValue, objPolicyScoresValue, objOverridesValue, objMitigationsValue, objFallbackValue, diags
}

func riskPolicyScoresCompactNameFromReferenceOk(v *string, ok bool) basetypes.StringValue {
	return riskPolicyCompactNameFromReferenceOk(v, ok, true)
}

func riskPolicyWeightsCompactNameFromReferenceOk(v *string, ok bool) basetypes.StringValue {
	return riskPolicyCompactNameFromReferenceOk(v, ok, false)
}

func riskPolicyOverrideCompactNameFromReferenceOk(v *string, ok bool) basetypes.StringValue {
	return riskPolicyCompactNameFromReferenceOk(v, ok, true)
}

func riskPolicyCompactNameFromReferenceOk(v *string, ok, useScores bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	}

	if useScores {
		return types.StringValue(strings.ReplaceAll(strings.ReplaceAll(*v, "${details.", ""), ".level}", ""))
	} else {
		return types.StringValue(strings.ReplaceAll(strings.ReplaceAll(*v, "${details.aggregatedWeights.", ""), "}", ""))
	}
}

// mitigationPolicyIDToTF converts an optional policy-ID string (from a mitigation inner)
// to a pingonetypes.ResourceIDValue, returning null when the value is absent.
func mitigationPolicyIDToTF(v *string, ok bool) pingonetypes.ResourceIDValue {
	if !ok || v == nil {
		return pingonetypes.NewResourceIDNull()
	}
	return framework.PingOneResourceIDToTF(*v)
}
