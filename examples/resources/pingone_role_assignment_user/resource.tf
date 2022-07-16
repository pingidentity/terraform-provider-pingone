resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  license_id  = var.license_id
  default_population {}
  service {}
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_role_assignment_user" "foo" {
  environment_id = pingone_environment.my_environment.id
  user_id        = var.user_id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
