data "pingone_risk_predictor" "example" {
  environment_id = var.environment_id
  name           = "Risk Predictor X"
}

data "pingone_risk_predictor" "example_by_id" {
  environment_id    = var.environment_id
  risk_predictor_id = var.risk_predictor_id
}
