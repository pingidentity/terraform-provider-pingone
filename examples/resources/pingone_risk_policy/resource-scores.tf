resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
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
        compact_name = pingone_risk_predictor.my_awesome_user_location_predictor.compact_name
        score        = 50
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "LOW"
      }

      condition = {
        type = "IP_RANGE"
        ip_range = [
          "10.0.0.0/8",
        ]
      }
    },

    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        equals       = "HIGH"
      }
    }
  ]
}