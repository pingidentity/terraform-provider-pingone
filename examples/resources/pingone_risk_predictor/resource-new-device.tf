resource "pingone_risk_predictor" "my_awesome_new_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome New Device Predictor"
  compact_name   = "my_awesome_new_device_predictor"

  default_result {
    weight    = ""
    score     = ""
    evaluated = ""
    result    = ""
  }

  predictor_new_device {

  }
}