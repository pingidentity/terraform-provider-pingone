resource "pingone_davinci_connector_instance" "connectorCircleAccess" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCircleAccess"
  }
  name = "My awesome connectorCircleAccess"
  properties = jsonencode({
    "appKey" = var.connectorcircleaccess_property_app_key
    "customAuth" = jsonencode({})
    "loginUrl" = var.connectorcircleaccess_property_login_url
    "readKey" = var.connectorcircleaccess_property_read_key
    "returnToUrl" = var.connectorcircleaccess_property_return_to_url
    "writeKey" = var.connectorcircleaccess_property_write_key
  })
}
