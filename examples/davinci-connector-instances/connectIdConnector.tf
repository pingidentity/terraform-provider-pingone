resource "pingone_davinci_connector_instance" "connectIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectIdConnector"
  }
  name = "My awesome connectIdConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
