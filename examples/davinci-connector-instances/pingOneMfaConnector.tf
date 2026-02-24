resource "pingone_davinci_connector_instance" "pingOneMfaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneMfaConnector"
  }
  name = "My awesome pingOneMfaConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "policyId" = var.pingonemfaconnector_property_policy_id
    "region" = var.pingonemfaconnector_property_region
  })
}
