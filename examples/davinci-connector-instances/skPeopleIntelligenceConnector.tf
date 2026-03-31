resource "pingone_davinci_connector_instance" "skPeopleIntelligenceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skPeopleIntelligenceConnector"
  }
  name = "My awesome skPeopleIntelligenceConnector"
  property {
    name  = "address"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_address
  }
  property {
    name  = "ageMax"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_age_max
  }
  property {
    name  = "ageMin"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_age_min
  }
  property {
    name  = "authUrl"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_auth_url
  }
  property {
    name  = "bankRuptcyCaseNum"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_bank_ruptcy_case_num
  }
  property {
    name  = "businessName"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_business_name
  }
  property {
    name  = "city"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_city
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_client_secret
  }
  property {
    name  = "criminalCaseNum"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_criminal_case_num
  }
  property {
    name  = "criminalCounty"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_criminal_county
  }
  property {
    name  = "criminalSlate"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_criminal_slate
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_dob
  }
  property {
    name  = "dppa"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_dppa
  }
  property {
    name  = "duns"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_duns
  }
  property {
    name  = "email"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_email
  }
  property {
    name  = "emailHash"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_email_hash
  }
  property {
    name  = "fields"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_fields
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_first_name
  }
  property {
    name  = "glba"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_glba
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_ip
  }
  property {
    name  = "judgementCaseNum"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_judgement_case_num
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_last_name
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_middle_name
  }
  property {
    name  = "naicsCodes"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_naics_codes
  }
  property {
    name  = "nickNameSearch"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_nick_name_search
  }
  property {
    name  = "nnumber"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_nnumber
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_phone
  }
  property {
    name  = "pidList"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_pid_list
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_postal_code
  }
  property {
    name  = "searchUrl"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_search_url
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_state
  }
  property {
    name  = "taxId"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_tax_id
  }
  property {
    name  = "vehicleCounty"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_vehicle_county
  }
  property {
    name  = "vehiclePlate"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_vehicle_plate
  }
  property {
    name  = "vehicleSlate"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_vehicle_slate
  }
  property {
    name  = "vehicleVin"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_vehicle_vin
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.skpeopleintelligenceconnector_property_zip
  }
}
