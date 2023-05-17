resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Location Predictor"
  compact_name   = "my_awesome_user_location_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_location_anomaly = {
    radius = {
      distance = 100
      unit     = "miles"
    }
  }
}