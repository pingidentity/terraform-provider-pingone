resource "pingone_davinci_connector_instance" "singpassLoginConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "singpassLoginConnector"
  }
  name = "My awesome singpassLoginConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.singpassloginconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.singpassloginconnector_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.singpassloginconnector_property_show_powered_by
  }
  property {
    name  = "skWebhookUri"
    type  = "string"
    value = var.singpassloginconnector_property_sk_webhook_uri
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.singpassloginconnector_property_skip_button_press
  }
}
