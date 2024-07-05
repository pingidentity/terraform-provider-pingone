package mfa

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceMFAPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA Policies in a PingOne Environment.",

		DeprecationMessage: "Use the pingone_mfa_device_policy resource going forward.  This resource will be removed in the next major release.",

		CreateContext: resourceMFAPolicyCreate,
		ReadContext:   resourceMFAPolicyRead,
		UpdateContext: resourceMFAPolicyUpdate,
		DeleteContext: resourceMFAPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMFAPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the MFA device policy in.",
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
			"device_selection": {
				Description:      fmt.Sprintf("A string that defines the device selection method. Options are `%s` (this is the default setting for new environments), `%s` and `%s`.", string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST), string(mfa.ENUMMFADEVICEPOLICYSELECTION_PROMPT_TO_SELECT), string(mfa.ENUMMFADEVICEPOLICYSELECTION_ALWAYS_DISPLAY_DEVICES)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST), string(mfa.ENUMMFASETTINGSDEVICESELECTION_PROMPT_TO_SELECT), string(mfa.ENUMMFADEVICEPOLICYSELECTION_ALWAYS_DISPLAY_DEVICES)}, false)),
			},
			"new_device_notification": {
				Description:      fmt.Sprintf("A string that defines whether a user should be notified if a new authentication method has been added to their account. Options are `%s` (the default), `%s` and `%s`.", string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE), string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_EMAIL_THEN_SMS), string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_SMS_THEN_EMAIL)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE), string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_EMAIL_THEN_SMS), string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_SMS_THEN_EMAIL)}, false)),
			},
			"sms": {
				Description: "SMS OTP authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("sms.0"),
			},
			"voice": {
				Description: "Voice OTP authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("voice.0"),
			},
			"email": {
				Description: "Email OTP authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        offlineDeviceResourceSchema("email.0"),
			},
			"mobile": {
				Description: "Mobile authenticator device policy settings.  This factor requires embedding the PingOne MFA SDK into a customer facing mobile application, and configuring as a Native application using the `pingone_application` resource.",
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
									"pairing_disabled": {
										Description: "You can set this parameter to `true` to prevent users from pairing new devices with the relevant method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"push_enabled": {
										Description: "Specifies whether push notification is enabled or disabled for the policy.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"push_timeout_duration": {
										Description:      "An integer that defines the amount of time (in seconds) a user has to respond to a push notification before it expires. Minimum is 40 seconds and maximum is 150 seconds. If this parameter is not provided, the duration is set to 40 seconds.",
										Type:             schema.TypeInt,
										Optional:         true,
										Default:          40,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(40, 150)),
									},
									"push_timeout_timeunit": {
										Description: "The time unit for the `push_timeout_duration` parameter. Currently, the only permitted value is `SECONDS`.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"pairing_key_lifetime_duration": {
										Description: "The amount of time an issued pairing key can be used until it expires. Minimum is 1 minute and maximum is 48 hours. If this parameter is not provided, the duration is set to 10 minutes.",
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     10,
									},
									"pairing_key_lifetime_timeunit": {
										Description:      fmt.Sprintf("The time unit for the `pairing_key_lifetime_duration` parameter.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_HOURS), string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_MINUTES)),
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_HOURS), string(mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_MINUTES)}, false)),
										Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
									},
									"push_limit_count": {
										Description:      "The number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is 1 and the maximum value is 50. If this parameter is not provided, the default value is 5.",
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 50)),
										Default:          5,
									},
									"push_limit_lock_duration": {
										Description:      "The length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is 1 minute and the maximum value is 120 minutes. If this parameter is not provided, the default value is 30 minutes.",
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
										Default:          30,
									},
									"push_limit_lock_duration_timeunit": {
										Description:      fmt.Sprintf("The time unit for the `push_limit_lock_duration` parameter.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
										Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
									},
									"push_limit_time_period_duration": {
										Description:      "The time period in which the push notifications are counted towards the defined limit. The minimum value is 1 minute and the maximum value is 120 minutes. If this parameter is not provided, the default value is 10 minutes.",
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
										Default:          10,
									},
									"push_limit_time_period_timeunit": {
										Description:      fmt.Sprintf("The time unit for the `push_limit_time_period_duration` parameter.  Options are `%s` or `%s`.", string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)),
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMTIMEUNIT_MINUTES), string(mfa.ENUMTIMEUNIT_SECONDS)}, false)),
										Default:          string(mfa.ENUMTIMEUNIT_MINUTES),
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
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"permissive", "restrictive"}, false)),
									},
								},
							},
						},
					},
				},
			},
			"totp": {
				Description: "TOTP authenticator policy settings.",
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
						"pairing_disabled": {
							Description: "You can set this parameter to `true` to prevent users from pairing new devices with the relevant method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
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
			"fido2": {
				Description: "FIDO2 device authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"security_key",
					"platform",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled or disabled in the policy.",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"pairing_disabled": {
							Description: "You can set this parameter to `true` to prevent users from pairing new devices with the relevant method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"fido2_policy_id": {
							Description:      "Specifies the UUID that represents the FIDO2 policy in PingOne. This property can be null. When null, the environment's default FIDO2 Policy is used.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
					},
				},
			},
			"security_key": {
				Description: "**Deprecation Notice** The `security_key` FIDO device type is deprecated and needs to be replaced with the `fido2` device type.  `security_key` will not be configurable for newly created environments, or existing environments that have not had their environment upgraded to use the latest FIDO2 policies. Security key (FIDO2) authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"fido2",
				},
				Deprecated: "The `security_key` FIDO device type is deprecated and needs to be replaced with the `fido2` device type.  `security_key` will not be configurable for newly created environments, or existing environments that have not had their environment upgraded to use the latest FIDO2 policies.",
				Elem:       fidoDeviceResourceSchema(),
			},
			"platform": {
				Description: "**Deprecation Notice** The `platform` FIDO device type is deprecated and needs to be replaced with the `fido2` device type.  `platform` will not be configurable for newly created environments, or existing environments that have not had their environment upgraded to use the latest FIDO2 policies. Platform biometrics authentication policy settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ConflictsWith: []string{
					"fido2",
				},
				Deprecated: "The `platform` FIDO device type is deprecated and needs to be replaced with the `fido2` device type.  `platform` will not be configurable for newly created environments, or existing environments that have not had their environment upgraded to use the latest FIDO2 policies.",
				Elem:       fidoDeviceResourceSchema(),
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
			"pairing_disabled": {
				Description: "You can set this parameter to `true` to prevent users from pairing new devices with the relevant method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
			"pairing_disabled": {
				Description: "You can set this parameter to `true` to prevent users from pairing new devices with the relevant method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Resource deprecated",
		Detail:   "The `pingone_mfa_policy` resource is deprecated and cannot be used to create new MFA device policies. Use the `pingone_mfa_device_policy` resource instead.",
	})

	return diags
}

func resourceMFAPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, p1Client.API.ManagementAPIClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"ReadOneDeviceAuthenticationPolicy",
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

	respObject := resp.(*mfa.DeviceAuthenticationPolicy)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetAuthenticationOk(); ok {
		d.Set("device_selection", v.GetDeviceSelection())
	} else {
		d.Set("device_selection", nil)
	}

	if v, ok := respObject.GetNewDeviceNotificationOk(); ok {
		d.Set("new_device_notification", string(*v))
	} else {
		d.Set("new_device_notification", nil)
	}

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

	if v, ok := respObject.GetFido2Ok(); ok {
		d.Set("fido2", flattenMFAPolicyFIDO2Device(v))
	} else {
		d.Set("fido2", nil)
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

	managementApiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	mfaPolicy, diags := expandMFAPolicy(ctx, managementApiClient, d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).DeviceAuthenticationPolicy(*mfaPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, p1Client.API.ManagementAPIClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"UpdateMFAPolicy",
		mfaPolicyCreateUpdateCustomErrorHandler,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceMFAPolicyRead(ctx, d, meta)
}

func resourceMFAPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	timeout := 30 * time.Second
	_, diags = sdk.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, p1Client.API.ManagementAPIClient, d.Get("environment_id").(string), nil, fR, fErr)
		},
		"DeleteDeviceAuthenticationPolicy",
		sdk.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1Error *model.P1Error) bool {

			if p1Error != nil && r.StatusCode == http.StatusConflict {
				if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
					if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "signOnPolicies" {
						tflog.Warn(ctx, "Conflict detected when deleting MFA policy.  This is likely due to the policy being referenced by a sign-on policy.  Retrying..")

						return true
					}
				}
			}

			return false
		},
		timeout,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceMFAPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "mfa_device_policy_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.SetId(attributes["mfa_device_policy_id"])

	resourceMFAPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

