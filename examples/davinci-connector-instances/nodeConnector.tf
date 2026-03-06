resource "pingone_davinci_connector_instance" "nodeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nodeConnector"
  }
  name = "My awesome nodeConnector"
}
