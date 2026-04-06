resource "pingone_davinci_connector_instance" "connectorSignicat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSignicat"
  }
  name = "My awesome connectorSignicat"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorsignicat_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorsignicat_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorsignicat_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorsignicat_property_skip_button_press
  }
}
