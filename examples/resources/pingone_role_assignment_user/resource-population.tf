resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_role_assignment_user" "population_identity_data_admin_to_user" {
  environment_id = pingone_environment.my_environment.id
  user_id        = var.user_id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
