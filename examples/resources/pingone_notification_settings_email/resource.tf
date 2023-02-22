resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_settings" "mfa_settings" {
  environment_id = pingone_environment.my_environment.id

  host     = "smtp-example.bxretail.org"
  port     = 25
  username = var.smtp_server_username
  password = var.smtp_server_password

  from {
    email_address = "services@bxretail.org"
    name          = "Customer Services"
  }
}