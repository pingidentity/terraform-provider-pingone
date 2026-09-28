data "pingone_phone_delivery_settings" "example_by_id" {
  environment_id             = var.environment_id
  phone_delivery_settings_id = var.phone_delivery_settings_id
}

data "pingone_phone_delivery_settings" "example_by_name" {
  environment_id = var.environment_id

  name = "My awesome custom notifications provider"
}
