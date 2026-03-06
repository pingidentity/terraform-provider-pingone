resource "pingone_davinci_connector_instance" "samlIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlIdpConnector"
  }
  name = "My awesome samlIdpConnector"
  properties = jsonencode({
    "saml" = jsonencode({})
  })
}
