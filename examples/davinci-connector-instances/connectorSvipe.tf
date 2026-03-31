resource "pingone_davinci_connector_instance" "connectorSvipe" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSvipe"
  }
  name = "My awesome connectorSvipe"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorsvipe_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorsvipe_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorsvipe_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorsvipe_property_skip_button_press
  }
}
