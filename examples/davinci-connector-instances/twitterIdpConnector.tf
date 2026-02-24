resource "pingone_davinci_connector_instance" "twitterIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "twitterIdpConnector"
  }
  name = "My awesome twitterIdpConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
