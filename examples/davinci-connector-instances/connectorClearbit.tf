resource "pingone_davinci_connector_instance" "connectorClearbit" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClearbit"
  }
  name = "My awesome connectorClearbit"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorclearbit_property_api_key
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorclearbit_property_body
  }
  property {
    name  = "domainName"
    type  = "string"
    value = var.connectorclearbit_property_domain_name
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorclearbit_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorclearbit_property_endpoint
  }
  property {
    name  = "fullName"
    type  = "string"
    value = var.connectorclearbit_property_full_name
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorclearbit_property_headers
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.connectorclearbit_property_ip_address
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorclearbit_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorclearbit_property_query_parameters
  }
  property {
    name  = "riskApiVersion"
    type  = "string"
    value = var.connectorclearbit_property_risk_api_version
  }
  property {
    name  = "version"
    type  = "string"
    value = var.connectorclearbit_property_version
  }
}
