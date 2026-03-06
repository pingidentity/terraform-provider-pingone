resource "pingone_davinci_connector_instance" "securIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "securIdConnector"
  }
  name = "My awesome securIdConnector"
  properties = jsonencode({
    "apiUrl" = var.securidconnector_property_api_url
    "clientKey" = var.securidconnector_property_client_key
  })
}
