resource "pingone_davinci_connector_instance" "samlIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlIdpConnector"
  }
  name = "My awesome samlIdpConnector"
  property {
    name  = "allowUnencryptedAssertion"
    type  = "string"
    value = var.samlidpconnector_property_allow_unencrypted_assertion
  }
  property {
    name  = "audience"
    type  = "string"
    value = var.samlidpconnector_property_audience
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.samlidpconnector_property_auth_type
  }
  property {
    name  = "authnContextClassRef"
    type  = "string"
    value = var.samlidpconnector_property_authn_context_class_ref
  }
  property {
    name  = "button"
    type  = "string"
    value = var.samlidpconnector_property_button
  }
  property {
    name  = "forceAuthn"
    type  = "string"
    value = var.samlidpconnector_property_force_authn
  }
  property {
    name  = "nameIdFormat"
    type  = "string"
    value = var.samlidpconnector_property_name_id_format
  }
  property {
    name  = "notBeforeSkew"
    type  = "string"
    value = var.samlidpconnector_property_not_before_skew
  }
  property {
    name  = "relayState"
    type  = "string"
    value = var.samlidpconnector_property_relay_state
  }
  property {
    name  = "requireSessionIndex"
    type  = "string"
    value = var.samlidpconnector_property_require_session_index
  }
  property {
    name  = "saml"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.samlidpconnector_property_show_powered_by
  }
  property {
    name  = "signRequest"
    type  = "string"
    value = var.samlidpconnector_property_sign_request
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.samlidpconnector_property_skip_button_press
  }
  property {
    name  = "specificConnectionId"
    type  = "string"
    value = var.samlidpconnector_property_specific_connection_id
  }
}
