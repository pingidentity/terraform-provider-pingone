package sso

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func resourceSignOnPolicyActionSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"environment_id": {
			Description:      "The ID of the environment to create the sign on policy action in.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
		},
		"sign_on_policy_id": {
			Description:      "The ID of the sign on policy to associate the sign on policy action to.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
		},
		"priority": {
			Description:      "An integer that specifies the order in which the policy referenced by this assignment is evaluated during an authentication flow relative to other policies. An assignment with a lower priority will be evaluated first.",
			Type:             schema.TypeInt,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		},
		"conditions": {
			Description:   "Conditions to apply to the sign on policy action.  Applies to policy actions of type `agreement`, `identifier_first`, `identity_provider`, `login`, `mfa`, `progressive_profiling`.",
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_sign_on_older_than_seconds": {
						Description:      "Set the number of seconds by which the user will not be prompted for this action following the last successful authentication.  Applies to policy actions of type `identifier_first`, `identity_provider`, `login`, `mfa`.",
						Type:             schema.TypeInt,
						Optional:         true,
						ConflictsWith:    []string{"progressive_profiling", "agreement", "conditions.0.last_sign_on_older_than_seconds_mfa"},
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						AtLeastOneOf:     []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"last_sign_on_older_than_seconds_mfa": {
						Description:      "Set the number of seconds by which the user will not be prompted for this action following the last successful authentication of an MFA authenticator device.  Applies to policy actions of type `mfa`.",
						Type:             schema.TypeInt,
						Optional:         true,
						ConflictsWith:    []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider", "conditions.0.last_sign_on_older_than_seconds"},
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						AtLeastOneOf:     []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"user_is_member_of_any_population_id": {
						Description:   "Activate this action only for users within the specified list of population IDs.  Applies to policy actions of type `identifier_first`, `login`, `mfa`, but cannot be set on policy actions where the priority is `1`.",
						Type:          schema.TypeSet,
						MaxItems:      100,
						Optional:      true,
						ConflictsWith: []string{"progressive_profiling", "agreement", "identity_provider"},
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"user_attribute_equals": {
						Description:   "One or more conditions where an attribute on the user's profile must match the configured value.  Applies to policy actions of type `identifier_first`, `login`, `mfa`, but cannot be set on policy actions where the priority is `1`.",
						Type:          schema.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"progressive_profiling", "agreement", "identity_provider"},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attribute_reference": {
									Description:      "Specifies the user attribute used in the condition. Only string core, standard, and custom attributes are supported. For complex attribute types, you must reference the sub-attribute (`$${user.name.firstName}`).  Note values that begin with a dollar sign (`$`) must be prefixed with an additional dollar sign.  E.g. `${name.given}` should be configured as `$${name.given}`.  When configured, one of `value` (for attributes of type `STRING` or `INTEGER`) or `value_boolean` (for attributes of type `BOOLEAN`) must be provided.",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"value": {
									Description:      "The string or integer (as string) value of the attribute (declared in `attribute_reference`) on the user profile that should be matched.  This value parameter should be used where the data type of the schema attribute in `attribute_reference` is of type `STRING` or `INTEGER`.  Conflicts with `value_boolean`.",
									Type:             schema.TypeString,
									Optional:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"value_boolean": {
									Description: "The boolean value of the attribute (declared in `attribute_reference`) on the user profile that should be matched.  This value parameter should be used where the data type of the schema attribute in `attribute_reference` is of type `BOOLEAN` (e.g `$${user.emailVerified}`, `$${user.verified}` and `$${user.mfaEnabled}`).  Conflicts with `value`.",
									Type:        schema.TypeBool,
									Optional:    true,
								},
							},
						},
						AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"ip_out_of_range_cidr": {
						Description:   "A list of strings that specifies the supported network IP addresses expressed as classless inter-domain routing (CIDR) strings.  Applies to policy actions of type `mfa`.",
						Type:          schema.TypeSet,
						MaxItems:      100,
						Optional:      true,
						ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
						AtLeastOneOf: []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"ip_reputation_high_risk": {
						Description:   "A boolean that specifies whether the user's IP risk should be used when evaluating this policy action.  A value of `HIGH` will prompt the user to authenticate with this action.  Applies to policy actions of type `mfa`.",
						Type:          schema.TypeBool,
						Optional:      true,
						Default:       false,
						ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
						AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"geovelocity_anomaly_detected": {
						Description:   "A boolean that specifies whether the user should be prompted for re-authentication on this action based on a detected geovelocity anomaly.  Applies to policy actions of type `mfa`.",
						Type:          schema.TypeBool,
						Optional:      true,
						Default:       false,
						ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
						AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"anonymous_network_detected": {
						Description:   "A boolean that specifies whether the user should be prompted for re-authentication on this action based on a detected anonymous network.  Applies to policy actions of type `mfa`.",
						Type:          schema.TypeBool,
						Optional:      true,
						Default:       false,
						ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
						AtLeastOneOf:  []string{"conditions.0.last_sign_on_older_than_seconds", "conditions.0.last_sign_on_older_than_seconds_mfa", "conditions.0.user_is_member_of_any_population_id", "conditions.0.user_attribute_equals", "conditions.0.ip_out_of_range_cidr", "conditions.0.ip_reputation_high_risk", "conditions.0.geovelocity_anomaly_detected", "conditions.0.anonymous_network_detected"},
					},
					"anonymous_network_detected_allowed_cidr": {
						Description:   "A list of allowed CIDR when an anonymous network is detected.  Applies to policy actions of type `mfa`.",
						Type:          schema.TypeSet,
						MaxItems:      100,
						Optional:      true,
						ConflictsWith: []string{"identifier_first", "login", "progressive_profiling", "agreement", "identity_provider"},
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
					},
				},
			},
		},
		"registration_external_href": {
			Description:      "A string that specifies the link to the external identity provider's identity store. This property is set when the administrator chooses to have users register in an external identity store. This attribute can be set only when the registration.enabled property is set to false.",
			Type:             schema.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"registration_local_population_id", "agreement", "identity_provider", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"registration_local_population_id": {
			Description:      "A string that specifies the population ID associated with the newly registered user. Setting this enables local registration features.",
			Type:             schema.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"registration_external_href", "agreement", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
		},
		"registration_confirm_user_attributes": {
			Description:   "A boolean that specifies whether users must confirm data returned from an identity provider prior to registration. Users can modify the data and omit non-required attributes. Modified attributes are added to the user's profile during account creation.",
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"registration_external_href", "agreement", "mfa", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
		},
		"social_provider_ids": {
			Description:   "One or more IDs of the identity providers that can be used for the social login sign-on flow.",
			Type:          schema.TypeSet,
			MaxItems:      100,
			Optional:      true,
			ConflictsWith: []string{"agreement", "mfa", "identity_provider", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
		},
		"enforce_lockout_for_identity_providers": {
			Description:   "A boolean that if set to true and if the user's account is locked (the account.canAuthenticate attribute is set to false), then social sign on with an external identity provider is prevented.",
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"agreement", "mfa", "identity_provider", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
		},
		"agreement": {
			Description:  "Options specific to the **Agreements** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"agreement_id": {
						Description:      "A string that specifies the ID of the agreement to which the user must consent.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
					},
					"show_decline_option": {
						Description: "When enabled, the `Do Not Accept` button will terminate the Flow and display an error message to the user.",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},
		"identifier_first": {
			Description:  "Options specific to the **Identifier First** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"discovery_rule": {
						Description: "One or more IDP discovery rules invoked when no user is associated with the user identifier. The condition on which this identity provider is used to authenticate the user is expressed using the PingOne policy condition language.",
						Type:        schema.TypeSet,
						MaxItems:    100,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attribute_contains_text": {
									Description:      "Text to match on a user's username. Any users that don't match a discovery rule will authenticate against PingOne.  E.g `@pingidentity.com`",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"identity_provider_id": {
									Description:      "The ID that specifies the identity provider that will be used to authenticate the user if the condition is matched.",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
								},
							},
						},
					},
					"recovery_enabled": recoveryEnabledSchema(),
				},
			},
		},
		"identity_provider": {
			Description:  "Options specific to the **Identity Provider** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acr_values": {
						Description:      "A string that designates the sign-on policies included in the authorization flow request. Options can include the PingOne predefined sign-on policies, Single_Factor and Multi_Factor, or any custom defined sign-on policy names. Sign-on policy names should be listed in order of preference, and they must be assigned to the application. This property can be configured on the identity provider action and is passed to the identity provider if the identity provider is of type `SAML` or `OPENID_CONNECT`.",
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"identity_provider_id": {
						Description:      "A string that specifies the ID of the external identity provider to which the user is redirected for sign-on.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
					},
					"pass_user_context": {
						Description: "A boolean that specifies whether to pass in a login hint to the identity provider on the sign on request. Based on user context, the login hint is set if (1) the user is set on the flow, and (2) the user already has an account link for the identity provider. If both of these conditions are true, then the user is sent to the identity provider with a login hint equal to their externalId for the identity provider (saved on the account link). If these conditions are not true, then the API checks see if there is an OIDC login hint on the flow. If so, that login hint is used. If none of these conditions are true, the login hint parameter is not included on the authorization request to the identity provider.",
						Type:        schema.TypeBool,
						Optional:    true,
					},
				},
			},
		},
		"login": {
			Description:  "Options specific to the **Login** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"recovery_enabled": recoveryEnabledSchema(),
				},
			},
		},
		"mfa": {
			Description:  "Options specific to the **Multi-factor Authentication** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_sign_on_policy_id": {
						Description:      "The ID of the MFA policy that should be used.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
					},
					"no_device_mode": {
						Description:      "A string that specifies the device mode for the MFA flow. Options are `BYPASS` to allow MFA without a specified device, or `BLOCK` to block the MFA flow if no device is specified. To use this configuration option, the authorize request must include a signed `login_hint_token` property. For more information, see Authorize (Browserless and MFA Only Flows).",
						Type:             schema.TypeString,
						Optional:         true,
						Default:          string(management.ENUMSIGNONPOLICYNODEVICEMODE_BLOCK),
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMSIGNONPOLICYNODEVICEMODE_BLOCK), string(management.ENUMSIGNONPOLICYNODEVICEMODE_BYPASS)}, false)),
					},
				},
			},
		},
		"progressive_profiling": {
			Description:  "Options specific to the **Progressive Profiling** policy action.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"attribute": {
						Description: "One or more attribute(s) that the user should be prompted to complete as part of the progressive profiling action.",
						Type:        schema.TypeSet,
						Required:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:      "A string that specifies the name and path of the user profile attribute as defined in the user schema (for example, email or address.postalCode).",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"required": {
									Description: "A boolean that specifies whether the user is required to provide a value for the attribute.",
									Type:        schema.TypeBool,
									Required:    true,
								},
							},
						},
					},
					"prevent_multiple_prompts_per_flow": {
						Description: "A boolean that specifies whether the progressive profiling action will not be executed if another progressive profiling action has already been executed during the flow.",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"prompt_interval_seconds": {
						Description:      "An integer that specifies how often to prompt the user to provide profile data for the configured attributes for which they do not have values.",
						Type:             schema.TypeInt,
						Optional:         true,
						Default:          7776000,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
					},
					"prompt_text": {
						Description:      "A string that specifies text to display to the user when prompting for attribute values.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
				},
			},
		},
		"pingid": {
			Description:  "Options specific to the **PingID** policy action.  This action can only be applied to Workforce solution context environments that have the PingID and SSO services enabled.",
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
			Description:  "Options specific to the **PingID Windows Login Passwordless** policy action.  This action can only be applied to Workforce solution context environments that have the PingID and SSO services enabled.",
			Type:         schema.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"identifier_first", "login", "mfa", "identity_provider", "agreement", "progressive_profiling", "pingid", "pingid_windows_login_passwordless"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"unique_user_attribute_name": {
						Description:      "A string that specifies the schema attribute to match against the provided identifier when searching for a user in the directory. Only unique attributes in the directory schema may be configured.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"offline_mode_enabled": {
						Description: "A boolean that specifies whether to allow users to log in when PingOne and or PingID are not available.",
						Type:        schema.TypeBool,
						Required:    true,
					},
				},
			},
		},
	}
}

func recoveryEnabledSchema() *schema.Schema {
	return &schema.Schema{
		Description: "A boolean that specifies whether account recovery features are active on the policy action.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	}
}
