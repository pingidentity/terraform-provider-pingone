resource "pingone_davinci_connector_instance" "connectorBTps" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTps"
  }
  name = "My awesome connectorBTps"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorbtps_property_api_key
  }
  property {
    name  = "apiUser"
    type  = "string"
    value = var.connectorbtps_property_api_user
  }
  property {
    name  = "domain"
    type  = "string"
    value = var.connectorbtps_property_domain
  }
  property {
    name  = "hostname"
    type  = "string"
    value = var.connectorbtps_property_hostname
  }
  property {
    name  = "reason"
    type  = "string"
    value = var.connectorbtps_property_reason
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectorbtps_property_username
  }
}
