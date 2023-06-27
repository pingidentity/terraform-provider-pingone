package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationSettingsResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type NotificationSettingsResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	EnvironmentId         types.String `tfsdk:"environment_id"`
	DeliveryMode          types.String `tfsdk:"delivery_mode"`
	ProviderFallbackChain types.List   `tfsdk:"provider_fallback_chain"`
	From                  types.Object `tfsdk:"from"`
	ReplyTo               types.Object `tfsdk:"reply_to"`
	AllowedList           types.Set    `tfsdk:"allowed_list"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

type NotificationSettingsAllowedListResourceModel struct {
	UserID types.String `tfsdk:"user_id"`
}

var (
	allowedListTFObjectTypes = map[string]attr.Type{
		"user_id": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &NotificationSettingsResource{}
	_ resource.ResourceWithConfigure   = &NotificationSettingsResource{}
	_ resource.ResourceWithImportState = &NotificationSettingsResource{}
)

// New Object
func NewNotificationSettingsResource() resource.Resource {
	return &NotificationSettingsResource{}
}

// Metadata
func (r *NotificationSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_settings"
}

// Schema
func (r *NotificationSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5

	deliveryModeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the delivery mode that the settings apply for.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsDeliveryModeEnumValues).DefaultValue(string(management.ENUMNOTIFICATIONSSETTINGSDELIVERYMODE_ALL))

	providerFallbackChainDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of strings that which represents the execution order of different SMS/Voice providers configured for the environment. The providers and their accounts’ configurations are represented in the list by the ID of the corresponding `pingone_phone_delivery_settings` resource. The only provider which is not represented by the `pingone_phone_delivery_settings.id` value is the PingOne Twilio provider. The PingOne Twilio provider is represented by the `PINGONE_TWILIO` string. If this parameter's list is empty, an SMS or voice message will be sent using the default Ping Twilio account. Otherwise, an SMS or voice message will be sent using the first provider in the list. If the server fails to queue the message using that provider, it will use the next provider in the list to try to send the message. This process will go on until there are no more providers in the list. If the server failed to send the message using all providers, the notification status is set to `FAILED`.",
	)

	allowedListDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that represent actors that are exempt from any delivery restrictions.",
	)

	fromNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email sender's \"from\" name.",
	).DefaultValue("PingOne")

	fromEmailAddressDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email sender's \"from\" email address.",
	).DefaultValue("noreply@pingidentity.com")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the notifications settings in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure notifications settings in."),
			),

			"delivery_mode": schema.StringAttribute{
				Description:         deliveryModeDescription.Description,
				MarkdownDescription: deliveryModeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(management.ENUMNOTIFICATIONSSETTINGSDELIVERYMODE_ALL)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsDeliveryModeEnumValues)...),
				},
			},

			"provider_fallback_chain": schema.ListAttribute{
				Description:         providerFallbackChainDescription.Description,
				MarkdownDescription: providerFallbackChainDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),

				ElementType: types.StringType,

				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.Any(
							verify.P1ResourceIDValidator(),
							stringvalidator.OneOf([]string{"PINGONE_TWILIO"}...),
						),
					),
				},
			},

			"allowed_list": schema.SetNestedAttribute{
				Description:         allowedListDescription.Description,
				MarkdownDescription: allowedListDescription.MarkdownDescription,
				Optional:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the user ID to add to the allowed list.  Must be a valid PingOne resource ID.").Description,
							Required:    true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
						},
					},
				},
			},

			"from": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A required single block that specifies the email sender's \"from\" name and email address.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"name":          types.StringValue("PingOne"),
						"email_address": types.StringValue("noreply@pingidentity.com"),
					}

					objValue, d := types.ObjectValue(emailSourceTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         fromNameDescription.Description,
						MarkdownDescription: fromNameDescription.MarkdownDescription,
						Optional:            true,
					},

					"email_address": schema.StringAttribute{
						Description:         fromEmailAddressDescription.Description,
						MarkdownDescription: fromEmailAddressDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(emailAddressMaxLength),
						},
					},
				},
			},

			"reply_to": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A required single block that specifies the email sender's \"reply to\" name and email address.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the email sender's \"reply to\" name.").Description,
						Optional:    true,
					},

					"email_address": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the email sender's \"reply to\" email address.").Description,
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(emailAddressMaxLength),
						},
					},
				},
			},

			"updated_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the time the resource was last updated in RFC3339 format.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *NotificationSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationSettingsResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsSettingsApi.UpdateNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettings(*notificationSettings).Execute()
		},
		"UpdateNotificationsSettings-Create",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsSettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NotificationSettingsResourceModel

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
			return r.client.NotificationsSettingsApi.ReadNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadNotificationsSettings",
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
	resp.Diagnostics.Append(data.toState(response.(*management.NotificationsSettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationSettingsResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.NotificationsSettingsApi.UpdateNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettings(*notificationSettings).Execute()
		},
		"UpdateNotificationsSettings-Update",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsSettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NotificationSettingsResourceModel

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
			return r.client.NotificationsSettingsApi.DeleteNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"DeleteNotificationsSettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (p *NotificationSettingsResourceModel) expand(ctx context.Context) (*management.NotificationsSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewNotificationsSettings()

	if !p.DeliveryMode.IsNull() && !p.DeliveryMode.IsUnknown() {
		data.SetDeliveryMode(management.EnumNotificationsSettingsDeliveryMode(p.DeliveryMode.ValueString()))
	}

	if !p.ProviderFallbackChain.IsNull() && !p.ProviderFallbackChain.IsUnknown() {

		var plan []string
		d := p.ProviderFallbackChain.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetSmsProvidersFallbackChain(plan)
	}

	if !p.AllowedList.IsNull() && !p.AllowedList.IsUnknown() {

		var plan []NotificationSettingsAllowedListResourceModel
		d := p.AllowedList.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		allowedList := make([]management.NotificationsSettingsWhitelistInner, 0)
		for _, v := range plan {
			allowedItem := *management.NewNotificationsSettingsWhitelistInner()

			if !v.UserID.IsNull() && !v.UserID.IsUnknown() {
				allowedItem.SetUser(
					*management.NewApplicationAccessControlGroupGroupsInner(v.UserID.ValueString()),
				)
			}

			allowedList = append(allowedList, allowedItem)
		}

		data.SetWhitelist(allowedList)
	}

	if !p.From.IsNull() && !p.From.IsUnknown() {
		var plan EmailSourceModel
		diags.Append(p.From.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)

		from := management.NewNotificationsSettingsFrom()

		if !plan.EmailAddress.IsNull() && !plan.EmailAddress.IsUnknown() {
			from.SetAddress(plan.EmailAddress.ValueString())
		}

		if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
			from.SetName(plan.Name.ValueString())
		}

		data.SetFrom(*from)
	}

	if !p.ReplyTo.IsNull() && !p.ReplyTo.IsUnknown() {
		var plan EmailSourceModel
		diags.Append(p.ReplyTo.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)

		replyTo := management.NewNotificationsSettingsReplyTo()

		if !plan.EmailAddress.IsNull() && !plan.EmailAddress.IsUnknown() {
			replyTo.SetAddress(plan.EmailAddress.ValueString())
		}

		if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
			replyTo.SetName(plan.Name.ValueString())
		}

		data.SetReplyTo(*replyTo)
	}

	return data, diags
}

func (p *NotificationSettingsResourceModel) toState(apiObject *management.NotificationsSettings) diag.Diagnostics {
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

	p.DeliveryMode = framework.EnumOkToTF(apiObject.GetDeliveryModeOk())
	p.ProviderFallbackChain = framework.StringListOkToTF(apiObject.GetSmsProvidersFallbackChainOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	var d diag.Diagnostics
	p.AllowedList, d = notificationsSettingsAllowedListOkToTF(apiObject.GetWhitelistOk())
	diags.Append(d...)

	from, d := toStateEmailSourceObject(apiObject.GetFromOk())
	diags.Append(d...)
	p.From = from

	replyTo, d := toStateEmailSourceObject(apiObject.GetReplyToOk())
	diags.Append(d...)
	p.ReplyTo = replyTo

	return diags
}

func notificationsSettingsAllowedListOkToTF(apiObject []management.NotificationsSettingsWhitelistInner, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: allowedListTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {
		if user, ok := v.GetUserOk(); ok {

			objMap := map[string]attr.Value{
				"user_id": framework.StringOkToTF(user.GetIdOk()),
			}

			flattenedObj, d := types.ObjectValue(allowedListTFObjectTypes, objMap)
			diags.Append(d...)

			flattenedList = append(flattenedList, flattenedObj)
		}
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateEmailSourceObject(emailSource interface{}, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || emailSource == nil {
		return types.ObjectNull(emailSourceTFObjectTypes), diags
	}

	var emailSourceMap map[string]attr.Value

	switch t := emailSource.(type) {
	case *management.NotificationsSettingsEmailDeliverySettingsFrom:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsEmailDeliverySettingsReplyTo:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsFrom:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsReplyTo:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
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

		return types.ObjectNull(emailSourceTFObjectTypes), diags
	}

	returnVar, d := types.ObjectValue(emailSourceTFObjectTypes, emailSourceMap)
	diags.Append(d...)

	return returnVar, diags

}
