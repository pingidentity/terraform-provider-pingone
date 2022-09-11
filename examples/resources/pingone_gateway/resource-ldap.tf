resource "pingone_environment" "my_environment" {
  # ...
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

  user_type {
    name               = "User Set 1"
    password_authority = "LDAP"
    search_base_dn     = "ou=users,dc=bxretail,dc=org"

    user_link_attributes = ["objectGUID", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

      population_id = pingone_environment.my_environment.default_population_id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.sAMAccountName}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }
    }

    push_password_changes_to_ldap = true
  }

}