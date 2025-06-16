// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationSettingsEmailResource serviceClientType

type requestsModel struct {
	Body           types.String `tfsdk:"body"`
	DeliveryMethod types.String `tfsdk:"delivery_method"`
	Headers        types.Map    `tfsdk:"headers"`
	Method         types.String `tfsdk:"method"`
	URL            types.String `tfsdk:"url"`
}

type notificationSettingsEmailResourceModelV1 struct {
	Id                 pingonetypes.ResourceIDValue `tfsdk:"id"`
	AuthToken          types.String                 `tfsdk:"auth_token"`
	CustomProviderName types.String                 `tfsdk:"custom_provider_name"`
	EnvironmentId      pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	From               types.Object                 `tfsdk:"from"`
	Host               types.String                 `tfsdk:"host"`
	Password           types.String                 `tfsdk:"password"`
	Protocol           types.String                 `tfsdk:"protocol"`
	Port               types.Int32                  `tfsdk:"port"`
	ProviderType       types.String                 `tfsdk:"provider_type"`
	ReplyTo            types.Object                 `tfsdk:"reply_to"`
	Requests           types.Set                    `tfsdk:"requests"`
	Username           types.String                 `tfsdk:"username"`
}

type emailSourceModelV1 struct {
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
	_ resource.Resource                   = &NotificationSettingsEmailResource{}
	_ resource.ResourceWithConfigure      = &NotificationSettingsEmailResource{}
	_ resource.ResourceWithImportState    = &NotificationSettingsEmailResource{}
	_ resource.ResourceWithValidateConfig = &NotificationSettingsEmailResource{}
)

// New Object
func NewNotificationSettingsEmailResource() resource.Resource {
	return &NotificationSettingsEmailResource{}
}

// Metadata
func (r *NotificationSettingsEmailResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_settings_email"
}

// ValidateConfig
func (p *NotificationSettingsEmailResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data notificationSettingsEmailResourceModelV1

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate either host or custom provider name is set
	if data.Host.IsNull() && data.Host.IsUnknown() && data.CustomProviderName.IsNull() && data.CustomProviderName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Host or Custom Provider Name Required",
			"Either `host` or `custom_provider_name` must be set to configure email settings.",
		)
	}

	if !data.Host.IsNull() && !data.Host.IsUnknown() && !data.CustomProviderName.IsNull() && !data.CustomProviderName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Host and Custom Provider Name Conflict",
			"`host` and `custom_provider_name` cannot be set at the same time. Please set only one of them.",
		)
	}

	// Custom provider validation
	if !data.CustomProviderName.IsNull() && !data.CustomProviderName.IsUnknown() {
		for _, request := range data.Requests.Elements() {
			reqObj, ok := request.(types.Object)
			if !ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("requests"),
					"Invalid Request Object",
					"The `requests` attribute must be a valid object. Please ensure the requests are defined correctly.",
				)
			}

			// Check if the method is set to POST and body is provided
			// If method is not POST, body should be empty
			methodAttr := reqObj.Attributes()["method"].(types.String).ValueString()
			bodyAttr := reqObj.Attributes()["body"].(types.String)

			if methodAttr != string(management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMREQUESTSMETHOD_POST) &&
				bodyAttr.ValueString() != "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("requests").AtName("body"),
					"Invalid Request Body",
					"The `body` in `requests` must be empty if the `method` is not set to `POST`. Please ensure the body is empty or the method is set to `POST`.",
				)
			}
		}
	}

	// Authentication validation
	if (data.AuthToken.IsNull() || data.AuthToken.IsUnknown()) &&
		(data.Username.IsNull() || data.Username.IsUnknown()) &&
		(data.Password.IsNull() || data.Password.IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_token"),
			"Authentication Conflict",
			"Either `auth_token` or `username` and `password` must be set.)",
		)
	}

	// SMTP provider validation
	if !data.Host.IsNull() && !data.Host.IsUnknown() {
		if !data.Protocol.IsNull() && !data.Protocol.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("protocol"),
				"Protocol Conflict",
				"`protocol` cannot be set when configuring a SMTP provider. Please remove the `protocol` attribute.",
			)
		}
	}
}

