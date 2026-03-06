resource "pingone_davinci_connector_instance" "annotationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "annotationConnector"
  }
  name = "My awesome annotationConnector"
}
