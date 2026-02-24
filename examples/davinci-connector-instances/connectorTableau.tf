resource "pingone_davinci_connector_instance" "connectorTableau" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTableau"
  }
  name = "My awesome connectorTableau"
  properties = jsonencode({
    "addFlowPermissionsRequestBody" = var.connectortableau_property_add_flow_permissions_request_body
    "addUsertoSiteRequestBody" = var.connectortableau_property_add_userto_site_request_body
    "apiVersion" = var.connectortableau_property_api_version
    "authId" = var.connectortableau_property_auth_id
    "createScheduleBody" = var.connectortableau_property_create_schedule_body
    "datasourceId" = var.connectortableau_property_datasource_id
    "flowId" = var.connectortableau_property_flow_id
    "groupId" = var.connectortableau_property_group_id
    "jobId" = var.connectortableau_property_job_id
    "scheduleId" = var.connectortableau_property_schedule_id
    "serverUrl" = var.connectortableau_property_server_url
    "siteId" = var.connectortableau_property_site_id
    "taskId" = var.connectortableau_property_task_id
    "updateScheduleRequestBody" = var.connectortableau_property_update_schedule_request_body
    "updateUserRequestBody" = var.connectortableau_property_update_user_request_body
    "userId" = var.connectortableau_property_user_id
    "workbookId" = var.connectortableau_property_workbook_id
  })
}
