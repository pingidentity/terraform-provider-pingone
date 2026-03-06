resource "pingone_davinci_connector_instance" "idMeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idMeConnector"
  }
  name = "My awesome idMeConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
