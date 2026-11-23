resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

# A mitigation-only risk policy (no policy_weights / policy_scores blocks), as
# supported on MFA-only licensed environments.
resource "pingone_risk_policy" "my_awesome_mitigations_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Mitigations-based Risk Policy"

  mitigations = [
    {
      action        = "CUSTOM"
      custom_action = "Prompt for additional authentication"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        equals       = "HIGH"
      }
    },

    {
      action = "DENY"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        equals       = "HIGH"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }

  targets = {
    condition = {
      and = [
        {
          list     = ["AUTHENTICATION", "AUTHORIZATION"]
          contains = "$${event.flow.type}"
        },
      ]
    }
  }
}
