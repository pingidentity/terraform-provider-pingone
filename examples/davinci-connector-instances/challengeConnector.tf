resource "pingone_davinci_connector_instance" "challengeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "challengeConnector"
  }
  name = "My awesome challengeConnector"
}
