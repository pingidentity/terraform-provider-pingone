resource "pingone_risk_predictor" "my_traffic_anomaly_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Traffic Anomaly Predictor"
  compact_name   = "myAwesomeTrafficAnomalyPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_traffic_anomaly = {
    rules = [
      {
        type = "UNIQUE_USERS_PER_DEVICE"
        threshold = {
          medium = 3
          high   = 4
        }
        interval = {
          unit     = "DAY"
          quantity = 1
        }
        enabled = true
      }
    ]
  }
}
