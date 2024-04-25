package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationPolicyResource serviceClientType

type NotificationPolicyResourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Default       types.Bool                   `tfsdk:"default"`
	CountryLimit  types.Object                 `tfsdk:"country_limit"`
	Quota         types.Set                    `tfsdk:"quota"`
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type NotificationPolicyQuotaResourceModel struct {
	Type            types.String `tfsdk:"type"`
	DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Total           types.Int64  `tfsdk:"total"`
	Used            types.Int64  `tfsdk:"used"`
	Unused          types.Int64  `tfsdk:"unused"`
}

type NotificationPolicyCountryLimitResourceModel struct {
	Type            types.String `tfsdk:"type"`
	DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Countries       types.Set    `tfsdk:"countries"`
}

var (
	quotaTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
		"delivery_methods": types.SetType{
			ElemType: types.StringType,
		},
		"total":  types.Int64Type,
		"used":   types.Int64Type,
		"unused": types.Int64Type,
	}

	countryLimitTFObjectTypes = map[string]attr.Type{
		"type":             types.StringType,
		"delivery_methods": types.SetType{ElemType: types.StringType},
		"countries":        types.SetType{ElemType: types.StringType},
	}
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

						"total": schema.Int64Attribute{
							Description:         quotaTotalDescription.Description,
							MarkdownDescription: quotaTotalDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("used")),
								int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("unused")),
							},
						},

						"used": schema.Int64Attribute{
							Description:         quotaUsedDescription.Description,
							MarkdownDescription: quotaUsedDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("total")),
								int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("unused")),
							},
						},

						"unused": schema.Int64Attribute{
							Description:         quotaUnusedDescription.Description,
							MarkdownDescription: quotaUnusedDescription.MarkdownDescription,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("total")),
								int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("used")),
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

// ModifyPlan
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

func (r *NotificationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.CreateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateNotificationsPolicy",
		framework.DefaultCustomError,
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

	if r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.ReadOneNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneNotificationsPolicy",
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
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.UpdateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateNotificationsPolicy",
		framework.DefaultCustomError,
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

	if r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.NotificationsPoliciesApi.DeleteNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteNotificationsPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
			quota.SetTotal(int32(v.Total.ValueInt64()))
		}

		if !v.Used.IsNull() && !v.Used.IsUnknown() {
			quota.SetClaimed(int32(v.Used.ValueInt64()))
		}

		if !v.Unused.IsNull() && !v.Unused.IsUnknown() {
			quota.SetUnclaimed(int32(v.Unused.ValueInt64()))
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

		var countries []string
		diags.Append(countryLimitPlan.Countries.ElementsAs(ctx, &countries, false)...)
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

	return data, diags
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
