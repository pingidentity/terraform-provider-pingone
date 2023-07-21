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

func ResourceMFASettings() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne Environment's MFA Settings",

		CreateContext: resourceMFASettingsCreate,
		ReadContext:   resourceMFASettingsRead,
		UpdateContext: resourceMFASettingsUpdate,
		DeleteContext: resourceMFASettingsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMFASettingsImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the sign on policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"phone_extensions_enabled": {
				Description: "A boolean when set to `true` allows one-time passwords to be delivered via voice to phone numbers that include extensions. Set to `false` to disable support for extensions.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"pairing": {
				Description: "An object that contains pairing settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_allowed_devices": {
							Description:      "An integer that defines the maximum number of MFA devices each user can have. This can be any number up to 15. The default value is 5.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          5,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 15)),
						},
						"pairing_key_format": {
							Description:      fmt.Sprintf("String that controls the type of pairing key issued. The valid values are `%s` (12-digit key) and `%s` (16-character alphanumeric key).", string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_NUMERIC), string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_ALPHANUMERIC)),
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_NUMERIC), string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_ALPHANUMERIC)}, false)),
						},
					},
				},
			},
			"lockout": {
				Description: "An object that contains lockout settings.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_count": {
							Description:      "An integer that defines the maximum number of incorrect authentication attempts before the account is locked.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"duration_seconds": {
							Description:      "An integer that defines the number of seconds to keep the account in a locked state.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
					},
				},
			},
			"authentication": {
				Description: "**This property is deprecated.**  Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource. An object that contains the device selection settings.",
				Type:        schema.TypeList,
				Deprecated:  "Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource.",
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_selection": {
							Description:      fmt.Sprintf("**This property is deprecated.**  Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource.  A string that defines the device selection method. Options are `%s` (this is the default setting for new environments) and `%s`.", string(mfa.ENUMMFASETTINGSDEVICESELECTION_DEFAULT_TO_FIRST), string(mfa.ENUMMFASETTINGSDEVICESELECTION_PROMPT_TO_SELECT)),
							Type:             schema.TypeString,
							Deprecated:       "Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMMFASETTINGSDEVICESELECTION_DEFAULT_TO_FIRST), string(mfa.ENUMMFASETTINGSDEVICESELECTION_PROMPT_TO_SELECT)}, false)),
						},
					},
				},
			},
		},
	}
}

func resourceMFASettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	mfaSettings := *mfa.NewMFASettings(expandMFASettingsPairing(d.Get("pairing").([]interface{})))

	if v, ok := d.GetOk("authentication"); ok {
		mfaSettings.SetAuthentication(expandMFASettingsAuthentication(v.([]interface{})))
	}

	if v, ok := d.GetOk("lockout"); ok {
		mfaSettings.SetLockout(expandMFASettingsLockout(v.([]interface{})))
	}

	if v, ok := d.GetOk("phone_extensions_enabled"); ok {
		mfaSettings.SetPhoneExtensions(expandMFASettingsPhoneExtensions(v))
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.MFASettingsApi.UpdateMFASettings(ctx, d.Get("environment_id").(string)).MFASettings(mfaSettings).Execute()
		},
		"UpdateMFASettings",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*mfa.MFASettings)

	d.SetId(*respObject.GetEnvironment().Id)

	return resourceMFASettingsRead(ctx, d, meta)
}

func resourceMFASettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.MFASettingsApi.ReadMFASettings(ctx, d.Get("environment_id").(string)).Execute()
		},
		"ReadMFASettings",
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

	respObject := resp.(*mfa.MFASettings)

	d.Set("pairing", flattenMFASettingPairing(respObject.GetPairing()))

	if v, ok := respObject.GetLockoutOk(); ok {
		d.Set("lockout", flattenMFASettingLockout(*v))
	} else {
		d.Set("lockout", nil)
	}

	if v, ok := respObject.GetPhoneExtensionsOk(); ok {
		d.Set("phone_extensions_enabled", v.GetEnabled())
	} else {
		d.Set("phone_extensions_enabled", nil)
	}

	if v, ok := respObject.GetAuthenticationOk(); ok {
		d.Set("authentication", flattenMFASettingAuthentication(*v))
	} else {
		d.Set("authentication", nil)
	}

	return diags
}

func resourceMFASettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	mfaSettings := *mfa.NewMFASettings(expandMFASettingsPairing(d.Get("pairing").([]interface{})))

	if v, ok := d.GetOk("authentication"); ok {
		mfaSettings.SetAuthentication(expandMFASettingsAuthentication(v.([]interface{})))
	}

	if v, ok := d.GetOk("lockout"); ok {
		mfaSettings.SetLockout(expandMFASettingsLockout(v.([]interface{})))
	}

	if v, ok := d.GetOk("phone_extensions_enabled"); ok {
		mfaSettings.SetPhoneExtensions(expandMFASettingsPhoneExtensions(v))
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.MFASettingsApi.UpdateMFASettings(ctx, d.Get("environment_id").(string)).MFASettings(mfaSettings).Execute()
		},
		"UpdateMFASettings",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceMFASettingsRead(ctx, d, meta)
}

func resourceMFASettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.MFASettingsApi.ResetMFASettings(ctx, d.Get("environment_id").(string)).Execute()
		},
		"ResetMFASettings",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceMFASettingsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 1
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID\"", d.Id())
	}

	environmentID := attributes[0]

	d.SetId(environmentID)

	resourceMFASettingsRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandMFASettingsPairing(v []interface{}) mfa.MFASettingsPairing {
	obj := v[0].(map[string]interface{})

	return *mfa.NewMFASettingsPairing(int32(obj["max_allowed_devices"].(int)), mfa.EnumMFASettingsPairingKeyFormat(obj["pairing_key_format"].(string)))
}

func expandMFASettingsLockout(v []interface{}) mfa.MFASettingsLockout {
	obj := v[0].(map[string]interface{})

	mfa := *mfa.NewMFASettingsLockout(int32(obj["failure_count"].(int)))

	if v, ok := obj["duration_seconds"].(int); ok && v > 0 {
		mfa.SetDurationSeconds(int32(v))
	}

	return mfa
}

func expandMFASettingsPhoneExtensions(v interface{}) mfa.MFASettingsPhoneExtensions {
	mfa := *mfa.NewMFASettingsPhoneExtensions()
	mfa.SetEnabled(v.(bool))

	return mfa
}

func expandMFASettingsAuthentication(v []interface{}) mfa.MFASettingsAuthentication {
	obj := v[0].(map[string]interface{})

	return *mfa.NewMFASettingsAuthentication(mfa.EnumMFASettingsDeviceSelection(obj["device_selection"].(string)))
}

func flattenMFASettingAuthentication(v mfa.MFASettingsAuthentication) []map[string]interface{} {
	c := make([]map[string]interface{}, 0)
	return append(c, map[string]interface{}{
		"device_selection": string(v.GetDeviceSelection()),
	})
}

func flattenMFASettingLockout(v mfa.MFASettingsLockout) []map[string]interface{} {
	c := make([]map[string]interface{}, 0)
	return append(c, map[string]interface{}{
		"failure_count":    v.GetFailureCount(),
		"duration_seconds": v.GetDurationSeconds(),
	})
}

func flattenMFASettingPairing(v mfa.MFASettingsPairing) []map[string]interface{} {
	c := make([]map[string]interface{}, 0)
	return append(c, map[string]interface{}{
		"max_allowed_devices": v.GetMaxAllowedDevices(),
		"pairing_key_format":  string(v.GetPairingKeyFormat()),
	})
}
