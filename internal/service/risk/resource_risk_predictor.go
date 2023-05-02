package risk

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	Id                        types.String `tfsdk:"id"`
	EnvironmentId             types.String `tfsdk:"environment_id"`
	Name                      types.String `tfsdk:"name"`
	CompactName               types.String `tfsdk:"compact_name"`
	Description               types.String `tfsdk:"description"`
	Type                      types.String `tfsdk:"type"`
	DefaultDecisionValue      types.String `tfsdk:"default_decision_value"`
	Licensed                  types.Bool   `tfsdk:"licensed"`
	Deletable                 types.Bool   `tfsdk:"deletable"`
	PredictorAnonymousNetwork types.List   `tfsdk:"predictor_anonymous_network"`
	// PredictorComposite           types.List `tfsdk:"predictor_composite"`
	// PredictorCustom       types.List `tfsdk:"predictor_custom"`
	PredictorGeovelocity         types.List `tfsdk:"predictor_geovelocity"`
	PredictorIPReputation        types.List `tfsdk:"predictor_ip_reputation"`
	PredictorNewDevice           types.List `tfsdk:"predictor_new_device"`
	PredictorUserLocationAnomaly types.List `tfsdk:"predictor_user_location_anomaly"`
	// PredictorUEBA                types.List `tfsdk:"predictor_user_risk_behavior"`
	// PredictorVelocity            types.List `tfsdk:"predictor_velocity"`
}

type predictorCompositeModel struct { // TODO
}

type predictorCustomModel struct {
	AttributeMapping types.String `tfsdk:"attribute_mapping"`
	MapIPRangeValues types.List   `tfsdk:"map_ip_range_values"`
	MapRangeValues   types.List   `tfsdk:"map_range_values"`
	MapListValues    types.List   `tfsdk:"map_list_values"`
}

type predictorCustomMapModel struct {
	High   types.List `tfsdk:"high"`
	Medium types.List `tfsdk:"medium"`
	Low    types.List `tfsdk:"low"`
}

type predictorCustomMapIPRangeModel struct {
	CIDRRangeList types.Set `tfsdk:"cidr_range_list"`
}

type predictorCustomMapBetweenModel struct {
	MinimumValue types.List `tfsdk:"minimum_value"`
	MaximumValue types.List `tfsdk:"maximum_value"`
}

type predictorCustomMapListModel struct {
	ListItems types.List `tfsdk:"list_items"`
}

// anonymous network, geovelocity, IP reputation
type predictorMinimalAllowedCIDRModel struct {
	AllowedCIDRList types.Set `tfsdk:"allowed_cidr_list"`
}

type predictorNewDeviceModel struct {
	ActivationAt types.String `tfsdk:"activation_at"`
	Detect       types.String `tfsdk:"detect"`
}

type predictorUserLocationModel struct {
	RadiusDistanceUnit types.String `tfsdk:"radius_distance_unit"`
	RadiusDistance     types.Int64  `tfsdk:"radius_distance"`
}

type predictorUEBAModel struct {
	PredictionModel types.String `tfsdk:"prediction_model"`
}

type predictorVelocityModel struct {
	By            types.List   `tfsdk:"by"`
	Every         types.List   `tfsdk:"every"`
	Fallback      types.List   `tfsdk:"fallback"`
	MaxDelay      types.List   `tfsdk:"max_delay"`
	Measure       types.String `tfsdk:"measure"`
	Of            types.String `tfsdk:"of"`
	SlidingWindow types.List   `tfsdk:"sliding_window"`
	Use           types.List   `tfsdk:"use"`
}

type predictorVelocityEveryModel struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	MinSample types.Int64  `tfsdk:"min_sample"`
}

type predictorVelocityFallbackModel struct {
	Strategy types.String `tfsdk:"strategy"`
	High     types.Int64  `tfsdk:"high"`
	Medium   types.Int64  `tfsdk:"medium"`
}

type predictorVelocityMaxDelayModel struct {
	Unit     types.String `tfsdk:"unit"`
	Quantity types.String `tfsdk:"quantity"`
}

