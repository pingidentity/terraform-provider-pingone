resource "pingone_davinci_connector_instance" "kbaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kbaConnector"
  }
  name = "My awesome kbaConnector"
  properties = jsonencode({
    "authDescription" = var.kbaconnector_property_auth_description
    "connectorName" = var.kbaconnector_property_connector_name
    "description" = var.kbaconnector_property_description
    "details1" = var.kbaconnector_property_details1
    "details2" = var.kbaconnector_property_details2
    "formFieldsList" = var.kbaconnector_property_form_fields_list
    "iconUrl" = var.kbaconnector_property_icon_url
    "iconUrlPng" = var.kbaconnector_property_icon_url_png
    "showCredAddedOn" = var.kbaconnector_property_show_cred_added_on
    "showCredAddedVia" = var.kbaconnector_property_show_cred_added_via
    "title" = var.kbaconnector_property_title
    "toolTip" = var.kbaconnector_property_tool_tip
  })
}
