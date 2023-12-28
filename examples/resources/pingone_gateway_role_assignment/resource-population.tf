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

  scope_population_id = pingone_population.my_population.id
}
