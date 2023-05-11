package risk

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type RiskPredictorResource struct {
	client *risk.APIClient
	region model.RegionMapping
}

type riskPredictorResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	CompactName   types.String `tfsdk:"compact_name"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
	Default       types.Object `tfsdk:"default"`
	Licensed      types.Bool   `tfsdk:"licensed"`
	Deletable     types.Bool   `tfsdk:"deletable"`
	// Anonymous network, IP reputation, geovelocity
	AllowedCIDRList types.Set `tfsdk:"allowed_cidr_list"`
	// Custom map
	CustomMap types.Object `tfsdk:"custom_map"`
	// New device
	ActivationAt types.String `tfsdk:"activation_at"`
	Detect       types.String `tfsdk:"detect"`
	// User Location Anomaly
	Radius types.Object `tfsdk:"radius"`
	Days   types.Int64  `tfsdk:"days"`
	// User Risk Behavior
	PredictionModel types.Object `tfsdk:"prediction_model"`
	// Velocity
	By            types.Set    `tfsdk:"by"`
	Every         types.Object `tfsdk:"every"`
	Fallback      types.Object `tfsdk:"fallback"`
	Measure       types.String `tfsdk:"measure"`
	Of            types.String `tfsdk:"of"`
	SlidingWindow types.Object `tfsdk:"sliding_window"`
	Use           types.Object `tfsdk:"use"`
}

type predictorDefault struct {
	Weight types.Int64  `tfsdk:"weight"`
	Result types.Object `tfsdk:"result"`
}

type predictorDefaultResult struct {
	ResultType types.String `tfsdk:"type"`
	Level      types.String `tfsdk:"level"`
}

type predictorCustomMap struct {
	Contains      types.String `tfsdk:"contains"`
	Type          types.String `tfsdk:"type"`
	BetweenRanges types.Object `tfsdk:"between_ranges"`
	IPRanges      types.Object `tfsdk:"ip_ranges"`
	StringList    types.Object `tfsdk:"string_list"`
}

type predictorCustomMapHML struct {
	High   types.Object `tfsdk:"high"`
	Medium types.Object `tfsdk:"medium"`
	Low    types.Object `tfsdk:"low"`
}

type predictorCustomMapHMLBetweenRanges struct {
	MinScore types.Float64 `tfsdk:"min_score"`
	MaxScore types.Float64 `tfsdk:"max_score"`
}

type predictorCustomMapHMLList struct {
	Values types.Set `tfsdk:"values"`
}

type predictorUserLocationAnomalyRadius struct {
	Distance types.Int64  `tfsdk:"distance"`
	Unit     types.String `tfsdk:"unit"`
}

type predictorUserRiskBehaviorPredictionModel struct {
	Name types.String `tfsdk:"name"`
}

type predictorVelocityEvery struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	MinSample types.Int64  `tfsdk:"min_sample"`
}

type predictorVelocityFallback struct {
	Strategy types.String  `tfsdk:"strategy"`
	High     types.Float64 `tfsdk:"high"`
	Medium   types.Float64 `tfsdk:"medium"`
}

type predictorVelocitySlidingWindow struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	MinSample types.Int64  `tfsdk:"min_sample"`
}

type predictorVelocityUse struct {
	UseType types.String  `tfsdk:"type"`
	Medium  types.Float64 `tfsdk:"medium"`
	High    types.Float64 `tfsdk:"high"`
}

var (
	defaultTFObjectTypes = map[string]attr.Type{
		"weight": types.Int64Type,
		"result": types.ObjectType{
			AttrTypes: defaultResultTFObjectTypes,
		},
	}

	defaultResultTFObjectTypes = map[string]attr.Type{
		"type":  types.StringType,
		"level": types.StringType,
	}

	predictorCustomMapTFObjectTypes = map[string]attr.Type{
		"contains": types.StringType,
		"type":     types.StringType,
		"between_ranges": types.ObjectType{
			AttrTypes: predictorCustomMapBetweenHMLTFObjectTypes,
		},
		"ip_ranges": types.ObjectType{
			AttrTypes: predictorCustomMapIPRangesHMLTFObjectTypes,
		},
		"string_list": types.ObjectType{
			AttrTypes: predictorCustomMapStringListHMLTFObjectTypes,
		},
	}

	predictorCustomMapBetweenHMLTFObjectTypes = map[string]attr.Type{
		"high":   types.ObjectType{AttrTypes: predictorCustomMapHMLBetweenRangesTFObjectTypes},
		"medium": types.ObjectType{AttrTypes: predictorCustomMapHMLBetweenRangesTFObjectTypes},
		"low":    types.ObjectType{AttrTypes: predictorCustomMapHMLBetweenRangesTFObjectTypes},
	}

	predictorCustomMapIPRangesHMLTFObjectTypes = map[string]attr.Type{
		"high":   types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
		"medium": types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
		"low":    types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
	}

	predictorCustomMapStringListHMLTFObjectTypes = map[string]attr.Type{
		"high":   types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
		"medium": types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
		"low":    types.ObjectType{AttrTypes: predictorCustomMapHMLListTFObjectTypes},
	}

	predictorCustomMapHMLBetweenRangesTFObjectTypes = map[string]attr.Type{
		"min_score": types.Float64Type,
		"max_score": types.Float64Type,
	}

	predictorCustomMapHMLListTFObjectTypes = map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}

	predictorUserLocationAnomalyRadiusTFObjectTypes = map[string]attr.Type{
		"distance": types.Int64Type,
		"unit":     types.StringType,
	}

	predictorUserRiskBehaviorPredictionModelTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	predictorVelocityEveryTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int64Type,
		"min_sample": types.Int64Type,
	}

	predictorVelocityFallbackTFObjectTypes = map[string]attr.Type{
		"strategy": types.StringType,
		"high":     types.Float64Type,
		"medium":   types.Float64Type,
	}

	predictorVelocitySlidingWindowTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int64Type,
		"min_sample": types.Int64Type,
	}

	predictorVelocityUseTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"medium": types.Float64Type,
		"high":   types.Float64Type,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &RiskPredictorResource{}
	_ resource.ResourceWithConfigure   = &RiskPredictorResource{}
	_ resource.ResourceWithImportState = &RiskPredictorResource{}
)

// New Object
func NewRiskPredictorResource() resource.Resource {
	return &RiskPredictorResource{}
}

// Metadata
func (r *RiskPredictorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_predictor"
}

// Schema
func (r *RiskPredictorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5
	const attrDescriptionMaxLength = 1024

	typeDescriptionFmt := "A string that specifies the type of the risk predictor.  This can be either `ANONYMOUS_NETWORK`, `COMPOSITE`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `DEVICE`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR` or `VELOCITY`."
	typeDescription := framework.SchemaDescription{
		MarkdownDescription: typeDescriptionFmt,
		Description:         strings.ReplaceAll(typeDescriptionFmt, "`", "\""),
	}

	// resultLevelDescriptionFmt := "A string that identifies the risk level. Options are `HIGH`, `MEDIUM`, and `LOW`."
	// resultLevelDescription := framework.SchemaDescription{
	// 	MarkdownDescription: resultLevelDescriptionFmt,
	// 	Description:         strings.ReplaceAll(resultLevelDescriptionFmt, "`", "\""),
	// }

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage risk predictors in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to configure the risk predictor in."},
			),

			"name": schema.StringAttribute{
				Description: "A string that specifies the unique, friendly name for the predictor. This name is displayed in the Risk Policies UI, when the admin is asked to define the overrides and weights in policy configuration.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"compact_name": schema.StringAttribute{
				Description: "A string that specifies the unique name for the predictor for use in risk evaluation request/response payloads. This property is immutable; it cannot be modified after initial creation. The value must be alpha-numeric, with no special characters or spaces. This name is used in the API both for policy configuration, and in the Risk Evaluation response (under details).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "The value must be alpha-numeric, with no special characters or spaces."),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the risk predictor. Maximum length is 1024 characters.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(attrDescriptionMaxLength),
				},
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(func() []string {
						strings := make([]string, 0)
						for _, v := range risk.AllowedEnumPredictorTypeEnumValues {
							strings = append(strings, string(v))
						}
						return strings
					}()...),
				},
			},

			"default": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"weight": schema.Int64Attribute{
						Description: "A number that specifies the default weight for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
						Optional:    true,
						Computed:    true,
					},

					"result": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description:         typeDescription.Description,
								MarkdownDescription: typeDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(func() []string {
										strings := make([]string, 0)
										for _, v := range risk.AllowedEnumResultTypeEnumValues {
											strings = append(strings, string(v))
										}
										return strings
									}()...),
								},
							},

							"level": schema.StringAttribute{
								Description:         typeDescription.Description,
								MarkdownDescription: typeDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(func() []string {
										strings := make([]string, 0)
										for _, v := range risk.AllowedEnumRiskLevelEnumValues {
											strings = append(strings, string(v))
										}
										return strings
									}()...),
								},
							},
						},
					},
				},
			},

			"licensed": schema.BoolAttribute{
				Description: "A boolean that indicates whether PingOne Risk is licensed for the environment.",
				Computed:    true,
			},

			"deletable": schema.BoolAttribute{
				Description: "A boolean that indicates the PingOne Risk predictor can be deleted or not.",
				Computed:    true,
			},

			// Anonymous network, IP reputation, geovelocity
			"allowed_cidr_list": schema.SetAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`), "Values must be valid CIDR format."),
					),
				},
			},

			"custom_map": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"contains": schema.StringAttribute{
						Description: "A string that specifies the value to match in the custom map. Maximum length is 1024 characters.",
						Required:    true,
					},

					"type": schema.StringAttribute{
						Description: typeDescription.Description,
						Computed:    true,
					},

					"between_ranges": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"high": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"min_score": schema.Float64Attribute{
										Description: "A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},

									"max_score": schema.Float64Attribute{
										Description: "A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},
								},
							},

							"medium": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"min_score": schema.Float64Attribute{
										Description: "A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},

									"max_score": schema.Float64Attribute{
										Description: "A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},
								},
							},

							"low": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"min_score": schema.Float64Attribute{
										Description: "A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},

									"max_score": schema.Float64Attribute{
										Description: "A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
										Required:    true,
									},
								},
							},
						},
					},

					"ip_ranges": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"high": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},

							"medium": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},

							"low": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},

					"string_list": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"high": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},

							"medium": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},

							"low": schema.SingleNestedAttribute{
								Description: "",
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},

			// New device
			"detect": schema.StringAttribute{
				Optional: true,
				Computed: true,

				Validators: []validator.String{
					stringvalidator.OneOf(func() []string {
						strings := make([]string, 0)
						for _, v := range risk.AllowedEnumPredictorNewDeviceDetectTypeEnumValues {
							strings = append(strings, string(v))
						}
						return strings
					}()...),
				},
			},

			"activation_at": schema.StringAttribute{
				Description: "You can use the `activation_at` parameter to specify a date on which the learning process for the predictor should be restarted. This can be used in conjunction with the fallback setting (`default.result.level`) to force strong authentication when moving the predictor to production. The date should be in an RFC3339 format. Note that activation date uses UTC time.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.RFC3339Regexp, "Attribute must be a valid RFC3339 date/time string."),
				},
			},

			// User Location Anomaly
			"radius": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"distance": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},

					"unit": schema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,

						Validators: []validator.String{
							stringvalidator.OneOf(func() []string {
								strings := make([]string, 0)
								for _, v := range risk.AllowedEnumDistanceUnitEnumValues {
									strings = append(strings, string(v))
								}
								return strings
							}()...),
						},
					},
				},
			},

			"days": schema.Int64Attribute{
				Description: "",
				Optional:    true,
				Computed:    true,
			},

			// User Risk Behavior
			"prediction_model": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "",
						Required:    true,

						Validators: []validator.String{
							stringvalidator.OneOf(func() []string {
								strings := make([]string, 0)
								for _, v := range risk.AllowedEnumUserRiskBehaviorRiskModelEnumValues {
									strings = append(strings, string(v))
								}
								return strings
							}()...),
						},
					},
				},
			},

			// Velocity
			"measure": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(func() []string {
						strings := make([]string, 0)
						for _, v := range risk.AllowedEnumPredictorVelocityMeasureEnumValues {
							strings = append(strings, string(v))
						}
						return strings
					}()...),
				},
			},

			"of": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("${event.ip}", "${event.user.id}"),
				},
			},

			"by": schema.SetAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("${event.user.id}", "${event.ip}"),
					),
				},
			},

			"use": schema.SingleNestedAttribute{
				Description: "",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of the risk predictor.",
						Computed:    true,
					},

					"medium": schema.Int64Attribute{
						Description: "The medium risk level.",
						Computed:    true,
					},

					"high": schema.Int64Attribute{
						Description: "The high risk level.",
						Computed:    true,
					},
				},
			},

			"fallback": schema.SingleNestedAttribute{
				Description: "An object that contains configuration values for the fallback risk predictor type.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"strategy": schema.StringAttribute{
						Description: "The strategy to use when the risk predictor is not able to determine a risk level.",
						Computed:    true,
					},

					"high": schema.Int64Attribute{
						Description: "The high risk level.",
						Computed:    true,
					},

					"medium": schema.Int64Attribute{
						Description: "The medium risk level.",
						Computed:    true,
					},
				},
			},

			"every": schema.SingleNestedAttribute{
				Description: "An object that contains configuration values for the every risk predictor type.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"unit": schema.StringAttribute{
						Description: "The unit of measurement for the `interval` parameter.",
						Computed:    true,
					},

					"quantity": schema.Int64Attribute{
						Description: "The number of `unit` intervals to use for the risk predictor.",
						Computed:    true,
					},

					"min_sample": schema.Int64Attribute{
						Description: "The minimum number of samples to use for the risk predictor.",
						Computed:    true,
					},
				},
			},

			"sliding_window": schema.SingleNestedAttribute{
				Description: "An object that contains configuration values for the sliding window risk predictor type.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"unit": schema.StringAttribute{
						Description: "The unit of measurement for the `interval` parameter.",
						Computed:    true,
					},

					"quantity": schema.Int64Attribute{
						Description: "The number of `unit` intervals to use for the risk predictor.",
						Computed:    true,
					},

					"min_sample": schema.Int64Attribute{
						Description: "The minimum number of samples to use for the risk predictor.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *RiskPredictorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RiskPredictorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state riskPredictorResourceModel

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
	riskPredictor, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskAdvancedPredictorsApi.CreateRiskPredictor(ctx, plan.EnvironmentId.ValueString()).RiskPredictor(*riskPredictor).Execute()
		},
		"CreateRiskPredictor",
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
	resp.Diagnostics.Append(state.toState(ctx, response.(*risk.RiskPredictor))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPredictorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *riskPredictorResourceModel

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
			return r.client.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneRiskPredictor",
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
	resp.Diagnostics.Append(data.toState(ctx, response.(*risk.RiskPredictor))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RiskPredictorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state riskPredictorResourceModel

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
	riskPredictor, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.RiskAdvancedPredictorsApi.UpdateRiskPredictor(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).RiskPredictor(*riskPredictor).Execute()
		},
		"UpdateRiskPredictor",
		framework.DefaultCustomError,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response.(*risk.RiskPredictor))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPredictorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *riskPredictorResourceModel

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
			r, err := r.client.RiskAdvancedPredictorsApi.DeleteRiskAdvancedPredictor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteRiskAdvancedPredictor",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RiskPredictorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (p *riskPredictorResourceModel) expand(ctx context.Context) (*risk.RiskPredictor, diag.Diagnostics) {
	var diags diag.Diagnostics

	riskPredictor := &risk.RiskPredictor{}
	var d diag.Diagnostics

	if !p.Type.IsNull() {

		data := risk.NewRiskPredictorCommon(p.Name.ValueString(), p.CompactName.ValueString(), risk.EnumPredictorType(p.Type.ValueString()))

		if !p.Description.IsNull() && !p.Description.IsUnknown() {
			data.SetDescription(p.Description.ValueString())
		}

		if !p.Default.IsNull() && !p.Default.IsUnknown() {
			var defaultPlan predictorDefault
			d := p.Default.As(ctx, &defaultPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			dataDefault := risk.NewRiskPredictorCommonDefault(int32(defaultPlan.Weight.ValueInt64()))

			if !defaultPlan.Result.IsNull() && !defaultPlan.Result.IsUnknown() {
				var defaultResultPlan predictorDefaultResult
				d := defaultPlan.Result.As(ctx, &defaultResultPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				dataDefaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(defaultResultPlan.Level.ValueString()))
				dataDefaultResult.SetType(risk.EnumResultType(defaultResultPlan.ResultType.ValueString()))
				dataDefault.SetResult(*dataDefaultResult)

				data.SetDefault(*dataDefault)
			}

			data.SetDefault(*dataDefault)
		}

		switch p.Type.ValueString() {
		case string(risk.ENUMPREDICTORTYPE_ANONYMOUS_NETWORK):
			riskPredictor.RiskPredictorAnonymousNetwork, d = p.expandPredictorAnonymousNetwork(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_COMPOSITE):
			riskPredictor.RiskPredictorComposite, d = p.expandPredictorComposite(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_MAP):
			riskPredictor.RiskPredictorCustom, d = p.expandPredictorCustom(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_GEO_VELOCITY):
			riskPredictor.RiskPredictorGeovelocity, d = p.expandPredictorGeovelocity(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_IP_REPUTATION):
			riskPredictor.RiskPredictorIPReputation, d = p.expandPredictorIPReputation(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_DEVICE):
			riskPredictor.RiskPredictorDevice, d = p.expandPredictorDevice(data)
		case string(risk.ENUMPREDICTORTYPE_USER_RISK_BEHAVIOR):
			riskPredictor.RiskPredictorUserRiskBehavior, d = p.expandPredictorUserRiskBehavior(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_USER_LOCATION_ANOMALY):
			riskPredictor.RiskPredictorUserLocationAnomaly, d = p.expandPredictorUserLocationAnomaly(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_VELOCITY):
			riskPredictor.RiskPredictorVelocity, d = p.expandPredictorVelocity(ctx, data)
		}

		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	return riskPredictor, diags
}

func (p *riskPredictorResourceModel) expandPredictorAnonymousNetwork(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorAnonymousNetwork, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorAnonymousNetwork{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.AllowedCIDRList.IsNull() && !p.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := p.AllowedCIDRList.ToSetValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, allowedCIDRListSet)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorComposite(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorComposite, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorComposite

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorCustom(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorCustom, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorCustom{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.CustomMap.IsNull() && !p.CustomMap.IsUnknown() {

		var plan predictorCustomMap
		d := p.CustomMap.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		var contains string
		if !plan.Contains.IsNull() && !plan.Contains.IsUnknown() {
			contains = plan.Contains.ValueString()
		}

		setHigh := false
		high := risk.RiskPredictorCustomItem{}
		setMedium := false
		medium := risk.RiskPredictorCustomItem{}
		setLow := false
		low := risk.RiskPredictorCustomItem{}

		if !plan.BetweenRanges.IsNull() && !plan.BetweenRanges.IsUnknown() {
			var hmlPlan predictorCustomMapHML
			d := plan.BetweenRanges.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			// High
			if !hmlPlan.High.IsNull() && !hmlPlan.High.IsUnknown() {
				var highHmlPlan predictorCustomMapHMLBetweenRanges
				d := hmlPlan.High.As(ctx, &highHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				v := risk.NewRiskPredictorCustomItemBetween(
					contains,
					*risk.NewRiskPredictorCustomItemBetweenBetween(
						float32(highHmlPlan.MinScore.ValueFloat64()),
						float32(highHmlPlan.MaxScore.ValueFloat64()),
					),
				)

				high.RiskPredictorCustomItemBetween = v
				setHigh = true
			}

			// Medium
			if !hmlPlan.Medium.IsNull() && !hmlPlan.Medium.IsUnknown() {
				var mediumHmlPlan predictorCustomMapHMLBetweenRanges
				d := hmlPlan.Medium.As(ctx, &mediumHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				v := risk.NewRiskPredictorCustomItemBetween(
					contains,
					*risk.NewRiskPredictorCustomItemBetweenBetween(
						float32(mediumHmlPlan.MinScore.ValueFloat64()),
						float32(mediumHmlPlan.MaxScore.ValueFloat64()),
					),
				)

				medium.RiskPredictorCustomItemBetween = v
				setMedium = true
			}

			// Low
			if !hmlPlan.Low.IsNull() && !hmlPlan.Low.IsUnknown() {
				var lowHmlPlan predictorCustomMapHMLBetweenRanges
				d := hmlPlan.Low.As(ctx, &lowHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				v := risk.NewRiskPredictorCustomItemBetween(
					contains,
					*risk.NewRiskPredictorCustomItemBetweenBetween(
						float32(lowHmlPlan.MinScore.ValueFloat64()),
						float32(lowHmlPlan.MaxScore.ValueFloat64()),
					),
				)

				low.RiskPredictorCustomItemBetween = v
				setLow = true
			}
		}

		if !plan.IPRanges.IsNull() && !plan.IPRanges.IsUnknown() {
			var hmlPlan predictorCustomMapHML
			d := plan.IPRanges.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			// High
			if !hmlPlan.High.IsNull() && !hmlPlan.High.IsUnknown() {
				var highHmlPlan predictorCustomMapHMLList
				d := hmlPlan.High.As(ctx, &highHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := highHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemIPRange(
					contains,
					valuesSlice,
				)

				high.RiskPredictorCustomItemIPRange = v
				setHigh = true
			}

			// Medium
			if !hmlPlan.Medium.IsNull() && !hmlPlan.Medium.IsUnknown() {
				var mediumHmlPlan predictorCustomMapHMLList
				d := hmlPlan.Medium.As(ctx, &mediumHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := mediumHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemIPRange(
					contains,
					valuesSlice,
				)

				medium.RiskPredictorCustomItemIPRange = v
				setMedium = true
			}

			// Low
			if !hmlPlan.Low.IsNull() && !hmlPlan.Low.IsUnknown() {
				var lowHmlPlan predictorCustomMapHMLList
				d := hmlPlan.Low.As(ctx, &lowHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := lowHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemIPRange(
					contains,
					valuesSlice,
				)

				low.RiskPredictorCustomItemIPRange = v
				setLow = true
			}
		}

		if !plan.StringList.IsNull() && !plan.StringList.IsUnknown() {
			var hmlPlan predictorCustomMapHML
			d := plan.StringList.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			// High
			if !hmlPlan.High.IsNull() && !hmlPlan.High.IsUnknown() {
				var highHmlPlan predictorCustomMapHMLList
				d := hmlPlan.High.As(ctx, &highHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := highHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemList(
					contains,
					valuesSlice,
				)

				high.RiskPredictorCustomItemList = v
				setHigh = true
			}

			// Medium
			if !hmlPlan.Medium.IsNull() && !hmlPlan.Medium.IsUnknown() {
				var mediumHmlPlan predictorCustomMapHMLList
				d := hmlPlan.Medium.As(ctx, &mediumHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := mediumHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemList(
					contains,
					valuesSlice,
				)

				medium.RiskPredictorCustomItemList = v
				setMedium = true
			}

			// Low
			if !hmlPlan.Low.IsNull() && !hmlPlan.Low.IsUnknown() {
				var lowHmlPlan predictorCustomMapHMLList
				d := hmlPlan.Low.As(ctx, &lowHmlPlan, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    false,
					UnhandledUnknownAsEmpty: false,
				})
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				valuesSlice := make([]string, 0)
				valueSet, d := lowHmlPlan.Values.ToSetValue(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
				pointerSlice := framework.TFSetToStringSlice(ctx, valueSet)

				if len(pointerSlice) > 0 {

					for i := range pointerSlice {
						valuesSlice = append(valuesSlice, *pointerSlice[i])
					}
				}

				v := risk.NewRiskPredictorCustomItemList(
					contains,
					valuesSlice,
				)

				low.RiskPredictorCustomItemList = v
				setLow = true
			}
		}

		customMap := risk.NewRiskPredictorCustomAllOfMap()

		if setHigh {
			customMap.SetHigh(high)
		}
		if setMedium {
			customMap.SetMedium(medium)
		}
		if setLow {
			customMap.SetLow(low)
		}

		data.SetMap(*customMap)
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorGeovelocity(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorGeovelocity, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorGeovelocity{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.AllowedCIDRList.IsNull() && !p.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := p.AllowedCIDRList.ToSetValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, allowedCIDRListSet)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorIPReputation(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorIPReputation, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorIPReputation{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.AllowedCIDRList.IsNull() && !p.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := p.AllowedCIDRList.ToSetValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, allowedCIDRListSet)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorDevice(riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorDevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorDevice{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if p.Detect.IsNull() || p.Detect.IsUnknown() {
		p.Detect = types.StringValue(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE))
	}

	data.SetDetect(risk.EnumPredictorNewDeviceDetectType(p.Detect.ValueString()))

	if !p.ActivationAt.IsNull() && !p.ActivationAt.IsUnknown() {
		t, e := time.Parse(time.RFC3339, p.ActivationAt.ValueString())
		if e != nil {
			diags.AddError(
				"Invalid data format",
				"Cannot convert activation_at to a date/time.  Please check the format is a valid RFC3339 date time format.")
			return nil, diags
		}

		data.SetActivationAt(t)
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUserLocationAnomaly(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorUserLocationAnomaly, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorUserLocationAnomaly{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.Radius.IsNull() && !p.Radius.IsUnknown() {
		var radiusPlan predictorUserLocationAnomalyRadius
		d := p.Radius.As(ctx, &radiusPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		radiusPlanUnit := risk.ENUMDISTANCEUNIT_KILOMETERS

		if !radiusPlan.Unit.IsNull() && !radiusPlan.Unit.IsUnknown() {
			radiusPlanUnit = risk.EnumDistanceUnit(radiusPlan.Unit.ValueString())
		}

		radius := risk.NewRiskPredictorUserLocationAnomalyAllOfRadius(int32(radiusPlan.Distance.ValueInt64()), radiusPlanUnit)

		data.SetRadius(*radius)
	}

	data.SetDays(50)

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUserRiskBehavior(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorUserRiskBehavior, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorUserRiskBehavior{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	if !p.PredictionModel.IsNull() && !p.PredictionModel.IsUnknown() {
		var plan predictorUserRiskBehaviorPredictionModel
		d := p.PredictionModel.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		predictionModel := risk.NewRiskPredictorUserRiskBehaviorAllOfPredictionModel(risk.EnumUserRiskBehaviorRiskModel(plan.Name.ValueString()))

		data.SetPredictionModel(*predictionModel)
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorVelocity(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorVelocity, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorVelocity{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
	}

	// Of
	if !p.Of.IsNull() && !p.Of.IsUnknown() {
		data.SetOf(p.Of.ValueString())
	}

	// By
	if !p.By.IsNull() && !p.By.IsUnknown() {
		bySet, d := p.By.ToSetValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		byPointerSlice := framework.TFSetToStringSlice(ctx, bySet)

		if len(byPointerSlice) > 0 {
			bySlice := make([]string, 0)
			for i := range byPointerSlice {
				bySlice = append(bySlice, *byPointerSlice[i])
			}
			data.SetBy(bySlice)
		}
	} else {
		if p.Of.Equal(types.StringValue("${event.ip}")) {
			data.SetBy([]string{"${event.user.id}"})
		}

		if p.Of.Equal(types.StringValue("${event.user.id}")) {
			data.SetBy([]string{"${event.ip}"})
		}
	}

	// Every
	if !p.Every.IsNull() && !p.Every.IsUnknown() {
		var plan predictorVelocityEvery
		d := p.Every.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		every := risk.NewRiskPredictorVelocityAllOfEvery()

		if !plan.Unit.IsNull() && !plan.Unit.IsUnknown() {
			every.SetUnit(risk.EnumPredictorUnit(plan.Unit.ValueString()))
		}

		if !plan.Quantity.IsNull() && !plan.Quantity.IsUnknown() {
			every.SetQuantity(int32(plan.Quantity.ValueInt64()))
		}

		if !plan.MinSample.IsNull() && !plan.MinSample.IsUnknown() {
			every.SetMinSample(int32(plan.MinSample.ValueInt64()))
		}

		data.SetEvery(*every)
	} else {
		every := risk.NewRiskPredictorVelocityAllOfEvery()
		every.SetUnit(risk.ENUMPREDICTORUNIT_HOUR)
		every.SetQuantity(int32(1))
		every.SetMinSample(int32(5))
		data.SetEvery(*every)
	}

	// Fallback
	if !p.Fallback.IsNull() && !p.Fallback.IsUnknown() {
		var plan predictorVelocityFallback
		d := p.Fallback.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		fallback := risk.NewRiskPredictorVelocityAllOfFallback()

		if !plan.Strategy.IsNull() && !plan.Strategy.IsUnknown() {
			fallback.SetStrategy(risk.EnumPredictorVelocityFallbackStrategy(plan.Strategy.ValueString()))
		}

		if !plan.High.IsNull() && !plan.High.IsUnknown() {
			fallback.SetHigh(float32(plan.High.ValueFloat64()))
		}

		if !plan.Medium.IsNull() && !plan.Medium.IsUnknown() {
			fallback.SetMedium(float32(plan.Medium.ValueFloat64()))
		}

		data.SetFallback(*fallback)
	} else {
		fallback := risk.NewRiskPredictorVelocityAllOfFallback()
		fallback.SetStrategy(risk.ENUMPREDICTORVELOCITYFALLBACKSTRATEGY_ENVIRONMENT_MAX)

		if p.Of.Equal(types.StringValue("${event.ip}")) {
			fallback.SetHigh(float32(30))
			fallback.SetMedium(float32(20))
		}

		if p.Of.Equal(types.StringValue("${event.user.id}")) {
			fallback.SetHigh(float32(3500))
			fallback.SetMedium(float32(2500))
		}

		data.SetFallback(*fallback)
	}

	// Measure
	if !p.Measure.IsNull() && !p.Measure.IsUnknown() {
		data.SetMeasure(risk.EnumPredictorVelocityMeasure(p.Measure.ValueString()))
	} else {
		data.SetMeasure(risk.ENUMPREDICTORVELOCITYMEASURE_DISTINCT_COUNT)
	}

	// SlidingWindow
	if !p.SlidingWindow.IsNull() && !p.SlidingWindow.IsUnknown() {
		var plan predictorVelocitySlidingWindow
		d := p.SlidingWindow.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		slidingWindow := risk.NewRiskPredictorVelocityAllOfSlidingWindow()

		if !plan.Unit.IsNull() && !plan.Unit.IsUnknown() {
			slidingWindow.SetUnit(risk.EnumPredictorUnit(plan.Unit.ValueString()))
		}

		if !plan.Quantity.IsNull() && !plan.Quantity.IsUnknown() {
			slidingWindow.SetQuantity(int32(plan.Quantity.ValueInt64()))
		}

		if !plan.MinSample.IsNull() && !plan.MinSample.IsUnknown() {
			slidingWindow.SetMinSample(int32(plan.MinSample.ValueInt64()))
		}

		data.SetSlidingWindow(*slidingWindow)
	} else {
		slidingWindow := risk.NewRiskPredictorVelocityAllOfSlidingWindow()
		slidingWindow.SetUnit(risk.ENUMPREDICTORUNIT_DAY)
		slidingWindow.SetQuantity(int32(7))
		slidingWindow.SetMinSample(int32(3))
		data.SetSlidingWindow(*slidingWindow)
	}

	// Use
	if !p.Use.IsNull() && !p.Use.IsUnknown() {
		var plan predictorVelocityUse
		d := p.Use.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		use := risk.NewRiskPredictorVelocityAllOfUse()

		if !plan.UseType.IsNull() && !plan.UseType.IsUnknown() {
			use.SetType(risk.EnumPredictorVelocityUseType(plan.UseType.ValueString()))
		}

		if !plan.Medium.IsNull() && !plan.Medium.IsUnknown() {
			use.SetMedium(float32(plan.Medium.ValueFloat64()))
		}

		if !plan.High.IsNull() && !plan.High.IsUnknown() {
			use.SetHigh(float32(plan.High.ValueFloat64()))
		}

		data.SetUse(*use)
	} else {
		use := risk.NewRiskPredictorVelocityAllOfUse()
		use.SetType(risk.ENUMPREDICTORVELOCITYUSETYPE_POISSON_WITH_MAX)
		use.SetMedium(float32(2.0))
		use.SetHigh(float32(4.0))
		data.SetUse(*use)
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) toState(ctx context.Context, apiObject *risk.RiskPredictor) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := risk.RiskPredictorCommon{}

	if apiObject.RiskPredictorAnonymousNetwork != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorAnonymousNetwork.Id,
			Name:        apiObject.RiskPredictorAnonymousNetwork.Name,
			CompactName: apiObject.RiskPredictorAnonymousNetwork.CompactName,
			Description: apiObject.RiskPredictorAnonymousNetwork.Description,
			Type:        apiObject.RiskPredictorAnonymousNetwork.Type,
			Default:     apiObject.RiskPredictorAnonymousNetwork.Default,
			Licensed:    apiObject.RiskPredictorAnonymousNetwork.Licensed,
			Deletable:   apiObject.RiskPredictorAnonymousNetwork.Deletable,
		}
	}

	if apiObject.RiskPredictorComposite != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorComposite.Id,
			Name:        apiObject.RiskPredictorComposite.Name,
			CompactName: apiObject.RiskPredictorComposite.CompactName,
			Description: apiObject.RiskPredictorComposite.Description,
			Type:        apiObject.RiskPredictorComposite.Type,
			Default:     apiObject.RiskPredictorComposite.Default,
			Licensed:    apiObject.RiskPredictorComposite.Licensed,
			Deletable:   apiObject.RiskPredictorComposite.Deletable,
		}
	}

	if apiObject.RiskPredictorCustom != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorCustom.Id,
			Name:        apiObject.RiskPredictorCustom.Name,
			CompactName: apiObject.RiskPredictorCustom.CompactName,
			Description: apiObject.RiskPredictorCustom.Description,
			Type:        apiObject.RiskPredictorCustom.Type,
			Default:     apiObject.RiskPredictorCustom.Default,
			Licensed:    apiObject.RiskPredictorCustom.Licensed,
			Deletable:   apiObject.RiskPredictorCustom.Deletable,
		}
	}

	if apiObject.RiskPredictorGeovelocity != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorGeovelocity.Id,
			Name:        apiObject.RiskPredictorGeovelocity.Name,
			CompactName: apiObject.RiskPredictorGeovelocity.CompactName,
			Description: apiObject.RiskPredictorGeovelocity.Description,
			Type:        apiObject.RiskPredictorGeovelocity.Type,
			Default:     apiObject.RiskPredictorGeovelocity.Default,
			Licensed:    apiObject.RiskPredictorGeovelocity.Licensed,
			Deletable:   apiObject.RiskPredictorGeovelocity.Deletable,
		}
	}

	if apiObject.RiskPredictorIPReputation != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorIPReputation.Id,
			Name:        apiObject.RiskPredictorIPReputation.Name,
			CompactName: apiObject.RiskPredictorIPReputation.CompactName,
			Description: apiObject.RiskPredictorIPReputation.Description,
			Type:        apiObject.RiskPredictorIPReputation.Type,
			Default:     apiObject.RiskPredictorIPReputation.Default,
			Licensed:    apiObject.RiskPredictorIPReputation.Licensed,
			Deletable:   apiObject.RiskPredictorIPReputation.Deletable,
		}
	}

	if apiObject.RiskPredictorDevice != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorDevice.Id,
			Name:        apiObject.RiskPredictorDevice.Name,
			CompactName: apiObject.RiskPredictorDevice.CompactName,
			Description: apiObject.RiskPredictorDevice.Description,
			Type:        apiObject.RiskPredictorDevice.Type,
			Default:     apiObject.RiskPredictorDevice.Default,
			Licensed:    apiObject.RiskPredictorDevice.Licensed,
			Deletable:   apiObject.RiskPredictorDevice.Deletable,
		}
	}

	if apiObject.RiskPredictorUserRiskBehavior != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorUserRiskBehavior.Id,
			Name:        apiObject.RiskPredictorUserRiskBehavior.Name,
			CompactName: apiObject.RiskPredictorUserRiskBehavior.CompactName,
			Description: apiObject.RiskPredictorUserRiskBehavior.Description,
			Type:        apiObject.RiskPredictorUserRiskBehavior.Type,
			Default:     apiObject.RiskPredictorUserRiskBehavior.Default,
			Licensed:    apiObject.RiskPredictorUserRiskBehavior.Licensed,
			Deletable:   apiObject.RiskPredictorUserRiskBehavior.Deletable,
		}
	}

	if apiObject.RiskPredictorUserLocationAnomaly != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorUserLocationAnomaly.Id,
			Name:        apiObject.RiskPredictorUserLocationAnomaly.Name,
			CompactName: apiObject.RiskPredictorUserLocationAnomaly.CompactName,
			Description: apiObject.RiskPredictorUserLocationAnomaly.Description,
			Type:        apiObject.RiskPredictorUserLocationAnomaly.Type,
			Default:     apiObject.RiskPredictorUserLocationAnomaly.Default,
			Licensed:    apiObject.RiskPredictorUserLocationAnomaly.Licensed,
			Deletable:   apiObject.RiskPredictorUserLocationAnomaly.Deletable,
		}
	}

	if apiObject.RiskPredictorVelocity != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorVelocity.Id,
			Name:        apiObject.RiskPredictorVelocity.Name,
			CompactName: apiObject.RiskPredictorVelocity.CompactName,
			Description: apiObject.RiskPredictorVelocity.Description,
			Type:        apiObject.RiskPredictorVelocity.Type,
			Default:     apiObject.RiskPredictorVelocity.Default,
			Licensed:    apiObject.RiskPredictorVelocity.Licensed,
			Deletable:   apiObject.RiskPredictorVelocity.Deletable,
		}
	}

	p.Id = framework.StringToTF(apiObjectCommon.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObjectCommon.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObjectCommon.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObjectCommon.GetTypeOk())
	p.Licensed = framework.BoolOkToTF(apiObjectCommon.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObjectCommon.GetDeletableOk())

	// Default block
	p.Default = types.ObjectNull(defaultTFObjectTypes)
	if v, ok := apiObjectCommon.GetDefaultOk(); ok {
		var d diag.Diagnostics

		defaultResultObj := types.ObjectNull(defaultResultTFObjectTypes)
		if v1, ok := v.GetResultOk(); ok {
			o := map[string]attr.Value{
				"type":  enumRiskPredictorResultTypeOkToTF(v1.GetTypeOk()),
				"level": enumRiskPredictorDefaultResultLevelOkToTF(v1.GetLevelOk()),
			}

			defaultResultObj, d = types.ObjectValue(defaultResultTFObjectTypes, o)
			diags.Append(d...)
		}

		o := map[string]attr.Value{
			"weight": framework.Int32OkToTF(v.GetWeightOk()),
			"result": defaultResultObj,
		}

		objValue, d := types.ObjectValue(defaultTFObjectTypes, o)
		diags.Append(d...)

		p.Default = objValue
	}

	p.Licensed = framework.BoolOkToTF(apiObjectCommon.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObjectCommon.GetDeletableOk())

	// Set all the predictor specific fields all to null before we overwrite them with a value
	p.AllowedCIDRList = types.SetNull(types.StringType)
	p.ActivationAt = types.StringNull()
	p.Detect = types.StringNull()
	p.Radius = types.ObjectNull(predictorUserLocationAnomalyRadiusTFObjectTypes)
	p.Days = types.Int64Null()
	p.CustomMap = types.ObjectNull(predictorCustomMapTFObjectTypes)
	p.PredictionModel = types.ObjectNull(predictorUserRiskBehaviorPredictionModelTFObjectTypes)
	p.By = types.SetNull(types.StringType)
	p.Every = types.ObjectNull(predictorVelocityEveryTFObjectTypes)
	p.Fallback = types.ObjectNull(predictorVelocityFallbackTFObjectTypes)
	p.Measure = types.StringNull()
	p.Of = types.StringNull()
	p.SlidingWindow = types.ObjectNull(predictorVelocitySlidingWindowTFObjectTypes)
	p.Use = types.ObjectNull(predictorVelocityUseTFObjectTypes)

	// Set the predictor specific fields by object type
	if apiObject.RiskPredictorAnonymousNetwork != nil && apiObject.RiskPredictorAnonymousNetwork.GetId() != "" {
		diags.Append(p.toStateRiskPredictorAnonymousNetwork(apiObject.RiskPredictorAnonymousNetwork)...)
	}

	if apiObject.RiskPredictorComposite != nil && apiObject.RiskPredictorComposite.GetId() != "" {
		diags.Append(p.toStateRiskPredictorComposite(apiObject.RiskPredictorComposite)...)
	}

	if apiObject.RiskPredictorCustom != nil && apiObject.RiskPredictorCustom.GetId() != "" {
		diags.Append(p.toStateRiskPredictorCustom(apiObject.RiskPredictorCustom)...)
	}

	if apiObject.RiskPredictorGeovelocity != nil && apiObject.RiskPredictorGeovelocity.GetId() != "" {
		diags.Append(p.toStateRiskPredictorGeovelocity(apiObject.RiskPredictorGeovelocity)...)
	}

	if apiObject.RiskPredictorIPReputation != nil && apiObject.RiskPredictorIPReputation.GetId() != "" {
		diags.Append(p.toStateRiskPredictorIPReputation(apiObject.RiskPredictorIPReputation)...)
	}

	if apiObject.RiskPredictorDevice != nil && apiObject.RiskPredictorDevice.GetId() != "" {
		diags.Append(p.toStateRiskPredictorDevice(apiObject.RiskPredictorDevice)...)
	}

	if apiObject.RiskPredictorUserRiskBehavior != nil && apiObject.RiskPredictorUserRiskBehavior.GetId() != "" {
		diags.Append(p.toStateRiskPredictorUserRiskBehavior(apiObject.RiskPredictorUserRiskBehavior)...)
	}

	if apiObject.RiskPredictorUserLocationAnomaly != nil && apiObject.RiskPredictorUserLocationAnomaly.GetId() != "" {
		diags.Append(p.toStateRiskPredictorUserLocationAnomaly(apiObject.RiskPredictorUserLocationAnomaly)...)
	}

	if apiObject.RiskPredictorVelocity != nil && apiObject.RiskPredictorVelocity.GetId() != "" {
		diags.Append(p.toStateRiskPredictorVelocity(apiObject.RiskPredictorVelocity)...)
	}

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorAnonymousNetwork(apiObject *risk.RiskPredictorAnonymousNetwork) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.AllowedCIDRList = framework.StringSetOkToTF(apiObject.GetWhiteListOk())

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorComposite(apiObject *risk.RiskPredictorComposite) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorCustom(apiObject *risk.RiskPredictorCustom) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.CustomMap = types.ObjectNull(predictorCustomMapTFObjectTypes)

	if v, ok := apiObject.GetMapOk(); ok {
		var d diag.Diagnostics

		// Set all to null before we overwrite them with a value
		betweenRangesObjValue := types.ObjectNull(predictorCustomMapBetweenHMLTFObjectTypes)
		ipRangesObjValue := types.ObjectNull(predictorCustomMapIPRangesHMLTFObjectTypes)
		stringListObjValue := types.ObjectNull(predictorCustomMapStringListHMLTFObjectTypes)

		o := map[string]attr.Value{
			"contains":       types.StringNull(),
			"type":           types.StringNull(),
			"between_ranges": betweenRangesObjValue,
			"ip_ranges":      ipRangesObjValue,
			"string_list":    stringListObjValue,
		}

		setBetweenRanges := false
		betweenObj := map[string]attr.Value{
			"high":   types.ObjectNull(predictorCustomMapHMLBetweenRangesTFObjectTypes),
			"medium": types.ObjectNull(predictorCustomMapHMLBetweenRangesTFObjectTypes),
			"low":    types.ObjectNull(predictorCustomMapHMLBetweenRangesTFObjectTypes),
		}

		setIpRanges := false
		ipRangesObj := map[string]attr.Value{
			"high":   types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
			"medium": types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
			"low":    types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
		}

		setStringList := false
		stringListObj := map[string]attr.Value{
			"high":   types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
			"medium": types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
			"low":    types.ObjectNull(predictorCustomMapHMLListTFObjectTypes),
		}

		if high, ok := v.GetHighOk(); ok {
			// Between
			if v1 := high.RiskPredictorCustomItemBetween; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_score": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_score": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLBetweenRangesTFObjectTypes, levelObj)
				diags.Append(d...)

				betweenObj["high"] = levelObjValue
			}

			// IP Range
			if v1 := high.RiskPredictorCustomItemIPRange; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// IP Ranges
				setIpRanges = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetIpRangeOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				ipRangesObj["high"] = levelObjValue
			}

			// String list
			if v1 := high.RiskPredictorCustomItemList; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// String list
				setStringList = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetListOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				stringListObj["high"] = levelObjValue
			}
		}

		if medium, ok := v.GetMediumOk(); ok {
			// Between
			if v1 := medium.RiskPredictorCustomItemBetween; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_score": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_score": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLBetweenRangesTFObjectTypes, levelObj)
				diags.Append(d...)

				betweenObj["medium"] = levelObjValue
			}

			// IP Range
			if v1 := medium.RiskPredictorCustomItemIPRange; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// IP Ranges
				setIpRanges = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetIpRangeOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				ipRangesObj["medium"] = levelObjValue
			}

			// String list
			if v1 := medium.RiskPredictorCustomItemList; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// String list
				setStringList = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetListOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				stringListObj["medium"] = levelObjValue
			}
		}

		if low, ok := v.GetLowOk(); ok {
			// Between
			if v1 := low.RiskPredictorCustomItemBetween; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_score": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_score": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLBetweenRangesTFObjectTypes, levelObj)
				diags.Append(d...)

				betweenObj["low"] = levelObjValue
			}

			// IP Range
			if v1 := low.RiskPredictorCustomItemIPRange; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// IP Ranges
				setIpRanges = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetIpRangeOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				ipRangesObj["low"] = levelObjValue
			}

			// String list
			if v1 := low.RiskPredictorCustomItemList; v1 != nil {
				o["type"] = framework.StringOkToTF(v1.GetTypeOk())

				// Contains
				contains := framework.StringOkToTF(v1.GetContainsOk())

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return diags
				}

				o["contains"] = contains

				// String list
				setStringList = true

				levelObj := map[string]attr.Value{
					"values": framework.StringSetOkToTF(v1.GetListOk()),
				}
				levelObjValue, d := types.ObjectValue(predictorCustomMapHMLListTFObjectTypes, levelObj)
				diags.Append(d...)

				stringListObj["low"] = levelObjValue
			}
		}

		if setBetweenRanges {
			betweenRangesObjValue, d := types.ObjectValue(predictorCustomMapBetweenHMLTFObjectTypes, betweenObj)
			diags.Append(d...)
			o["between_ranges"] = betweenRangesObjValue
		}

		if setIpRanges {
			ipRangesObjValue, d := types.ObjectValue(predictorCustomMapIPRangesHMLTFObjectTypes, ipRangesObj)
			diags.Append(d...)
			o["ip_ranges"] = ipRangesObjValue
		}

		if setStringList {
			stringListObjValue, d := types.ObjectValue(predictorCustomMapStringListHMLTFObjectTypes, stringListObj)
			diags.Append(d...)
			o["string_list"] = stringListObjValue
		}

		objValue, d := types.ObjectValue(predictorCustomMapTFObjectTypes, o)
		diags.Append(d...)

		p.CustomMap = objValue
	}

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorGeovelocity(apiObject *risk.RiskPredictorGeovelocity) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.AllowedCIDRList = framework.StringSetOkToTF(apiObject.GetWhiteListOk())

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorIPReputation(apiObject *risk.RiskPredictorIPReputation) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.AllowedCIDRList = framework.StringSetOkToTF(apiObject.GetWhiteListOk())

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorDevice(apiObject *risk.RiskPredictorDevice) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.ActivationAt = framework.TimeOkToTF(apiObject.GetActivationAtOk())
	p.Detect = enumRiskPredictorNewDeviceDetectOkToTF(apiObject.GetDetectOk())

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorUserRiskBehavior(apiObject *risk.RiskPredictorUserRiskBehavior) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.PredictionModel = types.ObjectNull(predictorUserRiskBehaviorPredictionModelTFObjectTypes)

	if v, ok := apiObject.GetPredictionModelOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"name": enumRiskPredictorUserRiskBehaviorRiskModelOkToTF(v.GetNameOk()),
		}

		objValue, d := types.ObjectValue(predictorUserRiskBehaviorPredictionModelTFObjectTypes, o)
		diags.Append(d...)

		p.PredictionModel = objValue
	}

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorUserLocationAnomaly(apiObject *risk.RiskPredictorUserLocationAnomaly) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Radius = types.ObjectNull(predictorUserLocationAnomalyRadiusTFObjectTypes)

	if v, ok := apiObject.GetRadiusOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"distance": framework.Int32OkToTF(v.GetDistanceOk()),
			"unit":     enumRiskPredictorDistanceUnitOkToTF(v.GetUnitOk()),
		}

		objValue, d := types.ObjectValue(predictorUserLocationAnomalyRadiusTFObjectTypes, o)
		diags.Append(d...)

		p.Radius = objValue
	}

	p.Days = framework.Int32OkToTF(apiObject.GetDaysOk())

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorVelocity(apiObject *risk.RiskPredictorVelocity) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.By = framework.StringSetOkToTF(apiObject.GetByOk())

	// Every
	p.Every = types.ObjectNull(predictorVelocityEveryTFObjectTypes)

	if v, ok := apiObject.GetEveryOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"unit":       enumRiskPredictorUnitOkToTF(v.GetUnitOk()),
			"quantity":   framework.Int32OkToTF(v.GetQuantityOk()),
			"min_sample": framework.Int32OkToTF(v.GetMinSampleOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityEveryTFObjectTypes, o)
		diags.Append(d...)

		p.Every = objValue
	}

	// Fallback
	p.Fallback = types.ObjectNull(predictorVelocityFallbackTFObjectTypes)

	if v, ok := apiObject.GetFallbackOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"strategy": enumRiskPredictorVelocityFallbackStrategyOkToTF(v.GetStrategyOk()),
			"high":     framework.Float32OkToTF(v.GetHighOk()),
			"medium":   framework.Float32OkToTF(v.GetMediumOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityFallbackTFObjectTypes, o)
		diags.Append(d...)

		p.Fallback = objValue
	}

	p.Measure = enumRiskPredictorVelocityMeasureOkToTF(apiObject.GetMeasureOk())
	p.Of = framework.StringOkToTF(apiObject.GetOfOk())

	// SlidingWindow
	p.SlidingWindow = types.ObjectNull(predictorVelocitySlidingWindowTFObjectTypes)

	if v, ok := apiObject.GetSlidingWindowOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"unit":       enumRiskPredictorUnitOkToTF(v.GetUnitOk()),
			"quantity":   framework.Int32OkToTF(v.GetQuantityOk()),
			"min_sample": framework.Int32OkToTF(v.GetMinSampleOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocitySlidingWindowTFObjectTypes, o)
		diags.Append(d...)

		p.SlidingWindow = objValue
	}

	// Use
	p.Use = types.ObjectNull(predictorVelocityUseTFObjectTypes)

	if v, ok := apiObject.GetUseOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"type":   enumRiskPredictorVelocityUseTypeOkToTF(v.GetTypeOk()),
			"medium": framework.Float32OkToTF(v.GetMediumOk()),
			"high":   framework.Float32OkToTF(v.GetHighOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityUseTFObjectTypes, o)
		diags.Append(d...)

		p.Use = objValue
	}

	return diags
}

func enumRiskPredictorResultTypeOkToTF(v *risk.EnumResultType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorTypeOkToTF(v *risk.EnumPredictorType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorUnitOkToTF(v *risk.EnumPredictorUnit, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorDefaultResultLevelOkToTF(v *risk.EnumRiskLevel, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorNewDeviceDetectOkToTF(v *risk.EnumPredictorNewDeviceDetectType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorDistanceUnitOkToTF(v *risk.EnumDistanceUnit, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorUserRiskBehaviorRiskModelOkToTF(v *risk.EnumUserRiskBehaviorRiskModel, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorVelocityFallbackStrategyOkToTF(v *risk.EnumPredictorVelocityFallbackStrategy, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorVelocityMeasureOkToTF(v *risk.EnumPredictorVelocityMeasure, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}

func enumRiskPredictorVelocityUseTypeOkToTF(v *risk.EnumPredictorVelocityUseType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		if sv := string(*v); sv == "" {
			return types.StringNull()
		} else {
			return types.StringValue(sv)
		}
	}
}
