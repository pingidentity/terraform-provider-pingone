resource "pingone_davinci_connector_instance" "analyticsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "analyticsConnector"
  }
  name = "My awesome analyticsConnector"
  property {
    name  = "customTimestamp"
    type  = "string"
    value = var.analyticsconnector_property_custom_timestamp
  }
  property {
    name  = "outcomeDescription"
    type  = "string"
    value = var.analyticsconnector_property_outcome_description
  }
  property {
    name  = "outcomeStatus"
    type  = "string"
    value = var.analyticsconnector_property_outcome_status
  }
  property {
    name  = "outcomeStatusDetail"
    type  = "string"
    value = var.analyticsconnector_property_outcome_status_detail
  }
  property {
    name  = "outcomeType"
    type  = "string"
    value = var.analyticsconnector_property_outcome_type
  }
  property {
    name  = "shouldContinueOnError"
    type  = "string"
    value = var.analyticsconnector_property_should_continue_on_error
  }
}
