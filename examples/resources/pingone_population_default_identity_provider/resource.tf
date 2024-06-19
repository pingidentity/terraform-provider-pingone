resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_awesome_population" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome population"
  description = "My new population for awesome people"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
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

resource "pingone_population_default_identity_provider" "my_awesome_population" {
  environment_id = pingone_environment.my_environment.id
  population_id  = pingone_population.my_awesome_population.id

  identity_provider_id = pingone_identity_provider.google.id
}
