resource "pingone_davinci_connector_instance" "connectorSaviyntFlow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSaviyntFlow"
  }
  name = "My awesome connectorSaviyntFlow"
  property {
    name  = "additionalUserProperties"
    type  = "string"
    value = var.connectorsaviyntflow_property_additional_user_properties
  }
  property {
    name  = "advsearchcriteria"
    type  = "string"
    value = var.connectorsaviyntflow_property_advsearchcriteria
  }
  property {
    name  = "city"
    type  = "string"
    value = var.connectorsaviyntflow_property_city
  }
  property {
    name  = "comments"
    type  = "string"
    value = var.connectorsaviyntflow_property_comments
  }
  property {
    name  = "companyname"
    type  = "string"
    value = var.connectorsaviyntflow_property_companyname
  }
  property {
    name  = "country"
    type  = "string"
    value = var.connectorsaviyntflow_property_country
  }
  property {
    name  = "designation"
    type  = "string"
    value = var.connectorsaviyntflow_property_designation
  }
  property {
    name  = "domainName"
    type  = "string"
    value = var.connectorsaviyntflow_property_domain_name
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorsaviyntflow_property_email
  }
  property {
    name  = "employeeType"
    type  = "string"
    value = var.connectorsaviyntflow_property_employee_type
  }
  property {
    name  = "enabled"
    type  = "string"
    value = var.connectorsaviyntflow_property_enabled
  }
  property {
    name  = "enddate"
    type  = "string"
    value = var.connectorsaviyntflow_property_enddate
  }
  property {
    name  = "filtercriteria"
    type  = "string"
    value = var.connectorsaviyntflow_property_filtercriteria
  }
  property {
    name  = "firstname"
    type  = "string"
    value = var.connectorsaviyntflow_property_firstname
  }
  property {
    name  = "lastname"
    type  = "string"
    value = var.connectorsaviyntflow_property_lastname
  }
  property {
    name  = "manager"
    type  = "string"
    value = var.connectorsaviyntflow_property_manager
  }
  property {
    name  = "max"
    type  = "string"
    value = var.connectorsaviyntflow_property_max
  }
  property {
    name  = "middlename"
    type  = "string"
    value = var.connectorsaviyntflow_property_middlename
  }
  property {
    name  = "name"
    type  = "string"
    value = var.connectorsaviyntflow_property_name
  }
  property {
    name  = "offset"
    type  = "string"
    value = var.connectorsaviyntflow_property_offset
  }
  property {
    name  = "order"
    type  = "string"
    value = var.connectorsaviyntflow_property_order
  }
  property {
    name  = "path"
    type  = "string"
    value = var.connectorsaviyntflow_property_path
  }
  property {
    name  = "phonenumber"
    type  = "string"
    value = var.connectorsaviyntflow_property_phonenumber
  }
  property {
    name  = "rank"
    type  = "string"
    value = var.connectorsaviyntflow_property_rank
  }
  property {
    name  = "reqaction"
    type  = "string"
    value = var.connectorsaviyntflow_property_reqaction
  }
  property {
    name  = "requestKey"
    type  = "string"
    value = var.connectorsaviyntflow_property_request_key
  }
  property {
    name  = "requestid"
    type  = "string"
    value = var.connectorsaviyntflow_property_requestid
  }
  property {
    name  = "responsefields"
    type  = "string"
    value = var.connectorsaviyntflow_property_responsefields
  }
  property {
    name  = "rolename"
    type  = "string"
    value = var.connectorsaviyntflow_property_rolename
  }
  property {
    name  = "roles"
    type  = "string"
    value = var.connectorsaviyntflow_property_roles
  }
  property {
    name  = "saviyntPassword"
    type  = "string"
    value = var.connectorsaviyntflow_property_saviynt_password
  }
  property {
    name  = "saviyntUserName"
    type  = "string"
    value = var.connectorsaviyntflow_property_saviynt_user_name
  }
  property {
    name  = "searchCriteria"
    type  = "string"
    value = var.connectorsaviyntflow_property_search_criteria
  }
  property {
    name  = "showsecurityanswers"
    type  = "string"
    value = var.connectorsaviyntflow_property_showsecurityanswers
  }
  property {
    name  = "sort"
    type  = "string"
    value = var.connectorsaviyntflow_property_sort
  }
  property {
    name  = "startdate"
    type  = "string"
    value = var.connectorsaviyntflow_property_startdate
  }
  property {
    name  = "state"
    type  = "string"
    value = var.connectorsaviyntflow_property_state
  }
  property {
    name  = "status"
    type  = "string"
    value = var.connectorsaviyntflow_property_status
  }
  property {
    name  = "statuskey"
    type  = "string"
    value = var.connectorsaviyntflow_property_statuskey
  }
  property {
    name  = "systemusername"
    type  = "string"
    value = var.connectorsaviyntflow_property_systemusername
  }
  property {
    name  = "type"
    type  = "string"
    value = var.connectorsaviyntflow_property_type
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectorsaviyntflow_property_username
  }
  property {
    name  = "value"
    type  = "string"
    value = var.connectorsaviyntflow_property_value
  }
}
