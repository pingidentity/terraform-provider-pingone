resource "pingone_davinci_connector_instance" "idrampOidcConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idrampOidcConnector"
  }
  name = "My awesome idrampOidcConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
