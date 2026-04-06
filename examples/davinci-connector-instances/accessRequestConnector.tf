resource "pingone_davinci_connector_instance" "accessRequestConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "accessRequestConnector"
  }
  name = "My awesome accessRequestConnector"
  property {
    name  = "apiRequestType"
    type  = "string"
    value = var.accessrequestconnector_property_api_request_type
  }
  property {
    name  = "applicationIds"
    type  = "string"
    value = var.accessrequestconnector_property_application_ids
  }
  property {
    name  = "baseURL"
    type  = "string"
    value = var.accessrequestconnector_property_base_url
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.accessrequestconnector_property_body_data
  }
  property {
    name  = "comment"
    type  = "string"
    value = var.accessrequestconnector_property_comment
  }
  property {
    name  = "customFilter"
    type  = "string"
    value = var.accessrequestconnector_property_custom_filter
  }
  property {
    name  = "endUserClientId"
    type  = "string"
    value = var.accessrequestconnector_property_end_user_client_id
  }
  property {
    name  = "endUserClientPrivateKey"
    type  = "string"
    value = var.accessrequestconnector_property_end_user_client_private_key
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.accessrequestconnector_property_endpoint
  }
  property {
    name  = "entitlementIds"
    type  = "string"
    value = var.accessrequestconnector_property_entitlement_ids
  }
  property {
    name  = "expiryDate"
    type  = "string"
    value = var.accessrequestconnector_property_expiry_date
  }
  property {
    name  = "givenName"
    type  = "string"
    value = var.accessrequestconnector_property_given_name
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.accessrequestconnector_property_headers
  }
  property {
    name  = "justification"
    type  = "string"
    value = var.accessrequestconnector_property_justification
  }
  property {
    name  = "mail"
    type  = "string"
    value = var.accessrequestconnector_property_mail
  }
  property {
    name  = "matchAttribute"
    type  = "string"
    value = var.accessrequestconnector_property_match_attribute
  }
  property {
    name  = "matchAttributes"
    type  = "string"
    value = var.accessrequestconnector_property_match_attributes
  }
  property {
    name  = "method"
    type  = "string"
    value = var.accessrequestconnector_property_method
  }
  property {
    name  = "password"
    type  = "string"
    value = var.accessrequestconnector_property_password
  }
  property {
    name  = "phaseName"
    type  = "string"
    value = var.accessrequestconnector_property_phase_name
  }
  property {
    name  = "priority"
    type  = "string"
    value = var.accessrequestconnector_property_priority
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.accessrequestconnector_property_query_parameters
  }
  property {
    name  = "realm"
    type  = "string"
    value = var.accessrequestconnector_property_realm
  }
  property {
    name  = "requestAction"
    type  = "string"
    value = var.accessrequestconnector_property_request_action
  }
  property {
    name  = "requestId"
    type  = "string"
    value = var.accessrequestconnector_property_request_id
  }
  property {
    name  = "requestType"
    type  = "string"
    value = var.accessrequestconnector_property_request_type
  }
  property {
    name  = "roleIds"
    type  = "string"
    value = var.accessrequestconnector_property_role_ids
  }
  property {
    name  = "serviceAccountId"
    type  = "string"
    value = var.accessrequestconnector_property_service_account_id
  }
  property {
    name  = "serviceAccountPrivateKey"
    type  = "string"
    value = var.accessrequestconnector_property_service_account_private_key
  }
  property {
    name  = "sn"
    type  = "string"
    value = var.accessrequestconnector_property_sn
  }
  property {
    name  = "useCustomFilter"
    type  = "string"
    value = var.accessrequestconnector_property_use_custom_filter
  }
  property {
    name  = "userAttributes"
    type  = "string"
    value = var.accessrequestconnector_property_user_attributes
  }
  property {
    name  = "userIdentifier"
    type  = "string"
    value = var.accessrequestconnector_property_user_identifier
  }
  property {
    name  = "userIds"
    type  = "string"
    value = var.accessrequestconnector_property_user_ids
  }
  property {
    name  = "userName"
    type  = "string"
    value = var.accessrequestconnector_property_user_name
  }
}
