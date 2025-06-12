// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type NotificationTemplateContentResource serviceClientType

type notificationTemplateContentResourceModelV1 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	TemplateName  types.String                 `tfsdk:"template_name"`
	Locale        types.String                 `tfsdk:"locale"`
	Default       types.Bool                   `tfsdk:"default"`
	Variant       types.String                 `tfsdk:"variant"`
	Email         types.Object                 `tfsdk:"email"`
	Push          types.Object                 `tfsdk:"push"`
	Sms           types.Object                 `tfsdk:"sms"`
	Voice         types.Object                 `tfsdk:"voice"`
}

type notificationTemplateContentEmailResourceModelV1 struct {
	Body         types.String `tfsdk:"body"`
	From         types.Object `tfsdk:"from"`
	Subject      types.String `tfsdk:"subject"`
	ReplyTo      types.Object `tfsdk:"reply_to"`
	CharacterSet types.String `tfsdk:"character_set"`
	ContentType  types.String `tfsdk:"content_type"`
}

type notificationTemplateContentEmailAddressResourceModelV1 struct {
	Name    types.String `tfsdk:"name"`
	Address types.String `tfsdk:"address"`
}

type notificationTemplateContentPushResourceModelV1 struct {
	Category types.String `tfsdk:"category"`
	Body     types.String `tfsdk:"body"`
	Title    types.String `tfsdk:"title"`
}

type notificationTemplateContentSmsResourceModelV1 struct {
	Content types.String `tfsdk:"content"`
	Sender  types.String `tfsdk:"sender"`
}

type notificationTemplateContentVoiceResourceModelV1 struct {
	Content types.String `tfsdk:"content"`
	Type    types.String `tfsdk:"type"`
}

