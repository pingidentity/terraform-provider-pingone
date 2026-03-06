resource "pingone_davinci_connector_instance" "connectorBadge" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBadge"
  }
  name = "My awesome connectorBadge"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
