resource "pingone_davinci_connector_instance" "connector-oai-authomizeapireference" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-authomizeapireference"
  }
  name = "My awesome connector-oai-authomizeapireference"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-authomizeapireference_property_auth_api_key
    "basePath" = var.connector-oai-authomizeapireference_property_base_path
  })
}
