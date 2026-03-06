resource "pingone_davinci_connector_instance" "connectorAmazonAwsSecretsManager" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonAwsSecretsManager"
  }
  name = "My awesome connectorAmazonAwsSecretsManager"
  properties = jsonencode({
    "accessKeyId" = var.connectoramazonawssecretsmanager_property_access_key_id
    "region" = "eu-west-1"
    "secretAccessKey" = var.connectoramazonawssecretsmanager_property_secret_access_key
  })
}
