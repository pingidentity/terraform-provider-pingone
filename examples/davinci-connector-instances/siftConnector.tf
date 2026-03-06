resource "pingone_davinci_connector_instance" "siftConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "siftConnector"
  }
  name = "My awesome siftConnector"
  properties = jsonencode({
    "apiKey" = var.siftconnector_property_api_key
  })
}
