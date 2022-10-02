data "pingone_language" "example_by_locale" {
  environment_id = var.environment_id

  locale = "fr-FR"
}

data "pingone_language" "example_by_id" {
  environment_id = var.environment_id

  language_id = var.language_id
}