resource "pingone_davinci_connector_instance" "connector-oai-sapidentityapis" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-sapidentityapis"
  }
  name = "My awesome connector-oai-sapidentityapis"
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_auth_api_key
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_base_path
  }
  property {
    name  = "schemasIdGet_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_schemas_id_get_id
  }
  property {
    name  = "usersGet_attributes"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_get_attributes
  }
  property {
    name  = "usersGet_count"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_get_count
  }
  property {
    name  = "usersGet_cursor"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_get_cursor
  }
  property {
    name  = "usersGet_excluded_attributes"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_get_excluded_attributes
  }
  property {
    name  = "usersIdDelete_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_delete_id
  }
  property {
    name  = "usersIdGet_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_get_id
  }
  property {
    name  = "usersIdPatch_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_patch_id
  }
  property {
    name  = "usersIdPatch_patchBody_PatchBody_Operations"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_patch_patch_body_patch_body_operations
  }
  property {
    name  = "usersIdPatch_patchBody_PatchBody_schemas"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_patch_patch_body_patch_body_schemas
  }
  property {
    name  = "usersIdPut_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_id
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preference24Hour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference24_hour
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceCurrencySymbolLocation"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_currency_symbol_location
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceDateFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_date_format
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceDefaultCalView"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_default_cal_view
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceDistance"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_distance
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceEndDayViewHour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_end_day_view_hour
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceFirstDayOfWeek"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_first_day_of_week
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceHourMinuteSeparator"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_hour_minute_separator
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceNegativeCurrencyFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_negative_currency_format
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceNegativeNumberFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_negative_number_format
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceNumberFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_number_format
  }
  property {
    name  = "usersIdPut_user_UserLocaleOverrides_preferenceStartDayViewHour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_locale_overrides_preference_start_day_view_hour
  }
  property {
    name  = "usersIdPut_user_UserMeta_created"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_meta_created
  }
  property {
    name  = "usersIdPut_user_UserMeta_lastModified"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_meta_last_modified
  }
  property {
    name  = "usersIdPut_user_UserMeta_location"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_meta_location
  }
  property {
    name  = "usersIdPut_user_UserMeta_resourceType"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_meta_resource_type
  }
  property {
    name  = "usersIdPut_user_UserMeta_version"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_meta_version
  }
  property {
    name  = "usersIdPut_user_UserName_academicTitle"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_academic_title
  }
  property {
    name  = "usersIdPut_user_UserName_familyName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_family_name
  }
  property {
    name  = "usersIdPut_user_UserName_familyNamePrefix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_family_name_prefix
  }
  property {
    name  = "usersIdPut_user_UserName_formatted"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_formatted
  }
  property {
    name  = "usersIdPut_user_UserName_givenName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_given_name
  }
  property {
    name  = "usersIdPut_user_UserName_honorificPrefix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_honorific_prefix
  }
  property {
    name  = "usersIdPut_user_UserName_honorificSuffix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_honorific_suffix
  }
  property {
    name  = "usersIdPut_user_UserName_legalName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_legal_name
  }
  property {
    name  = "usersIdPut_user_UserName_middleInitial"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_middle_initial
  }
  property {
    name  = "usersIdPut_user_UserName_middleName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_name_middle_name
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_$ref"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager__ref
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_displayName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_display_name
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_employeeNumber"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_employee_number
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_value"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_value
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_companyId"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_company_id
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_costCenter"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_cost_center
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_department"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_department
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_division"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_division
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_employeeNumber"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_employee_number
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_leavesOfAbsence"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_leaves_of_absence
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_legalEntity"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_legal_entity
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_organization"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_organization
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_startDate"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_start_date
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_terminationDate"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_termination_date
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionSap20User_userUuid"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_sap20_user_user_uuid
  }
  property {
    name  = "usersIdPut_user_UserUrnIetfParamsScimSchemasExtensionSettings20User_theme"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_urn_ietf_params_scim_schemas_extension_settings20_user_theme
  }
  property {
    name  = "usersIdPut_user_User_active"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_active
  }
  property {
    name  = "usersIdPut_user_User_addresses"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_addresses
  }
  property {
    name  = "usersIdPut_user_User_dateOfBirth"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_date_of_birth
  }
  property {
    name  = "usersIdPut_user_User_displayName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_display_name
  }
  property {
    name  = "usersIdPut_user_User_emails"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_emails
  }
  property {
    name  = "usersIdPut_user_User_emergencyContacts"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_emergency_contacts
  }
  property {
    name  = "usersIdPut_user_User_entitlements"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_entitlements
  }
  property {
    name  = "usersIdPut_user_User_externalId"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_external_id
  }
  property {
    name  = "usersIdPut_user_User_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_id
  }
  property {
    name  = "usersIdPut_user_User_nickName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_nick_name
  }
  property {
    name  = "usersIdPut_user_User_phoneNumbers"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_phone_numbers
  }
  property {
    name  = "usersIdPut_user_User_preferredLanguage"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_preferred_language
  }
  property {
    name  = "usersIdPut_user_User_schemas"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_schemas
  }
  property {
    name  = "usersIdPut_user_User_timezone"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_timezone
  }
  property {
    name  = "usersIdPut_user_User_title"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_title
  }
  property {
    name  = "usersIdPut_user_User_userName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_id_put_user_user_user_name
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preference24Hour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference24_hour
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceCurrencySymbolLocation"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_currency_symbol_location
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceDateFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_date_format
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceDefaultCalView"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_default_cal_view
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceDistance"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_distance
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceEndDayViewHour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_end_day_view_hour
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceFirstDayOfWeek"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_first_day_of_week
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceHourMinuteSeparator"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_hour_minute_separator
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceNegativeCurrencyFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_negative_currency_format
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceNegativeNumberFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_negative_number_format
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceNumberFormat"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_number_format
  }
  property {
    name  = "usersPost_user_UserLocaleOverrides_preferenceStartDayViewHour"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_locale_overrides_preference_start_day_view_hour
  }
  property {
    name  = "usersPost_user_UserMeta_created"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_meta_created
  }
  property {
    name  = "usersPost_user_UserMeta_lastModified"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_meta_last_modified
  }
  property {
    name  = "usersPost_user_UserMeta_location"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_meta_location
  }
  property {
    name  = "usersPost_user_UserMeta_resourceType"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_meta_resource_type
  }
  property {
    name  = "usersPost_user_UserMeta_version"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_meta_version
  }
  property {
    name  = "usersPost_user_UserName_academicTitle"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_academic_title
  }
  property {
    name  = "usersPost_user_UserName_familyName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_family_name
  }
  property {
    name  = "usersPost_user_UserName_familyNamePrefix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_family_name_prefix
  }
  property {
    name  = "usersPost_user_UserName_formatted"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_formatted
  }
  property {
    name  = "usersPost_user_UserName_givenName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_given_name
  }
  property {
    name  = "usersPost_user_UserName_honorificPrefix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_honorific_prefix
  }
  property {
    name  = "usersPost_user_UserName_honorificSuffix"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_honorific_suffix
  }
  property {
    name  = "usersPost_user_UserName_legalName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_legal_name
  }
  property {
    name  = "usersPost_user_UserName_middleInitial"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_middle_initial
  }
  property {
    name  = "usersPost_user_UserName_middleName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_name_middle_name
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_$ref"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager__ref
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_displayName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_display_name
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_employeeNumber"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_employee_number
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20UserManager_value"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_manager_value
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_companyId"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_company_id
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_costCenter"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_cost_center
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_department"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_department
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_division"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_division
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_employeeNumber"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_employee_number
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_leavesOfAbsence"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_leaves_of_absence
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_legalEntity"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_legal_entity
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_organization"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_organization
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_startDate"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_start_date
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionEnterprise20User_terminationDate"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_enterprise20_user_termination_date
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionSap20User_userUuid"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_sap20_user_user_uuid
  }
  property {
    name  = "usersPost_user_UserUrnIetfParamsScimSchemasExtensionSettings20User_theme"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_urn_ietf_params_scim_schemas_extension_settings20_user_theme
  }
  property {
    name  = "usersPost_user_User_active"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_active
  }
  property {
    name  = "usersPost_user_User_addresses"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_addresses
  }
  property {
    name  = "usersPost_user_User_dateOfBirth"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_date_of_birth
  }
  property {
    name  = "usersPost_user_User_displayName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_display_name
  }
  property {
    name  = "usersPost_user_User_emails"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_emails
  }
  property {
    name  = "usersPost_user_User_emergencyContacts"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_emergency_contacts
  }
  property {
    name  = "usersPost_user_User_entitlements"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_entitlements
  }
  property {
    name  = "usersPost_user_User_externalId"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_external_id
  }
  property {
    name  = "usersPost_user_User_id"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_id
  }
  property {
    name  = "usersPost_user_User_nickName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_nick_name
  }
  property {
    name  = "usersPost_user_User_phoneNumbers"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_phone_numbers
  }
  property {
    name  = "usersPost_user_User_preferredLanguage"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_preferred_language
  }
  property {
    name  = "usersPost_user_User_schemas"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_schemas
  }
  property {
    name  = "usersPost_user_User_timezone"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_timezone
  }
  property {
    name  = "usersPost_user_User_title"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_title
  }
  property {
    name  = "usersPost_user_User_userName"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_post_user_user_user_name
  }
  property {
    name  = "usersSearchPost_searchRequest_SearchRequest_attributes"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_search_post_search_request_search_request_attributes
  }
  property {
    name  = "usersSearchPost_searchRequest_SearchRequest_count"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_search_post_search_request_search_request_count
  }
  property {
    name  = "usersSearchPost_searchRequest_SearchRequest_cursor"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_search_post_search_request_search_request_cursor
  }
  property {
    name  = "usersSearchPost_searchRequest_SearchRequest_excludedAttributes"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_search_post_search_request_search_request_excluded_attributes
  }
  property {
    name  = "usersSearchPost_searchRequest_SearchRequest_filter"
    type  = "string"
    value = var.connector-oai-sapidentityapis_property_users_search_post_search_request_search_request_filter
  }
}
