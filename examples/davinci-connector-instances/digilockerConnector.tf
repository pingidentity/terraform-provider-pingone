resource "pingone_davinci_connector_instance" "digilockerConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "digilockerConnector"
  }
  name = "My awesome digilockerConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
