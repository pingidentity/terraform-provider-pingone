resource "pingone_davinci_connector_instance" "connectorWhatsAppBusiness" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorWhatsAppBusiness"
  }
  name = "My awesome connectorWhatsAppBusiness"
  properties = jsonencode({
    "accessToken" = var.connectorwhatsappbusiness_property_access_token
    "appSecret" = var.connectorwhatsappbusiness_property_app_secret
    "skWebhookUri" = var.connectorwhatsappbusiness_property_sk_webhook_uri
    "verifyToken" = var.connectorwhatsappbusiness_property_verify_token
    "version" = var.connectorwhatsappbusiness_property_version
  })
}
