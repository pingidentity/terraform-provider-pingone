resource "pingone_davinci_connector_instance" "jumioConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jumioConnector"
  }
  name = "My awesome jumioConnector"
  properties = jsonencode({
    "apiKey" = var.jumioconnector_property_api_key
    "authDescription" = var.jumioconnector_property_auth_description
    "authUrl" = var.jumioconnector_property_auth_url
    "authorizationTokenLifetime" = var.jumioconnector_property_authorization_token_lifetime
    "baseColor" = var.jumioconnector_property_base_color
    "bgColor" = var.jumioconnector_property_bg_color
    "callbackUrl" = var.jumioconnector_property_callback_url
    "clientSecret" = var.jumioconnector_property_client_secret
    "connectorName" = var.jumioconnector_property_connector_name
    "description" = var.jumioconnector_property_description
    "details1" = var.jumioconnector_property_details1
    "details2" = var.jumioconnector_property_details2
    "doNotShowInIframe" = var.jumioconnector_property_do_not_show_in_iframe
    "docVerificationUrl" = var.jumioconnector_property_doc_verification_url
    "headerImageUrl" = var.jumioconnector_property_header_image_url
    "iconUrl" = var.jumioconnector_property_icon_url
    "iconUrlPng" = var.jumioconnector_property_icon_url_png
    "locale" = var.jumioconnector_property_locale
    "showCredAddedOn" = var.jumioconnector_property_show_cred_added_on
    "showCredAddedVia" = var.jumioconnector_property_show_cred_added_via
    "title" = var.jumioconnector_property_title
    "toolTip" = var.jumioconnector_property_tool_tip
  })
}
