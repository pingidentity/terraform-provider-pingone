resource "pingone_davinci_connector_instance" "constellaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "constellaConnector"
  }
  name = "My awesome constellaConnector"
  property {
    name  = "appToken"
    type  = "string"
    value = var.constellaconnector_property_app_token
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.constellaconnector_property_base_url
  }
  property {
    name  = "emailAddress"
    type  = "string"
    value = var.constellaconnector_property_email_address
  }
  property {
    name  = "password"
    type  = "string"
    value = var.constellaconnector_property_password
  }
  property {
    name  = "profile"
    type  = "string"
    value = var.constellaconnector_property_profile
  }
  property {
    name  = "token"
    type  = "string"
    value = var.constellaconnector_property_token
  }
  property {
    name  = "username"
    type  = "string"
    value = var.constellaconnector_property_username
  }
}
