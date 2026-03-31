resource "pingone_davinci_connector_instance" "seonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "seonConnector"
  }
  name = "My awesome seonConnector"
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "body"
    type  = "string"
    value = var.seonconnector_property_body
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.seonconnector_property_endpoint
  }
  property {
    name  = "fraudScoreConfigType"
    type  = "string"
    value = var.seonconnector_property_fraud_score_config_type
  }
  property {
    name  = "fraudScoreJsonConfig"
    type  = "string"
    value = var.seonconnector_property_fraud_score_json_config
  }
  property {
    name  = "fraudScoreMultiSelectConfig"
    type  = "string"
    value = var.seonconnector_property_fraud_score_multi_select_config
  }
  property {
    name  = "fraudScoreMultiSelectResponseFields"
    type  = "string"
    value = var.seonconnector_property_fraud_score_multi_select_response_fields
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.seonconnector_property_headers
  }
  property {
    name  = "licenseKey"
    type  = "string"
    value = var.seonconnector_property_license_key
  }
  property {
    name  = "method"
    type  = "string"
    value = var.seonconnector_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.seonconnector_property_query_parameters
  }
  property {
    name  = "sendFeedbackTransactions"
    type  = "string"
    value = var.seonconnector_property_send_feedback_transactions
  }
}
