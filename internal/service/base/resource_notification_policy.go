package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
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
	Quota         types.List   `tfsdk:"quota"`
	Id            types.String `tfsdk:"id"`
}

type QuotaModel struct {
	Type types.String `tfsdk:"type"`
	// To enable when the platform supports individual configuration
	// DeliveryMethods types.Set    `tfsdk:"delivery_methods"`
	Total  types.Int64 `tfsdk:"total"`
	Used   types.Int64 `tfsdk:"used"`
	Unused types.Int64 `tfsdk:"unused"`
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

	quotaDescriptionFmt := "A single object block that define the SMS/Voice limits."
	quotaDescription := framework.SchemaDescription{
		MarkdownDescription: quotaDescriptionFmt,
		Description:         strings.ReplaceAll(quotaDescriptionFmt, "`", "\""),
	}

	defaultDescriptionFmt := "A boolean to provide an indication of whether this policy is the default notification policy for the environment. If the parameter is not provided, the value used is `false`."
	defaultDescription := framework.SchemaDescription{
		MarkdownDescription: defaultDescriptionFmt,
		Description:         strings.ReplaceAll(defaultDescriptionFmt, "`", "\""),
	}

	quotaTypeAllowedValues := make([]string, 0)
	for _, v := range management.AllowedEnumNotificationsPolicyQuotaItemTypeEnumValues {
		quotaTypeAllowedValues = append(quotaTypeAllowedValues, string(v))
	}

	quotaTypeDescriptionFmt := fmt.Sprintf("A string to specify whether the limit defined is per-user or per environment. Allowed values: `%s`.", strings.Join(quotaTypeAllowedValues, "`, `"))
	quotaTypeDescription := framework.SchemaDescription{
		MarkdownDescription: quotaTypeDescriptionFmt,
		Description:         strings.ReplaceAll(quotaTypeDescriptionFmt, "`", "\""),
	}

	quotaTotalDescriptionFmt := "The maximum number of notifications allowed per day.  Cannot be set with `used` and `unused`."
	quotaTotalDescription := framework.SchemaDescription{

		MarkdownDescription: quotaTotalDescriptionFmt,
		Description:         strings.ReplaceAll(quotaTotalDescriptionFmt, "`", "\""),
	}

	quotaUsedDescriptionFmt := "The maximum number of notifications that can be received and responded to each day. Must be configured with `unused` and cannot be configured with `total`."
	quotaUsedDescription := framework.SchemaDescription{

		MarkdownDescription: quotaUsedDescriptionFmt,
		Description:         strings.ReplaceAll(quotaUsedDescriptionFmt, "`", "\""),
	}

	quotaUnusedDescriptionFmt := "The maximum number of notifications that can be received and not responded to each day. Must be configured with `used` and cannot be configured with `total`."
	quotaUnusedDescription := framework.SchemaDescription{

		MarkdownDescription: quotaUnusedDescriptionFmt,
		Description:         strings.ReplaceAll(quotaUnusedDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage notification policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_EnvironmentID(framework.SchemaDescription{
				Description: "The ID of the environment to associate the notification policy with."},
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
								stringvalidator.OneOf(quotaTypeAllowedValues...),
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

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

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

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

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

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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

	var quotaPlan []QuotaModel
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

	return data, diags
}

func (p *NotificationPolicyResourceModel) toState(v *management.NotificationsPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = types.StringValue(v.GetId())
	p.Name = types.StringValue(v.GetName())
	p.Default = types.BoolValue(v.GetDefault())

	quota, d := toStateQuota(v.GetQuotas())
	diags.Append(d...)
	p.Quota = quota

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
		// deliveryMethods, d := framework.StringSliceToTF(deliveryMethodsToStringSlice(v.GetDeliveryMethods()))
		// diags.Append(d...)

		quota := map[string]attr.Value{
			"type": framework.StringToTF(string(v.GetType())),
			// To enable when the platform supports individual configuration
			// "delivery_method": deliveryMethods,
		}

		if i, ok := v.GetTotalOk(); ok {
			quota["total"] = framework.Int32ToTF(*i)
		} else {
			quota["total"] = types.Int64Null()
		}

		if i, ok := v.GetClaimedOk(); ok {
			quota["used"] = framework.Int32ToTF(*i)
		} else {
			quota["used"] = types.Int64Null()
		}

		if i, ok := v.GetUnclaimedOk(); ok {
			quota["unused"] = framework.Int32ToTF(*i)
		} else {
			quota["unused"] = types.Int64Null()
		}

		flattenedObj, d := types.ObjectValue(quotaTFObjectTypes, quota)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}

// To enable when the platform supports individual configuration
// func deliveryMethodsToStringSlice(methods []management.EnumNotificationsPolicyQuotaDeliveryMethods) []string {

// 	returnSlice := make([]string, 0)
// 	for _, v := range methods {
// 		returnSlice = append(returnSlice, string(v))
// 	}

// 	return returnSlice
// }
