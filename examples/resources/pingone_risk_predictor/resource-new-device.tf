resource "pingone_risk_predictor" "my_awesome_new_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome New Device Predictor"
  compact_name   = "my_awesome_new_device_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = "2023-05-01T00:00:00Z"
  }
}