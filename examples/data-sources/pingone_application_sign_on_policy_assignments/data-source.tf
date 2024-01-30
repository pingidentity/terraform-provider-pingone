data "pingone_application_sign_on_policy_assignments" "all_sop_assignments_by_app" {
  environment_id = var.environment_id

  application_id = pingone_application.my_awesome_application.id
}