var mfaPolicyCreateUpdateCustomErrorHandler = func(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Invalid FIDO2 combination
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if message, ok := details[0].GetMessageOk(); ok {
			m, _ := regexp.MatchString("^Deprecated Fido Settings", *message)

			if m {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid FIDO device types.",
					Detail:   "The `security_key` and `platform` parameters are deprecated and cannot be configured for new environments or environments that have been upgraded to use the latest FIDO2 policies.  Use the `fido2` parameter instead.",
				})

				return diags
			}
		}
	}

	return nil
}

func expandMFAPolicyPost(ctx context.Context, apiClient *management.APIClient, d *schema.ResourceData) (*mfa.DeviceAuthenticationPolicyPost, diag.Diagnostics) {
	var diags diag.Diagnostics
	mfaPolicyPost := mfa.DeviceAuthenticationPolicyPost{}
	mfaPolicyPost.DeviceAuthenticationPolicy, diags = expandMFAPolicy(ctx, apiClient, d)

	return &mfaPolicyPost, diags
}

func expandMFAPolicy(ctx context.Context, apiClient *management.APIClient, d *schema.ResourceData) (*mfa.DeviceAuthenticationPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	mobile, diags := expandMFAPolicyMobileDevice(d.Get("mobile").([]interface{})[0], ctx, apiClient, d.Get("environment_id").(string))
	if diags.HasError() {
		return nil, diags
	}

	item := mfa.NewDeviceAuthenticationPolicy(
		d.Get("name").(string),
		*expandMFAPolicyOfflineDevice(d.Get("sms").([]interface{})[0]),
		*expandMFAPolicyOfflineDevice(d.Get("voice").([]interface{})[0]),
		*expandMFAPolicyOfflineDevice(d.Get("email").([]interface{})[0]),
		*mobile,
		*expandMFAPolicyTOTPDevice(d.Get("totp").([]interface{})[0]),
		false,
		false,
	)

	if v, ok := d.GetOk("fido2"); ok {
		item.SetFido2(*expandMFAPolicyFIDO2Device(v.([]interface{})[0]))
	}

	if v, ok := d.GetOk("security_key"); ok {
		item.SetSecurityKey(*expandMFAPolicyFIDODevice(v.([]interface{})[0]))
	}

	if v, ok := d.GetOk("platform"); ok {
		item.SetPlatform(*expandMFAPolicyFIDODevice(v.([]interface{})[0]))
	}

	if v, ok := d.GetOk("device_selection"); ok {
		item.SetAuthentication(*mfa.NewDeviceAuthenticationPolicyAuthentication(mfa.EnumMFADevicePolicySelection(v.(string))))
	}

	if v, ok := d.GetOk("new_device_notification"); ok {
		item.SetNewDeviceNotification(mfa.EnumMFADevicePolicyNewDeviceNotification(v.(string)))
	}

	return item, diags
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

	if v, ok := obj["pairing_disabled"]; ok {
		item.SetPairingDisabled(v.(bool))
	}

	return item
}

