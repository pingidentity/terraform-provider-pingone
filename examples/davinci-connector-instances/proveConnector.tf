resource "pingone_davinci_connector_instance" "proveConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "proveConnector"
  }
  name = "My awesome proveConnector"
  properties = jsonencode({
    "baseUrl" = var.proveconnector_property_base_url
    "clientId" = var.proveconnector_property_client_id
    "grantType" = var.proveconnector_property_grant_type
    "password" = var.proveconnector_property_password
    "username" = var.proveconnector_property_username
  })
}
