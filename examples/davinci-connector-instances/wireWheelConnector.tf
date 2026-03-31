resource "pingone_davinci_connector_instance" "wireWheelConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "wireWheelConnector"
  }
  name = "My awesome wireWheelConnector"
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "button"
    type  = "string"
    value = var.wirewheelconnector_property_button
  }
  property {
    name  = "channelId"
    type  = "string"
    value = var.wirewheelconnector_property_channel_id
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.wirewheelconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.wirewheelconnector_property_client_secret
  }
  property {
    name  = "consentForm"
    type  = "string"
    value = var.wirewheelconnector_property_consent_form
  }
  property {
    name  = "consentPayloadJSON"
    type  = "string"
    value = var.wirewheelconnector_property_consent_payload_json
  }
  property {
    name  = "formFieldsList"
    type  = "string"
    value = var.wirewheelconnector_property_form_fields_list
  }
  property {
    name  = "getExistingConsent"
    type  = "string"
    value = var.wirewheelconnector_property_get_existing_consent
  }
  property {
    name  = "issuerId"
    type  = "string"
    value = var.wirewheelconnector_property_issuer_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.wirewheelconnector_property_username
  }
}
