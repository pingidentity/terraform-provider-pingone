resource "pingone_davinci_connector_instance" "connectorClear" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClear"
  }
  name = "My awesome connectorClear"
  properties = jsonencode({
    "customAuth" = var.connectorclear_property_custom_auth
  })
}
