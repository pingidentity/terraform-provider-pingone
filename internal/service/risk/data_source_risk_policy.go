// Copyright © 2026 Ping Identity Corporation

package risk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type RiskPolicyDataSource serviceClientType

type riskPolicyDataSourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	RiskPolicyId        pingonetypes.ResourceIDValue `tfsdk:"risk_policy_id"`
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

// Framework interfaces
var (
	_ datasource.DataSource              = &RiskPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &RiskPolicyDataSource{}
)

func NewRiskPolicyDataSource() datasource.DataSource {
	return &RiskPolicyDataSource{}
}

// Metadata
func (r *RiskPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_policy"
}

// Schema
func (r *RiskPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const attrMinLength = 1

	riskPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the risk policy to retrieve.",
	).ExactlyOneOf([]string{"risk_policy_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the risk policy. Matching is exact and case-sensitive. If multiple risk policies share the same name, the data source returns an error; use `risk_policy_id` instead.",
	).ExactlyOneOf([]string{"risk_policy_id", "name"})

	resp.Schema = schema.Schema{
		Description: "Data source to retrieve a PingOne Risk Policy.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to retrieve the risk policy from.",
				Required:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"risk_policy_id": schema.StringAttribute{
				Description:         riskPolicyIdDescription.Description,
				MarkdownDescription: riskPolicyIdDescription.MarkdownDescription,
				Optional:            true,
				CustomType:          pingonetypes.ResourceIDType{},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
				},
			},
			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
			"default_result": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies the default result value for the risk policy.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"level": schema.StringAttribute{
						Description: "The default result level.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The default result type.",
						Computed:    true,
					},
				},
			},
			"default": schema.BoolAttribute{
				Description: "A boolean that indicates whether this risk policy set is the environment's default risk policy set. This is used whenever an explicit policy set ID is not specified in a risk evaluation request.",
				Computed:    true,
			},
			"evaluated_predictors": schema.SetAttribute{
				Description: "A set of IDs for the predictors to evaluate in this policy set.  If omitted, if this property is null, all of the licensed predictors are used.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"policy_weights": policySchemaDataSource(
				framework.SchemaAttributeDescriptionFromMarkdown(
					"An object that describes settings for a risk policy using a weighted average calculation, with a final result being a risk score between `0` and `10`.",
				),
				false,
			),
			"policy_scores": policySchemaDataSource(
				framework.SchemaAttributeDescriptionFromMarkdown(
					"An object that describes settings for a risk policy calculated by aggregating score values, with a final result being the sum of score values from each of the configured predictors.",
				),
				true,
			),
			"overrides": schema.ListNestedAttribute{
				Description: "An ordered list of policy overrides to apply to the policy.  The ordering of the overrides is important as it determines the priority of the policy override during policy evaluation.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "A string that represents the name of the overriding risk policy in the set.",
							Computed:    true,
						},
						"priority": schema.Int32Attribute{
							Description: "An integer that indicates the order in which the override is applied during risk policy evaluation.  The lower the value, the higher the priority.",
							Computed:    true,
						},
						"result": schema.SingleNestedAttribute{
							Description: "A single object that contains the risk result that should be applied to the policy evaluation result when the override condition is met.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"value": schema.StringAttribute{
									Description: "An administrator defined string value that is applied to the policy evaluation result when the override condition is met.",
									Computed:    true,
								},
								"level": schema.StringAttribute{
									Description: "A string that specifies the risk level that should be applied to the policy evalution result when the override condition is met.",
									Computed:    true,
								},
								"type": schema.StringAttribute{
									Description: "A string that specifies the type of the risk result should be applied to the policy evalution result when the override condition is met.",
									Computed:    true,
								},
							},
						},
						"condition": schema.SingleNestedAttribute{
							Description: "A single object that contains the conditions to evaluate that determine whether the override result will be applied to the risk policy evaluation.",
							Computed:    true,
							Attributes:  overridesConditionAttributesDataSource(),
						},
					},
				},
			},
			"mitigations": schema.ListNestedAttribute{
				Description: "An ordered list of mitigation-style policy entries to apply to the policy. Each entry pairs a condition with a single mitigation action.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: mitigationEntryAttributesDataSource(),
				},
			},
			"fallback": schema.SingleNestedAttribute{
				Description: "A single object that specifies the required catch-all fallback mitigation action, applied when no mitigation condition matches.",
				Computed:    true,
				Attributes:  mitigationActionAttributesDataSource(),
			},
			"targets": schema.SingleNestedAttribute{
				Description: "A single object that scopes this policy set to a subset of events (targeted policy).",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"condition": schema.SingleNestedAttribute{
						Description: "A single object that specifies the AND-of-sub-conditions targeting condition. All sub-conditions in `and` must be satisfied for the policy set to be selected.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"and": schema.ListNestedAttribute{
								Description: "An ordered list of sub-conditions that are combined with AND logic. Each entry pairs a `list` of values with the event attribute (`contains`) to check against.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "A read-only string that identifies the sub-condition kind. Inferred by the API from `contains`.",
											Computed:    true,
										},
										"list": schema.ListAttribute{
											Description: "A list of values to match against the event attribute specified in `contains`. Transaction types are one or more of the supported flow types. User groups are group names. Applications are PingOne application IDs.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"contains": schema.StringAttribute{
											Description: "The event attribute checked against `list`. For transaction types use `${event.flow.type}`; for user groups use `${event.user.groups}`; for applications use `${event.targetResource.id}`.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func policyThresholdAttributesDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"min_score": schema.Int32Attribute{
			Description: "An integer that specifies the minimum score to use as the lower bound value of the policy threshold.",
			Computed:    true,
		},
		"max_score": schema.Int32Attribute{
			Description: "An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.",
			Computed:    true,
		},
	}
}

