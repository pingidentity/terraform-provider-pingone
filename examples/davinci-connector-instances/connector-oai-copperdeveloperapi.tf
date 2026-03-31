resource "pingone_davinci_connector_instance" "connector-oai-copperdeveloperapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-copperdeveloperapi"
  }
  name = "My awesome connector-oai-copperdeveloperapi"
  property {
    name  = "activitiesDeleteActivityIdDelete_delete_activity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_activities_delete_activity_id_delete_delete_activity_id
  }
  property {
    name  = "activitiesExampleActivityIdGet_example_activity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_activities_example_activity_id_get_example_activity_id
  }
  property {
    name  = "activitiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_activities_post_body
  }
  property {
    name  = "activitiesSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_activities_search_post_body
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_base_path
  }
  property {
    name  = "companiesDeleteCompanyIdDelete_delete_company_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_delete_company_id_delete_delete_company_id
  }
  property {
    name  = "companiesExampleCompanyIdActivitiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_example_company_id_activities_post_body
  }
  property {
    name  = "companiesExampleCompanyIdActivitiesPost_example_company_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_example_company_id_activities_post_example_company_id
  }
  property {
    name  = "companiesExampleCompanyIdGet_example_company_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_example_company_id_get_example_company_id
  }
  property {
    name  = "companiesExampleCompanyIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_example_company_id_put_body
  }
  property {
    name  = "companiesExampleCompanyIdPut_example_company_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_example_company_id_put_example_company_id
  }
  property {
    name  = "companiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_post_body
  }
  property {
    name  = "companiesSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_companies_search_post_body
  }
  property {
    name  = "contentType"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_content_type
  }
  property {
    name  = "customActivityTypesCustomActivityTypeIdGet_custom_activity_type_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_activity_types_custom_activity_type_id_get_custom_activity_type_id
  }
  property {
    name  = "customActivityTypesCustomActivityTypeIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_activity_types_custom_activity_type_id_put_body
  }
  property {
    name  = "customActivityTypesCustomActivityTypeIdPut_custom_activity_type_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_activity_types_custom_activity_type_id_put_custom_activity_type_id
  }
  property {
    name  = "customActivityTypesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_activity_types_post_body
  }
  property {
    name  = "customFieldDefinitionsCustomFieldDefinitionIdDelete_custom_field_definition_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_field_definitions_custom_field_definition_id_delete_custom_field_definition_id
  }
  property {
    name  = "customFieldDefinitionsCustomFieldDefinitionIdGet_custom_field_definition_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_field_definitions_custom_field_definition_id_get_custom_field_definition_id
  }
  property {
    name  = "customFieldDefinitionsCustomFieldDefinitionIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_field_definitions_custom_field_definition_id_put_body
  }
  property {
    name  = "customFieldDefinitionsCustomFieldDefinitionIdPut_custom_field_definition_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_field_definitions_custom_field_definition_id_put_custom_field_definition_id
  }
  property {
    name  = "customFieldDefinitionsPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_custom_field_definitions_post_body
  }
  property {
    name  = "entityEntityIdFilesFileIdGet_entity"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_files_file_id_get_entity
  }
  property {
    name  = "entityEntityIdFilesFileIdGet_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_files_file_id_get_entity_id
  }
  property {
    name  = "entityEntityIdFilesFileIdGet_file_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_files_file_id_get_file_id
  }
  property {
    name  = "entityEntityIdRelatedDelete_entity"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_delete_entity
  }
  property {
    name  = "entityEntityIdRelatedDelete_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_delete_entity_id
  }
  property {
    name  = "entityEntityIdRelatedGet_entity"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_get_entity
  }
  property {
    name  = "entityEntityIdRelatedGet_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_get_entity_id
  }
  property {
    name  = "entityEntityIdRelatedPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_post_body
  }
  property {
    name  = "entityEntityIdRelatedPost_entity"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_post_entity
  }
  property {
    name  = "entityEntityIdRelatedPost_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_post_entity_id
  }
  property {
    name  = "entityEntityIdRelatedRelatedEntityNameGet_entity"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_related_entity_name_get_entity
  }
  property {
    name  = "entityEntityIdRelatedRelatedEntityNameGet_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_related_entity_name_get_entity_id
  }
  property {
    name  = "entityEntityIdRelatedRelatedEntityNameGet_related_entity_name"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_entity_id_related_related_entity_name_get_related_entity_name
  }
  property {
    name  = "entityNameInPluralSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_name_in_plural_search_post_body
  }
  property {
    name  = "entityNameInPluralSearchPost_entity_name_in_plural"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_name_in_plural_search_post_entity_name_in_plural
  }
  property {
    name  = "entityTypeEntityIdFilesGet_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_get_entity_id
  }
  property {
    name  = "entityTypeEntityIdFilesGet_entity_type"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_get_entity_type
  }
  property {
    name  = "entityTypeEntityIdFilesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_post_body
  }
  property {
    name  = "entityTypeEntityIdFilesPost_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_post_entity_id
  }
  property {
    name  = "entityTypeEntityIdFilesPost_entity_type"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_post_entity_type
  }
  property {
    name  = "entityTypeEntityIdFilesS3SignedUrlGet_entity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_s3_signed_url_get_entity_id
  }
  property {
    name  = "entityTypeEntityIdFilesS3SignedUrlGet_entity_type"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_entity_type_entity_id_files_s3_signed_url_get_entity_type
  }
  property {
    name  = "leadsExampleLeadIdActivitiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_activities_post_body
  }
  property {
    name  = "leadsExampleLeadIdActivitiesPost_example_lead_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_activities_post_example_lead_id
  }
  property {
    name  = "leadsExampleLeadIdDelete_example_lead_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_delete_example_lead_id
  }
  property {
    name  = "leadsExampleLeadIdGet_example_lead_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_get_example_lead_id
  }
  property {
    name  = "leadsExampleLeadIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_put_body
  }
  property {
    name  = "leadsExampleLeadIdPut_example_lead_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_lead_id_put_example_lead_id
  }
  property {
    name  = "leadsExampleLeadconvertIdConvertPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_leadconvert_id_convert_post_body
  }
  property {
    name  = "leadsExampleLeadconvertIdConvertPost_example_leadconvert_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_example_leadconvert_id_convert_post_example_leadconvert_id
  }
  property {
    name  = "leadsPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_post_body
  }
  property {
    name  = "leadsSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_search_post_body
  }
  property {
    name  = "leadsUpsertPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_leads_upsert_put_body
  }
  property {
    name  = "opportunitiesDeleteOpportunityIdDelete_delete_opportunity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_delete_opportunity_id_delete_delete_opportunity_id
  }
  property {
    name  = "opportunitiesExampleOpportunityIdGet_example_opportunity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_example_opportunity_id_get_example_opportunity_id
  }
  property {
    name  = "opportunitiesExampleOpportunityIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_example_opportunity_id_put_body
  }
  property {
    name  = "opportunitiesExampleOpportunityIdPut_example_opportunity_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_example_opportunity_id_put_example_opportunity_id
  }
  property {
    name  = "opportunitiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_post_body
  }
  property {
    name  = "opportunitiesSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_opportunities_search_post_body
  }
  property {
    name  = "peopleExamplePersonIdActivitiesPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_activities_post_body
  }
  property {
    name  = "peopleExamplePersonIdActivitiesPost_example_person_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_activities_post_example_person_id
  }
  property {
    name  = "peopleExamplePersonIdDelete_example_person_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_delete_example_person_id
  }
  property {
    name  = "peopleExamplePersonIdGet_example_person_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_get_example_person_id
  }
  property {
    name  = "peopleExamplePersonIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_put_body
  }
  property {
    name  = "peopleExamplePersonIdPut_example_person_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_example_person_id_put_example_person_id
  }
  property {
    name  = "peopleFetchByEmailPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_fetch_by_email_post_body
  }
  property {
    name  = "peoplePost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_post_body
  }
  property {
    name  = "peopleSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_people_search_post_body
  }
  property {
    name  = "pipelineStagesPipelinePipelineIdGet_pipeline_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_pipeline_stages_pipeline_pipeline_id_get_pipeline_id
  }
  property {
    name  = "projectsDeleteProjectIdDelete_delete_project_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_delete_project_id_delete_delete_project_id
  }
  property {
    name  = "projectsExampleProjectIdGet_example_project_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_example_project_id_get_example_project_id
  }
  property {
    name  = "projectsExampleProjectIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_example_project_id_put_body
  }
  property {
    name  = "projectsExampleProjectIdPut_example_project_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_example_project_id_put_example_project_id
  }
  property {
    name  = "projectsPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_post_body
  }
  property {
    name  = "projectsSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_projects_search_post_body
  }
  property {
    name  = "relatedLinksConnectionIdDelete_connection_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_related_links_connection_id_delete_connection_id
  }
  property {
    name  = "relatedLinksGet_custom_field_definition_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_related_links_get_custom_field_definition_id
  }
  property {
    name  = "relatedLinksGet_source_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_related_links_get_source_id
  }
  property {
    name  = "relatedLinksGet_source_type"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_related_links_get_source_type
  }
  property {
    name  = "relatedLinksPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_related_links_post_body
  }
  property {
    name  = "rootPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_root_post_body
  }
  property {
    name  = "tasksDeleteTaskIdDelete_delete_task_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_delete_task_id_delete_delete_task_id
  }
  property {
    name  = "tasksExampleTaskIdGet_example_task_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_example_task_id_get_example_task_id
  }
  property {
    name  = "tasksExampleTaskIdPut_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_example_task_id_put_body
  }
  property {
    name  = "tasksExampleTaskIdPut_example_task_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_example_task_id_put_example_task_id
  }
  property {
    name  = "tasksPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_post_body
  }
  property {
    name  = "tasksSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_tasks_search_post_body
  }
  property {
    name  = "usersExampleUserIdGet_example_user_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_users_example_user_id_get_example_user_id
  }
  property {
    name  = "usersSearchPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_users_search_post_body
  }
  property {
    name  = "webhooksExampleWebhookIdDelete_example_webhook_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_webhooks_example_webhook_id_delete_example_webhook_id
  }
  property {
    name  = "webhooksPost_body"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_webhooks_post_body
  }
  property {
    name  = "webhooksexampleWebhookIdGet_example_webhook_id"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_webhooksexample_webhook_id_get_example_webhook_id
  }
  property {
    name  = "xPWAccessToken"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_x_pwaccess_token
  }
  property {
    name  = "xPWApplication"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_x_pwapplication
  }
  property {
    name  = "xPWUserEmail"
    type  = "string"
    value = var.connector-oai-copperdeveloperapi_property_x_pwuser_email
  }
}
