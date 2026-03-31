resource "pingone_davinci_connector_instance" "connectorSpycloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpycloud"
  }
  name = "My awesome connectorSpycloud"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorspycloud_property_api_key
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorspycloud_property_email
  }
  property {
    name  = "inboundPassword"
    type  = "string"
    value = var.connectorspycloud_property_inbound_password
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectorspycloud_property_username
  }
}
