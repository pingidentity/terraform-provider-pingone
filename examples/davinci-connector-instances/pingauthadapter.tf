resource "pingone_davinci_connector_instance" "pingauthadapter" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingauthadapter"
  }
  name = "My awesome pingauthadapter"
  property {
    name  = "accessTokenClaims"
    type  = "string"
    value = var.pingauthadapter_property_access_token_claims
  }
  property {
    name  = "apiGatewayCredentials"
    type  = "string"
    value = var.pingauthadapter_property_api_gateway_credentials
  }
  property {
    name  = "apiServiceUrl"
    type  = "string"
    value = var.pingauthadapter_property_api_service_url
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.pingauthadapter_property_headers
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.pingauthadapter_property_ip
  }
  property {
    name  = "requestBody"
    type  = "string"
    value = var.pingauthadapter_property_request_body
  }
  property {
    name  = "requestMethod"
    type  = "string"
    value = var.pingauthadapter_property_request_method
  }
  property {
    name  = "serviceUrl"
    type  = "string"
    value = var.pingauthadapter_property_service_url
  }
}
