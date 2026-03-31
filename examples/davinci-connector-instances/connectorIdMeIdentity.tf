resource "pingone_davinci_connector_instance" "connectorIdMeIdentity" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdMeIdentity"
  }
  name = "My awesome connectorIdMeIdentity"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectoridmeidentity_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectoridmeidentity_property_button
  }
  property {
    name  = "openId"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectoridmeidentity_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectoridmeidentity_property_skip_button_press
  }
}
