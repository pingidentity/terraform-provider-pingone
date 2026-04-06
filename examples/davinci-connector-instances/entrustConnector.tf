resource "pingone_davinci_connector_instance" "entrustConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "entrustConnector"
  }
  name = "My awesome entrustConnector"
  property {
    name  = "applicationId"
    type  = "string"
    value = var.entrustconnector_property_application_id
  }
  property {
    name  = "htmlConfig0"
    type  = "string"
    value = var.entrustconnector_property_html_config0
  }
  property {
    name  = "htmlConfig1"
    type  = "string"
    value = var.entrustconnector_property_html_config1
  }
  property {
    name  = "htmlConfig2"
    type  = "string"
    value = var.entrustconnector_property_html_config2
  }
  property {
    name  = "serviceDomain"
    type  = "string"
    value = var.entrustconnector_property_service_domain
  }
  property {
    name  = "transactionDetails"
    type  = "string"
    value = var.entrustconnector_property_transaction_details
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.entrustconnector_property_user_id
  }
}
