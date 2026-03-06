resource "pingone_davinci_connector_instance" "pingOneAuthenticationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthenticationConnector"
  }
  name = "My awesome pingOneAuthenticationConnector"
}
