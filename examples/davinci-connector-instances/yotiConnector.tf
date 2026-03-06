resource "pingone_davinci_connector_instance" "yotiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "yotiConnector"
  }
  name = "My awesome yotiConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
