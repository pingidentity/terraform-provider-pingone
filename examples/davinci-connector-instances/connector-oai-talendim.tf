resource "pingone_davinci_connector_instance" "connector-oai-talendim" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-talendim"
  }
  name = "My awesome connector-oai-talendim"
  property {
    name  = "addUserToGroups_id"
    type  = "string"
    value = var.connector-oai-talendim_property_add_user_to_groups_id
  }
  property {
    name  = "addUserToGroups_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_add_user_to_groups_request_body
  }
  property {
    name  = "assignRoleToUsers_id"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_role_to_users_id
  }
  property {
    name  = "assignRoleToUsers_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_role_to_users_request_body
  }
  property {
    name  = "assignRoles_id"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_roles_id
  }
  property {
    name  = "assignRoles_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_roles_request_body
  }
  property {
    name  = "assignUser1_id"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_user1_id
  }
  property {
    name  = "assignUser_id"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_user_id
  }
  property {
    name  = "assignUsersToAGroup_id"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_users_to_agroup_id
  }
  property {
    name  = "assignUsersToAGroup_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_assign_users_to_agroup_request_body
  }
  property {
    name  = "authBearerToken"
    type  = "string"
    value = var.connector-oai-talendim_property_auth_bearer_token
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-talendim_property_base_path
  }
  property {
    name  = "createANewUser1_admin"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_admin
  }
  property {
    name  = "createANewUser1_invite"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_invite
  }
  property {
    name  = "createANewUser1_user_User_active"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_active
  }
  property {
    name  = "createANewUser1_user_User_email"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_email
  }
  property {
    name  = "createANewUser1_user_User_firstName"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_first_name
  }
  property {
    name  = "createANewUser1_user_User_lastName"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_last_name
  }
  property {
    name  = "createANewUser1_user_User_login"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_login
  }
  property {
    name  = "createANewUser1_user_User_password"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_password
  }
  property {
    name  = "createANewUser1_user_User_phone"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_phone
  }
  property {
    name  = "createANewUser1_user_User_preferredLanguage"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_preferred_language
  }
  property {
    name  = "createANewUser1_user_User_roleIds"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_role_ids
  }
  property {
    name  = "createANewUser1_user_User_timezone"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_timezone
  }
  property {
    name  = "createANewUser1_user_User_title"
    type  = "string"
    value = var.connector-oai-talendim_property_create_anew_user1_user_user_title
  }
  property {
    name  = "createPermissions1_permission"
    type  = "string"
    value = var.connector-oai-talendim_property_create_permissions1_permission
  }
  property {
    name  = "createPermissions2_permission"
    type  = "string"
    value = var.connector-oai-talendim_property_create_permissions2_permission
  }
  property {
    name  = "createWorkspacePermissions1_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions1_request_body
  }
  property {
    name  = "createWorkspacePermissions1_service_account_id"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions1_service_account_id
  }
  property {
    name  = "createWorkspacePermissions1_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions1_workspace_id
  }
  property {
    name  = "createWorkspacePermissions2_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions2_request_body
  }
  property {
    name  = "createWorkspacePermissions2_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions2_user_id
  }
  property {
    name  = "createWorkspacePermissions2_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_create_workspace_permissions2_workspace_id
  }
  property {
    name  = "deleteGroup_id"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_group_id
  }
  property {
    name  = "deletePermissions1_service_account_ids"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_permissions1_service_account_ids
  }
  property {
    name  = "deletePermissions1_workspace_ids"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_permissions1_workspace_ids
  }
  property {
    name  = "deletePermissions2_user_ids"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_permissions2_user_ids
  }
  property {
    name  = "deletePermissions2_workspace_ids"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_permissions2_workspace_ids
  }
  property {
    name  = "deleteWorkspacePermissions1_service_account_id"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_workspace_permissions1_service_account_id
  }
  property {
    name  = "deleteWorkspacePermissions1_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_workspace_permissions1_workspace_id
  }
  property {
    name  = "deleteWorkspacePermissions2_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_workspace_permissions2_user_id
  }
  property {
    name  = "deleteWorkspacePermissions2_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_delete_workspace_permissions2_workspace_id
  }
  property {
    name  = "getPermissions1_environment_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions1_environment_id
  }
  property {
    name  = "getPermissions1_service_account_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions1_service_account_id
  }
  property {
    name  = "getPermissions1_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions1_workspace_id
  }
  property {
    name  = "getPermissions2_environment_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions2_environment_id
  }
  property {
    name  = "getPermissions2_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions2_user_id
  }
  property {
    name  = "getPermissions2_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_permissions2_workspace_id
  }
  property {
    name  = "getRoleById_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_role_by_id_id
  }
  property {
    name  = "getWorkspacePermission1_service_account_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_workspace_permission1_service_account_id
  }
  property {
    name  = "getWorkspacePermission1_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_workspace_permission1_workspace_id
  }
  property {
    name  = "getWorkspacePermission2_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_workspace_permission2_user_id
  }
  property {
    name  = "getWorkspacePermission2_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_get_workspace_permission2_workspace_id
  }
  property {
    name  = "listUserGroups_id"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_groups_id
  }
  property {
    name  = "listUserGroups_name_filter"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_groups_name_filter
  }
  property {
    name  = "listUserGroups_page"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_groups_page
  }
  property {
    name  = "listUserGroups_size"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_groups_size
  }
  property {
    name  = "listUserRoles_id"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_roles_id
  }
  property {
    name  = "listUserRoles_page"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_roles_page
  }
  property {
    name  = "listUserRoles_size"
    type  = "string"
    value = var.connector-oai-talendim_property_list_user_roles_size
  }
  property {
    name  = "registerANewGroup_group_Group_name"
    type  = "string"
    value = var.connector-oai-talendim_property_register_anew_group_group_group_name
  }
  property {
    name  = "registerANewGroup_group_Group_userIds"
    type  = "string"
    value = var.connector-oai-talendim_property_register_anew_group_group_group_user_ids
  }
  property {
    name  = "registerANewRole_role_Role_name"
    type  = "string"
    value = var.connector-oai-talendim_property_register_anew_role_role_role_name
  }
  property {
    name  = "registerANewRole_role_Role_permissions"
    type  = "string"
    value = var.connector-oai-talendim_property_register_anew_role_role_role_permissions
  }
  property {
    name  = "removeRole1_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_role1_id
  }
  property {
    name  = "removeRole1_role_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_role1_role_id
  }
  property {
    name  = "removeRole_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_role_id
  }
  property {
    name  = "removeRole_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_role_user_id
  }
  property {
    name  = "removeUserFromGroup_group_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_user_from_group_group_id
  }
  property {
    name  = "removeUserFromGroup_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_user_from_group_id
  }
  property {
    name  = "removeUser_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_user_id
  }
  property {
    name  = "removeUser_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_remove_user_user_id
  }
  property {
    name  = "renameGroup_id"
    type  = "string"
    value = var.connector-oai-talendim_property_rename_group_id
  }
  property {
    name  = "renameGroup_renameGroupRequest_RenameGroupRequest_name"
    type  = "string"
    value = var.connector-oai-talendim_property_rename_group_rename_group_request_rename_group_request_name
  }
  property {
    name  = "retrieveGroupUsers_id"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_group_users_id
  }
  property {
    name  = "retrieveGroupUsers_page"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_group_users_page
  }
  property {
    name  = "retrieveGroupUsers_size"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_group_users_size
  }
  property {
    name  = "retrieveGroup_id"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_group_id
  }
  property {
    name  = "retrieveGroups_name_filter"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_groups_name_filter
  }
  property {
    name  = "retrieveGroups_page"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_groups_page
  }
  property {
    name  = "retrieveGroups_size"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_groups_size
  }
  property {
    name  = "retrieveRoles_name_filter"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_roles_name_filter
  }
  property {
    name  = "retrieveRoles_page"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_roles_page
  }
  property {
    name  = "retrieveRoles_size"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_roles_size
  }
  property {
    name  = "retrieveUserDetails_id"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_user_details_id
  }
  property {
    name  = "retrieveUsers1_page"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_users1_page
  }
  property {
    name  = "retrieveUsers1_size"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_users1_size
  }
  property {
    name  = "retrieveUsers_id"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_users_id
  }
  property {
    name  = "retrieveUsers_page"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_users_page
  }
  property {
    name  = "retrieveUsers_size"
    type  = "string"
    value = var.connector-oai-talendim_property_retrieve_users_size
  }
  property {
    name  = "securityRoleMappingsCustomerRoleNameDelete_customer_role_name"
    type  = "string"
    value = var.connector-oai-talendim_property_security_role_mappings_customer_role_name_delete_customer_role_name
  }
  property {
    name  = "securityRoleMappingsCustomerRoleNameGet_customer_role_name"
    type  = "string"
    value = var.connector-oai-talendim_property_security_role_mappings_customer_role_name_get_customer_role_name
  }
  property {
    name  = "updatePermissions1_permission"
    type  = "string"
    value = var.connector-oai-talendim_property_update_permissions1_permission
  }
  property {
    name  = "updatePermissions2_permission"
    type  = "string"
    value = var.connector-oai-talendim_property_update_permissions2_permission
  }
  property {
    name  = "updateRole_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_role_id
  }
  property {
    name  = "updateRole_role_Role_name"
    type  = "string"
    value = var.connector-oai-talendim_property_update_role_role_role_name
  }
  property {
    name  = "updateRole_role_Role_permissions"
    type  = "string"
    value = var.connector-oai-talendim_property_update_role_role_role_permissions
  }
  property {
    name  = "updateUser1_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_id
  }
  property {
    name  = "updateUser1_user_User_active"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_active
  }
  property {
    name  = "updateUser1_user_User_email"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_email
  }
  property {
    name  = "updateUser1_user_User_firstName"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_first_name
  }
  property {
    name  = "updateUser1_user_User_lastName"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_last_name
  }
  property {
    name  = "updateUser1_user_User_login"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_login
  }
  property {
    name  = "updateUser1_user_User_password"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_password
  }
  property {
    name  = "updateUser1_user_User_phone"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_phone
  }
  property {
    name  = "updateUser1_user_User_preferredLanguage"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_preferred_language
  }
  property {
    name  = "updateUser1_user_User_roleIds"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_role_ids
  }
  property {
    name  = "updateUser1_user_User_timezone"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_timezone
  }
  property {
    name  = "updateUser1_user_User_title"
    type  = "string"
    value = var.connector-oai-talendim_property_update_user1_user_user_title
  }
  property {
    name  = "updateWorkspacePermissions1_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions1_request_body
  }
  property {
    name  = "updateWorkspacePermissions1_service_account_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions1_service_account_id
  }
  property {
    name  = "updateWorkspacePermissions1_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions1_workspace_id
  }
  property {
    name  = "updateWorkspacePermissions2_request_body"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions2_request_body
  }
  property {
    name  = "updateWorkspacePermissions2_user_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions2_user_id
  }
  property {
    name  = "updateWorkspacePermissions2_workspace_id"
    type  = "string"
    value = var.connector-oai-talendim_property_update_workspace_permissions2_workspace_id
  }
}
