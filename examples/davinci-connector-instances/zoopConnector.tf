resource "pingone_davinci_connector_instance" "zoopConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "zoopConnector"
  }
  name = "My awesome zoopConnector"
  properties = jsonencode({
    "agencyId" = var.zoopconnector_property_agency_id
    "apiKey" = var.zoopconnector_property_api_key
    "apiUrl" = var.zoopconnector_property_api_url
  })
}
