resource "pingone_davinci_connector_instance" "forterConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "forterConnector"
  }
  name = "My awesome forterConnector"
  properties = jsonencode({
    "apiVersion" = var.forterconnector_property_api_version
    "secretKey" = var.forterconnector_property_secret_key
    "siteId" = var.forterconnector_property_site_id
  })
}
