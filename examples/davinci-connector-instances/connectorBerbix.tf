resource "pingone_davinci_connector_instance" "connectorBerbix" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBerbix"
  }
  name = "My awesome connectorBerbix"
  property {
    name  = "body"
    type  = "string"
    value = var.connectorberbix_property_body
  }
  property {
    name  = "completionEmail"
    type  = "string"
    value = var.connectorberbix_property_completion_email
  }
  property {
    name  = "consentsToAutomatedFacialRecognition"
    type  = "string"
    value = var.connectorberbix_property_consents_to_automated_facial_recognition
  }
  property {
    name  = "customerUID"
    type  = "string"
    value = var.connectorberbix_property_customer_uid
  }
  property {
    name  = "domainName"
    type  = "string"
    value = var.connectorberbix_property_domain_name
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorberbix_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorberbix_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorberbix_property_headers
  }
  property {
    name  = "idCountry"
    type  = "string"
    value = var.connectorberbix_property_id_country
  }
  property {
    name  = "idType"
    type  = "string"
    value = var.connectorberbix_property_id_type
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorberbix_property_method
  }
  property {
    name  = "path"
    type  = "string"
    value = var.connectorberbix_property_path
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectorberbix_property_phone
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorberbix_property_query_parameters
  }
  property {
    name  = "redirectURL"
    type  = "string"
    value = var.connectorberbix_property_redirect_url
  }
  property {
    name  = "templateKey"
    type  = "string"
    value = var.connectorberbix_property_template_key
  }
  property {
    name  = "token"
    type  = "string"
    value = var.connectorberbix_property_token
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectorberbix_property_username
  }
}
