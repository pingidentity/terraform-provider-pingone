resource "pingone_davinci_connector_instance" "connector-oai-talendscim" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-talendscim"
  }
  name = "My awesome connector-oai-talendscim"
  property {
    name  = "_delete_id"
    type  = "string"
    value = var.connector-oai-talendscim_property__delete_id
  }
  property {
    name  = "authBearerToken"
    type  = "string"
    value = var.connector-oai-talendscim_property_auth_bearer_token
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-talendscim_property_base_path
  }
  property {
    name  = "create1_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_attributes
  }
  property {
    name  = "create1_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_excluded_attributes
  }
  property {
    name  = "create1_roleResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_meta_created
  }
  property {
    name  = "create1_roleResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_meta_last_modified
  }
  property {
    name  = "create1_roleResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_meta_location
  }
  property {
    name  = "create1_roleResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_meta_resource_type
  }
  property {
    name  = "create1_roleResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_meta_version
  }
  property {
    name  = "create1_roleResource_RoleResource_entitlements"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_role_resource_entitlements
  }
  property {
    name  = "create1_roleResource_RoleResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_role_resource_external_id
  }
  property {
    name  = "create1_roleResource_RoleResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_role_resource_id
  }
  property {
    name  = "create1_roleResource_RoleResource_name"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_role_resource_name
  }
  property {
    name  = "create1_roleResource_RoleResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_create1_role_resource_role_resource_schemas
  }
  property {
    name  = "create2_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_attributes
  }
  property {
    name  = "create2_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_excluded_attributes
  }
  property {
    name  = "create2_groupResource_GroupResource_displayName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_group_resource_display_name
  }
  property {
    name  = "create2_groupResource_GroupResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_group_resource_external_id
  }
  property {
    name  = "create2_groupResource_GroupResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_group_resource_id
  }
  property {
    name  = "create2_groupResource_GroupResource_members"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_group_resource_members
  }
  property {
    name  = "create2_groupResource_GroupResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_group_resource_schemas
  }
  property {
    name  = "create2_groupResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_meta_created
  }
  property {
    name  = "create2_groupResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_meta_last_modified
  }
  property {
    name  = "create2_groupResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_meta_location
  }
  property {
    name  = "create2_groupResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_meta_resource_type
  }
  property {
    name  = "create2_groupResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_create2_group_resource_meta_version
  }
  property {
    name  = "create_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_attributes
  }
  property {
    name  = "create_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_excluded_attributes
  }
  property {
    name  = "create_userResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_meta_created
  }
  property {
    name  = "create_userResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_meta_last_modified
  }
  property {
    name  = "create_userResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_meta_location
  }
  property {
    name  = "create_userResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_meta_resource_type
  }
  property {
    name  = "create_userResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_meta_version
  }
  property {
    name  = "create_userResource_Name_familyName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_family_name
  }
  property {
    name  = "create_userResource_Name_formatted"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_formatted
  }
  property {
    name  = "create_userResource_Name_givenName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_given_name
  }
  property {
    name  = "create_userResource_Name_honorificPrefix"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_honorific_prefix
  }
  property {
    name  = "create_userResource_Name_honorificSuffix"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_honorific_suffix
  }
  property {
    name  = "create_userResource_Name_middleName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_name_middle_name
  }
  property {
    name  = "create_userResource_UserResource_active"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_active
  }
  property {
    name  = "create_userResource_UserResource_addresses"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_addresses
  }
  property {
    name  = "create_userResource_UserResource_displayName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_display_name
  }
  property {
    name  = "create_userResource_UserResource_emails"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_emails
  }
  property {
    name  = "create_userResource_UserResource_entitlements"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_entitlements
  }
  property {
    name  = "create_userResource_UserResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_external_id
  }
  property {
    name  = "create_userResource_UserResource_groups"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_groups
  }
  property {
    name  = "create_userResource_UserResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_id
  }
  property {
    name  = "create_userResource_UserResource_ims"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_ims
  }
  property {
    name  = "create_userResource_UserResource_locale"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_locale
  }
  property {
    name  = "create_userResource_UserResource_nickName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_nick_name
  }
  property {
    name  = "create_userResource_UserResource_password"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_password
  }
  property {
    name  = "create_userResource_UserResource_phoneNumbers"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_phone_numbers
  }
  property {
    name  = "create_userResource_UserResource_photos"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_photos
  }
  property {
    name  = "create_userResource_UserResource_preferredLanguage"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_preferred_language
  }
  property {
    name  = "create_userResource_UserResource_profileUrl"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_profile_url
  }
  property {
    name  = "create_userResource_UserResource_roles"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_roles
  }
  property {
    name  = "create_userResource_UserResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_schemas
  }
  property {
    name  = "create_userResource_UserResource_timezone"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_timezone
  }
  property {
    name  = "create_userResource_UserResource_title"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_title
  }
  property {
    name  = "create_userResource_UserResource_userName"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_user_name
  }
  property {
    name  = "create_userResource_UserResource_userType"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_user_type
  }
  property {
    name  = "create_userResource_UserResource_x509Certificates"
    type  = "string"
    value = var.connector-oai-talendscim_property_create_user_resource_user_resource_x509_certificates
  }
  property {
    name  = "delete1_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_delete1_id
  }
  property {
    name  = "delete2_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_delete2_id
  }
  property {
    name  = "get1_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get1_attributes
  }
  property {
    name  = "get1_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get1_excluded_attributes
  }
  property {
    name  = "get1_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_get1_id
  }
  property {
    name  = "get2_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get2_attributes
  }
  property {
    name  = "get2_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get2_excluded_attributes
  }
  property {
    name  = "get2_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_get2_id
  }
  property {
    name  = "get4_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_get4_id
  }
  property {
    name  = "get5_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_get5_id
  }
  property {
    name  = "getAll1_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_attributes
  }
  property {
    name  = "getAll1_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_count
  }
  property {
    name  = "getAll1_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_excluded_attributes
  }
  property {
    name  = "getAll1_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_filter
  }
  property {
    name  = "getAll1_sort_by"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_sort_by
  }
  property {
    name  = "getAll1_sort_order"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_sort_order
  }
  property {
    name  = "getAll1_start_index"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all1_start_index
  }
  property {
    name  = "getAll2_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_attributes
  }
  property {
    name  = "getAll2_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_count
  }
  property {
    name  = "getAll2_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_excluded_attributes
  }
  property {
    name  = "getAll2_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_filter
  }
  property {
    name  = "getAll2_sort_by"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_sort_by
  }
  property {
    name  = "getAll2_sort_order"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_sort_order
  }
  property {
    name  = "getAll2_start_index"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all2_start_index
  }
  property {
    name  = "getAll_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_attributes
  }
  property {
    name  = "getAll_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_count
  }
  property {
    name  = "getAll_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_excluded_attributes
  }
  property {
    name  = "getAll_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_filter
  }
  property {
    name  = "getAll_sort_by"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_sort_by
  }
  property {
    name  = "getAll_sort_order"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_sort_order
  }
  property {
    name  = "getAll_start_index"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_all_start_index
  }
  property {
    name  = "get_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_attributes
  }
  property {
    name  = "get_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_excluded_attributes
  }
  property {
    name  = "get_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_get_id
  }
  property {
    name  = "patch1_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_attributes
  }
  property {
    name  = "patch1_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_excluded_attributes
  }
  property {
    name  = "patch1_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_id
  }
  property {
    name  = "patch1_patchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_meta_created
  }
  property {
    name  = "patch1_patchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_meta_last_modified
  }
  property {
    name  = "patch1_patchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_meta_location
  }
  property {
    name  = "patch1_patchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_meta_resource_type
  }
  property {
    name  = "patch1_patchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_meta_version
  }
  property {
    name  = "patch1_patchRequest_PatchRequest_Operations"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_patch_request_operations
  }
  property {
    name  = "patch1_patchRequest_PatchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_patch_request_external_id
  }
  property {
    name  = "patch1_patchRequest_PatchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_patch_request_id
  }
  property {
    name  = "patch1_patchRequest_PatchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch1_patch_request_patch_request_schemas
  }
  property {
    name  = "patch2_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_attributes
  }
  property {
    name  = "patch2_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_excluded_attributes
  }
  property {
    name  = "patch2_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_id
  }
  property {
    name  = "patch2_patchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_meta_created
  }
  property {
    name  = "patch2_patchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_meta_last_modified
  }
  property {
    name  = "patch2_patchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_meta_location
  }
  property {
    name  = "patch2_patchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_meta_resource_type
  }
  property {
    name  = "patch2_patchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_meta_version
  }
  property {
    name  = "patch2_patchRequest_PatchRequest_Operations"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_patch_request_operations
  }
  property {
    name  = "patch2_patchRequest_PatchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_patch_request_external_id
  }
  property {
    name  = "patch2_patchRequest_PatchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_patch_request_id
  }
  property {
    name  = "patch2_patchRequest_PatchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch2_patch_request_patch_request_schemas
  }
  property {
    name  = "patch_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_attributes
  }
  property {
    name  = "patch_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_excluded_attributes
  }
  property {
    name  = "patch_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_id
  }
  property {
    name  = "patch_patchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_meta_created
  }
  property {
    name  = "patch_patchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_meta_last_modified
  }
  property {
    name  = "patch_patchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_meta_location
  }
  property {
    name  = "patch_patchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_meta_resource_type
  }
  property {
    name  = "patch_patchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_meta_version
  }
  property {
    name  = "patch_patchRequest_PatchRequest_Operations"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_patch_request_operations
  }
  property {
    name  = "patch_patchRequest_PatchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_patch_request_external_id
  }
  property {
    name  = "patch_patchRequest_PatchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_patch_request_id
  }
  property {
    name  = "patch_patchRequest_PatchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_patch_patch_request_patch_request_schemas
  }
  property {
    name  = "search1_searchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_meta_created
  }
  property {
    name  = "search1_searchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_meta_last_modified
  }
  property {
    name  = "search1_searchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_meta_location
  }
  property {
    name  = "search1_searchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_meta_resource_type
  }
  property {
    name  = "search1_searchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_meta_version
  }
  property {
    name  = "search1_searchRequest_SearchRequest_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_attributes
  }
  property {
    name  = "search1_searchRequest_SearchRequest_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_count
  }
  property {
    name  = "search1_searchRequest_SearchRequest_excludedAttributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_excluded_attributes
  }
  property {
    name  = "search1_searchRequest_SearchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_external_id
  }
  property {
    name  = "search1_searchRequest_SearchRequest_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_filter
  }
  property {
    name  = "search1_searchRequest_SearchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_id
  }
  property {
    name  = "search1_searchRequest_SearchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_schemas
  }
  property {
    name  = "search1_searchRequest_SearchRequest_sortBy"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_sort_by
  }
  property {
    name  = "search1_searchRequest_SearchRequest_sortOrder"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_sort_order
  }
  property {
    name  = "search1_searchRequest_SearchRequest_startIndex"
    type  = "string"
    value = var.connector-oai-talendscim_property_search1_search_request_search_request_start_index
  }
  property {
    name  = "search2_searchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_meta_created
  }
  property {
    name  = "search2_searchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_meta_last_modified
  }
  property {
    name  = "search2_searchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_meta_location
  }
  property {
    name  = "search2_searchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_meta_resource_type
  }
  property {
    name  = "search2_searchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_meta_version
  }
  property {
    name  = "search2_searchRequest_SearchRequest_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_attributes
  }
  property {
    name  = "search2_searchRequest_SearchRequest_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_count
  }
  property {
    name  = "search2_searchRequest_SearchRequest_excludedAttributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_excluded_attributes
  }
  property {
    name  = "search2_searchRequest_SearchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_external_id
  }
  property {
    name  = "search2_searchRequest_SearchRequest_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_filter
  }
  property {
    name  = "search2_searchRequest_SearchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_id
  }
  property {
    name  = "search2_searchRequest_SearchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_schemas
  }
  property {
    name  = "search2_searchRequest_SearchRequest_sortBy"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_sort_by
  }
  property {
    name  = "search2_searchRequest_SearchRequest_sortOrder"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_sort_order
  }
  property {
    name  = "search2_searchRequest_SearchRequest_startIndex"
    type  = "string"
    value = var.connector-oai-talendscim_property_search2_search_request_search_request_start_index
  }
  property {
    name  = "search3_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_search3_filter
  }
  property {
    name  = "search4_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_search4_filter
  }
  property {
    name  = "search_searchRequest_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_meta_created
  }
  property {
    name  = "search_searchRequest_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_meta_last_modified
  }
  property {
    name  = "search_searchRequest_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_meta_location
  }
  property {
    name  = "search_searchRequest_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_meta_resource_type
  }
  property {
    name  = "search_searchRequest_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_meta_version
  }
  property {
    name  = "search_searchRequest_SearchRequest_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_attributes
  }
  property {
    name  = "search_searchRequest_SearchRequest_count"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_count
  }
  property {
    name  = "search_searchRequest_SearchRequest_excludedAttributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_excluded_attributes
  }
  property {
    name  = "search_searchRequest_SearchRequest_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_external_id
  }
  property {
    name  = "search_searchRequest_SearchRequest_filter"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_filter
  }
  property {
    name  = "search_searchRequest_SearchRequest_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_id
  }
  property {
    name  = "search_searchRequest_SearchRequest_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_schemas
  }
  property {
    name  = "search_searchRequest_SearchRequest_sortBy"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_sort_by
  }
  property {
    name  = "search_searchRequest_SearchRequest_sortOrder"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_sort_order
  }
  property {
    name  = "search_searchRequest_SearchRequest_startIndex"
    type  = "string"
    value = var.connector-oai-talendscim_property_search_search_request_search_request_start_index
  }
  property {
    name  = "update1_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_attributes
  }
  property {
    name  = "update1_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_excluded_attributes
  }
  property {
    name  = "update1_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_id
  }
  property {
    name  = "update1_roleResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_meta_created
  }
  property {
    name  = "update1_roleResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_meta_last_modified
  }
  property {
    name  = "update1_roleResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_meta_location
  }
  property {
    name  = "update1_roleResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_meta_resource_type
  }
  property {
    name  = "update1_roleResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_meta_version
  }
  property {
    name  = "update1_roleResource_RoleResource_entitlements"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_role_resource_entitlements
  }
  property {
    name  = "update1_roleResource_RoleResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_role_resource_external_id
  }
  property {
    name  = "update1_roleResource_RoleResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_role_resource_id
  }
  property {
    name  = "update1_roleResource_RoleResource_name"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_role_resource_name
  }
  property {
    name  = "update1_roleResource_RoleResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_update1_role_resource_role_resource_schemas
  }
  property {
    name  = "update2_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_attributes
  }
  property {
    name  = "update2_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_excluded_attributes
  }
  property {
    name  = "update2_groupResource_GroupResource_displayName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_group_resource_display_name
  }
  property {
    name  = "update2_groupResource_GroupResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_group_resource_external_id
  }
  property {
    name  = "update2_groupResource_GroupResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_group_resource_id
  }
  property {
    name  = "update2_groupResource_GroupResource_members"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_group_resource_members
  }
  property {
    name  = "update2_groupResource_GroupResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_group_resource_schemas
  }
  property {
    name  = "update2_groupResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_meta_created
  }
  property {
    name  = "update2_groupResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_meta_last_modified
  }
  property {
    name  = "update2_groupResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_meta_location
  }
  property {
    name  = "update2_groupResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_meta_resource_type
  }
  property {
    name  = "update2_groupResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_group_resource_meta_version
  }
  property {
    name  = "update2_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update2_id
  }
  property {
    name  = "update_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_attributes
  }
  property {
    name  = "update_excluded_attributes"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_excluded_attributes
  }
  property {
    name  = "update_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_id
  }
  property {
    name  = "update_userResource_Meta_created"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_meta_created
  }
  property {
    name  = "update_userResource_Meta_lastModified"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_meta_last_modified
  }
  property {
    name  = "update_userResource_Meta_location"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_meta_location
  }
  property {
    name  = "update_userResource_Meta_resourceType"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_meta_resource_type
  }
  property {
    name  = "update_userResource_Meta_version"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_meta_version
  }
  property {
    name  = "update_userResource_Name_familyName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_family_name
  }
  property {
    name  = "update_userResource_Name_formatted"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_formatted
  }
  property {
    name  = "update_userResource_Name_givenName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_given_name
  }
  property {
    name  = "update_userResource_Name_honorificPrefix"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_honorific_prefix
  }
  property {
    name  = "update_userResource_Name_honorificSuffix"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_honorific_suffix
  }
  property {
    name  = "update_userResource_Name_middleName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_name_middle_name
  }
  property {
    name  = "update_userResource_UserResource_active"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_active
  }
  property {
    name  = "update_userResource_UserResource_addresses"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_addresses
  }
  property {
    name  = "update_userResource_UserResource_displayName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_display_name
  }
  property {
    name  = "update_userResource_UserResource_emails"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_emails
  }
  property {
    name  = "update_userResource_UserResource_entitlements"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_entitlements
  }
  property {
    name  = "update_userResource_UserResource_externalId"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_external_id
  }
  property {
    name  = "update_userResource_UserResource_groups"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_groups
  }
  property {
    name  = "update_userResource_UserResource_id"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_id
  }
  property {
    name  = "update_userResource_UserResource_ims"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_ims
  }
  property {
    name  = "update_userResource_UserResource_locale"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_locale
  }
  property {
    name  = "update_userResource_UserResource_nickName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_nick_name
  }
  property {
    name  = "update_userResource_UserResource_password"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_password
  }
  property {
    name  = "update_userResource_UserResource_phoneNumbers"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_phone_numbers
  }
  property {
    name  = "update_userResource_UserResource_photos"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_photos
  }
  property {
    name  = "update_userResource_UserResource_preferredLanguage"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_preferred_language
  }
  property {
    name  = "update_userResource_UserResource_profileUrl"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_profile_url
  }
  property {
    name  = "update_userResource_UserResource_roles"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_roles
  }
  property {
    name  = "update_userResource_UserResource_schemas"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_schemas
  }
  property {
    name  = "update_userResource_UserResource_timezone"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_timezone
  }
  property {
    name  = "update_userResource_UserResource_title"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_title
  }
  property {
    name  = "update_userResource_UserResource_userName"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_user_name
  }
  property {
    name  = "update_userResource_UserResource_userType"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_user_type
  }
  property {
    name  = "update_userResource_UserResource_x509Certificates"
    type  = "string"
    value = var.connector-oai-talendscim_property_update_user_resource_user_resource_x509_certificates
  }
}
