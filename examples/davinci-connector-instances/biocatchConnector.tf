resource "pingone_davinci_connector_instance" "biocatchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "biocatchConnector"
  }
  name = "My awesome biocatchConnector"
  property {
    name  = "accountBalance"
    type  = "string"
    value = var.biocatchconnector_property_account_balance
  }
  property {
    name  = "accountID"
    type  = "string"
    value = var.biocatchconnector_property_account_id
  }
  property {
    name  = "accountOpenDate"
    type  = "string"
    value = var.biocatchconnector_property_account_open_date
  }
  property {
    name  = "activityAmount"
    type  = "string"
    value = var.biocatchconnector_property_activity_amount
  }
  property {
    name  = "activityAmountTotal"
    type  = "string"
    value = var.biocatchconnector_property_activity_amount_total
  }
  property {
    name  = "activityCategory"
    type  = "string"
    value = var.biocatchconnector_property_activity_category
  }
  property {
    name  = "activityName"
    type  = "string"
    value = var.biocatchconnector_property_activity_name
  }
  property {
    name  = "activityType"
    type  = "string"
    value = var.biocatchconnector_property_activity_type
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.biocatchconnector_property_api_url
  }
  property {
    name  = "authMethodUsed"
    type  = "string"
    value = var.biocatchconnector_property_auth_method_used
  }
  property {
    name  = "authenticationResult"
    type  = "string"
    value = var.biocatchconnector_property_authentication_result
  }
  property {
    name  = "biocatchAction"
    type  = "string"
    value = var.biocatchconnector_property_biocatch_action
  }
  property {
    name  = "biocatchFraudType"
    type  = "string"
    value = var.biocatchconnector_property_biocatch_fraud_type
  }
  property {
    name  = "biocatchSessionResult"
    type  = "string"
    value = var.biocatchconnector_property_biocatch_session_result
  }
  property {
    name  = "biocatchSessionStatus"
    type  = "string"
    value = var.biocatchconnector_property_biocatch_session_status
  }
  property {
    name  = "biometricLogin"
    type  = "string"
    value = var.biocatchconnector_property_biometric_login
  }
  property {
    name  = "brand"
    type  = "string"
    value = var.biocatchconnector_property_brand
  }
  property {
    name  = "cardType"
    type  = "string"
    value = var.biocatchconnector_property_card_type
  }
  property {
    name  = "cellularAreaNumbers"
    type  = "string"
    value = var.biocatchconnector_property_cellular_area_numbers
  }
  property {
    name  = "comment"
    type  = "string"
    value = var.biocatchconnector_property_comment
  }
  property {
    name  = "contextName"
    type  = "string"
    value = var.biocatchconnector_property_context_name
  }
  property {
    name  = "customFraudType"
    type  = "string"
    value = var.biocatchconnector_property_custom_fraud_type
  }
  property {
    name  = "customerId"
    type  = "string"
    value = var.biocatchconnector_property_customer_id
  }
  property {
    name  = "customerPassword"
    type  = "string"
    value = var.biocatchconnector_property_customer_password
  }
  property {
    name  = "customerSessionId"
    type  = "string"
    value = var.biocatchconnector_property_customer_session_id
  }
  property {
    name  = "dateOfCreation"
    type  = "string"
    value = var.biocatchconnector_property_date_of_creation
  }
  property {
    name  = "deviceID"
    type  = "string"
    value = var.biocatchconnector_property_device_id
  }
  property {
    name  = "deviceIpStatus"
    type  = "string"
    value = var.biocatchconnector_property_device_ip_status
  }
  property {
    name  = "deviceModel"
    type  = "string"
    value = var.biocatchconnector_property_device_model
  }
  property {
    name  = "income"
    type  = "string"
    value = var.biocatchconnector_property_income
  }
  property {
    name  = "invoiceNumber"
    type  = "string"
    value = var.biocatchconnector_property_invoice_number
  }
  property {
    name  = "isLoginSuccess"
    type  = "string"
    value = var.biocatchconnector_property_is_login_success
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.biocatchconnector_property_javascript_cdn_url
  }
  property {
    name  = "jobDescription"
    type  = "string"
    value = var.biocatchconnector_property_job_description
  }
  property {
    name  = "language"
    type  = "string"
    value = var.biocatchconnector_property_language
  }
  property {
    name  = "lastWithdrawalAmount"
    type  = "string"
    value = var.biocatchconnector_property_last_withdrawal_amount
  }
  property {
    name  = "lastWithdrawalDate"
    type  = "string"
    value = var.biocatchconnector_property_last_withdrawal_date
  }
  property {
    name  = "loadingText"
    type  = "string"
    value = var.biocatchconnector_property_loading_text
  }
  property {
    name  = "maritalStatus"
    type  = "string"
    value = var.biocatchconnector_property_marital_status
  }
  property {
    name  = "market"
    type  = "string"
    value = var.biocatchconnector_property_market
  }
  property {
    name  = "membershipID"
    type  = "string"
    value = var.biocatchconnector_property_membership_id
  }
  property {
    name  = "numOfFailedAuth"
    type  = "string"
    value = var.biocatchconnector_property_num_of_failed_auth
  }
  property {
    name  = "onlineAccountOpenDate"
    type  = "string"
    value = var.biocatchconnector_property_online_account_open_date
  }
  property {
    name  = "passwordChangeDate"
    type  = "string"
    value = var.biocatchconnector_property_password_change_date
  }
  property {
    name  = "payeeAccountType"
    type  = "string"
    value = var.biocatchconnector_property_payee_account_type
  }
  property {
    name  = "payeeAge"
    type  = "string"
    value = var.biocatchconnector_property_payee_age
  }
  property {
    name  = "payeeValue"
    type  = "string"
    value = var.biocatchconnector_property_payee_value
  }
  property {
    name  = "payerAccountType"
    type  = "string"
    value = var.biocatchconnector_property_payer_account_type
  }
  property {
    name  = "payerValue"
    type  = "string"
    value = var.biocatchconnector_property_payer_value
  }
  property {
    name  = "paymentMethod"
    type  = "string"
    value = var.biocatchconnector_property_payment_method
  }
  property {
    name  = "permissionsRequested"
    type  = "string"
    value = var.biocatchconnector_property_permissions_requested
  }
  property {
    name  = "platformType"
    type  = "string"
    value = var.biocatchconnector_property_platform_type
  }
  property {
    name  = "postcode"
    type  = "string"
    value = var.biocatchconnector_property_postcode
  }
  property {
    name  = "product"
    type  = "string"
    value = var.biocatchconnector_property_product
  }
  property {
    name  = "recurrentPayee"
    type  = "string"
    value = var.biocatchconnector_property_recurrent_payee
  }
  property {
    name  = "sdkToken"
    type  = "string"
    value = var.biocatchconnector_property_sdk_token
  }
  property {
    name  = "serviceCategory"
    type  = "string"
    value = var.biocatchconnector_property_service_category
  }
  property {
    name  = "sessionType"
    type  = "string"
    value = var.biocatchconnector_property_session_type
  }
  property {
    name  = "shipmentCountry"
    type  = "string"
    value = var.biocatchconnector_property_shipment_country
  }
  property {
    name  = "shipmentType"
    type  = "string"
    value = var.biocatchconnector_property_shipment_type
  }
  property {
    name  = "solution"
    type  = "string"
    value = var.biocatchconnector_property_solution
  }
  property {
    name  = "thirdPartyFlow"
    type  = "string"
    value = var.biocatchconnector_property_third_party_flow
  }
  property {
    name  = "timeOut"
    type  = "string"
    value = var.biocatchconnector_property_time_out
  }
  property {
    name  = "timestamp"
    type  = "string"
    value = var.biocatchconnector_property_timestamp
  }
  property {
    name  = "transactionID"
    type  = "string"
    value = var.biocatchconnector_property_transaction_id
  }
  property {
    name  = "truthApiKey"
    type  = "string"
    value = var.biocatchconnector_property_truth_api_key
  }
  property {
    name  = "truthApiUrl"
    type  = "string"
    value = var.biocatchconnector_property_truth_api_url
  }
  property {
    name  = "userType"
    type  = "string"
    value = var.biocatchconnector_property_user_type
  }
  property {
    name  = "uuid"
    type  = "string"
    value = var.biocatchconnector_property_uuid
  }
  property {
    name  = "web_journey"
    type  = "string"
    value = var.biocatchconnector_property_web_journey
  }
  property {
    name  = "yearOfBirth"
    type  = "string"
    value = var.biocatchconnector_property_year_of_birth
  }
  property {
    name  = "yob"
    type  = "string"
    value = var.biocatchconnector_property_yob
  }
}
