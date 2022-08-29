resource "pingone_environment" "my_environment" {
  # ...
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

resource "pingone_identity_provider_attribute" "apple_email" {
  environment_id       = pingone_environment.my_environment.id
  identity_provider_id = pingone_identity_provider.apple.id

  name   = "email"
  update = "EMPTY_ONLY"
  value  = "$${providerAttributes.user.email}"
}
