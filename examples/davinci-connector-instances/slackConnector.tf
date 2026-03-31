resource "pingone_davinci_connector_instance" "slackConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "slackConnector"
  }
  name = "My awesome slackConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.slackconnector_property_auth_type
  }
  property {
    name  = "blocks"
    type  = "string"
    value = var.slackconnector_property_blocks
  }
  property {
    name  = "button"
    type  = "string"
    value = var.slackconnector_property_button
  }
  property {
    name  = "channelId"
    type  = "string"
    value = var.slackconnector_property_channel_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.slackconnector_property_email
  }
  property {
    name  = "messageText"
    type  = "string"
    value = var.slackconnector_property_message_text
  }
  property {
    name  = "oauth2"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.slackconnector_property_show_powered_by
  }
  property {
    name  = "unFurlLinks"
    type  = "string"
    value = var.slackconnector_property_un_furl_links
  }
  property {
    name  = "unFurlMedia"
    type  = "string"
    value = var.slackconnector_property_un_furl_media
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.slackconnector_property_user_id
  }
}
