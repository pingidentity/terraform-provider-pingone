resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
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
      min_score = 50
    }

    policy_threshold_high = {
      min_score = 60
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        weight       = 5
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_user_location_predictor.compact_name
        weight       = 5
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