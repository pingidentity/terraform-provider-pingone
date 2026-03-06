resource "pingone_davinci_connector_instance" "iproovV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iproovV2Connector"
  }
  name = "My awesome iproovV2Connector"
  properties = jsonencode({
    "apiKey" = var.iproovv2connector_property_api_key
    "secret" = var.iproovv2connector_property_secret
    "tenant" = var.iproovv2connector_property_tenant
  })
}
