resource "pingone_davinci_connector_instance" "credovaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "credovaConnector"
  }
  name = "My awesome credovaConnector"
  property {
    name  = "bankAccountInstitutionName"
    type  = "string"
    value = var.credovaconnector_property_bank_account_institution_name
  }
  property {
    name  = "bankAccountNameOnAccount"
    type  = "string"
    value = var.credovaconnector_property_bank_account_name_on_account
  }
  property {
    name  = "bankAccountNumber"
    type  = "string"
    value = var.credovaconnector_property_bank_account_number
  }
  property {
    name  = "bankAccountRoutingNumber"
    type  = "string"
    value = var.credovaconnector_property_bank_account_routing_number
  }
  property {
    name  = "bankAccountType"
    type  = "string"
    value = var.credovaconnector_property_bank_account_type
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.credovaconnector_property_base_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.credovaconnector_property_city
  }
  property {
    name  = "dateOfBirth"
    type  = "string"
    value = var.credovaconnector_property_date_of_birth
  }
  property {
    name  = "email"
    type  = "string"
    value = var.credovaconnector_property_email
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.credovaconnector_property_first_name
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.credovaconnector_property_last_name
  }
  property {
    name  = "middleInitial"
    type  = "string"
    value = var.credovaconnector_property_middle_initial
  }
  property {
    name  = "mobilePhone"
    type  = "string"
    value = var.credovaconnector_property_mobile_phone
  }
  property {
    name  = "offerId"
    type  = "string"
    value = var.credovaconnector_property_offer_id
  }
  property {
    name  = "password"
    type  = "string"
    value = var.credovaconnector_property_password
  }
  property {
    name  = "publicId"
    type  = "string"
    value = var.credovaconnector_property_public_id
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.credovaconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.credovaconnector_property_state
  }
  property {
    name  = "storeCode"
    type  = "string"
    value = var.credovaconnector_property_store_code
  }
  property {
    name  = "street"
    type  = "string"
    value = var.credovaconnector_property_street
  }
  property {
    name  = "suiteApartment"
    type  = "string"
    value = var.credovaconnector_property_suite_apartment
  }
  property {
    name  = "username"
    type  = "string"
    value = var.credovaconnector_property_username
  }
  property {
    name  = "zipCode"
    type  = "string"
    value = var.credovaconnector_property_zip_code
  }
}
