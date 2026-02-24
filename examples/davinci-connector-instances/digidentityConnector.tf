resource "pingone_davinci_connector_instance" "digidentityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "digidentityConnector"
  }
  name = "My awesome digidentityConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
