resource "pingone_davinci_connector_instance" "forterConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "forterConnector"
  }
  name = "My awesome forterConnector"
  property {
    name  = "accessRequestType"
    type  = "string"
    value = var.forterconnector_property_access_request_type
  }
  property {
    name  = "accountId"
    type  = "string"
    value = var.forterconnector_property_account_id
  }
  property {
    name  = "accountIdSignup"
    type  = "string"
    value = var.forterconnector_property_account_id_signup
  }
  property {
    name  = "accountType"
    type  = "string"
    value = var.forterconnector_property_account_type
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.forterconnector_property_api_version
  }
  property {
    name  = "body"
    type  = "string"
    value = var.forterconnector_property_body
  }
  property {
    name  = "correlationId"
    type  = "string"
    value = var.forterconnector_property_correlation_id
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.forterconnector_property_endpoint
  }
  property {
    name  = "eventTimestamp"
    type  = "string"
    value = var.forterconnector_property_event_timestamp
  }
  property {
    name  = "forterTokenCookie"
    type  = "string"
    value = var.forterconnector_property_forter_token_cookie
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.forterconnector_property_headers
  }
  property {
    name  = "inputType"
    type  = "string"
    value = var.forterconnector_property_input_type
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.forterconnector_property_ip_address
  }
  property {
    name  = "loginMethodType"
    type  = "string"
    value = var.forterconnector_property_login_method_type
  }
  property {
    name  = "loginStatus"
    type  = "string"
    value = var.forterconnector_property_login_status
  }
  property {
    name  = "method"
    type  = "string"
    value = var.forterconnector_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.forterconnector_property_query_parameters
  }
  property {
    name  = "secretKey"
    type  = "string"
    value = var.forterconnector_property_secret_key
  }
  property {
    name  = "siteId"
    type  = "string"
    value = var.forterconnector_property_site_id
  }
  property {
    name  = "status"
    type  = "string"
    value = var.forterconnector_property_status
  }
  property {
    name  = "statusChangeBy"
    type  = "string"
    value = var.forterconnector_property_status_change_by
  }
  property {
    name  = "statusChangeReason"
    type  = "string"
    value = var.forterconnector_property_status_change_reason
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.forterconnector_property_user_agent
  }
  property {
    name  = "verificationOutcome"
    type  = "string"
    value = var.forterconnector_property_verification_outcome
  }
}
