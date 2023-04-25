package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceNotificationTemplateContent() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne notification template contents for push, SMS, email and voice notifications.",

		CreateContext: resourceNotificationTemplateContentCreate,
		ReadContext:   resourceNotificationTemplateContentRead,
		UpdateContext: resourceNotificationTemplateContentUpdate,
		DeleteContext: resourceNotificationTemplateContentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNotificationTemplateContentImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to manage notification template contents in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"template_name": {
				Description:      "The ID of the template to manage localised contents for.  Options are `email_verification_admin`, `email_verification_user`, `general`, `transaction`, `verification_code_template`, `recovery_code_template`, `device_pairing`, `strong_authentication`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"email_verification_admin", "email_verification_user", "general", "transaction", "verification_code_template", "recovery_code_template", "device_pairing", "strong_authentication"}, false)),
				ForceNew:         true,
			},
			"locale": {
				Description:      "An ISO standard language code. For more information about standard language codes, see [ISO Language Code Table](http://www.lingoes.net/en/translator/langcode.htm).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(verify.FullIsoList(), false)),
				ForceNew:         true,
			},
			"default": {
				Description: "Specifies whether the template is a predefined default template.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"variant": {
				Description:      "Holds the unique user-defined name for each content variant that uses the same template + `deliveryMethod` + `locale` combination.  This property is case insensitive and has a limit of 100 characters.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 100)),
			},
			"email": {
				Description:  "A block that specifies the content settings for the `email` delivery method.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "push", "sms", "voice"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body": {
							Description:      "A string representing the email body. Email text can contain HTML but cannot be larger than 100 kB.  Use of variables is supported.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 100000)),
						},
						"from": {
							Description: "A block that specifies the sender settings for the `email` delivery method.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "The email's sender name.  If the environment uses the Ping Identity email sender, the name `PingOne` is used. You can configure other email sender names per environment.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "PingOne",
									},
									"address": {
										Description: "The sender email address. If the environment uses the Ping Identity email sender, or if the address field is empty, the address `noreply@pingidentity.com` is used.  You can configure other email sender addresses per environment.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "noreply@pingidentity.com",
									},
								},
							},
						},
						"subject": {
							Description:      "The email's subject line. Cannot exceed 256 characters. Can include variables.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
						},
						"reply_to": {
							Description: "A block that specifies the reply-to settings for the `email` delivery method.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "The email's \"reply to\" name.  If the environment uses the Ping Identity email sender, the name `PingOne` is used.  You can configure other email \"reply to\" names per environment.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"address": {
										Description: "The \"reply to\" email address.  If the environment uses the Ping Identity email sender, or if the address field is empty, the address `noreply@pingidentity.com` is used.  You can configure other email \"reply to\" addresses per environment.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
						"character_set": {
							Description: "The email's character set.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "UTF-8",
							// TODO: ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
						},
						"content_type": {
							Description: "The email's content-type.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "text/html",
							// TODO: ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
						},
					},
				},
			},
			"push": {
				Description:  "A block that specifies the content settings for the mobile `push` delivery method.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "push", "sms", "voice"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Description:      fmt.Sprintf("For Push content, you can specify what type of banner should be displayed to the user. The available options are `%s` (the banner contains both Approve and Deny buttons), `%s` (when the user clicks the banner, they are taken to an application that contains the necessary approval controls), `%s` (when the Approve button is clicked, authentication is completed and the user is taken to the relevant application).  If this parameter is not provided, the default is `%s`. Note that to use the non-default push banners, you must implement them in your application code, using the PingOne SDK. For details, see the [README for iOS](https://github.com/pingidentity/pingone-mobile-sdk-ios/#171-push-notifications-categories) and the [README for Android](https://github.com/pingidentity/pingone-mobile-sdk-android).", string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS), string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_WITHOUT_BANNER_BUTTONS), string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_APPROVE_AND_OPEN_APP), string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS), string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_WITHOUT_BANNER_BUTTONS), string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_APPROVE_AND_OPEN_APP)}, false)),
						},
						"body": {
							Description:      "The push notification text. This can include variables.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 400)),
						},
						"title": {
							Description:      "The push notification title. This can include variables.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 200)),
						},
					},
				},
			},
			"sms": {
				Description:  "A block that specifies the content settings for the `SMS` delivery method.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "push", "sms", "voice"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Description:      "The SMS text. UC-2 encoding is used for text that contains non GSM-7 characters. UC-2 encoded text cannot exceed 67 characters. GSM-7 encoded text cannot exceed 153 characters. This can include variables.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 153)),
						},
						"sender": {
							Description: "The SMS sender ID. This property can contain only alphanumeric characters and spaces, and its length cannot exceed 11 characters. In some countries, it is impossible to send an SMS with an alphanumeric sender ID. For those countries, the sender ID must be empty. For SMS recipients in specific countries, refer to Twilio's documentation on [International support for Alphanumeric Sender ID](https://support.twilio.com/hc/en-us/articles/223133767-International-support-for-Alphanumeric-Sender-ID).",
							Type:        schema.TypeString,
							Optional:    true,
							// TODO: ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 400)),
						},
					},
				},
			},
			"voice": {
				Description:  "A block that specifies the content settings for the `voice` delivery method.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"email", "push", "sms", "voice"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Description:      "The voice text to read.  This can include variables.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"type": {
							Description: "The voice type desired for the message. Out of the box options include `Man`, `Woman`, `Alice` (Twilio only), `Amazon Polly`, or your own user-defined custom string. In the case that the selected voice type is not supported by the provider in the desired locale, another voice type will be automatically selected. Additional charges may be incurred for these selections, as determined by the sender.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Alice",
						},
					},
				},
			},
		},
	}
}

func resourceNotificationTemplateContentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	templateContent, diags := expandNotificationTemplateContent(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.NotificationsTemplatesApi.CreateContent(ctx, d.Get("environment_id").(string), d.Get("template_name").(string)).TemplateContent(*templateContent).Execute()
		},
		"CreateContent",
		notificationTemplateCustomWriteError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.TemplateContent)

	if respObject.TemplateContentEmail != nil && respObject.TemplateContentEmail.GetId() != "" {
		d.SetId(respObject.TemplateContentEmail.GetId())
	} else if respObject.TemplateContentPush != nil && respObject.TemplateContentPush.GetId() != "" {
		d.SetId(respObject.TemplateContentPush.GetId())
	} else if respObject.TemplateContentSMS != nil && respObject.TemplateContentSMS.GetId() != "" {
		d.SetId(respObject.TemplateContentSMS.GetId())
	} else if respObject.TemplateContentVoice != nil && respObject.TemplateContentVoice.GetId() != "" {
		d.SetId(respObject.TemplateContentVoice.GetId())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Notification Template Content Delivery method type not supported in the provider.  Please raise an issue with the provider maintainers.",
		})

		return diags
	}

	return resourceNotificationTemplateContentRead(ctx, d, meta)
}

func resourceNotificationTemplateContentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.NotificationsTemplatesApi.ReadOneContent(ctx, d.Get("environment_id").(string), d.Get("template_name").(string), d.Id()).Execute()
		},
		"ReadOneContent",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.TemplateContent)

	if respObject.TemplateContentEmail != nil && respObject.TemplateContentEmail.GetId() != "" {

		if v, ok := respObject.TemplateContentEmail.GetDefaultOk(); ok {
			d.Set("default", v)
		} else {
			d.Set("default", nil)
		}

		if v, ok := respObject.TemplateContentEmail.GetVariantOk(); ok {
			d.Set("variant", v)
		} else {
			d.Set("variant", nil)
		}

		d.Set("email", flattenNotificationTemplateContentDeliveryMethodEmail(respObject.TemplateContentEmail))
		d.Set("sms", nil)
		d.Set("push", nil)
		d.Set("voice", nil)

	} else if respObject.TemplateContentPush != nil && respObject.TemplateContentPush.GetId() != "" {

		if v, ok := respObject.TemplateContentPush.GetDefaultOk(); ok {
			d.Set("default", v)
		} else {
			d.Set("default", nil)
		}

		if v, ok := respObject.TemplateContentPush.GetVariantOk(); ok {
			d.Set("variant", v)
		} else {
			d.Set("variant", nil)
		}

		d.Set("email", nil)
		d.Set("sms", nil)
		d.Set("push", flattenNotificationTemplateContentDeliveryMethodPush(respObject.TemplateContentPush))
		d.Set("voice", nil)

	} else if respObject.TemplateContentSMS != nil && respObject.TemplateContentSMS.GetId() != "" {

		if v, ok := respObject.TemplateContentSMS.GetDefaultOk(); ok {
			d.Set("default", v)
		} else {
			d.Set("default", nil)
		}

		if v, ok := respObject.TemplateContentSMS.GetVariantOk(); ok {
			d.Set("variant", v)
		} else {
			d.Set("variant", nil)
		}

		d.Set("email", nil)
		d.Set("sms", flattenNotificationTemplateContentDeliveryMethodSMS(respObject.TemplateContentSMS))
		d.Set("push", nil)
		d.Set("voice", nil)

	} else if respObject.TemplateContentVoice != nil && respObject.TemplateContentVoice.GetId() != "" {

		if v, ok := respObject.TemplateContentVoice.GetDefaultOk(); ok {
			d.Set("default", v)
		} else {
			d.Set("default", nil)
		}

		if v, ok := respObject.TemplateContentVoice.GetVariantOk(); ok {
			d.Set("variant", v)
		} else {
			d.Set("variant", nil)
		}

		d.Set("email", nil)
		d.Set("sms", nil)
		d.Set("push", nil)
		d.Set("voice", flattenNotificationTemplateContentDeliveryMethodVoice(respObject.TemplateContentVoice))

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Notification Template Content Delivery method type not supported in the provider.  Please raise an issue with the provider maintainers.",
		})

		return diags
	}

	return diags
}

func resourceNotificationTemplateContentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	templateContent, diags := expandNotificationTemplateContent(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.NotificationsTemplatesApi.UpdateContent(ctx, d.Get("environment_id").(string), d.Get("template_name").(string), d.Id()).TemplateContent(*templateContent).Execute()
		},
		"UpdateContent",
		notificationTemplateCustomWriteError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceNotificationTemplateContentRead(ctx, d, meta)
}

func resourceNotificationTemplateContentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.NotificationsTemplatesApi.DeleteContent(ctx, d.Get("environment_id").(string), d.Get("template_name").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteContent",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceNotificationTemplateContentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/templateName/notificationTemplateContentID\"", d.Id())
	}

	environmentID, templateName, notificationTemplateContentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("template_name", templateName)
	d.SetId(notificationTemplateContentID)

	resourceNotificationTemplateContentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func notificationTemplateCustomWriteError(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

		// Delivery method not applicable to the template
		if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "deliveryMethod" {
			diags = diag.FromErr(fmt.Errorf("The configured delivery method does not apply to the selected template."))

			return diags
		}

		// Language not likely added
		if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "language" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "The locale is not valid for the environment.",
				Detail:   "Please ensure that the associated language for the locale been created with the `pingone_language` resource.",
			})

			return diags
		}

		// Not all variables set
		if message, ok := details[0].GetMessageOk(); ok && details[0].GetCode() == "REQUIRED_VALUE" {
			diags = diag.FromErr(fmt.Errorf(*message))

			return diags
		}

		// Custom notification content already exists
		if _, ok := details[0].GetMessageOk(); ok && details[0].GetCode() == "UNIQUENESS_VIOLATION" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Customized content for the template, locale and variant combination already exists.",
				Detail:   "Please ensure that:\n\t1.\tThe notification content for the template, locale and variant is not being managed by another process and is conflicting.\n\t2.\tAny custom content for the combination has been restored to default values. See [Editing a notification](https://docs.pingidentity.com/r/en-us/pingone/p1_c_edit_notification) for more details.",
			})

			return diags
		}
	}

	return nil
}

