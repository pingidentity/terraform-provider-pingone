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
  default     = true
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

resource "time_static" "now" {}

resource "pingone_agreement_localization_revision" "my_agreement_fr_now" {
  environment_id            = pingone_environment.my_environment.id
  agreement_id              = pingone_agreement.my_agreement.id
  agreement_localization_id = pingone_agreement_localization.my_agreement_fr.id

  content_type      = "text/html"
  effective_at      = time_static.now.id
  require_reconsent = true
  text              = <<EOT
<h1>Conditions de service</h1>

Veuillez accepter les termes et conditions.

<h2>Utilisation des données</h2>

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

<h2>Soutien</h2>

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
EOT
}

resource "pingone_agreement_localization_enable" "my_agreement_fr_enable" {
  environment_id            = pingone_environment.my_environment.id
  agreement_id              = pingone_agreement.my_agreement.id
  agreement_localization_id = pingone_agreement_localization.my_agreement_fr.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_revision.my_agreement_fr_now
  ]
}
