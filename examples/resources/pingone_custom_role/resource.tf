resource "pingone_environment" "my_environment" {
  # ...
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
	    # Default "Custom Roles Admin" administrator role id
      id = "6f770b08-793f-4393-b2aa-b1d1587a0324"
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
