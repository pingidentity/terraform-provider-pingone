resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.my_environment.id
  custom_provider_name = "My Custom Email Provider"

  auth_url      = var.custom_provider_auth_url
  grant_type    = "CLIENT_CREDENTIALS"
  client_id     = var.custom_provider_client_id
  client_secret = var.custom_provider_client_secret
  scopes        = ["mail.send"]
  protocol      = "HTTP"

  from = {
    name          = "From Services"
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    name          = "Reply To Services"
    email_address = "services@bxretail.org"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject" : "$${subject}"
        "reply-to" : "$${reply_to}"
        "from" : "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.bxretail.org/send-email"
    }
  ]
}
