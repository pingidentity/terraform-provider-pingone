resource "pingone_davinci_connector_instance" "typingdnaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "typingdnaConnector"
  }
  name = "My awesome typingdnaConnector"
  properties = jsonencode({
    "customAuth" = var.typingdnaconnector_property_custom_auth
  })
}
