resource "pingone_davinci_connector_instance" "connector-oai-talendscim" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-talendscim"
  }
  name = "My awesome connector-oai-talendscim"
  properties = jsonencode({
    "authBearerToken" = var.connector-oai-talendscim_property_auth_bearer_token
    "basePath" = var.connector-oai-talendscim_property_base_path
  })
}
