resource "pingone_davinci_connector_instance" "connectorIsland" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIsland"
  }
  name = "My awesome connectorIsland"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorisland_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorisland_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.connectorisland_property_custom_auth
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorisland_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorisland_property_skip_button_press
  }
}
