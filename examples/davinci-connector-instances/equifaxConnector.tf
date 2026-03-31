resource "pingone_davinci_connector_instance" "equifaxConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "equifaxConnector"
  }
  name = "My awesome equifaxConnector"
  property {
    name  = "accountNumber"
    type  = "string"
    value = var.equifaxconnector_property_account_number
  }
  property {
    name  = "accountType"
    type  = "string"
    value = var.equifaxconnector_property_account_type
  }
  property {
    name  = "addressLine1"
    type  = "string"
    value = var.equifaxconnector_property_address_line1
  }
  property {
    name  = "addressLine2"
    type  = "string"
    value = var.equifaxconnector_property_address_line2
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.equifaxconnector_property_base_url
  }
  property {
    name  = "cid"
    type  = "string"
    value = var.equifaxconnector_property_cid
  }
  property {
    name  = "city"
    type  = "string"
    value = var.equifaxconnector_property_city
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.equifaxconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.equifaxconnector_property_client_secret
  }
  property {
    name  = "cnx"
    type  = "string"
    value = var.equifaxconnector_property_cnx
  }
  property {
    name  = "currentAddressLine"
    type  = "string"
    value = var.equifaxconnector_property_current_address_line
  }
  property {
    name  = "currentCity"
    type  = "string"
    value = var.equifaxconnector_property_current_city
  }
  property {
    name  = "currentState"
    type  = "string"
    value = var.equifaxconnector_property_current_state
  }
  property {
    name  = "currentZip"
    type  = "string"
    value = var.equifaxconnector_property_current_zip
  }
  property {
    name  = "customerId"
    type  = "string"
    value = var.equifaxconnector_property_customer_id
  }
  property {
    name  = "deliveryChannel"
    type  = "string"
    value = var.equifaxconnector_property_delivery_channel
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.equifaxconnector_property_dob
  }
  property {
    name  = "dobDay"
    type  = "string"
    value = var.equifaxconnector_property_dob_day
  }
  property {
    name  = "dobMonth"
    type  = "string"
    value = var.equifaxconnector_property_dob_month
  }
  property {
    name  = "dobYear"
    type  = "string"
    value = var.equifaxconnector_property_dob_year
  }
  property {
    name  = "driversLicenseNumber"
    type  = "string"
    value = var.equifaxconnector_property_drivers_license_number
  }
  property {
    name  = "driversLicenseState"
    type  = "string"
    value = var.equifaxconnector_property_drivers_license_state
  }
  property {
    name  = "efxClientCorrelationId"
    type  = "string"
    value = var.equifaxconnector_property_efx_client_correlation_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.equifaxconnector_property_email
  }
  property {
    name  = "equifaxQuery"
    type  = "string"
    value = var.equifaxconnector_property_equifax_query
  }
  property {
    name  = "equifaxSoapApiEnvironment"
    type  = "string"
    value = var.equifaxconnector_property_equifax_soap_api_environment
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.equifaxconnector_property_first_name
  }
  property {
    name  = "formerAddressLine"
    type  = "string"
    value = var.equifaxconnector_property_former_address_line
  }
  property {
    name  = "formerCity"
    type  = "string"
    value = var.equifaxconnector_property_former_city
  }
  property {
    name  = "formerState"
    type  = "string"
    value = var.equifaxconnector_property_former_state
  }
  property {
    name  = "formerZip"
    type  = "string"
    value = var.equifaxconnector_property_former_zip
  }
  property {
    name  = "hitCode"
    type  = "string"
    value = var.equifaxconnector_property_hit_code
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.equifaxconnector_property_last_name
  }
  property {
    name  = "memberNumber"
    type  = "string"
    value = var.equifaxconnector_property_member_number
  }
  property {
    name  = "middleInitial"
    type  = "string"
    value = var.equifaxconnector_property_middle_initial
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.equifaxconnector_property_middle_name
  }
  property {
    name  = "nameSuffix"
    type  = "string"
    value = var.equifaxconnector_property_name_suffix
  }
  property {
    name  = "orchestrationCode"
    type  = "string"
    value = var.equifaxconnector_property_orchestration_code
  }
  property {
    name  = "password"
    type  = "string"
    value = var.equifaxconnector_property_password
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.equifaxconnector_property_phone_number
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.equifaxconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.equifaxconnector_property_state
  }
  property {
    name  = "synthetic2RulesCategory"
    type  = "string"
    value = var.equifaxconnector_property_synthetic2_rules_category
  }
  property {
    name  = "transactionTimestamp"
    type  = "string"
    value = var.equifaxconnector_property_transaction_timestamp
  }
  property {
    name  = "username"
    type  = "string"
    value = var.equifaxconnector_property_username
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.equifaxconnector_property_zip
  }
}
