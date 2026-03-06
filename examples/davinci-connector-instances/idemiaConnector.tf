resource "pingone_davinci_connector_instance" "idemiaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idemiaConnector"
  }
  name = "My awesome idemiaConnector"
  properties = jsonencode({
    "apikey" = var.idemiaconnector_property_apikey
    "baseUrl" = var.idemiaconnector_property_base_url
  })
}
