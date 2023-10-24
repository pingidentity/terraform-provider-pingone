resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_group" {
  # ...
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_group_role_assignment" "single_environment_admin_to_group" {
  environment_id = pingone_environment.my_environment.id
  group_id = pingone_group.my_group.id
  role_id        = data.pingone_role.environment_admin.id

  scope_environment_id = pingone_environment.my_environment.id
}
