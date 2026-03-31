resource "pingone_davinci_connector_instance" "connectorMicrosoftIntune" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMicrosoftIntune"
  }
  name = "My awesome connectorMicrosoftIntune"
  property {
    name  = "body"
    type  = "string"
    value = var.connectormicrosoftintune_property_body
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.connectormicrosoftintune_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectormicrosoftintune_property_client_secret
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectormicrosoftintune_property_endpoint
  }
  property {
    name  = "grantType"
    type  = "string"
    value = var.connectormicrosoftintune_property_grant_type
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectormicrosoftintune_property_headers
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectormicrosoftintune_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectormicrosoftintune_property_query_parameters
  }
  property {
    name  = "scope"
    type  = "string"
    value = var.connectormicrosoftintune_property_scope
  }
  property {
    name  = "serialNumber"
    type  = "string"
    value = var.connectormicrosoftintune_property_serial_number
  }
  property {
    name  = "tenant"
    type  = "string"
    value = var.connectormicrosoftintune_property_tenant
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.connectormicrosoftintune_property_user_id
  }
  property {
    name  = "userPrincipalName"
    type  = "string"
    value = var.connectormicrosoftintune_property_user_principal_name
  }
}
