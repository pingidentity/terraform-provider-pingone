resource "pingone_davinci_connector_instance" "connectorShopify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorShopify"
  }
  name = "My awesome connectorShopify"
  properties = jsonencode({
    "accessToken" = var.connectorshopify_property_access_token
    "apiVersion" = var.connectorshopify_property_api_version
    "multipassSecret" = var.connectorshopify_property_multipass_secret
    "multipassStoreDomain" = var.connectorshopify_property_multipass_store_domain
    "yourStoreName" = var.connectorshopify_property_your_store_name
  })
}
