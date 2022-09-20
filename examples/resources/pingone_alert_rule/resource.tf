resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_gateway" "my_awesome_alert_rule" {
  environment_id = pingone_environment.my_environment.id

  addresses          = ["myaddress@bxretail.org"]
  include_severities = ["INFO", "WARNING", "ERROR"]
  include_alert_types = [
    "KEY_PAIR_EXPIRING",
    "KEY_PAIR_EXPIRED",
    "CERTIFICATE_EXPIRED",
    "CERTIFICATE_EXPIRING",
    "GATEWAY_VERSION_DEPRECATED",
    "GATEWAY_VERSION_DEPRECATING"
  ]
}