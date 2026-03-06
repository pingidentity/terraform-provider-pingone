resource "pingone_davinci_connector_instance" "credovaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "credovaConnector"
  }
  name = "My awesome credovaConnector"
  properties = jsonencode({
    "baseUrl" = var.credovaconnector_property_base_url
    "password" = var.credovaconnector_property_password
    "username" = var.credovaconnector_property_username
  })
}