// Schema
func (r *NotificationSettingsEmailResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5

	portDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the port used by the organization's SMTP server to send emails (default: `465`). Note that the protocol used depends upon the port specified. If you specify port `25`, `587`, or `2525`, SMTP with `STARTTLS` is used. Otherwise, `SMTPS` is used.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the email sender settings in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure email settings in."),
			),

			"host": schema.StringAttribute{
				Description: "A string that specifies the organization's SMTP server.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(path.MatchRoot("auth_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("custom_provider_name")),
					stringvalidator.ConflictsWith(path.MatchRoot("requests")),
					stringvalidator.AlsoRequires(path.MatchRoot("from").AtName("email_address")),
				},
			},

			"port": schema.Int32Attribute{
				MarkdownDescription: portDescription.MarkdownDescription,
				Description:         portDescription.Description,
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.AtLeast(attrMinLength),
					int32validator.ConflictsWith(path.MatchRoot("custom_provider_name")),
					int32validator.ConflictsWith(path.MatchRoot("requests")),
				},
			},

			"protocol": schema.StringAttribute{
				Description: "A string that specifies the current protocol in use.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsEmailDeliverySettingsProtocolEnumValues)...),
				},
			},

			"username": schema.StringAttribute{
				Description: "A string that specifies the organization's server's username.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(path.MatchRoot("auth_token")),
				},
			},

			"password": schema.StringAttribute{
				Description: "A string that specifies the organization's server's password.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(path.MatchRoot("auth_token")),
				},
			},

			"from": schema.SingleNestedAttribute{
				Description: "A single block that specifies the email sender's \"from\" name and email address.",
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "A string that specifies the email sender's \"from\" name.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
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

			"reply_to": schema.SingleNestedAttribute{
				Description: "A single block that specifies the email sender's \"reply to\" name and email address.",
				Optional:    true,

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

			"auth_token": schema.StringAttribute{
				Description: "A string that specifies the authentication token when using a Custom Provider.",
				Sensitive:   true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
				},
			},

			"custom_provider_name": schema.StringAttribute{
				Description: "A string to use to identify the provider.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.AlsoRequires(path.MatchRoot("protocol")),
					stringvalidator.AlsoRequires(path.MatchRoot("requests")),
					stringvalidator.ConflictsWith(path.MatchRoot("host")),
					stringvalidator.ConflictsWith(path.MatchRoot("port")),
				},
			},

			"provider_type": schema.StringAttribute{
				Description: "A string that spefifies the provider type.",
				Computed:    true,
			},

			"requests": schema.SetNestedAttribute{
				Description: "A list of objects that is used to configure the API requests sent to the custom email provider.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"body": schema.StringAttribute{
							Description: "Required if method is set to `POST`. Use body to provide the content of the body for the request sent to the email provider.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},
						"delivery_method": schema.StringAttribute{
							Description: "A string that specifies the delivery method for the request.",
							Computed:    true,
						},
						"headers": schema.MapAttribute{
							Description: "A map of key-value pairs to specify the headers that your email provider's API expects.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"method": schema.StringAttribute{
							Description: "Use method to specify the type of API request the email provider requires. Valid values are `GET` and `POST`.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsEmailDeliverySettingsCustomRequestsMethodEnumValues)...),
							},
						},
						"url": schema.StringAttribute{
							Description: "A string that specifies the endpoint for your email provider.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.AlsoRequires(path.MatchRoot("custom_provider_name")),
					setvalidator.ConflictsWith(path.MatchRoot("port")),
					setvalidator.ConflictsWith(path.MatchRoot("host")),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *NotificationSettingsEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state notificationSettingsEmailResourceModelV1

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
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	var data *notificationSettingsEmailResourceModelV1

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
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsSMTPApi.ReadEmailNotificationsSettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
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
	var plan, state notificationSettingsEmailResourceModelV1

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
	var response *management.NotificationsSettingsEmailDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsSMTPApi.UpdateEmailNotificationsSettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsEmailDeliverySettings(*notificationSettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	var data *notificationSettingsEmailResourceModelV1

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.NotificationsSettingsSMTPApi.DeleteEmailDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
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

func (p *notificationSettingsEmailResourceModelV1) expand(ctx context.Context) (*management.NotificationsSettingsEmailDeliverySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	// SMTP settings
	if !p.Host.IsNull() && !p.Host.IsUnknown() {
		data := management.NewNotificationsSettingsEmailDeliverySettingsSMTP()

		if !p.Host.IsNull() && !p.Host.IsUnknown() {
			data.SetHost(p.Host.ValueString())
		}

		if !p.Port.IsNull() && !p.Port.IsUnknown() {
			data.SetPort(p.Port.ValueInt32())
		}

		if !p.Username.IsNull() && !p.Username.IsUnknown() {
			data.SetUsername(p.Username.ValueString())
		}

		if !p.Password.IsNull() && !p.Password.IsUnknown() {
			data.SetPassword(p.Password.ValueString())
		}

		if !p.From.IsNull() && !p.From.IsUnknown() {
			var plan emailSourceModelV1
			d := p.From.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)

			from := management.NewNotificationsSettingsEmailDeliverySettingsSMTPAllOfFrom(plan.EmailAddress.ValueString())

			if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
				from.SetName(plan.Name.ValueString())
			}

			data.SetFrom(*from)
		}

		if !p.ReplyTo.IsNull() && !p.ReplyTo.IsUnknown() {
			var plan emailSourceModelV1
			d := p.ReplyTo.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)

			replyTo := management.NewNotificationsSettingsEmailDeliverySettingsSMTPAllOfReplyTo()

			if !plan.EmailAddress.IsNull() && !plan.EmailAddress.IsUnknown() {
				replyTo.SetAddress(plan.EmailAddress.ValueString())
			}

			if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
				replyTo.SetName(plan.Name.ValueString())
			}

			data.SetReplyTo(*replyTo)
		}

		return &management.NotificationsSettingsEmailDeliverySettings{
			NotificationsSettingsEmailDeliverySettingsSMTP: data,
		}, diags
	}

	// Custom provider settings
	if !p.CustomProviderName.IsNull() && !p.CustomProviderName.IsUnknown() {
		requests := []management.NotificationsSettingsEmailDeliverySettingsCustomAllOfRequests{}
		if !p.Requests.IsNull() && !p.Requests.IsUnknown() {
			var requestModels []requestsModel
			d := p.Requests.ElementsAs(ctx, &requestModels, true)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			deliveryMethod, err := management.NewEnumNotificationsSettingsEmailDeliverySettingsCustomRequestsDeliveryMethodFromValue(string(management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMREQUESTSDELIVERYMETHOD_EMAIL))
			if err != nil {
				diags.AddError(
					"Invalid Delivery Method",
					fmt.Sprintf("Error: '%v'\n Verify the value used is use one of the allowed values: %v", err, management.AllowedEnumNotificationsSettingsEmailDeliverySettingsCustomRequestsDeliveryMethodEnumValues),
				)
				return nil, diags
			}

			for _, request := range requestModels {
				req := management.NewNotificationsSettingsEmailDeliverySettingsCustomAllOfRequests(
					*deliveryMethod,
					management.EnumNotificationsSettingsEmailDeliverySettingsCustomRequestsMethod(request.Method.ValueString()),
					request.URL.ValueString(),
				)

				if !request.Headers.IsNull() && !request.Headers.IsUnknown() {
					headers := make(map[string]string, len(request.Headers.Elements()))
					diags.Append(request.Headers.ElementsAs(ctx, &headers, true)...)
					if diags.HasError() {
						return nil, diags
					}
					req.SetHeaders(headers)
				}

				method, err := management.NewEnumNotificationsSettingsEmailDeliverySettingsCustomRequestsMethodFromValue(request.Method.ValueString())
				if err != nil {
					diags.AddError(
						"Invalid Method",
						fmt.Sprintf("The method '%s' is not valid. Please use one of the allowed values: %v", request.Method.ValueString(), management.AllowedEnumNotificationsSettingsEmailDeliverySettingsCustomRequestsMethodEnumValues),
					)
					return nil, diags
				}

				req.SetMethod(*method)
				req.SetUrl(request.URL.ValueString())

				if !request.Body.IsNull() && !request.Body.IsUnknown() {
					// Check if the body contains the required variables - ${to}, ${message}
					if !strings.Contains(request.Body.ValueString(), "${to}") || !strings.Contains(request.Body.ValueString(), "${message}") {
						diags.AddError(
							"Invalid Body Content",
							"The body must contain the variables `${to}` and `${message}`. Please ensure these variables are included in the body content.",
						)
						return nil, diags
					}

					req.SetBody(request.Body.ValueString())
				}

				requests = append(requests, *req)
			}
		}

		protocolEnum, err := management.NewEnumNotificationsSettingsEmailDeliverySettingsProtocolFromValue(p.Protocol.ValueString())
		if err != nil {
			diags.AddError(
				"Invalid Protocol",
				fmt.Sprintf("The protocol '%s' is not valid. Please use one of the allowed values: %v", p.Protocol.ValueString(), management.AllowedEnumNotificationsSettingsEmailDeliverySettingsProtocolEnumValues),
			)
			return nil, diags
		}
		var authMethod, authToken, username, password string
		if !p.Username.IsNull() && !p.Username.IsUnknown() && !p.Password.IsNull() && !p.Password.IsUnknown() {
			username = p.Username.ValueString()
			password = p.Password.ValueString()
			authMethod = string(management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMAUTHENTICATIONMETHOD_BASIC)
		} else if !p.AuthToken.IsNull() && !p.AuthToken.IsUnknown() {
			authMethod = string(management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMAUTHENTICATIONMETHOD_BEARER)
		}

		authenticationMethod, err := management.NewEnumNotificationsSettingsEmailDeliverySettingsCustomAuthenticationMethodFromValue(authMethod)
		if err != nil {
			diags.AddError(
				"Invalid Authentication Method",
				fmt.Sprintf("The authentication method '%s' is not valid. Please use one of the allowed values: %v", authMethod, management.AllowedEnumNotificationsSettingsEmailDeliverySettingsCustomAuthenticationMethodEnumValues),
			)
			return nil, diags
		}

		authentication := management.NewNotificationsSettingsEmailDeliverySettingsCustomAllOfAuthentication(*authenticationMethod)

		data := management.NewNotificationsSettingsEmailDeliverySettingsCustom(
			*protocolEnum,
			*authentication,
			p.CustomProviderName.ValueString(),
			requests,
		)

		if !p.From.IsNull() && !p.From.IsUnknown() {
			var plan emailSourceModelV1
			d := p.From.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)

			from := management.NewNotificationsSettingsEmailDeliverySettingsCustomAllOfFrom()

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
			d := p.ReplyTo.As(ctx, &plan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			replyTo := management.NewNotificationsSettingsEmailDeliverySettingsCustomAllOfReplyTo()

			if !plan.EmailAddress.IsNull() && !plan.EmailAddress.IsUnknown() {
				replyTo.SetAddress(plan.EmailAddress.ValueString())
			}

			if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
				replyTo.SetName(plan.Name.ValueString())
			}

			data.SetReplyTo(*replyTo)
		}

		switch *authenticationMethod {
		case management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMAUTHENTICATIONMETHOD_BEARER:
			data.Authentication.SetAuthToken(authToken)
		case management.ENUMNOTIFICATIONSSETTINGSEMAILDELIVERYSETTINGSCUSTOMAUTHENTICATIONMETHOD_BASIC:
			data.Authentication.SetUsername(username)
			data.Authentication.SetPassword(password)
		default:
			diags.AddError(
				"Unsupported Authentication Method",
				fmt.Sprintf("The authentication method '%s' is not supported. Please use one of the allowed values: %v", utils.EnumToString(*authenticationMethod), management.AllowedEnumNotificationsSettingsEmailDeliverySettingsCustomAuthenticationMethodEnumValues),
			)
			// Return early if the authentication method is not supported
			return nil, diags
		}

		p.Host = types.StringNull()
		p.Port = types.Int32Null()

		return &management.NotificationsSettingsEmailDeliverySettings{
			NotificationsSettingsEmailDeliverySettingsCustom: data,
		}, diags
	}

	return nil, diags
}

