resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "user" {
  environment_id = pingone_environment.my_environment.id

  name = "User Quota SMS and Voice"

  quota {
    type  = "USER"
    total = 30
  }
}