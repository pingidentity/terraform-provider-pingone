resource "pingone_davinci_connector_instance" "connectorCircleAccess" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCircleAccess"
  }
  name = "My awesome connectorCircleAccess"
  property {
    name  = "appKey"
    type  = "string"
    value = var.connectorcircleaccess_property_app_key
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorcircleaccess_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorcircleaccess_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "loginUrl"
    type  = "string"
    value = var.connectorcircleaccess_property_login_url
  }
  property {
    name  = "readKey"
    type  = "string"
    value = var.connectorcircleaccess_property_read_key
  }
  property {
    name  = "returnToUrl"
    type  = "string"
    value = var.connectorcircleaccess_property_return_to_url
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorcircleaccess_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorcircleaccess_property_skip_button_press
  }
  property {
    name  = "writeKey"
    type  = "string"
    value = var.connectorcircleaccess_property_write_key
  }
}
