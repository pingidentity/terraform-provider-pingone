// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationSettingsResource serviceClientType

type NotificationSettingsResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	DeliveryMode          types.String                 `tfsdk:"delivery_mode"`
	ProviderFallbackChain types.List                   `tfsdk:"provider_fallback_chain"`
	From                  types.Object                 `tfsdk:"from"`
	ReplyTo               types.Object                 `tfsdk:"reply_to"`
	AllowedList           types.Set                    `tfsdk:"allowed_list"`
	UpdatedAt             timetypes.RFC3339            `tfsdk:"updated_at"`
}

type NotificationSettingsAllowedListResourceModel struct {
	UserID pingonetypes.ResourceIDValue `tfsdk:"user_id"`
}

var (
	allowedListTFObjectTypes = map[string]attr.Type{
		"user_id": pingonetypes.ResourceIDType{},
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

							CustomType: pingonetypes.ResourceIDType{},
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

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *NotificationSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationSettingsResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsSettings
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsApi.UpdateNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettings(*notificationSettings).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateNotificationsSettings-Create",
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

func (r *NotificationSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NotificationSettingsResourceModel

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
	var response *management.NotificationsSettings
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsApi.ReadNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadNotificationsSettings",
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

func (r *NotificationSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationSettingsResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsSettings
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsApi.UpdateNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettings(*notificationSettings).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateNotificationsSettings-Update",
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

func (r *NotificationSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NotificationSettingsResourceModel

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
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsApi.DeleteNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"DeleteNotificationsSettings",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:     "environment_id",
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes["environment_id"])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes["environment_id"])...)
}

func (p *NotificationSettingsResourceModel) expand(ctx context.Context) (*management.NotificationsSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewNotificationsSettings()

	if !p.DeliveryMode.IsNull() && !p.DeliveryMode.IsUnknown() {
		data.SetDeliveryMode(management.EnumNotificationsSettingsDeliveryMode(p.DeliveryMode.ValueString()))
	}

	if !p.ProviderFallbackChain.IsNull() && !p.ProviderFallbackChain.IsUnknown() {

		var plan []types.String
		d := p.ProviderFallbackChain.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		smsProvidersFallbackChain, d := framework.TFTypeStringSliceToStringSlice(plan, path.Root("provider_fallback_chain"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetSmsProvidersFallbackChain(smsProvidersFallbackChain)
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
		var plan emailSourceModelV1
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
		var plan emailSourceModelV1
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

	p.Id = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)

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
				"user_id": framework.PingOneResourceIDOkToTF(user.GetIdOk()),
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
