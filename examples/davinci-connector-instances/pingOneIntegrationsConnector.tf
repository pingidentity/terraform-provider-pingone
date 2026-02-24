resource "pingone_davinci_connector_instance" "pingOneIntegrationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneIntegrationsConnector"
  }
  name = "My awesome pingOneIntegrationsConnector"
}
