resource "pingone_davinci_connector_instance" "bitbucketIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "bitbucketIdpConnector"
  }
  name = "My awesome bitbucketIdpConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.bitbucketidpconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.bitbucketidpconnector_property_button
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.bitbucketidpconnector_property_show_powered_by
  }
}
