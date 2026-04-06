resource "pingone_davinci_connector_instance" "screenConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "screenConnector"
  }
  name = "My awesome screenConnector"
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.screenconnector_property_screen0_config
  }
}
