resource "pingone_davinci_connector_instance" "oneTrustConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "oneTrustConnector"
  }
  name = "My awesome oneTrustConnector"
  properties = jsonencode({
    "clientId" = var.onetrustconnector_property_client_id
    "clientSecret" = var.onetrustconnector_property_client_secret
  })
}
