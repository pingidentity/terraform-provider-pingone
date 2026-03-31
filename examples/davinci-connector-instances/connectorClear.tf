resource "pingone_davinci_connector_instance" "connectorClear" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClear"
  }
  name = "My awesome connectorClear"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorclear_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorclear_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.connectorclear_property_custom_auth
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorclear_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorclear_property_skip_button_press
  }
  property {
    name  = "url"
    type  = "string"
    value = var.connectorclear_property_url
  }
}
