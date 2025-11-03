// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationPolicyResource serviceClientType

type NotificationPolicyResourceModel struct {
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                  types.String                 `tfsdk:"name"`
	Default               types.Bool                   `tfsdk:"default"`
	CountryLimit          types.Object                 `tfsdk:"country_limit"`
	CooldownConfiguration types.Object                 `tfsdk:"cooldown_configuration"`
	ProviderConfiguration types.Object                 `tfsdk:"provider_configuration"`
	Quota                 types.Set                    `tfsdk:"quota"`
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type NotificationPolicyQuotaResourceModel struct {
	Type            types.String `tfsdk:"type"`
	DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Total           types.Int32  `tfsdk:"total"`
	Used            types.Int32  `tfsdk:"used"`
	Unused          types.Int32  `tfsdk:"unused"`
}

type NotificationPolicyCountryLimitResourceModel struct {
	Type            types.String `tfsdk:"type"`
	DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Countries       types.Set    `tfsdk:"countries"`
}

type NotificationPolicyCooldownConfigurationResourceModel struct {
	Email    types.Object `tfsdk:"email"`
	Sms      types.Object `tfsdk:"sms"`
	Voice    types.Object `tfsdk:"voice"`
	WhatsApp types.Object `tfsdk:"whats_app"`
}

type NotificationPolicyCooldownConfigurationMethodResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Periods     types.List   `tfsdk:"periods"`
	GroupBy     types.String `tfsdk:"group_by"`
	ResendLimit types.Int32  `tfsdk:"resend_limit"`
}

type NotificationPolicyCooldownConfigurationMethodPeriodResourceModel struct {
	Duration types.Int32  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type NotificationPolicyProviderConfigurationResourceModel struct {
	Conditions types.List `tfsdk:"conditions"`
}

type NotificationPolicyProviderConfigurationConditionResourceModel struct {
	DeliveryMethods types.Set  `tfsdk:"delivery_methods"`
	Countries       types.Set  `tfsdk:"countries"`
	FallbackChain   types.List `tfsdk:"fallback_chain"`
}

type NotificationPolicyProviderConfigurationFallbackChainItemResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

var (
	quotaTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
		"delivery_methods": types.SetType{
			ElemType: types.StringType,
		},
		"total":  types.Int32Type,
		"used":   types.Int32Type,
		"unused": types.Int32Type,
	}

	countryLimitTFObjectTypes = map[string]attr.Type{
		"type":             types.StringType,
		"delivery_methods": types.SetType{ElemType: types.StringType},
		"countries":        types.SetType{ElemType: types.StringType},
	}

	cooldownConfigurationPeriodTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int32Type,
		"time_unit": types.StringType,
	}

	cooldownConfigurationMethodTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
		"periods": types.ListType{
			ElemType: types.ObjectType{AttrTypes: cooldownConfigurationPeriodTFObjectTypes},
		},
		"group_by":     types.StringType,
		"resend_limit": types.Int32Type,
	}

	cooldownConfigurationTFObjectTypes = map[string]attr.Type{
		"email":     types.ObjectType{AttrTypes: cooldownConfigurationMethodTFObjectTypes},
		"sms":       types.ObjectType{AttrTypes: cooldownConfigurationMethodTFObjectTypes},
		"voice":     types.ObjectType{AttrTypes: cooldownConfigurationMethodTFObjectTypes},
		"whats_app": types.ObjectType{AttrTypes: cooldownConfigurationMethodTFObjectTypes},
	}

	fallbackChainItemTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	providerConfigurationConditionTFObjectTypes = map[string]attr.Type{
		"delivery_methods": types.SetType{ElemType: types.StringType},
		"countries":        types.SetType{ElemType: types.StringType},
		"fallback_chain": types.ListType{
			ElemType: types.ObjectType{AttrTypes: fallbackChainItemTFObjectTypes},
		},
	}

	providerConfigurationTFObjectTypes = map[string]attr.Type{
		"conditions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: providerConfigurationConditionTFObjectTypes},
		},
	}

	cooldownConfigurationMethodDefault = types.ObjectValueMust(
		cooldownConfigurationMethodTFObjectTypes,
		map[string]attr.Value{
			"enabled":      types.BoolValue(false),
			"periods":      types.ListNull(types.ObjectType{AttrTypes: cooldownConfigurationPeriodTFObjectTypes}),
			"group_by":     types.StringNull(),
			"resend_limit": types.Int32Null(),
		},
	)

	cooldownConfigurationDefault = types.ObjectValueMust(
		cooldownConfigurationTFObjectTypes,
		map[string]attr.Value{
			"email":     cooldownConfigurationMethodDefault,
			"sms":       cooldownConfigurationMethodDefault,
			"voice":     cooldownConfigurationMethodDefault,
			"whats_app": cooldownConfigurationMethodDefault,
		},
	)
)

