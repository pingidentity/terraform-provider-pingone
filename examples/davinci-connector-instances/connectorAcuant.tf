resource "pingone_davinci_connector_instance" "connectorAcuant" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAcuant"
  }
  name = "My awesome connectorAcuant"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
