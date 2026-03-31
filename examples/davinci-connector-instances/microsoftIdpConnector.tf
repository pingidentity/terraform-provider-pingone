resource "pingone_davinci_connector_instance" "microsoftIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftIdpConnector"
  }
  name = "My awesome microsoftIdpConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.microsoftidpconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.microsoftidpconnector_property_button
  }
  property {
    name  = "openId"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.microsoftidpconnector_property_show_powered_by
  }
}
