resource "pingone_davinci_connector_instance" "salesforceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "salesforceConnector"
  }
  name = "My awesome salesforceConnector"
  property {
    name  = "adminUsername"
    type  = "string"
    value = var.salesforceconnector_property_admin_username
  }
  property {
    name  = "alias"
    type  = "string"
    value = var.salesforceconnector_property_alias
  }
  property {
    name  = "body"
    type  = "string"
    value = var.salesforceconnector_property_body
  }
  property {
    name  = "consumerKey"
    type  = "string"
    value = var.salesforceconnector_property_consumer_key
  }
  property {
    name  = "documentDescription"
    type  = "string"
    value = var.salesforceconnector_property_document_description
  }
  property {
    name  = "documentFields"
    type  = "string"
    value = var.salesforceconnector_property_document_fields
  }
  property {
    name  = "documentFolder"
    type  = "string"
    value = var.salesforceconnector_property_document_folder
  }
  property {
    name  = "documentFolderId"
    type  = "string"
    value = var.salesforceconnector_property_document_folder_id
  }
  property {
    name  = "documentKeywords"
    type  = "string"
    value = var.salesforceconnector_property_document_keywords
  }
  property {
    name  = "documentName"
    type  = "string"
    value = var.salesforceconnector_property_document_name
  }
  property {
    name  = "domainName"
    type  = "string"
    value = var.salesforceconnector_property_domain_name
  }
  property {
    name  = "email"
    type  = "string"
    value = var.salesforceconnector_property_email
  }
  property {
    name  = "emailEncodingKey"
    type  = "string"
    value = var.salesforceconnector_property_email_encoding_key
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.salesforceconnector_property_endpoint
  }
  property {
    name  = "environment"
    type  = "string"
    value = var.salesforceconnector_property_environment
  }
  property {
    name  = "file"
    type  = "string"
    value = var.salesforceconnector_property_file
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.salesforceconnector_property_first_name
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.salesforceconnector_property_headers
  }
  property {
    name  = "languageLocaleKey"
    type  = "string"
    value = var.salesforceconnector_property_language_locale_key
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.salesforceconnector_property_last_name
  }
  property {
    name  = "localeSidKey"
    type  = "string"
    value = var.salesforceconnector_property_locale_sid_key
  }
  property {
    name  = "method"
    type  = "string"
    value = var.salesforceconnector_property_method
  }
  property {
    name  = "nextRecordsUrl"
    type  = "string"
    value = var.salesforceconnector_property_next_records_url
  }
  property {
    name  = "objectType"
    type  = "string"
    value = var.salesforceconnector_property_object_type
  }
  property {
    name  = "privateKey"
    type  = "string"
    value = var.salesforceconnector_property_private_key
  }
  property {
    name  = "profile"
    type  = "string"
    value = var.salesforceconnector_property_profile
  }
  property {
    name  = "profileId"
    type  = "string"
    value = var.salesforceconnector_property_profile_id
  }
  property {
    name  = "query"
    type  = "string"
    value = var.salesforceconnector_property_query
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.salesforceconnector_property_query_parameters
  }
  property {
    name  = "recordFieldsRead"
    type  = "string"
    value = var.salesforceconnector_property_record_fields_read
  }
  property {
    name  = "recordFieldsUpdate"
    type  = "string"
    value = var.salesforceconnector_property_record_fields_update
  }
  property {
    name  = "recordId"
    type  = "string"
    value = var.salesforceconnector_property_record_id
  }
  property {
    name  = "timeZoneSidKey"
    type  = "string"
    value = var.salesforceconnector_property_time_zone_sid_key
  }
  property {
    name  = "userFieldsCreate"
    type  = "string"
    value = var.salesforceconnector_property_user_fields_create
  }
  property {
    name  = "userFieldsRead"
    type  = "string"
    value = var.salesforceconnector_property_user_fields_read
  }
  property {
    name  = "userFieldsUpdate"
    type  = "string"
    value = var.salesforceconnector_property_user_fields_update
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.salesforceconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.salesforceconnector_property_username
  }
}
