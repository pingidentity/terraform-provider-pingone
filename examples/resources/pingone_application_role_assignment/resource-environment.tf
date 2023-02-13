resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_application_role_assignment" "single_environment_admin_to_application" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id
  role_id        = data.pingone_role.environment_admin.id

  scope_environment_id = pingone_environment.my_environment.id
}
