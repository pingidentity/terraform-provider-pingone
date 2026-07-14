resource "pingone_davinci_connector_instance" "mfaUseCaseConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mfaUseCaseConnector"
  }
  name = "My awesome mfaUseCaseConnector"
  property {
    name  = "activationStatus"
    type  = "string"
    value = var.mfausecaseconnector_property_activation_status
  }
  property {
    name  = "applicationId"
    type  = "string"
    value = var.mfausecaseconnector_property_application_id
  }
  property {
    name  = "applicationIds"
    type  = "string"
    value = var.mfausecaseconnector_property_application_ids
  }
  property {
    name  = "assertion"
    type  = "string"
    value = var.mfausecaseconnector_property_assertion
  }
  property {
    name  = "attestation"
    type  = "string"
    value = var.mfausecaseconnector_property_attestation
  }
  property {
    name  = "authTemplateName"
    type  = "string"
    value = var.mfausecaseconnector_property_auth_template_name
  }
  property {
    name  = "authenticatingApplicationId"
    type  = "string"
    value = var.mfausecaseconnector_property_authenticating_application_id
  }
  property {
    name  = "clientContext"
    type  = "string"
    value = var.mfausecaseconnector_property_client_context
  }
  property {
    name  = "compatibility"
    type  = "string"
    value = var.mfausecaseconnector_property_compatibility
  }
  property {
    name  = "customApplicationId"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_application_id
  }
  property {
    name  = "customAuthenticatingApplicationId"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_authenticating_application_id
  }
  property {
    name  = "customDeviceAuthenticationPolicyId"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_device_authentication_policy_id
  }
  property {
    name  = "customDeviceType"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_device_type
  }
  property {
    name  = "customNotificationPolicyId"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_notification_policy_id
  }
  property {
    name  = "customTemplateVariant"
    type  = "string"
    value = var.mfausecaseconnector_property_custom_template_variant
  }
  property {
    name  = "device"
    type  = "string"
    value = var.mfausecaseconnector_property_device
  }
  property {
    name  = "deviceAuthenRpId"
    type  = "string"
    value = var.mfausecaseconnector_property_device_authen_rp_id
  }
  property {
    name  = "deviceAuthenticationId"
    type  = "string"
    value = var.mfausecaseconnector_property_device_authentication_id
  }
  property {
    name  = "deviceAuthenticationPolicyId"
    type  = "string"
    value = var.mfausecaseconnector_property_device_authentication_policy_id
  }
  property {
    name  = "deviceId"
    type  = "string"
    value = var.mfausecaseconnector_property_device_id
  }
  property {
    name  = "deviceType"
    type  = "string"
    value = var.mfausecaseconnector_property_device_type
  }
  property {
    name  = "devices"
    type  = "string"
    value = var.mfausecaseconnector_property_devices
  }
  property {
    name  = "email"
    type  = "string"
    value = var.mfausecaseconnector_property_email
  }
  property {
    name  = "extension"
    type  = "string"
    value = var.mfausecaseconnector_property_extension
  }
  property {
    name  = "fido2Settings"
    type  = "boolean"
    value = var.mfausecaseconnector_property_fido2_settings
  }
  property {
    name  = "fidoCompatibility"
    type  = "string"
    value = var.mfausecaseconnector_property_fido_compatibility
  }
  property {
    name  = "mfaEnabled"
    type  = "boolean"
    value = var.mfausecaseconnector_property_mfa_enabled
  }
  property {
    name  = "mfaPolicyId"
    type  = "string"
    value = var.mfausecaseconnector_property_mfa_policy_id
  }
  property {
    name  = "mfaSettings"
    type  = "boolean"
    value = var.mfausecaseconnector_property_mfa_settings
  }
  property {
    name  = "mobilePayload"
    type  = "string"
    value = var.mfausecaseconnector_property_mobile_payload
  }
  property {
    name  = "nickname"
    type  = "string"
    value = var.mfausecaseconnector_property_nickname
  }
  property {
    name  = "notificationPolicyId"
    type  = "string"
    value = var.mfausecaseconnector_property_notification_policy_id
  }
  property {
    name  = "notificationSettings"
    type  = "boolean"
    value = var.mfausecaseconnector_property_notification_settings
  }
  property {
    name  = "origin"
    type  = "string"
    value = var.mfausecaseconnector_property_origin
  }
  property {
    name  = "otp"
    type  = "string"
    value = var.mfausecaseconnector_property_otp
  }
  property {
    name  = "otpSettings"
    type  = "boolean"
    value = var.mfausecaseconnector_property_otp_settings
  }
  property {
    name  = "pairingKeyId"
    type  = "string"
    value = var.mfausecaseconnector_property_pairing_key_id
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.mfausecaseconnector_property_phone
  }
  property {
    name  = "pushSettings"
    type  = "boolean"
    value = var.mfausecaseconnector_property_push_settings
  }
  property {
    name  = "reason"
    type  = "string"
    value = var.mfausecaseconnector_property_reason
  }
  property {
    name  = "rpId"
    type  = "string"
    value = var.mfausecaseconnector_property_rp_id
  }
  property {
    name  = "rpName"
    type  = "string"
    value = var.mfausecaseconnector_property_rp_name
  }
  property {
    name  = "sectionLabelDeviceDetails"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_device_details
  }
  property {
    name  = "sectionLabelFido"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_fido
  }
  property {
    name  = "sectionLabelMfaSettings"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_mfa_settings
  }
  property {
    name  = "sectionLabelNotification"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_notification
  }
  property {
    name  = "sectionLabelOtp"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_otp
  }
  property {
    name  = "sectionLabelPush"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_push
  }
  property {
    name  = "sectionLabelUserDetails"
    type  = "string"
    value = var.mfausecaseconnector_property_section_label_user_details
  }
  property {
    name  = "selectedDeviceOtp"
    type  = "string"
    value = var.mfausecaseconnector_property_selected_device_otp
  }
  property {
    name  = "serialNumber"
    type  = "string"
    value = var.mfausecaseconnector_property_serial_number
  }
  property {
    name  = "setDeviceOrder"
    type  = "string"
    value = var.mfausecaseconnector_property_set_device_order
  }
  property {
    name  = "setFilterFlag"
    type  = "boolean"
    value = var.mfausecaseconnector_property_set_filter_flag
  }
  property {
    name  = "showAdvancedFields"
    type  = "boolean"
    value = var.mfausecaseconnector_property_show_advanced_fields
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.mfausecaseconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.mfausecaseconnector_property_skip_button_press
  }
  property {
    name  = "statusFilter"
    type  = "string"
    value = var.mfausecaseconnector_property_status_filter
  }
  property {
    name  = "templateLocale"
    type  = "string"
    value = var.mfausecaseconnector_property_template_locale
  }
  property {
    name  = "templateVariables"
    type  = "string"
    value = var.mfausecaseconnector_property_template_variables
  }
  property {
    name  = "templateVariant"
    type  = "string"
    value = var.mfausecaseconnector_property_template_variant
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.mfausecaseconnector_property_user_agent
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.mfausecaseconnector_property_user_id
  }
}
