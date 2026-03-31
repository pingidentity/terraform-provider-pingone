resource "pingone_davinci_connector_instance" "connectorIPGeolocationio" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPGeolocationio"
  }
  name = "My awesome connectorIPGeolocationio"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectoripgeolocationio_property_api_key
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.connectoripgeolocationio_property_ip
  }
  property {
    name  = "lang"
    type  = "string"
    value = var.connectoripgeolocationio_property_lang
  }
}
