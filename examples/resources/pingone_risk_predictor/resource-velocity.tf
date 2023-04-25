resource "pingone_risk_predictor" "my_awesome_velocity_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor"
  compact_name   = "my_awesome_velocity_predictor"

  default_result {
    weight    = ""
    score     = ""
    evaluated = ""
    result    = ""
  }

  predictor_velocity {

  }
}