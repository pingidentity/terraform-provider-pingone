resource "pingone_davinci_connector_instance" "errorConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "errorConnector"
  }
  name = "My awesome errorConnector"
  property {
    name  = "errorCallbackSuppress"
    type  = "string"
    value = var.errorconnector_property_error_callback_suppress
  }
  property {
    name  = "errorCode"
    type  = "string"
    value = var.errorconnector_property_error_code
  }
  property {
    name  = "errorDescription"
    type  = "string"
    value = var.errorconnector_property_error_description
  }
  property {
    name  = "errorMessage"
    type  = "string"
    value = var.errorconnector_property_error_message
  }
  property {
    name  = "errorReason"
    type  = "string"
    value = var.errorconnector_property_error_reason
  }
}
