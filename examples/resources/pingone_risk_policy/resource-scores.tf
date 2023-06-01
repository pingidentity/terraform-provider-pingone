resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

resource "pingone_risk_policy" "my_awesome_scores_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Scores-based Risk Policy"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 40
    }

    policy_threshold_high = {
      min_score = 75
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        score        = 50
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        score        = 50
      }
    ]
  }
}