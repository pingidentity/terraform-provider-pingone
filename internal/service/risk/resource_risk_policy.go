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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	int64validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int64validator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type RiskPolicyResource serviceClientType

type riskPolicyResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	EnvironmentId       types.String `tfsdk:"environment_id"`
	Name                types.String `tfsdk:"name"`
	DefaultResult       types.Object `tfsdk:"default_result"`
	Default             types.Bool   `tfsdk:"default"`
	EvaluatedPredictors types.Set    `tfsdk:"evaluated_predictors"`
	PolicyWeights       types.Object `tfsdk:"policy_weights"`
	PolicyScores        types.Object `tfsdk:"policy_scores"`
	Overrides           types.List   `tfsdk:"overrides"`
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
	MinScore types.Int64 `tfsdk:"min_score"`
	MaxScore types.Int64 `tfsdk:"max_score"`
}

type riskPolicyResourcePolicyWeightsPredictorModel struct {
	CompactName             types.String `tfsdk:"compact_name"`
	PredictorReferenceValue types.String `tfsdk:"predictor_reference_value"`
	Weight                  types.Int64  `tfsdk:"weight"`
}

type riskPolicyResourcePolicyScoresPredictorModel struct {
	CompactName             types.String `tfsdk:"compact_name"`
	PredictorReferenceValue types.String `tfsdk:"predictor_reference_value"`
	Score                   types.Int64  `tfsdk:"score"`
}

type riskPolicyResourcePolicyOverrideModel struct {
	Name      types.String `tfsdk:"name"`
	Priority  types.Int64  `tfsdk:"priority"`
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

var (
	policyThresholdsTFObjectTypes = map[string]attr.Type{
		"min_score": types.Int64Type,
		"max_score": types.Int64Type,
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
		"weight":                    types.Int64Type,
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
		"score":                     types.Int64Type,
	}

	overridesTFObjectTypes = map[string]attr.Type{
		"name":     types.StringType,
		"priority": types.Int64Type,
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
	).AllowedValuesEnum(risk.AllowedEnumResultTypeEnumValues)

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
	).AllowedValuesEnum(risk.AllowedEnumResultTypeEnumValues).DefaultValue(string(risk.ENUMRESULTTYPE_VALUE))

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
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}#\/.'_\s-]+$`), " Valid characters consist of any Unicode letter, mark (such as, accent, umlaut), # (numeric), / (forward slash), . (period), ' (apostrophe), _ (underscore), space, or - (hyphen)."),
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
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"evaluated_predictors": schema.SetAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of IDs for the predictors to evaluate in this policy set.  If omitted, if this property is null, all of the licensed predictors are used.").Description,
				Optional:    true,
				Computed:    true,

				ElementType: types.StringType,

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
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
						[]validator.Int64{
							int64validatorinternal.IsLessThanPathValue(
								path.MatchRoot("policy_weights").AtName("policy_threshold_high").AtName("min_score"),
							),
						},
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						false,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold score values that define the high risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int64{
							int64validatorinternal.IsGreaterThanPathValue(
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

								"weight": schema.Int64Attribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the weight to apply to the predictor when calculating the overall risk score.").Description,
									Required:    true,

									Validators: []validator.Int64{
										int64validator.AtLeast(weightMinimumDefault),
										int64validator.AtMost(weightMaximumDefault),
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
						[]validator.Int64{
							int64validatorinternal.IsLessThanPathValue(
								path.MatchRoot("policy_scores").AtName("policy_threshold_high").AtName("min_score"),
							),
						},
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						true,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation.",
						),
						[]validator.Int64{
							int64validatorinternal.IsGreaterThanPathValue(
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

								"score": schema.Int64Attribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the score to apply to the High risk / true outcome of the predictor, to apply to the overall risk calculation.").Description,
									Required:    true,

									Validators: []validator.Int64{
										int64validator.AtLeast(scoreMinimumDefault),
										int64validator.AtMost(scoreMaximumDefault),
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

						"priority": schema.Int64Attribute{
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
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumResultTypeEnumValues)...),
									},

									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
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
				},
			},
		},
	}
}

