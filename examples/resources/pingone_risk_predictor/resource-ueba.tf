resource "pingone_risk_predictor" "my_awesome_ueba_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome UEBA Predictor"
  compact_name   = "my_awesome_ueba_predictor"

  default_result {
    weight    = ""
    score     = ""
    evaluated = ""
    result    = ""
  }

  predictor_ueba {

  }
}