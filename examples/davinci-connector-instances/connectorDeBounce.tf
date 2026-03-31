resource "pingone_davinci_connector_instance" "connectorDeBounce" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorDeBounce"
  }
  name = "My awesome connectorDeBounce"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectordebounce_property_api_key
  }
  property {
    name  = "append"
    type  = "string"
    value = var.connectordebounce_property_append
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectordebounce_property_email
  }
  property {
    name  = "endDate"
    type  = "string"
    value = var.connectordebounce_property_end_date
  }
  property {
    name  = "gsuite"
    type  = "string"
    value = var.connectordebounce_property_gsuite
  }
  property {
    name  = "listID"
    type  = "string"
    value = var.connectordebounce_property_list_id
  }
  property {
    name  = "listURL"
    type  = "string"
    value = var.connectordebounce_property_list_url
  }
  property {
    name  = "photo"
    type  = "string"
    value = var.connectordebounce_property_photo
  }
  property {
    name  = "startDate"
    type  = "string"
    value = var.connectordebounce_property_start_date
  }
}
