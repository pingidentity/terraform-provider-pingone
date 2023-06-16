resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "CUSTOM_SYNIVERSE"

  auth_token = var.syniverse_auth_token
}