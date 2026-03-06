resource "pingone_davinci_connector_instance" "sentilinkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sentilinkConnector"
  }
  name = "My awesome sentilinkConnector"
  properties = jsonencode({
    "account" = var.sentilinkconnector_property_account
    "apiUrl" = var.sentilinkconnector_property_api_url
    "javascriptCdnUrl" = var.sentilinkconnector_property_javascript_cdn_url
    "token" = var.sentilinkconnector_property_token
  })
}
