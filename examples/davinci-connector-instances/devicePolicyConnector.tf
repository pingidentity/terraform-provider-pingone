resource "pingone_davinci_connector_instance" "devicePolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "devicePolicyConnector"
  }
  name = "My awesome devicePolicyConnector"
  property {
    name  = "browserNames"
    type  = "string"
    value = var.devicepolicyconnector_property_browser_names
  }
  property {
    name  = "deviceOSName"
    type  = "string"
    value = var.devicepolicyconnector_property_device_osname
  }
}
