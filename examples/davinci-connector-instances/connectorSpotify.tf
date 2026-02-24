resource "pingone_davinci_connector_instance" "connectorSpotify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpotify"
  }
  name = "My awesome connectorSpotify"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
