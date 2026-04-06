resource "pingone_davinci_connector_instance" "connectIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectIdConnector"
  }
  name = "My awesome connectIdConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectidconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectidconnector_property_button
  }
  property {
    name  = "claims"
    type  = "string"
    value = var.connectidconnector_property_claims
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectidconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectidconnector_property_skip_button_press
  }
  property {
    name  = "wellknownEndpoint"
    type  = "string"
    value = var.connectidconnector_property_wellknown_endpoint
  }
}
