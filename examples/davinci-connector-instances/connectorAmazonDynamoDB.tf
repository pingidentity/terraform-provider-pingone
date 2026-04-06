resource "pingone_davinci_connector_instance" "connectorAmazonDynamoDB" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonDynamoDB"
  }
  name = "My awesome connectorAmazonDynamoDB"
  property {
    name  = "attrDefinitions"
    type  = "string"
    value = var.connectoramazondynamodb_property_attr_definitions
  }
  property {
    name  = "attrKeySchema"
    type  = "string"
    value = var.connectoramazondynamodb_property_attr_key_schema
  }
  property {
    name  = "awsAccessKey"
    type  = "string"
    value = var.connectoramazondynamodb_property_aws_access_key
  }
  property {
    name  = "awsAccessSecret"
    type  = "string"
    value = var.connectoramazondynamodb_property_aws_access_secret
  }
  property {
    name  = "awsRegion"
    type  = "string"
    value = "eu-west-1"
  }
  property {
    name  = "condtionalExpression"
    type  = "string"
    value = var.connectoramazondynamodb_property_condtional_expression
  }
  property {
    name  = "convertClassInstanceToMap"
    type  = "string"
    value = var.connectoramazondynamodb_property_convert_class_instance_to_map
  }
  property {
    name  = "convertEmptyValues"
    type  = "string"
    value = var.connectoramazondynamodb_property_convert_empty_values
  }
  property {
    name  = "expressionAttributeNames"
    type  = "string"
    value = var.connectoramazondynamodb_property_expression_attribute_names
  }
  property {
    name  = "expressionAttributeValues"
    type  = "string"
    value = var.connectoramazondynamodb_property_expression_attribute_values
  }
  property {
    name  = "filterExpression"
    type  = "string"
    value = var.connectoramazondynamodb_property_filter_expression
  }
  property {
    name  = "indexName"
    type  = "string"
    value = var.connectoramazondynamodb_property_index_name
  }
  property {
    name  = "itemsJSON"
    type  = "string"
    value = var.connectoramazondynamodb_property_items_json
  }
  property {
    name  = "itemsKV"
    type  = "string"
    value = var.connectoramazondynamodb_property_items_kv
  }
  property {
    name  = "keyConditionExpression"
    type  = "string"
    value = var.connectoramazondynamodb_property_key_condition_expression
  }
  property {
    name  = "keyJSON"
    type  = "string"
    value = var.connectoramazondynamodb_property_key_json
  }
  property {
    name  = "keyKV"
    type  = "string"
    value = var.connectoramazondynamodb_property_key_kv
  }
  property {
    name  = "limit"
    type  = "string"
    value = var.connectoramazondynamodb_property_limit
  }
  property {
    name  = "projectionExpression"
    type  = "string"
    value = var.connectoramazondynamodb_property_projection_expression
  }
  property {
    name  = "removeUndefinedValues"
    type  = "string"
    value = var.connectoramazondynamodb_property_remove_undefined_values
  }
  property {
    name  = "sdkCommandName"
    type  = "string"
    value = var.connectoramazondynamodb_property_sdk_command_name
  }
  property {
    name  = "sdkParameters"
    type  = "string"
    value = var.connectoramazondynamodb_property_sdk_parameters
  }
  property {
    name  = "tableName"
    type  = "string"
    value = var.connectoramazondynamodb_property_table_name
  }
  property {
    name  = "updateExpression"
    type  = "string"
    value = var.connectoramazondynamodb_property_update_expression
  }
  property {
    name  = "wrapNumbers"
    type  = "string"
    value = var.connectoramazondynamodb_property_wrap_numbers
  }
}
