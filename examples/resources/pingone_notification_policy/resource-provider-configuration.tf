resource "pingone_phone_delivery_settings" "my_custom_provider_us" {
  # ...
}

resource "pingone_phone_delivery_settings" "my_custom_provider_uk" {
  # ...
}

resource "pingone_phone_delivery_settings" "my_fallback_provider" {
  # ...
}

resource "pingone_notification_policy" "my_notification_policy_with_provider_config" {
  environment_id = pingone_environment.my_environment.id

  name = "My notification policy with custom providers"

  quota = [
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["SMS", "Voice"]
      total            = 100
    }
  ]

  provider_configuration = {
    conditions = [
      {
        delivery_methods = ["SMS", "VOICE"]
        countries        = ["US", "CA"]

        fallback_chain = [
          {
            id = pingone_phone_delivery_settings.my_custom_provider_us.id
          },
          {
            id = pingone_phone_delivery_settings.my_fallback_provider.id
          }
        ]
      },
      {
        delivery_methods = ["SMS"]
        countries        = ["GB"]

        fallback_chain = [
          {
            id = pingone_phone_delivery_settings.my_custom_provider_uk.id
          },
          {
            id = pingone_phone_delivery_settings.my_fallback_provider.id
          }
        ]
      },
      {
        delivery_methods = ["SMS", "VOICE"]

        fallback_chain = [
          {
            id = pingone_phone_delivery_settings.my_fallback_provider.id
          }
        ]
      }
    ]
  }
}
