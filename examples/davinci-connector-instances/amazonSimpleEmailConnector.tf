resource "pingone_davinci_connector_instance" "amazonSimpleEmailConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "amazonSimpleEmailConnector"
  }
  name = "My awesome amazonSimpleEmailConnector"
  properties = jsonencode({
    "awsAccessKey" = var.amazonsimpleemailconnector_property_aws_access_key
    "awsAccessSecret" = var.amazonsimpleemailconnector_property_aws_access_secret
    "awsRegion" = "eu-west-1"
    "from" = "support@bxretail.org"
  })
}
