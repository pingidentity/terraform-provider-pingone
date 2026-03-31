resource "pingone_davinci_connector_instance" "connectorAcuant" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAcuant"
  }
  name = "My awesome connectorAcuant"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectoracuant_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectoracuant_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "prePopulatedFields"
    type  = "string"
    value = var.connectoracuant_property_pre_populated_fields
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectoracuant_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectoracuant_property_skip_button_press
  }
  property {
    name  = "url"
    type  = "string"
    value = var.connectoracuant_property_url
  }
}
