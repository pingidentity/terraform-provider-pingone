resource "pingone_phone_delivery_settings" "my_awesome_custom_provider_custom_header" {
  environment_id = pingone_environment.my_environment.id

  provider_custom = {
    name = "My awesome custom notifications provider"

    authentication = {
      method       = "CUSTOM_HEADER"
      header_name  = "X-Api-Key"
      header_value = var.custom_provider_api_key
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
