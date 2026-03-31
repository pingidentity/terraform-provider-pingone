resource "pingone_davinci_connector_instance" "googleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleConnector"
  }
  name = "My awesome googleConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.googleconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.googleconnector_property_button
  }
  property {
    name  = "openId"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.googleconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.googleconnector_property_skip_button_press
  }
}
