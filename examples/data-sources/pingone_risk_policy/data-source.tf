data "pingone_risk_policy" "example" {
  environment_id = var.environment_id
  name           = "Risk Policy X"
}

data "pingone_risk_policy" "example_by_id" {
  environment_id = var.environment_id
  risk_policy_id = var.risk_policy_id
}
