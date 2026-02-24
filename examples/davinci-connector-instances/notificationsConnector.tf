resource "pingone_davinci_connector_instance" "notificationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "notificationsConnector"
  }
  name = "My awesome notificationsConnector"
  properties = jsonencode({
    "clientId" = var.notificationsconnector_property_client_id
    "clientSecret" = var.notificationsconnector_property_client_secret
    "envId" = var.notificationsconnector_property_env_id
    "notificationPolicyId" = var.notificationsconnector_property_notification_policy_id
    "region" = var.notificationsconnector_property_region
  })
}
