package risk

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type RiskPredictorResource struct {
	client *risk.APIClient
	region model.RegionMapping
}

type riskPredictorResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	EnvironmentId                types.String `tfsdk:"environment_id"`
	Name                         types.String `tfsdk:"name"`
	CompactName                  types.String `tfsdk:"compact_name"`
	Description                  types.String `tfsdk:"description"`
	Type                         types.String `tfsdk:"type"`
	DefaultValues                types.List   `tfsdk:"default_values"`
	Licensed                     types.Bool   `tfsdk:"licensed"`
	PredictorAnonymousNetwork    types.List   `tfsdk:"predictor_anonymous_network"`
	PredictorComposite           types.List   `tfsdk:"predictor_composite"`
	PredictorCustom              types.List   `tfsdk:"predictor_custom"`
	PredictorGeovelocity         types.List   `tfsdk:"predictor_geovelocity"`
	PredictorIPReputation        types.List   `tfsdk:"predictor_ip_reputation"`
	PredictorNewDevice           types.List   `tfsdk:"predictor_new_device"`
	PredictorUserLocationAnomaly types.List   `tfsdk:"predictor_user_location_anomaly"`
	PredictorUserRiskBehavior    types.List   `tfsdk:"predictor_user_risk_behavior"`
	PredictorVelocity            types.List   `tfsdk:"predictor_velocity"`
}

type defaultValuesModel struct {
	Weight    types.Int64 `tfsdk:"weight"`
	Score     types.Int64 `tfsdk:"score"`
	Evaluated types.Bool  `tfsdk:"evaluated"`
	Result    types.List  `tfsdk:"result"`
}

type defaultValuesResultModel struct {
	Level types.String `tfsdk:"level"`
	Type  types.String `tfsdk:"type"`
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

type predictorMinimalAllowedCIDRModel struct {
	AllowedCIDRList types.Set `tfsdk:"allowed_cidr_list"`
}

type predictorNewDeviceModel struct {
	ActivationAt types.String `tfsdk:"activation_at"`
}

type predictorUserLocationModel struct {
	Days           types.Int64 `tfsdk:"days"`
	RadiusDistance types.List  `tfsdk:"radius_distance"`
}

type predictorUserLocationRadiusModel struct {
	Distance     types.Int64  `tfsdk:"distance"`
	DistanceUnit types.String `tfsdk:"distance_unit"`
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
	emailSourceTFObjectTypes = map[string]attr.Type{
		"name":          types.StringType,
		"email_address": types.StringType,
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

	typeDescriptionFmt := "A string that specifies the type of the risk predictor.  This can be either `ANONYMOUS_NETWORK`, `COMPOSITE`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `NEW_DEVICE`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR` or `VELOCITY`."
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

			"licensed": schema.BoolAttribute{
				Description: "A boolean that indicates whether PingOne Risk is licensed for the environment.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"default_values": schema.ListNestedBlock{
				Description: "A single block that contains the default values used for the risk predictor.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"weight": schema.Int64Attribute{
							Description: "An integer type that specifies the weight assigned to the risk predictor in a new policy by default.",
							Required:    true,
						},
						"score": schema.Int64Attribute{
							Description: "An integer type that specifies the score assigned to the risk predictor in a new policy by default.",
							Optional:    true,
						},
						"evaluated": schema.BoolAttribute{
							Description: "A boolean type.", // TODO
							Optional:    true,
						},
					},

					Blocks: map[string]schema.Block{
						"result": schema.ListNestedBlock{
							Description: "A single block that specifies the result assigned to the predictor if the predictor could not be calculated during the risk evaluation. If this field is not provided, and the predictor could not be calculated during risk evaluation, the following options are 1) If the predictor is used in an override, the override is skipped and 2) In the weighted policy, the predictor will have a weight of 0.",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"level": schema.StringAttribute{
										Description:         resultLevelDescription.Description,
										MarkdownDescription: resultLevelDescription.MarkdownDescription,
										Required:            true,
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

									"type": schema.StringAttribute{
										Description: "A string that specifies the risk evaluation result type. The only available option is `VALUE`.",
										Computed:    true,
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"predictor_anonymous_network": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the Anonymous Network risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_cidr_list": schema.SetAttribute{
							Description:         resultLevelDescription.Description,
										MarkdownDescription: resultLevelDescription.MarkdownDescription,
										Optional: true,
										ElementType: types.StringType,
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_composite": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the composite risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_custom": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the custom mapping risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attribute_mapping"
						"map_ip_range_values"
						"map_range_values"
						"map_list_values"
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_geovelocity": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the Geovelocity risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_cidr_list": schema.ListAttribute{},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_ip_reputation": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the IP reputation risk predictor type.",

				Attributes: map[string]schema.Attribute{
					"allowed_cidr_list": schema.ListAttribute{},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_new_device": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the new device risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"activation_at": schema.ListAttribute{},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_user_location_anomaly": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the user location anomaly risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"days": schema.ListAttribute{},
						"radius_distance": schema.ListAttribute{},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_user_risk_behavior": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the user risk behavior risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"prediction_model": schema.ListAttribute{},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
				},
			},

			"predictor_velocity": schema.ListNestedBlock{
				Description: "A single block that contains configuration values for the IP/user velocity risk predictor type.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"by": schema.ListAttribute{},
						"every": schema.ListAttribute{},
						"fallback": schema.ListAttribute{},
						"max_delay": schema.ListAttribute{},
						"measure": schema.ListAttribute{},
						"of": schema.ListAttribute{},
						"sliding_window": schema.ListAttribute{},
						"use": schema.ListAttribute{},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_anonymous_network")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_composite")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_custom")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_geovelocity")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_ip_reputation")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_new_device")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_location_anomaly")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_user_risk_behavior")),
					listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("predictor_velocity")),
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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
		},
		"UpdateEmailNotificationsSettings-Create",
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
	resp.Diagnostics.Append(state.toState(response.(*risk.NotificationsSettingsEmailDeliverySettings))...)
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
			return r.client.NotificationsSettingsSMTPApi.ReadEmailNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadEmailNotificationsSettings",
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
	resp.Diagnostics.Append(data.toState(response.(*risk.NotificationsSettingsEmailDeliverySettings))...)
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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
		},
		"UpdateEmailNotificationsSettings-Create",
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
	resp.Diagnostics.Append(state.toState(response.(*risk.NotificationsSettingsEmailDeliverySettings))...)
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
			r, err := r.client.NotificationsSettingsSMTPApi.DeleteEmailDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteEmailDeliverySettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RiskPredictorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 1
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[0])...)
}

