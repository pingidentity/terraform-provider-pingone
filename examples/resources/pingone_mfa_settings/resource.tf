resource "pingone_mfa_settings" "mfa_settings" {
  environment_id = pingone_environment.my_environment.id

  pairing {
    max_allowed_devices = 5
    pairing_key_format  = "ALPHANUMERIC"
  }

  lockout {
    failure_count    = 5
    duration_seconds = 600
  }

}