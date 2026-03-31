resource "pingone_davinci_connector_instance" "webhookConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "webhookConnector"
  }
  name = "My awesome webhookConnector"
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.webhookconnector_property_claims_name_value_pairs
  }
  property {
    name  = "urls"
    type  = "string"
    value = var.webhookconnector_property_urls
  }
  property {
    name  = "urlsSelections"
    type  = "string"
    value = var.webhookconnector_property_urls_selections
  }
  property {
    name  = "webhookUrl"
    type  = "string"
    value = var.webhookconnector_property_webhook_url
  }
}
