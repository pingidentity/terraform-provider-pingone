resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_mobile_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "Mobile App"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "org.bxretail.mobileapp"
      package_name = "org.bxretail.mobileapp"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = true

        cache_duration {
          amount = 30
          units  = "HOURS"
        }

        google_play {
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

  fcm {
    key = var.fcm_key
  }
}

resource "pingone_mfa_application_push_credential" "example_apns" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_mobile_application.id

  apns {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}

resource "pingone_mfa_policy" "my_awesome_mfa_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  depends_on = [
    pingone_mfa_application_push_credential.example_fcm,
    pingone_mfa_application_push_credential.example_apns,
  ]

  mobile {
    enabled = true

    otp_failure_count = 3

    application {
      id = pingone_application.my_mobile_application.id

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true

      integrity_detection = "restrictive"
    }
  }

  totp {
    enabled = true
  }

  security_key {
    enabled = true
  }

  platform {
    enabled = true
  }

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

}
