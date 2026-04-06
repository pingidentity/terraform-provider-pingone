resource "pingone_davinci_connector_instance" "dataZooConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "dataZooConnector"
  }
  name = "My awesome dataZooConnector"
  property {
    name  = "addressElement1"
    type  = "string"
    value = var.datazooconnector_property_address_element1
  }
  property {
    name  = "addressElement2"
    type  = "string"
    value = var.datazooconnector_property_address_element2
  }
  property {
    name  = "addressElement3"
    type  = "string"
    value = var.datazooconnector_property_address_element3
  }
  property {
    name  = "addressElement4"
    type  = "string"
    value = var.datazooconnector_property_address_element4
  }
  property {
    name  = "addressElement5"
    type  = "string"
    value = var.datazooconnector_property_address_element5
  }
  property {
    name  = "australiaService"
    type  = "string"
    value = var.datazooconnector_property_australia_service
  }
  property {
    name  = "austriaService"
    type  = "string"
    value = var.datazooconnector_property_austria_service
  }
  property {
    name  = "birthCertificateNumber"
    type  = "string"
    value = var.datazooconnector_property_birth_certificate_number
  }
  property {
    name  = "birthRegistrationDate"
    type  = "string"
    value = var.datazooconnector_property_birth_registration_date
  }
  property {
    name  = "birthRegistrationNumber"
    type  = "string"
    value = var.datazooconnector_property_birth_registration_number
  }
  property {
    name  = "birthRegistrationState"
    type  = "string"
    value = var.datazooconnector_property_birth_registration_state
  }
  property {
    name  = "centrelinkCardExpiry"
    type  = "string"
    value = var.datazooconnector_property_centrelink_card_expiry
  }
  property {
    name  = "centrelinkCardType"
    type  = "string"
    value = var.datazooconnector_property_centrelink_card_type
  }
  property {
    name  = "centrelinkCustomerReferenceNo"
    type  = "string"
    value = var.datazooconnector_property_centrelink_customer_reference_no
  }
  property {
    name  = "clientReference"
    type  = "string"
    value = var.datazooconnector_property_client_reference
  }
  property {
    name  = "company"
    type  = "string"
    value = var.datazooconnector_property_company
  }
  property {
    name  = "consent"
    type  = "string"
    value = var.datazooconnector_property_consent
  }
  property {
    name  = "country"
    type  = "string"
    value = var.datazooconnector_property_country
  }
  property {
    name  = "countryCode"
    type  = "string"
    value = var.datazooconnector_property_country_code
  }
  property {
    name  = "datazooCountryCodes"
    type  = "string"
    value = var.datazooconnector_property_datazoo_country_codes
  }
  property {
    name  = "dateOfIssue"
    type  = "string"
    value = var.datazooconnector_property_date_of_issue
  }
  property {
    name  = "denmarkService"
    type  = "string"
    value = var.datazooconnector_property_denmark_service
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.datazooconnector_property_dob
  }
  property {
    name  = "driversLicenceExpiryDate"
    type  = "string"
    value = var.datazooconnector_property_drivers_licence_expiry_date
  }
  property {
    name  = "driversLicenceNo"
    type  = "string"
    value = var.datazooconnector_property_drivers_licence_no
  }
  property {
    name  = "driversLicenceStateOfIssue"
    type  = "string"
    value = var.datazooconnector_property_drivers_licence_state_of_issue
  }
  property {
    name  = "email"
    type  = "string"
    value = var.datazooconnector_property_email
  }
  property {
    name  = "epic"
    type  = "string"
    value = var.datazooconnector_property_epic
  }
  property {
    name  = "fileNo"
    type  = "string"
    value = var.datazooconnector_property_file_no
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.datazooconnector_property_first_name
  }
  property {
    name  = "franceService"
    type  = "string"
    value = var.datazooconnector_property_france_service
  }
  property {
    name  = "fullName"
    type  = "string"
    value = var.datazooconnector_property_full_name
  }
  property {
    name  = "genderSelect"
    type  = "string"
    value = var.datazooconnector_property_gender_select
  }
  property {
    name  = "germanyService"
    type  = "string"
    value = var.datazooconnector_property_germany_service
  }
  property {
    name  = "greeceService"
    type  = "string"
    value = var.datazooconnector_property_greece_service
  }
  property {
    name  = "indiaService"
    type  = "string"
    value = var.datazooconnector_property_india_service
  }
  property {
    name  = "italyService"
    type  = "string"
    value = var.datazooconnector_property_italy_service
  }
  property {
    name  = "landlineNo"
    type  = "string"
    value = var.datazooconnector_property_landline_no
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.datazooconnector_property_last_name
  }
  property {
    name  = "medicareCardExpiry"
    type  = "string"
    value = var.datazooconnector_property_medicare_card_expiry
  }
  property {
    name  = "medicareCardNo"
    type  = "string"
    value = var.datazooconnector_property_medicare_card_no
  }
  property {
    name  = "medicareCardType"
    type  = "string"
    value = var.datazooconnector_property_medicare_card_type
  }
  property {
    name  = "medicareReferenceNo"
    type  = "string"
    value = var.datazooconnector_property_medicare_reference_no
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.datazooconnector_property_middle_name
  }
  property {
    name  = "mobile"
    type  = "string"
    value = var.datazooconnector_property_mobile
  }
  property {
    name  = "netherlandsService"
    type  = "string"
    value = var.datazooconnector_property_netherlands_service
  }
  property {
    name  = "norwayService"
    type  = "string"
    value = var.datazooconnector_property_norway_service
  }
  property {
    name  = "panNumber"
    type  = "string"
    value = var.datazooconnector_property_pan_number
  }
  property {
    name  = "passportCountry"
    type  = "string"
    value = var.datazooconnector_property_passport_country
  }
  property {
    name  = "passportExpiry"
    type  = "string"
    value = var.datazooconnector_property_passport_expiry
  }
  property {
    name  = "passportIssuerCountry"
    type  = "string"
    value = var.datazooconnector_property_passport_issuer_country
  }
  property {
    name  = "passportNo"
    type  = "string"
    value = var.datazooconnector_property_passport_no
  }
  property {
    name  = "password"
    type  = "string"
    value = var.datazooconnector_property_password
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.datazooconnector_property_phone
  }
  property {
    name  = "polandService"
    type  = "string"
    value = var.datazooconnector_property_poland_service
  }
  property {
    name  = "portugalService"
    type  = "string"
    value = var.datazooconnector_property_portugal_service
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.datazooconnector_property_postal_code
  }
  property {
    name  = "romaniaService"
    type  = "string"
    value = var.datazooconnector_property_romania_service
  }
  property {
    name  = "school"
    type  = "string"
    value = var.datazooconnector_property_school
  }
  property {
    name  = "spainService"
    type  = "string"
    value = var.datazooconnector_property_spain_service
  }
  property {
    name  = "switzerlandService"
    type  = "string"
    value = var.datazooconnector_property_switzerland_service
  }
  property {
    name  = "taxIDNo"
    type  = "string"
    value = var.datazooconnector_property_tax_idno
  }
  property {
    name  = "username"
    type  = "string"
    value = var.datazooconnector_property_username
  }
}
