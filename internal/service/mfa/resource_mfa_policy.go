package mfa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceMFAPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA Policies in a PingOne Environment.",

		CreateContext: resourceMFAPolicyCreate,
		ReadContext:   resourceMFAPolicyRead,
		UpdateContext: resourceMFAPolicyUpdate,
		DeleteContext: resourceMFAPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMFAPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the sign on policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "A string that specifies the MFA policy's name.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"sms": {
				Description: "SMS device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("sms.0"),
			},
			"voice": {
				Description: "Voice device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("voice.0"),
			},
			"email": {
				Description: "Email device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("email.0"),
			},
			"mobile": {
				Description: "Mobile device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled or disabled in the policy.",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"otp_failure_count": {
							Description: "An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
						},
						"otp_failure_cooldown_duration": {
							Description: "An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     2,
						},
						"otp_failure_cooldown_timeunit": {
							Description:      fmt.Sprintf("The type of time unit for `otp_failure_cooldown_duration`.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
						},
						"application": {
							Description: "Settings for a configured Mobile Application.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description:      "The mobile application's ID.  Mobile applications are configured with the `pingone_application` resource, as an OIDC `NATIVE` type.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
									},
									"push_enabled": {
										Description: "Specifies whether push notification is enabled or disabled for the policy.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"otp_enabled": {
										Description: "Specifies whether OTP authentication is enabled or disabled for the policy.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"device_authorization_enabled": {
										Description: "Specifies the enabled or disabled state of automatic MFA for native devices paired with the user, for the specified application.",
										Type:        schema.TypeBool,
										Optional:    true,
									},
									"device_authorization_extra_verification": {
										Description:      "Specifies the level of further verification when `device_authorization_enabled` is true. The PingOne platform performs an extra verification check by sending a \"silent\" push notification to the customer native application, and receives a confirmation in return.  Extra verification can be one of the following levels: `permissive`: The PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as successfully completed.  `restrictive`: The PingOne platform performs the extra verification check.The PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as failed.",
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"permissive", "restrictive"}, false)),
									},
									"auto_enrollment_enabled": {
										Description: "Set to `true` if you want the application to allow Auto Enrollment. Auto Enrollment means that the user can authenticate for the first time from an unpaired device, and the successful authentication will result in the pairing of the device for MFA.",
										Type:        schema.TypeBool,
										Optional:    true,
									},
									"integrity_detection": {
										Description:      "Controls how authentication or registration attempts should proceed if a device integrity check does not receive a response. Set the value to `permissive` if you want to allow the process to continue. Set the value to `restrictive` if you want to block the user in such situations.",
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "restrictive",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"permissive", "restrictive"}, false)),
									},
								},
							},
						},
					},
				},
			},
			"totp": {
				Description: "TOTP device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled or disabled in the policy.",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"otp_failure_count": {
							Description: "An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
						},
						"otp_failure_cooldown_duration": {
							Description: "An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     2,
						},
						"otp_failure_cooldown_timeunit": {
							Description:      fmt.Sprintf("The type of time unit for `otp_failure_cooldown_duration`.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
						},
					},
				},
			},
			"security_key": {
				Description: "Security key device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        fidoDeviceResourceSchema(),
			},
			"platform": {
				Description: "Platform device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        fidoDeviceResourceSchema(),
			},
		},
	}
}

func offlineDeviceResourceSchema(resourcePrefix string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Enabled or disabled in the policy.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"otp_lifetime_duration": {
				Description: "An integer that defines turation (number of time units) that the passcode is valid before it expires.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
			},
			"otp_lifetime_timeunit": {
				Description:      fmt.Sprintf("The type of time unit for `otp_lifetime_duration`.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
				RequiredWith:     []string{fmt.Sprintf("%s.otp_lifetime_duration", resourcePrefix), fmt.Sprintf("%s.otp_lifetime_timeunit", resourcePrefix)},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
			},
			"otp_failure_count": {
				Description: "An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
			},
			"otp_failure_cooldown_duration": {
				Description:  "An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Note that when using the \"onetime authentication\" feature, the user is not blocked after the maximum number of failures even if you specified a block duration.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				RequiredWith: []string{fmt.Sprintf("%s.otp_failure_cooldown_duration", resourcePrefix), fmt.Sprintf("%s.otp_failure_cooldown_timeunit", resourcePrefix)},
			},
			"otp_failure_cooldown_timeunit": {
				Description:      fmt.Sprintf("The type of time unit for `otp_failure_cooldown_duration`.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
				RequiredWith:     []string{fmt.Sprintf("%s.otp_failure_cooldown_duration", resourcePrefix), fmt.Sprintf("%s.otp_failure_cooldown_timeunit", resourcePrefix)},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
			},
		},
	}
}

func fidoDeviceResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Enabled or disabled in the policy.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"fido_policy_id": {
				Description:      "Specifies the FIDO policy ID. This property can be null. When null, the environment's default FIDO Policy is used.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
		},
	}
}

func resourceMFAPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	mfaPolicy := expandMFAPolicy(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.DeviceAuthenticationPolicyApi.CreateDeviceAuthenticationPolicies(ctx, d.Get("environment_id").(string)).DeviceAuthenticationPolicy(*mfaPolicy).Execute()
		},
		"CreateDeviceAuthenticationPolicies",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*mfa.DeviceAuthenticationPolicy)

	d.SetId(respObject.GetId())

	return resourceMFAPolicyRead(ctx, d, meta)
}

func resourceMFAPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneDeviceAuthenticationPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*mfa.DeviceAuthenticationPolicy)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetSmsOk(); ok {
		d.Set("sms", flattenMFAPolicyOfflineDevice(v))
	} else {
		d.Set("sms", nil)
	}

	if v, ok := respObject.GetVoiceOk(); ok {
		d.Set("voice", flattenMFAPolicyOfflineDevice(v))
	} else {
		d.Set("voice", nil)
	}

	if v, ok := respObject.GetEmailOk(); ok {
		d.Set("email", flattenMFAPolicyOfflineDevice(v))
	} else {
		d.Set("email", nil)
	}

	if v, ok := respObject.GetMobileOk(); ok {
		d.Set("mobile", flattenMFAPolicyMobile(v))
	} else {
		d.Set("mobile", nil)
	}

	if v, ok := respObject.GetTotpOk(); ok {
		d.Set("totp", flattenMFAPolicyTotp(v))
	} else {
		d.Set("totp", nil)
	}

	if v, ok := respObject.GetSecurityKeyOk(); ok {
		d.Set("security_key", flattenMFAPolicyFIDODevice(v))
	} else {
		d.Set("security_key", nil)
	}

	if v, ok := respObject.GetPlatformOk(); ok {
		d.Set("platform", flattenMFAPolicyFIDODevice(v))
	} else {
		d.Set("platform", nil)
	}

	return diags
}

func resourceMFAPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	mfaPolicy := expandMFAPolicy(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).DeviceAuthenticationPolicy(*mfaPolicy).Execute()
		},
		"UpdateMFAPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceMFAPolicyRead(ctx, d, meta)
}

func resourceMFAPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteDeviceAuthenticationPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceMFAPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/mfaPolicyID\"", d.Id())
	}

	environmentID, mfaPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(mfaPolicyID)

	resourceMFAPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandMFAPolicy(d *schema.ResourceData) *mfa.DeviceAuthenticationPolicy {

	item := mfa.NewDeviceAuthenticationPolicy(
		d.Get("name").(string),
		*expandMFAPolicyOfflineDevice(d.Get("sms").([]interface{})[0]),
		*expandMFAPolicyOfflineDevice(d.Get("voice").([]interface{})[0]),
		*expandMFAPolicyOfflineDevice(d.Get("email").([]interface{})[0]),
		*expandMFAPolicyMobileDevice(d.Get("mobile").([]interface{})[0]),
		*expandMFAPolicyTOTPDevice(d.Get("totp").([]interface{})[0]),
		*expandMFAPolicyFIDODevice(d.Get("security_key").([]interface{})[0]),
		*expandMFAPolicyFIDODevice(d.Get("platform").([]interface{})[0]),
		false,
		false,
	)

	return item
}

