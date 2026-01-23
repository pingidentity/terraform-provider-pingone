// Copyright © 2026 Ping Identity Corporation

package risk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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
type RiskPredictorDataSource serviceClientType

type riskPredictorDataSourceModel struct {
	Id                            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                 pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	RiskPredictorId               pingonetypes.ResourceIDValue `tfsdk:"risk_predictor_id"`
	Name                          types.String                 `tfsdk:"name"`
	CompactName                   types.String                 `tfsdk:"compact_name"`
	Description                   types.String                 `tfsdk:"description"`
	Type                          types.String                 `tfsdk:"type"`
	Default                       types.Object                 `tfsdk:"default"`
	Licensed                      types.Bool                   `tfsdk:"licensed"`
	Deletable                     types.Bool                   `tfsdk:"deletable"`
	PredictorAdversaryInTheMiddle types.Object                 `tfsdk:"predictor_adversary_in_the_middle"`
	PredictorAnonymousNetwork     types.Object                 `tfsdk:"predictor_anonymous_network"`
	PredictorBotDetection         types.Object                 `tfsdk:"predictor_bot_detection"`
	PredictorComposite            types.Object                 `tfsdk:"predictor_composite"`
	PredictorCustomMap            types.Object                 `tfsdk:"predictor_custom_map"`
	PredictorDevice               types.Object                 `tfsdk:"predictor_device"`
	PredictorEmailReputation      types.Object                 `tfsdk:"predictor_email_reputation"`
	PredictorGeoVelocity          types.Object                 `tfsdk:"predictor_geovelocity"`
	PredictorIPReputation         types.Object                 `tfsdk:"predictor_ip_reputation"`
	PredictorTrafficAnomaly       types.Object                 `tfsdk:"predictor_traffic_anomaly"`
	PredictorUserLocationAnomaly  types.Object                 `tfsdk:"predictor_user_location_anomaly"`
	PredictorUserRiskBehavior     types.Object                 `tfsdk:"predictor_user_risk_behavior"`
	PredictorVelocity             types.Object                 `tfsdk:"predictor_velocity"`
}

// Framework interfaces
var (
	_ datasource.DataSource              = &RiskPredictorDataSource{}
	_ datasource.DataSourceWithConfigure = &RiskPredictorDataSource{}
)

func NewRiskPredictorDataSource() datasource.DataSource {
	return &RiskPredictorDataSource{}
}

// Metadata
func (r *RiskPredictorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_predictor"
}