func expandNotificationTemplateContent(d *schema.ResourceData) (*management.TemplateContent, diag.Diagnostics) {
	var diags diag.Diagnostics

	templateContentRequest := &management.TemplateContent{}

	if v, ok := d.GetOk("email"); ok {
		var templateContent *management.TemplateContentEmail
		common := management.NewTemplateContentCommon(d.Get("locale").(string), management.ENUMTEMPLATECONTENTDELIVERYMETHOD_EMAIL)

		if v1, ok := d.Get("variant").(string); ok {
			common.SetVariant(v1)
		}

		templateContent, diags = expandNotificationTemplateContentEmail(v.([]interface{}), common)
		templateContentRequest.TemplateContentEmail = templateContent
	}

	if v, ok := d.GetOk("push"); ok {
		var templateContent *management.TemplateContentPush
		common := management.NewTemplateContentCommon(d.Get("locale").(string), management.ENUMTEMPLATECONTENTDELIVERYMETHOD_PUSH)

		if v1, ok := d.Get("variant").(string); ok {
			common.SetVariant(v1)
		}

		templateContent, diags = expandNotificationTemplateContentPush(v.([]interface{}), common)
		templateContentRequest.TemplateContentPush = templateContent
	}

	if v, ok := d.GetOk("sms"); ok {
		var templateContent *management.TemplateContentSMS
		common := management.NewTemplateContentCommon(d.Get("locale").(string), management.ENUMTEMPLATECONTENTDELIVERYMETHOD_SMS)

		if v1, ok := d.Get("variant").(string); ok {
			common.SetVariant(v1)
		}

		templateContent, diags = expandNotificationTemplateContentSMS(v.([]interface{}), common)
		templateContentRequest.TemplateContentSMS = templateContent
	}

	if v, ok := d.GetOk("voice"); ok {
		var templateContent *management.TemplateContentVoice
		common := management.NewTemplateContentCommon(d.Get("locale").(string), management.ENUMTEMPLATECONTENTDELIVERYMETHOD_VOICE)

		if v1, ok := d.Get("variant").(string); ok {
			common.SetVariant(v1)
		}

		templateContent, diags = expandNotificationTemplateContentVoice(v.([]interface{}), common)
		templateContentRequest.TemplateContentVoice = templateContent
	}
	if diags.HasError() {
		return nil, diags
	}

	return templateContentRequest, diags
}

func expandNotificationTemplateContentEmail(d []interface{}, common *management.TemplateContentCommon) (*management.TemplateContentEmail, diag.Diagnostics) {
	var diags diag.Diagnostics

	var templateContent management.TemplateContentEmail

	if len(d) > 0 && d[0] != nil {
		options := d[0].(map[string]interface{})

		templateContent = *management.NewTemplateContentEmail(common.GetLocale(), common.GetDeliveryMethod(), options["body"].(string))

		// From common
		if v1, ok := common.GetVariantOk(); ok {
			templateContent.SetVariant(*v1)
		}

		// From the block
		if v, ok := options["from"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			fromOptions := v[0].(map[string]interface{})

			from := management.NewTemplateContentEmailAllOfFrom()

			if v1, ok := fromOptions["name"].(string); ok && v1 != "" {
				from.SetName(v1)
			}

			if v1, ok := fromOptions["address"].(string); ok && v1 != "" {
				from.SetAddress(v1)
			}

			templateContent.SetFrom(*from)
		}

		if v, ok := options["subject"].(string); ok && v != "" {
			templateContent.SetSubject(v)
		}

		if v, ok := options["reply_to"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			replyToOptions := v[0].(map[string]interface{})

			replyTo := management.NewTemplateContentEmailAllOfReplyTo()

			if v1, ok := replyToOptions["name"].(string); ok && v1 != "" {
				replyTo.SetName(v1)
			}

			if v1, ok := replyToOptions["address"].(string); ok && v1 != "" {
				replyTo.SetAddress(v1)
			}

			templateContent.SetReplyTo(*replyTo)
		}

		if v, ok := options["character_set"].(string); ok && v != "" {
			templateContent.SetCharset(v)
		}

		if v, ok := options["content_type"].(string); ok && v != "" {
			templateContent.SetEmailContentType(v)
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s notification template content options not available", common.GetDeliveryMethod()),
		})

		return nil, diags
	}

	return &templateContent, diags
}

func expandNotificationTemplateContentPush(d []interface{}, common *management.TemplateContentCommon) (*management.TemplateContentPush, diag.Diagnostics) {
	var diags diag.Diagnostics

	var templateContent management.TemplateContentPush

	if len(d) > 0 && d[0] != nil {
		options := d[0].(map[string]interface{})

		templateContent = *management.NewTemplateContentPush(common.GetLocale(), common.GetDeliveryMethod(), options["title"].(string), options["body"].(string))

		// From common
		if v1, ok := common.GetVariantOk(); ok {
			templateContent.SetVariant(*v1)
		}

		// From the block
		if v, ok := options["category"].(string); ok && v != "" {
			templateContent.SetPushCategory(management.EnumTemplateContentPushCategory(v))
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s notification template content options not available", common.GetDeliveryMethod()),
		})

		return nil, diags
	}

	return &templateContent, diags
}