func expandMFAPolicyMobileDevice(v interface{}, ctx context.Context, apiClient *management.APIClient, environmentID string) (*mfa.DeviceAuthenticationPolicyMobile, diag.Diagnostics) {
	var diags diag.Diagnostics

	obj := v.(map[string]interface{})

	item := mfa.NewDeviceAuthenticationPolicyMobile(
		obj["enabled"].(bool),
		*mfa.NewDeviceAuthenticationPolicyMobileOtp(
			*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailure(
				int32(obj["otp_failure_count"].(int)),
				*mfa.NewDeviceAuthenticationPolicyOfflineDeviceOtpFailureCoolDown(int32(obj["otp_failure_cooldown_duration"].(int)), mfa.EnumTimeUnit(obj["otp_failure_cooldown_timeunit"].(string))),
			),
		),
	)

	if c, ok := obj["application"].(*schema.Set); ok && c != nil && len(c.List()) > 0 && c.List()[0] != nil {

		items := make([]mfa.DeviceAuthenticationPolicyMobileApplicationsInner, 0)

		for _, cn := range c.List() {

			c2 := cn.(map[string]interface{})

			item := *mfa.NewDeviceAuthenticationPolicyMobileApplicationsInner(c2["id"].(string))

			application, diags := checkApplicationForMobileAppSDKv2(ctx, apiClient, environmentID, c2["id"].(string))
			if diags.HasError() {
				return nil, diags
			}

			if c3, ok := c2["pairing_disabled"]; ok {
				item.SetPairingDisabled(c3.(bool))
			}

			if c3, ok := c2["push_enabled"].(bool); ok {
				item.SetPush(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPush(c3))
			}

			if c3, ok := c2["push_timeout_duration"].(int); ok {
				item.SetPushTimeout(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushTimeout(int32(c3), mfa.ENUMTIMEUNITPUSHTIMEOUT_SECONDS))
			}

			pairingKeyLifetimeDuration, pairingKeyLifetimeDurationOk := c2["pairing_key_lifetime_duration"].(int)
			pairingKeyLifetimeTimeunit, pairingKeyLifetimeTimeunitOk := c2["pairing_key_lifetime_timeunit"].(string)
			if pairingKeyLifetimeDurationOk && pairingKeyLifetimeTimeunitOk {
				item.SetPairingKeyLifetime(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPairingKeyLifetime(int32(pairingKeyLifetimeDuration), mfa.EnumTimeUnitPairingKeyLifetime(pairingKeyLifetimeTimeunit)))
			}

			pushLimit := mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimit()

			if c3, ok := c2["push_limit_count"].(int); ok {
				pushLimit.SetCount(int32(c3))
			}

			lockDuration, lockDurationOk := c2["push_limit_lock_duration"].(int)
			lockDurationTimeunit, lockDurationTimeunitOk := c2["push_limit_lock_duration_timeunit"].(string)
			if lockDurationOk && lockDurationTimeunitOk {
				pushLimit.SetLockDuration(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimitLockDuration(int32(lockDuration), mfa.EnumTimeUnit(lockDurationTimeunit)))
			}

			timePeriodDuration, timePeriodOk := c2["push_limit_time_period_duration"].(int)
			timePeriodTimeunit, timePeriodTimeunitOk := c2["push_limit_time_period_timeunit"].(string)
			if timePeriodOk && timePeriodTimeunitOk {
				pushLimit.SetTimePeriod(*mfa.NewDeviceAuthenticationPolicyMobileApplicationsInnerPushLimitTimePeriod(int32(timePeriodDuration), mfa.EnumTimeUnit(timePeriodTimeunit)))
			}

			item.SetPushLimit(*pushLimit)

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

			c3, ok := c2["integrity_detection"].(string)
			if application.GetMobile().IntegrityDetection.GetMode() == management.ENUMENABLEDSTATUS_ENABLED {

				if ok && c3 != "" {
					item.SetIntegrityDetection(mfa.EnumMFADevicePolicyMobileIntegrityDetection(c3))
				} else {
					// error - this must be set
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Integrity detection (`mobile.application.integrity_detection`) must be set when the Application resource has integrity detection enabled",
						Detail:   "The referenced mobile application (`mobile.application.id`) has integrity detection enabled. This policy must specify the level of integrity detection in the `mobile.application.integrity_detection` parameter.",
					})
					return nil, diags
				}
			} else {
				if ok && c3 != "" {
					// error - this has no effect
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Integrity detection (`mobile.application.integrity_detection`) has no effect when the Application resource has integrity detection disabled",
						Detail:   "The referenced mobile application (`mobile.application.id`) has integrity detection disabled. Setting the `mobile.application.integrity_detection` parameter has no effect.",
					})
					return nil, diags
				}
			}

			items = append(items, item)
		}

		item.SetApplications(items)
	}

	return item, diags
}

func checkApplicationForMobileAppSDKv2(ctx context.Context, apiClient *management.APIClient, environmentID, appID string) (*management.ApplicationOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ApplicationsApi.ReadOneApplication(ctx, environmentID, appID).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"ReadOneApplication",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return nil, diags
	}

	if resp == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Application referenced in `mobile.application.id` does not exist",
		})
		return nil, diags
	}

	respObject := resp.(*management.ReadOneApplication200Response)

	var oidcObject *management.ApplicationOIDC

	// check if oidc
	if respObject.ApplicationOIDC == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Application referenced in `mobile.application.id` is not of type OIDC",
			Detail:   "To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		})
		return nil, diags
	} else {
		oidcObject = respObject.ApplicationOIDC
	}

	// check if native
	if respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_NATIVE_APP && respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_CUSTOM_APP {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Application referenced in `mobile.application.id` is OIDC, but is not the required `Native` OIDC application type",
			Detail:   "To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		})
		return nil, diags
	}

	// check if mobile set and package/bundle set
	if _, ok := respObject.ApplicationOIDC.GetMobileOk(); !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Application referenced in `mobile.application.id` does not contain mobile application configuration",
			Detail:   "To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		})
		return nil, diags
	}

	if v, ok := respObject.ApplicationOIDC.GetMobileOk(); ok {

		_, bundleIDOk := v.GetBundleIdOk()
		_, packageNameOk := v.GetPackageNameOk()

		if !bundleIDOk && !packageNameOk {

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Application referenced in `mobile.application.id` does not contain mobile application configuration",
				Detail:   "To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
			})
			return nil, diags
		}
	}

	return oidcObject, diags
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

	if v, ok := obj["pairing_disabled"]; ok {
		item.SetPairingDisabled(v.(bool))
	}

	return item
}

