resource "pingone_davinci_connector_instance" "functionsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "functionsConnector"
  }
  name = "My awesome functionsConnector"
}
