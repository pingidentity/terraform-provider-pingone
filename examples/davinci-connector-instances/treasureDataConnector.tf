resource "pingone_davinci_connector_instance" "treasureDataConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "treasureDataConnector"
  }
  name = "My awesome treasureDataConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.treasuredataconnector_property_api_key
  }
  property {
    name  = "attributes"
    type  = "string"
    value = var.treasuredataconnector_property_attributes
  }
  property {
    name  = "audienceId"
    type  = "string"
    value = var.treasuredataconnector_property_audience_id
  }
  property {
    name  = "body"
    type  = "string"
    value = var.treasuredataconnector_property_body
  }
  property {
    name  = "database"
    type  = "string"
    value = var.treasuredataconnector_property_database
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.treasuredataconnector_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.treasuredataconnector_property_headers
  }
  property {
    name  = "id"
    type  = "string"
    value = var.treasuredataconnector_property_id
  }
  property {
    name  = "method"
    type  = "string"
    value = var.treasuredataconnector_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.treasuredataconnector_property_query_parameters
  }
  property {
    name  = "recordData"
    type  = "string"
    value = var.treasuredataconnector_property_record_data
  }
  property {
    name  = "regionRecord"
    type  = "string"
    value = var.treasuredataconnector_property_region_record
  }
  property {
    name  = "regionToken"
    type  = "string"
    value = var.treasuredataconnector_property_region_token
  }
  property {
    name  = "relationships"
    type  = "string"
    value = var.treasuredataconnector_property_relationships
  }
  property {
    name  = "table"
    type  = "string"
    value = var.treasuredataconnector_property_table
  }
  property {
    name  = "type"
    type  = "string"
    value = var.treasuredataconnector_property_type
  }
}
