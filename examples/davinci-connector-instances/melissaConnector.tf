resource "pingone_davinci_connector_instance" "melissaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "melissaConnector"
  }
  name = "My awesome melissaConnector"
  properties = jsonencode({
    "apiKey" = var.melissaconnector_property_api_key
  })
}
