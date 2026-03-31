resource "pingone_davinci_connector_instance" "socureConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "socureConnector"
  }
  name = "My awesome socureConnector"
  property {
    name  = "accountCreationDate"
    type  = "string"
    value = var.socureconnector_property_account_creation_date
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.socureconnector_property_api_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.socureconnector_property_base_url
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.socureconnector_property_body_data
  }
  property {
    name  = "bodyHeaderText"
    type  = "string"
    value = var.socureconnector_property_body_header_text
  }
  property {
    name  = "city"
    type  = "string"
    value = var.socureconnector_property_city
  }
  property {
    name  = "companyName"
    type  = "string"
    value = var.socureconnector_property_company_name
  }
  property {
    name  = "country"
    type  = "string"
    value = var.socureconnector_property_country
  }
  property {
    name  = "customApiUrl"
    type  = "string"
    value = var.socureconnector_property_custom_api_url
  }
  property {
    name  = "customerUserId"
    type  = "string"
    value = var.socureconnector_property_customer_user_id
  }
  property {
    name  = "deviceInterface"
    type  = "string"
    value = var.socureconnector_property_device_interface
  }
  property {
    name  = "deviceSessionId"
    type  = "string"
    value = var.socureconnector_property_device_session_id
  }
  property {
    name  = "deviceType"
    type  = "string"
    value = var.socureconnector_property_device_type
  }
  property {
    name  = "disbursementType"
    type  = "string"
    value = var.socureconnector_property_disbursement_type
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.socureconnector_property_dob
  }
  property {
    name  = "documentUuid"
    type  = "string"
    value = var.socureconnector_property_document_uuid
  }
  property {
    name  = "email"
    type  = "string"
    value = var.socureconnector_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.socureconnector_property_endpoint
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.socureconnector_property_first_name
  }
  property {
    name  = "geocodes"
    type  = "string"
    value = var.socureconnector_property_geocodes
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.socureconnector_property_headers
  }
  property {
    name  = "lastOrderDate"
    type  = "string"
    value = var.socureconnector_property_last_order_date
  }
  property {
    name  = "method"
    type  = "string"
    value = var.socureconnector_property_method
  }
  property {
    name  = "mobileNumber"
    type  = "string"
    value = var.socureconnector_property_mobile_number
  }
  property {
    name  = "modules"
    type  = "string"
    value = var.socureconnector_property_modules
  }
  property {
    name  = "nationalId"
    type  = "string"
    value = var.socureconnector_property_national_id
  }
  property {
    name  = "operatingSystem"
    type  = "string"
    value = var.socureconnector_property_operating_system
  }
  property {
    name  = "orderAmount"
    type  = "string"
    value = var.socureconnector_property_order_amount
  }
  property {
    name  = "orderChannel"
    type  = "string"
    value = var.socureconnector_property_order_channel
  }
  property {
    name  = "paymentType"
    type  = "string"
    value = var.socureconnector_property_payment_type
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.socureconnector_property_phone_number
  }
  property {
    name  = "physicalAddress"
    type  = "string"
    value = var.socureconnector_property_physical_address
  }
  property {
    name  = "physicalAddress2"
    type  = "string"
    value = var.socureconnector_property_physical_address2
  }
  property {
    name  = "prevOrderCount"
    type  = "string"
    value = var.socureconnector_property_prev_order_count
  }
  property {
    name  = "queryParams"
    type  = "string"
    value = var.socureconnector_property_query_params
  }
  property {
    name  = "recipientCountry"
    type  = "string"
    value = var.socureconnector_property_recipient_country
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.socureconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.socureconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.socureconnector_property_screen2_config
  }
  property {
    name  = "sdkKey"
    type  = "string"
    value = var.socureconnector_property_sdk_key
  }
  property {
    name  = "skWebhookUri"
    type  = "string"
    value = var.socureconnector_property_sk_webhook_uri
  }
  property {
    name  = "state"
    type  = "string"
    value = var.socureconnector_property_state
  }
  property {
    name  = "submissionDate"
    type  = "string"
    value = var.socureconnector_property_submission_date
  }
  property {
    name  = "surName"
    type  = "string"
    value = var.socureconnector_property_sur_name
  }
  property {
    name  = "title"
    type  = "string"
    value = var.socureconnector_property_title
  }
  property {
    name  = "verificationCountry"
    type  = "string"
    value = var.socureconnector_property_verification_country
  }
  property {
    name  = "verificationLevel"
    type  = "string"
    value = var.socureconnector_property_verification_level
  }
  property {
    name  = "watchlistFilters"
    type  = "string"
    value = var.socureconnector_property_watchlist_filters
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.socureconnector_property_zip
  }
}