func expandMFAPolicyOfflineDevice(v interface{}) *mfa.DeviceAuthenticationPolicyOfflineDevice {

	obj := v.(map[string]interface{})

	otp := *mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtp(
		*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpLifeTime(int32(obj["otp_lifetime_duration"].(int)), mfa.EnumTimeUnit(obj["otp_lifetime_timeunit"].(string))),
		*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
			int32(obj["otp_failure_count"].(int)),
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(int32(obj["otp_failure_cooldown_duration"].(int)), mfa.EnumTimeUnit(obj["otp_failure_cooldown_timeunit"].(string))),
		),
	)

	item := mfa.NewDeviceAuthenticationPolicyOfflineDevice(obj["enabled"].(bool), otp)

	return item
}

func expandMFAPolicyMobileDevice(v interface{}) *mfa.DeviceAuthenticationPolicyMobile {

	obj := v.(map[string]interface{})

	otpStepSizeDuration := 30

	item := mfa.NewDeviceAuthenticationPolicyMobile(
		obj["enabled"].(bool),
		*mfa.NewDeviceAuthenticationPolicyMobileOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(obj["otp_failure_count"].(int)),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(int32(obj["otp_failure_cooldown_duration"].(int)), mfa.EnumTimeUnit(obj["otp_failure_cooldown_timeunit"].(string))),
			),
			*mfa.NewDeviceAuthenticationPolicyMobileOtpWindow(
				*mfa.NewDeviceAuthenticationPolicyMobileOtpWindowStepSize(
					int32(otpStepSizeDuration),
					mfa.ENUMTIMEUNIT_SECONDS,
				),
			),
		),
	)

	if c, ok := obj["application"].(*schema.Set); ok && c != nil && len(c.List()) > 0 && c.List()[0] != nil {

		items := make([]mfa.DeviceAuthenticationPolicyMobileApplicationsInner, 0)

		for _, cn := range c.List() {

			c2 := cn.(map[string]interface{})

			item := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInner(c2["id"].(string))

			if c3, ok := c2["push_enabled"].(bool); ok {
				item.SetPush(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPush(c3))
			}

			if c3, ok := c2["otp_enabled"].(bool); ok {
				item.SetOtp(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerOtp(c3))
			}

			deviceAuthz := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerDeviceAuthorization(c2["device_authorization_enabled"].(bool))

			if c3, ok := c2["device_authorization_extra_verification"].(string); ok && c3 != "" {
				deviceAuthz.SetExtraVerification(mfa.EnumMFADevicePolicyMobileExtraVerification(c3))
			}

			item.SetDeviceAuthorization(deviceAuthz)

			if c3, ok := c2["auto_enrollment_enabled"].(bool); ok {
				item.SetAutoEnrollment(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerAutoEnrollment(c3))
			}

			if c3, ok := c2["integrity_detection"].(string); ok && c3 != "" {
				item.SetIntegrityDetection(mfa.EnumMFADevicePolicyMobileIntegrityDetection(c3))
			}

			items = append(items, item)
		}

		item.SetApplications(items)
	}

	return item
}

func expandMFAPolicyTOTPDevice(v interface{}) *mfa.DeviceAuthenticationPolicyTotp {

	obj := v.(map[string]interface{})

	item := mfa.NewDeviceAuthenticationPolicyTotp(
		obj["enabled"].(bool),
		*mfa.NewDeviceAuthenticationPolicyTotpOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(obj["otp_failure_count"].(int)),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(int32(obj["otp_failure_cooldown_duration"].(int)), mfa.EnumTimeUnit(obj["otp_failure_cooldown_timeunit"].(string))),
			),
		),
	)

	return item
}

func expandMFAPolicyFIDODevice(v interface{}) *mfa.DeviceAuthenticationPolicyFIDODevice {

	obj := v.(map[string]interface{})

	item := mfa.NewDeviceAuthenticationPolicyFIDODevice(obj["enabled"].(bool))

	if v, ok := obj["fido_policy_id"].(string); ok {
		item.SetFidoPolicyId(v)
	}

	return item
}

