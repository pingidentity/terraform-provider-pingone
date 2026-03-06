resource "pingone_davinci_connector_instance" "connector-oai-venafi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-venafi"
  }
  name = "My awesome connector-oai-venafi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-venafi_property_auth_api_key
    "basePath" = var.connector-oai-venafi_property_base_path
  })
}
