package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	objectplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectplanmodifier"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	riskservicehelpers "github.com/pingidentity/terraform-provider-pingone/internal/service/risk/helpers"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type RiskPredictorResource serviceClientType

type riskPredictorResourceModel struct {
	Id                            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                 pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
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
	PredictorUserLocationAnomaly  types.Object                 `tfsdk:"predictor_user_location_anomaly"`
	PredictorUserRiskBehavior     types.Object                 `tfsdk:"predictor_user_risk_behavior"`
	PredictorVelocity             types.Object                 `tfsdk:"predictor_velocity"`
}

// Anonymous network, IP reputation, geovelocity
type predictorGenericAllowedCIDR struct {
	AllowedCIDRList types.Set `tfsdk:"allowed_cidr_list"`
}

// Adversary in the middle
type predictorGenericAllowedDomain struct {
	AllowedDomainList types.Set `tfsdk:"allowed_domain_list"`
}

// Bot Detection
type predictorBotDetection struct {
	IncludeRepeatedEventsWithoutSDK types.Bool `tfsdk:"include_repeated_events_without_sdk"`
}

// Composite
type predictorComposite struct {
	Compositions types.List `tfsdk:"compositions"`
}

type predictorComposition struct {
	ConditionJSON jsontypes.Normalized `tfsdk:"condition_json"`
	Condition     jsontypes.Normalized `tfsdk:"condition"`
	Level         types.String         `tfsdk:"level"`
}

// Custom Map
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
	MinScore types.Float32 `tfsdk:"min_value"`
	MaxScore types.Float32 `tfsdk:"max_value"`
}

type predictorCustomMapHMLList struct {
	Values types.Set `tfsdk:"values"`
}

// New device
type predictorDevice struct {
	ActivationAt                   timetypes.RFC3339 `tfsdk:"activation_at"`
	Detect                         types.String      `tfsdk:"detect"`
	ShouldValidatePayloadSignature types.Bool        `tfsdk:"should_validate_payload_signature"`
}

// User Location Anomaly
type predictorUserLocationAnomaly struct {
	Radius types.Object `tfsdk:"radius"`
	Days   types.Int32  `tfsdk:"days"`
}

type predictorUserLocationAnomalyRadius struct {
	Distance types.Int32  `tfsdk:"distance"`
	Unit     types.String `tfsdk:"unit"`
}

// User Risk Behavior
type predictorUserRiskBehavior struct {
	PredictionModel types.Object `tfsdk:"prediction_model"`
}

type predictorUserRiskBehaviorPredictionModel struct {
	Name types.String `tfsdk:"name"`
}

// Velocity
type predictorVelocity struct {
	By            types.Set    `tfsdk:"by"`
	Every         types.Object `tfsdk:"every"`
	Fallback      types.Object `tfsdk:"fallback"`
	Measure       types.String `tfsdk:"measure"`
	Of            types.String `tfsdk:"of"`
	SlidingWindow types.Object `tfsdk:"sliding_window"`
	Use           types.Object `tfsdk:"use"`
}

type predictorVelocityEvery struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int32  `tfsdk:"quantity"`
	MinSample types.Int32  `tfsdk:"min_sample"`
}

type predictorVelocityFallback struct {
	Strategy types.String  `tfsdk:"strategy"`
	High     types.Float32 `tfsdk:"high"`
	Medium   types.Float32 `tfsdk:"medium"`
}

type predictorVelocitySlidingWindow struct {
	Unit      types.String `tfsdk:"unit"`
	Quantity  types.Int32  `tfsdk:"quantity"`
	MinSample types.Int32  `tfsdk:"min_sample"`
}

type predictorVelocityUse struct {
	UseType types.String  `tfsdk:"type"`
	Medium  types.Float32 `tfsdk:"medium"`
	High    types.Float32 `tfsdk:"high"`
}

// Default
type predictorDefault struct {
	Weight types.Int32  `tfsdk:"weight"`
	Result types.Object `tfsdk:"result"`
}

type predictorDefaultResult struct {
	ResultType types.String `tfsdk:"type"`
	Level      types.String `tfsdk:"level"`
}

