resource "pingone_davinci_connector_instance" "incodeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "incodeConnector"
  }
  name = "My awesome incodeConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
