resource "pingone_davinci_connector_instance" "locationPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "locationPolicyConnector"
  }
  name = "My awesome locationPolicyConnector"
}
