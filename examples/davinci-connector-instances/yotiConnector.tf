resource "pingone_davinci_connector_instance" "yotiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "yotiConnector"
  }
  name = "My awesome yotiConnector"
  property {
    name  = "addressLine1"
    type  = "string"
    value = var.yoticonnector_property_address_line1
  }
  property {
    name  = "addressLine2"
    type  = "string"
    value = var.yoticonnector_property_address_line2
  }
  property {
    name  = "addressLine3"
    type  = "string"
    value = var.yoticonnector_property_address_line3
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.yoticonnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.yoticonnector_property_button
  }
  property {
    name  = "challenge"
    type  = "string"
    value = var.yoticonnector_property_challenge
  }
  property {
    name  = "country"
    type  = "string"
    value = var.yoticonnector_property_country
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.yoticonnector_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.yoticonnector_property_custom_html
  }
  property {
    name  = "dateOfBirth"
    type  = "string"
    value = var.yoticonnector_property_date_of_birth
  }
  property {
    name  = "estimationLevel"
    type  = "string"
    value = var.yoticonnector_property_estimation_level
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.yoticonnector_property_first_name
  }
  property {
    name  = "identityVerificationAuthenticity"
    type  = "string"
    value = var.yoticonnector_property_identity_verification_authenticity
  }
  property {
    name  = "identityVerificationIssuingCountry"
    type  = "string"
    value = var.yoticonnector_property_identity_verification_issuing_country
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.yoticonnector_property_last_name
  }
  property {
    name  = "message"
    type  = "string"
    value = var.yoticonnector_property_message
  }
  property {
    name  = "messageIcon"
    type  = "string"
    value = var.yoticonnector_property_message_icon
  }
  property {
    name  = "messageIconHeight"
    type  = "string"
    value = var.yoticonnector_property_message_icon_height
  }
  property {
    name  = "messageTitle"
    type  = "string"
    value = var.yoticonnector_property_message_title
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.yoticonnector_property_postal_code
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.yoticonnector_property_session_id
  }
  property {
    name  = "showFooter"
    type  = "string"
    value = var.yoticonnector_property_show_footer
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.yoticonnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.yoticonnector_property_skip_button_press
  }
  property {
    name  = "threshold"
    type  = "string"
    value = var.yoticonnector_property_threshold
  }
  property {
    name  = "type"
    type  = "string"
    value = var.yoticonnector_property_type
  }
}
