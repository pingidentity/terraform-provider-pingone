resource "pingone_davinci_connector_instance" "connectorIPregistry" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPregistry"
  }
  name = "My awesome connectorIPregistry"
  properties = jsonencode({
    "apiKey" = var.connectoripregistry_property_api_key
  })
}
