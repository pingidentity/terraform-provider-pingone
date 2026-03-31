resource "pingone_davinci_connector_instance" "googleWorkSpaceAdminConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleWorkSpaceAdminConnector"
  }
  name = "My awesome googleWorkSpaceAdminConnector"
  property {
    name  = "action"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_action
  }
  property {
    name  = "additionalUserProperties"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_additional_user_properties
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_body_data
  }
  property {
    name  = "customerId"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_customer_id
  }
  property {
    name  = "deviceProjection"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_device_projection
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_endpoint
  }
  property {
    name  = "familyName"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_family_name
  }
  property {
    name  = "givenName"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_given_name
  }
  property {
    name  = "groupKey"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_group_key
  }
  property {
    name  = "headersForm"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_headers_form
  }
  property {
    name  = "iss"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_iss
  }
  property {
    name  = "licenseUserId"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_license_user_id
  }
  property {
    name  = "maxResults"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_max_results
  }
  property {
    name  = "memberDeliverySettings"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_member_delivery_settings
  }
  property {
    name  = "memberEmail"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_member_email
  }
  property {
    name  = "memberKey"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_member_key
  }
  property {
    name  = "memberRole"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_member_role
  }
  property {
    name  = "memberType"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_member_type
  }
  property {
    name  = "method"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_method
  }
  property {
    name  = "orderBy"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_order_by
  }
  property {
    name  = "pageToken"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_page_token
  }
  property {
    name  = "paramsForm"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_params_form
  }
  property {
    name  = "password"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_password
  }
  property {
    name  = "primaryEmail"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_primary_email
  }
  property {
    name  = "privateKey"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_private_key
  }
  property {
    name  = "productId"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_product_id
  }
  property {
    name  = "projection"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_projection
  }
  property {
    name  = "query"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_query
  }
  property {
    name  = "resourceId"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_resource_id
  }
  property {
    name  = "skuId"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_sku_id
  }
  property {
    name  = "sortOrder"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_sort_order
  }
  property {
    name  = "sub"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_sub
  }
  property {
    name  = "userCustomFieldMask"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_user_custom_field_mask
  }
  property {
    name  = "userFieldMask"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_user_field_mask
  }
  property {
    name  = "userKey"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_user_key
  }
  property {
    name  = "viewType"
    type  = "string"
    value = var.googleworkspaceadminconnector_property_view_type
  }
}
