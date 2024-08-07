resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_mobile_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "Mobile App"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "org.bxretail.mobileapp"
      package_name = "org.bxretail.mobileapp"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 30
          units  = "HOURS"
        }

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_integrity_api_decryption_key
          verification_key  = var.google_play_integrity_api_verification_key
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "example_fcm" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_mobile_application.id

  fcm = {
    google_service_account_credentials = var.google_service_account_credentials_json
  }
}

resource "pingone_mfa_application_push_credential" "example_apns" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_mobile_application.id

  apns = {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}

resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA device policy"

  depends_on = [
    pingone_mfa_application_push_credential.example_fcm,
    pingone_mfa_application_push_credential.example_apns,
  ]

  mobile = {
    enabled = true

    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 5
          time_unit = "MINUTES"
        }
      }
    }

    applications = {
      (pingone_application.my_mobile_application.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }

        integrity_detection = "restrictive"
      }
    }
  }

  totp = {
    enabled = true
  }

  fido2 = {
    enabled = true
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }
}
