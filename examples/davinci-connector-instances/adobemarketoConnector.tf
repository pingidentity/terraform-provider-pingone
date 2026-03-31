resource "pingone_davinci_connector_instance" "adobemarketoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "adobemarketoConnector"
  }
  name = "My awesome adobemarketoConnector"
  property {
    name  = "addLeadsArray"
    type  = "string"
    value = var.adobemarketoconnector_property_add_leads_array
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.adobemarketoconnector_property_body_data
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.adobemarketoconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.adobemarketoconnector_property_client_secret
  }
  property {
    name  = "customAttributes"
    type  = "string"
    value = var.adobemarketoconnector_property_custom_attributes
  }
  property {
    name  = "customEndpoint"
    type  = "string"
    value = var.adobemarketoconnector_property_custom_endpoint
  }
  property {
    name  = "customQueryParams"
    type  = "string"
    value = var.adobemarketoconnector_property_custom_query_params
  }
  property {
    name  = "email"
    type  = "string"
    value = var.adobemarketoconnector_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.adobemarketoconnector_property_endpoint
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.adobemarketoconnector_property_first_name
  }
  property {
    name  = "folderList"
    type  = "string"
    value = var.adobemarketoconnector_property_folder_list
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.adobemarketoconnector_property_headers
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.adobemarketoconnector_property_last_name
  }
  property {
    name  = "leadId"
    type  = "string"
    value = var.adobemarketoconnector_property_lead_id
  }
  property {
    name  = "listId"
    type  = "string"
    value = var.adobemarketoconnector_property_list_id
  }
  property {
    name  = "listName"
    type  = "string"
    value = var.adobemarketoconnector_property_list_name
  }
  property {
    name  = "method"
    type  = "string"
    value = var.adobemarketoconnector_property_method
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.adobemarketoconnector_property_middle_name
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.adobemarketoconnector_property_phone
  }
  property {
    name  = "searchColumn"
    type  = "string"
    value = var.adobemarketoconnector_property_search_column
  }
  property {
    name  = "searchLeadArray"
    type  = "string"
    value = var.adobemarketoconnector_property_search_lead_array
  }
  property {
    name  = "searchListId"
    type  = "string"
    value = var.adobemarketoconnector_property_search_list_id
  }
  property {
    name  = "searchListMethod"
    type  = "string"
    value = var.adobemarketoconnector_property_search_list_method
  }
  property {
    name  = "searchListName"
    type  = "string"
    value = var.adobemarketoconnector_property_search_list_name
  }
  property {
    name  = "smartListToggle"
    type  = "string"
    value = var.adobemarketoconnector_property_smart_list_toggle
  }
}
