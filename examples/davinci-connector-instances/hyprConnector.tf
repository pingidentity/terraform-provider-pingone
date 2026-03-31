resource "pingone_davinci_connector_instance" "hyprConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "hyprConnector"
  }
  name = "My awesome hyprConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.hyprconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.hyprconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.hyprconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.hyprconnector_property_skip_button_press
  }
  property {
    name  = "username"
    type  = "string"
    value = var.hyprconnector_property_username
  }
}
