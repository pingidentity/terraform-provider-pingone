resource "pingone_davinci_connector_instance" "connectorZscaler" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZscaler"
  }
  name = "My awesome connectorZscaler"
  property {
    name  = "adminUser"
    type  = "string"
    value = var.connectorzscaler_property_admin_user
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connectorzscaler_property_base_path
  }
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "comments"
    type  = "string"
    value = var.connectorzscaler_property_comments
  }
  property {
    name  = "departmentComments"
    type  = "string"
    value = var.connectorzscaler_property_department_comments
  }
  property {
    name  = "departmentDeleted"
    type  = "string"
    value = var.connectorzscaler_property_department_deleted
  }
  property {
    name  = "departmentID"
    type  = "string"
    value = var.connectorzscaler_property_department_id
  }
  property {
    name  = "departmentIDPID"
    type  = "string"
    value = var.connectorzscaler_property_department_idpid
  }
  property {
    name  = "departmentName"
    type  = "string"
    value = var.connectorzscaler_property_department_name
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorzscaler_property_email
  }
  property {
    name  = "groupsComments"
    type  = "string"
    value = var.connectorzscaler_property_groups_comments
  }
  property {
    name  = "groupsID"
    type  = "string"
    value = var.connectorzscaler_property_groups_id
  }
  property {
    name  = "groupsIDPID"
    type  = "string"
    value = var.connectorzscaler_property_groups_idpid
  }
  property {
    name  = "groupsName"
    type  = "string"
    value = var.connectorzscaler_property_groups_name
  }
  property {
    name  = "limitSearch"
    type  = "string"
    value = var.connectorzscaler_property_limit_search
  }
  property {
    name  = "name"
    type  = "string"
    value = var.connectorzscaler_property_name
  }
  property {
    name  = "page"
    type  = "string"
    value = var.connectorzscaler_property_page
  }
  property {
    name  = "pageSize"
    type  = "string"
    value = var.connectorzscaler_property_page_size
  }
  property {
    name  = "password"
    type  = "string"
    value = var.connectorzscaler_property_password
  }
  property {
    name  = "search"
    type  = "string"
    value = var.connectorzscaler_property_search
  }
  property {
    name  = "tempAuthEmail"
    type  = "string"
    value = var.connectorzscaler_property_temp_auth_email
  }
  property {
    name  = "type"
    type  = "string"
    value = var.connectorzscaler_property_type
  }
  property {
    name  = "userID"
    type  = "string"
    value = var.connectorzscaler_property_user_id
  }
  property {
    name  = "zscalerAPIkey"
    type  = "string"
    value = var.zscaler_api_key
  }
  property {
    name  = "zscalerPassword"
    type  = "string"
    value = var.connectorzscaler_property_zscaler_password
  }
  property {
    name  = "zscalerUsername"
    type  = "string"
    value = var.connectorzscaler_property_zscaler_username
  }
}
