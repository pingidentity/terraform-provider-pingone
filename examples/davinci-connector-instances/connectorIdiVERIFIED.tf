resource "pingone_davinci_connector_instance" "connectorIdiVERIFIED" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdiVERIFIED"
  }
  name = "My awesome connectorIdiVERIFIED"
  properties = jsonencode({
    "apiSecret" = var.connectoridiverified_property_api_secret
    "companyKey" = var.connectoridiverified_property_company_key
    "idiEnv" = var.connectoridiverified_property_idi_env
    "siteKey" = var.connectoridiverified_property_site_key
    "uniqueUrl" = var.connectoridiverified_property_unique_url
  })
}
