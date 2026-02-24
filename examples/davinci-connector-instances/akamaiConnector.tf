resource "pingone_davinci_connector_instance" "akamaiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "akamaiConnector"
  }
  name = "My awesome akamaiConnector"
  properties = jsonencode({
    "customAuth" = var.akamaiconnector_property_custom_auth
  })
}