func expandMFAPolicyFIDODevice(v interface{}) *mfa.DeviceAuthenticationPolicyFIDODevice {

	obj := v.(map[string]interface{})

	item := mfa.NewDeviceAuthenticationPolicyFIDODevice(obj["enabled"].(bool))

	if v, ok := obj["fido_policy_id"].(string); ok {
		item.SetFidoPolicyId(v)
	}

	if v, ok := obj["pairing_disabled"]; ok {
		item.SetPairingDisabled(v.(bool))
	}

	return item
}

func expandMFAPolicyFIDO2Device(v interface{}) *mfa.DeviceAuthenticationPolicyFido2 {

	obj := v.(map[string]interface{})

	item := mfa.NewDeviceAuthenticationPolicyFido2(obj["enabled"].(bool))

	if v, ok := obj["fido2_policy_id"].(string); ok {
		item.SetFido2PolicyId(v)
	}

	if v, ok := obj["pairing_disabled"]; ok {
		item.SetPairingDisabled(v.(bool))
	}

	return item
}

func flattenMFAPolicyOfflineDevice(c *mfa.DeviceAuthenticationPolicyOfflineDevice) []map[string]interface{} {
	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetPairingDisabledOk(); ok {
		item["pairing_disabled"] = *v
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

		if v1, ok := v.GetPairingDisabledOk(); ok {
			item["pairing_disabled"] = *v1
		}

		if v1, ok := v.GetDeviceAuthorizationOk(); ok {

			if v2, ok := v1.GetEnabledOk(); ok {
				item["device_authorization_enabled"] = v2
			}

			if v2, ok := v1.GetExtraVerificationOk(); ok {
				item["device_authorization_extra_verification"] = v2
			}
		}

		if v1, ok := v.GetPushTimeoutOk(); ok {

			if v2, ok := v1.GetDurationOk(); ok {
				item["push_timeout_duration"] = v2
			}

			if v2, ok := v1.GetTimeUnitOk(); ok {
				item["push_timeout_timeunit"] = v2
			}
		}

		if v1, ok := v.GetPairingKeyLifetimeOk(); ok {
			if v2, ok := v1.GetDurationOk(); ok {
				item["pairing_key_lifetime_duration"] = v2
			}

			if v2, ok := v1.GetTimeUnitOk(); ok {
				item["pairing_key_lifetime_timeunit"] = v2
			}
		}

		if v1, ok := v.GetPushLimitOk(); ok {

			if v2, ok := v1.GetCountOk(); ok {
				item["push_limit_count"] = v2
			}

			if v2, ok := v1.GetLockDurationOk(); ok {
				if v3, ok := v2.GetDurationOk(); ok {
					item["push_limit_lock_duration"] = v3
				}

				if v3, ok := v2.GetTimeUnitOk(); ok {
					item["push_limit_lock_duration_timeunit"] = v3
				}
			}

			if v2, ok := v1.GetTimePeriodOk(); ok {
				if v3, ok := v2.GetDurationOk(); ok {
					item["push_limit_time_period_duration"] = v3
				}

				if v3, ok := v2.GetTimeUnitOk(); ok {
					item["push_limit_time_period_timeunit"] = v3
				}
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

	if v, ok := c.GetPairingDisabledOk(); ok {
		item["pairing_disabled"] = *v
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

	if v, ok := c.GetPairingDisabledOk(); ok {
		item["pairing_disabled"] = *v
	}

	if v, ok := c.GetFidoPolicyIdOk(); ok {
		item["fido_policy_id"] = v
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenMFAPolicyFIDO2Device(c *mfa.DeviceAuthenticationPolicyFido2) []map[string]interface{} {

	item := map[string]interface{}{
		"enabled": c.GetEnabled(),
	}

	if v, ok := c.GetPairingDisabledOk(); ok {
		item["pairing_disabled"] = *v
	}

	if v, ok := c.GetFido2PolicyIdOk(); ok && *v != "" {
		item["fido2_policy_id"] = v
	}

	return append(make([]map[string]interface{}, 0), item)
}
