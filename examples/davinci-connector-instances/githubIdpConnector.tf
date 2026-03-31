resource "pingone_davinci_connector_instance" "githubIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "githubIdpConnector"
  }
  name = "My awesome githubIdpConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.githubidpconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.githubidpconnector_property_button
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.githubidpconnector_property_show_powered_by
  }
}
