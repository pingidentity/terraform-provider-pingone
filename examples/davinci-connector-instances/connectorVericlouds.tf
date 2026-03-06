resource "pingone_davinci_connector_instance" "connectorVericlouds" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorVericlouds"
  }
  name = "My awesome connectorVericlouds"
  properties = jsonencode({
    "apiSecret" = var.connectorvericlouds_property_api_secret
    "apikey" = var.connectorvericlouds_property_apikey
  })
}
