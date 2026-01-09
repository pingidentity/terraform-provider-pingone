resource "pingone_mfa_device_policy_default" "my_awesome_mfa_device_policy_default" {
  environment_id = pingone_environment.my_environment.id
  policy_type    = "PING_ONE_ID"

  name = "My Awesome PingID Device Policy"

  authentication = {
    device_selection = "PROMPT_TO_SELECT"
  }

  new_device_notification = "SMS_THEN_EMAIL"

  sms = {
    enabled = true
  }

  voice = {
    enabled = true
  }

  email = {
    enabled = true
  }

  mobile = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }

    applications = {
      (data.pingone_application.my_native_app.id) = {
        type               = "pingIdAppConfig"
        biometrics_enabled = true
        ip_pairing_configuration = {
          any_ip_address = true
        }
        new_request_duration_configuration = {
          device_timeout = {
            duration  = 25
            time_unit = "SECONDS"
          }
          total_timeout = {
            duration  = 40
            time_unit = "SECONDS"
          }
        }
        otp = {
          enabled = true
        }
        push = {
          enabled = true
        }
        pairing_disabled = false
      }
    }
  }

  totp = {
    enabled = true
  }

  desktop = {
    enabled = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 3
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }

  yubikey = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }

  oath_token = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }
}

data "pingone_application" "my_native_app" {
  # ...
}

resource "pingone_environment" "my_environment" {
  # ...
}
