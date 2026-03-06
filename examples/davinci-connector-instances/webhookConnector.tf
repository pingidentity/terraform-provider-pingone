resource "pingone_davinci_connector_instance" "webhookConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "webhookConnector"
  }
  name = "My awesome webhookConnector"
  properties = jsonencode({
    "urls" = var.webhookconnector_property_urls
  })
}
