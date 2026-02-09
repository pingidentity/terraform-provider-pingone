// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type NotificationPolicyDataSource serviceClientType

type NotificationPolicyDataSourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	NotificationPolicyId  pingonetypes.ResourceIDValue `tfsdk:"notification_policy_id"`
	Name                  types.String                 `tfsdk:"name"`
	Default               types.Bool                   `tfsdk:"default"`
	CountryLimit          types.Object                 `tfsdk:"country_limit"`
	CooldownConfiguration types.Object                 `tfsdk:"cooldown_configuration"`
	ProviderConfiguration types.Object                 `tfsdk:"provider_configuration"`
	Quota                 types.Set                    `tfsdk:"quota"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &NotificationPolicyDataSource{}
)

// New Object
func NewNotificationPolicyDataSource() datasource.DataSource {
	return &NotificationPolicyDataSource{}
}

// Metadata
func (r *NotificationPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_policy"
}

// Schema
func (r *NotificationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const attrMinLength = 1

	dataSourceExactlyOneOfRelativePaths := []string{
		"notification_policy_id",
		"name",
		"default",
	}

	notificationPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the notification policy to retrieve configuration for.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	notificationPolicyNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the notification policy to retrieve configuration for.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Set value to `true` to return the default notification policy. There is only one default policy per environment.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

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
		"The delivery methods that the defined limitation should be applied to. Content of the array can be `SMS`, `Voice`, or both.",
	)

	countryLimitCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The countries where the specified methods should be allowed or denied. Use two-letter country codes from ISO 3166-1.",
	)

	quotaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that define the SMS/Voice limits.  A maximum of two quota objects can be defined, one for SMS and/or Voice quota, and one for Email quota.",
	)

	quotaTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string to specify whether the limit defined is per-user or per environment.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsPolicyQuotaItemTypeEnumValues)

	quotaDeliveryMethodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The delivery methods for which the limit is being defined.",
	).AppendMarkdownString("This limits defined in this block are configured as two groups, Voice/SMS, or Email.  Email cannot be configured with Voice and/or SMS limits.").AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS):   fmt.Sprintf("configuration of SMS limits and can be set alongside `%s`, but not `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL)),
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE): fmt.Sprintf("configuration of Voice limits and can be set alongside `%s`, but not `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL)),
		string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_EMAIL): fmt.Sprintf("configuration of Email limits but can not be set alongside `%s` or `%s`", string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS), string(management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE)),
	})

	quotaTotalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications allowed per day.",
	)

	quotaUsedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications that can be received and responded to each day.",
	)

	quotaUnusedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of notifications that can be received and not responded to each day.",
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
		"The ID of a custom provider.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne notification policy in an environment by ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the notification policy.  Must be a valid PingOne resource ID.").Description,
				Required:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"notification_policy_id": schema.StringAttribute{
				Description:         notificationPolicyIdDescription.Description,
				MarkdownDescription: notificationPolicyIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("default"),
					),
				},
			},

			"name": schema.StringAttribute{
				Description:         notificationPolicyNameDescription.Description,
				MarkdownDescription: notificationPolicyNameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"default": schema.BoolAttribute{
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Description:         defaultDescription.Description,
				Optional:            true,
				Computed:            true,
			},

			"country_limit": schema.SingleNestedAttribute{
				Description:         countryLimitDescription.Description,
				MarkdownDescription: countryLimitDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         countryLimitTypeDescription.Description,
						MarkdownDescription: countryLimitTypeDescription.MarkdownDescription,
						Computed:            true,
					},

					"delivery_methods": schema.SetAttribute{
						Description:         countryLimitDeliveryMethodsDescription.Description,
						MarkdownDescription: countryLimitDeliveryMethodsDescription.MarkdownDescription,
						Computed:            true,

						ElementType: types.StringType,
					},

					"countries": schema.SetAttribute{
						Description:         countryLimitCountriesDescription.Description,
						MarkdownDescription: countryLimitCountriesDescription.MarkdownDescription,
						Computed:            true,

						ElementType: types.StringType,
					},
				},
			},

			"cooldown_configuration": schema.SingleNestedAttribute{
				Description: "A single object to specify a period of time that users must wait before requesting an additional notification such as an additional OTP.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"email": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for email notifications.",
						Computed:    true,
						Attributes:  dataSourceCooldownConfigurationMethodSchema(),
					},

					"sms": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for SMS notifications.",
						Computed:    true,
						Attributes:  dataSourceCooldownConfigurationMethodSchema(),
					},

					"voice": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for voice notifications.",
						Computed:    true,
						Attributes:  dataSourceCooldownConfigurationMethodSchema(),
					},

					"whats_app": schema.SingleNestedAttribute{
						Description: "Contains the notification cooldown period settings for WhatsApp notifications.",
						Computed:    true,
						Attributes:  dataSourceCooldownConfigurationMethodSchema(),
					},
				},
			},

			"provider_configuration": schema.SingleNestedAttribute{
				Description:         providerConfigurationDescription.Description,
				MarkdownDescription: providerConfigurationDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"conditions": schema.ListNestedAttribute{
						Description:         providerConfigurationConditionsDescription.Description,
						MarkdownDescription: providerConfigurationConditionsDescription.MarkdownDescription,
						Computed:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"delivery_methods": schema.SetAttribute{
									Description:         providerConfigurationConditionsDeliveryMethodsDescription.Description,
									MarkdownDescription: providerConfigurationConditionsDeliveryMethodsDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"countries": schema.SetAttribute{
									Description:         providerConfigurationConditionsCountriesDescription.Description,
									MarkdownDescription: providerConfigurationConditionsCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"fallback_chain": schema.ListNestedAttribute{
									Description:         providerConfigurationConditionsFallbackChainDescription.Description,
									MarkdownDescription: providerConfigurationConditionsFallbackChainDescription.MarkdownDescription,
									Computed:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         providerConfigurationConditionsFallbackChainIdDescription.Description,
												MarkdownDescription: providerConfigurationConditionsFallbackChainIdDescription.MarkdownDescription,
												Computed:            true,

												CustomType: pingonetypes.ResourceIDType{},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"quota": schema.SetNestedAttribute{
				Description:         quotaDescription.Description,
				MarkdownDescription: quotaDescription.MarkdownDescription,
				Computed:            true,

				NestedObject: schema.NestedAttributeObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         quotaTypeDescription.Description,
							MarkdownDescription: quotaTypeDescription.MarkdownDescription,
							Computed:            true,
						},

						"delivery_methods": schema.SetAttribute{
							Description:         quotaDeliveryMethodsDescription.Description,
							MarkdownDescription: quotaDeliveryMethodsDescription.MarkdownDescription,
							Computed:            true,

							ElementType: types.StringType,
						},

						"total": schema.Int32Attribute{
							Description:         quotaTotalDescription.Description,
							MarkdownDescription: quotaTotalDescription.MarkdownDescription,
							Computed:            true,
						},

						"used": schema.Int32Attribute{
							Description:         quotaUsedDescription.Description,
							MarkdownDescription: quotaUsedDescription.MarkdownDescription,
							Computed:            true,
						},

						"unused": schema.Int32Attribute{
							Description:         quotaUnusedDescription.Description,
							MarkdownDescription: quotaUnusedDescription.MarkdownDescription,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCooldownConfigurationMethodSchema() map[string]schema.Attribute {

	cooldownEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Set to `true` if notification cooldown periods are configured for the authentication method. Set to `false` if notification cooldown periods are not configured for this authentication method.",
	)

	cooldownPeriodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array of objects that specify the amount of time the user has to wait before requesting another notification such as another OTP. The array contains three objects: the time to wait before the first retry, the time to wait before the second retry, and the time to wait before any subsequent retries.",
	)

	cooldownPeriodDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Used in conjunction with `time_unit` to specify the waiting period.",
	)

	cooldownPeriodTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Used in conjunction with `duration` to specify the waiting period.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsPolicyCooldownConfigurationMethodPeriodTimeUnitEnumValues)

	cooldownGroupByDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"If the settings are applied at the single-user level for the address/number, this value is set to `USER_ID`.",
	)

	cooldownResendLimitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The maximum number of requests that a user can send to receive another notification, such as another OTP, before they are blocked for 30 minutes.",
	)

	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Description:         cooldownEnabledDescription.Description,
			MarkdownDescription: cooldownEnabledDescription.MarkdownDescription,
			Computed:            true,
		},

		"periods": schema.ListNestedAttribute{
			Description:         cooldownPeriodsDescription.Description,
			MarkdownDescription: cooldownPeriodsDescription.MarkdownDescription,
			Computed:            true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"duration": schema.Int32Attribute{
						Description:         cooldownPeriodDurationDescription.Description,
						MarkdownDescription: cooldownPeriodDurationDescription.MarkdownDescription,
						Computed:            true,
					},

					"time_unit": schema.StringAttribute{
						Description:         cooldownPeriodTimeUnitDescription.Description,
						MarkdownDescription: cooldownPeriodTimeUnitDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},
		},

		"group_by": schema.StringAttribute{
			Description:         cooldownGroupByDescription.Description,
			MarkdownDescription: cooldownGroupByDescription.MarkdownDescription,
			Computed:            true,
		},

		"resend_limit": schema.Int32Attribute{
			Description:         cooldownResendLimitDescription.Description,
			MarkdownDescription: cooldownResendLimitDescription.MarkdownDescription,
			Computed:            true,
		},
	}
}

