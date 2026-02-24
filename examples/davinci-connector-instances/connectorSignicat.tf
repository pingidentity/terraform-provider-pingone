resource "pingone_davinci_connector_instance" "connectorSignicat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSignicat"
  }
  name = "My awesome connectorSignicat"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
