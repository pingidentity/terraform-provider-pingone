resource "pingone_risk_predictor" "my_awesome_ip_reputation_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome IP Reputation Predictor"
  compact_name   = "my_awesome_ip_reputation_predictor"

  default_result {
    weight    = ""
    score     = ""
    evaluated = ""
    result    = ""
  }

  predictor_ip_reputation {

  }
}