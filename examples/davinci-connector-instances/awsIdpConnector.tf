resource "pingone_davinci_connector_instance" "awsIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "awsIdpConnector"
  }
  name = "My awesome awsIdpConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.awsidpconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.awsidpconnector_property_button
  }
  property {
    name  = "openId"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.awsidpconnector_property_show_powered_by
  }
}
