resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_language" "example_by_locale" {
  environment_id = var.environment_id

  locale = "fr"
}

resource "pingone_language_update" "my_customers_language" {
  environment_id = pingone_environment.my_environment.id

  language_id = data.pingone_language.example_by_locale.id
  enabled     = true
  default     = true
}