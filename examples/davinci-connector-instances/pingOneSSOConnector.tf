resource "pingone_davinci_connector_instance" "pingOneSSOConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneSSOConnector"
  }
  name = "My awesome pingOneSSOConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "envRegionInfo" = var.pingonessoconnector_property_env_region_info
    "region" = var.pingonessoconnector_property_region
  })
}
