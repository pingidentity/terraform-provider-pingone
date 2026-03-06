resource "pingone_davinci_connector_instance" "connectorKeyless" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorKeyless"
  }
  name = "My awesome connectorKeyless"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
