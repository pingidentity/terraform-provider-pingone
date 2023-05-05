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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	// New device
	ActivationAt types.String `tfsdk:"activation_at"`
	Detect       types.String `tfsdk:"detect"`
	// User Location Anomaly
	Radius types.Object `tfsdk:"radius"`
	Days   types.Int64  `tfsdk:"days"`
	// Velocity
	By            types.Set    `tfsdk:"by"`
	Every         types.Object `tfsdk:"every"`
	Fallback      types.Object `tfsdk:"fallback"`
	MaxDelay      types.Object `tfsdk:"max_delay"`
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

type predictorUserLocationAnomalyRadius struct {
	Distance types.Int64  `tfsdk:"distance"`
	Unit     types.String `tfsdk:"unit"`
}

type predictorVelocityEvery struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	MinSample types.Int64  `tfsdk:"min_sample"`
}

type predictorVelocityFallback struct {
	Strategy types.String `tfsdk:"strategy"`
	High     types.Int64  `tfsdk:"high"`
	Medium   types.Int64  `tfsdk:"medium"`
}

type predictorVelocityMaxDelay struct {
	Unit     types.String `tfsdk:"unit"`
	Quantity types.Int64  `tfsdk:"quantity"`
}

type predictorVelocitySlidingWindow struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	MinSample types.Int64  `tfsdk:"min_sample"`
}

