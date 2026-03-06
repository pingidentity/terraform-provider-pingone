resource "pingone_davinci_connector_instance" "seonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "seonConnector"
  }
  name = "My awesome seonConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "licenseKey" = var.seonconnector_property_license_key
  })
}
