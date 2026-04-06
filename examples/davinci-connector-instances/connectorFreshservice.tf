resource "pingone_davinci_connector_instance" "connectorFreshservice" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshservice"
  }
  name = "My awesome connectorFreshservice"
  property {
    name  = "address"
    type  = "string"
    value = var.connectorfreshservice_property_address
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorfreshservice_property_api_key
  }
  property {
    name  = "backgroundInformation"
    type  = "string"
    value = var.connectorfreshservice_property_background_information
  }
  property {
    name  = "canSeeAllTicketsFromAssociatedDepts"
    type  = "string"
    value = var.connectorfreshservice_property_can_see_all_tickets_from_associated_depts
  }
  property {
    name  = "ccEmails"
    type  = "string"
    value = var.connectorfreshservice_property_cc_emails
  }
  property {
    name  = "departmentIds"
    type  = "string"
    value = var.connectorfreshservice_property_department_ids
  }
  property {
    name  = "description"
    type  = "string"
    value = var.connectorfreshservice_property_description
  }
  property {
    name  = "domain"
    type  = "string"
    value = var.connectorfreshservice_property_domain
  }
  property {
    name  = "dueBy"
    type  = "string"
    value = var.connectorfreshservice_property_due_by
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorfreshservice_property_email
  }
  property {
    name  = "emailConfigId"
    type  = "string"
    value = var.connectorfreshservice_property_email_config_id
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.connectorfreshservice_property_first_name
  }
  property {
    name  = "frDueBy"
    type  = "string"
    value = var.connectorfreshservice_property_fr_due_by
  }
  property {
    name  = "groupId"
    type  = "string"
    value = var.connectorfreshservice_property_group_id
  }
  property {
    name  = "impact"
    type  = "string"
    value = var.connectorfreshservice_property_impact
  }
  property {
    name  = "jobTitle"
    type  = "string"
    value = var.connectorfreshservice_property_job_title
  }
  property {
    name  = "language"
    type  = "string"
    value = var.connectorfreshservice_property_language
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.connectorfreshservice_property_last_name
  }
  property {
    name  = "locationId"
    type  = "string"
    value = var.connectorfreshservice_property_location_id
  }
  property {
    name  = "mobilePhoneNumber"
    type  = "string"
    value = var.connectorfreshservice_property_mobile_phone_number
  }
  property {
    name  = "name"
    type  = "string"
    value = var.connectorfreshservice_property_name
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectorfreshservice_property_phone
  }
  property {
    name  = "primaryEmail"
    type  = "string"
    value = var.connectorfreshservice_property_primary_email
  }
  property {
    name  = "priority"
    type  = "string"
    value = var.connectorfreshservice_property_priority
  }
  property {
    name  = "reportingManagerId"
    type  = "string"
    value = var.connectorfreshservice_property_reporting_manager_id
  }
  property {
    name  = "requesterId"
    type  = "string"
    value = var.connectorfreshservice_property_requester_id
  }
  property {
    name  = "responderId"
    type  = "string"
    value = var.connectorfreshservice_property_responder_id
  }
  property {
    name  = "secondaryEmails"
    type  = "string"
    value = var.connectorfreshservice_property_secondary_emails
  }
  property {
    name  = "source"
    type  = "string"
    value = var.connectorfreshservice_property_source
  }
  property {
    name  = "status"
    type  = "string"
    value = var.connectorfreshservice_property_status
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.connectorfreshservice_property_subject
  }
  property {
    name  = "tags"
    type  = "string"
    value = var.connectorfreshservice_property_tags
  }
  property {
    name  = "ticketId"
    type  = "string"
    value = var.connectorfreshservice_property_ticket_id
  }
  property {
    name  = "timeFormat"
    type  = "string"
    value = var.connectorfreshservice_property_time_format
  }
  property {
    name  = "timeZone"
    type  = "string"
    value = var.connectorfreshservice_property_time_zone
  }
  property {
    name  = "urgency"
    type  = "string"
    value = var.connectorfreshservice_property_urgency
  }
  property {
    name  = "workPhoneNumber"
    type  = "string"
    value = var.connectorfreshservice_property_work_phone_number
  }
}
