resource "pingone_davinci_connector_instance" "biocatchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "biocatchConnector"
  }
  name = "My awesome biocatchConnector"
  properties = jsonencode({
    "apiUrl" = var.biocatchconnector_property_api_url
    "customerId" = var.biocatchconnector_property_customer_id
    "javascriptCdnUrl" = var.biocatchconnector_property_javascript_cdn_url
    "sdkToken" = var.biocatchconnector_property_sdk_token
    "truthApiKey" = var.biocatchconnector_property_truth_api_key
    "truthApiUrl" = var.biocatchconnector_property_truth_api_url
  })
}
