resource "pingone_davinci_connector_instance" "connectorTruid" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTruid"
  }
  name = "My awesome connectorTruid"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
