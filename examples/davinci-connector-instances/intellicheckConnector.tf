resource "pingone_davinci_connector_instance" "intellicheckConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "intellicheckConnector"
  }
  name = "My awesome intellicheckConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.intellicheckconnector_property_api_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.intellicheckconnector_property_base_url
  }
  property {
    name  = "customerId"
    type  = "string"
    value = var.intellicheckconnector_property_customer_id
  }
  property {
    name  = "deviceValidation"
    type  = "string"
    value = var.intellicheckconnector_property_device_validation
  }
  property {
    name  = "documentTypeSelection"
    type  = "string"
    value = var.intellicheckconnector_property_document_type_selection
  }
  property {
    name  = "intellicheckErrorRedirectUrl"
    type  = "string"
    value = var.intellicheckconnector_property_intellicheck_error_redirect_url
  }
  property {
    name  = "intellicheckRedirectUrl"
    type  = "string"
    value = var.intellicheckconnector_property_intellicheck_redirect_url
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.intellicheckconnector_property_phone_number
  }
  property {
    name  = "signals"
    type  = "string"
    value = var.intellicheckconnector_property_signals
  }
  property {
    name  = "transactionId"
    type  = "string"
    value = var.intellicheckconnector_property_transaction_id
  }
  property {
    name  = "ttl"
    type  = "string"
    value = var.intellicheckconnector_property_ttl
  }
}
