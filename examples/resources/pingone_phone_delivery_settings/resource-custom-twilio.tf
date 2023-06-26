resource "pingone_phone_delivery_settings" "my_awesome_custom_twilio_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_custom_twilio = {
    auth_token = var.twilio_auth_token
    sid        = var.twilio_sid

    selected_numbers = [
      {
        number = var.my_twilio_number
        type   = "PHONE_NUMBER"
      }
    ]
  }
}