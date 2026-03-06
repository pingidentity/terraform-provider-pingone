resource "pingone_davinci_connector_instance" "dataZooConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "dataZooConnector"
  }
  name = "My awesome dataZooConnector"
  properties = jsonencode({
    "password" = var.datazooconnector_property_password
    "username" = var.datazooconnector_property_username
  })
}
