resource "pingone_davinci_connector_instance" "connectorIdentityNow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdentityNow"
  }
  name = "My awesome connectorIdentityNow"
  property {
    name  = "accessToken"
    type  = "string"
    value = var.connectoridentitynow_property_access_token
  }
  property {
    name  = "accountName"
    type  = "string"
    value = var.connectoridentitynow_property_account_name
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.connectoridentitynow_property_api_version
  }
  property {
    name  = "approvalItemId"
    type  = "string"
    value = var.connectoridentitynow_property_approval_item_id
  }
  property {
    name  = "authMethod"
    type  = "string"
    value = var.connectoridentitynow_property_auth_method
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectoridentitynow_property_body
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.connectoridentitynow_property_client_id
  }
  property {
    name  = "clientIdPAT"
    type  = "string"
    value = var.connectoridentitynow_property_client_id_pat
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectoridentitynow_property_client_secret
  }
  property {
    name  = "clientSecretPAT"
    type  = "string"
    value = var.connectoridentitynow_property_client_secret_pat
  }
  property {
    name  = "comment"
    type  = "string"
    value = var.connectoridentitynow_property_comment
  }
  property {
    name  = "data"
    type  = "string"
    value = var.connectoridentitynow_property_data
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectoridentitynow_property_email
  }
  property {
    name  = "endDate"
    type  = "string"
    value = var.connectoridentitynow_property_end_date
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectoridentitynow_property_endpoint
  }
  property {
    name  = "familyName"
    type  = "string"
    value = var.connectoridentitynow_property_family_name
  }
  property {
    name  = "filters"
    type  = "string"
    value = var.connectoridentitynow_property_filters
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.connectoridentitynow_property_first_name
  }
  property {
    name  = "givenName"
    type  = "string"
    value = var.connectoridentitynow_property_given_name
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectoridentitynow_property_headers
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.connectoridentitynow_property_last_name
  }
  property {
    name  = "limit"
    type  = "string"
    value = var.connectoridentitynow_property_limit
  }
  property {
    name  = "manager"
    type  = "string"
    value = var.connectoridentitynow_property_manager
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectoridentitynow_property_method
  }
  property {
    name  = "offset"
    type  = "string"
    value = var.connectoridentitynow_property_offset
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectoridentitynow_property_phone
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectoridentitynow_property_query_parameters
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.connectoridentitynow_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.connectoridentitynow_property_skip_button_press
  }
  property {
    name  = "sourceId"
    type  = "string"
    value = var.connectoridentitynow_property_source_id
  }
  property {
    name  = "sourceName"
    type  = "string"
    value = var.connectoridentitynow_property_source_name
  }
  property {
    name  = "startDate"
    type  = "string"
    value = var.connectoridentitynow_property_start_date
  }
  property {
    name  = "tenant"
    type  = "string"
    value = var.connectoridentitynow_property_tenant
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.connectoridentitynow_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectoridentitynow_property_username
  }
}
