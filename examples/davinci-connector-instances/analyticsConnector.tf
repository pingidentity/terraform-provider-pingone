resource "pingone_davinci_connector_instance" "analyticsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "analyticsConnector"
  }
  name = "My awesome analyticsConnector"
}
