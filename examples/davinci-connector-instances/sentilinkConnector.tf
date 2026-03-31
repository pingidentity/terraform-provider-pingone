resource "pingone_davinci_connector_instance" "sentilinkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sentilinkConnector"
  }
  name = "My awesome sentilinkConnector"
  property {
    name  = "account"
    type  = "string"
    value = var.sentilinkconnector_property_account
  }
  property {
    name  = "addressLine1"
    type  = "string"
    value = var.sentilinkconnector_property_address_line1
  }
  property {
    name  = "addressLine2"
    type  = "string"
    value = var.sentilinkconnector_property_address_line2
  }
  property {
    name  = "altAddressLine1"
    type  = "string"
    value = var.sentilinkconnector_property_alt_address_line1
  }
  property {
    name  = "altAddressLine2"
    type  = "string"
    value = var.sentilinkconnector_property_alt_address_line2
  }
  property {
    name  = "altCity"
    type  = "string"
    value = var.sentilinkconnector_property_alt_city
  }
  property {
    name  = "altCountryCode"
    type  = "string"
    value = var.sentilinkconnector_property_alt_country_code
  }
  property {
    name  = "altEmail"
    type  = "string"
    value = var.sentilinkconnector_property_alt_email
  }
  property {
    name  = "altPhone"
    type  = "string"
    value = var.sentilinkconnector_property_alt_phone
  }
  property {
    name  = "altState"
    type  = "string"
    value = var.sentilinkconnector_property_alt_state
  }
  property {
    name  = "altZip"
    type  = "string"
    value = var.sentilinkconnector_property_alt_zip
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.sentilinkconnector_property_api_url
  }
  property {
    name  = "applicationCreated"
    type  = "string"
    value = var.sentilinkconnector_property_application_created
  }
  property {
    name  = "applicationId"
    type  = "string"
    value = var.sentilinkconnector_property_application_id
  }
  property {
    name  = "businessName"
    type  = "string"
    value = var.sentilinkconnector_property_business_name
  }
  property {
    name  = "businessUrl"
    type  = "string"
    value = var.sentilinkconnector_property_business_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.sentilinkconnector_property_city
  }
  property {
    name  = "consentValidDays"
    type  = "string"
    value = var.sentilinkconnector_property_consent_valid_days
  }
  property {
    name  = "countryCode"
    type  = "string"
    value = var.sentilinkconnector_property_country_code
  }
  property {
    name  = "deviceId"
    type  = "string"
    value = var.sentilinkconnector_property_device_id
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.sentilinkconnector_property_dob
  }
  property {
    name  = "ein"
    type  = "string"
    value = var.sentilinkconnector_property_ein
  }
  property {
    name  = "email"
    type  = "string"
    value = var.sentilinkconnector_property_email
  }
  property {
    name  = "extraDataClusters"
    type  = "string"
    value = var.sentilinkconnector_property_extra_data_clusters
  }
  property {
    name  = "extraDataManifest"
    type  = "string"
    value = var.sentilinkconnector_property_extra_data_manifest
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.sentilinkconnector_property_first_name
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.sentilinkconnector_property_ip
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.sentilinkconnector_property_javascript_cdn_url
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.sentilinkconnector_property_last_name
  }
  property {
    name  = "leadType"
    type  = "string"
    value = var.sentilinkconnector_property_lead_type
  }
  property {
    name  = "loanAmount"
    type  = "string"
    value = var.sentilinkconnector_property_loan_amount
  }
  property {
    name  = "loanCurrency"
    type  = "string"
    value = var.sentilinkconnector_property_loan_currency
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.sentilinkconnector_property_phone
  }
  property {
    name  = "scores"
    type  = "string"
    value = var.sentilinkconnector_property_scores
  }
  property {
    name  = "signatureType"
    type  = "string"
    value = var.sentilinkconnector_property_signature_type
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.sentilinkconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.sentilinkconnector_property_state
  }
  property {
    name  = "token"
    type  = "string"
    value = var.sentilinkconnector_property_token
  }
  property {
    name  = "userCreated"
    type  = "string"
    value = var.sentilinkconnector_property_user_created
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.sentilinkconnector_property_user_id
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.sentilinkconnector_property_zip
  }
}
