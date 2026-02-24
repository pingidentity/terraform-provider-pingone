resource "pingone_davinci_connector_instance" "pingOneRiskConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneRiskConnector"
  }
  name = "My awesome pingOneRiskConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingoneriskconnector_property_region
  })
}
