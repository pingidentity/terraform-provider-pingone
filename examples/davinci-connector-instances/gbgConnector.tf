resource "pingone_davinci_connector_instance" "gbgConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "gbgConnector"
  }
  name = "My awesome gbgConnector"
  property {
    name  = "buildingNumber"
    type  = "string"
    value = var.gbgconnector_property_building_number
  }
  property {
    name  = "city"
    type  = "string"
    value = var.gbgconnector_property_city
  }
  property {
    name  = "clientReference"
    type  = "string"
    value = var.gbgconnector_property_client_reference
  }
  property {
    name  = "countryOfBirth"
    type  = "string"
    value = var.gbgconnector_property_country_of_birth
  }
  property {
    name  = "countryOfOrigin"
    type  = "string"
    value = var.gbgconnector_property_country_of_origin
  }
  property {
    name  = "cpfNumber"
    type  = "string"
    value = var.gbgconnector_property_cpf_number
  }
  property {
    name  = "customerReference"
    type  = "string"
    value = var.gbgconnector_property_customer_reference
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.gbgconnector_property_dob
  }
  property {
    name  = "driversLicenceNo"
    type  = "string"
    value = var.gbgconnector_property_drivers_licence_no
  }
  property {
    name  = "electricityNumber"
    type  = "string"
    value = var.gbgconnector_property_electricity_number
  }
  property {
    name  = "email"
    type  = "string"
    value = var.gbgconnector_property_email
  }
  property {
    name  = "europeanIDCountryOfIssue"
    type  = "string"
    value = var.gbgconnector_property_european_idcountry_of_issue
  }
  property {
    name  = "europeanIDCountryOfNationality"
    type  = "string"
    value = var.gbgconnector_property_european_idcountry_of_nationality
  }
  property {
    name  = "europeanIDExpiryDate"
    type  = "string"
    value = var.gbgconnector_property_european_idexpiry_date
  }
  property {
    name  = "europeanIDLine1"
    type  = "string"
    value = var.gbgconnector_property_european_idline1
  }
  property {
    name  = "europeanIDLine2"
    type  = "string"
    value = var.gbgconnector_property_european_idline2
  }
  property {
    name  = "europeanIDLine3"
    type  = "string"
    value = var.gbgconnector_property_european_idline3
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.gbgconnector_property_first_name
  }
  property {
    name  = "gbgCountry"
    type  = "string"
    value = var.gbgconnector_property_gbg_country
  }
  property {
    name  = "gender"
    type  = "string"
    value = var.gbgconnector_property_gender
  }
  property {
    name  = "idCountry"
    type  = "string"
    value = var.gbgconnector_property_id_country
  }
  property {
    name  = "idNumber"
    type  = "string"
    value = var.gbgconnector_property_id_number
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.gbgconnector_property_last_name
  }
  property {
    name  = "medicareCardNo"
    type  = "string"
    value = var.gbgconnector_property_medicare_card_no
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.gbgconnector_property_middle_name
  }
  property {
    name  = "mobile"
    type  = "string"
    value = var.gbgconnector_property_mobile
  }
  property {
    name  = "passportExpiry"
    type  = "string"
    value = var.gbgconnector_property_passport_expiry
  }
  property {
    name  = "passportIssueDate"
    type  = "string"
    value = var.gbgconnector_property_passport_issue_date
  }
  property {
    name  = "passportNo"
    type  = "string"
    value = var.gbgconnector_property_passport_no
  }
  property {
    name  = "password"
    type  = "string"
    value = var.gbgconnector_property_password
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.gbgconnector_property_phone
  }
  property {
    name  = "placeOfBirth"
    type  = "string"
    value = var.gbgconnector_property_place_of_birth
  }
  property {
    name  = "poBox"
    type  = "string"
    value = var.gbgconnector_property_po_box
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.gbgconnector_property_postal_code
  }
  property {
    name  = "premise"
    type  = "string"
    value = var.gbgconnector_property_premise
  }
  property {
    name  = "principality"
    type  = "string"
    value = var.gbgconnector_property_principality
  }
  property {
    name  = "profileId"
    type  = "string"
    value = var.gbgconnector_property_profile_id
  }
  property {
    name  = "provinceOfBirth"
    type  = "string"
    value = var.gbgconnector_property_province_of_birth
  }
  property {
    name  = "region"
    type  = "string"
    value = var.gbgconnector_property_region
  }
  property {
    name  = "requestUrl"
    type  = "string"
    value = var.gbgconnector_property_request_url
  }
  property {
    name  = "secondLastName"
    type  = "string"
    value = var.gbgconnector_property_second_last_name
  }
  property {
    name  = "sin"
    type  = "string"
    value = var.gbgconnector_property_sin
  }
  property {
    name  = "soapAction"
    type  = "string"
    value = var.gbgconnector_property_soap_action
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.gbgconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.gbgconnector_property_state
  }
  property {
    name  = "streetAddress"
    type  = "string"
    value = var.gbgconnector_property_street_address
  }
  property {
    name  = "subBuilding"
    type  = "string"
    value = var.gbgconnector_property_sub_building
  }
  property {
    name  = "subCity"
    type  = "string"
    value = var.gbgconnector_property_sub_city
  }
  property {
    name  = "subStreet"
    type  = "string"
    value = var.gbgconnector_property_sub_street
  }
  property {
    name  = "taxIDNo"
    type  = "string"
    value = var.gbgconnector_property_tax_idno
  }
  property {
    name  = "title"
    type  = "string"
    value = var.gbgconnector_property_title
  }
  property {
    name  = "ukDriverLicenseExpiryDate"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_expiry_date
  }
  property {
    name  = "ukDriverLicenseIssueDate"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_issue_date
  }
  property {
    name  = "ukDriverLicenseIssueNo"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_issue_no
  }
  property {
    name  = "ukDriverLicenseMailSort"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_mail_sort
  }
  property {
    name  = "ukDriverLicenseMicrofiche"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_microfiche
  }
  property {
    name  = "ukDriverLicenseNo"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_no
  }
  property {
    name  = "ukDriverLicensePostcode"
    type  = "string"
    value = var.gbgconnector_property_uk_driver_license_postcode
  }
  property {
    name  = "ukPassportExpiryDate"
    type  = "string"
    value = var.gbgconnector_property_uk_passport_expiry_date
  }
  property {
    name  = "ukPassportNo"
    type  = "string"
    value = var.gbgconnector_property_uk_passport_no
  }
  property {
    name  = "usDriverLicenseNo"
    type  = "string"
    value = var.gbgconnector_property_us_driver_license_no
  }
  property {
    name  = "usDriverLicenseState"
    type  = "string"
    value = var.gbgconnector_property_us_driver_license_state
  }
  property {
    name  = "username"
    type  = "string"
    value = var.gbgconnector_property_username
  }
  property {
    name  = "workPhone"
    type  = "string"
    value = var.gbgconnector_property_work_phone
  }
}
