resource "pingone_davinci_connector_instance" "treasureDataConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "treasureDataConnector"
  }
  name = "My awesome treasureDataConnector"
  properties = jsonencode({
    "apiKey" = var.treasuredataconnector_property_api_key
  })
}
