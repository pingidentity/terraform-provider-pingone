resource "pingone_davinci_connector_instance" "azureUserManagementConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "azureUserManagementConnector"
  }
  name = "My awesome azureUserManagementConnector"
  properties = jsonencode({
    "baseUrl" = var.azureusermanagementconnector_property_base_url
    "customApiUrl" = var.azureusermanagementconnector_property_custom_api_url
    "customAuth" = jsonencode({})
  })
}
