resource "pingone_davinci_connector_instance" "connectorIPregistry" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPregistry"
  }
  name = "My awesome connectorIPregistry"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectoripregistry_property_api_key
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.connectoripregistry_property_ip_address
  }
}
