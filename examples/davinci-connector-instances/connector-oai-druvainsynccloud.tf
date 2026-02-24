resource "pingone_davinci_connector_instance" "connector-oai-druvainsynccloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-druvainsynccloud"
  }
  name = "My awesome connector-oai-druvainsynccloud"
  properties = jsonencode({
    "authClientId" = var.connector-oai-druvainsynccloud_property_auth_client_id
    "authClientSecret" = var.connector-oai-druvainsynccloud_property_auth_client_secret
    "authTokenUrl" = var.connector-oai-druvainsynccloud_property_auth_token_url
    "basePath" = var.connector-oai-druvainsynccloud_property_base_path
  })
}
