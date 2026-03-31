resource "pingone_davinci_connector_instance" "connector-oai-druvainsynccloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-druvainsynccloud"
  }
  name = "My awesome connector-oai-druvainsynccloud"
  property {
    name  = "authClientId"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_auth_client_id
  }
  property {
    name  = "authClientSecret"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_auth_client_secret
  }
  property {
    name  = "authTokenUrl"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_auth_token_url
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_base_path
  }
  property {
    name  = "clientsGet_connection_status"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_connection_status
  }
  property {
    name  = "clientsGet_page_number"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_page_number
  }
  property {
    name  = "clientsGet_page_size"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_page_size
  }
  property {
    name  = "clientsGet_search_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_search_by
  }
  property {
    name  = "clientsGet_sort_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_sort_by
  }
  property {
    name  = "clientsGet_sort_order"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_clients_get_sort_order
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_clientName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_client_name
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_createdBy"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_created_by
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_custodianEmails"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_custodian_emails
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_downloadLocation"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_download_location
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_downloadOption"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_download_option
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_enable_data_integrity"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_enable_data_integrity
  }
  property {
    name  = "createNewJob_createDownloadJobRequest_CreateDownloadJobRequest_legalholdId"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_create_new_job_create_download_job_request_create_download_job_request_legalhold_id
  }
  property {
    name  = "endpointsV1BackupsGet_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_device_id
  }
  property {
    name  = "endpointsV1BackupsGet_last_successful"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_last_successful
  }
  property {
    name  = "endpointsV1BackupsGet_max_backup_end_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_max_backup_end_time
  }
  property {
    name  = "endpointsV1BackupsGet_max_backup_start_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_max_backup_start_time
  }
  property {
    name  = "endpointsV1BackupsGet_min_backup_end_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_min_backup_end_time
  }
  property {
    name  = "endpointsV1BackupsGet_min_backup_start_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_min_backup_start_time
  }
  property {
    name  = "endpointsV1BackupsGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_page_token
  }
  property {
    name  = "endpointsV1BackupsGet_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_backups_get_user_id
  }
  property {
    name  = "endpointsV1DevicemappingsGet_device_identifier"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_get_device_identifier
  }
  property {
    name  = "endpointsV1DevicemappingsGet_device_identifier_type"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_get_device_identifier_type
  }
  property {
    name  = "endpointsV1DevicemappingsGet_email_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_get_email_id
  }
  property {
    name  = "endpointsV1DevicemappingsGet_restore_data"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_get_restore_data
  }
  property {
    name  = "endpointsV1DevicemappingsMappingIDDelete_mapping_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_mapping_iddelete_mapping_id
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_deviceIdentifier"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_device_identifier
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_deviceIdentifierType"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_device_identifier_type
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_deviceName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_device_name
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_emailID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_email_id
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_oldDeviceName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_old_device_name
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_restoreData"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_restore_data
  }
  property {
    name  = "endpointsV1DevicemappingsPost_endpointsV1DevicemappingsGetRequest_EndpointsV1DevicemappingsGetRequest_userName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devicemappings_post_endpoints_v1_devicemappings_get_request_endpoints_v1_devicemappings_get_request_user_name
  }
  property {
    name  = "endpointsV1DevicesDeviceIDDecommissionPost_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_iddecommission_post_device_id
  }
  property {
    name  = "endpointsV1DevicesDeviceIDDelete_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_iddelete_device_id
  }
  property {
    name  = "endpointsV1DevicesDeviceIDDelete_endpointsV1DevicesDeviceIDDeleteRequest_EndpointsV1DevicesDeviceIDDeleteRequest_deletionReason"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_iddelete_endpoints_v1_devices_device_iddelete_request_endpoints_v1_devices_device_iddelete_request_deletion_reason
  }
  property {
    name  = "endpointsV1DevicesDeviceIDDisablePost_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_iddisable_post_device_id
  }
  property {
    name  = "endpointsV1DevicesDeviceIDEnablePost_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_idenable_post_device_id
  }
  property {
    name  = "endpointsV1DevicesDeviceIDGet_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_idget_device_id
  }
  property {
    name  = "endpointsV1DevicesDeviceIDUpgradeClientPost_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_device_idupgrade_client_post_device_id
  }
  property {
    name  = "endpointsV1DevicesGet_client_version"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_client_version
  }
  property {
    name  = "endpointsV1DevicesGet_device_ids"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_device_ids
  }
  property {
    name  = "endpointsV1DevicesGet_device_marked_inactive"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_device_marked_inactive
  }
  property {
    name  = "endpointsV1DevicesGet_device_status"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_device_status
  }
  property {
    name  = "endpointsV1DevicesGet_max_added_on"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_max_added_on
  }
  property {
    name  = "endpointsV1DevicesGet_max_last_connected"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_max_last_connected
  }
  property {
    name  = "endpointsV1DevicesGet_min_added_on"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_min_added_on
  }
  property {
    name  = "endpointsV1DevicesGet_min_last_connected"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_min_last_connected
  }
  property {
    name  = "endpointsV1DevicesGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_page_token
  }
  property {
    name  = "endpointsV1DevicesGet_platform_os"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_platform_os
  }
  property {
    name  = "endpointsV1DevicesGet_profile_ids"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_profile_ids
  }
  property {
    name  = "endpointsV1DevicesGet_search_prefix_device_name"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_search_prefix_device_name
  }
  property {
    name  = "endpointsV1DevicesGet_serial_number"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_serial_number
  }
  property {
    name  = "endpointsV1DevicesGet_upgrade_state"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_upgrade_state
  }
  property {
    name  = "endpointsV1DevicesGet_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_user_id
  }
  property {
    name  = "endpointsV1DevicesGet_user_ids"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_user_ids
  }
  property {
    name  = "endpointsV1DevicesGet_uuid"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_devices_get_uuid
  }
  property {
    name  = "endpointsV1RestoresGet_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_device_id
  }
  property {
    name  = "endpointsV1RestoresGet_max_start_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_max_start_time
  }
  property {
    name  = "endpointsV1RestoresGet_min_start_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_min_start_time
  }
  property {
    name  = "endpointsV1RestoresGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_page_token
  }
  property {
    name  = "endpointsV1RestoresGet_target_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_target_device_id
  }
  property {
    name  = "endpointsV1RestoresGet_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_get_user_id
  }
  property {
    name  = "endpointsV1RestoresPost_endpointsV1RestoresGetRequest_EndpointsV1RestoresGetRequest_deviceID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_post_endpoints_v1_restores_get_request_endpoints_v1_restores_get_request_device_id
  }
  property {
    name  = "endpointsV1RestoresPost_endpointsV1RestoresGetRequest_EndpointsV1RestoresGetRequest_enableAVScan"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_post_endpoints_v1_restores_get_request_endpoints_v1_restores_get_request_enable_avscan
  }
  property {
    name  = "endpointsV1RestoresPost_endpointsV1RestoresGetRequest_EndpointsV1RestoresGetRequest_restoreLocation"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_post_endpoints_v1_restores_get_request_endpoints_v1_restores_get_request_restore_location
  }
  property {
    name  = "endpointsV1RestoresPost_endpointsV1RestoresGetRequest_EndpointsV1RestoresGetRequest_snapshotName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_post_endpoints_v1_restores_get_request_endpoints_v1_restores_get_request_snapshot_name
  }
  property {
    name  = "endpointsV1RestoresPost_endpointsV1RestoresGetRequest_EndpointsV1RestoresGetRequest_targetDeviceID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_post_endpoints_v1_restores_get_request_endpoints_v1_restores_get_request_target_device_id
  }
  property {
    name  = "endpointsV1RestoresRestoreIDGet_restore_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_endpoints_v1_restores_restore_idget_restore_id
  }
  property {
    name  = "eventmanagementV2EventsGet_tracker"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_eventmanagement_v2_events_get_tracker
  }
  property {
    name  = "getClient_client_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_client_client_id
  }
  property {
    name  = "getCustomerJobs_custodian"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_custodian
  }
  property {
    name  = "getCustomerJobs_job_start_time"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_job_start_time
  }
  property {
    name  = "getCustomerJobs_job_status"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_job_status
  }
  property {
    name  = "getCustomerJobs_legalhold_name"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_legalhold_name
  }
  property {
    name  = "getCustomerJobs_page_number"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_page_number
  }
  property {
    name  = "getCustomerJobs_page_size"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_page_size
  }
  property {
    name  = "getCustomerJobs_search_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_search_by
  }
  property {
    name  = "getCustomerJobs_sort_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_sort_by
  }
  property {
    name  = "getCustomerJobs_sort_order"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_customer_jobs_sort_order
  }
  property {
    name  = "getDetailsOfJob_job_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_details_of_job_job_id
  }
  property {
    name  = "getLegalHoldDetails_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_get_legal_hold_details_policy_id
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdDelete_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_delete_policy_id
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdGet_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_get_policy_id
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdUsersGet_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_users_get_policy_id
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdUsersPost_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_users_post_policy_id
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdUsersPost_requestOperateOnCustodians_RequestOperateOnCustodians_action"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_users_post_request_operate_on_custodians_request_operate_on_custodians_action
  }
  property {
    name  = "legalholdsV3PoliciesPolicyIdUsersPost_requestOperateOnCustodians_RequestOperateOnCustodians_custodians"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_policy_id_users_post_request_operate_on_custodians_request_operate_on_custodians_custodians
  }
  property {
    name  = "legalholdsV3PoliciesPost_requestCreateLHPolicy_RequestCreateLHPolicy_cullingType"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_post_request_create_lhpolicy_request_create_lhpolicy_culling_type
  }
  property {
    name  = "legalholdsV3PoliciesPost_requestCreateLHPolicy_RequestCreateLHPolicy_endDate"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_post_request_create_lhpolicy_request_create_lhpolicy_end_date
  }
  property {
    name  = "legalholdsV3PoliciesPost_requestCreateLHPolicy_RequestCreateLHPolicy_name"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_post_request_create_lhpolicy_request_create_lhpolicy_name
  }
  property {
    name  = "legalholdsV3PoliciesPost_requestCreateLHPolicy_RequestCreateLHPolicy_startDate"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_post_request_create_lhpolicy_request_create_lhpolicy_start_date
  }
  property {
    name  = "legalholdsV3PoliciesPost_requestCreateLHPolicy_RequestCreateLHPolicy_type"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_policies_post_request_create_lhpolicy_request_create_lhpolicy_type
  }
  property {
    name  = "legalholdsV3UsersPoliciesPost_requestListPoliciesByUsers_RequestListPoliciesByUsers_action"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_users_policies_post_request_list_policies_by_users_request_list_policies_by_users_action
  }
  property {
    name  = "legalholdsV3UsersPoliciesPost_requestListPoliciesByUsers_RequestListPoliciesByUsers_custodians"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v3_users_policies_post_request_list_policies_by_users_request_list_policies_by_users_custodians
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdDelete_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_delete_policy_id
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_page_number"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_page_number
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_page_size"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_page_size
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_policy_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_policy_id
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_search_prefix"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_search_prefix
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_sort_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_sort_by
  }
  property {
    name  = "legalholdsV4PoliciesPolicyIdUsersGet_sort_order"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_legalholds_v4_policies_policy_id_users_get_sort_order
  }
  property {
    name  = "policiesGet_custodians"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_get_custodians
  }
  property {
    name  = "policiesGet_legal_hold_types"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_get_legal_hold_types
  }
  property {
    name  = "policiesGet_search_string"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_get_search_string
  }
  property {
    name  = "policiesGet_sort_by"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_get_sort_by
  }
  property {
    name  = "policiesGet_sort_order"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_get_sort_order
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_accessAdditionalData"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_access_additional_data
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_accessBackupData"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_access_backup_data
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_collectionFrequency"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_collection_frequency
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_deletedFilesOnly"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_deleted_files_only
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_endDate"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_end_date
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_isADCEnabled"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_is_adcenabled
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_legalHoldType"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_legal_hold_type
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_name"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_name
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_startDate"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_start_date
  }
  property {
    name  = "policiesPost_createLegalHoldPolicyRequest_CreateLegalHoldPolicyRequest_webDavAccessType"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_policies_post_create_legal_hold_policy_request_create_legal_hold_policy_request_web_dav_access_type
  }
  property {
    name  = "profilemanagementV1ProfilesGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_profilemanagement_v1_profiles_get_page_token
  }
  property {
    name  = "profilemanagementV1ProfilesProfileIDGet_profile_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_profilemanagement_v1_profiles_profile_idget_profile_id
  }
  property {
    name  = "referenceGettingStarted_grant_type"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_reference_getting_started_grant_type
  }
  property {
    name  = "referenceGettingStarted_scope"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_reference_getting_started_scope
  }
  property {
    name  = "storagemanagementV1CloudcachesCloudCacheIDGet_cloud_cache_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_storagemanagement_v1_cloudcaches_cloud_cache_idget_cloud_cache_id
  }
  property {
    name  = "storagemanagementV1CloudcachesGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_storagemanagement_v1_cloudcaches_get_page_token
  }
  property {
    name  = "storagemanagementV1StoragesGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_storagemanagement_v1_storages_get_page_token
  }
  property {
    name  = "storagemanagementV1StoragesStorageIDGet_storage_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_storagemanagement_v1_storages_storage_idget_storage_id
  }
  property {
    name  = "takeActionOnJob_takeActionOnJobRequest_TakeActionOnJobRequest_action"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_take_action_on_job_take_action_on_job_request_take_action_on_job_request_action
  }
  property {
    name  = "takeActionOnJob_takeActionOnJobRequest_TakeActionOnJobRequest_jobIds"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_take_action_on_job_take_action_on_job_request_take_action_on_job_request_job_ids
  }
  property {
    name  = "usermanagementV1UsersGet_cache_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_cache_id
  }
  property {
    name  = "usermanagementV1UsersGet_email_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_email_id
  }
  property {
    name  = "usermanagementV1UsersGet_ldap_guid"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_ldap_guid
  }
  property {
    name  = "usermanagementV1UsersGet_max_added_on"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_max_added_on
  }
  property {
    name  = "usermanagementV1UsersGet_min_added_on"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_min_added_on
  }
  property {
    name  = "usermanagementV1UsersGet_page_token"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_page_token
  }
  property {
    name  = "usermanagementV1UsersGet_privacy_settings_enabled"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_privacy_settings_enabled
  }
  property {
    name  = "usermanagementV1UsersGet_profile_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_profile_id
  }
  property {
    name  = "usermanagementV1UsersGet_search_prefix_email_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_search_prefix_email_id
  }
  property {
    name  = "usermanagementV1UsersGet_search_prefix_user_name"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_search_prefix_user_name
  }
  property {
    name  = "usermanagementV1UsersGet_status"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_status
  }
  property {
    name  = "usermanagementV1UsersGet_storage_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_storage_id
  }
  property {
    name  = "usermanagementV1UsersGet_user_ids"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_get_user_ids
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_emailID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_email_id
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_profileID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_profile_id
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_quota"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_quota
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_quotaUnit"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_quota_unit
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_storageID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_storage_id
  }
  property {
    name  = "usermanagementV1UsersPost_usermanagementV1UsersGetRequest_UsermanagementV1UsersGetRequest_userName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_post_usermanagement_v1_users_get_request_usermanagement_v1_users_get_request_user_name
  }
  property {
    name  = "usermanagementV1UsersUserIDActivatePost_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idactivate_post_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDDelete_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_iddelete_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDDelete_usermanagementV1UsersUserIDDeleteRequest_UsermanagementV1UsersUserIDDeleteRequest_deletionReason"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_iddelete_usermanagement_v1_users_user_iddelete_request_usermanagement_v1_users_user_iddelete_request_deletion_reason
  }
  property {
    name  = "usermanagementV1UsersUserIDDownloadUserAuthKeyGet_device_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_iddownload_user_auth_key_get_device_id
  }
  property {
    name  = "usermanagementV1UsersUserIDDownloadUserAuthKeyGet_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_iddownload_user_auth_key_get_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDGet_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idget_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDPatch_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idpatch_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDPatch_usermanagementV1UsersUserIDDeleteRequest1_UsermanagementV1UsersUserIDDeleteRequest1_profileID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idpatch_usermanagement_v1_users_user_iddelete_request1_usermanagement_v1_users_user_iddelete_request1_profile_id
  }
  property {
    name  = "usermanagementV1UsersUserIDPatch_usermanagementV1UsersUserIDDeleteRequest1_UsermanagementV1UsersUserIDDeleteRequest1_storageID"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idpatch_usermanagement_v1_users_user_iddelete_request1_usermanagement_v1_users_user_iddelete_request1_storage_id
  }
  property {
    name  = "usermanagementV1UsersUserIDPatch_usermanagementV1UsersUserIDDeleteRequest1_UsermanagementV1UsersUserIDDeleteRequest1_userName"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idpatch_usermanagement_v1_users_user_iddelete_request1_usermanagement_v1_users_user_iddelete_request1_user_name
  }
  property {
    name  = "usermanagementV1UsersUserIDPreservePost_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idpreserve_post_user_id
  }
  property {
    name  = "usermanagementV1UsersUserIDResetPasswordPost_user_id"
    type  = "string"
    value = var.connector-oai-druvainsynccloud_property_usermanagement_v1_users_user_idreset_password_post_user_id
  }
}
