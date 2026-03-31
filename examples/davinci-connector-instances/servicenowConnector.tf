resource "pingone_davinci_connector_instance" "servicenowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "servicenowConnector"
  }
  name = "My awesome servicenowConnector"
  property {
    name  = "active"
    type  = "string"
    value = var.servicenowconnector_property_active
  }
  property {
    name  = "adminUsername"
    type  = "string"
    value = var.servicenowconnector_property_admin_username
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.servicenowconnector_property_api_url
  }
  property {
    name  = "callerId"
    type  = "string"
    value = var.servicenowconnector_property_caller_id
  }
  property {
    name  = "city"
    type  = "string"
    value = var.servicenowconnector_property_city
  }
  property {
    name  = "customAttributes"
    type  = "string"
    value = var.servicenowconnector_property_custom_attributes
  }
  property {
    name  = "email"
    type  = "string"
    value = var.servicenowconnector_property_email
  }
  property {
    name  = "emailOptional"
    type  = "string"
    value = var.servicenowconnector_property_email_optional
  }
  property {
    name  = "employeeNumber"
    type  = "string"
    value = var.servicenowconnector_property_employee_number
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.servicenowconnector_property_first_name
  }
  property {
    name  = "firstNameOptional"
    type  = "string"
    value = var.servicenowconnector_property_first_name_optional
  }
  property {
    name  = "group"
    type  = "string"
    value = var.servicenowconnector_property_group
  }
  property {
    name  = "incidentBusinessService"
    type  = "string"
    value = var.servicenowconnector_property_incident_business_service
  }
  property {
    name  = "incidentCategory"
    type  = "string"
    value = var.servicenowconnector_property_incident_category
  }
  property {
    name  = "incidentConfigItem"
    type  = "string"
    value = var.servicenowconnector_property_incident_config_item
  }
  property {
    name  = "incidentContactType"
    type  = "string"
    value = var.servicenowconnector_property_incident_contact_type
  }
  property {
    name  = "incidentDescription"
    type  = "string"
    value = var.servicenowconnector_property_incident_description
  }
  property {
    name  = "incidentId"
    type  = "string"
    value = var.servicenowconnector_property_incident_id
  }
  property {
    name  = "incidentImpact"
    type  = "string"
    value = var.servicenowconnector_property_incident_impact
  }
  property {
    name  = "incidentNumber"
    type  = "string"
    value = var.servicenowconnector_property_incident_number
  }
  property {
    name  = "incidentShortDescription"
    type  = "string"
    value = var.servicenowconnector_property_incident_short_description
  }
  property {
    name  = "incidentState"
    type  = "string"
    value = var.servicenowconnector_property_incident_state
  }
  property {
    name  = "incidentSubcategory"
    type  = "string"
    value = var.servicenowconnector_property_incident_subcategory
  }
  property {
    name  = "incidentSwitcher"
    type  = "string"
    value = var.servicenowconnector_property_incident_switcher
  }
  property {
    name  = "incidentUrgency"
    type  = "string"
    value = var.servicenowconnector_property_incident_urgency
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.servicenowconnector_property_last_name
  }
  property {
    name  = "lastNameOptional"
    type  = "string"
    value = var.servicenowconnector_property_last_name_optional
  }
  property {
    name  = "lockedOut"
    type  = "string"
    value = var.servicenowconnector_property_locked_out
  }
  property {
    name  = "middleName"
    type  = "string"
    value = var.servicenowconnector_property_middle_name
  }
  property {
    name  = "middleNameOptional"
    type  = "string"
    value = var.servicenowconnector_property_middle_name_optional
  }
  property {
    name  = "password"
    type  = "string"
    value = var.servicenowconnector_property_password
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.servicenowconnector_property_phone
  }
  property {
    name  = "queryParams"
    type  = "string"
    value = var.servicenowconnector_property_query_params
  }
  property {
    name  = "readIncidentId"
    type  = "string"
    value = var.servicenowconnector_property_read_incident_id
  }
  property {
    name  = "sysParmLimit"
    type  = "string"
    value = var.servicenowconnector_property_sys_parm_limit
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.servicenowconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.servicenowconnector_property_username
  }
  property {
    name  = "usernameOptional"
    type  = "string"
    value = var.servicenowconnector_property_username_optional
  }
}
