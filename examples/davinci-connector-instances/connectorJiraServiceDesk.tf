resource "pingone_davinci_connector_instance" "connectorJiraServiceDesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJiraServiceDesk"
  }
  name = "My awesome connectorJiraServiceDesk"
  property {
    name  = "JIRAServiceDeskAuth"
    type  = "string"
    value = var.jira_service_desk_auth
  }
  property {
    name  = "JIRAServiceDeskCreateData"
    type  = "string"
    value = var.jira_service_desk_create_data
  }
  property {
    name  = "JIRAServiceDeskIssueID"
    type  = "string"
    value = var.connectorjiraservicedesk_property_jiraservice_desk_issue_id
  }
  property {
    name  = "JIRAServiceDeskURL"
    type  = "string"
    value = var.jira_service_desk_url
  }
  property {
    name  = "JIRAServiceDeskUpdateData"
    type  = "string"
    value = var.jira_service_desk_update_data
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorjiraservicedesk_property_body
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorjiraservicedesk_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorjiraservicedesk_property_headers
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorjiraservicedesk_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorjiraservicedesk_property_query_parameters
  }
}
