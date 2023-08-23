package base

import (
	"context"
	"fmt"
	"net/http"

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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationSettingsEmailResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type NotificationSettingsEmailResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Host          types.String `tfsdk:"host"`
	Port          types.Int64  `tfsdk:"port"`
	Protocol      types.String `tfsdk:"protocol"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	From          types.List   `tfsdk:"from"`
	ReplyTo       types.List   `tfsdk:"reply_to"`
}

type EmailSourceModel struct {
	Name         types.String `tfsdk:"name"`
	EmailAddress types.String `tfsdk:"email_address"`
}

var (
	emailSourceTFObjectTypes = map[string]attr.Type{
		"name":          types.StringType,
		"email_address": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &NotificationSettingsEmailResource{}
	_ resource.ResourceWithConfigure   = &NotificationSettingsEmailResource{}
	_ resource.ResourceWithImportState = &NotificationSettingsEmailResource{}
)

// New Object
func NewNotificationSettingsEmailResource() resource.Resource {
	return &NotificationSettingsEmailResource{}
}

// Metadata
func (r *NotificationSettingsEmailResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_settings_email"
}

// Schema
func (r *NotificationSettingsEmailResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5

	portDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the port used by the organization's SMTP server to send emails (default: `465`). Note that the protocol used depends upon the port specified. If you specify port `25`, `587`, or `2525`, SMTP with `STARTTLS` is used. Otherwise, `SMTPS` is used.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the email sender settings in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure email settings in."),
			),

			"host": schema.StringAttribute{
				Description: "A string that specifies the organization's SMTP server.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"port": schema.Int64Attribute{
				MarkdownDescription: portDescription.MarkdownDescription,
				Description:         portDescription.Description,
				Required:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(attrMinLength),
				},
			},

			"protocol": schema.StringAttribute{
				Description: "A string that specifies the current protocol in use.",
				Computed:    true,
			},

			"username": schema.StringAttribute{
				Description: "A string that specifies the organization's SMTP server's username.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"password": schema.StringAttribute{
				Description: "A string that specifies the organization's SMTP server's password.",
				Required:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
		},

		Blocks: map[string]schema.Block{
			"from": schema.ListNestedBlock{
				Description: "A required single block that specifies the email sender's \"from\" name and email address.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "A string that specifies the email sender's \"from\" name.",
							Optional:    true,
						},
						"email_address": schema.StringAttribute{
							Description: "A string that specifies the email sender's \"from\" email address.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(emailAddressMaxLength),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtMost(1),
				},
			},
			"reply_to": schema.ListNestedBlock{
				Description: "A single block that specifies the email sender's \"reply to\" name and email address.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "A string that specifies the email sender's \"reply to\" name.",
							Optional:    true,
						},
						"email_address": schema.StringAttribute{
							Description: "A string that specifies the email sender's \"reply to\" email address.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(emailAddressMaxLength),
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

func (r *NotificationSettingsEmailResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *NotificationSettingsEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state NotificationSettingsEmailResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
		},
		"UpdateEmailNotificationsSettings-Create",
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

func (r *NotificationSettingsEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NotificationSettingsEmailResourceModel

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
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.NotificationsSettingsSMTPApi.ReadEmailNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadEmailNotificationsSettings",
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

func (r *NotificationSettingsEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NotificationSettingsEmailResourceModel

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
	notificationSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
		},
		"UpdateEmailNotificationsSettings-Create",
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

func (r *NotificationSettingsEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NotificationSettingsEmailResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := r.client.NotificationsSettingsSMTPApi.DeleteEmailDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteEmailDeliverySettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationSettingsEmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (p *NotificationSettingsEmailResourceModel) expand(ctx context.Context) (*management.NotificationsSettingsEmailDeliverySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewNotificationsSettingsEmailDeliverySettings()

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

		from := management.NewNotificationsSettingsEmailDeliverySettingsFrom(plan[0].EmailAddress.ValueString())

		if !plan[0].Name.IsNull() && !plan[0].Name.IsUnknown() {
			from.SetName(plan[0].Name.ValueString())
		}

		data.SetFrom(*from)
	}

	if !p.ReplyTo.IsNull() && !p.ReplyTo.IsUnknown() {
		var plan []EmailSourceModel
		d := p.ReplyTo.ElementsAs(ctx, &plan, false)
		diags.Append(d...)

		replyTo := management.NewNotificationsSettingsEmailDeliverySettingsReplyTo()

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

func (p *NotificationSettingsEmailResourceModel) toState(apiObject *management.NotificationsSettingsEmailDeliverySettings) diag.Diagnostics {
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
	case *management.NotificationsSettingsEmailDeliverySettingsFrom:
		if t.GetAddress() == "" {
			return types.ListValueMust(tfObjType, []attr.Value{}), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsEmailDeliverySettingsReplyTo:
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
