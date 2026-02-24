resource "pingone_davinci_connector_instance" "connector-oai-mailjetapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-mailjetapi"
  }
  name = "My awesome connector-oai-mailjetapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-mailjetapi_property_auth_password
    "authUsername" = var.connector-oai-mailjetapi_property_auth_username
    "basePath" = var.connector-oai-mailjetapi_property_base_path
  })
}
