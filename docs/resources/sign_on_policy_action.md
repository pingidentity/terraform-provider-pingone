---
page_title: "pingone_sign_on_policy_action Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to create and manage PingOne sign on policy actions.
---

# pingone_sign_on_policy_action (Resource)

Resource to create and manage PingOne sign on policy actions.

~> A warning will be issued following `terraform apply` when attempting to remove the final sign-on policy action from an associated sign-on policy.  When removing the final sign-on policy action from a sign-on policy, it's recommended to also remove the associated sign-on policy at the same time.  Further information can be found [here](https://github.com/pingidentity/terraform-provider-pingone/issues/68).

~> Some policy action conditions, such as `conditions.user_attribute_equals` and `conditions.user_is_member_of_any_population_id` conditions, are not available where the `priority` of a policy action is `1`.  Please refer to the schema documentation for more information.

## Example Usage - First Factor (Username/Password)

```terraform
resource "pingone_sign_on_policy_action" "my_policy_first_factor" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  login {
    recovery_enabled = true
  }
}
```

## Example Usage - First Factor (Username/Password) with New User Provisioning Gateway

```terraform
resource "pingone_gateway" "my_awesome_ldap_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome LDAP Gateway"

  # ...

  user_types = {
    "User Set 1" = {
      # ...
    }
  }
}

resource "pingone_gateway_credential" "my_awesome_ldap_gateway" {
  environment_id = pingone_environment.my_environment.id
  gateway_id     = pingone_gateway.my_awesome_ldap_gateway.id
}

resource "pingone_sign_on_policy_action" "my_policy_first_factor" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  login {
    recovery_enabled = true

    new_user_provisioning {
      gateway {
        id           = pingone_gateway.my_awesome_ldap_gateway.id
        user_type_id = pingone_gateway.my_awesome_ldap_gateway.user_types["User Set 1"].id
      }
    }
  }

  depends_on = [
    pingone_gateway_credential.my_awesome_ldap_gateway
  ]
}
```

## Example Usage - Identifier First

```terraform
resource "pingone_sign_on_policy_action" "my_policy_identifier_first" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  identifier_first {
    recovery_enabled = true

    discovery_rule {
      attribute_contains_text = "@pingidentity.com"
      identity_provider_id    = pingone_identity_provider.my_identity_provider.id
    }
  }
}
```

## Example Usage - MFA

```terraform
resource "pingone_sign_on_policy_action" "my_policy_mfa" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 2

  conditions {
    last_sign_on_older_than_seconds_mfa = 86400 // 24 hours

    ip_reputation_high_risk      = true
    geovelocity_anomaly_detected = true
    anonymous_network_detected   = true

    user_attribute_equals {
      attribute_reference = "$${user.mfaEnabled}"
      value_boolean       = true
    }

    user_attribute_equals {
      attribute_reference = "$${user.lifecycle.status}"
      value               = "ACCOUNT_OK"
    }
  }

  mfa {
    device_sign_on_policy_id = var.my_device_sign_on_policy_id
    no_device_mode           = "BYPASS"
  }
}
```

## Example Usage - Identity Provider

```terraform
resource "pingone_sign_on_policy_action" "my_policy_identity_provider" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  identity_provider {
    identity_provider_id = pingone_identity_provider.my_identity_provider.id

    acr_values        = "MFA"
    pass_user_context = true
  }
}
```

## Example Usage - Progressive Profiling

```terraform
resource "pingone_sign_on_policy_action" "my_policy_progressive_profiling" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 3

  progressive_profiling {

    attribute {
      name     = "name.given"
      required = false
    }

    attribute {
      name     = "name.family"
      required = true
    }

    prompt_text = "For the best experience, we need a couple things from you."

  }
}
```

## Example Usage - Agreement

```terraform
resource "pingone_sign_on_policy_action" "my_policy_agreement" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 3

  agreement {
    agreement_id        = var.my_agreement_id
    show_decline_option = false
  }
}
```

## Example Usage - PingID Windows Login Passwordless (Workforce Environments)

```terraform
resource "pingone_sign_on_policy_action" "my_policy_pingid_windows_login_passwordless" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  pingid_windows_login_passwordless {
    unique_user_attribute_name = "externalId"
    offline_mode_enabled       = true
  }
}
```

