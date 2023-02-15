resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_gateway" "my_awesome_pingfederate_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "Advanced Services SSO"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_gateway_role_assignment" "single_environment_admin_to_application" {
  environment_id = pingone_environment.my_environment.id
  gateway_id     = pingone_gateway.my_awesome_pingfederate_gateway.id
  role_id        = data.pingone_role.environment_admin.id

  scope_environment_id = pingone_environment.my_environment.id
}
