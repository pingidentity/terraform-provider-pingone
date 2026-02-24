resource "pingone_davinci_connector_instance" "pingOneFormsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneFormsConnector"
  }
  name = "My awesome pingOneFormsConnector"
}