func policySchemaDataSource(description framework.SchemaAttributeDescription, useScores bool) schema.SingleNestedAttribute {
	thresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the lower and upper bound threshold score values that define the medium risk outcome as a result of the policy evaluation.",
	)

	predictorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes a predictor to apply to the risk policy and its associated weight or score value for the overall risk calculation.",
	)

	predictorAttributes := map[string]schema.Attribute{
		"compact_name": schema.StringAttribute{
			Description: "A string that specifies the compact name of the predictor to apply to the risk policy.",
			Computed:    true,
		},
		"predictor_reference_value": schema.StringAttribute{
			Description: "A string that specifies the attribute reference of the level to evaluate.",
			Computed:    true,
		},
	}

	if useScores {
		predictorAttributes["score"] = schema.Int32Attribute{
			Description: "An integer that specifies the score to apply to the High risk / true outcome of the predictor, to apply to the overall risk calculation.",
			Computed:    true,
		}
	} else {
		predictorAttributes["weight"] = schema.Int32Attribute{
			Description: "An integer that specifies the weight to apply to the predictor when calculating the overall risk score.",
			Computed:    true,
		}
	}

	return schema.SingleNestedAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"policy_threshold_medium": schema.SingleNestedAttribute{
				Description:         thresholdDescription.Description,
				MarkdownDescription: thresholdDescription.MarkdownDescription,
				Computed:            true,
				Attributes:          policyThresholdAttributesDataSource(),
			},
			"policy_threshold_high": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"An object that specifies the lower and upper bound threshold score values that define the high risk outcome as a result of the policy evaluation.",
				).Description,
				Computed:   true,
				Attributes: policyThresholdAttributesDataSource(),
			},
			"predictors": schema.SetNestedAttribute{
				Description:         predictorDescription.Description,
				MarkdownDescription: predictorDescription.MarkdownDescription,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: predictorAttributes,
				},
			},
		},
	}
}

func overridesConditionAttributesDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "A string that specifies the type of the override condition to evaluate.",
			Computed:    true,
		},
		"equals": schema.StringAttribute{
			Description: "A string that specifies the value of the `predictor_reference_value` that must be matched for the override result to be applied to the policy evaluation.",
			Computed:    true,
		},
		"compact_name": schema.StringAttribute{
			Description: "A string that specifies the compact name of the predictor to apply to the override condition.",
			Computed:    true,
		},
		"predictor_reference_value": schema.StringAttribute{
			Description: "A string that specifies the attribute reference of the value to evaluate.",
			Computed:    true,
		},
		"ip_range": schema.SetAttribute{
			Description: "A set of strings that specifies the CIDR ranges that should be evaluated against the value of the `predictor_reference_contains` attribute.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"predictor_reference_contains": schema.StringAttribute{
			Description: "A string that specifies the attribute reference of the collection to evaluate.",
			Computed:    true,
		},
	}
}

func mitigationActionAttributesDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Description: "A string that specifies the mitigation action to apply when the condition is met.",
			Computed:    true,
		},
		"custom_action": schema.StringAttribute{
			Description: "A string that specifies the custom action name. Set when `action` is `CUSTOM`.",
			Computed:    true,
		},
		"mfa_authentication_policy_id": schema.StringAttribute{
			Description: "The ID of the MFA (sign-on/authentication) policy to apply. Set when `action` is `MFA`.",
			Computed:    true,
			CustomType:  pingonetypes.ResourceIDType{},
		},
		"mfa_registration_policy_id": schema.StringAttribute{
			Description: "The ID of the MFA registration policy to apply. Applies to MFA registration flows when `action` is `MFA`.",
			Computed:    true,
			CustomType:  pingonetypes.ResourceIDType{},
		},
		"verify_policy_id": schema.StringAttribute{
			Description: "The ID of the PingOne Verify policy to apply. Set when `action` is `VERIFY`.",
			Computed:    true,
			CustomType:  pingonetypes.ResourceIDType{},
		},
	}
}

