resource "pingone_davinci_connector_instance" "akamaiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "akamaiConnector"
  }
  name = "My awesome akamaiConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.akamaiconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.akamaiconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.akamaiconnector_property_custom_auth
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.akamaiconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.akamaiconnector_property_skip_button_press
  }
  property {
    name  = "username"
    type  = "string"
    value = var.akamaiconnector_property_username
  }
}
