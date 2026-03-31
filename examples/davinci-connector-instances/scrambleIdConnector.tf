resource "pingone_davinci_connector_instance" "scrambleIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "scrambleIdConnector"
  }
  name = "My awesome scrambleIdConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.scrambleidconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.scrambleidconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.scrambleidconnector_property_custom_auth
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.scrambleidconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.scrambleidconnector_property_skip_button_press
  }
}
