resource "pingone_davinci_connector_instance" "githubIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "githubIdpConnector"
  }
  name = "My awesome githubIdpConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
