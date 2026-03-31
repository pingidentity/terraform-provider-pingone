resource "pingone_davinci_connector_instance" "connector443id" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector443id"
  }
  name = "My awesome connector443id"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connector443id_property_api_key
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connector443id_property_email
  }
  property {
    name  = "eventType"
    type  = "string"
    value = var.connector443id_property_event_type
  }
  property {
    name  = "identifiers"
    type  = "string"
    value = var.connector443id_property_identifiers
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.connector443id_property_ip
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connector443id_property_phone
  }
  property {
    name  = "ruleUuid"
    type  = "string"
    value = var.connector443id_property_rule_uuid
  }
  property {
    name  = "timestamp"
    type  = "string"
    value = var.connector443id_property_timestamp
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.connector443id_property_user_agent
  }
}
