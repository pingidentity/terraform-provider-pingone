resource "pingone_davinci_connector_instance" "connectorOpswat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorOpswat"
  }
  name = "My awesome connectorOpswat"
  property {
    name  = "clientID"
    type  = "string"
    value = var.connectoropswat_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectoropswat_property_client_secret
  }
  property {
    name  = "crossDomainApiPort"
    type  = "string"
    value = var.connectoropswat_property_cross_domain_api_port
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.connectoropswat_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.connectoropswat_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.connectoropswat_property_custom_script
  }
  property {
    name  = "deviceId"
    type  = "string"
    value = var.connectoropswat_property_device_id
  }
  property {
    name  = "maDomain"
    type  = "string"
    value = var.connectoropswat_property_ma_domain
  }
}
