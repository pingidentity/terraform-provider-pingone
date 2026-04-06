resource "pingone_davinci_connector_instance" "siftConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "siftConnector"
  }
  name = "My awesome siftConnector"
  property {
    name  = "abuseTypes"
    type  = "string"
    value = var.siftconnector_property_abuse_types
  }
  property {
    name  = "acceptLanguage"
    type  = "string"
    value = var.siftconnector_property_accept_language
  }
  property {
    name  = "accountId"
    type  = "string"
    value = var.siftconnector_property_account_id
  }
  property {
    name  = "address1Billing"
    type  = "string"
    value = var.siftconnector_property_address1_billing
  }
  property {
    name  = "address1Shipping"
    type  = "string"
    value = var.siftconnector_property_address1_shipping
  }
  property {
    name  = "address2Billing"
    type  = "string"
    value = var.siftconnector_property_address2_billing
  }
  property {
    name  = "address2Shipping"
    type  = "string"
    value = var.siftconnector_property_address2_shipping
  }
  property {
    name  = "analyst"
    type  = "string"
    value = var.siftconnector_property_analyst
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.siftconnector_property_api_key
  }
  property {
    name  = "body"
    type  = "string"
    value = var.siftconnector_property_body
  }
  property {
    name  = "changedPassword"
    type  = "string"
    value = var.siftconnector_property_changed_password
  }
  property {
    name  = "cityBilling"
    type  = "string"
    value = var.siftconnector_property_city_billing
  }
  property {
    name  = "cityShipping"
    type  = "string"
    value = var.siftconnector_property_city_shipping
  }
  property {
    name  = "contentLanguage"
    type  = "string"
    value = var.siftconnector_property_content_language
  }
  property {
    name  = "countryBilling"
    type  = "string"
    value = var.siftconnector_property_country_billing
  }
  property {
    name  = "countryShipping"
    type  = "string"
    value = var.siftconnector_property_country_shipping
  }
  property {
    name  = "decisionId"
    type  = "string"
    value = var.siftconnector_property_decision_id
  }
  property {
    name  = "decisionType"
    type  = "string"
    value = var.siftconnector_property_decision_type
  }
  property {
    name  = "description"
    type  = "string"
    value = var.siftconnector_property_description
  }
  property {
    name  = "email"
    type  = "string"
    value = var.siftconnector_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.siftconnector_property_endpoint
  }
  property {
    name  = "failureReason"
    type  = "string"
    value = var.siftconnector_property_failure_reason
  }
  property {
    name  = "fullName"
    type  = "string"
    value = var.siftconnector_property_full_name
  }
  property {
    name  = "fullNameBilling"
    type  = "string"
    value = var.siftconnector_property_full_name_billing
  }
  property {
    name  = "fullNameShipping"
    type  = "string"
    value = var.siftconnector_property_full_name_shipping
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.siftconnector_property_headers
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.siftconnector_property_ip_address
  }
  property {
    name  = "loginStatus"
    type  = "string"
    value = var.siftconnector_property_login_status
  }
  property {
    name  = "method"
    type  = "string"
    value = var.siftconnector_property_method
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.siftconnector_property_phone_number
  }
  property {
    name  = "phoneNumberBilling"
    type  = "string"
    value = var.siftconnector_property_phone_number_billing
  }
  property {
    name  = "phoneNumberPrimary"
    type  = "string"
    value = var.siftconnector_property_phone_number_primary
  }
  property {
    name  = "phoneNumberShipping"
    type  = "string"
    value = var.siftconnector_property_phone_number_shipping
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.siftconnector_property_query_parameters
  }
  property {
    name  = "reason"
    type  = "string"
    value = var.siftconnector_property_reason
  }
  property {
    name  = "regionBilling"
    type  = "string"
    value = var.siftconnector_property_region_billing
  }
  property {
    name  = "regionShipping"
    type  = "string"
    value = var.siftconnector_property_region_shipping
  }
  property {
    name  = "returnScore"
    type  = "string"
    value = var.siftconnector_property_return_score
  }
  property {
    name  = "scorePercentiles"
    type  = "string"
    value = var.siftconnector_property_score_percentiles
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.siftconnector_property_session_id
  }
  property {
    name  = "sessionIdDecision"
    type  = "string"
    value = var.siftconnector_property_session_id_decision
  }
  property {
    name  = "sessionIdOptional"
    type  = "string"
    value = var.siftconnector_property_session_id_optional
  }
  property {
    name  = "source"
    type  = "string"
    value = var.siftconnector_property_source
  }
  property {
    name  = "status"
    type  = "string"
    value = var.siftconnector_property_status
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.siftconnector_property_user_agent
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.siftconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.siftconnector_property_username
  }
  property {
    name  = "verificationType"
    type  = "string"
    value = var.siftconnector_property_verification_type
  }
  property {
    name  = "verifiedReason"
    type  = "string"
    value = var.siftconnector_property_verified_reason
  }
  property {
    name  = "verifiedValue"
    type  = "string"
    value = var.siftconnector_property_verified_value
  }
  property {
    name  = "workflowStatus"
    type  = "string"
    value = var.siftconnector_property_workflow_status
  }
  property {
    name  = "zipcodeBilling"
    type  = "string"
    value = var.siftconnector_property_zipcode_billing
  }
  property {
    name  = "zipcodeShipping"
    type  = "string"
    value = var.siftconnector_property_zipcode_shipping
  }
}
