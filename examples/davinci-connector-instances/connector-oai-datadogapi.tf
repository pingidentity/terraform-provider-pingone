resource "pingone_davinci_connector_instance" "connector-oai-datadogapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-datadogapi"
  }
  name = "My awesome connector-oai-datadogapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-datadogapi_property_auth_api_key
    "authApplicationKey" = var.connector-oai-datadogapi_property_auth_application_key
    "basePath" = var.connector-oai-datadogapi_property_base_path
  })
}
