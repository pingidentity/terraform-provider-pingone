resource "pingone_davinci_connector_instance" "notificationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "notificationsConnector"
  }
  name = "My awesome notificationsConnector"
  property {
    name  = "clientId"
    type  = "string"
    value = var.notificationsconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.notificationsconnector_property_client_secret
  }
  property {
    name  = "customTemplateVariant"
    type  = "string"
    value = var.notificationsconnector_property_custom_template_variant
  }
  property {
    name  = "email"
    type  = "string"
    value = var.notificationsconnector_property_email
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.notificationsconnector_property_env_id
  }
  property {
    name  = "notificationPolicyId"
    type  = "string"
    value = var.notificationsconnector_property_notification_policy_id
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.notificationsconnector_property_phone
  }
  property {
    name  = "region"
    type  = "string"
    value = var.notificationsconnector_property_region
  }
  property {
    name  = "sendSync"
    type  = "boolean"
    value = var.notificationsconnector_property_send_sync
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.notificationsconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.notificationsconnector_property_skip_button_press
  }
  property {
    name  = "templateLocale"
    type  = "string"
    value = var.notificationsconnector_property_template_locale
  }
  property {
    name  = "templateVariables"
    type  = "string"
    value = var.notificationsconnector_property_template_variables
  }
  property {
    name  = "templateVariant"
    type  = "string"
    value = var.notificationsconnector_property_template_variant
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.notificationsconnector_property_user_agent
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.notificationsconnector_property_user_id
  }
}
