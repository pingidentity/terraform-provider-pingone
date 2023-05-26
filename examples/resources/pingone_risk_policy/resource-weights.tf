resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

resource "pingone_risk_policy" "my_awesome_weights_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Weights-based Risk Policy"

  policy_weights = {
    policy_threshold_medium = {
      min_score = 2
    }

    policy_threshold_high = {
      min_score = 5
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        weight       = 5
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        weight       = 5
      }
    ]
  }
}