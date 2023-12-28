resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_forms_recaptcha_v2" "my_awesome_recaptcha_config" {
  environment_id = pingone_environment.my_environment.id

  site_key   = var.google_recaptcha_site_key
  secret_key = var.google_recaptcha_secret_key
}
