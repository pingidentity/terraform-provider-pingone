resource "pingone_davinci_connector_instance" "connectorSalesforceMarketingCloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSalesforceMarketingCloud"
  }
  name = "My awesome connectorSalesforceMarketingCloud"
  property {
    name  = "SalesforceMarketingCloudURL"
    type  = "string"
    value = var.salesforce_marketing_cloud_url
  }
  property {
    name  = "Subscribers"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_subscribers
  }
  property {
    name  = "accountId"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_account_id
  }
  property {
    name  = "attributeSets"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_attribute_sets
  }
  property {
    name  = "attributes"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_attributes
  }
  property {
    name  = "blackoutWindow"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_blackout_window
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_client_secret
  }
  property {
    name  = "contactKey"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_contact_key
  }
  property {
    name  = "contactKeyArray"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_contact_key_array
  }
  property {
    name  = "content"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_content
  }
  property {
    name  = "contentURL"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_content_url
  }
  property {
    name  = "definitionKey"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_definition_key
  }
  property {
    name  = "from"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_from
  }
  property {
    name  = "id"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_id
  }
  property {
    name  = "keyword"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_keyword
  }
  property {
    name  = "messageText"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_message_text
  }
  property {
    name  = "objectIDorKey"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_object_idor_key
  }
  property {
    name  = "options"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_options
  }
  property {
    name  = "optionsASYNCandSYNC"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_options_asyncand_sync
  }
  property {
    name  = "override"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_override
  }
  property {
    name  = "recipients"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_recipients
  }
  property {
    name  = "resubscribe"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_resubscribe
  }
  property {
    name  = "scope"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_scope
  }
  property {
    name  = "sendTime"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_send_time
  }
  property {
    name  = "subscribe"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_subscribe
  }
  property {
    name  = "to"
    type  = "string"
    value = var.connectorsalesforcemarketingcloud_property_to
  }
}
