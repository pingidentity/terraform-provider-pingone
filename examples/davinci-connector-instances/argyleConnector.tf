resource "pingone_davinci_connector_instance" "argyleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "argyleConnector"
  }
  name = "My awesome argyleConnector"
  properties = jsonencode({
    "apiUrl" = var.argyleconnector_property_api_url
    "clientId" = var.argyleconnector_property_client_id
    "clientSecret" = var.argyleconnector_property_client_secret
    "javascriptWebUrl" = var.argyleconnector_property_javascript_web_url
    "pluginKey" = var.argyleconnector_property_plugin_key
  })
}
