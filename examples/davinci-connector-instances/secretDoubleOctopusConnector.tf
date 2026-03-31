resource "pingone_davinci_connector_instance" "secretDoubleOctopusConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "secretDoubleOctopusConnector"
  }
  name = "My awesome secretDoubleOctopusConnector"
  property {
    name  = "apiToken"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_api_token
  }
  property {
    name  = "authId"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_auth_id
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_base_url
  }
  property {
    name  = "message"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_message
  }
  property {
    name  = "password"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_password
  }
  property {
    name  = "serviceId"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_service_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_username
  }
  property {
    name  = "x509Certificate"
    type  = "string"
    value = var.secretdoubleoctopusconnector_property_x509_certificate
  }
}