func riskPolicyThresholdSchema(useScores bool, policyThresholdsDescription framework.SchemaAttributeDescription, validators []validator.Int64) schema.SingleNestedAttribute {

	validators = append(validators, int64validator.AtLeast(1))

	policyThresholdScoresMediumScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum score to use as the lower bound value of the policy threshold.",
	)

	policyThresholdScoresHighScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.",
	)

	if !useScores {
		maxAllowedValue := 100
		denominator := 10
		validators = append(validators, int64validator.AtMost(int64(maxAllowedValue)))
		validators = append(validators, int64validatorinternal.IsDivisibleBy(int64(denominator)))

		policyThresholdScoresMediumScoreDescription = policyThresholdScoresMediumScoreDescription.AppendMarkdownString(fmt.Sprintf("For weights policies, the score values should be 10x the desired risk value in the console. For example, a risk score of `5` in the console should be entered as `50`.  The provided score must be exactly divisible by 10.  Maximum value allowed is `%d`", maxAllowedValue))
	} else {
		maxAllowedValue := 1000
		validators = append(validators, int64validator.AtMost(int64(maxAllowedValue)))
		policyThresholdScoresMediumScoreDescription = policyThresholdScoresMediumScoreDescription.AppendMarkdownString(fmt.Sprintf("Maximum value allowed is `%d`", maxAllowedValue))
	}

	return schema.SingleNestedAttribute{
		Description:         policyThresholdsDescription.Description,
		MarkdownDescription: policyThresholdsDescription.MarkdownDescription,
		Required:            true,

		Attributes: map[string]schema.Attribute{
			"min_score": schema.Int64Attribute{
				Description:         policyThresholdScoresMediumScoreDescription.Description,
				MarkdownDescription: policyThresholdScoresMediumScoreDescription.MarkdownDescription,
				Required:            true,

				Validators: validators,
			},

			"max_score": schema.Int64Attribute{
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
	var policyThresholdHighMinValue *int64
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("min_score"), &policyThresholdHighMinValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_medium").AtName("max_score"), types.Int64Value(*policyThresholdHighMinValue))
	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("max_score"), types.Int64Value(int64(maxScore)))

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
				"priority":  types.Int64Value(int64(priorityCount)),
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

	if setEvaluatedPredictorsToUnknown {
		resp.Plan.SetAttribute(ctx, path.Root("evaluated_predictors"), types.SetUnknown(types.StringType))
	}
}

func (r *RiskPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RiskPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state riskPolicyResourceModel

	if r.Client.RiskAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.CreateRiskPolicySet(ctx, plan.EnvironmentId.ValueString()).RiskPolicySet(*riskPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), createResponse.GetId()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPolicySet",
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

func (r *RiskPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *riskPolicyResourceModel

	if r.Client.RiskAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPolicySet",
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

func (r *RiskPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state riskPolicyResourceModel

	if r.Client.RiskAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.UpdateRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).RiskPolicySet(*riskPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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

	if r.Client.RiskAPIClient == nil {
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
			fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.DeleteRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteRiskPolicySet",
		framework.CustomErrorResourceNotFoundWarning,
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
		},
		Refresh: func() (interface{}, string, error) {
			base := 10

			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			resp, r, err := framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil {
				if r.StatusCode == 404 {
					return risk.RiskPolicySet{}, strconv.FormatInt(int64(r.StatusCode), base), nil
				}
				return nil, strconv.FormatInt(int64(r.StatusCode), base), err
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

func riskPolicyCreateUpdateCustomErrorHandler(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Invalid composition
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if target, ok := details[0].GetTargetOk(); ok && *target == "composition.condition" {
			diags.AddError(
				"Invalid \"composition.condition\" policy JSON.",
				"Please check the \"composition.condition\" policy JSON structure and contents and try again.",
			)

			return diags
		}
	}

	return nil
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
				var conditionIPRangePlan []string
				diags.Append(conditionPlan.IPRange.ElementsAs(ctx, &conditionIPRangePlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				condition.SetIpRange(conditionIPRangePlan)
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

			result := risk.NewRiskPolicyResult(risk.EnumRiskLevel(resultPlan.Level.ValueString()))

			if !resultPlan.Value.IsNull() && !resultPlan.Value.IsUnknown() {
				result.SetValue(resultPlan.Value.ValueString())
			}

			if !resultPlan.Type.IsNull() && !resultPlan.Type.IsUnknown() {
				result.SetType(risk.EnumResultType(resultPlan.Type.ValueString()))
			}

			op := *risk.NewRiskPolicy(
				*condition,
				overridePlan.Name.ValueString(),
				*result,
			)

			op.SetPriority(int32(overridePlan.Priority.ValueInt64()))

			riskPolicies = append(riskPolicies, op)
		}
	}

	// Medium Weighted Policy
	mwp := *risk.NewRiskPolicy(
		*mediumPolicyCondition,
		"MEDIUM_WEIGHTED_POLICY",
		*risk.NewRiskPolicyResult(risk.ENUMRISKLEVEL_MEDIUM),
	)
	mwp.SetPriority(int32(len(riskPolicies)) + 1)
	riskPolicies = append(riskPolicies, mwp)

	// High Weighted Policy
	hwp := *risk.NewRiskPolicy(
		*highPolicyCondition,
		"HIGH_WEIGHTED_POLICY",
		*risk.NewRiskPolicyResult(risk.ENUMRISKLEVEL_HIGH),
	)
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
		var plan []string
		diags.Append(p.EvaluatedPredictors.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		configuredEvaluatedPredictorIDs := make(map[string]bool)

		for _, predictorID := range plan {
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
				int32(plan.MinScore.ValueInt64()),
				int32(plan.MaxScore.ValueInt64()),
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
				int32(plan.MinScore.ValueInt64()),
				int32(plan.MaxScore.ValueInt64()),
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
						int32(predictor.Score.ValueInt64()),
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
						int32(predictor.Weight.ValueInt64()),
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

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
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

	p.PolicyWeights, p.PolicyScores, p.Overrides, d = p.toStatePolicy(r, ok)
	diags.Append(d...)

	return diags
}

func (p *riskPolicyResourceModel) toStatePolicy(riskPolicies []risk.RiskPolicy, ok bool) (basetypes.ObjectValue, basetypes.ObjectValue, basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objPolicyWeightsValue := types.ObjectUnknown(policyWeightsTFObjectTypes)
	objPolicyScoresValue := types.ObjectUnknown(policyScoresTFObjectTypes)
	objOverridesValue := types.ListUnknown(types.ObjectType{AttrTypes: overridesTFObjectTypes})

	useScores := false
	useWeights := false

	if !ok || riskPolicies == nil || len(riskPolicies) < 1 {
		return objPolicyWeightsValue, objPolicyScoresValue, objOverridesValue, diags
	}

	highMediumPolicy := map[string]attr.Value{}
	overrides := []attr.Value{}

	setOverride := false

	for _, policy := range riskPolicies {
		// First build the high and medium outcome policies

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

	return objPolicyWeightsValue, objPolicyScoresValue, objOverridesValue, diags
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
		return types.StringValue(strings.Replace(strings.Replace(*v, "${details.", "", -1), ".level}", "", -1))
	} else {
		return types.StringValue(strings.Replace(strings.Replace(*v, "${details.aggregatedWeights.", "", -1), "}", "", -1))
	}
}
