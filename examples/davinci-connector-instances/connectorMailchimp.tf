resource "pingone_davinci_connector_instance" "connectorMailchimp" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMailchimp"
  }
  name = "My awesome connectorMailchimp"
  property {
    name  = "allowlistEmail"
    type  = "string"
    value = var.connectormailchimp_property_allowlist_email
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectormailchimp_property_body
  }
  property {
    name  = "denylistEmail"
    type  = "string"
    value = var.connectormailchimp_property_denylist_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectormailchimp_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectormailchimp_property_headers
  }
  property {
    name  = "messageFromEmail"
    type  = "string"
    value = var.connectormailchimp_property_message_from_email
  }
  property {
    name  = "messageSubject"
    type  = "string"
    value = var.connectormailchimp_property_message_subject
  }
  property {
    name  = "messageText"
    type  = "string"
    value = var.connectormailchimp_property_message_text
  }
  property {
    name  = "messageToEmail"
    type  = "string"
    value = var.connectormailchimp_property_message_to_email
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectormailchimp_property_method
  }
  property {
    name  = "optionalComment"
    type  = "string"
    value = var.connectormailchimp_property_optional_comment
  }
  property {
    name  = "removeAllowlistEmail"
    type  = "string"
    value = var.connectormailchimp_property_remove_allowlist_email
  }
  property {
    name  = "storeId"
    type  = "string"
    value = var.connectormailchimp_property_store_id
  }
  property {
    name  = "transactionalApiKey"
    type  = "string"
    value = var.connectormailchimp_property_transactional_api_key
  }
  property {
    name  = "transactionalApiVersion"
    type  = "string"
    value = var.connectormailchimp_property_transactional_api_version
  }
}
