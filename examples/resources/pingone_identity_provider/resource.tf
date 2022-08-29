resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_identity_provider" "facebook" {
  environment_id = pingone_environment.my_environment.id

  name    = "Facebook"
  enabled = true

  facebook {
    app_id     = var.facebook_app_id
    app_secret = var.facebook_app_secret
  }
}

resource "pingone_identity_provider" "google" {
  environment_id = pingone_environment.my_environment.id

  name    = "Google"
  enabled = true

  google {
    client_id     = var.google_client_id
    client_secret = var.google_client_secret
  }
}

resource "pingone_identity_provider" "apple" {
  environment_id = pingone_environment.my_environment.id

  name    = "Apple"
  enabled = true

  apple {
    client_id                 = var.apple_client_id
    client_secret_signing_key = var.apple_client_secret_signing_key
    key_id                    = var.apple_key_id
    team_id                   = var.apple_team_id
  }
}