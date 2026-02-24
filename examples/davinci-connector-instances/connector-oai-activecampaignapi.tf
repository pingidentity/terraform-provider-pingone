resource "pingone_davinci_connector_instance" "connector-oai-activecampaignapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-activecampaignapi"
  }
  name = "My awesome connector-oai-activecampaignapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-activecampaignapi_property_auth_api_key
    "authApiVersion" = var.connector-oai-activecampaignapi_property_auth_api_version
    "basePath" = var.connector-oai-activecampaignapi_property_base_path
  })
}
