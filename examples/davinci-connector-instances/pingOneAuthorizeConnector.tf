resource "pingone_davinci_connector_instance" "pingOneAuthorizeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthorizeConnector"
  }
  name = "My awesome pingOneAuthorizeConnector"
  properties = jsonencode({
    "clientId" = var.pingoneauthorizeconnector_property_client_id
    "clientSecret" = var.pingoneauthorizeconnector_property_client_secret
    "endpointURL" = var.endpoint_url
  })
}
