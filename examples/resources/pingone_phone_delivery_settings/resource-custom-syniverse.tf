resource "pingone_phone_delivery_settings" "my_awesome_custom_syniverse_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_custom_syniverse = {
    auth_token = var.syniverse_auth_token

    selected_numbers = [
      {
        number = var.my_syniverse_number
        type   = "PHONE_NUMBER"
      }
    ]
  }
}