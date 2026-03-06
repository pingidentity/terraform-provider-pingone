resource "pingone_davinci_connector_instance" "connector-oai-launchdarklyrestapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-launchdarklyrestapi"
  }
  name = "My awesome connector-oai-launchdarklyrestapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-launchdarklyrestapi_property_auth_api_key
    "basePath" = var.connector-oai-launchdarklyrestapi_property_base_path
  })
}
