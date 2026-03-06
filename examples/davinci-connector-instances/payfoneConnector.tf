resource "pingone_davinci_connector_instance" "payfoneConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "payfoneConnector"
  }
  name = "My awesome payfoneConnector"
  properties = jsonencode({
    "appClientId" = var.payfoneconnector_property_app_client_id
    "baseUrl" = var.payfoneconnector_property_base_url
    "clientId" = var.payfoneconnector_property_client_id
    "password" = var.payfoneconnector_property_password
    "simulatorMode" = var.payfoneconnector_property_simulator_mode
    "simulatorPhoneNumber" = var.payfoneconnector_property_simulator_phone_number
    "skCallbackBaseUrl" = var.payfoneconnector_property_sk_callback_base_url
    "username" = var.payfoneconnector_property_username
  })
}
