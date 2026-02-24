resource "pingone_davinci_connector_instance" "microsoftTeamsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftTeamsConnector"
  }
  name = "My awesome microsoftTeamsConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
