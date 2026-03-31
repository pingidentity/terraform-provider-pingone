resource "pingone_davinci_connector_instance" "connectorFreshdesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshdesk"
  }
  name = "My awesome connectorFreshdesk"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorfreshdesk_property_api_key
  }
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "contactAddress"
    type  = "string"
    value = var.connectorfreshdesk_property_contact_address
  }
  property {
    name  = "contactId"
    type  = "string"
    value = var.connectorfreshdesk_property_contact_id
  }
  property {
    name  = "contactJobTitle"
    type  = "string"
    value = var.connectorfreshdesk_property_contact_job_title
  }
  property {
    name  = "contactName"
    type  = "string"
    value = var.connectorfreshdesk_property_contact_name
  }
  property {
    name  = "contactPhone"
    type  = "string"
    value = var.connectorfreshdesk_property_contact_phone
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorfreshdesk_property_email
  }
  property {
    name  = "ticketDescription"
    type  = "string"
    value = var.connectorfreshdesk_property_ticket_description
  }
  property {
    name  = "ticketId"
    type  = "string"
    value = var.connectorfreshdesk_property_ticket_id
  }
  property {
    name  = "ticketPriority"
    type  = "string"
    value = var.connectorfreshdesk_property_ticket_priority
  }
  property {
    name  = "ticketStatus"
    type  = "string"
    value = var.connectorfreshdesk_property_ticket_status
  }
  property {
    name  = "ticketSubject"
    type  = "string"
    value = var.connectorfreshdesk_property_ticket_subject
  }
  property {
    name  = "version"
    type  = "string"
    value = var.connectorfreshdesk_property_version
  }
}
