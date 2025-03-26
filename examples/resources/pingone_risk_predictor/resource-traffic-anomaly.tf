resource "pingone_risk_predictor" "my_traffic_anomaly_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "my test traffic anomaly predictor"
  compact_name   = "myAwesometrafficMPredictor"

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
