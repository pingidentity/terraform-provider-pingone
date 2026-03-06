resource "pingone_davinci_connector_instance" "mailchainConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mailchainConnector"
  }
  name = "My awesome mailchainConnector"
  properties = jsonencode({
    "apiKey" = var.mailchainconnector_property_api_key
    "version" = var.mailchainconnector_property_version
  })
}
