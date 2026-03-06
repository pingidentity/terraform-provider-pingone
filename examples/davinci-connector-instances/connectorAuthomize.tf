resource "pingone_davinci_connector_instance" "connectorAuthomize" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAuthomize"
  }
  name = "My awesome connectorAuthomize"
  properties = jsonencode({
    "apiKey" = var.connectorauthomize_property_api_key
  })
}
