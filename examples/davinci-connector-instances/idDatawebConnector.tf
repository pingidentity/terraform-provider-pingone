resource "pingone_davinci_connector_instance" "idDatawebConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idDatawebConnector"
  }
  name = "My awesome idDatawebConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.iddatawebconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.iddatawebconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "piiParams"
    type  = "string"
    value = var.iddatawebconnector_property_pii_params
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.iddatawebconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.iddatawebconnector_property_skip_button_press
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.iddatawebconnector_property_subject
  }
}
