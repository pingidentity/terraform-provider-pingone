resource "pingone_davinci_connector_instance" "connectorFreshdesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshdesk"
  }
  name = "My awesome connectorFreshdesk"
  properties = jsonencode({
    "apiKey" = var.connectorfreshdesk_property_api_key
    "baseURL" = var.base_url
    "version" = var.connectorfreshdesk_property_version
  })
}
