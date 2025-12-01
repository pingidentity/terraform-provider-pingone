resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "myAwesomeAnonymousNetworkPredictor"

  predictor_anonymous_network = {}
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "myAwesomeGeovelocityPredictor"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "my_awesome_composite_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Composite Predictor"
  compact_name   = "myAwesomeCompositePredictor"

  predictor_composite = {
    compositions = [
      {
        level = "HIGH"

        condition_json = jsonencode({
          "not" : {
            "or" : [{
              "equals" : 0,
              "value" : "$${details.counters.predictorLevels.medium}",
              "type" : "VALUE_COMPARISON"
              }, {
              "equals" : "High",
              "value" : "$${details.${pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name}.level}",
              "type" : "VALUE_COMPARISON"
              }, {
              "startsWith" : "admin",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
              }, {
              "endsWith" : "@contractor.example.com",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
              }, {
              "and" : [{
                "equals" : "High",
                "value" : "$${details.${pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name}.level}",
                "type" : "VALUE_COMPARISON"
                },
                {
                  "list" : ["Group Name"],
                  "contains" : "$${event.user.groups}",
                  "type" : "GROUPS_INTERSECTION"
              }],
              "type" : "AND"
            }],
            "type" : "OR"
          },
          "type" : "NOT"
        })
      }
    ]
  }
}