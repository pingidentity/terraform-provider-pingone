resource "pingone_davinci_connector_instance" "connectorHello" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHello"
  }
  name = "My awesome connectorHello"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
