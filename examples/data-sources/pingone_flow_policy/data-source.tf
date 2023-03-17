data "pingone_flow_policy" "example_by_id" {
  environment_id = var.environment_id

  flow_policy_id = var.davinci_flow_policy_id
}