var (
	// Default
	defaultTFObjectTypes = map[string]attr.Type{
		"weight": types.Int32Type,
		"result": types.ObjectType{
			AttrTypes: defaultResultTFObjectTypes,
		},
	}

	defaultResultTFObjectTypes = map[string]attr.Type{
		"type":  types.StringType,
		"level": types.StringType,
	}

	// Anonymous network, IP reputation, geovelocity
	predictorGenericAllowedCIDRTFObjectTypes = map[string]attr.Type{
		"allowed_cidr_list": types.SetType{ElemType: types.StringType},
	}

	// Adversary in the middle
	predictorGenericAllowedDomainTFObjectTypes = map[string]attr.Type{
		"allowed_domain_list": types.SetType{ElemType: types.StringType},
	}

	// Bot Detection
	predictorBotDetectionTFObjectTypes = map[string]attr.Type{
		"include_repeated_events_without_sdk": types.BoolType,
	}

	// Composite
	predictorCompositeTFObjectTypes = map[string]attr.Type{
		"compositions": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: predictorCompositionTFObjectTypes,
			},
		},
	}

	predictorCompositionTFObjectTypes = map[string]attr.Type{
		"condition_json": jsontypes.NormalizedType{},
		"condition":      jsontypes.NormalizedType{},
		"level":          types.StringType,
	}

	// Custom Map
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
		"min_value": types.Float32Type,
		"max_value": types.Float32Type,
	}

	predictorCustomMapHMLListTFObjectTypes = map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}

	// Device
	predictorDeviceTFObjectTypes = map[string]attr.Type{
		"activation_at":                     timetypes.RFC3339Type{},
		"detect":                            types.StringType,
		"should_validate_payload_signature": types.BoolType,
	}

	// User Location Anomaly
	predictorUserLocationAnomalyTFObjectTypes = map[string]attr.Type{
		"radius": types.ObjectType{
			AttrTypes: predictorUserLocationAnomalyRadiusTFObjectTypes,
		},
		"days": types.Int32Type,
	}

	predictorUserLocationAnomalyRadiusTFObjectTypes = map[string]attr.Type{
		"distance": types.Int32Type,
		"unit":     types.StringType,
	}

	// User Risk Behavior
	predictorUserRiskBehaviorTFObjectTypes = map[string]attr.Type{
		"prediction_model": types.ObjectType{
			AttrTypes: predictorUserRiskBehaviorPredictionModelTFObjectTypes,
		},
	}

	predictorUserRiskBehaviorPredictionModelTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	// Velocity
	predictorVelocityTFObjectTypes = map[string]attr.Type{
		"by": types.SetType{ElemType: types.StringType},
		"every": types.ObjectType{
			AttrTypes: predictorVelocityEveryTFObjectTypes,
		},
		"fallback": types.ObjectType{
			AttrTypes: predictorVelocityFallbackTFObjectTypes,
		},
		"measure": types.StringType,
		"of":      types.StringType,
		"sliding_window": types.ObjectType{
			AttrTypes: predictorVelocitySlidingWindowTFObjectTypes,
		},
		"use": types.ObjectType{
			AttrTypes: predictorVelocityUseTFObjectTypes,
		},
	}

	predictorVelocityEveryTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int32Type,
		"min_sample": types.Int32Type,
	}

	predictorVelocityFallbackTFObjectTypes = map[string]attr.Type{
		"strategy": types.StringType,
		"high":     types.Float32Type,
		"medium":   types.Float32Type,
	}

	predictorVelocitySlidingWindowTFObjectTypes = map[string]attr.Type{
		"unit":       types.StringType,
		"quantity":   types.Int32Type,
		"min_sample": types.Int32Type,
	}

	predictorVelocityUseTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"medium": types.Float32Type,
		"high":   types.Float32Type,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &RiskPredictorResource{}
	_ resource.ResourceWithConfigure   = &RiskPredictorResource{}
	_ resource.ResourceWithImportState = &RiskPredictorResource{}
	_ resource.ResourceWithModifyPlan  = &RiskPredictorResource{}
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
	const defaultWeightValue = 5

	// General attrs
	compactNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the unique name for the predictor for use in risk evaluation request/response payloads. The value must be alpha-numeric, with no special characters or spaces. This name is used in the API both for policy configuration, and in the Risk Evaluation response (under `details`).  If the value used for `compact_name` relates to a built-in predictor (a predictor that cannot be deleted), then this resource will attempt to overwrite the predictor's configuration.",
	).RequiresReplace()

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the risk predictor.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorTypeEnumValues)

	// Default nested attr
	defaultWeightDescription := framework.SchemaAttributeDescriptionFromMarkdown("A number that specifies the default weight for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.").DefaultValue(fmt.Sprint(defaultWeightValue))

	defaultResultDescription := framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that contains the result assigned to the predictor if the predictor could not be calculated during the risk evaluation. If this field is not provided, and the predictor could not be calculated during risk evaluation, the behavior is: 1) If the predictor is used in an override, the override is skipped; 2) In the weighted policy, the predictor will have a `weight` of `0`.")

	defaultResultTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The default result type.",
	).AllowedValuesComplex(map[string]string{
		string(risk.ENUMRESULTTYPE_VALUE): "any custom attribute value that's defined",
	})

	defaultResultLevelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The default result level.",
	).AllowedValuesEnum(risk.AllowedEnumRiskLevelEnumValues)

	// Adversary In The Middle
	predictorAdversaryInTheMiddleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Adversary-In-The-Middle (AitM) predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// Anonymous network predictor
	predictorAnonymousNetworkDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Anonymous Network predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// Bot Detection predictor
	predictorBotDetectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Bot Detection predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// Composite Predictor
	predictorCompositeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Composite predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	const minCompositions = 1
	const maxCompositions = 3
	predictorCompositeCompositionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A list of compositions of risk factors you want to use, and the condition logic that determines when or whether a risk factor is applied.  The minimum number of compositions is %d and the maximum number of compositions is %d.", minCompositions, maxCompositions),
	)

	predictorCompositeCompositionLevelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the risk level for the composite risk predictor.",
	).AllowedValuesEnum(risk.AllowedEnumRiskLevelEnumValues)

	// Custom Map Predictor
	predictorCustomMapDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Custom Map predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	predictorCustomMapContainsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the attribute reference that contains the value to match in the custom map.  The attribute reference should come from either the incoming event (`${event.*}`) or the evaluation details (`${details.*}`).  When defining attribute references in Terraform, the leading `$` needs to be escaped with an additional `$` character, e.g. `contains = \"$${event.myattribute}\"`.",
	)

	predictorCustomMapBetweenRangesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that describes the upper and lower bounds of ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
	)

	predictorCustomMapIPRangesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that describes IP CIDR ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
	)

	predictorCustomMapStringsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that describes the string values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results.",
	)

	// Device Predictor
	predictorDeviceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Device predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// Email reputation
	predictorEmailReputationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Email reputation predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// Geovelocity Predictor
	predictorGeovelocityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Geovelocity predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	// IP reputation Predictor
	predictorIPReputationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the IP reputation predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	predictorDeviceDetectDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents the type of device detection to use.",
	).AllowedValuesComplex(map[string]string{
		string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE):        "to configure a model based on new devices",
		string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_SUSPICIOUS_DEVICE): "to configure a model based on detection of suspicious devices",
	}).DefaultValue(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE))

	predictorDeviceActivationAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that represents a date on which the learning process for the device predictor should be restarted.  Can only be configured where the `detect` parameter is `%s`. This can be used in conjunction with the fallback setting (`default.result.level`) to force strong authentication when moving the predictor to production. The date should be in an RFC3339 format. Note that activation date uses UTC time.", string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)),
	)

	predictorDeviceShouldValidatePayloadSignatureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Relevant only for Suspicious Device predictors. A boolean that, if set to `true`, then any risk policies that include this predictor will require that the Signals SDK payload be provided as a signed JWT whose signature will be verified before proceeding with risk evaluation. You instruct the Signals SDK to provide the payload as a signed JWT by using the `universalDeviceIdentification` flag during initialization of the SDK, or by selecting the relevant setting for the `skrisk` component in DaVinci flows.",
	)

	// User location Predictor
	predictorUserLocationAnomalyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the User Location Anomaly predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	predictorUserLocationAnomalyDistanceUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the unit of distance to apply to the predictor distance.",
	).AllowedValuesEnum(risk.AllowedEnumDistanceUnitEnumValues).DefaultValue(string(risk.ENUMDISTANCEUNIT_KILOMETERS))

	// User Risk Behavior Predictor
	predictorUserRiskBehaviorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the User Risk Behavior predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	predictorUserRiskBehaviorPredictionModelNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the prediction model to apply to the predictor evaluation.",
	).AllowedValuesComplex(map[string]string{
		string(risk.ENUMUSERRISKBEHAVIORRISKMODEL_POINTS):                  "to configure the user-based risk model",
		string(risk.ENUMUSERRISKBEHAVIORRISKMODEL_LOGIN_ANOMALY_STATISTIC): "to configure the organisation based risk model",
	})

	// Velocity Predictor
	predictorVelocityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested object that specifies options for the Velocity predictor.",
	).ExactlyOneOf(descriptionPredictorObjectPaths)

	predictorVelocityMeasureDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the type of measure to use for the predictor.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorVelocityMeasureEnumValues).DefaultValue(string(risk.ENUMPREDICTORVELOCITYMEASURE_DISTINCT_COUNT))

	predictorVelocityOfDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the attribute reference for the value to aggregate when calculating velocity metrics.",
	).AllowedValuesComplex(map[string]string{
		"${event.ip}":      "to configure IP address velocity by user ID",
		"${event.user.id}": "to configure user velocity by IP address",
	}).AppendMarkdownString("When defining attribute references in Terraform, the leading `$` needs to be escaped with an additional `$` character, e.g. `of = \"$${event.ip}\"`.")

	predictorVelocityByDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of string values that specifies the attribute references that denote the subject of the velocity metric.",
	).AllowedValuesComplex(map[string]string{
		"${event.ip}":      "denotes the velocity metric is calculated by IP address",
		"${event.user.id}": "denotes the velocity metric is calculated by user ID",
	})

	predictorVelocityUseTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the type of velocity algorithm to use.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorVelocityUseTypeEnumValues)

	predictorVelocityFallbackTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the type of fallback strategy algorithm to use.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorVelocityFallbackStrategyEnumValues)

	predictorVelocityEveryUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the time unit to use when sampling data.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorUnitEnumValues)

	predictorVelocitySlidingWindowUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string value that specifies the time unit to use when sampling data over time.",
	).AllowedValuesEnum(risk.AllowedEnumPredictorUnitEnumValues)

	// Schema
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage Risk predictors in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the risk predictor in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique, friendly name for the predictor. This name is displayed in the Risk Policies UI, when the admin is asked to define the overrides and weights in policy configuration and is unique per environment.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"compact_name": schema.StringAttribute{
				Description:         compactNameDescription.Description,
				MarkdownDescription: compactNameDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "The value must be alpha-numeric, with no special characters or spaces."),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the risk predictor. Maximum length is 1024 characters.").Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(attrDescriptionMaxLength),
				},
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"default": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies the default configuration values for the risk predictor.").Description,
				Optional:    true,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"weight": schema.Int32Attribute{
						Description:         defaultWeightDescription.Description,
						MarkdownDescription: defaultWeightDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int32default.StaticInt32(defaultWeightValue),
					},

					"result": schema.SingleNestedAttribute{
						Description:         defaultResultDescription.Description,
						MarkdownDescription: defaultResultDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

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

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumRiskLevelEnumValues)...),
								},
							},
						},
					},
				},
			},

			"licensed": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that indicates whether PingOne Risk is licensed for the environment.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"deletable": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that indicates the PingOne Risk predictor can be deleted or not.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"predictor_adversary_in_the_middle": schema.SingleNestedAttribute{
				Description:         predictorAdversaryInTheMiddleDescription.Description,
				MarkdownDescription: predictorAdversaryInTheMiddleDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"allowed_domain_list": schema.SetAttribute{
						Description: "A set of domains that are ignored for the predictor results.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,

						Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),

						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.RegexMatches(verify.IsDomain, "Values must be valid domains."),
							),
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_anonymous_network": schema.SingleNestedAttribute{
				Description:         predictorAnonymousNetworkDescription.Description,
				MarkdownDescription: predictorAnonymousNetworkDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": allowedCIDRSchemaAttribute(),
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_bot_detection": schema.SingleNestedAttribute{
				Description:         predictorBotDetectionDescription.Description,
				MarkdownDescription: predictorBotDetectionDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"include_repeated_events_without_sdk": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to expand the range of bot activity that PingOne Protect can detect.").Description,
						Optional:    true,
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_composite": schema.SingleNestedAttribute{
				Description:         predictorCompositeDescription.Description,
				MarkdownDescription: predictorCompositeDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"compositions": schema.ListNestedAttribute{
						Description:         predictorCompositeCompositionsDescription.Description,
						MarkdownDescription: predictorCompositeCompositionsDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.List{
							listvalidator.SizeAtLeast(minCompositions),
							listvalidator.SizeAtMost(maxCompositions),
						},

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"condition_json": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the condition logic for the composite risk predictor. The value must be a valid JSON string.").Description,
									Required:    true,

									CustomType: jsontypes.NormalizedType{},
								},

								"condition": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the condition logic for the composite risk predictor as applied to the service.").Description,
									Computed:    true,

									CustomType: jsontypes.NormalizedType{},
								},

								"level": schema.StringAttribute{
									Description:         predictorCompositeCompositionLevelDescription.Description,
									MarkdownDescription: predictorCompositeCompositionLevelDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumRiskLevelEnumValues)...),
									},
								},
							},
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_custom_map": schema.SingleNestedAttribute{
				Description:         predictorCustomMapDescription.Description,
				MarkdownDescription: predictorCustomMapDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"contains": schema.StringAttribute{
						Description:         predictorCustomMapContainsDescription.Description,
						MarkdownDescription: predictorCustomMapContainsDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(`^\$\{(event|details)\.[a-zA-Z0-9.]+\}$`), `Value must match the regex "^\$\{(event|details)\.[a-zA-Z0-9.]+\}$\".`),
						},
					},

					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of custom map predictor.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},

					"between_ranges": schema.SingleNestedAttribute{
						Description:         predictorCustomMapBetweenRangesDescription.Description,
						MarkdownDescription: predictorCustomMapBetweenRangesDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"high":   customMapBetweenRangesBoundSchema("high"),
							"medium": customMapBetweenRangesBoundSchema("medium"),
							"low":    customMapBetweenRangesBoundSchema("low"),
						},

						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("between_ranges"),
								path.MatchRelative().AtParent().AtName("ip_ranges"),
								path.MatchRelative().AtParent().AtName("string_list"),
							),
						},
					},

					"ip_ranges": schema.SingleNestedAttribute{
						Description:         predictorCustomMapIPRangesDescription.Description,
						MarkdownDescription: predictorCustomMapIPRangesDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"high":   customMapIpRangesBoundSchema("high"),
							"medium": customMapIpRangesBoundSchema("medium"),
							"low":    customMapIpRangesBoundSchema("low"),
						},

						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("between_ranges"),
								path.MatchRelative().AtParent().AtName("ip_ranges"),
								path.MatchRelative().AtParent().AtName("string_list"),
							),
						},
					},

					"string_list": schema.SingleNestedAttribute{
						Description:         predictorCustomMapStringsDescription.Description,
						MarkdownDescription: predictorCustomMapStringsDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"high":   customMapStringValuesSchema("high"),
							"medium": customMapStringValuesSchema("medium"),
							"low":    customMapStringValuesSchema("low"),
						},

						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("between_ranges"),
								path.MatchRelative().AtParent().AtName("ip_ranges"),
								path.MatchRelative().AtParent().AtName("string_list"),
							),
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_device": schema.SingleNestedAttribute{
				Description:         predictorDeviceDescription.Description,
				MarkdownDescription: predictorDeviceDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"detect": schema.StringAttribute{
						Description:         predictorDeviceDetectDescription.Description,
						MarkdownDescription: predictorDeviceDetectDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumPredictorNewDeviceDetectTypeEnumValues)...),
						},

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"activation_at": schema.StringAttribute{
						Description:         predictorDeviceActivationAtDescription.Description,
						MarkdownDescription: predictorDeviceActivationAtDescription.MarkdownDescription,
						Optional:            true,

						CustomType: timetypes.RFC3339Type{},

						Validators: []validator.String{
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_SUSPICIOUS_DEVICE)),
								path.MatchRelative().AtParent().AtName("detect"),
							),
						},
					},

					"should_validate_payload_signature": schema.BoolAttribute{
						Description:         predictorDeviceShouldValidatePayloadSignatureDescription.Description,
						MarkdownDescription: predictorDeviceShouldValidatePayloadSignatureDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Bool{
							boolvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)),
								path.MatchRelative().AtParent().AtName("detect"),
							),
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_email_reputation": schema.SingleNestedAttribute{
				Description:         predictorEmailReputationDescription.Description,
				MarkdownDescription: predictorEmailReputationDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_geovelocity": schema.SingleNestedAttribute{
				Description:         predictorGeovelocityDescription.Description,
				MarkdownDescription: predictorGeovelocityDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": allowedCIDRSchemaAttribute(),
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_ip_reputation": schema.SingleNestedAttribute{
				Description:         predictorIPReputationDescription.Description,
				MarkdownDescription: predictorIPReputationDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": allowedCIDRSchemaAttribute(),
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_user_location_anomaly": schema.SingleNestedAttribute{
				Description:         predictorUserLocationAnomalyDescription.Description,
				MarkdownDescription: predictorUserLocationAnomalyDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"radius": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the radius to apply to the predictor evaluation").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"distance": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the distance to apply to the predictor evaluation.").Description,
								Required:    true,
							},

							"unit": schema.StringAttribute{
								Description:         predictorUserLocationAnomalyDistanceUnitDescription.Description,
								MarkdownDescription: predictorUserLocationAnomalyDistanceUnitDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Default: stringdefault.StaticString(string(risk.ENUMDISTANCEUNIT_KILOMETERS)),

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumDistanceUnitEnumValues)...),
								},
							},
						},
					},

					"days": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of days to apply to the predictor evaluation.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_user_risk_behavior": schema.SingleNestedAttribute{
				Description:         predictorUserRiskBehaviorDescription.Description,
				MarkdownDescription: predictorUserRiskBehaviorDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"prediction_model": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the prediction model to apply to the predictor evaluation.").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description:         predictorUserRiskBehaviorPredictionModelNameDescription.Description,
								MarkdownDescription: predictorUserRiskBehaviorPredictionModelNameDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumUserRiskBehaviorRiskModelEnumValues)...),
								},
							},
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},

			"predictor_velocity": schema.SingleNestedAttribute{
				Description:         predictorVelocityDescription.Description,
				MarkdownDescription: predictorVelocityDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"measure": schema.StringAttribute{
						Description:         predictorVelocityMeasureDescription.Description,
						MarkdownDescription: predictorVelocityMeasureDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(risk.ENUMPREDICTORVELOCITYMEASURE_DISTINCT_COUNT)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(risk.AllowedEnumPredictorVelocityMeasureEnumValues)...),
						},
					},

					"of": schema.StringAttribute{
						Description:         predictorVelocityOfDescription.Description,
						MarkdownDescription: predictorVelocityOfDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf("${event.ip}", "${event.user.id}"),
						},

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"by": schema.SetAttribute{
						Description:         predictorVelocityByDescription.Description,
						MarkdownDescription: predictorVelocityByDescription.MarkdownDescription,
						Computed:            true,
						ElementType:         types.StringType,

						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("${event.user.id}", "${event.ip}"),
							),
						},

						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},

					"use": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the velocity algorithm.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description:         predictorVelocityUseTypeDescription.Description,
								MarkdownDescription: predictorVelocityUseTypeDescription.MarkdownDescription,
								Computed:            true,

								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},

							"medium": schema.Float32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A floating point value that specifies a medium risk threshold for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Float32{
									float32planmodifier.UseStateForUnknown(),
								},
							},

							"high": schema.Float32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A floating point value that specifies a high risk threshold for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Float32{
									float32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},

					"fallback": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the predictor fallback strategy.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},

						Attributes: map[string]schema.Attribute{
							"strategy": schema.StringAttribute{
								Description:         predictorVelocityFallbackTypeDescription.Description,
								MarkdownDescription: predictorVelocityFallbackTypeDescription.MarkdownDescription,
								Computed:            true,

								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},

							"high": schema.Float32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A floating point value that specifies a high risk threshold for the fallback strategy.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Float32{
									float32planmodifier.UseStateForUnknown(),
								},
							},

							"medium": schema.Float32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A floating point value that specifies a medium risk threshold for the fallback strategy.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Float32{
									float32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},

					"every": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the granularlity of data sampling.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},

						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description:         predictorVelocityEveryUnitDescription.Description,
								MarkdownDescription: predictorVelocityEveryUnitDescription.MarkdownDescription,
								Computed:            true,

								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},

							"quantity": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that denotes the quantity of unit intervals to use for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},

							"min_sample": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that denotes the minimum sample of data to use for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},

					"sliding_window": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies options for the distribution of data that is compared against to detect anomaly.").Description,
						Computed:    true,

						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},

						Attributes: map[string]schema.Attribute{
							"unit": schema.StringAttribute{
								Description:         predictorVelocitySlidingWindowUnitDescription.Description,
								MarkdownDescription: predictorVelocitySlidingWindowUnitDescription.MarkdownDescription,
								Computed:            true,

								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},

							"quantity": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that denotes the quantity of unit intervals to use for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},

							"min_sample": schema.Int32Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that denotes the minimum sample of data to use for the velocity algorithm.").Description,
								Computed:    true,

								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},

				Validators: predictorObjectValidators,

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},
			},
		},
	}
}

func allowedCIDRSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		Description: "A set of IP addresses (CIDRs) that are ignored for the predictor results. The list can include IPs in IPv4 format and IPs in IPv6 format.",
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,

		Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),

		Validators: []validator.Set{
			setvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(verify.IPv4IPv6Regexp, "Values must be valid IPv4 or IPv6 CIDR format."),
			),
		},
	}
}

var (
	hmlValidators = []validator.Object{
		objectvalidator.AtLeastOneOf(
			path.MatchRelative().AtParent().AtName("high"),
			path.MatchRelative().AtParent().AtName("medium"),
			path.MatchRelative().AtParent().AtName("low"),
		),
	}

	predictorObjectValidators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("predictor_adversary_in_the_middle"),
			path.MatchRelative().AtParent().AtName("predictor_anonymous_network"),
			path.MatchRelative().AtParent().AtName("predictor_bot_detection"),
			path.MatchRelative().AtParent().AtName("predictor_composite"),
			path.MatchRelative().AtParent().AtName("predictor_custom_map"),
			path.MatchRelative().AtParent().AtName("predictor_device"),
			path.MatchRelative().AtParent().AtName("predictor_email_reputation"),
			path.MatchRelative().AtParent().AtName("predictor_geovelocity"),
			path.MatchRelative().AtParent().AtName("predictor_ip_reputation"),
			path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly"),
			path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior"),
			path.MatchRelative().AtParent().AtName("predictor_velocity"),
		),
	}

	descriptionPredictorObjectPaths = []string{
		"predictor_adversary_in_the_middle",
		"predictor_anonymous_network",
		"predictor_bot_detection",
		"predictor_composite",
		"predictor_custom_map",
		"predictor_device",
		"predictor_email_reputation",
		"predictor_geovelocity",
		"predictor_ip_reputation",
		"predictor_user_location_anomaly",
		"predictor_user_risk_behavior",
		"predictor_velocity",
	}
)

