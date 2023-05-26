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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
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

type riskPolicyResourceThresholdScoresModel struct {
	High   types.Int64 `tfsdk:"high"`
	Medium types.Int64 `tfsdk:"medium"`
}

type riskPolicyResourcePredictorScoreModel struct {
	PredictorReference types.String `tfsdk:"predictor_reference_value"`
	Score              types.Int64  `tfsdk:"score"`
}

type riskPolicyResourceOverrideModel struct {
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
	)

	policyWeightedAveragePredictor := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes a predictor to apply to the risk policy and its associated weight value for the overall weighted average risk calculation.",
	)

	// Scores policy
	policyScoresDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes settings for a risk policy calculated by aggregating score values, with a final result being the sum of score values from each of the configured predictors.",
	)

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

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         defaultResultTypeDescription.Description,
						MarkdownDescription: defaultResultTypeDescription.MarkdownDescription,
						Computed:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},

					"level": schema.StringAttribute{
						Description:         defaultResultLevelDescription.Description,
						MarkdownDescription: defaultResultLevelDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(risk.ENUMRISKLEVEL_LOW)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumRiskLevelEnumValues)...),
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
				Computed:    true,

				ElementType: types.StringType,
			},

			"policy_weights": schema.SingleNestedAttribute{
				Description:         policyWeightedAverageDescription.Description,
				MarkdownDescription: policyWeightedAverageDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"policy_threshold_medium": riskPolicyThresholdSchema(
						false,
						1,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation.",
						),
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						false,
						2,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation.",
						),
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
										int64validator.AtLeast(1),
										int64validator.AtMost(10),
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
						40,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation.",
						),
					),

					"policy_threshold_high": riskPolicyThresholdSchema(
						true,
						75,
						framework.SchemaAttributeDescriptionFromMarkdown(
							"An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation.",
						),
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
										int64validator.AtLeast(1),
										int64validator.AtMost(100),
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
		},
	}
}

func riskPolicyThresholdSchema(useScores bool, defaultPolicyThresholdMinScore int64, policyThresholdsDescription framework.SchemaAttributeDescription) schema.SingleNestedAttribute {

	validators := []validator.Int64{
		int64validator.AtLeast(1),
		// TODO medium must be less than high rule
	}

	if !useScores {
		validators = append(validators, int64validator.AtMost(10))
	}

	policyThresholdScoresMediumScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum score to use as the lower bound value of the policy threshold.",
	).DefaultValue(fmt.Sprint(defaultPolicyThresholdMinScore))

	policyThresholdScoresHighScoreDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.",
	)

	return schema.SingleNestedAttribute{
		Description:         policyThresholdsDescription.Description,
		MarkdownDescription: policyThresholdsDescription.MarkdownDescription,
		Optional:            true,
		Computed:            true,

		Attributes: map[string]schema.Attribute{
			"min_score": schema.Int64Attribute{
				Description:         policyThresholdScoresMediumScoreDescription.Description,
				MarkdownDescription: policyThresholdScoresMediumScoreDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(defaultPolicyThresholdMinScore),

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

	if !plan.PolicyWeights.IsNull() && !plan.PolicyWeights.IsUnknown() {
		rootPath = "policy_weights"
		maxScore = 10
		referenceValueFmt = "${details.aggregatedWeights.%s}"

		var predictorsPlan []riskPolicyResourcePolicyWeightsPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, predictor := range predictorsPlan {
			tflog.Debug(ctx, "HERE!!!", map[string]interface{}{
				"predictor": fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString()),
			})
		}
	}

	if !plan.PolicyScores.IsNull() && !plan.PolicyScores.IsUnknown() {
		rootPath = "policy_scores"
		maxScore = 1000
		referenceValueFmt = "${details.%s.level}"

		var predictorsPlan []riskPolicyResourcePolicyScoresPredictorModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("predictors"), &predictorsPlan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, predictor := range predictorsPlan {
			tflog.Debug(ctx, "HERE!!!", map[string]interface{}{
				"predictor": fmt.Sprintf(referenceValueFmt, predictor.CompactName.ValueString()),
			})
		}
	}

	var policyThresholdHighMinValue *int64
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("min_score"), &policyThresholdHighMinValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_medium").AtName("max_score"), types.Int64Value(*policyThresholdHighMinValue))
	resp.Plan.SetAttribute(ctx, path.Root(rootPath).AtName("policy_threshold_high").AtName("max_score"), int64(maxScore))
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
	riskPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
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

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response.(*risk.RiskPolicySet))...)
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
	resp.Diagnostics.Append(data.toState(ctx, response.(*risk.RiskPolicySet))...)
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
	riskPolicy, d := plan.expand(ctx)
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
	resp.Diagnostics.Append(state.toState(ctx, response.(*risk.RiskPolicySet))...)
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

