resource "pingone_davinci_connector_instance" "connectorZscaler" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZscaler"
  }
  name = "My awesome connectorZscaler"
  properties = jsonencode({
    "basePath" = var.connectorzscaler_property_base_path
    "baseURL" = var.base_url
    "zscalerAPIkey" = var.zscaler_api_key
    "zscalerPassword" = var.connectorzscaler_property_zscaler_password
    "zscalerUsername" = var.connectorzscaler_property_zscaler_username
  })
}
