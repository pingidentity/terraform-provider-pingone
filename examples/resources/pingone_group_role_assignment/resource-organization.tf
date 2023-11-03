resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_group" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Group"
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_group_role_assignment" "organization_environment_admin_to_group" {
  environment_id = pingone_environment.my_environment.id
  group_id       = pingone_group.my_group.id
  role_id        = data.pingone_role.environment_admin.id

  scope_organization_id = var.organization_id
}
