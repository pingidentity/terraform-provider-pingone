resource "pingone_davinci_connector_instance" "telesignConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "telesignConnector"
  }
  name = "My awesome telesignConnector"
  properties = jsonencode({
    "authDescription" = var.telesignconnector_property_auth_description
    "connectorName" = var.telesignconnector_property_connector_name
    "description" = var.telesignconnector_property_description
    "details1" = var.telesignconnector_property_details1
    "details2" = var.telesignconnector_property_details2
    "iconUrl" = var.telesignconnector_property_icon_url
    "iconUrlPng" = var.telesignconnector_property_icon_url_png
    "password" = var.telesignconnector_property_password
    "providerName" = var.telesignconnector_property_provider_name
    "showCredAddedOn" = var.telesignconnector_property_show_cred_added_on
    "showCredAddedVia" = var.telesignconnector_property_show_cred_added_via
    "title" = var.telesignconnector_property_title
    "toolTip" = var.telesignconnector_property_tool_tip
    "username" = var.telesignconnector_property_username
  })
}
