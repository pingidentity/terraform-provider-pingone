resource "pingone_davinci_connector_instance" "variablesConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "variablesConnector"
  }
  name = "My awesome variablesConnector"
}
