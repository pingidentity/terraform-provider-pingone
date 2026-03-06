resource "pingone_davinci_connector_instance" "duoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "duoConnector"
  }
  name = "My awesome duoConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
