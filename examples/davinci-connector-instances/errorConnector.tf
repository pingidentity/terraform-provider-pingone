resource "pingone_davinci_connector_instance" "errorConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "errorConnector"
  }
  name = "My awesome errorConnector"
}
