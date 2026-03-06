resource "pingone_davinci_connector_instance" "genericConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "genericConnector"
  }
  name = "My awesome genericConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
