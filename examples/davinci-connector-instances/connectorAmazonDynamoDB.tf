resource "pingone_davinci_connector_instance" "connectorAmazonDynamoDB" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonDynamoDB"
  }
  name = "My awesome connectorAmazonDynamoDB"
  properties = jsonencode({
    "awsAccessKey" = var.connectoramazondynamodb_property_aws_access_key
    "awsAccessSecret" = var.connectoramazondynamodb_property_aws_access_secret
    "awsRegion" = "eu-west-1"
  })
}
