resource "pingone_davinci_connector_instance" "onfidoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "onfidoConnector"
  }
  name = "My awesome onfidoConnector"
  property {
    name  = "address1"
    type  = "string"
    value = var.onfidoconnector_property_address1
  }
  property {
    name  = "address2"
    type  = "string"
    value = var.onfidoconnector_property_address2
  }
  property {
    name  = "androidPackageName"
    type  = "string"
    value = var.onfidoconnector_property_android_package_name
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.onfidoconnector_property_api_key
  }
  property {
    name  = "authDescription"
    type  = "string"
    value = var.onfidoconnector_property_auth_description
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.onfidoconnector_property_base_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.onfidoconnector_property_city
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.onfidoconnector_property_connector_name
  }
  property {
    name  = "countryCode"
    type  = "string"
    value = var.onfidoconnector_property_country_code
  }
  property {
    name  = "customizeSteps"
    type  = "string"
    value = var.onfidoconnector_property_customize_steps
  }
  property {
    name  = "description"
    type  = "string"
    value = var.onfidoconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.onfidoconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.onfidoconnector_property_details2
  }
  property {
    name  = "dob"
    type  = "string"
    value = var.onfidoconnector_property_dob
  }
  property {
    name  = "email"
    type  = "string"
    value = var.onfidoconnector_property_email
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.onfidoconnector_property_first_name
  }
  property {
    name  = "iOSBundleId"
    type  = "string"
    value = var.onfidoconnector_property_i_osbundle_id
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.onfidoconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.onfidoconnector_property_icon_url_png
  }
  property {
    name  = "javascriptCSSUrl"
    type  = "string"
    value = var.javascript_css_url
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.onfidoconnector_property_javascript_cdn_url
  }
  property {
    name  = "language"
    type  = "string"
    value = var.onfidoconnector_property_language
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.onfidoconnector_property_last_name
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.onfidoconnector_property_phone_number
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.onfidoconnector_property_postal_code
  }
  property {
    name  = "referenceStepsList"
    type  = "string"
    value = var.onfidoconnector_property_reference_steps_list
  }
  property {
    name  = "referrerUrl"
    type  = "string"
    value = var.onfidoconnector_property_referrer_url
  }
  property {
    name  = "reportTypes"
    type  = "string"
    value = var.onfidoconnector_property_report_types
  }
  property {
    name  = "retrieveReports"
    type  = "string"
    value = var.onfidoconnector_property_retrieve_reports
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.onfidoconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.onfidoconnector_property_screen2_config
  }
  property {
    name  = "shouldCloseOnOverlayClick"
    type  = "string"
    value = var.onfidoconnector_property_should_close_on_overlay_click
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.onfidoconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.onfidoconnector_property_show_cred_added_via
  }
  property {
    name  = "stepsList"
    type  = "string"
    value = var.onfidoconnector_property_steps_list
  }
  property {
    name  = "title"
    type  = "string"
    value = var.onfidoconnector_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.onfidoconnector_property_tool_tip
  }
  property {
    name  = "useLanguage"
    type  = "string"
    value = var.onfidoconnector_property_use_language
  }
  property {
    name  = "useModal"
    type  = "string"
    value = var.onfidoconnector_property_use_modal
  }
  property {
    name  = "viewDescriptions"
    type  = "string"
    value = var.onfidoconnector_property_view_descriptions
  }
  property {
    name  = "viewTitle"
    type  = "string"
    value = var.onfidoconnector_property_view_title
  }
}
