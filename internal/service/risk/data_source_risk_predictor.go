// Copyright © 2025 Ping Identity Corporation

package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
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
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne Risk Predictor.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to create the risk predictor in.",
				Required:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
			"risk_predictor_id": schema.StringAttribute{
				Description: "The ID of the risk predictor to retrieve.",
				Optional:    true,
				CustomType:  pingonetypes.ResourceIDType{},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the risk predictor to retrieve.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"compact_name": schema.StringAttribute{
				Description: "A string that specifies the compact name of the risk predictor. This is a unique name for the predictor.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the risk predictor.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "A string that specifies the type of the risk predictor.",
				Computed:    true,
			},
			"default": schema.SingleNestedAttribute{
				Description: "A single object that specifies the default result schema for the risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"weight": schema.Int32Attribute{
						Description: "The weight of the risk predictor.",
						Computed:    true,
					},
					"result": schema.SingleNestedAttribute{
						Description: "The result of the risk predictor.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "The type of the result.",
								Computed:    true,
							},
							"level": schema.StringAttribute{
								Description: "The level of the result.",
								Computed:    true,
							},
						},
					},
				},
			},
			"licensed": schema.BoolAttribute{
				Description: "A boolean that specifies whether the risk predictor is licensed.",
				Computed:    true,
			},
			"deletable": schema.BoolAttribute{
				Description: "A boolean that specifies whether the risk predictor is deletable.",
				Computed:    true,
			},
			"predictor_adversary_in_the_middle": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Adversary In The Middle risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allowed_domain_list": schema.SetAttribute{
						Description: "A set of strings that specifies the allowed domain list.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"predictor_anonymous_network": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Anonymous Network risk predictor.",
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
				Description: "A single object that specifies the configuration for the Bot Detection risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"include_repeated_events_without_sdk": schema.BoolAttribute{
						Description: "A boolean that specifies whether to include repeated events without SDK.",
						Computed:    true,
					},
				},
			},
			"predictor_composite": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Composite risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"compositions": schema.ListNestedAttribute{
						Description: "A list of objects that specifies the compositions.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"condition_json": schema.StringAttribute{
									Description: "The condition JSON.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"condition": schema.StringAttribute{
									Description: "The condition.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"level": schema.StringAttribute{
									Description: "The level.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"predictor_custom_map": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Custom Map risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.StringAttribute{
						Description: "The contains value.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type.",
						Computed:    true,
					},
					"between_ranges": schema.SingleNestedAttribute{
						Description: "The between ranges configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapBetweenRangesBoundSchemaDataSource(),
							"medium": customMapBetweenRangesBoundSchemaDataSource(),
							"low":    customMapBetweenRangesBoundSchemaDataSource(),
						},
					},
					"ip_ranges": schema.SingleNestedAttribute{
						Description: "The ip ranges configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapIpRangesBoundSchemaDataSource(),
							"medium": customMapIpRangesBoundSchemaDataSource(),
							"low":    customMapIpRangesBoundSchemaDataSource(),
						},
					},
					"string_list": schema.SingleNestedAttribute{
						Description: "The string list configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"high":   customMapStringValuesSchemaDataSource(),
							"medium": customMapStringValuesSchemaDataSource(),
							"low":    customMapStringValuesSchemaDataSource(),
						},
					},
				},
			},
			"predictor_device": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Device risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"activation_at": schema.StringAttribute{
						Description: "The activation time.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"detect": schema.StringAttribute{
						Description: "The detect mode.",
						Computed:    true,
					},
					"should_validate_payload_signature": schema.BoolAttribute{
						Description: "A boolean that specifies whether to validate payload signature.",
						Computed:    true,
					},
				},
			},
			"predictor_email_reputation": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Email Reputation risk predictor.",
				Computed:    true,
			},
			"predictor_geovelocity": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Geovelocity risk predictor.",
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
				Description: "A single object that specifies the configuration for the IP Reputation risk predictor.",
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
				Description: "A single object that specifies the configuration for the Traffic Anomaly risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"rules": schema.ListNestedAttribute{
						Description: "A list of objects that specifies the rules.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "A boolean that specifies whether the rule is enabled.",
									Computed:    true,
								},
								"interval": schema.SingleNestedAttribute{
									Description: "The interval configuration.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"unit": schema.StringAttribute{
											Description: "The unit.",
											Computed:    true,
										},
										"quantity": schema.Int32Attribute{
											Description: "The quantity.",
											Computed:    true,
										},
									},
								},
								"threshold": schema.SingleNestedAttribute{
									Description: "The threshold configuration.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"high": schema.Float32Attribute{
											Description: "The high threshold.",
											Computed:    true,
										},
										"medium": schema.Float32Attribute{
											Description: "The medium threshold.",
											Computed:    true,
										},
									},
								},
								"type": schema.StringAttribute{
									Description: "The type.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"predictor_user_location_anomaly": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the User Location Anomaly risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"radius": schema.SingleNestedAttribute{
						Description: "The radius configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"distance": schema.Int32Attribute{
								Description: "The distance.",
								Computed:    true,
							},
							"unit": schema.StringAttribute{
								Description: "The unit.",
								Computed:    true,
							},
						},
					},
					"days": schema.Int32Attribute{
						Description: "The days.",
						Computed:    true,
					},
				},
			},
			"predictor_user_risk_behavior": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the User Risk Behavior risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"prediction_model": schema.SingleNestedAttribute{
						Description: "The prediction model configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "The name.",
								Computed:    true,
							},
						},
					},
				},
			},
			"predictor_velocity": schema.SingleNestedAttribute{
				Description: "A single object that specifies the configuration for the Velocity risk predictor.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"by": schema.SetAttribute{
						Description: "The by configuration.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"every": schema.SingleNestedAttribute{
						Description: "The every configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description: "The unit.",
								Computed:    true,
							},
							"quantity": schema.Int32Attribute{
								Description: "The quantity.",
								Computed:    true,
							},
							"min_sample": schema.Int32Attribute{
								Description: "The min sample.",
								Computed:    true,
							},
						},
					},
					"fallback": schema.SingleNestedAttribute{
						Description: "The fallback configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"strategy": schema.StringAttribute{
								Description: "The strategy.",
								Computed:    true,
							},
							"high": schema.Float32Attribute{
								Description: "The high.",
								Computed:    true,
							},
							"medium": schema.Float32Attribute{
								Description: "The medium.",
								Computed:    true,
							},
						},
					},
					"measure": schema.StringAttribute{
						Description: "The measure.",
						Computed:    true,
					},
					"of": schema.StringAttribute{
						Description: "The of.",
						Computed:    true,
					},
					"sliding_window": schema.SingleNestedAttribute{
						Description: "The sliding window configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description: "The unit.",
								Computed:    true,
							},
							"quantity": schema.Int32Attribute{
								Description: "The quantity.",
								Computed:    true,
							},
							"min_sample": schema.Int32Attribute{
								Description: "The min sample.",
								Computed:    true,
							},
						},
					},
					"use": schema.SingleNestedAttribute{
						Description: "The use configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "The type.",
								Computed:    true,
							},
							"medium": schema.Float32Attribute{
								Description: "The medium.",
								Computed:    true,
							},
							"high": schema.Float32Attribute{
								Description: "The high.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func customMapBetweenRangesBoundSchemaDataSource() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"min_value": schema.Float32Attribute{
				Computed: true,
			},
			"max_value": schema.Float32Attribute{
				Computed: true,
			},
		},
		Computed: true,
	}
}

