resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_application_role_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
