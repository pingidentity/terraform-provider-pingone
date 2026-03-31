resource "pingone_davinci_connector_instance" "twitterIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "twitterIdpConnector"
  }
  name = "My awesome twitterIdpConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.twitteridpconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.twitteridpconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.twitteridpconnector_property_show_powered_by
  }
}
