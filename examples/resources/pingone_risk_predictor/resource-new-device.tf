resource "time_static" "my_awesome_new_device_predictor_activation" {}

resource "pingone_risk_predictor" "my_awesome_new_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome New Device Predictor"
  compact_name   = "myAwesomeNewDevicePredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = format("%sT00:00:00Z", formatdate("YYYY-MM-DD", time_static.my_awesome_new_device_predictor_activation.rfc3339))
  }
}
