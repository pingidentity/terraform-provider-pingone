resource "pingone_environment" "my_environment" {
  # ...
}

# Example custom SMTP server configuration
resource "pingone_notification_settings_email" "my_awesome_smtp_settings" {
  environment_id = pingone_environment.my_environment.id

  host     = "smtp-example.bxretail.org"
  port     = 25
  username = var.smtp_server_username
  password = var.smtp_server_password

  from = {
    email_address = "services@bxretail.org"
    name          = "Customer Services"
  }
}

# Example custom email provider configuration
resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.test_language_env.id
  custom_provider_name = "My Custom Email Provider"

  username = var.custom_provider_username
  password = var.custom_provider_password
  protocol = "HTTP"

  from = {
    name          = "Customer Services"
    email_address = "services@bxretail.org"
  }

  reply_to = {
    name          = "Customer Services"
    email_address = "customerservices@bxretail.org"
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