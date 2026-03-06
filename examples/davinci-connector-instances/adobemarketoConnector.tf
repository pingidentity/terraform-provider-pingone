resource "pingone_davinci_connector_instance" "adobemarketoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "adobemarketoConnector"
  }
  name = "My awesome adobemarketoConnector"
  properties = jsonencode({
    "clientId" = var.adobemarketoconnector_property_client_id
    "clientSecret" = var.adobemarketoconnector_property_client_secret
    "endpoint" = var.adobemarketoconnector_property_endpoint
  })
}
