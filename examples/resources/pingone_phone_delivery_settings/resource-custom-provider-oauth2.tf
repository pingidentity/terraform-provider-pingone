resource "pingone_phone_delivery_settings" "my_awesome_custom_provider_oauth2" {
  environment_id = pingone_environment.my_environment.id

  provider_custom = {
    name = "My awesome custom notifications provider"

    authentication = {
      method        = "OAUTH2"
      auth_url      = "https://auth.my-sms-gateway.com/oauth2/token"
      grant_type    = "CLIENT_CREDENTIALS"
      client_id     = var.custom_provider_client_id
      client_secret = var.custom_provider_client_secret
      scopes        = ["sms:send", "voice:send"]
    }

    requests = [
      {
        delivery_method = "SMS"
        method          = "POST"
        url             = "https://api.my-sms-gateway.com/send-sms"

        headers = {
          "Content-Type" = "application/json"
        }

        body = jsonencode({
          to      = "$${to}"
          from    = "$${from}"
          message = "$${message}"
        })
      }
    ]
  }
}
