resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  auth_token = var.twilio_auth_token
  sid        = var.twilio_sid
}