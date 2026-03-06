resource "pingone_davinci_connector_instance" "devicePolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "devicePolicyConnector"
  }
  name = "My awesome devicePolicyConnector"
}
