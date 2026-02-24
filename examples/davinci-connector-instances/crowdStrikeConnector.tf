resource "pingone_davinci_connector_instance" "crowdStrikeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "crowdStrikeConnector"
  }
  name = "My awesome crowdStrikeConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "clientId" = var.crowdstrikeconnector_property_client_id
    "clientSecret" = var.crowdstrikeconnector_property_client_secret
  })
}
