resource "pingone_davinci_connector_instance" "neoeyedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "neoeyedConnector"
  }
  name = "My awesome neoeyedConnector"
  properties = jsonencode({
    "appKey" = var.neoeyedconnector_property_app_key
    "javascriptCdnUrl" = var.neoeyedconnector_property_javascript_cdn_url
  })
}
