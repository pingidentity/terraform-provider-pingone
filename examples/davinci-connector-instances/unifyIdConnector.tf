resource "pingone_davinci_connector_instance" "unifyIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "unifyIdConnector"
  }
  name = "My awesome unifyIdConnector"
  properties = jsonencode({
    "accountId" = var.unifyidconnector_property_account_id
    "apiKey" = var.unifyidconnector_property_api_key
    "connectorName" = var.unifyidconnector_property_connector_name
    "details1" = var.unifyidconnector_property_details1
    "details2" = var.unifyidconnector_property_details2
    "iconUrl" = var.unifyidconnector_property_icon_url
    "iconUrlPng" = var.unifyidconnector_property_icon_url_png
    "sdkToken" = var.unifyidconnector_property_sdk_token
    "showCredAddedOn" = var.unifyidconnector_property_show_cred_added_on
    "showCredAddedVia" = var.unifyidconnector_property_show_cred_added_via
    "toolTip" = var.unifyidconnector_property_tool_tip
  })
}
