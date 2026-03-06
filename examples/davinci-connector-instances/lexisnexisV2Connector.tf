resource "pingone_davinci_connector_instance" "lexisnexisV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "lexisnexisV2Connector"
  }
  name = "My awesome lexisnexisV2Connector"
  properties = jsonencode({
    "apiKey" = var.lexisnexisv2connector_property_api_key
    "apiUrl" = var.lexisnexisv2connector_property_api_url
    "orgId" = var.lexisnexisv2connector_property_org_id
    "useCustomApiURL" = var.use_custom_api_url
  })
}
