resource "pingone_davinci_connector_instance" "connector1Kosmos" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector1Kosmos"
  }
  name = "My awesome connector1Kosmos"
  property {
    name  = "authType"
    type  = "string"
    value = var.connector1kosmos_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connector1kosmos_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connector1kosmos_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connector1kosmos_property_skip_button_press
  }
}
