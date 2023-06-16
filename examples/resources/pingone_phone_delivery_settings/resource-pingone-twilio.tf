resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "PINGONE_TWILIO"
}