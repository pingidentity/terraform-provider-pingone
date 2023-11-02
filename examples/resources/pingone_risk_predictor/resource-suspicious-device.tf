resource "pingone_risk_predictor" "my_awesome_suspicious_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Suspicious Device Predictor"
  compact_name   = "myAwesomeSuspiciousDevicePredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect = "SUSPICIOUS_DEVICE"
  }
}