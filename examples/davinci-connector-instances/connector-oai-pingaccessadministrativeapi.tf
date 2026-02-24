resource "pingone_davinci_connector_instance" "connector-oai-pingaccessadministrativeapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pingaccessadministrativeapi"
  }
  name = "My awesome connector-oai-pingaccessadministrativeapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-pingaccessadministrativeapi_property_auth_password
    "authUsername" = var.connector-oai-pingaccessadministrativeapi_property_auth_username
    "basePath" = var.connector-oai-pingaccessadministrativeapi_property_base_path
    "sslVerification" = var.connector-oai-pingaccessadministrativeapi_property_ssl_verification
  })
}
