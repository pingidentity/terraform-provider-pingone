package risk

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	int64validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int64validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type RiskPolicyResource struct {
	client *risk.APIClient
	region model.RegionMapping
}

type riskPolicyResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	EnvironmentId       types.String `tfsdk:"environment_id"`
	Name                types.String `tfsdk:"name"`
	DefaultResult       types.Object `tfsdk:"default_result"`
	Default             types.Bool   `tfsdk:"default"`
	EvaluatedPredictors types.Set    `tfsdk:"evaluated_predictors"`
	PolicyWeights       types.Object `tfsdk:"policy_weights"`
	PolicyScores        types.Object `tfsdk:"policy_scores"`
	//Overrides             types.Set    `tfsdk:"override"`
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

	var plan riskPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the max threshold score
	var rootPath, referenceValueFmt string
	var maxScore int
	flattenedList := []attr.Value{}
	var predictorAttrType map[string]attr.Type

	if !plan.PolicyWeights.IsNull() && !plan.PolicyWeights.IsUnknown() {
		rootPath = "policy_weights"
		maxScore = 100
		referenceValueFmt = "${details.aggregatedWeights.%s}"
		predictorAttrType = policyWeightsPredictorTFObjectTypes

		var predictorsPlan []riskPolicyResourcePolicyWeightsPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, predictor := range predictorsPlan {
			predictorObj := map[string]attr.Value{
				"predictor_reference_value": framework.StringToTF(fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString())),
				"compact_name":              predictor.CompactName,
				"weight":                    predictor.Weight,
			}

			flattenedObj, d := types.ObjectValue(policyWeightsPredictorTFObjectTypes, predictorObj)
			resp.Diagnostics.Append(d...)

			flattenedList = append(flattenedList, flattenedObj)
		}
	}

	if !plan.PolicyScores.IsNull() && !plan.PolicyScores.IsUnknown() {
		rootPath = "policy_scores"
		maxScore = 1000
		referenceValueFmt = "${details.%s.level}"
		predictorAttrType = policyScoresPredictorTFObjectTypes

		var predictorsPlan []riskPolicyResourcePolicyScoresPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, predictor := range predictorsPlan {
			predictorObj := map[string]attr.Value{
				"predictor_reference_value": framework.StringToTF(fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString())),
				"compact_name":              predictor.CompactName,
				"score":                     predictor.Score,
			}

			flattenedObj, d := types.ObjectValue(policyScoresPredictorTFObjectTypes, predictorObj)
			resp.Diagnostics.Append(d...)

			flattenedList = append(flattenedList, flattenedObj)
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
	plannedPredictors, d := types.SetValue(types.ObjectType{AttrTypes: predictorAttrType}, flattenedList)
	resp.Diagnostics.Append(d...)
	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("predictors"), plannedPredictors)
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
}

func (r *RiskPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state riskPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	riskPolicy, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	createResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskPoliciesApi.CreateRiskPolicySet(ctx, plan.EnvironmentId.ValueString()).RiskPolicySet(*riskPolicy).Execute()
		},
		"CreateRiskPolicySet",
		riskPolicyCreateUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We have to read it back because the API does not return the full state object on create
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), createResponse.(*risk.RiskPolicySet).GetId()).Execute()
		},
		"ReadOneRiskPolicySet",
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
	resp.Diagnostics.Append(state.toState(response.(*risk.RiskPolicySet))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *riskPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneRiskPolicySet",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
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
	resp.Diagnostics.Append(data.toState(response.(*risk.RiskPolicySet))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RiskPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state riskPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	riskPolicy, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskPoliciesApi.UpdateRiskPolicySet(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).RiskPolicySet(*riskPolicy).Execute()
		},
		"UpdateRiskPolicySet",
		riskPolicyCreateUpdateCustomErrorHandler,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*risk.RiskPolicySet))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *riskPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
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
			r, err := r.client.RiskPoliciesApi.DeleteRiskPolicySet(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteRiskPolicySet",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RiskPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/risk_predictor_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
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

func (p *riskPolicyResourceModel) expand(ctx context.Context, apiClient *risk.APIClient) (*risk.RiskPolicySet, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPolicySet(p.Name.ValueString())
	data.SetDefault(false)

	if !p.DefaultResult.IsNull() && !p.DefaultResult.IsUnknown() {
		var plan riskPolicyResourceDefaultResultModel
		d := p.DefaultResult.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
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

		d = p.PolicyWeights.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
	}

	if !p.PolicyScores.IsNull() && !p.PolicyScores.IsUnknown() {
		highPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES)
		mediumPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES)
		useScores = true

		d = p.PolicyScores.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
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

	// High Weighted Policy
	riskPolicies = append(riskPolicies, *risk.NewRiskPolicy(
		*highPolicyCondition,
		"HIGH_WEIGHTED_POLICY",
		*risk.NewRiskPolicyResult(risk.ENUMRISKLEVEL_HIGH),
	))

	// Medium Weighted Policy
	riskPolicies = append(riskPolicies, *risk.NewRiskPolicy(
		*mediumPolicyCondition,
		"MEDIUM_WEIGHTED_POLICY",
		*risk.NewRiskPolicyResult(risk.ENUMRISKLEVEL_MEDIUM),
	))

	riskPolicyPredictorsIDs, d := riskPredictorFetchIDsFromCompactNames(ctx, apiClient, p.EnvironmentId.ValueString(), predictorCompactNames)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	evaluatedPredictors := make([]risk.RiskPolicySetEvaluatedPredictorsInner, 0)

	if !p.EvaluatedPredictors.IsNull() && !p.EvaluatedPredictors.IsUnknown() {
		var plan []string
		d := p.EvaluatedPredictors.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
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

		data.SetEvaluatedPredictors(evaluatedPredictors)
	}

	// TODO: Overrides

	data.SetRiskPolicies(riskPolicies)

	return data, diags
}