func (r *NotificationPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NotificationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NotificationPolicyDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var notificationPolicy *management.NotificationsPolicy

	if !data.NotificationPolicyId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.ReadOneNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.NotificationPolicyId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneNotificationsPolicy",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&notificationPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.NotificationsPoliciesApi.ReadAllNotificationsPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if notificationPolicies, ok := pageCursor.EntityArray.Embedded.GetNotificationsPoliciesOk(); ok {
						for _, notificationPolicy := range notificationPolicies {
							if strings.EqualFold(notificationPolicy.GetName(), data.Name.ValueString()) {
								return &notificationPolicy, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllNotificationsPolicies",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&notificationPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if notificationPolicy == nil {
			resp.Diagnostics.AddError(
				"Cannot find notification policy from name",
				fmt.Sprintf("The notification policy name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else if data.Default.ValueBool() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.NotificationsPoliciesApi.ReadAllNotificationsPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if notificationPolicies, ok := pageCursor.EntityArray.Embedded.GetNotificationsPoliciesOk(); ok {
						for _, notificationPolicy := range notificationPolicies {
							if notificationPolicy.GetDefault() {
								return &notificationPolicy, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllNotificationsPolicies",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&notificationPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if notificationPolicy == nil {
			resp.Diagnostics.AddError(
				"Cannot find default notification policy",
				fmt.Sprintf("The default notification policy for environment %s cannot be found", data.EnvironmentId.String()),
			)
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Notification Policy: notification_policy_id, name, or default argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(notificationPolicy)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *NotificationPolicyDataSourceModel) toState(apiObject *management.NotificationsPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.NotificationPolicyId = framework.PingOneResourceIDToTF(apiObject.GetId())
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
