resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "my_awesome_anonymous_network_predictor"

  predictor_anonymous_network = {}
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "my_awesome_geovelocity_predictor"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "my_awesome_composite_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Composite Predictor"
  compact_name   = "my_awesome_composite_predictor"

  predictor_composite = {
    composition = {
      level = "HIGH"

      condition_json = jsonencode({
        "not" : {
          "or" : [{
            "equals" : 0,
            "value" : "$${details.counters.predictorLevels.medium}",
            "type" : "VALUE_COMPARISON"
            }, {
            "equals" : "High",
            "value" : "$${details.${pingone_risk_predictor.my_awesome_geovelocity_predictor.compact_name}.level}",
            "type" : "VALUE_COMPARISON"
            }, {
            "and" : [{
              "equals" : "High",
              "value" : "$${details.${pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name}.level}",
              "type" : "VALUE_COMPARISON"
            }],
            "type" : "AND"
          }],
          "type" : "OR"
        },
        "type" : "NOT"
      })
    }
  }
}