func customMapBetweenRangesBoundSchema(riskResult string) schema.SingleNestedAttribute {
	predictorCustomMapBetweenRangesMinValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.",
	)

	return schema.SingleNestedAttribute{
		Description: fmt.Sprintf("A single nested object that describes the upper and lower bounds of ranges that map to a %s risk result.", riskResult),
		Optional:    true,

		Attributes: map[string]schema.Attribute{
			"min_value": schema.Float32Attribute{
				Description:         predictorCustomMapBetweenRangesMinValueDescription.Description,
				MarkdownDescription: predictorCustomMapBetweenRangesMinValueDescription.MarkdownDescription,
				Required:            true,
			},

			"max_value": schema.Float32Attribute{
				Description:         predictorCustomMapBetweenRangesMinValueDescription.Description,
				MarkdownDescription: predictorCustomMapBetweenRangesMinValueDescription.MarkdownDescription,
				Required:            true,
			},
		},

		Validators: hmlValidators,
	}
}

func customMapIpRangesBoundSchema(riskResult string) schema.SingleNestedAttribute {

	attributeDescription := framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single nested object that describes the IP CIDR ranges that map to a %s risk result.", riskResult))

	predictorCustomMapIPRangeValuesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings, in CIDR format, that describe the CIDR ranges that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.",
	)

	return customMapGenericValuesSchema(
		attributeDescription,
		predictorCustomMapIPRangeValuesDescription,
		hmlValidators,
		[]validator.String{
			stringvalidator.RegexMatches(verify.IPv4IPv6Regexp, "Values must be valid IPv4 or IPv6 CIDR format."),
		},
	)
}