func flattenMFAPolicyOfflineDevice(c *mfa.DeviceAuthenticationPolicyOfflineDevice) []map[string]interface{} {
	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetOtpOk(); ok {

		if v1, ok := v.GetLifeTimeOk(); ok {

			if v2, ok := v1.GetDurationOk(); ok {
				item["otp_lifetime_duration"] = int(*v2)
			}

			if v2, ok := v1.GetTimeUnitOk(); ok {
				item["otp_lifetime_timeunit"] = string(*v2)
			}

		}

		if v1, ok := v.GetFailureOk(); ok {

			if v2, ok := v1.GetCountOk(); ok {
				item["otp_failure_count"] = int(*v2)
			}

			if v2, ok := v1.GetCoolDownOk(); ok {

				if v3, ok := v2.GetDurationOk(); ok {
					item["otp_failure_cooldown_duration"] = int(*v3)
				}

				if v3, ok := v2.GetTimeUnitOk(); ok {
					item["otp_failure_cooldown_timeunit"] = string(*v3)
				}
			}
		}

	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenMFAPolicyMobile(c *mfa.DeviceAuthenticationPolicyMobile) []map[string]interface{} {

	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetOtpOk(); ok {

		if v1, ok := v.GetFailureOk(); ok {

			if v2, ok := v1.GetCountOk(); ok {
				item["otp_failure_count"] = int(*v2)
			}

			if v2, ok := v1.GetCoolDownOk(); ok {

				if v3, ok := v2.GetDurationOk(); ok {
					item["otp_failure_cooldown_duration"] = int(*v3)
				}

				if v3, ok := v2.GetTimeUnitOk(); ok {
					item["otp_failure_cooldown_timeunit"] = string(*v3)
				}
			}
		}
	}

	if v, ok := c.GetApplicationsOk(); ok {
		item["application"] = expandMFAPolicyMobileApplication(v)
	}

	return append(make([]map[string]interface{}, 0), item)
}

func expandMFAPolicyMobileApplication(c []mfa.DeviceAuthenticationPolicyMobileApplicationsInner) []map[string]interface{} {

	items := make([]map[string]interface{}, 0)

	for _, v := range c {

		item := map[string]interface{}{
			"id":           v.GetId(),
			"push_enabled": v.GetPush().Enabled,
			"otp_enabled":  v.GetOtp().Enabled,
		}

		if v1, ok := v.GetDeviceAuthorizationOk(); ok {

			if v2, ok := v1.GetEnabledOk(); ok {
				item["device_authorization_enabled"] = v2
			}

			if v2, ok := v1.GetExtraVerificationOk(); ok {
				item["device_authorization_extra_verification"] = v2
			}
		}

		if v1, ok := v.GetAutoEnrollmentOk(); ok {
			item["auto_enrollment_enabled"] = v1.GetEnabled()
		}

		if v1, ok := v.GetIntegrityDetectionOk(); ok {
			item["integrity_detection"] = string(*v1)
		}

		items = append(items, item)

	}

	return items

}

func flattenMFAPolicyTotp(c *mfa.DeviceAuthenticationPolicyTotp) []map[string]interface{} {

	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetOtpOk(); ok {

		if v1, ok := v.GetFailureOk(); ok {

			if v2, ok := v1.GetCountOk(); ok {
				item["otp_failure_count"] = int(*v2)
			}

			if v2, ok := v1.GetCoolDownOk(); ok {

				if v3, ok := v2.GetDurationOk(); ok {
					item["otp_failure_cooldown_duration"] = int(*v3)
				}

				if v3, ok := v2.GetTimeUnitOk(); ok {
					item["otp_failure_cooldown_timeunit"] = string(*v3)
				}
			}
		}

	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenMFAPolicyFIDODevice(c *mfa.DeviceAuthenticationPolicyFIDODevice) []map[string]interface{} {

	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetFidoPolicyIdOk(); ok {
		item["fido_policy_id"] = v
	}

	return append(make([]map[string]interface{}, 0), item)
}
