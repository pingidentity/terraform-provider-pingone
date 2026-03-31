resource "pingone_davinci_connector_instance" "variablesConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "variablesConnector"
  }
  name = "My awesome variablesConnector"
  property {
    name  = "incrementCounter"
    type  = "string"
    value = var.variablesconnector_property_increment_counter
  }
  property {
    name  = "locale"
    type  = "string"
    value = var.variablesconnector_property_locale
  }
  property {
    name  = "saveCompanyVariables"
    type  = "string"
    value = var.variablesconnector_property_save_company_variables
  }
  property {
    name  = "saveFlowVariables"
    type  = "string"
    value = var.variablesconnector_property_save_flow_variables
  }
  property {
    name  = "saveVariables"
    type  = "string"
    value = var.variablesconnector_property_save_variables
  }
  property {
    name  = "variable"
    type  = "string"
    value = var.variablesconnector_property_variable
  }
}
