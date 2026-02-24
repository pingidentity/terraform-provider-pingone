resource "pingone_davinci_connector_instance" "connectorSvipe" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSvipe"
  }
  name = "My awesome connectorSvipe"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
