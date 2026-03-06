resource "pingone_davinci_connector_instance" "connector-oai-github" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-github"
  }
  name = "My awesome connector-oai-github"
  properties = jsonencode({
    "apiVersion" = var.connector-oai-github_property_api_version
    "authBearerToken" = var.connector-oai-github_property_auth_bearer_token
    "basePath" = var.connector-oai-github_property_base_path
  })
}
