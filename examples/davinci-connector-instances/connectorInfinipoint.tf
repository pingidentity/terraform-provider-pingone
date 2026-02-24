resource "pingone_davinci_connector_instance" "connectorInfinipoint" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorInfinipoint"
  }
  name = "My awesome connectorInfinipoint"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