type predictorVelocitySlidingWindowModel struct {
	MinSample types.Int64  `tfsdk:"min_sample"`
	Quantity  types.Int64  `tfsdk:"quantity"`
	Unit      types.String `tfsdk:"unit"`
}

type predictorVelocityUseModel struct {
	High   types.Int64  `tfsdk:"high"`
	Medium types.Int64  `tfsdk:"medium"`
	Type   types.String `tfsdk:"type"`
}

var (
	riskPredictorCompositeModelTFObjectTypes = map[string]attr.Type{} // TODO

	riskPredictorCustomModelTFObjectTypes = map[string]attr.Type{
		"attribute_mapping": types.StringType,
		"map_ip_range_values": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorCustomMapIPRangeTFObjectTypes,
			},
		},
		"map_range_values": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorCustomMapBetweenTFObjectTypes,
			},
		},
		"map_list_values": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorCustomMapListTFObjectTypes,
			},
		},
	}

	riskPredictorCustomMapTFObjectTypes = map[string]attr.Type{
		"high":   types.ListType{ElemType: types.StringType},
		"medium": types.ListType{ElemType: types.StringType},
		"low":    types.ListType{ElemType: types.StringType},
	}

	riskPredictorCustomMapIPRangeTFObjectTypes = map[string]attr.Type{
		"cidr_range_list": types.SetType{ElemType: types.StringType},
	}

	riskPredictorCustomMapBetweenTFObjectTypes = map[string]attr.Type{
		"minimum_value": types.ListType{ElemType: types.StringType},
		"maximum_value": types.ListType{ElemType: types.StringType},
	}

	riskPredictorCustomMapListTFObjectTypes = map[string]attr.Type{
		"list_items": types.ListType{ElemType: types.StringType},
	}

	riskPredictorNewDeviceModelTFObjectTypes = map[string]attr.Type{
		"activation_at": types.StringType,
		"detect":        types.StringType,
	}

	// anonymous network, geovelocity, IP reputation
	riskPredictorMinimalAllowedCIDRModelTFObjectTypes = map[string]attr.Type{
		"allowed_cidr_list": types.SetType{ElemType: types.StringType},
	}

	riskPredictorUEBAModelTFObjectTypes = map[string]attr.Type{
		"prediction_model": types.StringType,
	}

	riskPredictorUserLocationModelTFObjectTypes = map[string]attr.Type{
		"radius_distance":      types.Int64Type,
		"radius_distance_unit": types.StringType,
	}

	riskPredictorVelocityModelTFObjectTypes = map[string]attr.Type{
		"by": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocityByModelTFObjectTypes,
			},
		},
		"every": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocityEveryModelTFObjectTypes,
			},
		},
		"fallback": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocityFallbackModelTFObjectTypes,
			},
		},
		"max_delay": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocityMaxDelayModelTFObjectTypes,
			},
		},
		"sliding_window": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocitySlidingWindowModelTFObjectTypes,
			},
		},
		"use": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: riskPredictorVelocityUseModelTFObjectTypes,
			},
		},
	}

	riskPredictorVelocityByModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.BoolType,
		"username":   types.BoolType,
	}

	riskPredictorVelocityEveryModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.Int64Type,
		"username":   types.Int64Type,
	}

	riskPredictorVelocityFallbackModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.StringType,
		"username":   types.StringType,
	}

	riskPredictorVelocityMaxDelayModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.Int64Type,
		"username":   types.Int64Type,
	}

	riskPredictorVelocitySlidingWindowModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.Int64Type,
		"username":   types.Int64Type,
	}

	riskPredictorVelocityUseModelTFObjectTypes = map[string]attr.Type{
		"ip_address": types.BoolType,
		"username":   types.BoolType,
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

	resultLevelDescriptionFmt := "A string that identifies the risk level. Options are `HIGH`, `MEDIUM`, and `LOW`."
	resultLevelDescription := framework.SchemaDescription{
		MarkdownDescription: resultLevelDescriptionFmt,
		Description:         strings.ReplaceAll(resultLevelDescriptionFmt, "`", "\""),
	}

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
				Computed:            true,
			},

			"default_decision_value": schema.StringAttribute{
				Description:         resultLevelDescription.Description,
				MarkdownDescription: resultLevelDescription.MarkdownDescription,
				Optional:            true,
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

			"licensed": schema.BoolAttribute{
				Description: "A boolean that indicates whether PingOne Risk is licensed for the environment.",
				Computed:    true,
			},

			"deletable": schema.BoolAttribute{
				Description: "A boolean that indicates the PingOne Risk predictor can be deleted or not.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"predictor_anonymous_network": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the Anonymous Network risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_cidr_list": schema.SetAttribute{
							Description:         resultLevelDescription.Description,
							MarkdownDescription: resultLevelDescription.MarkdownDescription,
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`), "Values must be valid CIDR format."),
								),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
						// path.MatchRelative().AtParent().AtName("predictor_composite"),
						// path.MatchRelative().AtParent().AtName("predictor_custom"),
						path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
						path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
						path.MatchRelative().AtParent().AtName("predictor_new_device"),
						path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
					// path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
					// path.MatchRelative().AtParent().AtName("predictor_velocity"),
					),
				},
			},

			// "predictor_composite": schema.ListNestedBlock{
			// 	Description: "A single block that contains configuration values for the composite risk predictor type.",

			// 	NestedObject: schema.NestedBlockObject{
			// 		Attributes: map[string]schema.Attribute{},
			// 	},

			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
			// 	},
			// },

			// "predictor_custom": schema.ListNestedBlock{
			// 	Description: "A single block that contains configuration values for the custom mapping risk predictor type.",

			// 	NestedObject: schema.NestedBlockObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"attribute_mapping": schema.StringAttribute{
			// 				Required: true,
			// 				Validators: []validator.String{
			// 					stringvalidator.LengthAtLeast(attrMinLength),
			// 					stringvalidator.RegexMatches(regexp.MustCompile(`^\${(event|details)[a-z\.]+}$`), `Attribute mapping must match regex pattern "^\${(event|details)[a-z\.]+}$"`),
			// 				},
			// 			},

			// 			"map_ip_range_values": schema.SetAttribute{
			// 				Description: "The mapping of risk levels for the IP ranges specified.",
			// 				Optional:    true,
			// 				ElementType: types.StringType,
			// 				Validators: []validator.Set{
			// 					setvalidator.SizeAtLeast(attrMinLength),
			// 					setvalidator.ValueStringsAre(
			// 						stringvalidator.RegexMatches(regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`), `IP CIDR range must match regex pattern "^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$"`),
			// 					),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_ip_range_values")),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_range_values")),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_list_values")),
			// 				},
			// 			},

			// 			"map_range_values": schema.SetAttribute{
			// 				Description: "The mapping of risk levels for numerical values in a minimum, maxiumum boundary.",
			// 				Optional:    true,
			// 				ElementType: types.StringType,
			// 				Validators: []validator.Set{
			// 					setvalidator.SizeAtLeast(attrMinLength),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_ip_range_values")),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_range_values")),
			// 					setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("map_list_values")),
			// 				},
			// 			},
			// 		},

			// 		Blocks: map[string]schema.Block{
			// 			"map_range_values": schema.ListNestedBlock{
			// 			}
			// 	},

			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
			// 		// listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
			// 		// listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
			// 		// listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
			// 		// listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
			// 	},
			// },

			"predictor_geovelocity": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the Geovelocity risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_cidr_list": schema.SetAttribute{
							Description:         resultLevelDescription.Description,
							MarkdownDescription: resultLevelDescription.MarkdownDescription,
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`), "Values must be valid CIDR format."),
								),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
						// path.MatchRelative().AtParent().AtName("predictor_composite"),
						// path.MatchRelative().AtParent().AtName("predictor_custom"),
						path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
						path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
						path.MatchRelative().AtParent().AtName("predictor_new_device"),
						path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
					// path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
					// path.MatchRelative().AtParent().AtName("predictor_velocity"),
					),
				},
			},

			"predictor_ip_reputation": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the IP Reputation risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_cidr_list": schema.SetAttribute{
							Description:         resultLevelDescription.Description,
							MarkdownDescription: resultLevelDescription.MarkdownDescription,
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`), "Values must be valid CIDR format."),
								),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
						// path.MatchRelative().AtParent().AtName("predictor_composite"),
						// path.MatchRelative().AtParent().AtName("predictor_custom"),
						path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
						path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
						path.MatchRelative().AtParent().AtName("predictor_new_device"),
						path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
					// path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
					// path.MatchRelative().AtParent().AtName("predictor_velocity"),
					),
				},
			},

			"predictor_new_device": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the new device risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"detect": schema.StringAttribute{
							Required: true,
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
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
						// path.MatchRelative().AtParent().AtName("predictor_composite"),
						// path.MatchRelative().AtParent().AtName("predictor_custom"),
						path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
						path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
						path.MatchRelative().AtParent().AtName("predictor_new_device"),
						path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
					// path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
					// path.MatchRelative().AtParent().AtName("predictor_velocity"),
					),
				},
			},

			"predictor_user_location_anomaly": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the user location anomaly risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"radius_distance_unit": schema.StringAttribute{
							Description: "The unit of measurement for the `radius_distance` parameter.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(string(risk.ENUMDISTANCEUNIT_KILOMETERS)),
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
						"radius_distance": schema.Int64Attribute{
							Description: "The radius distance above which anomalies are detected.",
							Required:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(10),
								int64validator.AtMost(160),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
						// path.MatchRelative().AtParent().AtName("predictor_composite"),
						// path.MatchRelative().AtParent().AtName("predictor_custom"),
						path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
						path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
						path.MatchRelative().AtParent().AtName("predictor_new_device"),
						path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
					// path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
					// path.MatchRelative().AtParent().AtName("predictor_velocity"),
					),
				},
			},

			// "predictor_user_risk_behavior": schema.ListNestedBlock{
			// 	Description: "A single block that contains configuration values for the user risk behavior risk predictor type.",

			// 	NestedObject: schema.NestedBlockObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"prediction_model": schema.ListAttribute{},
			// 		},
			// 	},

			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
			// 	},
			// },

			// "predictor_velocity": schema.ListNestedBlock{
			// 	Description: "A single block that contains configuration values for the IP/user velocity risk predictor type.",

			// 	NestedObject: schema.NestedBlockObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"by":             schema.ListAttribute{},
			// 			"every":          schema.ListAttribute{},
			// 			"fallback":       schema.ListAttribute{},
			// 			"max_delay":      schema.ListAttribute{},
			// 			"measure":        schema.ListAttribute{},
			// 			"of":             schema.ListAttribute{},
			// 			"sliding_window": schema.ListAttribute{},
			// 			"use":            schema.ListAttribute{},
			// 		},
			// 	},

			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
			// 		listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
			// 	},
			// },
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
	resp.Diagnostics.Append(state.toState(response.(*risk.RiskPredictor))...)
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
	resp.Diagnostics.Append(data.toState(response.(*risk.RiskPredictor))...)
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
	resp.Diagnostics.Append(state.toState(response.(*risk.RiskPredictor))...)
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

	if !p.PredictorAnonymousNetwork.IsNull() {
		riskPredictor.RiskPredictorAnonymousNetwork, d = p.expandPredictorAnonymousNetwork(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// if !p.PredictorComposite.IsNull() {
	// 	riskPredictor.RiskPredictorComposite, d = p.expandPredictorComposite(ctx)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}
	// }

	// if !p.PredictorCustom.IsNull() {
	// 	riskPredictor.RiskPredictorCustom, d = p.expandPredictorCustom(ctx)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}
	// }

	if !p.PredictorGeovelocity.IsNull() {
		riskPredictor.RiskPredictorGeovelocity, d = p.expandPredictorGeovelocity(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	if !p.PredictorIPReputation.IsNull() {
		riskPredictor.RiskPredictorIPReputation, d = p.expandPredictorIPReputation(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	if !p.PredictorNewDevice.IsNull() {
		riskPredictor.RiskPredictorNewDevice, d = p.expandPredictorNewDevice(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	if !p.PredictorUserLocationAnomaly.IsNull() {
		riskPredictor.RiskPredictorUserLocationAnomaly, d = p.expandPredictorUserLocationAnomaly(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// if !p.PredictorUEBA.IsNull() {
	// 	riskPredictor.RiskPredictorUEBA, d = p.expandPredictorUEBA(ctx)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}
	// }

	// if !p.PredictorVelocity.IsNull() {
	// 	riskPredictor.RiskPredictorVelocity, d = p.expandPredictorVelocity(ctx)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}
	// }

	return riskPredictor, diags
}

func (p *riskPredictorResourceModel) expandPredictorAnonymousNetwork(ctx context.Context) (*risk.RiskPredictorAnonymousNetwork, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPredictorAnonymousNetwork(p.Name.ValueString(), p.CompactName.ValueString(), risk.ENUMPREDICTORTYPE_ANONYMOUS_NETWORK)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.DefaultDecisionValue.IsNull() && !p.DefaultDecisionValue.IsUnknown() {
		defaultModel := risk.NewRiskPredictorCommonDefault(int32(5))
		defaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(p.DefaultDecisionValue.ValueString()))

		defaultModel.SetResult(*defaultResult)

		data.SetDefault(*defaultModel)
	}

	if !p.PredictorAnonymousNetwork.IsNull() && !p.PredictorAnonymousNetwork.IsUnknown() {
		var plan []predictorMinimalAllowedCIDRModel
		d := p.PredictorAnonymousNetwork.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, plan[0].AllowedCIDRList)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return data, diags
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

func (p *riskPredictorResourceModel) expandPredictorGeovelocity(ctx context.Context) (*risk.RiskPredictorGeovelocity, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPredictorGeovelocity(p.Name.ValueString(), p.CompactName.ValueString(), risk.ENUMPREDICTORTYPE_GEO_VELOCITY)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.DefaultDecisionValue.IsNull() && !p.DefaultDecisionValue.IsUnknown() {
		defaultModel := risk.NewRiskPredictorCommonDefault(int32(5))
		defaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(p.DefaultDecisionValue.ValueString()))

		defaultModel.SetResult(*defaultResult)

		data.SetDefault(*defaultModel)
	}

	if !p.PredictorGeovelocity.IsNull() && !p.PredictorGeovelocity.IsUnknown() {
		var plan []predictorMinimalAllowedCIDRModel
		d := p.PredictorGeovelocity.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, plan[0].AllowedCIDRList)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorIPReputation(ctx context.Context) (*risk.RiskPredictorIPReputation, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPredictorIPReputation(p.Name.ValueString(), p.CompactName.ValueString(), risk.ENUMPREDICTORTYPE_IP_REPUTATION)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.DefaultDecisionValue.IsNull() && !p.DefaultDecisionValue.IsUnknown() {
		defaultModel := risk.NewRiskPredictorCommonDefault(int32(5))
		defaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(p.DefaultDecisionValue.ValueString()))

		defaultModel.SetResult(*defaultResult)

		data.SetDefault(*defaultModel)
	}

	if !p.PredictorIPReputation.IsNull() && !p.PredictorIPReputation.IsUnknown() {
		var plan []predictorMinimalAllowedCIDRModel
		d := p.PredictorIPReputation.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, plan[0].AllowedCIDRList)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetWhiteList(valuesSlice)
		}
	}

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorNewDevice(ctx context.Context) (*risk.RiskPredictorNewDevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewRiskPredictorNewDevice(p.Name.ValueString(), p.CompactName.ValueString(), risk.ENUMPREDICTORTYPE_DEVICE, risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.DefaultDecisionValue.IsNull() && !p.DefaultDecisionValue.IsUnknown() {
		defaultModel := risk.NewRiskPredictorCommonDefault(int32(5))
		defaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(p.DefaultDecisionValue.ValueString()))

		defaultModel.SetResult(*defaultResult)

		data.SetDefault(*defaultModel)
	}

	if !p.PredictorNewDevice.IsNull() && !p.PredictorNewDevice.IsUnknown() {
		var plan []predictorNewDeviceModel
		d := p.PredictorNewDevice.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		data.SetDetect(risk.EnumPredictorNewDeviceDetectType(plan[0].Detect.ValueString()))

		if !plan[0].ActivationAt.IsNull() && !plan[0].ActivationAt.IsUnknown() {
			t, e := time.Parse(time.RFC3339, plan[0].ActivationAt.ValueString())
			if e != nil {
				diags.AddError(
					"Invalid data format",
					"Cannot convert activation_at to a date/time.  Please check the format is a valid RFC3339 date time format.")
				return nil, diags
			}

			data.SetActivationAt(t)
		}
	}

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUserLocationAnomaly(ctx context.Context) (*risk.RiskPredictorUserLocationAnomaly, diag.Diagnostics) {
	var diags diag.Diagnostics

	var radius *risk.RiskPredictorUserLocationAnomalyAllOfRadius

	if !p.PredictorUserLocationAnomaly.IsNull() && !p.PredictorUserLocationAnomaly.IsUnknown() {
		var plan []predictorUserLocationModel
		d := p.PredictorUserLocationAnomaly.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		radius = risk.NewRiskPredictorUserLocationAnomalyAllOfRadius(int32(plan[0].RadiusDistance.ValueInt64()), risk.EnumDistanceUnit(plan[0].RadiusDistanceUnit.ValueString()))
	}

	data := risk.NewRiskPredictorUserLocationAnomaly(p.Name.ValueString(), p.CompactName.ValueString(), risk.ENUMPREDICTORTYPE_USER_LOCATION_ANOMALY, 50, *radius)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.DefaultDecisionValue.IsNull() && !p.DefaultDecisionValue.IsUnknown() {
		defaultModel := risk.NewRiskPredictorCommonDefault(int32(5))
		defaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumRiskLevel(p.DefaultDecisionValue.ValueString()))

		defaultModel.SetResult(*defaultResult)

		data.SetDefault(*defaultModel)
	}

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUEBA(ctx context.Context) (*risk.RiskPredictorUEBA, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorUEBA

	return data, diags
}

func (p *riskPredictorResourceModel) expandPredictorVelocity(ctx context.Context) (*risk.RiskPredictorVelocity, diag.Diagnostics) {
	var diags diag.Diagnostics

	var data *risk.RiskPredictorVelocity

	return data, diags
}

func (p *riskPredictorResourceModel) toState(apiObject *risk.RiskPredictor) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	if apiObject.RiskPredictorAnonymousNetwork != nil && apiObject.RiskPredictorAnonymousNetwork.GetId() != "" {
		return p.toStateRiskPredictorAnonymousNetwork(apiObject.RiskPredictorAnonymousNetwork)
	}

	if apiObject.RiskPredictorComposite != nil && apiObject.RiskPredictorComposite.GetId() != "" {
		return p.toStateRiskPredictorComposite(apiObject.RiskPredictorComposite)
	}

	if apiObject.RiskPredictorCustom != nil && apiObject.RiskPredictorCustom.GetId() != "" {
		return p.toStateRiskPredictorCustom(apiObject.RiskPredictorCustom)
	}

	if apiObject.RiskPredictorGeovelocity != nil && apiObject.RiskPredictorGeovelocity.GetId() != "" {
		return p.toStateRiskPredictorGeovelocity(apiObject.RiskPredictorGeovelocity)
	}

	if apiObject.RiskPredictorIPReputation != nil && apiObject.RiskPredictorIPReputation.GetId() != "" {
		return p.toStateRiskPredictorIPReputation(apiObject.RiskPredictorIPReputation)
	}

	if apiObject.RiskPredictorNewDevice != nil && apiObject.RiskPredictorNewDevice.GetId() != "" {
		return p.toStateRiskPredictorNewDevice(apiObject.RiskPredictorNewDevice)
	}

	if apiObject.RiskPredictorUEBA != nil && apiObject.RiskPredictorUEBA.GetId() != "" {
		return p.toStateRiskPredictorUEBA(apiObject.RiskPredictorUEBA)
	}

	if apiObject.RiskPredictorUserLocationAnomaly != nil && apiObject.RiskPredictorUserLocationAnomaly.GetId() != "" {
		return p.toStateRiskPredictorUserLocationAnomaly(apiObject.RiskPredictorUserLocationAnomaly)
	}

	if apiObject.RiskPredictorVelocity != nil && apiObject.RiskPredictorVelocity.GetId() != "" {
		return p.toStateRiskPredictorVelocity(apiObject.RiskPredictorVelocity)
	}

	diags.AddError(
		"Data object missing",
		"Cannot convert the data object to state as the predictor type is not supported.  Please report this to the provider maintainers.",
	)

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

	p.Id = framework.StringToTF(apiObject.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObject.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObject.GetTypeOk())

	p.DefaultDecisionValue = types.StringNull()
	if v, ok := apiObject.GetDefaultOk(); ok {
		if c, ok := v.GetResultOk(); ok {
			p.DefaultDecisionValue = enumRiskPredictorDefaultResultLevelOkToTF(c.GetLevelOk())
		}
	}

	p.Licensed = framework.BoolOkToTF(apiObject.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObject.GetDeletableOk())

	tfObjType := types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes}
	blockObject := map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	}

	flattenedObj, d := types.ObjectValue(riskPredictorMinimalAllowedCIDRModelTFObjectTypes, blockObject)
	diags.Append(d...)

	predictorAnonymousNetwork, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	p.PredictorAnonymousNetwork = predictorAnonymousNetwork
	// p.PredictorComposite = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCompositeModelTFObjectTypes})
	// p.PredictorCustom = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCustomModelTFObjectTypes})
	p.PredictorGeovelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorIPReputation = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorNewDevice = types.ListNull(types.ObjectType{AttrTypes: riskPredictorNewDeviceModelTFObjectTypes})
	// p.PredictorUEBA = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUEBAModelTFObjectTypes})
	p.PredictorUserLocationAnomaly = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUserLocationModelTFObjectTypes})
	// p.PredictorVelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorVelocityModelTFObjectTypes})

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

	p.Id = framework.StringToTF(apiObject.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObject.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObject.GetTypeOk())

	p.DefaultDecisionValue = types.StringNull()
	if v, ok := apiObject.GetDefaultOk(); ok {
		if c, ok := v.GetResultOk(); ok {
			p.DefaultDecisionValue = enumRiskPredictorDefaultResultLevelOkToTF(c.GetLevelOk())
		}
	}

	p.Licensed = framework.BoolOkToTF(apiObject.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObject.GetDeletableOk())

	tfObjType := types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes}
	blockObject := map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	}

	flattenedObj, d := types.ObjectValue(riskPredictorMinimalAllowedCIDRModelTFObjectTypes, blockObject)
	diags.Append(d...)

	predictorGeovelocity, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	p.PredictorAnonymousNetwork = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	// p.PredictorComposite = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCompositeModelTFObjectTypes})
	// p.PredictorCustom = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCustomModelTFObjectTypes})
	p.PredictorGeovelocity = predictorGeovelocity
	p.PredictorIPReputation = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorNewDevice = types.ListNull(types.ObjectType{AttrTypes: riskPredictorNewDeviceModelTFObjectTypes})
	// p.PredictorUEBA = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUEBAModelTFObjectTypes})
	p.PredictorUserLocationAnomaly = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUserLocationModelTFObjectTypes})
	// p.PredictorVelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorVelocityModelTFObjectTypes})

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

	p.Id = framework.StringToTF(apiObject.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObject.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObject.GetTypeOk())

	p.DefaultDecisionValue = types.StringNull()
	if v, ok := apiObject.GetDefaultOk(); ok {
		if c, ok := v.GetResultOk(); ok {
			p.DefaultDecisionValue = enumRiskPredictorDefaultResultLevelOkToTF(c.GetLevelOk())
		}
	}

	p.Licensed = framework.BoolOkToTF(apiObject.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObject.GetDeletableOk())

	tfObjType := types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes}
	blockObject := map[string]attr.Value{
		"allowed_cidr_list": framework.StringSetOkToTF(apiObject.GetWhiteListOk()),
	}

	flattenedObj, d := types.ObjectValue(riskPredictorMinimalAllowedCIDRModelTFObjectTypes, blockObject)
	diags.Append(d...)

	predictorIPReputation, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	p.PredictorAnonymousNetwork = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	// p.PredictorComposite = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCompositeModelTFObjectTypes})
	// p.PredictorCustom = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCustomModelTFObjectTypes})
	p.PredictorGeovelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorIPReputation = predictorIPReputation
	p.PredictorNewDevice = types.ListNull(types.ObjectType{AttrTypes: riskPredictorNewDeviceModelTFObjectTypes})
	// p.PredictorUEBA = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUEBAModelTFObjectTypes})
	p.PredictorUserLocationAnomaly = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUserLocationModelTFObjectTypes})
	// p.PredictorVelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorVelocityModelTFObjectTypes})

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

	p.Id = framework.StringToTF(apiObject.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObject.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObject.GetTypeOk())

	p.DefaultDecisionValue = types.StringNull()
	if v, ok := apiObject.GetDefaultOk(); ok {
		if c, ok := v.GetResultOk(); ok {
			p.DefaultDecisionValue = enumRiskPredictorDefaultResultLevelOkToTF(c.GetLevelOk())
		}
	}

	p.Licensed = framework.BoolOkToTF(apiObject.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObject.GetDeletableOk())

	tfObjType := types.ObjectType{AttrTypes: riskPredictorNewDeviceModelTFObjectTypes}
	blockObject := map[string]attr.Value{
		"activation_at": framework.TimeOkToTF(apiObject.GetActivationAtOk()),
		"detect":        enumRiskPredictorNewDeviceDetectOkToTF(apiObject.GetDetectOk()),
	}

	flattenedObj, d := types.ObjectValue(riskPredictorNewDeviceModelTFObjectTypes, blockObject)
	diags.Append(d...)

	predictorNewDevice, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	p.PredictorAnonymousNetwork = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	// p.PredictorComposite = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCompositeModelTFObjectTypes})
	// p.PredictorCustom = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCustomModelTFObjectTypes})
	p.PredictorGeovelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorIPReputation = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorNewDevice = predictorNewDevice
	// p.PredictorUEBA = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUEBAModelTFObjectTypes})
	p.PredictorUserLocationAnomaly = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUserLocationModelTFObjectTypes})
	// p.PredictorVelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorVelocityModelTFObjectTypes})

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

	p.Id = framework.StringToTF(apiObject.GetId())
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.CompactName = framework.StringOkToTF(apiObject.GetCompactNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = enumRiskPredictorTypeOkToTF(apiObject.GetTypeOk())

	p.DefaultDecisionValue = types.StringNull()
	if v, ok := apiObject.GetDefaultOk(); ok {
		if c, ok := v.GetResultOk(); ok {
			p.DefaultDecisionValue = enumRiskPredictorDefaultResultLevelOkToTF(c.GetLevelOk())
		}
	}

	p.Licensed = framework.BoolOkToTF(apiObject.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObject.GetDeletableOk())

	tfObjType := types.ObjectType{AttrTypes: riskPredictorUserLocationModelTFObjectTypes}
	blockObject := map[string]attr.Value{
		"radius_distance":      framework.Int32OkToTF(apiObject.Radius.GetDistanceOk()),
		"radius_distance_unit": enumRiskPredictorDistanceUnitOkToTF(apiObject.Radius.GetUnitOk()),
	}

	flattenedObj, d := types.ObjectValue(riskPredictorUserLocationModelTFObjectTypes, blockObject)
	diags.Append(d...)

	predictorUserLocationAnomaly, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	p.PredictorAnonymousNetwork = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	// p.PredictorComposite = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCompositeModelTFObjectTypes})
	// p.PredictorCustom = types.ListNull(types.ObjectType{AttrTypes: riskPredictorCustomModelTFObjectTypes})
	p.PredictorGeovelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorIPReputation = types.ListNull(types.ObjectType{AttrTypes: riskPredictorMinimalAllowedCIDRModelTFObjectTypes})
	p.PredictorNewDevice = types.ListNull(types.ObjectType{AttrTypes: riskPredictorNewDeviceModelTFObjectTypes})
	// p.PredictorUEBA = types.ListNull(types.ObjectType{AttrTypes: riskPredictorUEBAModelTFObjectTypes})
	p.PredictorUserLocationAnomaly = predictorUserLocationAnomaly
	// p.PredictorVelocity = types.ListNull(types.ObjectType{AttrTypes: riskPredictorVelocityModelTFObjectTypes})

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

	return diags
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
