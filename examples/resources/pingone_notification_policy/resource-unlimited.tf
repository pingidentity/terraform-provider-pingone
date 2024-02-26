resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "unlimited" {
  environment_id = pingone_environment.my_environment.id

  name = "Unlimited Quotas SMS Voice and Email"
}