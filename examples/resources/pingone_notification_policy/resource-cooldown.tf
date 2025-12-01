resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "cooldown" {
  environment_id = pingone_environment.my_environment.id

  name = "Policy with Cooldown Configuration"

  quota = [
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["SMS", "Voice"]
      total            = 100
    }
  ]

  cooldown_configuration = {
    email = {
      enabled      = true
      resend_limit = 5
      group_by     = "USER_ID"

      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 1
          time_unit = "MINUTES"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled      = true
      resend_limit = 3
      group_by     = "USER_ID"

      periods = [
        {
          duration  = 45
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        },
        {
          duration  = 5
          time_unit = "MINUTES"
        }
      ]
    }

    voice = {
      enabled      = true
      resend_limit = 3

      periods = [
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 3
          time_unit = "MINUTES"
        },
        {
          duration  = 5
          time_unit = "MINUTES"
        }
      ]
    }

    whats_app = {
      enabled = false
    }
  }
}
