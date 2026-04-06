resource "pingone_davinci_connector_instance" "connectorSpotify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpotify"
  }
  name = "My awesome connectorSpotify"
  property {
    name  = "accessToken"
    type  = "string"
    value = var.connectorspotify_property_access_token
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorspotify_property_auth_type
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorspotify_property_body
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorspotify_property_button
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorspotify_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorspotify_property_headers
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorspotify_property_method
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorspotify_property_query_parameters
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorspotify_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorspotify_property_skip_button_press
  }
}
