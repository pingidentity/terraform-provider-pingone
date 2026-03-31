resource "pingone_davinci_connector_instance" "mailchainConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mailchainConnector"
  }
  name = "My awesome mailchainConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.mailchainconnector_property_api_key
  }
  property {
    name  = "bodyCredential"
    type  = "string"
    value = var.mailchainconnector_property_body_credential
  }
  property {
    name  = "bodyPresentation"
    type  = "string"
    value = var.mailchainconnector_property_body_presentation
  }
  property {
    name  = "did"
    type  = "string"
    value = var.mailchainconnector_property_did
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.mailchainconnector_property_endpoint
  }
  property {
    name  = "endpointVerify"
    type  = "string"
    value = var.mailchainconnector_property_endpoint_verify
  }
  property {
    name  = "optionalParams"
    type  = "string"
    value = var.mailchainconnector_property_optional_params
  }
  property {
    name  = "version"
    type  = "string"
    value = var.mailchainconnector_property_version
  }
}
