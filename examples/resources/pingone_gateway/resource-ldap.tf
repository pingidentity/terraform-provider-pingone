resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_population" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Population"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}

resource "pingone_gateway" "my_ldap_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Active Directory"
  enabled        = true
  type           = "LDAP"

  bind_dn       = var.bind_dn
  bind_password = var.bind_password

  connection_security = "TLS"
  vendor              = "Microsoft Active Directory"

  servers = [
    "ds1.bxretail.org:636",
    "ds2.bxretail.org:636",
    "ds3.bxretail.org:636",
  ]

  user_types = {
    "User Set 1" = {
      password_authority = "LDAP"
      search_base_dn     = "ou=users,dc=bxretail,dc=org"

      user_link_attributes = ["objectGUID", "objectSid"]

      new_user_lookup = {
        ldap_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

        population_id = pingone_population.my_population.id

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.sAMAccountName}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          }
        ]
      }

      update_user_on_successful_authentication = true
    }
  }

}