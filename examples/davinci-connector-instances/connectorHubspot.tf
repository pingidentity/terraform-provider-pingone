resource "pingone_davinci_connector_instance" "connectorHubspot" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHubspot"
  }
  name = "My awesome connectorHubspot"
  properties = jsonencode({
    "bearerToken" = var.connectorhubspot_property_bearer_token
  })
}
