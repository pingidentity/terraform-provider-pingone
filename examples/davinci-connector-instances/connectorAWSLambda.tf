resource "pingone_davinci_connector_instance" "connectorAWSLambda" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAWSLambda"
  }
  name = "My awesome connectorAWSLambda"
  property {
    name  = "accessKeyId"
    type  = "string"
    value = var.connectorawslambda_property_access_key_id
  }
  property {
    name  = "functionName"
    type  = "string"
    value = var.connectorawslambda_property_function_name
  }
  property {
    name  = "payload"
    type  = "string"
    value = var.connectorawslambda_property_payload
  }
  property {
    name  = "raw"
    type  = "string"
    value = var.connectorawslambda_property_raw
  }
  property {
    name  = "region"
    type  = "string"
    value = "eu-west-1"
  }
  property {
    name  = "secretAccessKey"
    type  = "string"
    value = var.connectorawslambda_property_secret_access_key
  }
}