var (
	notificationTemplateContentEmailTFObjectTypes = map[string]attr.Type{
		"body": types.StringType,
		"from": types.ObjectType{
			AttrTypes: notificationTemplateContentEmailAddressTFObjectTypes,
		},
		"subject": types.StringType,
		"reply_to": types.ObjectType{
			AttrTypes: notificationTemplateContentEmailAddressTFObjectTypes,
		},
		"character_set": types.StringType,
		"content_type":  types.StringType,
	}

	notificationTemplateContentEmailAddressTFObjectTypes = map[string]attr.Type{
		"name":    types.StringType,
		"address": types.StringType,
	}

	notificationTemplateContentPushTFObjectTypes = map[string]attr.Type{
		"category": types.StringType,
		"body":     types.StringType,
		"title":    types.StringType,
	}

	notificationTemplateContentSmsTFObjectTypes = map[string]attr.Type{
		"content": types.StringType,
		"sender":  types.StringType,
	}

	notificationTemplateContentVoiceTFObjectTypes = map[string]attr.Type{
		"content": types.StringType,
		"type":    types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &NotificationTemplateContentResource{}
	_ resource.ResourceWithConfigure   = &NotificationTemplateContentResource{}
	_ resource.ResourceWithImportState = &NotificationTemplateContentResource{}
)

// New Object
func NewNotificationTemplateContentResource() resource.Resource {
	return &NotificationTemplateContentResource{}
}

// Metadata
func (r *NotificationTemplateContentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_template_content"
}

// Schema.
func (r *NotificationTemplateContentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	templateNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the template to manage localised contents for.",
	).AllowedValuesEnum(management.AllowedEnumTemplateNameEnumValues).RequiresReplace()

	localeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies an ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).",
	).AllowedValuesEnum(verify.FullIsoList()).RequiresReplace()

	const variantMinLength = 1
	const variantMaxLength = 100
	variantDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the unique user-defined name for each content variant that uses the same template + `deliveryMethod` + `locale` combination.  This property is case insensitive and has a limit of %d characters.", variantMaxLength),
	)

	emailDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the `email` delivery method.  Exactly one of `email`, `push`, `sms` or `voice` must be specified.",
	)

	const emailBodyMinLength = 1
	const emailBodyMaxLength = 100000

	emailFromDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the email sender.",
	)

	emailFromNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email's sender name.  If the environment uses the Ping Identity email sender, the name `PingOne` is used. You can configure other email sender names per environment.",
	)

	emailFromAddressDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the sender email address. If the environment uses the Ping Identity email sender, or if the address field is empty, the address `noreply@pingidentity.com` is used.  You can configure other email sender addresses per environment.",
	)

	const emailSubjectMinLength = 1
	const emailSubjectMaxLength = 256

	emailReplyToDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the email \"reply to\" address.",
	)

	emailReplyToNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email's \"reply to\" name.  If the environment uses the Ping Identity email sender, the name `PingOne` is used.  You can configure other email \"reply to\" names per environment.",
	)

	emailReplyToAddressDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the \"reply to\" email address.  If the environment uses the Ping Identity email sender, or if the address field is empty, the address `noreply@pingidentity.com` is used.  You can configure other email \"reply to\" addresses per environment.",
	)

	emailCharacterSetDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email's character set.",
	).DefaultValue("UTF-8")

	emailContentTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the email's content-type.",
	).DefaultValue("text/html")

	pushDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the `push` delivery method.  Exactly one of `email`, `push`, `sms` or `voice` must be specified.",
	)

	pushCategoryDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies what type of banner should be displayed to the user.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS):         "the banner contains both Approve and Deny buttons",
		string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_WITHOUT_BANNER_BUTTONS): "when the user clicks the banner, they are taken to an application that contains the necessary approval controls",
		string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_APPROVE_AND_OPEN_APP):   "when the Approve button is clicked, authentication is completed and the user is taken to the relevant application",
	}).DefaultValue(string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS)).AppendMarkdownString("Note that to use the non-default push banners, you must implement them in your application code, using the PingOne SDK. For details, see the [README for iOS](https://github.com/pingidentity/pingone-mobile-sdk-ios/#171-push-notifications-categories) and the [README for Android](https://github.com/pingidentity/pingone-mobile-sdk-android).")

	const pushBodyMinLength = 1
	const pushBodyMaxLength = 400

	const pushTitleMinLength = 1
	const pushTitleMaxLength = 200

	smsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the `sms` delivery method.  Exactly one of `email`, `push`, `sms` or `voice` must be specified.",
	)

	smsSenderDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the SMS sender ID. This property can contain only alphanumeric characters and spaces, and its length cannot exceed 11 characters. In some countries, it is impossible to send an SMS with an alphanumeric sender ID. For those countries, the sender ID must be empty. For SMS recipients in specific countries, refer to Twilio's documentation on [International support for Alphanumeric Sender ID](https://support.twilio.com/hc/en-us/articles/223133767-International-support-for-Alphanumeric-Sender-ID).",
	)

	voiceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies properties for the `voice` delivery method.  Exactly one of `email`, `push`, `sms` or `voice` must be specified.",
	)

	const voiceContentMinLength = 1
	const voiceContentMaxLength = 1024

	voiceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the voice type desired for the message. Out of the box options include `Man`, `Woman`, `Alice` (Twilio only), `Amazon Polly`, or your own user-defined custom string. In the case that the selected voice type is not supported by the provider in the desired locale, another voice type will be automatically selected. Additional charges may be incurred for these selections, as determined by the sender.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne notification template contents for push, SMS, email and voice notifications in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage notification template contents in."),
			),

			"template_name": schema.StringAttribute{
				Description:         templateNameDescription.Description,
				MarkdownDescription: templateNameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumTemplateNameEnumValues)...),
				},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"locale": schema.StringAttribute{
				Description:         localeDescription.Description,
				MarkdownDescription: localeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(verify.FullIsoList()...),
				},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"default": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the template is a predefined default template.").Description,
				Computed:    true,
			},

			"variant": schema.StringAttribute{
				Description:         variantDescription.Description,
				MarkdownDescription: variantDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthBetween(variantMinLength, variantMaxLength),
				},
			},

			"email": schema.SingleNestedAttribute{
				Description:         emailDescription.Description,
				MarkdownDescription: emailDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string representing the email body. Email text can contain HTML but cannot be larger than 100 kB.  Use of variables is supported.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(emailBodyMinLength, emailBodyMaxLength),
						},
					},

					"from": schema.SingleNestedAttribute{
						Description:         emailFromDescription.Description,
						MarkdownDescription: emailFromDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description:         emailFromNameDescription.Description,
								MarkdownDescription: emailFromNameDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Default: stringdefault.StaticString("PingOne"),

								Validators: []validator.String{
									stringvalidator.AtLeastOneOf(
										path.MatchRelative().AtParent().AtName("name"),
										path.MatchRelative().AtParent().AtName("address"),
									),
								},
							},

							"address": schema.StringAttribute{
								Description:         emailFromAddressDescription.Description,
								MarkdownDescription: emailFromAddressDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Default: stringdefault.StaticString("noreply@pingidentity.com"),

								Validators: []validator.String{
									stringvalidator.AtLeastOneOf(
										path.MatchRelative().AtParent().AtName("name"),
										path.MatchRelative().AtParent().AtName("address"),
									),
								},
							},
						},
					},

					"subject": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string representing the email's subject line. Cannot exceed 256 characters. Can include variables.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(emailSubjectMinLength, emailSubjectMaxLength),
						},
					},

					"reply_to": schema.SingleNestedAttribute{
						Description:         emailReplyToDescription.Description,
						MarkdownDescription: emailReplyToDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description:         emailReplyToNameDescription.Description,
								MarkdownDescription: emailReplyToNameDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Default: stringdefault.StaticString("PingOne"),

								Validators: []validator.String{
									stringvalidator.AtLeastOneOf(
										path.MatchRelative().AtParent().AtName("name"),
										path.MatchRelative().AtParent().AtName("address"),
									),
								},
							},

							"address": schema.StringAttribute{
								Description:         emailReplyToAddressDescription.Description,
								MarkdownDescription: emailReplyToAddressDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,

								Default: stringdefault.StaticString("noreply@pingidentity.com"),

								Validators: []validator.String{
									stringvalidator.AtLeastOneOf(
										path.MatchRelative().AtParent().AtName("name"),
										path.MatchRelative().AtParent().AtName("address"),
									),
								},
							},
						},
					},

					"character_set": schema.StringAttribute{
						Description:         emailCharacterSetDescription.Description,
						MarkdownDescription: emailCharacterSetDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString("UTF-8"),
					},

					"content_type": schema.StringAttribute{
						Description:         emailContentTypeDescription.Description,
						MarkdownDescription: emailContentTypeDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString("text/html"),
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("email"),
						path.MatchRelative().AtParent().AtName("push"),
						path.MatchRelative().AtParent().AtName("sms"),
						path.MatchRelative().AtParent().AtName("voice"),
					),
				},
			},

			"push": schema.SingleNestedAttribute{
				Description:         pushDescription.Description,
				MarkdownDescription: pushDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"category": schema.StringAttribute{
						Description:         pushCategoryDescription.Description,
						MarkdownDescription: pushCategoryDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumTemplateContentPushCategoryEnumValues)...),
						},
					},

					"body": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the push notification text. This can include variables.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(pushBodyMinLength, pushBodyMaxLength),
						},
					},

					"title": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the push notification title. This can include variables.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(pushTitleMinLength, pushTitleMaxLength),
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("email"),
						path.MatchRelative().AtParent().AtName("push"),
						path.MatchRelative().AtParent().AtName("sms"),
						path.MatchRelative().AtParent().AtName("voice"),
					),
				},
			},

			"sms": schema.SingleNestedAttribute{
				Description:         smsDescription.Description,
				MarkdownDescription: smsDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"content": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the SMS text. UC-2 encoding is used for text that contains non GSM-7 characters. UC-2 encoded text cannot exceed 67 characters. GSM-7 encoded text cannot exceed 153 characters. This can include variables.").Description,
						Required:    true,
					},

					"sender": schema.StringAttribute{
						Description:         smsSenderDescription.Description,
						MarkdownDescription: smsSenderDescription.MarkdownDescription,
						Optional:            true,
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("email"),
						path.MatchRelative().AtParent().AtName("push"),
						path.MatchRelative().AtParent().AtName("sms"),
						path.MatchRelative().AtParent().AtName("voice"),
					),
				},
			},

			"voice": schema.SingleNestedAttribute{
				Description:         voiceDescription.Description,
				MarkdownDescription: voiceDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"content": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the voice text to read. This can include variables.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(voiceContentMinLength, voiceContentMaxLength),
						},
					},

					"type": schema.StringAttribute{
						Description:         voiceTypeDescription.Description,
						MarkdownDescription: voiceTypeDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString("Alice"),
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("email"),
						path.MatchRelative().AtParent().AtName("push"),
						path.MatchRelative().AtParent().AtName("sms"),
						path.MatchRelative().AtParent().AtName("voice"),
					),
				},
			},
		},
	}
}

