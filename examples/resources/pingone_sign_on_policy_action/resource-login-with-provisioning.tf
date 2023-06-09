resource "pingone_gateway" "my_awesome_ldap_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome LDAP Gateway"

  # ...

  user_type {
    name = "User Set 1"

    # ...
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
        user_type_id = pingone_gateway.my_awesome_ldap_gateway.user_type.* [index(pingone_gateway.my_awesome_ldap_gateway.user_type[*].name, "User Set 1")].id
      }
    }
  }

  depends_on = [
    pingone_gateway_credential.my_awesome_ldap_gateway
  ]
}
