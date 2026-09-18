resource "pingone_davinci_connector_instance" "sentineloneConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sentineloneConnector"
  }
  name = "My awesome sentineloneConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.sentineloneconnector_property_api_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.sentineloneconnector_property_base_url
  }
  property {
    name  = "serialNumber"
    type  = "string"
    value = var.sentineloneconnector_property_serial_number
  }
}
