resource "pingone_davinci_connector_instance" "pingOneCredentialsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneCredentialsConnector"
  }
  name = "My awesome pingOneCredentialsConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "digitalWalletApplicationId" = var.pingonecredentialsconnector_property_digital_wallet_application_id
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingonecredentialsconnector_property_region
  })
}
