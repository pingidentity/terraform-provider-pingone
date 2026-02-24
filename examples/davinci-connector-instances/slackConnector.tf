resource "pingone_davinci_connector_instance" "slackConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "slackConnector"
  }
  name = "My awesome slackConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
