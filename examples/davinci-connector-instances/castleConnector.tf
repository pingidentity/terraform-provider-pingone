resource "pingone_davinci_connector_instance" "castleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "castleConnector"
  }
  name = "My awesome castleConnector"
  properties = jsonencode({
    "apiSecret" = var.castleconnector_property_api_secret
  })
}
