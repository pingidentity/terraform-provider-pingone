resource "pingone_davinci_connector_instance" "onfidoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "onfidoConnector"
  }
  name = "My awesome onfidoConnector"
  properties = jsonencode({
    "androidPackageName" = var.onfidoconnector_property_android_package_name
    "apiKey" = var.onfidoconnector_property_api_key
    "authDescription" = var.onfidoconnector_property_auth_description
    "baseUrl" = var.onfidoconnector_property_base_url
    "connectorName" = var.onfidoconnector_property_connector_name
    "customizeSteps" = var.onfidoconnector_property_customize_steps
    "description" = var.onfidoconnector_property_description
    "details1" = var.onfidoconnector_property_details1
    "details2" = var.onfidoconnector_property_details2
    "iOSBundleId" = var.onfidoconnector_property_i_o_s_bundle_id
    "iconUrl" = var.onfidoconnector_property_icon_url
    "iconUrlPng" = var.onfidoconnector_property_icon_url_png
    "javascriptCSSUrl" = var.javascript_css_url
    "javascriptCdnUrl" = var.onfidoconnector_property_javascript_cdn_url
    "language" = var.onfidoconnector_property_language
    "referenceStepsList" = var.onfidoconnector_property_reference_steps_list
    "referrerUrl" = var.onfidoconnector_property_referrer_url
    "retrieveReports" = var.onfidoconnector_property_retrieve_reports
    "shouldCloseOnOverlayClick" = var.onfidoconnector_property_should_close_on_overlay_click
    "showCredAddedOn" = var.onfidoconnector_property_show_cred_added_on
    "showCredAddedVia" = var.onfidoconnector_property_show_cred_added_via
    "stepsList" = var.onfidoconnector_property_steps_list
    "title" = var.onfidoconnector_property_title
    "toolTip" = var.onfidoconnector_property_tool_tip
    "useLanguage" = var.onfidoconnector_property_use_language
    "useModal" = var.onfidoconnector_property_use_modal
    "viewDescriptions" = var.onfidoconnector_property_view_descriptions
    "viewTitle" = var.onfidoconnector_property_view_title
  })
}