// Schema
func (r *RiskPredictorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	riskPredictorIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the risk predictor to retrieve.",
	).ExactlyOneOf([]string{"risk_predictor_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the risk predictor.",
	).ExactlyOneOf([]string{"risk_predictor_id", "name"})

	resp.Schema = schema.Schema{
		Description: "Data source to retrieve a PingOne Risk Predictor.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to retrieve the risk predictor from.",
				Required:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"risk_predictor_id": schema.StringAttribute{
				Description:         riskPredictorIdDescription.Description,
				MarkdownDescription: riskPredictorIdDescription.MarkdownDescription,
				Optional:            true,
				CustomType:          pingonetypes.ResourceIDType{},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},
			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"compact_name": schema.StringAttribute{
				Description: "A string that specifies the unique name for the predictor for use in risk evaluation request/response payloads. The value must be alpha-numeric, with no special characters or spaces. This name is used in the API both for policy configuration, and in the Risk Evaluation response (under `details`). If the value used for `compact_name` relates to a built-in predictor (a predictor that cannot be deleted), then this resource will attempt to overwrite the predictor's configuration.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the risk predictor. Maximum length is 1024 characters.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "A string that specifies the type of the risk predictor.",
				Computed:    true,
			},
			"default": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies the default configuration values for the risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"weight": schema.Int32Attribute{
						Description: "A number that specifies the default weight for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.",
						Computed:    true,
					},
					"result": schema.SingleNestedAttribute{
						Description: "A single nested object that contains the result assigned to the predictor if the predictor could not be calculated during the risk evaluation. If this field is not provided, and the predictor could not be calculated during risk evaluation, the behavior is: 1) If the predictor is used in an override, the override is skipped; 2) In the weighted policy, the predictor will have a `weight` of `0`.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "The default result type.",
								Computed:    true,
							},
							"level": schema.StringAttribute{
								Description: "The default result level.",
								Computed:    true,
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
			"predictor_adversary_in_the_middle": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Adversary-In-The-Middle (AitM) predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allowed_domain_list": schema.SetAttribute{
						Description: "A set of domains that are ignored for the predictor results.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"predictor_anonymous_network": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Anonymous Network predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": schema.SetAttribute{
						Description: "A set of strings that specifies the allowed CIDR list.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"predictor_bot_detection": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Bot Detection predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"include_repeated_events_without_sdk": schema.BoolAttribute{
						Description: "A boolean that specifies whether to expand the range of bot activity that PingOne Protect can detect.",
						Computed:    true,
					},
				},
			},
			"predictor_composite": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Composite predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"compositions": schema.ListNestedAttribute{
						Description: "A list of compositions of risk factors you want to use, and the condition logic that determines when or whether a risk factor is applied. The minimum number of compositions is 1 and the maximum number of compositions is 3.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"condition_json": schema.StringAttribute{
									Description: "A string that specifies the condition logic for the composite risk predictor. The value must be a valid JSON string.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"condition": schema.StringAttribute{
									Description: "A string that specifies the condition logic for the composite risk predictor as applied to the service.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"level": schema.StringAttribute{
									Description: "A string that specifies the risk level for the composite risk predictor.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"predictor_custom_map": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Custom Map predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.StringAttribute{
						Description: "A string that specifies the attribute reference that contains the value to match in the custom map.  The attribute reference should come from either the incoming event (`${event.*}`) or the evaluation details (`${details.*}`).",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "A string that specifies the type of custom map predictor.",
						Computed:    true,
					},
					"between_ranges": schema.SingleNestedAttribute{
						Description: "A single nested object that describes the upper and lower bounds of ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapBetweenRangesBoundSchemaDataSource("High"),
							"medium": customMapBetweenRangesBoundSchemaDataSource("Medium"),
							"low":    customMapBetweenRangesBoundSchemaDataSource("Low"),
						},
					},
					"ip_ranges": schema.SingleNestedAttribute{
						Description: "A single nested object that describes IP CIDR ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapIpRangesBoundSchemaDataSource("High"),
							"medium": customMapIpRangesBoundSchemaDataSource("Medium"),
							"low":    customMapIpRangesBoundSchemaDataSource("Low"),
						},
					},
					"string_list": schema.SingleNestedAttribute{
						Description: "A single nested object that describes the string values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapStringValuesSchemaDataSource("High"),
							"medium": customMapStringValuesSchemaDataSource("Medium"),
							"low":    customMapStringValuesSchemaDataSource("Low"),
						},
					},
				},
			},
			"predictor_device": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Device predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"activation_at": schema.StringAttribute{
						Description: "A string that represents a date on which the learning process for the device predictor should be restarted.  Can only be configured where the `detect` parameter is `NEW_DEVICE`. This can be used in conjunction with the fallback setting (`default.result.level`) to force strong authentication when moving the predictor to production. The date should be in an RFC3339 format. Note that activation date uses UTC time.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"detect": schema.StringAttribute{
						Description: "A string that represents the type of device detection to use.",
						Computed:    true,
					},
					"should_validate_payload_signature": schema.BoolAttribute{
						Description: "Relevant only for Suspicious Device predictors. A boolean that, if set to `true`, then any risk policies that include this predictor will require that the Signals SDK payload be provided as a signed JWT whose signature will be verified before proceeding with risk evaluation. You instruct the Signals SDK to provide the payload as a signed JWT by using the `universalDeviceIdentification` flag during initialization of the SDK, or by selecting the relevant setting for the `skrisk` component in DaVinci flows.",
						Computed:    true,
					},
				},
			},
			"predictor_email_reputation": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Email reputation predictor.",
				Computed:    true,
			},
			"predictor_geovelocity": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Geovelocity predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": schema.SetAttribute{
						Description: "A set of strings that specifies the allowed CIDR list.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"predictor_ip_reputation": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the IP reputation predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": schema.SetAttribute{
						Description: "A set of strings that specifies the allowed CIDR list.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"predictor_traffic_anomaly": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Traffic Anomaly predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"rules": schema.ListNestedAttribute{
						Description: "A collection with a single rule to use for this traffic anomaly predictor.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "A boolean to use the defined rule in the predictor.",
									Computed:    true,
								},
								"interval": schema.SingleNestedAttribute{
									Description: "A single nested object that contains the fields used to define the timeframe to consider. The timeframe can be between 1 hour and 14 days.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"unit": schema.StringAttribute{
											Description: "A string that specifies time unit for defining the timeframe for tracking number of users on the device.",
											Computed:    true,
										},
										"quantity": schema.Int32Attribute{
											Description: "An integer that specifies the number of days or hours for the timeframe for tracking number of users on the device.",
											Computed:    true,
										},
									},
								},
								"threshold": schema.SingleNestedAttribute{
									Description: "A single nested object that contains the fields used to define the risk thresholds.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"high": schema.Float32Attribute{
											Description: "A float that specifies the number of users during the defined timeframe that will be considered High risk.",
											Computed:    true,
										},
										"medium": schema.Float32Attribute{
											Description: "A float that specifies the number of users during the defined timeframe that will be considered Medium risk.",
											Computed:    true,
										},
									},
								},
								"type": schema.StringAttribute{
									Description: "A string that specifies the type of velocity algorithm to use.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"predictor_user_location_anomaly": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the User Location Anomaly predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"radius": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the radius to apply to the predictor evaluation",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"distance": schema.Int32Attribute{
								Description: "An integer that specifies the distance to apply to the predictor evaluation.",
								Computed:    true,
							},
							"unit": schema.StringAttribute{
								Description: "A string that specifies the unit of distance to apply to the predictor distance.",
								Computed:    true,
							},
						},
					},
					"days": schema.Int32Attribute{
						Description: "An integer that specifies the number of days to apply to the predictor evaluation.",
						Computed:    true,
					},
				},
			},
			"predictor_user_risk_behavior": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the User Risk Behavior predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"prediction_model": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the prediction model to apply to the predictor evaluation.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "A string that specifies the name of the prediction model to apply to the predictor evaluation.",
								Computed:    true,
							},
						},
					},
				},
			},
			"predictor_velocity": schema.SingleNestedAttribute{
				Description: "A single nested object that specifies options for the Velocity predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"by": schema.SetAttribute{
						Description: "A set of string values that specifies the attribute references that denote the subject of the velocity metric.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"every": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the velocity algorithm.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description: "A string value that specifies the time unit to use when sampling data.",
								Computed:    true,
							},
							"quantity": schema.Int32Attribute{
								Description: "An integer constant that specifies the quantity of time units to use when sampling data.",
								Computed:    true,
							},
							"min_sample": schema.Int32Attribute{
								Description: "An integer constant that specifies the minimum number of data points to use when sampling data.",
								Computed:    true,
							},
						},
					},
					"fallback": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the fallback strategy.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"strategy": schema.StringAttribute{
								Description: "A string value that specifies the type of fallback strategy algorithm to use.",
								Computed:    true,
							},
							"high": schema.Float32Attribute{
								Description: "A float that specifies the high risk threshold.",
								Computed:    true,
							},
							"medium": schema.Float32Attribute{
								Description: "A float that specifies the medium risk threshold.",
								Computed:    true,
							},
						},
					},
					"measure": schema.StringAttribute{
						Description: "A string value that specifies the type of measure to use for the predictor.",
						Computed:    true,
					},
					"of": schema.StringAttribute{
						Description: "A string value that specifies the attribute reference for the value to aggregate when calculating velocity metrics.",
						Computed:    true,
					},
					"sliding_window": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the sliding window.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description: "A string value that specifies the time unit to use when sampling data over time.",
								Computed:    true,
							},
							"quantity": schema.Int32Attribute{
								Description: "An integer constant that specifies the quantity of time units to use when sampling data over time.",
								Computed:    true,
							},
							"min_sample": schema.Int32Attribute{
								Description: "An integer constant that specifies the minimum number of data points to use when sampling data over time.",
								Computed:    true,
							},
						},
					},
					"use": schema.SingleNestedAttribute{
						Description: "A single nested object that specifies options for the velocity algorithm.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "A string value that specifies the type of velocity algorithm to use.",
								Computed:    true,
							},
							"medium": schema.Float32Attribute{
								Description: "A float that specifies the medium risk threshold.",
								Computed:    true,
							},
							"high": schema.Float32Attribute{
								Description: "A float that specifies the high risk threshold.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func customMapBetweenRangesBoundSchemaDataSource(riskResult string) schema.SingleNestedAttribute {
	desc := fmt.Sprintf("A single nested object that describes the upper and lower bounds of ranges that map to a %s risk result.", riskResult)
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"min_value": schema.Float32Attribute{
				Description: "A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`. This represents the lower bound of this risk result range.",
				Computed:    true,
			},
			"max_value": schema.Float32Attribute{
				Description: "A number that specifies the maximum value of the attribute named in `predictor_custom_map.contains`. This represents the upper bound of this risk result range.",
				Computed:    true,
			},
		},
		Computed: true,
	}
}

