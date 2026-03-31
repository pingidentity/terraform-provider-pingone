resource "pingone_davinci_connector_instance" "connectorZendesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZendesk"
  }
  name = "My awesome connectorZendesk"
  property {
    name  = "alias"
    type  = "string"
    value = var.connectorzendesk_property_alias
  }
  property {
    name  = "apiToken"
    type  = "string"
    value = var.connectorzendesk_property_api_token
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorzendesk_property_body
  }
  property {
    name  = "chatOnly"
    type  = "string"
    value = var.connectorzendesk_property_chat_only
  }
  property {
    name  = "collaborators"
    type  = "string"
    value = var.connectorzendesk_property_collaborators
  }
  property {
    name  = "comment"
    type  = "string"
    value = var.connectorzendesk_property_comment
  }
  property {
    name  = "commentBody"
    type  = "string"
    value = var.connectorzendesk_property_comment_body
  }
  property {
    name  = "customFields"
    type  = "string"
    value = var.connectorzendesk_property_custom_fields
  }
  property {
    name  = "customRoleID"
    type  = "string"
    value = var.connectorzendesk_property_custom_role_id
  }
  property {
    name  = "details"
    type  = "string"
    value = var.connectorzendesk_property_details
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorzendesk_property_email
  }
  property {
    name  = "emailCCs"
    type  = "string"
    value = var.connectorzendesk_property_email_ccs
  }
  property {
    name  = "emailUsername"
    type  = "string"
    value = var.connectorzendesk_property_email_username
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorzendesk_property_endpoint
  }
  property {
    name  = "externalID"
    type  = "string"
    value = var.connectorzendesk_property_external_id
  }
  property {
    name  = "fileName"
    type  = "string"
    value = var.connectorzendesk_property_file_name
  }
  property {
    name  = "groupID"
    type  = "string"
    value = var.connectorzendesk_property_group_id
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorzendesk_property_headers
  }
  property {
    name  = "isPublic"
    type  = "string"
    value = var.connectorzendesk_property_is_public
  }
  property {
    name  = "locale"
    type  = "string"
    value = var.connectorzendesk_property_locale
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorzendesk_property_method
  }
  property {
    name  = "moderator"
    type  = "string"
    value = var.connectorzendesk_property_moderator
  }
  property {
    name  = "name"
    type  = "string"
    value = var.connectorzendesk_property_name
  }
  property {
    name  = "notes"
    type  = "string"
    value = var.connectorzendesk_property_notes
  }
  property {
    name  = "organizationID"
    type  = "string"
    value = var.connectorzendesk_property_organization_id
  }
  property {
    name  = "organizationMembershipID"
    type  = "string"
    value = var.connectorzendesk_property_organization_membership_id
  }
  property {
    name  = "organizationName"
    type  = "string"
    value = var.connectorzendesk_property_organization_name
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectorzendesk_property_phone
  }
  property {
    name  = "priority"
    type  = "string"
    value = var.connectorzendesk_property_priority
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorzendesk_property_query_parameters
  }
  property {
    name  = "recipient"
    type  = "string"
    value = var.connectorzendesk_property_recipient
  }
  property {
    name  = "role"
    type  = "string"
    value = var.connectorzendesk_property_role
  }
  property {
    name  = "safeUpdate"
    type  = "string"
    value = var.connectorzendesk_property_safe_update
  }
  property {
    name  = "status"
    type  = "string"
    value = var.connectorzendesk_property_status
  }
  property {
    name  = "subdomain"
    type  = "string"
    value = var.connectorzendesk_property_subdomain
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.connectorzendesk_property_subject
  }
  property {
    name  = "tags"
    type  = "string"
    value = var.connectorzendesk_property_tags
  }
  property {
    name  = "ticketFormID"
    type  = "string"
    value = var.connectorzendesk_property_ticket_form_id
  }
  property {
    name  = "ticketID"
    type  = "string"
    value = var.connectorzendesk_property_ticket_id
  }
  property {
    name  = "ticketRestriction"
    type  = "string"
    value = var.connectorzendesk_property_ticket_restriction
  }
  property {
    name  = "type"
    type  = "string"
    value = var.connectorzendesk_property_type
  }
  property {
    name  = "userFields"
    type  = "string"
    value = var.connectorzendesk_property_user_fields
  }
  property {
    name  = "userID"
    type  = "string"
    value = var.connectorzendesk_property_user_id
  }
  property {
    name  = "verified"
    type  = "string"
    value = var.connectorzendesk_property_verified
  }
  property {
    name  = "viewTickets"
    type  = "string"
    value = var.connectorzendesk_property_view_tickets
  }
}
