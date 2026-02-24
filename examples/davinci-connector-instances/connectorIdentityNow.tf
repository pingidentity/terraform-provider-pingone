resource "pingone_davinci_connector_instance" "connectorIdentityNow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdentityNow"
  }
  name = "My awesome connectorIdentityNow"
  properties = jsonencode({
    "clientId" = var.connectoridentitynow_property_client_id
    "clientSecret" = var.connectoridentitynow_property_client_secret
    "tenant" = var.connectoridentitynow_property_tenant
  })
}
