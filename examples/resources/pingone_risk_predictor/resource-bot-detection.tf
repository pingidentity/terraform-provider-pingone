resource "pingone_risk_predictor" "my_awesome_bot_detection_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Bot Detection Predictor"
  compact_name   = "myAwesomeBotDetectionPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_bot_detection = {}
}