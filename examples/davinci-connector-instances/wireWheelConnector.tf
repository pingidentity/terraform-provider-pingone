resource "pingone_davinci_connector_instance" "wireWheelConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "wireWheelConnector"
  }
  name = "My awesome wireWheelConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "clientId" = var.wirewheelconnector_property_client_id
    "clientSecret" = var.wirewheelconnector_property_client_secret
    "issuerId" = var.wirewheelconnector_property_issuer_id
  })
}
