resource "pingone_davinci_connector_instance" "connectorMicrosoftIntune" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMicrosoftIntune"
  }
  name = "My awesome connectorMicrosoftIntune"
  properties = jsonencode({
    "clientId" = var.connectormicrosoftintune_property_client_id
    "clientSecret" = var.connectormicrosoftintune_property_client_secret
    "grantType" = var.connectormicrosoftintune_property_grant_type
    "scope" = var.connectormicrosoftintune_property_scope
    "tenant" = var.connectormicrosoftintune_property_tenant
  })
}
