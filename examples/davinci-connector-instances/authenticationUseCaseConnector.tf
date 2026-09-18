resource "pingone_davinci_connector_instance" "authenticationUseCaseConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "authenticationUseCaseConnector"
  }
  name = "My awesome authenticationUseCaseConnector"
  property {
    name  = "accountCreatedCustomTemplateVariant"
    type  = "string"
    value = var.authenticationusecaseconnector_property_account_created_custom_template_variant
  }
  property {
    name  = "accountCreatedNotificationSettingsLabel"
    type  = "string"
    value = var.authenticationusecaseconnector_property_account_created_notification_settings_label
  }
  property {
    name  = "accountCreatedShowAdvancedFields"
    type  = "string"
    value = var.authenticationusecaseconnector_property_account_created_show_advanced_fields
  }
  property {
    name  = "accountCreatedTemplateVariables"
    type  = "string"
    value = var.authenticationusecaseconnector_property_account_created_template_variables
  }
  property {
    name  = "additionalUserProperties"
    type  = "string"
    value = var.authenticationusecaseconnector_property_additional_user_properties
  }
  property {
    name  = "agreement"
    type  = "string"
    value = var.authenticationusecaseconnector_property_agreement
  }
  property {
    name  = "agreementId"
    type  = "string"
    value = var.authenticationusecaseconnector_property_agreement_id
  }
  property {
    name  = "agreementLabel"
    type  = "string"
    value = var.authenticationusecaseconnector_property_agreement_label
  }
  property {
    name  = "changePasswordCustomTemplateVariant"
    type  = "string"
    value = var.authenticationusecaseconnector_property_change_password_custom_template_variant
  }
  property {
    name  = "changePasswordNotificationSettingsLabel"
    type  = "string"
    value = var.authenticationusecaseconnector_property_change_password_notification_settings_label
  }
  property {
    name  = "changePasswordShowAdvancedFields"
    type  = "string"
    value = var.authenticationusecaseconnector_property_change_password_show_advanced_fields
  }
  property {
    name  = "changePasswordTemplateVariables"
    type  = "string"
    value = var.authenticationusecaseconnector_property_change_password_template_variables
  }
  property {
    name  = "checkAccountStatusAgreement"
    type  = "string"
    value = var.authenticationusecaseconnector_property_check_account_status_agreement
  }
  property {
    name  = "checkAccountStatusAgreementId"
    type  = "string"
    value = var.authenticationusecaseconnector_property_check_account_status_agreement_id
  }
  property {
    name  = "checkAccountVerificationStatus"
    type  = "string"
    value = var.authenticationusecaseconnector_property_check_account_verification_status
  }
  property {
    name  = "checkAgreementConsent"
    type  = "string"
    value = var.authenticationusecaseconnector_property_check_agreement_consent
  }
  property {
    name  = "currentPassword"
    type  = "string"
    value = var.authenticationusecaseconnector_property_current_password
  }
  property {
    name  = "email"
    type  = "string"
    value = var.authenticationusecaseconnector_property_email
  }
  property {
    name  = "family"
    type  = "string"
    value = var.authenticationusecaseconnector_property_family
  }
  property {
    name  = "given"
    type  = "string"
    value = var.authenticationusecaseconnector_property_given
  }
  property {
    name  = "labelSendVerificationCode"
    type  = "string"
    value = var.authenticationusecaseconnector_property_label_send_verification_code
  }
  property {
    name  = "mobilePhone"
    type  = "string"
    value = var.authenticationusecaseconnector_property_mobile_phone
  }
  property {
    name  = "newPassword"
    type  = "string"
    value = var.authenticationusecaseconnector_property_new_password
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.authenticationusecaseconnector_property_next_event
  }
  property {
    name  = "notificationSettingsLabel"
    type  = "string"
    value = var.authenticationusecaseconnector_property_notification_settings_label
  }
  property {
    name  = "password"
    type  = "string"
    value = var.authenticationusecaseconnector_property_password
  }
  property {
    name  = "population"
    type  = "string"
    value = var.authenticationusecaseconnector_property_population
  }
  property {
    name  = "populationId"
    type  = "string"
    value = var.authenticationusecaseconnector_property_population_id
  }
  property {
    name  = "recoveryCode"
    type  = "string"
    value = var.authenticationusecaseconnector_property_recovery_code
  }
  property {
    name  = "requireUserToVerifyEmail"
    type  = "string"
    value = var.authenticationusecaseconnector_property_require_user_to_verify_email
  }
  property {
    name  = "resendPasswordRecoveryCode"
    type  = "string"
    value = var.authenticationusecaseconnector_property_resend_password_recovery_code
  }
  property {
    name  = "sectionLabelAccountVerification"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_account_verification
  }
  property {
    name  = "sectionLabelAgreementConsent"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_agreement_consent
  }
  property {
    name  = "sectionLabelNotification"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_notification
  }
  property {
    name  = "sectionLabelPwd"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_pwd
  }
  property {
    name  = "sectionLabelPwdAuthn"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_pwd_authn
  }
  property {
    name  = "sectionLabelUserDetails"
    type  = "string"
    value = var.authenticationusecaseconnector_property_section_label_user_details
  }
  property {
    name  = "showAdvancedFields"
    type  = "string"
    value = var.authenticationusecaseconnector_property_show_advanced_fields
  }
  property {
    name  = "templateVariantAccountCreated"
    type  = "string"
    value = var.authenticationusecaseconnector_property_template_variant_account_created
  }
  property {
    name  = "templateVariantChangePassword"
    type  = "string"
    value = var.authenticationusecaseconnector_property_template_variant_change_password
  }
  property {
    name  = "templateVariantVerificationCode"
    type  = "string"
    value = var.authenticationusecaseconnector_property_template_variant_verification_code
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.authenticationusecaseconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.authenticationusecaseconnector_property_username
  }
  property {
    name  = "verificationCode"
    type  = "string"
    value = var.authenticationusecaseconnector_property_verification_code
  }
  property {
    name  = "verificationCodeCustomTemplateVariant"
    type  = "string"
    value = var.authenticationusecaseconnector_property_verification_code_custom_template_variant
  }
  property {
    name  = "verificationCodeNotificationSettingsLabel"
    type  = "string"
    value = var.authenticationusecaseconnector_property_verification_code_notification_settings_label
  }
  property {
    name  = "verificationCodeShowAdvancedFields"
    type  = "string"
    value = var.authenticationusecaseconnector_property_verification_code_show_advanced_fields
  }
  property {
    name  = "verificationCodeTemplateVariables"
    type  = "string"
    value = var.authenticationusecaseconnector_property_verification_code_template_variables
  }
}
