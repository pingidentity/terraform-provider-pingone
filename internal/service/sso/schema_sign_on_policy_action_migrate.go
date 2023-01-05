package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func resourceSignOnPolicyActionSchemaV0() *schema.Resource {
	// Ignore "XS001 schema should configure Description" lint
	//lintignore:XS001
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sign_on_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"conditions": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_sign_on_older_than_seconds": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"progressive_profiling", "agreement"},
							AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"user_is_member_of_any_population_id": {
							Type:          schema.TypeSet,
							MaxItems:      100,
							Optional:      true,
							ConflictsWith: []string{"progressive_profiling", "agreement", "identity_provider"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"user_attribute_equals": {
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"progressive_profiling", "agreement", "identity_provider"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute_reference": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"ip_out_of_range_cidr": {
							Type:          schema.TypeSet,
							MaxItems:      100,
							Optional:      true,
							ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"ip_reputation_high_risk": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
							AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"geovelocity_anomaly_detected": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
							AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"anonymous_network_detected": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
							AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
						},
						"anonymous_network_detected_allowed_cidr": {
							Type:          schema.TypeSet,
							MaxItems:      100,
							Optional:      true,
							ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"registration_external_href": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"registration_local_population_id", "agreement", "identity_provider", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			},
			"registration_local_population_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"registration_external_href", "agreement", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			},
			"registration_confirm_user_attributes": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"registration_external_href", "agreement", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			},
			"social_provider_ids": {
				Type:          schema.TypeSet,
				MaxItems:      100,
				Optional:      true,
				ConflictsWith: []string{"agreement", "mfa", "identity_provider", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enforce_lockout_for_identity_providers": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"agreement", "mfa", "identity_provider", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			},
			"agreement": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agreement_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"show_decline_option": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"identifier_first": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"discovery_rule": {
							Type:     schema.TypeSet,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute_contains_text": {
										Type:     schema.TypeString,
										Required: true,
									},
									"identity_provider_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"recovery_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"identity_provider": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acr_values": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"identity_provider_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"pass_user_context": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"login": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recovery_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"mfa": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_sign_on_policy_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"no_device_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(management.ENUMSIGNONPOLICYNODEVICEMODE_BLOCK),
						},
					},
				},
			},
			"progressive_profiling": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"required": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"prevent_multiple_prompts_per_flow": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"prompt_interval_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  7776000,
						},
						"prompt_text": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"pingid": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"pingid_windows_login_passwordless": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_user_attribute_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"offline_mode_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ResourceSignOnPolicyActionStateUpgradeV0(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return nil, nil
	}

	if conditions, ok := rawState["conditions"].([]map[string]interface{}); ok && len(conditions) > 0 && conditions[0] != nil {

		if v, ok := conditions[0]["user_attribute_equals"].([]map[string]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {

			newUserAttributeEqualsMap := make([]map[string]interface{}, 0)

			for _, attributeMap := range v {

				newUserAttributeEqualsMap = append(newUserAttributeEqualsMap, map[string]interface{}{
					"attribute_reference": attributeMap["attribute_reference"].(string),
					"value_string":        attributeMap["value"].(string),
				})
			}

			rawState["conditions"].([]map[string]interface{})[0]["user_attribute_equals"] = newUserAttributeEqualsMap

		}

	}

	return rawState, nil
}
