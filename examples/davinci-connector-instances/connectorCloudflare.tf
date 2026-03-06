resource "pingone_davinci_connector_instance" "connectorCloudflare" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCloudflare"
  }
  name = "My awesome connectorCloudflare"
  properties = jsonencode({
    "accountId" = var.connectorcloudflare_property_account_id
    "apiToken" = var.connectorcloudflare_property_api_token
  })
}
