resource "pingone_davinci_connector_instance" "azureUserManagementConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "azureUserManagementConnector"
  }
  name = "My awesome azureUserManagementConnector"
  property {
    name  = "accountEnabled"
    type  = "string"
    value = var.azureusermanagementconnector_property_account_enabled
  }
  property {
    name  = "additionalUserProperties"
    type  = "string"
    value = var.azureusermanagementconnector_property_additional_user_properties
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.azureusermanagementconnector_property_auth_type
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.azureusermanagementconnector_property_base_url
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.azureusermanagementconnector_property_body_data
  }
  property {
    name  = "city"
    type  = "string"
    value = var.azureusermanagementconnector_property_city
  }
  property {
    name  = "country"
    type  = "string"
    value = var.azureusermanagementconnector_property_country
  }
  property {
    name  = "createUserAccountEnabled"
    type  = "string"
    value = var.azureusermanagementconnector_property_create_user_account_enabled
  }
  property {
    name  = "createUserDisplayName"
    type  = "string"
    value = var.azureusermanagementconnector_property_create_user_display_name
  }
  property {
    name  = "createUserMailNickname"
    type  = "string"
    value = var.azureusermanagementconnector_property_create_user_mail_nickname
  }
  property {
    name  = "createUserPassword"
    type  = "string"
    value = var.azureusermanagementconnector_property_create_user_password
  }
  property {
    name  = "createUserPrincipalName"
    type  = "string"
    value = var.azureusermanagementconnector_property_create_user_principal_name
  }
  property {
    name  = "customApiUrl"
    type  = "string"
    value = var.azureusermanagementconnector_property_custom_api_url
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "customQueryParams"
    type  = "string"
    value = var.azureusermanagementconnector_property_custom_query_params
  }
  property {
    name  = "department"
    type  = "string"
    value = var.azureusermanagementconnector_property_department
  }
  property {
    name  = "disabledPlans"
    type  = "string"
    value = var.azureusermanagementconnector_property_disabled_plans
  }
  property {
    name  = "displayName"
    type  = "string"
    value = var.azureusermanagementconnector_property_display_name
  }
  property {
    name  = "employeeId"
    type  = "string"
    value = var.azureusermanagementconnector_property_employee_id
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.azureusermanagementconnector_property_endpoint
  }
  property {
    name  = "forceChangePasswordNextSignIn"
    type  = "string"
    value = var.azureusermanagementconnector_property_force_change_password_next_sign_in
  }
  property {
    name  = "forceChangePasswordNextSignInWithMfa"
    type  = "string"
    value = var.azureusermanagementconnector_property_force_change_password_next_sign_in_with_mfa
  }
  property {
    name  = "givenName"
    type  = "string"
    value = var.azureusermanagementconnector_property_given_name
  }
  property {
    name  = "groupId"
    type  = "string"
    value = var.azureusermanagementconnector_property_group_id
  }
  property {
    name  = "groupUserId"
    type  = "string"
    value = var.azureusermanagementconnector_property_group_user_id
  }
  property {
    name  = "groups"
    type  = "string"
    value = var.azureusermanagementconnector_property_groups
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.azureusermanagementconnector_property_headers
  }
  property {
    name  = "mail"
    type  = "string"
    value = var.azureusermanagementconnector_property_mail
  }
  property {
    name  = "mailNickname"
    type  = "string"
    value = var.azureusermanagementconnector_property_mail_nickname
  }
  property {
    name  = "method"
    type  = "string"
    value = var.azureusermanagementconnector_property_method
  }
  property {
    name  = "password"
    type  = "string"
    value = var.azureusermanagementconnector_property_password
  }
  property {
    name  = "removeLicenses"
    type  = "string"
    value = var.azureusermanagementconnector_property_remove_licenses
  }
  property {
    name  = "securityEnabledOnly"
    type  = "string"
    value = var.azureusermanagementconnector_property_security_enabled_only
  }
  property {
    name  = "skuId"
    type  = "string"
    value = var.azureusermanagementconnector_property_sku_id
  }
  property {
    name  = "state"
    type  = "string"
    value = var.azureusermanagementconnector_property_state
  }
  property {
    name  = "surname"
    type  = "string"
    value = var.azureusermanagementconnector_property_surname
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.azureusermanagementconnector_property_user_id
  }
  property {
    name  = "userPrincipalName"
    type  = "string"
    value = var.azureusermanagementconnector_property_user_principal_name
  }
  property {
    name  = "userQuery"
    type  = "string"
    value = var.azureusermanagementconnector_property_user_query
  }
  property {
    name  = "users"
    type  = "string"
    value = var.azureusermanagementconnector_property_users
  }
}
