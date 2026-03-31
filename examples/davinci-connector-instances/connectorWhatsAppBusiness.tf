resource "pingone_davinci_connector_instance" "connectorWhatsAppBusiness" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorWhatsAppBusiness"
  }
  name = "My awesome connectorWhatsAppBusiness"
  property {
    name  = "accessToken"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_access_token
  }
  property {
    name  = "appSecret"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_app_secret
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_body
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_endpoint
  }
  property {
    name  = "fromId"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_from_id
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_headers
  }
  property {
    name  = "interactiveMessagetemplateName"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_interactive_messagetemplate_name
  }
  property {
    name  = "mediaMessageTemplateName"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_media_message_template_name
  }
  property {
    name  = "mediaMessageUrl"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_media_message_url
  }
  property {
    name  = "messageLanguageCode"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_message_language_code
  }
  property {
    name  = "messageType"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_message_type
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_query_parameters
  }
  property {
    name  = "skWebhookUri"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_sk_webhook_uri
  }
  property {
    name  = "templateArgumentList"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_template_argument_list
  }
  property {
    name  = "textMessageTemplateName"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_text_message_template_name
  }
  property {
    name  = "to"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_to
  }
  property {
    name  = "verifyToken"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_verify_token
  }
  property {
    name  = "version"
    type  = "string"
    value = var.connectorwhatsappbusiness_property_version
  }
}