type predictorVelocityUse struct {
	UseType types.String `tfsdk:"type"`
	Medium  types.Int64  `tfsdk:"medium"`
	High    types.Int64  `tfsdk:"high"`
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

	predictorUserLocationAnomalyRadiusTFObjectTypes = map[string]attr.Type{
		"distance": types.Int64Type,
		"unit":     types.StringType,
	}

	predictorVelocityEveryTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int64Type,
		"min_sample": types.Int64Type,
	}

	predictorVelocityFallbackTFObjectTypes = map[string]attr.Type{
		"strategy": types.StringType,
		"high":     types.Int64Type,
		"medium":   types.Int64Type,
	}

	predictorVelocityMaxDelayTFObjectTypes = map[string]attr.Type{
		"unit":     types.StringType,
		"quantity": types.Int64Type,
	}

	predictorVelocitySlidingWindowTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int64Type,
		"min_sample": types.Int64Type,
	}

	predictorVelocityUseTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"medium": types.Int64Type,
		"high":   types.Int64Type,
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

			// New device
			"detect": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)),

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
						Required:    true,

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

			// Velocity
			"measure": schema.StringAttribute{
				Optional: true,
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

			"max_delay": schema.SingleNestedAttribute{
				Description: "An object that contains configuration values for the max delay risk predictor type.",
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
		case string(risk.ENUMPREDICTORTYPE_GEO_VELOCITY):
			riskPredictor.RiskPredictorGeovelocity, d = p.expandPredictorGeovelocity(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_IP_REPUTATION):
			riskPredictor.RiskPredictorIPReputation, d = p.expandPredictorIPReputation(ctx, data)
		case string(risk.ENUMPREDICTORTYPE_DEVICE):
			riskPredictor.RiskPredictorNewDevice, d = p.expandPredictorNewDevice(data)
		case string(risk.ENUMPREDICTORTYPE_USER_LOCATION_ANOMALY):
			riskPredictor.RiskPredictorUserLocationAnomaly, d = p.expandPredictorUserLocationAnomaly(ctx, data)
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

func (p *riskPredictorResourceModel) expandPredictorComposite(ctx context.Context) (*risk.RiskPredictorComposite, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorComposite

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorCustom(ctx context.Context) (*risk.RiskPredictorCustom, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorCustom

	return data, diags
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

func (p *riskPredictorResourceModel) expandPredictorNewDevice(riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorNewDevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorNewDevice{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        riskPredictorCommon.Type,
		Default:     riskPredictorCommon.Default,
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

		radius := risk.NewRiskPredictorUserLocationAnomalyAllOfRadius(int32(radiusPlan.Distance.ValueInt64()), risk.EnumDistanceUnit(radiusPlan.Unit.ValueString()))

		data.SetRadius(*radius)
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUEBA(ctx context.Context) (*risk.RiskPredictorUEBA, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorUEBA

	return data, diags
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
			fallback.SetHigh(int32(plan.High.ValueInt64()))
		}

		if !plan.Medium.IsNull() && !plan.Medium.IsUnknown() {
			fallback.SetMedium(int32(plan.Medium.ValueInt64()))
		}

		data.SetFallback(*fallback)
	}

	// MaxDelay
	if !p.MaxDelay.IsNull() && !p.MaxDelay.IsUnknown() {
		var plan predictorVelocityMaxDelay
		d := p.MaxDelay.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		maxDelay := risk.NewRiskPredictorVelocityAllOfMaxDelay()

		if !plan.Unit.IsNull() && !plan.Unit.IsUnknown() {
			maxDelay.SetUnit(risk.EnumPredictorUnit(plan.Unit.ValueString()))
		}

		if !plan.Quantity.IsNull() && !plan.Quantity.IsUnknown() {
			maxDelay.SetQuantity(int32(plan.Quantity.ValueInt64()))
		}

		data.SetMaxDelay(*maxDelay)
	}

	// Measure
	if !p.Measure.IsNull() && !p.Measure.IsUnknown() {
		data.SetMeasure(risk.EnumPredictorVelocityMeasure(p.Measure.ValueString()))
	}

	// Of
	if !p.Of.IsNull() && !p.Of.IsUnknown() {
		data.SetOf(p.Of.ValueString())
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
			use.SetMedium(int32(plan.Medium.ValueInt64()))
		}

		if !plan.High.IsNull() && !plan.High.IsUnknown() {
			use.SetHigh(int32(plan.High.ValueInt64()))
		}

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

	if apiObject.RiskPredictorNewDevice != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorNewDevice.Id,
			Name:        apiObject.RiskPredictorNewDevice.Name,
			CompactName: apiObject.RiskPredictorNewDevice.CompactName,
			Description: apiObject.RiskPredictorNewDevice.Description,
			Type:        apiObject.RiskPredictorNewDevice.Type,
			Default:     apiObject.RiskPredictorNewDevice.Default,
			Licensed:    apiObject.RiskPredictorNewDevice.Licensed,
			Deletable:   apiObject.RiskPredictorNewDevice.Deletable,
		}
	}

	if apiObject.RiskPredictorUEBA != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorUEBA.Id,
			Name:        apiObject.RiskPredictorUEBA.Name,
			CompactName: apiObject.RiskPredictorUEBA.CompactName,
			Description: apiObject.RiskPredictorUEBA.Description,
			Type:        apiObject.RiskPredictorUEBA.Type,
			Default:     apiObject.RiskPredictorUEBA.Default,
			Licensed:    apiObject.RiskPredictorUEBA.Licensed,
			Deletable:   apiObject.RiskPredictorUEBA.Deletable,
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
	p.By = types.SetNull(types.StringType)
	p.Every = types.ObjectNull(predictorVelocityEveryTFObjectTypes)
	p.Fallback = types.ObjectNull(predictorVelocityFallbackTFObjectTypes)
	p.MaxDelay = types.ObjectNull(predictorVelocityMaxDelayTFObjectTypes)
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

	if apiObject.RiskPredictorNewDevice != nil && apiObject.RiskPredictorNewDevice.GetId() != "" {
		diags.Append(p.toStateRiskPredictorNewDevice(apiObject.RiskPredictorNewDevice)...)
	}

	if apiObject.RiskPredictorUEBA != nil && apiObject.RiskPredictorUEBA.GetId() != "" {
		diags.Append(p.toStateRiskPredictorUEBA(apiObject.RiskPredictorUEBA)...)
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

func (p *riskPredictorResourceModel) toStateRiskPredictorNewDevice(apiObject *risk.RiskPredictorNewDevice) diag.Diagnostics {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorUEBA(apiObject *risk.RiskPredictorUEBA) diag.Diagnostics {
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
			"high":     framework.Int32OkToTF(v.GetHighOk()),
			"medium":   framework.Int32OkToTF(v.GetMediumOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityFallbackTFObjectTypes, o)
		diags.Append(d...)

		p.Fallback = objValue
	}

	// MaxDelay
	p.MaxDelay = types.ObjectNull(predictorVelocityMaxDelayTFObjectTypes)

	if v, ok := apiObject.GetMaxDelayOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"unit":     enumRiskPredictorUnitOkToTF(v.GetUnitOk()),
			"quantity": framework.Int32OkToTF(v.GetQuantityOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityMaxDelayTFObjectTypes, o)
		diags.Append(d...)

		p.MaxDelay = objValue
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
			"medium": framework.Int32OkToTF(v.GetMediumOk()),
			"high":   framework.Int32OkToTF(v.GetHighOk()),
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
