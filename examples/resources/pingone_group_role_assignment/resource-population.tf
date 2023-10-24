resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_group" {
  # ...
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_group_role_assignment" "population_identity_data_admin_to_group" {
  environment_id = pingone_environment.my_environment.id
  group_id = pingone_group.my_group.id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
