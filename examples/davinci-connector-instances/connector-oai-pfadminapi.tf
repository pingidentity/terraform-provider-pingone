resource "pingone_davinci_connector_instance" "connector-oai-pfadminapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pfadminapi"
  }
  name = "My awesome connector-oai-pfadminapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-pfadminapi_property_auth_password
    "authUsername" = var.connector-oai-pfadminapi_property_auth_username
    "basePath" = var.connector-oai-pfadminapi_property_base_path
    "sslVerification" = var.connector-oai-pfadminapi_property_ssl_verification
  })
}
