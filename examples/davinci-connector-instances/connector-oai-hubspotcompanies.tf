resource "pingone_davinci_connector_instance" "connector-oai-hubspotcompanies" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-hubspotcompanies"
  }
  name = "My awesome connector-oai-hubspotcompanies"
  property {
    name  = "authBearerToken"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_auth_bearer_token
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_base_path
  }
  property {
    name  = "deleteCrmV3ObjectsCompaniesCompanyIdArchive_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_delete_crm_v3_objects_companies_company_id_archive_company_id
  }
  property {
    name  = "deleteCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdArchive_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_delete_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_archive_company_id
  }
  property {
    name  = "deleteCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdArchive_to_object_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_delete_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_archive_to_object_id
  }
  property {
    name  = "deleteCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdArchive_to_object_type"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_delete_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_archive_to_object_type
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_archived"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_archived
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_associations"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_associations
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_company_id
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_id_property"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_id_property
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_properties
  }
  property {
    name  = "getCrmV3ObjectsCompaniesCompanyIdGetById_properties_with_history"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_company_id_get_by_id_properties_with_history
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_after"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_after
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_archived"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_archived
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_associations"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_associations
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_limit"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_limit
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_properties
  }
  property {
    name  = "getCrmV3ObjectsCompaniesGetPage_properties_with_history"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v3_objects_companies_get_page_properties_with_history
  }
  property {
    name  = "getCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeGetAll_after"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v4_objects_companies_company_id_associations_to_object_type_get_all_after
  }
  property {
    name  = "getCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeGetAll_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v4_objects_companies_company_id_associations_to_object_type_get_all_company_id
  }
  property {
    name  = "getCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeGetAll_limit"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v4_objects_companies_company_id_associations_to_object_type_get_all_limit
  }
  property {
    name  = "getCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeGetAll_to_object_type"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_get_crm_v4_objects_companies_company_id_associations_to_object_type_get_all_to_object_type
  }
  property {
    name  = "patchCrmV3ObjectsCompaniesCompanyIdUpdate_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_patch_crm_v3_objects_companies_company_id_update_company_id
  }
  property {
    name  = "patchCrmV3ObjectsCompaniesCompanyIdUpdate_id_property"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_patch_crm_v3_objects_companies_company_id_update_id_property
  }
  property {
    name  = "patchCrmV3ObjectsCompaniesCompanyIdUpdate_simplePublicObjectInput_SimplePublicObjectInput_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_patch_crm_v3_objects_companies_company_id_update_simple_public_object_input_simple_public_object_input_properties
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchArchiveArchive_batchInputSimplePublicObjectId_BatchInputSimplePublicObjectId_inputs"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_archive_archive_batch_input_simple_public_object_id_batch_input_simple_public_object_id_inputs
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchCreateCreate_batchInputSimplePublicObjectInput_BatchInputSimplePublicObjectInput_inputs"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_create_create_batch_input_simple_public_object_input_batch_input_simple_public_object_input_inputs
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchReadRead_archived"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_read_read_archived
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchReadRead_batchReadInputSimplePublicObjectId_BatchReadInputSimplePublicObjectId_idProperty"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_read_read_batch_read_input_simple_public_object_id_batch_read_input_simple_public_object_id_id_property
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchReadRead_batchReadInputSimplePublicObjectId_BatchReadInputSimplePublicObjectId_inputs"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_read_read_batch_read_input_simple_public_object_id_batch_read_input_simple_public_object_id_inputs
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchReadRead_batchReadInputSimplePublicObjectId_BatchReadInputSimplePublicObjectId_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_read_read_batch_read_input_simple_public_object_id_batch_read_input_simple_public_object_id_properties
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchReadRead_batchReadInputSimplePublicObjectId_BatchReadInputSimplePublicObjectId_propertiesWithHistory"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_read_read_batch_read_input_simple_public_object_id_batch_read_input_simple_public_object_id_properties_with_history
  }
  property {
    name  = "postCrmV3ObjectsCompaniesBatchUpdateUpdate_batchInputSimplePublicObjectBatchInput_BatchInputSimplePublicObjectBatchInput_inputs"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_batch_update_update_batch_input_simple_public_object_batch_input_batch_input_simple_public_object_batch_input_inputs
  }
  property {
    name  = "postCrmV3ObjectsCompaniesCreate_simplePublicObjectInput_SimplePublicObjectInput_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_create_simple_public_object_input_simple_public_object_input_properties
  }
  property {
    name  = "postCrmV3ObjectsCompaniesMergeMerge_publicMergeInput_PublicMergeInput_objectIdToMerge"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_merge_merge_public_merge_input_public_merge_input_object_id_to_merge
  }
  property {
    name  = "postCrmV3ObjectsCompaniesMergeMerge_publicMergeInput_PublicMergeInput_primaryObjectId"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_merge_merge_public_merge_input_public_merge_input_primary_object_id
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_after"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_after
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_filterGroups"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_filter_groups
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_limit"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_limit
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_properties"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_properties
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_query"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_query
  }
  property {
    name  = "postCrmV3ObjectsCompaniesSearchDoSearch_publicObjectSearchRequest_PublicObjectSearchRequest_sorts"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_post_crm_v3_objects_companies_search_do_search_public_object_search_request_public_object_search_request_sorts
  }
  property {
    name  = "putCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdCreate_association_spec"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_put_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_create_association_spec
  }
  property {
    name  = "putCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdCreate_company_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_put_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_create_company_id
  }
  property {
    name  = "putCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdCreate_to_object_id"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_put_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_create_to_object_id
  }
  property {
    name  = "putCrmV4ObjectsCompaniesCompanyIdAssociationsToObjectTypeToObjectIdCreate_to_object_type"
    type  = "string"
    value = var.connector-oai-hubspotcompanies_property_put_crm_v4_objects_companies_company_id_associations_to_object_type_to_object_id_create_to_object_type
  }
}
