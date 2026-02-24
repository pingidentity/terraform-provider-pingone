resource "pingone_davinci_connector_instance" "connectorIPGeolocationio" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPGeolocationio"
  }
  name = "My awesome connectorIPGeolocationio"
  properties = jsonencode({
    "apiKey" = var.connectoripgeolocationio_property_api_key
  })
}
