package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourcePasswordPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne password policy data",

		ReadContext: datasourcePingOnePasswordPolicyRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"password_policy_id": {
				Description:      "The ID of the password policy.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the password policy.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password_policy_id"},
			},
			"description": {
				Description: "A description to apply to the password policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment_default": {
				Description: "Indicates whether this password policy is enforced within the environment. When set to true, all other password policies are set to false. Note: this may cause state management conflicts if more than one password policy is set as default.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"bypass_policy": {
				Description: "Determines whether the password policy for a user will be ignored.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"exclude_commonly_used_passwords": {
				Description: "Set this to true to ensure the password is not one of the commonly used passwords.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"exclude_profile_data": {
				Description: "Set this to true to ensure the password is not an exact match for the value of any attribute in the user’s profile, such as name, phone number, or address.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"password_history": {
				Description: "Settings to control the users password history.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prior_password_count": {
							Description: "Specifies the number of prior passwords to keep for prevention of password re-use. The value must be a positive, non-zero integer.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"retention_days": {
							Description: "The length of time to keep recent passwords for prevention of password re-use. The value must be a positive, non-zero integer.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"password_length": {
				Description: "Settings to control the user's password length.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Description: "The maximum number of characters allowed for the password. Defaults to 255. This property is not enforced when not present.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"min": {
							Description: "The minimum number of characters required for the password. This can be from 8 to 32 (inclusive). Defaults to 8 characters. This property is not enforced when not present.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"account_lockout": {
				Description: "Settings to control the user's lockout on unsuccessful authentication attempts.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration_seconds": {
							Description: "The length of time before a password is automatically moved out of the lock out state. The value must be a positive, non-zero integer.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"fail_count": {
							Description: "The number of tries before a password is placed in the lockout state. The value must be a positive, non-zero integer.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"min_characters": {
				Description: "Sets of characters that can be included, and the value is the minimum number of times one of the characters must appear in the password. The only allowed key values are `ABCDEFGHIJKLMNOPQRSTUVWXYZ`, `abcdefghijklmnopqrstuvwxyz`, `0123456789`, and `~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`. This property is not enforced when not present.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alphabetical_uppercase": {
							Description: "Count of alphabetical uppercase characters (`ABCDEFGHIJKLMNOPQRSTUVWXYZ`) that should feature in the user's password.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"alphabetical_lowercase": {
							Description: "Count of alphabetical uppercase characters (`abcdefghijklmnopqrstuvwxyz`) that should feature in the user's password.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"numeric": {
							Description: "Count of numeric characters (`0123456789`) that should feature in the user's password.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"special_characters": {
							Description: "Count of special characters (`~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`) that should feature in the user's password.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"password_age": {
				Description: "Settings to control the user's password age.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Description: "The maximum number of days the same password can be used before it must be changed. The value must be a positive, non-zero integer.  The value must be greater than the sum of `min` (if set) + 21 (the expiration warning interval for passwords).",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"min": {
							Description: "The minimum number of days a password must be used before changing. The value must be a positive, non-zero integer. This property is not enforced when not present.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"max_repeated_characters": {
				Description: "The maximum number of repeated characters allowed. This property is not enforced when not present.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"min_complexity": {
				Description: "The minimum complexity of the password based on the concept of password haystacks. The value is the number of days required to exhaust the entire search space during a brute force attack. This property is not enforced when not present.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"min_unique_characters": {
				Description: "The minimum number of unique characters required. This property is not enforced when not present.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"not_similar_to_current": {
				Description: "Set this to true to ensure that the proposed password is not too similar to the user's current password based on the Levenshtein distance algorithm. The value of this parameter is evaluated only for password change actions in which the user enters both the current and the new password. By design, PingOne does not know the user's current password.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"population_count": {
				Description: "The number of populations associated with the password policy.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOnePasswordPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp management.PasswordPolicy

	if v, ok := d.GetOk("name"); ok {

		respList, diags := sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := apiClient.PasswordPoliciesApi.ReadAllPasswordPolicies(ctx, d.Get("environment_id").(string)).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
			},
			"ReadAllPasswordPolicies",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		if passwordPolicies, ok := respList.(*management.EntityArray).Embedded.GetPasswordPoliciesOk(); ok {

			found := false
			for _, passwordPolicy := range passwordPolicies {

				if passwordPolicy.GetName() == v.(string) {
					resp = passwordPolicy
					found = true
					break
				}
			}

			if !found {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Cannot find password policy %s", v),
				})

				return diags
			}

		}

	} else if v, ok2 := d.GetOk("password_policy_id"); ok2 {

		passwordPolicyResp, diags := sdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := apiClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, d.Get("environment_id").(string), v.(string)).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
			},
			"ReadOnePasswordPolicy",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		resp = *passwordPolicyResp.(*management.PasswordPolicy)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither password_policy_id or name are set",
			Detail:   "Neither password_policy_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("password_policy_id", resp.GetId())
	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := resp.GetDefaultOk(); ok {
		d.Set("environment_default", v)
	} else {
		d.Set("environment_default", nil)
	}

	if v, ok := resp.GetBypassPolicyOk(); ok {
		d.Set("bypass_policy", v)
	} else {
		d.Set("bypass_policy", nil)
	}

	if v, ok := resp.GetExcludesCommonlyUsedOk(); ok {
		d.Set("exclude_commonly_used_passwords", v)
	} else {
		d.Set("exclude_commonly_used_passwords", nil)
	}

	if v, ok := resp.GetExcludesProfileDataOk(); ok {
		d.Set("exclude_profile_data", v)
	} else {
		d.Set("exclude_profile_data", nil)
	}

	if v, ok := resp.GetHistoryOk(); ok {
		flattenedVal := flattenPasswordHistory(v)
		d.Set("password_history", flattenedVal)
	} else {
		d.Set("password_history", nil)
	}

	if v, ok := resp.GetLengthOk(); ok {
		flattenedVal := flattenPasswordLength(v)
		d.Set("password_length", flattenedVal)
	} else {
		d.Set("password_length", nil)
	}

	if v, ok := resp.GetLockoutOk(); ok {
		flattenedVal := flattenUserLockout(v)
		d.Set("account_lockout", flattenedVal)
	} else {
		d.Set("account_lockout", nil)
	}

	if v, ok := resp.GetMinCharactersOk(); ok {
		flattenedVal := flattenMinCharacters(v)
		d.Set("min_characters", flattenedVal)
	} else {
		d.Set("min_characters", nil)
	}

	passwordAgeMaxV, passwordAgeMaxOk := resp.GetMaxAgeDaysOk()
	passwordAgeMinV, passwordAgeMinOk := resp.GetMinAgeDaysOk()

	if passwordAgeMaxOk || passwordAgeMinOk {
		flattenedVal := flattenPasswordAge(passwordAgeMaxV, passwordAgeMinV)
		d.Set("password_age", flattenedVal)
	} else {
		d.Set("password_age", nil)
	}

	if v, ok := resp.GetMaxRepeatedCharactersOk(); ok {
		d.Set("max_repeated_characters", v)
	} else {
		d.Set("max_repeated_characters", nil)
	}

	if v, ok := resp.GetMinComplexityOk(); ok {
		d.Set("min_complexity", v)
	} else {
		d.Set("min_complexity", nil)
	}

	if v, ok := resp.GetMinUniqueCharactersOk(); ok {
		d.Set("min_unique_characters", v)
	} else {
		d.Set("min_unique_characters", nil)
	}

	if v, ok := resp.GetNotSimilarToCurrentOk(); ok {
		d.Set("not_similar_to_current", v)
	} else {
		d.Set("not_similar_to_current", nil)
	}

	if v, ok := resp.GetPopulationCountOk(); ok {
		d.Set("population_count", v)
	} else {
		d.Set("population_count", nil)
	}

	return diags
}
