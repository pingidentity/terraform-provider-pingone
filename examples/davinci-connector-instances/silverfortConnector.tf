resource "pingone_davinci_connector_instance" "silverfortConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "silverfortConnector"
  }
  name = "My awesome silverfortConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.silverfortconnector_property_api_key
  }
  property {
    name  = "appUserSecret"
    type  = "string"
    value = var.silverfortconnector_property_app_user_secret
  }
  property {
    name  = "consoleApi"
    type  = "string"
    value = var.silverfortconnector_property_console_api
  }
  property {
    name  = "description"
    type  = "string"
    value = var.silverfortconnector_property_description
  }
  property {
    name  = "severity"
    type  = "string"
    value = var.silverfortconnector_property_severity
  }
  property {
    name  = "shortName"
    type  = "string"
    value = var.silverfortconnector_property_short_name
  }
  property {
    name  = "userIdentifier"
    type  = "string"
    value = var.silverfortconnector_property_user_identifier
  }
  property {
    name  = "userParameter"
    type  = "string"
    value = var.silverfortconnector_property_user_parameter
  }
  property {
    name  = "validFor"
    type  = "string"
    value = var.silverfortconnector_property_valid_for
  }
}