func (p *riskPolicyResourceModel) expand(ctx context.Context) (*risk.RiskPolicySet, diag.Diagnostics) {
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
			data.SetDefaultResult(*risk.NewRiskPolicyResult(risk.EnumRiskLevel(plan.Level.ValueString())))
		}
	}

	highPolicyCondition := risk.NewRiskPolicyCondition()
	mediumPolicyCondition := risk.NewRiskPolicyCondition()

	if !p.PolicyWeights.IsNull() && !p.PolicyWeights.IsUnknown() {
		highPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_WEIGHTS)

		var plan riskPolicyResourcePolicyModel
		d := p.PolicyWeights.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		highPolicyCondition, mediumPolicyCondition, d = plan.expand(ctx, false, highPolicyCondition, mediumPolicyCondition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	if !p.PolicyScores.IsNull() && !p.PolicyScores.IsUnknown() {
		highPolicyCondition.SetType(risk.ENUMRISKPOLICYCONDITIONTYPE_AGGREGATED_SCORES)

		var plan riskPolicyResourcePolicyModel
		d := p.PolicyScores.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		highPolicyCondition, mediumPolicyCondition, d = plan.expand(ctx, true, highPolicyCondition, mediumPolicyCondition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
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

	// TODO: Overrides

	data.SetRiskPolicies(riskPolicies)

	return data, diags
}

func (p *riskPolicyResourcePolicyModel) expand(ctx context.Context, useScores bool, highPolicyCondition, mediumPolicyCondition *risk.RiskPolicyCondition) (*risk.RiskPolicyCondition, *risk.RiskPolicyCondition, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !p.PolicyThresholdMedium.IsNull() && !p.PolicyThresholdMedium.IsUnknown() {
		var plan riskPolicyResourcePolicyThresholdScoreBetweenModel
		d := p.PolicyThresholdMedium.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}

		mediumPolicyCondition.SetBetween(
			*risk.NewRiskPolicyConditionBetween(
				int32(plan.MinScore.ValueInt64()),
				int32(plan.MaxScore.ValueInt64()),
			),
		)
	}

	if !p.PolicyThresholdHigh.IsNull() && !p.PolicyThresholdHigh.IsUnknown() {
		var plan riskPolicyResourcePolicyThresholdScoreBetweenModel
		d := p.PolicyThresholdHigh.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
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
				return nil, nil, diags
			}

			for _, predictor := range predictorsPlan {
				aggregatedScores = append(
					aggregatedScores,
					*risk.NewRiskPolicyConditionAggregatedScoresInner(
						predictor.PredictorReferenceValue.ValueString(),
						int32(predictor.Score.ValueInt64()),
					),
				)
			}

			mediumPolicyCondition.SetAggregatedScores(aggregatedScores)
			highPolicyCondition.SetAggregatedScores(aggregatedScores)

		} else {
			aggregatedWeights := make([]risk.RiskPolicyConditionAggregatedWeightsInner, 0)

			var predictorsPlan []riskPolicyResourcePolicyWeightsPredictorModel
			d := p.Predictors.ElementsAs(ctx, &predictorsPlan, false)
			diags.Append(d...)
			if diags.HasError() {
				return nil, nil, diags
			}

			for _, predictor := range predictorsPlan {
				aggregatedWeights = append(
					aggregatedWeights,
					*risk.NewRiskPolicyConditionAggregatedWeightsInner(
						predictor.PredictorReferenceValue.ValueString(),
						int32(predictor.Weight.ValueInt64()),
					),
				)
			}

			mediumPolicyCondition.SetAggregatedWeights(aggregatedWeights)
			highPolicyCondition.SetAggregatedWeights(aggregatedWeights)
		}
	}

	return highPolicyCondition, mediumPolicyCondition, diags
}

func (p *riskPolicyResourceModel) toState(ctx context.Context, apiObject *risk.RiskPolicySet) diag.Diagnostics {
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

	p.PolicyWeights = types.ObjectNull(policyWeightsTFObjectTypes)

	p.PolicyScores = types.ObjectNull(policyScoresTFObjectTypes)

	return diags
}
