resource "pingone_davinci_connector_instance" "connectorKeyri" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorKeyri"
  }
  name = "My awesome connectorKeyri"
}
