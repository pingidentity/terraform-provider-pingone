resource "pingone_davinci_connector_instance" "connector-oai-mailjetapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-mailjetapi"
  }
  name = "My awesome connector-oai-mailjetapi"
  property {
    name  = "authPassword"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_auth_password
  }
  property {
    name  = "authUsername"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_auth_username
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_base_path
  }
  property {
    name  = "v31SendPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v31_send_post_body
  }
  property {
    name  = "v31SendPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v31_send_post_content_type
  }
  property {
    name  = "v3RESTApikeyApikeyIDGet_apikey_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_apikey_idget_apikey_id
  }
  property {
    name  = "v3RESTApikeyApikeyIDPut_apikey_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_apikey_idput_apikey_id
  }
  property {
    name  = "v3RESTApikeyApikeyIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_apikey_idput_body
  }
  property {
    name  = "v3RESTApikeyApikeyIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_apikey_idput_content_type
  }
  property {
    name  = "v3RESTApikeyPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_post_body
  }
  property {
    name  = "v3RESTApikeyPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restapikey_post_content_type
  }
  property {
    name  = "v3RESTBouncestatisticsMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restbouncestatistics_message_idget_message_id
  }
  property {
    name  = "v3RESTCampaignCampaignIDGet_campaign_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaign_campaign_idget_campaign_id
  }
  property {
    name  = "v3RESTCampaignCampaignIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaign_campaign_idput_body
  }
  property {
    name  = "v3RESTCampaignCampaignIDPut_campaign_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaign_campaign_idput_campaign_id
  }
  property {
    name  = "v3RESTCampaignCampaignIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaign_campaign_idput_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDDetailcontentGet_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_iddetailcontent_get_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDDetailcontentPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_iddetailcontent_post_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDDetailcontentPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_iddetailcontent_post_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDDetailcontentPost_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_iddetailcontent_post_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDGet_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idget_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idput_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idput_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDPut_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idput_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDScheduleDelete_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_delete_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDScheduleGet_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_get_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_post_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_post_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePost_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_post_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_put_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_put_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSchedulePut_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idschedule_put_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSendPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idsend_post_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDSendPost_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idsend_post_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDStatusGet_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idstatus_get_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDTestPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idtest_post_body
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDTestPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idtest_post_content_type
  }
  property {
    name  = "v3RESTCampaigndraftDraftIDTestPost_draft_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_draft_idtest_post_draft_id
  }
  property {
    name  = "v3RESTCampaigndraftPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_post_body
  }
  property {
    name  = "v3RESTCampaigndraftPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigndraft_post_content_type
  }
  property {
    name  = "v3RESTCampaigngraphstatisticsGet_i_ds"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaigngraphstatistics_get_i_ds
  }
  property {
    name  = "v3RESTCampaignoverviewIDTypeIDGet_i_d"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaignoverview_idtype_idget_i_d
  }
  property {
    name  = "v3RESTCampaignoverviewIDTypeIDGet_i_d_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcampaignoverview_idtype_idget_i_d_type
  }
  property {
    name  = "v3RESTContactContactIDGet_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idget_contact_id
  }
  property {
    name  = "v3RESTContactContactIDGetcontactslistsGet_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idgetcontactslists_get_contact_id
  }
  property {
    name  = "v3RESTContactContactIDManagecontactslistsPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idmanagecontactslists_post_body
  }
  property {
    name  = "v3RESTContactContactIDManagecontactslistsPost_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idmanagecontactslists_post_contact_id
  }
  property {
    name  = "v3RESTContactContactIDManagecontactslistsPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idmanagecontactslists_post_content_type
  }
  property {
    name  = "v3RESTContactContactIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idput_body
  }
  property {
    name  = "v3RESTContactContactIDPut_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idput_contact_id
  }
  property {
    name  = "v3RESTContactContactIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_contact_idput_content_type
  }
  property {
    name  = "v3RESTContactManagemanycontactsJobIDGet_job_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_managemanycontacts_job_idget_job_id
  }
  property {
    name  = "v3RESTContactManagemanycontactsPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_managemanycontacts_post_body
  }
  property {
    name  = "v3RESTContactManagemanycontactsPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_managemanycontacts_post_content_type
  }
  property {
    name  = "v3RESTContactPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_post_body
  }
  property {
    name  = "v3RESTContactPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontact_post_content_type
  }
  property {
    name  = "v3RESTContactdataContactIDDelete_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_contact_iddelete_contact_id
  }
  property {
    name  = "v3RESTContactdataContactIDGet_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_contact_idget_contact_id
  }
  property {
    name  = "v3RESTContactdataContactIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_contact_idput_body
  }
  property {
    name  = "v3RESTContactdataContactIDPut_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_contact_idput_contact_id
  }
  property {
    name  = "v3RESTContactdataContactIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_contact_idput_content_type
  }
  property {
    name  = "v3RESTContactdataPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_post_body
  }
  property {
    name  = "v3RESTContactdataPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactdata_post_content_type
  }
  property {
    name  = "v3RESTContactfilterContactfilterIDDelete_contactfilter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_contactfilter_iddelete_contactfilter_id
  }
  property {
    name  = "v3RESTContactfilterContactfilterIDGet_contactfilter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_contactfilter_idget_contactfilter_id
  }
  property {
    name  = "v3RESTContactfilterContactfilterIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_contactfilter_idput_body
  }
  property {
    name  = "v3RESTContactfilterContactfilterIDPut_contactfilter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_contactfilter_idput_contactfilter_id
  }
  property {
    name  = "v3RESTContactfilterContactfilterIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_contactfilter_idput_content_type
  }
  property {
    name  = "v3RESTContactfilterPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_post_body
  }
  property {
    name  = "v3RESTContactfilterPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactfilter_post_content_type
  }
  property {
    name  = "v3RESTContactmetadataContactmetadataIDDelete_contactmetadata_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_contactmetadata_iddelete_contactmetadata_id
  }
  property {
    name  = "v3RESTContactmetadataContactmetadataIDGet_contactmetadata_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_contactmetadata_idget_contactmetadata_id
  }
  property {
    name  = "v3RESTContactmetadataContactmetadataIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_contactmetadata_idput_body
  }
  property {
    name  = "v3RESTContactmetadataContactmetadataIDPut_contactmetadata_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_contactmetadata_idput_contactmetadata_id
  }
  property {
    name  = "v3RESTContactmetadataContactmetadataIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_contactmetadata_idput_content_type
  }
  property {
    name  = "v3RESTContactmetadataPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_post_body
  }
  property {
    name  = "v3RESTContactmetadataPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactmetadata_post_content_type
  }
  property {
    name  = "v3RESTContactslistListIDDelete_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_iddelete_list_id
  }
  property {
    name  = "v3RESTContactslistListIDGet_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idget_list_id
  }
  property {
    name  = "v3RESTContactslistListIDImportlistJobIDGet_job_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idimportlist_job_idget_job_id
  }
  property {
    name  = "v3RESTContactslistListIDImportlistJobIDGet_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idimportlist_job_idget_list_id
  }
  property {
    name  = "v3RESTContactslistListIDImportlistPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idimportlist_post_body
  }
  property {
    name  = "v3RESTContactslistListIDImportlistPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idimportlist_post_content_type
  }
  property {
    name  = "v3RESTContactslistListIDImportlistPost_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idimportlist_post_list_id
  }
  property {
    name  = "v3RESTContactslistListIDManagecontactPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagecontact_post_body
  }
  property {
    name  = "v3RESTContactslistListIDManagecontactPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagecontact_post_content_type
  }
  property {
    name  = "v3RESTContactslistListIDManagecontactPost_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagecontact_post_list_id
  }
  property {
    name  = "v3RESTContactslistListIDManagemanycontactsJobIDGet_job_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagemanycontacts_job_idget_job_id
  }
  property {
    name  = "v3RESTContactslistListIDManagemanycontactsJobIDGet_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagemanycontacts_job_idget_list_id
  }
  property {
    name  = "v3RESTContactslistListIDManagemanycontactsPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagemanycontacts_post_body
  }
  property {
    name  = "v3RESTContactslistListIDManagemanycontactsPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagemanycontacts_post_content_type
  }
  property {
    name  = "v3RESTContactslistListIDManagemanycontactsPost_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idmanagemanycontacts_post_list_id
  }
  property {
    name  = "v3RESTContactslistListIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idput_body
  }
  property {
    name  = "v3RESTContactslistListIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idput_content_type
  }
  property {
    name  = "v3RESTContactslistListIDPut_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_list_idput_list_id
  }
  property {
    name  = "v3RESTContactslistPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_post_body
  }
  property {
    name  = "v3RESTContactslistPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslist_post_content_type
  }
  property {
    name  = "v3RESTContactslistsignupSignuprequestIDGet_signuprequest_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactslistsignup_signuprequest_idget_signuprequest_id
  }
  property {
    name  = "v3RESTContactstatisticsContactIDGet_contact_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcontactstatistics_contact_idget_contact_id
  }
  property {
    name  = "v3RESTCsvimportImportjobIDGet_importjob_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_importjob_idget_importjob_id
  }
  property {
    name  = "v3RESTCsvimportImportjobIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_importjob_idput_body
  }
  property {
    name  = "v3RESTCsvimportImportjobIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_importjob_idput_content_type
  }
  property {
    name  = "v3RESTCsvimportImportjobIDPut_importjob_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_importjob_idput_importjob_id
  }
  property {
    name  = "v3RESTCsvimportPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_post_body
  }
  property {
    name  = "v3RESTCsvimportPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restcsvimport_post_content_type
  }
  property {
    name  = "v3RESTDnsDnsIDCheckPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restdns_dns_idcheck_post_body
  }
  property {
    name  = "v3RESTDnsDnsIDCheckPost_dns_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restdns_dns_idcheck_post_dns_id
  }
  property {
    name  = "v3RESTDnsDnsIDGet_dns_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restdns_dns_idget_dns_id
  }
  property {
    name  = "v3RESTEventcallbackurlPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_post_body
  }
  property {
    name  = "v3RESTEventcallbackurlPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_post_content_type
  }
  property {
    name  = "v3RESTEventcallbackurlUrlIDDelete_url_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_url_iddelete_url_id
  }
  property {
    name  = "v3RESTEventcallbackurlUrlIDGet_url_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_url_idget_url_id
  }
  property {
    name  = "v3RESTEventcallbackurlUrlIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_url_idput_body
  }
  property {
    name  = "v3RESTEventcallbackurlUrlIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_url_idput_content_type
  }
  property {
    name  = "v3RESTEventcallbackurlUrlIDPut_url_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resteventcallbackurl_url_idput_url_id
  }
  property {
    name  = "v3RESTGraphstatisticsGet_scale"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restgraphstatistics_get_scale
  }
  property {
    name  = "v3RESTListrecipientListrecipientIDDelete_listrecipient_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_listrecipient_iddelete_listrecipient_id
  }
  property {
    name  = "v3RESTListrecipientListrecipientIDGet_listrecipient_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_listrecipient_idget_listrecipient_id
  }
  property {
    name  = "v3RESTListrecipientListrecipientIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_listrecipient_idput_body
  }
  property {
    name  = "v3RESTListrecipientListrecipientIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_listrecipient_idput_content_type
  }
  property {
    name  = "v3RESTListrecipientListrecipientIDPut_listrecipient_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_listrecipient_idput_listrecipient_id
  }
  property {
    name  = "v3RESTListrecipientPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_post_body
  }
  property {
    name  = "v3RESTListrecipientPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipient_post_content_type
  }
  property {
    name  = "v3RESTListrecipientstatisticsListrecipientIDGet_listrecipient_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restlistrecipientstatistics_listrecipient_idget_listrecipient_id
  }
  property {
    name  = "v3RESTListstatisticsListIDGet_list_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restliststatistics_list_idget_list_id
  }
  property {
    name  = "v3RESTMessageMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmessage_message_idget_message_id
  }
  property {
    name  = "v3RESTMessagehistoryMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmessagehistory_message_idget_message_id
  }
  property {
    name  = "v3RESTMessageinformationMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmessageinformation_message_idget_message_id
  }
  property {
    name  = "v3RESTMessagesentstatisticsMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmessagesentstatistics_message_idget_message_id
  }
  property {
    name  = "v3RESTMetasenderMetasenderIDGet_metasender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_metasender_idget_metasender_id
  }
  property {
    name  = "v3RESTMetasenderMetasenderIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_metasender_idput_body
  }
  property {
    name  = "v3RESTMetasenderMetasenderIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_metasender_idput_content_type
  }
  property {
    name  = "v3RESTMetasenderMetasenderIDPut_metasender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_metasender_idput_metasender_id
  }
  property {
    name  = "v3RESTMetasenderPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_post_body
  }
  property {
    name  = "v3RESTMetasenderPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmetasender_post_content_type
  }
  property {
    name  = "v3RESTMyprofilePut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmyprofile_put_body
  }
  property {
    name  = "v3RESTMyprofilePut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restmyprofile_put_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentDelete_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_delete_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentGet_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_get_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_post_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_post_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPost_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_post_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_put_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_put_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDDetailcontentPut_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_iddetailcontent_put_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDGet_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idget_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idput_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idput_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDPut_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idput_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDScheduleDelete_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_delete_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDScheduleGet_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_get_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_post_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_post_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePost_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_post_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_put_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_put_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSchedulePut_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idschedule_put_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSendPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idsend_post_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDSendPost_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idsend_post_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDStatusGet_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idstatus_get_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDTestPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idtest_post_body
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDTestPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idtest_post_content_type
  }
  property {
    name  = "v3RESTNewsletterNewsletterIDTestPost_newsletter_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_newsletter_idtest_post_newsletter_id
  }
  property {
    name  = "v3RESTNewsletterPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_post_body
  }
  property {
    name  = "v3RESTNewsletterPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restnewsletter_post_content_type
  }
  property {
    name  = "v3RESTOpeninformationMessageIDGet_message_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restopeninformation_message_idget_message_id
  }
  property {
    name  = "v3RESTParserouteParserouteIDDelete_parseroute_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_parseroute_iddelete_parseroute_id
  }
  property {
    name  = "v3RESTParserouteParserouteIDGet_parseroute_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_parseroute_idget_parseroute_id
  }
  property {
    name  = "v3RESTParserouteParserouteIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_parseroute_idput_body
  }
  property {
    name  = "v3RESTParserouteParserouteIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_parseroute_idput_content_type
  }
  property {
    name  = "v3RESTParserouteParserouteIDPut_parseroute_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_parseroute_idput_parseroute_id
  }
  property {
    name  = "v3RESTParseroutePost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_post_body
  }
  property {
    name  = "v3RESTParseroutePost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restparseroute_post_content_type
  }
  property {
    name  = "v3RESTSenderPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_post_body
  }
  property {
    name  = "v3RESTSenderPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_post_content_type
  }
  property {
    name  = "v3RESTSenderSenderIDDelete_sender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_iddelete_sender_id
  }
  property {
    name  = "v3RESTSenderSenderIDGet_sender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idget_sender_id
  }
  property {
    name  = "v3RESTSenderSenderIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idput_body
  }
  property {
    name  = "v3RESTSenderSenderIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idput_content_type
  }
  property {
    name  = "v3RESTSenderSenderIDPut_sender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idput_sender_id
  }
  property {
    name  = "v3RESTSenderSenderIDValidatePost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idvalidate_post_body
  }
  property {
    name  = "v3RESTSenderSenderIDValidatePost_sender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsender_sender_idvalidate_post_sender_id
  }
  property {
    name  = "v3RESTSenderstatisticsSenderIDGet_sender_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restsenderstatistics_sender_idget_sender_id
  }
  property {
    name  = "v3RESTStatcountersGet_counter_resolution"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_counter_resolution
  }
  property {
    name  = "v3RESTStatcountersGet_counter_source"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_counter_source
  }
  property {
    name  = "v3RESTStatcountersGet_counter_timing"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_counter_timing
  }
  property {
    name  = "v3RESTStatcountersGet_from_ts"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_from_ts
  }
  property {
    name  = "v3RESTStatcountersGet_source_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_source_id
  }
  property {
    name  = "v3RESTStatcountersGet_to_ts"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatcounters_get_to_ts
  }
  property {
    name  = "v3RESTStatisticsLinkClickGet_campaign_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatistics_link_click_get_campaign_id
  }
  property {
    name  = "v3RESTStatisticsRecipientEspGet_campaign_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_reststatistics_recipient_esp_get_campaign_id
  }
  property {
    name  = "v3RESTTemplatePost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_post_body
  }
  property {
    name  = "v3RESTTemplatePost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_post_content_type
  }
  property {
    name  = "v3RESTTemplateTemplateIDDelete_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddelete_template_id
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentGet_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_get_template_id
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_post_body
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_post_content_type
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPost_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_post_template_id
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_put_body
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_put_content_type
  }
  property {
    name  = "v3RESTTemplateTemplateIDDetailcontentPut_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_iddetailcontent_put_template_id
  }
  property {
    name  = "v3RESTTemplateTemplateIDGet_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_idget_template_id
  }
  property {
    name  = "v3RESTTemplateTemplateIDPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_idput_body
  }
  property {
    name  = "v3RESTTemplateTemplateIDPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_idput_content_type
  }
  property {
    name  = "v3RESTTemplateTemplateIDPut_template_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_resttemplate_template_idput_template_id
  }
  property {
    name  = "v3RESTUserPut_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restuser_put_body
  }
  property {
    name  = "v3RESTUserPut_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_restuser_put_content_type
  }
  property {
    name  = "v3SendPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_send_post_body
  }
  property {
    name  = "v3SendPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v3_send_post_content_type
  }
  property {
    name  = "v4SmsExportJobIDGet_job_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_export_job_idget_job_id
  }
  property {
    name  = "v4SmsExportPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_export_post_body
  }
  property {
    name  = "v4SmsExportPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_export_post_content_type
  }
  property {
    name  = "v4SmsSendPost_body"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_send_post_body
  }
  property {
    name  = "v4SmsSendPost_content_type"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_send_post_content_type
  }
  property {
    name  = "v4SmsSmsIDGet_sms_id"
    type  = "string"
    value = var.connector-oai-mailjetapi_property_v4_sms_sms_idget_sms_id
  }
}
