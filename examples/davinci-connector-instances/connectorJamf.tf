resource "pingone_davinci_connector_instance" "connectorJamf" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJamf"
  }
  name = "My awesome connectorJamf"
  properties = jsonencode({
    "jamfPassword" = var.connectorjamf_property_jamf_password
    "jamfUsername" = var.connectorjamf_property_jamf_username
    "serverName" = var.connectorjamf_property_server_name
  })
}
