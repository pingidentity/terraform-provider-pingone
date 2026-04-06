resource "pingone_davinci_connector_instance" "telesignConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "telesignConnector"
  }
  name = "My awesome telesignConnector"
  property {
    name  = "address"
    type  = "string"
    value = var.telesignconnector_property_address
  }
  property {
    name  = "authDescription"
    type  = "string"
    value = var.telesignconnector_property_auth_description
  }
  property {
    name  = "authDescriptionDetail"
    type  = "string"
    value = var.telesignconnector_property_auth_description_detail
  }
  property {
    name  = "billing_postal_code"
    type  = "string"
    value = var.telesignconnector_property_billing_postal_code
  }
  property {
    name  = "city"
    type  = "string"
    value = var.telesignconnector_property_city
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.telesignconnector_property_connector_name
  }
  property {
    name  = "consentMethod"
    type  = "string"
    value = var.telesignconnector_property_consent_method
  }
  property {
    name  = "consentTimeStamp"
    type  = "string"
    value = var.telesignconnector_property_consent_time_stamp
  }
  property {
    name  = "country"
    type  = "string"
    value = var.telesignconnector_property_country
  }
  property {
    name  = "date_of_birth"
    type  = "string"
    value = var.telesignconnector_property_date_of_birth
  }
  property {
    name  = "description"
    type  = "string"
    value = var.telesignconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.telesignconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.telesignconnector_property_details2
  }
  property {
    name  = "first_name"
    type  = "string"
    value = var.telesignconnector_property_first_name
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.telesignconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.telesignconnector_property_icon_url_png
  }
  property {
    name  = "last_name"
    type  = "string"
    value = var.telesignconnector_property_last_name
  }
  property {
    name  = "last_verified"
    type  = "string"
    value = var.telesignconnector_property_last_verified
  }
  property {
    name  = "message"
    type  = "string"
    value = var.telesignconnector_property_message
  }
  property {
    name  = "otp"
    type  = "string"
    value = var.telesignconnector_property_otp
  }
  property {
    name  = "password"
    type  = "string"
    value = var.telesignconnector_property_password
  }
  property {
    name  = "past_x_days"
    type  = "string"
    value = var.telesignconnector_property_past_x_days
  }
  property {
    name  = "phoneIdAddOns"
    type  = "string"
    value = var.telesignconnector_property_phone_id_add_ons
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.telesignconnector_property_phone_number
  }
  property {
    name  = "postal_code"
    type  = "string"
    value = var.telesignconnector_property_postal_code
  }
  property {
    name  = "providerName"
    type  = "string"
    value = var.telesignconnector_property_provider_name
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.telesignconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.telesignconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.telesignconnector_property_screen2_config
  }
  property {
    name  = "screen3Config"
    type  = "string"
    value = var.telesignconnector_property_screen3_config
  }
  property {
    name  = "screen4Config"
    type  = "string"
    value = var.telesignconnector_property_screen4_config
  }
  property {
    name  = "screen5Config"
    type  = "string"
    value = var.telesignconnector_property_screen5_config
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.telesignconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.telesignconnector_property_show_cred_added_via
  }
  property {
    name  = "state"
    type  = "string"
    value = var.telesignconnector_property_state
  }
  property {
    name  = "title"
    type  = "string"
    value = var.telesignconnector_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.telesignconnector_property_tool_tip
  }
  property {
    name  = "useScore20"
    type  = "string"
    value = var.telesignconnector_property_use_score20
  }
  property {
    name  = "username"
    type  = "string"
    value = var.telesignconnector_property_username
  }
}
