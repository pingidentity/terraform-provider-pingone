resource "pingone_davinci_connector_instance" "neoeyedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "neoeyedConnector"
  }
  name = "My awesome neoeyedConnector"
  property {
    name  = "appKey"
    type  = "string"
    value = var.neoeyedconnector_property_app_key
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.neoeyedconnector_property_javascript_cdn_url
  }
  property {
    name  = "loadingText"
    type  = "string"
    value = var.neoeyedconnector_property_loading_text
  }
}
