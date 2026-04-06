resource "pingone_davinci_connector_instance" "connector-oai-activecampaignapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-activecampaignapi"
  }
  name = "My awesome connector-oai-activecampaignapi"
  property {
    name  = "apiApiVersionAccountsIdDelete_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_delete_id
  }
  property {
    name  = "apiApiVersionAccountsIdGet_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_get_id
  }
  property {
    name  = "apiApiVersionAccountsIdNotesNoteIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_notes_note_id_put_body
  }
  property {
    name  = "apiApiVersionAccountsIdNotesNoteIdPut_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_notes_note_id_put_id
  }
  property {
    name  = "apiApiVersionAccountsIdNotesNoteIdPut_note_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_notes_note_id_put_note_id
  }
  property {
    name  = "apiApiVersionAccountsIdNotesPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_notes_post_body
  }
  property {
    name  = "apiApiVersionAccountsIdNotesPost_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_notes_post_id
  }
  property {
    name  = "apiApiVersionAccountsIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_put_body
  }
  property {
    name  = "apiApiVersionAccountsIdPut_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_id_put_id
  }
  property {
    name  = "apiApiVersionAccountsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_accounts_post_body
  }
  property {
    name  = "apiApiVersionActivitiesGet_contact"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_activities_get_contact
  }
  property {
    name  = "apiApiVersionCampaignsCampaignIDGet_campaign_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_campaigns_campaign_idget_campaign_id
  }
  property {
    name  = "apiApiVersionCampaignsCampaignIDLinksGet_campaign_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_campaigns_campaign_idlinks_get_campaign_id
  }
  property {
    name  = "apiApiVersionContactAutomationsContactAutomationIDDelete_contact_automation_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_automations_contact_automation_iddelete_contact_automation_id
  }
  property {
    name  = "apiApiVersionContactAutomationsContactAutomationIDGet_contact_automation_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_automations_contact_automation_idget_contact_automation_id
  }
  property {
    name  = "apiApiVersionContactAutomationsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_automations_post_body
  }
  property {
    name  = "apiApiVersionContactListsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_lists_post_body
  }
  property {
    name  = "apiApiVersionContactSyncPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_sync_post_body
  }
  property {
    name  = "apiApiVersionContactTagsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contact_tags_post_body
  }
  property {
    name  = "apiApiVersionContactsContactIdAccountContactsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_account_contacts_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdAutomationEntryCountsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_automation_entry_counts_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdBounceLogsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_bounce_logs_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactAutomationsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_automations_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactDataGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_data_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactDealsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_deals_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactGoalsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_goals_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactListsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_lists_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactLogsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_logs_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdContactTagsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_contact_tags_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdDealsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_deals_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdDelete_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_delete_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdFieldValuesGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_field_values_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdGeoIpsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_geo_ips_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdNotesGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_notes_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdOrganizationGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_organization_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdPlusAppendGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_plus_append_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsContactIdTrackingLogsGet_contact_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_contact_id_tracking_logs_get_contact_id
  }
  property {
    name  = "apiApiVersionContactsGet_filters_email"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_get_filters_email
  }
  property {
    name  = "apiApiVersionContactsGet_include"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_get_include
  }
  property {
    name  = "apiApiVersionContactsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_contacts_post_body
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIDGet_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_idget_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdExternalExternalIdDelete_external_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_external_external_id_delete_external_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdExternalExternalIdDelete_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_external_external_id_delete_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdExternalExternalIdGet_external_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_external_external_id_get_external_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdExternalExternalIdGet_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_external_external_id_get_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdGet_filters_relationships_primary_contact_eq"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_get_filters_relationships_primary_contact_eq
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdGet_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_get_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_post_body
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdPost_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_post_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdRecordIdDelete_record_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_record_id_delete_record_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdRecordIdDelete_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_record_id_delete_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdRecordIdGet_record_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_record_id_get_record_id
  }
  property {
    name  = "apiApiVersionCustomObjectsRecordsSchemaIdRecordIdGet_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_records_schema_id_record_id_get_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasParentSchemaIdChildPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_parent_schema_id_child_post_body
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasParentSchemaIdChildPost_parent_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_parent_schema_id_child_post_parent_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_post_body
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasPublicPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_public_post_body
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasSchemaIdDelete_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_schema_id_delete_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasSchemaIdGet_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_schema_id_get_schema_id
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasSchemaIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_schema_id_put_body
  }
  property {
    name  = "apiApiVersionCustomObjectsSchemasSchemaIdPut_schema_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_custom_objects_schemas_schema_id_put_schema_id
  }
  property {
    name  = "apiApiVersionDealRolesPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deal_roles_post_body
  }
  property {
    name  = "apiApiVersionDealsBulkUpdatePatch_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_bulk_update_patch_body
  }
  property {
    name  = "apiApiVersionDealsDealIdNotesNoteIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_deal_id_notes_note_id_put_body
  }
  property {
    name  = "apiApiVersionDealsDealIdNotesNoteIdPut_deal_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_deal_id_notes_note_id_put_deal_id
  }
  property {
    name  = "apiApiVersionDealsDealIdNotesNoteIdPut_note_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_deal_id_notes_note_id_put_note_id
  }
  property {
    name  = "apiApiVersionDealsIdDealActivitiesGet_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_deal_activities_get_id
  }
  property {
    name  = "apiApiVersionDealsIdDelete_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_delete_id
  }
  property {
    name  = "apiApiVersionDealsIdGet_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_get_id
  }
  property {
    name  = "apiApiVersionDealsIdNotesPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_notes_post_body
  }
  property {
    name  = "apiApiVersionDealsIdNotesPost_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_notes_post_id
  }
  property {
    name  = "apiApiVersionDealsIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_put_body
  }
  property {
    name  = "apiApiVersionDealsIdPut_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_id_put_id
  }
  property {
    name  = "apiApiVersionDealsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_deals_post_body
  }
  property {
    name  = "apiApiVersionEventTrackingEventEventNameDelete_event_name"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_event_tracking_event_event_name_delete_event_name
  }
  property {
    name  = "apiApiVersionEventTrackingEventsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_event_tracking_events_post_body
  }
  property {
    name  = "apiApiVersionEventTrackingPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_event_tracking_put_body
  }
  property {
    name  = "apiApiVersionFieldOptionBulkPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_field_option_bulk_post_body
  }
  property {
    name  = "apiApiVersionFieldRelsFieldRelIdDelete_field_rel_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_field_rels_field_rel_id_delete_field_rel_id
  }
  property {
    name  = "apiApiVersionFieldRelsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_field_rels_post_body
  }
  property {
    name  = "apiApiVersionFieldValuesPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_field_values_post_body
  }
  property {
    name  = "apiApiVersionFieldsFieldIDDelete_field_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_fields_field_iddelete_field_id
  }
  property {
    name  = "apiApiVersionFieldsFieldIDPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_fields_field_idput_body
  }
  property {
    name  = "apiApiVersionFieldsFieldIDPut_field_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_fields_field_idput_field_id
  }
  property {
    name  = "apiApiVersionFieldsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_fields_post_body
  }
  property {
    name  = "apiApiVersionGroupsGroupIdGet_group_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_groups_group_id_get_group_id
  }
  property {
    name  = "apiApiVersionGroupsGroupIdUserGroupsGet_group_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_groups_group_id_user_groups_get_group_id
  }
  property {
    name  = "apiApiVersionGroupsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_groups_post_body
  }
  property {
    name  = "apiApiVersionImportBulkImportPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_import_bulk_import_post_body
  }
  property {
    name  = "apiApiVersionImportInfoGet_batch_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_import_info_get_batch_id
  }
  property {
    name  = "apiApiVersionListsGet_limit"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_lists_get_limit
  }
  property {
    name  = "apiApiVersionListsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_lists_post_body
  }
  property {
    name  = "apiApiVersionMessagesIdDelete_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_messages_id_delete_id
  }
  property {
    name  = "apiApiVersionMessagesIdGet_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_messages_id_get_id
  }
  property {
    name  = "apiApiVersionMessagesIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_messages_id_put_body
  }
  property {
    name  = "apiApiVersionMessagesIdPut_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_messages_id_put_id
  }
  property {
    name  = "apiApiVersionMessagesPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_messages_post_body
  }
  property {
    name  = "apiApiVersionSiteTrackingDomainsNameDelete_name"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_site_tracking_domains_name_delete_name
  }
  property {
    name  = "apiApiVersionSiteTrackingDomainsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_site_tracking_domains_post_body
  }
  property {
    name  = "apiApiVersionSiteTrackingPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_site_tracking_put_body
  }
  property {
    name  = "apiApiVersionTagsPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_tags_post_body
  }
  property {
    name  = "apiApiVersionTagsTagIdDelete_tag_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_tags_tag_id_delete_tag_id
  }
  property {
    name  = "apiApiVersionTagsTagIdPut_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_tags_tag_id_put_body
  }
  property {
    name  = "apiApiVersionTagsTagIdPut_tag_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_tags_tag_id_put_tag_id
  }
  property {
    name  = "apiApiVersionTrackingLogsEventIdContactGet_event_id"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_tracking_logs_event_id_contact_get_event_id
  }
  property {
    name  = "apiApiVersionUsersPost_body"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_api_api_version_users_post_body
  }
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_auth_api_key
  }
  property {
    name  = "authApiVersion"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_auth_api_version
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_base_path
  }
  property {
    name  = "eventPost_actid"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_event_post_actid
  }
  property {
    name  = "eventPost_event"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_event_post_event
  }
  property {
    name  = "eventPost_eventdata"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_event_post_eventdata
  }
  property {
    name  = "eventPost_key"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_event_post_key
  }
  property {
    name  = "eventPost_visit"
    type  = "string"
    value = var.connector-oai-activecampaignapi_property_event_post_visit
  }
}
