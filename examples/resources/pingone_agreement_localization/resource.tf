resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_language" "fr" {
  environment_id = pingone_environment.my_environment.id

  locale = "fr"
}

resource "pingone_language_update" "fr" {
  environment_id = pingone_environment.my_environment.id

  language_id = data.pingone_language.fr.id
  default     = false
  enabled     = true
}

resource "pingone_agreement" "my_agreement" {
  environment_id = pingone_environment.my_environment.id

  name        = "Terms and Conditions"
  description = "An agreement for general Terms and Conditions"
}

resource "pingone_agreement_localization" "my_agreement_fr" {
  environment_id = pingone_environment.my_environment.id
  agreement_id   = pingone_agreement.my_agreement.id
  language_id    = pingone_language_update.fr.id

  display_name = "Terms and Conditions - French Locale"
}