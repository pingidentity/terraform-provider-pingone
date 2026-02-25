resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings" "my_awesome_ping_smtp_notification_settings" {
  environment_id = pingone_environment.my_environment.id

  from = {
    email_address = "noreply@pingidentity.com"
  }

  reply_to = {
    email_address = "noreply@pingidentity.com"
  }
}
