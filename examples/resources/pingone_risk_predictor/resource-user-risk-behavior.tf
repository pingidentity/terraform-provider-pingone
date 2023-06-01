resource "pingone_risk_predictor" "my_awesome_user_risk_behavior_predictor_by_user" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Risk Behavior Predictor By User"
  compact_name   = "myAwesomeUserRiskBehaviorPredictorByUser"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "points"
    }
  }
}

resource "pingone_risk_predictor" "my_awesome_user_risk_behavior_predictor_by_organization" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Risk Behavior Predictor By Organization"
  compact_name   = "myAwesomeUserRiskBehaviorPredictorByOrganization"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "login_anomaly_statistic"
    }
  }
}