func customMapStringValuesSchema(riskResult string) schema.SingleNestedAttribute {

	attributeDescription := framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single nested object that describes the string values that map to a %s risk result.", riskResult))

	predictorCustomMapStringValuesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.",
	)

	return customMapGenericValuesSchema(
		attributeDescription,
		predictorCustomMapStringValuesDescription,
		hmlValidators,
		[]validator.String{},
	)
}

func customMapGenericValuesSchema(attributeDescription framework.SchemaAttributeDescription, attributeValuesDescription framework.SchemaAttributeDescription, validators []validator.Object, valuesValidators []validator.String) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description:         attributeDescription.Description,
		MarkdownDescription: attributeDescription.MarkdownDescription,
		Optional:            true,

		Attributes: map[string]schema.Attribute{
			"values": schema.SetAttribute{
				Description:         attributeValuesDescription.Description,
				MarkdownDescription: attributeValuesDescription.MarkdownDescription,
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(valuesValidators...),
				},
			},
		},

		Validators: validators,
	}
}

// ModifyPlan
func (r *RiskPredictorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var compositions []predictorComposition
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("predictor_composite").AtName("compositions"), &compositions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(compositions) > 0 {

		normalisedCompositions := make([]attr.Value, 0, len(compositions))

		for _, composition := range compositions {
			// Check the structure of the composite condition
			resp.Diagnostics.Append(riskservicehelpers.CheckCompositeConditionStructure(ctx, composition.ConditionJSON.ValueString())...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Normalise the composite condition with what we expect the API will do
			normalisedJSON, d := riskservicehelpers.NormaliseCompositeCondition(ctx, composition.ConditionJSON.ValueString())
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}

			composition.Condition = jsontypes.NewNormalizedValue(*normalisedJSON)

			objValue, d := types.ObjectValue(predictorCompositionTFObjectTypes, map[string]attr.Value{
				"condition_json": composition.ConditionJSON,
				"condition":      composition.Condition,
				"level":          composition.Level,
			})
			resp.Diagnostics.Append(d...)
			if !resp.Diagnostics.HasError() {
				normalisedCompositions = append(normalisedCompositions, objValue)
			}
		}

		resp.Plan.SetAttribute(ctx, path.Root("predictor_composite").AtName("compositions"), types.ListValueMust(
			types.ObjectType{
				AttrTypes: predictorCompositionTFObjectTypes,
			},
			normalisedCompositions,
		))

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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *RiskPredictorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state riskPredictorResourceModel

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
	riskPredictor, predefinedPredictorId, d := plan.expand(ctx, r.Client.RiskAPIClient, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *risk.RiskPredictor
	if predefinedPredictorId == nil {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.CreateRiskPredictor(ctx, plan.EnvironmentId.ValueString()).RiskPredictor(*riskPredictor).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateRiskPredictor",
			riskPredictorCreateUpdateCustomErrorHandler,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.UpdateRiskPredictor(ctx, plan.EnvironmentId.ValueString(), *predefinedPredictorId).RiskPredictor(*riskPredictor).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateRiskPredictor",
			riskPredictorCreateUpdateCustomErrorHandler,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPredictorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *riskPredictorResourceModel

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
	var response *risk.RiskPredictor
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRiskPredictor",
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
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RiskPredictorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state riskPredictorResourceModel

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
	riskPredictor, _, d := plan.expand(ctx, r.Client.RiskAPIClient, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *risk.RiskPredictor
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.UpdateRiskPredictor(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).RiskPredictor(*riskPredictor).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateRiskPredictor",
		riskPredictorCreateUpdateCustomErrorHandler,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RiskPredictorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *riskPredictorResourceModel

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

	if data.Deletable.ValueBool() {
		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.DeleteRiskAdvancedPredictor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
			},
			"DeleteRiskAdvancedPredictor",
			framework.CustomErrorResourceNotFoundWarning,
			nil,
			nil,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {

		if v, ok := riskservicehelpers.BootstrapPredictorValues[data.CompactName.ValueString()]; ok {

			// Run the API call
			resp.Diagnostics.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					_, r, err := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.UpdateRiskPredictor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).RiskPredictor(v).Execute()
					return nil, r, err
				},
				"UpdateRiskPredictor",
				framework.CustomErrorResourceNotFoundWarning,
				nil,
				nil,
			)...)
		}

		resp.Diagnostics.AddWarning(
			"Risk Predictor not deletable",
			fmt.Sprintf("The risk predictor with id \"%s\" cannot be deleted due to API limitation.  The risk predictor has been left in place but is no longer managed by the provider.", data.Id.ValueString()),
		)
	}
}

func (r *RiskPredictorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "risk_predictor_id",
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

func riskPredictorCreateUpdateCustomErrorHandler(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
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
	}

	return diags
}

