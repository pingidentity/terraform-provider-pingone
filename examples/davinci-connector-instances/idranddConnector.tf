resource "pingone_davinci_connector_instance" "idranddConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idranddConnector"
  }
  name = "My awesome idranddConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.idranddconnector_property_api_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.idranddconnector_property_api_url
  }
  property {
    name  = "bodyHeaderText"
    type  = "string"
    value = var.idranddconnector_property_body_header_text
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.idranddconnector_property_next_button_text
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.idranddconnector_property_next_event
  }
  property {
    name  = "title"
    type  = "string"
    value = var.idranddconnector_property_title
  }
}
