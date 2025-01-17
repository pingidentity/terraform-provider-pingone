resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_custom_role" "my_custom_role" {
  environment_id = pingone_environment.my_environment.id

  name        = "My custom role"
  description = "My custom role for reading role assignments"

  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]

  can_be_assigned_by = [
    {
      id = pingone_role.environment_admin.id
    }
  ]

  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    },
    {
      id = "permissions:read:groupRoleAssignments"
    },
  ]

}
