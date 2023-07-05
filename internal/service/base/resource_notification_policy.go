package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationPolicyResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type NotificationPolicyResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	Default       types.Bool   `tfsdk:"default"`
	CountryLimit  types.Object `tfsdk:"country_limit"`
	Quota         types.List   `tfsdk:"quota"`
	Id            types.String `tfsdk:"id"`
}

type NotificationPolicyQuotaResourceModel struct {
	Type types.String `tfsdk:"type"`
	// To enable when the platform supports individual configuration
	// DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Total  types.Int64 `tfsdk:"total"`
	Used   types.Int64 `tfsdk:"used"`
	Unused types.Int64 `tfsdk:"unused"`
}

type NotificationPolicyCountryLimitResourceModel struct {
	Type            types.String `tfsdk:"type"`
	DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Countries       types.Set    `tfsdk:"countries"`
}

var (
	quotaTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
		// To enable when the platform supports individual configuration
		// "delivery_method": types.SetType{
		// 	ElemType: types.StringType,
		// },
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

	quotaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object block that define the SMS/Voice limits.",
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

	quotaTotalDescriptionFmt := "The maximum number of notifications allowed per day.  Cannot be set with `used` and `unused`."
	quotaTotalDescription := framework.SchemaAttributeDescription{

		MarkdownDescription: quotaTotalDescriptionFmt,
		Description:         strings.ReplaceAll(quotaTotalDescriptionFmt, "`", "\""),
	}

	quotaUsedDescriptionFmt := "The maximum number of notifications that can be received and responded to each day. Must be configured with `unused` and cannot be configured with `total`."
	quotaUsedDescription := framework.SchemaAttributeDescription{

		MarkdownDescription: quotaUsedDescriptionFmt,
		Description:         strings.ReplaceAll(quotaUsedDescriptionFmt, "`", "\""),
	}

	quotaUnusedDescriptionFmt := "The maximum number of notifications that can be received and not responded to each day. Must be configured with `used` and cannot be configured with `total`."
	quotaUnusedDescription := framework.SchemaAttributeDescription{

		MarkdownDescription: quotaUnusedDescriptionFmt,
		Description:         strings.ReplaceAll(quotaUnusedDescriptionFmt, "`", "\""),
	}

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
		},

		Blocks: map[string]schema.Block{
			"quota": schema.ListNestedBlock{
				Description:         quotaDescription.Description,
				MarkdownDescription: quotaDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         quotaTypeDescription.Description,
							MarkdownDescription: quotaTypeDescription.MarkdownDescription,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsPolicyQuotaItemTypeEnumValues)...),
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

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
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

func (r *NotificationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.client == nil {
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
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsPoliciesApi.CreateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
		},
		"CreateNotificationsPolicy",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NotificationPolicyResourceModel

	if r.client == nil {
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
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsPoliciesApi.ReadOneNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneNotificationsPolicy",
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
	resp.Diagnostics.Append(data.toState(response.(*management.NotificationsPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationPolicyResourceModel

	if r.client == nil {
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
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsPoliciesApi.UpdateNotificationsPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).NotificationsPolicy(*notificationPolicy).Execute()
		},
		"UpdateNotificationsPolicy",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NotificationPolicyResourceModel

	if r.client == nil {
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
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.NotificationsPoliciesApi.DeleteNotificationsPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteNotificationsPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/notification_policy_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
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
		quota := *management.NewNotificationsPolicyQuotasInner(
			management.EnumNotificationsPolicyQuotaItemType(v.Type.ValueString()),

			// These are always set this way, otherwise the platform will reject
			[]management.EnumNotificationsPolicyQuotaDeliveryMethods{management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_SMS, management.ENUMNOTIFICATIONSPOLICYQUOTADELIVERYMETHODS_VOICE},
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
	data.SetDefault(false)

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

	p.Id = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())

	var d diag.Diagnostics

	p.Quota, d = toStateQuota(apiObject.GetQuotas())
	diags.Append(d...)

	p.CountryLimit, d = toStateCountryLimit(apiObject.GetCountryLimitOk())
	diags.Append(d...)

	return diags
}

func toStateQuota(quotas []management.NotificationsPolicyQuotasInner) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: quotaTFObjectTypes}

	if len(quotas) == 0 {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range quotas {

		// To enable when the platform supports individual configuration
		// deliveryMethods, d := framework.EnumSetOkToTF(v.GetDeliveryMethodsOk()))
		// diags.Append(d...)

		quota := map[string]attr.Value{
			"type":   framework.EnumOkToTF(v.GetTypeOk()),
			"total":  framework.Int32OkToTF(v.GetTotalOk()),
			"used":   framework.Int32OkToTF(v.GetClaimedOk()),
			"unused": framework.Int32OkToTF(v.GetUnclaimedOk()),
		}

		// To enable when the platform supports individual configuration
		// "delivery_method": deliveryMethods,

		flattenedObj, d := types.ObjectValue(quotaTFObjectTypes, quota)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
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
