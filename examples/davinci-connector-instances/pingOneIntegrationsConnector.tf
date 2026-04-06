resource "pingone_davinci_connector_instance" "pingOneIntegrationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneIntegrationsConnector"
  }
  name = "My awesome pingOneIntegrationsConnector"
  property {
    name  = "attributeMapping"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_attribute_mapping
  }
  property {
    name  = "attributesRulesList"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_attributes_rules_list
  }
  property {
    name  = "description"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_description
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_next_event
  }
  property {
    name  = "responseType"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_response_type
  }
  property {
    name  = "responseValue"
    type  = "string"
    value = var.pingoneintegrationsconnector_property_response_value
  }
}