func (p *riskPredictorResourceModel) expand(ctx context.Context) (*risk.NotificationsSettingsEmailDeliverySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := risk.NewNotificationsSettingsEmailDeliverySettings()

	if !p.Host.IsNull() && !p.Host.IsUnknown() {
		data.SetHost(p.Host.ValueString())
	}

	if !p.Port.IsNull() && !p.Port.IsUnknown() {
		data.SetPort(int32(p.Port.ValueInt64()))
	}

	if !p.Username.IsNull() && !p.Username.IsUnknown() {
		data.SetUsername(p.Username.ValueString())
	}

	if !p.Password.IsNull() && !p.Password.IsUnknown() {
		data.SetPassword(p.Password.ValueString())
	}

	if !p.From.IsNull() && !p.From.IsUnknown() {
		var plan []EmailSourceModel
		d := p.From.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		from := risk.NewNotificationsSettingsEmailDeliverySettingsFrom(plan[0].EmailAddress.ValueString())

		if !plan[0].Name.IsNull() && !plan[0].Name.IsUnknown() {
			from.SetName(plan[0].Name.ValueString())
		}

		data.SetFrom(*from)
	}

	if !p.ReplyTo.IsNull() && !p.ReplyTo.IsUnknown() {
		var plan []EmailSourceModel
		d := p.ReplyTo.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		replyTo := risk.NewNotificationsSettingsEmailDeliverySettingsReplyTo()

		if !plan[0].EmailAddress.IsNull() && !plan[0].EmailAddress.IsUnknown() {
			replyTo.SetAddress(plan[0].EmailAddress.ValueString())
		}

		if !plan[0].Name.IsNull() && !plan[0].Name.IsUnknown() {
			replyTo.SetName(plan[0].Name.ValueString())
		}

		data.SetReplyTo(*replyTo)
	}

	return data, diags
}

func (p *riskPredictorResourceModel) toState(apiObject *risk.NotificationsSettingsEmailDeliverySettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)

	p.Host = framework.StringOkToTF(apiObject.GetHostOk())
	p.Port = framework.Int32OkToTF(apiObject.GetPortOk())
	p.Protocol = framework.StringOkToTF(apiObject.GetProtocolOk())
	p.Username = framework.StringOkToTF(apiObject.GetUsernameOk())

	from, d := toStateEmailSource(apiObject.GetFromOk())
	diags.Append(d...)
	p.From = from

	replyTo, d := toStateEmailSource(apiObject.GetReplyToOk())
	diags.Append(d...)
	p.ReplyTo = replyTo

	return diags
}

func toStateEmailSource(emailSource interface{}, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: emailSourceTFObjectTypes}

	if !ok || emailSource == nil {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	var emailSourceMap map[string]attr.Value

	switch t := emailSource.(type) {
	case *risk.NotificationsSettingsEmailDeliverySettingsFrom:
		if t.GetAddress() == "" {
			return types.ListValueMust(tfObjType, []attr.Value{}), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *risk.NotificationsSettingsEmailDeliverySettingsReplyTo:
		if t.GetAddress() == "" {
			return types.ListValueMust(tfObjType, []attr.Value{}), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	default:
		diags.AddError(
			"Unexpected Email Source Type",
			fmt.Sprintf("Expected an email type object, got: %T. Please report this issue to the provider maintainers.", t),
		)

		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	flattenedObj, d := types.ObjectValue(emailSourceTFObjectTypes, emailSourceMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

}