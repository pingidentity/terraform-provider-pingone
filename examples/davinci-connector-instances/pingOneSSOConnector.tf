resource "pingone_davinci_connector_instance" "pingOneSSOConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneSSOConnector"
  }
  name = "My awesome pingOneSSOConnector"
  property {
    name  = "acccountCreatedTemplateLocale"
    type  = "string"
    value = var.pingonessoconnector_property_acccount_created_template_locale
  }
  property {
    name  = "acccountCreatedTemplateVariables"
    type  = "string"
    value = var.pingonessoconnector_property_acccount_created_template_variables
  }
  property {
    name  = "acceptLanguage"
    type  = "string"
    value = var.pingonessoconnector_property_accept_language
  }
  property {
    name  = "accountCreatedCustomTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_account_created_custom_template_variant
  }
  property {
    name  = "accountLinkId"
    type  = "string"
    value = var.pingonessoconnector_property_account_link_id
  }
  property {
    name  = "additionalGatewayUserTypeList"
    type  = "string"
    value = var.pingonessoconnector_property_additional_gateway_user_type_list
  }
  property {
    name  = "additionalUserProperties"
    type  = "string"
    value = var.pingonessoconnector_property_additional_user_properties
  }
  property {
    name  = "agreement"
    type  = "string"
    value = var.pingonessoconnector_property_agreement
  }
  property {
    name  = "agreementId"
    type  = "string"
    value = var.pingonessoconnector_property_agreement_id
  }
  property {
    name  = "agreementPresentationId"
    type  = "string"
    value = var.pingonessoconnector_property_agreement_presentation_id
  }
  property {
    name  = "alternativeIdentifier"
    type  = "string"
    value = var.pingonessoconnector_property_alternative_identifier
  }
  property {
    name  = "bypassPolicy"
    type  = "string"
    value = var.pingonessoconnector_property_bypass_policy
  }
  property {
    name  = "changePasswordUserCustomTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_change_password_user_custom_template_variant
  }
  property {
    name  = "changePasswordUserTemplateLocale"
    type  = "string"
    value = var.pingonessoconnector_property_change_password_user_template_locale
  }
  property {
    name  = "changePasswordUserTemplateVariables"
    type  = "string"
    value = var.pingonessoconnector_property_change_password_user_template_variables
  }
  property {
    name  = "changePasswordUserTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_change_password_user_template_variant
  }
  property {
    name  = "clearUserAttributes"
    type  = "string"
    value = var.pingonessoconnector_property_clear_user_attributes
  }
  property {
    name  = "clearUserAttributesLink"
    type  = "string"
    value = var.pingonessoconnector_property_clear_user_attributes_link
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
    name  = "createUserGroups"
    type  = "string"
    value = var.pingonessoconnector_property_create_user_groups
  }
  property {
    name  = "createUserGroupsInput"
    type  = "string"
    value = var.pingonessoconnector_property_create_user_groups_input
  }
  property {
    name  = "createUserGroupsLink"
    type  = "string"
    value = var.pingonessoconnector_property_create_user_groups_link
  }
  property {
    name  = "createUserGroupsSource"
    type  = "string"
    value = var.pingonessoconnector_property_create_user_groups_source
  }
  property {
    name  = "createUserIfNotFound"
    type  = "string"
    value = var.pingonessoconnector_property_create_user_if_not_found
  }
  property {
    name  = "currentPassword"
    type  = "string"
    value = var.pingonessoconnector_property_current_password
  }
  property {
    name  = "customDataArray"
    type  = "string"
    value = var.pingonessoconnector_property_custom_data_array
  }
  property {
    name  = "customDataText"
    type  = "string"
    value = var.pingonessoconnector_property_custom_data_text
  }
  property {
    name  = "customTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_custom_template_variant
  }
  property {
    name  = "description"
    type  = "string"
    value = var.pingonessoconnector_property_description
  }
  property {
    name  = "email"
    type  = "string"
    value = var.pingonessoconnector_property_email
  }
  property {
    name  = "enabled"
    type  = "string"
    value = var.pingonessoconnector_property_enabled
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "envRegionInfo"
    type  = "string"
    value = var.pingonessoconnector_property_env_region_info
  }
  property {
    name  = "escapeIdentifier"
    type  = "string"
    value = var.pingonessoconnector_property_escape_identifier
  }
  property {
    name  = "externalId"
    type  = "string"
    value = var.pingonessoconnector_property_external_id
  }
  property {
    name  = "family"
    type  = "string"
    value = var.pingonessoconnector_property_family
  }
  property {
    name  = "filterCondition"
    type  = "string"
    value = var.pingonessoconnector_property_filter_condition
  }
  property {
    name  = "forceChange"
    type  = "string"
    value = var.pingonessoconnector_property_force_change
  }
  property {
    name  = "gatewayId"
    type  = "string"
    value = var.pingonessoconnector_property_gateway_id
  }
  property {
    name  = "gatewayUserTypeList"
    type  = "string"
    value = var.pingonessoconnector_property_gateway_user_type_list
  }
  property {
    name  = "given"
    type  = "string"
    value = var.pingonessoconnector_property_given
  }
  property {
    name  = "group"
    type  = "string"
    value = var.pingonessoconnector_property_group
  }
  property {
    name  = "groupId"
    type  = "string"
    value = var.pingonessoconnector_property_group_id
  }
  property {
    name  = "groupIds"
    type  = "string"
    value = var.pingonessoconnector_property_group_ids
  }
  property {
    name  = "groupName"
    type  = "string"
    value = var.pingonessoconnector_property_group_name
  }
  property {
    name  = "identifier"
    type  = "string"
    value = var.pingonessoconnector_property_identifier
  }
  property {
    name  = "identityProvider"
    type  = "string"
    value = var.pingonessoconnector_property_identity_provider
  }
  property {
    name  = "identityProviderId"
    type  = "string"
    value = var.pingonessoconnector_property_identity_provider_id
  }
  property {
    name  = "lifecycleStatus"
    type  = "string"
    value = var.pingonessoconnector_property_lifecycle_status
  }
  property {
    name  = "locale"
    type  = "string"
    value = var.pingonessoconnector_property_locale
  }
  property {
    name  = "matchAttribute"
    type  = "string"
    value = var.pingonessoconnector_property_match_attribute
  }
  property {
    name  = "matchAttributes"
    type  = "string"
    value = var.pingonessoconnector_property_match_attributes
  }
  property {
    name  = "memberGroupRelationship"
    type  = "string"
    value = var.pingonessoconnector_property_member_group_relationship
  }
  property {
    name  = "mobilePhone"
    type  = "string"
    value = var.pingonessoconnector_property_mobile_phone
  }
  property {
    name  = "newPassword"
    type  = "string"
    value = var.pingonessoconnector_property_new_password
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.pingonessoconnector_property_next_event
  }
  property {
    name  = "nextPageLink"
    type  = "string"
    value = var.pingonessoconnector_property_next_page_link
  }
  property {
    name  = "overrideUser"
    type  = "string"
    value = var.pingonessoconnector_property_override_user
  }
  property {
    name  = "password"
    type  = "string"
    value = var.pingonessoconnector_property_password
  }
  property {
    name  = "passwordForCreateUser"
    type  = "string"
    value = var.pingonessoconnector_property_password_for_create_user
  }
  property {
    name  = "passwordGateway"
    type  = "string"
    value = var.pingonessoconnector_property_password_gateway
  }
  property {
    name  = "passwordValue"
    type  = "string"
    value = var.pingonessoconnector_property_password_value
  }
  property {
    name  = "population"
    type  = "string"
    value = var.pingonessoconnector_property_population
  }
  property {
    name  = "populationId"
    type  = "string"
    value = var.pingonessoconnector_property_population_id
  }
  property {
    name  = "populationIds"
    type  = "string"
    value = var.pingonessoconnector_property_population_ids
  }
  property {
    name  = "preferredLanguage"
    type  = "string"
    value = var.pingonessoconnector_property_preferred_language
  }
  property {
    name  = "primaryPhone"
    type  = "string"
    value = var.pingonessoconnector_property_primary_phone
  }
  property {
    name  = "recoveryCode"
    type  = "string"
    value = var.pingonessoconnector_property_recovery_code
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingonessoconnector_property_region
  }
  property {
    name  = "returnUserPasswordStatus"
    type  = "string"
    value = var.pingonessoconnector_property_return_user_password_status
  }
  property {
    name  = "scimFilter"
    type  = "string"
    value = var.pingonessoconnector_property_scim_filter
  }
  property {
    name  = "setPasswordAdminCustomTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_set_password_admin_custom_template_variant
  }
  property {
    name  = "setPasswordAdminTemplateLocale"
    type  = "string"
    value = var.pingonessoconnector_property_set_password_admin_template_locale
  }
  property {
    name  = "setPasswordAdminTemplateVariables"
    type  = "string"
    value = var.pingonessoconnector_property_set_password_admin_template_variables
  }
  property {
    name  = "setPasswordAdminTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_set_password_admin_template_variant
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingonessoconnector_property_show_powered_by
  }
  property {
    name  = "singleFactorSignOnFormId"
    type  = "string"
    value = var.pingonessoconnector_property_single_factor_sign_on_form_id
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingonessoconnector_property_skip_button_press
  }
  property {
    name  = "subFlowVersion"
    type  = "string"
    value = var.pingonessoconnector_property_sub_flow_version
  }
  property {
    name  = "subflow"
    type  = "string"
    value = var.pingonessoconnector_property_subflow
  }
  property {
    name  = "templateLocale"
    type  = "string"
    value = var.pingonessoconnector_property_template_locale
  }
  property {
    name  = "templateVariables"
    type  = "string"
    value = var.pingonessoconnector_property_template_variables
  }
  property {
    name  = "templateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_template_variant
  }
  property {
    name  = "templateVariantAccountCreated"
    type  = "string"
    value = var.pingonessoconnector_property_template_variant_account_created
  }
  property {
    name  = "templateVariantVerificationRequired"
    type  = "string"
    value = var.pingonessoconnector_property_template_variant_verification_required
  }
  property {
    name  = "theme"
    type  = "string"
    value = var.pingonessoconnector_property_theme
  }
  property {
    name  = "themeId"
    type  = "string"
    value = var.pingonessoconnector_property_theme_id
  }
  property {
    name  = "useCustomDataText"
    type  = "string"
    value = var.pingonessoconnector_property_use_custom_data_text
  }
  property {
    name  = "useCustomSCIMFilter"
    type  = "string"
    value = var.pingonessoconnector_property_use_custom_scimfilter
  }
  property {
    name  = "userFilter"
    type  = "string"
    value = var.pingonessoconnector_property_user_filter
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingonessoconnector_property_user_id
  }
  property {
    name  = "userIdentifierForFindUser"
    type  = "string"
    value = var.pingonessoconnector_property_user_identifier_for_find_user
  }
  property {
    name  = "userLocale"
    type  = "string"
    value = var.pingonessoconnector_property_user_locale
  }
  property {
    name  = "userTypeId"
    type  = "string"
    value = var.pingonessoconnector_property_user_type_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.pingonessoconnector_property_username
  }
  property {
    name  = "usernameGateway"
    type  = "string"
    value = var.pingonessoconnector_property_username_gateway
  }
  property {
    name  = "verificationCode"
    type  = "string"
    value = var.pingonessoconnector_property_verification_code
  }
  property {
    name  = "verificationCodeTemplateLocale"
    type  = "string"
    value = var.pingonessoconnector_property_verification_code_template_locale
  }
  property {
    name  = "verificationCodeTemplateVariables"
    type  = "string"
    value = var.pingonessoconnector_property_verification_code_template_variables
  }
  property {
    name  = "verificationRequiredCustomTemplateVariant"
    type  = "string"
    value = var.pingonessoconnector_property_verification_required_custom_template_variant
  }
}
