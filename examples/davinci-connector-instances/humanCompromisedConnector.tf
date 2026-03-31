resource "pingone_davinci_connector_instance" "humanCompromisedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "humanCompromisedConnector"
  }
  name = "My awesome humanCompromisedConnector"
  property {
    name  = "appId"
    type  = "string"
    value = var.humancompromisedconnector_property_app_id
  }
  property {
    name  = "authToken"
    type  = "string"
    value = var.humancompromisedconnector_property_auth_token
  }
  property {
    name  = "password"
    type  = "string"
    value = var.humancompromisedconnector_property_password
  }
  property {
    name  = "username"
    type  = "string"
    value = var.humancompromisedconnector_property_username
  }
}
