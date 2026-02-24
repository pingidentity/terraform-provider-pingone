resource "pingone_davinci_connector_instance" "idDatawebConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idDatawebConnector"
  }
  name = "My awesome idDatawebConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
