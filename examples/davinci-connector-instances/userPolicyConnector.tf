resource "pingone_davinci_connector_instance" "userPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "userPolicyConnector"
  }
  name = "My awesome userPolicyConnector"
  properties = jsonencode({
    "passwordExpiryInDays" = var.userpolicyconnector_property_password_expiry_in_days
    "passwordExpiryNotification" = var.userpolicyconnector_property_password_expiry_notification
    "passwordLengthMax" = var.userpolicyconnector_property_password_length_max
    "passwordLengthMin" = var.userpolicyconnector_property_password_length_min
    "passwordLockoutAttempts" = var.userpolicyconnector_property_password_lockout_attempts
    "passwordPreviousXPasswords" = var.userpolicyconnector_property_password_previous_x_passwords
    "passwordRequireLowercase" = var.userpolicyconnector_property_password_require_lowercase
    "passwordRequireNumbers" = var.userpolicyconnector_property_password_require_numbers
    "passwordRequireSpecial" = var.userpolicyconnector_property_password_require_special
    "passwordRequireUppercase" = var.userpolicyconnector_property_password_require_uppercase
    "passwordSpacesOk" = var.userpolicyconnector_property_password_spaces_ok
    "passwordsEnabled" = var.userpolicyconnector_property_passwords_enabled
    "temporaryPasswordExpiryInDays" = var.userpolicyconnector_property_temporary_password_expiry_in_days
  })
}
