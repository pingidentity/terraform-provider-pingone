resource "pingone_davinci_connector_instance" "sinchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sinchConnector"
  }
  name = "My awesome sinchConnector"
  property {
    name  = "acceptLanguage"
    type  = "string"
    value = var.sinchconnector_property_accept_language
  }
  property {
    name  = "applicationKey"
    type  = "string"
    value = var.sinchconnector_property_application_key
  }
  property {
    name  = "code"
    type  = "string"
    value = var.sinchconnector_property_code
  }
  property {
    name  = "codeType"
    type  = "string"
    value = var.sinchconnector_property_code_type
  }
  property {
    name  = "expiry"
    type  = "string"
    value = var.sinchconnector_property_expiry
  }
  property {
    name  = "id"
    type  = "string"
    value = var.sinchconnector_property_id
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.sinchconnector_property_phone_number
  }
  property {
    name  = "secretKey"
    type  = "string"
    value = var.sinchconnector_property_secret_key
  }
  property {
    name  = "verificationType"
    type  = "string"
    value = var.sinchconnector_property_verification_type
  }
}
