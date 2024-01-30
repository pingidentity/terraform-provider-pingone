data "pingone_application_flow_policy_assignments" "all_flow_policy_assignments_by_app" {
  environment_id = var.environment_id

  application_id = pingone_application.my_awesome_application.id
}
