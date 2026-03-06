resource "pingone_davinci_connector_instance" "iproovConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iproovConnector"
  }
  name = "My awesome iproovConnector"
  properties = jsonencode({
    "allowLandscape" = var.iproovconnector_property_allow_landscape
    "apiKey" = var.iproovconnector_property_api_key
    "authDescription" = var.iproovconnector_property_auth_description
    "baseUrl" = var.iproovconnector_property_base_url
    "color1" = var.iproovconnector_property_color1
    "color2" = var.iproovconnector_property_color2
    "color3" = var.iproovconnector_property_color3
    "color4" = var.iproovconnector_property_color4
    "connectorName" = var.iproovconnector_property_connector_name
    "customTitle" = var.iproovconnector_property_custom_title
    "description" = var.iproovconnector_property_description
    "details1" = var.iproovconnector_property_details1
    "details2" = var.iproovconnector_property_details2
    "enableCameraSelector" = var.iproovconnector_property_enable_camera_selector
    "iconUrl" = var.iproovconnector_property_icon_url
    "iconUrlPng" = var.iproovconnector_property_icon_url_png
    "javascriptCSSUrl" = var.javascript_css_url
    "javascriptCdnUrl" = var.iproovconnector_property_javascript_cdn_url
    "kioskMode" = var.iproovconnector_property_kiosk_mode
    "logo" = var.iproovconnector_property_logo
    "password" = var.iproovconnector_property_password
    "secret" = var.iproovconnector_property_secret
    "showCountdown" = var.iproovconnector_property_show_countdown
    "showCredAddedOn" = var.iproovconnector_property_show_cred_added_on
    "showCredAddedVia" = var.iproovconnector_property_show_cred_added_via
    "startScreenTitle" = var.iproovconnector_property_start_screen_title
    "title" = var.iproovconnector_property_title
    "toolTip" = var.iproovconnector_property_tool_tip
    "username" = var.iproovconnector_property_username
  })
}
