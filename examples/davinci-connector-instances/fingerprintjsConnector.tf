resource "pingone_davinci_connector_instance" "fingerprintjsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "fingerprintjsConnector"
  }
  name = "My awesome fingerprintjsConnector"
  properties = jsonencode({
    "apiToken" = var.fingerprintjsconnector_property_api_token
    "javascriptCdnUrl" = var.fingerprintjsconnector_property_javascript_cdn_url
    "token" = var.fingerprintjsconnector_property_token
  })
}
