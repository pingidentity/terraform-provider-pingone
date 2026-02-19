data "pingone_application_role_assignments" "all_role_assignments_by_app" {
  environment_id = var.environment_id

  application_id = pingone_application.my_awesome_application.id
}
