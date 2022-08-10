package sso

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func resourceAuthenticationPolicyActionSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"action_type": {
			Description:      fmt.Sprintf("A string that specifies the type of policy action.  Options are `%s`, `%s`, `%s`, `%s`, `%s` and `%s`", string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT), string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST), string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER), string(management.ENUMSIGNONPOLICYTYPE_LOGIN), string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION), string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING)),
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT), string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST), string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER), string(management.ENUMSIGNONPOLICYTYPE_LOGIN), string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION), string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING)}, false)),
		},
		"conditions": {
			Description: "Conditions to apply to the authentication policy action.",
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_sign_on_older_than": {
						Description:      "",
						Type:             schema.TypeInt,
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
					},
					"user_is_member_of_any_population_id": {
						Description: "",
						Type:        schema.TypeList,
						MaxItems:    100,
						Optional:    true,
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
					"user_attribute_equals": {
						Description: "",
						Type:        schema.TypeSet,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attribute_reference": {
									Description:      "",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"value": {
									Description:      "",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
							},
						},
					},
					"ip_out_of_range_cidr": {
						Description: "",
						Type:        schema.TypeList,
						MaxItems:    100,
						Optional:    true,
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
					},
					"ip_reputation_high_risk": {
						Description: "",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"geovelocity_anomaly_detected": {
						Description: "",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"anonymous_network_detected": {
						Description: "",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"anonymous_network_detected_allowed_cidr": {
						Description: "",
						Type:        schema.TypeList,
						MaxItems:    100,
						Optional:    true,
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
					},
				},
			},
		},
		"identifier_first_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"confirm_identity_provider_attributes": confirmIdentityProviderAttributesSchema(),
					"discovery_rule": {
						Description: "An IDP discovery rule invoked when no user is associated with the user identifier. The condition on which this identity provider is used to authenticate the user is expressed using the PingOne policy condition language.",
						Type:        schema.TypeList,
						MaxItems:    100,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"condition": {
									Description: "Set the authentication policy action to be the `Identifier First` type.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Required:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"contains": {
												Description:      "",
												Type:             schema.TypeString,
												Required:         true,
												ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
											},
											"value": {
												Description:      "",
												Type:             schema.TypeString,
												Required:         true,
												ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
											},
										},
									},
								},
								"identity_provider_id": {
									Description:      "The ID that specifies the identity provider that will be used to authenticate the user if the condition is matched.",
									Type:             schema.TypeString,
									Required:         true,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
							},
						},
					},
					"enforce_lockout_for_identity_providers": enforceLockoutForIdPSchema(),
					"recovery":                               recoverySchema(),
					"registration":                           registrationSchema(),
					"social_providers":                       socialProvidersSchema(),
				},
			},
		},
		"login_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_LOGIN)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"confirm_identity_provider_attributes":   confirmIdentityProviderAttributesSchema(),
					"enforce_lockout_for_identity_providers": enforceLockoutForIdPSchema(),
					"recovery":                               recoverySchema(),
					"registration":                           registrationSchema(),
					"social_providers":                       socialProvidersSchema(),
				},
			},
		},
		"mfa_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_authentication_policy_id": {
						Description:      "The ID of the MFA policy that should be used.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"no_device_mode": {
						Description:      "A string that specifies the device mode for the MFA flow. Options are `BYPASS` to allow MFA without a specified device, or `BLOCK` to block the MFA flow if no device is specified. To use this configuration option, the authorize request must include a signed `login_hint_token` property. For more information, see Authorize (Browserless and MFA Only Flows).",
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BLOCK", "BYPASS"}, false)),
					},
				},
			},
		},
		"identity_provider_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
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
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"pass_user_context": {
						Description: "A boolean that specifies whether to pass in a login hint to the identity provider on the authentication request. Based on user context, the login hint is set if (1) the user is set on the flow, and (2) the user already has an account link for the identity provider. If both of these conditions are true, then the user is sent to the identity provider with a login hint equal to their externalId for the identity provider (saved on the account link). If these conditions are not true, then the API checks see if there is an OIDC login hint on the flow. If so, that login hint is used. If none of these conditions are true, the login hint parameter is not included on the authorization request to the identity provider.",
						Type:        schema.TypeBool,
						Optional:    true,
					},
					"registration": registrationSchema(),
				},
			},
		},
		"agreement_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"agreement_id": {
						Description:      "A string that specifies the ID of the agreement to which the user must consent.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"show_decline_option": {
						Description:      "When enabled, the `Do Not Accept` button will terminate the Flow and display an error message to the user.",
						Type:             schema.TypeBool,
						Optional:         true,
						Default:          true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
				},
			},
		},
		"progressive_profiling_options": {
			Description: fmt.Sprintf("Options specific to the **%s** policy action.", string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING)),
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"attribute": {
						Description: "",
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
	}
}

func confirmIdentityProviderAttributesSchema() *schema.Schema {
	return &schema.Schema{
		Description: "A boolean that specifies whether users must confirm data returned from an identity provider prior to registration. Users can modify the data and omit non-required attributes. Modified attributes are added to the user's profile during account creation.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	}
}

func enforceLockoutForIdPSchema() *schema.Schema {
	return &schema.Schema{
		Description: "A boolean that if set to true and if the user's account is locked (the account.canAuthenticate attribute is set to false), then social sign on with an external identity provider is prevented.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	}
}

func recoverySchema() *schema.Schema {
	return &schema.Schema{
		Description: "A boolean that specifies the enabled/disabled state of the account recovery action. The default is disabled when creating a new policy. When enabled, it allows the use of the forgot password flow.",
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Description: "A boolean that specifies the enabled/disabled state of the policy action. The default is disabled when creating a new policy. When enabled, it allows the use of the new user registration flow. This attribute should be set to true when implementing the social login sign-on policy option.",
					Type:        schema.TypeBool,
					Required:    true,
				},
			},
		},
	}
}

func registrationSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Specifies the account registration options.",
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Description: "A boolean that specifies the enabled/disabled state of the policy action. The default is disabled when creating a new policy. When enabled, it allows the use of the new user registration flow. This attribute should be set to true when implementing the social login sign-on policy option.",
					Type:        schema.TypeBool,
					Required:    true,
				},
				"external_href": {
					Description:      "A string that specifies the link to the external identity provider's identity store. This property is set when the administrator chooses to have users register in an external identity store. This attribute can be set only when the registration.enabled property is set to false.",
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
				"population_id": {
					Description:      "A string that specifies the population ID associated with the newly registered user.",
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
			},
		},
	}
}

func socialProvidersSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The IDs of the identity providers that can be used for the social login sign-on flow.",
		Type:        schema.TypeList,
		MaxItems:    100,
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}