func (r *NotificationTemplateContentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationTemplateContentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state notificationTemplateContentResourceModelV1

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
	notificationTemplateContent, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API calls
	var response *management.TemplateContent
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsTemplatesApi.CreateContent(ctx, plan.EnvironmentId.ValueString(), management.EnumTemplateName(plan.TemplateName.ValueString())).TemplateContent(*notificationTemplateContent).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateContent",
		notificationTemplateCustomWriteError,
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

func (r *NotificationTemplateContentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *notificationTemplateContentResourceModelV1

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
	var response *management.TemplateContent
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsTemplatesApi.ReadOneContent(ctx, data.EnvironmentId.ValueString(), management.EnumTemplateName(data.TemplateName.ValueString()), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneContent",
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

func (r *NotificationTemplateContentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state notificationTemplateContentResourceModelV1

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
	notificationTemplateContent, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.TemplateContent
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.NotificationsTemplatesApi.UpdateContent(ctx, plan.EnvironmentId.ValueString(), management.EnumTemplateName(plan.TemplateName.ValueString()), plan.Id.ValueString()).TemplateContent(*notificationTemplateContent).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateContent",
		notificationTemplateCustomWriteError,
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

func (r *NotificationTemplateContentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *notificationTemplateContentResourceModelV1

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.NotificationsTemplatesApi.DeleteContent(ctx, data.EnvironmentId.ValueString(), management.EnumTemplateName(data.TemplateName.ValueString()), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteContent",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NotificationTemplateContentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "template_name",
			Regexp: regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(utils.EnumSliceToStringSlice(management.AllowedEnumTemplateNameEnumValues), "|"))),
		},
		{
			Label:     "template_content_id",
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

func notificationTemplateCustomWriteError(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			// Delivery method not applicable to the template
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "deliveryMethod" {
				diags.AddError(
					"The configured delivery method does not apply to the selected template.",
					"Please ensure that the delivery method (`email`, `sms`, `push`, `voice`) is applicable to the selected template.",
				)

				return diags
			}

			// Language not likely added
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "language" {
				diags.AddError(
					"The locale is not valid for the environment.",
					"Please ensure that the associated language for the locale been created with the `pingone_language` resource.",
				)

				return diags
			}

			// Not all variables set
			if message, ok := details[0].GetMessageOk(); ok && details[0].GetCode() == "REQUIRED_VALUE" {
				diags.AddError(
					"Content body is missing a required value.",
					*message,
				)

				return diags
			}

			// Custom notification content already exists
			if _, ok := details[0].GetMessageOk(); ok && details[0].GetCode() == "UNIQUENESS_VIOLATION" {
				diags.AddError(
					"Customized content for the template, locale and variant combination already exists.",
					"Please ensure that:\n\t1.\tThe notification content for the template, locale and variant is not being managed by another process and is conflicting.\n\t2.\tAny custom content for the combination has been restored to default values. See [Editing a notification](https://docs.pingidentity.com/r/en-us/pingone/p1_c_edit_notification) for more details.",
				)

				return diags
			}
		}
	}

	return diags
}

func (p *notificationTemplateContentResourceModelV1) expand(ctx context.Context) (*management.TemplateContent, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.TemplateContent{}

	if !p.Email.IsNull() && !p.Email.IsUnknown() {

		var providerPlan notificationTemplateContentEmailResourceModelV1
		diags.Append(p.Email.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		email := management.NewTemplateContentEmail(
			p.Locale.ValueString(),
			management.ENUMTEMPLATECONTENTDELIVERYMETHOD_EMAIL,
			providerPlan.Body.ValueString(),
		)

		if !p.Variant.IsNull() && !p.Variant.IsUnknown() {
			email.SetVariant(p.Variant.ValueString())
		}

		// Email specific
		if !providerPlan.From.IsNull() && !providerPlan.From.IsUnknown() {
			var fromPlan notificationTemplateContentEmailAddressResourceModelV1
			diags.Append(providerPlan.From.As(ctx, &fromPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			from := fromPlan.expandFrom()

			email.SetFrom(*from)
		}

		if !providerPlan.Subject.IsNull() && !providerPlan.Subject.IsUnknown() {
			email.SetSubject(providerPlan.Subject.ValueString())
		}

		if !providerPlan.ReplyTo.IsNull() && !providerPlan.ReplyTo.IsUnknown() {
			var replyToPlan notificationTemplateContentEmailAddressResourceModelV1
			diags.Append(providerPlan.ReplyTo.As(ctx, &replyToPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			replyTo := replyToPlan.expandReplyTo()

			email.SetReplyTo(*replyTo)
		}

		if !providerPlan.CharacterSet.IsNull() && !providerPlan.CharacterSet.IsUnknown() {
			email.SetCharset(providerPlan.CharacterSet.ValueString())
		}

		if !providerPlan.ContentType.IsNull() && !providerPlan.ContentType.IsUnknown() {
			email.SetEmailContentType(providerPlan.ContentType.ValueString())
		}

		data.TemplateContentEmail = email
	}

	if !p.Push.IsNull() && !p.Push.IsUnknown() {

		var providerPlan notificationTemplateContentPushResourceModelV1
		diags.Append(p.Push.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		push := management.NewTemplateContentPush(
			p.Locale.ValueString(),
			management.ENUMTEMPLATECONTENTDELIVERYMETHOD_PUSH,
			providerPlan.Title.ValueString(),
			providerPlan.Body.ValueString(),
		)

		if !p.Variant.IsNull() && !p.Variant.IsUnknown() {
			push.SetVariant(p.Variant.ValueString())
		}

		// Push specific
		if !providerPlan.Category.IsNull() && !providerPlan.Category.IsUnknown() {
			push.SetPushCategory(management.EnumTemplateContentPushCategory(providerPlan.Category.ValueString()))
		}

		data.TemplateContentPush = push
	}

	if !p.Sms.IsNull() && !p.Sms.IsUnknown() {

		var providerPlan notificationTemplateContentSmsResourceModelV1
		diags.Append(p.Sms.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		sms := management.NewTemplateContentSMS(
			p.Locale.ValueString(),
			management.ENUMTEMPLATECONTENTDELIVERYMETHOD_SMS,
			providerPlan.Content.ValueString(),
		)

		if !p.Variant.IsNull() && !p.Variant.IsUnknown() {
			sms.SetVariant(p.Variant.ValueString())
		}

		// SMS specific
		if !providerPlan.Sender.IsNull() && !providerPlan.Sender.IsUnknown() {
			sms.SetSender(providerPlan.Sender.ValueString())
		}

		data.TemplateContentSMS = sms
	}

	if !p.Voice.IsNull() && !p.Voice.IsUnknown() {

		var providerPlan notificationTemplateContentVoiceResourceModelV1
		diags.Append(p.Voice.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		voice := management.NewTemplateContentVoice(
			p.Locale.ValueString(),
			management.ENUMTEMPLATECONTENTDELIVERYMETHOD_VOICE,
			providerPlan.Content.ValueString(),
		)

		if !p.Variant.IsNull() && !p.Variant.IsUnknown() {
			voice.SetVariant(p.Variant.ValueString())
		}

		// Voice specific
		if !providerPlan.Type.IsNull() && !providerPlan.Type.IsUnknown() {
			voice.SetVoice(providerPlan.Type.ValueString())
		}

		data.TemplateContentVoice = voice
	}

	return &data, diags
}

func (p *notificationTemplateContentEmailAddressResourceModelV1) expandFrom() *management.TemplateContentEmailAllOfFrom {

	data := management.NewTemplateContentEmailAllOfFrom()

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.Address.IsNull() && !p.Address.IsUnknown() {
		data.SetAddress(p.Address.ValueString())
	}

	return data
}

func (p *notificationTemplateContentEmailAddressResourceModelV1) expandReplyTo() *management.TemplateContentEmailAllOfReplyTo {

	data := management.NewTemplateContentEmailAllOfReplyTo()

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.Address.IsNull() && !p.Address.IsUnknown() {
		data.SetAddress(p.Address.ValueString())
	}

	return data
}

func (p *notificationTemplateContentResourceModelV1) toState(apiObject *management.TemplateContent) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := management.TemplateContentCommon{}

	if v := apiObject.TemplateContentEmail; v != nil {
		apiObjectCommon = management.TemplateContentCommon{
			Id:             v.Id,
			Default:        v.Default,
			Locale:         v.Locale,
			DeliveryMethod: v.DeliveryMethod,
			Variant:        v.Variant,
		}
	}

	if v := apiObject.TemplateContentPush; v != nil {
		apiObjectCommon = management.TemplateContentCommon{
			Id:             v.Id,
			Default:        v.Default,
			Locale:         v.Locale,
			DeliveryMethod: v.DeliveryMethod,
			Variant:        v.Variant,
		}
	}

	if v := apiObject.TemplateContentSMS; v != nil {
		apiObjectCommon = management.TemplateContentCommon{
			Id:             v.Id,
			Default:        v.Default,
			Locale:         v.Locale,
			DeliveryMethod: v.DeliveryMethod,
			Variant:        v.Variant,
		}
	}

	if v := apiObject.TemplateContentVoice; v != nil {
		apiObjectCommon = management.TemplateContentCommon{
			Id:             v.Id,
			Default:        v.Default,
			Locale:         v.Locale,
			DeliveryMethod: v.DeliveryMethod,
			Variant:        v.Variant,
		}
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObjectCommon.GetIdOk())
	p.Locale = framework.StringOkToTF(apiObjectCommon.GetLocaleOk())
	p.Default = framework.BoolOkToTF(apiObjectCommon.GetDefaultOk())
	p.Variant = framework.StringOkToTF(apiObjectCommon.GetVariantOk())

	var d diag.Diagnostics

	p.Email, d = toStateNotificationTemplateContentEmailToTF(apiObject.TemplateContentEmail)
	diags.Append(d...)

	p.Push, d = toStateNotificationTemplateContentPushToTF(apiObject.TemplateContentPush)
	diags.Append(d...)

	p.Sms, d = toStateNotificationTemplateContentSmsToTF(apiObject.TemplateContentSMS)
	diags.Append(d...)

	p.Voice, d = toStateNotificationTemplateContentVoiceToTF(apiObject.TemplateContentVoice)
	diags.Append(d...)

	return diags
}

func toStateNotificationTemplateContentEmailToTF(apiObject *management.TemplateContentEmail) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(notificationTemplateContentEmailTFObjectTypes), diags
	}

	from, d := toStateNotificationTemplateContentEmailFromToTF(apiObject.GetFromOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(notificationTemplateContentEmailTFObjectTypes), diags
	}

	replyTo, d := toStateNotificationTemplateContentEmailReplyToToTF(apiObject.GetReplyToOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(notificationTemplateContentEmailTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"body":          framework.StringOkToTF(apiObject.GetBodyOk()),
		"from":          from,
		"subject":       framework.StringOkToTF(apiObject.GetSubjectOk()),
		"reply_to":      replyTo,
		"character_set": framework.StringOkToTF(apiObject.GetCharsetOk()),
		"content_type":  framework.StringOkToTF(apiObject.GetEmailContentTypeOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentEmailTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateNotificationTemplateContentEmailFromToTF(apiObject *management.TemplateContentEmailAllOfFrom, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(notificationTemplateContentEmailAddressTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"name":    framework.StringOkToTF(apiObject.GetNameOk()),
		"address": framework.StringOkToTF(apiObject.GetAddressOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentEmailAddressTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateNotificationTemplateContentEmailReplyToToTF(apiObject *management.TemplateContentEmailAllOfReplyTo, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(notificationTemplateContentEmailAddressTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"name":    framework.StringOkToTF(apiObject.GetNameOk()),
		"address": framework.StringOkToTF(apiObject.GetAddressOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentEmailAddressTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateNotificationTemplateContentPushToTF(apiObject *management.TemplateContentPush) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(notificationTemplateContentPushTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"category": framework.EnumOkToTF(apiObject.GetPushCategoryOk()),
		"body":     framework.StringOkToTF(apiObject.GetBodyOk()),
		"title":    framework.StringOkToTF(apiObject.GetTitleOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentPushTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateNotificationTemplateContentSmsToTF(apiObject *management.TemplateContentSMS) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(notificationTemplateContentSmsTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"content": framework.StringOkToTF(apiObject.GetContentOk()),
		"sender":  framework.EnumOkToTF(apiObject.GetSenderOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentSmsTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func toStateNotificationTemplateContentVoiceToTF(apiObject *management.TemplateContentVoice) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(notificationTemplateContentVoiceTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"content": framework.StringOkToTF(apiObject.GetContentOk()),
		"type":    framework.EnumOkToTF(apiObject.GetVoiceOk()),
	}

	returnVar, d := types.ObjectValue(notificationTemplateContentVoiceTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}
