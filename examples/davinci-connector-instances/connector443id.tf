resource "pingone_davinci_connector_instance" "connector443id" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector443id"
  }
  name = "My awesome connector443id"
  properties = jsonencode({
    "apiKey" = var.connector443id_property_api_key
  })
}
