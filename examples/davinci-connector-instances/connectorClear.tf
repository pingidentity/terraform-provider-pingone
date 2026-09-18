resource "pingone_davinci_connector_instance" "connectorClear" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClear"
  }
  name = "My awesome connectorClear"
  property {
    name  = "authType"
    type  = "string"
    value = var.connectorclear_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.connectorclear_property_button
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = var.connectorclear_property_custom_auth
  }
  property {
    name  = "customFields"
    type  = "string"
    value = var.connectorclear_property_custom_fields
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectorclear_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectorclear_property_skip_button_press
  }
  property {
    name  = "url"
    type  = "string"
    value = var.connectorclear_property_url
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.connectorclear_property_user_id
  }
  property {
    name  = "userProfileAddressCity"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_city
  }
  property {
    name  = "userProfileAddressCountry"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_country
  }
  property {
    name  = "userProfileAddressLine1"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_line1
  }
  property {
    name  = "userProfileAddressLine2"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_line2
  }
  property {
    name  = "userProfileAddressPostalCode"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_postal_code
  }
  property {
    name  = "userProfileAddressState"
    type  = "string"
    value = var.connectorclear_property_user_profile_address_state
  }
  property {
    name  = "userProfileDob"
    type  = "string"
    value = var.connectorclear_property_user_profile_dob
  }
  property {
    name  = "userProfileEmail"
    type  = "string"
    value = var.connectorclear_property_user_profile_email
  }
  property {
    name  = "userProfileFirstName"
    type  = "string"
    value = var.connectorclear_property_user_profile_first_name
  }
  property {
    name  = "userProfileLastName"
    type  = "string"
    value = var.connectorclear_property_user_profile_last_name
  }
  property {
    name  = "userProfilePhone"
    type  = "string"
    value = var.connectorclear_property_user_profile_phone
  }
}