func (p *notificationSettingsEmailResourceModelV1) toState(apiObject *management.NotificationsSettingsEmailDeliverySettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	reqAttrTypes := map[string]attr.Type{
		"body":            types.StringType,
		"delivery_method": types.StringType,
		"headers":         types.MapType{ElemType: types.StringType},
		"method":          types.StringType,
		"url":             types.StringType,
	}

	p.Id = p.EnvironmentId

	switch t := apiObject.GetActualInstance().(type) {
	case *management.NotificationsSettingsEmailDeliverySettingsSMTP:
		p.Host = framework.StringOkToTF(t.GetHostOk())
		p.Port = framework.Int32OkToTF(t.GetPortOk())
		p.Protocol = framework.EnumOkToTF(t.GetProtocolOk())
		p.Username = framework.StringOkToTF(t.GetUsernameOk())

		from, d := toStateEmailSource(t.GetFromOk())
		diags.Append(d...)
		p.From = from

		replyTo, d := toStateEmailSource(t.GetReplyToOk())
		diags.Append(d...)
		p.ReplyTo = replyTo

		p.AuthToken = types.StringNull()
		p.CustomProviderName = types.StringNull()
		p.ProviderType = types.StringNull()
		p.Requests = types.SetNull(types.ObjectType{AttrTypes: reqAttrTypes})

	case *management.NotificationsSettingsEmailDeliverySettingsCustom:

		if t.Authentication.AuthToken != nil {
			p.AuthToken = framework.StringOkToTF(t.Authentication.GetAuthTokenOk())
		}

		if t.Authentication.Username != nil {
			p.Username = framework.StringOkToTF(t.Authentication.GetUsernameOk())
		}

		p.CustomProviderName = framework.StringOkToTF(t.GetNameOk())
		p.ProviderType = framework.EnumOkToTF(t.GetProviderOk())
		p.Protocol = framework.EnumOkToTF(t.GetProtocolOk())

		from, d := toStateEmailSource(t.GetFromOk())
		diags.Append(d...)
		p.From = from

		replyTo, d := toStateEmailSource(t.GetReplyToOk())
		diags.Append(d...)
		p.ReplyTo = replyTo

		if t.Requests != nil {
			requests := make([]attr.Value, 0, len(t.Requests))
			for _, request := range t.Requests {
				req := map[string]attr.Value{
					"delivery_method": framework.EnumOkToTF(request.GetDeliveryMethodOk()),
					"headers":         framework.StringMapOkToTF(request.GetHeadersOk()),
					"method":          framework.EnumOkToTF(request.GetMethodOk()),
					"url":             framework.StringOkToTF(request.GetUrlOk()),
				}

				if request.GetBody() != "" {
					req["body"] = framework.StringOkToTF(request.GetBodyOk())
				} else {
					req["body"] = types.StringNull()
				}

				reqValue, d := types.ObjectValue(reqAttrTypes, req)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				requests = append(requests, reqValue)
			}
			p.Requests, diags = types.SetValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"body":            types.StringType,
						"delivery_method": types.StringType,
						"headers":         types.MapType{ElemType: types.StringType},
						"method":          types.StringType,
						"url":             types.StringType,
					},
				},
				requests,
			)
			if diags.HasError() {
				return diags
			}
		} else {
			p.Requests = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"body":            types.StringType,
					"delivery_method": types.StringType,
					"headers":         types.MapType{ElemType: types.StringType},
					"method":          types.StringType,
					"url":             types.StringType,
				},
			})
		}

	}

	return diags
}

func toStateEmailSource(emailSource interface{}, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || emailSource == nil {
		return types.ObjectNull(emailSourceTFObjectTypes), diags
	}

	var emailSourceMap map[string]attr.Value

	switch t := emailSource.(type) {
	case *management.NotificationsSettingsEmailDeliverySettingsSMTPAllOfFrom:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsEmailDeliverySettingsSMTPAllOfReplyTo:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}

		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}

		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsEmailDeliverySettingsCustomAllOfFrom:
		if t.GetAddress() == "" {
			return types.ObjectNull(emailSourceTFObjectTypes), diags
		}
		emailSourceMap = map[string]attr.Value{
			"email_address": framework.StringOkToTF(t.GetAddressOk()),
		}
		emailSourceMap["name"] = framework.StringOkToTF(t.GetNameOk())

	case *management.NotificationsSettingsEmailDeliverySettingsCustomAllOfReplyTo:
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
