data "pingone_user_role_assignments" "awesome_admin_user" {
  environment_id = var.admin_users_environment_id

  user_id = var.awesome_admin_user_id
}
