resource "pingone_davinci_connector_instance" "pingOneAuthorizeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthorizeConnector"
  }
  name = "My awesome pingOneAuthorizeConnector"
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_client_secret
  }
  property {
    name  = "code"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_code
  }
  property {
    name  = "endpointURL"
    type  = "string"
    value = var.endpoint_url
  }
  property {
    name  = "parameters"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_parameters
  }
  property {
    name  = "statements"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_statements
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingoneauthorizeconnector_property_user_id
  }
}
