resource "pingone_davinci_connector_instance" "jiraConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jiraConnector"
  }
  name = "My awesome jiraConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.jiraconnector_property_api_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.jiraconnector_property_api_url
  }
  property {
    name  = "assignee"
    type  = "string"
    value = var.jiraconnector_property_assignee
  }
  property {
    name  = "assigneeId"
    type  = "string"
    value = var.jiraconnector_property_assignee_id
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.jiraconnector_property_body_data
  }
  property {
    name  = "components"
    type  = "string"
    value = var.jiraconnector_property_components
  }
  property {
    name  = "customQueryParams"
    type  = "string"
    value = var.jiraconnector_property_custom_query_params
  }
  property {
    name  = "description"
    type  = "string"
    value = var.jiraconnector_property_description
  }
  property {
    name  = "dueDate"
    type  = "string"
    value = var.jiraconnector_property_due_date
  }
  property {
    name  = "email"
    type  = "string"
    value = var.jiraconnector_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.jiraconnector_property_endpoint
  }
  property {
    name  = "expand"
    type  = "string"
    value = var.jiraconnector_property_expand
  }
  property {
    name  = "fields"
    type  = "string"
    value = var.jiraconnector_property_fields
  }
  property {
    name  = "fieldsByKeys"
    type  = "string"
    value = var.jiraconnector_property_fields_by_keys
  }
  property {
    name  = "fixVersions"
    type  = "string"
    value = var.jiraconnector_property_fix_versions
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.jiraconnector_property_headers
  }
  property {
    name  = "issueId"
    type  = "string"
    value = var.jiraconnector_property_issue_id
  }
  property {
    name  = "issueType"
    type  = "string"
    value = var.jiraconnector_property_issue_type
  }
  property {
    name  = "jql"
    type  = "string"
    value = var.jiraconnector_property_jql
  }
  property {
    name  = "labels"
    type  = "string"
    value = var.jiraconnector_property_labels
  }
  property {
    name  = "loopingTransition"
    type  = "string"
    value = var.jiraconnector_property_looping_transition
  }
  property {
    name  = "maxResults"
    type  = "number"
    value = var.jiraconnector_property_max_results
  }
  property {
    name  = "method"
    type  = "string"
    value = var.jiraconnector_property_method
  }
  property {
    name  = "otherAttributes"
    type  = "string"
    value = var.jiraconnector_property_other_attributes
  }
  property {
    name  = "parent"
    type  = "string"
    value = var.jiraconnector_property_parent
  }
  property {
    name  = "priority"
    type  = "string"
    value = var.jiraconnector_property_priority
  }
  property {
    name  = "project"
    type  = "string"
    value = var.jiraconnector_property_project
  }
  property {
    name  = "projectKey"
    type  = "string"
    value = var.jiraconnector_property_project_key
  }
  property {
    name  = "queryParams"
    type  = "string"
    value = var.jiraconnector_property_query_params
  }
  property {
    name  = "reporter"
    type  = "string"
    value = var.jiraconnector_property_reporter
  }
  property {
    name  = "reporterId"
    type  = "string"
    value = var.jiraconnector_property_reporter_id
  }
  property {
    name  = "searchProperties"
    type  = "string"
    value = var.jiraconnector_property_search_properties
  }
  property {
    name  = "startAt"
    type  = "number"
    value = var.jiraconnector_property_start_at
  }
  property {
    name  = "summary"
    type  = "string"
    value = var.jiraconnector_property_summary
  }
  property {
    name  = "transitionId"
    type  = "string"
    value = var.jiraconnector_property_transition_id
  }
  property {
    name  = "updateIssueType"
    type  = "string"
    value = var.jiraconnector_property_update_issue_type
  }
  property {
    name  = "updateSummary"
    type  = "string"
    value = var.jiraconnector_property_update_summary
  }
  property {
    name  = "validateQuery"
    type  = "string"
    value = var.jiraconnector_property_validate_query
  }
  property {
    name  = "versions"
    type  = "string"
    value = var.jiraconnector_property_versions
  }
}
