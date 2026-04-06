resource "pingone_davinci_connector_instance" "zoopConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "zoopConnector"
  }
  name = "My awesome zoopConnector"
  property {
    name  = "agencyId"
    type  = "string"
    value = var.zoopconnector_property_agency_id
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.zoopconnector_property_api_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.zoopconnector_property_api_url
  }
  property {
    name  = "consent"
    type  = "string"
    value = var.zoopconnector_property_consent
  }
  property {
    name  = "consentText"
    type  = "string"
    value = var.zoopconnector_property_consent_text
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.zoopconnector_property_dob
  }
  property {
    name  = "driverLicenseNumber"
    type  = "string"
    value = var.zoopconnector_property_driver_license_number
  }
  property {
    name  = "email"
    type  = "string"
    value = var.zoopconnector_property_email
  }
  property {
    name  = "fileBase64"
    type  = "string"
    value = var.zoopconnector_property_file_base64
  }
  property {
    name  = "getDetailedAddress"
    type  = "string"
    value = var.zoopconnector_property_get_detailed_address
  }
  property {
    name  = "gstNumber"
    type  = "string"
    value = var.zoopconnector_property_gst_number
  }
  property {
    name  = "mobile"
    type  = "string"
    value = var.zoopconnector_property_mobile
  }
  property {
    name  = "mode"
    type  = "string"
    value = var.zoopconnector_property_mode
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.zoopconnector_property_next_event
  }
  property {
    name  = "panNumber"
    type  = "string"
    value = var.zoopconnector_property_pan_number
  }
  property {
    name  = "password"
    type  = "string"
    value = var.zoopconnector_property_password
  }
  property {
    name  = "purpose"
    type  = "string"
    value = var.zoopconnector_property_purpose
  }
  property {
    name  = "voterIdNumber"
    type  = "string"
    value = var.zoopconnector_property_voter_id_number
  }
}
