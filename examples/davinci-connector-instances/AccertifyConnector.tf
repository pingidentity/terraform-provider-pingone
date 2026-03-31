resource "pingone_davinci_connector_instance" "AccertifyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "AccertifyConnector"
  }
  name = "My awesome AccertifyConnector"
  property {
    name  = "accountCreateFailureCode"
    type  = "string"
    value = var.accertifyconnector_property_account_create_failure_code
  }
  property {
    name  = "accountID"
    type  = "string"
    value = var.accertifyconnector_property_account_id
  }
  property {
    name  = "authenticationMethod"
    type  = "string"
    value = var.accertifyconnector_property_authentication_method
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.accertifyconnector_property_base_url
  }
  property {
    name  = "basicAuthPassword"
    type  = "string"
    value = var.accertifyconnector_property_basic_auth_password
  }
  property {
    name  = "basicAuthUsername"
    type  = "string"
    value = var.accertifyconnector_property_basic_auth_username
  }
  property {
    name  = "devicePayload"
    type  = "string"
    value = var.accertifyconnector_property_device_payload
  }
  property {
    name  = "deviceTransactionID"
    type  = "string"
    value = var.accertifyconnector_property_device_transaction_id
  }
  property {
    name  = "emailAddress"
    type  = "string"
    value = var.accertifyconnector_property_email_address
  }
  property {
    name  = "eventSource"
    type  = "string"
    value = var.accertifyconnector_property_event_source
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.accertifyconnector_property_first_name
  }
  property {
    name  = "hashedPassword"
    type  = "string"
    value = var.accertifyconnector_property_hashed_password
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.accertifyconnector_property_ip_address
  }
  property {
    name  = "isMobileAppInstalled"
    type  = "string"
    value = var.accertifyconnector_property_is_mobile_app_installed
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.accertifyconnector_property_last_name
  }
  property {
    name  = "loginEventFailureCode"
    type  = "string"
    value = var.accertifyconnector_property_login_event_failure_code
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.accertifyconnector_property_middle_name
  }
  property {
    name  = "pageID"
    type  = "string"
    value = var.accertifyconnector_property_page_id
  }
  property {
    name  = "previousAvsResult"
    type  = "string"
    value = var.accertifyconnector_property_previous_avs_result
  }
  property {
    name  = "previousCardBin"
    type  = "string"
    value = var.accertifyconnector_property_previous_card_bin
  }
  property {
    name  = "previousCardLastFour"
    type  = "string"
    value = var.accertifyconnector_property_previous_card_last_four
  }
  property {
    name  = "previousCardNumber"
    type  = "string"
    value = var.accertifyconnector_property_previous_card_number
  }
  property {
    name  = "previousCvvResult"
    type  = "string"
    value = var.accertifyconnector_property_previous_cvv_result
  }
  property {
    name  = "previousEmailAddress"
    type  = "string"
    value = var.accertifyconnector_property_previous_email_address
  }
  property {
    name  = "previousExpirationDay"
    type  = "string"
    value = var.accertifyconnector_property_previous_expiration_day
  }
  property {
    name  = "previousExpirationMonth"
    type  = "string"
    value = var.accertifyconnector_property_previous_expiration_month
  }
  property {
    name  = "previousExpirationYear"
    type  = "string"
    value = var.accertifyconnector_property_previous_expiration_year
  }
  property {
    name  = "previousNameOnCreditCard"
    type  = "string"
    value = var.accertifyconnector_property_previous_name_on_credit_card
  }
  property {
    name  = "previousPaymentType"
    type  = "string"
    value = var.accertifyconnector_property_previous_payment_type
  }
  property {
    name  = "success"
    type  = "string"
    value = var.accertifyconnector_property_success
  }
  property {
    name  = "ubaEvents"
    type  = "string"
    value = var.accertifyconnector_property_uba_events
  }
  property {
    name  = "ubaID"
    type  = "string"
    value = var.accertifyconnector_property_uba_id
  }
  property {
    name  = "ubaSessionID"
    type  = "string"
    value = var.accertifyconnector_property_uba_session_id
  }
  property {
    name  = "updateEventID"
    type  = "string"
    value = var.accertifyconnector_property_update_event_id
  }
  property {
    name  = "updateEventRecommendationDetail"
    type  = "string"
    value = var.accertifyconnector_property_update_event_recommendation_detail
  }
  property {
    name  = "updateEventType"
    type  = "string"
    value = var.accertifyconnector_property_update_event_type
  }
  property {
    name  = "updateTrigger"
    type  = "string"
    value = var.accertifyconnector_property_update_trigger
  }
  property {
    name  = "updatedAvsResult"
    type  = "string"
    value = var.accertifyconnector_property_updated_avs_result
  }
  property {
    name  = "updatedCardBin"
    type  = "string"
    value = var.accertifyconnector_property_updated_card_bin
  }
  property {
    name  = "updatedCardLastFour"
    type  = "string"
    value = var.accertifyconnector_property_updated_card_last_four
  }
  property {
    name  = "updatedCardNumber"
    type  = "string"
    value = var.accertifyconnector_property_updated_card_number
  }
  property {
    name  = "updatedCvvResult"
    type  = "string"
    value = var.accertifyconnector_property_updated_cvv_result
  }
  property {
    name  = "updatedEmailAddress"
    type  = "string"
    value = var.accertifyconnector_property_updated_email_address
  }
  property {
    name  = "updatedExpirationDay"
    type  = "string"
    value = var.accertifyconnector_property_updated_expiration_day
  }
  property {
    name  = "updatedExpirationMonth"
    type  = "string"
    value = var.accertifyconnector_property_updated_expiration_month
  }
  property {
    name  = "updatedExpirationYear"
    type  = "string"
    value = var.accertifyconnector_property_updated_expiration_year
  }
  property {
    name  = "updatedNameOnCreditCard"
    type  = "string"
    value = var.accertifyconnector_property_updated_name_on_credit_card
  }
  property {
    name  = "updatedPaymentType"
    type  = "string"
    value = var.accertifyconnector_property_updated_payment_type
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.accertifyconnector_property_user_agent
  }
  property {
    name  = "username"
    type  = "string"
    value = var.accertifyconnector_property_username
  }
  property {
    name  = "verificationAttempts"
    type  = "string"
    value = var.accertifyconnector_property_verification_attempts
  }
  property {
    name  = "verificationFailureCode"
    type  = "string"
    value = var.accertifyconnector_property_verification_failure_code
  }
  property {
    name  = "verificationStatus"
    type  = "string"
    value = var.accertifyconnector_property_verification_status
  }
  property {
    name  = "verificationType"
    type  = "string"
    value = var.accertifyconnector_property_verification_type
  }
}
