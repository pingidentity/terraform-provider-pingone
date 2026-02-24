resource "pingone_davinci_connector_instance" "nuanceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nuanceConnector"
  }
  name = "My awesome nuanceConnector"
  properties = jsonencode({
    "authDescription" = var.nuanceconnector_property_auth_description
    "configSetName" = var.nuanceconnector_property_config_set_name
    "connectorName" = var.nuanceconnector_property_connector_name
    "description" = var.nuanceconnector_property_description
    "details1" = var.nuanceconnector_property_details1
    "details2" = var.nuanceconnector_property_details2
    "iconUrl" = var.nuanceconnector_property_icon_url
    "iconUrlPng" = var.nuanceconnector_property_icon_url_png
    "passphrase1" = var.nuanceconnector_property_passphrase1
    "passphrase2" = var.nuanceconnector_property_passphrase2
    "passphrase3" = var.nuanceconnector_property_passphrase3
    "passphrase4" = var.nuanceconnector_property_passphrase4
    "passphrase5" = var.nuanceconnector_property_passphrase5
    "showCredAddedOn" = var.nuanceconnector_property_show_cred_added_on
    "showCredAddedVia" = var.nuanceconnector_property_show_cred_added_via
    "title" = var.nuanceconnector_property_title
    "toolTip" = var.nuanceconnector_property_tool_tip
  })
}
