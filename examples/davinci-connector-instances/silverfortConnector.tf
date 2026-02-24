resource "pingone_davinci_connector_instance" "silverfortConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "silverfortConnector"
  }
  name = "My awesome silverfortConnector"
  properties = jsonencode({
    "apiKey" = var.silverfortconnector_property_api_key
    "appUserSecret" = var.silverfortconnector_property_app_user_secret
    "consoleApi" = var.silverfortconnector_property_console_api
  })
}
