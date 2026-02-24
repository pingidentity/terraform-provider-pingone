resource "pingone_davinci_connector_instance" "complyAdvatangeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "complyAdvatangeConnector"
  }
  name = "My awesome complyAdvatangeConnector"
  properties = jsonencode({
    "apiKey" = var.complyadvatangeconnector_property_api_key
    "baseUrl" = var.complyadvatangeconnector_property_base_url
  })
}
