resource "pingone_davinci_connector_instance" "veriffConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "veriffConnector"
  }
  name = "My awesome veriffConnector"
  property {
    name  = "access_token"
    type  = "string"
    value = var.veriffconnector_property_access_token
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.veriffconnector_property_base_url
  }
  property {
    name  = "body"
    type  = "string"
    value = var.veriffconnector_property_body
  }
  property {
    name  = "dateOfBirth"
    type  = "string"
    value = var.veriffconnector_property_date_of_birth
  }
  property {
    name  = "documentBack"
    type  = "string"
    value = var.veriffconnector_property_document_back
  }
  property {
    name  = "documentFace"
    type  = "string"
    value = var.veriffconnector_property_document_face
  }
  property {
    name  = "documentFront"
    type  = "string"
    value = var.veriffconnector_property_document_front
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.veriffconnector_property_endpoint
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.veriffconnector_property_first_name
  }
  property {
    name  = "gender"
    type  = "string"
    value = var.veriffconnector_property_gender
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.veriffconnector_property_headers
  }
  property {
    name  = "identificationNumber"
    type  = "string"
    value = var.veriffconnector_property_identification_number
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.veriffconnector_property_last_name
  }
  property {
    name  = "method"
    type  = "string"
    value = var.veriffconnector_property_method
  }
  property {
    name  = "password"
    type  = "string"
    value = var.veriffconnector_property_password
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.veriffconnector_property_query_parameters
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.veriffconnector_property_session_id
  }
  property {
    name  = "vendorData"
    type  = "string"
    value = var.veriffconnector_property_vendor_data
  }
}
