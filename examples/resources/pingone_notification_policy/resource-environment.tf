resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "environment" {
  environment_id = pingone_environment.my_environment.id

  name = "Environment Quota SMS and Voice"

  quota {
    type  = "ENVIRONMENT"
    total = 100
  }
}