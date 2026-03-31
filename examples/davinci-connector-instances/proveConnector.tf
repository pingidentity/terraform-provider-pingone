resource "pingone_davinci_connector_instance" "proveConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "proveConnector"
  }
  name = "My awesome proveConnector"
  property {
    name  = "address"
    type  = "string"
    value = var.proveconnector_property_address
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.proveconnector_property_base_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.proveconnector_property_city
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.proveconnector_property_client_id
  }
  property {
    name  = "consentStat"
    type  = "string"
    value = var.proveconnector_property_consent_stat
  }
  property {
    name  = "country"
    type  = "string"
    value = var.proveconnector_property_country
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.proveconnector_property_dob
  }
  property {
    name  = "extendedAddress"
    type  = "string"
    value = var.proveconnector_property_extended_address
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.proveconnector_property_first_name
  }
  property {
    name  = "grantType"
    type  = "string"
    value = var.proveconnector_property_grant_type
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.proveconnector_property_last_name
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.proveconnector_property_middle_name
  }
  property {
    name  = "password"
    type  = "string"
    value = var.proveconnector_property_password
  }
  property {
    name  = "payfoneId"
    type  = "string"
    value = var.proveconnector_property_payfone_id
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.proveconnector_property_phone
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.proveconnector_property_postal_code
  }
  property {
    name  = "region"
    type  = "string"
    value = var.proveconnector_property_region
  }
  property {
    name  = "requestId"
    type  = "string"
    value = var.proveconnector_property_request_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.proveconnector_property_username
  }
}
