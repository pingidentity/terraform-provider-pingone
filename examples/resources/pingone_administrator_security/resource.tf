resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_identity_provider" "google" {
  environment_id = pingone_environment.my_environment.id

  name    = "Google"
  enabled = true

  google = {
    client_id     = var.google_client_id
    client_secret = var.google_client_secret
  }
}

resource "pingone_administrator_security" "my_administrator_security" {
  environment_id = pingone_environment.my_environment.id

  allowed_methods = {
    email = jsonencode(
      {
        enabled = true
      }
    )
    fido2 = jsonencode(
      {
        enabled = true
      }
    )
    totp = jsonencode(
      {
        enabled = true
      }
    )
  }
  authentication_method = "HYBRID"
  mfa_status            = "ENFORCE"
  identity_provider = {
    id = pingone_identity_provider.google.id
  }
  recovery = true
}