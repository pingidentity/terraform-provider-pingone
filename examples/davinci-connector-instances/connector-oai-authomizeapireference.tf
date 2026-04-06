resource "pingone_davinci_connector_instance" "connector-oai-authomizeapireference" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-authomizeapireference"
  }
  name = "My awesome connector-oai-authomizeapireference"
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_auth_api_key
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_base_path
  }
  property {
    name  = "createAccountsAssociationV2AppsAppIdAssociationAccountsPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_accounts_association_v2_apps_app_id_association_accounts_post_app_id
  }
  property {
    name  = "createAccountsAssociationV2AppsAppIdAssociationAccountsPost_newAccountsAssociationsListRequestSchema_NewAccountsAssociationsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_accounts_association_v2_apps_app_id_association_accounts_post_new_accounts_associations_list_request_schema_new_accounts_associations_list_request_schema_data
  }
  property {
    name  = "createAssetsInheritanceV2AppsAppIdAssetsInheritancePost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_assets_inheritance_v2_apps_app_id_assets_inheritance_post_app_id
  }
  property {
    name  = "createAssetsInheritanceV2AppsAppIdAssetsInheritancePost_newAssetsInheritanceListRequestSchema_NewAssetsInheritanceListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_assets_inheritance_v2_apps_app_id_assets_inheritance_post_new_assets_inheritance_list_request_schema_new_assets_inheritance_list_request_schema_data
  }
  property {
    name  = "createAssetsV2AppsAppIdAssetsPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_assets_v2_apps_app_id_assets_post_app_id
  }
  property {
    name  = "createAssetsV2AppsAppIdAssetsPost_newAssetsListRequestSchema_NewAssetsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_assets_v2_apps_app_id_assets_post_new_assets_list_request_schema_new_assets_list_request_schema_data
  }
  property {
    name  = "createGroupingsAssociationV2AppsAppIdAssociationGroupingsPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_groupings_association_v2_apps_app_id_association_groupings_post_app_id
  }
  property {
    name  = "createGroupingsAssociationV2AppsAppIdAssociationGroupingsPost_newGroupingsAssociationsListRequestSchema_NewGroupingsAssociationsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_groupings_association_v2_apps_app_id_association_groupings_post_new_groupings_associations_list_request_schema_new_groupings_associations_list_request_schema_data
  }
  property {
    name  = "createGroupingsV2AppsAppIdAccessGroupingPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_groupings_v2_apps_app_id_access_grouping_post_app_id
  }
  property {
    name  = "createGroupingsV2AppsAppIdAccessGroupingPost_newGroupingsListRequestSchema_NewGroupingsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_groupings_v2_apps_app_id_access_grouping_post_new_groupings_list_request_schema_new_groupings_list_request_schema_data
  }
  property {
    name  = "createIdentitiesV2AppsAppIdIdentitiesPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_identities_v2_apps_app_id_identities_post_app_id
  }
  property {
    name  = "createIdentitiesV2AppsAppIdIdentitiesPost_newIdentitiesListRequestSchema_NewIdentitiesListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_identities_v2_apps_app_id_identities_post_new_identities_list_request_schema_new_identities_list_request_schema_data
  }
  property {
    name  = "createPermissionsV2AppsAppIdAccessPermissionsPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_permissions_v2_apps_app_id_access_permissions_post_app_id
  }
  property {
    name  = "createPermissionsV2AppsAppIdAccessPermissionsPost_newPermissionsListRequestSchema_NewPermissionsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_permissions_v2_apps_app_id_access_permissions_post_new_permissions_list_request_schema_new_permissions_list_request_schema_data
  }
  property {
    name  = "createPrivilegesGrantsV2AppsAppIdPrivilegesGrantsPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_privileges_grants_v2_apps_app_id_privileges_grants_post_app_id
  }
  property {
    name  = "createPrivilegesGrantsV2AppsAppIdPrivilegesGrantsPost_newPrivilegesGrantsListRequestSchema_NewPrivilegesGrantsListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_privileges_grants_v2_apps_app_id_privileges_grants_post_new_privileges_grants_list_request_schema_new_privileges_grants_list_request_schema_data
  }
  property {
    name  = "createPrivilegesV2AppsAppIdPrivilegesPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_privileges_v2_apps_app_id_privileges_post_app_id
  }
  property {
    name  = "createPrivilegesV2AppsAppIdPrivilegesPost_newPrivilegesListRequestSchema_NewPrivilegesListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_privileges_v2_apps_app_id_privileges_post_new_privileges_list_request_schema_new_privileges_list_request_schema_data
  }
  property {
    name  = "createUsersV2AppsAppIdAccountsUsersPost_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_users_v2_apps_app_id_accounts_users_post_app_id
  }
  property {
    name  = "createUsersV2AppsAppIdAccountsUsersPost_newUsersListRequestSchema_NewUsersListRequestSchema_data"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_create_users_v2_apps_app_id_accounts_users_post_new_users_list_request_schema_new_users_list_request_schema_data
  }
  property {
    name  = "deleteApplicationDataV2AppsAppIdDataDelete_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_delete_application_data_v2_apps_app_id_data_delete_app_id
  }
  property {
    name  = "searchAccountsAssociationV2AppsAppIdAssociationAccountsGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_accounts_association_v2_apps_app_id_association_accounts_get_app_id
  }
  property {
    name  = "searchAccountsAssociationV2AppsAppIdAssociationAccountsGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_accounts_association_v2_apps_app_id_association_accounts_get_start_date
  }
  property {
    name  = "searchAssetsInheritanceV2AppsAppIdAssetsInheritanceGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_assets_inheritance_v2_apps_app_id_assets_inheritance_get_app_id
  }
  property {
    name  = "searchAssetsInheritanceV2AppsAppIdAssetsInheritanceGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_assets_inheritance_v2_apps_app_id_assets_inheritance_get_start_date
  }
  property {
    name  = "searchAssetsV2AppsAppIdAssetsGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_assets_v2_apps_app_id_assets_get_app_id
  }
  property {
    name  = "searchAssetsV2AppsAppIdAssetsGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_assets_v2_apps_app_id_assets_get_start_date
  }
  property {
    name  = "searchGroupingsAssociationV2AppsAppIdAssociationGroupingsGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_groupings_association_v2_apps_app_id_association_groupings_get_app_id
  }
  property {
    name  = "searchGroupingsAssociationV2AppsAppIdAssociationGroupingsGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_groupings_association_v2_apps_app_id_association_groupings_get_start_date
  }
  property {
    name  = "searchGroupingsV2AppsAppIdAccessGroupingGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_groupings_v2_apps_app_id_access_grouping_get_app_id
  }
  property {
    name  = "searchGroupingsV2AppsAppIdAccessGroupingGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_groupings_v2_apps_app_id_access_grouping_get_start_date
  }
  property {
    name  = "searchIdentitiesV2AppsAppIdIdentitiesGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_identities_v2_apps_app_id_identities_get_app_id
  }
  property {
    name  = "searchIdentitiesV2AppsAppIdIdentitiesGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_identities_v2_apps_app_id_identities_get_start_date
  }
  property {
    name  = "searchPermissionsV2AppsAppIdAccessPermissionsGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_permissions_v2_apps_app_id_access_permissions_get_app_id
  }
  property {
    name  = "searchPermissionsV2AppsAppIdAccessPermissionsGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_permissions_v2_apps_app_id_access_permissions_get_start_date
  }
  property {
    name  = "searchPrivilegesGrantsV2AppsAppIdPrivilegesGrantsGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_privileges_grants_v2_apps_app_id_privileges_grants_get_app_id
  }
  property {
    name  = "searchPrivilegesGrantsV2AppsAppIdPrivilegesGrantsGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_privileges_grants_v2_apps_app_id_privileges_grants_get_start_date
  }
  property {
    name  = "searchPrivilegesV2AppsAppIdPrivilegesGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_privileges_v2_apps_app_id_privileges_get_app_id
  }
  property {
    name  = "searchPrivilegesV2AppsAppIdPrivilegesGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_privileges_v2_apps_app_id_privileges_get_start_date
  }
  property {
    name  = "searchUsersV2AppsAppIdAccountsUsersGet_app_id"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_users_v2_apps_app_id_accounts_users_get_app_id
  }
  property {
    name  = "searchUsersV2AppsAppIdAccountsUsersGet_start_date"
    type  = "string"
    value = var.connector-oai-authomizeapireference_property_search_users_v2_apps_app_id_accounts_users_get_start_date
  }
}