// Framework interfaces
var (
	_ resource.Resource                = &NotificationPolicyResource{}
	_ resource.ResourceWithConfigure   = &NotificationPolicyResource{}
	_ resource.ResourceWithImportState = &NotificationPolicyResource{}
	_ resource.ResourceWithModifyPlan  = &NotificationPolicyResource{}
)

// New Object
func NewNotificationPolicyResource() resource.Resource {
	return &NotificationPolicyResource{}
}

// Metadata
func (r *NotificationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_policy"
}

// Schema
func (r *NotificationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5
	const maxQuotaLimit = 2

	quotaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that define the SMS/Voice limits.  A maximum of two quota objects can be defined, one for SMS and/or Voice quota, and one for Email quota.",
	)

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean to provide an indication of whether this policy is the default notification policy for the environment. If the parameter is not provided, the value used is `false`.",
	)

	countryLimitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object to limit the countries where you can send SMS and voice notifications.",
	)

	countryLimitTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the kind of limitation being defined.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_NONE):    "no limitation is defined",
		string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_ALLOWED): "allows notifications only for the countries specified in the `countries` parameter",
		string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_DENIED):  "denies notifications only for the countries specified in the `countries` parameter",
	})

	countryLimitDeliveryMethodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The delivery methods that the defined limitation should be applied to. Content of the array can be `SMS`, `Voice`, or both. If the parameter is not provided, the default is `SMS` and `Voice`.",
	)

	countryLimitCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The countries where the specified methods should be allowed or denied. Use two-letter country codes from ISO 3166-1.  Required when `type` is not `%s`.", string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_NONE)),
	)

	quotaTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string to specify whether the limit defined is per-user or per environment.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsPolicyQuotaItemTypeEnumValues)

	quotaCountryLimitDeliveryMethodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The delivery methods for which the limit is being defined.",
	).AppendMarkdownString("This limits defined in this block are configured as two groups, Voice/SMS, or Email.  Email cannot be configured with Voice and/or SMS limits.").AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS):   fmt.Sprintf("configuration of SMS limits and can be set alongside `%s`, but not `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL)),
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE): fmt.Sprintf("configuration of Voice limits and can be set alongside `%s`, but not `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL)),
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL): fmt.Sprintf("configuration of Email limits but can not be set alongside `%s` or `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE)),
	}).DefaultValue(fmt.Sprintf(`["%s", "%s"]`, string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE)))

	quotaTotalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications allowed per day.  Cannot be set with `used` and `unused`.",
	)

	quotaUsedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications that can be received and responded to each day. Must be configured with `unused` and cannot be configured with `total`.",
	)

	quotaUnusedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications that can be received and not responded to each day. Must be configured with `used` and cannot be configured with `total`.",
	)

	providerConfigurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object to specify the custom notification providers to use for different countries and delivery methods (SMS and Voice).",
	)

	providerConfigurationConditionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of condition objects that define the provider fallback order to use for specific groups of countries and delivery methods. The **last condition** in the list must not have the `countries` field configured, which makes it serve as the default fallback order for all countries not specified in the preceding conditions.",
	)

	providerConfigurationConditionsDeliveryMethodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The delivery methods for which the fallback order should be applied.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsPolicyProviderConfigurationConditionsDeliveryMethodsEnumValues)

	providerConfigurationConditionsCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The countries for which the fallback order should be used. Use the two-letter country codes from ISO 3166-1.",
	)

	providerConfigurationConditionsFallbackChainDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of custom provider IDs in the order they should be used if available.",
	)

	providerConfigurationConditionsFallbackChainIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of a custom provider. Reference the `id` attribute of a `pingone_phone_delivery_settings` resource.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage notification policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to associate the notification policy with."),
			),

			"name": schema.StringAttribute{
				Description: "The name to use for the notification policy.  Must be unique among the notification policies in the environment.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"default": schema.BoolAttribute{
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Description:         defaultDescription.Description,
				Computed:            true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"country_limit": schema.SingleNestedAttribute{
				Description:         countryLimitDescription.Description,
				MarkdownDescription: countryLimitDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					countryLimitTFObjectTypes,
					map[string]attr.Value{
						"type":             types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_NONE)),
						"delivery_methods": types.SetNull(types.StringType),
						"countries":        types.SetNull(types.StringType),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         countryLimitTypeDescription.Description,
						MarkdownDescription: countryLimitTypeDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyCountryLimitTypeEnumValues)...),
						},
					},

					"delivery_methods": schema.SetAttribute{
						Description:         countryLimitDeliveryMethodsDescription.Description,
						MarkdownDescription: countryLimitDeliveryMethodsDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyCountryLimitDeliveryMethodEnumValues)...),
							),
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_NONE)),
								path.MatchRelative().AtParent().AtName("type"),
							),
						},
					},

					"countries": schema.SetAttribute{
						Description:         countryLimitCountriesDescription.Description,
						MarkdownDescription: countryLimitCountriesDescription.MarkdownDescription,
						Optional:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, "must be a valid two character country code"),
							),
							setvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_ALLOWED)),
								path.MatchRelative().AtParent().AtName("type"),
							),
							setvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_DENIED)),
								path.MatchRelative().AtParent().AtName("type"),
							),
						},
					},
				},
			},

			"cooldown_configuration": schema.SingleNestedAttribute{
				Description: "A single object to specify a period of time that users must wait before requesting an additional notification such as an additional OTP.",
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(cooldownConfigurationDefault),

				Attributes: map[string]schema.Attribute{
					"email": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for email notifications.",
						Required:    true,
						Attributes:  cooldownConfigurationMethodSchema(),
					},

					"sms": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for SMS notifications.",
						Required:    true,
						Attributes:  cooldownConfigurationMethodSchema(),
					},

					"voice": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for voice notifications.",
						Required:    true,
						Attributes:  cooldownConfigurationMethodSchema(),
					},

					"whats_app": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for WhatsApp notifications.",
						Required:    true,
						Attributes:  cooldownConfigurationMethodSchema(),
					},
				},
			},

			"provider_configuration": schema.SingleNestedAttribute{
				Description:         providerConfigurationDescription.Description,
				MarkdownDescription: providerConfigurationDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"conditions": schema.ListNestedAttribute{
						Description:         providerConfigurationConditionsDescription.Description,
						MarkdownDescription: providerConfigurationConditionsDescription.MarkdownDescription,
						Required:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"delivery_methods": schema.SetAttribute{
									Description:         providerConfigurationConditionsDeliveryMethodsDescription.Description,
									MarkdownDescription: providerConfigurationConditionsDeliveryMethodsDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.SizeAtLeast(attrMinLength),
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyProviderConfigurationConditionsDeliveryMethodsEnumValues)...),
										),
									},
								},

								"countries": schema.SetAttribute{
									Description:         providerConfigurationConditionsCountriesDescription.Description,
									MarkdownDescription: providerConfigurationConditionsCountriesDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.SizeAtLeast(attrMinLength),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, "must be a valid two character country code"),
										),
									},
								},

								"fallback_chain": schema.ListNestedAttribute{
									Description:         providerConfigurationConditionsFallbackChainDescription.Description,
									MarkdownDescription: providerConfigurationConditionsFallbackChainDescription.MarkdownDescription,
									Required:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         providerConfigurationConditionsFallbackChainIdDescription.Description,
												MarkdownDescription: providerConfigurationConditionsFallbackChainIdDescription.MarkdownDescription,
												Required:            true,

												CustomType: pingonetypes.ResourceIDType{},

												Validators: []validator.String{
													stringvalidator.LengthAtLeast(attrMinLength),
												},
											},
										},
									},

									Validators: []validator.List{
										listvalidator.SizeAtLeast(attrMinLength),
									},
								},
							},
						}, Validators: []validator.List{
							listvalidator.SizeAtLeast(attrMinLength),
						},
					},
				},
			},

			"quota": schema.SetNestedAttribute{
				Description:         quotaDescription.Description,
				MarkdownDescription: quotaDescription.MarkdownDescription,
				Optional:            true,

				NestedObject: schema.NestedAttributeObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         quotaTypeDescription.Description,
							MarkdownDescription: quotaTypeDescription.MarkdownDescription,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyQuotaItemTypeEnumValues)...),
							},
						},

						"delivery_methods": schema.SetAttribute{
							Description:         quotaCountryLimitDeliveryMethodsDescription.Description,
							MarkdownDescription: quotaCountryLimitDeliveryMethodsDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: setdefault.StaticValue(types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS)),
									types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE)),
								},
							)),

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.Any(
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE)),
									),
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL)),
									),
								),
							},
						},

						"total": schema.Int32Attribute{
							Description:         quotaTotalDescription.Description,
							MarkdownDescription: quotaTotalDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int32{
								int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("used")),
								int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("unused")),
							},
						},

						"used": schema.Int32Attribute{
							Description:         quotaUsedDescription.Description,
							MarkdownDescription: quotaUsedDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int32{
								int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("total")),
								int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("unused")),
							},
						},

						"unused": schema.Int32Attribute{
							Description:         quotaUnusedDescription.Description,
							MarkdownDescription: quotaUnusedDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int32{
								int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("total")),
								int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("used")),
							},
						},
					},
				},

				Validators: []validator.Set{
					setvalidator.SizeAtMost(maxQuotaLimit),
				},
			},
		},
	}
}

