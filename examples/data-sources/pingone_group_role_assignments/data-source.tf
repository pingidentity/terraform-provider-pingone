data "pingone_group_role_assignments" "all_role_assignments_by_group" {
  environment_id = var.environment_id

  groupid = pingone_group.my_awesome_group.id
}
