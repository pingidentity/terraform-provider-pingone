resource "pingone_davinci_connector_instance" "pingoneAdvancedIdentityCloudLoginConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingoneAdvancedIdentityCloudLoginConnector"
  }
  name = "My awesome pingoneAdvancedIdentityCloudLoginConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.pingoneadvancedidentitycloudloginconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.pingoneadvancedidentitycloudloginconnector_property_button
  }
  property {
    name  = "openId"
    type  = "string"
    value = var.pingoneadvancedidentitycloudloginconnector_property_open_id
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingoneadvancedidentitycloudloginconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingoneadvancedidentitycloudloginconnector_property_skip_button_press
  }
}
