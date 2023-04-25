resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "my_awesome_anonymous_network_predictor"

  default_result {
    weight    = ""
    score     = ""
    evaluated = ""
    result    = ""
  }

  predictor_anonymous_network {
    allowed_cidr_list = []
  }
}