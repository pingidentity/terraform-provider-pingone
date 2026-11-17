resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_risk_settings" "my_risk_settings" {
  environment_id = pingone_environment.my_environment.id

  data_retention = {
    new_device_in_days               = 30
    user_based_risk_behavior_in_days = 90
    user_location_anomaly_in_days    = 45
  }
}
