resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_language" "my_customers_language" {
  environment_id = pingone_environment.my_environment.id

  locale = "fr-FR"
}
