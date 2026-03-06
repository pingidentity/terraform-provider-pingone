resource "pingone_davinci_connector_instance" "codeSnippetConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "codeSnippetConnector"
  }
  name = "My awesome codeSnippetConnector"
  properties = jsonencode({
    "code" = var.codesnippetconnector_property_code
    "inputSchema" = var.codesnippetconnector_property_input_schema
    "outputSchema" = var.codesnippetconnector_property_output_schema
  })
}
