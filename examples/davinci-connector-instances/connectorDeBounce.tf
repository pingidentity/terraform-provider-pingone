resource "pingone_davinci_connector_instance" "connectorDeBounce" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorDeBounce"
  }
  name = "My awesome connectorDeBounce"
  properties = jsonencode({
    "apiKey" = var.connectordebounce_property_api_key
  })
}
