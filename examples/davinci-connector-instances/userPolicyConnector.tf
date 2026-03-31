resource "pingone_davinci_connector_instance" "userPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "userPolicyConnector"
  }
  name = "My awesome userPolicyConnector"
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.userpolicyconnector_property_claims_name_value_pairs
  }
  property {
    name  = "connectionId"
    type  = "string"
    value = var.userpolicyconnector_property_connection_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.userpolicyconnector_property_email
  }
  property {
    name  = "errorIfUserExists"
    type  = "string"
    value = var.userpolicyconnector_property_error_if_user_exists
  }
  property {
    name  = "name"
    type  = "string"
    value = var.userpolicyconnector_property_name
  }
  property {
    name  = "newPassword"
    type  = "string"
    value = var.userpolicyconnector_property_new_password
  }
  property {
    name  = "password"
    type  = "string"
    value = var.userpolicyconnector_property_password
  }
  property {
    name  = "passwordLengthMax"
    type  = "string"
    value = var.userpolicyconnector_property_password_length_max
  }
  property {
    name  = "passwordLengthMin"
    type  = "string"
    value = var.userpolicyconnector_property_password_length_min
  }
  property {
    name  = "passwordPreviousXPasswords"
    type  = "string"
    value = var.userpolicyconnector_property_password_previous_xpasswords
  }
  property {
    name  = "passwordRequireLowercase"
    type  = "string"
    value = var.userpolicyconnector_property_password_require_lowercase
  }
  property {
    name  = "passwordRequireNumbers"
    type  = "string"
    value = var.userpolicyconnector_property_password_require_numbers
  }
  property {
    name  = "passwordRequireSpecial"
    type  = "string"
    value = var.userpolicyconnector_property_password_require_special
  }
  property {
    name  = "passwordRequireUppercase"
    type  = "string"
    value = var.userpolicyconnector_property_password_require_uppercase
  }
  property {
    name  = "passwordSpacesOk"
    type  = "string"
    value = var.userpolicyconnector_property_password_spaces_ok
  }
  property {
    name  = "passwordsEnabled"
    type  = "string"
    value = var.userpolicyconnector_property_passwords_enabled
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.userpolicyconnector_property_phone_number
  }
  property {
    name  = "userAlias"
    type  = "string"
    value = var.userpolicyconnector_property_user_alias
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.userpolicyconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.userpolicyconnector_property_username
  }
}
