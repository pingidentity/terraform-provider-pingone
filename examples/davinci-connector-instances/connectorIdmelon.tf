resource "pingone_davinci_connector_instance" "connectorIdmelon" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdmelon"
  }
  name = "My awesome connectorIdmelon"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
