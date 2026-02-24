resource "pingone_davinci_connector_instance" "bitbucketIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "bitbucketIdpConnector"
  }
  name = "My awesome bitbucketIdpConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
