resource "pingone_davinci_connector_instance" "connectorTableau" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTableau"
  }
  name = "My awesome connectorTableau"
  property {
    name  = "addFlowPermissionsRequestBody"
    type  = "string"
    value = var.connectortableau_property_add_flow_permissions_request_body
  }
  property {
    name  = "addUsertoSiteRequestBody"
    type  = "string"
    value = var.connectortableau_property_add_userto_site_request_body
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.connectortableau_property_api_version
  }
  property {
    name  = "authId"
    type  = "string"
    value = var.connectortableau_property_auth_id
  }
  property {
    name  = "createScheduleBody"
    type  = "string"
    value = var.connectortableau_property_create_schedule_body
  }
  property {
    name  = "datasourceId"
    type  = "string"
    value = var.connectortableau_property_datasource_id
  }
  property {
    name  = "flowId"
    type  = "string"
    value = var.connectortableau_property_flow_id
  }
  property {
    name  = "groupId"
    type  = "string"
    value = var.connectortableau_property_group_id
  }
  property {
    name  = "jobId"
    type  = "string"
    value = var.connectortableau_property_job_id
  }
  property {
    name  = "scheduleId"
    type  = "string"
    value = var.connectortableau_property_schedule_id
  }
  property {
    name  = "serverUrl"
    type  = "string"
    value = var.connectortableau_property_server_url
  }
  property {
    name  = "siteId"
    type  = "string"
    value = var.connectortableau_property_site_id
  }
  property {
    name  = "taskId"
    type  = "string"
    value = var.connectortableau_property_task_id
  }
  property {
    name  = "updateScheduleRequestBody"
    type  = "string"
    value = var.connectortableau_property_update_schedule_request_body
  }
  property {
    name  = "updateUserRequestBody"
    type  = "string"
    value = var.connectortableau_property_update_user_request_body
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.connectortableau_property_user_id
  }
  property {
    name  = "workbookId"
    type  = "string"
    value = var.connectortableau_property_workbook_id
  }
}
