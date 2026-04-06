resource "pingone_davinci_connector_instance" "fapiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "fapiConnector"
  }
  name = "My awesome fapiConnector"
  property {
    name  = "alg"
    type  = "string"
    value = var.fapiconnector_property_alg
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.fapiconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.fapiconnector_property_button
  }
  property {
    name  = "capabilityClientId"
    type  = "string"
    value = var.fapiconnector_property_capability_client_id
  }
  property {
    name  = "capabilityScope"
    type  = "string"
    value = var.fapiconnector_property_capability_scope
  }
  property {
    name  = "claims"
    type  = "string"
    value = var.fapiconnector_property_claims
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.fapiconnector_property_custom_auth
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.fapiconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.fapiconnector_property_skip_button_press
  }
  property {
    name  = "wellknownEndpoint"
    type  = "string"
    value = var.fapiconnector_property_wellknown_endpoint
  }
}