func (p *riskPolicyResourcePolicyModel) expand(ctx context.Context, useScores bool, highPolicyCondition, mediumPolicyCondition *risk.RiskPolicyCondition) (*risk.RiskPolicyCondition, *risk.RiskPolicyCondition, []string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !p.PolicyThresholdMedium.IsNull() && !p.PolicyThresholdMedium.IsUnknown() {
		var plan riskPolicyResourcePolicyThresholdScoreBetweenModel
		d := p.PolicyThresholdMedium.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
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
		d := p.PolicyThresholdHigh.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
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
			d := p.Predictors.ElementsAs(ctx, &predictorsPlan, false)
			diags.Append(d...)
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
			d := p.Predictors.ElementsAs(ctx, &predictorsPlan, false)
			diags.Append(d...)
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

	p.PolicyWeights, p.PolicyScores, d = p.toStatePolicy(r, ok)
	diags.Append(d...)

	return diags
}

func (p *riskPolicyResourceModel) toStatePolicy(riskPolicies []risk.RiskPolicy, ok bool) (basetypes.ObjectValue, basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objPolicyWeightsValue := types.ObjectNull(policyWeightsTFObjectTypes)
	objPolicyScoresValue := types.ObjectNull(policyScoresTFObjectTypes)

	useScores := false
	useWeights := false

	if !ok || riskPolicies == nil || len(riskPolicies) < 1 {
		return objPolicyWeightsValue, objPolicyScoresValue, diags
	}

	o := map[string]attr.Value{}

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
						o["policy_threshold_medium"] = thresholdObj
					}

					if policy.Result.GetLevel() == risk.ENUMRISKLEVEL_HIGH {
						o["policy_threshold_high"] = thresholdObj
					}

				}

				var d diag.Diagnostics
				// Predictors
				if scores, ok := condition.GetAggregatedScoresOk(); ok && *v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES {
					useScores = true

					tfObjType := types.ObjectType{AttrTypes: policyScoresPredictorTFObjectTypes}

					if len(scores) == 0 {
						o["predictors"] = types.SetValueMust(tfObjType, []attr.Value{})
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

					o["predictors"], d = types.SetValue(tfObjType, flattenedList)
					diags.Append(d...)
				}

				if weights, ok := condition.GetAggregatedWeightsOk(); ok && *v == risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS {
					useWeights = true

					tfObjType := types.ObjectType{AttrTypes: policyWeightsPredictorTFObjectTypes}

					if len(weights) == 0 {
						o["predictors"] = types.SetValueMust(tfObjType, []attr.Value{})
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

					o["predictors"], d = types.SetValue(tfObjType, flattenedList)
					diags.Append(d...)
				}
			}
		}
	}

	var d diag.Diagnostics
	if useScores {
		objPolicyScoresValue, d = types.ObjectValue(policyScoresTFObjectTypes, o)
		diags.Append(d...)
	}

	if useWeights {
		objPolicyWeightsValue, d = types.ObjectValue(policyWeightsTFObjectTypes, o)
		diags.Append(d...)
	}

	return objPolicyWeightsValue, objPolicyScoresValue, diags
}

func riskPolicyScoresCompactNameFromReferenceOk(v *string, ok bool) basetypes.StringValue {
	return riskPolicyCompactNameFromReferenceOk(v, ok, true)
}

func riskPolicyWeightsCompactNameFromReferenceOk(v *string, ok bool) basetypes.StringValue {
	return riskPolicyCompactNameFromReferenceOk(v, ok, false)
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
