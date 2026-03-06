resource "pingone_davinci_connector_instance" "privateidConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "privateidConnector"
  }
  name = "My awesome privateidConnector"
  properties = jsonencode({
    "customAuth" = var.privateidconnector_property_custom_auth
  })
}