func customMapIpRangesBoundSchemaDataSource() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"values": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
		Computed: true,
	}
}

func customMapStringValuesSchemaDataSource() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"values": schema.SetAttribute{
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
		var response *http.Response
		var err error
		riskPredictor, response, err = r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, data.EnvironmentId.ValueString(), data.RiskPredictorId.ValueString()).Execute()
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				resp.Diagnostics.AddError(
					"Resource Failure",
					fmt.Sprintf("Risk Predictor with ID %s not found", data.RiskPredictorId.ValueString()),
				)
			} else {
				resp.Diagnostics.AddError(
					"Resource Failure",
					fmt.Sprintf("Unable to read Risk Predictor with ID %s: %s", data.RiskPredictorId.ValueString(), err),
				)
			}
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		pagedIterator := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, data.EnvironmentId.ValueString()).Execute()

		var found bool
		for pageCursor, err := range pagedIterator {
			if err != nil {
				resp.Diagnostics.AddError(
					"Resource Failure",
					fmt.Sprintf("Unable to read Risk Predictors: %s", err),
				)
				return
			}

			if riskPredictors, ok := pageCursor.EntityArray.Embedded.GetRiskPredictorsOk(); ok {
				for _, rp := range riskPredictors {
					// Use toState to extract the name and check if it matches
					tempModel := &riskPredictorDataSourceModel{}
					tempModel.toState(&rp)

					if tempModel.Name.ValueString() == data.Name.ValueString() {
						// Found it
						// We need to take a copy because rp is reused in loop? No, it's value.
						// But check pointer.
						val := rp
						riskPredictor = &val
						found = true
						break
					}
				}
			}

			if found {
				break
			}
		}

		if !found {
			resp.Diagnostics.AddError(
				"Resource Failure",
				fmt.Sprintf("Risk Predictor with name %s not found", data.Name.ValueString()),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing configuration",
			"One of 'risk_predictor_id' or 'name' must be configured.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(riskPredictor)...)
	resp.State.Set(ctx, &data)
}

func (p *riskPredictorDataSourceModel) toState(apiObject *risk.RiskPredictor) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := risk.RiskPredictorCommon{}

	if apiObject.RiskPredictorAdversaryInTheMiddle != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorAdversaryInTheMiddle.Id,
			Name:        apiObject.RiskPredictorAdversaryInTheMiddle.Name,
			CompactName: apiObject.RiskPredictorAdversaryInTheMiddle.CompactName,
			Description: apiObject.RiskPredictorAdversaryInTheMiddle.Description,
			Type:        apiObject.RiskPredictorAdversaryInTheMiddle.Type,
			Default:     apiObject.RiskPredictorAdversaryInTheMiddle.Default,
			Licensed:    apiObject.RiskPredictorAdversaryInTheMiddle.Licensed,
			Deletable:   apiObject.RiskPredictorAdversaryInTheMiddle.Deletable,
		}
	}

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

	if apiObject.RiskPredictorBotDetection != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorBotDetection.Id,
			Name:        apiObject.RiskPredictorBotDetection.Name,
			CompactName: apiObject.RiskPredictorBotDetection.CompactName,
			Description: apiObject.RiskPredictorBotDetection.Description,
			Type:        apiObject.RiskPredictorBotDetection.Type,
			Default:     apiObject.RiskPredictorBotDetection.Default,
			Licensed:    apiObject.RiskPredictorBotDetection.Licensed,
			Deletable:   apiObject.RiskPredictorBotDetection.Deletable,
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

	if apiObject.RiskPredictorEmailReputation != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorEmailReputation.Id,
			Name:        apiObject.RiskPredictorEmailReputation.Name,
			CompactName: apiObject.RiskPredictorEmailReputation.CompactName,
			Description: apiObject.RiskPredictorEmailReputation.Description,
			Type:        apiObject.RiskPredictorEmailReputation.Type,
			Default:     apiObject.RiskPredictorEmailReputation.Default,
			Licensed:    apiObject.RiskPredictorEmailReputation.Licensed,
			Deletable:   apiObject.RiskPredictorEmailReputation.Deletable,
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

	if apiObject.RiskPredictorTrafficAnomaly != nil {
		apiObjectCommon = risk.RiskPredictorCommon{
			Id:          apiObject.RiskPredictorTrafficAnomaly.Id,
			Name:        apiObject.RiskPredictorTrafficAnomaly.Name,
			CompactName: apiObject.RiskPredictorTrafficAnomaly.CompactName,
			Description: apiObject.RiskPredictorTrafficAnomaly.Description,
			Type:        apiObject.RiskPredictorTrafficAnomaly.Type,
			Default:     apiObject.RiskPredictorTrafficAnomaly.Default,
			Licensed:    apiObject.RiskPredictorTrafficAnomaly.Licensed,
			Deletable:   apiObject.RiskPredictorTrafficAnomaly.Deletable,
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

	p.Id = framework.PingOneResourceIDToTF(apiObjectCommon.GetId())
	p.RiskPredictorId = framework.PingOneResourceIDToTF(apiObjectCommon.GetId())
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObjectCommon.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObjectCommon.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(apiObjectCommon.GetTypeOk())
	p.Licensed = framework.BoolOkToTF(apiObjectCommon.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObjectCommon.GetDeletableOk())

	// Default block
	p.Default = types.ObjectNull(defaultTFObjectTypes)
	if v, ok := apiObjectCommon.GetDefaultOk(); ok {
		var d diag.Diagnostics

		defaultResultObj := types.ObjectNull(defaultResultTFObjectTypes)
		if v1, ok := v.GetResultOk(); ok {
			o := map[string]attr.Value{
				"type":  framework.EnumOkToTF(v1.GetTypeOk()),
				"level": framework.EnumOkToTF(v1.GetLevelOk()),
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

	// Set the predictor specific fields by object type
	var d diag.Diagnostics
	p.PredictorAdversaryInTheMiddle, d = p.toStateRiskPredictorAdversaryInTheMiddle(apiObject.RiskPredictorAdversaryInTheMiddle)
	diags.Append(d...)

	p.PredictorAnonymousNetwork, d = p.toStateRiskPredictorAnonymousNetwork(apiObject.RiskPredictorAnonymousNetwork)
	diags.Append(d...)

	p.PredictorBotDetection, d = p.toStateRiskPredictorBotDetection(apiObject.RiskPredictorBotDetection)
	diags.Append(d...)

	p.PredictorComposite, d = p.toStateRiskPredictorComposite(apiObject.RiskPredictorComposite)
	diags.Append(d...)

	p.PredictorCustomMap, d = p.toStateRiskPredictorCustom(apiObject.RiskPredictorCustom)
	diags.Append(d...)

	p.PredictorDevice, d = p.toStateRiskPredictorDevice(apiObject.RiskPredictorDevice)
	diags.Append(d...)

	p.PredictorEmailReputation, d = p.toStateRiskPredictorEmailReputation(apiObject.RiskPredictorEmailReputation)
	diags.Append(d...)

	p.PredictorGeoVelocity, d = p.toStateRiskPredictorGeovelocity(apiObject.RiskPredictorGeovelocity)
	diags.Append(d...)

	p.PredictorIPReputation, d = p.toStateRiskPredictorIPReputation(apiObject.RiskPredictorIPReputation)
	diags.Append(d...)

	p.PredictorTrafficAnomaly, d = p.toStateRiskPredictorTrafficAnomaly(apiObject.RiskPredictorTrafficAnomaly)
	diags.Append(d...)

	p.PredictorUserRiskBehavior, d = p.toStateRiskPredictorUserRiskBehavior(apiObject.RiskPredictorUserRiskBehavior)
	diags.Append(d...)

	p.PredictorUserLocationAnomaly, d = p.toStateRiskPredictorUserLocationAnomaly(apiObject.RiskPredictorUserLocationAnomaly)
	diags.Append(d...)

	p.PredictorVelocity, d = p.toStateRiskPredictorVelocity(apiObject.RiskPredictorVelocity)
	diags.Append(d...)

	return diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorAdversaryInTheMiddle(apiObject *risk.RiskPredictorAdversaryInTheMiddle) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorGenericAllowedDomainTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorGenericAllowedDomainTFObjectTypes, map[string]attr.Value{
		"allowed_domain_list": framework.StringSetOkToTF(apiObject.GetDomainWhiteListOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorAnonymousNetwork(apiObject *risk.RiskPredictorAnonymousNetwork) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorGenericAllowedCIDRTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorGenericAllowedCIDRTFObjectTypes, map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorBotDetection(apiObject *risk.RiskPredictorBotDetection) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorBotDetectionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorBotDetectionTFObjectTypes, map[string]attr.Value{
		"include_repeated_events_without_sdk": framework.BoolOkToTF(apiObject.GetIncludeRepeatedEventsWithoutSdkOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorComposite(apiObject *risk.RiskPredictorComposite) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorCompositeTFObjectTypes), diags
	}

	compositeObject := map[string]attr.Value{
		"compositions": types.ListNull(types.ObjectType{AttrTypes: predictorCompositionTFObjectTypes}),
	}

	v, ok := apiObject.GetCompositionsOk()
	compositions, d := p.riskPredictorCompositeConditionsOkToTF(v, ok)
	diags.Append(d...)

	compositeObject["compositions"] = compositions

	objValue, d := types.ObjectValue(predictorCompositeTFObjectTypes, compositeObject)
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) riskPredictorCompositeConditionsOkToTF(apiObject []risk.RiskPredictorCompositeAllOfCompositionsInner, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: predictorCompositionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectAttrTypes := []attr.Value{}
	for _, v := range apiObject {

		o := map[string]attr.Value{
			"level":          framework.EnumOkToTF(v.GetLevelOk()),
			"condition_json": types.StringNull(),
			"condition":      types.StringNull(),
		}

		if v1, ok := v.GetConditionOk(); ok {
			jsonString, err := json.Marshal(v1)
			if err != nil {
				diags.AddError(
					"Cannot convert map object to JSON string",
					"The provider cannot convert the `composite` map object to JSON.  Please report this to the provider maintainers.",
				)

				continue
			}

			conditionNormalized := jsontypes.NewNormalizedValue(string(jsonString))

			o["condition_json"] = conditionNormalized
			o["condition"] = conditionNormalized
		}

		objValue, d := types.ObjectValue(predictorCompositionTFObjectTypes, o)
		diags.Append(d...)

		objectAttrTypes = append(objectAttrTypes, objValue)
	}

	returnVar, d := types.ListValue(tfObjType, objectAttrTypes)
	diags.Append(d...)

	return returnVar, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorCustom(apiObject *risk.RiskPredictorCustom) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
	}

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

	if v, ok := apiObject.GetMapOk(); ok {
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

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_value": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_value": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
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

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_value": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_value": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
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

				o["contains"] = contains

				// Between Ranges
				setBetweenRanges = true

				levelObj := map[string]attr.Value{
					"min_value": framework.Float32OkToTF(v1.Between.GetMinScoreOk()),
					"max_value": framework.Float32OkToTF(v1.Between.GetMaxScoreOk()),
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
	}

	objValue, d := types.ObjectValue(predictorCustomMapTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorDevice(apiObject *risk.RiskPredictorDevice) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorDeviceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorDeviceTFObjectTypes, map[string]attr.Value{
		"activation_at":                     framework.TimeOkToTF(apiObject.GetActivationAtOk()),
		"detect":                            framework.EnumOkToTF(apiObject.GetDetectOk()),
		"should_validate_payload_signature": framework.BoolOkToTF(apiObject.GetShouldValidatePayloadSignatureOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorEmailReputation(apiObject *risk.RiskPredictorEmailReputation) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(map[string]attr.Type{}), diags
	}

	objValue, d := types.ObjectValue(map[string]attr.Type{}, map[string]attr.Value{})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorGeovelocity(apiObject *risk.RiskPredictorGeovelocity) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorGenericAllowedCIDRTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorGenericAllowedCIDRTFObjectTypes, map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorIPReputation(apiObject *risk.RiskPredictorIPReputation) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorGenericAllowedCIDRTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(predictorGenericAllowedCIDRTFObjectTypes, map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorTrafficAnomaly(apiObject *risk.RiskPredictorTrafficAnomaly) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorTrafficAnomalyTFObjectTypes), diags
	}

	rulesList, d := toStateRiskPredictorTrafficAnomalyRulesDataSource(apiObject.GetRulesOk())
	diags.Append(d...)

	o := map[string]attr.Value{
		"rules": rulesList,
	}

	objValue, d := types.ObjectValue(predictorTrafficAnomalyTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateRiskPredictorTrafficAnomalyRulesDataSource(apiObject []risk.RiskPredictorTrafficAnomalyAllOfRules, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: predictorTrafficAnomalyRulesTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	objectAttrTypes := []attr.Value{}
	for _, v := range apiObject {

		threshold, d := toStateRiskPredictorTrafficAnomalyRulesThresholdDataSource(v.GetThresholdOk())
		diags.Append(d...)

		interval, d := toStateRiskPredictorTrafficAnomalyRulesIntervalDataSource(v.GetIntervalOk())
		diags.Append(d...)

		o := map[string]attr.Value{
			"type":      framework.EnumOkToTF(v.GetTypeOk()),
			"enabled":   framework.BoolOkToTF(v.GetEnabledOk()),
			"threshold": threshold,
			"interval":  interval,
		}

		objValue, d := types.ObjectValue(predictorTrafficAnomalyRulesTFObjectTypes, o)
		diags.Append(d...)

		objectAttrTypes = append(objectAttrTypes, objValue)
	}

	returnVar, d := types.ListValue(tfObjType, objectAttrTypes)
	diags.Append(d...)

	return returnVar, diags
}

func toStateRiskPredictorTrafficAnomalyRulesThresholdDataSource(apiObject *risk.RiskPredictorTrafficAnomalyAllOfThreshold, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(predictorTrafficAnomalyRulesThresholdTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"high":   framework.Float32OkToTF(apiObject.GetHighOk()),
		"medium": framework.Float32OkToTF(apiObject.GetMediumOk()),
	}

	objValue, d := types.ObjectValue(predictorTrafficAnomalyRulesThresholdTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateRiskPredictorTrafficAnomalyRulesIntervalDataSource(apiObject *risk.RiskPredictorTrafficAnomalyAllOfInterval, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(predictorTrafficAnomalyRulesIntervalTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"unit":     framework.EnumOkToTF(apiObject.GetUnitOk()),
		"quantity": framework.Int32OkToTF(apiObject.GetQuantityOk()),
	}

	objValue, d := types.ObjectValue(predictorTrafficAnomalyRulesIntervalTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorUserRiskBehavior(apiObject *risk.RiskPredictorUserRiskBehavior) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorUserRiskBehaviorTFObjectTypes), diags
	}

	predictionModelObject := types.ObjectNull(predictorUserRiskBehaviorPredictionModelTFObjectTypes)

	if v, ok := apiObject.GetPredictionModelOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"name": framework.EnumOkToTF(v.GetNameOk()),
		}

		objValue, d := types.ObjectValue(predictorUserRiskBehaviorPredictionModelTFObjectTypes, o)
		diags.Append(d...)

		predictionModelObject = objValue
	}

	objValue, d := types.ObjectValue(predictorUserRiskBehaviorTFObjectTypes, map[string]attr.Value{
		"prediction_model": predictionModelObject,
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorUserLocationAnomaly(apiObject *risk.RiskPredictorUserLocationAnomaly) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorUserLocationAnomalyTFObjectTypes), diags
	}

	predictionRadiusObject := types.ObjectNull(predictorUserLocationAnomalyRadiusTFObjectTypes)

	if v, ok := apiObject.GetRadiusOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"distance": framework.Int32OkToTF(v.GetDistanceOk()),
			"unit":     framework.EnumOkToTF(v.GetUnitOk()),
		}

		objValue, d := types.ObjectValue(predictorUserLocationAnomalyRadiusTFObjectTypes, o)
		diags.Append(d...)

		predictionRadiusObject = objValue
	}

	objValue, d := types.ObjectValue(predictorUserLocationAnomalyTFObjectTypes, map[string]attr.Value{
		"radius": predictionRadiusObject,
		"days":   framework.Int32OkToTF(apiObject.GetDaysOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorDataSourceModel) toStateRiskPredictorVelocity(apiObject *risk.RiskPredictorVelocity) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorVelocityTFObjectTypes), diags
	}

	// Every
	modelEvery := types.ObjectNull(predictorVelocityEveryTFObjectTypes)

	if v, ok := apiObject.GetEveryOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"unit":       framework.EnumOkToTF(v.GetUnitOk()),
			"quantity":   framework.Int32OkToTF(v.GetQuantityOk()),
			"min_sample": framework.Int32OkToTF(v.GetMinSampleOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityEveryTFObjectTypes, o)
		diags.Append(d...)

		modelEvery = objValue
	}

	// Fallback
	modelFallback := types.ObjectNull(predictorVelocityFallbackTFObjectTypes)

	if v, ok := apiObject.GetFallbackOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"strategy": framework.EnumOkToTF(v.GetStrategyOk()),
			"high":     framework.Float32OkToTF(v.GetHighOk()),
			"medium":   framework.Float32OkToTF(v.GetMediumOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityFallbackTFObjectTypes, o)
		diags.Append(d...)

		modelFallback = objValue
	}

	// SlidingWindow
	modelSlidingWindow := types.ObjectNull(predictorVelocitySlidingWindowTFObjectTypes)

	if v, ok := apiObject.GetSlidingWindowOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"unit":       framework.EnumOkToTF(v.GetUnitOk()),
			"quantity":   framework.Int32OkToTF(v.GetQuantityOk()),
			"min_sample": framework.Int32OkToTF(v.GetMinSampleOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocitySlidingWindowTFObjectTypes, o)
		diags.Append(d...)

		modelSlidingWindow = objValue
	}

	// Use
	modelUse := types.ObjectNull(predictorVelocityUseTFObjectTypes)

	if v, ok := apiObject.GetUseOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"type":   framework.EnumOkToTF(v.GetTypeOk()),
			"medium": framework.Float32OkToTF(v.GetMediumOk()),
			"high":   framework.Float32OkToTF(v.GetHighOk()),
		}

		objValue, d := types.ObjectValue(predictorVelocityUseTFObjectTypes, o)
		diags.Append(d...)

		modelUse = objValue
	}

	objValue, d := types.ObjectValue(predictorVelocityTFObjectTypes, map[string]attr.Value{
		"by":             framework.StringSetOkToTF(apiObject.GetByOk()),
		"every":          modelEvery,
		"fallback":       modelFallback,
		"measure":        framework.EnumOkToTF(apiObject.GetMeasureOk()),
		"of":             framework.StringOkToTF(apiObject.GetOfOk()),
		"sliding_window": modelSlidingWindow,
		"use":            modelUse,
	})
	diags.Append(d...)

	return objValue, diags
}
