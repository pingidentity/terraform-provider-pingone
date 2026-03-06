resource "pingone_davinci_connector_instance" "connectorFreshservice" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshservice"
  }
  name = "My awesome connectorFreshservice"
  properties = jsonencode({
    "apiKey" = var.connectorfreshservice_property_api_key
    "domain" = var.connectorfreshservice_property_domain
  })
}
