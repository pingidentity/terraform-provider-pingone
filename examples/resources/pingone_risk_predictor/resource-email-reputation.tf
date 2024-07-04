resource "pingone_risk_predictor" "my_awesome_email_reputation_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Email Reputation Predictor"
  compact_name   = "myAwesomeEmailReputationPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_email_reputation = {}
}