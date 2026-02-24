resource "pingone_davinci_connector_instance" "kyxstartConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kyxstartConnector"
  }
  name = "My awesome kyxstartConnector"
  properties = jsonencode({
    "clientId" = var.kyxstartconnector_property_client_id
    "clientSecret" = var.kyxstartconnector_property_client_secret
    "tenantName" = var.kyxstartconnector_property_tenant_name
  })
}
