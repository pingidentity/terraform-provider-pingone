resource "pingone_davinci_connector_instance" "stringsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "stringsConnector"
  }
  name = "My awesome stringsConnector"
}
