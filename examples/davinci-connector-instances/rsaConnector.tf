resource "pingone_davinci_connector_instance" "rsaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "rsaConnector"
  }
  name = "My awesome rsaConnector"
  properties = jsonencode({
    "accessId" = var.rsaconnector_property_access_id
    "accessKey" = var.rsaconnector_property_access_key
    "baseUrl" = var.rsaconnector_property_base_url
  })
}
