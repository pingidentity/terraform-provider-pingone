resource "pingone_davinci_connector_instance" "connector-oai-venafi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-venafi"
  }
  name = "My awesome connector-oai-venafi"
  property {
    name  = "_delete_id"
    type  = "string"
    value = var.connector-oai-venafi_property__delete_id
  }
  property {
    name  = "addMember_id"
    type  = "string"
    value = var.connector-oai-venafi_property_add_member_id
  }
  property {
    name  = "addMember_teamMembersRequest_TeamMembersRequest_members"
    type  = "string"
    value = var.connector-oai-venafi_property_add_member_team_members_request_team_members_request_members
  }
  property {
    name  = "addOwner_id"
    type  = "string"
    value = var.connector-oai-venafi_property_add_owner_id
  }
  property {
    name  = "addOwner_teamOwnersRequest_TeamOwnersRequest_owners"
    type  = "string"
    value = var.connector-oai-venafi_property_add_owner_team_owners_request_team_owners_request_owners
  }
  property {
    name  = "apikeyCreate_apiKeyRequest_ApiKeyRequest_apiVersion"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_create_api_key_request_api_key_request_api_version
  }
  property {
    name  = "apikeyCreate_apiKeyRequest_ApiKeyRequest_validityDays"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_create_api_key_request_api_key_request_validity_days
  }
  property {
    name  = "apikeyCreate_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_create_user_id
  }
  property {
    name  = "apikeyGetAll_api_key_status"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_all_api_key_status
  }
  property {
    name  = "apikeyGetAll_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_all_user_id
  }
  property {
    name  = "apikeyGetByKey1_key"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key1_key
  }
  property {
    name  = "apikeyGetByKey1_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key1_user_id
  }
  property {
    name  = "apikeyGetByKey2_apiKeyRequest_ApiKeyRequest_apiVersion"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key2_api_key_request_api_key_request_api_version
  }
  property {
    name  = "apikeyGetByKey2_apiKeyRequest_ApiKeyRequest_validityDays"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key2_api_key_request_api_key_request_validity_days
  }
  property {
    name  = "apikeyGetByKey2_key"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key2_key
  }
  property {
    name  = "apikeyGetByKey2_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key2_user_id
  }
  property {
    name  = "apikeyGetByKey3_key"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key3_key
  }
  property {
    name  = "apikeyGetByKey3_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key3_user_id
  }
  property {
    name  = "apikeyGetByKey_apiKeyRequest_ApiKeyRequest_apiVersion"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key_api_key_request_api_key_request_api_version
  }
  property {
    name  = "apikeyGetByKey_apiKeyRequest_ApiKeyRequest_validityDays"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key_api_key_request_api_key_request_validity_days
  }
  property {
    name  = "apikeyGetByKey_key"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key_key
  }
  property {
    name  = "apikeyGetByKey_user_id"
    type  = "string"
    value = var.connector-oai-venafi_property_apikey_get_by_key_user_id
  }
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-venafi_property_auth_api_key
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-venafi_property_base_path
  }
  property {
    name  = "create1_createTeamRequest_CreateTeamRequest_members"
    type  = "string"
    value = var.connector-oai-venafi_property_create1_create_team_request_create_team_request_members
  }
  property {
    name  = "create1_createTeamRequest_CreateTeamRequest_name"
    type  = "string"
    value = var.connector-oai-venafi_property_create1_create_team_request_create_team_request_name
  }
  property {
    name  = "create1_createTeamRequest_CreateTeamRequest_owners"
    type  = "string"
    value = var.connector-oai-venafi_property_create1_create_team_request_create_team_request_owners
  }
  property {
    name  = "create1_createTeamRequest_CreateTeamRequest_role"
    type  = "string"
    value = var.connector-oai-venafi_property_create1_create_team_request_create_team_request_role
  }
  property {
    name  = "create1_createTeamRequest_CreateTeamRequest_userMatchingRules"
    type  = "string"
    value = var.connector-oai-venafi_property_create1_create_team_request_create_team_request_user_matching_rules
  }
  property {
    name  = "create_ssoConfigurationRequest_SsoConfigurationRequest_clientId"
    type  = "string"
    value = var.connector-oai-venafi_property_create_sso_configuration_request_sso_configuration_request_client_id
  }
  property {
    name  = "create_ssoConfigurationRequest_SsoConfigurationRequest_clientSecret"
    type  = "string"
    value = var.connector-oai-venafi_property_create_sso_configuration_request_sso_configuration_request_client_secret
  }
  property {
    name  = "create_ssoConfigurationRequest_SsoConfigurationRequest_issuerUrl"
    type  = "string"
    value = var.connector-oai-venafi_property_create_sso_configuration_request_sso_configuration_request_issuer_url
  }
  property {
    name  = "create_ssoConfigurationRequest_SsoConfigurationRequest_scopes"
    type  = "string"
    value = var.connector-oai-venafi_property_create_sso_configuration_request_sso_configuration_request_scopes
  }
  property {
    name  = "delete1_id"
    type  = "string"
    value = var.connector-oai-venafi_property_delete1_id
  }
  property {
    name  = "get1_id"
    type  = "string"
    value = var.connector-oai-venafi_property_get1_id
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_firstname"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_firstname
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_grecaptchaResponse"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_grecaptcha_response
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_invitationId"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_invitation_id
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_lastname"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_lastname
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_password"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_password
  }
  property {
    name  = "invitationsConfirm_invitationConfirmationRequest_InvitationConfirmationRequest_username"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_confirm_invitation_confirmation_request_invitation_confirmation_request_username
  }
  property {
    name  = "invitationsCreate_invitationRequest_InvitationRequest_properties"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_create_invitation_request_invitation_request_properties
  }
  property {
    name  = "invitationsCreate_invitationRequest_InvitationRequest_roles"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_create_invitation_request_invitation_request_roles
  }
  property {
    name  = "invitationsCreate_invitationRequest_InvitationRequest_teams"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_create_invitation_request_invitation_request_teams
  }
  property {
    name  = "invitationsGetById_id"
    type  = "string"
    value = var.connector-oai-venafi_property_invitations_get_by_id_id
  }
  property {
    name  = "notificationsCreate_notificationConfigurationRequest_NotificationConfigurationRequest_recipients"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_create_notification_configuration_request_notification_configuration_request_recipients
  }
  property {
    name  = "notificationsCreate_notificationConfigurationRequest_NotificationConfigurationRequest_type"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_create_notification_configuration_request_notification_configuration_request_type
  }
  property {
    name  = "notificationsCreate_notificationConfigurationRequest_RecurrencePatternInformation_recurrenceType"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_create_notification_configuration_request_recurrence_pattern_information_recurrence_type
  }
  property {
    name  = "notificationsCreate_notificationConfigurationRequest_RecurrencePatternInformation_recurrenceValues"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_create_notification_configuration_request_recurrence_pattern_information_recurrence_values
  }
  property {
    name  = "notificationsDelete_id"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_delete_id
  }
  property {
    name  = "notificationsGetById_id"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_get_by_id_id
  }
  property {
    name  = "notificationsGetByType_type"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_get_by_type_type
  }
  property {
    name  = "notificationsUnsubscribe_id"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_unsubscribe_id
  }
  property {
    name  = "notificationsUnsubscribe_recipient_token"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_unsubscribe_recipient_token
  }
  property {
    name  = "notificationsUpdate_id"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_update_id
  }
  property {
    name  = "notificationsUpdate_notificationConfigurationRequest_NotificationConfigurationRequest_recipients"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_update_notification_configuration_request_notification_configuration_request_recipients
  }
  property {
    name  = "notificationsUpdate_notificationConfigurationRequest_NotificationConfigurationRequest_type"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_update_notification_configuration_request_notification_configuration_request_type
  }
  property {
    name  = "notificationsUpdate_notificationConfigurationRequest_RecurrencePatternInformation_recurrenceType"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_update_notification_configuration_request_recurrence_pattern_information_recurrence_type
  }
  property {
    name  = "notificationsUpdate_notificationConfigurationRequest_RecurrencePatternInformation_recurrenceValues"
    type  = "string"
    value = var.connector-oai-venafi_property_notifications_update_notification_configuration_request_recurrence_pattern_information_recurrence_values
  }
  property {
    name  = "patch_id"
    type  = "string"
    value = var.connector-oai-venafi_property_patch_id
  }
  property {
    name  = "patch_updateTeamRequest_UpdateTeamRequest_name"
    type  = "string"
    value = var.connector-oai-venafi_property_patch_update_team_request_update_team_request_name
  }
  property {
    name  = "patch_updateTeamRequest_UpdateTeamRequest_role"
    type  = "string"
    value = var.connector-oai-venafi_property_patch_update_team_request_update_team_request_role
  }
  property {
    name  = "patch_updateTeamRequest_UpdateTeamRequest_userMatchingRules"
    type  = "string"
    value = var.connector-oai-venafi_property_patch_update_team_request_update_team_request_user_matching_rules
  }
  property {
    name  = "preferencesCreate_userPreferenceRequest_UserPreferenceRequest_name"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_create_user_preference_request_user_preference_request_name
  }
  property {
    name  = "preferencesCreate_userPreferenceRequest_UserPreferenceRequest_value"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_create_user_preference_request_user_preference_request_value
  }
  property {
    name  = "preferencesDelete_id"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_delete_id
  }
  property {
    name  = "preferencesGetById_id"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_get_by_id_id
  }
  property {
    name  = "preferencesGetByName_name"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_get_by_name_name
  }
  property {
    name  = "preferencesUpdate_id"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_update_id
  }
  property {
    name  = "preferencesUpdate_userPreferenceRequest_UserPreferenceRequest_name"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_update_user_preference_request_user_preference_request_name
  }
  property {
    name  = "preferencesUpdate_userPreferenceRequest_UserPreferenceRequest_value"
    type  = "string"
    value = var.connector-oai-venafi_property_preferences_update_user_preference_request_user_preference_request_value
  }
  property {
    name  = "removeMember_id"
    type  = "string"
    value = var.connector-oai-venafi_property_remove_member_id
  }
  property {
    name  = "removeMember_teamMembersRequest_TeamMembersRequest_members"
    type  = "string"
    value = var.connector-oai-venafi_property_remove_member_team_members_request_team_members_request_members
  }
  property {
    name  = "removeOwner_id"
    type  = "string"
    value = var.connector-oai-venafi_property_remove_owner_id
  }
  property {
    name  = "removeOwner_teamOwnersRequest_TeamOwnersRequest_owners"
    type  = "string"
    value = var.connector-oai-venafi_property_remove_owner_team_owners_request_team_owners_request_owners
  }
  property {
    name  = "update_id"
    type  = "string"
    value = var.connector-oai-venafi_property_update_id
  }
  property {
    name  = "update_ssoConfigurationRequest_SsoConfigurationRequest_clientId"
    type  = "string"
    value = var.connector-oai-venafi_property_update_sso_configuration_request_sso_configuration_request_client_id
  }
  property {
    name  = "update_ssoConfigurationRequest_SsoConfigurationRequest_clientSecret"
    type  = "string"
    value = var.connector-oai-venafi_property_update_sso_configuration_request_sso_configuration_request_client_secret
  }
  property {
    name  = "update_ssoConfigurationRequest_SsoConfigurationRequest_issuerUrl"
    type  = "string"
    value = var.connector-oai-venafi_property_update_sso_configuration_request_sso_configuration_request_issuer_url
  }
  property {
    name  = "update_ssoConfigurationRequest_SsoConfigurationRequest_scopes"
    type  = "string"
    value = var.connector-oai-venafi_property_update_sso_configuration_request_sso_configuration_request_scopes
  }
  property {
    name  = "useraccountsActivate_k"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_activate_k
  }
  property {
    name  = "useraccountsActivate_v"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_activate_v
  }
  property {
    name  = "useraccountsCheckResetPasswordToken_token"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_check_reset_password_token_token
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_companyId"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_company_id
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_companyName"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_company_name
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_firstname"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_firstname
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_grecaptchaResponse"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_grecaptcha_response
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_lastname"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_lastname
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_marketoAttributes"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_marketo_attributes
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_password"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_password
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_productEntitlements"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_product_entitlements
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_referralPartner"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_referral_partner
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_userAccountType"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_user_account_type
  }
  property {
    name  = "useraccountsCreate_userAccountRequest_UserAccountRequest_username"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_create_user_account_request_user_account_request_username
  }
  property {
    name  = "useraccountsResendActivation_resendActivationRequest_ResendActivationRequest_email"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_resend_activation_resend_activation_request_resend_activation_request_email
  }
  property {
    name  = "useraccountsResetPassword1_changePasswordRequest_ChangePasswordRequest_currentPassword"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_reset_password1_change_password_request_change_password_request_current_password
  }
  property {
    name  = "useraccountsResetPassword1_changePasswordRequest_ChangePasswordRequest_newPassword"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_reset_password1_change_password_request_change_password_request_new_password
  }
  property {
    name  = "useraccountsResetPassword_resetPasswordRequest_ResetPasswordRequest_email"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_reset_password_reset_password_request_reset_password_request_email
  }
  property {
    name  = "useraccountsRotateApiKey_k"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_rotate_api_key_k
  }
  property {
    name  = "useraccountsRotateApiKey_v"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_rotate_api_key_v
  }
  property {
    name  = "useraccountsUpdatePassword_updatePasswordRequest_UpdatePasswordRequest_newPassword"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_update_password_update_password_request_update_password_request_new_password
  }
  property {
    name  = "useraccountsUpdatePassword_updatePasswordRequest_UpdatePasswordRequest_token"
    type  = "string"
    value = var.connector-oai-venafi_property_useraccounts_update_password_update_password_request_update_password_request_token
  }
  property {
    name  = "usersDelete_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_delete_id
  }
  property {
    name  = "usersGetAll_user_status"
    type  = "string"
    value = var.connector-oai-venafi_property_users_get_all_user_status
  }
  property {
    name  = "usersGetAll_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_get_all_username
  }
  property {
    name  = "usersGetById_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_get_by_id_id
  }
  property {
    name  = "usersGetByUsername_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_get_by_username_username
  }
  property {
    name  = "usersGetUserLoginConfig_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_get_user_login_config_username
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUserPk_encryptedUsername"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_pk_encrypted_username
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUserPk_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_pk_username
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_activationEmailDelayedSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_activation_email_delayed_send_date
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_activationEmailSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_activation_email_send_date
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_activationKey"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_activation_key
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_companyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_company_id
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_creationDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_creation_date
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_dataEncryptionKeyContainerId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_data_encryption_key_container_id
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_dataEncryptionKeyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_data_encryption_key_id
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_emailAddress"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_email_address
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_encryptorDecryptors"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_encryptor_decryptors
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_eulaAcceptDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_eula_accept_date
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_failedLoginCount"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_failed_login_count
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_firstLoginDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_first_login_date
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_firstname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_firstname
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_id
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_lastname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_lastname
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_localLoginDisabled"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_local_login_disabled
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_marketoAttributes"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_marketo_attributes
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_memberedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_membered_teams
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_ownedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_owned_teams
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_passwordHash"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_password_hash
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_productRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_product_roles
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_requestedEntitlements"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_requested_entitlements
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_ssoStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_sso_status
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_systemRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_system_roles
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_userAccountType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_user_account_type
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_userFullName"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_user_full_name
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_userStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_user_status
  }
  property {
    name  = "usersUpdateAccountType_cUser_CUser_userType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_c_user_cuser_user_type
  }
  property {
    name  = "usersUpdateAccountType_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_account_type_id
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUserPk_encryptedUsername"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_pk_encrypted_username
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUserPk_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_pk_username
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_activationEmailDelayedSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_activation_email_delayed_send_date
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_activationEmailSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_activation_email_send_date
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_activationKey"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_activation_key
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_companyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_company_id
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_creationDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_creation_date
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_dataEncryptionKeyContainerId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_data_encryption_key_container_id
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_dataEncryptionKeyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_data_encryption_key_id
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_emailAddress"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_email_address
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_encryptorDecryptors"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_encryptor_decryptors
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_eulaAcceptDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_eula_accept_date
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_failedLoginCount"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_failed_login_count
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_firstLoginDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_first_login_date
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_firstname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_firstname
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_id
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_lastname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_lastname
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_localLoginDisabled"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_local_login_disabled
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_marketoAttributes"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_marketo_attributes
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_memberedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_membered_teams
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_ownedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_owned_teams
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_passwordHash"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_password_hash
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_productRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_product_roles
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_requestedEntitlements"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_requested_entitlements
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_ssoStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_sso_status
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_systemRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_system_roles
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_userAccountType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_user_account_type
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_userFullName"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_user_full_name
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_userStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_user_status
  }
  property {
    name  = "usersUpdateLocalLogin_cUser_CUser_userType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_c_user_cuser_user_type
  }
  property {
    name  = "usersUpdateLocalLogin_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_local_login_id
  }
  property {
    name  = "usersUpdateRoles_cUser_CUserPk_encryptedUsername"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_pk_encrypted_username
  }
  property {
    name  = "usersUpdateRoles_cUser_CUserPk_username"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_pk_username
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_activationEmailDelayedSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_activation_email_delayed_send_date
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_activationEmailSendDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_activation_email_send_date
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_activationKey"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_activation_key
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_companyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_company_id
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_creationDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_creation_date
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_dataEncryptionKeyContainerId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_data_encryption_key_container_id
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_dataEncryptionKeyId"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_data_encryption_key_id
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_emailAddress"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_email_address
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_encryptorDecryptors"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_encryptor_decryptors
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_eulaAcceptDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_eula_accept_date
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_failedLoginCount"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_failed_login_count
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_firstLoginDate"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_first_login_date
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_firstname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_firstname
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_id
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_lastname"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_lastname
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_localLoginDisabled"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_local_login_disabled
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_marketoAttributes"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_marketo_attributes
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_memberedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_membered_teams
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_ownedTeams"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_owned_teams
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_passwordHash"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_password_hash
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_productRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_product_roles
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_requestedEntitlements"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_requested_entitlements
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_ssoStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_sso_status
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_systemRoles"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_system_roles
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_userAccountType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_user_account_type
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_userFullName"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_user_full_name
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_userStatus"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_user_status
  }
  property {
    name  = "usersUpdateRoles_cUser_CUser_userType"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_c_user_cuser_user_type
  }
  property {
    name  = "usersUpdateRoles_id"
    type  = "string"
    value = var.connector-oai-venafi_property_users_update_roles_id
  }
}