## Example Usage - PingID (Workforce Environments)

```terraform
resource "pingone_sign_on_policy_action" "my_policy_pingid" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  pingid {}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to create the sign on policy action in.
- `priority` (Number) An integer that specifies the order in which the policy referenced by this assignment is evaluated during an authentication flow relative to other policies. An assignment with a lower priority will be evaluated first.
- `sign_on_policy_id` (String) The ID of the sign on policy to associate the sign on policy action to.

### Optional

- `agreement` (Block List, Max: 1) Options specific to the **Agreements** policy action. (see [below for nested schema](#nestedblock--agreement))
- `conditions` (Block List, Max: 1) Conditions to apply to the sign on policy action.  Applies to policy actions of type `agreement`, `identifier_first`, `identity_provider`, `login`, `mfa`, `progressive_profiling`. (see [below for nested schema](#nestedblock--conditions))
- `enforce_lockout_for_identity_providers` (Boolean) A boolean that if set to true and if the user's account is locked (the account.canAuthenticate attribute is set to false), then social sign on with an external identity provider is prevented. Defaults to `false`.
- `identifier_first` (Block List, Max: 1) Options specific to the **Identifier First** policy action. (see [below for nested schema](#nestedblock--identifier_first))
- `identity_provider` (Block List, Max: 1) Options specific to the **Identity Provider** policy action. (see [below for nested schema](#nestedblock--identity_provider))
- `login` (Block List, Max: 1) Options specific to the **Login** policy action. (see [below for nested schema](#nestedblock--login))
- `mfa` (Block List, Max: 1) Options specific to the **Multi-factor Authentication** policy action. (see [below for nested schema](#nestedblock--mfa))
- `pingid` (Block List, Max: 1) Options specific to the **PingID** policy action.  This action can only be applied to Workforce solution context environments that have the PingID and SSO services enabled. (see [below for nested schema](#nestedblock--pingid))
- `pingid_windows_login_passwordless` (Block List, Max: 1) Options specific to the **PingID Windows Login Passwordless** policy action.  This action can only be applied to Workforce solution context environments that have the PingID and SSO services enabled. (see [below for nested schema](#nestedblock--pingid_windows_login_passwordless))
- `progressive_profiling` (Block List, Max: 1) Options specific to the **Progressive Profiling** policy action. (see [below for nested schema](#nestedblock--progressive_profiling))
- `registration_confirm_user_attributes` (Boolean) A boolean that specifies whether users must confirm data returned from an identity provider prior to registration. Users can modify the data and omit non-required attributes. Modified attributes are added to the user's profile during account creation. Defaults to `false`.
- `registration_external_href` (String) A string that specifies the link to the external identity provider's identity store. This property is set when the administrator chooses to have users register in an external identity store. This attribute can be set only when the registration.enabled property is set to false.
- `registration_local_population_id` (String) A string that specifies the population ID associated with the newly registered user. Setting this enables local registration features.
- `social_provider_ids` (Set of String) One or more IDs of the identity providers that can be used for the social login sign-on flow.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--agreement"></a>
### Nested Schema for `agreement`

Required:

- `agreement_id` (String) A string that specifies the ID of the agreement to which the user must consent.

Optional:

- `show_decline_option` (Boolean) When enabled, the `Do Not Accept` button will terminate the Flow and display an error message to the user. Defaults to `true`.


<a id="nestedblock--conditions"></a>
### Nested Schema for `conditions`

Optional:

- `anonymous_network_detected` (Boolean) A boolean that specifies whether the user should be prompted for re-authentication on this action based on a detected anonymous network.  Applies to policy actions of type `mfa`. Defaults to `false`.
- `anonymous_network_detected_allowed_cidr` (Set of String) A list of allowed CIDR when an anonymous network is detected.  Applies to policy actions of type `mfa`.
- `geovelocity_anomaly_detected` (Boolean) A boolean that specifies whether the user should be prompted for re-authentication on this action based on a detected geovelocity anomaly.  Applies to policy actions of type `mfa`. Defaults to `false`.
- `ip_out_of_range_cidr` (Set of String) A list of strings that specifies the supported network IP addresses expressed as classless inter-domain routing (CIDR) strings.  Applies to policy actions of type `mfa`.
- `ip_reputation_high_risk` (Boolean) A boolean that specifies whether the user's IP risk should be used when evaluating this policy action.  A value of `HIGH` will prompt the user to authenticate with this action.  Applies to policy actions of type `mfa`. Defaults to `false`.
- `last_sign_on_older_than_seconds` (Number) Set the number of seconds by which the user will not be prompted for this action following the last successful authentication.  Applies to policy actions of type `identifier_first`, `identity_provider`, `login`, `mfa`.
- `last_sign_on_older_than_seconds_mfa` (Number) Set the number of seconds by which the user will not be prompted for this action following the last successful authentication of an MFA authenticator device.  Applies to policy actions of type `mfa`.
- `user_attribute_equals` (Block Set) One or more conditions where an attribute on the user's profile must match the configured value.  Applies to policy actions of type `identifier_first`, `login`, `mfa`, but cannot be set on policy actions where the priority is `1`. (see [below for nested schema](#nestedblock--conditions--user_attribute_equals))
- `user_is_member_of_any_population_id` (Set of String) Activate this action only for users within the specified list of population IDs.  Applies to policy actions of type `identifier_first`, `login`, `mfa`, but cannot be set on policy actions where the priority is `1`.

<a id="nestedblock--conditions--user_attribute_equals"></a>
### Nested Schema for `conditions.user_attribute_equals`

Required:

- `attribute_reference` (String) Specifies the user attribute used in the condition. Only string core, standard, and custom attributes are supported. For complex attribute types, you must reference the sub-attribute (`$${user.name.firstName}`).  Note values that begin with a dollar sign (`$`) must be prefixed with an additional dollar sign.  E.g. `${name.given}` should be configured as `$${name.given}`.  When configured, one of `value` (for attributes of type `STRING` or `INTEGER`) or `value_boolean` (for attributes of type `BOOLEAN`) must be provided.

Optional:

- `value` (String) The string or integer (as string) value of the attribute (declared in `attribute_reference`) on the user profile that should be matched.  This value parameter should be used where the data type of the schema attribute in `attribute_reference` is of type `STRING` or `INTEGER`.  Conflicts with `value_boolean`.
- `value_boolean` (Boolean) The boolean value of the attribute (declared in `attribute_reference`) on the user profile that should be matched.  This value parameter should be used where the data type of the schema attribute in `attribute_reference` is of type `BOOLEAN` (e.g `$${user.emailVerified}`, `$${user.verified}` and `$${user.mfaEnabled}`).  Conflicts with `value`.



<a id="nestedblock--identifier_first"></a>
### Nested Schema for `identifier_first`

Optional:

- `discovery_rule` (Block Set, Max: 100) One or more IDP discovery rules invoked when no user is associated with the user identifier. The condition on which this identity provider is used to authenticate the user is expressed using the PingOne policy condition language. (see [below for nested schema](#nestedblock--identifier_first--discovery_rule))
- `recovery_enabled` (Boolean) A boolean that specifies whether account recovery features are active on the policy action. Defaults to `true`.

<a id="nestedblock--identifier_first--discovery_rule"></a>
### Nested Schema for `identifier_first.discovery_rule`

Required:

- `attribute_contains_text` (String) Text to match on a user's username. Any users that don't match a discovery rule will authenticate against PingOne.  E.g `@pingidentity.com`
- `identity_provider_id` (String) The ID that specifies the identity provider that will be used to authenticate the user if the condition is matched.



<a id="nestedblock--identity_provider"></a>
### Nested Schema for `identity_provider`

Required:

- `identity_provider_id` (String) A string that specifies the ID of the external identity provider to which the user is redirected for sign-on.

Optional:

- `acr_values` (String) A string that designates the sign-on policies included in the authorization flow request. Options can include the PingOne predefined sign-on policies, Single_Factor and Multi_Factor, or any custom defined sign-on policy names. Sign-on policy names should be listed in order of preference, and they must be assigned to the application. This property can be configured on the identity provider action and is passed to the identity provider if the identity provider is of type `SAML` or `OPENID_CONNECT`.
- `pass_user_context` (Boolean) A boolean that specifies whether to pass in a login hint to the identity provider on the sign on request. Based on user context, the login hint is set if (1) the user is set on the flow, and (2) the user already has an account link for the identity provider. If both of these conditions are true, then the user is sent to the identity provider with a login hint equal to their externalId for the identity provider (saved on the account link). If these conditions are not true, then the API checks see if there is an OIDC login hint on the flow. If so, that login hint is used. If none of these conditions are true, the login hint parameter is not included on the authorization request to the identity provider.


<a id="nestedblock--login"></a>
### Nested Schema for `login`

Optional:

- `new_user_provisioning` (Block List, Max: 1) Enables user entries existing outside of PingOne to be provisioned during login, using an external integration solution (such as a Gateway). (see [below for nested schema](#nestedblock--login--new_user_provisioning))
- `recovery_enabled` (Boolean) A boolean that specifies whether account recovery features are active on the policy action. Defaults to `true`.

<a id="nestedblock--login--new_user_provisioning"></a>
### Nested Schema for `login.new_user_provisioning`

Required:

- `gateway` (Block Set, Min: 1) One or more blocks that describe a preconfigured gateway and user type that are specified in the Gateway Management schema to determine how to find and migrate user entries existing in an external directory. (see [below for nested schema](#nestedblock--login--new_user_provisioning--gateway))

<a id="nestedblock--login--new_user_provisioning--gateway"></a>
### Nested Schema for `login.new_user_provisioning.gateway`

Required:

- `id` (String) A string that specifies the UUID ID of the gateway instance.  The ID may come from the `id` parameter of the `pingone_gateway` resource.  Must be a valid PingOne resource ID.
- `user_type_id` (String) A string that specifies the UUID ID of the user type within the gateway instance.  The ID may come from the `user_type[*].id` parameter of the `pingone_gateway` resource.  Must be a valid PingOne resource ID.

Optional:

- `type` (String) A string that specifies the type of the gateway. Currently, only `LDAP` is supported. Defaults to `LDAP`.




<a id="nestedblock--mfa"></a>
### Nested Schema for `mfa`

Required:

- `device_sign_on_policy_id` (String) The ID of the MFA policy that should be used.

Optional:

- `no_device_mode` (String) A string that specifies the device mode for the MFA flow. Options are `BYPASS` to allow MFA without a specified device, or `BLOCK` to block the MFA flow if no device is specified. To use this configuration option, the authorize request must include a signed `login_hint_token` property. For more information, see Authorize (Browserless and MFA Only Flows). Defaults to `BLOCK`.


<a id="nestedblock--pingid"></a>
### Nested Schema for `pingid`


<a id="nestedblock--pingid_windows_login_passwordless"></a>
### Nested Schema for `pingid_windows_login_passwordless`

Required:

- `offline_mode_enabled` (Boolean) A boolean that specifies whether to allow users to log in when PingOne and or PingID are not available.
- `unique_user_attribute_name` (String) A string that specifies the schema attribute to match against the provided identifier when searching for a user in the directory. Only unique attributes in the directory schema may be configured.


<a id="nestedblock--progressive_profiling"></a>
### Nested Schema for `progressive_profiling`

Required:

- `attribute` (Block Set, Min: 1) One or more attribute(s) that the user should be prompted to complete as part of the progressive profiling action. (see [below for nested schema](#nestedblock--progressive_profiling--attribute))
- `prompt_text` (String) A string that specifies text to display to the user when prompting for attribute values.

Optional:

- `prevent_multiple_prompts_per_flow` (Boolean) A boolean that specifies whether the progressive profiling action will not be executed if another progressive profiling action has already been executed during the flow. Defaults to `true`.
- `prompt_interval_seconds` (Number) An integer that specifies how often to prompt the user to provide profile data for the configured attributes for which they do not have values. Defaults to `7776000`.

<a id="nestedblock--progressive_profiling--attribute"></a>
### Nested Schema for `progressive_profiling.attribute`

Required:

- `name` (String) A string that specifies the name and path of the user profile attribute as defined in the user schema (for example, email or address.postalCode).
- `required` (Boolean) A boolean that specifies whether the user is required to provide a value for the attribute.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_sign_on_policy_action.example <environment_id>/<sign_on_policy_id>/<sign_on_policy_action_id>
```
