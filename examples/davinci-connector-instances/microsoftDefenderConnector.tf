resource "pingone_davinci_connector_instance" "microsoftDefenderConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftDefenderConnector"
  }
  name = "My awesome microsoftDefenderConnector"
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.microsoftdefenderconnector_property_base_url
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.microsoftdefenderconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.microsoftdefenderconnector_property_client_secret
  }
  property {
    name  = "comment"
    type  = "string"
    value = var.microsoftdefenderconnector_property_comment
  }
  property {
    name  = "deviceId"
    type  = "string"
    value = var.microsoftdefenderconnector_property_device_id
  }
  property {
    name  = "isolationType"
    type  = "string"
    value = var.microsoftdefenderconnector_property_isolation_type
  }
  property {
    name  = "tenantId"
    type  = "string"
    value = var.microsoftdefenderconnector_property_tenant_id
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.microsoftdefenderconnector_property_user_id
  }
}