func (p *riskPredictorResourceModel) expand(ctx context.Context, apiClient *risk.APIClient, managementApiClient *management.APIClient) (*risk.RiskPredictor, *string, diag.Diagnostics) {
	var diags diag.Diagnostics

	riskPredictor := &risk.RiskPredictor{}
	var overwriteRiskPredictorId *string
	var d diag.Diagnostics
	var riskPredictorCommonData *risk.RiskPredictorCommon

	// Check if this is attempting to overwrite an existing predictor.  We'll only allows overwriting where deletable = false
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, p.EnvironmentId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, managementApiClient, p.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if predictors, ok := pageCursor.EntityArray.Embedded.GetRiskPredictorsOk(); ok {

					for _, predictor := range predictors {
						predictorObject := predictor.GetActualInstance()

						var predictorId string
						var predictorCompactName string
						var predictorDeletable bool

						switch t := predictorObject.(type) {
						case *risk.RiskPredictorAdversaryInTheMiddle:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorAnonymousNetwork:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorBotDetection:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorComposite:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorCustom:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorDevice:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorEmailReputation:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorGeovelocity:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorIPReputation:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorUserLocationAnomaly:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorUserRiskBehavior:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						case *risk.RiskPredictorVelocity:
							predictorId = t.GetId()
							predictorCompactName = t.GetCompactName()
							predictorDeletable = t.GetDeletable()
						}

						if strings.EqualFold(predictorCompactName, p.CompactName.ValueString()) && !predictorDeletable {
							return &predictorId, pageCursor.HTTPResponse, nil
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllRiskPredictors",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&overwriteRiskPredictorId,
	)...)
	if diags.HasError() {
		return nil, nil, diags
	}

	riskPredictorCommonData = risk.NewRiskPredictorCommon(p.Name.ValueString(), p.CompactName.ValueString(), risk.EnumPredictorType(p.Type.ValueString()))

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		riskPredictorCommonData.SetDescription(p.Description.ValueString())
	}

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		var defaultPlan predictorDefault
		d := p.Default.As(ctx, &defaultPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}

		dataDefault := risk.NewRiskPredictorCommonDefault(defaultPlan.Weight.ValueInt32())

		if !defaultPlan.Result.IsNull() && !defaultPlan.Result.IsUnknown() {
			var defaultResultPlan predictorDefaultResult
			d := defaultPlan.Result.As(ctx, &defaultResultPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			if diags.HasError() {
				return nil, nil, diags
			}

			dataDefaultResult := risk.NewRiskPredictorCommonDefaultResult(risk.EnumResultType(defaultResultPlan.ResultType.ValueString()))
			dataDefaultResult.SetLevel(risk.EnumRiskLevel(defaultResultPlan.Level.ValueString()))
			dataDefault.SetResult(*dataDefaultResult)

			riskPredictorCommonData.SetDefault(*dataDefault)
		}

		riskPredictorCommonData.SetDefault(*dataDefault)
	}

	if !p.PredictorAdversaryInTheMiddle.IsNull() && !p.PredictorAdversaryInTheMiddle.IsUnknown() {
		riskPredictor.RiskPredictorAdversaryInTheMiddle, d = p.expandPredictorAdversaryInTheMiddle(ctx, riskPredictorCommonData)
	}

	if !p.PredictorAnonymousNetwork.IsNull() && !p.PredictorAnonymousNetwork.IsUnknown() {
		riskPredictor.RiskPredictorAnonymousNetwork, d = p.expandPredictorAnonymousNetwork(ctx, riskPredictorCommonData)
	}

	if !p.PredictorBotDetection.IsNull() && !p.PredictorBotDetection.IsUnknown() {
		riskPredictor.RiskPredictorBotDetection, d = p.expandPredictorBotDetection(ctx, riskPredictorCommonData)
	}

	if !p.PredictorComposite.IsNull() && !p.PredictorComposite.IsUnknown() {
		riskPredictor.RiskPredictorComposite, d = p.expandPredictorComposite(ctx, riskPredictorCommonData)
	}

	if !p.PredictorCustomMap.IsNull() && !p.PredictorCustomMap.IsUnknown() {
		riskPredictor.RiskPredictorCustom, d = p.expandPredictorCustom(ctx, riskPredictorCommonData)
	}

	if !p.PredictorDevice.IsNull() && !p.PredictorDevice.IsUnknown() {
		riskPredictor.RiskPredictorDevice, d = p.expandPredictorDevice(ctx, riskPredictorCommonData)
	}

	if !p.PredictorEmailReputation.IsNull() && !p.PredictorEmailReputation.IsUnknown() {
		riskPredictor.RiskPredictorEmailReputation = p.expandPredictorEmailReputation(riskPredictorCommonData)
	}

	if !p.PredictorGeoVelocity.IsNull() && !p.PredictorGeoVelocity.IsUnknown() {
		riskPredictor.RiskPredictorGeovelocity, d = p.expandPredictorGeovelocity(ctx, riskPredictorCommonData)
	}

	if !p.PredictorIPReputation.IsNull() && !p.PredictorIPReputation.IsUnknown() {
		riskPredictor.RiskPredictorIPReputation, d = p.expandPredictorIPReputation(ctx, riskPredictorCommonData)
	}

	if !p.PredictorUserRiskBehavior.IsNull() && !p.PredictorUserRiskBehavior.IsUnknown() {
		riskPredictor.RiskPredictorUserRiskBehavior, d = p.expandPredictorUserRiskBehavior(ctx, riskPredictorCommonData)
	}

	if !p.PredictorUserLocationAnomaly.IsNull() && !p.PredictorUserLocationAnomaly.IsUnknown() {
		riskPredictor.RiskPredictorUserLocationAnomaly, d = p.expandPredictorUserLocationAnomaly(ctx, riskPredictorCommonData)
	}

	if !p.PredictorVelocity.IsNull() && !p.PredictorVelocity.IsUnknown() {
		riskPredictor.RiskPredictorVelocity, d = p.expandPredictorVelocity(ctx, riskPredictorCommonData)
	}

	diags.Append(d...)
	if diags.HasError() {
		return nil, nil, diags
	}

	return riskPredictor, overwriteRiskPredictorId, diags
}

