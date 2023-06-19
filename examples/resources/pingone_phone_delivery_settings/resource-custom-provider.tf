resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_custom = {
    name = "My awesome custom notifications provider"

    authentication = {
      method   = "BASIC"
      username = var.custom_provider_username
      password = var.custom_provider_password
    }

    requests = [
      {
        delivery_method = "SMS"
        method          = "GET"
        url             = "https://api.my-sms-gateway.com/send-sms.json?to=$${to}&from=$${from}&message=$${message}"
      },
      {
        delivery_method = "Voice"
        method          = "POST"
        url             = "https://api.my-voice-gateway.com/send-voice"

        headers = {
          "Content-Type" = "application/json"
        }

        body = jsonencode({
          to      = "$${to}"
          from    = "$${from}"
          message = "$${message}"
        })

        after_tag  = "</Say> <Pause length=\"1\"/>"
        before_tag = "<Say>"
      }
    ]
  }
}