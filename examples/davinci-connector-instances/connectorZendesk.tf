resource "pingone_davinci_connector_instance" "connectorZendesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZendesk"
  }
  name = "My awesome connectorZendesk"
  properties = jsonencode({
    "apiToken" = var.connectorzendesk_property_api_token
    "emailUsername" = var.connectorzendesk_property_email_username
    "subdomain" = var.connectorzendesk_property_subdomain
  })
}
