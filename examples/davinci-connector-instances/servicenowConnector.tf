resource "pingone_davinci_connector_instance" "servicenowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "servicenowConnector"
  }
  name = "My awesome servicenowConnector"
  properties = jsonencode({
    "adminUsername" = var.servicenowconnector_property_admin_username
    "apiUrl" = var.servicenowconnector_property_api_url
    "password" = var.servicenowconnector_property_password
  })
}
