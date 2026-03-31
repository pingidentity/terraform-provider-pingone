resource "pingone_davinci_connector_instance" "fingerprintjsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "fingerprintjsConnector"
  }
  name = "My awesome fingerprintjsConnector"
  property {
    name  = "apiToken"
    type  = "string"
    value = var.fingerprintjsconnector_property_api_token
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.fingerprintjsconnector_property_javascript_cdn_url
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.fingerprintjsconnector_property_next_event
  }
  property {
    name  = "token"
    type  = "string"
    value = var.fingerprintjsconnector_property_token
  }
  property {
    name  = "visitorId"
    type  = "string"
    value = var.fingerprintjsconnector_property_visitor_id
  }
}
