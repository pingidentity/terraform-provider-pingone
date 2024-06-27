resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_alert_channel" "my_awesome_alert_channel" {
  environment_id = pingone_environment.my_environment.id

  alert_name = "My awesome alert channel"

  addresses = [
    "iam_license_admins@bxretail.org",
  ]

  channel_type = "EMAIL"

  include_alert_types = [
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_ROTATED",
  ]

  include_severities = [
    "INFO",
    "WARNING",
    "ERROR",
  ]
}
