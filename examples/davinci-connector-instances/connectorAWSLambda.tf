resource "pingone_davinci_connector_instance" "connectorAWSLambda" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAWSLambda"
  }
  name = "My awesome connectorAWSLambda"
  properties = jsonencode({
    "accessKeyId" = var.connectorawslambda_property_access_key_id
    "region" = "eu-west-1"
    "secretAccessKey" = var.connectorawslambda_property_secret_access_key
  })
}