func customMapIpRangesBoundSchemaDataSource(riskResult string) schema.SingleNestedAttribute {
	desc := fmt.Sprintf("A single nested object that describes the IP CIDR ranges that map to a %s risk result.", riskResult)
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"values": schema.SetAttribute{
				Description: "A set of strings, in CIDR format, that describe the CIDR ranges that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
		Computed: true,
	}
}

func customMapStringValuesSchemaDataSource(riskResult string) schema.SingleNestedAttribute {
	desc := fmt.Sprintf("A single nested object that describes the string values that map to a %s risk result.", riskResult)
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"values": schema.SetAttribute{
				Description: "A set of strings that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
		Computed: true,
	}
}

func (r *RiskPredictorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *RiskPredictorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data riskPredictorDataSourceModel

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

	var riskPredictor *risk.RiskPredictor

	if !data.RiskPredictorId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, data.EnvironmentId.ValueString(), data.RiskPredictorId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneRiskPredictor",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&riskPredictor,
		)...)

	} else if !data.Name.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if riskPredictors, ok := pageCursor.EntityArray.Embedded.GetRiskPredictorsOk(); ok {
						for _, rp := range riskPredictors {
							// Use toState to extract the name and check if it matches
							tempModel := &riskPredictorDataSourceModel{}
							tempModel.toState(ctx, &rp)

							if tempModel.Name.ValueString() == data.Name.ValueString() {
								val := rp
								return &val, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, fmt.Errorf("Risk Predictor with name %s not found", data.Name.ValueString())
			},
			"ReadAllRiskPredictors",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&riskPredictor,
		)...)

	} else {
		resp.Diagnostics.AddError(
			"Missing configuration",
			"One of 'risk_predictor_id' or 'name' must be configured.",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, riskPredictor)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *riskPredictorDataSourceModel) toState(ctx context.Context, apiObject *risk.RiskPredictor) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	resourceModel := &riskPredictorResourceModel{}
	diags.Append(resourceModel.toState(ctx, apiObject)...)
	if diags.HasError() {
		return diags
	}

	p.Id = resourceModel.Id
	p.RiskPredictorId = resourceModel.Id
	p.Name = resourceModel.Name
	p.CompactName = resourceModel.CompactName
	p.Description = resourceModel.Description
	p.Type = resourceModel.Type
	p.Licensed = resourceModel.Licensed
	p.Deletable = resourceModel.Deletable
	p.Default = resourceModel.Default

	p.PredictorAdversaryInTheMiddle = resourceModel.PredictorAdversaryInTheMiddle
	p.PredictorAnonymousNetwork = resourceModel.PredictorAnonymousNetwork
	p.PredictorBotDetection = resourceModel.PredictorBotDetection
	p.PredictorComposite = resourceModel.PredictorComposite
	p.PredictorCustomMap = resourceModel.PredictorCustomMap
	p.PredictorDevice = resourceModel.PredictorDevice
	p.PredictorEmailReputation = resourceModel.PredictorEmailReputation
	p.PredictorGeoVelocity = resourceModel.PredictorGeoVelocity
	p.PredictorIPReputation = resourceModel.PredictorIPReputation
	p.PredictorTrafficAnomaly = resourceModel.PredictorTrafficAnomaly
	p.PredictorUserRiskBehavior = resourceModel.PredictorUserRiskBehavior
	p.PredictorUserLocationAnomaly = resourceModel.PredictorUserLocationAnomaly
	p.PredictorVelocity = resourceModel.PredictorVelocity

	return diags
}
