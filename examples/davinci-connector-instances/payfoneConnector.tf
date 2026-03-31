resource "pingone_davinci_connector_instance" "payfoneConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "payfoneConnector"
  }
  name = "My awesome payfoneConnector"
  property {
    name  = "address"
    type  = "string"
    value = var.payfoneconnector_property_address
  }
  property {
    name  = "address1"
    type  = "string"
    value = var.payfoneconnector_property_address1
  }
  property {
    name  = "address2"
    type  = "string"
    value = var.payfoneconnector_property_address2
  }
  property {
    name  = "appClientId"
    type  = "string"
    value = var.payfoneconnector_property_app_client_id
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.payfoneconnector_property_auth_type
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.payfoneconnector_property_base_url
  }
  property {
    name  = "button"
    type  = "string"
    value = var.payfoneconnector_property_button
  }
  property {
    name  = "city"
    type  = "string"
    value = var.payfoneconnector_property_city
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.payfoneconnector_property_client_id
  }
  property {
    name  = "consentStatus"
    type  = "string"
    value = var.payfoneconnector_property_consent_status
  }
  property {
    name  = "countryCode"
    type  = "string"
    value = var.payfoneconnector_property_country_code
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.payfoneconnector_property_custom_auth
  }
  property {
    name  = "detailsFlag"
    type  = "string"
    value = var.payfoneconnector_property_details_flag
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.payfoneconnector_property_dob
  }
  property {
    name  = "emailAddress"
    type  = "string"
    value = var.payfoneconnector_property_email_address
  }
  property {
    name  = "extendedAddress"
    type  = "string"
    value = var.payfoneconnector_property_extended_address
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.payfoneconnector_property_first_name
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.payfoneconnector_property_last_name
  }
  property {
    name  = "lastVerified"
    type  = "string"
    value = var.payfoneconnector_property_last_verified
  }
  property {
    name  = "numberOfAddresses"
    type  = "string"
    value = var.payfoneconnector_property_number_of_addresses
  }
  property {
    name  = "numberOfEmails"
    type  = "string"
    value = var.payfoneconnector_property_number_of_emails
  }
  property {
    name  = "password"
    type  = "string"
    value = var.payfoneconnector_property_password
  }
  property {
    name  = "payfoneAlias"
    type  = "string"
    value = var.payfoneconnector_property_payfone_alias
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.payfoneconnector_property_phone_number
  }
  property {
    name  = "phoneUpdateFlag"
    type  = "string"
    value = var.payfoneconnector_property_phone_update_flag
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.payfoneconnector_property_postal_code
  }
  property {
    name  = "region"
    type  = "string"
    value = var.payfoneconnector_property_region
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.payfoneconnector_property_show_powered_by
  }
  property {
    name  = "simulatorMode"
    type  = "string"
    value = var.payfoneconnector_property_simulator_mode
  }
  property {
    name  = "simulatorPhoneNumber"
    type  = "string"
    value = var.payfoneconnector_property_simulator_phone_number
  }
  property {
    name  = "skCallbackBaseUrl"
    type  = "string"
    value = var.payfoneconnector_property_sk_callback_base_url
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.payfoneconnector_property_ssn
  }
  property {
    name  = "ssnLast4"
    type  = "string"
    value = var.payfoneconnector_property_ssn_last4
  }
  property {
    name  = "trustField"
    type  = "string"
    value = var.payfoneconnector_property_trust_field
  }
  property {
    name  = "trustScoreFlag"
    type  = "string"
    value = var.payfoneconnector_property_trust_score_flag
  }
  property {
    name  = "username"
    type  = "string"
    value = var.payfoneconnector_property_username
  }
  property {
    name  = "vfp"
    type  = "string"
    value = var.payfoneconnector_property_vfp
  }
}
