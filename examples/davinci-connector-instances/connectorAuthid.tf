resource "pingone_davinci_connector_instance" "connectorAuthid" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAuthid"
  }
  name = "My awesome connectorAuthid"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
