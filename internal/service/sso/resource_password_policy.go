package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourcePasswordPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne password policies",

		CreateContext: resourcePasswordPolicyCreate,
		ReadContext:   resourcePasswordPolicyRead,
		UpdateContext: resourcePasswordPolicyUpdate,
		DeleteContext: resourcePasswordPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePasswordPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the password policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the password policy.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the password policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"environment_default": {
				Description: "Indicates whether this password policy is enforced within the environment. When set to true, all other password policies are set to false. Note: this may cause state management conflicts if more than one password policy is set as default.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"bypass_policy": {
				Description: "Determines whether the password policy for a user will be ignored.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"exclude_commonly_used_passwords": {
				Description: "Set this to true to ensure the password is not one of the commonly used passwords.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"exclude_profile_data": {
				Description: "Set this to true to ensure the password is not an exact match for the value of any attribute in the user’s profile, such as name, phone number, or address.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"password_history": {
				Description: "Settings to control the users password history.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prior_password_count": {
							Description:      "Specifies the number of prior passwords to keep for prevention of password re-use. The value must be a positive, non-zero integer.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
						"retention_days": {
							Description:      "The length of time to keep recent passwords for prevention of password re-use. The value must be a positive, non-zero integer.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
					},
				},
			},
			"password_length": {
				Description: "Settings to control the user's password length.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Description:      "The maximum number of characters allowed for the password. Defaults to 255. This property is not enforced when not present.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(255, 255)),
						},
						"min": {
							Description:      "The minimum number of characters required for the password. Defaults to 8 characters. This property is not enforced when not present.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(8, 8)),
						},
					},
				},
			},
			"account_lockout": {
				Description: "Settings to control the user's lockout on unsuccessful authentication attempts.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration_seconds": {
							Description:      "The length of time before a password is automatically moved out of the lock out state. The value must be a positive, non-zero integer.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
						"fail_count": {
							Description:      "The number of tries before a password is placed in the lockout state. The value must be a positive, non-zero integer.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
					},
				},
			},
			"min_characters": {
				Description: "Sets of characters that can be included, and the value is the minimum number of times one of the characters must appear in the password. The only allowed key values are `ABCDEFGHIJKLMNOPQRSTUVWXYZ`, `abcdefghijklmnopqrstuvwxyz`, `0123456789`, and `~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`. This property is not enforced when not present.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alphabetical_uppercase": {
							Description:      "Count of alphabetical uppercase characters (`ABCDEFGHIJKLMNOPQRSTUVWXYZ`) that should feature in the user's password.  Fixed value of 1.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 1)),
							Optional:         true,
						},
						"alphabetical_lowercase": {
							Description:      "Count of alphabetical uppercase characters (`abcdefghijklmnopqrstuvwxyz`) that should feature in the user's password.  Fixed value of 1.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 1)),
							Optional:         true,
						},
						"numeric": {
							Description:      "Count of numeric characters (`0123456789`) that should feature in the user's password.  Fixed value of 1.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 1)),
							Optional:         true,
						},
						"special_characters": {
							Description:      "Count of special characters (`~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`) that should feature in the user's password.  Fixed value of 1.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 1)),
							Optional:         true,
						},
					},
				},
			},
			"password_age": {
				Description: "Settings to control the user's password age.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Description:      "The maximum number of days the same password can be used before it must be changed. The value must be a positive, non-zero integer.  The value must be greater than the sum of minAgeDays (if set) + 21 (the expiration warning interval for passwords).",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
						"min": {
							Description:      "The minimum number of days a password must be used before changing. The value must be a positive, non-zero integer. This property is not enforced when not present.",
							Type:             schema.TypeInt,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Optional:         true,
						},
					},
				},
			},
			"max_repeated_characters": {
				Description:      "The maximum number of repeated characters allowed. This property is not enforced when not present.",
				Type:             schema.TypeInt,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(2, 2)),
				Optional:         true,
			},
			"min_complexity": {
				Description:      "The minimum complexity of the password based on the concept of password haystacks. The value is the number of days required to exhaust the entire search space during a brute force attack. This property is not enforced when not present.",
				Type:             schema.TypeInt,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(7, 7)),
				Optional:         true,
			},
			"min_unique_characters": {
				Description:      "The minimum number of unique characters required. This property is not enforced when not present.",
				Type:             schema.TypeInt,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 5)),
				Optional:         true,
			},
			"not_similar_to_current": {
				Description: "Set this to true to ensure that the proposed password is not too similar to the user's current password based on the Levenshtein distance algorithm. The value of this parameter is evaluated only for password change actions in which the user enters both the current and the new password. By design, PingOne does not know the user's current password.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"population_count": {
				Description: "The number of populations associated with the password policy.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourcePasswordPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	passwordPolicy := expandPasswordPolicy(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PasswordPoliciesApi.CreatePasswordPolicy(ctx, d.Get("environment_id").(string)).PasswordPolicy(passwordPolicy.(management.PasswordPolicy)).Execute()
		},
		"CreatePasswordPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.PasswordPolicy)

	d.SetId(respObject.GetId())

	return resourcePasswordPolicyRead(ctx, d, meta)
}

func resourcePasswordPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOnePasswordPolicy",
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

	respObject := resp.(*management.PasswordPolicy)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := respObject.GetDefaultOk(); ok {
		d.Set("environment_default", v)
	} else {
		d.Set("environment_default", nil)
	}

	if v, ok := respObject.GetBypassPolicyOk(); ok {
		d.Set("bypass_policy", v)
	} else {
		d.Set("bypass_policy", nil)
	}

	if v, ok := respObject.GetExcludesCommonlyUsedOk(); ok {
		d.Set("exclude_commonly_used_passwords", v)
	} else {
		d.Set("exclude_commonly_used_passwords", nil)
	}

	if v, ok := respObject.GetExcludesProfileDataOk(); ok {
		d.Set("exclude_profile_data", v)
	} else {
		d.Set("exclude_profile_data", nil)
	}

	if v, ok := respObject.GetHistoryOk(); ok {
		flattenedVal := flattenPasswordHistory(v)
		d.Set("password_history", flattenedVal)
	} else {
		d.Set("password_history", nil)
	}

	if v, ok := respObject.GetLengthOk(); ok {
		flattenedVal := flattenPasswordLength(v)
		d.Set("password_length", flattenedVal)
	} else {
		d.Set("password_length", nil)
	}

	if v, ok := respObject.GetLockoutOk(); ok {
		flattenedVal := flattenUserLockout(v)
		d.Set("account_lockout", flattenedVal)
	} else {
		d.Set("account_lockout", nil)
	}

	if v, ok := respObject.GetMinCharactersOk(); ok {
		flattenedVal := flattenMinCharacters(v)
		d.Set("min_characters", flattenedVal)
	} else {
		d.Set("min_characters", nil)
	}

	passwordAgeMaxV, passwordAgeMaxOk := respObject.GetMaxAgeDaysOk()
	passwordAgeMinV, passwordAgeMinOk := respObject.GetMinAgeDaysOk()

	if passwordAgeMaxOk || passwordAgeMinOk {
		flattenedVal := flattenPasswordAge(passwordAgeMaxV, passwordAgeMinV)
		d.Set("password_age", flattenedVal)
	} else {
		d.Set("password_age", nil)
	}

	if v, ok := respObject.GetMaxRepeatedCharactersOk(); ok {
		d.Set("max_repeated_characters", v)
	} else {
		d.Set("max_repeated_characters", nil)
	}

	if v, ok := respObject.GetMinComplexityOk(); ok {
		d.Set("min_complexity", v)
	} else {
		d.Set("min_complexity", nil)
	}

	if v, ok := respObject.GetMinUniqueCharactersOk(); ok {
		d.Set("min_unique_characters", v)
	} else {
		d.Set("min_unique_characters", nil)
	}

	if v, ok := respObject.GetNotSimilarToCurrentOk(); ok {
		d.Set("not_similar_to_current", v)
	} else {
		d.Set("not_similar_to_current", nil)
	}

	if v, ok := respObject.GetPopulationCountOk(); ok {
		d.Set("population_count", v)
	} else {
		d.Set("population_count", nil)
	}

	return diags
}

func resourcePasswordPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	passwordPolicy := expandPasswordPolicy(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PasswordPoliciesApi.UpdatePasswordPolicy(ctx, d.Get("environment_id").(string), d.Id()).PasswordPolicy(passwordPolicy.(management.PasswordPolicy)).Execute()
		},
		"UpdatePasswordPolicy",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourcePasswordPolicyRead(ctx, d, meta)
}

func resourcePasswordPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.PasswordPoliciesApi.DeletePasswordPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeletePasswordPolicy",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePasswordPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/passwordPolicyID\"", d.Id())
	}

	environmentID, passwordPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(passwordPolicyID)

	resourcePasswordPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandPasswordPolicy(d *schema.ResourceData) interface{} {

	passwordPolicy := *management.NewPasswordPolicy(d.Get("exclude_commonly_used_passwords").(bool), d.Get("exclude_profile_data").(bool), d.Get("name").(string), d.Get("not_similar_to_current").(bool)) // PasswordPolicy |  (optional)

	if v, ok := d.GetOk("description"); ok {
		passwordPolicy.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("environment_default"); ok {
		passwordPolicy.SetDefault(v.(bool))
	}

	if v, ok := d.GetOk("bypass_policy"); ok {
		passwordPolicy.SetBypassPolicy(v.(bool))
	}

	if v, ok := d.GetOk("password_history"); ok {

		priorCount := v.([]interface{})[0].(map[string]interface{})["prior_password_count"]
		retentionDays := v.([]interface{})[0].(map[string]interface{})["retention_days"]

		if priorCount != nil || retentionDays != nil {

			passwordPolicyHistory := *management.NewPasswordPolicyHistory()

			if priorCount != nil {
				passwordPolicyHistory.SetCount(int32(priorCount.(int)))
			}

			if retentionDays != nil {
				passwordPolicyHistory.SetRetentionDays(int32(retentionDays.(int)))
			}

			passwordPolicy.SetHistory(passwordPolicyHistory)
		}

	}

	if v, ok := d.GetOk("password_length"); ok {

		max := v.([]interface{})[0].(map[string]interface{})["max"]
		min := v.([]interface{})[0].(map[string]interface{})["min"]

		if max != nil || min != nil {

			passwordPolicyLength := *management.NewPasswordPolicyLength()

			if max != nil {
				passwordPolicyLength.SetMax(int32(max.(int)))
			}

			if min != nil {
				passwordPolicyLength.SetMin(int32(min.(int)))
			}

			passwordPolicy.SetLength(passwordPolicyLength)
		}

	}

	if v, ok := d.GetOk("account_lockout"); ok {

		duration := v.([]interface{})[0].(map[string]interface{})["duration_seconds"]
		failCount := v.([]interface{})[0].(map[string]interface{})["fail_count"]

		if duration != nil || failCount != nil {

			passwordPolicyLockout := *management.NewPasswordPolicyLockout()

			if duration != nil {
				passwordPolicyLockout.SetDurationSeconds(int32(duration.(int)))
			}

			if failCount != nil {
				passwordPolicyLockout.SetFailureCount(int32(failCount.(int)))
			}

			passwordPolicy.SetLockout(passwordPolicyLockout)
		}

	}

	if v, ok := d.GetOk("min_characters"); ok {

		alphaUpper := v.([]interface{})[0].(map[string]interface{})["alphabetical_uppercase"]
		alphaLower := v.([]interface{})[0].(map[string]interface{})["alphabetical_lowercase"]
		numeric := v.([]interface{})[0].(map[string]interface{})["numeric"]
		special := v.([]interface{})[0].(map[string]interface{})["special_characters"]

		if alphaUpper != nil || alphaLower != nil {

			passwordPolicyMinChars := *management.NewPasswordPolicyMinCharacters()

			if alphaUpper != nil {
				passwordPolicyMinChars.SetABCDEFGHIJKLMNOPQRSTUVWXYZ(int32(alphaUpper.(int)))
			}

			if alphaLower != nil {
				passwordPolicyMinChars.SetAbcdefghijklmnopqrstuvwxyz(int32(alphaLower.(int)))
			}

			if numeric != nil {
				passwordPolicyMinChars.SetVar0123456789(int32(numeric.(int)))
			}

			if special != nil {
				passwordPolicyMinChars.SetSpecialChar(int32(special.(int)))
			}

			passwordPolicy.SetMinCharacters(passwordPolicyMinChars)
		}

	}

	if v, ok := d.GetOk("password_age"); ok {

		max := v.([]interface{})[0].(map[string]interface{})["max"]
		min := v.([]interface{})[0].(map[string]interface{})["min"]

		if max != nil || min != nil {

			if max != nil {
				passwordPolicy.SetMaxAgeDays(int32(max.(int)))
			}

			if min != nil {
				passwordPolicy.SetMinAgeDays(int32(min.(int)))
			}

		}

	}

	if v, ok := d.GetOk("max_repeated_characters"); ok {
		passwordPolicy.SetMaxRepeatedCharacters(int32(v.(int)))
	}

	if v, ok := d.GetOk("min_complexity"); ok {
		passwordPolicy.SetMinComplexity(int32(v.(int)))
	}

	if v, ok := d.GetOk("min_unique_characters"); ok {
		passwordPolicy.SetMinUniqueCharacters(int32(v.(int)))
	}

	return passwordPolicy
}

func flattenPasswordHistory(passwordPolicyHistory *management.PasswordPolicyHistory) []interface{} {

	item := make(map[string]interface{})

	if v, ok := passwordPolicyHistory.GetCountOk(); ok {
		item["prior_password_count"] = v
	}

	if v, ok := passwordPolicyHistory.GetRetentionDaysOk(); ok {
		item["retention_days"] = v
	}

	items := make([]interface{}, 0)
	items = append(items, item)

	return items
}

func flattenPasswordLength(passwordPolicyLength *management.PasswordPolicyLength) []interface{} {

	item := make(map[string]interface{})

	if v, ok := passwordPolicyLength.GetMaxOk(); ok {
		item["max"] = v
	}

	if v, ok := passwordPolicyLength.GetMinOk(); ok {
		item["min"] = v
	}

	items := make([]interface{}, 0)
	items = append(items, item)

	return items

}

func flattenUserLockout(passwordPolicyLockout *management.PasswordPolicyLockout) []interface{} {

	item := make(map[string]interface{})

	if v, ok := passwordPolicyLockout.GetDurationSecondsOk(); ok {
		item["duration_seconds"] = v
	}

	if v, ok := passwordPolicyLockout.GetFailureCountOk(); ok {
		item["fail_count"] = v
	}

	items := make([]interface{}, 0)
	items = append(items, item)

	return items

}

func flattenMinCharacters(passwordPolicyMinChars *management.PasswordPolicyMinCharacters) []interface{} {

	item := make(map[string]interface{})

	if v, ok := passwordPolicyMinChars.GetABCDEFGHIJKLMNOPQRSTUVWXYZOk(); ok {
		item["alphabetical_uppercase"] = v
	}

	if v, ok := passwordPolicyMinChars.GetAbcdefghijklmnopqrstuvwxyzOk(); ok {
		item["alphabetical_lowercase"] = v
	}

	if v, ok := passwordPolicyMinChars.GetVar0123456789Ok(); ok {
		item["numeric"] = v
	}

	if v, ok := passwordPolicyMinChars.GetSpecialCharOk(); ok {
		item["special_characters"] = v
	}

	items := make([]interface{}, 0)
	items = append(items, item)

	return items

}

func flattenPasswordAge(max, min *int32) []interface{} {

	item := make(map[string]interface{})

	if max != nil {
		item["max"] = max
	}

	if min != nil {
		item["min"] = min
	}

	items := make([]interface{}, 0)
	items = append(items, item)

	return items

}
