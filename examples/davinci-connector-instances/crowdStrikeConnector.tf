resource "pingone_davinci_connector_instance" "crowdStrikeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "crowdStrikeConnector"
  }
  name = "My awesome crowdStrikeConnector"
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.crowdstrikeconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.crowdstrikeconnector_property_client_secret
  }
  property {
    name  = "deviceIdDeviceManaged"
    type  = "string"
    value = var.crowdstrikeconnector_property_device_id_device_managed
  }
  property {
    name  = "deviceIdIncidentScore"
    type  = "string"
    value = var.crowdstrikeconnector_property_device_id_incident_score
  }
  property {
    name  = "deviceIds"
    type  = "string"
    value = var.crowdstrikeconnector_property_device_ids
  }
  property {
    name  = "domainForAnalysis"
    type  = "string"
    value = var.crowdstrikeconnector_property_domain_for_analysis
  }
  property {
    name  = "email"
    type  = "string"
    value = var.crowdstrikeconnector_property_email
  }
  property {
    name  = "filter"
    type  = "string"
    value = var.crowdstrikeconnector_property_filter
  }
  property {
    name  = "incidentIds"
    type  = "string"
    value = var.crowdstrikeconnector_property_incident_ids
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.crowdstrikeconnector_property_ip
  }
  property {
    name  = "lastSeenDays"
    type  = "string"
    value = var.crowdstrikeconnector_property_last_seen_days
  }
  property {
    name  = "limit"
    type  = "string"
    value = var.crowdstrikeconnector_property_limit
  }
  property {
    name  = "offset"
    type  = "string"
    value = var.crowdstrikeconnector_property_offset
  }
  property {
    name  = "searchLoginDays"
    type  = "string"
    value = var.crowdstrikeconnector_property_search_login_days
  }
  property {
    name  = "username"
    type  = "string"
    value = var.crowdstrikeconnector_property_username
  }
  property {
    name  = "usernameForAnalysis"
    type  = "string"
    value = var.crowdstrikeconnector_property_username_for_analysis
  }
}
