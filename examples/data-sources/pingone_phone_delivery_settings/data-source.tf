resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = var.environment_id

  provider_custom = {
    name = "My awesome custom notifications provider"

    authentication = {
      method     = "BEARER"
      auth_token = var.custom_provider_auth_token
    }

    requests = [
      {
        delivery_method = "SMS"
        method          = "GET"
        url             = "https://api.my-sms-gateway.com/send-sms.json?to=$${to}&from=$${from}&message=$${message}"
      }
    ]
  }
}

data "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = var.environment_id

  name = "My awesome custom notifications provider"

  depends_on = [
    pingone_phone_delivery_settings.my_awesome_custom_provider,
  ]
}
