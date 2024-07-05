resource "pingone_application" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Native Mobile App"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris = [
      "https://demo.bxretail.org/app/callback",
      "org.bxretail.app://callback"
    ]

    mobile_app = {
      bundle_id           = var.apple_bundle_id
      package_name        = var.android_package_name
      huawei_app_id       = var.huawei_app_id
      huawei_package_name = var.huawei_package_name

      universal_app_link = "https://demo.bxretail.org"

      passcode_refresh_seconds = 30

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 24
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

resource "time_rotating" "my_awesome_native_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_native_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_native_app_secret_rotation.rotation_rfc3339,
  }
}
