resource "pingone_davinci_connector_instance" "tutloxpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tutloxpConnector"
  }
  name = "My awesome tutloxpConnector"
  property {
    name  = "addressLineOne"
    type  = "string"
    value = var.tutloxpconnector_property_address_line_one
  }
  property {
    name  = "addressLineThree"
    type  = "string"
    value = var.tutloxpconnector_property_address_line_three
  }
  property {
    name  = "addressLineTwo"
    type  = "string"
    value = var.tutloxpconnector_property_address_line_two
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.tutloxpconnector_property_api_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.tutloxpconnector_property_city
  }
  property {
    name  = "country"
    type  = "string"
    value = var.tutloxpconnector_property_country
  }
  property {
    name  = "doNotModifySearch"
    type  = "string"
    value = var.tutloxpconnector_property_do_not_modify_search
  }
  property {
    name  = "dobDay"
    type  = "string"
    value = var.tutloxpconnector_property_dob_day
  }
  property {
    name  = "dobMonth"
    type  = "string"
    value = var.tutloxpconnector_property_dob_month
  }
  property {
    name  = "dobYear"
    type  = "string"
    value = var.tutloxpconnector_property_dob_year
  }
  property {
    name  = "dppaCode"
    type  = "string"
    value = var.tutloxpconnector_property_dppa_code
  }
  property {
    name  = "echoTestInput"
    type  = "string"
    value = var.tutloxpconnector_property_echo_test_input
  }
  property {
    name  = "findBestMatch"
    type  = "string"
    value = var.tutloxpconnector_property_find_best_match
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.tutloxpconnector_property_first_name
  }
  property {
    name  = "glbCode"
    type  = "string"
    value = var.tutloxpconnector_property_glb_code
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.tutloxpconnector_property_last_name
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.tutloxpconnector_property_middle_name
  }
  property {
    name  = "password"
    type  = "string"
    value = var.tutloxpconnector_property_password
  }
  property {
    name  = "reportToken"
    type  = "string"
    value = var.tutloxpconnector_property_report_token
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.tutloxpconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.tutloxpconnector_property_state
  }
  property {
    name  = "username"
    type  = "string"
    value = var.tutloxpconnector_property_username
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.tutloxpconnector_property_zip
  }
}
