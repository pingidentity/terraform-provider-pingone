resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_user_role_assignment" "organization_environment_admin_to_user" {
  environment_id = pingone_environment.my_environment.id
  user_id        = var.user_id
  role_id        = data.pingone_role.environment_admin.id

  scope_organization_id = var.organization_id
}
