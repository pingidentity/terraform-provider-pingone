resource "pingone_davinci_connector_instance" "socureConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "socureConnector"
  }
  name = "My awesome socureConnector"
  properties = jsonencode({
    "apiKey" = var.socureconnector_property_api_key
    "baseUrl" = var.socureconnector_property_base_url
    "customApiUrl" = var.socureconnector_property_custom_api_url
    "sdkKey" = var.socureconnector_property_sdk_key
    "skWebhookUri" = var.socureconnector_property_sk_webhook_uri
  })
}