func mitigationEntryAttributesDataSource() map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "A string that represents the name of the mitigation policy entry.",
			Computed:    true,
		},
		"priority": schema.Int32Attribute{
			Description: "An integer that indicates the order in which the mitigation entry is applied during risk policy evaluation. The lower the value, the higher the priority.",
			Computed:    true,
		},
		"condition": schema.SingleNestedAttribute{
			Description: "A single object that contains the conditions to evaluate that determine whether the mitigation action will be applied to the risk policy evaluation.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "A string that specifies the type of the override condition to evaluate.",
					Computed:    true,
				},
				"equals": schema.StringAttribute{
					Description: "A string that specifies the value of the `predictor_reference_value` that must be matched for the override result to be applied to the policy evaluation.",
					Computed:    true,
				},
				"compact_name": schema.StringAttribute{
					Description: "A string that specifies the compact name of the predictor to apply to the override condition.",
					Computed:    true,
				},
				"predictor_reference_value": schema.StringAttribute{
					Description: "A string that specifies the attribute reference of the value to evaluate.",
					Computed:    true,
				},
				"predictor_reference_contains": schema.StringAttribute{
					Description: "A string that specifies the attribute reference of the collection to evaluate.",
					Computed:    true,
				},
			},
		},
	}

	for k, v := range mitigationActionAttributesDataSource() {
		attributes[k] = v
	}

	return attributes
}

// Configure
func (r *RiskPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this to the provider maintainers.", req.ProviderData),
		)

		return
	}

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this to the provider maintainers.",
		)
		return
	}
}

// Read
func (r *RiskPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data riskPolicyDataSourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this to the provider maintainers.",
		)
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	riskPolicyID := data.RiskPolicyId.ValueString()

	if data.RiskPolicyId.IsNull() {
		// Resolve the risk policy ID from the configured name
		var found *risk.RiskPolicySet
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.findRiskPolicyByName(ctx, data.EnvironmentId.ValueString(), data.Name.ValueString())
			},
			"ReadRiskPolicySets",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&found,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		riskPolicyID = found.GetId()
	}

	// Run the API call
	var riskPolicy *risk.RiskPolicySet
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, data.EnvironmentId.ValueString(), riskPolicyID).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPolicySet",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&riskPolicy,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(riskPolicy)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RiskPolicyDataSource) findRiskPolicyByName(ctx context.Context, environmentID, name string) (*risk.RiskPolicySet, *http.Response, error) {
	notFoundMessage := fmt.Sprintf("Risk Policy with name %s not found", name)
	duplicateMessage := fmt.Sprintf("Multiple risk policies found with name %s.  Use risk_policy_id to select one", name)

	pagedIterator := r.Client.RiskAPIClient.RiskPoliciesApi.ReadRiskPolicySets(ctx, environmentID).Execute()

	var initialHttpResponse *http.Response

	var found *risk.RiskPolicySet

	for pageCursor, err := range pagedIterator {
		if err != nil {
			_, resp, err := legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, pageCursor.HTTPResponse, err)
			return nil, resp, err
		}

		if initialHttpResponse == nil {
			initialHttpResponse = pageCursor.HTTPResponse
		}

		if riskPolicySets, ok := pageCursor.EntityArray.Embedded.GetRiskPolicySetsOk(); ok {
			for i := range riskPolicySets {
				if riskPolicySets[i].GetName() == name {
					if found != nil {
						return nil, pageCursor.HTTPResponse, errors.New(duplicateMessage)
					}

					found = &riskPolicySets[i]
				}
			}
		}
	}

	if found == nil {
		return nil, initialHttpResponse, errors.New(notFoundMessage)
	}

	return found, initialHttpResponse, nil
}

func (p *riskPolicyDataSourceModel) toState(apiObject *risk.RiskPolicySet) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	resourceModel := &riskPolicyResourceModel{}
	diags.Append(resourceModel.toState(apiObject)...)
	if diags.HasError() {
		return diags
	}

	p.Id = resourceModel.Id
	p.RiskPolicyId = resourceModel.Id
	p.Name = resourceModel.Name
	p.DefaultResult = resourceModel.DefaultResult
	p.Default = resourceModel.Default
	p.EvaluatedPredictors = resourceModel.EvaluatedPredictors
	p.PolicyWeights = resourceModel.PolicyWeights
	p.PolicyScores = resourceModel.PolicyScores
	p.Overrides = resourceModel.Overrides
	p.Mitigations = resourceModel.Mitigations
	p.Fallback = resourceModel.Fallback
	p.Targets = resourceModel.Targets

	// The resource toStatePolicy initialises these as unknown when the API returns
	// no inner riskPolicies; unknown values are not valid in data source state.
	if p.PolicyWeights.IsUnknown() {
		p.PolicyWeights = types.ObjectNull(policyWeightsTFObjectTypes)
	}

	if p.PolicyScores.IsUnknown() {
		p.PolicyScores = types.ObjectNull(policyScoresTFObjectTypes)
	}

	if p.Overrides.IsUnknown() {
		p.Overrides = types.ListNull(types.ObjectType{AttrTypes: overridesTFObjectTypes})
	}

	return diags
}
