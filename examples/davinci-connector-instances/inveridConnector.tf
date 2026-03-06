resource "pingone_davinci_connector_instance" "inveridConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "inveridConnector"
  }
  name = "My awesome inveridConnector"
  properties = jsonencode({
    "getApiKey" = var.inveridconnector_property_get_api_key
    "host" = var.inveridconnector_property_host
    "postApiKey" = var.inveridconnector_property_post_api_key
    "skWebhookUri" = var.inveridconnector_property_sk_webhook_uri
    "timeToLive" = var.inveridconnector_property_time_to_live
  })
}
