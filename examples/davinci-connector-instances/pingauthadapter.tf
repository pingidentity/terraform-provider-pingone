resource "pingone_davinci_connector_instance" "pingauthadapter" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingauthadapter"
  }
  name = "My awesome pingauthadapter"
}
