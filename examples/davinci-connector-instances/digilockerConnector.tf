resource "pingone_davinci_connector_instance" "digilockerConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "digilockerConnector"
  }
  name = "My awesome digilockerConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.digilockerconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.digilockerconnector_property_button
  }
  property {
    name  = "doctype"
    type  = "string"
    value = var.digilockerconnector_property_doctype
  }
  property {
    name  = "fileType"
    type  = "string"
    value = var.digilockerconnector_property_file_type
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.digilockerconnector_property_show_powered_by
  }
  property {
    name  = "token"
    type  = "string"
    value = var.digilockerconnector_property_token
  }
}
