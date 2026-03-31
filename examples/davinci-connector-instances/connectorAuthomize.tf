resource "pingone_davinci_connector_instance" "connectorAuthomize" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAuthomize"
  }
  name = "My awesome connectorAuthomize"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorauthomize_property_api_key
  }
  property {
    name  = "incidentComment"
    type  = "string"
    value = var.connectorauthomize_property_incident_comment
  }
  property {
    name  = "incidentID"
    type  = "string"
    value = var.connectorauthomize_property_incident_id
  }
  property {
    name  = "incidentSeverity"
    type  = "string"
    value = var.connectorauthomize_property_incident_severity
  }
  property {
    name  = "incidentStatusUpdate"
    type  = "string"
    value = var.connectorauthomize_property_incident_status_update
  }
  property {
    name  = "searchPolicyId"
    type  = "string"
    value = var.connectorauthomize_property_search_policy_id
  }
  property {
    name  = "statusUpdate"
    type  = "string"
    value = var.connectorauthomize_property_status_update
  }
}
