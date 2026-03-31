resource "pingone_davinci_connector_instance" "connectorTrulioo" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTrulioo"
  }
  name = "My awesome connectorTrulioo"
  property {
    name  = "body"
    type  = "string"
    value = var.connectortrulioo_property_body
  }
  property {
    name  = "clientID"
    type  = "string"
    value = var.connectortrulioo_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectortrulioo_property_client_secret
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectortrulioo_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectortrulioo_property_headers
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectortrulioo_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectortrulioo_property_query_parameters
  }
}
