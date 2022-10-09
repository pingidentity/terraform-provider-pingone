resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_language" "my_customers_language" {
  environment_id = pingone_environment.my_environment.id

  locale = "fr-FR"
}

resource "pingone_language_update" "my_customers_language" {
  environment_id = pingone_environment.my_environment.id

  language_id = pingone_language.my_customers_language.id
  enabled     = true
  default     = true
}