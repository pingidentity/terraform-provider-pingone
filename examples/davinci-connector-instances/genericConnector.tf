resource "pingone_davinci_connector_instance" "genericConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "genericConnector"
  }
  name = "My awesome genericConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.genericconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.genericconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "password"
    type  = "string"
    value = var.genericconnector_property_password
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.genericconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.genericconnector_property_skip_button_press
  }
  property {
    name  = "username"
    type  = "string"
    value = var.genericconnector_property_username
  }
}
