resource "pingone_davinci_connector_instance" "pingOneMfaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneMfaConnector"
  }
  name = "My awesome pingOneMfaConnector"
  property {
    name  = "applicationId"
    type  = "string"
    value = var.pingonemfaconnector_property_application_id
  }
  property {
    name  = "applicationIds"
    type  = "string"
    value = var.pingonemfaconnector_property_application_ids
  }
  property {
    name  = "applicationName"
    type  = "string"
    value = var.pingonemfaconnector_property_application_name
  }
  property {
    name  = "assertion"
    type  = "string"
    value = var.pingonemfaconnector_property_assertion
  }
  property {
    name  = "attestation"
    type  = "string"
    value = var.pingonemfaconnector_property_attestation
  }
  property {
    name  = "authTemplateName"
    type  = "string"
    value = var.pingonemfaconnector_property_auth_template_name
  }
  property {
    name  = "authenticatingApplicationId"
    type  = "string"
    value = var.pingonemfaconnector_property_authenticating_application_id
  }
  property {
    name  = "authenticationCodeId"
    type  = "string"
    value = var.pingonemfaconnector_property_authentication_code_id
  }
  property {
    name  = "clientContext"
    type  = "string"
    value = var.pingonemfaconnector_property_client_context
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingone_worker_app_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingone_worker_app_client_secret
  }
  property {
    name  = "compatibility"
    type  = "string"
    value = var.pingonemfaconnector_property_compatibility
  }
  property {
    name  = "createDeviceTestMode"
    type  = "string"
    value = var.pingonemfaconnector_property_create_device_test_mode
  }
  property {
    name  = "customApplicationId"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_application_id
  }
  property {
    name  = "customAuthenticatingApplicationId"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_authenticating_application_id
  }
  property {
    name  = "customDeviceAuthenticationPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_device_authentication_policy_id
  }
  property {
    name  = "customDeviceType"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_device_type
  }
  property {
    name  = "customNotificationPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_notification_policy_id
  }
  property {
    name  = "customTemplateVariant"
    type  = "string"
    value = var.pingonemfaconnector_property_custom_template_variant
  }
  property {
    name  = "defaultDeviceId"
    type  = "string"
    value = var.pingonemfaconnector_property_default_device_id
  }
  property {
    name  = "device"
    type  = "string"
    value = var.pingonemfaconnector_property_device
  }
  property {
    name  = "deviceAuthenRpId"
    type  = "string"
    value = var.pingonemfaconnector_property_device_authen_rp_id
  }
  property {
    name  = "deviceAuthenticationId"
    type  = "string"
    value = var.pingonemfaconnector_property_device_authentication_id
  }
  property {
    name  = "deviceAuthenticationPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_device_authentication_policy_id
  }
  property {
    name  = "deviceId"
    type  = "string"
    value = var.pingonemfaconnector_property_device_id
  }
  property {
    name  = "deviceOrderList"
    type  = "string"
    value = var.pingonemfaconnector_property_device_order_list
  }
  property {
    name  = "deviceType"
    type  = "string"
    value = var.pingonemfaconnector_property_device_type
  }
  property {
    name  = "deviceTypes"
    type  = "string"
    value = var.pingonemfaconnector_property_device_types
  }
  property {
    name  = "devices"
    type  = "string"
    value = var.pingonemfaconnector_property_devices
  }
  property {
    name  = "duration"
    type  = "string"
    value = var.pingonemfaconnector_property_duration
  }
  property {
    name  = "email"
    type  = "string"
    value = var.pingonemfaconnector_property_email
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "extension"
    type  = "string"
    value = var.pingonemfaconnector_property_extension
  }
  property {
    name  = "fidoCompatibility"
    type  = "string"
    value = var.pingonemfaconnector_property_fido_compatibility
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.pingonemfaconnector_property_ip
  }
  property {
    name  = "jsonAttributes"
    type  = "string"
    value = var.pingonemfaconnector_property_json_attributes
  }
  property {
    name  = "lastAuthenticationMethod"
    type  = "string"
    value = var.pingonemfaconnector_property_last_authentication_method
  }
  property {
    name  = "lastMFATransaction"
    type  = "string"
    value = var.pingonemfaconnector_property_last_mfatransaction
  }
  property {
    name  = "mfaEnabled"
    type  = "boolean"
    value = var.pingonemfaconnector_property_mfa_enabled
  }
  property {
    name  = "mfaPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_mfa_policy_id
  }
  property {
    name  = "mobilePayload"
    type  = "string"
    value = var.pingonemfaconnector_property_mobile_payload
  }
  property {
    name  = "nickname"
    type  = "string"
    value = var.pingonemfaconnector_property_nickname
  }
  property {
    name  = "notificationPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_notification_policy_id
  }
  property {
    name  = "oathResync"
    type  = "boolean"
    value = var.pingonemfaconnector_property_oath_resync
  }
  property {
    name  = "oneTimeDeviceTestMode"
    type  = "string"
    value = var.pingonemfaconnector_property_one_time_device_test_mode
  }
  property {
    name  = "oneTimeDeviceType"
    type  = "string"
    value = var.pingonemfaconnector_property_one_time_device_type
  }
  property {
    name  = "oneTimeEmailDevice"
    type  = "string"
    value = var.pingonemfaconnector_property_one_time_email_device
  }
  property {
    name  = "oneTimeSmsDevice"
    type  = "string"
    value = var.pingonemfaconnector_property_one_time_sms_device
  }
  property {
    name  = "oneTimeVoiceDevice"
    type  = "string"
    value = var.pingonemfaconnector_property_one_time_voice_device
  }
  property {
    name  = "origin"
    type  = "string"
    value = var.pingonemfaconnector_property_origin
  }
  property {
    name  = "otp"
    type  = "string"
    value = var.pingonemfaconnector_property_otp
  }
  property {
    name  = "pairingKeyId"
    type  = "string"
    value = var.pingonemfaconnector_property_pairing_key_id
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.pingonemfaconnector_property_phone
  }
  property {
    name  = "policyFactsApplicationId"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_facts_application_id
  }
  property {
    name  = "policyFactsApplicationName"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_facts_application_name
  }
  property {
    name  = "policyFactsGroups"
    type  = "json"
    value = var.pingonemfaconnector_property_policy_facts_groups
  }
  property {
    name  = "policyFactsIp"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_facts_ip
  }
  property {
    name  = "policyFactsRiskLevel"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_facts_risk_level
  }
  property {
    name  = "policyFactsUserAgent"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_facts_user_agent
  }
  property {
    name  = "policyId"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_id
  }
  property {
    name  = "policyRecentAuthMethod"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_recent_auth_method
  }
  property {
    name  = "policyRecentAuthTimestamp"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_recent_auth_timestamp
  }
  property {
    name  = "policyRecentIp"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_recent_ip
  }
  property {
    name  = "policyRecentPingOneCookie"
    type  = "string"
    value = var.pingonemfaconnector_property_policy_recent_ping_one_cookie
  }
  property {
    name  = "reason"
    type  = "string"
    value = var.pingonemfaconnector_property_reason
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingonemfaconnector_property_region
  }
  property {
    name  = "rememberMeCookie"
    type  = "string"
    value = var.pingonemfaconnector_property_remember_me_cookie
  }
  property {
    name  = "rememberMeDeviceType"
    type  = "string"
    value = var.pingonemfaconnector_property_remember_me_device_type
  }
  property {
    name  = "rememberMePayload"
    type  = "string"
    value = var.pingonemfaconnector_property_remember_me_payload
  }
  property {
    name  = "rpId"
    type  = "string"
    value = var.pingonemfaconnector_property_rp_id
  }
  property {
    name  = "rpName"
    type  = "string"
    value = var.pingonemfaconnector_property_rp_name
  }
  property {
    name  = "selectedDevice"
    type  = "string"
    value = var.pingonemfaconnector_property_selected_device
  }
  property {
    name  = "selectedDeviceId"
    type  = "string"
    value = var.pingonemfaconnector_property_selected_device_id
  }
  property {
    name  = "selectedDeviceOtp"
    type  = "string"
    value = var.pingonemfaconnector_property_selected_device_otp
  }
  property {
    name  = "serialNumber"
    type  = "string"
    value = var.pingonemfaconnector_property_serial_number
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.pingonemfaconnector_property_session_id
  }
  property {
    name  = "setDeviceOrder"
    type  = "string"
    value = var.pingonemfaconnector_property_set_device_order
  }
  property {
    name  = "setFilterFlag"
    type  = "boolean"
    value = var.pingonemfaconnector_property_set_filter_flag
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingonemfaconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingonemfaconnector_property_skip_button_press
  }
  property {
    name  = "status"
    type  = "string"
    value = var.pingonemfaconnector_property_status
  }
  property {
    name  = "statusFilter"
    type  = "string"
    value = var.pingonemfaconnector_property_status_filter
  }
  property {
    name  = "templateLocale"
    type  = "string"
    value = var.pingonemfaconnector_property_template_locale
  }
  property {
    name  = "templateVariables"
    type  = "string"
    value = var.pingonemfaconnector_property_template_variables
  }
  property {
    name  = "templateVariant"
    type  = "string"
    value = var.pingonemfaconnector_property_template_variant
  }
  property {
    name  = "timeUnit"
    type  = "string"
    value = var.pingonemfaconnector_property_time_unit
  }
  property {
    name  = "useDeviceOrderJsonAttributes"
    type  = "string"
    value = var.pingonemfaconnector_property_use_device_order_json_attributes
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.pingonemfaconnector_property_user_agent
  }
  property {
    name  = "userApproval"
    type  = "string"
    value = var.pingonemfaconnector_property_user_approval
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingonemfaconnector_property_user_id
  }
  property {
    name  = "usernameless"
    type  = "boolean"
    value = var.pingonemfaconnector_property_usernameless
  }
  property {
    name  = "webAuthnChallenge"
    type  = "string"
    value = var.pingonemfaconnector_property_web_authn_challenge
  }
  property {
    name  = "webAuthnChallengeForPairing"
    type  = "string"
    value = var.pingonemfaconnector_property_web_authn_challenge_for_pairing
  }
  property {
    name  = "workforceDeviceType"
    type  = "string"
    value = var.pingonemfaconnector_property_workforce_device_type
  }
  property {
    name  = "workforceDeviceTypes"
    type  = "string"
    value = var.pingonemfaconnector_property_workforce_device_types
  }
  property {
    name  = "workforcePolicyMfaPolicyId"
    type  = "string"
    value = var.pingonemfaconnector_property_workforce_policy_mfa_policy_id
  }
  property {
    name  = "yubiKey"
    type  = "string"
    value = var.pingonemfaconnector_property_yubi_key
  }
}
