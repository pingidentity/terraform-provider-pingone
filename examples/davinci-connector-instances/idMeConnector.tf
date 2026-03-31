resource "pingone_davinci_connector_instance" "idMeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idMeConnector"
  }
  name = "My awesome idMeConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.idmeconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.idmeconnector_property_button
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.idmeconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.idmeconnector_property_skip_button_press
  }
}
