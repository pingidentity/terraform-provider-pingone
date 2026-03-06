resource "pingone_davinci_connector_instance" "singpassLoginConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "singpassLoginConnector"
  }
  name = "My awesome singpassLoginConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
