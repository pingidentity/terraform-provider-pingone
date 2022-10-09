resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "pingone_mfa_application_push_credential" "example_fcm" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  fcm {
    key = var.fcm_key
  }
}

resource "pingone_mfa_application_push_credential" "example_apns" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  apns {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}