func cooldownConfigurationMethodSchema() map[string]schema.Attribute {
	const cooldownPeriodMinDurationSeconds = 10
	const cooldownPeriodMaxDurationSeconds = 600
	const cooldownPeriodMinDurationMinutes = 1
	const cooldownPeriodMaxDurationMinutes = 10
	const cooldownPeriodsArraySize = 3
	const cooldownResendLimitMin = 1
	const cooldownResendLimitMax = 10

	cooldownEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Set to `true` if you want to specify notification cooldown periods for the authentication method. Set to `false` if you don't want notification cooldown periods for this authentication method.",
	)

	cooldownPeriodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Use the periods array to specify the amount of time the user has to wait before requesting another notification such as another OTP. The array should contain three objects: the time to wait before the first retry, the time to wait before the second retry, and the time to wait before any subsequent retries. Required when `enabled` is `true`.",
	)

	cooldownPeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Used in conjunction with `time_unit` to specify the waiting period.",
	).AppendMarkdownString(fmt.Sprintf("If `time_unit` is `%s`, the allowed range is %d-%d.", string(management.ENUMNOTIFICATIONSPOLICYCOOLDOWNCONFIGURATIONMETHODPERIODTIMEUNIT_SECONDS), cooldownPeriodMinDurationSeconds, cooldownPeriodMaxDurationSeconds)).AppendMarkdownString(fmt.Sprintf("If `time_unit` is `%s`, the allowed range is %d-%d.", string(management.ENUMNOTIFICATIONSPOLICYCOOLDOWNCONFIGURATIONMETHODPERIODTIMEUNIT_MINUTES), cooldownPeriodMinDurationMinutes, cooldownPeriodMaxDurationMinutes))

	cooldownPeriodTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Used in conjunction with `duration` to specify the waiting period.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsPolicyCooldownConfigurationMethodPeriodTimeUnitEnumValues)

	cooldownGroupByDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"If you want the settings to be applied at the single-user level for the address/number, set this to `USER_ID`.",
	)

	cooldownResendLimitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of requests that a user can send to receive another notification, such as another OTP, before they are blocked for 30 minutes. Required when `enabled` is `true`.",
	)

	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Description:         cooldownEnabledDescription.Description,
			MarkdownDescription: cooldownEnabledDescription.MarkdownDescription,
			Required:            true,
		},

		"periods": schema.ListNestedAttribute{
			Description:         cooldownPeriodsDescription.Description,
			MarkdownDescription: cooldownPeriodsDescription.MarkdownDescription,
			Optional:            true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"duration": schema.Int32Attribute{
						Description:         cooldownPeriodDurationDescription.Description,
						MarkdownDescription: cooldownPeriodDurationDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.Int32{
							int32validator.Any(
								int32validator.All(
									int32validator.Between(cooldownPeriodMinDurationSeconds, cooldownPeriodMaxDurationSeconds),
									int32validatorinternal.RegexMatchesPathValue(
										regexp.MustCompile(`SECONDS`),
										fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", cooldownPeriodMinDurationSeconds, cooldownPeriodMaxDurationSeconds),
										path.MatchRelative().AtParent().AtName("time_unit"),
									),
								),
								int32validator.All(
									int32validator.Between(cooldownPeriodMinDurationMinutes, cooldownPeriodMaxDurationMinutes),
									int32validatorinternal.RegexMatchesPathValue(
										regexp.MustCompile(`MINUTES`),
										fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", cooldownPeriodMinDurationMinutes, cooldownPeriodMaxDurationMinutes),
										path.MatchRelative().AtParent().AtName("time_unit"),
									),
								),
							),
						},
					},

					"time_unit": schema.StringAttribute{
						Description:         cooldownPeriodTimeUnitDescription.Description,
						MarkdownDescription: cooldownPeriodTimeUnitDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyCooldownConfigurationMethodPeriodTimeUnitEnumValues)...),
						},
					},
				},
			},

			Validators: []validator.List{
				listvalidator.SizeBetween(cooldownPeriodsArraySize, cooldownPeriodsArraySize),
				listvalidatorinternal.IsRequiredIfMatchesPathBoolValue(
					types.BoolValue(true),
					path.MatchRelative().AtParent().AtName("enabled"),
				),
			},
		},

		"group_by": schema.StringAttribute{
			Description:         cooldownGroupByDescription.Description,
			MarkdownDescription: cooldownGroupByDescription.MarkdownDescription,
			Optional:            true,
		},

		"resend_limit": schema.Int32Attribute{
			Description:         cooldownResendLimitDescription.Description,
			MarkdownDescription: cooldownResendLimitDescription.MarkdownDescription,
			Optional:            true,
			Validators: []validator.Int32{
				int32validator.AtLeast(cooldownResendLimitMin),
				int32validator.AtMost(cooldownResendLimitMax),
				int32validatorinternal.IsRequiredIfMatchesPathBoolValue(
					types.BoolValue(true),
					path.MatchRelative().AtParent().AtName("enabled"),
				),
			},
		},
	}
} // ModifyPlan
func (r *NotificationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	var plan *NotificationPolicyCountryLimitResourceModel
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("country_limit"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan == nil {
		return
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() && plan.DeliveryMethods.IsUnknown() {

		if plan.Type.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSPOLICYCOUNTRYLIMITTYPE_NONE))) {
			resp.Plan.SetAttribute(ctx, path.Root("country_limit").AtName("delivery_methods"), types.SetNull(types.StringType))
		} else {
			setObj, d := types.SetValueFrom(ctx, types.StringType, []string{"Voice", "SMS"})
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Plan.SetAttribute(ctx, path.Root("country_limit").AtName("delivery_methods"), setObj)
		}
	}

}

