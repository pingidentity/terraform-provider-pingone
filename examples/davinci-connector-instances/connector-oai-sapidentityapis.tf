resource "pingone_davinci_connector_instance" "connector-oai-sapidentityapis" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-sapidentityapis"
  }
  name = "My awesome connector-oai-sapidentityapis"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-sapidentityapis_property_auth_api_key
    "basePath" = var.connector-oai-sapidentityapis_property_base_path
  })
}
