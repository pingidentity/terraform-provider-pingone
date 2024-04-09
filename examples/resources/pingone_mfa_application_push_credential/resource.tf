resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_mobile_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Mobile App"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "NONE"
    pkce_enforcement            = "S256_REQUIRED"

    mobile_app {

      // Apple
      bundle_id = "org.bxretail.mybundle"

      // Android
      package_name = "org.bxretail.mypackage"

      // Huawei
      huawei_app_id       = "12345679"
      huawei_package_name = "org.bxretail.huaweipackage"
    }
  }
}

// Android
resource "pingone_mfa_application_push_credential" "example_fcm" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_mobile_application.id

  fcm = {
    google_service_account_credentials = var.google_service_account_credentials_json
  }
}

// Apple
resource "pingone_mfa_application_push_credential" "example_apns" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_mobile_application.id

  apns = {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}

// Huawei
resource "pingone_mfa_application_push_credential" "example_hms" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_mobile_application.id

  hms = {
    client_id     = var.hms_client_id
    client_secret = var.hms_client_secret
  }
}