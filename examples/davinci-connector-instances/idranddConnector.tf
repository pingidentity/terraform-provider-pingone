resource "pingone_davinci_connector_instance" "idranddConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idranddConnector"
  }
  name = "My awesome idranddConnector"
  properties = jsonencode({
    "apiKey" = var.idranddconnector_property_api_key
    "apiUrl" = var.idranddconnector_property_api_url
  })
}
