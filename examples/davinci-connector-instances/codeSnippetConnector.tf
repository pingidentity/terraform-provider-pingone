resource "pingone_davinci_connector_instance" "codeSnippetConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "codeSnippetConnector"
  }
  name = "My awesome codeSnippetConnector"
  property {
    name  = "code"
    type  = "string"
    value = var.codesnippetconnector_property_code
  }
  property {
    name  = "functionArgumentList"
    type  = "string"
    value = var.codesnippetconnector_property_function_argument_list
  }
  property {
    name  = "inputSchema"
    type  = "string"
    value = var.codesnippetconnector_property_input_schema
  }
  property {
    name  = "outputSchema"
    type  = "string"
    value = var.codesnippetconnector_property_output_schema
  }
}
