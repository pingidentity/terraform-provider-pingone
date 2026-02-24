resource "pingone_davinci_connector_instance" "haveIBeenPwnedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "haveIBeenPwnedConnector"
  }
  name = "My awesome haveIBeenPwnedConnector"
  properties = jsonencode({
    "apiKey" = var.haveibeenpwnedconnector_property_api_key
    "apiUrl" = var.haveibeenpwnedconnector_property_api_url
    "userAgent" = var.haveibeenpwnedconnector_property_user_agent
  })
}