func (r *NotificationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *NotificationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	notificationPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.CreateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateNotificationsPolicy",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NotificationPolicyResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var response *management.NotificationsPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.ReadOneNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneNotificationsPolicy",
		legacysdk.CustomErrorResourceNotFoundWarning,
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
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	notificationPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.UpdateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateNotificationsPolicy",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NotificationPolicyResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.DeleteNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteNotificationsPolicy",
		notificationPolicyDeleteCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

var notificationPolicyDeleteCustomError = func(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Undeletable default notifications policy
		if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
			if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
				if match, _ := regexp.MatchString("remove default notifications policy", v[0].GetMessage()); match {

					diags.AddWarning("Cannot delete the default notifications policy", "Due to API restrictions, the provider cannot delete the default notifications policy for an environment.  The policy has been removed from Terraform state but has been left in place in the PingOne service.")

					return diags
				}
			}
		}
	}

	diags.Append(legacysdk.CustomErrorResourceNotFoundWarning(r, p1Error)...)
	return diags
}

func (r *NotificationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "notification_policy_id",
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

func (p *NotificationPolicyResourceModel) expand(ctx context.Context) (*management.NotificationsPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	var quotaPlan []NotificationPolicyQuotaResourceModel
	diags.Append(p.Quota.ElementsAs(ctx, &quotaPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	quotas := make([]management.NotificationsPolicyQuotasInner, 0)
	for _, v := range quotaPlan {

		var deliveryMethodsPlan []management.EnumNotificationsPolicyQuotaDeliveryMethods
		diags.Append(v.DeliveryMethods.ElementsAs(ctx, &deliveryMethodsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		quota := *management.NewNotificationsPolicyQuotasInner(
			management.EnumNotificationsPolicyQuotaItemType(v.Type.ValueString()),
			deliveryMethodsPlan,
		)

		if !v.Total.IsNull() && !v.Total.IsUnknown() {
			quota.SetTotal(v.Total.ValueInt32())
		}

		if !v.Used.IsNull() && !v.Used.IsUnknown() {
			quota.SetClaimed(v.Used.ValueInt32())
		}

		if !v.Unused.IsNull() && !v.Unused.IsUnknown() {
			quota.SetUnclaimed(v.Unused.ValueInt32())
		}

		if management.EnumNotificationsPolicyQuotaItemType(v.Type.ValueString()) == management.ENUMNOTIFICATIONSPOLICYQUOTAITEMTYPE_USER &&
			(quota.GetTotal() > 50 || quota.GetClaimed() > 50 || quota.GetUnclaimed() > 50) {
			diags.AddError(
				"Invalid parameter",
				"User quota (parameters \"total\", \"used\" and \"unused\") for paid environment must be maximum of 50")
		}

		quotas = append(quotas, quota)
	}

	data := management.NewNotificationsPolicy(p.Name.ValueString(), quotas)

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	} else {
		data.SetDefault(false)
	}

	if !p.CountryLimit.IsNull() && !p.CountryLimit.IsUnknown() {
		var countryLimitPlan NotificationPolicyCountryLimitResourceModel
		diags.Append(p.CountryLimit.As(ctx, &countryLimitPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		var countriesPlan []types.String
		diags.Append(countryLimitPlan.Countries.ElementsAs(ctx, &countriesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		countries, d := framework.TFTypeStringSliceToStringSlice(countriesPlan, path.Root("country_limit").AtName("countries"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		countryLimit := *management.NewNotificationsPolicyCountryLimit(
			management.EnumNotificationsPolicyCountryLimitType(countryLimitPlan.Type.ValueString()),
			countries,
		)

		if !countryLimitPlan.DeliveryMethods.IsNull() && !countryLimitPlan.DeliveryMethods.IsUnknown() {
			var deliveryMethods []management.EnumNotificationsPolicyCountryLimitDeliveryMethod
			diags.Append(countryLimitPlan.DeliveryMethods.ElementsAs(ctx, &deliveryMethods, false)...)
			if diags.HasError() {
				return nil, diags
			}

			countryLimit.SetDeliveryMethods(deliveryMethods)
		}

		data.SetCountryLimit(countryLimit)
	}

	if !p.CooldownConfiguration.IsNull() && !p.CooldownConfiguration.IsUnknown() {
		var cooldownConfigPlan NotificationPolicyCooldownConfigurationResourceModel
		diags.Append(p.CooldownConfiguration.As(ctx, &cooldownConfigPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		cooldownConfig, d := expandCooldownConfiguration(ctx, &cooldownConfigPlan)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCooldownConfiguration(*cooldownConfig)
	}

	if !p.ProviderConfiguration.IsNull() && !p.ProviderConfiguration.IsUnknown() {
		var providerConfigPlan NotificationPolicyProviderConfigurationResourceModel
		diags.Append(p.ProviderConfiguration.As(ctx, &providerConfigPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		providerConfig, d := expandProviderConfiguration(ctx, &providerConfigPlan)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProviderConfiguration(*providerConfig)
	}

	return data, diags
}

func expandCooldownConfiguration(ctx context.Context, plan *NotificationPolicyCooldownConfigurationResourceModel) (*management.NotificationsPolicyCooldownConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	email, d := expandCooldownConfigurationMethod(ctx, plan.Email)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	sms, d := expandCooldownConfigurationMethod(ctx, plan.Sms)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	voice, d := expandCooldownConfigurationMethod(ctx, plan.Voice)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	whatsApp, d := expandCooldownConfigurationMethod(ctx, plan.WhatsApp)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	cooldownConfig := management.NewNotificationsPolicyCooldownConfiguration(*email, *sms, *voice, *whatsApp)

	return cooldownConfig, diags
}

func expandCooldownConfigurationMethod(ctx context.Context, methodObj types.Object) (*management.NotificationsPolicyCooldownConfigurationMethod, diag.Diagnostics) {
	var diags diag.Diagnostics

	var methodPlan NotificationPolicyCooldownConfigurationMethodResourceModel
	diags.Append(methodObj.As(ctx, &methodPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	method := &management.NotificationsPolicyCooldownConfigurationMethod{
		Enabled: methodPlan.Enabled.ValueBool(),
	}

	if methodPlan.Enabled.ValueBool() {
		var periodsPlan []NotificationPolicyCooldownConfigurationMethodPeriodResourceModel
		diags.Append(methodPlan.Periods.ElementsAs(ctx, &periodsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		periods := make([]management.NotificationsPolicyCooldownConfigurationMethodPeriodsInner, 0, len(periodsPlan))
		for _, p := range periodsPlan {
			period := management.NewNotificationsPolicyCooldownConfigurationMethodPeriodsInner(
				p.Duration.ValueInt32(),
				management.EnumNotificationsPolicyCooldownConfigurationMethodPeriodTimeUnit(p.TimeUnit.ValueString()),
			)
			periods = append(periods, *period)
		}
		method.Periods = periods

		if !methodPlan.ResendLimit.IsNull() && !methodPlan.ResendLimit.IsUnknown() {
			method.SetResendLimit(methodPlan.ResendLimit.ValueInt32())
		}
	}

	if !methodPlan.GroupBy.IsNull() && !methodPlan.GroupBy.IsUnknown() {
		method.SetGroupBy(methodPlan.GroupBy.ValueString())
	}

	return method, diags
}

func expandProviderConfiguration(ctx context.Context, plan *NotificationPolicyProviderConfigurationResourceModel) (*management.NotificationsPolicyProviderConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	providerConfig := management.NewNotificationsPolicyProviderConfiguration()

	if !plan.Conditions.IsNull() && !plan.Conditions.IsUnknown() {
		var conditionsPlan []NotificationPolicyProviderConfigurationConditionResourceModel
		diags.Append(plan.Conditions.ElementsAs(ctx, &conditionsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		conditions := make([]management.NotificationsPolicyProviderConfigurationConditionsInner, 0, len(conditionsPlan))
		for _, c := range conditionsPlan {
			condition := management.NewNotificationsPolicyProviderConfigurationConditionsInner()

			// Delivery methods
			if !c.DeliveryMethods.IsNull() && !c.DeliveryMethods.IsUnknown() {
				var deliveryMethodsPlan []types.String
				diags.Append(c.DeliveryMethods.ElementsAs(ctx, &deliveryMethodsPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				deliveryMethods := make([]management.EnumNotificationsPolicyProviderConfigurationConditionsDeliveryMethods, 0, len(deliveryMethodsPlan))
				for _, dm := range deliveryMethodsPlan {
					deliveryMethods = append(deliveryMethods, management.EnumNotificationsPolicyProviderConfigurationConditionsDeliveryMethods(dm.ValueString()))
				}
				condition.SetDeliveryMethods(deliveryMethods)
			}

			// Countries
			if !c.Countries.IsNull() && !c.Countries.IsUnknown() {
				var countriesPlan []types.String
				diags.Append(c.Countries.ElementsAs(ctx, &countriesPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				countries, d := framework.TFTypeStringSliceToStringSlice(countriesPlan, path.Root("provider_configuration").AtName("conditions").AtListIndex(0).AtName("countries"))
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				condition.SetCountries(countries)
			}

			// Fallback chain
			if !c.FallbackChain.IsNull() && !c.FallbackChain.IsUnknown() {
				var fallbackChainPlan []NotificationPolicyProviderConfigurationFallbackChainItemResourceModel
				diags.Append(c.FallbackChain.ElementsAs(ctx, &fallbackChainPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				fallbackChain := make([]management.NotificationsPolicyProviderConfigurationConditionsInnerFallbackChainInner, 0, len(fallbackChainPlan))
				for _, fb := range fallbackChainPlan {
					fallbackItem := management.NewNotificationsPolicyProviderConfigurationConditionsInnerFallbackChainInner(fb.Id.ValueString())
					fallbackChain = append(fallbackChain, *fallbackItem)
				}
				condition.SetFallbackChain(fallbackChain)
			}

			conditions = append(conditions, *condition)
		}

		providerConfig.SetConditions(conditions)
	}

	return providerConfig, diags
}

func (p *NotificationPolicyResourceModel) toState(apiObject *management.NotificationsPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())

	var d diag.Diagnostics

	p.Quota, d = toStateQuota(apiObject.GetQuotas())
	diags.Append(d...)

	p.CountryLimit, d = toStateCountryLimit(apiObject.GetCountryLimitOk())
	diags.Append(d...)

	p.CooldownConfiguration, d = toStateCooldownConfiguration(apiObject.GetCooldownConfigurationOk())
	diags.Append(d...)

	p.ProviderConfiguration, d = toStateProviderConfiguration(apiObject.GetProviderConfigurationOk())
	diags.Append(d...)

	return diags
}

func toStateQuota(quotas []management.NotificationsPolicyQuotasInner) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: quotaTFObjectTypes}

	if len(quotas) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range quotas {

		quota := map[string]attr.Value{
			"type":             framework.EnumOkToTF(v.GetTypeOk()),
			"delivery_methods": framework.EnumSetOkToTF(v.GetDeliveryMethodsOk()),
			"total":            framework.Int32OkToTF(v.GetTotalOk()),
			"used":             framework.Int32OkToTF(v.GetClaimedOk()),
			"unused":           framework.Int32OkToTF(v.GetUnclaimedOk()),
		}

		flattenedObj, d := types.ObjectValue(quotaTFObjectTypes, quota)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}

func toStateCountryLimit(apiObject *management.NotificationsPolicyCountryLimit, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(countryLimitTFObjectTypes), diags
	}

	countryLimitMap := map[string]attr.Value{
		"type":             framework.EnumOkToTF(apiObject.GetTypeOk()),
		"delivery_methods": framework.EnumSetOkToTF(apiObject.GetDeliveryMethodsOk()),
		"countries":        framework.StringSetOkToTF(apiObject.GetCountriesOk()),
	}

	returnVar, d := types.ObjectValue(countryLimitTFObjectTypes, countryLimitMap)
	diags.Append(d...)

	return returnVar, diags

}

func toStateCooldownConfiguration(apiObject *management.NotificationsPolicyCooldownConfiguration, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(cooldownConfigurationTFObjectTypes), diags
	}

	email, d := toStateCooldownConfigurationMethod(apiObject.GetEmailOk())
	diags.Append(d...)

	sms, d := toStateCooldownConfigurationMethod(apiObject.GetSmsOk())
	diags.Append(d...)

	voice, d := toStateCooldownConfigurationMethod(apiObject.GetVoiceOk())
	diags.Append(d...)

	whatsApp, d := toStateCooldownConfigurationMethod(apiObject.GetWhatsAppOk())
	diags.Append(d...)

	cooldownConfigMap := map[string]attr.Value{
		"email":     email,
		"sms":       sms,
		"voice":     voice,
		"whats_app": whatsApp,
	}

	returnVar, d := types.ObjectValue(cooldownConfigurationTFObjectTypes, cooldownConfigMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateCooldownConfigurationMethod(apiObject *management.NotificationsPolicyCooldownConfigurationMethod, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(cooldownConfigurationMethodTFObjectTypes), diags
	}

	periods, d := toStateCooldownConfigurationMethodPeriods(apiObject.GetPeriodsOk())
	diags.Append(d...)

	methodMap := map[string]attr.Value{
		"enabled":      framework.BoolOkToTF(apiObject.GetEnabledOk()),
		"periods":      periods,
		"group_by":     framework.StringOkToTF(apiObject.GetGroupByOk()),
		"resend_limit": framework.Int32OkToTF(apiObject.GetResendLimitOk()),
	}

	returnVar, d := types.ObjectValue(cooldownConfigurationMethodTFObjectTypes, methodMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateCooldownConfigurationMethodPeriods(periods []management.NotificationsPolicyCooldownConfigurationMethodPeriodsInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: cooldownConfigurationPeriodTFObjectTypes}

	if !ok || len(periods) == 0 {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, p := range periods {
		periodMap := map[string]attr.Value{
			"duration":  framework.Int32OkToTF(p.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(p.GetTimeUnitOk()),
		}

		flattenedObj, d := types.ObjectValue(cooldownConfigurationPeriodTFObjectTypes, periodMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateProviderConfiguration(apiObject *management.NotificationsPolicyProviderConfiguration, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(providerConfigurationTFObjectTypes), diags
	}

	conditions, d := toStateProviderConfigurationConditions(apiObject.GetConditionsOk())
	diags.Append(d...)

	objMap := map[string]attr.Value{
		"conditions": conditions,
	}

	objValue, d := types.ObjectValue(providerConfigurationTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func toStateProviderConfigurationConditions(conditions []management.NotificationsPolicyProviderConfigurationConditionsInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: providerConfigurationConditionTFObjectTypes}

	if !ok || len(conditions) == 0 {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, c := range conditions {
		deliveryMethods, d := toStateProviderConfigurationDeliveryMethods(c.GetDeliveryMethodsOk())
		diags.Append(d...)

		countries, d := toStateProviderConfigurationCountries(c.GetCountriesOk())
		diags.Append(d...)

		fallbackChain, d := toStateProviderConfigurationFallbackChain(c.GetFallbackChainOk())
		diags.Append(d...)

		conditionMap := map[string]attr.Value{
			"delivery_methods": deliveryMethods,
			"countries":        countries,
			"fallback_chain":   fallbackChain,
		}

		flattenedObj, d := types.ObjectValue(providerConfigurationConditionTFObjectTypes, conditionMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateProviderConfigurationDeliveryMethods(deliveryMethods []management.EnumNotificationsPolicyProviderConfigurationConditionsDeliveryMethods, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || len(deliveryMethods) == 0 {
		return types.SetNull(types.StringType), diags
	}

	flattenedList := []attr.Value{}
	for _, dm := range deliveryMethods {
		flattenedList = append(flattenedList, types.StringValue(string(dm)))
	}

	returnVar, d := types.SetValue(types.StringType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateProviderConfigurationCountries(countries []string, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || len(countries) == 0 {
		return types.SetNull(types.StringType), diags
	}

	flattenedList := []attr.Value{}
	for _, country := range countries {
		flattenedList = append(flattenedList, types.StringValue(country))
	}

	returnVar, d := types.SetValue(types.StringType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateProviderConfigurationFallbackChain(fallbackChain []management.NotificationsPolicyProviderConfigurationConditionsInnerFallbackChainInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: fallbackChainItemTFObjectTypes}

	if !ok || len(fallbackChain) == 0 {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, fb := range fallbackChain {
		fallbackMap := map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(fb.GetIdOk()),
		}

		flattenedObj, d := types.ObjectValue(fallbackChainItemTFObjectTypes, fallbackMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}
