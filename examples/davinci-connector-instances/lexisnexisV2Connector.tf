resource "pingone_davinci_connector_instance" "lexisnexisV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "lexisnexisV2Connector"
  }
  name = "My awesome lexisnexisV2Connector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.lexisnexisv2connector_property_api_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.lexisnexisv2connector_property_api_url
  }
  property {
    name  = "customAttribute"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_attribute
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_css
  }
  property {
    name  = "customDeliveryMethod"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_delivery_method
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_html
  }
  property {
    name  = "customProfilingUrl"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_profiling_url
  }
  property {
    name  = "customRequiredParameter"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_required_parameter
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.lexisnexisv2connector_property_custom_script
  }
  property {
    name  = "ddpProfilingUrl"
    type  = "string"
    value = var.lexisnexisv2connector_property_ddp_profiling_url
  }
  property {
    name  = "deliveryMethod"
    type  = "string"
    value = var.lexisnexisv2connector_property_delivery_method
  }
  property {
    name  = "email"
    type  = "string"
    value = var.lexisnexisv2connector_property_email
  }
  property {
    name  = "emailBody"
    type  = "string"
    value = var.lexisnexisv2connector_property_email_body
  }
  property {
    name  = "emailTitle"
    type  = "string"
    value = var.lexisnexisv2connector_property_email_title
  }
  property {
    name  = "events_type"
    type  = "string"
    value = var.lexisnexisv2connector_property_events_type
  }
  property {
    name  = "htmlConfig"
    type  = "string"
    value = var.lexisnexisv2connector_property_html_config
  }
  property {
    name  = "htmlConfig1"
    type  = "string"
    value = var.lexisnexisv2connector_property_html_config1
  }
  property {
    name  = "loadingText"
    type  = "string"
    value = var.lexisnexisv2connector_property_loading_text
  }
  property {
    name  = "orgId"
    type  = "string"
    value = var.lexisnexisv2connector_property_org_id
  }
  property {
    name  = "otpLength"
    type  = "string"
    value = var.lexisnexisv2connector_property_otp_length
  }
  property {
    name  = "otpTimeout"
    type  = "string"
    value = var.lexisnexisv2connector_property_otp_timeout
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.lexisnexisv2connector_property_phone_number
  }
  property {
    name  = "policy"
    type  = "string"
    value = var.lexisnexisv2connector_property_policy
  }
  property {
    name  = "queryType"
    type  = "string"
    value = var.lexisnexisv2connector_property_query_type
  }
  property {
    name  = "requiredParameter"
    type  = "string"
    value = var.lexisnexisv2connector_property_required_parameter
  }
  property {
    name  = "requiredParameterValue"
    type  = "string"
    value = var.lexisnexisv2connector_property_required_parameter_value
  }
  property {
    name  = "review_status"
    type  = "string"
    value = var.lexisnexisv2connector_property_review_status
  }
  property {
    name  = "service_type"
    type  = "string"
    value = var.lexisnexisv2connector_property_service_type
  }
  property {
    name  = "session_id"
    type  = "string"
    value = var.lexisnexisv2connector_property_session_id
  }
  property {
    name  = "smsMessage"
    type  = "string"
    value = var.lexisnexisv2connector_property_sms_message
  }
  property {
    name  = "timeOut"
    type  = "string"
    value = var.lexisnexisv2connector_property_time_out
  }
  property {
    name  = "useCustomApiURL"
    type  = "string"
    value = var.use_custom_api_url
  }
  property {
    name  = "validationRules"
    type  = "string"
    value = var.lexisnexisv2connector_property_validation_rules
  }
}
