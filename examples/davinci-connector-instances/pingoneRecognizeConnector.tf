resource "pingone_davinci_connector_instance" "pingoneRecognizeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingoneRecognizeConnector"
  }
  name = "My awesome pingoneRecognizeConnector"
}
