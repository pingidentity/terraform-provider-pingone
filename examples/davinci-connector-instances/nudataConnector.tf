resource "pingone_davinci_connector_instance" "nudataConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nudataConnector"
  }
  name = "My awesome nudataConnector"
}
