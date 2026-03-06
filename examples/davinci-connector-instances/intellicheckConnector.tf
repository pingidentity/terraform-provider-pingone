resource "pingone_davinci_connector_instance" "intellicheckConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "intellicheckConnector"
  }
  name = "My awesome intellicheckConnector"
  properties = jsonencode({
    "apiKey" = var.intellicheckconnector_property_api_key
    "baseUrl" = var.intellicheckconnector_property_base_url
    "customerId" = var.intellicheckconnector_property_customer_id
  })
}
