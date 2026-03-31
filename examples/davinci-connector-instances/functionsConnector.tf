resource "pingone_davinci_connector_instance" "functionsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "functionsConnector"
  }
  name = "My awesome functionsConnector"
  property {
    name  = "caseSensitive"
    type  = "string"
    value = var.functionsconnector_property_case_sensitive
  }
  property {
    name  = "checkNullORUndefined"
    type  = "string"
    value = var.functionsconnector_property_check_null_orundefined
  }
  property {
    name  = "code"
    type  = "string"
    value = var.functionsconnector_property_code
  }
  property {
    name  = "customSaltValue"
    type  = "string"
    value = var.functionsconnector_property_custom_salt_value
  }
  property {
    name  = "digestAlgorithm"
    type  = "string"
    value = var.functionsconnector_property_digest_algorithm
  }
  property {
    name  = "inputContains"
    type  = "string"
    value = var.functionsconnector_property_input_contains
  }
  property {
    name  = "leftValueA"
    type  = "string"
    value = var.functionsconnector_property_left_value_a
  }
  property {
    name  = "message"
    type  = "string"
    value = var.functionsconnector_property_message
  }
  property {
    name  = "outputEncodingFormat"
    type  = "string"
    value = var.functionsconnector_property_output_encoding_format
  }
  property {
    name  = "outputSchema"
    type  = "string"
    value = var.functionsconnector_property_output_schema
  }
  property {
    name  = "rightValueB"
    type  = "string"
    value = var.functionsconnector_property_right_value_b
  }
  property {
    name  = "rightValueC"
    type  = "string"
    value = var.functionsconnector_property_right_value_c
  }
  property {
    name  = "rightValueMultiple"
    type  = "string"
    value = var.functionsconnector_property_right_value_multiple
  }
  property {
    name  = "saltMode"
    type  = "string"
    value = var.functionsconnector_property_salt_mode
  }
  property {
    name  = "type"
    type  = "string"
    value = var.functionsconnector_property_type
  }
  property {
    name  = "variableInputList"
    type  = "string"
    value = var.functionsconnector_property_variable_input_list
  }
}
