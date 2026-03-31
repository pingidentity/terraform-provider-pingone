resource "pingone_davinci_connector_instance" "connectorJamf" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJamf"
  }
  name = "My awesome connectorJamf"
  property {
    name  = "data"
    type  = "string"
    value = var.connectorjamf_property_data
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorjamf_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorjamf_property_headers
  }
  property {
    name  = "identifier"
    type  = "string"
    value = var.connectorjamf_property_identifier
  }
  property {
    name  = "jamfPassword"
    type  = "string"
    value = var.connectorjamf_property_jamf_password
  }
  property {
    name  = "jamfUsername"
    type  = "string"
    value = var.connectorjamf_property_jamf_username
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorjamf_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorjamf_property_query_parameters
  }
  property {
    name  = "searchAttribute"
    type  = "string"
    value = var.connectorjamf_property_search_attribute
  }
  property {
    name  = "serverName"
    type  = "string"
    value = var.connectorjamf_property_server_name
  }
}
