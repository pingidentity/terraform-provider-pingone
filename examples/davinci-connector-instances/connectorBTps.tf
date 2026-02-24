resource "pingone_davinci_connector_instance" "connectorBTps" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTps"
  }
  name = "My awesome connectorBTps"
  properties = jsonencode({
    "apiKey" = var.connectorbtps_property_api_key
    "apiUser" = var.connectorbtps_property_api_user
    "domain" = var.connectorbtps_property_domain
  })
}
