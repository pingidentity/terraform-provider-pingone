resource "pingone_mfa_device_policy_default" "my_awesome_mfa_device_policy_default" {
  environment_id = pingone_environment.my_environment.id
  policy_type    = "PING_ONE_MFA"

  name = "My Awesome Default MFA Device Policy"

  authentication = {
    device_selection = "DEFAULT_TO_FIRST"
  }

  new_device_notification = "SMS_THEN_EMAIL"
  ignore_user_lock        = true

  notifications_policy = {
    id = pingone_notification_policy.my_notification_policy.id
  }

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 60
        time_unit = "MINUTES"
      }
    }
  }

  sms = {
    enabled                        = true
    pairing_disabled               = true
    prompt_for_nickname_on_pairing = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }
      otp_length = 7
    }
  }

  email = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 6
    }
  }

  voice = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 6
    }
  }

  mobile = {
    enabled                        = true
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }

    applications = [
      {
        id = pingone_application.my_native_app.id
        auto_enrollment = {
          enabled = true
        }
        device_authorization = {
          enabled            = true
          extra_verification = "PERMISSIVE"
        }
        integrity_detection = "PERMISSIVE"
        otp = {
          enabled = true
        }
        push = {
          enabled = true
          number_matching = {
            enabled = true
          }
        }
        pairing_disabled = false
        pairing_key_lifetime = {
          duration  = 10
          time_unit = "MINUTES"
        }
        push_limit = {
          count = 5
          lock_duration = {
            duration  = 30
            time_unit = "MINUTES"
          }
          time_period = {
            duration  = 10
            time_unit = "MINUTES"
          }
        }
        push_timeout = {
          duration  = 30
          time_unit = "SECONDS"
        }
      }
    ]
  }

  totp = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }

  fido2 = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
  }
}

resource "pingone_application" "my_native_app" {
  # ...
}

resource "pingone_notification_policy" "my_notification_policy" {
  # ...
}

resource "pingone_environment" "my_environment" {
  # ...
}
