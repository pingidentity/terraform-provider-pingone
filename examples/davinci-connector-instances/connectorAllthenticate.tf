resource "pingone_davinci_connector_instance" "connectorAllthenticate" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAllthenticate"
  }
  name = "My awesome connectorAllthenticate"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
