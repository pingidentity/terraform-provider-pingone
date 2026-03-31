resource "pingone_davinci_connector_instance" "jumioConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jumioConnector"
  }
  name = "My awesome jumioConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.jumioconnector_property_api_key
  }
  property {
    name  = "authDescription"
    type  = "string"
    value = var.jumioconnector_property_auth_description
  }
  property {
    name  = "authUrl"
    type  = "string"
    value = var.jumioconnector_property_auth_url
  }
  property {
    name  = "authorizationTokenLifetime"
    type  = "string"
    value = var.jumioconnector_property_authorization_token_lifetime
  }
  property {
    name  = "baseColor"
    type  = "string"
    value = var.jumioconnector_property_base_color
  }
  property {
    name  = "bgColor"
    type  = "string"
    value = var.jumioconnector_property_bg_color
  }
  property {
    name  = "callbackUrl"
    type  = "string"
    value = var.jumioconnector_property_callback_url
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.jumioconnector_property_client_secret
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.jumioconnector_property_connector_name
  }
  property {
    name  = "country"
    type  = "string"
    value = var.jumioconnector_property_country
  }
  property {
    name  = "countryCode3"
    type  = "string"
    value = var.jumioconnector_property_country_code3
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.jumioconnector_property_custom_css
  }
  property {
    name  = "customDocumentCode"
    type  = "string"
    value = var.jumioconnector_property_custom_document_code
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.jumioconnector_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.jumioconnector_property_custom_script
  }
  property {
    name  = "description"
    type  = "string"
    value = var.jumioconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.jumioconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.jumioconnector_property_details2
  }
  property {
    name  = "doNotShowInIframe"
    type  = "string"
    value = var.jumioconnector_property_do_not_show_in_iframe
  }
  property {
    name  = "docVerificationUrl"
    type  = "string"
    value = var.jumioconnector_property_doc_verification_url
  }
  property {
    name  = "documentType"
    type  = "string"
    value = var.jumioconnector_property_document_type
  }
  property {
    name  = "enableExtraction"
    type  = "string"
    value = var.jumioconnector_property_enable_extraction
  }
  property {
    name  = "headerImageUrl"
    type  = "string"
    value = var.jumioconnector_property_header_image_url
  }
  property {
    name  = "htmlConfig"
    type  = "string"
    value = var.jumioconnector_property_html_config
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.jumioconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.jumioconnector_property_icon_url_png
  }
  property {
    name  = "jumioIdTypes"
    type  = "string"
    value = var.jumioconnector_property_jumio_id_types
  }
  property {
    name  = "locale"
    type  = "string"
    value = var.jumioconnector_property_locale
  }
  property {
    name  = "merchantReportingCriteria"
    type  = "string"
    value = var.jumioconnector_property_merchant_reporting_criteria
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.jumioconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.jumioconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.jumioconnector_property_screen2_config
  }
  property {
    name  = "screen3Config"
    type  = "string"
    value = var.jumioconnector_property_screen3_config
  }
  property {
    name  = "screen4Config"
    type  = "string"
    value = var.jumioconnector_property_screen4_config
  }
  property {
    name  = "screen5Config"
    type  = "string"
    value = var.jumioconnector_property_screen5_config
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.jumioconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.jumioconnector_property_show_cred_added_via
  }
  property {
    name  = "title"
    type  = "string"
    value = var.jumioconnector_property_title
  }
  property {
    name  = "tokenLifetimeInMinutes"
    type  = "number"
    value = var.jumioconnector_property_token_lifetime_in_minutes
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.jumioconnector_property_tool_tip
  }
  property {
    name  = "useCustomScreens"
    type  = "string"
    value = var.jumioconnector_property_use_custom_screens
  }
  property {
    name  = "workflowId"
    type  = "string"
    value = var.jumioconnector_property_workflow_id
  }
}
