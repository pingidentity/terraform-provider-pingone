resource "pingone_davinci_connector_instance" "connectorSpycloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpycloud"
  }
  name = "My awesome connectorSpycloud"
  properties = jsonencode({
    "apiKey" = var.connectorspycloud_property_api_key
  })
}
