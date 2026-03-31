resource "pingone_davinci_connector_instance" "finicityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "finicityConnector"
  }
  name = "My awesome finicityConnector"
  property {
    name  = "address"
    type  = "string"
    value = var.finicityconnector_property_address
  }
  property {
    name  = "appKey"
    type  = "string"
    value = var.finicityconnector_property_app_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.finicityconnector_property_base_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.finicityconnector_property_city
  }
  property {
    name  = "dayOfMonth"
    type  = "string"
    value = var.finicityconnector_property_day_of_month
  }
  property {
    name  = "finicityConnectType"
    type  = "string"
    value = var.finicityconnector_property_finicity_connect_type
  }
  property {
    name  = "finicityCustomerType"
    type  = "string"
    value = var.finicityconnector_property_finicity_customer_type
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.finicityconnector_property_first_name
  }
  property {
    name  = "fromDate"
    type  = "string"
    value = var.finicityconnector_property_from_date
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.finicityconnector_property_last_name
  }
  property {
    name  = "mainHeaderText"
    type  = "string"
    value = var.finicityconnector_property_main_header_text
  }
  property {
    name  = "month"
    type  = "string"
    value = var.finicityconnector_property_month
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.finicityconnector_property_next_button_text
  }
  property {
    name  = "partnerId"
    type  = "string"
    value = var.finicityconnector_property_partner_id
  }
  property {
    name  = "partnerSecret"
    type  = "string"
    value = var.finicityconnector_property_partner_secret
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.finicityconnector_property_phone
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.finicityconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.finicityconnector_property_screen1_config
  }
  property {
    name  = "ssn"
    type  = "string"
    value = var.finicityconnector_property_ssn
  }
  property {
    name  = "state"
    type  = "string"
    value = var.finicityconnector_property_state
  }
  property {
    name  = "title"
    type  = "string"
    value = var.finicityconnector_property_title
  }
  property {
    name  = "year"
    type  = "string"
    value = var.finicityconnector_property_year
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.finicityconnector_property_zip
  }
}
