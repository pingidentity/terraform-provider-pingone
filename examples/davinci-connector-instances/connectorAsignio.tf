resource "pingone_davinci_connector_instance" "connectorAsignio" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAsignio"
  }
  name = "My awesome connectorAsignio"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
