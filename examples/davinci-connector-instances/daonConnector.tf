resource "pingone_davinci_connector_instance" "daonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "daonConnector"
  }
  name = "My awesome daonConnector"
  properties = jsonencode({
    "apiUrl" = var.daonconnector_property_api_url
    "password" = var.daonconnector_property_password
    "username" = var.daonconnector_property_username
  })
}