func (p *riskPredictorResourceModel) expandPredictorAdversaryInTheMiddle(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorAdversaryInTheMiddle, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorAdversaryInTheMiddle{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_ADVERSARY_IN_THE_MIDDLE,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorGenericAllowedDomain
	d := p.PredictorAdversaryInTheMiddle.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.AllowedDomainList.IsNull() && !predictorPlan.AllowedDomainList.IsUnknown() {
		allowedDomainListSet, d := predictorPlan.AllowedDomainList.ToSetValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		valuesPointerSlice := framework.TFSetToStringSlice(ctx, allowedDomainListSet)
		if len(valuesPointerSlice) > 0 {
			valuesSlice := make([]string, 0)
			for i := range valuesPointerSlice {
				valuesSlice = append(valuesSlice, *valuesPointerSlice[i])
			}
			data.SetDomainWhiteList(valuesSlice)
		}
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorAnonymousNetwork(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorAnonymousNetwork, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorAnonymousNetwork{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_ANONYMOUS_NETWORK,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorGenericAllowedCIDR
	d := p.PredictorAnonymousNetwork.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.AllowedCIDRList.IsNull() && !predictorPlan.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := predictorPlan.AllowedCIDRList.ToSetValue(ctx)
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

func (p *riskPredictorResourceModel) expandPredictorBotDetection(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorBotDetection, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorBotDetection{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_BOT,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorBotDetection
	d := p.PredictorBotDetection.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.IncludeRepeatedEventsWithoutSDK.IsNull() && !predictorPlan.IncludeRepeatedEventsWithoutSDK.IsUnknown() {
		data.SetIncludeRepeatedEventsWithoutSdk(predictorPlan.IncludeRepeatedEventsWithoutSDK.ValueBool())
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorComposite(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorComposite, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorComposite{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_COMPOSITE,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorComposite
	d := p.PredictorComposite.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.Compositions.IsNull() && !predictorPlan.Compositions.IsUnknown() {
		var plan []predictorComposition
		d := predictorPlan.Compositions.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		dataCompositons := make([]risk.RiskPredictorCompositeAllOfCompositionsInner, 0)
		for _, composition := range plan {
			dataComposition, diags := composition.expandPredictorCompositeComposition(ctx)
			if diags.HasError() {
				return nil, diags
			}

			dataCompositons = append(dataCompositons, *dataComposition)
		}
		data.SetCompositions(dataCompositons)
	}

	return &data, diags
}

func (p *predictorComposition) expandPredictorCompositeComposition(ctx context.Context) (*risk.RiskPredictorCompositeAllOfCompositionsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	var level risk.EnumRiskLevel
	if !p.Level.IsNull() && !p.Level.IsUnknown() {
		level = risk.EnumRiskLevel(p.Level.ValueString())
	}

	var condition risk.RiskPredictorCompositeConditionBase
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {
		err := json.Unmarshal([]byte(p.Condition.ValueString()), &condition)
		if err != nil {
			tflog.Error(ctx, "Cannot parse the `condition` JSON", map[string]interface{}{
				"err": err,
			})
			diags.AddError(
				"Cannot parse the `condition` JSON",
				"The JSON string passed to the `condition` parameter cannot be parsed as JSON.  Please check the policy is a valid JSON structure.",
			)
			return nil, diags
		}

	}

	return risk.NewRiskPredictorCompositeAllOfCompositionsInner(
		condition,
		level,
	), diags
}

func (p *riskPredictorResourceModel) expandPredictorCustom(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorCustom, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorCustom{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_MAP,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorCustomMap
	d := p.PredictorCustomMap.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	var contains string
	if !predictorPlan.Contains.IsNull() && !predictorPlan.Contains.IsUnknown() {
		contains = predictorPlan.Contains.ValueString()
	}

	setHigh := false
	high := risk.RiskPredictorCustomItem{}
	setMedium := false
	medium := risk.RiskPredictorCustomItem{}
	setLow := false
	low := risk.RiskPredictorCustomItem{}

	if !predictorPlan.BetweenRanges.IsNull() && !predictorPlan.BetweenRanges.IsUnknown() {
		var hmlPlan predictorCustomMapHML
		d := predictorPlan.BetweenRanges.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
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
					highHmlPlan.MinScore.ValueFloat32(),
					highHmlPlan.MaxScore.ValueFloat32(),
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
					mediumHmlPlan.MinScore.ValueFloat32(),
					mediumHmlPlan.MaxScore.ValueFloat32(),
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
					lowHmlPlan.MinScore.ValueFloat32(),
					lowHmlPlan.MaxScore.ValueFloat32(),
				),
			)

			low.RiskPredictorCustomItemBetween = v
			setLow = true
		}
	}

	if !predictorPlan.IPRanges.IsNull() && !predictorPlan.IPRanges.IsUnknown() {
		var hmlPlan predictorCustomMapHML
		d := predictorPlan.IPRanges.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
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

	if !predictorPlan.StringList.IsNull() && !predictorPlan.StringList.IsUnknown() {
		var hmlPlan predictorCustomMapHML
		d := predictorPlan.StringList.As(ctx, &hmlPlan, basetypes.ObjectAsOptions{
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

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorDevice(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorDevice, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorDevice{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_DEVICE,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorDevice
	d := p.PredictorDevice.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if predictorPlan.Detect.IsNull() || predictorPlan.Detect.IsUnknown() {
		predictorPlan.Detect = framework.EnumToTF(risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE)
	}

	data.SetDetect(risk.EnumPredictorNewDeviceDetectType(predictorPlan.Detect.ValueString()))

	if !predictorPlan.ActivationAt.IsNull() && !predictorPlan.ActivationAt.IsUnknown() {
		t, d := predictorPlan.ActivationAt.ValueRFC3339Time()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetActivationAt(t)
	}

	if !predictorPlan.ShouldValidatePayloadSignature.IsNull() && !predictorPlan.ShouldValidatePayloadSignature.IsUnknown() {
		data.SetShouldValidatePayloadSignature(predictorPlan.ShouldValidatePayloadSignature.ValueBool())
	}

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorEmailReputation(riskPredictorCommon *risk.RiskPredictorCommon) *risk.RiskPredictorEmailReputation {

	data := risk.RiskPredictorEmailReputation{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_EMAIL_REPUTATION,
		Default:     riskPredictorCommon.Default,
	}

	return &data
}

func (p *riskPredictorResourceModel) expandPredictorGeovelocity(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorGeovelocity, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorGeovelocity{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_GEO_VELOCITY,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorGenericAllowedCIDR
	d := p.PredictorGeoVelocity.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.AllowedCIDRList.IsNull() && !predictorPlan.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := predictorPlan.AllowedCIDRList.ToSetValue(ctx)
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
		Type:        risk.ENUMPREDICTORTYPE_IP_REPUTATION,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorGenericAllowedCIDR
	d := p.PredictorIPReputation.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.AllowedCIDRList.IsNull() && !predictorPlan.AllowedCIDRList.IsUnknown() {
		allowedCIDRListSet, d := predictorPlan.AllowedCIDRList.ToSetValue(ctx)
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

func (p *riskPredictorResourceModel) expandPredictorUserLocationAnomaly(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorUserLocationAnomaly, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorUserLocationAnomaly{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_USER_LOCATION_ANOMALY,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorUserLocationAnomaly
	d := p.PredictorUserLocationAnomaly.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.Radius.IsNull() && !predictorPlan.Radius.IsUnknown() {
		var radiusPlan predictorUserLocationAnomalyRadius
		d := predictorPlan.Radius.As(ctx, &radiusPlan, basetypes.ObjectAsOptions{
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

		radius := risk.NewRiskPredictorUserLocationAnomalyAllOfRadius(radiusPlan.Distance.ValueInt32(), radiusPlanUnit)

		data.SetRadius(*radius)
	}

	days := 50
	data.SetDays(int32(days))

	return &data, diags
}

func (p *riskPredictorResourceModel) expandPredictorUserRiskBehavior(ctx context.Context, riskPredictorCommon *risk.RiskPredictorCommon) (*risk.RiskPredictorUserRiskBehavior, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.RiskPredictorUserRiskBehavior{
		Name:        riskPredictorCommon.Name,
		CompactName: riskPredictorCommon.CompactName,
		Description: riskPredictorCommon.Description,
		Type:        risk.ENUMPREDICTORTYPE_USER_RISK_BEHAVIOR,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorUserRiskBehavior
	d := p.PredictorUserRiskBehavior.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !predictorPlan.PredictionModel.IsNull() && !predictorPlan.PredictionModel.IsUnknown() {
		var plan predictorUserRiskBehaviorPredictionModel
		d := predictorPlan.PredictionModel.As(ctx, &plan, basetypes.ObjectAsOptions{
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
		Type:        risk.ENUMPREDICTORTYPE_VELOCITY,
		Default:     riskPredictorCommon.Default,
	}

	var predictorPlan predictorVelocity
	d := p.PredictorVelocity.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Of
	if !predictorPlan.Of.IsNull() && !predictorPlan.Of.IsUnknown() {
		data.SetOf(predictorPlan.Of.ValueString())
	}

	// By
	if !predictorPlan.By.IsNull() && !predictorPlan.By.IsUnknown() {
		bySet, d := predictorPlan.By.ToSetValue(ctx)
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
		if predictorPlan.Of.Equal(types.StringValue("${event.ip}")) {
			data.SetBy([]string{"${event.user.id}"})
		}

		if predictorPlan.Of.Equal(types.StringValue("${event.user.id}")) {
			data.SetBy([]string{"${event.ip}"})
		}
	}

	// Every
	if !predictorPlan.Every.IsNull() && !predictorPlan.Every.IsUnknown() {
		var plan predictorVelocityEvery
		d := predictorPlan.Every.As(ctx, &plan, basetypes.ObjectAsOptions{
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
			every.SetQuantity(plan.Quantity.ValueInt32())
		}

		if !plan.MinSample.IsNull() && !plan.MinSample.IsUnknown() {
			every.SetMinSample(plan.MinSample.ValueInt32())
		}

		data.SetEvery(*every)
	} else {
		every := risk.NewRiskPredictorVelocityAllOfEvery()
		every.SetUnit(risk.ENUMPREDICTORUNIT_HOUR)
		quantity := 1
		every.SetQuantity(int32(quantity))
		minSample := 5
		every.SetMinSample(int32(minSample))
		data.SetEvery(*every)
	}

	// Fallback
	if !predictorPlan.Fallback.IsNull() && !predictorPlan.Fallback.IsUnknown() {
		var plan predictorVelocityFallback
		d := predictorPlan.Fallback.As(ctx, &plan, basetypes.ObjectAsOptions{
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
			fallback.SetHigh(plan.High.ValueFloat32())
		}

		if !plan.Medium.IsNull() && !plan.Medium.IsUnknown() {
			fallback.SetMedium(plan.Medium.ValueFloat32())
		}

		data.SetFallback(*fallback)
	} else {
		fallback := risk.NewRiskPredictorVelocityAllOfFallback()
		fallback.SetStrategy(risk.ENUMPREDICTORVELOCITYFALLBACKSTRATEGY_ENVIRONMENT_MAX)

		if predictorPlan.Of.Equal(types.StringValue("${event.ip}")) {
			high := 30
			fallback.SetHigh(float32(high))
			medium := 20
			fallback.SetMedium(float32(medium))
		}

		if predictorPlan.Of.Equal(types.StringValue("${event.user.id}")) {
			high := 3500
			fallback.SetHigh(float32(high))
			medium := 2500
			fallback.SetMedium(float32(medium))
		}

		data.SetFallback(*fallback)
	}

	// Measure
	if !predictorPlan.Measure.IsNull() && !predictorPlan.Measure.IsUnknown() {
		data.SetMeasure(risk.EnumPredictorVelocityMeasure(predictorPlan.Measure.ValueString()))
	} else {
		data.SetMeasure(risk.ENUMPREDICTORVELOCITYMEASURE_DISTINCT_COUNT)
	}

	// SlidingWindow
	if !predictorPlan.SlidingWindow.IsNull() && !predictorPlan.SlidingWindow.IsUnknown() {
		var plan predictorVelocitySlidingWindow
		d := predictorPlan.SlidingWindow.As(ctx, &plan, basetypes.ObjectAsOptions{
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
			slidingWindow.SetQuantity(plan.Quantity.ValueInt32())
		}

		if !plan.MinSample.IsNull() && !plan.MinSample.IsUnknown() {
			slidingWindow.SetMinSample(plan.MinSample.ValueInt32())
		}

		data.SetSlidingWindow(*slidingWindow)
	} else {
		slidingWindow := risk.NewRiskPredictorVelocityAllOfSlidingWindow()
		slidingWindow.SetUnit(risk.ENUMPREDICTORUNIT_DAY)
		quantity := 7
		slidingWindow.SetQuantity(int32(quantity))
		minSample := 3
		slidingWindow.SetMinSample(int32(minSample))
		data.SetSlidingWindow(*slidingWindow)
	}

	// Use
	if !predictorPlan.Use.IsNull() && !predictorPlan.Use.IsUnknown() {
		var plan predictorVelocityUse
		d := predictorPlan.Use.As(ctx, &plan, basetypes.ObjectAsOptions{
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
			use.SetMedium(plan.Medium.ValueFloat32())
		}

		if !plan.High.IsNull() && !plan.High.IsUnknown() {
			use.SetHigh(plan.High.ValueFloat32())
		}

		data.SetUse(*use)
	} else {
		use := risk.NewRiskPredictorVelocityAllOfUse()
		use.SetType(risk.ENUMPREDICTORVELOCITYUSETYPE_POISSON_WITH_MAX)
		medium := 2
		use.SetMedium(float32(medium))
		high := 4
		use.SetHigh(float32(high))
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
	// p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
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

	p.Licensed = framework.BoolOkToTF(apiObjectCommon.GetLicensedOk())
	p.Deletable = framework.BoolOkToTF(apiObjectCommon.GetDeletableOk())

	// Set the predictor specific fields by object type
	var d diag.Diagnostics
	p.PredictorAdversaryInTheMiddle, d = p.toStateRiskPredictorAdversaryInTheMiddle(apiObject.RiskPredictorAdversaryInTheMiddle)
	diags.Append(d...)

	p.PredictorAnonymousNetwork, d = p.toStateRiskPredictorAnonymousNetwork(apiObject.RiskPredictorAnonymousNetwork)
	diags.Append(d...)

	p.PredictorBotDetection, d = p.toStateRiskPredictorBotDetection(apiObject.RiskPredictorBotDetection)
	diags.Append(d...)

	p.PredictorComposite, d = p.toStateRiskPredictorComposite(ctx, apiObject.RiskPredictorComposite)
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

	p.PredictorUserRiskBehavior, d = p.toStateRiskPredictorUserRiskBehavior(apiObject.RiskPredictorUserRiskBehavior)
	diags.Append(d...)

	p.PredictorUserLocationAnomaly, d = p.toStateRiskPredictorUserLocationAnomaly(apiObject.RiskPredictorUserLocationAnomaly)
	diags.Append(d...)

	p.PredictorVelocity, d = p.toStateRiskPredictorVelocity(apiObject.RiskPredictorVelocity)
	diags.Append(d...)

	return diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorAdversaryInTheMiddle(apiObject *risk.RiskPredictorAdversaryInTheMiddle) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorAnonymousNetwork(apiObject *risk.RiskPredictorAnonymousNetwork) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorBotDetection(apiObject *risk.RiskPredictorBotDetection) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorComposite(ctx context.Context, apiObject *risk.RiskPredictorComposite) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(predictorCompositeTFObjectTypes), diags
	}

	compositeObject := map[string]attr.Value{
		"compositions": types.ListNull(types.ObjectType{AttrTypes: predictorCompositionTFObjectTypes}),
	}

	// The JSON fields the admins define are direct-to-state (because the API changes the JSON structure on output), but if the resource is imported, we have to copy the API value to the JSON input fields
	v, ok := apiObject.GetCompositionsOk()
	compositions, d := p.riskPredictorCompositeConditionsOkToTF(ctx, v, ok)
	diags.Append(d...)

	compositeObject["compositions"] = compositions

	objValue, d := types.ObjectValue(predictorCompositeTFObjectTypes, compositeObject)
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorResourceModel) riskPredictorCompositeConditionsOkToTF(ctx context.Context, apiObject []risk.RiskPredictorCompositeAllOfCompositionsInner, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: predictorCompositionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	planJSONElements := make([]jsontypes.Normalized, 0)

	if !p.PredictorComposite.IsNull() && !p.PredictorComposite.IsUnknown() {

		var predictorPlan predictorComposite
		d := p.PredictorComposite.As(ctx, &predictorPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)

		if !predictorPlan.Compositions.IsNull() && !predictorPlan.Compositions.IsUnknown() {
			var plan []predictorComposition
			d := predictorPlan.Compositions.ElementsAs(ctx, &plan, false)
			diags.Append(d...)

			for _, planElement := range plan {
				planJSONElements = append(planJSONElements, planElement.ConditionJSON)
			}
		}
	}

	objectAttrTypes := []attr.Value{}
	for i, v := range apiObject {

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

			if i >= 0 && i < len(planJSONElements) {
				o["condition_json"] = planJSONElements[i]
			} else {
				o["condition_json"] = conditionNormalized
			}
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

func (p *riskPredictorResourceModel) toStateRiskPredictorCustom(apiObject *risk.RiskPredictorCustom) (basetypes.ObjectValue, diag.Diagnostics) {
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

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
				}

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

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
				}

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

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
				}

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

				if !o["contains"].IsNull() && !contains.Equal(o["contains"]) {
					diags.AddError(
						"Data object inconsistent",
						"Cannot convert the data object to state as the data object is inconsistent (\"contains\" value).  Please report this to the provider maintainers.",
					)

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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

					return types.ObjectNull(predictorCustomMapTFObjectTypes), diags
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
	}

	objValue, d := types.ObjectValue(predictorCustomMapTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorDevice(apiObject *risk.RiskPredictorDevice) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorEmailReputation(apiObject *risk.RiskPredictorEmailReputation) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(map[string]attr.Type{}), diags
	}

	objValue, d := types.ObjectValue(map[string]attr.Type{}, map[string]attr.Value{})
	diags.Append(d...)

	return objValue, diags
}

func (p *riskPredictorResourceModel) toStateRiskPredictorGeovelocity(apiObject *risk.RiskPredictorGeovelocity) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorIPReputation(apiObject *risk.RiskPredictorIPReputation) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorUserRiskBehavior(apiObject *risk.RiskPredictorUserRiskBehavior) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorUserLocationAnomaly(apiObject *risk.RiskPredictorUserLocationAnomaly) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (p *riskPredictorResourceModel) toStateRiskPredictorVelocity(apiObject *risk.RiskPredictorVelocity) (basetypes.ObjectValue, diag.Diagnostics) {
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