func expandNotificationTemplateContentSMS(d []interface{}, common *management.TemplateContentCommon) (*management.TemplateContentSMS, diag.Diagnostics) {
	var diags diag.Diagnostics

	var templateContent management.TemplateContentSMS

	if len(d) > 0 && d[0] != nil {
		options := d[0].(map[string]interface{})

		templateContent = *management.NewTemplateContentSMS(common.GetLocale(), common.GetDeliveryMethod(), options["content"].(string))

		// From common
		if v1, ok := common.GetVariantOk(); ok {
			templateContent.SetVariant(*v1)
		}

		// From the block
		if v, ok := options["sender"].(string); ok && v != "" {
			templateContent.SetSender(v)
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s notification template content options not available", common.GetDeliveryMethod()),
		})

		return nil, diags
	}

	return &templateContent, diags
}

func expandNotificationTemplateContentVoice(d []interface{}, common *management.TemplateContentCommon) (*management.TemplateContentVoice, diag.Diagnostics) {
	var diags diag.Diagnostics

	var templateContent management.TemplateContentVoice

	if len(d) > 0 && d[0] != nil {
		options := d[0].(map[string]interface{})

		templateContent = *management.NewTemplateContentVoice(common.GetLocale(), common.GetDeliveryMethod(), options["content"].(string))

		// From common
		if v1, ok := common.GetVariantOk(); ok {
			templateContent.SetVariant(*v1)
		}

		// From the block
		if v, ok := options["type"].(string); ok && v != "" {
			templateContent.SetVoice(v)
		}

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s notification template content options not available", common.GetDeliveryMethod()),
		})

		return nil, diags
	}

	return &templateContent, diags
}

func flattenNotificationTemplateContentDeliveryMethodEmail(d *management.TemplateContentEmail) []interface{} {
	// Required
	item := map[string]interface{}{
		"body":    d.GetBody(),
		"subject": d.GetSubject(),
	}

	// Optional
	if v, ok := d.GetFromOk(); ok {

		from := map[string]interface{}{
			"name":    nil,
			"address": nil,
		}

		if c, ok := v.GetNameOk(); ok {
			from["name"] = c
		}

		if c, ok := v.GetAddressOk(); ok {
			from["address"] = c
		}

		fromItems := make([]interface{}, 0)
		item["from"] = append(fromItems, from)
	} else {
		item["from"] = nil
	}

	if v, ok := d.GetReplyToOk(); ok {
		replyTo := map[string]interface{}{
			"name":    nil,
			"address": nil,
		}

		if c, ok := v.GetNameOk(); ok {
			replyTo["name"] = c
		}

		if c, ok := v.GetAddressOk(); ok {
			replyTo["address"] = c
		}

		replyToItems := make([]interface{}, 0)
		item["reply_to"] = append(replyToItems, replyTo)
	} else {
		item["reply_to"] = nil
	}

	if v, ok := d.GetCharsetOk(); ok {
		item["character_set"] = v
	} else {
		item["character_set"] = nil
	}

	if v, ok := d.GetEmailContentTypeOk(); ok {
		item["content_type"] = v
	} else {
		item["content_type"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenNotificationTemplateContentDeliveryMethodPush(d *management.TemplateContentPush) []interface{} {
	// Required
	item := map[string]interface{}{
		"body":  d.GetBody(),
		"title": d.GetTitle(),
	}

	// Optional
	if v, ok := d.GetPushCategoryOk(); ok {
		item["category"] = string(*v)
	} else {
		item["category"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenNotificationTemplateContentDeliveryMethodSMS(d *management.TemplateContentSMS) []interface{} {
	// Required
	item := map[string]interface{}{
		"content": d.GetContent(),
	}

	// Optional
	if v, ok := d.GetSenderOk(); ok {
		item["sender"] = v
	} else {
		item["sender"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenNotificationTemplateContentDeliveryMethodVoice(d *management.TemplateContentVoice) []interface{} {
	// Required
	item := map[string]interface{}{
		"content": d.GetContent(),
	}

	// Optional
	if v, ok := d.GetVoiceOk(); ok {
		item["type"] = v
	} else {
		item["type"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}
