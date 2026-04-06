resource "pingone_davinci_connector_instance" "connectorMailgun" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMailgun"
  }
  name = "My awesome connectorMailgun"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectormailgun_property_api_key
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.connectormailgun_property_api_version
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectormailgun_property_body
  }
  property {
    name  = "domainEvent"
    type  = "string"
    value = var.connectormailgun_property_domain_event
  }
  property {
    name  = "domainEventDuration"
    type  = "string"
    value = var.connectormailgun_property_domain_event_duration
  }
  property {
    name  = "emailBcc"
    type  = "string"
    value = var.connectormailgun_property_email_bcc
  }
  property {
    name  = "emailCc"
    type  = "string"
    value = var.connectormailgun_property_email_cc
  }
  property {
    name  = "emailSubject"
    type  = "string"
    value = var.connectormailgun_property_email_subject
  }
  property {
    name  = "emailTag"
    type  = "string"
    value = var.connectormailgun_property_email_tag
  }
  property {
    name  = "emailText"
    type  = "string"
    value = var.connectormailgun_property_email_text
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectormailgun_property_endpoint
  }
  property {
    name  = "fromAddress"
    type  = "string"
    value = var.connectormailgun_property_from_address
  }
  property {
    name  = "fromName"
    type  = "string"
    value = var.connectormailgun_property_from_name
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectormailgun_property_headers
  }
  property {
    name  = "mailgunDomain"
    type  = "string"
    value = var.connectormailgun_property_mailgun_domain
  }
  property {
    name  = "mailingListAccessLevel"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_access_level
  }
  property {
    name  = "mailingListAddress"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_address
  }
  property {
    name  = "mailingListDescription"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_description
  }
  property {
    name  = "mailingListMemberAddress"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_member_address
  }
  property {
    name  = "mailingListMemberAttributes"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_member_attributes
  }
  property {
    name  = "mailingListMemberName"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_member_name
  }
  property {
    name  = "mailingListMemberSubscriptionStatus"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_member_subscription_status
  }
  property {
    name  = "mailingListMemberUpsert"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_member_upsert
  }
  property {
    name  = "mailingListName"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_name
  }
  property {
    name  = "mailingListReplyPreference"
    type  = "string"
    value = var.connectormailgun_property_mailing_list_reply_preference
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectormailgun_property_method
  }
  property {
    name  = "newMailingListMemberAddress"
    type  = "string"
    value = var.connectormailgun_property_new_mailing_list_member_address
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectormailgun_property_query_parameters
  }
  property {
    name  = "toAddress"
    type  = "string"
    value = var.connectormailgun_property_to_address
  }
}
