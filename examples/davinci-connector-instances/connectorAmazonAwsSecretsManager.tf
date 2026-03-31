resource "pingone_davinci_connector_instance" "connectorAmazonAwsSecretsManager" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonAwsSecretsManager"
  }
  name = "My awesome connectorAmazonAwsSecretsManager"
  property {
    name  = "accessKeyId"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_access_key_id
  }
  property {
    name  = "kmsKeyId"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_kms_key_id
  }
  property {
    name  = "overwriteIfExists"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_overwrite_if_exists
  }
  property {
    name  = "recoverFromSoftDelete"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_recover_from_soft_delete
  }
  property {
    name  = "region"
    type  = "string"
    value = "eu-west-1"
  }
  property {
    name  = "secondaryAccessKeyId"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secondary_access_key_id
  }
  property {
    name  = "secondaryRegion"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secondary_region
  }
  property {
    name  = "secondarySecretAccessKey"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secondary_secret_access_key
  }
  property {
    name  = "secretAccessKey"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secret_access_key
  }
  property {
    name  = "secretDescription"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secret_description
  }
  property {
    name  = "secretName"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secret_name
  }
  property {
    name  = "secretValue"
    type  = "string"
    value = var.connectoramazonawssecretsmanager_property_secret_value
  }
}
