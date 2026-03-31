resource "pingone_davinci_connector_instance" "stringsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "stringsConnector"
  }
  name = "My awesome stringsConnector"
  property {
    name  = "concatenateDelimiter"
    type  = "string"
    value = var.stringsconnector_property_concatenate_delimiter
  }
  property {
    name  = "concatenateInput"
    type  = "string"
    value = var.stringsconnector_property_concatenate_input
  }
  property {
    name  = "dataTypeTo"
    type  = "string"
    value = var.stringsconnector_property_data_type_to
  }
  property {
    name  = "decodeMethod"
    type  = "string"
    value = var.stringsconnector_property_decode_method
  }
  property {
    name  = "delimiter"
    type  = "string"
    value = var.stringsconnector_property_delimiter
  }
  property {
    name  = "encodeMethod"
    type  = "string"
    value = var.stringsconnector_property_encode_method
  }
  property {
    name  = "finalDelimiter"
    type  = "string"
    value = var.stringsconnector_property_final_delimiter
  }
  property {
    name  = "finder"
    type  = "string"
    value = var.stringsconnector_property_finder
  }
  property {
    name  = "inputValue"
    type  = "string"
    value = var.stringsconnector_property_input_value
  }
  property {
    name  = "isAlphaNumeric"
    type  = "string"
    value = var.stringsconnector_property_is_alpha_numeric
  }
  property {
    name  = "isRegex"
    type  = "string"
    value = var.stringsconnector_property_is_regex
  }
  property {
    name  = "length"
    type  = "number"
    value = var.stringsconnector_property_length
  }
  property {
    name  = "newToken"
    type  = "string"
    value = var.stringsconnector_property_new_token
  }
  property {
    name  = "oldToken"
    type  = "string"
    value = var.stringsconnector_property_old_token
  }
  property {
    name  = "originalValue"
    type  = "string"
    value = var.stringsconnector_property_original_value
  }
  property {
    name  = "replaceMode"
    type  = "string"
    value = var.stringsconnector_property_replace_mode
  }
  property {
    name  = "shouldTrim"
    type  = "string"
    value = var.stringsconnector_property_should_trim
  }
}
