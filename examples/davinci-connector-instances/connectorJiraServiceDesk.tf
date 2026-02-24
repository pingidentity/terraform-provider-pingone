resource "pingone_davinci_connector_instance" "connectorJiraServiceDesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJiraServiceDesk"
  }
  name = "My awesome connectorJiraServiceDesk"
  properties = jsonencode({
    "JIRAServiceDeskAuth" = var.jira_service_desk_auth
    "JIRAServiceDeskCreateData" = var.jira_service_desk_create_data
    "JIRAServiceDeskURL" = var.jira_service_desk_url
    "JIRAServiceDeskUpdateData" = var.jira_service_desk_update_data
    "method" = var.connectorjiraservicedesk_property_method
  })
}
