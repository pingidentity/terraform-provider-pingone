resource "pingone_davinci_connector_instance" "connectorVericlouds" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorVericlouds"
  }
  name = "My awesome connectorVericlouds"
  property {
    name  = "apiSecret"
    type  = "string"
    value = var.connectorvericlouds_property_api_secret
  }
  property {
    name  = "apikey"
    type  = "string"
    value = var.connectorvericlouds_property_apikey
  }
  property {
    name  = "identifier"
    type  = "string"
    value = var.connectorvericlouds_property_identifier
  }
  property {
    name  = "identifierType"
    type  = "string"
    value = var.connectorvericlouds_property_identifier_type
  }
  property {
    name  = "password"
    type  = "string"
    value = var.connectorvericlouds_property_password
  }
}
