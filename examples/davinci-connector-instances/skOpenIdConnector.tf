resource "pingone_davinci_connector_instance" "skOpenIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skOpenIdConnector"
  }
  name = "My awesome skOpenIdConnector"
}
