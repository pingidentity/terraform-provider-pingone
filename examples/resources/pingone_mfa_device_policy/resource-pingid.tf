resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  environment_id = pingone_environment.my_environment.id
  policy_type    = "PING_ONE_ID"

  name = "My Awesome PingID Device Policy"

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
      (pingone_application.my_native_app.id) = {
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

  fido2 = {
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

resource "pingone_application" "my_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Native App"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_environment" "my_environment" {
  # ...
}
