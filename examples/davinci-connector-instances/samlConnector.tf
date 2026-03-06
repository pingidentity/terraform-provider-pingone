resource "pingone_davinci_connector_instance" "samlConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlConnector"
  }
  name = "My awesome samlConnector"
}
