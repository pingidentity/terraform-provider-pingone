resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_gateway" "my_awesome_pingfederate_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "Advanced Services SSO"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_gateway_role_assignment" "population_identity_data_admin_to_gateway" {
  environment_id = pingone_environment.my_environment.id
  gateway_id     = pingone_gateway.my_awesome_pingfederate_gateway.id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
