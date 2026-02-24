resource "pingone_davinci_connector_instance" "veriffConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "veriffConnector"
  }
  name = "My awesome veriffConnector"
  properties = jsonencode({
    "access_token" = var.veriffconnector_property_access_token
    "baseUrl" = var.veriffconnector_property_base_url
    "password" = var.veriffconnector_property_password
  })
}
