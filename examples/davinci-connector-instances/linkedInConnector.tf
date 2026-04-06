resource "pingone_davinci_connector_instance" "linkedInConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "linkedInConnector"
  }
  name = "My awesome linkedInConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.linkedinconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.linkedinconnector_property_button
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.linkedinconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.linkedinconnector_property_skip_button_press
  }
}
